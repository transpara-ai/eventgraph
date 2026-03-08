use std::collections::BTreeMap;

use serde_json::Value;

use eventgraph::actor::{Actor, ActorStatus, ActorStore, ActorType, InMemoryActorStore};
use eventgraph::decision::AuthorityLevel;
use eventgraph::event::Signer;
use eventgraph::graph::{Graph, GraphConfig};
use eventgraph::store::{InMemoryStore, Store};
use eventgraph::types::{
    ActorId, ConversationId, EventType, PublicKey, Signature,
};

// ── Helpers ──────────────────────────────────────────────────────────

struct TestSigner;

impl Signer for TestSigner {
    fn sign(&self, data: &[u8]) -> Signature {
        let mut buf = vec![0u8; 64];
        let len = data.len().min(64);
        buf[..len].copy_from_slice(&data[..len]);
        Signature::new(buf).unwrap()
    }
}

fn test_actor(name: &str, key_byte: u8) -> Actor {
    let mut key = [0u8; 32];
    key[0] = key_byte;
    Actor::new(
        ActorId::new(name).unwrap(),
        PublicKey::new(key),
        name.to_string(),
        ActorType::Human,
        BTreeMap::new(),
        1_000_000_000,
        ActorStatus::Active,
    )
}

fn new_test_graph() -> (Graph, ActorId) {
    let store = InMemoryStore::new();
    let actor_store = InMemoryActorStore::new();
    let g = Graph::new(store, actor_store);
    g.start().expect("start should succeed");
    let actor_id = ActorId::new("actor_system0000000000000000000001").unwrap();
    (g, actor_id)
}

// ── 1. TestNewGraph ─────────────────────────────────────────────────

#[test]
fn test_new_graph() {
    let store = InMemoryStore::new();
    let actor_store = InMemoryActorStore::new();
    let g = Graph::new(store, actor_store);

    assert_eq!(g.store().count(), 0);
    assert_eq!(g.actor_store().actor_count(), 0);
}

// ── 2. TestBootstrap ────────────────────────────────────────────────

#[test]
fn test_bootstrap() {
    let (mut g, actor_id) = new_test_graph();

    let ev = g.bootstrap(actor_id.clone(), Some(&TestSigner)).unwrap();
    assert_eq!(ev.event_type.value(), "system.bootstrapped");
    assert_eq!(ev.source.value(), actor_id.value());
    assert_eq!(g.store().count(), 1);
}

// ── 3. TestBootstrapDoubleReject ────────────────────────────────────

#[test]
fn test_bootstrap_double_reject() {
    let (mut g, actor_id) = new_test_graph();

    g.bootstrap(actor_id.clone(), Some(&TestSigner)).unwrap();
    let result = g.bootstrap(actor_id, Some(&TestSigner));
    assert!(result.is_err(), "second bootstrap should fail");
}

// ── 4. TestRecord ───────────────────────────────────────────────────

#[test]
fn test_record() {
    let (mut g, actor_id) = new_test_graph();

    let bootstrap = g.bootstrap(actor_id.clone(), Some(&TestSigner)).unwrap();

    let mut content = BTreeMap::new();
    content.insert("key".to_string(), Value::String("value".to_string()));

    let ev = g
        .record(
            EventType::new("trust.updated").unwrap(),
            actor_id,
            content,
            vec![bootstrap.id.clone()],
            ConversationId::new("conv_test000000000000000000000001").unwrap(),
            Some(&TestSigner),
        )
        .unwrap();

    assert_eq!(ev.event_type.value(), "trust.updated");
    assert_eq!(g.store().count(), 2);
}

// ── 5. TestRecordAfterClose ─────────────────────────────────────────

#[test]
fn test_record_after_close() {
    let (mut g, actor_id) = new_test_graph();

    let bootstrap = g.bootstrap(actor_id.clone(), Some(&TestSigner)).unwrap();
    g.close();

    let result = g.record(
        EventType::new("trust.updated").unwrap(),
        actor_id,
        BTreeMap::new(),
        vec![bootstrap.id.clone()],
        ConversationId::new("conv_test000000000000000000000001").unwrap(),
        Some(&TestSigner),
    );
    assert!(result.is_err(), "record after close should fail");
}

// ── 6. TestEvaluate ─────────────────────────────────────────────────

#[test]
fn test_evaluate() {
    let (g, _actor_id) = new_test_graph();

    let alice = test_actor("alice", 1);
    let result = g.evaluate(&alice, "test.action").unwrap();

    // Default authority chain returns Notification for unknown actions
    assert_eq!(result.level, AuthorityLevel::Notification);
}

// ── 7. TestEvaluateAfterClose ───────────────────────────────────────

#[test]
fn test_evaluate_after_close() {
    let (mut g, _actor_id) = new_test_graph();
    g.close();

    let alice = test_actor("alice", 1);
    let result = g.evaluate(&alice, "test.action");
    assert!(result.is_err(), "evaluate after close should fail");
}

// ── 8. TestQuery ────────────────────────────────────────────────────

#[test]
fn test_query_event_count_and_recent() {
    let (mut g, actor_id) = new_test_graph();

    g.bootstrap(actor_id, Some(&TestSigner)).unwrap();

    let q = g.query().unwrap();
    assert_eq!(q.event_count(), 1);

    let recent = q.recent(10);
    assert_eq!(recent.len(), 1);
    assert_eq!(recent[0].event_type.value(), "system.bootstrapped");
}

// ── 9. TestQueryByType ──────────────────────────────────────────────

#[test]
fn test_query_by_type() {
    let (mut g, actor_id) = new_test_graph();
    g.bootstrap(actor_id, Some(&TestSigner)).unwrap();

    let q = g.query().unwrap();
    let et = EventType::new("system.bootstrapped").unwrap();
    let events = q.by_type(&et, 10);
    assert_eq!(events.len(), 1);
}

// ── 10. TestQueryBySource ───────────────────────────────────────────

#[test]
fn test_query_by_source() {
    let (mut g, actor_id) = new_test_graph();
    g.bootstrap(actor_id.clone(), Some(&TestSigner)).unwrap();

    let q = g.query().unwrap();
    let events = q.by_source(&actor_id, 10);
    assert_eq!(events.len(), 1);
}

// ── 11. TestStartIdempotent ─────────────────────────────────────────

#[test]
fn test_start_idempotent() {
    let store = InMemoryStore::new();
    let actor_store = InMemoryActorStore::new();
    let g = Graph::new(store, actor_store);

    g.start().unwrap();
    g.start().unwrap(); // second start should be a no-op
}

// ── 12. TestCloseIdempotent ─────────────────────────────────────────

#[test]
fn test_close_idempotent() {
    let store = InMemoryStore::new();
    let actor_store = InMemoryActorStore::new();
    let mut g = Graph::new(store, actor_store);

    g.close();
    g.close(); // second close should be a no-op
}

// ── 13. TestQueryAncestorsDescendants ───────────────────────────────

#[test]
fn test_query_ancestors_and_descendants() {
    let (mut g, actor_id) = new_test_graph();

    let bootstrap = g.bootstrap(actor_id.clone(), Some(&TestSigner)).unwrap();

    let mut content = BTreeMap::new();
    content.insert("key".to_string(), Value::String("val".to_string()));

    let child = g
        .record(
            EventType::new("trust.updated").unwrap(),
            actor_id,
            content,
            vec![bootstrap.id.clone()],
            ConversationId::new("conv_test000000000000000000000001").unwrap(),
            Some(&TestSigner),
        )
        .unwrap();

    let q = g.query().unwrap();

    // Ancestors of child should include bootstrap
    let ancestors = q.ancestors(&child.id, 10).unwrap();
    assert!(
        ancestors.iter().any(|e| e.id == bootstrap.id),
        "expected bootstrap in ancestors"
    );

    // Descendants of bootstrap should include child
    let descendants = q.descendants(&bootstrap.id, 10).unwrap();
    assert!(
        descendants.iter().any(|e| e.id == child.id),
        "expected child in descendants"
    );
}

// ── 14. TestQueryTrustScore ─────────────────────────────────────────

#[test]
fn test_query_trust_score() {
    let (g, _) = new_test_graph();

    let alice = test_actor("alice", 1);

    let q = g.query().unwrap();
    let metrics = q.trust_score(&alice).unwrap();
    assert_eq!(metrics.overall.value(), 0.0, "initial trust should be 0.0");
}

// ── 15. TestQueryTrustBetween ───────────────────────────────────────

#[test]
fn test_query_trust_between() {
    let (g, _) = new_test_graph();

    let alice = test_actor("alice", 1);
    let bob = test_actor("bob", 2);

    let q = g.query().unwrap();
    let metrics = q.trust_between(&alice, &bob).unwrap();
    assert_eq!(
        metrics.overall.value(),
        0.0,
        "initial trust between should be 0.0"
    );
}

// ── 16. TestQueryActor ──────────────────────────────────────────────

#[test]
fn test_query_actor() {
    let store = InMemoryStore::new();
    let mut actor_store = InMemoryActorStore::new();
    let pk = PublicKey::new([3u8; 32]);
    let registered = actor_store
        .register(pk, "Charlie", ActorType::Human)
        .unwrap();

    let g = Graph::new(store, actor_store);
    g.start().unwrap();

    let q = g.query().unwrap();
    let found = q.actor(registered.id()).unwrap();
    assert_eq!(found.id().value(), registered.id().value());
    assert_eq!(found.display_name(), "Charlie");
}

// ── 17. TestRecordBeforeStart ───────────────────────────────────────

#[test]
fn test_record_before_start() {
    let store = InMemoryStore::new();
    let actor_store = InMemoryActorStore::new();
    let mut g = Graph::new(store, actor_store);
    // Do NOT call start

    let actor_id = ActorId::new("actor_system0000000000000000000001").unwrap();
    let result = g.bootstrap(actor_id, Some(&TestSigner));
    assert!(result.is_err(), "bootstrap before start should fail");
}

// ── 18. TestWithCustomConfig ────────────────────────────────────────

#[test]
fn test_with_custom_config() {
    let store = InMemoryStore::new();
    let actor_store = InMemoryActorStore::new();
    let config = GraphConfig {
        subscriber_buffer_size: 512,
        fallback_to_mechanical: false,
    };
    let g = Graph::with_config(store, actor_store, config);
    assert_eq!(g.store().count(), 0);
}
