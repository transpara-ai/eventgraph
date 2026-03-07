// Package layer3 implements the Layer 3 Society primitives.
// Groups: Membership (Group, Role, Reputation, Exclusion),
// Collective Decision (Vote, Consensus, Dissent, Majority),
// Norms (Convention, Norm, Sanction, Forgiveness).
package layer3

import (
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer3 = types.MustLayer(3)
var cadence1 = types.MustCadence(1)

// --- Group 0: Membership ---

// GroupPrimitive forms and manages groups of actors.
type GroupPrimitive struct{}

func NewGroupPrimitive() *GroupPrimitive { return &GroupPrimitive{} }

func (p *GroupPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Group") }
func (p *GroupPrimitive) Layer() types.Layer               { return layer3 }
func (p *GroupPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *GroupPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *GroupPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("actor.*"),
		types.MustSubscriptionPattern("consent.*"),
	}
}

func (p *GroupPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RolePrimitive assigns roles within groups.
type RolePrimitive struct{}

func NewRolePrimitive() *RolePrimitive { return &RolePrimitive{} }

func (p *RolePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Role") }
func (p *RolePrimitive) Layer() types.Layer               { return layer3 }
func (p *RolePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RolePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RolePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("group.*"),
		types.MustSubscriptionPattern("delegation.*"),
	}
}

func (p *RolePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ReputationPrimitive tracks how the group perceives each member.
type ReputationPrimitive struct{}

func NewReputationPrimitive() *ReputationPrimitive { return &ReputationPrimitive{} }

func (p *ReputationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Reputation") }
func (p *ReputationPrimitive) Layer() types.Layer               { return layer3 }
func (p *ReputationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ReputationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ReputationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("trust.*"),
		types.MustSubscriptionPattern("commitment.*"),
		types.MustSubscriptionPattern("violation.*"),
		types.MustSubscriptionPattern("gratitude.*"),
	}
}

func (p *ReputationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	trustEvents := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "trust.") {
			trustEvents++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "trustEventsProcessed", Value: trustEvents},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ExclusionPrimitive handles when someone must leave a group.
type ExclusionPrimitive struct{}

func NewExclusionPrimitive() *ExclusionPrimitive { return &ExclusionPrimitive{} }

func (p *ExclusionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Exclusion") }
func (p *ExclusionPrimitive) Layer() types.Layer               { return layer3 }
func (p *ExclusionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ExclusionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ExclusionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("reputation.*"),
		types.MustSubscriptionPattern("violation.*"),
		types.MustSubscriptionPattern("quarantine.*"),
		types.MustSubscriptionPattern("dispute.*"),
	}
}

func (p *ExclusionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	violations := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "violation.") || strings.HasPrefix(ev.Type().Value(), "quarantine.") {
			violations++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "violationsObserved", Value: violations},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 1: Collective Decision ---

// VotePrimitive provides structured group decision-making.
type VotePrimitive struct{}

func NewVotePrimitive() *VotePrimitive { return &VotePrimitive{} }

func (p *VotePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Vote") }
func (p *VotePrimitive) Layer() types.Layer               { return layer3 }
func (p *VotePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *VotePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *VotePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("authority.requested"),
		types.MustSubscriptionPattern("group.*"),
	}
}

func (p *VotePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ConsensusPrimitive detects when a group naturally agrees.
type ConsensusPrimitive struct{}

func NewConsensusPrimitive() *ConsensusPrimitive { return &ConsensusPrimitive{} }

func (p *ConsensusPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Consensus") }
func (p *ConsensusPrimitive) Layer() types.Layer               { return layer3 }
func (p *ConsensusPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ConsensusPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ConsensusPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("message.*"),
		types.MustSubscriptionPattern("corroboration.*"),
		types.MustSubscriptionPattern("vote.result"),
	}
}

func (p *ConsensusPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// DissentPrimitive tracks and protects disagreement.
type DissentPrimitive struct{}

func NewDissentPrimitive() *DissentPrimitive { return &DissentPrimitive{} }

func (p *DissentPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Dissent") }
func (p *DissentPrimitive) Layer() types.Layer               { return layer3 }
func (p *DissentPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *DissentPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *DissentPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("vote.*"),
		types.MustSubscriptionPattern("consensus.*"),
		types.MustSubscriptionPattern("contradiction.found"),
	}
}

func (p *DissentPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// MajorityPrimitive handles the tyranny of the majority.
type MajorityPrimitive struct{}

func NewMajorityPrimitive() *MajorityPrimitive { return &MajorityPrimitive{} }

func (p *MajorityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Majority") }
func (p *MajorityPrimitive) Layer() types.Layer               { return layer3 }
func (p *MajorityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *MajorityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *MajorityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("vote.result"),
		types.MustSubscriptionPattern("dissent.*"),
		types.MustSubscriptionPattern("exclusion.*"),
	}
}

func (p *MajorityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 2: Norms ---

// ConventionPrimitive detects unwritten rules that emerge from behaviour.
type ConventionPrimitive struct{}

func NewConventionPrimitive() *ConventionPrimitive { return &ConventionPrimitive{} }

func (p *ConventionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Convention") }
func (p *ConventionPrimitive) Layer() types.Layer               { return layer3 }
func (p *ConventionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ConventionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ConventionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("pattern.detected"),
		types.MustSubscriptionPattern("*"),
	}
}

func (p *ConventionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// NormPrimitive manages explicit shared expectations with enforcement.
type NormPrimitive struct{}

func NewNormPrimitive() *NormPrimitive { return &NormPrimitive{} }

func (p *NormPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Norm") }
func (p *NormPrimitive) Layer() types.Layer               { return layer3 }
func (p *NormPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *NormPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *NormPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("convention.detected"),
		types.MustSubscriptionPattern("consensus.reached"),
	}
}

func (p *NormPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// SanctionPrimitive applies consequences for norm violations.
type SanctionPrimitive struct{}

func NewSanctionPrimitive() *SanctionPrimitive { return &SanctionPrimitive{} }

func (p *SanctionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Sanction") }
func (p *SanctionPrimitive) Layer() types.Layer               { return layer3 }
func (p *SanctionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *SanctionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *SanctionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("norm.violated"),
		types.MustSubscriptionPattern("violation.*"),
	}
}

func (p *SanctionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	violations := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "norm.") || strings.HasPrefix(ev.Type().Value(), "violation.") {
			violations++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "violationsProcessed", Value: violations},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ForgivenessPrimitive restores standing after violation and sanction.
type ForgivenessPrimitive struct{}

func NewForgivenessPrimitive() *ForgivenessPrimitive { return &ForgivenessPrimitive{} }

func (p *ForgivenessPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Forgiveness") }
func (p *ForgivenessPrimitive) Layer() types.Layer               { return layer3 }
func (p *ForgivenessPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ForgivenessPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ForgivenessPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("sanction.applied"),
		types.MustSubscriptionPattern("trust.*"),
		types.MustSubscriptionPattern("obligation.fulfilled"),
	}
}

func (p *ForgivenessPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
