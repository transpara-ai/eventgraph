package compositions_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/compositions"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

func TestSocialGrammar(t *testing.T) {
	t.Run("Norm", func(t *testing.T) {
		env := newTestEnv(t)
		social := compositions.NewSocialGrammar(env.grammar)
		alice := env.actor("Alice", 1, event.ActorTypeHuman)
		bob := env.actor("Bob", 2, event.ActorTypeHuman)

		norm, err := social.Norm(env.ctx, alice.ID(), bob.ID(),
			"all PRs require at least one review",
			types.MustDomainScope("engineering"),
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Norm: %v", err)
		}
		_ = norm
		env.verifyChain()
	})

	t.Run("Moderate", func(t *testing.T) {
		env := newTestEnv(t)
		social := compositions.NewSocialGrammar(env.grammar)
		mod := env.actor("Moderator", 1, event.ActorTypeHuman)

		content, _ := env.grammar.Emit(env.ctx, mod.ID(), "some content",
			env.convID, []types.EventID{env.boot.ID()}, signer)

		modAction, _ := social.Moderate(env.ctx, mod.ID(), content.ID(),
			"flagged: violates community guideline 3.2", env.convID, signer)

		ancestors := env.ancestors(modAction.ID(), 5)
		if !containsEvent(ancestors, content.ID()) {
			t.Error("moderation should trace to content")
		}
		env.verifyChain()
	})

	t.Run("WelcomeAndExile", func(t *testing.T) {
		env := newTestEnv(t)
		social := compositions.NewSocialGrammar(env.grammar)
		sponsor := env.actor("Sponsor", 1, event.ActorTypeHuman)
		newcomer := env.actor("Newcomer", 2, event.ActorTypeHuman)

		_, sub, err := social.Welcome(env.ctx, sponsor.ID(), newcomer.ID(),
			types.MustWeight(0.5), types.Some(types.MustDomainScope("community")),
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Welcome: %v", err)
		}

		edgeID, _ := types.NewEdgeID(sub.ID().Value())
		violation, _ := env.grammar.Emit(env.ctx, sponsor.ID(), "violation detected",
			env.convID, []types.EventID{sub.ID()}, signer)

		exile, _ := social.Exile(env.ctx, sponsor.ID(), edgeID, "violation detected",
			violation.ID(), env.convID, signer)
		_ = exile
		env.verifyChain()
	})

	t.Run("Elect", func(t *testing.T) {
		env := newTestEnv(t)
		social := compositions.NewSocialGrammar(env.grammar)
		community := env.actor("Community", 1, event.ActorTypeCommittee)
		candidate := env.actor("Candidate", 2, event.ActorTypeHuman)

		election, err := social.Elect(env.ctx, community.ID(), candidate.ID(),
			"governance lead", types.MustDomainScope("governance"),
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Elect: %v", err)
		}
		_ = election
		env.verifyChain()
	})

	t.Run("Poll", func(t *testing.T) {
		env := newTestEnv(t)
		social := compositions.NewSocialGrammar(env.grammar)
		proposer := env.actor("Proposer", 1, event.ActorTypeHuman)
		v1 := env.actor("Voter1", 2, event.ActorTypeHuman)
		v2 := env.actor("Voter2", 3, event.ActorTypeHuman)

		result, err := social.Poll(env.ctx, proposer.ID(),
			"should we adopt weekly standups?",
			[]types.ActorID{v1.ID(), v2.ID()},
			types.MustDomainScope("process"),
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Poll: %v", err)
		}
		if len(result.Votes) != 2 {
			t.Errorf("expected 2 votes, got %d", len(result.Votes))
		}
		env.verifyChain()
	})

	t.Run("Federation", func(t *testing.T) {
		env := newTestEnv(t)
		social := compositions.NewSocialGrammar(env.grammar)
		comA := env.actor("CommunityA", 1, event.ActorTypeCommittee)
		comB := env.actor("CommunityB", 2, event.ActorTypeCommittee)

		result, err := social.Federation(env.ctx, comA.ID(), comB.ID(),
			"shared moderation standards",
			types.MustDomainScope("moderation"), types.MustWeight(0.5),
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Federation: %v", err)
		}
		ancestors := env.ancestors(result.Delegation.ID(), 5)
		if !containsEvent(ancestors, result.Agreement.ID()) {
			t.Error("delegation should trace to agreement")
		}
		env.verifyChain()
	})
}
