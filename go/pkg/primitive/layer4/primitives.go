// Package layer4 implements the Layer 4 Legal primitives.
// Groups: Codification (Law, Right, Contract, Liability),
// Process (DueProcess, Adjudication, Remedy, Precedent),
// SovereignStructure (Jurisdiction, Sovereignty, Legitimacy, Treaty).
package layer4

import (
	"strings"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/primitive"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

var layer4 = types.MustLayer(4)
var cadence1 = types.MustCadence(1)

// --- Group A: Codification ---

// LawPrimitive manages formal, codified rules with explicit conditions and consequences.
type LawPrimitive struct{}

func NewLawPrimitive() *LawPrimitive { return &LawPrimitive{} }

func (p *LawPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Law") }
func (p *LawPrimitive) Layer() types.Layer              { return layer4 }
func (p *LawPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *LawPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *LawPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("norm.*"),
		types.MustSubscriptionPattern("governance.*"),
	}
}

func (p *LawPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "norm.") || strings.HasPrefix(t, "governance.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RightPrimitive manages fundamental protections that override other rules.
type RightPrimitive struct{}

func NewRightPrimitive() *RightPrimitive { return &RightPrimitive{} }

func (p *RightPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Right") }
func (p *RightPrimitive) Layer() types.Layer              { return layer4 }
func (p *RightPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *RightPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *RightPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("law.*"),
		types.MustSubscriptionPattern("sanction.*"),
		types.MustSubscriptionPattern("dueprocess.*"),
	}
}

func (p *RightPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "law.") || strings.HasPrefix(t, "sanction.") || strings.HasPrefix(t, "dueprocess.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ContractPrimitive manages binding bilateral agreements with legal enforcement.
type ContractPrimitive struct{}

func NewContractPrimitive() *ContractPrimitive { return &ContractPrimitive{} }

func (p *ContractPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Contract") }
func (p *ContractPrimitive) Layer() types.Layer              { return layer4 }
func (p *ContractPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *ContractPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *ContractPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("agreement.*"),
		types.MustSubscriptionPattern("law.*"),
		types.MustSubscriptionPattern("consent.*"),
	}
}

func (p *ContractPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "agreement.") || strings.HasPrefix(t, "law.") || strings.HasPrefix(t, "consent.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// LiabilityPrimitive assigns legal responsibility for breaches and harms.
type LiabilityPrimitive struct{}

func NewLiabilityPrimitive() *LiabilityPrimitive { return &LiabilityPrimitive{} }

func (p *LiabilityPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Liability") }
func (p *LiabilityPrimitive) Layer() types.Layer              { return layer4 }
func (p *LiabilityPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *LiabilityPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *LiabilityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("breach.*"),
		types.MustSubscriptionPattern("contract.*"),
		types.MustSubscriptionPattern("accountability.*"),
	}
}

func (p *LiabilityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "breach.") || strings.HasPrefix(t, "contract.") || strings.HasPrefix(t, "accountability.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group B: Process ---

// DueProcessPrimitive ensures procedural fairness in all legal actions.
type DueProcessPrimitive struct{}

func NewDueProcessPrimitive() *DueProcessPrimitive { return &DueProcessPrimitive{} }

func (p *DueProcessPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("DueProcess") }
func (p *DueProcessPrimitive) Layer() types.Layer              { return layer4 }
func (p *DueProcessPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *DueProcessPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *DueProcessPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("adjudication.*"),
		types.MustSubscriptionPattern("sanction.*"),
		types.MustSubscriptionPattern("right.*"),
	}
}

func (p *DueProcessPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "adjudication.") || strings.HasPrefix(t, "sanction.") || strings.HasPrefix(t, "right.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AdjudicationPrimitive handles formal dispute resolution by an authorised decision-maker.
type AdjudicationPrimitive struct{}

func NewAdjudicationPrimitive() *AdjudicationPrimitive { return &AdjudicationPrimitive{} }

func (p *AdjudicationPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Adjudication") }
func (p *AdjudicationPrimitive) Layer() types.Layer              { return layer4 }
func (p *AdjudicationPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *AdjudicationPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *AdjudicationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("breach.*"),
		types.MustSubscriptionPattern("liability.*"),
		types.MustSubscriptionPattern("law.*"),
		types.MustSubscriptionPattern("precedent.*"),
	}
}

func (p *AdjudicationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	disputes := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "breach.") || strings.HasPrefix(t, "liability.") {
			disputes++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "disputesProcessed", Value: disputes},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RemedyPrimitive provides corrective action to restore justice after a wrong.
type RemedyPrimitive struct{}

func NewRemedyPrimitive() *RemedyPrimitive { return &RemedyPrimitive{} }

func (p *RemedyPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Remedy") }
func (p *RemedyPrimitive) Layer() types.Layer              { return layer4 }
func (p *RemedyPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *RemedyPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *RemedyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("adjudication.*"),
		types.MustSubscriptionPattern("liability.*"),
		types.MustSubscriptionPattern("right.*"),
	}
}

func (p *RemedyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "adjudication.") || strings.HasPrefix(t, "liability.") || strings.HasPrefix(t, "right.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// PrecedentPrimitive tracks past decisions that inform future ones.
type PrecedentPrimitive struct{}

func NewPrecedentPrimitive() *PrecedentPrimitive { return &PrecedentPrimitive{} }

func (p *PrecedentPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Precedent") }
func (p *PrecedentPrimitive) Layer() types.Layer              { return layer4 }
func (p *PrecedentPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *PrecedentPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *PrecedentPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("adjudication.*"),
		types.MustSubscriptionPattern("decision.*"),
	}
}

func (p *PrecedentPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	decisions := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "adjudication.") || strings.HasPrefix(t, "decision.") {
			decisions++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "decisionsTracked", Value: decisions},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group C: Sovereign Structure ---

// JurisdictionPrimitive determines which laws apply where and to whom.
type JurisdictionPrimitive struct{}

func NewJurisdictionPrimitive() *JurisdictionPrimitive { return &JurisdictionPrimitive{} }

func (p *JurisdictionPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Jurisdiction") }
func (p *JurisdictionPrimitive) Layer() types.Layer              { return layer4 }
func (p *JurisdictionPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *JurisdictionPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *JurisdictionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("law.*"),
		types.MustSubscriptionPattern("group.*"),
		types.MustSubscriptionPattern("sovereignty.*"),
	}
}

func (p *JurisdictionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "law.") || strings.HasPrefix(t, "group.") || strings.HasPrefix(t, "sovereignty.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// SovereigntyPrimitive manages ultimate authority within a jurisdiction.
type SovereigntyPrimitive struct{}

func NewSovereigntyPrimitive() *SovereigntyPrimitive { return &SovereigntyPrimitive{} }

func (p *SovereigntyPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Sovereignty") }
func (p *SovereigntyPrimitive) Layer() types.Layer              { return layer4 }
func (p *SovereigntyPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *SovereigntyPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *SovereigntyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("jurisdiction.*"),
		types.MustSubscriptionPattern("authority.*"),
		types.MustSubscriptionPattern("legitimacy.*"),
	}
}

func (p *SovereigntyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "jurisdiction.") || strings.HasPrefix(t, "authority.") || strings.HasPrefix(t, "legitimacy.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// LegitimacyPrimitive assesses whether authority is rightfully held and exercised.
type LegitimacyPrimitive struct{}

func NewLegitimacyPrimitive() *LegitimacyPrimitive { return &LegitimacyPrimitive{} }

func (p *LegitimacyPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Legitimacy") }
func (p *LegitimacyPrimitive) Layer() types.Layer              { return layer4 }
func (p *LegitimacyPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *LegitimacyPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *LegitimacyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("sovereignty.*"),
		types.MustSubscriptionPattern("consent.*"),
		types.MustSubscriptionPattern("governance.*"),
	}
}

func (p *LegitimacyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "sovereignty.") || strings.HasPrefix(t, "consent.") || strings.HasPrefix(t, "governance.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// TreatyPrimitive manages formal agreements between sovereign jurisdictions.
type TreatyPrimitive struct{}

func NewTreatyPrimitive() *TreatyPrimitive { return &TreatyPrimitive{} }

func (p *TreatyPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Treaty") }
func (p *TreatyPrimitive) Layer() types.Layer              { return layer4 }
func (p *TreatyPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *TreatyPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *TreatyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("sovereignty.*"),
		types.MustSubscriptionPattern("jurisdiction.*"),
		types.MustSubscriptionPattern("agreement.*"),
	}
}

func (p *TreatyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "sovereignty.") || strings.HasPrefix(t, "jurisdiction.") || strings.HasPrefix(t, "agreement.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
