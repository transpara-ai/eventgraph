package compositions

import (
	"context"
	"fmt"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/grammar"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// MarketGrammar provides Layer 2 (Exchange) composition operations.
// 14 operations + 7 named functions for trust-based marketplaces.
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
	return m.g.Respond(ctx, source, "dispute: "+complaint, target, convID, signer)
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
