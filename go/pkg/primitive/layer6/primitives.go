// Package layer6 implements the Layer 6 Information primitives.
// Groups: Representation (Symbol, Language, Encoding, Record),
// Dynamics (Channel, Copy, Noise, Redundancy),
// Transformation (Data, Computation, Algorithm, Entropy).
package layer6

import (
	"strings"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/primitive"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

var layer6 = types.MustLayer(6)
var cadence1 = types.MustCadence(1)

// --- Group A: Representation ---

// SymbolPrimitive represents a token that stands for something else by convention.
type SymbolPrimitive struct{}

func NewSymbolPrimitive() *SymbolPrimitive { return &SymbolPrimitive{} }

func (p *SymbolPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Symbol") }
func (p *SymbolPrimitive) Layer() types.Layer               { return layer6 }
func (p *SymbolPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *SymbolPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *SymbolPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("language.*"),
		types.MustSubscriptionPattern("encoding.*"),
	}
}

func (p *SymbolPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "language.") || strings.HasPrefix(t, "encoding.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// LanguagePrimitive represents a structured system of symbols with grammar and semantics.
type LanguagePrimitive struct{}

func NewLanguagePrimitive() *LanguagePrimitive { return &LanguagePrimitive{} }

func (p *LanguagePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Language") }
func (p *LanguagePrimitive) Layer() types.Layer               { return layer6 }
func (p *LanguagePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *LanguagePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *LanguagePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("symbol.*"),
		types.MustSubscriptionPattern("record.*"),
	}
}

func (p *LanguagePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "symbol.") || strings.HasPrefix(t, "record.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// EncodingPrimitive represents the mapping between information and its physical representation.
type EncodingPrimitive struct{}

func NewEncodingPrimitive() *EncodingPrimitive { return &EncodingPrimitive{} }

func (p *EncodingPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Encoding") }
func (p *EncodingPrimitive) Layer() types.Layer               { return layer6 }
func (p *EncodingPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *EncodingPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *EncodingPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("symbol.*"),
		types.MustSubscriptionPattern("channel.*"),
	}
}

func (p *EncodingPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "symbol.") || strings.HasPrefix(t, "channel.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RecordPrimitive represents information persisted for retrieval across time.
type RecordPrimitive struct{}

func NewRecordPrimitive() *RecordPrimitive { return &RecordPrimitive{} }

func (p *RecordPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Record") }
func (p *RecordPrimitive) Layer() types.Layer               { return layer6 }
func (p *RecordPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RecordPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RecordPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("encoding.*"),
		types.MustSubscriptionPattern("data.*"),
	}
}

func (p *RecordPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "encoding.") || strings.HasPrefix(t, "data.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group B: Dynamics ---

// ChannelPrimitive represents a medium through which information flows between actors.
type ChannelPrimitive struct{}

func NewChannelPrimitive() *ChannelPrimitive { return &ChannelPrimitive{} }

func (p *ChannelPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Channel") }
func (p *ChannelPrimitive) Layer() types.Layer               { return layer6 }
func (p *ChannelPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ChannelPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ChannelPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("signal.*"),
		types.MustSubscriptionPattern("noise.*"),
	}
}

func (p *ChannelPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "signal.") || strings.HasPrefix(t, "noise.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// CopyPrimitive represents the reproduction of information in a new location.
type CopyPrimitive struct{}

func NewCopyPrimitive() *CopyPrimitive { return &CopyPrimitive{} }

func (p *CopyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Copy") }
func (p *CopyPrimitive) Layer() types.Layer               { return layer6 }
func (p *CopyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *CopyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *CopyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("channel.*"),
		types.MustSubscriptionPattern("redundancy.*"),
	}
}

func (p *CopyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "channel.") || strings.HasPrefix(t, "redundancy.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// NoisePrimitive represents unwanted distortion introduced during transmission or storage.
type NoisePrimitive struct{}

func NewNoisePrimitive() *NoisePrimitive { return &NoisePrimitive{} }

func (p *NoisePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Noise") }
func (p *NoisePrimitive) Layer() types.Layer               { return layer6 }
func (p *NoisePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *NoisePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *NoisePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("channel.*"),
		types.MustSubscriptionPattern("entropy.*"),
	}
}

func (p *NoisePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "channel.") || strings.HasPrefix(t, "entropy.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RedundancyPrimitive represents intentional duplication that protects against loss.
type RedundancyPrimitive struct{}

func NewRedundancyPrimitive() *RedundancyPrimitive { return &RedundancyPrimitive{} }

func (p *RedundancyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Redundancy") }
func (p *RedundancyPrimitive) Layer() types.Layer               { return layer6 }
func (p *RedundancyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RedundancyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RedundancyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("noise.*"),
		types.MustSubscriptionPattern("copy.*"),
	}
}

func (p *RedundancyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "noise.") || strings.HasPrefix(t, "copy.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group C: Transformation ---

// DataPrimitive represents raw distinguishable differences before interpretation.
type DataPrimitive struct{}

func NewDataPrimitive() *DataPrimitive { return &DataPrimitive{} }

func (p *DataPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Data") }
func (p *DataPrimitive) Layer() types.Layer               { return layer6 }
func (p *DataPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *DataPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *DataPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("record.*"),
		types.MustSubscriptionPattern("encoding.*"),
	}
}

func (p *DataPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "record.") || strings.HasPrefix(t, "encoding.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ComputationPrimitive represents the systematic transformation of data according to rules.
type ComputationPrimitive struct{}

func NewComputationPrimitive() *ComputationPrimitive { return &ComputationPrimitive{} }

func (p *ComputationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Computation") }
func (p *ComputationPrimitive) Layer() types.Layer               { return layer6 }
func (p *ComputationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ComputationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ComputationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("algorithm.*"),
		types.MustSubscriptionPattern("data.*"),
	}
}

func (p *ComputationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "algorithm.") || strings.HasPrefix(t, "data.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AlgorithmPrimitive represents a finite, deterministic sequence of steps that solves a class of problems.
type AlgorithmPrimitive struct{}

func NewAlgorithmPrimitive() *AlgorithmPrimitive { return &AlgorithmPrimitive{} }

func (p *AlgorithmPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Algorithm") }
func (p *AlgorithmPrimitive) Layer() types.Layer               { return layer6 }
func (p *AlgorithmPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AlgorithmPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AlgorithmPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("computation.*"),
		types.MustSubscriptionPattern("data.*"),
	}
}

func (p *AlgorithmPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "computation.") || strings.HasPrefix(t, "data.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// EntropyPrimitive measures the uncertainty or disorder in an information source.
type EntropyPrimitive struct{}

func NewEntropyPrimitive() *EntropyPrimitive { return &EntropyPrimitive{} }

func (p *EntropyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Entropy") }
func (p *EntropyPrimitive) Layer() types.Layer               { return layer6 }
func (p *EntropyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *EntropyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *EntropyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("noise.*"),
		types.MustSubscriptionPattern("data.*"),
	}
}

func (p *EntropyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "noise.") || strings.HasPrefix(t, "data.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
