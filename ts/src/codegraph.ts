/**
 * Code Graph event types and compositions — vocabulary/content for the Code Graph layer.
 */
import { EventType } from "./types.js";

// ── Code Graph Event Types ────────────────────────────────────────────────

// Data
export const CGEventTypeEntityDefined = new EventType("codegraph.entity.defined");
export const CGEventTypeEntityModified = new EventType("codegraph.entity.modified");
export const CGEventTypeEntityRelated = new EventType("codegraph.entity.related");
export const CGEventTypeStateTransitioned = new EventType("codegraph.state.transitioned");

// Logic
export const CGEventTypeTransformApplied = new EventType("codegraph.logic.transform.applied");
export const CGEventTypeConditionEvaluated = new EventType("codegraph.logic.condition.evaluated");
export const CGEventTypeSequenceExecuted = new EventType("codegraph.logic.sequence.executed");
export const CGEventTypeLoopIterated = new EventType("codegraph.logic.loop.iterated");
export const CGEventTypeTriggerFired = new EventType("codegraph.logic.trigger.fired");
export const CGEventTypeConstraintViolated = new EventType("codegraph.logic.constraint.violated");

// IO
export const CGEventTypeQueryExecuted = new EventType("codegraph.io.query.executed");
export const CGEventTypeCommandExecuted = new EventType("codegraph.io.command.executed");
export const CGEventTypeCommandRejected = new EventType("codegraph.io.command.rejected");
export const CGEventTypeSubscribeRegistered = new EventType("codegraph.io.subscribe.registered");
export const CGEventTypeAuthorizeGranted = new EventType("codegraph.io.authorize.granted");
export const CGEventTypeAuthorizeDenied = new EventType("codegraph.io.authorize.denied");
export const CGEventTypeSearchExecuted = new EventType("codegraph.io.search.executed");
export const CGEventTypeInteropSent = new EventType("codegraph.io.interop.sent");
export const CGEventTypeInteropReceived = new EventType("codegraph.io.interop.received");

// UI
export const CGEventTypeViewRendered = new EventType("codegraph.ui.view.rendered");
export const CGEventTypeActionInvoked = new EventType("codegraph.ui.action.invoked");
export const CGEventTypeNavigationTriggered = new EventType("codegraph.ui.navigation.triggered");
export const CGEventTypeFeedbackEmitted = new EventType("codegraph.ui.feedback.emitted");
export const CGEventTypeAlertDispatched = new EventType("codegraph.ui.alert.dispatched");
export const CGEventTypeDragCompleted = new EventType("codegraph.ui.drag.completed");
export const CGEventTypeSelectionChanged = new EventType("codegraph.ui.selection.changed");
export const CGEventTypeConfirmationResolved = new EventType("codegraph.ui.confirmation.resolved");

// Aesthetic
export const CGEventTypeSkinApplied = new EventType("codegraph.aesthetic.skin.applied");

// Temporal
export const CGEventTypeUndoExecuted = new EventType("codegraph.temporal.undo.executed");
export const CGEventTypeRetryAttempted = new EventType("codegraph.temporal.retry.attempted");

// Resilience
export const CGEventTypeOfflineEntered = new EventType("codegraph.resilience.offline.entered");
export const CGEventTypeOfflineSynced = new EventType("codegraph.resilience.offline.synced");

// Structural
export const CGEventTypeScopeDefined = new EventType("codegraph.structural.scope.defined");

// Social
export const CGEventTypePresenceUpdated = new EventType("codegraph.social.presence.updated");
export const CGEventTypeSalienceTriggered = new EventType("codegraph.social.salience.triggered");

export function allCodeGraphEventTypes(): EventType[] {
  return [
    // Data
    CGEventTypeEntityDefined, CGEventTypeEntityModified, CGEventTypeEntityRelated, CGEventTypeStateTransitioned,
    // Logic
    CGEventTypeTransformApplied, CGEventTypeConditionEvaluated, CGEventTypeSequenceExecuted,
    CGEventTypeLoopIterated, CGEventTypeTriggerFired, CGEventTypeConstraintViolated,
    // IO
    CGEventTypeQueryExecuted, CGEventTypeCommandExecuted, CGEventTypeCommandRejected,
    CGEventTypeSubscribeRegistered, CGEventTypeAuthorizeGranted, CGEventTypeAuthorizeDenied,
    CGEventTypeSearchExecuted, CGEventTypeInteropSent, CGEventTypeInteropReceived,
    // UI
    CGEventTypeViewRendered, CGEventTypeActionInvoked, CGEventTypeNavigationTriggered,
    CGEventTypeFeedbackEmitted, CGEventTypeAlertDispatched, CGEventTypeDragCompleted,
    CGEventTypeSelectionChanged, CGEventTypeConfirmationResolved,
    // Aesthetic
    CGEventTypeSkinApplied,
    // Temporal
    CGEventTypeUndoExecuted, CGEventTypeRetryAttempted,
    // Resilience
    CGEventTypeOfflineEntered, CGEventTypeOfflineSynced,
    // Structural
    CGEventTypeScopeDefined,
    // Social
    CGEventTypePresenceUpdated, CGEventTypeSalienceTriggered,
  ];
}

// ════════════════════════════════════════════════════════════════════════
// COMPOSITIONS (7)
// ════════════════════════════════════════════════════════════════════════

export interface CodeGraphComposition {
  name: string;
  primitives: string[];
}

/** Board — Kanban-style board view. */
export function boardComposition(): CodeGraphComposition {
  return {
    name: "Board",
    primitives: [
      "CGLayout", "CGList", "CGQuery", "CGDisplay", "CGDrag",
      "CGCommand", "CGEmpty", "CGAction", "CGLoop", "CGState",
    ],
  };
}

/** Detail — Entity detail view with history. */
export function detailComposition(): CodeGraphComposition {
  return {
    name: "Detail",
    primitives: [
      "CGLayout", "CGDisplay", "CGForm", "CGThread", "CGAudit",
      "CGHistory", "CGList", "CGAction", "CGNavigation",
    ],
  };
}

/** Feed — Chronological feed with pagination. */
export function feedComposition(): CodeGraphComposition {
  return {
    name: "Feed",
    primitives: [
      "CGList", "CGQuery", "CGDisplay", "CGAvatar", "CGSubscribe",
      "CGPagination", "CGRecency",
    ],
  };
}

/** Dashboard — Summary metrics and key indicators. */
export function dashboardComposition(): CodeGraphComposition {
  return {
    name: "Dashboard",
    primitives: ["CGLayout", "CGDisplay", "CGQuery", "CGTransform", "CGSalience"],
  };
}

/** Inbox — Prioritised notification list. */
export function inboxComposition(): CodeGraphComposition {
  return {
    name: "Inbox",
    primitives: [
      "CGList", "CGQuery", "CGDisplay", "CGAvatar", "CGSalience",
      "CGAction", "CGSelection", "CGEmpty",
    ],
  };
}

/** Wizard — Multi-step guided flow. */
export function wizardComposition(): CodeGraphComposition {
  return {
    name: "Wizard",
    primitives: [
      "CGSequence", "CGForm", "CGInput", "CGDisplay", "CGAction",
      "CGNavigation", "CGConsequencePreview", "CGConstraint",
    ],
  };
}

/** Skin — Visual theme composition. */
export function skinComposition(): CodeGraphComposition {
  return {
    name: "Skin",
    primitives: [
      "CGPalette", "CGTypography", "CGSpacing", "CGElevation",
      "CGMotion", "CGDensity", "CGShape",
    ],
  };
}

/** Returns all 7 named compositions. */
export function allCodeGraphCompositions(): CodeGraphComposition[] {
  return [
    boardComposition(),
    detailComposition(),
    feedComposition(),
    dashboardComposition(),
    inboxComposition(),
    wizardComposition(),
    skinComposition(),
  ];
}
