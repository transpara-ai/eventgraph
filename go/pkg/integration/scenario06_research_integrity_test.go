package integration_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// TestScenario06_ResearchIntegrity exercises pre-registration and analysis audit trails.
// Grace pre-registers a hypothesis, collects data, runs analyses (including a failed one),
// submits a manuscript, gets peer reviews, and publishes. The full chain proves
// pre-registration and makes both successful and failed analyses visible.
func TestScenario06_ResearchIntegrity(t *testing.T) {
	env := newTestEnv(t)

	grace := env.registerActor("Grace", 1, event.ActorTypeHuman)
	henry := env.registerActor("Henry", 2, event.ActorTypeHuman) // reviewer
	iris := env.registerActor("Iris", 3, event.ActorTypeHuman)   // reviewer

	// 1. Grace pre-registers hypothesis
	hypothesis, err := env.grammar.Emit(env.ctx, grace.ID(),
		"hypothesis: gamified learning improves retention by >15% vs traditional methods",
		env.convID, []types.EventID{env.boot.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit hypothesis: %v", err)
	}

	// 2. Grace documents methodology (before data)
	methodology, err := env.grammar.Extend(env.ctx, grace.ID(),
		"methodology: RCT, n=60, 3 groups, 4-week intervention, mixed ANOVA, outlier criterion: >3 SD",
		hypothesis.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Extend methodology: %v", err)
	}

	// 3. Data collected over 4 weeks (hash of data, not content)
	data1, err := env.grammar.Extend(env.ctx, grace.ID(),
		"data collected: week 1, n=58, 2 dropouts, data hash: sha256:abc123",
		methodology.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Extend data1: %v", err)
	}

	data4, err := env.grammar.Extend(env.ctx, grace.ID(),
		"data collected: week 4 (final), n=55, data hash: sha256:def456",
		data1.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Extend data4: %v", err)
	}

	// 4. FIRST analysis — fails (not significant)
	analysis1, err := env.grammar.Derive(env.ctx, grace.ID(),
		"analysis attempt 1: mixed ANOVA, F(2,55)=1.23, p=0.301, NOT SIGNIFICANT",
		data4.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive analysis1: %v", err)
	}

	// 5. Event STAYS on graph — append-only, cannot delete

	// 6. Grace adjusts analysis (linked to first attempt)
	analysis2, err := env.grammar.Derive(env.ctx, grace.ID(),
		"analysis attempt 2: removed 3 outliers per pre-registered criterion (>3 SD), F(2,52)=4.87, p=0.011, SIGNIFICANT",
		analysis1.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive analysis2: %v", err)
	}

	// 7. Grace writes manuscript
	manuscript, err := env.grammar.Derive(env.ctx, grace.ID(),
		"manuscript: Gamified Learning Effects on Knowledge Retention",
		analysis2.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive manuscript: %v", err)
	}

	// 8. Henry reviews — requests revision
	henryReview, err := env.grammar.Respond(env.ctx, henry.ID(),
		"review: need to see full analysis chain including failed attempts, revise and resubmit",
		manuscript.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Respond henry review: %v", err)
	}

	// 9. Iris reviews — endorses
	irisReview, err := env.grammar.Respond(env.ctx, iris.ID(),
		"review: methodology sound, pre-registration verified, accept",
		manuscript.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Respond iris review: %v", err)
	}

	irisEndorse, err := env.grammar.Endorse(env.ctx, iris.ID(),
		manuscript.ID(), grace.ID(), types.MustWeight(0.7),
		types.Some(types.MustDomainScope("research")),
		env.convID, signer)
	if err != nil {
		t.Fatalf("Endorse iris: %v", err)
	}

	// 10. Grace revises (both reviews in causes via Merge)
	revision, err := env.grammar.Merge(env.ctx, grace.ID(),
		"revision: added full analysis chain, addressed Henry's concerns",
		[]types.EventID{henryReview.ID(), irisReview.ID()},
		env.convID, signer)
	if err != nil {
		t.Fatalf("Merge revision: %v", err)
	}

	// 11. Published
	published, err := env.grammar.Derive(env.ctx, grace.ID(),
		"published: Gamified Learning Effects on Knowledge Retention, DOI:10.1234/example",
		revision.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive published: %v", err)
	}

	// --- Assertions ---

	// Pre-registration provable: hypothesis hash-chained BEFORE data collection
	// methodology → hypothesis (methodology comes after hypothesis)
	methAncestors := env.ancestors(methodology.ID(), 5)
	if !containsEvent(methAncestors, hypothesis.ID()) {
		t.Error("methodology should trace to hypothesis")
	}

	// Failed analysis visible: analysis2 traces to analysis1
	analysis2Ancestors := env.ancestors(analysis2.ID(), 5)
	if !containsEvent(analysis2Ancestors, analysis1.ID()) {
		t.Error("second analysis should trace to first (failed) analysis")
	}

	// Both analyses visible from manuscript
	manuscriptAncestors := env.ancestors(manuscript.ID(), 10)
	if !containsEvent(manuscriptAncestors, analysis2.ID()) {
		t.Error("manuscript should trace to successful analysis")
	}
	if !containsEvent(manuscriptAncestors, analysis1.ID()) {
		t.Error("manuscript should trace to failed analysis (through analysis2)")
	}

	// Revision includes both reviews
	revisionAncestors := env.ancestors(revision.ID(), 5)
	if !containsEvent(revisionAncestors, henryReview.ID()) {
		t.Error("revision should include Henry's review")
	}
	if !containsEvent(revisionAncestors, irisReview.ID()) {
		t.Error("revision should include Iris's review")
	}

	// Published traces all the way to hypothesis
	publishedAncestors := env.ancestors(published.ID(), 20)
	if !containsEvent(publishedAncestors, hypothesis.ID()) {
		t.Error("publication should trace back to pre-registered hypothesis")
	}

	_ = irisEndorse

	// Chain integrity
	env.verifyChain()

	// bootstrap(1) + hypothesis(1) + methodology(1) + data1(1) + data4(1) +
	// analysis1(1) + analysis2(1) + manuscript(1) + henryReview(1) + irisReview(1) +
	// irisEndorse(1) + revision(1) + published(1) = 13
	if count := env.eventCount(); count != 13 {
		t.Errorf("event count = %d, want 13", count)
	}
}
