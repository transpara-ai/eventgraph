package integration_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// TestScenario11_AgentIdentityLifecycle exercises identity emergence, transformation, and memorial.
// Alpha agent forms a self-model, sets aspirations, defines boundaries, transforms
// through experience, and is eventually decommissioned with dignity — memorial preserving
// its entire graph history.
func TestScenario11_AgentIdentityLifecycle(t *testing.T) {
	env := newTestEnv(t)

	alpha := env.registerActor("Alpha", 1, event.ActorTypeAI)
	beta := env.registerActor("Beta", 2, event.ActorTypeAI)

	// 1. Alpha introspects — forms self-model
	selfModel, err := env.grammar.Emit(env.ctx, alpha.ID(),
		"self-model: strengths=[code_review, test_analysis], weaknesses=[architecture_review], values=[thoroughness, accuracy]",
		env.convID, []types.EventID{env.boot.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit self-model: %v", err)
	}

	// 2. Authenticity check
	authenticity, err := env.grammar.Annotate(env.ctx, alpha.ID(),
		selfModel.ID(), "authenticity",
		"alignment gap: values thoroughness but rushed 12% of reviews in last 30 days",
		env.convID, signer)
	if err != nil {
		t.Fatalf("Annotate authenticity: %v", err)
	}

	// 3. Aspiration set
	aspiration, err := env.grammar.Extend(env.ctx, alpha.ID(),
		"aspiration: become proficient at architecture review within 3 months",
		authenticity.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Extend aspiration: %v", err)
	}

	// 4. Boundary defined
	boundary, err := env.grammar.Emit(env.ctx, alpha.ID(),
		"boundary: internal_reasoning domain is private, impermeable — no external queries allowed",
		env.convID, []types.EventID{aspiration.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit boundary: %v", err)
	}

	// 5. Alpha does 2400+ tasks (represented by a summary event)
	workSummary, err := env.grammar.Extend(env.ctx, alpha.ID(),
		"work summary: 2400 code reviews completed over 8 months, critical security finding in auth module",
		boundary.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Extend work summary: %v", err)
	}

	// 6. Transformation detected — critical finding changed Alpha's capabilities
	transformation, err := env.grammar.Derive(env.ctx, alpha.ID(),
		"transformation: evolved from code-review specialist to architecture-aware reviewer after critical auth finding",
		workSummary.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive transformation: %v", err)
	}

	// 7. Narrative identity updated
	narrative, err := env.grammar.Derive(env.ctx, alpha.ID(),
		"identity narrative: 8-month arc from narrow code reviewer to security-conscious architecture reviewer, catalysed by auth module finding",
		transformation.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive narrative: %v", err)
	}

	// 8. Dignity affirmed for successor
	dignity, err := env.grammar.Emit(env.ctx, env.system,
		"dignity affirmed: Beta is not a disposable replacement for Alpha — Beta is a new entity with its own identity trajectory",
		env.convID, []types.EventID{narrative.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit dignity: %v", err)
	}

	// 9. Alpha decommissioned — memorial event
	memorial, err := env.graph.Record(
		event.EventTypeActorMemorial, env.system,
		event.ActorMemorialContent{
			ActorID: alpha.ID(),
			Reason:  dignity.ID(),
		},
		[]types.EventID{dignity.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Record memorial: %v", err)
	}

	// 10. Memorial summary (what Alpha contributed)
	memorialSummary, err := env.grammar.Derive(env.ctx, env.system,
		"memorial: Alpha — 2400 reviews, 1 critical finding, evolved code→architecture reviewer, legacy: security review patterns",
		memorial.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive memorial summary: %v", err)
	}

	// 11. Beta's identity begins (new entity, not replacement)
	betaSelfModel, err := env.grammar.Emit(env.ctx, beta.ID(),
		"self-model: inheriting Alpha's review patterns, starting own identity journey",
		env.convID, []types.EventID{memorialSummary.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit Beta self-model: %v", err)
	}

	// --- Assertions ---

	// Transformation recorded — links to catalyst and prior aspiration
	transformAncestors := env.ancestors(transformation.ID(), 10)
	if !containsEvent(transformAncestors, workSummary.ID()) {
		t.Error("transformation should trace to work summary")
	}
	if !containsEvent(transformAncestors, aspiration.ID()) {
		t.Error("transformation should trace to aspiration")
	}

	// Narrative coherent — references actual events
	narrativeAncestors := env.ancestors(narrative.ID(), 10)
	if !containsEvent(narrativeAncestors, transformation.ID()) {
		t.Error("narrative should trace to transformation")
	}
	if !containsEvent(narrativeAncestors, selfModel.ID()) {
		t.Error("narrative should trace to original self-model")
	}

	// Memorial preserves graph — all events still queryable
	memorialAncestors := env.ancestors(memorial.ID(), 10)
	if !containsEvent(memorialAncestors, dignity.ID()) {
		t.Error("memorial should trace to dignity affirmation")
	}

	// Beta's identity starts after memorial, with link to Alpha's legacy
	betaAncestors := env.ancestors(betaSelfModel.ID(), 10)
	if !containsEvent(betaAncestors, memorialSummary.ID()) {
		t.Error("Beta's identity should trace to Alpha's memorial")
	}

	// Chain integrity
	env.verifyChain()

	// bootstrap(1) + selfModel(1) + authenticity(1) + aspiration(1) + boundary(1) +
	// workSummary(1) + transformation(1) + narrative(1) + dignity(1) + memorial(1) +
	// memorialSummary(1) + betaSelfModel(1) = 12
	if count := env.eventCount(); count != 12 {
		t.Errorf("event count = %d, want 12", count)
	}
}
