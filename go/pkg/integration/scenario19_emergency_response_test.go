package integration_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/compositions"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// TestScenario19_EmergencyResponse exercises crisis management:
// Security breach detected → triage incoming issues → emergency injunction →
// plea deal for minor actor → emergency migration to patched system.
// Crosses Work, Justice, and Build grammars.
func TestScenario19_EmergencyResponse(t *testing.T) {
	env := newTestEnv(t)
	work := compositions.NewWorkGrammar(env.grammar)
	justice := compositions.NewJusticeGrammar(env.grammar)
	build := compositions.NewBuildGrammar(env.grammar)

	secLead := env.registerActor("SecurityLead", 1, event.ActorTypeHuman)
	dev1 := env.registerActor("Dev1", 2, event.ActorTypeHuman)
	dev2 := env.registerActor("Dev2", 3, event.ActorTypeHuman)
	judge := env.registerActor("CISO", 4, event.ActorTypeHuman)
	executor := env.registerActor("OpsBot", 5, event.ActorTypeAI)
	minorActor := env.registerActor("ContractorBot", 6, event.ActorTypeAI)

	// 1. Security breach — multiple issues need triage
	issue1, _ := env.grammar.Emit(env.ctx, secLead.ID(),
		"CVE-2026-1234: auth bypass in API gateway",
		env.convID, []types.EventID{env.boot.ID()}, signer)
	issue2, _ := env.grammar.Emit(env.ctx, secLead.ID(),
		"CVE-2026-1235: SQL injection in search endpoint",
		env.convID, []types.EventID{env.boot.ID()}, signer)

	// 2. Triage — prioritize and assign both issues
	triage, err := work.Triage(env.ctx, secLead.ID(),
		[]types.EventID{issue1.ID(), issue2.ID()},
		[]string{"P0: auth bypass, actively exploited", "P1: SQL injection, no evidence of exploitation"},
		[]types.ActorID{dev1.ID(), dev2.ID()},
		[]types.DomainScope{types.MustDomainScope("auth"), types.MustDomainScope("search")},
		[]types.Weight{types.MustWeight(1.0), types.MustWeight(0.8)},
		env.convID, signer)
	if err != nil {
		t.Fatalf("Triage: %v", err)
	}
	if len(triage.Priorities) != 2 {
		t.Fatalf("expected 2 priorities, got %d", len(triage.Priorities))
	}

	// 3. Emergency injunction — block all API access until auth is patched
	injunction, err := justice.Injunction(env.ctx, secLead.ID(), judge.ID(), executor.ID(),
		"auth bypass allows unauthenticated access to all API endpoints",
		"block all external API traffic pending auth patch",
		types.MustDomainScope("api_access"), types.MustWeight(1.0),
		triage.Priorities[0].ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Injunction: %v", err)
	}

	// 4. Plea deal — contractor bot that introduced the vulnerability accepts reduced penalty
	plea, err := justice.Plea(env.ctx, minorActor.ID(), secLead.ID(), executor.ID(),
		"introduced auth bypass through misconfigured middleware",
		"accept restricted scope: read-only access for 30 days, mandatory security training",
		types.MustDomainScope("api_development"), types.MustWeight(0.3),
		injunction.Ruling.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Plea: %v", err)
	}

	// 5. Emergency migration — deploy patched auth system
	oldSystem, _ := env.grammar.Emit(env.ctx, dev1.ID(),
		"current auth system v2.3.1",
		env.convID, []types.EventID{injunction.Enforcement.ID()}, signer)

	migration, err := build.Migration(env.ctx, dev1.ID(),
		oldSystem.ID(),
		"migrate to auth v2.4.0 with CVE-2026-1234 fix",
		"v2.4.0",
		"deployed to production with zero-downtime rolling update",
		"all 156 auth tests pass, penetration test confirms fix",
		env.convID, signer)
	if err != nil {
		t.Fatalf("Migration: %v", err)
	}

	// --- Assertions ---

	// Migration traces back to triage
	migrationAncestors := env.ancestors(migration.Test.ID(), 20)
	if !containsEvent(migrationAncestors, triage.Priorities[0].ID()) {
		t.Error("migration should trace to triage priority")
	}

	// Plea traces through injunction
	pleaAncestors := env.ancestors(plea.Enforcement.ID(), 15)
	if !containsEvent(pleaAncestors, injunction.Filing.ID()) {
		t.Error("plea should trace to injunction filing")
	}

	// Injunction traces to triage
	injAncestors := env.ancestors(injunction.Enforcement.ID(), 10)
	if !containsEvent(injAncestors, triage.Priorities[0].ID()) {
		t.Error("injunction should trace to triage")
	}

	env.verifyChain()
}
