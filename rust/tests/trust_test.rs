use std::collections::BTreeMap;
use serde_json::Value;

use eventgraph::actor::{Actor, ActorStatus, ActorType};
use eventgraph::event::{create_event, NoopSigner};
use eventgraph::trust::{DefaultTrustModel, TrustConfig, TrustModel};
use eventgraph::types::*;

// ── Helpers ─────────────────────────────────────────────────────────────

fn test_actor(name: &str) -> Actor {
    let id = ActorId::new(name).unwrap();
    let pk = PublicKey::new([0u8; 32]);
    Actor::new(
        id,
        pk,
        name.to_string(),
        ActorType::AI,
        BTreeMap::new(),
        1_000_000_000,
        ActorStatus::Active,
    )
}

fn test_actor_with_key(name: &str, key_byte: u8) -> Actor {
    let id = ActorId::new(name).unwrap();
    let pk = PublicKey::new([key_byte; 32]);
    Actor::new(
        id,
        pk,
        name.to_string(),
        ActorType::AI,
        BTreeMap::new(),
        1_000_000_000,
        ActorStatus::Active,
    )
}

fn make_evidence(content: BTreeMap<String, Value>) -> eventgraph::event::Event {
    let source = ActorId::new("system").unwrap();
    let event_type = EventType::new("trust.updated").unwrap();
    let cause_id = eventgraph::event::new_event_id();
    let conv_id = ConversationId::new("conv_test").unwrap();
    let prev_hash = Hash::zero();
    let signer = NoopSigner;

    create_event(
        event_type,
        source,
        content,
        vec![cause_id],
        conv_id,
        prev_hash,
        &signer,
        1,
    )
}

fn evidence_with_current(current: f64) -> eventgraph::event::Event {
    let mut content = BTreeMap::new();
    content.insert("current".to_string(), Value::from(current));
    make_evidence(content)
}

fn plain_evidence() -> eventgraph::event::Event {
    make_evidence(BTreeMap::new())
}

// ── Tests ───────────────────────────────────────────────────────────────

#[test]
fn test_initial_score() {
    let model = DefaultTrustModel::new(TrustConfig::default());
    let actor = test_actor("alice");
    let metrics = model.score(&actor).unwrap();

    assert!((metrics.overall.value() - 0.0).abs() < f64::EPSILON);
    assert!((metrics.confidence.value() - 0.0).abs() < f64::EPSILON);
    assert!((metrics.trend.value() - 0.0).abs() < f64::EPSILON);
    assert!(metrics.evidence.is_empty());
}

#[test]
fn test_update_increases_trust() {
    let mut model = DefaultTrustModel::new(TrustConfig::default());
    let actor = test_actor("alice");

    // current=0.5 when score is 0.0 gives delta=0.5, clamped to max_adjustment=0.1
    let ev = evidence_with_current(0.5);
    let metrics = model.update(&actor, &ev).unwrap();

    assert!((metrics.overall.value() - 0.1).abs() < 1e-9);
    assert_eq!(metrics.evidence.len(), 1);
}

#[test]
fn test_update_decreases_trust() {
    let config = TrustConfig {
        initial_trust: Score::new(0.5).unwrap(),
        ..TrustConfig::default()
    };
    let mut model = DefaultTrustModel::new(config);
    let actor = test_actor("alice");

    // current=0.2 when score is 0.5 gives delta=-0.3, clamped to -0.1
    let ev = evidence_with_current(0.2);
    let metrics = model.update(&actor, &ev).unwrap();

    assert!((metrics.overall.value() - 0.4).abs() < 1e-9);
}

#[test]
fn test_update_clamped_to_max_adjustment() {
    let mut model = DefaultTrustModel::new(TrustConfig::default());
    let actor = test_actor("alice");

    // current=1.0 when score is 0.0 gives delta=1.0, clamped to 0.1
    let ev = evidence_with_current(1.0);
    let metrics = model.update(&actor, &ev).unwrap();

    assert!((metrics.overall.value() - 0.1).abs() < 1e-9);
}

#[test]
fn test_update_deduplication() {
    let mut model = DefaultTrustModel::new(TrustConfig::default());
    let actor = test_actor("alice");

    let ev = evidence_with_current(0.5);
    let m1 = model.update(&actor, &ev).unwrap();
    // Same event again should not change anything.
    let m2 = model.update(&actor, &ev).unwrap();

    assert!((m1.overall.value() - m2.overall.value()).abs() < f64::EPSILON);
    assert_eq!(m1.evidence.len(), m2.evidence.len());
}

#[test]
fn test_update_trend_positive() {
    let mut model = DefaultTrustModel::new(TrustConfig::default());
    let actor = test_actor("alice");

    let ev = evidence_with_current(0.5);
    let metrics = model.update(&actor, &ev).unwrap();

    assert!((metrics.trend.value() - 0.1).abs() < 1e-9);
}

#[test]
fn test_update_trend_negative() {
    let config = TrustConfig {
        initial_trust: Score::new(0.5).unwrap(),
        ..TrustConfig::default()
    };
    let mut model = DefaultTrustModel::new(config);
    let actor = test_actor("alice");

    let ev = evidence_with_current(0.2);
    let metrics = model.update(&actor, &ev).unwrap();

    assert!((metrics.trend.value() - (-0.1)).abs() < 1e-9);
}

#[test]
fn test_score_in_domain() {
    let mut model = DefaultTrustModel::new(TrustConfig::default());
    let actor = test_actor("alice");

    // First update to create state.
    let ev = plain_evidence();
    model.update(&actor, &ev).unwrap();

    let domain = DomainScope::new("code.review").unwrap();
    let metrics = model.score_in_domain(&actor, &domain).unwrap();

    // No domain-specific score, so falls back to overall with halved confidence.
    assert!((metrics.overall.value() - 0.01).abs() < 1e-9); // observed_event_delta
}

#[test]
fn test_score_in_domain_fallback() {
    let mut model = DefaultTrustModel::new(TrustConfig::default());
    let actor = test_actor("alice");

    // Create enough evidence for measurable confidence.
    for _ in 0..10 {
        let ev = plain_evidence();
        model.update(&actor, &ev).unwrap();
    }

    let domain = DomainScope::new("code.review").unwrap();
    let metrics_overall = model.score(&actor).unwrap();
    let metrics_domain = model.score_in_domain(&actor, &domain).unwrap();

    // Fallback: same overall score, halved confidence.
    assert!((metrics_domain.overall.value() - metrics_overall.overall.value()).abs() < 1e-9);
    assert!((metrics_domain.confidence.value() - metrics_overall.confidence.value() / 2.0).abs() < 1e-9);
}

#[test]
fn test_decay() {
    let config = TrustConfig {
        initial_trust: Score::new(0.5).unwrap(),
        decay_rate: Score::new(0.1).unwrap(),
        ..TrustConfig::default()
    };
    let mut model = DefaultTrustModel::new(config);
    let actor = test_actor("alice");

    // Force state creation.
    let ev = plain_evidence();
    model.update(&actor, &ev).unwrap();

    // Before decay, score should be around 0.51 (0.5 + 0.01 observed_event_delta).
    let before = model.score(&actor).unwrap();

    // Decay for 1 day (86400 seconds).
    let after = model.decay(&actor, 86400.0).unwrap();

    // decay = 0.1 * 1 day = 0.1
    assert!(after.overall.value() < before.overall.value());
    let expected = before.overall.value() - 0.1;
    assert!((after.overall.value() - expected).abs() < 1e-9);
}

#[test]
fn test_decay_negative_duration() {
    let config = TrustConfig {
        initial_trust: Score::new(0.5).unwrap(),
        ..TrustConfig::default()
    };
    let mut model = DefaultTrustModel::new(config);
    let actor = test_actor("alice");

    // Negative duration should return current score unchanged.
    let metrics = model.decay(&actor, -100.0).unwrap();
    assert!((metrics.overall.value() - 0.5).abs() < 1e-9);
}

#[test]
fn test_update_between() {
    let mut model = DefaultTrustModel::new(TrustConfig::default());
    let alice = test_actor_with_key("alice", 1);
    let bob = test_actor_with_key("bob", 2);

    let ev = evidence_with_current(0.5);
    let metrics = model.update_between(&alice, &bob, &ev).unwrap();

    assert!((metrics.overall.value() - 0.1).abs() < 1e-9);
    assert_eq!(metrics.actor, *bob.id());
}

#[test]
fn test_between_no_relationship() {
    let model = DefaultTrustModel::new(TrustConfig::default());
    let alice = test_actor_with_key("alice", 1);
    let bob = test_actor_with_key("bob", 2);

    let metrics = model.between(&alice, &bob).unwrap();

    assert!((metrics.overall.value() - 0.0).abs() < f64::EPSILON);
    assert!(metrics.evidence.is_empty());
}

#[test]
fn test_between_asymmetric() {
    let mut model = DefaultTrustModel::new(TrustConfig::default());
    let alice = test_actor_with_key("alice", 1);
    let bob = test_actor_with_key("bob", 2);

    // Alice trusts Bob.
    let ev = evidence_with_current(0.8);
    model.update_between(&alice, &bob, &ev).unwrap();

    // Alice -> Bob should have trust.
    let a_to_b = model.between(&alice, &bob).unwrap();
    assert!(a_to_b.overall.value() > 0.0);

    // Bob -> Alice should be initial (0.0), asymmetric.
    let b_to_a = model.between(&bob, &alice).unwrap();
    assert!((b_to_a.overall.value() - 0.0).abs() < f64::EPSILON);
}

#[test]
fn test_evidence_capped_at_100() {
    let mut model = DefaultTrustModel::new(TrustConfig::default());
    let actor = test_actor("alice");

    for _ in 0..110 {
        let ev = plain_evidence();
        model.update(&actor, &ev).unwrap();
    }

    let metrics = model.score(&actor).unwrap();
    assert_eq!(metrics.evidence.len(), 100);
}

#[test]
fn test_decay_directed_trust() {
    let config = TrustConfig {
        initial_trust: Score::new(0.0).unwrap(),
        decay_rate: Score::new(0.1).unwrap(),
        ..TrustConfig::default()
    };
    let mut model = DefaultTrustModel::new(config);
    let alice = test_actor_with_key("alice", 1);
    let bob = test_actor_with_key("bob", 2);

    // Build up directed trust.
    let ev = evidence_with_current(0.5);
    model.update_between(&alice, &bob, &ev).unwrap();

    let before = model.between(&alice, &bob).unwrap();

    // Decay bob's trust for 1 day — should also affect directed trust targeting bob.
    model.decay(&bob, 86400.0).unwrap();

    let after = model.between(&alice, &bob).unwrap();
    assert!(after.overall.value() < before.overall.value());
}
