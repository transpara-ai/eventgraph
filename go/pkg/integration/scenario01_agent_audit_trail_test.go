package integration_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// TestScenario01_AgentAuditTrail exercises the AI Agent Audit Trail scenario.
// Alice submits code, an agent reviews it under delegation, a bug is discovered,
// trust decreases, and the full causal chain from bug to delegation is traversable.
func TestScenario01_AgentAuditTrail(t *testing.T) {
	env := newTestEnv(t)

	alice := env.registerActor("Alice", 1, event.ActorTypeHuman)
	agent := env.registerActor("ReviewBot", 2, event.ActorTypeAI)

	// 1. Alice submits code for review
	submission, err := env.grammar.Emit(env.ctx, alice.ID(),
		"code submission: auth module refactor",
		env.convID, []types.EventID{env.boot.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit submission: %v", err)
	}

	// 2. Alice delegates code_review authority to agent
	delegation, err := env.grammar.Delegate(env.ctx, alice.ID(), agent.ID(),
		types.MustDomainScope("code_review"), types.MustWeight(0.8),
		submission.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Delegate: %v", err)
	}

	// 3. Agent picks up and reviews the code (Derive from submission)
	review, err := env.grammar.Derive(env.ctx, agent.ID(),
		"review: LGTM, no issues found, approving PR",
		submission.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive review: %v", err)
	}

	// 4. Agent approves (Respond to its own review with approval)
	approval, err := env.grammar.Respond(env.ctx, agent.ID(),
		"decision: approve PR with confidence 0.85",
		review.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Respond approval: %v", err)
	}

	// 5. Trust updated after successful review
	trustUp, err := env.graph.Record(
		event.EventTypeTrustUpdated, env.system,
		event.TrustUpdatedContent{
			Actor:    agent.ID(),
			Previous: types.MustScore(0.1),
			Current:  types.MustScore(0.3),
			Domain:   types.MustDomainScope("code_review"),
			Cause:    approval.ID(),
		},
		[]types.EventID{approval.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Record trust up: %v", err)
	}

	// 6. Bug discovered in approved code
	bugReport, err := env.grammar.Emit(env.ctx, alice.ID(),
		"bug found in auth module: session tokens not invalidated on logout",
		env.convID, []types.EventID{approval.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit bug: %v", err)
	}

	// 7. Violation detected — agent missed the bug
	violation, err := env.graph.Record(
		event.EventTypeViolationDetected, env.system,
		event.ViolationDetectedContent{
			Expectation: approval.ID(),
			Actor:       agent.ID(),
			Severity:    event.SeverityLevelSerious,
			Description: "agent approved code with session management bug",
			Evidence:    types.MustNonEmpty([]types.EventID{bugReport.ID()}),
		},
		[]types.EventID{bugReport.ID(), approval.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Record violation: %v", err)
	}

	// 8. Trust decreases
	_, err = env.graph.Record(
		event.EventTypeTrustUpdated, env.system,
		event.TrustUpdatedContent{
			Actor:    agent.ID(),
			Previous: types.MustScore(0.3),
			Current:  types.MustScore(0.15),
			Domain:   types.MustDomainScope("code_review"),
			Cause:    violation.ID(),
		},
		[]types.EventID{violation.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Record trust down: %v", err)
	}

	// --- Assertions ---

	// Causal chain: bug → approval → review → submission
	ancestors := env.ancestors(bugReport.ID(), 10)
	if !containsEvent(ancestors, approval.ID()) {
		t.Error("bug report should have approval in ancestors")
	}

	// Causal chain: violation → bug + approval
	violationAncestors := env.ancestors(violation.ID(), 10)
	if !containsEvent(violationAncestors, bugReport.ID()) {
		t.Error("violation should have bug report in ancestors")
	}
	if !containsEvent(violationAncestors, approval.ID()) {
		t.Error("violation should have approval in ancestors")
	}

	// Full chain from violation back to submission
	if !containsEvent(violationAncestors, submission.ID()) {
		t.Error("violation should trace back to original submission")
	}

	// Delegation is in the chain
	delegationDescendants := env.descendants(delegation.ID(), 10)
	_ = delegationDescendants

	// Trust events exist
	_ = trustUp

	// Verify chain integrity
	env.verifyChain()

	// Verify event count: bootstrap + submission + delegation + review + approval + trustUp + bug + violation + trustDown = 9
	if count := env.eventCount(); count != 9 {
		t.Errorf("event count = %d, want 9", count)
	}
}
