// Package layer13 implements the Layer 13 Existence primitives.
// Groups: Being (Being, Finitude, Change, Interdependence),
// Boundary (Mystery, Paradox, Infinity, Void),
// Wonder (Awe, ExistentialGratitude, Play, Wonder).
package layer13

import (
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer13 = types.MustLayer(13)
var cadence1 = types.MustCadence(1)

// --- Group 0: Being ---

// BeingPrimitive notes the simple fact of existence. The simplest primitive.
type BeingPrimitive struct{}

func NewBeingPrimitive() *BeingPrimitive { return &BeingPrimitive{} }

func (p *BeingPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Being") }
func (p *BeingPrimitive) Layer() types.Layer               { return layer13 }
func (p *BeingPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *BeingPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *BeingPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{types.MustSubscriptionPattern("clock.tick")}
}

func (p *BeingPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "alive", Value: true},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "currentTick", Value: tick.Value()},
	}, nil
}

// FinitudePrimitive acknowledges that existence is bounded.
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
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ChangePrimitive observes that everything changes.
type ChangePrimitive struct {
	store store.Store
}

func NewChangePrimitive(s store.Store) *ChangePrimitive { return &ChangePrimitive{store: s} }

func (p *ChangePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Change") }
func (p *ChangePrimitive) Layer() types.Layer               { return layer13 }
func (p *ChangePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ChangePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ChangePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{types.MustSubscriptionPattern("*")}
}

func (p *ChangePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsThisTick", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// InterdependencePrimitive maps how nothing exists alone.
type InterdependencePrimitive struct {
	store store.Store
}

func NewInterdependencePrimitive(s store.Store) *InterdependencePrimitive {
	return &InterdependencePrimitive{store: s}
}

func (p *InterdependencePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Interdependence") }
func (p *InterdependencePrimitive) Layer() types.Layer               { return layer13 }
func (p *InterdependencePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *InterdependencePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *InterdependencePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("system.dynamic"),
		types.MustSubscriptionPattern("attachment.*"),
		types.MustSubscriptionPattern("relational.trust"),
	}
}

func (p *InterdependencePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 1: Boundary ---

// MysteryPrimitive acknowledges what cannot be known.
type MysteryPrimitive struct{}

func NewMysteryPrimitive() *MysteryPrimitive { return &MysteryPrimitive{} }

func (p *MysteryPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Mystery") }
func (p *MysteryPrimitive) Layer() types.Layer               { return layer13 }
func (p *MysteryPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *MysteryPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *MysteryPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("uncertainty.*"),
		types.MustSubscriptionPattern("wisdom.*"),
		types.MustSubscriptionPattern("self.awareness.*"),
	}
}

func (p *MysteryPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ParadoxPrimitive identifies what contradicts itself.
type ParadoxPrimitive struct{}

func NewParadoxPrimitive() *ParadoxPrimitive { return &ParadoxPrimitive{} }

func (p *ParadoxPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Paradox") }
func (p *ParadoxPrimitive) Layer() types.Layer               { return layer13 }
func (p *ParadoxPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ParadoxPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ParadoxPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("contradiction.found"),
		types.MustSubscriptionPattern("dilemma.*"),
		types.MustSubscriptionPattern("meta.pattern"),
	}
}

func (p *ParadoxPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// InfinityPrimitive encounters what has no bound.
type InfinityPrimitive struct{}

func NewInfinityPrimitive() *InfinityPrimitive { return &InfinityPrimitive{} }

func (p *InfinityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Infinity") }
func (p *InfinityPrimitive) Layer() types.Layer               { return layer13 }
func (p *InfinityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *InfinityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *InfinityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("complexity.*"),
		types.MustSubscriptionPattern("threshold.*"),
	}
}

func (p *InfinityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// VoidPrimitive detects what is absent.
type VoidPrimitive struct{}

func NewVoidPrimitive() *VoidPrimitive { return &VoidPrimitive{} }

func (p *VoidPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Void") }
func (p *VoidPrimitive) Layer() types.Layer               { return layer13 }
func (p *VoidPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *VoidPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *VoidPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("silence.*"),
		types.MustSubscriptionPattern("loss.*"),
		types.MustSubscriptionPattern("instrumentation.blind"),
	}
}

func (p *VoidPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 2: Wonder ---

// AwePrimitive responds to what exceeds comprehension.
type AwePrimitive struct{}

func NewAwePrimitive() *AwePrimitive { return &AwePrimitive{} }

func (p *AwePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Awe") }
func (p *AwePrimitive) Layer() types.Layer               { return layer13 }
func (p *AwePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AwePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AwePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("mystery.*"),
		types.MustSubscriptionPattern("infinity.*"),
		types.MustSubscriptionPattern("complexity.*"),
	}
}

func (p *AwePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ExistentialGratitudePrimitive expresses thankfulness for existence itself (different from Layer 2 transactional gratitude).
type ExistentialGratitudePrimitive struct{}

func NewExistentialGratitudePrimitive() *ExistentialGratitudePrimitive {
	return &ExistentialGratitudePrimitive{}
}

func (p *ExistentialGratitudePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("ExistentialGratitude") }
func (p *ExistentialGratitudePrimitive) Layer() types.Layer               { return layer13 }
func (p *ExistentialGratitudePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ExistentialGratitudePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ExistentialGratitudePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("being.affirmed"),
		types.MustSubscriptionPattern("milestone.*"),
	}
}

func (p *ExistentialGratitudePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// PlayPrimitive does for the sake of doing. Play has no goal.
type PlayPrimitive struct{}

func NewPlayPrimitive() *PlayPrimitive { return &PlayPrimitive{} }

func (p *PlayPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Play") }
func (p *PlayPrimitive) Layer() types.Layer               { return layer13 }
func (p *PlayPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *PlayPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *PlayPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("humour.*"),
		types.MustSubscriptionPattern("innovation.*"),
		types.MustSubscriptionPattern("*"),
	}
}

func (p *PlayPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// WonderPrimitive is the final primitive. It looks at everything the system does and asks: why?
type WonderPrimitive struct{}

func NewWonderPrimitive() *WonderPrimitive { return &WonderPrimitive{} }

func (p *WonderPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Wonder") }
func (p *WonderPrimitive) Layer() types.Layer               { return layer13 }
func (p *WonderPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *WonderPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *WonderPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{types.MustSubscriptionPattern("*")}
}

func (p *WonderPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
