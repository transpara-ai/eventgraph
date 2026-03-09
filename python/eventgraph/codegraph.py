"""Code Graph primitives — 61 primitives, 7 compositions, 35 event types.

All Code Graph primitives operate at Layer 5 (Code Graph). They define how
application UIs, data models, logic, and interactions are represented as
primitives on the event graph.
"""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any

from .event import Event
from .primitive import (
    Mutation,
    Registry,
    Snapshot,
    UpdateState,
)
from .types import (
    Cadence,
    EventType,
    Layer,
    PrimitiveID,
    SubscriptionPattern,
)


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
# Code Graph Primitive Base
# ============================================================================

_CODEGRAPH_LAYER = Layer(5)
_CADENCE_1 = Cadence(1)


class _CodeGraphBase:
    """Common implementation for all 61 code graph primitives.

    Subclasses set _name and _subs at the class level.
    Each primitive's process() counts events and returns UpdateState mutations.
    """

    _name: str
    _subs: list[str]

    def id(self) -> PrimitiveID:
        return PrimitiveID(self._name)

    def layer(self) -> Layer:
        return _CODEGRAPH_LAYER

    def subscriptions(self) -> list[SubscriptionPattern]:
        return [SubscriptionPattern(s) for s in self._subs]

    def cadence(self) -> Cadence:
        return _CADENCE_1

    def process(
        self,
        tick: int,
        events: list[Event],
        snapshot: Snapshot,
    ) -> list[Mutation]:
        pid = self.id()
        return [
            UpdateState(primitive_id=pid, key="eventsProcessed", value=len(events)),
            UpdateState(primitive_id=pid, key="lastTick", value=tick),
        ]


def _cg_def(name: str, subs: list[str]) -> type:
    """Create a code graph primitive class with the given name and subscriptions."""
    cls = type(name, (_CodeGraphBase,), {
        "_name": f"CG{name}",
        "_subs": subs,
    })
    cls.__module__ = __name__
    cls.__qualname__ = name
    return cls


# ============================================================================
# DATA PRIMITIVES (6)
# ============================================================================

CGEntity = _cg_def("Entity", ["codegraph.entity.*", "codegraph.io.command.*"])
CGProperty = _cg_def("Property", ["codegraph.entity.*"])
CGRelation = _cg_def("Relation", ["codegraph.entity.*"])
CGCollection = _cg_def("Collection", ["codegraph.entity.*", "codegraph.io.query.*"])
CGState = _cg_def("State", ["codegraph.state.*", "codegraph.io.command.*"])
CGEvent = _cg_def("Event", ["codegraph.*"])

# ============================================================================
# LOGIC PRIMITIVES (6)
# ============================================================================

CGTransform = _cg_def("Transform", ["codegraph.logic.transform.*"])
CGCondition = _cg_def("Condition", ["codegraph.logic.condition.*"])
CGSequence = _cg_def("Sequence", ["codegraph.logic.sequence.*"])
CGLoop = _cg_def("Loop", ["codegraph.logic.loop.*"])
CGTrigger = _cg_def("Trigger", ["codegraph.logic.trigger.*", "codegraph.entity.*"])
CGConstraint = _cg_def("Constraint", ["codegraph.logic.constraint.*", "codegraph.io.command.*"])

# ============================================================================
# IO PRIMITIVES (6)
# ============================================================================

CGQuery = _cg_def("Query", ["codegraph.io.query.*"])
CGCommand = _cg_def("Command", ["codegraph.io.command.*"])
CGSubscribe = _cg_def("Subscribe", ["codegraph.io.subscribe.*"])
CGAuthorize = _cg_def("Authorize", ["codegraph.io.authorize.*", "authority.*"])
CGSearch = _cg_def("Search", ["codegraph.io.search.*"])
CGInterop = _cg_def("Interop", ["codegraph.io.interop.*"])

# ============================================================================
# UI PRIMITIVES (19)
# ============================================================================

CGDisplay = _cg_def("Display", ["codegraph.ui.*"])
CGInput = _cg_def("Input", ["codegraph.ui.*"])
CGLayout = _cg_def("Layout", ["codegraph.ui.*"])
CGList = _cg_def("List", ["codegraph.ui.*", "codegraph.io.query.*"])
CGForm = _cg_def("Form", ["codegraph.ui.*", "codegraph.io.command.*"])
CGAction = _cg_def("Action", ["codegraph.ui.action.*"])
CGNavigation = _cg_def("Navigation", ["codegraph.ui.navigation.*"])
CGView = _cg_def("View", ["codegraph.ui.view.*"])
CGFeedback = _cg_def("Feedback", ["codegraph.ui.feedback.*"])
CGAlert = _cg_def("Alert", ["codegraph.ui.alert.*"])
CGThread = _cg_def("Thread", ["codegraph.ui.*", "codegraph.entity.*"])
CGAvatar = _cg_def("Avatar", ["codegraph.ui.*"])
CGAudit = _cg_def("Audit", ["codegraph.*"])
CGDrag = _cg_def("Drag", ["codegraph.ui.drag.*"])
CGSelection = _cg_def("Selection", ["codegraph.ui.selection.*"])
CGConfirmation = _cg_def("Confirmation", ["codegraph.ui.confirmation.*"])
CGEmpty = _cg_def("Empty", ["codegraph.ui.*"])
CGLoading = _cg_def("Loading", ["codegraph.ui.*"])
CGPagination = _cg_def("Pagination", ["codegraph.ui.*", "codegraph.io.query.*"])

# ============================================================================
# AESTHETIC PRIMITIVES (7)
# ============================================================================

CGPalette = _cg_def("Palette", ["codegraph.aesthetic.*"])
CGTypography = _cg_def("Typography", ["codegraph.aesthetic.*"])
CGSpacing = _cg_def("Spacing", ["codegraph.aesthetic.*"])
CGElevation = _cg_def("Elevation", ["codegraph.aesthetic.*"])
CGMotion = _cg_def("Motion", ["codegraph.aesthetic.*"])
CGDensity = _cg_def("Density", ["codegraph.aesthetic.*"])
CGShape = _cg_def("Shape", ["codegraph.aesthetic.*"])

# ============================================================================
# ACCESSIBILITY PRIMITIVES (4)
# ============================================================================

CGAnnounce = _cg_def("Announce", ["codegraph.ui.*", "codegraph.aesthetic.*"])
CGFocus = _cg_def("Focus", ["codegraph.ui.*"])
CGContrast = _cg_def("Contrast", ["codegraph.aesthetic.*"])
CGSimplify = _cg_def("Simplify", ["codegraph.ui.*", "codegraph.aesthetic.*"])

# ============================================================================
# TEMPORAL PRIMITIVES (3)
# ============================================================================

CGRecency = _cg_def("Recency", ["codegraph.entity.*", "codegraph.temporal.*"])
CGHistory = _cg_def("History", ["codegraph.entity.*", "codegraph.temporal.*"])
CGLiveness = _cg_def("Liveness", ["codegraph.io.subscribe.*", "codegraph.social.presence.*"])

# ============================================================================
# RESILIENCE PRIMITIVES (4)
# ============================================================================

CGUndo = _cg_def("Undo", ["codegraph.temporal.undo.*", "codegraph.io.command.*"])
CGRetry = _cg_def("Retry", ["codegraph.temporal.retry.*", "codegraph.io.command.*"])
CGFallback = _cg_def("Fallback", ["codegraph.resilience.*", "codegraph.ui.*"])
CGOffline = _cg_def("Offline", ["codegraph.resilience.offline.*"])

# ============================================================================
# STRUCTURAL PRIMITIVES (3)
# ============================================================================

CGScope = _cg_def("Scope", ["codegraph.structural.*"])
CGFormat = _cg_def("Format", ["codegraph.io.*", "codegraph.structural.*"])
CGGesture = _cg_def("Gesture", ["codegraph.ui.*"])

# ============================================================================
# SOCIAL PRIMITIVES (3)
# ============================================================================

CGPresence = _cg_def("Presence", ["codegraph.social.presence.*"])
CGSalience = _cg_def("Salience", ["codegraph.social.salience.*", "codegraph.entity.*"])
CGConsequencePreview = _cg_def("ConsequencePreview", ["codegraph.ui.confirmation.*", "codegraph.io.command.*"])


# ============================================================================
# ALL PRIMITIVE CLASSES
# ============================================================================

ALL_CODEGRAPH_PRIMITIVE_CLASSES: list[type] = [
    # Data (6)
    CGEntity, CGProperty, CGRelation, CGCollection, CGState, CGEvent,
    # Logic (6)
    CGTransform, CGCondition, CGSequence, CGLoop, CGTrigger, CGConstraint,
    # IO (6)
    CGQuery, CGCommand, CGSubscribe, CGAuthorize, CGSearch, CGInterop,
    # UI (19)
    CGDisplay, CGInput, CGLayout, CGList, CGForm, CGAction, CGNavigation,
    CGView, CGFeedback, CGAlert, CGThread, CGAvatar, CGAudit, CGDrag,
    CGSelection, CGConfirmation, CGEmpty, CGLoading, CGPagination,
    # Aesthetic (7)
    CGPalette, CGTypography, CGSpacing, CGElevation, CGMotion, CGDensity, CGShape,
    # Accessibility (4)
    CGAnnounce, CGFocus, CGContrast, CGSimplify,
    # Temporal (3)
    CGRecency, CGHistory, CGLiveness,
    # Resilience (4)
    CGUndo, CGRetry, CGFallback, CGOffline,
    # Structural (3)
    CGScope, CGFormat, CGGesture,
    # Social (3)
    CGPresence, CGSalience, CGConsequencePreview,
]


def all_codegraph_primitives() -> list[_CodeGraphBase]:
    """Return all 61 code graph primitives."""
    return [cls() for cls in ALL_CODEGRAPH_PRIMITIVE_CLASSES]


def register_all_codegraph(reg: Registry) -> None:
    """Register and activate all 61 code graph primitives with the given registry."""
    for p in all_codegraph_primitives():
        reg.register(p)
        reg.activate(p.id())


def is_codegraph_primitive(pid: PrimitiveID) -> bool:
    """Return True if the primitive ID belongs to the code graph layer."""
    return pid.value.startswith("CG")


# ============================================================================
# COMPOSITIONS (7) — Named groupings of code graph primitives
# ============================================================================

@dataclass(frozen=True)
class CodeGraphComposition:
    """A named grouping of code graph primitives."""

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
