package compositions_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/compositions"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

func TestWorkGrammar(t *testing.T) {
	t.Run("Intend", func(t *testing.T) {
		env := newTestEnv(t)
		work := compositions.NewWorkGrammar(env.grammar)
		dev := env.actor("Dev", 1, event.ActorTypeHuman)

		goal, err := work.Intend(env.ctx, dev.ID(),
			"review all pending PRs before release",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("Intend: %v", err)
		}
		if goal.Source() != dev.ID() {
			t.Error("goal source should be dev")
		}
		env.verifyChain()
	})

	t.Run("Decompose", func(t *testing.T) {
		env := newTestEnv(t)
		work := compositions.NewWorkGrammar(env.grammar)
		dev := env.actor("Dev", 1, event.ActorTypeHuman)

		goal, _ := work.Intend(env.ctx, dev.ID(), "ship v2.0",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		sub1, _ := work.Decompose(env.ctx, dev.ID(), "update auth module", goal.ID(), env.convID, signer)
		sub2, _ := work.Decompose(env.ctx, dev.ID(), "write migration guide", goal.ID(), env.convID, signer)

		ancestors := env.ancestors(sub1.ID(), 5)
		if !containsEvent(ancestors, goal.ID()) {
			t.Error("subtask should trace to goal")
		}
		_ = sub2
		env.verifyChain()
	})

	t.Run("AssignAndClaim", func(t *testing.T) {
		env := newTestEnv(t)
		work := compositions.NewWorkGrammar(env.grammar)
		lead := env.actor("Lead", 1, event.ActorTypeHuman)
		dev := env.actor("Dev", 2, event.ActorTypeHuman)

		goal, _ := work.Intend(env.ctx, lead.ID(), "fix auth bug",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		assign, _ := work.Assign(env.ctx, lead.ID(), dev.ID(),
			types.MustDomainScope("auth"), types.MustWeight(0.5),
			goal.ID(), env.convID, signer)

		claim, _ := work.Claim(env.ctx, dev.ID(), "taking auth bug fix",
			[]types.EventID{assign.ID()}, env.convID, signer)

		_ = claim
		env.verifyChain()
	})

	t.Run("BlockAndUnblock", func(t *testing.T) {
		env := newTestEnv(t)
		work := compositions.NewWorkGrammar(env.grammar)
		dev := env.actor("Dev", 1, event.ActorTypeHuman)

		task, _ := work.Intend(env.ctx, dev.ID(), "deploy to staging",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		block, _ := work.Block(env.ctx, dev.ID(), task.ID(),
			"CI pipeline broken", env.convID, signer)

		unblock, _ := work.Unblock(env.ctx, dev.ID(), "CI pipeline fixed",
			[]types.EventID{block.ID()}, env.convID, signer)

		ancestors := env.ancestors(unblock.ID(), 10)
		if !containsEvent(ancestors, task.ID()) {
			t.Error("unblock should trace to original task")
		}
		env.verifyChain()
	})

	t.Run("ProgressAndComplete", func(t *testing.T) {
		env := newTestEnv(t)
		work := compositions.NewWorkGrammar(env.grammar)
		dev := env.actor("Dev", 1, event.ActorTypeHuman)

		task, _ := work.Intend(env.ctx, dev.ID(), "implement search",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		p1, _ := work.Progress(env.ctx, dev.ID(), "basic search working", task.ID(), env.convID, signer)
		p2, _ := work.Progress(env.ctx, dev.ID(), "added fuzzy matching", p1.ID(), env.convID, signer)

		complete, _ := work.Complete(env.ctx, dev.ID(), "search with fuzzy matching",
			[]types.EventID{p2.ID()}, env.convID, signer)

		ancestors := env.ancestors(complete.ID(), 10)
		if !containsEvent(ancestors, task.ID()) {
			t.Error("completion should trace to task")
		}
		env.verifyChain()
	})

	t.Run("Review", func(t *testing.T) {
		env := newTestEnv(t)
		work := compositions.NewWorkGrammar(env.grammar)
		dev := env.actor("Dev", 1, event.ActorTypeHuman)
		reviewer := env.actor("Reviewer", 2, event.ActorTypeHuman)

		complete, _ := work.Complete(env.ctx, dev.ID(), "auth module done",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		review, _ := work.Review(env.ctx, reviewer.ID(), "approved, clean implementation",
			complete.ID(), env.convID, signer)

		ancestors := env.ancestors(review.ID(), 5)
		if !containsEvent(ancestors, complete.ID()) {
			t.Error("review should trace to completion")
		}
		env.verifyChain()
	})

	t.Run("Sprint", func(t *testing.T) {
		env := newTestEnv(t)
		work := compositions.NewWorkGrammar(env.grammar)
		lead := env.actor("Lead", 1, event.ActorTypeHuman)
		dev1 := env.actor("Dev1", 2, event.ActorTypeHuman)
		dev2 := env.actor("Dev2", 3, event.ActorTypeHuman)

		result, err := work.Sprint(env.ctx, lead.ID(), "Sprint 7: auth hardening",
			[]string{"add rate limiting", "add 2FA"},
			[]types.ActorID{dev1.ID(), dev2.ID()},
			[]types.DomainScope{types.MustDomainScope("rate_limiting"), types.MustDomainScope("two_factor")},
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("Sprint: %v", err)
		}

		if len(result.Subtasks) != 2 {
			t.Errorf("expected 2 subtasks, got %d", len(result.Subtasks))
		}
		if len(result.Assignments) != 2 {
			t.Errorf("expected 2 assignments, got %d", len(result.Assignments))
		}

		ancestors := env.ancestors(result.Assignments[1].ID(), 10)
		if !containsEvent(ancestors, result.Intent.ID()) {
			t.Error("assignment should trace to sprint intent")
		}
		env.verifyChain()
	})

	t.Run("Escalate", func(t *testing.T) {
		env := newTestEnv(t)
		work := compositions.NewWorkGrammar(env.grammar)
		dev := env.actor("Dev", 1, event.ActorTypeHuman)
		lead := env.actor("Lead", 2, event.ActorTypeHuman)

		task, _ := work.Intend(env.ctx, dev.ID(), "migrate database",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		result, err := work.Escalate(env.ctx, dev.ID(),
			"need DBA approval for schema change",
			task.ID(), lead.ID(), types.MustDomainScope("database"),
			env.convID, signer)
		if err != nil {
			t.Fatalf("Escalate: %v", err)
		}

		ancestors := env.ancestors(result.HandoffEvent.ID(), 10)
		if !containsEvent(ancestors, task.ID()) {
			t.Error("escalation should trace to original task")
		}
		env.verifyChain()
	})

	t.Run("DelegateAndVerify", func(t *testing.T) {
		env := newTestEnv(t)
		work := compositions.NewWorkGrammar(env.grammar)
		lead := env.actor("Lead", 1, event.ActorTypeHuman)
		agent := env.actor("Agent", 2, event.ActorTypeAI)

		result, err := work.DelegateAndVerify(env.ctx, lead.ID(), agent.ID(),
			types.MustDomainScope("code_review"), types.MustWeight(0.7),
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("DelegateAndVerify: %v", err)
		}

		ancestors := env.ancestors(result.ScopeEvent.ID(), 5)
		if !containsEvent(ancestors, result.AssignEvent.ID()) {
			t.Error("scope should trace to assignment")
		}
		env.verifyChain()
	})

	t.Run("Standup", func(t *testing.T) {
		env := newTestEnv(t)
		work := compositions.NewWorkGrammar(env.grammar)
		dev1 := env.actor("Dev1", 1, event.ActorTypeHuman)
		dev2 := env.actor("Dev2", 2, event.ActorTypeHuman)
		lead := env.actor("Lead", 3, event.ActorTypeHuman)

		result, err := work.Standup(env.ctx,
			[]types.ActorID{dev1.ID(), dev2.ID()},
			[]string{"finished auth module", "started API tests"},
			lead.ID(), "focus on API coverage today",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("Standup: %v", err)
		}

		if len(result.Updates) != 2 {
			t.Errorf("expected 2 updates, got %d", len(result.Updates))
		}

		ancestors := env.ancestors(result.Priority.ID(), 10)
		if !containsEvent(ancestors, result.Updates[0].ID()) {
			t.Error("priority should trace to updates")
		}
		env.verifyChain()
	})

	t.Run("Retrospective", func(t *testing.T) {
		env := newTestEnv(t)
		work := compositions.NewWorkGrammar(env.grammar)
		dev1 := env.actor("Dev1", 1, event.ActorTypeHuman)
		dev2 := env.actor("Dev2", 2, event.ActorTypeHuman)
		lead := env.actor("Lead", 3, event.ActorTypeHuman)

		task, _ := work.Intend(env.ctx, lead.ID(), "sprint 5 work",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		result, err := work.Retrospective(env.ctx,
			[]types.ActorID{dev1.ID(), dev2.ID()},
			[]string{"CI was slow", "pairing worked well"},
			lead.ID(), "invest in CI pipeline speed",
			task.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Retrospective: %v", err)
		}

		if len(result.Reviews) != 2 {
			t.Errorf("expected 2 reviews, got %d", len(result.Reviews))
		}

		ancestors := env.ancestors(result.Improvement.ID(), 10)
		if !containsEvent(ancestors, task.ID()) {
			t.Error("improvement should trace to reviewed task")
		}
		env.verifyChain()
	})

	t.Run("Triage", func(t *testing.T) {
		env := newTestEnv(t)
		work := compositions.NewWorkGrammar(env.grammar)
		lead := env.actor("Lead", 1, event.ActorTypeHuman)
		dev1 := env.actor("Dev1", 2, event.ActorTypeHuman)
		dev2 := env.actor("Dev2", 3, event.ActorTypeHuman)

		t1, _ := work.Intend(env.ctx, lead.ID(), "fix login bug",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		t2, _ := work.Intend(env.ctx, lead.ID(), "add rate limiting",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		result, err := work.Triage(env.ctx, lead.ID(),
			[]types.EventID{t1.ID(), t2.ID()},
			[]string{"critical", "high"},
			[]types.ActorID{dev1.ID(), dev2.ID()},
			[]types.DomainScope{types.MustDomainScope("auth"), types.MustDomainScope("rate_limiting")},
			[]types.Weight{types.MustWeight(0.9), types.MustWeight(0.6)},
			env.convID, signer)
		if err != nil {
			t.Fatalf("Triage: %v", err)
		}

		if len(result.Priorities) != 2 {
			t.Errorf("expected 2 priorities, got %d", len(result.Priorities))
		}
		if len(result.Assignments) != 2 {
			t.Errorf("expected 2 assignments, got %d", len(result.Assignments))
		}
		if len(result.Scopes) != 2 {
			t.Errorf("expected 2 scopes, got %d", len(result.Scopes))
		}

		ancestors := env.ancestors(result.Scopes[1].ID(), 10)
		if !containsEvent(ancestors, result.Assignments[1].ID()) {
			t.Error("scope should trace to assignment")
		}
		env.verifyChain()
	})
}
