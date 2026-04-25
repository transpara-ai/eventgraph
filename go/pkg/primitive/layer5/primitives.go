// Package layer5 implements the Layer 5 Technology primitives.
// Groups: Investigation (Method, Measurement, Knowledge, Model),
// Creation (Tool, Technique, Invention, Abstraction),
// Systems (Infrastructure, Standard, Efficiency, Automation).
package layer5

import (
	"strings"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/primitive"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

var layer5 = types.MustLayer(5)
var cadence1 = types.MustCadence(1)

// --- Group A: Investigation ---

// MethodPrimitive represents a systematic procedure for acquiring knowledge or solving problems.
type MethodPrimitive struct{}

func NewMethodPrimitive() *MethodPrimitive { return &MethodPrimitive{} }

func (p *MethodPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Method") }
func (p *MethodPrimitive) Layer() types.Layer               { return layer5 }
func (p *MethodPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *MethodPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *MethodPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("knowledge.*"),
		types.MustSubscriptionPattern("measurement.*"),
	}
}

func (p *MethodPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "knowledge.") || strings.HasPrefix(t, "measurement.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// MeasurementPrimitive quantifies observable properties through systematic observation.
type MeasurementPrimitive struct{}

func NewMeasurementPrimitive() *MeasurementPrimitive { return &MeasurementPrimitive{} }

func (p *MeasurementPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Measurement") }
func (p *MeasurementPrimitive) Layer() types.Layer               { return layer5 }
func (p *MeasurementPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *MeasurementPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *MeasurementPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("method.*"),
		types.MustSubscriptionPattern("standard.*"),
	}
}

func (p *MeasurementPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "method.") || strings.HasPrefix(t, "standard.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// KnowledgePrimitive represents verified understanding produced by investigation.
type KnowledgePrimitive struct{}

func NewKnowledgePrimitive() *KnowledgePrimitive { return &KnowledgePrimitive{} }

func (p *KnowledgePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Knowledge") }
func (p *KnowledgePrimitive) Layer() types.Layer               { return layer5 }
func (p *KnowledgePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *KnowledgePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *KnowledgePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("method.*"),
		types.MustSubscriptionPattern("measurement.*"),
		types.MustSubscriptionPattern("model.*"),
	}
}

func (p *KnowledgePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "method.") || strings.HasPrefix(t, "measurement.") || strings.HasPrefix(t, "model.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ModelPrimitive represents a simplified representation that predicts or explains phenomena.
type ModelPrimitive struct{}

func NewModelPrimitive() *ModelPrimitive { return &ModelPrimitive{} }

func (p *ModelPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Model") }
func (p *ModelPrimitive) Layer() types.Layer               { return layer5 }
func (p *ModelPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ModelPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ModelPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("knowledge.*"),
		types.MustSubscriptionPattern("abstraction.*"),
	}
}

func (p *ModelPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "knowledge.") || strings.HasPrefix(t, "abstraction.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group B: Creation ---

// ToolPrimitive represents an artefact that extends an actor's capabilities.
type ToolPrimitive struct{}

func NewToolPrimitive() *ToolPrimitive { return &ToolPrimitive{} }

func (p *ToolPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Tool") }
func (p *ToolPrimitive) Layer() types.Layer               { return layer5 }
func (p *ToolPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ToolPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ToolPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("invention.*"),
		types.MustSubscriptionPattern("technique.*"),
	}
}

func (p *ToolPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "invention.") || strings.HasPrefix(t, "technique.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// TechniquePrimitive represents a learned procedure for applying tools effectively.
type TechniquePrimitive struct{}

func NewTechniquePrimitive() *TechniquePrimitive { return &TechniquePrimitive{} }

func (p *TechniquePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Technique") }
func (p *TechniquePrimitive) Layer() types.Layer               { return layer5 }
func (p *TechniquePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *TechniquePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *TechniquePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("method.*"),
		types.MustSubscriptionPattern("tool.*"),
	}
}

func (p *TechniquePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "method.") || strings.HasPrefix(t, "tool.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// InventionPrimitive represents a novel combination of knowledge and technique that creates something new.
type InventionPrimitive struct{}

func NewInventionPrimitive() *InventionPrimitive { return &InventionPrimitive{} }

func (p *InventionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Invention") }
func (p *InventionPrimitive) Layer() types.Layer               { return layer5 }
func (p *InventionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *InventionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *InventionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("knowledge.*"),
		types.MustSubscriptionPattern("technique.*"),
		types.MustSubscriptionPattern("tool.*"),
	}
}

func (p *InventionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "knowledge.") || strings.HasPrefix(t, "technique.") || strings.HasPrefix(t, "tool.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AbstractionPrimitive represents the removal of inessential detail to reveal reusable structure.
type AbstractionPrimitive struct{}

func NewAbstractionPrimitive() *AbstractionPrimitive { return &AbstractionPrimitive{} }

func (p *AbstractionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Abstraction") }
func (p *AbstractionPrimitive) Layer() types.Layer               { return layer5 }
func (p *AbstractionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AbstractionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AbstractionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("model.*"),
		types.MustSubscriptionPattern("knowledge.*"),
	}
}

func (p *AbstractionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "model.") || strings.HasPrefix(t, "knowledge.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group C: Systems ---

// InfrastructurePrimitive represents the shared foundation that tools and techniques depend on.
type InfrastructurePrimitive struct{}

func NewInfrastructurePrimitive() *InfrastructurePrimitive { return &InfrastructurePrimitive{} }

func (p *InfrastructurePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Infrastructure") }
func (p *InfrastructurePrimitive) Layer() types.Layer               { return layer5 }
func (p *InfrastructurePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *InfrastructurePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *InfrastructurePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("tool.*"),
		types.MustSubscriptionPattern("standard.*"),
	}
}

func (p *InfrastructurePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "tool.") || strings.HasPrefix(t, "standard.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// StandardPrimitive represents an agreed-upon specification that enables interoperability.
type StandardPrimitive struct{}

func NewStandardPrimitive() *StandardPrimitive { return &StandardPrimitive{} }

func (p *StandardPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Standard") }
func (p *StandardPrimitive) Layer() types.Layer               { return layer5 }
func (p *StandardPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *StandardPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *StandardPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("measurement.*"),
		types.MustSubscriptionPattern("infrastructure.*"),
	}
}

func (p *StandardPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "measurement.") || strings.HasPrefix(t, "infrastructure.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// EfficiencyPrimitive measures the ratio of useful output to total input for a process.
type EfficiencyPrimitive struct{}

func NewEfficiencyPrimitive() *EfficiencyPrimitive { return &EfficiencyPrimitive{} }

func (p *EfficiencyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Efficiency") }
func (p *EfficiencyPrimitive) Layer() types.Layer               { return layer5 }
func (p *EfficiencyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *EfficiencyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *EfficiencyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("measurement.*"),
		types.MustSubscriptionPattern("automation.*"),
	}
}

func (p *EfficiencyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "measurement.") || strings.HasPrefix(t, "automation.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AutomationPrimitive converts repeatable manual processes into self-executing ones.
type AutomationPrimitive struct{}

func NewAutomationPrimitive() *AutomationPrimitive { return &AutomationPrimitive{} }

func (p *AutomationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Automation") }
func (p *AutomationPrimitive) Layer() types.Layer               { return layer5 }
func (p *AutomationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AutomationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AutomationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("technique.*"),
		types.MustSubscriptionPattern("efficiency.*"),
	}
}

func (p *AutomationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "technique.") || strings.HasPrefix(t, "efficiency.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
