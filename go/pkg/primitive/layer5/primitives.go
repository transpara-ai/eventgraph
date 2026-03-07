// Package layer5 implements the Layer 5 Technology primitives.
// Groups: Artefact (Create, Tool, Quality, Deprecation),
// Process (Workflow, Automation, Testing, Review),
// Improvement (Feedback, Iteration, Innovation, Legacy).
package layer5

import (
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer5 = types.MustLayer(5)
var cadence1 = types.MustCadence(1)

// --- Group 0: Artefact ---

// CreatePrimitive handles making new things.
type CreatePrimitive struct{}

func NewCreatePrimitive() *CreatePrimitive { return &CreatePrimitive{} }

func (p *CreatePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Create") }
func (p *CreatePrimitive) Layer() types.Layer               { return layer5 }
func (p *CreatePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *CreatePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *CreatePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("plan.*"),
		types.MustSubscriptionPattern("goal.*"),
	}
}

func (p *CreatePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ToolPrimitive manages artefacts that extend actor capabilities.
type ToolPrimitive struct{}

func NewToolPrimitive() *ToolPrimitive { return &ToolPrimitive{} }

func (p *ToolPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Tool") }
func (p *ToolPrimitive) Layer() types.Layer               { return layer5 }
func (p *ToolPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ToolPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ToolPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("artefact.created"),
		types.MustSubscriptionPattern("capability.*"),
	}
}

func (p *ToolPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// QualityPrimitive assesses how well something was made.
type QualityPrimitive struct{}

func NewQualityPrimitive() *QualityPrimitive { return &QualityPrimitive{} }

func (p *QualityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Quality") }
func (p *QualityPrimitive) Layer() types.Layer               { return layer5 }
func (p *QualityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *QualityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *QualityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("artefact.*"),
		types.MustSubscriptionPattern("tool.used"),
	}
}

func (p *QualityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// DeprecationPrimitive handles when artefacts should no longer be used.
type DeprecationPrimitive struct{}

func NewDeprecationPrimitive() *DeprecationPrimitive { return &DeprecationPrimitive{} }

func (p *DeprecationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Deprecation") }
func (p *DeprecationPrimitive) Layer() types.Layer               { return layer5 }
func (p *DeprecationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *DeprecationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *DeprecationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("quality.*"),
		types.MustSubscriptionPattern("artefact.version"),
	}
}

func (p *DeprecationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 1: Process ---

// WorkflowPrimitive defines repeatable processes.
type WorkflowPrimitive struct{}

func NewWorkflowPrimitive() *WorkflowPrimitive { return &WorkflowPrimitive{} }

func (p *WorkflowPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Workflow") }
func (p *WorkflowPrimitive) Layer() types.Layer               { return layer5 }
func (p *WorkflowPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *WorkflowPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *WorkflowPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("plan.*"),
		types.MustSubscriptionPattern("convention.detected"),
	}
}

func (p *WorkflowPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AutomationPrimitive converts manual workflows to mechanical ones.
type AutomationPrimitive struct{}

func NewAutomationPrimitive() *AutomationPrimitive { return &AutomationPrimitive{} }

func (p *AutomationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Automation") }
func (p *AutomationPrimitive) Layer() types.Layer               { return layer5 }
func (p *AutomationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AutomationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AutomationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("workflow.executed"),
		types.MustSubscriptionPattern("pattern.detected"),
	}
}

func (p *AutomationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	patterns := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "pattern.") {
			patterns++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "patternsDetected", Value: patterns},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// TestingPrimitive verifies that artefacts and processes work correctly.
type TestingPrimitive struct{}

func NewTestingPrimitive() *TestingPrimitive { return &TestingPrimitive{} }

func (p *TestingPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Testing") }
func (p *TestingPrimitive) Layer() types.Layer               { return layer5 }
func (p *TestingPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *TestingPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *TestingPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("artefact.*"),
		types.MustSubscriptionPattern("workflow.*"),
		types.MustSubscriptionPattern("automation.*"),
	}
}

func (p *TestingPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ReviewPrimitive provides peer assessment of artefacts and decisions.
type ReviewPrimitive struct{}

func NewReviewPrimitive() *ReviewPrimitive { return &ReviewPrimitive{} }

func (p *ReviewPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Review") }
func (p *ReviewPrimitive) Layer() types.Layer               { return layer5 }
func (p *ReviewPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ReviewPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ReviewPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("artefact.*"),
		types.MustSubscriptionPattern("decision.*"),
	}
}

func (p *ReviewPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 2: Improvement ---

// FeedbackPrimitive handles structured input on outcomes.
type FeedbackPrimitive struct{}

func NewFeedbackPrimitive() *FeedbackPrimitive { return &FeedbackPrimitive{} }

func (p *FeedbackPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Feedback") }
func (p *FeedbackPrimitive) Layer() types.Layer               { return layer5 }
func (p *FeedbackPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *FeedbackPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *FeedbackPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{types.MustSubscriptionPattern("*")}
}

func (p *FeedbackPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// IterationPrimitive improves through repeated cycles.
type IterationPrimitive struct{}

func NewIterationPrimitive() *IterationPrimitive { return &IterationPrimitive{} }

func (p *IterationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Iteration") }
func (p *IterationPrimitive) Layer() types.Layer               { return layer5 }
func (p *IterationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *IterationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *IterationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("feedback.*"),
		types.MustSubscriptionPattern("test.*"),
		types.MustSubscriptionPattern("review.*"),
	}
}

func (p *IterationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// InnovationPrimitive detects genuinely new creations.
type InnovationPrimitive struct{}

func NewInnovationPrimitive() *InnovationPrimitive { return &InnovationPrimitive{} }

func (p *InnovationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Innovation") }
func (p *InnovationPrimitive) Layer() types.Layer               { return layer5 }
func (p *InnovationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *InnovationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *InnovationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("artefact.*"),
		types.MustSubscriptionPattern("pattern.detected"),
	}
}

func (p *InnovationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// LegacyPrimitive assesses what persists after an artefact is deprecated.
type LegacyPrimitive struct{}

func NewLegacyPrimitive() *LegacyPrimitive { return &LegacyPrimitive{} }

func (p *LegacyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Legacy") }
func (p *LegacyPrimitive) Layer() types.Layer               { return layer5 }
func (p *LegacyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *LegacyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *LegacyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("deprecation.*"),
		types.MustSubscriptionPattern("artefact.*"),
	}
}

func (p *LegacyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
