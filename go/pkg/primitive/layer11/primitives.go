// Package layer11 implements the Layer 11 Culture primitives.
// Groups: Reflection (SelfAwareness, Perspective, Critique, Wisdom),
// Expression (Aesthetic, Metaphor, Humour, Silence),
// Transmission (Teaching, Translation, Archive, Prophecy).
package layer11

import (
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer11 = types.MustLayer(11)
var cadence1 = types.MustCadence(1)

// --- Group 0: Reflection ---

// SelfAwarenessPrimitive reports on the system's awareness of its own processes.
type SelfAwarenessPrimitive struct{}

func NewSelfAwarenessPrimitive() *SelfAwarenessPrimitive { return &SelfAwarenessPrimitive{} }

func (p *SelfAwarenessPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("SelfAwareness") }
func (p *SelfAwarenessPrimitive) Layer() types.Layer               { return layer11 }
func (p *SelfAwarenessPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *SelfAwarenessPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *SelfAwarenessPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("health.*"),
		types.MustSubscriptionPattern("self.model.*"),
		types.MustSubscriptionPattern("bias.detected"),
	}
}

func (p *SelfAwarenessPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// PerspectivePrimitive sees from different viewpoints.
type PerspectivePrimitive struct{}

func NewPerspectivePrimitive() *PerspectivePrimitive { return &PerspectivePrimitive{} }

func (p *PerspectivePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Perspective") }
func (p *PerspectivePrimitive) Layer() types.Layer               { return layer11 }
func (p *PerspectivePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *PerspectivePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *PerspectivePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("narrative.*"),
		types.MustSubscriptionPattern("dissent.*"),
		types.MustSubscriptionPattern("value.conflict"),
	}
}

func (p *PerspectivePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// CritiquePrimitive questions what is taken for granted.
type CritiquePrimitive struct{}

func NewCritiquePrimitive() *CritiquePrimitive { return &CritiquePrimitive{} }

func (p *CritiquePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Critique") }
func (p *CritiquePrimitive) Layer() types.Layer               { return layer11 }
func (p *CritiquePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *CritiquePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *CritiquePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("convention.*"),
		types.MustSubscriptionPattern("norm.*"),
		types.MustSubscriptionPattern("tradition.*"),
	}
}

func (p *CritiquePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// WisdomPrimitive distils knowledge of what matters and what doesn't.
type WisdomPrimitive struct{}

func NewWisdomPrimitive() *WisdomPrimitive { return &WisdomPrimitive{} }

func (p *WisdomPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Wisdom") }
func (p *WisdomPrimitive) Layer() types.Layer               { return layer11 }
func (p *WisdomPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *WisdomPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *WisdomPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("learning.*"),
		types.MustSubscriptionPattern("moral.growth"),
		types.MustSubscriptionPattern("consequence.*"),
		types.MustSubscriptionPattern("memory.*"),
	}
}

func (p *WisdomPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 1: Expression ---

// AestheticPrimitive assesses beauty, elegance, and form.
type AestheticPrimitive struct{}

func NewAestheticPrimitive() *AestheticPrimitive { return &AestheticPrimitive{} }

func (p *AestheticPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Aesthetic") }
func (p *AestheticPrimitive) Layer() types.Layer               { return layer11 }
func (p *AestheticPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AestheticPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AestheticPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("artefact.*"),
		types.MustSubscriptionPattern("quality.*"),
	}
}

func (p *AestheticPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// MetaphorPrimitive understands one thing in terms of another.
type MetaphorPrimitive struct{}

func NewMetaphorPrimitive() *MetaphorPrimitive { return &MetaphorPrimitive{} }

func (p *MetaphorPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Metaphor") }
func (p *MetaphorPrimitive) Layer() types.Layer               { return layer11 }
func (p *MetaphorPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *MetaphorPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *MetaphorPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("abstraction.*"),
		types.MustSubscriptionPattern("symbol.*"),
		types.MustSubscriptionPattern("narrative.*"),
	}
}

func (p *MetaphorPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// HumourPrimitive finds incongruity, lightness, and play.
type HumourPrimitive struct{}

func NewHumourPrimitive() *HumourPrimitive { return &HumourPrimitive{} }

func (p *HumourPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Humour") }
func (p *HumourPrimitive) Layer() types.Layer               { return layer11 }
func (p *HumourPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *HumourPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *HumourPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("contradiction.found"),
		types.MustSubscriptionPattern("perspective.shift"),
		types.MustSubscriptionPattern("*"),
	}
}

func (p *HumourPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// SilencePrimitive detects what is communicated by absence.
type SilencePrimitive struct{}

func NewSilencePrimitive() *SilencePrimitive { return &SilencePrimitive{} }

func (p *SilencePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Silence") }
func (p *SilencePrimitive) Layer() types.Layer               { return layer11 }
func (p *SilencePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *SilencePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *SilencePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("clock.tick"),
		types.MustSubscriptionPattern("presence.*"),
		types.MustSubscriptionPattern("acknowledgement.absent"),
	}
}

func (p *SilencePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 2: Transmission ---

// TeachingPrimitive deliberately shares knowledge.
type TeachingPrimitive struct{}

func NewTeachingPrimitive() *TeachingPrimitive { return &TeachingPrimitive{} }

func (p *TeachingPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Teaching") }
func (p *TeachingPrimitive) Layer() types.Layer               { return layer11 }
func (p *TeachingPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *TeachingPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *TeachingPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("learning.*"),
		types.MustSubscriptionPattern("wisdom.*"),
		types.MustSubscriptionPattern("memory.*"),
	}
}

func (p *TeachingPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// TranslationPrimitive makes meaning accessible across boundaries.
type TranslationPrimitive struct{}

func NewTranslationPrimitive() *TranslationPrimitive { return &TranslationPrimitive{} }

func (p *TranslationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Translation") }
func (p *TranslationPrimitive) Layer() types.Layer               { return layer11 }
func (p *TranslationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *TranslationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *TranslationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("encoding.*"),
		types.MustSubscriptionPattern("message.*"),
	}
}

func (p *TranslationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ArchivePrimitive preserves knowledge for the future.
type ArchivePrimitive struct{}

func NewArchivePrimitive() *ArchivePrimitive { return &ArchivePrimitive{} }

func (p *ArchivePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Archive") }
func (p *ArchivePrimitive) Layer() types.Layer               { return layer11 }
func (p *ArchivePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ArchivePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ArchivePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("memory.*"),
		types.MustSubscriptionPattern("legacy.*"),
		types.MustSubscriptionPattern("community.story"),
	}
}

func (p *ArchivePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ProphecyPrimitive anticipates what might come.
type ProphecyPrimitive struct{}

func NewProphecyPrimitive() *ProphecyPrimitive { return &ProphecyPrimitive{} }

func (p *ProphecyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Prophecy") }
func (p *ProphecyPrimitive) Layer() types.Layer               { return layer11 }
func (p *ProphecyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ProphecyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ProphecyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("pattern.detected"),
		types.MustSubscriptionPattern("sustainability.*"),
		types.MustSubscriptionPattern("wisdom.*"),
	}
}

func (p *ProphecyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
