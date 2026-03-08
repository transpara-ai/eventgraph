use sha2::{Sha256, Digest};
use std::collections::BTreeMap;
use std::fmt;
use std::sync::RwLock;

use crate::errors::{EventGraphError, Result};
use crate::types::{ActorId, PublicKey};

// ── ActorType ─────────────────────────────────────────────────────────

/// What kind of decision-maker an actor is.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ActorType {
    Human,
    AI,
    System,
    Committee,
    RulesEngine,
}

impl fmt::Display for ActorType {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::Human => f.write_str("Human"),
            Self::AI => f.write_str("AI"),
            Self::System => f.write_str("System"),
            Self::Committee => f.write_str("Committee"),
            Self::RulesEngine => f.write_str("RulesEngine"),
        }
    }
}

// ── ActorStatus ───────────────────────────────────────────────────────

/// Actor lifecycle status. Memorial is terminal.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ActorStatus {
    Active,
    Suspended,
    Memorial,
}

impl ActorStatus {
    /// Attempt to transition to the target status.
    /// Returns the new status on success, or an error if the transition is invalid.
    pub fn transition_to(&self, target: ActorStatus) -> Result<ActorStatus> {
        let valid = self.valid_transitions();
        if valid.contains(&target) {
            Ok(target)
        } else {
            Err(EventGraphError::InvalidTransition {
                from: self.to_string(),
                to: target.to_string(),
            })
        }
    }

    /// Returns the list of valid target statuses from this status.
    pub fn valid_transitions(&self) -> Vec<ActorStatus> {
        match self {
            Self::Active => vec![Self::Suspended, Self::Memorial],
            Self::Suspended => vec![Self::Active, Self::Memorial],
            Self::Memorial => vec![],
        }
    }
}

impl fmt::Display for ActorStatus {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::Active => f.write_str("Active"),
            Self::Suspended => f.write_str("Suspended"),
            Self::Memorial => f.write_str("Memorial"),
        }
    }
}

// ── Actor ─────────────────────────────────────────────────────────────

/// An immutable identity in the system.
#[derive(Debug, Clone)]
pub struct Actor {
    id: ActorId,
    public_key: PublicKey,
    display_name: String,
    actor_type: ActorType,
    metadata: BTreeMap<String, serde_json::Value>,
    created_at_nanos: u64,
    status: ActorStatus,
}

impl Actor {
    /// Creates a new Actor. Metadata is deep-cloned on construction.
    pub fn new(
        id: ActorId,
        public_key: PublicKey,
        display_name: String,
        actor_type: ActorType,
        metadata: BTreeMap<String, serde_json::Value>,
        created_at_nanos: u64,
        status: ActorStatus,
    ) -> Self {
        Self {
            id,
            public_key,
            display_name,
            actor_type,
            metadata,
            created_at_nanos,
            status,
        }
    }

    pub fn id(&self) -> &ActorId { &self.id }
    pub fn public_key(&self) -> &PublicKey { &self.public_key }
    pub fn display_name(&self) -> &str { &self.display_name }
    pub fn actor_type(&self) -> ActorType { self.actor_type }
    pub fn metadata(&self) -> BTreeMap<String, serde_json::Value> { self.metadata.clone() }
    pub fn created_at_nanos(&self) -> u64 { self.created_at_nanos }
    pub fn status(&self) -> ActorStatus { self.status }

    fn with_status(&self, status: ActorStatus) -> Self {
        Self {
            id: self.id.clone(),
            public_key: self.public_key.clone(),
            display_name: self.display_name.clone(),
            actor_type: self.actor_type,
            metadata: self.metadata.clone(),
            created_at_nanos: self.created_at_nanos,
            status,
        }
    }

    fn with_updates(&self, update: &ActorUpdate) -> Self {
        let mut result = Self {
            id: self.id.clone(),
            public_key: self.public_key.clone(),
            display_name: self.display_name.clone(),
            actor_type: self.actor_type,
            metadata: self.metadata.clone(),
            created_at_nanos: self.created_at_nanos,
            status: self.status,
        };
        if let Some(ref name) = update.display_name {
            result.display_name = name.clone();
        }
        if let Some(ref md) = update.metadata {
            for (k, v) in md {
                result.metadata.insert(k.clone(), v.clone());
            }
        }
        result
    }
}

// ── ActorUpdate ───────────────────────────────────────────────────────

/// Describes updates to apply to an actor.
pub struct ActorUpdate {
    pub display_name: Option<String>,
    pub metadata: Option<BTreeMap<String, serde_json::Value>>,
}

// ── ActorFilter ───────────────────────────────────────────────────────

/// Describes criteria for listing actors.
pub struct ActorFilter {
    pub status: Option<ActorStatus>,
    pub actor_type: Option<ActorType>,
    pub limit: usize,
    pub after: Option<String>,
}

// ── ActorPage ─────────────────────────────────────────────────────────

/// A page of actor results with cursor-based pagination.
pub struct ActorPage {
    pub items: Vec<Actor>,
    pub cursor: Option<String>,
    pub has_more: bool,
}

// ── ActorStore trait ──────────────────────────────────────────────────

/// Actor persistence interface.
pub trait ActorStore {
    fn register(&mut self, public_key: PublicKey, display_name: &str, actor_type: ActorType) -> Result<Actor>;
    fn get(&self, id: &ActorId) -> Result<Actor>;
    fn get_by_public_key(&self, public_key: &PublicKey) -> Result<Actor>;
    fn update(&mut self, id: &ActorId, updates: &ActorUpdate) -> Result<Actor>;
    fn list(&self, filter: &ActorFilter) -> Result<ActorPage>;
    fn suspend(&mut self, id: &ActorId, reason: &str) -> Result<Actor>;
    fn reactivate(&mut self, id: &ActorId, reason: &str) -> Result<Actor>;
    fn memorial(&mut self, id: &ActorId, reason: &str) -> Result<Actor>;
}

// ── InMemoryActorStore ───────────────────────────────────────────────

/// In-memory implementation of ActorStore. Thread-safe via RwLock.
pub struct InMemoryActorStore {
    inner: RwLock<InMemoryActorStoreInner>,
}

struct InMemoryActorStoreInner {
    actors: BTreeMap<String, Actor>,      // actor_id string → Actor
    by_key: BTreeMap<String, String>,     // hex(publicKey) → actor_id string
    ordered: Vec<String>,                 // insertion order for pagination
}

fn pub_key_hex(pk: &PublicKey) -> String {
    pk.bytes().iter().map(|b| format!("{b:02x}")).collect()
}

fn derive_actor_id(pk: &PublicKey) -> ActorId {
    let mut hasher = Sha256::new();
    hasher.update(pk.bytes());
    let result = hasher.finalize();
    let hex: String = result[..16].iter().map(|b| format!("{b:02x}")).collect();
    let id = format!("actor_{hex}");
    ActorId::new(id).expect("derived actor ID is always valid")
}

fn now_nanos() -> u64 {
    std::time::SystemTime::now()
        .duration_since(std::time::UNIX_EPOCH)
        .unwrap()
        .as_nanos() as u64
}

impl InMemoryActorStore {
    /// Creates a new empty InMemoryActorStore.
    pub fn new() -> Self {
        Self {
            inner: RwLock::new(InMemoryActorStoreInner {
                actors: BTreeMap::new(),
                by_key: BTreeMap::new(),
                ordered: Vec::new(),
            }),
        }
    }

    /// Returns the number of registered actors. For testing.
    pub fn actor_count(&self) -> usize {
        let inner = self.inner.read().unwrap();
        inner.actors.len()
    }
}

impl Default for InMemoryActorStore {
    fn default() -> Self { Self::new() }
}

impl ActorStore for InMemoryActorStore {
    fn register(&mut self, public_key: PublicKey, display_name: &str, actor_type: ActorType) -> Result<Actor> {
        let inner = self.inner.get_mut().unwrap();
        let key_hex = pub_key_hex(&public_key);

        // Idempotent: return existing actor if public key is already registered
        if let Some(existing_id) = inner.by_key.get(&key_hex) {
            let actor = inner.actors.get(existing_id).unwrap().clone();
            return Ok(actor);
        }

        let id = derive_actor_id(&public_key);
        let actor = Actor::new(
            id.clone(),
            public_key,
            display_name.to_string(),
            actor_type,
            BTreeMap::new(),
            now_nanos(),
            ActorStatus::Active,
        );

        let id_str = id.value().to_string();
        inner.actors.insert(id_str.clone(), actor.clone());
        inner.by_key.insert(key_hex, id_str.clone());
        inner.ordered.push(id_str);

        Ok(actor)
    }

    fn get(&self, id: &ActorId) -> Result<Actor> {
        let inner = self.inner.read().unwrap();
        inner.actors.get(id.value())
            .cloned()
            .ok_or_else(|| EventGraphError::ActorNotFound {
                actor_id: id.value().to_string(),
            })
    }

    fn get_by_public_key(&self, public_key: &PublicKey) -> Result<Actor> {
        let inner = self.inner.read().unwrap();
        let key_hex = pub_key_hex(public_key);
        let id = inner.by_key.get(&key_hex)
            .ok_or_else(|| EventGraphError::ActorKeyNotFound {
                key_hex: key_hex.clone(),
            })?;
        Ok(inner.actors.get(id).unwrap().clone())
    }

    fn update(&mut self, id: &ActorId, updates: &ActorUpdate) -> Result<Actor> {
        let inner = self.inner.get_mut().unwrap();
        let id_str = id.value().to_string();
        let actor = inner.actors.get(&id_str)
            .ok_or_else(|| EventGraphError::ActorNotFound {
                actor_id: id_str.clone(),
            })?;
        let updated = actor.with_updates(updates);
        inner.actors.insert(id_str, updated.clone());
        Ok(updated)
    }

    fn list(&self, filter: &ActorFilter) -> Result<ActorPage> {
        let inner = self.inner.read().unwrap();

        let limit = if filter.limit == 0 { 100 } else { filter.limit };

        // Find start position
        let start_idx = if let Some(ref after) = filter.after {
            let pos = inner.ordered.iter().position(|id| id == after);
            match pos {
                Some(i) => i + 1,
                None => return Ok(ActorPage { items: vec![], cursor: None, has_more: false }),
            }
        } else {
            0
        };

        let mut items = Vec::new();
        for i in start_idx..inner.ordered.len() {
            if items.len() >= limit {
                break;
            }
            let actor = inner.actors.get(&inner.ordered[i]).unwrap();
            if let Some(ref status) = filter.status {
                if actor.status() != *status {
                    continue;
                }
            }
            if let Some(ref at) = filter.actor_type {
                if actor.actor_type() != *at {
                    continue;
                }
            }
            items.push(actor.clone());
        }

        let mut has_more = false;
        let mut cursor = None;

        if items.len() == limit {
            // Check if there are more matching items after the last one
            let last_id = items.last().unwrap().id().value();
            let last_idx = inner.ordered.iter().position(|id| id == last_id).unwrap();
            for i in (last_idx + 1)..inner.ordered.len() {
                let actor = inner.actors.get(&inner.ordered[i]).unwrap();
                if let Some(ref status) = filter.status {
                    if actor.status() != *status {
                        continue;
                    }
                }
                if let Some(ref at) = filter.actor_type {
                    if actor.actor_type() != *at {
                        continue;
                    }
                }
                has_more = true;
                break;
            }
            if has_more {
                cursor = Some(last_id.to_string());
            }
        }

        Ok(ActorPage { items, cursor, has_more })
    }

    fn suspend(&mut self, id: &ActorId, _reason: &str) -> Result<Actor> {
        let inner = self.inner.get_mut().unwrap();
        let id_str = id.value().to_string();
        let actor = inner.actors.get(&id_str)
            .ok_or_else(|| EventGraphError::ActorNotFound {
                actor_id: id_str.clone(),
            })?;
        let new_status = actor.status().transition_to(ActorStatus::Suspended)?;
        let updated = actor.with_status(new_status);
        inner.actors.insert(id_str, updated.clone());
        Ok(updated)
    }

    fn reactivate(&mut self, id: &ActorId, _reason: &str) -> Result<Actor> {
        let inner = self.inner.get_mut().unwrap();
        let id_str = id.value().to_string();
        let actor = inner.actors.get(&id_str)
            .ok_or_else(|| EventGraphError::ActorNotFound {
                actor_id: id_str.clone(),
            })?;
        let new_status = actor.status().transition_to(ActorStatus::Active)?;
        let updated = actor.with_status(new_status);
        inner.actors.insert(id_str, updated.clone());
        Ok(updated)
    }

    fn memorial(&mut self, id: &ActorId, _reason: &str) -> Result<Actor> {
        let inner = self.inner.get_mut().unwrap();
        let id_str = id.value().to_string();
        let actor = inner.actors.get(&id_str)
            .ok_or_else(|| EventGraphError::ActorNotFound {
                actor_id: id_str.clone(),
            })?;
        let new_status = actor.status().transition_to(ActorStatus::Memorial)?;
        let updated = actor.with_status(new_status);
        inner.actors.insert(id_str, updated.clone());
        Ok(updated)
    }
}
