// Package layer12 implements the Layer 12 Emergence primitives.
// Groups: PrinciplesOfComplexity (Emergence, SelfOrganization, Feedback, Complexity),
// LimitsAndSelfReference (Consciousness, Recursion, Paradox, Incompleteness),
// DynamicArchitecture (PhaseTransition, DownwardCausation, Autopoiesis, CoEvolution).
package layer12

import (
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer12 = types.MustLayer(12)
var cadence1 = types.MustCadence(1)

// --- Group A: Principles of Complexity ---

// EmergencePrimitive detects properties arising from interactions that are not present in individual components.
type EmergencePrimitive struct{}

func NewEmergencePrimitive() *EmergencePrimitive { return &EmergencePrimitive{} }

func (p *EmergencePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Emergence") }
func (p *EmergencePrimitive) Layer() types.Layer               { return layer12 }
func (p *EmergencePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *EmergencePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *EmergencePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("feedback.*"),
		types.MustSubscriptionPattern("complexity.*"),
		types.MustSubscriptionPattern("pattern.detected"),
	}
}

func (p *EmergencePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "feedback.") || strings.HasPrefix(t, "complexity.") || strings.HasPrefix(t, "pattern.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// SelfOrganizationPrimitive detects order arising spontaneously from local interactions without central control.
type SelfOrganizationPrimitive struct{}

func NewSelfOrganizationPrimitive() *SelfOrganizationPrimitive { return &SelfOrganizationPrimitive{} }

func (p *SelfOrganizationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("SelfOrganization") }
func (p *SelfOrganizationPrimitive) Layer() types.Layer               { return layer12 }
func (p *SelfOrganizationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *SelfOrganizationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *SelfOrganizationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("emergence.*"),
		types.MustSubscriptionPattern("actor.*"),
		types.MustSubscriptionPattern("coordination.*"),
	}
}

func (p *SelfOrganizationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "emergence.") || strings.HasPrefix(t, "actor.") || strings.HasPrefix(t, "coordination.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// FeedbackPrimitive detects self-reinforcing or self-correcting cycles in the system.
type FeedbackPrimitive struct{}

func NewFeedbackPrimitive() *FeedbackPrimitive { return &FeedbackPrimitive{} }

func (p *FeedbackPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Feedback") }
func (p *FeedbackPrimitive) Layer() types.Layer               { return layer12 }
func (p *FeedbackPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *FeedbackPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *FeedbackPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("consequence.*"),
		types.MustSubscriptionPattern("system.dynamic"),
		types.MustSubscriptionPattern("pattern.detected"),
	}
}

func (p *FeedbackPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "consequence.") || strings.HasPrefix(t, "system.") || strings.HasPrefix(t, "pattern.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ComplexityPrimitive measures the irreducible richness of system behaviour beyond simple aggregation.
type ComplexityPrimitive struct{}

func NewComplexityPrimitive() *ComplexityPrimitive { return &ComplexityPrimitive{} }

func (p *ComplexityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Complexity") }
func (p *ComplexityPrimitive) Layer() types.Layer               { return layer12 }
func (p *ComplexityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ComplexityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ComplexityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("emergence.*"),
		types.MustSubscriptionPattern("feedback.*"),
		types.MustSubscriptionPattern("selforganization.*"),
	}
}

func (p *ComplexityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "emergence.") || strings.HasPrefix(t, "feedback.") || strings.HasPrefix(t, "selforganization.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group B: Limits and Self-Reference ---

// ConsciousnessPrimitive models the system's awareness of its own processes and states.
type ConsciousnessPrimitive struct{}

func NewConsciousnessPrimitive() *ConsciousnessPrimitive { return &ConsciousnessPrimitive{} }

func (p *ConsciousnessPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Consciousness") }
func (p *ConsciousnessPrimitive) Layer() types.Layer               { return layer12 }
func (p *ConsciousnessPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ConsciousnessPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ConsciousnessPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("self.model.*"),
		types.MustSubscriptionPattern("recursion.*"),
		types.MustSubscriptionPattern("health.*"),
	}
}

func (p *ConsciousnessPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "self.model.") || strings.HasPrefix(t, "recursion.") || strings.HasPrefix(t, "health.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RecursionPrimitive detects self-referential structures where processes operate on themselves.
type RecursionPrimitive struct{}

func NewRecursionPrimitive() *RecursionPrimitive { return &RecursionPrimitive{} }

func (p *RecursionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Recursion") }
func (p *RecursionPrimitive) Layer() types.Layer               { return layer12 }
func (p *RecursionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RecursionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RecursionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("consciousness.*"),
		types.MustSubscriptionPattern("meta.pattern"),
		types.MustSubscriptionPattern("self.model.*"),
	}
}

func (p *RecursionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "consciousness.") || strings.HasPrefix(t, "meta.") || strings.HasPrefix(t, "self.model.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ParadoxPrimitive identifies contradictions that cannot be resolved without changing the frame of reference.
type ParadoxPrimitive struct{}

func NewParadoxPrimitive() *ParadoxPrimitive { return &ParadoxPrimitive{} }

func (p *ParadoxPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Paradox") }
func (p *ParadoxPrimitive) Layer() types.Layer               { return layer12 }
func (p *ParadoxPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ParadoxPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ParadoxPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("contradiction.found"),
		types.MustSubscriptionPattern("recursion.*"),
		types.MustSubscriptionPattern("dilemma.*"),
	}
}

func (p *ParadoxPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "contradiction.") || strings.HasPrefix(t, "recursion.") || strings.HasPrefix(t, "dilemma.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// IncompletenessePrimitive recognises that no system can fully describe itself from within.
type IncompletenesPrimitive struct{}

func NewIncompletenesPrimitive() *IncompletenesPrimitive { return &IncompletenesPrimitive{} }

func (p *IncompletenesPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Incompleteness") }
func (p *IncompletenesPrimitive) Layer() types.Layer               { return layer12 }
func (p *IncompletenesPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *IncompletenesPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *IncompletenesPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("paradox.*"),
		types.MustSubscriptionPattern("recursion.*"),
		types.MustSubscriptionPattern("uncertainty.*"),
	}
}

func (p *IncompletenesPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "paradox.") || strings.HasPrefix(t, "recursion.") || strings.HasPrefix(t, "uncertainty.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group C: Dynamic Architecture ---

// PhaseTransitionPrimitive detects points where quantitative change becomes qualitative — regime shifts.
type PhaseTransitionPrimitive struct{}

func NewPhaseTransitionPrimitive() *PhaseTransitionPrimitive { return &PhaseTransitionPrimitive{} }

func (p *PhaseTransitionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("PhaseTransition") }
func (p *PhaseTransitionPrimitive) Layer() types.Layer               { return layer12 }
func (p *PhaseTransitionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *PhaseTransitionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *PhaseTransitionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("threshold.*"),
		types.MustSubscriptionPattern("complexity.*"),
		types.MustSubscriptionPattern("feedback.*"),
	}
}

func (p *PhaseTransitionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "threshold.") || strings.HasPrefix(t, "complexity.") || strings.HasPrefix(t, "feedback.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// DownwardCausationPrimitive detects higher-level patterns constraining lower-level components.
type DownwardCausationPrimitive struct{}

func NewDownwardCausationPrimitive() *DownwardCausationPrimitive {
	return &DownwardCausationPrimitive{}
}

func (p *DownwardCausationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("DownwardCausation") }
func (p *DownwardCausationPrimitive) Layer() types.Layer               { return layer12 }
func (p *DownwardCausationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *DownwardCausationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *DownwardCausationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("emergence.*"),
		types.MustSubscriptionPattern("constraint.*"),
		types.MustSubscriptionPattern("norm.*"),
	}
}

func (p *DownwardCausationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "emergence.") || strings.HasPrefix(t, "constraint.") || strings.HasPrefix(t, "norm.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AutopoiesisPrimitive detects self-producing systems that maintain their own organisation.
type AutopoiesisPrimitive struct{}

func NewAutopoiesisPrimitive() *AutopoiesisPrimitive { return &AutopoiesisPrimitive{} }

func (p *AutopoiesisPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Autopoiesis") }
func (p *AutopoiesisPrimitive) Layer() types.Layer               { return layer12 }
func (p *AutopoiesisPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AutopoiesisPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AutopoiesisPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("selforganization.*"),
		types.MustSubscriptionPattern("sustainability.*"),
		types.MustSubscriptionPattern("health.*"),
	}
}

func (p *AutopoiesisPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "selforganization.") || strings.HasPrefix(t, "sustainability.") || strings.HasPrefix(t, "health.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// CoEvolutionPrimitive tracks how coupled systems shape each other's development over time.
type CoEvolutionPrimitive struct{}

func NewCoEvolutionPrimitive() *CoEvolutionPrimitive { return &CoEvolutionPrimitive{} }

func (p *CoEvolutionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("CoEvolution") }
func (p *CoEvolutionPrimitive) Layer() types.Layer               { return layer12 }
func (p *CoEvolutionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *CoEvolutionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *CoEvolutionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("adaptation.*"),
		types.MustSubscriptionPattern("autopoiesis.*"),
		types.MustSubscriptionPattern("feedback.*"),
	}
}

func (p *CoEvolutionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "adaptation.") || strings.HasPrefix(t, "autopoiesis.") || strings.HasPrefix(t, "feedback.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
