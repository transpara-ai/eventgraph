use std::collections::{BTreeMap, HashMap};
use std::sync::RwLock;

use crate::actor::Actor;
use crate::errors::Result;
use crate::event::Event;
use crate::types::{ActorId, DomainScope, EventId, Score, Weight};

// ── TrustMetrics ────────────────────────────────────────────────────────

/// Trust metrics for an actor, optionally scoped to a domain.
#[derive(Debug, Clone)]
pub struct TrustMetrics {
    pub actor: ActorId,
    pub overall: Score,
    pub by_domain: BTreeMap<String, Score>,
    pub confidence: Score,
    pub trend: Weight,
    pub evidence: Vec<EventId>,
    pub last_updated_nanos: u64,
    pub decay_rate: Score,
}

// ── TrustConfig ─────────────────────────────────────────────────────────

/// Configuration for the default trust model.
#[derive(Debug, Clone)]
pub struct TrustConfig {
    pub initial_trust: Score,
    pub decay_rate: Score,
    pub max_adjustment: Weight,
    pub observed_event_delta: f64,
    pub trend_decay_rate: f64,
}

impl Default for TrustConfig {
    fn default() -> Self {
        Self {
            initial_trust: Score::new(0.0).unwrap(),
            decay_rate: Score::new(0.01).unwrap(),
            max_adjustment: Weight::new(0.1).unwrap(),
            observed_event_delta: 0.01,
            trend_decay_rate: 0.01,
        }
    }
}

// ── TrustModel trait ────────────────────────────────────────────────────

/// Extension point for trust computation strategies.
pub trait TrustModel {
    /// Returns the current trust metrics for an actor.
    fn score(&self, actor: &Actor) -> Result<TrustMetrics>;

    /// Returns trust metrics for an actor scoped to a domain.
    fn score_in_domain(&self, actor: &Actor, domain: &DomainScope) -> Result<TrustMetrics>;

    /// Updates trust for an actor based on observed evidence.
    fn update(&mut self, actor: &Actor, evidence: &Event) -> Result<TrustMetrics>;

    /// Updates directed trust from one actor to another based on evidence.
    fn update_between(&mut self, from: &Actor, to: &Actor, evidence: &Event) -> Result<TrustMetrics>;

    /// Applies time-based decay to an actor's trust.
    fn decay(&mut self, actor: &Actor, elapsed_seconds: f64) -> Result<TrustMetrics>;

    /// Returns the directed trust metrics from one actor to another.
    fn between(&self, from: &Actor, to: &Actor) -> Result<TrustMetrics>;
}

// ── Internal state ──────────────────────────────────────────────────────

#[derive(Debug, Clone)]
struct TrustState {
    score: f64,
    trend: f64,
    evidence: Vec<EventId>,
    last_updated_nanos: u64,
    by_domain: BTreeMap<String, f64>,
}

impl TrustState {
    fn new(initial: f64) -> Self {
        Self {
            score: initial,
            trend: 0.0,
            evidence: Vec::new(),
            last_updated_nanos: 0,
            by_domain: BTreeMap::new(),
        }
    }
}

// ── DefaultTrustModel ───────────────────────────────────────────────────

/// Thread-safe default implementation of the trust model.
///
/// Uses linear decay, clamped adjustments, evidence deduplication,
/// and directed trust relationships.
pub struct DefaultTrustModel {
    config: TrustConfig,
    state: RwLock<TrustModelState>,
}

struct TrustModelState {
    /// Per-actor trust state, keyed by actor ID string.
    per_actor: HashMap<String, TrustState>,
    /// Directed trust state, keyed by (from_id, to_id).
    directed: HashMap<(String, String), TrustState>,
}

impl DefaultTrustModel {
    /// Creates a new DefaultTrustModel with the given configuration.
    pub fn new(config: TrustConfig) -> Self {
        Self {
            config,
            state: RwLock::new(TrustModelState {
                per_actor: HashMap::new(),
                directed: HashMap::new(),
            }),
        }
    }

    fn metrics_from_state(&self, actor_id: &ActorId, ts: &TrustState) -> TrustMetrics {
        let evidence_count = ts.evidence.len();
        let confidence_val = (evidence_count as f64 / 50.0).min(1.0);

        let by_domain: BTreeMap<String, Score> = ts
            .by_domain
            .iter()
            .map(|(k, v)| (k.clone(), Score::new(*v).unwrap()))
            .collect();

        TrustMetrics {
            actor: actor_id.clone(),
            overall: Score::new(ts.score).unwrap(),
            by_domain,
            confidence: Score::new(confidence_val).unwrap(),
            trend: Weight::new(ts.trend).unwrap(),
            evidence: ts.evidence.clone(),
            last_updated_nanos: ts.last_updated_nanos,
            decay_rate: self.config.decay_rate,
        }
    }

    fn get_or_init<'a>(map: &'a mut HashMap<String, TrustState>, key: &str, initial: f64) -> &'a mut TrustState {
        map.entry(key.to_string()).or_insert_with(|| TrustState::new(initial))
    }

    fn apply_evidence(ts: &mut TrustState, evidence: &Event, config: &TrustConfig) {
        // Deduplication: skip if this event was already recorded.
        let eid = evidence.id.clone();
        if ts.evidence.iter().any(|e| e == &eid) {
            return;
        }

        // Determine delta from evidence content.
        let content = evidence.content();
        let delta = if let Some(current_val) = content.get("current") {
            if let Some(current_f64) = current_val.as_f64() {
                current_f64 - ts.score
            } else {
                config.observed_event_delta
            }
        } else {
            config.observed_event_delta
        };

        // Clamp to max_adjustment.
        let max_adj = config.max_adjustment.value();
        let clamped = delta.clamp(-max_adj, max_adj);

        // Apply and clamp score to [0, 1].
        ts.score = (ts.score + clamped).clamp(0.0, 1.0);

        // Update trend: nudge by ±0.1 based on delta sign, clamped to [-1, 1].
        if delta > 0.0 {
            ts.trend = (ts.trend + 0.1).clamp(-1.0, 1.0);
        } else if delta < 0.0 {
            ts.trend = (ts.trend - 0.1).clamp(-1.0, 1.0);
        }

        ts.last_updated_nanos = evidence.timestamp_nanos;

        // Cap evidence at 100.
        if ts.evidence.len() >= 100 {
            ts.evidence.remove(0);
        }
        ts.evidence.push(eid);
    }
}

impl TrustModel for DefaultTrustModel {
    fn score(&self, actor: &Actor) -> Result<TrustMetrics> {
        let state = self.state.read().unwrap();
        let initial = self.config.initial_trust.value();
        let ts = state.per_actor.get(actor.id().value());
        match ts {
            Some(ts) => Ok(self.metrics_from_state(actor.id(), ts)),
            None => {
                let default = TrustState::new(initial);
                Ok(self.metrics_from_state(actor.id(), &default))
            }
        }
    }

    fn score_in_domain(&self, actor: &Actor, domain: &DomainScope) -> Result<TrustMetrics> {
        let state = self.state.read().unwrap();
        let initial = self.config.initial_trust.value();
        let ts = state.per_actor.get(actor.id().value());

        match ts {
            Some(ts) => {
                let domain_key = domain.value().to_string();
                if let Some(&domain_score) = ts.by_domain.get(&domain_key) {
                    // Domain-specific score exists.
                    let evidence_count = ts.evidence.len();
                    let confidence_val = (evidence_count as f64 / 50.0).min(1.0);

                    Ok(TrustMetrics {
                        actor: actor.id().clone(),
                        overall: Score::new(domain_score).unwrap(),
                        by_domain: ts.by_domain.iter()
                            .map(|(k, v)| (k.clone(), Score::new(*v).unwrap()))
                            .collect(),
                        confidence: Score::new(confidence_val).unwrap(),
                        trend: Weight::new(ts.trend).unwrap(),
                        evidence: ts.evidence.clone(),
                        last_updated_nanos: ts.last_updated_nanos,
                        decay_rate: self.config.decay_rate,
                    })
                } else {
                    // Fallback: use overall score with halved confidence.
                    let evidence_count = ts.evidence.len();
                    let confidence_val = ((evidence_count as f64 / 50.0).min(1.0)) / 2.0;

                    Ok(TrustMetrics {
                        actor: actor.id().clone(),
                        overall: Score::new(ts.score).unwrap(),
                        by_domain: ts.by_domain.iter()
                            .map(|(k, v)| (k.clone(), Score::new(*v).unwrap()))
                            .collect(),
                        confidence: Score::new(confidence_val).unwrap(),
                        trend: Weight::new(ts.trend).unwrap(),
                        evidence: ts.evidence.clone(),
                        last_updated_nanos: ts.last_updated_nanos,
                        decay_rate: self.config.decay_rate,
                    })
                }
            }
            None => {
                // No state at all: return initial with zero confidence.
                let default = TrustState::new(initial);
                Ok(self.metrics_from_state(actor.id(), &default))
            }
        }
    }

    fn update(&mut self, actor: &Actor, evidence: &Event) -> Result<TrustMetrics> {
        let mut state = self.state.write().unwrap();
        let initial = self.config.initial_trust.value();
        let config = self.config.clone();
        let ts = Self::get_or_init(&mut state.per_actor, actor.id().value(), initial);
        Self::apply_evidence(ts, evidence, &config);
        let ts_clone = ts.clone();
        drop(state);

        Ok(self.metrics_from_state(actor.id(), &ts_clone))
    }

    fn update_between(&mut self, from: &Actor, to: &Actor, evidence: &Event) -> Result<TrustMetrics> {
        let mut state = self.state.write().unwrap();
        let initial = self.config.initial_trust.value();
        let config = self.config.clone();
        let key = (from.id().value().to_string(), to.id().value().to_string());
        let ts = state.directed.entry(key).or_insert_with(|| TrustState::new(initial));
        Self::apply_evidence(ts, evidence, &config);
        let ts_clone = ts.clone();
        drop(state);

        Ok(self.metrics_from_state(to.id(), &ts_clone))
    }

    fn decay(&mut self, actor: &Actor, elapsed_seconds: f64) -> Result<TrustMetrics> {
        if elapsed_seconds < 0.0 {
            return self.score(actor);
        }

        let mut state = self.state.write().unwrap();
        let initial = self.config.initial_trust.value();
        let decay_rate = self.config.decay_rate.value();
        let days = elapsed_seconds / 86400.0;

        // Decay per-actor trust.
        {
            let ts = Self::get_or_init(&mut state.per_actor, actor.id().value(), initial);
            ts.score = (ts.score - decay_rate * days).clamp(0.0, 1.0);
        }

        // Also decay directed trust where this actor is the target.
        let actor_id_str = actor.id().value().to_string();
        for ((_, to), dts) in state.directed.iter_mut() {
            if to == &actor_id_str {
                dts.score = (dts.score - decay_rate * days).clamp(0.0, 1.0);
            }
        }

        let ts_clone = state.per_actor.get(actor.id().value()).unwrap().clone();
        drop(state);

        Ok(self.metrics_from_state(actor.id(), &ts_clone))
    }

    fn between(&self, from: &Actor, to: &Actor) -> Result<TrustMetrics> {
        let state = self.state.read().unwrap();
        let initial = self.config.initial_trust.value();
        let key = (from.id().value().to_string(), to.id().value().to_string());
        match state.directed.get(&key) {
            Some(ts) => Ok(self.metrics_from_state(to.id(), ts)),
            None => {
                let default = TrustState::new(initial);
                Ok(self.metrics_from_state(to.id(), &default))
            }
        }
    }
}
