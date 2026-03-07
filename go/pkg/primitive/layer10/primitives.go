// Package layer10 implements the Layer 10 Community primitives.
// Groups: Belonging (Home, Contribution, Inclusion, Tradition),
// Stewardship (Commons, Sustainability, Succession, Renewal),
// Celebration (Milestone, Ceremony, Story, Gift).
package layer10

import (
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer10 = types.MustLayer(10)
var cadence1 = types.MustCadence(1)

// --- Group 0: Belonging ---

// HomePrimitive identifies the sense of place in a community.
type HomePrimitive struct{}

func NewHomePrimitive() *HomePrimitive { return &HomePrimitive{} }

func (p *HomePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Home") }
func (p *HomePrimitive) Layer() types.Layer               { return layer10 }
func (p *HomePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *HomePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *HomePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("group.*"),
		types.MustSubscriptionPattern("attachment.*"),
		types.MustSubscriptionPattern("presence.*"),
	}
}

func (p *HomePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ContributionPrimitive records what each member gives.
type ContributionPrimitive struct{}

func NewContributionPrimitive() *ContributionPrimitive { return &ContributionPrimitive{} }

func (p *ContributionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Contribution") }
func (p *ContributionPrimitive) Layer() types.Layer               { return layer10 }
func (p *ContributionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ContributionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ContributionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("artefact.created"),
		types.MustSubscriptionPattern("review.*"),
		types.MustSubscriptionPattern("care.action"),
	}
}

func (p *ContributionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// InclusionPrimitive actively ensures everyone can participate.
type InclusionPrimitive struct{}

func NewInclusionPrimitive() *InclusionPrimitive { return &InclusionPrimitive{} }

func (p *InclusionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Inclusion") }
func (p *InclusionPrimitive) Layer() types.Layer               { return layer10 }
func (p *InclusionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *InclusionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *InclusionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("group.*"),
		types.MustSubscriptionPattern("exclusion.*"),
		types.MustSubscriptionPattern("fairness.*"),
	}
}

func (p *InclusionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// TraditionPrimitive identifies practices that define a community.
type TraditionPrimitive struct{}

func NewTraditionPrimitive() *TraditionPrimitive { return &TraditionPrimitive{} }

func (p *TraditionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Tradition") }
func (p *TraditionPrimitive) Layer() types.Layer               { return layer10 }
func (p *TraditionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *TraditionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *TraditionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("convention.detected"),
		types.MustSubscriptionPattern("heritage.*"),
		types.MustSubscriptionPattern("pattern.detected"),
	}
}

func (p *TraditionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 1: Stewardship ---

// CommonsPrimitive identifies shared resources belonging to the community.
type CommonsPrimitive struct{}

func NewCommonsPrimitive() *CommonsPrimitive { return &CommonsPrimitive{} }

func (p *CommonsPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Commons") }
func (p *CommonsPrimitive) Layer() types.Layer               { return layer10 }
func (p *CommonsPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *CommonsPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *CommonsPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("artefact.*"),
		types.MustSubscriptionPattern("group.*"),
	}
}

func (p *CommonsPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// SustainabilityPrimitive assesses whether the community can continue.
type SustainabilityPrimitive struct{}

func NewSustainabilityPrimitive() *SustainabilityPrimitive { return &SustainabilityPrimitive{} }

func (p *SustainabilityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Sustainability") }
func (p *SustainabilityPrimitive) Layer() types.Layer               { return layer10 }
func (p *SustainabilityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *SustainabilityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *SustainabilityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("health.*"),
		types.MustSubscriptionPattern("commons.*"),
		types.MustSubscriptionPattern("contribution.*"),
	}
}

func (p *SustainabilityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// SuccessionPrimitive manages passing stewardship to the next generation.
type SuccessionPrimitive struct{}

func NewSuccessionPrimitive() *SuccessionPrimitive { return &SuccessionPrimitive{} }

func (p *SuccessionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Succession") }
func (p *SuccessionPrimitive) Layer() types.Layer               { return layer10 }
func (p *SuccessionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *SuccessionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *SuccessionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("delegation.*"),
		types.MustSubscriptionPattern("actor.memorial"),
		types.MustSubscriptionPattern("role.*"),
	}
}

func (p *SuccessionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RenewalPrimitive handles how communities regenerate.
type RenewalPrimitive struct{}

func NewRenewalPrimitive() *RenewalPrimitive { return &RenewalPrimitive{} }

func (p *RenewalPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Renewal") }
func (p *RenewalPrimitive) Layer() types.Layer               { return layer10 }
func (p *RenewalPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RenewalPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RenewalPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("sustainability.*"),
		types.MustSubscriptionPattern("innovation.*"),
		types.MustSubscriptionPattern("tradition.evolved"),
	}
}

func (p *RenewalPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 2: Celebration ---

// MilestonePrimitive recognises significant achievements.
type MilestonePrimitive struct{}

func NewMilestonePrimitive() *MilestonePrimitive { return &MilestonePrimitive{} }

func (p *MilestonePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Milestone") }
func (p *MilestonePrimitive) Layer() types.Layer               { return layer10 }
func (p *MilestonePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *MilestonePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *MilestonePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("goal.achieved"),
		types.MustSubscriptionPattern("innovation.*"),
		types.MustSubscriptionPattern("reconciliation.completed"),
	}
}

func (p *MilestonePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// CeremonyPrimitive handles formal recognition events.
type CeremonyPrimitive struct{}

func NewCeremonyPrimitive() *CeremonyPrimitive { return &CeremonyPrimitive{} }

func (p *CeremonyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Ceremony") }
func (p *CeremonyPrimitive) Layer() types.Layer               { return layer10 }
func (p *CeremonyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *CeremonyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *CeremonyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("milestone.*"),
		types.MustSubscriptionPattern("succession.*"),
		types.MustSubscriptionPattern("actor.memorial"),
	}
}

func (p *CeremonyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// StoryPrimitive maintains the community's shared narrative.
type StoryPrimitive struct{}

func NewStoryPrimitive() *StoryPrimitive { return &StoryPrimitive{} }

func (p *StoryPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Story") }
func (p *StoryPrimitive) Layer() types.Layer               { return layer10 }
func (p *StoryPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *StoryPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *StoryPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("milestone.*"),
		types.MustSubscriptionPattern("ceremony.*"),
		types.MustSubscriptionPattern("tradition.*"),
		types.MustSubscriptionPattern("memorial.created"),
	}
}

func (p *StoryPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// GiftPrimitive handles giving without expectation of return.
type GiftPrimitive struct{}

func NewGiftPrimitive() *GiftPrimitive { return &GiftPrimitive{} }

func (p *GiftPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Gift") }
func (p *GiftPrimitive) Layer() types.Layer               { return layer10 }
func (p *GiftPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *GiftPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *GiftPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("contribution.*"),
		types.MustSubscriptionPattern("gratitude.*"),
	}
}

func (p *GiftPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
