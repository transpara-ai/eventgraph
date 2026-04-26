package compositions

import (
	"context"
	"fmt"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/grammar"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// MarketGrammar provides Layer 2 (Exchange) composition operations.
// 14 operations + 2 named functions for trust-based marketplaces.
type MarketGrammar struct {
	g *grammar.Grammar
}

// NewMarketGrammar creates a MarketGrammar bound to the given base grammar.
func NewMarketGrammar(g *grammar.Grammar) *MarketGrammar {
	return &MarketGrammar{g: g}
}

// --- Operations (14) ---

// List publishes an offering to the market. (Offer + Emit)
func (m *MarketGrammar) List(
	ctx context.Context, source types.ActorID, offering string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Emit(ctx, source, "list: "+offering, convID, causes, signer)
}

// Bid makes a counter-offer on a listing. (Offer + Respond)
func (m *MarketGrammar) Bid(
	ctx context.Context, source types.ActorID, offer string,
	listing types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Respond(ctx, source, "bid: "+offer, listing, convID, signer)
}

// Inquire asks for clarification about an offering. (Clarification + Respond)
func (m *MarketGrammar) Inquire(
	ctx context.Context, source types.ActorID, question string,
	listing types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Respond(ctx, source, "inquire: "+question, listing, convID, signer)
}

// Negotiate opens a channel for refining terms. (Negotiation + Channel)
func (m *MarketGrammar) Negotiate(
	ctx context.Context, source types.ActorID, counterparty types.ActorID,
	scope types.Option[types.DomainScope],
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Channel(ctx, source, counterparty, scope, cause, convID, signer)
}

// Accept accepts terms, creating mutual obligation. (Acceptance + Consent)
func (m *MarketGrammar) Accept(
	ctx context.Context, buyer types.ActorID, seller types.ActorID,
	terms string, scope types.DomainScope,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Consent(ctx, buyer, seller, "accept: "+terms, scope, cause, convID, signer)
}

// Decline rejects an offer, closing negotiation. (Acceptance rejected + Emit)
func (m *MarketGrammar) Decline(
	ctx context.Context, source types.ActorID, reason string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Emit(ctx, source, "decline: "+reason, convID, causes, signer)
}

// Invoice formalizes a payment obligation. (Obligation + Emit)
func (m *MarketGrammar) Invoice(
	ctx context.Context, source types.ActorID, description string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Emit(ctx, source, "invoice: "+description, convID, causes, signer)
}

// Pay satisfies a financial obligation. (Obligation fulfilled + Emit)
func (m *MarketGrammar) Pay(
	ctx context.Context, source types.ActorID, description string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Emit(ctx, source, "pay: "+description, convID, causes, signer)
}

// Deliver satisfies a service/goods obligation. (Obligation fulfilled + Emit)
func (m *MarketGrammar) Deliver(
	ctx context.Context, source types.ActorID, description string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Emit(ctx, source, "deliver: "+description, convID, causes, signer)
}

// Confirm acknowledges receipt and satisfaction. (Acknowledgement + Emit)
func (m *MarketGrammar) Confirm(
	ctx context.Context, source types.ActorID, confirmation string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Emit(ctx, source, "confirm: "+confirmation, convID, causes, signer)
}

// Rate provides structured feedback on an exchange. (Gratitude + Endorse)
func (m *MarketGrammar) Rate(
	ctx context.Context, source types.ActorID,
	target types.EventID, targetActor types.ActorID,
	weight types.Weight, scope types.Option[types.DomainScope],
	convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Endorse(ctx, source, target, targetActor, weight, scope, convID, signer)
}

// Dispute flags a failed obligation. (Breach + Challenge)
func (m *MarketGrammar) Dispute(
	ctx context.Context, source types.ActorID, complaint string,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	_, flag, err := m.g.Challenge(ctx, source, "dispute: "+complaint, target, convID, signer)
	if err != nil {
		return event.Event{}, err
	}
	return flag, nil
}

// Escrow holds value pending conditions. (Obligation + Delegate)
func (m *MarketGrammar) Escrow(
	ctx context.Context, source types.ActorID, escrowActor types.ActorID,
	scope types.DomainScope, weight types.Weight,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Delegate(ctx, source, escrowActor, scope, weight, cause, convID, signer)
}

// Release releases escrowed value on condition. (Resolution + Consent)
func (m *MarketGrammar) Release(
	ctx context.Context, partyA types.ActorID, partyB types.ActorID,
	terms string, scope types.DomainScope,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Consent(ctx, partyA, partyB, "release: "+terms, scope, cause, convID, signer)
}

// --- Named Functions (7) ---

// AuctionResult holds the events produced by an Auction.
type AuctionResult struct {
	Listing    event.Event
	Bids       []event.Event
	Acceptance event.Event
}

// Auction runs competitive bidding: List + Bid (multiple) + Accept (highest).
func (m *MarketGrammar) Auction(
	ctx context.Context, seller types.ActorID, offering string,
	bidders []types.ActorID, bids []string,
	winnerIdx int, scope types.DomainScope,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (AuctionResult, error) {
	if len(bidders) != len(bids) {
		return AuctionResult{}, fmt.Errorf("auction: bidders and bids must have equal length")
	}
	if winnerIdx < 0 || winnerIdx >= len(bidders) {
		return AuctionResult{}, fmt.Errorf("auction: winnerIdx out of range")
	}

	listing, err := m.List(ctx, seller, offering, causes, convID, signer)
	if err != nil {
		return AuctionResult{}, fmt.Errorf("auction/list: %w", err)
	}

	result := AuctionResult{Listing: listing}
	for i, bidder := range bidders {
		bid, err := m.Bid(ctx, bidder, bids[i], listing.ID(), convID, signer)
		if err != nil {
			return AuctionResult{}, fmt.Errorf("auction/bid[%d]: %w", i, err)
		}
		result.Bids = append(result.Bids, bid)
	}

	acceptance, err := m.Accept(ctx, bidders[winnerIdx], seller,
		"auction won: "+bids[winnerIdx], scope,
		result.Bids[winnerIdx].ID(), convID, signer)
	if err != nil {
		return AuctionResult{}, fmt.Errorf("auction/accept: %w", err)
	}
	result.Acceptance = acceptance

	return result, nil
}

// MilestoneResult holds the events produced by a Milestone delivery.
type MilestoneResult struct {
	Acceptance event.Event
	Deliveries []event.Event
	Payments   []event.Event
}

// Milestone is staged delivery and payment: Accept + Deliver (partial) + Pay (partial).
func (m *MarketGrammar) Milestone(
	ctx context.Context, buyer types.ActorID, seller types.ActorID,
	terms string, milestones []string, payments []string,
	scope types.DomainScope,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (MilestoneResult, error) {
	if len(milestones) != len(payments) {
		return MilestoneResult{}, fmt.Errorf("milestone: milestones and payments must have equal length")
	}

	acceptance, err := m.Accept(ctx, buyer, seller, terms, scope, cause, convID, signer)
	if err != nil {
		return MilestoneResult{}, fmt.Errorf("milestone/accept: %w", err)
	}

	result := MilestoneResult{Acceptance: acceptance}
	prev := acceptance.ID()
	for i := range milestones {
		delivery, err := m.Deliver(ctx, seller, milestones[i], []types.EventID{prev}, convID, signer)
		if err != nil {
			return MilestoneResult{}, fmt.Errorf("milestone/deliver[%d]: %w", i, err)
		}
		result.Deliveries = append(result.Deliveries, delivery)

		payment, err := m.Pay(ctx, buyer, payments[i], []types.EventID{delivery.ID()}, convID, signer)
		if err != nil {
			return MilestoneResult{}, fmt.Errorf("milestone/pay[%d]: %w", i, err)
		}
		result.Payments = append(result.Payments, payment)
		prev = payment.ID()
	}

	return result, nil
}

// BarterResult holds the events produced by a Barter.
type BarterResult struct {
	Listing      event.Event
	CounterOffer event.Event
	Acceptance   event.Event
}

// Barter exchanges goods for goods: List + Bid (goods) + Accept.
func (m *MarketGrammar) Barter(
	ctx context.Context, partyA types.ActorID, partyB types.ActorID,
	offerA string, offerB string, scope types.DomainScope,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (BarterResult, error) {
	listing, err := m.List(ctx, partyA, offerA, causes, convID, signer)
	if err != nil {
		return BarterResult{}, fmt.Errorf("barter/list: %w", err)
	}

	counter, err := m.Bid(ctx, partyB, offerB, listing.ID(), convID, signer)
	if err != nil {
		return BarterResult{}, fmt.Errorf("barter/bid: %w", err)
	}

	acceptance, err := m.Accept(ctx, partyA, partyB, "barter: "+offerA+" for "+offerB, scope, counter.ID(), convID, signer)
	if err != nil {
		return BarterResult{}, fmt.Errorf("barter/accept: %w", err)
	}

	return BarterResult{Listing: listing, CounterOffer: counter, Acceptance: acceptance}, nil
}

// SubscriptionResult holds the events produced by a Subscription.
type SubscriptionResult struct {
	Acceptance event.Event
	Payments   []event.Event
	Deliveries []event.Event
}

// Subscription creates recurring delivery and payment: Accept + Pay + Deliver (repeated).
func (m *MarketGrammar) Subscription(
	ctx context.Context, subscriber types.ActorID, provider types.ActorID,
	terms string, periods []string, deliveries []string,
	scope types.DomainScope,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (SubscriptionResult, error) {
	if len(periods) != len(deliveries) {
		return SubscriptionResult{}, fmt.Errorf("subscription: periods and deliveries must have equal length")
	}

	acceptance, err := m.Accept(ctx, subscriber, provider, terms, scope, cause, convID, signer)
	if err != nil {
		return SubscriptionResult{}, fmt.Errorf("subscription/accept: %w", err)
	}

	result := SubscriptionResult{Acceptance: acceptance}
	prev := acceptance.ID()
	for i := range periods {
		payment, err := m.Pay(ctx, subscriber, periods[i], []types.EventID{prev}, convID, signer)
		if err != nil {
			return SubscriptionResult{}, fmt.Errorf("subscription/pay[%d]: %w", i, err)
		}
		result.Payments = append(result.Payments, payment)

		delivery, err := m.Deliver(ctx, provider, deliveries[i], []types.EventID{payment.ID()}, convID, signer)
		if err != nil {
			return SubscriptionResult{}, fmt.Errorf("subscription/deliver[%d]: %w", i, err)
		}
		result.Deliveries = append(result.Deliveries, delivery)
		prev = delivery.ID()
	}

	return result, nil
}

// RefundResult holds the events produced by a Refund.
type RefundResult struct {
	Dispute    event.Event
	Resolution event.Event
	Reversal   event.Event
}

// Refund processes a return: Dispute + resolution + Pay (reversed).
func (m *MarketGrammar) Refund(
	ctx context.Context, buyer types.ActorID, seller types.ActorID,
	complaint string, resolution string, refundAmount string,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (RefundResult, error) {
	dispute, err := m.Dispute(ctx, buyer, complaint, target, convID, signer)
	if err != nil {
		return RefundResult{}, fmt.Errorf("refund/dispute: %w", err)
	}

	resolutionEv, err := m.g.Emit(ctx, seller, "resolution: "+resolution, convID, []types.EventID{dispute.ID()}, signer)
	if err != nil {
		return RefundResult{}, fmt.Errorf("refund/resolution: %w", err)
	}

	reversal, err := m.Pay(ctx, seller, "refund: "+refundAmount, []types.EventID{resolutionEv.ID()}, convID, signer)
	if err != nil {
		return RefundResult{}, fmt.Errorf("refund/pay: %w", err)
	}

	return RefundResult{Dispute: dispute, Resolution: resolutionEv, Reversal: reversal}, nil
}

// ReputationTransferResult holds the events produced by a ReputationTransfer.
type ReputationTransferResult struct {
	Ratings []event.Event
}

// ReputationTransfer collects ratings from multiple parties: Rate (batch).
func (m *MarketGrammar) ReputationTransfer(
	ctx context.Context,
	raters []types.ActorID, targets []types.EventID, targetActor types.ActorID,
	weights []types.Weight, scope types.Option[types.DomainScope],
	convID types.ConversationID, signer event.Signer,
) (ReputationTransferResult, error) {
	if len(raters) != len(targets) || len(raters) != len(weights) {
		return ReputationTransferResult{}, fmt.Errorf("reputation-transfer: raters, targets, and weights must have equal length")
	}

	result := ReputationTransferResult{}
	for i, rater := range raters {
		rating, err := m.Rate(ctx, rater, targets[i], targetActor, weights[i], scope, convID, signer)
		if err != nil {
			return ReputationTransferResult{}, fmt.Errorf("reputation-transfer/rate[%d]: %w", i, err)
		}
		result.Ratings = append(result.Ratings, rating)
	}

	return result, nil
}

// ArbitrationResult holds the events produced by an Arbitration.
type ArbitrationResult struct {
	Dispute event.Event
	Escrow  event.Event
	Release event.Event
}

// Arbitration resolves a dispute with escrow: Dispute + Escrow + Release.
func (m *MarketGrammar) Arbitration(
	ctx context.Context, plaintiff types.ActorID, defendant types.ActorID,
	arbiter types.ActorID, complaint string,
	scope types.DomainScope, weight types.Weight,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (ArbitrationResult, error) {
	dispute, err := m.Dispute(ctx, plaintiff, complaint, target, convID, signer)
	if err != nil {
		return ArbitrationResult{}, fmt.Errorf("arbitration/dispute: %w", err)
	}

	escrow, err := m.Escrow(ctx, defendant, arbiter, scope, weight, dispute.ID(), convID, signer)
	if err != nil {
		return ArbitrationResult{}, fmt.Errorf("arbitration/escrow: %w", err)
	}

	release, err := m.Release(ctx, arbiter, plaintiff, "arbitration resolved", scope, escrow.ID(), convID, signer)
	if err != nil {
		return ArbitrationResult{}, fmt.Errorf("arbitration/release: %w", err)
	}

	return ArbitrationResult{Dispute: dispute, Escrow: escrow, Release: release}, nil
}
