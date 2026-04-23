package compositions_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/compositions"
	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
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

	t.Run("DeclineAndInvoice", func(t *testing.T) {
		env := newTestEnv(t)
		market := compositions.NewMarketGrammar(env.grammar)
		seller := env.actor("Seller", 1, event.ActorTypeHuman)
		buyer := env.actor("Buyer", 2, event.ActorTypeHuman)

		listing, _ := market.List(env.ctx, seller.ID(), "consulting, $200/hr",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		decline, _ := market.Decline(env.ctx, buyer.ID(), "budget too high",
			[]types.EventID{listing.ID()}, env.convID, signer)
		invoice, _ := market.Invoice(env.ctx, seller.ID(), "$500 for initial assessment",
			[]types.EventID{listing.ID()}, env.convID, signer)

		ancestors := env.ancestors(decline.ID(), 5)
		if !containsEvent(ancestors, listing.ID()) {
			t.Error("decline should trace to listing")
		}
		_ = invoice
		env.verifyChain()
	})

	t.Run("EscrowAndRelease", func(t *testing.T) {
		env := newTestEnv(t)
		market := compositions.NewMarketGrammar(env.grammar)
		buyer := env.actor("Buyer", 1, event.ActorTypeHuman)
		seller := env.actor("Seller", 2, event.ActorTypeHuman)
		escrowAgent := env.actor("Escrow", 3, event.ActorTypeAI)

		acceptance, _ := market.Accept(env.ctx, buyer.ID(), seller.ID(),
			"$1000 engagement", types.MustDomainScope("consulting"),
			env.boot.ID(), env.convID, signer)
		escrow, _ := market.Escrow(env.ctx, buyer.ID(), escrowAgent.ID(),
			types.MustDomainScope("consulting"), types.MustWeight(0.5),
			acceptance.ID(), env.convID, signer)
		release, _ := market.Release(env.ctx, escrowAgent.ID(), seller.ID(),
			"work delivered and confirmed", types.MustDomainScope("consulting"),
			escrow.ID(), env.convID, signer)

		ancestors := env.ancestors(release.ID(), 10)
		if !containsEvent(ancestors, acceptance.ID()) {
			t.Error("release should trace to acceptance")
		}
		env.verifyChain()
	})

	t.Run("Barter", func(t *testing.T) {
		env := newTestEnv(t)
		market := compositions.NewMarketGrammar(env.grammar)
		partyA := env.actor("PartyA", 1, event.ActorTypeHuman)
		partyB := env.actor("PartyB", 2, event.ActorTypeHuman)

		result, err := market.Barter(env.ctx, partyA.ID(), partyB.ID(),
			"logo design", "website copywriting",
			types.MustDomainScope("trade"),
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("Barter: %v", err)
		}
		ancestors := env.ancestors(result.Acceptance.ID(), 10)
		if !containsEvent(ancestors, result.Listing.ID()) {
			t.Error("acceptance should trace to listing")
		}
		if !containsEvent(ancestors, result.CounterOffer.ID()) {
			t.Error("acceptance should trace to counter offer")
		}
		env.verifyChain()
	})

	t.Run("Subscription", func(t *testing.T) {
		env := newTestEnv(t)
		market := compositions.NewMarketGrammar(env.grammar)
		subscriber := env.actor("Subscriber", 1, event.ActorTypeHuman)
		provider := env.actor("Provider", 2, event.ActorTypeHuman)

		result, err := market.Subscription(env.ctx, subscriber.ID(), provider.ID(),
			"monthly code review", []string{"$100 Jan", "$100 Feb"},
			[]string{"Jan review", "Feb review"},
			types.MustDomainScope("review"),
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Subscription: %v", err)
		}
		if len(result.Payments) != 2 {
			t.Errorf("expected 2 payments, got %d", len(result.Payments))
		}
		if len(result.Deliveries) != 2 {
			t.Errorf("expected 2 deliveries, got %d", len(result.Deliveries))
		}
		ancestors := env.ancestors(result.Deliveries[1].ID(), 15)
		if !containsEvent(ancestors, result.Acceptance.ID()) {
			t.Error("final delivery should trace to acceptance")
		}
		env.verifyChain()
	})

	t.Run("Refund", func(t *testing.T) {
		env := newTestEnv(t)
		market := compositions.NewMarketGrammar(env.grammar)
		buyer := env.actor("Buyer", 1, event.ActorTypeHuman)
		seller := env.actor("Seller", 2, event.ActorTypeHuman)

		result, err := market.Refund(env.ctx, buyer.ID(), seller.ID(),
			"deliverable incomplete", "agreed to refund", "$500",
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Refund: %v", err)
		}
		ancestors := env.ancestors(result.Reversal.ID(), 10)
		if !containsEvent(ancestors, result.Dispute.ID()) {
			t.Error("reversal should trace to dispute")
		}
		if !containsEvent(ancestors, result.Resolution.ID()) {
			t.Error("reversal should trace to resolution")
		}
		env.verifyChain()
	})

	t.Run("ReputationTransfer", func(t *testing.T) {
		env := newTestEnv(t)
		market := compositions.NewMarketGrammar(env.grammar)
		rater1 := env.actor("Rater1", 1, event.ActorTypeHuman)
		rater2 := env.actor("Rater2", 2, event.ActorTypeHuman)
		target := env.actor("Target", 3, event.ActorTypeHuman)

		result, err := market.ReputationTransfer(env.ctx,
			[]types.ActorID{rater1.ID(), rater2.ID()},
			[]types.EventID{env.boot.ID(), env.boot.ID()},
			target.ID(),
			[]types.Weight{types.MustWeight(0.8), types.MustWeight(0.6)},
			types.Some(types.MustDomainScope("review")),
			env.convID, signer)
		if err != nil {
			t.Fatalf("ReputationTransfer: %v", err)
		}
		if len(result.Ratings) != 2 {
			t.Errorf("expected 2 ratings, got %d", len(result.Ratings))
		}
		env.verifyChain()
	})

	t.Run("Arbitration", func(t *testing.T) {
		env := newTestEnv(t)
		market := compositions.NewMarketGrammar(env.grammar)
		plaintiff := env.actor("Plaintiff", 1, event.ActorTypeHuman)
		defendant := env.actor("Defendant", 2, event.ActorTypeHuman)
		arbiter := env.actor("Arbiter", 3, event.ActorTypeAI)

		result, err := market.Arbitration(env.ctx, plaintiff.ID(), defendant.ID(),
			arbiter.ID(), "contract breach",
			types.MustDomainScope("consulting"), types.MustWeight(0.7),
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Arbitration: %v", err)
		}
		ancestors := env.ancestors(result.Release.ID(), 10)
		if !containsEvent(ancestors, result.Dispute.ID()) {
			t.Error("release should trace to dispute")
		}
		if !containsEvent(ancestors, result.Escrow.ID()) {
			t.Error("release should trace to escrow")
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
