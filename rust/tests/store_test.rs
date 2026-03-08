use eventgraph::event::*;
use eventgraph::store::*;
use eventgraph::types::*;

fn bootstrap() -> Event {
    create_bootstrap(ActorId::new("alice").unwrap(), &NoopSigner, 1)
}

#[test]
fn append_and_get() {
    let mut store = InMemoryStore::new();
    let ev = bootstrap();
    let id = ev.id.clone();
    store.append(ev).unwrap();
    assert_eq!(store.get(&id).unwrap().id, id);
}

#[test]
fn head_returns_latest() {
    let mut store = InMemoryStore::new();
    assert!(store.head().is_none());
    let ev = bootstrap();
    store.append(ev).unwrap();
    assert!(store.head().is_some());
}

#[test]
fn count_increments() {
    let mut store = InMemoryStore::new();
    assert_eq!(store.count(), 0);
    store.append(bootstrap()).unwrap();
    assert_eq!(store.count(), 1);
}

#[test]
fn chain_integrity_enforced() {
    let mut store = InMemoryStore::new();
    let boot = bootstrap();
    store.append(boot).unwrap();

    // Create a second bootstrap (wrong prev_hash)
    let bad = bootstrap();
    assert!(store.append(bad).is_err());
}

#[test]
fn verify_chain_valid() {
    let mut store = InMemoryStore::new();
    let boot = bootstrap();
    store.append(boot).unwrap();
    let v = store.verify_chain();
    assert!(v.valid);
    assert_eq!(v.length, 1);
}

#[test]
fn get_nonexistent_fails() {
    let store = InMemoryStore::new();
    let id = EventId::new("019462a0-0000-7000-8000-000000000001").unwrap();
    assert!(store.get(&id).is_err());
}

#[test]
fn chained_append() {
    let signer = NoopSigner;
    let mut store = InMemoryStore::new();
    let boot = bootstrap();
    let prev_hash = boot.hash.clone();
    let boot_id = boot.id.clone();
    store.append(boot).unwrap();

    let mut content = std::collections::BTreeMap::new();
    content.insert("score".to_string(), serde_json::Value::Number(serde_json::Number::from_f64(0.5).unwrap()));

    let ev = create_event(
        EventType::new("trust.updated").unwrap(),
        ActorId::new("alice").unwrap(),
        content,
        vec![boot_id],
        ConversationId::new("conv_1").unwrap(),
        prev_hash,
        &signer,
        1,
    );
    store.append(ev).unwrap();
    assert_eq!(store.count(), 2);
    let v = store.verify_chain();
    assert!(v.valid);
    assert_eq!(v.length, 2);
}

#[test]
fn recent_returns_reverse_order() {
    let mut store = InMemoryStore::new();
    store.append(bootstrap()).unwrap();
    let recent = store.recent(10);
    assert_eq!(recent.len(), 1);
}

// ── Helper to build a chained event ────────────────────────────────────

fn append_event(
    store: &mut InMemoryStore,
    event_type: &str,
    source: &str,
    conversation_id: &str,
    causes: Vec<EventId>,
) -> Event {
    let head = store.head().unwrap();
    let prev_hash = head.hash.clone();
    let ev = create_event(
        EventType::new(event_type).unwrap(),
        ActorId::new(source).unwrap(),
        std::collections::BTreeMap::new(),
        causes,
        ConversationId::new(conversation_id).unwrap(),
        prev_hash,
        &NoopSigner,
        1,
    );
    store.append(ev).unwrap()
}

// ── by_type ────────────────────────────────────────────────────────────

#[test]
fn by_type_filters_correctly() {
    let mut store = InMemoryStore::new();
    let boot = bootstrap();
    let boot_id = boot.id.clone();
    store.append(boot).unwrap();

    let ev1 = append_event(&mut store, "trust.updated", "alice", "conv_1", vec![boot_id.clone()]);
    let ev1_id = ev1.id.clone();
    append_event(&mut store, "message.sent", "alice", "conv_1", vec![ev1_id.clone()]);
    append_event(&mut store, "trust.updated", "alice", "conv_1", vec![ev1_id]);

    let trust_type = EventType::new("trust.updated").unwrap();
    let results = store.by_type(&trust_type, 10);
    assert_eq!(results.len(), 2);
    // Newest first
    assert_eq!(results[0].event_type.value(), "trust.updated");
    assert_eq!(results[1].event_type.value(), "trust.updated");
}

#[test]
fn by_type_respects_limit() {
    let mut store = InMemoryStore::new();
    let boot = bootstrap();
    let boot_id = boot.id.clone();
    store.append(boot).unwrap();

    let ev1 = append_event(&mut store, "trust.updated", "alice", "conv_1", vec![boot_id]);
    append_event(&mut store, "trust.updated", "alice", "conv_1", vec![ev1.id.clone()]);

    let trust_type = EventType::new("trust.updated").unwrap();
    let results = store.by_type(&trust_type, 1);
    assert_eq!(results.len(), 1);
}

#[test]
fn by_type_returns_empty_for_no_match() {
    let mut store = InMemoryStore::new();
    store.append(bootstrap()).unwrap();

    let t = EventType::new("nonexistent.type").unwrap();
    let results = store.by_type(&t, 10);
    assert!(results.is_empty());
}

// ── by_source ──────────────────────────────────────────────────────────

#[test]
fn by_source_filters_correctly() {
    let mut store = InMemoryStore::new();
    let boot = bootstrap(); // source = alice
    let boot_id = boot.id.clone();
    store.append(boot).unwrap();

    let ev1 = append_event(&mut store, "trust.updated", "bob", "conv_1", vec![boot_id.clone()]);
    append_event(&mut store, "trust.updated", "alice", "conv_1", vec![ev1.id.clone()]);

    let bob = ActorId::new("bob").unwrap();
    let results = store.by_source(&bob, 10);
    assert_eq!(results.len(), 1);
    assert_eq!(results[0].source.value(), "bob");
}

#[test]
fn by_source_respects_limit() {
    let mut store = InMemoryStore::new();
    let boot = bootstrap();
    let boot_id = boot.id.clone();
    store.append(boot).unwrap();

    let ev1 = append_event(&mut store, "trust.updated", "alice", "conv_1", vec![boot_id]);
    append_event(&mut store, "trust.updated", "alice", "conv_1", vec![ev1.id.clone()]);

    let alice = ActorId::new("alice").unwrap();
    // bootstrap + 2 events = 3 from alice, but limit to 2
    let results = store.by_source(&alice, 2);
    assert_eq!(results.len(), 2);
}

// ── by_conversation ────────────────────────────────────────────────────

#[test]
fn by_conversation_filters_correctly() {
    let mut store = InMemoryStore::new();
    let boot = bootstrap(); // conversation = conv_alice
    let boot_id = boot.id.clone();
    store.append(boot).unwrap();

    let ev1 = append_event(&mut store, "trust.updated", "alice", "conv_a", vec![boot_id.clone()]);
    append_event(&mut store, "trust.updated", "alice", "conv_b", vec![ev1.id.clone()]);

    let conv_a = ConversationId::new("conv_a").unwrap();
    let results = store.by_conversation(&conv_a, 10);
    assert_eq!(results.len(), 1);
    assert_eq!(results[0].conversation_id.value(), "conv_a");
}

#[test]
fn by_conversation_respects_limit() {
    let mut store = InMemoryStore::new();
    let boot = bootstrap();
    let boot_id = boot.id.clone();
    store.append(boot).unwrap();

    let ev1 = append_event(&mut store, "trust.updated", "alice", "conv_a", vec![boot_id]);
    append_event(&mut store, "trust.updated", "alice", "conv_a", vec![ev1.id.clone()]);

    let conv_a = ConversationId::new("conv_a").unwrap();
    let results = store.by_conversation(&conv_a, 1);
    assert_eq!(results.len(), 1);
}

// ── ancestors ──────────────────────────────────────────────────────────

#[test]
fn ancestors_returns_causal_parents() {
    let mut store = InMemoryStore::new();
    let boot = bootstrap();
    let boot_id = boot.id.clone();
    store.append(boot).unwrap();

    let ev1 = append_event(&mut store, "trust.updated", "alice", "conv_1", vec![boot_id.clone()]);
    let ev2 = append_event(&mut store, "trust.updated", "alice", "conv_1", vec![ev1.id.clone()]);

    // ancestors of ev2 at depth 1 should be ev1
    let ancestors = store.ancestors(&ev2.id, 1).unwrap();
    assert_eq!(ancestors.len(), 1);
    assert_eq!(ancestors[0].id, ev1.id);

    // ancestors of ev2 at depth 2 should be ev1 + boot
    let ancestors = store.ancestors(&ev2.id, 2).unwrap();
    assert_eq!(ancestors.len(), 2);
}

#[test]
fn ancestors_error_on_not_found() {
    let store = InMemoryStore::new();
    let id = EventId::new("019462a0-0000-7000-8000-000000000001").unwrap();
    assert!(store.ancestors(&id, 5).is_err());
}

#[test]
fn ancestors_depth_zero_returns_empty() {
    let mut store = InMemoryStore::new();
    let boot = bootstrap();
    let boot_id = boot.id.clone();
    store.append(boot).unwrap();

    let ancestors = store.ancestors(&boot_id, 0).unwrap();
    assert!(ancestors.is_empty());
}

// ── descendants ────────────────────────────────────────────────────────

#[test]
fn descendants_returns_causal_children() {
    let mut store = InMemoryStore::new();
    let boot = bootstrap();
    let boot_id = boot.id.clone();
    store.append(boot).unwrap();

    let ev1 = append_event(&mut store, "trust.updated", "alice", "conv_1", vec![boot_id.clone()]);
    append_event(&mut store, "trust.updated", "alice", "conv_1", vec![ev1.id.clone()]);

    // descendants of boot at depth 1 should include ev1 (ev1 references boot in causes)
    let desc = store.descendants(&boot_id, 1).unwrap();
    assert_eq!(desc.len(), 1);
    assert_eq!(desc[0].id, ev1.id);

    // descendants of boot at depth 2 should include ev1 and ev2
    let desc = store.descendants(&boot_id, 2).unwrap();
    assert_eq!(desc.len(), 2);
}

#[test]
fn descendants_error_on_not_found() {
    let store = InMemoryStore::new();
    let id = EventId::new("019462a0-0000-7000-8000-000000000001").unwrap();
    assert!(store.descendants(&id, 5).is_err());
}

#[test]
fn descendants_depth_zero_returns_empty() {
    let mut store = InMemoryStore::new();
    let boot = bootstrap();
    let boot_id = boot.id.clone();
    store.append(boot).unwrap();

    let desc = store.descendants(&boot_id, 0).unwrap();
    assert!(desc.is_empty());
}
