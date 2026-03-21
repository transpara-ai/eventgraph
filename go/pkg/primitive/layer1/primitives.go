// Package layer1 implements the Layer 1 Agency primitives.
// Groups: Volition (Value, Intent, Choice, Risk),
// Action (Act, Consequence, Capacity, Resource),
// Communication (Signal, Reception, Acknowledgment, Commitment).
package layer1

import (
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer1 = types.MustLayer(1)
var cadence1 = types.MustCadence(1)

// --- Group A: Volition (why act) ---

// ValuePrimitive measures importance relative to Self. What matters and how much.
type ValuePrimitive struct{}

func NewValuePrimitive() *ValuePrimitive { return &ValuePrimitive{} }

func (p *ValuePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Value") }
func (p *ValuePrimitive) Layer() types.Layer               { return layer1 }
func (p *ValuePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ValuePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ValuePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("decision.*"),
		types.MustSubscriptionPattern("actor.*"),
	}
}

func (p *ValuePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "decision.") || strings.HasPrefix(t, "actor.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// IntentPrimitive represents a desired future state the system seeks to bring about.
type IntentPrimitive struct{}

func NewIntentPrimitive() *IntentPrimitive { return &IntentPrimitive{} }

func (p *IntentPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Intent") }
func (p *IntentPrimitive) Layer() types.Layer               { return layer1 }
func (p *IntentPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *IntentPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *IntentPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("value.*"),
		types.MustSubscriptionPattern("expectation.*"),
	}
}

func (p *IntentPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "value.") || strings.HasPrefix(t, "expectation.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ChoicePrimitive selects among possible Acts based on Value and Confidence.
type ChoicePrimitive struct{}

func NewChoicePrimitive() *ChoicePrimitive { return &ChoicePrimitive{} }

func (p *ChoicePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Choice") }
func (p *ChoicePrimitive) Layer() types.Layer               { return layer1 }
func (p *ChoicePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ChoicePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ChoicePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("intent.*"),
		types.MustSubscriptionPattern("value.*"),
		types.MustSubscriptionPattern("confidence.*"),
	}
}

func (p *ChoicePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "intent.") || strings.HasPrefix(t, "value.") || strings.HasPrefix(t, "confidence.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RiskPrimitive assesses potential loss from an Act under Uncertainty.
type RiskPrimitive struct{}

func NewRiskPrimitive() *RiskPrimitive { return &RiskPrimitive{} }

func (p *RiskPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Risk") }
func (p *RiskPrimitive) Layer() types.Layer               { return layer1 }
func (p *RiskPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RiskPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RiskPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("intent.*"),
		types.MustSubscriptionPattern("uncertainty.*"),
		types.MustSubscriptionPattern("value.*"),
	}
}

func (p *RiskPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "intent.") || strings.HasPrefix(t, "uncertainty.") || strings.HasPrefix(t, "value.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group B: Action (doing things) ---

// ActPrimitive produces causally effective Events. Self becomes a FirstCause.
type ActPrimitive struct{}

func NewActPrimitive() *ActPrimitive { return &ActPrimitive{} }

func (p *ActPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Act") }
func (p *ActPrimitive) Layer() types.Layer               { return layer1 }
func (p *ActPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ActPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ActPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("choice.*"),
		types.MustSubscriptionPattern("intent.*"),
	}
}

func (p *ActPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "choice.") || strings.HasPrefix(t, "intent.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ConsequencePrimitive tracks effects of an Act attributed back to the actor.
type ConsequencePrimitive struct{}

func NewConsequencePrimitive() *ConsequencePrimitive { return &ConsequencePrimitive{} }

func (p *ConsequencePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Consequence") }
func (p *ConsequencePrimitive) Layer() types.Layer               { return layer1 }
func (p *ConsequencePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ConsequencePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ConsequencePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("act.*"),
		types.MustSubscriptionPattern("violation.*"),
	}
}

func (p *ConsequencePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "act.") || strings.HasPrefix(t, "violation.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// CapacityPrimitive tracks what the system is able to do.
type CapacityPrimitive struct{}

func NewCapacityPrimitive() *CapacityPrimitive { return &CapacityPrimitive{} }

func (p *CapacityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Capacity") }
func (p *CapacityPrimitive) Layer() types.Layer               { return layer1 }
func (p *CapacityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *CapacityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *CapacityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("actor.*"),
		types.MustSubscriptionPattern("resource.*"),
		types.MustSubscriptionPattern("trust.*"),
	}
}

func (p *CapacityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "actor.") || strings.HasPrefix(t, "resource.") || strings.HasPrefix(t, "trust.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ResourcePrimitive tracks finite things consumed or required by Acts.
type ResourcePrimitive struct{}

func NewResourcePrimitive() *ResourcePrimitive { return &ResourcePrimitive{} }

func (p *ResourcePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Resource") }
func (p *ResourcePrimitive) Layer() types.Layer               { return layer1 }
func (p *ResourcePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ResourcePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ResourcePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("act.*"),
		types.MustSubscriptionPattern("budget.*"),
	}
}

func (p *ResourcePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "act.") || strings.HasPrefix(t, "budget.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group C: Communication (exchanging with others) ---

// SignalPrimitive is an Act directed at a specific ActorID to convey information.
type SignalPrimitive struct{}

func NewSignalPrimitive() *SignalPrimitive { return &SignalPrimitive{} }

func (p *SignalPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Signal") }
func (p *SignalPrimitive) Layer() types.Layer               { return layer1 }
func (p *SignalPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *SignalPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *SignalPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("act.*"),
		types.MustSubscriptionPattern("actor.*"),
	}
}

func (p *SignalPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "act.") || strings.HasPrefix(t, "actor.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ReceptionPrimitive handles external Events entering Self's awareness.
type ReceptionPrimitive struct{}

func NewReceptionPrimitive() *ReceptionPrimitive { return &ReceptionPrimitive{} }

func (p *ReceptionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Reception") }
func (p *ReceptionPrimitive) Layer() types.Layer               { return layer1 }
func (p *ReceptionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ReceptionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ReceptionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{types.MustSubscriptionPattern("*")}
}

func (p *ReceptionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsReceived", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AcknowledgmentPrimitive confirms receipt of a prior Signal.
type AcknowledgmentPrimitive struct{}

func NewAcknowledgmentPrimitive() *AcknowledgmentPrimitive { return &AcknowledgmentPrimitive{} }

func (p *AcknowledgmentPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Acknowledgment") }
func (p *AcknowledgmentPrimitive) Layer() types.Layer               { return layer1 }
func (p *AcknowledgmentPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AcknowledgmentPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AcknowledgmentPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("signal.*"),
	}
}

func (p *AcknowledgmentPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	signals := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "signal.") {
			signals++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "signalsAcknowledged", Value: signals},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// CommitmentPrimitive is a Signal that binds future behavior.
type CommitmentPrimitive struct{}

func NewCommitmentPrimitive() *CommitmentPrimitive { return &CommitmentPrimitive{} }

func (p *CommitmentPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Commitment") }
func (p *CommitmentPrimitive) Layer() types.Layer               { return layer1 }
func (p *CommitmentPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *CommitmentPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *CommitmentPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("signal.*"),
		types.MustSubscriptionPattern("agreement.*"),
		types.MustSubscriptionPattern("intent.*"),
	}
}

func (p *CommitmentPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "signal.") || strings.HasPrefix(t, "agreement.") || strings.HasPrefix(t, "intent.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
