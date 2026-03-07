// Package layer7 implements the Layer 7 Ethics primitives.
// Groups: Value (Value, Harm, Fairness, Care),
// Judgement (Dilemma, Proportionality, Intention, Consequence),
// Accountability (Responsibility, Transparency, Redress, Growth).
package layer7

import (
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer7 = types.MustLayer(7)
var cadence1 = types.MustCadence(1)

// --- Group 0: Value ---

// ValuePrimitive identifies and weights values.
type ValuePrimitive struct{}

func NewValuePrimitive() *ValuePrimitive { return &ValuePrimitive{} }

func (p *ValuePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Value") }
func (p *ValuePrimitive) Layer() types.Layer               { return layer7 }
func (p *ValuePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ValuePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ValuePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("consensus.*"),
		types.MustSubscriptionPattern("norm.*"),
		types.MustSubscriptionPattern("right.*"),
	}
}

func (p *ValuePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// HarmPrimitive detects and measures harm.
type HarmPrimitive struct{}

func NewHarmPrimitive() *HarmPrimitive { return &HarmPrimitive{} }

func (p *HarmPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Harm") }
func (p *HarmPrimitive) Layer() types.Layer               { return layer7 }
func (p *HarmPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *HarmPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *HarmPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("violation.*"),
		types.MustSubscriptionPattern("right.violated"),
		types.MustSubscriptionPattern("exclusion.*"),
	}
}

func (p *HarmPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	violations := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "violation.") || strings.HasPrefix(ev.Type().Value(), "right.") {
			violations++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "harmsDetected", Value: violations},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// FairnessPrimitive evaluates equitable treatment.
type FairnessPrimitive struct{}

func NewFairnessPrimitive() *FairnessPrimitive { return &FairnessPrimitive{} }

func (p *FairnessPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Fairness") }
func (p *FairnessPrimitive) Layer() types.Layer               { return layer7 }
func (p *FairnessPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *FairnessPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *FairnessPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("decision.*"),
		types.MustSubscriptionPattern("sanction.*"),
		types.MustSubscriptionPattern("exclusion.*"),
		types.MustSubscriptionPattern("bias.detected"),
	}
}

func (p *FairnessPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// CarePrimitive prioritises wellbeing. The soul statement flows through this primitive.
type CarePrimitive struct{}

func NewCarePrimitive() *CarePrimitive { return &CarePrimitive{} }

func (p *CarePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Care") }
func (p *CarePrimitive) Layer() types.Layer               { return layer7 }
func (p *CarePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *CarePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *CarePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("harm.*"),
		types.MustSubscriptionPattern("health.*"),
		types.MustSubscriptionPattern("trust.*"),
	}
}

func (p *CarePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 1: Judgement ---

// DilemmaPrimitive detects situations where values conflict.
type DilemmaPrimitive struct{}

func NewDilemmaPrimitive() *DilemmaPrimitive { return &DilemmaPrimitive{} }

func (p *DilemmaPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Dilemma") }
func (p *DilemmaPrimitive) Layer() types.Layer               { return layer7 }
func (p *DilemmaPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *DilemmaPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *DilemmaPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("value.conflict"),
		types.MustSubscriptionPattern("decision.*"),
	}
}

func (p *DilemmaPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ProportionalityPrimitive ensures responses match the severity of the situation.
type ProportionalityPrimitive struct{}

func NewProportionalityPrimitive() *ProportionalityPrimitive { return &ProportionalityPrimitive{} }

func (p *ProportionalityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Proportionality") }
func (p *ProportionalityPrimitive) Layer() types.Layer               { return layer7 }
func (p *ProportionalityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ProportionalityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ProportionalityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("enforcement.*"),
		types.MustSubscriptionPattern("sanction.*"),
		types.MustSubscriptionPattern("harm.*"),
	}
}

func (p *ProportionalityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// IntentionPrimitive evaluates the purpose behind actions.
type IntentionPrimitive struct{}

func NewIntentionPrimitive() *IntentionPrimitive { return &IntentionPrimitive{} }

func (p *IntentionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Intention") }
func (p *IntentionPrimitive) Layer() types.Layer               { return layer7 }
func (p *IntentionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *IntentionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *IntentionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("decision.*"),
		types.MustSubscriptionPattern("goal.*"),
		types.MustSubscriptionPattern("initiative.*"),
	}
}

func (p *IntentionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ConsequencePrimitive evaluates outcomes of actions.
type ConsequencePrimitive struct{}

func NewConsequencePrimitive() *ConsequencePrimitive { return &ConsequencePrimitive{} }

func (p *ConsequencePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Consequence") }
func (p *ConsequencePrimitive) Layer() types.Layer               { return layer7 }
func (p *ConsequencePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ConsequencePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ConsequencePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("decision.*"),
		types.MustSubscriptionPattern("harm.*"),
		types.MustSubscriptionPattern("goal.achieved"),
		types.MustSubscriptionPattern("goal.abandoned"),
	}
}

func (p *ConsequencePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 2: Accountability ---

// ResponsibilityPrimitive assigns moral responsibility (not just causal).
type ResponsibilityPrimitive struct{}

func NewResponsibilityPrimitive() *ResponsibilityPrimitive { return &ResponsibilityPrimitive{} }

func (p *ResponsibilityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Responsibility") }
func (p *ResponsibilityPrimitive) Layer() types.Layer               { return layer7 }
func (p *ResponsibilityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ResponsibilityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ResponsibilityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("intention.*"),
		types.MustSubscriptionPattern("consequence.*"),
		types.MustSubscriptionPattern("accountability.traced"),
	}
}

func (p *ResponsibilityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// TransparencyPrimitive makes reasoning visible. The TRANSPARENT invariant flows through this.
type TransparencyPrimitive struct{}

func NewTransparencyPrimitive() *TransparencyPrimitive { return &TransparencyPrimitive{} }

func (p *TransparencyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Transparency") }
func (p *TransparencyPrimitive) Layer() types.Layer               { return layer7 }
func (p *TransparencyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *TransparencyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *TransparencyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("decision.*"),
		types.MustSubscriptionPattern("adjudication.*"),
	}
}

func (p *TransparencyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RedressPrimitive makes things right after harm.
type RedressPrimitive struct{}

func NewRedressPrimitive() *RedressPrimitive { return &RedressPrimitive{} }

func (p *RedressPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Redress") }
func (p *RedressPrimitive) Layer() types.Layer               { return layer7 }
func (p *RedressPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RedressPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RedressPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("harm.*"),
		types.MustSubscriptionPattern("responsibility.*"),
	}
}

func (p *RedressPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// GrowthPrimitive tracks learning from ethical failures.
type GrowthPrimitive struct{}

func NewGrowthPrimitive() *GrowthPrimitive { return &GrowthPrimitive{} }

func (p *GrowthPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Growth") }
func (p *GrowthPrimitive) Layer() types.Layer               { return layer7 }
func (p *GrowthPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *GrowthPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *GrowthPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("redress.*"),
		types.MustSubscriptionPattern("responsibility.*"),
		types.MustSubscriptionPattern("learning.*"),
	}
}

func (p *GrowthPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
