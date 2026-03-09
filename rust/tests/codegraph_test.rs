use std::collections::{BTreeMap, HashMap};

use eventgraph::codegraph::*;
use eventgraph::event::{create_event, new_event_id, NoopSigner};
use eventgraph::primitive::{Mutation, Registry, Snapshot};
use eventgraph::types::*;

// ── Helper: create an event with a given type ───────────────────────────

fn make_event(event_type_str: &str) -> eventgraph::event::Event {
    let et = EventType::new(event_type_str).unwrap();
    let source = ActorId::new("test-agent").unwrap();
    let conv = ConversationId::new("test-conv").unwrap();
    let cause = new_event_id();
    create_event(
        et, source, BTreeMap::new(),
        vec![cause], conv,
        Hash::zero(), &NoopSigner, 1,
    )
}

fn empty_snapshot() -> Snapshot {
    Snapshot {
        tick: 1,
        primitives: HashMap::new(),
        pending_events: vec![],
        recent_events: vec![],
    }
}

// ══════════════════════════════════════════════════════════════════════════
// Event type constants tests
// ══════════════════════════════════════════════════════════════════════════

#[test]
fn test_all_codegraph_event_types() {
    let types = all_codegraph_event_types();
    assert_eq!(types.len(), 35, "Expected 35 codegraph event types");
}

#[test]
fn test_all_codegraph_event_types_valid() {
    for et_str in all_codegraph_event_types() {
        assert!(
            EventType::new(et_str).is_ok(),
            "Event type '{}' should be a valid EventType",
            et_str
        );
    }
}

#[test]
fn test_all_codegraph_event_types_start_with_codegraph() {
    for et_str in all_codegraph_event_types() {
        assert!(
            et_str.starts_with("codegraph."),
            "Codegraph event type '{}' should start with 'codegraph.'",
            et_str
        );
    }
}

// ══════════════════════════════════════════════════════════════════════════
// Primitives tests
// ══════════════════════════════════════════════════════════════════════════

#[test]
fn test_all_codegraph_primitives_count() {
    let prims = all_codegraph_primitives();
    assert_eq!(prims.len(), 61, "Expected exactly 61 codegraph primitives");
}

#[test]
fn test_no_duplicate_ids() {
    let prims = all_codegraph_primitives();
    let mut ids: Vec<String> = prims.iter().map(|p| p.id().value().to_string()).collect();
    let total = ids.len();
    ids.sort();
    ids.dedup();
    assert_eq!(ids.len(), total, "All codegraph primitive IDs must be unique");
}

#[test]
fn test_all_layer_5() {
    let prims = all_codegraph_primitives();
    for p in &prims {
        assert_eq!(p.layer().value(), 5, "Codegraph primitive {} should be layer 5", p.id().value());
    }
}

#[test]
fn test_all_cadence_1() {
    let prims = all_codegraph_primitives();
    for p in &prims {
        assert_eq!(p.cadence().value(), 1, "Codegraph primitive {} should have cadence 1", p.id().value());
    }
}

#[test]
fn test_all_have_subscriptions() {
    let prims = all_codegraph_primitives();
    for p in &prims {
        let subs = p.subscriptions();
        assert!(!subs.is_empty(), "Codegraph primitive {} has no subscriptions", p.id().value());
    }
}

#[test]
fn test_is_codegraph_primitive() {
    let prims = all_codegraph_primitives();
    for p in &prims {
        assert!(
            is_codegraph_primitive(&p.id()),
            "{} should be identified as codegraph primitive",
            p.id().value()
        );
    }

    // Negative case: an agent primitive should not match
    let agent_id = PrimitiveId::new("agent.Identity").unwrap();
    assert!(!is_codegraph_primitive(&agent_id), "agent.Identity should not be a codegraph primitive");
}

#[test]
fn test_register_all() {
    let mut registry = Registry::new();
    register_all_codegraph(&mut registry).expect("register_all_codegraph should succeed");
    assert_eq!(registry.count(), 61, "Registry should contain 61 codegraph primitives");

    // All should be active
    let prims = all_codegraph_primitives();
    for p in &prims {
        assert_eq!(
            registry.get_lifecycle(&p.id()),
            LifecycleState::Active,
            "Primitive {} should be active after register_all_codegraph",
            p.id().value()
        );
    }
}

#[test]
fn test_register_all_no_duplicates() {
    let mut registry = Registry::new();
    register_all_codegraph(&mut registry).expect("first register_all_codegraph should succeed");
    // Second registration should fail due to duplicate IDs
    assert!(register_all_codegraph(&mut registry).is_err(), "Duplicate registration should fail");
}

#[test]
fn test_process_returns_mutations() {
    let prims = all_codegraph_primitives();
    let snap = empty_snapshot();

    for p in &prims {
        // With no events
        let mutations = p.process(1, &[], &snap);
        let has_last_tick = mutations.iter().any(|m| matches!(m,
            Mutation::UpdateState { key, .. } if key == "lastTick"
        ));
        assert!(has_last_tick, "Primitive {} should produce lastTick mutation", p.id().value());

        let has_events_processed = mutations.iter().any(|m| matches!(m,
            Mutation::UpdateState { key, .. } if key == "eventsProcessed"
        ));
        assert!(has_events_processed, "Primitive {} should produce eventsProcessed mutation", p.id().value());

        // With some events
        let events = vec![
            make_event(CODEGRAPH_ENTITY_DEFINED),
            make_event(CODEGRAPH_UI_VIEW_RENDERED),
        ];
        let mutations = p.process(5, &events, &snap);

        // Should report eventsProcessed = 2
        let ep = mutations.iter().find(|m| matches!(m,
            Mutation::UpdateState { key, .. } if key == "eventsProcessed"
        )).unwrap();
        if let Mutation::UpdateState { value, .. } = ep {
            assert_eq!(value.as_u64().unwrap(), 2);
        }

        // Should report lastTick = 5
        let lt = mutations.iter().find(|m| matches!(m,
            Mutation::UpdateState { key, .. } if key == "lastTick"
        )).unwrap();
        if let Mutation::UpdateState { value, .. } = lt {
            assert_eq!(value.as_u64().unwrap(), 5);
        }
    }
}

// ══════════════════════════════════════════════════════════════════════════
// Compositions tests
// ══════════════════════════════════════════════════════════════════════════

#[test]
fn test_all_compositions() {
    let comps = all_codegraph_compositions();
    assert_eq!(comps.len(), 7, "Expected exactly 7 compositions");
}

#[test]
fn test_composition_names() {
    let comps = all_codegraph_compositions();
    let names: Vec<&str> = comps.iter().map(|c| c.name).collect();
    assert_eq!(names, vec![
        "Board", "Detail", "Feed", "Dashboard",
        "Inbox", "Wizard", "Skin",
    ]);
}

#[test]
fn test_composition_primitives_exist() {
    let all_prim_ids: Vec<String> = all_codegraph_primitives()
        .iter()
        .map(|p| p.id().value().to_string())
        .collect();

    for comp in all_codegraph_compositions() {
        for prim_id in &comp.primitives {
            assert!(
                all_prim_ids.contains(&prim_id.to_string()),
                "Composition {} references unknown primitive {}",
                comp.name, prim_id
            );
        }
    }
}

#[test]
fn test_board_composition() {
    let b = board();
    assert_eq!(b.name, "Board");
    assert_eq!(b.primitives.len(), 10);
    assert!(b.primitives.contains(&"CGLayout"));
    assert!(b.primitives.contains(&"CGList"));
    assert!(b.primitives.contains(&"CGDrag"));
    assert!(b.primitives.contains(&"CGState"));
}

#[test]
fn test_detail_composition() {
    let d = detail();
    assert_eq!(d.name, "Detail");
    assert_eq!(d.primitives.len(), 9);
    assert!(d.primitives.contains(&"CGForm"));
    assert!(d.primitives.contains(&"CGThread"));
    assert!(d.primitives.contains(&"CGAudit"));
    assert!(d.primitives.contains(&"CGHistory"));
}

#[test]
fn test_feed_composition() {
    let f = feed();
    assert_eq!(f.name, "Feed");
    assert_eq!(f.primitives.len(), 7);
    assert!(f.primitives.contains(&"CGPagination"));
    assert!(f.primitives.contains(&"CGRecency"));
}

#[test]
fn test_wizard_composition() {
    let w = wizard();
    assert_eq!(w.name, "Wizard");
    assert_eq!(w.primitives.len(), 8);
    assert!(w.primitives.contains(&"CGSequence"));
    assert!(w.primitives.contains(&"CGConsequencePreview"));
    assert!(w.primitives.contains(&"CGConstraint"));
}

#[test]
fn test_skin_composition() {
    let s = skin();
    assert_eq!(s.name, "Skin");
    assert_eq!(s.primitives.len(), 7);
    assert!(s.primitives.contains(&"CGPalette"));
    assert!(s.primitives.contains(&"CGTypography"));
    assert!(s.primitives.contains(&"CGShape"));
}
