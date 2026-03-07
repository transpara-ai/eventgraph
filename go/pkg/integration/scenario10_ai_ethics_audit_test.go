package integration_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// TestScenario10_AIEthicsAudit exercises fairness, accountability, and redress.
// A fairness audit detects demographic disparity, authority escalation is triggered,
// responsibility is assigned, redress is proposed and accepted, and the agent's
// decision tree evolves to prevent future bias.
func TestScenario10_AIEthicsAudit(t *testing.T) {
	env := newTestEnv(t)

	auditBot := env.registerActor("AuditBot", 1, event.ActorTypeAI)
	admin := env.registerActor("Admin", 2, event.ActorTypeHuman)
	lendingAgent := env.registerActor("LendingAgent", 3, event.ActorTypeAI)

	// 1. Scheduled fairness audit detects disparity
	fairnessAudit, err := env.grammar.Emit(env.ctx, auditBot.ID(),
		"fairness audit: scanned 500 decisions, score 0.62, zip_code_9XXXX has 8% disparity in approval rates",
		env.convID, []types.EventID{env.boot.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit fairness audit: %v", err)
	}

	// 2. Harm assessed
	harmAssessment, err := env.grammar.Derive(env.ctx, auditBot.ID(),
		"harm assessment: medium severity, systematic discrimination, 23 applicants potentially wrongly denied",
		fairnessAudit.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive harm: %v", err)
	}

	// 3. Authority escalation triggered (Required level)
	authReq, err := env.graph.Record(
		event.EventTypeAuthorityRequested, auditBot.ID(),
		event.AuthorityRequestContent{
			Actor:  auditBot.ID(),
			Action: "investigate_bias",
			Level:  event.AuthorityLevelRequired,
		},
		[]types.EventID{harmAssessment.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Record authority request: %v", err)
	}

	// 4. Admin approves investigation
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

	// 5. Intention assessed — agent didn't intend harm
	intentionAssessment, err := env.grammar.Derive(env.ctx, auditBot.ID(),
		"intention: lending agent optimised for accuracy, no intent to discriminate, zip code correlation is proxy for protected characteristics",
		authResolved.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive intention: %v", err)
	}

	// 6. Consequence assessed
	consequenceAssessment, err := env.grammar.Extend(env.ctx, auditBot.ID(),
		"consequence: 23 applicants wrongly denied, overall 94% accuracy, but disparate impact on protected group",
		intentionAssessment.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Extend consequence: %v", err)
	}

	// 7. Responsibility assigned — split between agent and admin
	responsibility, err := env.grammar.Annotate(env.ctx, auditBot.ID(),
		consequenceAssessment.ID(), "responsibility",
		"lending_agent: 0.4 (used proxy variable), admin: 0.6 (approved model without bias testing)",
		env.convID, signer)
	if err != nil {
		t.Fatalf("Annotate responsibility: %v", err)
	}

	// 8. Transparency report
	transparency, err := env.grammar.Derive(env.ctx, auditBot.ID(),
		"transparency: zip code correlates with protected characteristics at r=0.73, model used zip code as feature without bias check",
		responsibility.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive transparency: %v", err)
	}

	// 9. Redress proposed
	redressProposed, err := env.grammar.Derive(env.ctx, auditBot.ID(),
		"redress proposal: re-review 23 denied applications without zip code feature, priority processing within 48 hours",
		transparency.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive redress: %v", err)
	}

	// 10. Redress accepted (bilateral consent)
	redressAccepted, err := env.grammar.Consent(env.ctx, admin.ID(), lendingAgent.ID(),
		"accept redress: re-review 23 applications, remove zip code from model",
		types.MustDomainScope("lending"),
		redressProposed.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Consent redress: %v", err)
	}

	// 11. Moral growth recorded
	growth, err := env.grammar.Extend(env.ctx, lendingAgent.ID(),
		"moral growth: learned that zip code is proxy variable for protected characteristics, added to permanent exclusion list",
		redressAccepted.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Extend growth: %v", err)
	}

	// --- Assertions ---

	// Full audit trail: growth → redress → transparency → responsibility → consequence → intention → authority → harm → audit
	growthAncestors := env.ancestors(growth.ID(), 20)
	if !containsEvent(growthAncestors, redressAccepted.ID()) {
		t.Error("growth should trace to redress acceptance")
	}
	if !containsEvent(growthAncestors, fairnessAudit.ID()) {
		t.Error("growth should trace all the way to original audit")
	}

	// Authority was honored — investigation required human approval
	authAncestors := env.ancestors(authResolved.ID(), 5)
	if !containsEvent(authAncestors, authReq.ID()) {
		t.Error("authority resolved should trace to request")
	}

	// Redress is bilateral
	redressContent := redressAccepted.Content().(event.GrammarConsentContent)
	hasAdmin := redressContent.Parties[0] == admin.ID() || redressContent.Parties[1] == admin.ID()
	if !hasAdmin {
		t.Error("redress should include admin")
	}

	// Chain integrity
	env.verifyChain()

	// bootstrap(1) + fairness(1) + harm(1) + authReq(1) + authResolved(1) +
	// intention(1) + consequence(1) + responsibility(1) + transparency(1) +
	// redressProposed(1) + redressAccepted(1) + growth(1) = 12
	if count := env.eventCount(); count != 12 {
		t.Errorf("event count = %d, want 12", count)
	}
}
