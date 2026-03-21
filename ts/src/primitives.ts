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


// -- Layer 1 (Agency) -- 12 primitives --
// Volition: Value, Intent, Choice, Risk
// Action: Act, Consequence, Capacity, Resource
// Communication: Signal, Reception, Acknowledgment, Commitment

export class ValuePrimitive implements Primitive {
  id() { return new PrimitiveId("Value"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("decision.*"), new SubscriptionPattern("actor.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class IntentPrimitive implements Primitive {
  id() { return new PrimitiveId("Intent"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("value.*"), new SubscriptionPattern("expectation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ChoicePrimitive implements Primitive {
  id() { return new PrimitiveId("Choice"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("intent.*"), new SubscriptionPattern("value.*"), new SubscriptionPattern("confidence.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RiskPrimitive implements Primitive {
  id() { return new PrimitiveId("Risk"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("intent.*"), new SubscriptionPattern("uncertainty.*"), new SubscriptionPattern("value.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ActPrimitive implements Primitive {
  id() { return new PrimitiveId("Act"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("choice.*"), new SubscriptionPattern("intent.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ConsequencePrimitive implements Primitive {
  id() { return new PrimitiveId("Consequence"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("act.*"), new SubscriptionPattern("violation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CapacityPrimitive implements Primitive {
  id() { return new PrimitiveId("Capacity"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("actor.*"), new SubscriptionPattern("resource.*"), new SubscriptionPattern("trust.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ResourcePrimitive implements Primitive {
  id() { return new PrimitiveId("Resource"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("act.*"), new SubscriptionPattern("budget.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SignalPrimitive implements Primitive {
  id() { return new PrimitiveId("Signal"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("act.*"), new SubscriptionPattern("actor.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ReceptionPrimitive implements Primitive {
  id() { return new PrimitiveId("Reception"); }
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

export class AcknowledgmentPrimitive implements Primitive {
  id() { return new PrimitiveId("Acknowledgment"); }
  layer() { return new Layer(1); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("signal.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("signal.*"), new SubscriptionPattern("agreement.*"), new SubscriptionPattern("intent.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 2 (Exchange) -- 12 primitives --

export class TermPrimitive implements Primitive {
  id() { return new PrimitiveId("Term"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("agreement.*"), new SubscriptionPattern("obligation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ProtocolPrimitive implements Primitive {
  id() { return new PrimitiveId("Protocol"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("term.*"), new SubscriptionPattern("exchange.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("intent.*"), new SubscriptionPattern("value.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("offer.*"), new SubscriptionPattern("choice.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AgreementPrimitive implements Primitive {
  id() { return new PrimitiveId("Agreement"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("offer.*"), new SubscriptionPattern("acceptance.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("agreement.*"), new SubscriptionPattern("commitment.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class FulfillmentPrimitive implements Primitive {
  id() { return new PrimitiveId("Fulfillment"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("obligation.*"), new SubscriptionPattern("act.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class BreachPrimitive implements Primitive {
  id() { return new PrimitiveId("Breach"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("obligation.*"), new SubscriptionPattern("violation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ExchangePrimitive implements Primitive {
  id() { return new PrimitiveId("Exchange"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("fulfillment.*"), new SubscriptionPattern("resource.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AccountabilityPrimitive implements Primitive {
  id() { return new PrimitiveId("Accountability"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("act.*"), new SubscriptionPattern("consequence.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DebtPrimitive implements Primitive {
  id() { return new PrimitiveId("Debt"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("obligation.*"), new SubscriptionPattern("breach.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ReciprocityPrimitive implements Primitive {
  id() { return new PrimitiveId("Reciprocity"); }
  layer() { return new Layer(2); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("exchange.*"), new SubscriptionPattern("obligation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 3 (Society) -- 12 primitives --

export class GroupPrimitive implements Primitive {
  id() { return new PrimitiveId("Group"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("actor.*"), new SubscriptionPattern("membership.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class MembershipPrimitive implements Primitive {
  id() { return new PrimitiveId("Membership"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("group.*"), new SubscriptionPattern("actor.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("group.*"), new SubscriptionPattern("membership.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ConsentPrimitive implements Primitive {
  id() { return new PrimitiveId("Consent"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("choice.*"), new SubscriptionPattern("agreement.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("group.*"), new SubscriptionPattern("consent.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("trust.*"), new SubscriptionPattern("act.*")]; }
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

export class AuthorityPrimitive implements Primitive {
  id() { return new PrimitiveId("Authority"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("role.*"), new SubscriptionPattern("consent.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PropertyPrimitive implements Primitive {
  id() { return new PrimitiveId("Property"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("resource.*"), new SubscriptionPattern("actor.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CommonsPrimitive implements Primitive {
  id() { return new PrimitiveId("Commons"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("group.*"), new SubscriptionPattern("resource.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class GovernancePrimitive implements Primitive {
  id() { return new PrimitiveId("Governance"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("authority.*"), new SubscriptionPattern("norm.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CollectiveActPrimitive implements Primitive {
  id() { return new PrimitiveId("CollectiveAct"); }
  layer() { return new Layer(3); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("group.*"), new SubscriptionPattern("act.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 4 (Legal) -- 12 primitives --

export class LawPrimitive implements Primitive {
  id() { return new PrimitiveId("Law"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("norm.*"), new SubscriptionPattern("governance.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RightPrimitive implements Primitive {
  id() { return new PrimitiveId("Right"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("law.*"), new SubscriptionPattern("actor.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ContractPrimitive implements Primitive {
  id() { return new PrimitiveId("Contract"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("agreement.*"), new SubscriptionPattern("obligation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class LiabilityPrimitive implements Primitive {
  id() { return new PrimitiveId("Liability"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("breach.*"), new SubscriptionPattern("consequence.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("right.*"), new SubscriptionPattern("law.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("dueprocess.*"), new SubscriptionPattern("liability.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RemedyPrimitive implements Primitive {
  id() { return new PrimitiveId("Remedy"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("adjudication.*"), new SubscriptionPattern("breach.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("adjudication.*"), new SubscriptionPattern("law.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("authority.*"), new SubscriptionPattern("law.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SovereigntyPrimitive implements Primitive {
  id() { return new PrimitiveId("Sovereignty"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("jurisdiction.*"), new SubscriptionPattern("authority.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class LegitimacyPrimitive implements Primitive {
  id() { return new PrimitiveId("Legitimacy"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("consent.*"), new SubscriptionPattern("authority.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class TreatyPrimitive implements Primitive {
  id() { return new PrimitiveId("Treaty"); }
  layer() { return new Layer(4); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("agreement.*"), new SubscriptionPattern("sovereignty.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 5 (Technology) -- 12 primitives --

export class MethodPrimitive implements Primitive {
  id() { return new PrimitiveId("Method"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("knowledge.*"), new SubscriptionPattern("technique.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class MeasurementPrimitive implements Primitive {
  id() { return new PrimitiveId("Measurement"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("method.*"), new SubscriptionPattern("standard.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class KnowledgePrimitive implements Primitive {
  id() { return new PrimitiveId("Knowledge"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("measurement.*"), new SubscriptionPattern("model.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ModelPrimitive implements Primitive {
  id() { return new PrimitiveId("Model"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("knowledge.*"), new SubscriptionPattern("abstraction.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("technique.*"), new SubscriptionPattern("capacity.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class TechniquePrimitive implements Primitive {
  id() { return new PrimitiveId("Technique"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("method.*"), new SubscriptionPattern("tool.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class InventionPrimitive implements Primitive {
  id() { return new PrimitiveId("Invention"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("knowledge.*"), new SubscriptionPattern("tool.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AbstractionPrimitive implements Primitive {
  id() { return new PrimitiveId("Abstraction"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("model.*"), new SubscriptionPattern("knowledge.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class InfrastructurePrimitive implements Primitive {
  id() { return new PrimitiveId("Infrastructure"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("tool.*"), new SubscriptionPattern("standard.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class StandardPrimitive implements Primitive {
  id() { return new PrimitiveId("Standard"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("measurement.*"), new SubscriptionPattern("norm.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class EfficiencyPrimitive implements Primitive {
  id() { return new PrimitiveId("Efficiency"); }
  layer() { return new Layer(5); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("measurement.*"), new SubscriptionPattern("resource.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("tool.*"), new SubscriptionPattern("efficiency.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 6 (Information) -- 12 primitives --

export class SymbolPrimitive implements Primitive {
  id() { return new PrimitiveId("Symbol"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("signal.*"), new SubscriptionPattern("encoding.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class LanguagePrimitive implements Primitive {
  id() { return new PrimitiveId("Language"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("symbol.*"), new SubscriptionPattern("channel.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("symbol.*"), new SubscriptionPattern("data.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RecordPrimitive implements Primitive {
  id() { return new PrimitiveId("Record"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("data.*"), new SubscriptionPattern("encoding.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ChannelPrimitive implements Primitive {
  id() { return new PrimitiveId("Channel"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("signal.*"), new SubscriptionPattern("noise.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CopyPrimitive implements Primitive {
  id() { return new PrimitiveId("Copy"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("record.*"), new SubscriptionPattern("redundancy.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class NoisePrimitive implements Primitive {
  id() { return new PrimitiveId("Noise"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("channel.*"), new SubscriptionPattern("entropy.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RedundancyPrimitive implements Primitive {
  id() { return new PrimitiveId("Redundancy"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("copy.*"), new SubscriptionPattern("noise.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DataPrimitive implements Primitive {
  id() { return new PrimitiveId("Data"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("record.*"), new SubscriptionPattern("encoding.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ComputationPrimitive implements Primitive {
  id() { return new PrimitiveId("Computation"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("data.*"), new SubscriptionPattern("algorithm.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AlgorithmPrimitive implements Primitive {
  id() { return new PrimitiveId("Algorithm"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("method.*"), new SubscriptionPattern("computation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class EntropyPrimitive implements Primitive {
  id() { return new PrimitiveId("Entropy"); }
  layer() { return new Layer(6); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("noise.*"), new SubscriptionPattern("data.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 7 (Ethics) -- 12 primitives --

export class MoralStatusPrimitive implements Primitive {
  id() { return new PrimitiveId("MoralStatus"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("actor.*"), new SubscriptionPattern("dignity.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DignityPrimitive implements Primitive {
  id() { return new PrimitiveId("Dignity"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("moralstatus.*"), new SubscriptionPattern("right.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AutonomyPrimitive implements Primitive {
  id() { return new PrimitiveId("Autonomy"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("choice.*"), new SubscriptionPattern("dignity.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class FlourishingPrimitive implements Primitive {
  id() { return new PrimitiveId("Flourishing"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("value.*"), new SubscriptionPattern("care.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DutyPrimitive implements Primitive {
  id() { return new PrimitiveId("Duty"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("obligation.*"), new SubscriptionPattern("moralstatus.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("consequence.*"), new SubscriptionPattern("violation.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("bond.*"), new SubscriptionPattern("duty.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class JusticePrimitive implements Primitive {
  id() { return new PrimitiveId("Justice"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("right.*"), new SubscriptionPattern("harm.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ConsciencePrimitive implements Primitive {
  id() { return new PrimitiveId("Conscience"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("motive.*"), new SubscriptionPattern("duty.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class VirtuePrimitive implements Primitive {
  id() { return new PrimitiveId("Virtue"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("conscience.*"), new SubscriptionPattern("flourishing.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("duty.*"), new SubscriptionPattern("consequence.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class MotivePrimitive implements Primitive {
  id() { return new PrimitiveId("Motive"); }
  layer() { return new Layer(7); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("intent.*"), new SubscriptionPattern("value.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 8 (Identity) -- 12 primitives --

export class NarrativePrimitive implements Primitive {
  id() { return new PrimitiveId("Narrative"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("memory.*"), new SubscriptionPattern("selfconcept.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SelfConceptPrimitive implements Primitive {
  id() { return new PrimitiveId("SelfConcept"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("reflection.*"), new SubscriptionPattern("narrative.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ReflectionPrimitive implements Primitive {
  id() { return new PrimitiveId("Reflection"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("act.*"), new SubscriptionPattern("consequence.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class MemoryPrimitive implements Primitive {
  id() { return new PrimitiveId("Memory"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("record.*"), new SubscriptionPattern("narrative.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PurposePrimitive implements Primitive {
  id() { return new PrimitiveId("Purpose"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("value.*"), new SubscriptionPattern("aspiration.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("purpose.*"), new SubscriptionPattern("growth.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("selfconcept.*"), new SubscriptionPattern("expression.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ExpressionPrimitive implements Primitive {
  id() { return new PrimitiveId("Expression"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("authenticity.*"), new SubscriptionPattern("signal.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class GrowthPrimitive implements Primitive {
  id() { return new PrimitiveId("Growth"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("reflection.*"), new SubscriptionPattern("aspiration.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ContinuityPrimitive implements Primitive {
  id() { return new PrimitiveId("Continuity"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("memory.*"), new SubscriptionPattern("selfconcept.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class IntegrationPrimitive implements Primitive {
  id() { return new PrimitiveId("Integration"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("narrative.*"), new SubscriptionPattern("continuity.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CrisisPrimitive implements Primitive {
  id() { return new PrimitiveId("Crisis"); }
  layer() { return new Layer(8); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("selfconcept.*"), new SubscriptionPattern("rupture.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 9 (Relationship) -- 12 primitives --

export class BondPrimitive implements Primitive {
  id() { return new PrimitiveId("Bond"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("attachment.*"), new SubscriptionPattern("recognition.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AttachmentPrimitive implements Primitive {
  id() { return new PrimitiveId("Attachment"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("bond.*"), new SubscriptionPattern("care.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RecognitionPrimitive implements Primitive {
  id() { return new PrimitiveId("Recognition"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("actor.*"), new SubscriptionPattern("dignity.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class IntimacyPrimitive implements Primitive {
  id() { return new PrimitiveId("Intimacy"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("bond.*"), new SubscriptionPattern("trust.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AttunementPrimitive implements Primitive {
  id() { return new PrimitiveId("Attunement"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("reception.*"), new SubscriptionPattern("care.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("bond.*"), new SubscriptionPattern("harm.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RepairPrimitive implements Primitive {
  id() { return new PrimitiveId("Repair"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("rupture.*"), new SubscriptionPattern("care.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class LoyaltyPrimitive implements Primitive {
  id() { return new PrimitiveId("Loyalty"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("bond.*"), new SubscriptionPattern("commitment.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class MutualConstitutionPrimitive implements Primitive {
  id() { return new PrimitiveId("MutualConstitution"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("bond.*"), new SubscriptionPattern("selfconcept.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RelationalObligationPrimitive implements Primitive {
  id() { return new PrimitiveId("RelationalObligation"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("bond.*"), new SubscriptionPattern("obligation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class GriefPrimitive implements Primitive {
  id() { return new PrimitiveId("Grief"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("bond.*"), new SubscriptionPattern("finitude.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ForgivenessPrimitive implements Primitive {
  id() { return new PrimitiveId("Forgiveness"); }
  layer() { return new Layer(9); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("rupture.*"), new SubscriptionPattern("repair.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 10 (Community) -- 12 primitives --

export class CulturePrimitive implements Primitive {
  id() { return new PrimitiveId("Culture"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("group.*"), new SubscriptionPattern("sharednarrative.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SharedNarrativePrimitive implements Primitive {
  id() { return new PrimitiveId("SharedNarrative"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("narrative.*"), new SubscriptionPattern("group.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class EthosPrimitive implements Primitive {
  id() { return new PrimitiveId("Ethos"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("culture.*"), new SubscriptionPattern("value.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SacredPrimitive implements Primitive {
  id() { return new PrimitiveId("Sacred"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("ethos.*"), new SubscriptionPattern("culture.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("culture.*"), new SubscriptionPattern("memory.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RitualPrimitive implements Primitive {
  id() { return new PrimitiveId("Ritual"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("tradition.*"), new SubscriptionPattern("sacred.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PracticePrimitive implements Primitive {
  id() { return new PrimitiveId("Practice"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("norm.*"), new SubscriptionPattern("tradition.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PlacePrimitive implements Primitive {
  id() { return new PrimitiveId("Place"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("group.*"), new SubscriptionPattern("property.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class BelongingPrimitive implements Primitive {
  id() { return new PrimitiveId("Belonging"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("membership.*"), new SubscriptionPattern("bond.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SolidarityPrimitive implements Primitive {
  id() { return new PrimitiveId("Solidarity"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("belonging.*"), new SubscriptionPattern("collectiveact.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class VoicePrimitive implements Primitive {
  id() { return new PrimitiveId("Voice"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("expression.*"), new SubscriptionPattern("belonging.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class WelcomePrimitive implements Primitive {
  id() { return new PrimitiveId("Welcome"); }
  layer() { return new Layer(10); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("membership.*"), new SubscriptionPattern("care.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 11 (Culture) -- 12 primitives --

export class ReflexivityPrimitive implements Primitive {
  id() { return new PrimitiveId("Reflexivity"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("reflection.*"), new SubscriptionPattern("culture.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class EncounterPrimitive implements Primitive {
  id() { return new PrimitiveId("Encounter"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("recognition.*"), new SubscriptionPattern("culture.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("language.*"), new SubscriptionPattern("encounter.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PluralismPrimitive implements Primitive {
  id() { return new PrimitiveId("Pluralism"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("culture.*"), new SubscriptionPattern("encounter.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CreativityPrimitive implements Primitive {
  id() { return new PrimitiveId("Creativity"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("expression.*"), new SubscriptionPattern("aesthetic.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("value.*"), new SubscriptionPattern("expression.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class InterpretationPrimitive implements Primitive {
  id() { return new PrimitiveId("Interpretation"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("symbol.*"), new SubscriptionPattern("narrative.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DialoguePrimitive implements Primitive {
  id() { return new PrimitiveId("Dialogue"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("signal.*"), new SubscriptionPattern("encounter.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SyncretismPrimitive implements Primitive {
  id() { return new PrimitiveId("Syncretism"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("pluralism.*"), new SubscriptionPattern("translation.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("reflexivity.*"), new SubscriptionPattern("interpretation.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class HegemonyPrimitive implements Primitive {
  id() { return new PrimitiveId("Hegemony"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("authority.*"), new SubscriptionPattern("culture.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CulturalEvolutionPrimitive implements Primitive {
  id() { return new PrimitiveId("CulturalEvolution"); }
  layer() { return new Layer(11); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("culture.*"), new SubscriptionPattern("creativity.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

// -- Layer 12 (Emergence) -- 12 primitives --

export class EmergencePrimitive implements Primitive {
  id() { return new PrimitiveId("Emergence"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("complexity.*"), new SubscriptionPattern("selforganization.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class SelfOrganizationPrimitive implements Primitive {
  id() { return new PrimitiveId("SelfOrganization"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("feedback.*"), new SubscriptionPattern("autopoiesis.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class FeedbackPrimitive implements Primitive {
  id() { return new PrimitiveId("Feedback"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("consequence.*"), new SubscriptionPattern("recursion.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ComplexityPrimitive implements Primitive {
  id() { return new PrimitiveId("Complexity"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("emergence.*"), new SubscriptionPattern("feedback.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ConsciousnessPrimitive implements Primitive {
  id() { return new PrimitiveId("Consciousness"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("selforganization.*"), new SubscriptionPattern("recursion.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class RecursionPrimitive implements Primitive {
  id() { return new PrimitiveId("Recursion"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("feedback.*"), new SubscriptionPattern("selforganization.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ParadoxPrimitive implements Primitive {
  id() { return new PrimitiveId("Paradox"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("recursion.*"), new SubscriptionPattern("incompleteness.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class IncompletenessPrimitive implements Primitive {
  id() { return new PrimitiveId("Incompleteness"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("paradox.*"), new SubscriptionPattern("complexity.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PhaseTransitionPrimitive implements Primitive {
  id() { return new PrimitiveId("PhaseTransition"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("complexity.*"), new SubscriptionPattern("emergence.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class DownwardCausationPrimitive implements Primitive {
  id() { return new PrimitiveId("DownwardCausation"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("emergence.*"), new SubscriptionPattern("feedback.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class AutopoiesisPrimitive implements Primitive {
  id() { return new PrimitiveId("Autopoiesis"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("selforganization.*"), new SubscriptionPattern("emergence.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class CoEvolutionPrimitive implements Primitive {
  id() { return new PrimitiveId("CoEvolution"); }
  layer() { return new Layer(12); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("feedback.*"), new SubscriptionPattern("culturalevolution.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("presence.*"), new SubscriptionPattern("finitude.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class NothingnessPrimitive implements Primitive {
  id() { return new PrimitiveId("Nothingness"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("being.*"), new SubscriptionPattern("groundlessness.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("being.*"), new SubscriptionPattern("contingency.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ContingencyPrimitive implements Primitive {
  id() { return new PrimitiveId("Contingency"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("finitude.*"), new SubscriptionPattern("groundlessness.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("mystery.*"), new SubscriptionPattern("being.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ExistentialAcceptancePrimitive implements Primitive {
  id() { return new PrimitiveId("ExistentialAcceptance"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("finitude.*"), new SubscriptionPattern("contingency.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class PresencePrimitive implements Primitive {
  id() { return new PrimitiveId("Presence"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("being.*"), new SubscriptionPattern("acceptance.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class GratitudePrimitive implements Primitive {
  id() { return new PrimitiveId("Gratitude"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("presence.*"), new SubscriptionPattern("wonder.*")]; }
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
  subscriptions() { return [new SubscriptionPattern("wonder.*"), new SubscriptionPattern("incompleteness.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class TranscendencePrimitive implements Primitive {
  id() { return new PrimitiveId("Transcendence"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("being.*"), new SubscriptionPattern("mystery.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class GroundlessnessPrimitive implements Primitive {
  id() { return new PrimitiveId("Groundlessness"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("contingency.*"), new SubscriptionPattern("nothingness.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export class ReturnPrimitive implements Primitive {
  id() { return new PrimitiveId("Return"); }
  layer() { return new Layer(13); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("acceptance.*"), new SubscriptionPattern("being.*")]; }
  process(tick: number, events: Event[], _snapshot: Snapshot): Mutation[] {
    return [
      { kind: "updateState", primitiveId: this.id(), key: "eventsProcessed", value: events.length },
      { kind: "updateState", primitiveId: this.id(), key: "lastTick", value: tick },
    ];
  }
}

export function createAllPrimitives(): Primitive[] {
  return [
    // Layer 0 — Foundation (45)
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
    // Layer 1 — Agency (12)
    new ValuePrimitive(),
    new IntentPrimitive(),
    new ChoicePrimitive(),
    new RiskPrimitive(),
    new ActPrimitive(),
    new ConsequencePrimitive(),
    new CapacityPrimitive(),
    new ResourcePrimitive(),
    new SignalPrimitive(),
    new ReceptionPrimitive(),
    new AcknowledgmentPrimitive(),
    new CommitmentPrimitive(),
    // Layer 2 — Exchange (12)
    new TermPrimitive(),
    new ProtocolPrimitive(),
    new OfferPrimitive(),
    new AcceptancePrimitive(),
    new AgreementPrimitive(),
    new ObligationPrimitive(),
    new FulfillmentPrimitive(),
    new BreachPrimitive(),
    new ExchangePrimitive(),
    new AccountabilityPrimitive(),
    new DebtPrimitive(),
    new ReciprocityPrimitive(),
    // Layer 3 — Society (12)
    new GroupPrimitive(),
    new MembershipPrimitive(),
    new RolePrimitive(),
    new ConsentPrimitive(),
    new NormPrimitive(),
    new ReputationPrimitive(),
    new SanctionPrimitive(),
    new AuthorityPrimitive(),
    new PropertyPrimitive(),
    new CommonsPrimitive(),
    new GovernancePrimitive(),
    new CollectiveActPrimitive(),
    // Layer 4 — Legal (12)
    new LawPrimitive(),
    new RightPrimitive(),
    new ContractPrimitive(),
    new LiabilityPrimitive(),
    new DueProcessPrimitive(),
    new AdjudicationPrimitive(),
    new RemedyPrimitive(),
    new PrecedentPrimitive(),
    new JurisdictionPrimitive(),
    new SovereigntyPrimitive(),
    new LegitimacyPrimitive(),
    new TreatyPrimitive(),
    // Layer 5 — Technology (12)
    new MethodPrimitive(),
    new MeasurementPrimitive(),
    new KnowledgePrimitive(),
    new ModelPrimitive(),
    new ToolPrimitive(),
    new TechniquePrimitive(),
    new InventionPrimitive(),
    new AbstractionPrimitive(),
    new InfrastructurePrimitive(),
    new StandardPrimitive(),
    new EfficiencyPrimitive(),
    new AutomationPrimitive(),
    // Layer 6 — Information (12)
    new SymbolPrimitive(),
    new LanguagePrimitive(),
    new EncodingPrimitive(),
    new RecordPrimitive(),
    new ChannelPrimitive(),
    new CopyPrimitive(),
    new NoisePrimitive(),
    new RedundancyPrimitive(),
    new DataPrimitive(),
    new ComputationPrimitive(),
    new AlgorithmPrimitive(),
    new EntropyPrimitive(),
    // Layer 7 — Ethics (12)
    new MoralStatusPrimitive(),
    new DignityPrimitive(),
    new AutonomyPrimitive(),
    new FlourishingPrimitive(),
    new DutyPrimitive(),
    new HarmPrimitive(),
    new CarePrimitive(),
    new JusticePrimitive(),
    new ConsciencePrimitive(),
    new VirtuePrimitive(),
    new ResponsibilityPrimitive(),
    new MotivePrimitive(),
    // Layer 8 — Identity (12)
    new NarrativePrimitive(),
    new SelfConceptPrimitive(),
    new ReflectionPrimitive(),
    new MemoryPrimitive(),
    new PurposePrimitive(),
    new AspirationPrimitive(),
    new AuthenticityPrimitive(),
    new ExpressionPrimitive(),
    new GrowthPrimitive(),
    new ContinuityPrimitive(),
    new IntegrationPrimitive(),
    new CrisisPrimitive(),
    // Layer 9 — Relationship (12)
    new BondPrimitive(),
    new AttachmentPrimitive(),
    new RecognitionPrimitive(),
    new IntimacyPrimitive(),
    new AttunementPrimitive(),
    new RupturePrimitive(),
    new RepairPrimitive(),
    new LoyaltyPrimitive(),
    new MutualConstitutionPrimitive(),
    new RelationalObligationPrimitive(),
    new GriefPrimitive(),
    new ForgivenessPrimitive(),
    // Layer 10 — Community (12)
    new CulturePrimitive(),
    new SharedNarrativePrimitive(),
    new EthosPrimitive(),
    new SacredPrimitive(),
    new TraditionPrimitive(),
    new RitualPrimitive(),
    new PracticePrimitive(),
    new PlacePrimitive(),
    new BelongingPrimitive(),
    new SolidarityPrimitive(),
    new VoicePrimitive(),
    new WelcomePrimitive(),
    // Layer 11 — Culture (12)
    new ReflexivityPrimitive(),
    new EncounterPrimitive(),
    new TranslationPrimitive(),
    new PluralismPrimitive(),
    new CreativityPrimitive(),
    new AestheticPrimitive(),
    new InterpretationPrimitive(),
    new DialoguePrimitive(),
    new SyncretismPrimitive(),
    new CritiquePrimitive(),
    new HegemonyPrimitive(),
    new CulturalEvolutionPrimitive(),
    // Layer 12 — Emergence (12)
    new EmergencePrimitive(),
    new SelfOrganizationPrimitive(),
    new FeedbackPrimitive(),
    new ComplexityPrimitive(),
    new ConsciousnessPrimitive(),
    new RecursionPrimitive(),
    new ParadoxPrimitive(),
    new IncompletenessPrimitive(),
    new PhaseTransitionPrimitive(),
    new DownwardCausationPrimitive(),
    new AutopoiesisPrimitive(),
    new CoEvolutionPrimitive(),
    // Layer 13 — Existence (12)
    new BeingPrimitive(),
    new NothingnessPrimitive(),
    new FinitudePrimitive(),
    new ContingencyPrimitive(),
    new WonderPrimitive(),
    new ExistentialAcceptancePrimitive(),
    new PresencePrimitive(),
    new GratitudePrimitive(),
    new MysteryPrimitive(),
    new TranscendencePrimitive(),
    new GroundlessnessPrimitive(),
    new ReturnPrimitive(),
  ];
}
