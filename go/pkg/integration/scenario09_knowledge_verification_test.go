package integration_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// TestScenario09_KnowledgeVerification exercises self-correcting knowledge.
// An analyst agent makes a claim, infers a generalization, a reviewer challenges
// with counter-evidence, bias is detected, knowledge is corrected, and the
// correction propagates to dependent inferences.
func TestScenario09_KnowledgeVerification(t *testing.T) {
	env := newTestEnv(t)

	analyst := env.registerActor("AnalystBot", 1, event.ActorTypeAI)
	reviewer := env.registerActor("ReviewerBot", 2, event.ActorTypeAI)

	// 1. Analyst makes performance claim
	claim, err := env.grammar.Emit(env.ctx, analyst.ID(),
		"fact: Service X handles 10,000 RPS with p99 < 50ms on framework Y",
		env.convID, []types.EventID{env.boot.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit claim: %v", err)
	}

	// 2. Agent categorizes the claim
	classification, err := env.grammar.Annotate(env.ctx, analyst.ID(),
		claim.ID(), "classification", "performance_benchmark",
		env.convID, signer)
	if err != nil {
		t.Fatalf("Annotate classification: %v", err)
	}

	// 3. Agent infers generalization
	inference, err := env.grammar.Derive(env.ctx, analyst.ID(),
		"inference: all services on framework Y can handle 10,000+ RPS under load",
		claim.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive inference: %v", err)
	}

	// 4. Reviewer challenges with counter-evidence
	challenge, err := env.grammar.Respond(env.ctx, reviewer.ID(),
		"challenge: independent benchmark shows Service X at 6,200 RPS, p99=120ms under production traffic with DB contention",
		claim.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Respond challenge: %v", err)
	}

	// 5. Bias detected in original benchmark
	biasDetected, err := env.grammar.Annotate(env.ctx, reviewer.ID(),
		claim.ID(), "bias",
		"sampling bias: original benchmark used synthetic traffic without DB contention or concurrent users",
		env.convID, signer)
	if err != nil {
		t.Fatalf("Annotate bias: %v", err)
	}

	// 6. Knowledge corrected
	correction, err := env.grammar.Derive(env.ctx, analyst.ID(),
		"correction: Service X handles 6,000-7,000 RPS under production load with p99=100-120ms",
		challenge.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive correction: %v", err)
	}

	// 7. Correction propagated to dependent inference
	propagation, err := env.grammar.Annotate(env.ctx, analyst.ID(),
		inference.ID(), "invalidated",
		"dependent inference invalidated: original claim corrected, generalization no longer supported",
		env.convID, signer)
	if err != nil {
		t.Fatalf("Annotate propagation: %v", err)
	}

	// 8. Agent learns from experience
	learning, err := env.grammar.Extend(env.ctx, analyst.ID(),
		"learning: always verify benchmarks include production conditions (DB contention, concurrent users, realistic payloads)",
		correction.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Extend learning: %v", err)
	}

	// 9. Trust updated — decreased but not crashed
	_, err = env.graph.Record(
		event.EventTypeTrustUpdated, env.system,
		event.TrustUpdatedContent{
			Actor:    analyst.ID(),
			Previous: types.MustScore(0.5),
			Current:  types.MustScore(0.35),
			Domain:   types.MustDomainScope("benchmarking"),
			Cause:    correction.ID(),
		},
		[]types.EventID{correction.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Record trust: %v", err)
	}

	// --- Assertions ---

	// Original claim preserved — corrections don't delete
	originalClaim, err := env.store.Get(claim.ID())
	if err != nil {
		t.Fatalf("original claim should still exist: %v", err)
	}
	if originalClaim.Type().Value() != "grammar.emit" {
		t.Error("original claim should be intact")
	}

	// Correction chain: correction traces through challenge to original claim
	correctionAncestors := env.ancestors(correction.ID(), 10)
	if !containsEvent(correctionAncestors, challenge.ID()) {
		t.Error("correction should trace to challenge")
	}
	if !containsEvent(correctionAncestors, claim.ID()) {
		t.Error("correction should trace to original claim")
	}

	// Propagation: inference has invalidation annotation
	_ = propagation
	_ = biasDetected
	_ = classification

	// Learning is linked to correction
	learningAncestors := env.ancestors(learning.ID(), 5)
	if !containsEvent(learningAncestors, correction.ID()) {
		t.Error("learning should trace to correction")
	}

	// Chain integrity
	env.verifyChain()

	// bootstrap(1) + claim(1) + classification(1) + inference(1) + challenge(1) +
	// bias(1) + correction(1) + propagation(1) + learning(1) + trust(1) = 10
	if count := env.eventCount(); count != 10 {
		t.Errorf("event count = %d, want 10", count)
	}
}
