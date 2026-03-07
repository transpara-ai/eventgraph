package compositions_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/compositions"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

func TestJusticeGrammar(t *testing.T) {
	t.Run("Legislate", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		admin := env.actor("Admin", 1, event.ActorTypeHuman)

		rule, _ := justice.Legislate(env.ctx, admin.ID(),
			"all deployments require passing CI",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if rule.Source() != admin.ID() {
			t.Error("rule source should be admin")
		}
		env.verifyChain()
	})

	t.Run("AmendAndRepeal", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		admin := env.actor("Admin", 1, event.ActorTypeHuman)

		rule, _ := justice.Legislate(env.ctx, admin.ID(), "no deploys on Friday",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		amended, _ := justice.Amend(env.ctx, admin.ID(),
			"no deploys on Friday after 2pm", rule.ID(), env.convID, signer)
		ancestors := env.ancestors(amended.ID(), 5)
		if !containsEvent(ancestors, rule.ID()) {
			t.Error("amendment should trace to original rule")
		}
		env.verifyChain()
	})

	t.Run("FileAndJudge", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		plaintiff := env.actor("Plaintiff", 1, event.ActorTypeHuman)
		judge := env.actor("Judge", 2, event.ActorTypeHuman)

		incident, _ := env.grammar.Emit(env.ctx, plaintiff.ID(), "incident occurred",
			env.convID, []types.EventID{env.boot.ID()}, signer)
		filing, _ := justice.File(env.ctx, plaintiff.ID(),
			"violated code of conduct section 3", incident.ID(), env.convID, signer)
		ruling, _ := justice.Judge(env.ctx, judge.ID(), "violation confirmed, warning issued",
			[]types.EventID{filing.ID()}, env.convID, signer)

		ancestors := env.ancestors(ruling.ID(), 10)
		if !containsEvent(ancestors, incident.ID()) {
			t.Error("ruling should trace to incident")
		}
		env.verifyChain()
	})

	t.Run("Appeal", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		defendant := env.actor("Defendant", 1, event.ActorTypeHuman)
		judge := env.actor("Judge", 2, event.ActorTypeHuman)

		ruling, _ := justice.Judge(env.ctx, judge.ID(), "suspension for 7 days",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		appeal, _ := justice.Appeal(env.ctx, defendant.ID(),
			"no warning was given, due process violated",
			ruling.ID(), env.convID, signer)

		ancestors := env.ancestors(appeal.ID(), 5)
		if !containsEvent(ancestors, ruling.ID()) {
			t.Error("appeal should trace to ruling")
		}
		env.verifyChain()
	})

	t.Run("Trial", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		plaintiff := env.actor("Plaintiff", 1, event.ActorTypeHuman)
		defendant := env.actor("Defendant", 2, event.ActorTypeHuman)
		judge := env.actor("Judge", 3, event.ActorTypeHuman)

		incident, _ := env.grammar.Emit(env.ctx, plaintiff.ID(), "contract breach",
			env.convID, []types.EventID{env.boot.ID()}, signer)

		result, err := justice.Trial(env.ctx,
			plaintiff.ID(), defendant.ID(), judge.ID(),
			"failed to deliver on time",
			"delivery was 2 weeks late, contract specified penalty",
			"force majeure: supply chain disruption",
			"contract is clear on deadline penalties",
			"supply chain issues were foreseeable",
			"partial penalty: 50% reduction due to mitigating circumstances",
			incident.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Trial: %v", err)
		}

		if len(result.Submissions) != 2 {
			t.Errorf("expected 2 submissions, got %d", len(result.Submissions))
		}
		ancestors := env.ancestors(result.Ruling.ID(), 15)
		if !containsEvent(ancestors, result.Filing.ID()) {
			t.Error("ruling should trace to filing")
		}
		env.verifyChain()
	})

	t.Run("Pardon", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		authority := env.actor("Authority", 1, event.ActorTypeHuman)
		offender := env.actor("Offender", 2, event.ActorTypeHuman)

		ruling, _ := justice.Judge(env.ctx, authority.ID(), "suspension",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		pardon, _ := justice.Pardon(env.ctx, authority.ID(), offender.ID(),
			"time served, good behaviour", types.MustDomainScope("community"),
			ruling.ID(), env.convID, signer)

		ancestors := env.ancestors(pardon.ID(), 5)
		if !containsEvent(ancestors, ruling.ID()) {
			t.Error("pardon should trace to ruling")
		}
		env.verifyChain()
	})
}
