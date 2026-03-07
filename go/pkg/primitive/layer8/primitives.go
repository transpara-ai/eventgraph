// Package layer8 implements the Layer 8 Identity primitives.
// Groups: Self-Knowledge (SelfModel, Authenticity, NarrativeIdentity, Boundary),
// Continuity (Persistence, Transformation, Heritage, Aspiration),
// Recognition (Dignity, Acknowledgement, Uniqueness, Memorial).
package layer8

import (
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer8 = types.MustLayer(8)
var cadence1 = types.MustCadence(1)

// --- Group 0: Self-Knowledge ---

// SelfModelPrimitive maintains an actor's model of itself.
type SelfModelPrimitive struct{}

func NewSelfModelPrimitive() *SelfModelPrimitive { return &SelfModelPrimitive{} }

func (p *SelfModelPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("SelfModel") }
func (p *SelfModelPrimitive) Layer() types.Layer               { return layer8 }
func (p *SelfModelPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *SelfModelPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *SelfModelPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("commitment.*"),
		types.MustSubscriptionPattern("learning.*"),
		types.MustSubscriptionPattern("moral.growth"),
		types.MustSubscriptionPattern("capability.*"),
	}
}

func (p *SelfModelPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AuthenticityPrimitive assesses alignment between self-model and behaviour.
type AuthenticityPrimitive struct{}

func NewAuthenticityPrimitive() *AuthenticityPrimitive { return &AuthenticityPrimitive{} }

func (p *AuthenticityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Authenticity") }
func (p *AuthenticityPrimitive) Layer() types.Layer               { return layer8 }
func (p *AuthenticityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AuthenticityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AuthenticityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("self.model.*"),
		types.MustSubscriptionPattern("decision.*"),
		types.MustSubscriptionPattern("value.*"),
	}
}

func (p *AuthenticityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// NarrativeIdentityPrimitive maintains the story an actor tells about itself.
type NarrativeIdentityPrimitive struct{}

func NewNarrativeIdentityPrimitive() *NarrativeIdentityPrimitive { return &NarrativeIdentityPrimitive{} }

func (p *NarrativeIdentityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("NarrativeIdentity") }
func (p *NarrativeIdentityPrimitive) Layer() types.Layer               { return layer8 }
func (p *NarrativeIdentityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *NarrativeIdentityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *NarrativeIdentityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("self.model.*"),
		types.MustSubscriptionPattern("narrative.*"),
		types.MustSubscriptionPattern("memory.*"),
	}
}

func (p *NarrativeIdentityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// BoundaryPrimitive defines where one actor ends and another begins.
type BoundaryPrimitive struct{}

func NewBoundaryPrimitive() *BoundaryPrimitive { return &BoundaryPrimitive{} }

func (p *BoundaryPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Boundary") }
func (p *BoundaryPrimitive) Layer() types.Layer               { return layer8 }
func (p *BoundaryPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *BoundaryPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *BoundaryPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("delegation.*"),
		types.MustSubscriptionPattern("group.*"),
		types.MustSubscriptionPattern("consent.*"),
	}
}

func (p *BoundaryPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 1: Continuity ---

// PersistencePrimitive tracks what stays the same as everything else changes.
type PersistencePrimitive struct{}

func NewPersistencePrimitive() *PersistencePrimitive { return &PersistencePrimitive{} }

func (p *PersistencePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Persistence") }
func (p *PersistencePrimitive) Layer() types.Layer               { return layer8 }
func (p *PersistencePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *PersistencePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *PersistencePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("self.model.*"),
		types.MustSubscriptionPattern("learning.*"),
	}
}

func (p *PersistencePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// TransformationPrimitive detects fundamental changes in identity.
type TransformationPrimitive struct{}

func NewTransformationPrimitive() *TransformationPrimitive { return &TransformationPrimitive{} }

func (p *TransformationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Transformation") }
func (p *TransformationPrimitive) Layer() types.Layer               { return layer8 }
func (p *TransformationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *TransformationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *TransformationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("self.model.*"),
		types.MustSubscriptionPattern("moral.growth"),
		types.MustSubscriptionPattern("learning.*"),
	}
}

func (p *TransformationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// HeritagePrimitive recognises what came before — identity from history.
type HeritagePrimitive struct{}

func NewHeritagePrimitive() *HeritagePrimitive { return &HeritagePrimitive{} }

func (p *HeritagePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Heritage") }
func (p *HeritagePrimitive) Layer() types.Layer               { return layer8 }
func (p *HeritagePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *HeritagePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *HeritagePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("memory.*"),
		types.MustSubscriptionPattern("legacy.*"),
		types.MustSubscriptionPattern("provenance.*"),
	}
}

func (p *HeritagePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AspirationPrimitive tracks who the actor wants to become.
type AspirationPrimitive struct{}

func NewAspirationPrimitive() *AspirationPrimitive { return &AspirationPrimitive{} }

func (p *AspirationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Aspiration") }
func (p *AspirationPrimitive) Layer() types.Layer               { return layer8 }
func (p *AspirationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AspirationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AspirationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("self.model.*"),
		types.MustSubscriptionPattern("goal.*"),
		types.MustSubscriptionPattern("value.*"),
	}
}

func (p *AspirationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 2: Recognition ---

// DignityPrimitive affirms the inherent worth of every actor. The DIGNITY invariant flows through this.
type DignityPrimitive struct{}

func NewDignityPrimitive() *DignityPrimitive { return &DignityPrimitive{} }

func (p *DignityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Dignity") }
func (p *DignityPrimitive) Layer() types.Layer               { return layer8 }
func (p *DignityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *DignityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *DignityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("exclusion.*"),
		types.MustSubscriptionPattern("harm.*"),
		types.MustSubscriptionPattern("right.violated"),
		types.MustSubscriptionPattern("actor.memorial"),
	}
}

func (p *DignityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AcknowledgementPrimitive tracks being seen and recognised by others.
type AcknowledgementPrimitive struct{}

func NewAcknowledgementPrimitive() *AcknowledgementPrimitive { return &AcknowledgementPrimitive{} }

func (p *AcknowledgementPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("IdentityAcknowledgement") }
func (p *AcknowledgementPrimitive) Layer() types.Layer               { return layer8 }
func (p *AcknowledgementPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AcknowledgementPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AcknowledgementPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("message.*"),
		types.MustSubscriptionPattern("gratitude.*"),
		types.MustSubscriptionPattern("reputation.*"),
	}
}

func (p *AcknowledgementPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// UniquenessPrimitive identifies what makes each actor distinct.
type UniquenessPrimitive struct{}

func NewUniquenessPrimitive() *UniquenessPrimitive { return &UniquenessPrimitive{} }

func (p *UniquenessPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Uniqueness") }
func (p *UniquenessPrimitive) Layer() types.Layer               { return layer8 }
func (p *UniquenessPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *UniquenessPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *UniquenessPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("self.model.*"),
		types.MustSubscriptionPattern("identity.narrative"),
		types.MustSubscriptionPattern("pattern.detected"),
	}
}

func (p *UniquenessPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// MemorialPrimitive honours actors who have left.
type MemorialPrimitive struct{}

func NewMemorialPrimitive() *MemorialPrimitive { return &MemorialPrimitive{} }

func (p *MemorialPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Memorial") }
func (p *MemorialPrimitive) Layer() types.Layer               { return layer8 }
func (p *MemorialPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *MemorialPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *MemorialPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{types.MustSubscriptionPattern("actor.memorial")}
}

func (p *MemorialPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "memorialsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
