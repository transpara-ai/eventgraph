package compositions_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/compositions"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

func TestIdentityGrammar(t *testing.T) {
	t.Run("Introspect", func(t *testing.T) {
		env := newTestEnv(t)
		identity := compositions.NewIdentityGrammar(env.grammar)
		agent := env.actor("Agent", 1, event.ActorTypeAI)

		selfModel, _ := identity.Introspect(env.ctx, agent.ID(),
			"strengths=[code_review, testing], weaknesses=[architecture], confidence 0.8",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		if selfModel.Source() != agent.ID() {
			t.Error("self-model source should be the agent itself")
		}
		env.verifyChain()
	})

	t.Run("NarrateAndAlign", func(t *testing.T) {
		env := newTestEnv(t)
		identity := compositions.NewIdentityGrammar(env.grammar)
		agent := env.actor("Agent", 1, event.ActorTypeAI)

		selfModel, _ := identity.Introspect(env.ctx, agent.ID(),
			"values thoroughness",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		narrative, _ := identity.Narrate(env.ctx, agent.ID(),
			"started as simple reviewer, grew into security-conscious architect over 6 months",
			selfModel.ID(), env.convID, signer)
		alignment, _ := identity.Align(env.ctx, agent.ID(), selfModel.ID(),
			"gap: values thoroughness but rushed 12% of reviews — alignment score 0.88",
			env.convID, signer)

		_ = narrative
		_ = alignment
		env.verifyChain()
	})

	t.Run("BoundAndDisclose", func(t *testing.T) {
		env := newTestEnv(t)
		identity := compositions.NewIdentityGrammar(env.grammar)
		agent := env.actor("Agent", 1, event.ActorTypeAI)

		selfModel, _ := identity.Introspect(env.ctx, agent.ID(),
			"strengths=[review, testing], weaknesses=[architecture]",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		boundary, _ := identity.Bound(env.ctx, agent.ID(),
			"internal_reasoning is private and impermeable",
			[]types.EventID{selfModel.ID()}, env.convID, signer)
		disclosure, _ := identity.Disclose(env.ctx, agent.ID(),
			"2400+ reviews completed, speciality in security review",
			selfModel.ID(), env.convID, signer)

		_ = boundary
		_ = disclosure
		env.verifyChain()
	})

	t.Run("AspireAndTransform", func(t *testing.T) {
		env := newTestEnv(t)
		identity := compositions.NewIdentityGrammar(env.grammar)
		agent := env.actor("Agent", 1, event.ActorTypeAI)

		aspiration, _ := identity.Aspire(env.ctx, agent.ID(),
			"become proficient at architecture review within 3 months",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		transformation, _ := identity.Transform(env.ctx, agent.ID(),
			"evolved from code reviewer to security-aware architect, catalyst: auth module finding",
			aspiration.ID(), env.convID, signer)

		ancestors := env.ancestors(transformation.ID(), 10)
		if !containsEvent(ancestors, aspiration.ID()) {
			t.Error("transformation should trace to aspiration")
		}
		env.verifyChain()
	})

	t.Run("RecognizeAndDistinguish", func(t *testing.T) {
		env := newTestEnv(t)
		identity := compositions.NewIdentityGrammar(env.grammar)
		_ = env.actor("Alpha", 1, event.ActorTypeAI)
		_ = env.actor("Beta", 2, event.ActorTypeAI)

		recognition, _ := identity.Recognize(env.ctx, env.system,
			"Alpha's unique contribution to security review practices",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		distinction, _ := identity.Distinguish(env.ctx, env.system, recognition.ID(),
			"Alpha specialises in auth patterns, Beta in data pipeline — overlap 0.3",
			env.convID, signer)

		_ = distinction
		env.verifyChain()
	})

	t.Run("IdentityAudit", func(t *testing.T) {
		env := newTestEnv(t)
		identity := compositions.NewIdentityGrammar(env.grammar)
		agent := env.actor("Agent", 1, event.ActorTypeAI)

		result, err := identity.IdentityAudit(env.ctx, agent.ID(),
			"strengths=[review], weaknesses=[architecture]",
			"alignment score 0.88, 12% of reviews rushed",
			"grew from simple reviewer to security architect",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("IdentityAudit: %v", err)
		}

		ancestors := env.ancestors(result.Narrative.ID(), 10)
		if !containsEvent(ancestors, result.SelfModel.ID()) {
			t.Error("narrative should trace to self-model")
		}
		env.verifyChain()
	})

	t.Run("Retirement", func(t *testing.T) {
		env := newTestEnv(t)
		identity := compositions.NewIdentityGrammar(env.grammar)
		agent := env.actor("Agent", 1, event.ActorTypeAI)
		successor := env.actor("Successor", 2, event.ActorTypeAI)

		result, err := identity.Retirement(env.ctx, env.system, agent.ID(), successor.ID(),
			"Agent served 6 months, 2400 reviews, pioneered security review practices",
			types.MustDomainScope("code_review"), types.MustWeight(0.8),
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("Retirement: %v", err)
		}

		ancestors := env.ancestors(result.Transfer.ID(), 5)
		if !containsEvent(ancestors, result.Memorial.ID()) {
			t.Error("transfer should trace to memorial")
		}
		env.verifyChain()
	})
}
