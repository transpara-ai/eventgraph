use std::collections::HashSet;

use eventgraph::codegraph::*;
use eventgraph::types::*;

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
// Compositions tests
// ══════════════════════════════════════════════════════════════════════════

#[test]
fn test_all_compositions() {
    let comps = all_codegraph_compositions();
    assert_eq!(comps.len(), 7, "Expected exactly 7 compositions");
}

#[test]
fn test_composition_names_unique() {
    let comps = all_codegraph_compositions();
    let names: HashSet<&str> = comps.iter().map(|c| c.name).collect();
    assert_eq!(names.len(), comps.len(), "All composition names must be unique");
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
