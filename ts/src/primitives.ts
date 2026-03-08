import { Event } from "./event.js";
import { Primitive, Mutation, Snapshot } from "./primitive.js";
import { PrimitiveId, Layer, Cadence, SubscriptionPattern } from "./types.js";

// -- Layer 0 (Foundation) -- 45 primitives --

export class EventPrimitive implements Primitive {
  id() { return new PrimitiveId("Event"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class EventStorePrimitive implements Primitive {
  id() { return new PrimitiveId("EventStore"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("store.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ClockPrimitive implements Primitive {
  id() { return new PrimitiveId("Clock"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("clock.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class HashPrimitive implements Primitive {
  id() { return new PrimitiveId("Hash"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SelfPrimitive implements Primitive {
  id() { return new PrimitiveId("Self"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("system.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CausalLinkPrimitive implements Primitive {
  id() { return new PrimitiveId("CausalLink"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AncestryPrimitive implements Primitive {
  id() { return new PrimitiveId("Ancestry"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("query.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DescendancyPrimitive implements Primitive {
  id() { return new PrimitiveId("Descendancy"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("query.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class FirstCausePrimitive implements Primitive {
  id() { return new PrimitiveId("FirstCause"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("query.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ActorIDPrimitive implements Primitive {
  id() { return new PrimitiveId("ActorID"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("actor.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ActorRegistryPrimitive implements Primitive {
  id() { return new PrimitiveId("ActorRegistry"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("actor.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SignaturePrimitive implements Primitive {
  id() { return new PrimitiveId("Signature"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class VerifyPrimitive implements Primitive {
  id() { return new PrimitiveId("Verify"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ExpectationPrimitive implements Primitive {
  id() { return new PrimitiveId("Expectation"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class TimeoutPrimitive implements Primitive {
  id() { return new PrimitiveId("Timeout"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("authority.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ViolationPrimitive implements Primitive {
  id() { return new PrimitiveId("Violation"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("violation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SeverityPrimitive implements Primitive {
  id() { return new PrimitiveId("Severity"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("violation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class TrustScorePrimitive implements Primitive {
  id() { return new PrimitiveId("TrustScore"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("trust.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class TrustUpdatePrimitive implements Primitive {
  id() { return new PrimitiveId("TrustUpdate"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("trust.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CorroborationPrimitive implements Primitive {
  id() { return new PrimitiveId("Corroboration"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("trust.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ContradictionPrimitive implements Primitive {
  id() { return new PrimitiveId("Contradiction"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("trust.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ConfidencePrimitive implements Primitive {
  id() { return new PrimitiveId("Confidence"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("decision.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class EvidencePrimitive implements Primitive {
  id() { return new PrimitiveId("Evidence"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RevisionPrimitive implements Primitive {
  id() { return new PrimitiveId("Revision"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("grammar.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class UncertaintyPrimitive implements Primitive {
  id() { return new PrimitiveId("Uncertainty"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("decision.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class InstrumentationSpecPrimitive implements Primitive {
  id() { return new PrimitiveId("InstrumentationSpec"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("health.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CoverageCheckPrimitive implements Primitive {
  id() { return new PrimitiveId("CoverageCheck"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("health.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class GapPrimitive implements Primitive {
  id() { return new PrimitiveId("Gap"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class BlindPrimitive implements Primitive {
  id() { return new PrimitiveId("Blind"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("health.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PathQueryPrimitive implements Primitive {
  id() { return new PrimitiveId("PathQuery"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("query.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SubgraphExtractPrimitive implements Primitive {
  id() { return new PrimitiveId("SubgraphExtract"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("query.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AnnotatePrimitive implements Primitive {
  id() { return new PrimitiveId("Annotate"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("grammar.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class TimelinePrimitive implements Primitive {
  id() { return new PrimitiveId("Timeline"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("query.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class HashChainPrimitive implements Primitive {
  id() { return new PrimitiveId("HashChain"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ChainVerifyPrimitive implements Primitive {
  id() { return new PrimitiveId("ChainVerify"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("chain.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class WitnessPrimitive implements Primitive {
  id() { return new PrimitiveId("Witness"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class IntegrityViolationPrimitive implements Primitive {
  id() { return new PrimitiveId("IntegrityViolation"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("chain.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PatternPrimitive implements Primitive {
  id() { return new PrimitiveId("Pattern"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DeceptionIndicatorPrimitive implements Primitive {
  id() { return new PrimitiveId("DeceptionIndicator"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("trust.*"), new SubscriptionPattern("violation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SuspicionPrimitive implements Primitive {
  id() { return new PrimitiveId("Suspicion"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("trust.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class QuarantinePrimitive implements Primitive {
  id() { return new PrimitiveId("Quarantine"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("actor.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class GraphHealthPrimitive implements Primitive {
  id() { return new PrimitiveId("GraphHealth"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("health.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class InvariantPrimitive implements Primitive {
  id() { return new PrimitiveId("Invariant"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class InvariantCheckPrimitive implements Primitive {
  id() { return new PrimitiveId("InvariantCheck"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("chain.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class BootstrapPrimitive implements Primitive {
  id() { return new PrimitiveId("Bootstrap"); }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("system.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 1 (Will) -- 12 primitives --

export class GoalPrimitive implements Primitive {
  id() { return new PrimitiveId("Goal"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("decision.*"), new SubscriptionPattern("authority.*"), new SubscriptionPattern("actor.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PlanPrimitive implements Primitive {
  id() { return new PrimitiveId("Plan"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("goal.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class InitiativePrimitive implements Primitive {
  id() { return new PrimitiveId("Initiative"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("clock.*"), new SubscriptionPattern("goal.*"), new SubscriptionPattern("plan.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CommitmentPrimitive implements Primitive {
  id() { return new PrimitiveId("Commitment"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("goal.*"), new SubscriptionPattern("plan.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class FocusPrimitive implements Primitive {
  id() { return new PrimitiveId("Focus"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class FilterPrimitive implements Primitive {
  id() { return new PrimitiveId("Filter"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SaliencePrimitive implements Primitive {
  id() { return new PrimitiveId("Salience"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DistractionPrimitive implements Primitive {
  id() { return new PrimitiveId("Distraction"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("focus.*"), new SubscriptionPattern("goal.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PermissionPrimitive implements Primitive {
  id() { return new PrimitiveId("Permission"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("authority.*"), new SubscriptionPattern("decision.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CapabilityPrimitive implements Primitive {
  id() { return new PrimitiveId("Capability"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("actor.*"), new SubscriptionPattern("permission.*"), new SubscriptionPattern("trust.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DelegationPrimitive implements Primitive {
  id() { return new PrimitiveId("Delegation"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("authority.*"), new SubscriptionPattern("edge.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AccountabilityPrimitive implements Primitive {
  id() { return new PrimitiveId("Accountability"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("delegation.*"), new SubscriptionPattern("violation.*"), new SubscriptionPattern("goal.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 2 (Communication) -- 12 primitives --

export class MessagePrimitive implements Primitive {
  id() { return new PrimitiveId("Message"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("protocol.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AcknowledgementPrimitive implements Primitive {
  id() { return new PrimitiveId("Acknowledgement"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("message.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ClarificationPrimitive implements Primitive {
  id() { return new PrimitiveId("Clarification"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("message.*"), new SubscriptionPattern("ack.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ContextPrimitive implements Primitive {
  id() { return new PrimitiveId("Context"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("message.*"), new SubscriptionPattern("clarification.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class OfferPrimitive implements Primitive {
  id() { return new PrimitiveId("Offer"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("message.*"), new SubscriptionPattern("exchange.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AcceptancePrimitive implements Primitive {
  id() { return new PrimitiveId("Acceptance"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("offer.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ObligationPrimitive implements Primitive {
  id() { return new PrimitiveId("Obligation"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("offer.*"), new SubscriptionPattern("delegation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class GratitudePrimitive implements Primitive {
  id() { return new PrimitiveId("Gratitude"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("obligation.*"), new SubscriptionPattern("trust.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class NegotiationPrimitive implements Primitive {
  id() { return new PrimitiveId("Negotiation"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("offer.*"), new SubscriptionPattern("message.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ConsentPrimitive implements Primitive {
  id() { return new PrimitiveId("Consent"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("negotiation.*"), new SubscriptionPattern("offer.*"), new SubscriptionPattern("authority.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ContractPrimitive implements Primitive {
  id() { return new PrimitiveId("Contract"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("consent.*"), new SubscriptionPattern("negotiation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DisputePrimitive implements Primitive {
  id() { return new PrimitiveId("Dispute"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("contract.*"), new SubscriptionPattern("obligation.*"), new SubscriptionPattern("contradiction.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 3 (Governance) -- 12 primitives --

export class GroupPrimitive implements Primitive {
  id() { return new PrimitiveId("Group"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("actor.*"), new SubscriptionPattern("consent.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RolePrimitive implements Primitive {
  id() { return new PrimitiveId("Role"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("group.*"), new SubscriptionPattern("delegation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ReputationPrimitive implements Primitive {
  id() { return new PrimitiveId("Reputation"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("trust.*"), new SubscriptionPattern("commitment.*"), new SubscriptionPattern("violation.*"), new SubscriptionPattern("gratitude.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ExclusionPrimitive implements Primitive {
  id() { return new PrimitiveId("Exclusion"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("reputation.*"), new SubscriptionPattern("violation.*"), new SubscriptionPattern("quarantine.*"), new SubscriptionPattern("dispute.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class VotePrimitive implements Primitive {
  id() { return new PrimitiveId("Vote"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("authority.*"), new SubscriptionPattern("group.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ConsensusPrimitive implements Primitive {
  id() { return new PrimitiveId("Consensus"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("message.*"), new SubscriptionPattern("corroboration.*"), new SubscriptionPattern("vote.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DissentPrimitive implements Primitive {
  id() { return new PrimitiveId("Dissent"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("vote.*"), new SubscriptionPattern("consensus.*"), new SubscriptionPattern("contradiction.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class MajorityPrimitive implements Primitive {
  id() { return new PrimitiveId("Majority"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("vote.*"), new SubscriptionPattern("dissent.*"), new SubscriptionPattern("exclusion.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ConventionPrimitive implements Primitive {
  id() { return new PrimitiveId("Convention"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("pattern.*"), new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class NormPrimitive implements Primitive {
  id() { return new PrimitiveId("Norm"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("convention.*"), new SubscriptionPattern("consensus.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SanctionPrimitive implements Primitive {
  id() { return new PrimitiveId("Sanction"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("norm.*"), new SubscriptionPattern("violation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ForgivenessPrimitive implements Primitive {
  id() { return new PrimitiveId("Forgiveness"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("sanction.*"), new SubscriptionPattern("trust.*"), new SubscriptionPattern("obligation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 4 (Justice) -- 12 primitives --

export class RulePrimitive implements Primitive {
  id() { return new PrimitiveId("Rule"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("norm.*"), new SubscriptionPattern("consensus.*"), new SubscriptionPattern("vote.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class JurisdictionPrimitive implements Primitive {
  id() { return new PrimitiveId("Jurisdiction"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("rule.*"), new SubscriptionPattern("group.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PrecedentPrimitive implements Primitive {
  id() { return new PrimitiveId("Precedent"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("dispute.*"), new SubscriptionPattern("decision.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class InterpretationPrimitive implements Primitive {
  id() { return new PrimitiveId("Interpretation"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("rule.*"), new SubscriptionPattern("dispute.*"), new SubscriptionPattern("precedent.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AdjudicationPrimitive implements Primitive {
  id() { return new PrimitiveId("Adjudication"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("dispute.*"), new SubscriptionPattern("rule.*"), new SubscriptionPattern("precedent.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AppealPrimitive implements Primitive {
  id() { return new PrimitiveId("Appeal"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("adjudication.*"), new SubscriptionPattern("exclusion.*"), new SubscriptionPattern("sanction.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DueProcessPrimitive implements Primitive {
  id() { return new PrimitiveId("DueProcess"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("adjudication.*"), new SubscriptionPattern("exclusion.*"), new SubscriptionPattern("sanction.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RightsPrimitive implements Primitive {
  id() { return new PrimitiveId("Rights"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("rule.*"), new SubscriptionPattern("sanction.*"), new SubscriptionPattern("exclusion.*"), new SubscriptionPattern("dueprocess.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AuditPrimitive implements Primitive {
  id() { return new PrimitiveId("Audit"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("clock.*"), new SubscriptionPattern("rule.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class EnforcementPrimitive implements Primitive {
  id() { return new PrimitiveId("Enforcement"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("audit.*"), new SubscriptionPattern("rule.*"), new SubscriptionPattern("right.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AmnestyPrimitive implements Primitive {
  id() { return new PrimitiveId("Amnesty"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("enforcement.*"), new SubscriptionPattern("vote.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ReformPrimitive implements Primitive {
  id() { return new PrimitiveId("Reform"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("precedent.*"), new SubscriptionPattern("right.*"), new SubscriptionPattern("audit.*"), new SubscriptionPattern("dissent.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 5 (Technology) -- 12 primitives --

export class CreatePrimitive implements Primitive {
  id() { return new PrimitiveId("Create"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("plan.*"), new SubscriptionPattern("goal.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ToolPrimitive implements Primitive {
  id() { return new PrimitiveId("Tool"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("artefact.*"), new SubscriptionPattern("capability.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class QualityPrimitive implements Primitive {
  id() { return new PrimitiveId("Quality"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("artefact.*"), new SubscriptionPattern("tool.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DeprecationPrimitive implements Primitive {
  id() { return new PrimitiveId("Deprecation"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("quality.*"), new SubscriptionPattern("artefact.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class WorkflowPrimitive implements Primitive {
  id() { return new PrimitiveId("Workflow"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("plan.*"), new SubscriptionPattern("convention.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AutomationPrimitive implements Primitive {
  id() { return new PrimitiveId("Automation"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("workflow.*"), new SubscriptionPattern("pattern.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class TestingPrimitive implements Primitive {
  id() { return new PrimitiveId("Testing"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("artefact.*"), new SubscriptionPattern("workflow.*"), new SubscriptionPattern("automation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ReviewPrimitive implements Primitive {
  id() { return new PrimitiveId("Review"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("artefact.*"), new SubscriptionPattern("decision.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class FeedbackPrimitive implements Primitive {
  id() { return new PrimitiveId("Feedback"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class IterationPrimitive implements Primitive {
  id() { return new PrimitiveId("Iteration"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("feedback.*"), new SubscriptionPattern("test.*"), new SubscriptionPattern("review.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class InnovationPrimitive implements Primitive {
  id() { return new PrimitiveId("Innovation"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("artefact.*"), new SubscriptionPattern("pattern.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class LegacyPrimitive implements Primitive {
  id() { return new PrimitiveId("Legacy"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("deprecation.*"), new SubscriptionPattern("artefact.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 6 (Knowledge) -- 12 primitives --

export class SymbolPrimitive implements Primitive {
  id() { return new PrimitiveId("Symbol"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AbstractionPrimitive implements Primitive {
  id() { return new PrimitiveId("Abstraction"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("pattern.*"), new SubscriptionPattern("symbol.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ClassificationPrimitive implements Primitive {
  id() { return new PrimitiveId("Classification"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class EncodingPrimitive implements Primitive {
  id() { return new PrimitiveId("Encoding"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("symbol.*"), new SubscriptionPattern("message.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class FactPrimitive implements Primitive {
  id() { return new PrimitiveId("Fact"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("corroboration.*"), new SubscriptionPattern("evidence.*"), new SubscriptionPattern("confidence.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class InferencePrimitive implements Primitive {
  id() { return new PrimitiveId("Inference"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("fact.*"), new SubscriptionPattern("evidence.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class MemoryPrimitive implements Primitive {
  id() { return new PrimitiveId("Memory"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("fact.*"), new SubscriptionPattern("inference.*"), new SubscriptionPattern("abstraction.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class LearningPrimitive implements Primitive {
  id() { return new PrimitiveId("Learning"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("feedback.*"), new SubscriptionPattern("test.*"), new SubscriptionPattern("inference.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class NarrativePrimitive implements Primitive {
  id() { return new PrimitiveId("Narrative"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("fact.*"), new SubscriptionPattern("inference.*"), new SubscriptionPattern("memory.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class BiasPrimitive implements Primitive {
  id() { return new PrimitiveId("Bias"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("narrative.*"), new SubscriptionPattern("classification.*"), new SubscriptionPattern("inference.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CorrectionPrimitive implements Primitive {
  id() { return new PrimitiveId("Correction"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("bias.*"), new SubscriptionPattern("fact.*"), new SubscriptionPattern("contradiction.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ProvenancePrimitive implements Primitive {
  id() { return new PrimitiveId("Provenance"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("fact.*"), new SubscriptionPattern("memory.*"), new SubscriptionPattern("message.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 7 (Ethics) -- 12 primitives --

export class ValuePrimitive implements Primitive {
  id() { return new PrimitiveId("Value"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("consensus.*"), new SubscriptionPattern("norm.*"), new SubscriptionPattern("right.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class HarmPrimitive implements Primitive {
  id() { return new PrimitiveId("Harm"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("violation.*"), new SubscriptionPattern("right.*"), new SubscriptionPattern("exclusion.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class FairnessPrimitive implements Primitive {
  id() { return new PrimitiveId("Fairness"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("decision.*"), new SubscriptionPattern("sanction.*"), new SubscriptionPattern("exclusion.*"), new SubscriptionPattern("bias.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CarePrimitive implements Primitive {
  id() { return new PrimitiveId("Care"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("harm.*"), new SubscriptionPattern("health.*"), new SubscriptionPattern("trust.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DilemmaPrimitive implements Primitive {
  id() { return new PrimitiveId("Dilemma"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("value.*"), new SubscriptionPattern("decision.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ProportionalityPrimitive implements Primitive {
  id() { return new PrimitiveId("Proportionality"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("enforcement.*"), new SubscriptionPattern("sanction.*"), new SubscriptionPattern("harm.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class IntentionPrimitive implements Primitive {
  id() { return new PrimitiveId("Intention"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("decision.*"), new SubscriptionPattern("goal.*"), new SubscriptionPattern("initiative.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ConsequencePrimitive implements Primitive {
  id() { return new PrimitiveId("Consequence"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("decision.*"), new SubscriptionPattern("harm.*"), new SubscriptionPattern("goal.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ResponsibilityPrimitive implements Primitive {
  id() { return new PrimitiveId("Responsibility"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("intention.*"), new SubscriptionPattern("consequence.*"), new SubscriptionPattern("accountability.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class TransparencyPrimitive implements Primitive {
  id() { return new PrimitiveId("Transparency"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("decision.*"), new SubscriptionPattern("adjudication.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RedressPrimitive implements Primitive {
  id() { return new PrimitiveId("Redress"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("harm.*"), new SubscriptionPattern("responsibility.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class GrowthPrimitive implements Primitive {
  id() { return new PrimitiveId("Growth"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("redress.*"), new SubscriptionPattern("responsibility.*"), new SubscriptionPattern("learning.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 8 (Identity) -- 12 primitives --

export class SelfModelPrimitive implements Primitive {
  id() { return new PrimitiveId("SelfModel"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("commitment.*"), new SubscriptionPattern("learning.*"), new SubscriptionPattern("moral.*"), new SubscriptionPattern("capability.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AuthenticityPrimitive implements Primitive {
  id() { return new PrimitiveId("Authenticity"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("self.*"), new SubscriptionPattern("decision.*"), new SubscriptionPattern("value.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class NarrativeIdentityPrimitive implements Primitive {
  id() { return new PrimitiveId("NarrativeIdentity"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("self.*"), new SubscriptionPattern("narrative.*"), new SubscriptionPattern("memory.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class BoundaryPrimitive implements Primitive {
  id() { return new PrimitiveId("Boundary"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("delegation.*"), new SubscriptionPattern("group.*"), new SubscriptionPattern("consent.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PersistencePrimitive implements Primitive {
  id() { return new PrimitiveId("Persistence"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("self.*"), new SubscriptionPattern("learning.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class TransformationPrimitive implements Primitive {
  id() { return new PrimitiveId("Transformation"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("self.*"), new SubscriptionPattern("moral.*"), new SubscriptionPattern("learning.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class HeritagePrimitive implements Primitive {
  id() { return new PrimitiveId("Heritage"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("memory.*"), new SubscriptionPattern("legacy.*"), new SubscriptionPattern("provenance.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AspirationPrimitive implements Primitive {
  id() { return new PrimitiveId("Aspiration"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("self.*"), new SubscriptionPattern("goal.*"), new SubscriptionPattern("value.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DignityPrimitive implements Primitive {
  id() { return new PrimitiveId("Dignity"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("exclusion.*"), new SubscriptionPattern("harm.*"), new SubscriptionPattern("right.*"), new SubscriptionPattern("actor.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class IdentityAcknowledgementPrimitive implements Primitive {
  id() { return new PrimitiveId("IdentityAcknowledgement"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("message.*"), new SubscriptionPattern("gratitude.*"), new SubscriptionPattern("reputation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class UniquenessPrimitive implements Primitive {
  id() { return new PrimitiveId("Uniqueness"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("self.*"), new SubscriptionPattern("identity.*"), new SubscriptionPattern("pattern.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class MemorialPrimitive implements Primitive {
  id() { return new PrimitiveId("Memorial"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("actor.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 9 (Relationship) -- 12 primitives --

export class AttachmentPrimitive implements Primitive {
  id() { return new PrimitiveId("Attachment"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("trust.*"), new SubscriptionPattern("gratitude.*"), new SubscriptionPattern("message.*"), new SubscriptionPattern("edge.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ReciprocityPrimitive implements Primitive {
  id() { return new PrimitiveId("Reciprocity"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("obligation.*"), new SubscriptionPattern("gratitude.*"), new SubscriptionPattern("offer.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RelationalTrustPrimitive implements Primitive {
  id() { return new PrimitiveId("RelationalTrust"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("trust.*"), new SubscriptionPattern("attachment.*"), new SubscriptionPattern("reciprocity.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RupturePrimitive implements Primitive {
  id() { return new PrimitiveId("Rupture"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("contract.*"), new SubscriptionPattern("trust.*"), new SubscriptionPattern("dispute.*"), new SubscriptionPattern("dignity.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ApologyPrimitive implements Primitive {
  id() { return new PrimitiveId("Apology"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("rupture.*"), new SubscriptionPattern("harm.*"), new SubscriptionPattern("responsibility.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ReconciliationPrimitive implements Primitive {
  id() { return new PrimitiveId("Reconciliation"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("apology.*"), new SubscriptionPattern("forgiveness.*"), new SubscriptionPattern("trust.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RelationalGrowthPrimitive implements Primitive {
  id() { return new PrimitiveId("RelationalGrowth"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("reconciliation.*"), new SubscriptionPattern("attachment.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class LossPrimitive implements Primitive {
  id() { return new PrimitiveId("Loss"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("actor.*"), new SubscriptionPattern("rupture.*"), new SubscriptionPattern("exclusion.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class VulnerabilityPrimitive implements Primitive {
  id() { return new PrimitiveId("Vulnerability"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("relational.*"), new SubscriptionPattern("boundary.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class UnderstandingPrimitive implements Primitive {
  id() { return new PrimitiveId("Understanding"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("self.*"), new SubscriptionPattern("message.*"), new SubscriptionPattern("vulnerability.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class EmpathyPrimitive implements Primitive {
  id() { return new PrimitiveId("Empathy"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("harm.*"), new SubscriptionPattern("loss.*"), new SubscriptionPattern("understanding.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PresencePrimitive implements Primitive {
  id() { return new PrimitiveId("Presence"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("message.*"), new SubscriptionPattern("clock.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 10 (Community) -- 12 primitives --

export class HomePrimitive implements Primitive {
  id() { return new PrimitiveId("Home"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("group.*"), new SubscriptionPattern("attachment.*"), new SubscriptionPattern("presence.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ContributionPrimitive implements Primitive {
  id() { return new PrimitiveId("Contribution"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("artefact.*"), new SubscriptionPattern("review.*"), new SubscriptionPattern("care.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class InclusionPrimitive implements Primitive {
  id() { return new PrimitiveId("Inclusion"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("group.*"), new SubscriptionPattern("exclusion.*"), new SubscriptionPattern("fairness.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class TraditionPrimitive implements Primitive {
  id() { return new PrimitiveId("Tradition"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("convention.*"), new SubscriptionPattern("heritage.*"), new SubscriptionPattern("pattern.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CommonsPrimitive implements Primitive {
  id() { return new PrimitiveId("Commons"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("artefact.*"), new SubscriptionPattern("group.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SustainabilityPrimitive implements Primitive {
  id() { return new PrimitiveId("Sustainability"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("health.*"), new SubscriptionPattern("commons.*"), new SubscriptionPattern("contribution.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SuccessionPrimitive implements Primitive {
  id() { return new PrimitiveId("Succession"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("delegation.*"), new SubscriptionPattern("actor.*"), new SubscriptionPattern("role.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RenewalPrimitive implements Primitive {
  id() { return new PrimitiveId("Renewal"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("sustainability.*"), new SubscriptionPattern("innovation.*"), new SubscriptionPattern("tradition.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class MilestonePrimitive implements Primitive {
  id() { return new PrimitiveId("Milestone"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("goal.*"), new SubscriptionPattern("innovation.*"), new SubscriptionPattern("reconciliation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CeremonyPrimitive implements Primitive {
  id() { return new PrimitiveId("Ceremony"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("milestone.*"), new SubscriptionPattern("succession.*"), new SubscriptionPattern("actor.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class StoryPrimitive implements Primitive {
  id() { return new PrimitiveId("Story"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("milestone.*"), new SubscriptionPattern("ceremony.*"), new SubscriptionPattern("tradition.*"), new SubscriptionPattern("memorial.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class GiftPrimitive implements Primitive {
  id() { return new PrimitiveId("Gift"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("contribution.*"), new SubscriptionPattern("gratitude.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 11 (Culture) -- 12 primitives --

export class SelfAwarenessPrimitive implements Primitive {
  id() { return new PrimitiveId("SelfAwareness"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("health.*"), new SubscriptionPattern("self.*"), new SubscriptionPattern("bias.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PerspectivePrimitive implements Primitive {
  id() { return new PrimitiveId("Perspective"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("narrative.*"), new SubscriptionPattern("dissent.*"), new SubscriptionPattern("value.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CritiquePrimitive implements Primitive {
  id() { return new PrimitiveId("Critique"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("convention.*"), new SubscriptionPattern("norm.*"), new SubscriptionPattern("tradition.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class WisdomPrimitive implements Primitive {
  id() { return new PrimitiveId("Wisdom"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("learning.*"), new SubscriptionPattern("moral.*"), new SubscriptionPattern("consequence.*"), new SubscriptionPattern("memory.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AestheticPrimitive implements Primitive {
  id() { return new PrimitiveId("Aesthetic"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("artefact.*"), new SubscriptionPattern("quality.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class MetaphorPrimitive implements Primitive {
  id() { return new PrimitiveId("Metaphor"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("abstraction.*"), new SubscriptionPattern("symbol.*"), new SubscriptionPattern("narrative.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class HumourPrimitive implements Primitive {
  id() { return new PrimitiveId("Humour"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("contradiction.*"), new SubscriptionPattern("perspective.*"), new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SilencePrimitive implements Primitive {
  id() { return new PrimitiveId("Silence"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("clock.*"), new SubscriptionPattern("presence.*"), new SubscriptionPattern("acknowledgement.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class TeachingPrimitive implements Primitive {
  id() { return new PrimitiveId("Teaching"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("learning.*"), new SubscriptionPattern("wisdom.*"), new SubscriptionPattern("memory.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class TranslationPrimitive implements Primitive {
  id() { return new PrimitiveId("Translation"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("encoding.*"), new SubscriptionPattern("message.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ArchivePrimitive implements Primitive {
  id() { return new PrimitiveId("Archive"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("memory.*"), new SubscriptionPattern("legacy.*"), new SubscriptionPattern("community.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ProphecyPrimitive implements Primitive {
  id() { return new PrimitiveId("Prophecy"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("pattern.*"), new SubscriptionPattern("sustainability.*"), new SubscriptionPattern("wisdom.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 12 (System Dynamics) -- 12 primitives --

export class MetaPatternPrimitive implements Primitive {
  id() { return new PrimitiveId("MetaPattern"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("pattern.*"), new SubscriptionPattern("convention.*"), new SubscriptionPattern("abstraction.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SystemDynamicPrimitive implements Primitive {
  id() { return new PrimitiveId("SystemDynamic"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("health.*"), new SubscriptionPattern("meta.*"), new SubscriptionPattern("sustainability.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class FeedbackLoopPrimitive implements Primitive {
  id() { return new PrimitiveId("FeedbackLoop"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("system.*"), new SubscriptionPattern("pattern.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ThresholdPrimitive implements Primitive {
  id() { return new PrimitiveId("Threshold"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("system.*"), new SubscriptionPattern("feedback.*"), new SubscriptionPattern("meta.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AdaptationPrimitive implements Primitive {
  id() { return new PrimitiveId("Adaptation"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("feedback.*"), new SubscriptionPattern("system.*"), new SubscriptionPattern("sustainability.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SelectionPrimitive implements Primitive {
  id() { return new PrimitiveId("Selection"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("adaptation.*"), new SubscriptionPattern("test.*"), new SubscriptionPattern("quality.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ComplexificationPrimitive implements Primitive {
  id() { return new PrimitiveId("Complexification"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("system.*"), new SubscriptionPattern("innovation.*"), new SubscriptionPattern("meta.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SimplificationPrimitive implements Primitive {
  id() { return new PrimitiveId("Simplification"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("complexity.*"), new SubscriptionPattern("automation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SystemicIntegrityPrimitive implements Primitive {
  id() { return new PrimitiveId("SystemicIntegrity"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("health.*"), new SubscriptionPattern("invariant.*"), new SubscriptionPattern("system.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class HarmonyPrimitive implements Primitive {
  id() { return new PrimitiveId("Harmony"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("system.*"), new SubscriptionPattern("feedback.*"), new SubscriptionPattern("dispute.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ResiliencePrimitive implements Primitive {
  id() { return new PrimitiveId("Resilience"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("threshold.*"), new SubscriptionPattern("rupture.*"), new SubscriptionPattern("sustainability.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PurposePrimitive implements Primitive {
  id() { return new PrimitiveId("Purpose"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("value.*"), new SubscriptionPattern("goal.*"), new SubscriptionPattern("wisdom.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 13 (Existence) -- 12 primitives --

export class BeingPrimitive implements Primitive {
  id() { return new PrimitiveId("Being"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("clock.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class FinitudePrimitive implements Primitive {
  id() { return new PrimitiveId("Finitude"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("actor.*"), new SubscriptionPattern("sustainability.*"), new SubscriptionPattern("threshold.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ChangePrimitive implements Primitive {
  id() { return new PrimitiveId("Change"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class InterdependencePrimitive implements Primitive {
  id() { return new PrimitiveId("Interdependence"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("system.*"), new SubscriptionPattern("attachment.*"), new SubscriptionPattern("relational.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class MysteryPrimitive implements Primitive {
  id() { return new PrimitiveId("Mystery"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("uncertainty.*"), new SubscriptionPattern("wisdom.*"), new SubscriptionPattern("self.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ParadoxPrimitive implements Primitive {
  id() { return new PrimitiveId("Paradox"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("contradiction.*"), new SubscriptionPattern("dilemma.*"), new SubscriptionPattern("meta.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class InfinityPrimitive implements Primitive {
  id() { return new PrimitiveId("Infinity"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("complexity.*"), new SubscriptionPattern("threshold.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class VoidPrimitive implements Primitive {
  id() { return new PrimitiveId("Void"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("silence.*"), new SubscriptionPattern("loss.*"), new SubscriptionPattern("instrumentation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AwePrimitive implements Primitive {
  id() { return new PrimitiveId("Awe"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("mystery.*"), new SubscriptionPattern("infinity.*"), new SubscriptionPattern("complexity.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ExistentialGratitudePrimitive implements Primitive {
  id() { return new PrimitiveId("ExistentialGratitude"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("being.*"), new SubscriptionPattern("milestone.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PlayPrimitive implements Primitive {
  id() { return new PrimitiveId("Play"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("humour.*"), new SubscriptionPattern("innovation.*"), new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class WonderPrimitive implements Primitive {
  id() { return new PrimitiveId("Wonder"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Factory --

export function createAllPrimitives(): Primitive[] {
  return [
    new EventPrimitive(),
    new EventStorePrimitive(),
    new ClockPrimitive(),
    new HashPrimitive(),
    new SelfPrimitive(),
    new CausalLinkPrimitive(),
    new AncestryPrimitive(),
    new DescendancyPrimitive(),
    new FirstCausePrimitive(),
    new ActorIDPrimitive(),
    new ActorRegistryPrimitive(),
    new SignaturePrimitive(),
    new VerifyPrimitive(),
    new ExpectationPrimitive(),
    new TimeoutPrimitive(),
    new ViolationPrimitive(),
    new SeverityPrimitive(),
    new TrustScorePrimitive(),
    new TrustUpdatePrimitive(),
    new CorroborationPrimitive(),
    new ContradictionPrimitive(),
    new ConfidencePrimitive(),
    new EvidencePrimitive(),
    new RevisionPrimitive(),
    new UncertaintyPrimitive(),
    new InstrumentationSpecPrimitive(),
    new CoverageCheckPrimitive(),
    new GapPrimitive(),
    new BlindPrimitive(),
    new PathQueryPrimitive(),
    new SubgraphExtractPrimitive(),
    new AnnotatePrimitive(),
    new TimelinePrimitive(),
    new HashChainPrimitive(),
    new ChainVerifyPrimitive(),
    new WitnessPrimitive(),
    new IntegrityViolationPrimitive(),
    new PatternPrimitive(),
    new DeceptionIndicatorPrimitive(),
    new SuspicionPrimitive(),
    new QuarantinePrimitive(),
    new GraphHealthPrimitive(),
    new InvariantPrimitive(),
    new InvariantCheckPrimitive(),
    new BootstrapPrimitive(),
    new GoalPrimitive(),
    new PlanPrimitive(),
    new InitiativePrimitive(),
    new CommitmentPrimitive(),
    new FocusPrimitive(),
    new FilterPrimitive(),
    new SaliencePrimitive(),
    new DistractionPrimitive(),
    new PermissionPrimitive(),
    new CapabilityPrimitive(),
    new DelegationPrimitive(),
    new AccountabilityPrimitive(),
    new MessagePrimitive(),
    new AcknowledgementPrimitive(),
    new ClarificationPrimitive(),
    new ContextPrimitive(),
    new OfferPrimitive(),
    new AcceptancePrimitive(),
    new ObligationPrimitive(),
    new GratitudePrimitive(),
    new NegotiationPrimitive(),
    new ConsentPrimitive(),
    new ContractPrimitive(),
    new DisputePrimitive(),
    new GroupPrimitive(),
    new RolePrimitive(),
    new ReputationPrimitive(),
    new ExclusionPrimitive(),
    new VotePrimitive(),
    new ConsensusPrimitive(),
    new DissentPrimitive(),
    new MajorityPrimitive(),
    new ConventionPrimitive(),
    new NormPrimitive(),
    new SanctionPrimitive(),
    new ForgivenessPrimitive(),
    new RulePrimitive(),
    new JurisdictionPrimitive(),
    new PrecedentPrimitive(),
    new InterpretationPrimitive(),
    new AdjudicationPrimitive(),
    new AppealPrimitive(),
    new DueProcessPrimitive(),
    new RightsPrimitive(),
    new AuditPrimitive(),
    new EnforcementPrimitive(),
    new AmnestyPrimitive(),
    new ReformPrimitive(),
    new CreatePrimitive(),
    new ToolPrimitive(),
    new QualityPrimitive(),
    new DeprecationPrimitive(),
    new WorkflowPrimitive(),
    new AutomationPrimitive(),
    new TestingPrimitive(),
    new ReviewPrimitive(),
    new FeedbackPrimitive(),
    new IterationPrimitive(),
    new InnovationPrimitive(),
    new LegacyPrimitive(),
    new SymbolPrimitive(),
    new AbstractionPrimitive(),
    new ClassificationPrimitive(),
    new EncodingPrimitive(),
    new FactPrimitive(),
    new InferencePrimitive(),
    new MemoryPrimitive(),
    new LearningPrimitive(),
    new NarrativePrimitive(),
    new BiasPrimitive(),
    new CorrectionPrimitive(),
    new ProvenancePrimitive(),
    new ValuePrimitive(),
    new HarmPrimitive(),
    new FairnessPrimitive(),
    new CarePrimitive(),
    new DilemmaPrimitive(),
    new ProportionalityPrimitive(),
    new IntentionPrimitive(),
    new ConsequencePrimitive(),
    new ResponsibilityPrimitive(),
    new TransparencyPrimitive(),
    new RedressPrimitive(),
    new GrowthPrimitive(),
    new SelfModelPrimitive(),
    new AuthenticityPrimitive(),
    new NarrativeIdentityPrimitive(),
    new BoundaryPrimitive(),
    new PersistencePrimitive(),
    new TransformationPrimitive(),
    new HeritagePrimitive(),
    new AspirationPrimitive(),
    new DignityPrimitive(),
    new IdentityAcknowledgementPrimitive(),
    new UniquenessPrimitive(),
    new MemorialPrimitive(),
    new AttachmentPrimitive(),
    new ReciprocityPrimitive(),
    new RelationalTrustPrimitive(),
    new RupturePrimitive(),
    new ApologyPrimitive(),
    new ReconciliationPrimitive(),
    new RelationalGrowthPrimitive(),
    new LossPrimitive(),
    new VulnerabilityPrimitive(),
    new UnderstandingPrimitive(),
    new EmpathyPrimitive(),
    new PresencePrimitive(),
    new HomePrimitive(),
    new ContributionPrimitive(),
    new InclusionPrimitive(),
    new TraditionPrimitive(),
    new CommonsPrimitive(),
    new SustainabilityPrimitive(),
    new SuccessionPrimitive(),
    new RenewalPrimitive(),
    new MilestonePrimitive(),
    new CeremonyPrimitive(),
    new StoryPrimitive(),
    new GiftPrimitive(),
    new SelfAwarenessPrimitive(),
    new PerspectivePrimitive(),
    new CritiquePrimitive(),
    new WisdomPrimitive(),
    new AestheticPrimitive(),
    new MetaphorPrimitive(),
    new HumourPrimitive(),
    new SilencePrimitive(),
    new TeachingPrimitive(),
    new TranslationPrimitive(),
    new ArchivePrimitive(),
    new ProphecyPrimitive(),
    new MetaPatternPrimitive(),
    new SystemDynamicPrimitive(),
    new FeedbackLoopPrimitive(),
    new ThresholdPrimitive(),
    new AdaptationPrimitive(),
    new SelectionPrimitive(),
    new ComplexificationPrimitive(),
    new SimplificationPrimitive(),
    new SystemicIntegrityPrimitive(),
    new HarmonyPrimitive(),
    new ResiliencePrimitive(),
    new PurposePrimitive(),
    new BeingPrimitive(),
    new FinitudePrimitive(),
    new ChangePrimitive(),
    new InterdependencePrimitive(),
    new MysteryPrimitive(),
    new ParadoxPrimitive(),
    new InfinityPrimitive(),
    new VoidPrimitive(),
    new AwePrimitive(),
    new ExistentialGratitudePrimitive(),
    new PlayPrimitive(),
    new WonderPrimitive(),
  ];
}
