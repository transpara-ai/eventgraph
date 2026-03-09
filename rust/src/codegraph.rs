// Code Graph vocabulary: 35 event types across 11 categories, 7 named compositions.
// Code Graph atoms are vocabulary/content, not tick engine participants.

// ── Code Graph Event Type Constants ─────────────────────────────────────

// Data events
pub const CODEGRAPH_ENTITY_DEFINED: &str = "codegraph.entity.defined";
pub const CODEGRAPH_ENTITY_MODIFIED: &str = "codegraph.entity.modified";
pub const CODEGRAPH_ENTITY_RELATED: &str = "codegraph.entity.related";
pub const CODEGRAPH_STATE_TRANSITIONED: &str = "codegraph.state.transitioned";

// Logic events
pub const CODEGRAPH_LOGIC_TRANSFORM_APPLIED: &str = "codegraph.logic.transform.applied";
pub const CODEGRAPH_LOGIC_CONDITION_EVALUATED: &str = "codegraph.logic.condition.evaluated";
pub const CODEGRAPH_LOGIC_SEQUENCE_EXECUTED: &str = "codegraph.logic.sequence.executed";
pub const CODEGRAPH_LOGIC_LOOP_ITERATED: &str = "codegraph.logic.loop.iterated";
pub const CODEGRAPH_LOGIC_TRIGGER_FIRED: &str = "codegraph.logic.trigger.fired";
pub const CODEGRAPH_LOGIC_CONSTRAINT_VIOLATED: &str = "codegraph.logic.constraint.violated";

// IO events
pub const CODEGRAPH_IO_QUERY_EXECUTED: &str = "codegraph.io.query.executed";
pub const CODEGRAPH_IO_COMMAND_EXECUTED: &str = "codegraph.io.command.executed";
pub const CODEGRAPH_IO_COMMAND_REJECTED: &str = "codegraph.io.command.rejected";
pub const CODEGRAPH_IO_SUBSCRIBE_REGISTERED: &str = "codegraph.io.subscribe.registered";
pub const CODEGRAPH_IO_AUTHORIZE_GRANTED: &str = "codegraph.io.authorize.granted";
pub const CODEGRAPH_IO_AUTHORIZE_DENIED: &str = "codegraph.io.authorize.denied";
pub const CODEGRAPH_IO_SEARCH_EXECUTED: &str = "codegraph.io.search.executed";
pub const CODEGRAPH_IO_INTEROP_SENT: &str = "codegraph.io.interop.sent";
pub const CODEGRAPH_IO_INTEROP_RECEIVED: &str = "codegraph.io.interop.received";

// UI events
pub const CODEGRAPH_UI_VIEW_RENDERED: &str = "codegraph.ui.view.rendered";
pub const CODEGRAPH_UI_ACTION_INVOKED: &str = "codegraph.ui.action.invoked";
pub const CODEGRAPH_UI_NAVIGATION_TRIGGERED: &str = "codegraph.ui.navigation.triggered";
pub const CODEGRAPH_UI_FEEDBACK_EMITTED: &str = "codegraph.ui.feedback.emitted";
pub const CODEGRAPH_UI_ALERT_DISPATCHED: &str = "codegraph.ui.alert.dispatched";
pub const CODEGRAPH_UI_DRAG_COMPLETED: &str = "codegraph.ui.drag.completed";
pub const CODEGRAPH_UI_SELECTION_CHANGED: &str = "codegraph.ui.selection.changed";
pub const CODEGRAPH_UI_CONFIRMATION_RESOLVED: &str = "codegraph.ui.confirmation.resolved";

// Aesthetic events
pub const CODEGRAPH_AESTHETIC_SKIN_APPLIED: &str = "codegraph.aesthetic.skin.applied";

// Temporal events
pub const CODEGRAPH_TEMPORAL_UNDO_EXECUTED: &str = "codegraph.temporal.undo.executed";
pub const CODEGRAPH_TEMPORAL_RETRY_ATTEMPTED: &str = "codegraph.temporal.retry.attempted";

// Resilience events
pub const CODEGRAPH_RESILIENCE_OFFLINE_ENTERED: &str = "codegraph.resilience.offline.entered";
pub const CODEGRAPH_RESILIENCE_OFFLINE_SYNCED: &str = "codegraph.resilience.offline.synced";

// Structural events
pub const CODEGRAPH_STRUCTURAL_SCOPE_DEFINED: &str = "codegraph.structural.scope.defined";

// Social events
pub const CODEGRAPH_SOCIAL_PRESENCE_UPDATED: &str = "codegraph.social.presence.updated";
pub const CODEGRAPH_SOCIAL_SALIENCE_TRIGGERED: &str = "codegraph.social.salience.triggered";

/// Returns all 35 code graph event type strings.
pub fn all_codegraph_event_types() -> Vec<&'static str> {
    vec![
        // Data
        CODEGRAPH_ENTITY_DEFINED, CODEGRAPH_ENTITY_MODIFIED,
        CODEGRAPH_ENTITY_RELATED, CODEGRAPH_STATE_TRANSITIONED,
        // Logic
        CODEGRAPH_LOGIC_TRANSFORM_APPLIED, CODEGRAPH_LOGIC_CONDITION_EVALUATED,
        CODEGRAPH_LOGIC_SEQUENCE_EXECUTED, CODEGRAPH_LOGIC_LOOP_ITERATED,
        CODEGRAPH_LOGIC_TRIGGER_FIRED, CODEGRAPH_LOGIC_CONSTRAINT_VIOLATED,
        // IO
        CODEGRAPH_IO_QUERY_EXECUTED, CODEGRAPH_IO_COMMAND_EXECUTED,
        CODEGRAPH_IO_COMMAND_REJECTED, CODEGRAPH_IO_SUBSCRIBE_REGISTERED,
        CODEGRAPH_IO_AUTHORIZE_GRANTED, CODEGRAPH_IO_AUTHORIZE_DENIED,
        CODEGRAPH_IO_SEARCH_EXECUTED, CODEGRAPH_IO_INTEROP_SENT,
        CODEGRAPH_IO_INTEROP_RECEIVED,
        // UI
        CODEGRAPH_UI_VIEW_RENDERED, CODEGRAPH_UI_ACTION_INVOKED,
        CODEGRAPH_UI_NAVIGATION_TRIGGERED, CODEGRAPH_UI_FEEDBACK_EMITTED,
        CODEGRAPH_UI_ALERT_DISPATCHED, CODEGRAPH_UI_DRAG_COMPLETED,
        CODEGRAPH_UI_SELECTION_CHANGED, CODEGRAPH_UI_CONFIRMATION_RESOLVED,
        // Aesthetic
        CODEGRAPH_AESTHETIC_SKIN_APPLIED,
        // Temporal
        CODEGRAPH_TEMPORAL_UNDO_EXECUTED, CODEGRAPH_TEMPORAL_RETRY_ATTEMPTED,
        // Resilience
        CODEGRAPH_RESILIENCE_OFFLINE_ENTERED, CODEGRAPH_RESILIENCE_OFFLINE_SYNCED,
        // Structural
        CODEGRAPH_STRUCTURAL_SCOPE_DEFINED,
        // Social
        CODEGRAPH_SOCIAL_PRESENCE_UPDATED, CODEGRAPH_SOCIAL_SALIENCE_TRIGGERED,
    ]
}

// ══════════════════════════════════════════════════════════════════════════
// COMPOSITIONS (7) — Named groupings of code graph atoms
// ══════════════════════════════════════════════════════════════════════════

/// Represents a named grouping of code graph atoms.
#[derive(Debug, Clone)]
pub struct Composition {
    pub name: &'static str,
    pub primitives: Vec<&'static str>,
    pub events: Vec<&'static str>,
}

/// Board — Kanban/task board view.
pub fn board() -> Composition {
    Composition {
        name: "Board",
        primitives: vec![
            "CGLayout", "CGList", "CGQuery", "CGDisplay", "CGDrag",
            "CGCommand", "CGEmpty", "CGAction", "CGLoop", "CGState",
        ],
        events: vec![],
    }
}

/// Detail — Single-entity detail view.
pub fn detail() -> Composition {
    Composition {
        name: "Detail",
        primitives: vec![
            "CGLayout", "CGDisplay", "CGForm", "CGThread", "CGAudit",
            "CGHistory", "CGList", "CGAction", "CGNavigation",
        ],
        events: vec![],
    }
}

/// Feed — Scrollable feed of items.
pub fn feed() -> Composition {
    Composition {
        name: "Feed",
        primitives: vec![
            "CGList", "CGQuery", "CGDisplay", "CGAvatar",
            "CGSubscribe", "CGPagination", "CGRecency",
        ],
        events: vec![],
    }
}

/// Dashboard — Summary/metrics view.
pub fn dashboard() -> Composition {
    Composition {
        name: "Dashboard",
        primitives: vec![
            "CGLayout", "CGDisplay", "CGQuery", "CGTransform", "CGSalience",
        ],
        events: vec![],
    }
}

/// Inbox — Actionable list of items requiring attention.
pub fn inbox() -> Composition {
    Composition {
        name: "Inbox",
        primitives: vec![
            "CGList", "CGQuery", "CGDisplay", "CGAvatar",
            "CGSalience", "CGAction", "CGSelection", "CGEmpty",
        ],
        events: vec![],
    }
}

/// Wizard — Multi-step guided flow.
pub fn wizard() -> Composition {
    Composition {
        name: "Wizard",
        primitives: vec![
            "CGSequence", "CGForm", "CGInput", "CGDisplay",
            "CGAction", "CGNavigation", "CGConsequencePreview", "CGConstraint",
        ],
        events: vec![],
    }
}

/// Skin — Visual theming composition.
pub fn skin() -> Composition {
    Composition {
        name: "Skin",
        primitives: vec![
            "CGPalette", "CGTypography", "CGSpacing",
            "CGElevation", "CGMotion", "CGDensity", "CGShape",
        ],
        events: vec![],
    }
}

/// Returns all 7 named compositions.
pub fn all_codegraph_compositions() -> Vec<Composition> {
    vec![
        board(), detail(), feed(), dashboard(),
        inbox(), wizard(), skin(),
    ]
}
