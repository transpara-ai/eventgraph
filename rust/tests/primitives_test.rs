use std::collections::HashMap;

use eventgraph::primitive::Snapshot;
use eventgraph::primitives::{create_all_primitives, layer_counts};

#[test]
fn test_all_201_primitives_instantiate() {
    let prims = create_all_primitives();
    assert_eq!(prims.len(), 201);
    for p in &prims {
        // Each primitive should have a non-empty id
        assert!(!p.id().value().is_empty(), "Primitive has empty id");
    }
}

#[test]
fn test_create_all_primitives_returns_201() {
    let prims = create_all_primitives();
    assert_eq!(prims.len(), 201, "Expected exactly 201 primitives");
}

#[test]
fn test_layer_counts_correct() {
    let prims = create_all_primitives();
    let mut counts: HashMap<u8, usize> = HashMap::new();
    for p in &prims {
        *counts.entry(p.layer().value()).or_insert(0) += 1;
    }

    for (layer, expected) in layer_counts() {
        let actual = counts.get(&layer).copied().unwrap_or(0);
        assert_eq!(
            actual, expected,
            "Layer {} expected {} primitives, got {}",
            layer, expected, actual
        );
    }

    // Verify all 14 layers are present
    assert_eq!(counts.len(), 14, "Expected 14 layers");
}

#[test]
fn test_process_returns_mutations() {
    let prims = create_all_primitives();
    let snapshot = Snapshot {
        tick: 1,
        primitives: HashMap::new(),
        pending_events: vec![],
        recent_events: vec![],
    };

    for p in &prims {
        let mutations = p.process(1, &[], &snapshot);
        assert_eq!(
            mutations.len(),
            2,
            "Primitive {} should return 2 mutations (eventsProcessed + lastTick)",
            p.id().value()
        );
    }
}

#[test]
fn test_subscriptions_non_empty() {
    let prims = create_all_primitives();
    for p in &prims {
        let subs = p.subscriptions();
        assert!(
            !subs.is_empty(),
            "Primitive {} has empty subscriptions",
            p.id().value()
        );
    }
}

#[test]
fn test_unique_primitive_ids() {
    let prims = create_all_primitives();
    let mut ids: Vec<String> = prims.iter().map(|p| p.id().value().to_string()).collect();
    let total = ids.len();
    ids.sort();
    ids.dedup();
    assert_eq!(ids.len(), total, "All primitive IDs must be unique");
}

#[test]
fn test_cadence_is_valid() {
    let prims = create_all_primitives();
    for p in &prims {
        assert!(
            p.cadence().value() >= 1,
            "Primitive {} has invalid cadence",
            p.id().value()
        );
    }
}

#[test]
fn test_layer_values_in_range() {
    let prims = create_all_primitives();
    for p in &prims {
        assert!(
            p.layer().value() <= 13,
            "Primitive {} has layer > 13",
            p.id().value()
        );
    }
}
