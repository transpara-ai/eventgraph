// Package layer3 implements the Layer 3 Society primitives.
// Groups: CollectiveIdentity (Group, Membership, Role, Consent),
// SocialOrder (Norm, Reputation, Sanction, Authority),
// CollectiveAgency (Property, Commons, Governance, CollectiveAct).
package layer3

import (
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer3 = types.MustLayer(3)
var cadence1 = types.MustCadence(1)

// --- Group A: Collective Identity ---

// GroupPrimitive forms and manages groups of actors.
type GroupPrimitive struct{}

func NewGroupPrimitive() *GroupPrimitive { return &GroupPrimitive{} }

func (p *GroupPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Group") }
func (p *GroupPrimitive) Layer() types.Layer              { return layer3 }
func (p *GroupPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *GroupPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *GroupPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("actor.*"),
		types.MustSubscriptionPattern("membership.*"),
	}
}

func (p *GroupPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "actor.") || strings.HasPrefix(t, "membership.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// MembershipPrimitive tracks which actors belong to which groups.
type MembershipPrimitive struct{}

func NewMembershipPrimitive() *MembershipPrimitive { return &MembershipPrimitive{} }

func (p *MembershipPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Membership") }
func (p *MembershipPrimitive) Layer() types.Layer              { return layer3 }
func (p *MembershipPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *MembershipPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *MembershipPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("group.*"),
		types.MustSubscriptionPattern("consent.*"),
	}
}

func (p *MembershipPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "group.") || strings.HasPrefix(t, "consent.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RolePrimitive assigns roles within groups.
type RolePrimitive struct{}

func NewRolePrimitive() *RolePrimitive { return &RolePrimitive{} }

func (p *RolePrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Role") }
func (p *RolePrimitive) Layer() types.Layer              { return layer3 }
func (p *RolePrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *RolePrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *RolePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("group.*"),
		types.MustSubscriptionPattern("membership.*"),
		types.MustSubscriptionPattern("delegation.*"),
	}
}

func (p *RolePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "group.") || strings.HasPrefix(t, "membership.") || strings.HasPrefix(t, "delegation.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ConsentPrimitive ensures explicit agreement from all affected parties.
type ConsentPrimitive struct{}

func NewConsentPrimitive() *ConsentPrimitive { return &ConsentPrimitive{} }

func (p *ConsentPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Consent") }
func (p *ConsentPrimitive) Layer() types.Layer              { return layer3 }
func (p *ConsentPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *ConsentPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *ConsentPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("membership.*"),
		types.MustSubscriptionPattern("authority.*"),
		types.MustSubscriptionPattern("governance.*"),
	}
}

func (p *ConsentPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "membership.") || strings.HasPrefix(t, "authority.") || strings.HasPrefix(t, "governance.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group B: Social Order ---

// NormPrimitive manages explicit shared expectations with enforcement.
type NormPrimitive struct{}

func NewNormPrimitive() *NormPrimitive { return &NormPrimitive{} }

func (p *NormPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Norm") }
func (p *NormPrimitive) Layer() types.Layer              { return layer3 }
func (p *NormPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *NormPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *NormPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("consensus.*"),
		types.MustSubscriptionPattern("governance.*"),
	}
}

func (p *NormPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "consensus.") || strings.HasPrefix(t, "governance.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ReputationPrimitive tracks how the group perceives each member.
type ReputationPrimitive struct{}

func NewReputationPrimitive() *ReputationPrimitive { return &ReputationPrimitive{} }

func (p *ReputationPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Reputation") }
func (p *ReputationPrimitive) Layer() types.Layer              { return layer3 }
func (p *ReputationPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *ReputationPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *ReputationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("trust.*"),
		types.MustSubscriptionPattern("fulfillment.*"),
		types.MustSubscriptionPattern("breach.*"),
		types.MustSubscriptionPattern("sanction.*"),
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

// SanctionPrimitive applies consequences for norm violations.
type SanctionPrimitive struct{}

func NewSanctionPrimitive() *SanctionPrimitive { return &SanctionPrimitive{} }

func (p *SanctionPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Sanction") }
func (p *SanctionPrimitive) Layer() types.Layer              { return layer3 }
func (p *SanctionPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *SanctionPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *SanctionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("norm.*"),
		types.MustSubscriptionPattern("violation.*"),
	}
}

func (p *SanctionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	violations := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "norm.") || strings.HasPrefix(t, "violation.") {
			violations++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "violationsProcessed", Value: violations},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AuthorityPrimitive manages the legitimate power to make decisions for the group.
type AuthorityPrimitive struct{}

func NewAuthorityPrimitive() *AuthorityPrimitive { return &AuthorityPrimitive{} }

func (p *AuthorityPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Authority") }
func (p *AuthorityPrimitive) Layer() types.Layer              { return layer3 }
func (p *AuthorityPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *AuthorityPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *AuthorityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("role.*"),
		types.MustSubscriptionPattern("consent.*"),
		types.MustSubscriptionPattern("governance.*"),
	}
}

func (p *AuthorityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "role.") || strings.HasPrefix(t, "consent.") || strings.HasPrefix(t, "governance.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group C: Collective Agency ---

// PropertyPrimitive tracks ownership and access rights over resources.
type PropertyPrimitive struct{}

func NewPropertyPrimitive() *PropertyPrimitive { return &PropertyPrimitive{} }

func (p *PropertyPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Property") }
func (p *PropertyPrimitive) Layer() types.Layer              { return layer3 }
func (p *PropertyPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *PropertyPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *PropertyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("resource.*"),
		types.MustSubscriptionPattern("exchange.*"),
		types.MustSubscriptionPattern("authority.*"),
	}
}

func (p *PropertyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "resource.") || strings.HasPrefix(t, "exchange.") || strings.HasPrefix(t, "authority.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// CommonsPrimitive manages shared resources that belong to the group collectively.
type CommonsPrimitive struct{}

func NewCommonsPrimitive() *CommonsPrimitive { return &CommonsPrimitive{} }

func (p *CommonsPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Commons") }
func (p *CommonsPrimitive) Layer() types.Layer              { return layer3 }
func (p *CommonsPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *CommonsPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *CommonsPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("property.*"),
		types.MustSubscriptionPattern("group.*"),
		types.MustSubscriptionPattern("resource.*"),
	}
}

func (p *CommonsPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "property.") || strings.HasPrefix(t, "group.") || strings.HasPrefix(t, "resource.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// GovernancePrimitive manages collective decision-making processes.
type GovernancePrimitive struct{}

func NewGovernancePrimitive() *GovernancePrimitive { return &GovernancePrimitive{} }

func (p *GovernancePrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Governance") }
func (p *GovernancePrimitive) Layer() types.Layer              { return layer3 }
func (p *GovernancePrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *GovernancePrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *GovernancePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("authority.*"),
		types.MustSubscriptionPattern("norm.*"),
		types.MustSubscriptionPattern("consent.*"),
	}
}

func (p *GovernancePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "authority.") || strings.HasPrefix(t, "norm.") || strings.HasPrefix(t, "consent.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// CollectiveActPrimitive represents coordinated action taken by a group as a whole.
type CollectiveActPrimitive struct{}

func NewCollectiveActPrimitive() *CollectiveActPrimitive { return &CollectiveActPrimitive{} }

func (p *CollectiveActPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("CollectiveAct") }
func (p *CollectiveActPrimitive) Layer() types.Layer              { return layer3 }
func (p *CollectiveActPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *CollectiveActPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *CollectiveActPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("governance.*"),
		types.MustSubscriptionPattern("consent.*"),
		types.MustSubscriptionPattern("act.*"),
	}
}

func (p *CollectiveActPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "governance.") || strings.HasPrefix(t, "consent.") || strings.HasPrefix(t, "act.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
