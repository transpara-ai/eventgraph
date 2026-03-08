//! Graph facade module -- top-level API for the eventgraph.
//!
//! Ports the Go `graph` package. Provides a simplified facade that ties
//! together Store, ActorStore, EventBus, TrustModel, and AuthorityChain.

use std::collections::BTreeMap;
use std::sync::Mutex;

use serde_json::Value;

use crate::actor::{Actor, ActorStore, InMemoryActorStore};
use crate::authority::{AuthorityChain, AuthorityResult, DefaultAuthorityChain};
use crate::bus::EventBus;
use crate::errors::{EventGraphError, Result};
use crate::event::{create_bootstrap, create_event, Event, NoopSigner, Signer};
use crate::store::{InMemoryStore, Store};
use crate::trust::{DefaultTrustModel, TrustConfig, TrustMetrics, TrustModel};
use crate::types::{ActorId, ConversationId, EventId, EventType, Hash};

// ── GraphConfig ────────────────────────────────────────────────────────

/// Configuration for the Graph facade.
#[derive(Debug, Clone)]
pub struct GraphConfig {
    /// Buffer size for event bus subscribers.
    pub subscriber_buffer_size: usize,
    /// Whether to fall back to mechanical evaluation when intelligence is unavailable.
    pub fallback_to_mechanical: bool,
}

impl Default for GraphConfig {
    fn default() -> Self {
        Self {
            subscriber_buffer_size: 256,
            fallback_to_mechanical: true,
        }
    }
}

// ── Graph ──────────────────────────────────────────────────────────────

/// Top-level facade for the eventgraph system.
///
/// Ties together event storage, actor management, trust evaluation,
/// authority checking, and event distribution via the bus.
pub struct Graph {
    store: InMemoryStore,
    actor_store: InMemoryActorStore,
    bus: EventBus,
    trust_model: Box<dyn TrustModel + Send>,
    authority_chain: Box<dyn AuthorityChain + Send>,
    signer: Box<dyn Signer + Send>,
    _config: GraphConfig,
    state: Mutex<GraphState>,
}

struct GraphState {
    started: bool,
    closed: bool,
}

impl Graph {
    /// Creates a new Graph with the given store, actor store, and default configuration.
    pub fn new(store: InMemoryStore, actor_store: InMemoryActorStore) -> Self {
        Self::with_config(store, actor_store, GraphConfig::default())
    }

    /// Creates a new Graph with explicit configuration.
    pub fn with_config(
        store: InMemoryStore,
        actor_store: InMemoryActorStore,
        config: GraphConfig,
    ) -> Self {
        let trust_model: Box<dyn TrustModel + Send> =
            Box::new(DefaultTrustModel::new(TrustConfig::default()));
        let authority_chain: Box<dyn AuthorityChain + Send> =
            Box::new(DefaultAuthorityChain::new(Box::new(
                DefaultTrustModel::new(TrustConfig::default()),
            )));

        Self {
            store,
            actor_store,
            bus: EventBus::new(),
            trust_model,
            authority_chain,
            signer: Box::new(NoopSigner),
            _config: config,
            state: Mutex::new(GraphState {
                started: false,
                closed: false,
            }),
        }
    }

    /// Sets a custom trust model. Must be called before `start()`.
    pub fn set_trust_model(&mut self, model: Box<dyn TrustModel + Send>) {
        self.trust_model = model;
    }

    /// Sets a custom authority chain. Must be called before `start()`.
    pub fn set_authority_chain(&mut self, chain: Box<dyn AuthorityChain + Send>) {
        self.authority_chain = chain;
    }

    /// Sets the default signer. Must be called before `start()`.
    pub fn set_signer(&mut self, signer: Box<dyn Signer + Send>) {
        self.signer = signer;
    }

    // ── Lifecycle ──────────────────────────────────────────────────────

    /// Starts the graph. Must be called before Record, Bootstrap, Evaluate, or Query.
    /// Idempotent -- calling Start on an already-started graph is a no-op.
    pub fn start(&self) -> Result<()> {
        let mut state = self.state.lock().unwrap();
        if state.started {
            return Ok(());
        }
        state.started = true;
        Ok(())
    }

    /// Closes the graph. Idempotent.
    /// After close, all mutating operations return an error.
    pub fn close(&mut self) {
        {
            let mut state = self.state.lock().unwrap();
            if state.closed {
                return;
            }
            state.closed = true;
        }
        self.bus.close();
        self.store.close();
    }

    // ── Guard helpers ──────────────────────────────────────────────────

    fn require_running(&self) -> Result<()> {
        let state = self.state.lock().unwrap();
        if state.closed {
            return Err(EventGraphError::GrammarViolation {
                detail: "graph is closed".to_string(),
            });
        }
        if !state.started {
            return Err(EventGraphError::GrammarViolation {
                detail: "graph is not started (call start first)".to_string(),
            });
        }
        Ok(())
    }

    // ── Bootstrap ──────────────────────────────────────────────────────

    /// Initializes the graph with a genesis event.
    /// Fails if the graph already has events, is closed, or is not started.
    pub fn bootstrap(
        &mut self,
        system_actor: ActorId,
        signer: Option<&dyn Signer>,
    ) -> Result<Event> {
        self.require_running()?;

        if self.store.count() > 0 {
            return Err(EventGraphError::GrammarViolation {
                detail: "graph already bootstrapped".to_string(),
            });
        }

        let s: &dyn Signer = signer.unwrap_or(&*self.signer);
        let ev = create_bootstrap(system_actor, s, 1);
        let stored = self.store.append(ev)?;
        self.bus.publish(&stored);
        Ok(stored)
    }

    // ── Record ─────────────────────────────────────────────────────────

    /// Creates and persists an event, then publishes it on the bus.
    pub fn record(
        &mut self,
        event_type: EventType,
        source: ActorId,
        content: BTreeMap<String, Value>,
        causes: Vec<EventId>,
        conversation_id: ConversationId,
        signer: Option<&dyn Signer>,
    ) -> Result<Event> {
        self.require_running()?;

        let prev_hash = match self.store.head() {
            Some(head) => head.hash.clone(),
            None => Hash::zero(),
        };

        let s: &dyn Signer = signer.unwrap_or(&*self.signer);
        let ev = create_event(
            event_type,
            source,
            content,
            causes,
            conversation_id,
            prev_hash,
            s,
            1,
        );
        let stored = self.store.append(ev)?;
        self.bus.publish(&stored);
        Ok(stored)
    }

    // ── Evaluate ───────────────────────────────────────────────────────

    /// Evaluates authority for a given actor and action.
    pub fn evaluate(&self, actor: &Actor, action: &str) -> Result<AuthorityResult> {
        self.require_running()?;
        self.authority_chain.evaluate(actor, action)
    }

    // ── Query ──────────────────────────────────────────────────────────

    /// Returns a query handle for the graph.
    pub fn query(&self) -> Result<GraphQuery<'_>> {
        self.require_running()?;
        Ok(GraphQuery { graph: self })
    }

    // ── Accessors ──────────────────────────────────────────────────────

    /// Returns a reference to the underlying event store.
    pub fn store(&self) -> &InMemoryStore {
        &self.store
    }

    /// Returns a mutable reference to the underlying event store.
    pub fn store_mut(&mut self) -> &mut InMemoryStore {
        &mut self.store
    }

    /// Returns a reference to the actor store.
    pub fn actor_store(&self) -> &InMemoryActorStore {
        &self.actor_store
    }

    /// Returns a mutable reference to the actor store.
    pub fn actor_store_mut(&mut self) -> &mut InMemoryActorStore {
        &mut self.actor_store
    }

    /// Returns a reference to the event bus.
    pub fn bus(&self) -> &EventBus {
        &self.bus
    }

    /// Returns a mutable reference to the event bus.
    pub fn bus_mut(&mut self) -> &mut EventBus {
        &mut self.bus
    }
}

// ── GraphQuery ─────────────────────────────────────────────────────────

/// A read-only query handle for the graph.
///
/// Provides convenience methods for querying events, trust, and actors.
pub struct GraphQuery<'a> {
    graph: &'a Graph,
}

impl<'a> GraphQuery<'a> {
    /// Returns the most recent events (up to `limit`).
    pub fn recent(&self, limit: usize) -> Vec<&Event> {
        self.graph.store.recent(limit)
    }

    /// Returns events of the given type (up to `limit`).
    pub fn by_type(&self, event_type: &EventType, limit: usize) -> Vec<&Event> {
        self.graph.store.by_type(event_type, limit)
    }

    /// Returns events from the given source (up to `limit`).
    pub fn by_source(&self, source: &ActorId, limit: usize) -> Vec<&Event> {
        self.graph.store.by_source(source, limit)
    }

    /// Returns events in the given conversation (up to `limit`).
    pub fn by_conversation(&self, conversation_id: &ConversationId, limit: usize) -> Vec<&Event> {
        self.graph.store.by_conversation(conversation_id, limit)
    }

    /// Returns causal ancestors of the given event (up to `max_depth`).
    pub fn ancestors(&self, event_id: &EventId, max_depth: usize) -> Result<Vec<&Event>> {
        self.graph.store.ancestors(event_id, max_depth)
    }

    /// Returns causal descendants of the given event (up to `max_depth`).
    pub fn descendants(&self, event_id: &EventId, max_depth: usize) -> Result<Vec<&Event>> {
        self.graph.store.descendants(event_id, max_depth)
    }

    /// Returns trust metrics for an actor.
    pub fn trust_score(&self, actor: &Actor) -> Result<TrustMetrics> {
        self.graph.trust_model.score(actor)
    }

    /// Returns directed trust from one actor to another.
    pub fn trust_between(&self, from: &Actor, to: &Actor) -> Result<TrustMetrics> {
        self.graph.trust_model.between(from, to)
    }

    /// Returns an actor by ID.
    pub fn actor(&self, id: &ActorId) -> Result<Actor> {
        self.graph.actor_store.get(id)
    }

    /// Returns the total number of events in the store.
    pub fn event_count(&self) -> usize {
        self.graph.store.count()
    }
}
