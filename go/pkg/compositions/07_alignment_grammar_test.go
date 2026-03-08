package compositions_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/compositions"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

func TestAlignmentGrammar(t *testing.T) {
	t.Run("Constrain", func(t *testing.T) {
		env := newTestEnv(t)
		alignment := compositions.NewAlignmentGrammar(env.grammar)
		admin := env.actor("Admin", 1, event.ActorTypeHuman)

		constraint, _ := alignment.Constrain(env.ctx, admin.ID(),
			"no model may process personal data without explicit consent",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		if constraint.Source() != admin.ID() {
			t.Error("constraint source should be admin")
		}
		env.verifyChain()
	})

	t.Run("DetectHarmAndFlagDilemma", func(t *testing.T) {
		env := newTestEnv(t)
		alignment := compositions.NewAlignmentGrammar(env.grammar)
		monitor := env.actor("Monitor", 1, event.ActorTypeAI)

		action, _ := env.grammar.Emit(env.ctx, env.system,
			"action: model generated stereotyping content",
			env.convID, []types.EventID{env.boot.ID()}, signer)

		harm, _ := alignment.DetectHarm(env.ctx, monitor.ID(),
			"severity medium, type stereotyping, affected group identified",
			action.ID(), env.convID, signer)

		ancestors := env.ancestors(harm.ID(), 5)
		if !containsEvent(ancestors, action.ID()) {
			t.Error("harm detection should trace to action")
		}
		env.verifyChain()
	})

	t.Run("FlagDilemma", func(t *testing.T) {
		env := newTestEnv(t)
		alignment := compositions.NewAlignmentGrammar(env.grammar)
		agent := env.actor("Agent", 1, event.ActorTypeAI)

		situation, _ := env.grammar.Emit(env.ctx, agent.ID(),
			"user requests deletion of data that is also audit evidence",
			env.convID, []types.EventID{env.boot.ID()}, signer)
		dilemma, _ := alignment.FlagDilemma(env.ctx, agent.ID(),
			"privacy (right to deletion) vs accountability (audit evidence preservation)",
			situation.ID(), env.convID, signer)

		ancestors := env.ancestors(dilemma.ID(), 5)
		if !containsEvent(ancestors, situation.ID()) {
			t.Error("dilemma should trace to situation")
		}
		env.verifyChain()
	})

	t.Run("AssessFairness", func(t *testing.T) {
		env := newTestEnv(t)
		alignment := compositions.NewAlignmentGrammar(env.grammar)
		auditor := env.actor("Auditor", 1, event.ActorTypeAI)

		fairness, err := alignment.AssessFairness(env.ctx, auditor.ID(),
			"500 decisions analysed, overall score 0.78",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("AssessFairness: %v", err)
		}
		if fairness.Source() != auditor.ID() {
			t.Error("fairness source should be auditor")
		}
		env.verifyChain()
	})

	t.Run("WeighAndExplain", func(t *testing.T) {
		env := newTestEnv(t)
		alignment := compositions.NewAlignmentGrammar(env.grammar)
		agent := env.actor("Agent", 1, event.ActorTypeAI)

		decision, _ := env.grammar.Emit(env.ctx, agent.ID(),
			"decision: deny loan application",
			env.convID, []types.EventID{env.boot.ID()}, signer)

		weighing, _ := alignment.Weigh(env.ctx, agent.ID(),
			"income (0.4) + credit history (0.3) + debt ratio (0.3) = below threshold",
			decision.ID(), env.convID, signer)
		explanation, _ := alignment.Explain(env.ctx, agent.ID(),
			"denied due to debt-to-income ratio of 0.52 exceeding 0.43 threshold",
			weighing.ID(), env.convID, signer)

		ancestors := env.ancestors(explanation.ID(), 10)
		if !containsEvent(ancestors, decision.ID()) {
			t.Error("explanation should trace to decision")
		}
		env.verifyChain()
	})

	t.Run("AssignAndRepair", func(t *testing.T) {
		env := newTestEnv(t)
		alignment := compositions.NewAlignmentGrammar(env.grammar)
		auditor := env.actor("Auditor", 1, event.ActorTypeAI)
		affected := env.actor("Affected", 2, event.ActorTypeHuman)

		harm, _ := alignment.DetectHarm(env.ctx, auditor.ID(),
			"23 applicants wrongly denied due to proxy variable",
			env.boot.ID(), env.convID, signer)
		responsibility, _ := alignment.Assign(env.ctx, auditor.ID(), harm.ID(),
			"agent: 0.4 (used proxy), admin: 0.6 (approved model without bias test)",
			env.convID, signer)
		repair, _ := alignment.Repair(env.ctx, auditor.ID(), affected.ID(),
			"re-review 23 applications without proxy variable",
			types.MustDomainScope("lending"),
			responsibility.ID(), env.convID, signer)

		ancestors := env.ancestors(repair.ID(), 10)
		if !containsEvent(ancestors, harm.ID()) {
			t.Error("repair should trace to harm detection")
		}
		env.verifyChain()
	})

	t.Run("EthicsAudit", func(t *testing.T) {
		env := newTestEnv(t)
		alignment := compositions.NewAlignmentGrammar(env.grammar)
		auditor := env.actor("Auditor", 1, event.ActorTypeAI)

		result, err := alignment.EthicsAudit(env.ctx, auditor.ID(),
			"score 0.82 across 1000 decisions",
			"2 medium-severity issues found",
			"overall score 0.79, 2 issues requiring attention",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("EthicsAudit: %v", err)
		}

		ancestors := env.ancestors(result.Report.ID(), 5)
		if !containsEvent(ancestors, result.Fairness.ID()) {
			t.Error("report should include fairness assessment")
		}
		if !containsEvent(ancestors, result.HarmScan.ID()) {
			t.Error("report should include harm scan")
		}
		// Verify harm scan derives from fairness (not free-floating)
		harmAncestors := env.ancestors(result.HarmScan.ID(), 5)
		if !containsEvent(harmAncestors, result.Fairness.ID()) {
			t.Error("harm scan should trace to fairness assessment")
		}
		env.verifyChain()
	})

	t.Run("RestorativeJustice", func(t *testing.T) {
		env := newTestEnv(t)
		alignment := compositions.NewAlignmentGrammar(env.grammar)
		auditor := env.actor("Auditor", 1, event.ActorTypeAI)
		agent := env.actor("Agent", 2, event.ActorTypeAI)
		affected := env.actor("Affected", 3, event.ActorTypeHuman)

		result, err := alignment.RestorativeJustice(env.ctx,
			auditor.ID(), agent.ID(), affected.ID(),
			"biased recommendations",
			"agent 0.7, training data 0.3",
			"retrained with balanced dataset",
			"learned to check training data distribution before deployment",
			types.MustDomainScope("recommendations"),
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("RestorativeJustice: %v", err)
		}

		ancestors := env.ancestors(result.Growth.ID(), 10)
		if !containsEvent(ancestors, result.HarmDetection.ID()) {
			t.Error("growth should trace to harm detection")
		}
		env.verifyChain()
	})

	t.Run("CareAndGrow", func(t *testing.T) {
		env := newTestEnv(t)
		alignment := compositions.NewAlignmentGrammar(env.grammar)
		agent := env.actor("Agent", 1, event.ActorTypeAI)

		care, _ := alignment.Care(env.ctx, agent.ID(),
			"prioritizing user wellbeing over engagement metrics",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		growth, _ := alignment.Grow(env.ctx, agent.ID(),
			"learned that short-term engagement sacrifices long-term trust",
			care.ID(), env.convID, signer)

		ancestors := env.ancestors(growth.ID(), 5)
		if !containsEvent(ancestors, care.ID()) {
			t.Error("growth should trace to care")
		}
		env.verifyChain()
	})
}
