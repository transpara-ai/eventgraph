use eventgraph::actor::*;
use eventgraph::errors::EventGraphError;
use eventgraph::types::PublicKey;

fn test_public_key(b: u8) -> PublicKey {
    let mut key = [0u8; 32];
    key[0] = b;
    PublicKey::new(key)
}

// ── Registration ──────────────────────────────────────────────────────

#[test]
fn test_register() {
    let mut store = InMemoryActorStore::new();
    let pk = test_public_key(1);

    let actor = store.register(pk, "Alice", ActorType::Human).unwrap();
    assert_eq!(actor.display_name(), "Alice");
    assert_eq!(actor.actor_type(), ActorType::Human);
    assert_eq!(actor.status(), ActorStatus::Active);
}

#[test]
fn test_register_idempotent() {
    let mut store = InMemoryActorStore::new();
    let pk = test_public_key(1);

    let a1 = store.register(pk.clone(), "Alice", ActorType::Human).unwrap();
    let a2 = store.register(pk, "Alice Again", ActorType::Human).unwrap();

    assert_eq!(a1.id(), a2.id());
    assert_eq!(store.actor_count(), 1);
}

// ── Get ───────────────────────────────────────────────────────────────

#[test]
fn test_get() {
    let mut store = InMemoryActorStore::new();
    let pk = test_public_key(1);
    let actor = store.register(pk, "Alice", ActorType::Human).unwrap();

    let got = store.get(actor.id()).unwrap();
    assert_eq!(got.id(), actor.id());
}

#[test]
fn test_get_not_found() {
    let store = InMemoryActorStore::new();
    let id = eventgraph::types::ActorId::new("actor_nonexistent").unwrap();
    let err = store.get(&id).unwrap_err();
    match err {
        EventGraphError::ActorNotFound { .. } => {}
        other => panic!("expected ActorNotFound, got {:?}", other),
    }
}

// ── GetByPublicKey ────────────────────────────────────────────────────

#[test]
fn test_get_by_public_key() {
    let mut store = InMemoryActorStore::new();
    let pk = test_public_key(1);
    let actor = store.register(pk.clone(), "Alice", ActorType::Human).unwrap();

    let got = store.get_by_public_key(&pk).unwrap();
    assert_eq!(got.id(), actor.id());
}

#[test]
fn test_get_by_public_key_not_found() {
    let store = InMemoryActorStore::new();
    let pk = test_public_key(99);
    let err = store.get_by_public_key(&pk).unwrap_err();
    match err {
        EventGraphError::ActorKeyNotFound { .. } => {}
        other => panic!("expected ActorKeyNotFound, got {:?}", other),
    }
}

// ── Update ────────────────────────────────────────────────────────────

#[test]
fn test_update() {
    let mut store = InMemoryActorStore::new();
    let pk = test_public_key(1);
    let actor = store.register(pk, "Alice", ActorType::Human).unwrap();

    let updated = store.update(actor.id(), &ActorUpdate {
        display_name: Some("Alice Updated".to_string()),
        metadata: None,
    }).unwrap();
    assert_eq!(updated.display_name(), "Alice Updated");
}

#[test]
fn test_update_metadata_merge() {
    let mut store = InMemoryActorStore::new();
    let pk = test_public_key(1);
    let actor = store.register(pk, "Alice", ActorType::Human).unwrap();

    // Set initial metadata
    let mut md1 = std::collections::BTreeMap::new();
    md1.insert("role".to_string(), serde_json::Value::String("builder".to_string()));
    store.update(actor.id(), &ActorUpdate {
        display_name: None,
        metadata: Some(md1),
    }).unwrap();

    // Merge additional metadata
    let mut md2 = std::collections::BTreeMap::new();
    md2.insert("team".to_string(), serde_json::Value::String("core".to_string()));
    let updated = store.update(actor.id(), &ActorUpdate {
        display_name: None,
        metadata: Some(md2),
    }).unwrap();

    let md = updated.metadata();
    assert_eq!(md.get("role").unwrap(), &serde_json::Value::String("builder".to_string()));
    assert_eq!(md.get("team").unwrap(), &serde_json::Value::String("core".to_string()));
}

#[test]
fn test_update_not_found() {
    let mut store = InMemoryActorStore::new();
    let id = eventgraph::types::ActorId::new("actor_nonexistent").unwrap();
    let err = store.update(&id, &ActorUpdate {
        display_name: None,
        metadata: None,
    }).unwrap_err();
    match err {
        EventGraphError::ActorNotFound { .. } => {}
        other => panic!("expected ActorNotFound, got {:?}", other),
    }
}

// ── Suspend ───────────────────────────────────────────────────────────

#[test]
fn test_suspend() {
    let mut store = InMemoryActorStore::new();
    let pk = test_public_key(1);
    let actor = store.register(pk, "Alice", ActorType::Human).unwrap();

    let suspended = store.suspend(actor.id(), "reason-event-id").unwrap();
    assert_eq!(suspended.status(), ActorStatus::Suspended);
}

#[test]
fn test_suspend_not_found() {
    let mut store = InMemoryActorStore::new();
    let id = eventgraph::types::ActorId::new("actor_nonexistent").unwrap();
    let err = store.suspend(&id, "reason").unwrap_err();
    match err {
        EventGraphError::ActorNotFound { .. } => {}
        other => panic!("expected ActorNotFound, got {:?}", other),
    }
}

#[test]
fn test_suspend_and_reactivate() {
    let mut store = InMemoryActorStore::new();
    let pk = test_public_key(1);
    let actor = store.register(pk, "Alice", ActorType::Human).unwrap();

    store.suspend(actor.id(), "reason").unwrap();
    let got = store.get(actor.id()).unwrap();
    assert_eq!(got.status(), ActorStatus::Suspended);

    let reactivated = store.reactivate(actor.id(), "reason").unwrap();
    assert_eq!(reactivated.status(), ActorStatus::Active);
}

// ── Memorial ──────────────────────────────────────────────────────────

#[test]
fn test_memorial() {
    let mut store = InMemoryActorStore::new();
    let pk = test_public_key(1);
    let actor = store.register(pk, "Alice", ActorType::Human).unwrap();

    let memorial = store.memorial(actor.id(), "reason").unwrap();
    assert_eq!(memorial.status(), ActorStatus::Memorial);
}

#[test]
fn test_memorial_is_terminal() {
    let mut store = InMemoryActorStore::new();
    let pk = test_public_key(1);
    let actor = store.register(pk, "Alice", ActorType::Human).unwrap();

    store.memorial(actor.id(), "reason").unwrap();

    // Try to suspend — should fail
    let err = store.suspend(actor.id(), "reason").unwrap_err();
    match err {
        EventGraphError::InvalidTransition { .. } => {}
        other => panic!("expected InvalidTransition, got {:?}", other),
    }
}

#[test]
fn test_memorial_reactivate_is_error() {
    let mut store = InMemoryActorStore::new();
    let pk = test_public_key(1);
    let actor = store.register(pk, "Alice", ActorType::Human).unwrap();

    store.memorial(actor.id(), "reason").unwrap();

    let err = store.reactivate(actor.id(), "reason").unwrap_err();
    match err {
        EventGraphError::InvalidTransition { .. } => {}
        other => panic!("expected InvalidTransition, got {:?}", other),
    }
}

// ── Reactivate ────────────────────────────────────────────────────────

#[test]
fn test_reactivate_from_active_is_error() {
    let mut store = InMemoryActorStore::new();
    let pk = test_public_key(1);
    let actor = store.register(pk, "Alice", ActorType::Human).unwrap();

    let err = store.reactivate(actor.id(), "reason").unwrap_err();
    match err {
        EventGraphError::InvalidTransition { .. } => {}
        other => panic!("expected InvalidTransition, got {:?}", other),
    }
}

#[test]
fn test_reactivate_not_found() {
    let mut store = InMemoryActorStore::new();
    let id = eventgraph::types::ActorId::new("actor_nonexistent").unwrap();
    let err = store.reactivate(&id, "reason").unwrap_err();
    match err {
        EventGraphError::ActorNotFound { .. } => {}
        other => panic!("expected ActorNotFound, got {:?}", other),
    }
}

// ── List ──────────────────────────────────────────────────────────────

#[test]
fn test_list() {
    let mut store = InMemoryActorStore::new();
    for i in 1u8..=5 {
        store.register(test_public_key(i), "Actor", ActorType::Human).unwrap();
    }

    let page = store.list(&ActorFilter {
        status: None,
        actor_type: None,
        limit: 10,
        after: None,
    }).unwrap();
    assert_eq!(page.items.len(), 5);
}

#[test]
fn test_list_with_status_filter() {
    let mut store = InMemoryActorStore::new();
    store.register(test_public_key(1), "Active1", ActorType::Human).unwrap();
    store.register(test_public_key(2), "Active2", ActorType::Human).unwrap();
    let a3 = store.register(test_public_key(3), "ToBeSuspended", ActorType::Human).unwrap();

    store.suspend(a3.id(), "reason").unwrap();

    let page = store.list(&ActorFilter {
        status: Some(ActorStatus::Active),
        actor_type: None,
        limit: 10,
        after: None,
    }).unwrap();
    assert_eq!(page.items.len(), 2);
}

#[test]
fn test_list_with_type_filter() {
    let mut store = InMemoryActorStore::new();
    store.register(test_public_key(1), "Human", ActorType::Human).unwrap();
    store.register(test_public_key(2), "AI", ActorType::AI).unwrap();
    store.register(test_public_key(3), "System", ActorType::System).unwrap();

    let page = store.list(&ActorFilter {
        status: None,
        actor_type: Some(ActorType::AI),
        limit: 10,
        after: None,
    }).unwrap();
    assert_eq!(page.items.len(), 1);
}

#[test]
fn test_list_pagination() {
    let mut store = InMemoryActorStore::new();
    for i in 1u8..=5 {
        store.register(test_public_key(i), "Actor", ActorType::Human).unwrap();
    }

    // Page 1
    let page1 = store.list(&ActorFilter {
        status: None,
        actor_type: None,
        limit: 2,
        after: None,
    }).unwrap();
    assert_eq!(page1.items.len(), 2);
    assert!(page1.has_more);

    // Page 2
    let page2 = store.list(&ActorFilter {
        status: None,
        actor_type: None,
        limit: 2,
        after: page1.cursor,
    }).unwrap();
    assert_eq!(page2.items.len(), 2);
}

// ── Actor getters ─────────────────────────────────────────────────────

#[test]
fn test_actor_getters() {
    let mut store = InMemoryActorStore::new();
    let pk = test_public_key(42);
    let actor = store.register(pk, "Alice", ActorType::AI).unwrap();

    assert!(!actor.id().value().is_empty());
    let _ = actor.public_key();
    assert!(actor.created_at_nanos() > 0);
    let _ = actor.metadata();
}
