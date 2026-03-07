// Package layer2 implements the Layer 2 Exchange primitives.
// Groups: Communication (Message, Acknowledgement, Clarification, Context),
// Reciprocity (Offer, Acceptance, Obligation, Gratitude),
// Agreement (Negotiation, Consent, Contract, Dispute).
package layer2

import (
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer2 = types.MustLayer(2)
var cadence1 = types.MustCadence(1)

// --- Group 0: Communication ---

// MessagePrimitive handles sending and receiving structured messages between actors.
type MessagePrimitive struct{}

func NewMessagePrimitive() *MessagePrimitive { return &MessagePrimitive{} }

func (p *MessagePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Message") }
func (p *MessagePrimitive) Layer() types.Layer               { return layer2 }
func (p *MessagePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *MessagePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *MessagePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{types.MustSubscriptionPattern("protocol.message.*")}
}

func (p *MessagePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "messagesProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AcknowledgementPrimitive confirms receipt and understanding of messages.
type AcknowledgementPrimitive struct{}

func NewAcknowledgementPrimitive() *AcknowledgementPrimitive { return &AcknowledgementPrimitive{} }

func (p *AcknowledgementPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Acknowledgement") }
func (p *AcknowledgementPrimitive) Layer() types.Layer               { return layer2 }
func (p *AcknowledgementPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AcknowledgementPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AcknowledgementPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("message.sent"),
		types.MustSubscriptionPattern("message.received"),
	}
}

func (p *AcknowledgementPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	pending := 0
	for _, ev := range events {
		if ev.Type().Value() == "message.sent" {
			pending++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "pendingAcks", Value: pending},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ClarificationPrimitive resolves ambiguity in communication.
type ClarificationPrimitive struct{}

func NewClarificationPrimitive() *ClarificationPrimitive { return &ClarificationPrimitive{} }

func (p *ClarificationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Clarification") }
func (p *ClarificationPrimitive) Layer() types.Layer               { return layer2 }
func (p *ClarificationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ClarificationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ClarificationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("message.*"),
		types.MustSubscriptionPattern("ack.*"),
	}
}

func (p *ClarificationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ContextPrimitive maintains shared conversational context between actors.
type ContextPrimitive struct{}

func NewContextPrimitive() *ContextPrimitive { return &ContextPrimitive{} }

func (p *ContextPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Context") }
func (p *ContextPrimitive) Layer() types.Layer               { return layer2 }
func (p *ContextPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ContextPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ContextPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("message.*"),
		types.MustSubscriptionPattern("clarification.*"),
	}
}

func (p *ContextPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "contextUpdates", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 1: Reciprocity ---

// OfferPrimitive proposes something to another actor.
type OfferPrimitive struct{}

func NewOfferPrimitive() *OfferPrimitive { return &OfferPrimitive{} }

func (p *OfferPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Offer") }
func (p *OfferPrimitive) Layer() types.Layer               { return layer2 }
func (p *OfferPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *OfferPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *OfferPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("message.*"),
		types.MustSubscriptionPattern("exchange.*"),
	}
}

func (p *OfferPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AcceptancePrimitive accepts or rejects offers.
type AcceptancePrimitive struct{}

func NewAcceptancePrimitive() *AcceptancePrimitive { return &AcceptancePrimitive{} }

func (p *AcceptancePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Acceptance") }
func (p *AcceptancePrimitive) Layer() types.Layer               { return layer2 }
func (p *AcceptancePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *AcceptancePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *AcceptancePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("offer.made"),
		types.MustSubscriptionPattern("offer.withdrawn"),
	}
}

func (p *AcceptancePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	offers := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "offer.") {
			offers++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "offersProcessed", Value: offers},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ObligationPrimitive tracks what actors owe each other.
type ObligationPrimitive struct{}

func NewObligationPrimitive() *ObligationPrimitive { return &ObligationPrimitive{} }

func (p *ObligationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Obligation") }
func (p *ObligationPrimitive) Layer() types.Layer               { return layer2 }
func (p *ObligationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ObligationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ObligationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("offer.accepted"),
		types.MustSubscriptionPattern("delegation.*"),
	}
}

func (p *ObligationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// GratitudePrimitive recognises when obligations are fulfilled beyond expectation.
type GratitudePrimitive struct{}

func NewGratitudePrimitive() *GratitudePrimitive { return &GratitudePrimitive{} }

func (p *GratitudePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Gratitude") }
func (p *GratitudePrimitive) Layer() types.Layer               { return layer2 }
func (p *GratitudePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *GratitudePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *GratitudePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("obligation.fulfilled"),
		types.MustSubscriptionPattern("trust.*"),
	}
}

func (p *GratitudePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	fulfilled := 0
	for _, ev := range events {
		if ev.Type().Value() == "obligation.fulfilled" {
			fulfilled++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "fulfilmentsObserved", Value: fulfilled},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group 2: Agreement ---

// NegotiationPrimitive works toward agreement through iterative proposals.
type NegotiationPrimitive struct{}

func NewNegotiationPrimitive() *NegotiationPrimitive { return &NegotiationPrimitive{} }

func (p *NegotiationPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Negotiation") }
func (p *NegotiationPrimitive) Layer() types.Layer               { return layer2 }
func (p *NegotiationPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *NegotiationPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *NegotiationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("offer.*"),
		types.MustSubscriptionPattern("message.*"),
	}
}

func (p *NegotiationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ConsentPrimitive ensures both parties explicitly agree.
type ConsentPrimitive struct{}

func NewConsentPrimitive() *ConsentPrimitive { return &ConsentPrimitive{} }

func (p *ConsentPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Consent") }
func (p *ConsentPrimitive) Layer() types.Layer               { return layer2 }
func (p *ConsentPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ConsentPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ConsentPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("negotiation.concluded"),
		types.MustSubscriptionPattern("offer.accepted"),
		types.MustSubscriptionPattern("authority.*"),
	}
}

func (p *ConsentPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ContractPrimitive manages binding bilateral agreements with terms and enforcement.
type ContractPrimitive struct{}

func NewContractPrimitive() *ContractPrimitive { return &ContractPrimitive{} }

func (p *ContractPrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Contract") }
func (p *ContractPrimitive) Layer() types.Layer               { return layer2 }
func (p *ContractPrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *ContractPrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *ContractPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("consent.given"),
		types.MustSubscriptionPattern("negotiation.concluded"),
	}
}

func (p *ContractPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// DisputePrimitive handles disagreements about contracts or obligations.
type DisputePrimitive struct{}

func NewDisputePrimitive() *DisputePrimitive { return &DisputePrimitive{} }

func (p *DisputePrimitive) ID() types.PrimitiveID           { return types.MustPrimitiveID("Dispute") }
func (p *DisputePrimitive) Layer() types.Layer               { return layer2 }
func (p *DisputePrimitive) Lifecycle() types.LifecycleState  { return types.LifecycleActive }
func (p *DisputePrimitive) Cadence() types.Cadence           { return cadence1 }
func (p *DisputePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("contract.breached"),
		types.MustSubscriptionPattern("obligation.defaulted"),
		types.MustSubscriptionPattern("contradiction.found"),
	}
}

func (p *DisputePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	breaches := 0
	for _, ev := range events {
		if strings.HasPrefix(ev.Type().Value(), "contract.") || strings.HasPrefix(ev.Type().Value(), "obligation.") {
			breaches++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "breachesDetected", Value: breaches},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
