package compositions_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/compositions"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

func TestMeaningGrammar(t *testing.T) {
	t.Run("ExamineAndReframe", func(t *testing.T) {
		env := newTestEnv(t)
		meaning := compositions.NewMeaningGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeHuman)

		exam, _ := meaning.Examine(env.ctx, analyst.ID(),
			"assumption: more features = more value — blind spot: user cognitive load",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		reframe, _ := meaning.Reframe(env.ctx, analyst.ID(),
			"less is more: removing friction > adding features",
			exam.ID(), env.convID, signer)

		ancestors := env.ancestors(reframe.ID(), 5)
		if !containsEvent(ancestors, exam.ID()) {
			t.Error("reframe should trace to examination")
		}
		env.verifyChain()
	})

	t.Run("QuestionAndDistill", func(t *testing.T) {
		env := newTestEnv(t)
		meaning := compositions.NewMeaningGrammar(env.grammar)
		critic := env.actor("Critic", 1, event.ActorTypeHuman)

		status_quo, _ := env.grammar.Emit(env.ctx, critic.ID(),
			"practice: deploy only on Tuesdays",
			env.convID, []types.EventID{env.boot.ID()}, signer)

		question, _ := meaning.Question(env.ctx, critic.ID(),
			"why Tuesdays? is this tradition or evidence-based?",
			status_quo.ID(), env.convID, signer)
		wisdom, _ := meaning.Distill(env.ctx, critic.ID(),
			"the day matters less than having a consistent cadence with buffer time",
			question.ID(), env.convID, signer)

		ancestors := env.ancestors(wisdom.ID(), 10)
		if !containsEvent(ancestors, status_quo.ID()) {
			t.Error("wisdom should trace to status quo")
		}
		env.verifyChain()
	})

	t.Run("BeautifyAndLiken", func(t *testing.T) {
		env := newTestEnv(t)
		meaning := compositions.NewMeaningGrammar(env.grammar)
		creator := env.actor("Creator", 1, event.ActorTypeHuman)

		beauty, _ := meaning.Beautify(env.ctx, creator.ID(),
			"the event graph as a living river — each event a drop, the chain a flowing current",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		metaphor, _ := meaning.Liken(env.ctx, creator.ID(),
			"distributed consensus is like starlings murmurating — local rules, global order",
			beauty.ID(), env.convID, signer)

		ancestors := env.ancestors(metaphor.ID(), 5)
		if !containsEvent(ancestors, beauty.ID()) {
			t.Error("metaphor should trace to beauty")
		}
		env.verifyChain()
	})

	t.Run("TeachAndTranslate", func(t *testing.T) {
		env := newTestEnv(t)
		meaning := compositions.NewMeaningGrammar(env.grammar)
		mentor := env.actor("Mentor", 1, event.ActorTypeHuman)
		student := env.actor("Student", 2, event.ActorTypeHuman)

		channel, _ := meaning.Teach(env.ctx, mentor.ID(), student.ID(),
			types.Some(types.MustDomainScope("architecture")),
			env.boot.ID(), env.convID, signer)
		lesson, _ := env.grammar.Emit(env.ctx, mentor.ID(),
			"CQRS separates reads from writes",
			env.convID, []types.EventID{channel.ID()}, signer)
		translation, _ := meaning.Translate(env.ctx, student.ID(),
			"so it's like having a library (read) and a journal (write) instead of one notebook",
			lesson.ID(), env.convID, signer)

		ancestors := env.ancestors(translation.ID(), 10)
		if !containsEvent(ancestors, channel.ID()) {
			t.Error("translation should trace to teaching channel")
		}
		env.verifyChain()
	})

	t.Run("LightenAndProphesy", func(t *testing.T) {
		env := newTestEnv(t)
		meaning := compositions.NewMeaningGrammar(env.grammar)
		observer := env.actor("Observer", 1, event.ActorTypeHuman)

		humour, _ := meaning.Lighten(env.ctx, observer.ID(),
			"a distributed system is just a system where you can't get work done because some computer you've never heard of has crashed",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		prophecy, _ := meaning.Prophesy(env.ctx, observer.ID(),
			"if current adoption trends continue, event sourcing will be table stakes in 3 years",
			[]types.EventID{humour.ID()}, env.convID, signer)

		_ = prophecy
		env.verifyChain()
	})

	t.Run("PostMortem", func(t *testing.T) {
		env := newTestEnv(t)
		meaning := compositions.NewMeaningGrammar(env.grammar)
		lead := env.actor("Lead", 1, event.ActorTypeHuman)

		incident, _ := env.grammar.Emit(env.ctx, lead.ID(), "production outage",
			env.convID, []types.EventID{env.boot.ID()}, signer)

		result, err := meaning.PostMortem(env.ctx, lead.ID(),
			"we optimised for speed over resilience",
			"why didn't our monitoring catch this earlier?",
			"resilience is not the opposite of speed — it enables sustained speed",
			incident.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("PostMortem: %v", err)
		}

		ancestors := env.ancestors(result.Wisdom.ID(), 10)
		if !containsEvent(ancestors, incident.ID()) {
			t.Error("wisdom should trace to incident")
		}
		env.verifyChain()
	})

	t.Run("Mentorship", func(t *testing.T) {
		env := newTestEnv(t)
		meaning := compositions.NewMeaningGrammar(env.grammar)
		mentor := env.actor("Mentor", 1, event.ActorTypeHuman)
		student := env.actor("Student", 2, event.ActorTypeHuman)

		result, err := meaning.Mentorship(env.ctx, mentor.ID(), student.ID(),
			"don't think of tests as overhead — they're your safety net for bold changes",
			"the best code is the code you don't write",
			"in my context: prefer composition over inheritance, keep interfaces small",
			types.Some(types.MustDomainScope("software_design")),
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Mentorship: %v", err)
		}

		ancestors := env.ancestors(result.Translation.ID(), 10)
		if !containsEvent(ancestors, result.Channel.ID()) {
			t.Error("translation should trace to teaching channel")
		}
		env.verifyChain()
	})
}
