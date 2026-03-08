package compositions_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/compositions"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

func TestBondGrammar(t *testing.T) {
	t.Run("Connect", func(t *testing.T) {
		env := newTestEnv(t)
		bond := compositions.NewBondGrammar(env.grammar)
		alice := env.actor("Alice", 1, event.ActorTypeHuman)
		bob := env.actor("Bob", 2, event.ActorTypeHuman)

		sub1, sub2, err := bond.Connect(env.ctx, alice.ID(), bob.ID(),
			types.Some(types.MustDomainScope("collaboration")),
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Connect: %v", err)
		}

		_ = sub1
		_ = sub2
		env.verifyChain()
	})

	t.Run("BalanceAndAttune", func(t *testing.T) {
		env := newTestEnv(t)
		bond := compositions.NewBondGrammar(env.grammar)
		alice := env.actor("Alice", 1, event.ActorTypeHuman)

		help, _ := env.grammar.Emit(env.ctx, alice.ID(), "reviewed Bob's PR",
			env.convID, []types.EventID{env.boot.ID()}, signer)

		balance, _ := bond.Balance(env.ctx, env.system, help.ID(),
			"give/take ratio: 0.0 (balanced), Alice gave 1, Bob gave 1",
			env.convID, signer)
		attunement, _ := bond.Attune(env.ctx, alice.ID(),
			"Bob prefers async communication, values autonomy",
			[]types.EventID{balance.ID()}, env.convID, signer)

		_ = attunement
		env.verifyChain()
	})

	t.Run("OpenAndFeelWith", func(t *testing.T) {
		env := newTestEnv(t)
		bond := compositions.NewBondGrammar(env.grammar)
		alice := env.actor("Alice", 1, event.ActorTypeHuman)
		bob := env.actor("Bob", 2, event.ActorTypeHuman)

		channel, _ := bond.Open(env.ctx, alice.ID(), bob.ID(),
			types.Some(types.MustDomainScope("personal")),
			env.boot.ID(), env.convID, signer)

		share, _ := env.grammar.Emit(env.ctx, alice.ID(),
			"struggling with imposter syndrome about the architecture role",
			env.convID, []types.EventID{channel.ID()}, signer)

		empathy, _ := bond.FeelWith(env.ctx, bob.ID(),
			"I felt the same when I first led a project — it gets easier",
			share.ID(), env.convID, signer)

		ancestors := env.ancestors(empathy.ID(), 10)
		if !containsEvent(ancestors, channel.ID()) {
			t.Error("empathy should trace to channel opening")
		}
		env.verifyChain()
	})

	t.Run("BreakAndApologize", func(t *testing.T) {
		env := newTestEnv(t)
		bond := compositions.NewBondGrammar(env.grammar)
		alice := env.actor("Alice", 1, event.ActorTypeHuman)
		bob := env.actor("Bob", 2, event.ActorTypeHuman)

		_, _, _ = bond.Connect(env.ctx, alice.ID(), bob.ID(),
			types.None[types.DomainScope](),
			env.boot.ID(), env.convID, signer)

		rupture, _ := bond.Break(env.ctx, alice.ID(),
			"Bob took credit for shared work in the team meeting",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		apology, _ := bond.Apologize(env.ctx, bob.ID(),
			"I should have credited you — it was our joint work",
			[]types.EventID{rupture.ID()}, env.convID, signer)

		ancestors := env.ancestors(apology.ID(), 5)
		if !containsEvent(ancestors, rupture.ID()) {
			t.Error("apology should trace to rupture")
		}
		env.verifyChain()
	})

	t.Run("BetrayalRepair", func(t *testing.T) {
		env := newTestEnv(t)
		bond := compositions.NewBondGrammar(env.grammar)
		alice := env.actor("Alice", 1, event.ActorTypeHuman)
		bob := env.actor("Bob", 2, event.ActorTypeHuman)

		result, err := bond.BetrayalRepair(env.ctx, alice.ID(), bob.ID(),
			"Bob shared private conversation externally",
			"I violated your trust by sharing our private conversation",
			"progress 0.5, new boundaries established",
			"relationship stronger after repair — trust rebuilt on new basis",
			types.MustDomainScope("personal"),
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("BetrayalRepair: %v", err)
		}

		ancestors := env.ancestors(result.Deepened.ID(), 10)
		if !containsEvent(ancestors, result.Rupture.ID()) {
			t.Error("deepened relationship should trace to original rupture")
		}
		env.verifyChain()
	})

	t.Run("CheckIn", func(t *testing.T) {
		env := newTestEnv(t)
		bond := compositions.NewBondGrammar(env.grammar)
		alice := env.actor("Alice", 1, event.ActorTypeHuman)

		interaction, _ := env.grammar.Emit(env.ctx, alice.ID(), "collaboration event",
			env.convID, []types.EventID{env.boot.ID()}, signer)

		result, err := bond.CheckIn(env.ctx, alice.ID(),
			interaction.ID(),
			"give/take ratio balanced, both contributing equally",
			"partner prefers detailed written feedback over verbal",
			"feeling grateful for consistent support",
			env.convID, signer)
		if err != nil {
			t.Fatalf("CheckIn: %v", err)
		}

		ancestors := env.ancestors(result.Empathy.ID(), 10)
		if !containsEvent(ancestors, result.Balance.ID()) {
			t.Error("empathy should trace to balance")
		}
		env.verifyChain()
	})

	t.Run("Mourn", func(t *testing.T) {
		env := newTestEnv(t)
		bond := compositions.NewBondGrammar(env.grammar)
		alice := env.actor("Alice", 1, event.ActorTypeHuman)

		departure, _ := env.grammar.Emit(env.ctx, alice.ID(),
			"Bob has left the organization permanently",
			env.convID, []types.EventID{env.boot.ID()}, signer)
		mourning, _ := bond.Mourn(env.ctx, alice.ID(),
			"processing the loss of daily collaboration with Bob",
			[]types.EventID{departure.ID()}, env.convID, signer)

		ancestors := env.ancestors(mourning.ID(), 5)
		if !containsEvent(ancestors, departure.ID()) {
			t.Error("mourning should trace to departure")
		}
		env.verifyChain()
	})

	t.Run("Forgive", func(t *testing.T) {
		env := newTestEnv(t)
		bond := compositions.NewBondGrammar(env.grammar)
		alice := env.actor("Alice", 1, event.ActorTypeHuman)
		bob := env.actor("Bob", 2, event.ActorTypeHuman)

		sub, _ := env.grammar.Subscribe(env.ctx, alice.ID(), bob.ID(),
			types.None[types.DomainScope](),
			env.boot.ID(), env.convID, signer)

		violation, _ := bond.Break(env.ctx, alice.ID(), "trust broken",
			[]types.EventID{sub.ID()}, env.convID, signer)

		edgeID, _ := types.NewEdgeID(sub.ID().Value())
		severEv, _ := env.grammar.Sever(env.ctx, alice.ID(),
			edgeID, violation.ID(), env.convID, signer)

		forgiveEv, err := bond.Forgive(env.ctx, alice.ID(),
			severEv.ID(), bob.ID(),
			types.None[types.DomainScope](),
			env.convID, signer)
		if err != nil {
			t.Fatalf("Forgive: %v", err)
		}

		ancestors := env.ancestors(forgiveEv.ID(), 10)
		if !containsEvent(ancestors, severEv.ID()) {
			t.Error("forgiveness should trace to sever")
		}
		env.verifyChain()
	})
}
