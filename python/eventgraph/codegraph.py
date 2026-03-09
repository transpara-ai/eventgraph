"""Code Graph vocabulary — 35 event types and 7 compositions.

Code Graph atoms are vocabulary/content definitions, not tick engine participants.
They define how application UIs, data models, logic, and interactions are
represented as event types on the event graph.
"""

from __future__ import annotations

from dataclasses import dataclass

from .types import EventType


# ============================================================================
# Code Graph Event Type Constants (35)
# ============================================================================

# Data events
CODEGRAPH_ENTITY_DEFINED = EventType("codegraph.entity.defined")
CODEGRAPH_ENTITY_MODIFIED = EventType("codegraph.entity.modified")
CODEGRAPH_ENTITY_RELATED = EventType("codegraph.entity.related")
CODEGRAPH_STATE_TRANSITIONED = EventType("codegraph.state.transitioned")

# Logic events
CODEGRAPH_LOGIC_TRANSFORM_APPLIED = EventType("codegraph.logic.transform.applied")
CODEGRAPH_LOGIC_CONDITION_EVALUATED = EventType("codegraph.logic.condition.evaluated")
CODEGRAPH_LOGIC_SEQUENCE_EXECUTED = EventType("codegraph.logic.sequence.executed")
CODEGRAPH_LOGIC_LOOP_ITERATED = EventType("codegraph.logic.loop.iterated")
CODEGRAPH_LOGIC_TRIGGER_FIRED = EventType("codegraph.logic.trigger.fired")
CODEGRAPH_LOGIC_CONSTRAINT_VIOLATED = EventType("codegraph.logic.constraint.violated")

# IO events
CODEGRAPH_IO_QUERY_EXECUTED = EventType("codegraph.io.query.executed")
CODEGRAPH_IO_COMMAND_EXECUTED = EventType("codegraph.io.command.executed")
CODEGRAPH_IO_COMMAND_REJECTED = EventType("codegraph.io.command.rejected")
CODEGRAPH_IO_SUBSCRIBE_REGISTERED = EventType("codegraph.io.subscribe.registered")
CODEGRAPH_IO_AUTHORIZE_GRANTED = EventType("codegraph.io.authorize.granted")
CODEGRAPH_IO_AUTHORIZE_DENIED = EventType("codegraph.io.authorize.denied")
CODEGRAPH_IO_SEARCH_EXECUTED = EventType("codegraph.io.search.executed")
CODEGRAPH_IO_INTEROP_SENT = EventType("codegraph.io.interop.sent")
CODEGRAPH_IO_INTEROP_RECEIVED = EventType("codegraph.io.interop.received")

# UI events
CODEGRAPH_UI_VIEW_RENDERED = EventType("codegraph.ui.view.rendered")
CODEGRAPH_UI_ACTION_INVOKED = EventType("codegraph.ui.action.invoked")
CODEGRAPH_UI_NAVIGATION_TRIGGERED = EventType("codegraph.ui.navigation.triggered")
CODEGRAPH_UI_FEEDBACK_EMITTED = EventType("codegraph.ui.feedback.emitted")
CODEGRAPH_UI_ALERT_DISPATCHED = EventType("codegraph.ui.alert.dispatched")
CODEGRAPH_UI_DRAG_COMPLETED = EventType("codegraph.ui.drag.completed")
CODEGRAPH_UI_SELECTION_CHANGED = EventType("codegraph.ui.selection.changed")
CODEGRAPH_UI_CONFIRMATION_RESOLVED = EventType("codegraph.ui.confirmation.resolved")

# Aesthetic events
CODEGRAPH_AESTHETIC_SKIN_APPLIED = EventType("codegraph.aesthetic.skin.applied")

# Temporal events
CODEGRAPH_TEMPORAL_UNDO_EXECUTED = EventType("codegraph.temporal.undo.executed")
CODEGRAPH_TEMPORAL_RETRY_ATTEMPTED = EventType("codegraph.temporal.retry.attempted")

# Resilience events
CODEGRAPH_RESILIENCE_OFFLINE_ENTERED = EventType("codegraph.resilience.offline.entered")
CODEGRAPH_RESILIENCE_OFFLINE_SYNCED = EventType("codegraph.resilience.offline.synced")

# Structural events
CODEGRAPH_STRUCTURAL_SCOPE_DEFINED = EventType("codegraph.structural.scope.defined")

# Social events
CODEGRAPH_SOCIAL_PRESENCE_UPDATED = EventType("codegraph.social.presence.updated")
CODEGRAPH_SOCIAL_SALIENCE_TRIGGERED = EventType("codegraph.social.salience.triggered")


def all_codegraph_event_types() -> list[EventType]:
    """Return all 35 code graph event types."""
    return [
        # Data
        CODEGRAPH_ENTITY_DEFINED, CODEGRAPH_ENTITY_MODIFIED,
        CODEGRAPH_ENTITY_RELATED, CODEGRAPH_STATE_TRANSITIONED,
        # Logic
        CODEGRAPH_LOGIC_TRANSFORM_APPLIED, CODEGRAPH_LOGIC_CONDITION_EVALUATED,
        CODEGRAPH_LOGIC_SEQUENCE_EXECUTED, CODEGRAPH_LOGIC_LOOP_ITERATED,
        CODEGRAPH_LOGIC_TRIGGER_FIRED, CODEGRAPH_LOGIC_CONSTRAINT_VIOLATED,
        # IO
        CODEGRAPH_IO_QUERY_EXECUTED, CODEGRAPH_IO_COMMAND_EXECUTED,
        CODEGRAPH_IO_COMMAND_REJECTED, CODEGRAPH_IO_SUBSCRIBE_REGISTERED,
        CODEGRAPH_IO_AUTHORIZE_GRANTED, CODEGRAPH_IO_AUTHORIZE_DENIED,
        CODEGRAPH_IO_SEARCH_EXECUTED, CODEGRAPH_IO_INTEROP_SENT,
        CODEGRAPH_IO_INTEROP_RECEIVED,
        # UI
        CODEGRAPH_UI_VIEW_RENDERED, CODEGRAPH_UI_ACTION_INVOKED,
        CODEGRAPH_UI_NAVIGATION_TRIGGERED, CODEGRAPH_UI_FEEDBACK_EMITTED,
        CODEGRAPH_UI_ALERT_DISPATCHED, CODEGRAPH_UI_DRAG_COMPLETED,
        CODEGRAPH_UI_SELECTION_CHANGED, CODEGRAPH_UI_CONFIRMATION_RESOLVED,
        # Aesthetic
        CODEGRAPH_AESTHETIC_SKIN_APPLIED,
        # Temporal
        CODEGRAPH_TEMPORAL_UNDO_EXECUTED, CODEGRAPH_TEMPORAL_RETRY_ATTEMPTED,
        # Resilience
        CODEGRAPH_RESILIENCE_OFFLINE_ENTERED, CODEGRAPH_RESILIENCE_OFFLINE_SYNCED,
        # Structural
        CODEGRAPH_STRUCTURAL_SCOPE_DEFINED,
        # Social
        CODEGRAPH_SOCIAL_PRESENCE_UPDATED, CODEGRAPH_SOCIAL_SALIENCE_TRIGGERED,
    ]


# ============================================================================
# COMPOSITIONS (7) — Named groupings of code graph atoms
# ============================================================================

@dataclass(frozen=True)
class CodeGraphComposition:
    """A named grouping of code graph atoms."""

    name: str
    primitives: list[str]


def board() -> CodeGraphComposition:
    """Board composition — kanban/card-based views."""
    return CodeGraphComposition(
        name="Board",
        primitives=[
            "CGLayout", "CGList", "CGQuery", "CGDisplay", "CGDrag",
            "CGCommand", "CGEmpty", "CGAction", "CGLoop", "CGState",
        ],
    )


def detail() -> CodeGraphComposition:
    """Detail composition — single-entity detail views."""
    return CodeGraphComposition(
        name="Detail",
        primitives=[
            "CGLayout", "CGDisplay", "CGForm", "CGThread", "CGAudit",
            "CGHistory", "CGList", "CGAction", "CGNavigation",
        ],
    )


def feed() -> CodeGraphComposition:
    """Feed composition — chronological stream views."""
    return CodeGraphComposition(
        name="Feed",
        primitives=[
            "CGList", "CGQuery", "CGDisplay", "CGAvatar", "CGSubscribe",
            "CGPagination", "CGRecency",
        ],
    )


def dashboard() -> CodeGraphComposition:
    """Dashboard composition — aggregated data overview."""
    return CodeGraphComposition(
        name="Dashboard",
        primitives=[
            "CGLayout", "CGDisplay", "CGQuery", "CGTransform", "CGSalience",
        ],
    )


def inbox() -> CodeGraphComposition:
    """Inbox composition — message/notification list."""
    return CodeGraphComposition(
        name="Inbox",
        primitives=[
            "CGList", "CGQuery", "CGDisplay", "CGAvatar", "CGSalience",
            "CGAction", "CGSelection", "CGEmpty",
        ],
    )


def wizard() -> CodeGraphComposition:
    """Wizard composition — multi-step guided flow."""
    return CodeGraphComposition(
        name="Wizard",
        primitives=[
            "CGSequence", "CGForm", "CGInput", "CGDisplay", "CGAction",
            "CGNavigation", "CGConsequencePreview", "CGConstraint",
        ],
    )


def skin() -> CodeGraphComposition:
    """Skin composition — visual theming."""
    return CodeGraphComposition(
        name="Skin",
        primitives=[
            "CGPalette", "CGTypography", "CGSpacing", "CGElevation",
            "CGMotion", "CGDensity", "CGShape",
        ],
    )


def all_codegraph_compositions() -> list[CodeGraphComposition]:
    """Return all 7 code graph compositions."""
    return [
        board(),
        detail(),
        feed(),
        dashboard(),
        inbox(),
        wizard(),
        skin(),
    ]
