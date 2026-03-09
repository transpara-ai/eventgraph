// Code Graph primitives: 61 primitives across 11 categories, 7 named compositions.
// All code graph primitives operate at Layer 5 (Code Graph).

use serde_json::Value;

use crate::event::Event;
use crate::primitive::{Mutation, Primitive, Registry, Snapshot};
use crate::types::*;

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

// ── Helpers ─────────────────────────────────────────────────────────────

fn update_state(id: &str, key: &str, value: Value) -> Mutation {
    Mutation::UpdateState {
        primitive_id: PrimitiveId::new(id).unwrap(),
        key: key.to_string(),
        value,
    }
}

fn int_val(n: usize) -> Value {
    Value::Number(serde_json::Number::from(n))
}

fn tick_val(tick: u64) -> Value {
    Value::Number(serde_json::Number::from(tick))
}

// ── Generic process function ────────────────────────────────────────────
// All CG primitives use the same process logic: count events, emit
// eventsProcessed and lastTick mutations.

fn process_cg(id: &str, tick: u64, events: &[Event]) -> Vec<Mutation> {
    vec![
        update_state(id, "eventsProcessed", int_val(events.len())),
        update_state(id, "lastTick", tick_val(tick)),
    ]
}

// ── CodeGraph Primitive definition and wrapper ──────────────────────────

struct CgPrimitiveDef {
    name: &'static str,
    subs: &'static [&'static str],
}

struct CgPrimitive {
    def: &'static CgPrimitiveDef,
}

impl Primitive for CgPrimitive {
    fn id(&self) -> PrimitiveId {
        PrimitiveId::new(self.def.name).unwrap()
    }

    fn layer(&self) -> Layer {
        Layer::new(5).unwrap()
    }

    fn subscriptions(&self) -> Vec<SubscriptionPattern> {
        self.def.subs.iter()
            .map(|s| SubscriptionPattern::new(*s).unwrap())
            .collect()
    }

    fn cadence(&self) -> Cadence {
        Cadence::new(1).unwrap()
    }

    fn process(&self, tick: u64, events: &[Event], _snapshot: &Snapshot) -> Vec<Mutation> {
        process_cg(self.def.name, tick, events)
    }
}

// ── Static definitions for all 61 primitives ────────────────────────────

static PRIMITIVE_DEFS: &[CgPrimitiveDef] = &[
    // Data (6)
    CgPrimitiveDef { name: "CGEntity",     subs: &["codegraph.entity.*", "codegraph.io.command.*"] },
    CgPrimitiveDef { name: "CGProperty",   subs: &["codegraph.entity.*"] },
    CgPrimitiveDef { name: "CGRelation",   subs: &["codegraph.entity.*"] },
    CgPrimitiveDef { name: "CGCollection", subs: &["codegraph.entity.*", "codegraph.io.query.*"] },
    CgPrimitiveDef { name: "CGState",      subs: &["codegraph.state.*", "codegraph.io.command.*"] },
    CgPrimitiveDef { name: "CGEvent",      subs: &["codegraph.*"] },
    // Logic (6)
    CgPrimitiveDef { name: "CGTransform",  subs: &["codegraph.logic.transform.*"] },
    CgPrimitiveDef { name: "CGCondition",  subs: &["codegraph.logic.condition.*"] },
    CgPrimitiveDef { name: "CGSequence",   subs: &["codegraph.logic.sequence.*"] },
    CgPrimitiveDef { name: "CGLoop",       subs: &["codegraph.logic.loop.*"] },
    CgPrimitiveDef { name: "CGTrigger",    subs: &["codegraph.logic.trigger.*", "codegraph.entity.*"] },
    CgPrimitiveDef { name: "CGConstraint", subs: &["codegraph.logic.constraint.*", "codegraph.io.command.*"] },
    // IO (6)
    CgPrimitiveDef { name: "CGQuery",      subs: &["codegraph.io.query.*"] },
    CgPrimitiveDef { name: "CGCommand",    subs: &["codegraph.io.command.*"] },
    CgPrimitiveDef { name: "CGSubscribe",  subs: &["codegraph.io.subscribe.*"] },
    CgPrimitiveDef { name: "CGAuthorize",  subs: &["codegraph.io.authorize.*", "authority.*"] },
    CgPrimitiveDef { name: "CGSearch",     subs: &["codegraph.io.search.*"] },
    CgPrimitiveDef { name: "CGInterop",    subs: &["codegraph.io.interop.*"] },
    // UI (19)
    CgPrimitiveDef { name: "CGDisplay",       subs: &["codegraph.ui.*"] },
    CgPrimitiveDef { name: "CGInput",         subs: &["codegraph.ui.*"] },
    CgPrimitiveDef { name: "CGLayout",        subs: &["codegraph.ui.*"] },
    CgPrimitiveDef { name: "CGList",          subs: &["codegraph.ui.*", "codegraph.io.query.*"] },
    CgPrimitiveDef { name: "CGForm",          subs: &["codegraph.ui.*", "codegraph.io.command.*"] },
    CgPrimitiveDef { name: "CGAction",        subs: &["codegraph.ui.action.*"] },
    CgPrimitiveDef { name: "CGNavigation",    subs: &["codegraph.ui.navigation.*"] },
    CgPrimitiveDef { name: "CGView",          subs: &["codegraph.ui.view.*"] },
    CgPrimitiveDef { name: "CGFeedback",      subs: &["codegraph.ui.feedback.*"] },
    CgPrimitiveDef { name: "CGAlert",         subs: &["codegraph.ui.alert.*"] },
    CgPrimitiveDef { name: "CGThread",        subs: &["codegraph.ui.*", "codegraph.entity.*"] },
    CgPrimitiveDef { name: "CGAvatar",        subs: &["codegraph.ui.*"] },
    CgPrimitiveDef { name: "CGAudit",         subs: &["codegraph.*"] },
    CgPrimitiveDef { name: "CGDrag",          subs: &["codegraph.ui.drag.*"] },
    CgPrimitiveDef { name: "CGSelection",     subs: &["codegraph.ui.selection.*"] },
    CgPrimitiveDef { name: "CGConfirmation",  subs: &["codegraph.ui.confirmation.*"] },
    CgPrimitiveDef { name: "CGEmpty",         subs: &["codegraph.ui.*"] },
    CgPrimitiveDef { name: "CGLoading",       subs: &["codegraph.ui.*"] },
    CgPrimitiveDef { name: "CGPagination",    subs: &["codegraph.ui.*", "codegraph.io.query.*"] },
    // Aesthetic (7)
    CgPrimitiveDef { name: "CGPalette",    subs: &["codegraph.aesthetic.*"] },
    CgPrimitiveDef { name: "CGTypography", subs: &["codegraph.aesthetic.*"] },
    CgPrimitiveDef { name: "CGSpacing",    subs: &["codegraph.aesthetic.*"] },
    CgPrimitiveDef { name: "CGElevation",  subs: &["codegraph.aesthetic.*"] },
    CgPrimitiveDef { name: "CGMotion",     subs: &["codegraph.aesthetic.*"] },
    CgPrimitiveDef { name: "CGDensity",    subs: &["codegraph.aesthetic.*"] },
    CgPrimitiveDef { name: "CGShape",      subs: &["codegraph.aesthetic.*"] },
    // Accessibility (4)
    CgPrimitiveDef { name: "CGAnnounce",   subs: &["codegraph.ui.*", "codegraph.aesthetic.*"] },
    CgPrimitiveDef { name: "CGFocus",      subs: &["codegraph.ui.*"] },
    CgPrimitiveDef { name: "CGContrast",   subs: &["codegraph.aesthetic.*"] },
    CgPrimitiveDef { name: "CGSimplify",   subs: &["codegraph.ui.*", "codegraph.aesthetic.*"] },
    // Temporal (3)
    CgPrimitiveDef { name: "CGRecency",    subs: &["codegraph.entity.*", "codegraph.temporal.*"] },
    CgPrimitiveDef { name: "CGHistory",    subs: &["codegraph.entity.*", "codegraph.temporal.*"] },
    CgPrimitiveDef { name: "CGLiveness",   subs: &["codegraph.io.subscribe.*", "codegraph.social.presence.*"] },
    // Resilience (4)
    CgPrimitiveDef { name: "CGUndo",       subs: &["codegraph.temporal.undo.*", "codegraph.io.command.*"] },
    CgPrimitiveDef { name: "CGRetry",      subs: &["codegraph.temporal.retry.*", "codegraph.io.command.*"] },
    CgPrimitiveDef { name: "CGFallback",   subs: &["codegraph.resilience.*", "codegraph.ui.*"] },
    CgPrimitiveDef { name: "CGOffline",    subs: &["codegraph.resilience.offline.*"] },
    // Structural (3)
    CgPrimitiveDef { name: "CGScope",      subs: &["codegraph.structural.*"] },
    CgPrimitiveDef { name: "CGFormat",     subs: &["codegraph.io.*", "codegraph.structural.*"] },
    CgPrimitiveDef { name: "CGGesture",    subs: &["codegraph.ui.*"] },
    // Social (3)
    CgPrimitiveDef { name: "CGPresence",           subs: &["codegraph.social.presence.*"] },
    CgPrimitiveDef { name: "CGSalience",           subs: &["codegraph.social.salience.*", "codegraph.entity.*"] },
    CgPrimitiveDef { name: "CGConsequencePreview", subs: &["codegraph.ui.confirmation.*", "codegraph.io.command.*"] },
];

// ── Public API ──────────────────────────────────────────────────────────

/// Returns all 61 code graph primitives as boxed trait objects.
pub fn all_codegraph_primitives() -> Vec<Box<dyn Primitive>> {
    PRIMITIVE_DEFS.iter()
        .map(|def| Box::new(CgPrimitive { def }) as Box<dyn Primitive>)
        .collect()
}

/// Registers all 61 code graph primitives with the given registry and activates them.
pub fn register_all_codegraph(registry: &mut Registry) -> crate::errors::Result<()> {
    for p in all_codegraph_primitives() {
        let id = p.id();
        registry.register(p)?;
        registry.activate(&id)?;
    }
    Ok(())
}

/// Returns true if the primitive ID belongs to the code graph layer.
pub fn is_codegraph_primitive(id: &PrimitiveId) -> bool {
    id.value().starts_with("CG")
}

// ══════════════════════════════════════════════════════════════════════════
// COMPOSITIONS (7) — Named groupings of code graph primitives
// ══════════════════════════════════════════════════════════════════════════

/// Represents a named grouping of code graph primitive operations.
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
