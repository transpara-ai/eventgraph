// Package layer2 implements the Layer 2 Exchange primitives.
// Groups: CommonGround (Term, Protocol, Offer, Acceptance),
// MutualBinding (Agreement, Obligation, Fulfillment, Breach),
// ValueTransfer (Exchange, Accountability, Debt, Reciprocity).
package layer2

import (
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer2 = types.MustLayer(2)
var cadence1 = types.MustCadence(1)

// --- Group A: Common Ground ---

// TermPrimitive defines shared vocabulary and conditions for exchange.
type TermPrimitive struct{}

func NewTermPrimitive() *TermPrimitive { return &TermPrimitive{} }

func (p *TermPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Term") }
func (p *TermPrimitive) Layer() types.Layer              { return layer2 }
func (p *TermPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *TermPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *TermPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("signal.*"),
		types.MustSubscriptionPattern("commitment.*"),
	}
}

func (p *TermPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "signal.") || strings.HasPrefix(t, "commitment.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ProtocolPrimitive establishes the rules and sequence of interaction.
type ProtocolPrimitive struct{}

func NewProtocolPrimitive() *ProtocolPrimitive { return &ProtocolPrimitive{} }

func (p *ProtocolPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Protocol") }
func (p *ProtocolPrimitive) Layer() types.Layer              { return layer2 }
func (p *ProtocolPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *ProtocolPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *ProtocolPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("term.*"),
		types.MustSubscriptionPattern("signal.*"),
	}
}

func (p *ProtocolPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "term.") || strings.HasPrefix(t, "signal.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// OfferPrimitive proposes an exchange to another actor.
type OfferPrimitive struct{}

func NewOfferPrimitive() *OfferPrimitive { return &OfferPrimitive{} }

func (p *OfferPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Offer") }
func (p *OfferPrimitive) Layer() types.Layer              { return layer2 }
func (p *OfferPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *OfferPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *OfferPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("term.*"),
		types.MustSubscriptionPattern("protocol.*"),
		types.MustSubscriptionPattern("intent.*"),
	}
}

func (p *OfferPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "term.") || strings.HasPrefix(t, "protocol.") || strings.HasPrefix(t, "intent.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AcceptancePrimitive records acceptance or rejection of an offer.
type AcceptancePrimitive struct{}

func NewAcceptancePrimitive() *AcceptancePrimitive { return &AcceptancePrimitive{} }

func (p *AcceptancePrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Acceptance") }
func (p *AcceptancePrimitive) Layer() types.Layer              { return layer2 }
func (p *AcceptancePrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *AcceptancePrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *AcceptancePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("offer.*"),
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

// --- Group B: Mutual Binding ---

// AgreementPrimitive represents a binding mutual commitment between parties.
type AgreementPrimitive struct{}

func NewAgreementPrimitive() *AgreementPrimitive { return &AgreementPrimitive{} }

func (p *AgreementPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Agreement") }
func (p *AgreementPrimitive) Layer() types.Layer              { return layer2 }
func (p *AgreementPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *AgreementPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *AgreementPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("offer.*"),
		types.MustSubscriptionPattern("acceptance.*"),
	}
}

func (p *AgreementPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "offer.") || strings.HasPrefix(t, "acceptance.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ObligationPrimitive tracks what actors owe each other under agreements.
type ObligationPrimitive struct{}

func NewObligationPrimitive() *ObligationPrimitive { return &ObligationPrimitive{} }

func (p *ObligationPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Obligation") }
func (p *ObligationPrimitive) Layer() types.Layer              { return layer2 }
func (p *ObligationPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *ObligationPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *ObligationPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("agreement.*"),
		types.MustSubscriptionPattern("commitment.*"),
	}
}

func (p *ObligationPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "agreement.") || strings.HasPrefix(t, "commitment.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// FulfillmentPrimitive records when an obligation has been met.
type FulfillmentPrimitive struct{}

func NewFulfillmentPrimitive() *FulfillmentPrimitive { return &FulfillmentPrimitive{} }

func (p *FulfillmentPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Fulfillment") }
func (p *FulfillmentPrimitive) Layer() types.Layer              { return layer2 }
func (p *FulfillmentPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *FulfillmentPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *FulfillmentPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("obligation.*"),
		types.MustSubscriptionPattern("act.*"),
	}
}

func (p *FulfillmentPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "obligation.") || strings.HasPrefix(t, "act.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// BreachPrimitive detects and records when an obligation is violated.
type BreachPrimitive struct{}

func NewBreachPrimitive() *BreachPrimitive { return &BreachPrimitive{} }

func (p *BreachPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Breach") }
func (p *BreachPrimitive) Layer() types.Layer              { return layer2 }
func (p *BreachPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *BreachPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *BreachPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("obligation.*"),
		types.MustSubscriptionPattern("violation.*"),
	}
}

func (p *BreachPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	breaches := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "obligation.") || strings.HasPrefix(t, "violation.") {
			breaches++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "breachesDetected", Value: breaches},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// --- Group C: Value Transfer ---

// ExchangePrimitive manages the actual transfer of value between parties.
type ExchangePrimitive struct{}

func NewExchangePrimitive() *ExchangePrimitive { return &ExchangePrimitive{} }

func (p *ExchangePrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Exchange") }
func (p *ExchangePrimitive) Layer() types.Layer              { return layer2 }
func (p *ExchangePrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *ExchangePrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *ExchangePrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("fulfillment.*"),
		types.MustSubscriptionPattern("agreement.*"),
		types.MustSubscriptionPattern("resource.*"),
	}
}

func (p *ExchangePrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "fulfillment.") || strings.HasPrefix(t, "agreement.") || strings.HasPrefix(t, "resource.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// AccountabilityPrimitive tracks who is responsible for what in an exchange.
type AccountabilityPrimitive struct{}

func NewAccountabilityPrimitive() *AccountabilityPrimitive { return &AccountabilityPrimitive{} }

func (p *AccountabilityPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Accountability") }
func (p *AccountabilityPrimitive) Layer() types.Layer              { return layer2 }
func (p *AccountabilityPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *AccountabilityPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *AccountabilityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("breach.*"),
		types.MustSubscriptionPattern("fulfillment.*"),
		types.MustSubscriptionPattern("consequence.*"),
	}
}

func (p *AccountabilityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "breach.") || strings.HasPrefix(t, "fulfillment.") || strings.HasPrefix(t, "consequence.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// DebtPrimitive tracks outstanding value owed between actors.
type DebtPrimitive struct{}

func NewDebtPrimitive() *DebtPrimitive { return &DebtPrimitive{} }

func (p *DebtPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Debt") }
func (p *DebtPrimitive) Layer() types.Layer              { return layer2 }
func (p *DebtPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *DebtPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *DebtPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("obligation.*"),
		types.MustSubscriptionPattern("exchange.*"),
		types.MustSubscriptionPattern("breach.*"),
	}
}

func (p *DebtPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "obligation.") || strings.HasPrefix(t, "exchange.") || strings.HasPrefix(t, "breach.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}

// ReciprocityPrimitive tracks the balance of give-and-take between actors over time.
type ReciprocityPrimitive struct{}

func NewReciprocityPrimitive() *ReciprocityPrimitive { return &ReciprocityPrimitive{} }

func (p *ReciprocityPrimitive) ID() types.PrimitiveID          { return types.MustPrimitiveID("Reciprocity") }
func (p *ReciprocityPrimitive) Layer() types.Layer              { return layer2 }
func (p *ReciprocityPrimitive) Lifecycle() types.LifecycleState { return types.LifecycleActive }
func (p *ReciprocityPrimitive) Cadence() types.Cadence          { return cadence1 }
func (p *ReciprocityPrimitive) Subscriptions() []types.SubscriptionPattern {
	return []types.SubscriptionPattern{
		types.MustSubscriptionPattern("exchange.*"),
		types.MustSubscriptionPattern("debt.*"),
		types.MustSubscriptionPattern("trust.*"),
	}
}

func (p *ReciprocityPrimitive) Process(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
	relevant := 0
	for _, ev := range events {
		t := ev.Type().Value()
		if strings.HasPrefix(t, "exchange.") || strings.HasPrefix(t, "debt.") || strings.HasPrefix(t, "trust.") {
			relevant++
		}
	}
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "eventsProcessed", Value: relevant},
		primitive.UpdateState{PrimitiveID: p.ID(), Key: "lastTick", Value: tick.Value()},
	}, nil
}
