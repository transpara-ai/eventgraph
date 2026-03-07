package compositions_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/compositions"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

func TestMarketGrammar(t *testing.T) {
	t.Run("ListAndBid", func(t *testing.T) {
		env := newTestEnv(t)
		market := compositions.NewMarketGrammar(env.grammar)
		seller := env.actor("Seller", 1, event.ActorTypeHuman)
		buyer := env.actor("Buyer", 2, event.ActorTypeHuman)

		listing, _ := market.List(env.ctx, seller.ID(), "code review, $100/hr",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		bid, _ := market.Bid(env.ctx, buyer.ID(), "$90/hr, available next week",
			listing.ID(), env.convID, signer)

		ancestors := env.ancestors(bid.ID(), 5)
		if !containsEvent(ancestors, listing.ID()) {
			t.Error("bid should trace to listing")
		}
		env.verifyChain()
	})

	t.Run("Inquire", func(t *testing.T) {
		env := newTestEnv(t)
		market := compositions.NewMarketGrammar(env.grammar)
		seller := env.actor("Seller", 1, event.ActorTypeHuman)
		buyer := env.actor("Buyer", 2, event.ActorTypeHuman)

		listing, _ := market.List(env.ctx, seller.ID(), "full-stack web app",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		question, _ := market.Inquire(env.ctx, buyer.ID(), "what's the expected timeline?",
			listing.ID(), env.convID, signer)

		ancestors := env.ancestors(question.ID(), 5)
		if !containsEvent(ancestors, listing.ID()) {
			t.Error("inquiry should trace to listing")
		}
		env.verifyChain()
	})

	t.Run("NegotiateAndAccept", func(t *testing.T) {
		env := newTestEnv(t)
		market := compositions.NewMarketGrammar(env.grammar)
		seller := env.actor("Seller", 1, event.ActorTypeHuman)
		buyer := env.actor("Buyer", 2, event.ActorTypeHuman)

		listing, _ := market.List(env.ctx, seller.ID(), "consulting, $150/hr",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		channel, _ := market.Negotiate(env.ctx, buyer.ID(), seller.ID(),
			types.Some(types.MustDomainScope("consulting")),
			listing.ID(), env.convID, signer)
		acceptance, _ := market.Accept(env.ctx, buyer.ID(), seller.ID(),
			"$130/hr, 20hr engagement", types.MustDomainScope("consulting"),
			channel.ID(), env.convID, signer)

		ancestors := env.ancestors(acceptance.ID(), 10)
		if !containsEvent(ancestors, listing.ID()) {
			t.Error("acceptance should trace to listing")
		}
		env.verifyChain()
	})

	t.Run("DeliverAndConfirm", func(t *testing.T) {
		env := newTestEnv(t)
		market := compositions.NewMarketGrammar(env.grammar)
		seller := env.actor("Seller", 1, event.ActorTypeHuman)
		buyer := env.actor("Buyer", 2, event.ActorTypeHuman)

		acceptance, _ := market.Accept(env.ctx, buyer.ID(), seller.ID(),
			"code review engagement", types.MustDomainScope("review"),
			env.boot.ID(), env.convID, signer)
		delivery, _ := market.Deliver(env.ctx, seller.ID(), "review of auth module complete",
			[]types.EventID{acceptance.ID()}, env.convID, signer)
		confirm, _ := market.Confirm(env.ctx, buyer.ID(), "quality excellent",
			[]types.EventID{delivery.ID()}, env.convID, signer)

		ancestors := env.ancestors(confirm.ID(), 10)
		if !containsEvent(ancestors, acceptance.ID()) {
			t.Error("confirmation should trace to acceptance")
		}
		env.verifyChain()
	})

	t.Run("Rate", func(t *testing.T) {
		env := newTestEnv(t)
		market := compositions.NewMarketGrammar(env.grammar)
		seller := env.actor("Seller", 1, event.ActorTypeHuman)
		buyer := env.actor("Buyer", 2, event.ActorTypeHuman)

		delivery, _ := market.Deliver(env.ctx, seller.ID(), "project complete",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		rating, _ := market.Rate(env.ctx, buyer.ID(), delivery.ID(), seller.ID(),
			types.MustWeight(0.9), types.Some(types.MustDomainScope("review")),
			env.convID, signer)

		ancestors := env.ancestors(rating.ID(), 5)
		if !containsEvent(ancestors, delivery.ID()) {
			t.Error("rating should trace to delivery")
		}
		env.verifyChain()
	})

	t.Run("Dispute", func(t *testing.T) {
		env := newTestEnv(t)
		market := compositions.NewMarketGrammar(env.grammar)
		seller := env.actor("Seller", 1, event.ActorTypeHuman)
		buyer := env.actor("Buyer", 2, event.ActorTypeHuman)

		delivery, _ := market.Deliver(env.ctx, seller.ID(), "code review delivered",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		dispute, _ := market.Dispute(env.ctx, buyer.ID(),
			"review missed critical SQL injection",
			delivery.ID(), env.convID, signer)

		ancestors := env.ancestors(dispute.ID(), 5)
		if !containsEvent(ancestors, delivery.ID()) {
			t.Error("dispute should trace to delivery")
		}
		env.verifyChain()
	})

	t.Run("Auction", func(t *testing.T) {
		env := newTestEnv(t)
		market := compositions.NewMarketGrammar(env.grammar)
		seller := env.actor("Seller", 1, event.ActorTypeHuman)
		b1 := env.actor("Bidder1", 2, event.ActorTypeHuman)
		b2 := env.actor("Bidder2", 3, event.ActorTypeHuman)

		result, err := market.Auction(env.ctx, seller.ID(), "consulting engagement",
			[]types.ActorID{b1.ID(), b2.ID()}, []string{"$500", "$750"},
			1, types.MustDomainScope("consulting"),
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("Auction: %v", err)
		}
		if len(result.Bids) != 2 {
			t.Errorf("expected 2 bids, got %d", len(result.Bids))
		}
		ancestors := env.ancestors(result.Acceptance.ID(), 10)
		if !containsEvent(ancestors, result.Listing.ID()) {
			t.Error("acceptance should trace to listing")
		}
		env.verifyChain()
	})

	t.Run("Milestone", func(t *testing.T) {
		env := newTestEnv(t)
		market := compositions.NewMarketGrammar(env.grammar)
		seller := env.actor("Seller", 1, event.ActorTypeHuman)
		buyer := env.actor("Buyer", 2, event.ActorTypeHuman)

		result, err := market.Milestone(env.ctx, buyer.ID(), seller.ID(),
			"3-phase build",
			[]string{"wireframes", "backend", "frontend"},
			[]string{"$1000", "$3000", "$2000"},
			types.MustDomainScope("webapp"),
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Milestone: %v", err)
		}
		if len(result.Deliveries) != 3 {
			t.Errorf("expected 3 deliveries, got %d", len(result.Deliveries))
		}
		ancestors := env.ancestors(result.Payments[2].ID(), 15)
		if !containsEvent(ancestors, result.Acceptance.ID()) {
			t.Error("final payment should trace to acceptance")
		}
		env.verifyChain()
	})
}
