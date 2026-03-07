// Package layer1 implements the Layer 1 Agency primitives.
// Groups: Intention (Goal, Plan, Initiative, Commitment),
// Attention (Focus, Filter, Salience, Distraction),
// Autonomy (Permission, Capability, Delegation, Accountability).
package layer1

import (
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer1 = types.MustLayer(1)
var cadence1 = types.MustCadence(1)

// --- Group 0: Intention ---

// GoalPrimitive sets and tracks objectives for actors.
type GoalPrimitive struct{}

func NewGoalPrimitive() *GoalPrimitive { return &GoalPrimitive{} }

func (p *GoalPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Goal") }
func (p *GoalPrimitive) Layer() types.Layer               { return layer1 }
func (p *GoalPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *GoalPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *GoalPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("decision.*"),
		types.MustSubscriptionPattern("authority.resolved"),
		types.MustSubscriptionPattern("actor.*"),
	}
}

func (p *GoalPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	var mutations []primitive.Mutation
	goalCount := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "decision.") || strings.HasPrefix(t, "actor.") {
			goalCount++
		}
	}
	mutations = append(mutations,
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "goalEventsProcessed", Value: goalCount},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	)
	return mutations, nil
}

// PlanPrimitive decomposes goals into steps.
type PlanPrimitive struct{}

func NewPlanPrimitive() *PlanPrimitive { return &PlanPrimitive{} }

func (p *PlanPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Plan") }
func (p *PlanPrimitive) Layer() types.Layer               { return layer1 }
func (p *PlanPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *PlanPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *PlanPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("goal.*"),
	}
}

func (p *PlanPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	planEvents := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "goal.") {
			planEvents++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "goalEventsReceived", Value: planEvents},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// InitiativePrimitive decides when to act without being asked.
type InitiativePrimitive struct{}

func NewInitiativePrimitive() *InitiativePrimitive { return &InitiativePrimitive{} }

func (p *InitiativePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Initiative") }
func (p *InitiativePrimitive) Layer() types.Layer               { return layer1 }
func (p *InitiativePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *InitiativePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *InitiativePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("clock.tick"),
		types.MustSubscriptionPattern("goal.*"),
		types.MustSubscriptionPattern("plan.*"),
	}
}

func (p *InitiativePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	triggers := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "goal.") || strings.HasPrefix(t, "plan.") {
			triggers++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "triggerCount", Value: triggers},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// CommitmentPrimitive tracks whether actors follow through on goals.
type CommitmentPrimitive struct{}

func NewCommitmentPrimitive() *CommitmentPrimitive { return &CommitmentPrimitive{} }

func (p *CommitmentPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Commitment") }
func (p *CommitmentPrimitive) Layer() types.Layer               { return layer1 }
func (p *CommitmentPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *CommitmentPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *CommitmentPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("goal.set"),
		types.MustSubscriptionPattern("goal.achieved"),
		types.MustSubscriptionPattern("goal.abandoned"),
		types.MustSubscriptionPattern("plan.step.completed"),
	}
}

func (p *CommitmentPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	achieved := 0
	abandoned := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "goal.achieved") {
			achieved++
		} else if strings.HasPrefix(t, "goal.abandoned") {
			abandoned++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "achieved", Value: achieved},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "abandoned", Value: abandoned},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 1: Attention ---

// FocusPrimitive directs processing resources to high-priority events.
type FocusPrimitive struct{}

func NewFocusPrimitive() *FocusPrimitive { return &FocusPrimitive{} }

func (p *FocusPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Focus") }
func (p *FocusPrimitive) Layer() types.Layer               { return layer1 }
func (p *FocusPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *FocusPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *FocusPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{types.MustSubscriptionPattern("*")}
}

func (p *FocusPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsInFocus", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// FilterPrimitive suppresses noise by deciding what NOT to process.
type FilterPrimitive struct{}

func NewFilterPrimitive() *FilterPrimitive { return &FilterPrimitive{} }

func (p *FilterPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Filter") }
func (p *FilterPrimitive) Layer() types.Layer               { return layer1 }
func (p *FilterPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *FilterPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *FilterPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{types.MustSubscriptionPattern("*")}
}

func (p *FilterPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	suppressed := 0
	for range events {
		// Mechanical: no filtering rules yet — count all events as passed
		_ = suppressed
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "totalEvents", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "suppressedCount", Value: suppressed},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// SaliencePrimitive detects what matters in the current context.
type SaliencePrimitive struct{}

func NewSaliencePrimitive() *SaliencePrimitive { return &SaliencePrimitive{} }

func (p *SaliencePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Salience") }
func (p *SaliencePrimitive) Layer() types.Layer               { return layer1 }
func (p *SaliencePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *SaliencePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *SaliencePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{types.MustSubscriptionPattern("*")}
}

func (p *SaliencePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsScored", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// DistractionPrimitive detects when attention is pulled away from goals.
type DistractionPrimitive struct{}

func NewDistractionPrimitive() *DistractionPrimitive { return &DistractionPrimitive{} }

func (p *DistractionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Distraction") }
func (p *DistractionPrimitive) Layer() types.Layer               { return layer1 }
func (p *DistractionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *DistractionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *DistractionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("focus.*"),
		types.MustSubscriptionPattern("goal.*"),
	}
}

func (p *DistractionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	focusShifts := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "focus.") {
			focusShifts++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "focusShifts", Value: focusShifts},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 2: Autonomy ---

// PermissionPrimitive requests and tracks permissions for actions.
type PermissionPrimitive struct{}

func NewPermissionPrimitive() *PermissionPrimitive { return &PermissionPrimitive{} }

func (p *PermissionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Permission") }
func (p *PermissionPrimitive) Layer() types.Layer               { return layer1 }
func (p *PermissionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *PermissionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *PermissionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("authority.*"),
		types.MustSubscriptionPattern("decision.*"),
	}
}

func (p *PermissionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	authorityEvents := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "authority.") {
			authorityEvents++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "authorityEventsProcessed", Value: authorityEvents},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// CapabilityPrimitive tracks what an actor can do.
type CapabilityPrimitive struct{}

func NewCapabilityPrimitive() *CapabilityPrimitive { return &CapabilityPrimitive{} }

func (p *CapabilityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Capability") }
func (p *CapabilityPrimitive) Layer() types.Layer               { return layer1 }
func (p *CapabilityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *CapabilityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *CapabilityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("actor.registered"),
		types.MustSubscriptionPattern("permission.*"),
		types.MustSubscriptionPattern("trust.*"),
	}
}

func (p *CapabilityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	actorEvents := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "actor.") || strings.HasPrefix(ev.Type().Value(), "permission.") {
			actorEvents++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "capabilityEventsProcessed", Value: actorEvents},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// DelegationPrimitive assigns tasks or authority to others.
type DelegationPrimitive struct{}

func NewDelegationPrimitive() *DelegationPrimitive { return &DelegationPrimitive{} }

func (p *DelegationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Delegation") }
func (p *DelegationPrimitive) Layer() types.Layer               { return layer1 }
func (p *DelegationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *DelegationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *DelegationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("authority.*"),
		types.MustSubscriptionPattern("edge.created"),
	}
}

func (p *DelegationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	delegations := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "edge.created") {
			delegations++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "delegationEvents", Value: delegations},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AccountabilityPrimitive traces who is responsible when things go wrong.
type AccountabilityPrimitive struct{}

func NewAccountabilityPrimitive() *AccountabilityPrimitive { return &AccountabilityPrimitive{} }

func (p *AccountabilityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Accountability") }
func (p *AccountabilityPrimitive) Layer() types.Layer               { return layer1 }
func (p *AccountabilityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AccountabilityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AccountabilityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("delegation.*"),
		types.MustSubscriptionPattern("violation.*"),
		types.MustSubscriptionPattern("goal.abandoned"),
	}
}

func (p *AccountabilityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	violations := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "violation.") {
			violations++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "violationsTracked", Value: violations},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
