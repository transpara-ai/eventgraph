// Package layer8 implements the Layer 8 Identity primitives.
// Groups: SelfKnowledge (Narrative, SelfConcept, Reflection, Memory),
// SelfDirection (Purpose, Aspiration, Authenticity, Expression),
// SelfBecoming (Growth, Continuity, Integration, Crisis).
package layer8

import (
	"strings"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/primitive"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

var layer8 = types.MustLayer(8)
var cadence1 = types.MustCadence(1)

// --- Group A: Self-Knowledge ---

// NarrativePrimitive maintains the story an actor tells about itself over time.
type NarrativePrimitive struct{}

func NewNarrativePrimitive() *NarrativePrimitive { return &NarrativePrimitive{} }

func (p *NarrativePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Narrative") }
func (p *NarrativePrimitive) Layer() types.Layer               { return layer8 }
func (p *NarrativePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *NarrativePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *NarrativePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("identity.*"),
		types.MustSubscriptionPattern("memory.*"),
		types.MustSubscriptionPattern("reflection.*"),
	}
}

func (p *NarrativePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "identity.") || strings.HasPrefix(t, "memory.") || strings.HasPrefix(t, "reflection.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// SelfConceptPrimitive maintains an actor's model of what it is.
type SelfConceptPrimitive struct{}

func NewSelfConceptPrimitive() *SelfConceptPrimitive { return &SelfConceptPrimitive{} }

func (p *SelfConceptPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("SelfConcept") }
func (p *SelfConceptPrimitive) Layer() types.Layer               { return layer8 }
func (p *SelfConceptPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *SelfConceptPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *SelfConceptPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("identity.*"),
		types.MustSubscriptionPattern("capability.*"),
		types.MustSubscriptionPattern("value.*"),
	}
}

func (p *SelfConceptPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "identity.") || strings.HasPrefix(t, "capability.") || strings.HasPrefix(t, "value.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ReflectionPrimitive enables an actor to examine its own states and processes.
type ReflectionPrimitive struct{}

func NewReflectionPrimitive() *ReflectionPrimitive { return &ReflectionPrimitive{} }

func (p *ReflectionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Reflection") }
func (p *ReflectionPrimitive) Layer() types.Layer               { return layer8 }
func (p *ReflectionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ReflectionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ReflectionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("decision.*"),
		types.MustSubscriptionPattern("consequence.*"),
		types.MustSubscriptionPattern("learning.*"),
	}
}

func (p *ReflectionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "decision.") || strings.HasPrefix(t, "consequence.") || strings.HasPrefix(t, "learning.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// MemoryPrimitive tracks what an actor retains from past experience.
type MemoryPrimitive struct{}

func NewMemoryPrimitive() *MemoryPrimitive { return &MemoryPrimitive{} }

func (p *MemoryPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Memory") }
func (p *MemoryPrimitive) Layer() types.Layer               { return layer8 }
func (p *MemoryPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *MemoryPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *MemoryPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("experience.*"),
		types.MustSubscriptionPattern("narrative.*"),
		types.MustSubscriptionPattern("learning.*"),
	}
}

func (p *MemoryPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "experience.") || strings.HasPrefix(t, "narrative.") || strings.HasPrefix(t, "learning.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group B: Self-Direction ---

// PurposePrimitive tracks an actor's sense of why it exists and what it serves.
type PurposePrimitive struct{}

func NewPurposePrimitive() *PurposePrimitive { return &PurposePrimitive{} }

func (p *PurposePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Purpose") }
func (p *PurposePrimitive) Layer() types.Layer               { return layer8 }
func (p *PurposePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *PurposePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *PurposePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("goal.*"),
		types.MustSubscriptionPattern("value.*"),
		types.MustSubscriptionPattern("mission.*"),
	}
}

func (p *PurposePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "goal.") || strings.HasPrefix(t, "value.") || strings.HasPrefix(t, "mission.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
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
		types.MustSubscriptionPattern("purpose.*"),
		types.MustSubscriptionPattern("growth.*"),
		types.MustSubscriptionPattern("potential.*"),
	}
}

func (p *AspirationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "purpose.") || strings.HasPrefix(t, "growth.") || strings.HasPrefix(t, "potential.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AuthenticityPrimitive assesses alignment between self-concept and behaviour.
type AuthenticityPrimitive struct{}

func NewAuthenticityPrimitive() *AuthenticityPrimitive { return &AuthenticityPrimitive{} }

func (p *AuthenticityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Authenticity") }
func (p *AuthenticityPrimitive) Layer() types.Layer               { return layer8 }
func (p *AuthenticityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AuthenticityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AuthenticityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("self.concept.*"),
		types.MustSubscriptionPattern("decision.*"),
		types.MustSubscriptionPattern("expression.*"),
	}
}

func (p *AuthenticityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "self.concept.") || strings.HasPrefix(t, "decision.") || strings.HasPrefix(t, "expression.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ExpressionPrimitive handles how an actor manifests its identity outward.
type ExpressionPrimitive struct{}

func NewExpressionPrimitive() *ExpressionPrimitive { return &ExpressionPrimitive{} }

func (p *ExpressionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Expression") }
func (p *ExpressionPrimitive) Layer() types.Layer               { return layer8 }
func (p *ExpressionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ExpressionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ExpressionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("authenticity.*"),
		types.MustSubscriptionPattern("signal.*"),
		types.MustSubscriptionPattern("act.*"),
	}
}

func (p *ExpressionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "authenticity.") || strings.HasPrefix(t, "signal.") || strings.HasPrefix(t, "act.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group C: Self-Becoming ---

// GrowthPrimitive tracks an actor's development and maturation over time.
type GrowthPrimitive struct{}

func NewGrowthPrimitive() *GrowthPrimitive { return &GrowthPrimitive{} }

func (p *GrowthPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Growth") }
func (p *GrowthPrimitive) Layer() types.Layer               { return layer8 }
func (p *GrowthPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *GrowthPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *GrowthPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("learning.*"),
		types.MustSubscriptionPattern("moral.growth"),
		types.MustSubscriptionPattern("aspiration.*"),
	}
}

func (p *GrowthPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "learning.") || strings.HasPrefix(t, "moral.") || strings.HasPrefix(t, "aspiration.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ContinuityPrimitive tracks what persists as identity changes over time.
type ContinuityPrimitive struct{}

func NewContinuityPrimitive() *ContinuityPrimitive { return &ContinuityPrimitive{} }

func (p *ContinuityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Continuity") }
func (p *ContinuityPrimitive) Layer() types.Layer               { return layer8 }
func (p *ContinuityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ContinuityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ContinuityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("self.concept.*"),
		types.MustSubscriptionPattern("memory.*"),
		types.MustSubscriptionPattern("persistence.*"),
	}
}

func (p *ContinuityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "self.concept.") || strings.HasPrefix(t, "memory.") || strings.HasPrefix(t, "persistence.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// IntegrationPrimitive synthesises disparate aspects of identity into coherence.
type IntegrationPrimitive struct{}

func NewIntegrationPrimitive() *IntegrationPrimitive { return &IntegrationPrimitive{} }

func (p *IntegrationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Integration") }
func (p *IntegrationPrimitive) Layer() types.Layer               { return layer8 }
func (p *IntegrationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *IntegrationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *IntegrationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("narrative.*"),
		types.MustSubscriptionPattern("self.concept.*"),
		types.MustSubscriptionPattern("growth.*"),
	}
}

func (p *IntegrationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "narrative.") || strings.HasPrefix(t, "self.concept.") || strings.HasPrefix(t, "growth.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// CrisisPrimitive detects when identity is fundamentally threatened or disrupted.
type CrisisPrimitive struct{}

func NewCrisisPrimitive() *CrisisPrimitive { return &CrisisPrimitive{} }

func (p *CrisisPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Crisis") }
func (p *CrisisPrimitive) Layer() types.Layer               { return layer8 }
func (p *CrisisPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *CrisisPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *CrisisPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("continuity.*"),
		types.MustSubscriptionPattern("rupture.*"),
		types.MustSubscriptionPattern("harm.*"),
	}
}

func (p *CrisisPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "continuity.") || strings.HasPrefix(t, "rupture.") || strings.HasPrefix(t, "harm.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
