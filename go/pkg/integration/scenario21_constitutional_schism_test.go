package integration_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/compositions"
	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// TestScenario21_ConstitutionalSchism exercises governance crisis and evolution:
// Constitutional amendment fails → community schism → splinter group barters
// for shared infrastructure → system prunes abandoned structures.
// Crosses Justice, Social, Market, and Evolution grammars.
func TestScenario21_ConstitutionalSchism(t *testing.T) {
	env := newTestEnv(t)
	justice := compositions.NewJusticeGrammar(env.grammar)
	social := compositions.NewSocialGrammar(env.grammar)
	market := compositions.NewMarketGrammar(env.grammar)
	evolution := compositions.NewEvolutionGrammar(env.grammar)

	founder := env.registerActor("Founder", 1, event.ActorTypeHuman)
	reformer := env.registerActor("Reformer", 2, event.ActorTypeHuman)
	conservative := env.registerActor("Conservative", 3, event.ActorTypeHuman)
	sysBot := env.registerActor("SystemBot", 4, event.ActorTypeAI)

	// 1. Establish initial law
	law, _ := justice.Legislate(env.ctx, founder.ID(),
		"all governance decisions require unanimous consent",
		[]types.EventID{env.boot.ID()}, env.convID, signer)

	// 2. Constitutional amendment proposed — change from unanimous to 2/3 majority
	amendment, err := justice.ConstitutionalAmendment(env.ctx, reformer.ID(),
		"unanimous consent blocks progress — propose 2/3 supermajority threshold",
		"governance decisions require 2/3 supermajority instead of unanimity",
		"rights preserved: individual veto retained for membership and expulsion decisions",
		law.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("ConstitutionalAmendment: %v", err)
	}

	// 3. Amendment passes but causes schism — conservative faction splits
	// First create a subscription to sever
	sub, _ := env.grammar.Subscribe(env.ctx, conservative.ID(), founder.ID(),
		types.Some(types.MustDomainScope("governance")),
		amendment.RightsCheck.ID(), env.convID, signer)
	edgeID, _ := types.NewEdgeID(sub.ID().Value())

	schism, err := social.Schism(env.ctx, conservative.ID(), founder.ID(),
		"reject supermajority — unanimity is the only legitimate standard",
		types.MustDomainScope("governance"),
		edgeID, "irreconcilable governance philosophy differences",
		amendment.RightsCheck.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Schism: %v", err)
	}

	// 4. Splinter community barters for shared infrastructure
	barter, err := market.Barter(env.ctx, conservative.ID(), founder.ID(),
		"continued access to shared event store for 6 months",
		"historical governance data export in standard format",
		types.MustDomainScope("infrastructure"),
		[]types.EventID{schism.NewCommunity.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Barter: %v", err)
	}

	// 5. System prunes abandoned governance structures
	prune, err := evolution.Prune(env.ctx, sysBot.ID(),
		"unanimous consent voting module — zero invocations since amendment",
		"removed unanimous consent module, replaced with supermajority",
		"all 34 governance tests pass without unanimous module",
		[]types.EventID{barter.Acceptance.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}

	// --- Assertions ---

	// Prune traces all the way back to the original law
	pruneAncestors := env.ancestors(prune.Verification.ID(), 25)
	if !containsEvent(pruneAncestors, law.ID()) {
		t.Error("prune should trace to original law")
	}

	// Barter traces through schism to amendment
	barterAncestors := env.ancestors(barter.Acceptance.ID(), 20)
	if !containsEvent(barterAncestors, amendment.Reform.ID()) {
		t.Error("barter should trace to constitutional amendment")
	}

	// Schism traces to amendment rights check
	schismAncestors := env.ancestors(schism.NewCommunity.ID(), 15)
	if !containsEvent(schismAncestors, amendment.RightsCheck.ID()) {
		t.Error("schism should trace to amendment rights check")
	}

	// Amendment traces to original law
	amendAncestors := env.ancestors(amendment.RightsCheck.ID(), 10)
	if !containsEvent(amendAncestors, law.ID()) {
		t.Error("amendment should trace to original law")
	}

	env.verifyChain()
}
