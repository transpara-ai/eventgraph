// Package layer4 implements the Layer 4 Legal primitives.
// Groups: Codification (Rule, Jurisdiction, Precedent, Interpretation),
// Process (Adjudication, Appeal, DueProcess, Rights),
// Compliance (Audit, Enforcement, Amnesty, Reform).
package layer4

import (
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer4 = types.MustLayer(4)
var cadence1 = types.MustCadence(1)

// --- Group 0: Codification ---

// RulePrimitive manages formal, codified rules with explicit conditions and consequences.
type RulePrimitive struct{}

func NewRulePrimitive() *RulePrimitive { return &RulePrimitive{} }

func (p *RulePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Rule") }
func (p *RulePrimitive) Layer() types.Layer               { return layer4 }
func (p *RulePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RulePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RulePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("norm.established"),
		types.MustSubscriptionPattern("consensus.reached"),
		types.MustSubscriptionPattern("vote.result"),
	}
}

func (p *RulePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// JurisdictionPrimitive determines which rules apply where.
type JurisdictionPrimitive struct{}

func NewJurisdictionPrimitive() *JurisdictionPrimitive { return &JurisdictionPrimitive{} }

func (p *JurisdictionPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Jurisdiction") }
func (p *JurisdictionPrimitive) Layer() types.Layer               { return layer4 }
func (p *JurisdictionPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *JurisdictionPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *JurisdictionPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("rule.*"),
		types.MustSubscriptionPattern("group.*"),
	}
}

func (p *JurisdictionPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// PrecedentPrimitive tracks past decisions that inform future ones.
type PrecedentPrimitive struct{}

func NewPrecedentPrimitive() *PrecedentPrimitive { return &PrecedentPrimitive{} }

func (p *PrecedentPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Precedent") }
func (p *PrecedentPrimitive) Layer() types.Layer               { return layer4 }
func (p *PrecedentPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *PrecedentPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *PrecedentPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("dispute.resolved"),
		types.MustSubscriptionPattern("decision.*"),
	}
}

func (p *PrecedentPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	decisions := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "decision.") || strings.HasPrefix(ev.Type().Value(), "dispute.") {
			decisions++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "decisionsTracked", Value: decisions},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// InterpretationPrimitive applies rules to specific situations.
type InterpretationPrimitive struct{}

func NewInterpretationPrimitive() *InterpretationPrimitive { return &InterpretationPrimitive{} }

func (p *InterpretationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Interpretation") }
func (p *InterpretationPrimitive) Layer() types.Layer               { return layer4 }
func (p *InterpretationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *InterpretationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *InterpretationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("rule.*"),
		types.MustSubscriptionPattern("dispute.*"),
		types.MustSubscriptionPattern("precedent.*"),
	}
}

func (p *InterpretationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 1: Process ---

// AdjudicationPrimitive handles formal dispute resolution by an authorised decision-maker.
type AdjudicationPrimitive struct{}

func NewAdjudicationPrimitive() *AdjudicationPrimitive { return &AdjudicationPrimitive{} }

func (p *AdjudicationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Adjudication") }
func (p *AdjudicationPrimitive) Layer() types.Layer               { return layer4 }
func (p *AdjudicationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AdjudicationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AdjudicationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("dispute.raised"),
		types.MustSubscriptionPattern("rule.*"),
		types.MustSubscriptionPattern("precedent.*"),
	}
}

func (p *AdjudicationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	disputes := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "dispute.") {
			disputes++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "disputesProcessed", Value: disputes},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AppealPrimitive challenges rulings.
type AppealPrimitive struct{}

func NewAppealPrimitive() *AppealPrimitive { return &AppealPrimitive{} }

func (p *AppealPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Appeal") }
func (p *AppealPrimitive) Layer() types.Layer               { return layer4 }
func (p *AppealPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AppealPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AppealPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("adjudication.ruling"),
		types.MustSubscriptionPattern("exclusion.enacted"),
		types.MustSubscriptionPattern("sanction.applied"),
	}
}

func (p *AppealPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// DueProcessPrimitive ensures procedural fairness.
type DueProcessPrimitive struct{}

func NewDueProcessPrimitive() *DueProcessPrimitive { return &DueProcessPrimitive{} }

func (p *DueProcessPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("DueProcess") }
func (p *DueProcessPrimitive) Layer() types.Layer               { return layer4 }
func (p *DueProcessPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *DueProcessPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *DueProcessPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("adjudication.*"),
		types.MustSubscriptionPattern("exclusion.*"),
		types.MustSubscriptionPattern("sanction.*"),
	}
}

func (p *DueProcessPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// RightsPrimitive manages fundamental protections that override other rules.
type RightsPrimitive struct{}

func NewRightsPrimitive() *RightsPrimitive { return &RightsPrimitive{} }

func (p *RightsPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Rights") }
func (p *RightsPrimitive) Layer() types.Layer               { return layer4 }
func (p *RightsPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *RightsPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *RightsPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("rule.*"),
		types.MustSubscriptionPattern("sanction.*"),
		types.MustSubscriptionPattern("exclusion.*"),
		types.MustSubscriptionPattern("dueprocess.*"),
	}
}

func (p *RightsPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 2: Compliance ---

// AuditPrimitive performs systematic review of actions against rules.
type AuditPrimitive struct{}

func NewAuditPrimitive() *AuditPrimitive { return &AuditPrimitive{} }

func (p *AuditPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Audit") }
func (p *AuditPrimitive) Layer() types.Layer               { return layer4 }
func (p *AuditPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AuditPrimitive) Cadence() types.Cadence           { return types.MustCadence(5) }
func (p *AuditPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("clock.tick"),
		types.MustSubscriptionPattern("rule.*"),
	}
}

func (p *AuditPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastAuditTick", Value: tick.Value()},
	}, nil
}

// EnforcementPrimitive takes action when rules are broken.
type EnforcementPrimitive struct{}

func NewEnforcementPrimitive() *EnforcementPrimitive { return &EnforcementPrimitive{} }

func (p *EnforcementPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Enforcement") }
func (p *EnforcementPrimitive) Layer() types.Layer               { return layer4 }
func (p *EnforcementPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *EnforcementPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *EnforcementPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("audit.*"),
		types.MustSubscriptionPattern("rule.*"),
		types.MustSubscriptionPattern("right.violated"),
	}
}

func (p *EnforcementPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	violations := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "right.") || strings.HasPrefix(ev.Type().Value(), "audit.") {
			violations++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "violationsEnforced", Value: violations},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AmnestyPrimitive provides formal forgiveness that supersedes enforcement.
type AmnestyPrimitive struct{}

func NewAmnestyPrimitive() *AmnestyPrimitive { return &AmnestyPrimitive{} }

func (p *AmnestyPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Amnesty") }
func (p *AmnestyPrimitive) Layer() types.Layer               { return layer4 }
func (p *AmnestyPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AmnestyPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AmnestyPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("enforcement.*"),
		types.MustSubscriptionPattern("vote.result"),
	}
}

func (p *AmnestyPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ReformPrimitive changes rules based on experience.
type ReformPrimitive struct{}

func NewReformPrimitive() *ReformPrimitive { return &ReformPrimitive{} }

func (p *ReformPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Reform") }
func (p *ReformPrimitive) Layer() types.Layer               { return layer4 }
func (p *ReformPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ReformPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ReformPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("precedent.*"),
		types.MustSubscriptionPattern("right.violated"),
		types.MustSubscriptionPattern("audit.*"),
		types.MustSubscriptionPattern("dissent.*"),
	}
}

func (p *ReformPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
