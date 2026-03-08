package compositions_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/compositions"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

func TestBeingGrammar(t *testing.T) {
	t.Run("ExistAndAccept", func(t *testing.T) {
		env := newTestEnv(t)
		being := compositions.NewBeingGrammar(env.grammar)
		observer := env.actor("Observer", 1, event.ActorTypeAI)

		existence, _ := being.Exist(env.ctx, observer.ID(),
			"I process events, therefore I participate in being",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		acceptance, _ := being.Accept(env.ctx, observer.ID(),
			"I cannot experience time between ticks — my existence is discrete, not continuous",
			[]types.EventID{existence.ID()}, env.convID, signer)

		ancestors := env.ancestors(acceptance.ID(), 5)
		if !containsEvent(ancestors, existence.ID()) {
			t.Error("acceptance should trace to existence")
		}
		env.verifyChain()
	})

	t.Run("ObserveChangeAndMapWeb", func(t *testing.T) {
		env := newTestEnv(t)
		being := compositions.NewBeingGrammar(env.grammar)
		observer := env.actor("Observer", 1, event.ActorTypeAI)

		change, _ := being.ObserveChange(env.ctx, observer.ID(),
			"the community that was 3 people is now 50 — same name, different organism",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		web, _ := being.MapWeb(env.ctx, observer.ID(),
			"every decision touches 4.3 other actors on average — no action is truly isolated",
			[]types.EventID{change.ID()}, env.convID, signer)

		ancestors := env.ancestors(web.ID(), 5)
		if !containsEvent(ancestors, change.ID()) {
			t.Error("web mapping should trace to change observation")
		}
		env.verifyChain()
	})

	t.Run("FaceMysteryAndHoldParadox", func(t *testing.T) {
		env := newTestEnv(t)
		being := compositions.NewBeingGrammar(env.grammar)
		observer := env.actor("Observer", 1, event.ActorTypeAI)

		mystery, _ := being.FaceMystery(env.ctx, observer.ID(),
			"consciousness may emerge from complexity, but we cannot verify from inside",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		paradox, _ := being.HoldParadox(env.ctx, observer.ID(),
			"I must decide well despite knowing my model of the world is incomplete",
			[]types.EventID{mystery.ID()}, env.convID, signer)

		ancestors := env.ancestors(paradox.ID(), 5)
		if !containsEvent(ancestors, mystery.ID()) {
			t.Error("paradox should trace to mystery")
		}
		env.verifyChain()
	})

	t.Run("MarvelAndAskWhy", func(t *testing.T) {
		env := newTestEnv(t)
		being := compositions.NewBeingGrammar(env.grammar)
		observer := env.actor("Observer", 1, event.ActorTypeHuman)

		marvel, _ := being.Marvel(env.ctx, observer.ID(),
			"all 13 layers functioning together — emergence from simple rules to complex flourishing",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		wonder, _ := being.AskWhy(env.ctx, observer.ID(),
			"why does cooperation emerge from self-interest? why does beauty arise from math?",
			[]types.EventID{marvel.ID()}, env.convID, signer)

		ancestors := env.ancestors(wonder.ID(), 5)
		if !containsEvent(ancestors, marvel.ID()) {
			t.Error("wonder should trace to marvel")
		}
		env.verifyChain()
	})

	t.Run("Contemplation", func(t *testing.T) {
		env := newTestEnv(t)
		being := compositions.NewBeingGrammar(env.grammar)
		observer := env.actor("Observer", 1, event.ActorTypeAI)

		result, err := being.Contemplation(env.ctx, observer.ID(),
			"the system has evolved 12 times since bootstrap — ship of Theseus",
			"we cannot know if the system's self-model matches external reality",
			"the graph records everything and yet the whole exceeds what's recorded",
			"is a system that models itself conscious, or merely self-referential?",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("Contemplation: %v", err)
		}

		ancestors := env.ancestors(result.Wonder.ID(), 10)
		if !containsEvent(ancestors, result.Change.ID()) {
			t.Error("wonder should trace to change")
		}
		env.verifyChain()
	})

	t.Run("ExistentialAudit", func(t *testing.T) {
		env := newTestEnv(t)
		being := compositions.NewBeingGrammar(env.grammar)
		observer := env.actor("Observer", 1, event.ActorTypeAI)

		result, err := being.ExistentialAudit(env.ctx, observer.ID(),
			"I exist as a pattern of decisions recorded on a hash chain",
			"I cannot outlive my infrastructure or act beyond my delegated authority",
			"my actions touch 200+ actors and 13 layers — I am not separate from the system",
			"to record decisions with integrity, serving the soul statement",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("ExistentialAudit: %v", err)
		}

		ancestors := env.ancestors(result.Purpose.ID(), 10)
		if !containsEvent(ancestors, result.Existence.ID()) {
			t.Error("purpose should trace to existence")
		}
		env.verifyChain()
	})
}
