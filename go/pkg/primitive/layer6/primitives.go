// Package layer6 implements the Layer 6 Information primitives.
// Groups: Representation (Symbol, Abstraction, Classification, Encoding),
// Knowledge (Fact, Inference, Memory, Learning),
// Truth (Narrative, Bias, Correction, Provenance).
package layer6

import (
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer6 = types.MustLayer(6)
var cadence1 = types.MustCadence(1)

// --- Group 0: Representation ---

// SymbolPrimitive creates and interprets symbolic representations.
type SymbolPrimitive struct{}

func NewSymbolPrimitive() *SymbolPrimitive { return &SymbolPrimitive{} }

func (p *SymbolPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Symbol") }
func (p *SymbolPrimitive) Layer() types.Layer               { return layer6 }
func (p *SymbolPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *SymbolPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *SymbolPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{types.MustSubscriptionPattern("*")}
}

func (p *SymbolPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AbstractionPrimitive generalises from specifics.
type AbstractionPrimitive struct{}

func NewAbstractionPrimitive() *AbstractionPrimitive { return &AbstractionPrimitive{} }

func (p *AbstractionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Abstraction") }
func (p *AbstractionPrimitive) Layer() types.Layer               { return layer6 }
func (p *AbstractionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AbstractionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AbstractionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("pattern.detected"),
		types.MustSubscriptionPattern("symbol.*"),
	}
}

func (p *AbstractionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ClassificationPrimitive organises information into categories.
type ClassificationPrimitive struct{}

func NewClassificationPrimitive() *ClassificationPrimitive { return &ClassificationPrimitive{} }

func (p *ClassificationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Classification") }
func (p *ClassificationPrimitive) Layer() types.Layer               { return layer6 }
func (p *ClassificationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ClassificationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ClassificationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{types.MustSubscriptionPattern("*")}
}

func (p *ClassificationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// EncodingPrimitive transforms information between representations.
type EncodingPrimitive struct{}

func NewEncodingPrimitive() *EncodingPrimitive { return &EncodingPrimitive{} }

func (p *EncodingPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Encoding") }
func (p *EncodingPrimitive) Layer() types.Layer               { return layer6 }
func (p *EncodingPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *EncodingPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *EncodingPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("symbol.*"),
		types.MustSubscriptionPattern("message.*"),
	}
}

func (p *EncodingPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 1: Knowledge ---

// FactPrimitive establishes verified claims.
type FactPrimitive struct{}

func NewFactPrimitive() *FactPrimitive { return &FactPrimitive{} }

func (p *FactPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Fact") }
func (p *FactPrimitive) Layer() types.Layer               { return layer6 }
func (p *FactPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *FactPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *FactPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("corroboration.*"),
		types.MustSubscriptionPattern("evidence.*"),
		types.MustSubscriptionPattern("confidence.*"),
	}
}

func (p *FactPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	evidenceCount := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "evidence.") || strings.HasPrefix(ev.Type().Value(), "corroboration.") {
			evidenceCount++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "evidenceProcessed", Value: evidenceCount},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// InferencePrimitive derives new knowledge from existing facts.
type InferencePrimitive struct{}

func NewInferencePrimitive() *InferencePrimitive { return &InferencePrimitive{} }

func (p *InferencePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Inference") }
func (p *InferencePrimitive) Layer() types.Layer               { return layer6 }
func (p *InferencePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *InferencePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *InferencePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("fact.*"),
		types.MustSubscriptionPattern("evidence.*"),
	}
}

func (p *InferencePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// MemoryPrimitive handles long-term knowledge retention and retrieval.
type MemoryPrimitive struct{}

func NewMemoryPrimitive() *MemoryPrimitive { return &MemoryPrimitive{} }

func (p *MemoryPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Memory") }
func (p *MemoryPrimitive) Layer() types.Layer               { return layer6 }
func (p *MemoryPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *MemoryPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *MemoryPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("fact.*"),
		types.MustSubscriptionPattern("inference.*"),
		types.MustSubscriptionPattern("abstraction.*"),
	}
}

func (p *MemoryPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// LearningPrimitive updates behaviour based on experience.
type LearningPrimitive struct{}

func NewLearningPrimitive() *LearningPrimitive { return &LearningPrimitive{} }

func (p *LearningPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Learning") }
func (p *LearningPrimitive) Layer() types.Layer               { return layer6 }
func (p *LearningPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *LearningPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *LearningPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("feedback.*"),
		types.MustSubscriptionPattern("test.*"),
		types.MustSubscriptionPattern("inference.*"),
	}
}

func (p *LearningPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 2: Truth ---

// NarrativePrimitive constructs coherent stories from events.
type NarrativePrimitive struct{}

func NewNarrativePrimitive() *NarrativePrimitive { return &NarrativePrimitive{} }

func (p *NarrativePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Narrative") }
func (p *NarrativePrimitive) Layer() types.Layer               { return layer6 }
func (p *NarrativePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *NarrativePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *NarrativePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("fact.*"),
		types.MustSubscriptionPattern("inference.*"),
		types.MustSubscriptionPattern("memory.*"),
	}
}

func (p *NarrativePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// BiasPrimitive detects systematic distortions in information.
type BiasPrimitive struct{}

func NewBiasPrimitive() *BiasPrimitive { return &BiasPrimitive{} }

func (p *BiasPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Bias") }
func (p *BiasPrimitive) Layer() types.Layer               { return layer6 }
func (p *BiasPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *BiasPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *BiasPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("narrative.*"),
		types.MustSubscriptionPattern("classification.*"),
		types.MustSubscriptionPattern("inference.*"),
	}
}

func (p *BiasPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// CorrectionPrimitive fixes errors in the knowledge base.
type CorrectionPrimitive struct{}

func NewCorrectionPrimitive() *CorrectionPrimitive { return &CorrectionPrimitive{} }

func (p *CorrectionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Correction") }
func (p *CorrectionPrimitive) Layer() types.Layer               { return layer6 }
func (p *CorrectionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *CorrectionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *CorrectionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("bias.detected"),
		types.MustSubscriptionPattern("fact.retracted"),
		types.MustSubscriptionPattern("contradiction.found"),
	}
}

func (p *CorrectionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	corrections := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "bias.") || strings.HasPrefix(ev.Type().Value(), "contradiction.") {
			corrections++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "correctionsNeeded", Value: corrections},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ProvenancePrimitive tracks the origin and chain of custody of information.
type ProvenancePrimitive struct{}

func NewProvenancePrimitive() *ProvenancePrimitive { return &ProvenancePrimitive{} }

func (p *ProvenancePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Provenance") }
func (p *ProvenancePrimitive) Layer() types.Layer               { return layer6 }
func (p *ProvenancePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ProvenancePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ProvenancePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("fact.*"),
		types.MustSubscriptionPattern("memory.*"),
		types.MustSubscriptionPattern("message.*"),
	}
}

func (p *ProvenancePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
