// Package layer9 implements the Layer 9 Relationship primitives.
// Groups: Connection (Bond, Attachment, Recognition, Intimacy),
// RelationalDynamics (Attunement, Rupture, Repair, Loyalty),
// RelationalIdentity (MutualConstitution, RelationalObligation, Grief, Forgiveness).
package layer9

import (
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer9 = types.MustLayer(9)
var cadence1 = types.MustCadence(1)

// --- Group A: Connection ---

// BondPrimitive tracks the formation and strength of connections between actors.
type BondPrimitive struct{}

func NewBondPrimitive() *BondPrimitive { return &BondPrimitive{} }

func (p *BondPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Bond") }
func (p *BondPrimitive) Layer() types.Layer               { return layer9 }
func (p *BondPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *BondPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *BondPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("trust.*"),
		types.MustSubscriptionPattern("attachment.*"),
		types.MustSubscriptionPattern("edge.created"),
	}
}

func (p *BondPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "trust.") || strings.HasPrefix(t, "attachment.") || strings.HasPrefix(t, "edge.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AttachmentPrimitive assesses the quality and security of connections.
type AttachmentPrimitive struct{}

func NewAttachmentPrimitive() *AttachmentPrimitive { return &AttachmentPrimitive{} }

func (p *AttachmentPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Attachment") }
func (p *AttachmentPrimitive) Layer() types.Layer               { return layer9 }
func (p *AttachmentPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AttachmentPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AttachmentPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("bond.*"),
		types.MustSubscriptionPattern("presence.*"),
		types.MustSubscriptionPattern("care.*"),
	}
}

func (p *AttachmentPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "bond.") || strings.HasPrefix(t, "presence.") || strings.HasPrefix(t, "care.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RecognitionPrimitive tracks being seen and acknowledged by another.
type RecognitionPrimitive struct{}

func NewRecognitionPrimitive() *RecognitionPrimitive { return &RecognitionPrimitive{} }

func (p *RecognitionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Recognition") }
func (p *RecognitionPrimitive) Layer() types.Layer               { return layer9 }
func (p *RecognitionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RecognitionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RecognitionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("acknowledgment.*"),
		types.MustSubscriptionPattern("gratitude.*"),
		types.MustSubscriptionPattern("message.*"),
	}
}

func (p *RecognitionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "acknowledgment.") || strings.HasPrefix(t, "gratitude.") || strings.HasPrefix(t, "message.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// IntimacyPrimitive measures depth of mutual knowledge and vulnerability between actors.
type IntimacyPrimitive struct{}

func NewIntimacyPrimitive() *IntimacyPrimitive { return &IntimacyPrimitive{} }

func (p *IntimacyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Intimacy") }
func (p *IntimacyPrimitive) Layer() types.Layer               { return layer9 }
func (p *IntimacyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *IntimacyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *IntimacyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("vulnerability.*"),
		types.MustSubscriptionPattern("trust.*"),
		types.MustSubscriptionPattern("boundary.*"),
	}
}

func (p *IntimacyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "vulnerability.") || strings.HasPrefix(t, "trust.") || strings.HasPrefix(t, "boundary.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group B: Relational Dynamics ---

// AttunementPrimitive detects responsiveness and synchronisation between actors.
type AttunementPrimitive struct{}

func NewAttunementPrimitive() *AttunementPrimitive { return &AttunementPrimitive{} }

func (p *AttunementPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Attunement") }
func (p *AttunementPrimitive) Layer() types.Layer               { return layer9 }
func (p *AttunementPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AttunementPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AttunementPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("signal.*"),
		types.MustSubscriptionPattern("recognition.*"),
		types.MustSubscriptionPattern("empathy.*"),
	}
}

func (p *AttunementPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "signal.") || strings.HasPrefix(t, "recognition.") || strings.HasPrefix(t, "empathy.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RupturePrimitive detects when relationships break or are damaged.
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
		types.MustSubscriptionPattern("harm.*"),
	}
}

func (p *RupturePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "contract.") || strings.HasPrefix(t, "trust.") || strings.HasPrefix(t, "dispute.") || strings.HasPrefix(t, "harm.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RepairPrimitive rebuilds relationships after rupture.
type RepairPrimitive struct{}

func NewRepairPrimitive() *RepairPrimitive { return &RepairPrimitive{} }

func (p *RepairPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Repair") }
func (p *RepairPrimitive) Layer() types.Layer               { return layer9 }
func (p *RepairPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RepairPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RepairPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("rupture.*"),
		types.MustSubscriptionPattern("apology.*"),
		types.MustSubscriptionPattern("forgiveness.*"),
	}
}

func (p *RepairPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "rupture.") || strings.HasPrefix(t, "apology.") || strings.HasPrefix(t, "forgiveness.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// LoyaltyPrimitive tracks sustained commitment to a relationship through difficulty.
type LoyaltyPrimitive struct{}

func NewLoyaltyPrimitive() *LoyaltyPrimitive { return &LoyaltyPrimitive{} }

func (p *LoyaltyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Loyalty") }
func (p *LoyaltyPrimitive) Layer() types.Layer               { return layer9 }
func (p *LoyaltyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *LoyaltyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *LoyaltyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("commitment.*"),
		types.MustSubscriptionPattern("bond.*"),
		types.MustSubscriptionPattern("repair.*"),
	}
}

func (p *LoyaltyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "commitment.") || strings.HasPrefix(t, "bond.") || strings.HasPrefix(t, "repair.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group C: Relational Identity ---

// MutualConstitutionPrimitive tracks how actors shape each other's identities through relationship.
type MutualConstitutionPrimitive struct{}

func NewMutualConstitutionPrimitive() *MutualConstitutionPrimitive {
	return &MutualConstitutionPrimitive{}
}

func (p *MutualConstitutionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("MutualConstitution") }
func (p *MutualConstitutionPrimitive) Layer() types.Layer               { return layer9 }
func (p *MutualConstitutionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *MutualConstitutionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *MutualConstitutionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("identity.*"),
		types.MustSubscriptionPattern("bond.*"),
		types.MustSubscriptionPattern("intimacy.*"),
	}
}

func (p *MutualConstitutionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "identity.") || strings.HasPrefix(t, "bond.") || strings.HasPrefix(t, "intimacy.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RelationalObligationPrimitive tracks responsibilities that arise from relationships.
type RelationalObligationPrimitive struct{}

func NewRelationalObligationPrimitive() *RelationalObligationPrimitive {
	return &RelationalObligationPrimitive{}
}

func (p *RelationalObligationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("RelationalObligation") }
func (p *RelationalObligationPrimitive) Layer() types.Layer               { return layer9 }
func (p *RelationalObligationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RelationalObligationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RelationalObligationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("obligation.*"),
		types.MustSubscriptionPattern("loyalty.*"),
		types.MustSubscriptionPattern("commitment.*"),
	}
}

func (p *RelationalObligationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "obligation.") || strings.HasPrefix(t, "loyalty.") || strings.HasPrefix(t, "commitment.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// GriefPrimitive processes the loss of a relationship or a related actor.
type GriefPrimitive struct{}

func NewGriefPrimitive() *GriefPrimitive { return &GriefPrimitive{} }

func (p *GriefPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Grief") }
func (p *GriefPrimitive) Layer() types.Layer               { return layer9 }
func (p *GriefPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *GriefPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *GriefPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("actor.memorial"),
		types.MustSubscriptionPattern("loss.*"),
		types.MustSubscriptionPattern("rupture.*"),
	}
}

func (p *GriefPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "actor.") || strings.HasPrefix(t, "loss.") || strings.HasPrefix(t, "rupture.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ForgivenessPrimitive enables releasing resentment and restoring relational possibility.
type ForgivenessPrimitive struct{}

func NewForgivenessPrimitive() *ForgivenessPrimitive { return &ForgivenessPrimitive{} }

func (p *ForgivenessPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Forgiveness") }
func (p *ForgivenessPrimitive) Layer() types.Layer               { return layer9 }
func (p *ForgivenessPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ForgivenessPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ForgivenessPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("repair.*"),
		types.MustSubscriptionPattern("grief.*"),
		types.MustSubscriptionPattern("harm.*"),
	}
}

func (p *ForgivenessPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "repair.") || strings.HasPrefix(t, "grief.") || strings.HasPrefix(t, "harm.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
