/**
 * Code Graph primitives — 61 primitives at Layer 5 (Code Graph) plus 7 named compositions.
 * Ported from the Go reference implementation.
 */
import { Event } from "./event.js";
import { Primitive, Mutation, Snapshot, Registry } from "./primitive.js";
import { PrimitiveId, Layer, Cadence, EventType, SubscriptionPattern } from "./types.js";

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

// ── Shared constants ─────────────────────────────────────────────────────

const CG_LAYER = new Layer(5);
const CADENCE_1 = new Cadence(1);

// ── Helper to build a simple counting primitive ──────────────────────────

function eventTypeValue(ev: Event): string {
  return ev.type.value;
}

// ════════════════════════════════════════════════════════════════════════
// DATA PRIMITIVES (6)
// ════════════════════════════════════════════════════════════════════════

export class CGEntityPrimitive implements Primitive {
  id() { return new PrimitiveId("CGEntity"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.entity.*"),
      new SubscriptionPattern("codegraph.io.command.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGPropertyPrimitive implements Primitive {
  id() { return new PrimitiveId("CGProperty"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.entity.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGRelationPrimitive implements Primitive {
  id() { return new PrimitiveId("CGRelation"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.entity.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGCollectionPrimitive implements Primitive {
  id() { return new PrimitiveId("CGCollection"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.entity.*"),
      new SubscriptionPattern("codegraph.io.query.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGStatePrimitive implements Primitive {
  id() { return new PrimitiveId("CGState"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.state.*"),
      new SubscriptionPattern("codegraph.io.command.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGEventPrimitive implements Primitive {
  id() { return new PrimitiveId("CGEvent"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// ════════════════════════════════════════════════════════════════════════
// LOGIC PRIMITIVES (6)
// ════════════════════════════════════════════════════════════════════════

export class CGTransformPrimitive implements Primitive {
  id() { return new PrimitiveId("CGTransform"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.logic.transform.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGConditionPrimitive implements Primitive {
  id() { return new PrimitiveId("CGCondition"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.logic.condition.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGSequencePrimitive implements Primitive {
  id() { return new PrimitiveId("CGSequence"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.logic.sequence.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGLoopPrimitive implements Primitive {
  id() { return new PrimitiveId("CGLoop"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.logic.loop.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGTriggerPrimitive implements Primitive {
  id() { return new PrimitiveId("CGTrigger"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.logic.trigger.*"),
      new SubscriptionPattern("codegraph.entity.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGConstraintPrimitive implements Primitive {
  id() { return new PrimitiveId("CGConstraint"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.logic.constraint.*"),
      new SubscriptionPattern("codegraph.io.command.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// ════════════════════════════════════════════════════════════════════════
// IO PRIMITIVES (6)
// ════════════════════════════════════════════════════════════════════════

export class CGQueryPrimitive implements Primitive {
  id() { return new PrimitiveId("CGQuery"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.io.query.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGCommandPrimitive implements Primitive {
  id() { return new PrimitiveId("CGCommand"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.io.command.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGSubscribePrimitive implements Primitive {
  id() { return new PrimitiveId("CGSubscribe"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.io.subscribe.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGAuthorizePrimitive implements Primitive {
  id() { return new PrimitiveId("CGAuthorize"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.io.authorize.*"),
      new SubscriptionPattern("authority.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGSearchPrimitive implements Primitive {
  id() { return new PrimitiveId("CGSearch"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.io.search.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGInteropPrimitive implements Primitive {
  id() { return new PrimitiveId("CGInterop"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.io.interop.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// ════════════════════════════════════════════════════════════════════════
// UI PRIMITIVES (19)
// ════════════════════════════════════════════════════════════════════════

export class CGDisplayPrimitive implements Primitive {
  id() { return new PrimitiveId("CGDisplay"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGInputPrimitive implements Primitive {
  id() { return new PrimitiveId("CGInput"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGLayoutPrimitive implements Primitive {
  id() { return new PrimitiveId("CGLayout"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGListPrimitive implements Primitive {
  id() { return new PrimitiveId("CGList"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.ui.*"),
      new SubscriptionPattern("codegraph.io.query.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGFormPrimitive implements Primitive {
  id() { return new PrimitiveId("CGForm"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.ui.*"),
      new SubscriptionPattern("codegraph.io.command.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGActionPrimitive implements Primitive {
  id() { return new PrimitiveId("CGAction"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.action.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGNavigationPrimitive implements Primitive {
  id() { return new PrimitiveId("CGNavigation"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.navigation.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGViewPrimitive implements Primitive {
  id() { return new PrimitiveId("CGView"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.view.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGFeedbackPrimitive implements Primitive {
  id() { return new PrimitiveId("CGFeedback"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.feedback.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGAlertPrimitive implements Primitive {
  id() { return new PrimitiveId("CGAlert"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.alert.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGThreadPrimitive implements Primitive {
  id() { return new PrimitiveId("CGThread"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.ui.*"),
      new SubscriptionPattern("codegraph.entity.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGAvatarPrimitive implements Primitive {
  id() { return new PrimitiveId("CGAvatar"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGAuditPrimitive implements Primitive {
  id() { return new PrimitiveId("CGAudit"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGDragPrimitive implements Primitive {
  id() { return new PrimitiveId("CGDrag"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.drag.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGSelectionPrimitive implements Primitive {
  id() { return new PrimitiveId("CGSelection"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.selection.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGConfirmationPrimitive implements Primitive {
  id() { return new PrimitiveId("CGConfirmation"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.confirmation.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGEmptyPrimitive implements Primitive {
  id() { return new PrimitiveId("CGEmpty"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGLoadingPrimitive implements Primitive {
  id() { return new PrimitiveId("CGLoading"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGPaginationPrimitive implements Primitive {
  id() { return new PrimitiveId("CGPagination"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.ui.*"),
      new SubscriptionPattern("codegraph.io.query.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// ════════════════════════════════════════════════════════════════════════
// AESTHETIC PRIMITIVES (7)
// ════════════════════════════════════════════════════════════════════════

export class CGPalettePrimitive implements Primitive {
  id() { return new PrimitiveId("CGPalette"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.aesthetic.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGTypographyPrimitive implements Primitive {
  id() { return new PrimitiveId("CGTypography"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.aesthetic.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGSpacingPrimitive implements Primitive {
  id() { return new PrimitiveId("CGSpacing"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.aesthetic.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGElevationPrimitive implements Primitive {
  id() { return new PrimitiveId("CGElevation"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.aesthetic.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGMotionPrimitive implements Primitive {
  id() { return new PrimitiveId("CGMotion"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.aesthetic.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGDensityPrimitive implements Primitive {
  id() { return new PrimitiveId("CGDensity"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.aesthetic.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGShapePrimitive implements Primitive {
  id() { return new PrimitiveId("CGShape"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.aesthetic.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// ════════════════════════════════════════════════════════════════════════
// ACCESSIBILITY PRIMITIVES (4)
// ════════════════════════════════════════════════════════════════════════

export class CGAnnouncePrimitive implements Primitive {
  id() { return new PrimitiveId("CGAnnounce"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.ui.*"),
      new SubscriptionPattern("codegraph.aesthetic.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGFocusPrimitive implements Primitive {
  id() { return new PrimitiveId("CGFocus"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGContrastPrimitive implements Primitive {
  id() { return new PrimitiveId("CGContrast"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.aesthetic.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGSimplifyPrimitive implements Primitive {
  id() { return new PrimitiveId("CGSimplify"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.ui.*"),
      new SubscriptionPattern("codegraph.aesthetic.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// ════════════════════════════════════════════════════════════════════════
// TEMPORAL PRIMITIVES (3)
// ════════════════════════════════════════════════════════════════════════

export class CGRecencyPrimitive implements Primitive {
  id() { return new PrimitiveId("CGRecency"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.entity.*"),
      new SubscriptionPattern("codegraph.temporal.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGHistoryPrimitive implements Primitive {
  id() { return new PrimitiveId("CGHistory"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.entity.*"),
      new SubscriptionPattern("codegraph.temporal.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGLivenessPrimitive implements Primitive {
  id() { return new PrimitiveId("CGLiveness"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.io.subscribe.*"),
      new SubscriptionPattern("codegraph.social.presence.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// ════════════════════════════════════════════════════════════════════════
// RESILIENCE PRIMITIVES (4)
// ════════════════════════════════════════════════════════════════════════

export class CGUndoPrimitive implements Primitive {
  id() { return new PrimitiveId("CGUndo"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.temporal.undo.*"),
      new SubscriptionPattern("codegraph.io.command.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGRetryPrimitive implements Primitive {
  id() { return new PrimitiveId("CGRetry"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.temporal.retry.*"),
      new SubscriptionPattern("codegraph.io.command.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGFallbackPrimitive implements Primitive {
  id() { return new PrimitiveId("CGFallback"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.resilience.*"),
      new SubscriptionPattern("codegraph.ui.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGOfflinePrimitive implements Primitive {
  id() { return new PrimitiveId("CGOffline"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.resilience.offline.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// ════════════════════════════════════════════════════════════════════════
// STRUCTURAL PRIMITIVES (3)
// ════════════════════════════════════════════════════════════════════════

export class CGScopePrimitive implements Primitive {
  id() { return new PrimitiveId("CGScope"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.structural.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGFormatPrimitive implements Primitive {
  id() { return new PrimitiveId("CGFormat"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.io.*"),
      new SubscriptionPattern("codegraph.structural.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGGesturePrimitive implements Primitive {
  id() { return new PrimitiveId("CGGesture"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.ui.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// ════════════════════════════════════════════════════════════════════════
// SOCIAL PRIMITIVES (3)
// ════════════════════════════════════════════════════════════════════════

export class CGPresencePrimitive implements Primitive {
  id() { return new PrimitiveId("CGPresence"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [new SubscriptionPattern("codegraph.social.presence.*")];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGSaliencePrimitive implements Primitive {
  id() { return new PrimitiveId("CGSalience"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.social.salience.*"),
      new SubscriptionPattern("codegraph.entity.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CGConsequencePreviewPrimitive implements Primitive {
  id() { return new PrimitiveId("CGConsequencePreview"); }
  layer() { return CG_LAYER; }
  cadence() { return CADENCE_1; }
  subscriptions() {
    return [
      new SubscriptionPattern("codegraph.ui.confirmation.*"),
      new SubscriptionPattern("codegraph.io.command.*"),
    ];
  }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// ════════════════════════════════════════════════════════════════════════
// REGISTRATION
// ════════════════════════════════════════════════════════════════════════

/** Returns all 61 code graph primitives. */
export function allCodeGraphPrimitives(): Primitive[] {
  return [
    // Data (6)
    new CGEntityPrimitive(),
    new CGPropertyPrimitive(),
    new CGRelationPrimitive(),
    new CGCollectionPrimitive(),
    new CGStatePrimitive(),
    new CGEventPrimitive(),
    // Logic (6)
    new CGTransformPrimitive(),
    new CGConditionPrimitive(),
    new CGSequencePrimitive(),
    new CGLoopPrimitive(),
    new CGTriggerPrimitive(),
    new CGConstraintPrimitive(),
    // IO (6)
    new CGQueryPrimitive(),
    new CGCommandPrimitive(),
    new CGSubscribePrimitive(),
    new CGAuthorizePrimitive(),
    new CGSearchPrimitive(),
    new CGInteropPrimitive(),
    // UI (19)
    new CGDisplayPrimitive(),
    new CGInputPrimitive(),
    new CGLayoutPrimitive(),
    new CGListPrimitive(),
    new CGFormPrimitive(),
    new CGActionPrimitive(),
    new CGNavigationPrimitive(),
    new CGViewPrimitive(),
    new CGFeedbackPrimitive(),
    new CGAlertPrimitive(),
    new CGThreadPrimitive(),
    new CGAvatarPrimitive(),
    new CGAuditPrimitive(),
    new CGDragPrimitive(),
    new CGSelectionPrimitive(),
    new CGConfirmationPrimitive(),
    new CGEmptyPrimitive(),
    new CGLoadingPrimitive(),
    new CGPaginationPrimitive(),
    // Aesthetic (7)
    new CGPalettePrimitive(),
    new CGTypographyPrimitive(),
    new CGSpacingPrimitive(),
    new CGElevationPrimitive(),
    new CGMotionPrimitive(),
    new CGDensityPrimitive(),
    new CGShapePrimitive(),
    // Accessibility (4)
    new CGAnnouncePrimitive(),
    new CGFocusPrimitive(),
    new CGContrastPrimitive(),
    new CGSimplifyPrimitive(),
    // Temporal (3)
    new CGRecencyPrimitive(),
    new CGHistoryPrimitive(),
    new CGLivenessPrimitive(),
    // Resilience (4)
    new CGUndoPrimitive(),
    new CGRetryPrimitive(),
    new CGFallbackPrimitive(),
    new CGOfflinePrimitive(),
    // Structural (3)
    new CGScopePrimitive(),
    new CGFormatPrimitive(),
    new CGGesturePrimitive(),
    // Social (3)
    new CGPresencePrimitive(),
    new CGSaliencePrimitive(),
    new CGConsequencePreviewPrimitive(),
  ];
}

/** Registers all 61 code graph primitives with the given registry and activates them. */
export function registerAllCodeGraph(registry: Registry): void {
  for (const p of allCodeGraphPrimitives()) {
    registry.register(p);
    registry.activate(p.id());
  }
}

/** Returns true if the primitive ID belongs to the code graph layer. */
export function isCodeGraphPrimitive(id: PrimitiveId): boolean {
  return id.value.startsWith("CG");
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
