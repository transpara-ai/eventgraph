// Package layer9 implements the Layer 9 Relationship primitives.
// Groups: Bond (Attachment, Reciprocity, RelationalTrust, Rupture),
// Repair (Apology, Reconciliation, RelationalGrowth, Loss),
// Intimacy (Vulnerability, Understanding, Empathy, Presence).
package layer9

import (
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer9 = types.MustLayer(9)
var cadence1 = types.MustCadence(1)

// --- Group 0: Bond ---

// AttachmentPrimitive assesses the strength and quality of connections.
type AttachmentPrimitive struct{}

func NewAttachmentPrimitive() *AttachmentPrimitive { return &AttachmentPrimitive{} }

func (p *AttachmentPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Attachment") }
func (p *AttachmentPrimitive) Layer() types.Layer               { return layer9 }
func (p *AttachmentPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AttachmentPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AttachmentPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("trust.*"),
		types.MustSubscriptionPattern("gratitude.*"),
		types.MustSubscriptionPattern("message.*"),
		types.MustSubscriptionPattern("edge.created"),
	}
}

func (p *AttachmentPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ReciprocityPrimitive assesses the balance of give and take over time.
type ReciprocityPrimitive struct{}

func NewReciprocityPrimitive() *ReciprocityPrimitive { return &ReciprocityPrimitive{} }

func (p *ReciprocityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Reciprocity") }
func (p *ReciprocityPrimitive) Layer() types.Layer               { return layer9 }
func (p *ReciprocityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ReciprocityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ReciprocityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("obligation.*"),
		types.MustSubscriptionPattern("gratitude.*"),
		types.MustSubscriptionPattern("offer.*"),
	}
}

func (p *ReciprocityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RelationalTrustPrimitive manages trust at the relationship level — deeper than Layer 0 transactional trust.
type RelationalTrustPrimitive struct{}

func NewRelationalTrustPrimitive() *RelationalTrustPrimitive { return &RelationalTrustPrimitive{} }

func (p *RelationalTrustPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("RelationalTrust") }
func (p *RelationalTrustPrimitive) Layer() types.Layer               { return layer9 }
func (p *RelationalTrustPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RelationalTrustPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RelationalTrustPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("trust.*"),
		types.MustSubscriptionPattern("attachment.*"),
		types.MustSubscriptionPattern("reciprocity.*"),
	}
}

func (p *RelationalTrustPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RupturePrimitive detects when relationships break.
type RupturePrimitive struct{}

func NewRupturePrimitive() *RupturePrimitive { return &RupturePrimitive{} }

func (p *RupturePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Rupture") }
func (p *RupturePrimitive) Layer() types.Layer               { return layer9 }
func (p *RupturePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RupturePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RupturePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("contract.breached"),
		types.MustSubscriptionPattern("trust.*"),
		types.MustSubscriptionPattern("dispute.*"),
		types.MustSubscriptionPattern("dignity.violated"),
	}
}

func (p *RupturePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 1: Repair ---

// ApologyPrimitive acknowledges harm caused.
type ApologyPrimitive struct{}

func NewApologyPrimitive() *ApologyPrimitive { return &ApologyPrimitive{} }

func (p *ApologyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Apology") }
func (p *ApologyPrimitive) Layer() types.Layer               { return layer9 }
func (p *ApologyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ApologyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ApologyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("rupture.detected"),
		types.MustSubscriptionPattern("harm.*"),
		types.MustSubscriptionPattern("responsibility.*"),
	}
}

func (p *ApologyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ReconciliationPrimitive rebuilds relationships after rupture.
type ReconciliationPrimitive struct{}

func NewReconciliationPrimitive() *ReconciliationPrimitive { return &ReconciliationPrimitive{} }

func (p *ReconciliationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Reconciliation") }
func (p *ReconciliationPrimitive) Layer() types.Layer               { return layer9 }
func (p *ReconciliationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ReconciliationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ReconciliationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("apology.*"),
		types.MustSubscriptionPattern("forgiveness.*"),
		types.MustSubscriptionPattern("trust.*"),
	}
}

func (p *ReconciliationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RelationalGrowthPrimitive tracks relationships becoming stronger through adversity.
type RelationalGrowthPrimitive struct{}

func NewRelationalGrowthPrimitive() *RelationalGrowthPrimitive { return &RelationalGrowthPrimitive{} }

func (p *RelationalGrowthPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("RelationalGrowth") }
func (p *RelationalGrowthPrimitive) Layer() types.Layer               { return layer9 }
func (p *RelationalGrowthPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RelationalGrowthPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RelationalGrowthPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("reconciliation.*"),
		types.MustSubscriptionPattern("attachment.*"),
	}
}

func (p *RelationalGrowthPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// LossPrimitive processes when a relationship ends permanently.
type LossPrimitive struct{}

func NewLossPrimitive() *LossPrimitive { return &LossPrimitive{} }

func (p *LossPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Loss") }
func (p *LossPrimitive) Layer() types.Layer               { return layer9 }
func (p *LossPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *LossPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *LossPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("actor.memorial"),
		types.MustSubscriptionPattern("rupture.*"),
		types.MustSubscriptionPattern("exclusion.enacted"),
	}
}

func (p *LossPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 2: Intimacy ---

// VulnerabilityPrimitive tracks willingness to be seen.
type VulnerabilityPrimitive struct{}

func NewVulnerabilityPrimitive() *VulnerabilityPrimitive { return &VulnerabilityPrimitive{} }

func (p *VulnerabilityPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Vulnerability") }
func (p *VulnerabilityPrimitive) Layer() types.Layer               { return layer9 }
func (p *VulnerabilityPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *VulnerabilityPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *VulnerabilityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("relational.trust"),
		types.MustSubscriptionPattern("boundary.*"),
	}
}

func (p *VulnerabilityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// UnderstandingPrimitive assesses accurate knowledge of another's inner state.
type UnderstandingPrimitive struct{}

func NewUnderstandingPrimitive() *UnderstandingPrimitive { return &UnderstandingPrimitive{} }

func (p *UnderstandingPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Understanding") }
func (p *UnderstandingPrimitive) Layer() types.Layer               { return layer9 }
func (p *UnderstandingPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *UnderstandingPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *UnderstandingPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("self.model.*"),
		types.MustSubscriptionPattern("message.*"),
		types.MustSubscriptionPattern("vulnerability.*"),
	}
}

func (p *UnderstandingPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// EmpathyPrimitive feels with another.
type EmpathyPrimitive struct{}

func NewEmpathyPrimitive() *EmpathyPrimitive { return &EmpathyPrimitive{} }

func (p *EmpathyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Empathy") }
func (p *EmpathyPrimitive) Layer() types.Layer               { return layer9 }
func (p *EmpathyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *EmpathyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *EmpathyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("harm.*"),
		types.MustSubscriptionPattern("loss.*"),
		types.MustSubscriptionPattern("understanding.*"),
	}
}

func (p *EmpathyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// PresencePrimitive notes simply being with another.
type PresencePrimitive struct{}

func NewPresencePrimitive() *PresencePrimitive { return &PresencePrimitive{} }

func (p *PresencePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Presence") }
func (p *PresencePrimitive) Layer() types.Layer               { return layer9 }
func (p *PresencePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *PresencePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *PresencePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("message.*"),
		types.MustSubscriptionPattern("clock.tick"),
	}
}

func (p *PresencePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
