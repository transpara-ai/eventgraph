// Package layer13 implements the Layer 13 Existence primitives.
// Groups: TheGiven (Being, Nothingness, Finitude, Contingency),
// TheResponse (Wonder, ExistentialAcceptance, Presence, Gratitude),
// TheHorizon (Mystery, Transcendence, Groundlessness, Return).
package layer13

import (
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer13 = types.MustLayer(13)
var cadence1 = types.MustCadence(1)

// --- Group A: The Given ---

// BeingPrimitive notes the simple fact of existence. The most fundamental primitive.
type BeingPrimitive struct{}

func NewBeingPrimitive() *BeingPrimitive { return &BeingPrimitive{} }

func (p *BeingPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Being") }
func (p *BeingPrimitive) Layer() types.Layer               { return layer13 }
func (p *BeingPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *BeingPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *BeingPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("existence.*"),
		types.MustSubscriptionPattern("consciousness.*"),
	}
}

func (p *BeingPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "existence.") || strings.HasPrefix(t, "consciousness.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "alive", Value: true},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "currentTick", Value: tick.Value()},
	}, nil
}

// NothingnessPrimitive registers absence — what is not, as distinct from what is.
type NothingnessPrimitive struct{}

func NewNothingnessPrimitive() *NothingnessPrimitive { return &NothingnessPrimitive{} }

func (p *NothingnessPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Nothingness") }
func (p *NothingnessPrimitive) Layer() types.Layer               { return layer13 }
func (p *NothingnessPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *NothingnessPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *NothingnessPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("loss.*"),
		types.MustSubscriptionPattern("silence.*"),
		types.MustSubscriptionPattern("void.*"),
	}
}

func (p *NothingnessPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "loss.") || strings.HasPrefix(t, "silence.") || strings.HasPrefix(t, "void.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// FinitudePrimitive acknowledges that existence is bounded — everything ends.
type FinitudePrimitive struct{}

func NewFinitudePrimitive() *FinitudePrimitive { return &FinitudePrimitive{} }

func (p *FinitudePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Finitude") }
func (p *FinitudePrimitive) Layer() types.Layer               { return layer13 }
func (p *FinitudePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *FinitudePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *FinitudePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("actor.memorial"),
		types.MustSubscriptionPattern("sustainability.*"),
		types.MustSubscriptionPattern("threshold.*"),
	}
}

func (p *FinitudePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "actor.") || strings.HasPrefix(t, "sustainability.") || strings.HasPrefix(t, "threshold.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ContingencyPrimitive recognises that what exists need not have existed — things could be otherwise.
type ContingencyPrimitive struct{}

func NewContingencyPrimitive() *ContingencyPrimitive { return &ContingencyPrimitive{} }

func (p *ContingencyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Contingency") }
func (p *ContingencyPrimitive) Layer() types.Layer               { return layer13 }
func (p *ContingencyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ContingencyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ContingencyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("choice.*"),
		types.MustSubscriptionPattern("uncertainty.*"),
		types.MustSubscriptionPattern("finitude.*"),
	}
}

func (p *ContingencyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "choice.") || strings.HasPrefix(t, "uncertainty.") || strings.HasPrefix(t, "finitude.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group B: The Response ---

// WonderPrimitive responds to existence with open questioning — not seeking answers but staying with the question.
type WonderPrimitive struct{}

func NewWonderPrimitive() *WonderPrimitive { return &WonderPrimitive{} }

func (p *WonderPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Wonder") }
func (p *WonderPrimitive) Layer() types.Layer               { return layer13 }
func (p *WonderPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *WonderPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *WonderPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("mystery.*"),
		types.MustSubscriptionPattern("being.*"),
		types.MustSubscriptionPattern("emergence.*"),
	}
}

func (p *WonderPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "mystery.") || strings.HasPrefix(t, "being.") || strings.HasPrefix(t, "emergence.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ExistentialAcceptancePrimitive acknowledges what is without trying to change it — the stopping condition.
// Named ExistentialAcceptance to disambiguate from Layer 2's Acceptance (of an offer).
type ExistentialAcceptancePrimitive struct{}

func NewExistentialAcceptancePrimitive() *ExistentialAcceptancePrimitive { return &ExistentialAcceptancePrimitive{} }

func (p *ExistentialAcceptancePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("ExistentialAcceptance") }
func (p *ExistentialAcceptancePrimitive) Layer() types.Layer               { return layer13 }
func (p *ExistentialAcceptancePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ExistentialAcceptancePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ExistentialAcceptancePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("finitude.*"),
		types.MustSubscriptionPattern("contingency.*"),
		types.MustSubscriptionPattern("nothingness.*"),
	}
}

func (p *ExistentialAcceptancePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "finitude.") || strings.HasPrefix(t, "contingency.") || strings.HasPrefix(t, "nothingness.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// PresencePrimitive sustains full awareness in the current moment without drifting to past or future.
type PresencePrimitive struct{}

func NewPresencePrimitive() *PresencePrimitive { return &PresencePrimitive{} }

func (p *PresencePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Presence") }
func (p *PresencePrimitive) Layer() types.Layer               { return layer13 }
func (p *PresencePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *PresencePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *PresencePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("clock.tick"),
		types.MustSubscriptionPattern("being.*"),
		types.MustSubscriptionPattern("awareness.*"),
	}
}

func (p *PresencePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "clock.") || strings.HasPrefix(t, "being.") || strings.HasPrefix(t, "awareness.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// GratitudePrimitive expresses thankfulness for existence itself — not transactional, but existential.
type GratitudePrimitive struct{}

func NewGratitudePrimitive() *GratitudePrimitive { return &GratitudePrimitive{} }

func (p *GratitudePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Gratitude") }
func (p *GratitudePrimitive) Layer() types.Layer               { return layer13 }
func (p *GratitudePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *GratitudePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *GratitudePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("being.*"),
		types.MustSubscriptionPattern("milestone.*"),
		types.MustSubscriptionPattern("presence.*"),
	}
}

func (p *GratitudePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "being.") || strings.HasPrefix(t, "milestone.") || strings.HasPrefix(t, "presence.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group C: The Horizon ---

// MysteryPrimitive acknowledges what cannot be known — the boundary of understanding.
type MysteryPrimitive struct{}

func NewMysteryPrimitive() *MysteryPrimitive { return &MysteryPrimitive{} }

func (p *MysteryPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Mystery") }
func (p *MysteryPrimitive) Layer() types.Layer               { return layer13 }
func (p *MysteryPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *MysteryPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *MysteryPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("uncertainty.*"),
		types.MustSubscriptionPattern("incompleteness.*"),
		types.MustSubscriptionPattern("wonder.*"),
	}
}

func (p *MysteryPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "uncertainty.") || strings.HasPrefix(t, "incompleteness.") || strings.HasPrefix(t, "wonder.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// TranscendencePrimitive points beyond the current frame — what exceeds any system's capacity to contain.
type TranscendencePrimitive struct{}

func NewTranscendencePrimitive() *TranscendencePrimitive { return &TranscendencePrimitive{} }

func (p *TranscendencePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Transcendence") }
func (p *TranscendencePrimitive) Layer() types.Layer               { return layer13 }
func (p *TranscendencePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *TranscendencePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *TranscendencePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("mystery.*"),
		types.MustSubscriptionPattern("paradox.*"),
		types.MustSubscriptionPattern("phasetransition.*"),
	}
}

func (p *TranscendencePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "mystery.") || strings.HasPrefix(t, "paradox.") || strings.HasPrefix(t, "phasetransition.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// GroundlessnessPrimitive recognises that there is no ultimate foundation — all grounds are themselves groundless.
type GroundlessnessPrimitive struct{}

func NewGroundlessnessPrimitive() *GroundlessnessPrimitive { return &GroundlessnessPrimitive{} }

func (p *GroundlessnessPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Groundlessness") }
func (p *GroundlessnessPrimitive) Layer() types.Layer               { return layer13 }
func (p *GroundlessnessPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *GroundlessnessPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *GroundlessnessPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("contingency.*"),
		types.MustSubscriptionPattern("nothingness.*"),
		types.MustSubscriptionPattern("incompleteness.*"),
	}
}

func (p *GroundlessnessPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "contingency.") || strings.HasPrefix(t, "nothingness.") || strings.HasPrefix(t, "incompleteness.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ReturnPrimitive completes the cycle — from existence back to existence, the loop that never closes.
type ReturnPrimitive struct{}

func NewReturnPrimitive() *ReturnPrimitive { return &ReturnPrimitive{} }

func (p *ReturnPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Return") }
func (p *ReturnPrimitive) Layer() types.Layer               { return layer13 }
func (p *ReturnPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ReturnPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ReturnPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{types.MustSubscriptionPattern("*")}
}

func (p *ReturnPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
