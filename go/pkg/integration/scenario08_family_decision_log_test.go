package integration_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// TestScenario08_FamilyDecisionLog exercises consensual domestic decision making.
// Maria proposes buying a house, an AI advisor researches under delegation,
// family discusses, both parents make a bilateral decision, and years later
// the process is queryable as precedent.
func TestScenario08_FamilyDecisionLog(t *testing.T) {
	env := newTestEnv(t)

	maria := env.registerActor("Maria", 1, event.ActorTypeHuman)
	james := env.registerActor("James", 2, event.ActorTypeHuman)
	sophie := env.registerActor("Sophie", 3, event.ActorTypeHuman)
	advisor := env.registerActor("AIAdvisor", 4, event.ActorTypeAI)

	// 1. Maria proposes buying a house
	proposal, err := env.grammar.Emit(env.ctx, maria.ID(),
		"proposal: buy a house in Eastside neighbourhood, budget $450K",
		env.convID, []types.EventID{env.boot.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit proposal: %v", err)
	}

	// 2. James delegates research to AI advisor
	delegation, err := env.grammar.Delegate(env.ctx, james.ID(), advisor.ID(),
		types.MustDomainScope("market_research"), types.MustWeight(0.7),
		proposal.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Delegate research: %v", err)
	}

	// 3. Advisor researches market (Derive under delegation)
	research, err := env.grammar.Derive(env.ctx, advisor.ID(),
		"research: Eastside median $440K, rent $2200/mo, mortgage $2400/mo at current rates, break-even 5 years, confidence 0.82",
		delegation.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive research: %v", err)
	}

	// 4. Sophie shares perspective
	sophieView, err := env.grammar.Respond(env.ctx, sophie.ID(),
		"I support it IF I get my own room. Current apartment sharing is hard for studying.",
		proposal.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Respond Sophie: %v", err)
	}

	// 5. James raises concern (informed by research)
	jamesConcern, err := env.grammar.Respond(env.ctx, james.ID(),
		"concern: mortgage is $200/mo more than rent, tight on single income months",
		research.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Respond James concern: %v", err)
	}

	// 6. Maria addresses concern
	mariaResponse, err := env.grammar.Respond(env.ctx, maria.ID(),
		"response: we can use the $15K savings buffer, and break-even is 5 years — we plan to stay 10+",
		jamesConcern.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Respond Maria: %v", err)
	}

	// 7. Maria + James make bilateral decision (Consent)
	decision, err := env.grammar.Consent(env.ctx, maria.ID(), james.ID(),
		"decision: buy house in Eastside, budget $450K, conditions: Sophie gets own room, maintain 3-month emergency fund",
		types.MustDomainScope("family_finance"),
		mariaResponse.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Consent decision: %v", err)
	}

	// --- Assertions ---

	// Consultation recorded: Sophie's opinion in decision's causal chain
	decisionAncestors := env.ancestors(decision.ID(), 10)
	// Decision traces through Maria's response → James's concern → research → delegation → proposal
	if !containsEvent(decisionAncestors, mariaResponse.ID()) {
		t.Error("decision should trace to Maria's response")
	}
	if !containsEvent(decisionAncestors, jamesConcern.ID()) {
		t.Error("decision should trace to James's concern")
	}
	if !containsEvent(decisionAncestors, research.ID()) {
		t.Error("decision should trace to research")
	}
	if !containsEvent(decisionAncestors, proposal.ID()) {
		t.Error("decision should trace to original proposal")
	}

	// Sophie's view is reachable through the proposal's descendants
	proposalDescendants := env.descendants(proposal.ID(), 5)
	if !containsEvent(proposalDescendants, sophieView.ID()) {
		t.Error("proposal descendants should include Sophie's perspective")
	}

	// AI contribution scoped — delegation has domain scope
	delegationContent := delegation.Content().(event.EdgeCreatedContent)
	if !delegationContent.Scope.IsSome() {
		t.Error("delegation should have domain scope")
	}
	scope := delegationContent.Scope.Unwrap()
	if scope.Value() != "market_research" {
		t.Errorf("delegation scope = %v, want market_research", scope.Value())
	}

	// Decision is bilateral — consent content has both parties
	consentContent := decision.Content().(event.GrammarConsentContent)
	parties := consentContent.Parties
	hasMaria := parties[0] == maria.ID() || parties[1] == maria.ID()
	hasJames := parties[0] == james.ID() || parties[1] == james.ID()
	if !hasMaria {
		t.Error("consent should include Maria")
	}
	if !hasJames {
		t.Error("consent should include James")
	}

	// Chain integrity
	env.verifyChain()

	// bootstrap(1) + proposal(1) + delegation(1) + research(1) + sophie(1) +
	// jamesConcern(1) + mariaResponse(1) + decision(1) = 8
	if count := env.eventCount(); count != 8 {
		t.Errorf("event count = %d, want 8", count)
	}
}
