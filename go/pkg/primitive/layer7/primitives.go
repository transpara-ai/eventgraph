// Package layer7 implements the Layer 7 Ethics primitives.
// Groups: MoralStanding (MoralStatus, Dignity, Autonomy, Flourishing),
// MoralObligation (Duty, Harm, Care, Justice),
// MoralAgency (Conscience, Virtue, Responsibility, Motive).
package layer7

import (
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer7 = types.MustLayer(7)
var cadence1 = types.MustCadence(1)

// --- Group A: Moral Standing ---

// MoralStatusPrimitive determines whether an entity deserves moral consideration.
type MoralStatusPrimitive struct{}

func NewMoralStatusPrimitive() *MoralStatusPrimitive { return &MoralStatusPrimitive{} }

func (p *MoralStatusPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("MoralStatus") }
func (p *MoralStatusPrimitive) Layer() types.Layer               { return layer7 }
func (p *MoralStatusPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *MoralStatusPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *MoralStatusPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("dignity.*"),
		types.MustSubscriptionPattern("harm.*"),
	}
}

func (p *MoralStatusPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "dignity.") || strings.HasPrefix(t, "harm.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// DignityPrimitive tracks the inherent worth of moral subjects that must not be instrumentalised.
type DignityPrimitive struct{}

func NewDignityPrimitive() *DignityPrimitive { return &DignityPrimitive{} }

func (p *DignityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Dignity") }
func (p *DignityPrimitive) Layer() types.Layer               { return layer7 }
func (p *DignityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *DignityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *DignityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("moralstatus.*"),
		types.MustSubscriptionPattern("care.*"),
	}
}

func (p *DignityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "moralstatus.") || strings.HasPrefix(t, "care.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AutonomyPrimitive represents the capacity for self-governance and free choice.
type AutonomyPrimitive struct{}

func NewAutonomyPrimitive() *AutonomyPrimitive { return &AutonomyPrimitive{} }

func (p *AutonomyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Autonomy") }
func (p *AutonomyPrimitive) Layer() types.Layer               { return layer7 }
func (p *AutonomyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AutonomyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AutonomyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("dignity.*"),
		types.MustSubscriptionPattern("conscience.*"),
	}
}

func (p *AutonomyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "dignity.") || strings.HasPrefix(t, "conscience.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// FlourishingPrimitive represents the full realisation of a moral subject's potential.
type FlourishingPrimitive struct{}

func NewFlourishingPrimitive() *FlourishingPrimitive { return &FlourishingPrimitive{} }

func (p *FlourishingPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Flourishing") }
func (p *FlourishingPrimitive) Layer() types.Layer               { return layer7 }
func (p *FlourishingPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *FlourishingPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *FlourishingPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("autonomy.*"),
		types.MustSubscriptionPattern("virtue.*"),
	}
}

func (p *FlourishingPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "autonomy.") || strings.HasPrefix(t, "virtue.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group B: Moral Obligation ---

// DutyPrimitive represents an obligation that binds regardless of desire or consequence.
type DutyPrimitive struct{}

func NewDutyPrimitive() *DutyPrimitive { return &DutyPrimitive{} }

func (p *DutyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Duty") }
func (p *DutyPrimitive) Layer() types.Layer               { return layer7 }
func (p *DutyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *DutyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *DutyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("responsibility.*"),
		types.MustSubscriptionPattern("justice.*"),
	}
}

func (p *DutyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "responsibility.") || strings.HasPrefix(t, "justice.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// HarmPrimitive detects and measures damage to moral subjects.
type HarmPrimitive struct{}

func NewHarmPrimitive() *HarmPrimitive { return &HarmPrimitive{} }

func (p *HarmPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Harm") }
func (p *HarmPrimitive) Layer() types.Layer               { return layer7 }
func (p *HarmPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *HarmPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *HarmPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("violation.*"),
		types.MustSubscriptionPattern("dignity.*"),
	}
}

func (p *HarmPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	violations := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "violation.") || strings.HasPrefix(t, "dignity.") {
			violations++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "harmsDetected", Value: violations},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// CarePrimitive prioritises the wellbeing of moral subjects. The soul statement flows through this primitive.
type CarePrimitive struct{}

func NewCarePrimitive() *CarePrimitive { return &CarePrimitive{} }

func (p *CarePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Care") }
func (p *CarePrimitive) Layer() types.Layer               { return layer7 }
func (p *CarePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *CarePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *CarePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("harm.*"),
		types.MustSubscriptionPattern("flourishing.*"),
	}
}

func (p *CarePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "harm.") || strings.HasPrefix(t, "flourishing.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// JusticePrimitive evaluates the fair distribution of benefits and burdens.
type JusticePrimitive struct{}

func NewJusticePrimitive() *JusticePrimitive { return &JusticePrimitive{} }

func (p *JusticePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Justice") }
func (p *JusticePrimitive) Layer() types.Layer               { return layer7 }
func (p *JusticePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *JusticePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *JusticePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("duty.*"),
		types.MustSubscriptionPattern("harm.*"),
		types.MustSubscriptionPattern("dignity.*"),
	}
}

func (p *JusticePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "duty.") || strings.HasPrefix(t, "harm.") || strings.HasPrefix(t, "dignity.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group C: Moral Agency ---

// ConsciencePrimitive represents the internal capacity to distinguish right from wrong.
type ConsciencePrimitive struct{}

func NewConsciencePrimitive() *ConsciencePrimitive { return &ConsciencePrimitive{} }

func (p *ConsciencePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Conscience") }
func (p *ConsciencePrimitive) Layer() types.Layer               { return layer7 }
func (p *ConsciencePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ConsciencePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ConsciencePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("duty.*"),
		types.MustSubscriptionPattern("harm.*"),
		types.MustSubscriptionPattern("motive.*"),
	}
}

func (p *ConsciencePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "duty.") || strings.HasPrefix(t, "harm.") || strings.HasPrefix(t, "motive.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// VirtuePrimitive represents a stable disposition toward morally good action.
type VirtuePrimitive struct{}

func NewVirtuePrimitive() *VirtuePrimitive { return &VirtuePrimitive{} }

func (p *VirtuePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Virtue") }
func (p *VirtuePrimitive) Layer() types.Layer               { return layer7 }
func (p *VirtuePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *VirtuePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *VirtuePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("conscience.*"),
		types.MustSubscriptionPattern("care.*"),
	}
}

func (p *VirtuePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "conscience.") || strings.HasPrefix(t, "care.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ResponsibilityPrimitive assigns moral accountability for actions and their consequences.
type ResponsibilityPrimitive struct{}

func NewResponsibilityPrimitive() *ResponsibilityPrimitive { return &ResponsibilityPrimitive{} }

func (p *ResponsibilityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Responsibility") }
func (p *ResponsibilityPrimitive) Layer() types.Layer               { return layer7 }
func (p *ResponsibilityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ResponsibilityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ResponsibilityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("conscience.*"),
		types.MustSubscriptionPattern("motive.*"),
		types.MustSubscriptionPattern("harm.*"),
	}
}

func (p *ResponsibilityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "conscience.") || strings.HasPrefix(t, "motive.") || strings.HasPrefix(t, "harm.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// MotivePrimitive represents the underlying reason or intention driving a moral action.
type MotivePrimitive struct{}

func NewMotivePrimitive() *MotivePrimitive { return &MotivePrimitive{} }

func (p *MotivePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Motive") }
func (p *MotivePrimitive) Layer() types.Layer               { return layer7 }
func (p *MotivePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *MotivePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *MotivePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("virtue.*"),
		types.MustSubscriptionPattern("duty.*"),
	}
}

func (p *MotivePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "virtue.") || strings.HasPrefix(t, "duty.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
