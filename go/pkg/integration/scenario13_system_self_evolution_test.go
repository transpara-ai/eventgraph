package integration_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// TestScenario13_SystemSelfEvolution exercises pattern detection, adaptation, and validation.
// The system detects that 97% of deploy approvals are rubber-stamped, identifies the
// root cause (test coverage < 80%), proposes a mechanical gate, validates it through
// a parallel run, updates the decision tree, and verifies purpose alignment.
func TestScenario13_SystemSelfEvolution(t *testing.T) {
	env := newTestEnv(t)

	patternBot := env.registerActor("PatternBot", 1, event.ActorTypeAI)
	admin := env.registerActor("Admin", 2, event.ActorTypeHuman)

	// 1. Pattern detected — 97% of deploy approvals are approved
	pattern, err := env.grammar.Emit(env.ctx, patternBot.ID(),
		"pattern: 194/200 deploy_staging authority requests approved over 30 days, 97% approval rate",
		env.convID, []types.EventID{env.boot.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit pattern: %v", err)
	}

	// 2. Meta-pattern: all 6 rejections correlate with test coverage < 80%
	metaPattern, err := env.grammar.Derive(env.ctx, patternBot.ID(),
		"meta-pattern: all 6 rejections correlate with test coverage < 80%, no other rejections in 200 requests",
		pattern.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive meta-pattern: %v", err)
	}

	// 3. System dynamic modelled
	systemDynamic, err := env.grammar.Extend(env.ctx, patternBot.ID(),
		"system dynamic: human approval adds 2-15 min latency per deploy, 97% of time the decision is purely mechanical",
		metaPattern.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Extend system dynamic: %v", err)
	}

	// 4. Feedback loop identified
	feedbackLoop, err := env.grammar.Extend(env.ctx, patternBot.ID(),
		"feedback loop (positive/harmful): slow deploys → backlog → cursory reviews → more issues → more reviews → slower deploys",
		systemDynamic.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Extend feedback loop: %v", err)
	}

	// 5. Threshold assessment
	threshold, err := env.grammar.Annotate(env.ctx, patternBot.ID(),
		feedbackLoop.ID(), "threshold",
		"approval rate 97%, threshold for mechanical conversion 98%, approaching safe to convert",
		env.convID, signer)
	if err != nil {
		t.Fatalf("Annotate threshold: %v", err)
	}

	// 6. Adaptation proposed
	adaptation, err := env.grammar.Derive(env.ctx, patternBot.ID(),
		"adaptation proposal: auto-approve deploy_staging when tests pass AND coverage >= 80%, reject otherwise",
		threshold.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive adaptation: %v", err)
	}

	// 7. Authority required for structural change
	authReq, err := env.graph.Record(
		event.EventTypeAuthorityRequested, patternBot.ID(),
		event.AuthorityRequestContent{
			Actor:  patternBot.ID(),
			Action: "modify_decision_tree",
			Level:  event.AuthorityLevelRequired,
		},
		[]types.EventID{adaptation.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Record authority request: %v", err)
	}

	// 8. Admin approves with parallel run condition
	authResolved, err := env.graph.Record(
		event.EventTypeAuthorityResolved, admin.ID(),
		event.AuthorityResolvedContent{
			RequestID: authReq.ID(),
			Approved:  true,
			Resolver:  admin.ID(),
			Reason:    types.None[string](),
		},
		[]types.EventID{authReq.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Record authority resolved: %v", err)
	}

	// 9. Parallel run validates (75 deploys over 3 weeks)
	validation, err := env.grammar.Derive(env.ctx, patternBot.ID(),
		"parallel run results: 75 deploys, mechanical matched human 74/75 cases, fitness 0.987, 1 edge case (empty test suite)",
		authResolved.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive validation: %v", err)
	}

	// 10. Decision tree updated
	treeUpdate, err := env.grammar.Derive(env.ctx, patternBot.ID(),
		"decision tree updated: added mechanical branch — deploy_staging: IF tests_pass AND coverage >= 80% THEN auto_approve ELSE require_human",
		validation.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive tree update: %v", err)
	}

	// 11. Simplification measured
	simplification, err := env.grammar.Extend(env.ctx, patternBot.ID(),
		"simplification: decision complexity reduced from 0.72 to 0.58, human review load reduced by 97%",
		treeUpdate.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Extend simplification: %v", err)
	}

	// 12. System integrity checked
	integrity, err := env.grammar.Annotate(env.ctx, patternBot.ID(),
		simplification.ID(), "integrity",
		"systemic integrity score 0.96, recommendation: monitor for coverage threshold gaming",
		env.convID, signer)
	if err != nil {
		t.Fatalf("Annotate integrity: %v", err)
	}

	// 13. Purpose alignment verified
	purpose, err := env.grammar.Derive(env.ctx, patternBot.ID(),
		"purpose check: system still accountable — mechanical gate is fully auditable, human oversight preserved for edge cases",
		integrity.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive purpose: %v", err)
	}

	// --- Assertions ---

	// Full evolution chain: purpose → integrity → simplification → tree update → validation → authority → adaptation → pattern
	purposeAncestors := env.ancestors(purpose.ID(), 20)
	if !containsEvent(purposeAncestors, integrity.ID()) {
		t.Error("purpose should trace to integrity check")
	}
	if !containsEvent(purposeAncestors, simplification.ID()) {
		t.Error("purpose should trace to simplification")
	}
	if !containsEvent(purposeAncestors, treeUpdate.ID()) {
		t.Error("purpose should trace to tree update")
	}
	if !containsEvent(purposeAncestors, validation.ID()) {
		t.Error("purpose should trace to validation")
	}
	if !containsEvent(purposeAncestors, authResolved.ID()) {
		t.Error("purpose should trace to authority resolution")
	}
	if !containsEvent(purposeAncestors, adaptation.ID()) {
		t.Error("purpose should trace to adaptation proposal")
	}
	if !containsEvent(purposeAncestors, pattern.ID()) {
		t.Error("purpose should trace all the way to original pattern detection")
	}

	// Meta-pattern identifies root cause
	metaAncestors := env.ancestors(metaPattern.ID(), 5)
	if !containsEvent(metaAncestors, pattern.ID()) {
		t.Error("meta-pattern should trace to pattern")
	}

	// Adaptation required authority
	adaptationDesc := env.descendants(adaptation.ID(), 5)
	if !containsEvent(adaptationDesc, authReq.ID()) {
		t.Error("adaptation should lead to authority request")
	}

	// Chain integrity
	env.verifyChain()

	// bootstrap(1) + pattern(1) + metaPattern(1) + systemDynamic(1) + feedbackLoop(1) +
	// threshold(1) + adaptation(1) + authReq(1) + authResolved(1) + validation(1) +
	// treeUpdate(1) + simplification(1) + integrity(1) + purpose(1) = 14
	if count := env.eventCount(); count != 14 {
		t.Errorf("event count = %d, want 14", count)
	}
}
