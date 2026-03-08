package integration_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/compositions"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// TestScenario20_KnowledgeEcosystem exercises knowledge management across cultures:
// Build a knowledge base → survey existing knowledge → transfer to new team →
// cultural onboarding for international expansion → design review of the system →
// forecast future needs.
// Crosses Knowledge and Meaning grammars.
func TestScenario20_KnowledgeEcosystem(t *testing.T) {
	env := newTestEnv(t)
	knowledge := compositions.NewKnowledgeGrammar(env.grammar)
	meaning := compositions.NewMeaningGrammar(env.grammar)

	architect := env.registerActor("Architect", 1, event.ActorTypeHuman)
	researcher := env.registerActor("Researcher", 2, event.ActorTypeAI)
	newcomer := env.registerActor("TokyoLead", 3, event.ActorTypeHuman)

	// 1. Build a knowledge base of architectural decisions
	kb, err := knowledge.KnowledgeBase(env.ctx, architect.ID(),
		[]string{
			"event sourcing chosen over CRUD for auditability",
			"Ed25519 chosen over RSA for signature performance",
			"append-only store prevents tampering",
		},
		[]string{"architecture.patterns", "architecture.security", "architecture.integrity"},
		"core architectural decisions Q1 2026",
		[]types.EventID{env.boot.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("KnowledgeBase: %v", err)
	}
	if len(kb.Claims) != 3 {
		t.Fatalf("expected 3 claims, got %d", len(kb.Claims))
	}

	// 2. Survey existing knowledge to find patterns
	survey, err := knowledge.Survey(env.ctx, researcher.ID(),
		[]string{
			"what patterns emerge from our architectural decisions?",
			"what security properties does the current design guarantee?",
			"what are the performance characteristics of our choices?",
		},
		"all decisions prioritize verifiability over convenience",
		"the architecture optimizes for trust minimization — every claim is independently verifiable",
		[]types.EventID{kb.Memory.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Survey: %v", err)
	}
	if len(survey.Recalls) != 3 {
		t.Fatalf("expected 3 recalls, got %d", len(survey.Recalls))
	}

	// 3. Transfer knowledge to new team context
	transfer, err := knowledge.Transfer(env.ctx, architect.ID(),
		"core architectural principles for new Tokyo office",
		"translated to Japanese engineering conventions, mapped to local compliance requirements",
		"Tokyo team now understands event sourcing in context of J-SOX compliance",
		[]types.EventID{survey.Synthesis.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Transfer: %v", err)
	}

	// 4. Cultural onboarding — prepare newcomer for different engineering culture
	onboarding, err := meaning.CulturalOnboarding(env.ctx, architect.ID(), newcomer.ID(),
		"Western direct feedback style → Japanese nemawashi consensus-building",
		types.Some(types.MustDomainScope("engineering_culture")),
		"the consensus process feels slower but produces more durable decisions",
		transfer.Learn.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("CulturalOnboarding: %v", err)
	}

	// 5. Design review of the knowledge transfer system itself
	designReview, err := meaning.DesignReview(env.ctx, architect.ID(),
		"the knowledge graph's self-referential structure is elegant — it documents its own architecture",
		"viewing knowledge transfer as a graph problem rather than a document problem",
		"does our transfer process preserve tacit knowledge or only explicit claims?",
		"explicit knowledge transfers well; tacit knowledge requires mentorship, not documents",
		onboarding.Examination.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("DesignReview: %v", err)
	}

	// 6. Forecast — where is the knowledge ecosystem heading?
	forecast, err := meaning.Forecast(env.ctx, researcher.ID(),
		"at current growth, knowledge base will reach 10k claims by Q3 — need automated categorization",
		"assumes linear claim growth and stable team size — may underestimate if Tokyo ramps faster",
		"high confidence: need automated categorization within 6 months, medium confidence: need multi-language support within 12",
		[]types.EventID{designReview.Wisdom.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Forecast: %v", err)
	}

	// --- Assertions ---

	// Forecast traces all the way back to knowledge base
	forecastAncestors := env.ancestors(forecast.Wisdom.ID(), 30)
	if !containsEvent(forecastAncestors, kb.Memory.ID()) {
		t.Error("forecast should trace to knowledge base")
	}

	// Design review traces through cultural onboarding to transfer
	reviewAncestors := env.ancestors(designReview.Wisdom.ID(), 20)
	if !containsEvent(reviewAncestors, transfer.Learn.ID()) {
		t.Error("design review should trace to knowledge transfer")
	}

	// Cultural onboarding traces through transfer to survey
	onboardAncestors := env.ancestors(onboarding.Examination.ID(), 20)
	if !containsEvent(onboardAncestors, survey.Synthesis.ID()) {
		t.Error("cultural onboarding should trace to survey synthesis")
	}

	// Survey traces to knowledge base
	surveyAncestors := env.ancestors(survey.Synthesis.ID(), 15)
	if !containsEvent(surveyAncestors, kb.Memory.ID()) {
		t.Error("survey should trace to knowledge base memory")
	}

	env.verifyChain()
}
