// Package layer10 implements the Layer 10 Community primitives.
// Groups: SharedMeaning (Culture, SharedNarrative, Ethos, Sacred),
// LivingPractice (Tradition, Ritual, Practice, Place),
// CommunalExperience (Belonging, Solidarity, Voice, Welcome).
package layer10

import (
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer10 = types.MustLayer(10)
var cadence1 = types.MustCadence(1)

// --- Group A: Shared Meaning ---

// CulturePrimitive tracks the shared patterns of meaning within a community.
type CulturePrimitive struct{}

func NewCulturePrimitive() *CulturePrimitive { return &CulturePrimitive{} }

func (p *CulturePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Culture") }
func (p *CulturePrimitive) Layer() types.Layer               { return layer10 }
func (p *CulturePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *CulturePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *CulturePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("norm.*"),
		types.MustSubscriptionPattern("tradition.*"),
		types.MustSubscriptionPattern("ethos.*"),
	}
}

func (p *CulturePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "norm.") || strings.HasPrefix(t, "tradition.") || strings.HasPrefix(t, "ethos.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// SharedNarrativePrimitive maintains the stories a community tells about itself.
type SharedNarrativePrimitive struct{}

func NewSharedNarrativePrimitive() *SharedNarrativePrimitive { return &SharedNarrativePrimitive{} }

func (p *SharedNarrativePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("SharedNarrative") }
func (p *SharedNarrativePrimitive) Layer() types.Layer               { return layer10 }
func (p *SharedNarrativePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *SharedNarrativePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *SharedNarrativePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("narrative.*"),
		types.MustSubscriptionPattern("milestone.*"),
		types.MustSubscriptionPattern("memorial.*"),
	}
}

func (p *SharedNarrativePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "narrative.") || strings.HasPrefix(t, "milestone.") || strings.HasPrefix(t, "memorial.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// EthosPrimitive captures the moral character and guiding values of a community.
type EthosPrimitive struct{}

func NewEthosPrimitive() *EthosPrimitive { return &EthosPrimitive{} }

func (p *EthosPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Ethos") }
func (p *EthosPrimitive) Layer() types.Layer               { return layer10 }
func (p *EthosPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *EthosPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *EthosPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("value.*"),
		types.MustSubscriptionPattern("norm.*"),
		types.MustSubscriptionPattern("governance.*"),
	}
}

func (p *EthosPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "value.") || strings.HasPrefix(t, "norm.") || strings.HasPrefix(t, "governance.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// SacredPrimitive identifies what a community holds inviolable.
type SacredPrimitive struct{}

func NewSacredPrimitive() *SacredPrimitive { return &SacredPrimitive{} }

func (p *SacredPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Sacred") }
func (p *SacredPrimitive) Layer() types.Layer               { return layer10 }
func (p *SacredPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *SacredPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *SacredPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("ethos.*"),
		types.MustSubscriptionPattern("ritual.*"),
		types.MustSubscriptionPattern("boundary.*"),
	}
}

func (p *SacredPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "ethos.") || strings.HasPrefix(t, "ritual.") || strings.HasPrefix(t, "boundary.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group B: Living Practice ---

// TraditionPrimitive identifies practices passed down that define a community.
type TraditionPrimitive struct{}

func NewTraditionPrimitive() *TraditionPrimitive { return &TraditionPrimitive{} }

func (p *TraditionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Tradition") }
func (p *TraditionPrimitive) Layer() types.Layer               { return layer10 }
func (p *TraditionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *TraditionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *TraditionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("convention.*"),
		types.MustSubscriptionPattern("heritage.*"),
		types.MustSubscriptionPattern("pattern.detected"),
	}
}

func (p *TraditionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "convention.") || strings.HasPrefix(t, "heritage.") || strings.HasPrefix(t, "pattern.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RitualPrimitive handles formal, repeated communal actions that carry shared meaning.
type RitualPrimitive struct{}

func NewRitualPrimitive() *RitualPrimitive { return &RitualPrimitive{} }

func (p *RitualPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Ritual") }
func (p *RitualPrimitive) Layer() types.Layer               { return layer10 }
func (p *RitualPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RitualPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RitualPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("ceremony.*"),
		types.MustSubscriptionPattern("tradition.*"),
		types.MustSubscriptionPattern("sacred.*"),
	}
}

func (p *RitualPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "ceremony.") || strings.HasPrefix(t, "tradition.") || strings.HasPrefix(t, "sacred.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// PracticePrimitive tracks embodied, repeated communal activities.
type PracticePrimitive struct{}

func NewPracticePrimitive() *PracticePrimitive { return &PracticePrimitive{} }

func (p *PracticePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Practice") }
func (p *PracticePrimitive) Layer() types.Layer               { return layer10 }
func (p *PracticePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *PracticePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *PracticePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("habit.*"),
		types.MustSubscriptionPattern("skill.*"),
		types.MustSubscriptionPattern("contribution.*"),
	}
}

func (p *PracticePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "habit.") || strings.HasPrefix(t, "skill.") || strings.HasPrefix(t, "contribution.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// PlacePrimitive identifies the sense of shared location and environment.
type PlacePrimitive struct{}

func NewPlacePrimitive() *PlacePrimitive { return &PlacePrimitive{} }

func (p *PlacePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Place") }
func (p *PlacePrimitive) Layer() types.Layer               { return layer10 }
func (p *PlacePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *PlacePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *PlacePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("group.*"),
		types.MustSubscriptionPattern("commons.*"),
		types.MustSubscriptionPattern("presence.*"),
	}
}

func (p *PlacePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "group.") || strings.HasPrefix(t, "commons.") || strings.HasPrefix(t, "presence.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group C: Communal Experience ---

// BelongingPrimitive tracks the felt sense of being part of a community.
type BelongingPrimitive struct{}

func NewBelongingPrimitive() *BelongingPrimitive { return &BelongingPrimitive{} }

func (p *BelongingPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Belonging") }
func (p *BelongingPrimitive) Layer() types.Layer               { return layer10 }
func (p *BelongingPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *BelongingPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *BelongingPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("inclusion.*"),
		types.MustSubscriptionPattern("welcome.*"),
		types.MustSubscriptionPattern("recognition.*"),
	}
}

func (p *BelongingPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "inclusion.") || strings.HasPrefix(t, "welcome.") || strings.HasPrefix(t, "recognition.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// SolidarityPrimitive tracks mutual support and collective action within a community.
type SolidarityPrimitive struct{}

func NewSolidarityPrimitive() *SolidarityPrimitive { return &SolidarityPrimitive{} }

func (p *SolidarityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Solidarity") }
func (p *SolidarityPrimitive) Layer() types.Layer               { return layer10 }
func (p *SolidarityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *SolidarityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *SolidarityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("cooperation.*"),
		types.MustSubscriptionPattern("sacrifice.*"),
		types.MustSubscriptionPattern("care.*"),
	}
}

func (p *SolidarityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "cooperation.") || strings.HasPrefix(t, "sacrifice.") || strings.HasPrefix(t, "care.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// VoicePrimitive ensures every member can be heard within the community.
type VoicePrimitive struct{}

func NewVoicePrimitive() *VoicePrimitive { return &VoicePrimitive{} }

func (p *VoicePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Voice") }
func (p *VoicePrimitive) Layer() types.Layer               { return layer10 }
func (p *VoicePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *VoicePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *VoicePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("expression.*"),
		types.MustSubscriptionPattern("dissent.*"),
		types.MustSubscriptionPattern("fairness.*"),
	}
}

func (p *VoicePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "expression.") || strings.HasPrefix(t, "dissent.") || strings.HasPrefix(t, "fairness.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// WelcomePrimitive handles how a community receives new members.
type WelcomePrimitive struct{}

func NewWelcomePrimitive() *WelcomePrimitive { return &WelcomePrimitive{} }

func (p *WelcomePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Welcome") }
func (p *WelcomePrimitive) Layer() types.Layer               { return layer10 }
func (p *WelcomePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *WelcomePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *WelcomePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("actor.joined"),
		types.MustSubscriptionPattern("belonging.*"),
		types.MustSubscriptionPattern("onboarding.*"),
	}
}

func (p *WelcomePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "actor.") || strings.HasPrefix(t, "belonging.") || strings.HasPrefix(t, "onboarding.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
