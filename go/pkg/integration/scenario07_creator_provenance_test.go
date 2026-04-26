package integration_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// TestScenario07_CreatorProvenance exercises human vs AI content distinction.
// Kai studies Luna's work, creates iterative drafts with feedback, and publishes.
// The rich Derive chain is structurally distinguishable from AI-generated content
// which has no creative history.
func TestScenario07_CreatorProvenance(t *testing.T) {
	env := newTestEnv(t)

	kai := env.registerActor("Kai", 1, event.ActorTypeHuman)
	luna := env.registerActor("Luna", 2, event.ActorTypeHuman)
	aiGen := env.registerActor("AIGenerator", 3, event.ActorTypeAI)

	// --- Human creative process ---

	// 1. Luna's existing work (already on graph)
	lunasWork, err := env.grammar.Emit(env.ctx, luna.ID(),
		"artwork: Digital landscape, watercolour technique, 2025",
		env.convID, []types.EventID{env.boot.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit Luna's work: %v", err)
	}

	// 2. Kai encounters and annotates Luna's work
	inspiration, err := env.grammar.Annotate(env.ctx, kai.ID(),
		lunasWork.ID(), "inspiration",
		"technique: layered transparency creates depth without weight",
		env.convID, signer)
	if err != nil {
		t.Fatalf("Annotate inspiration: %v", err)
	}

	// 3. Kai studies the technique (3 hours of practice)
	study, err := env.grammar.Derive(env.ctx, kai.ID(),
		"study: practiced layered transparency technique for 3 hours, 12 practice pieces",
		inspiration.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive study: %v", err)
	}

	// 4. Kai creates draft 1
	draft1, err := env.grammar.Derive(env.ctx, kai.ID(),
		"draft 1: mountain landscape using layered transparency, artifact hash: sha256:draft1abc",
		study.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive draft1: %v", err)
	}

	// 5. Kai requests feedback from Luna
	feedbackReq, err := env.grammar.Channel(env.ctx, kai.ID(), luna.ID(),
		types.Some(types.MustDomainScope("art")),
		draft1.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Channel feedback request: %v", err)
	}

	// 6. Luna responds with guidance
	feedback, err := env.grammar.Respond(env.ctx, luna.ID(),
		"feedback: the foreground layers are too opaque, try reducing opacity to 40% for depth",
		feedbackReq.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Respond feedback: %v", err)
	}

	// 7. Kai creates draft 2 incorporating feedback
	draft2, err := env.grammar.Derive(env.ctx, kai.ID(),
		"draft 2: revised with 40% opacity foreground, artifact hash: sha256:draft2def",
		feedback.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive draft2: %v", err)
	}

	// 8. Kai publishes final work
	published, err := env.grammar.Derive(env.ctx, kai.ID(),
		"published: Mountain Dawn, digital landscape, influenced by Luna's transparency technique",
		draft2.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive published: %v", err)
	}

	// 9. Luna endorses Kai's work
	_, err = env.grammar.Endorse(env.ctx, luna.ID(),
		published.ID(), kai.ID(), types.MustWeight(0.6),
		types.Some(types.MustDomainScope("art")),
		env.convID, signer)
	if err != nil {
		t.Fatalf("Endorse: %v", err)
	}

	// --- AI-generated content (contrast) ---

	// Single event with no creative history
	aiContent, err := env.grammar.Emit(env.ctx, aiGen.ID(),
		"generated: Mountain landscape, digital art",
		env.convID, []types.EventID{env.boot.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit AI content: %v", err)
	}

	// --- Assertions ---

	// Human creative chain: published ← draft2 ← feedback ← draft1 ← study ← inspiration ← Luna's work
	publishedAncestors := env.ancestors(published.ID(), 10)
	if !containsEvent(publishedAncestors, draft2.ID()) {
		t.Error("published should trace to draft2")
	}
	if !containsEvent(publishedAncestors, feedback.ID()) {
		t.Error("published should trace to Luna's feedback")
	}
	if !containsEvent(publishedAncestors, draft1.ID()) {
		t.Error("published should trace to draft1")
	}
	if !containsEvent(publishedAncestors, study.ID()) {
		t.Error("published should trace to study")
	}
	if !containsEvent(publishedAncestors, inspiration.ID()) {
		t.Error("published should trace to inspiration")
	}
	if !containsEvent(publishedAncestors, lunasWork.ID()) {
		t.Error("published should trace to Luna's original work")
	}

	// AI content has NO creative chain — only traces to bootstrap
	aiAncestors := env.ancestors(aiContent.ID(), 10)
	if len(aiAncestors) != 1 {
		t.Errorf("AI content ancestors = %d, want 1 (bootstrap only)", len(aiAncestors))
	}

	// Structurally distinguishable: human work has 6+ ancestors, AI has 1
	if len(publishedAncestors) <= len(aiAncestors) {
		t.Error("human creative work should have more ancestors than AI-generated content")
	}

	// Chain integrity
	env.verifyChain()

	// bootstrap(1) + lunasWork(1) + inspiration(1) + study(1) + draft1(1) +
	// feedbackReq(1) + feedback(1) + draft2(1) + published(1) + endorse(1) +
	// aiContent(1) = 11
	if count := env.eventCount(); count != 11 {
		t.Errorf("event count = %d, want 11", count)
	}
}
