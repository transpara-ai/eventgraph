package integration_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// TestScenario05_SupplyChain exercises supply chain provenance.
// NOTE: The full spec requires EGIP (inter-system protocol) which is not yet built.
// This test exercises single-system provenance: farm → factory → retailer → consumer trace.
func TestScenario05_SupplyChain(t *testing.T) {
	env := newTestEnv(t)

	farm := env.registerActor("Sunrise Farm", 1, event.ActorTypeHuman)
	factory := env.registerActor("GreenProcess Factory", 2, event.ActorTypeHuman)
	qaAgent := env.registerActor("QA Inspector", 3, event.ActorTypeAI)
	retailer := env.registerActor("FreshMart", 4, event.ActorTypeHuman)

	// 1. Farm records harvest
	harvest, err := env.grammar.Emit(env.ctx, farm.ID(),
		"harvest: 500kg organic tomatoes, lot #TOM-2026-0308, field B3",
		env.convID, []types.EventID{env.boot.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit harvest: %v", err)
	}

	// 2. Factory receives produce (Derive from harvest)
	received, err := env.grammar.Derive(env.ctx, factory.ID(),
		"received: 500kg tomatoes from Sunrise Farm, lot #TOM-2026-0308",
		harvest.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive received: %v", err)
	}

	// 3. QA agent inspects produce
	inspection, err := env.grammar.Derive(env.ctx, qaAgent.ID(),
		"qa inspection: pesticide-free verified, freshness grade A, confidence 0.95",
		received.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive inspection: %v", err)
	}

	// 4. Factory manufactures product
	product, err := env.grammar.Derive(env.ctx, factory.ID(),
		"manufactured: 200 jars organic tomato sauce, batch #SAU-2026-0308",
		inspection.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive product: %v", err)
	}

	// 5. Factory endorses farm quality
	_, err = env.grammar.Endorse(env.ctx, factory.ID(),
		harvest.ID(), farm.ID(), types.MustWeight(0.85),
		types.Some(types.MustDomainScope("produce_quality")),
		env.convID, signer)
	if err != nil {
		t.Fatalf("Endorse farm: %v", err)
	}

	// 6. Retailer receives product
	listed, err := env.grammar.Derive(env.ctx, retailer.ID(),
		"product listed: organic tomato sauce, batch #SAU-2026-0308, provenance verified",
		product.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive listed: %v", err)
	}

	// --- Assertions ---

	// Consumer traces provenance: listed → product → inspection → received → harvest
	listedAncestors := env.ancestors(listed.ID(), 10)
	if !containsEvent(listedAncestors, product.ID()) {
		t.Error("listed should trace to product")
	}
	if !containsEvent(listedAncestors, inspection.ID()) {
		t.Error("listed should trace to inspection")
	}
	if !containsEvent(listedAncestors, received.ID()) {
		t.Error("listed should trace to received")
	}
	if !containsEvent(listedAncestors, harvest.ID()) {
		t.Error("listed should trace to harvest (farm origin)")
	}

	// QA is auditable — inspection is on the chain
	inspectionAncestors := env.ancestors(inspection.ID(), 5)
	if !containsEvent(inspectionAncestors, received.ID()) {
		t.Error("inspection should trace to received produce")
	}

	// Chain integrity
	env.verifyChain()

	// bootstrap(1) + harvest(1) + received(1) + inspection(1) + product(1) +
	// endorse(1) + listed(1) = 7
	if count := env.eventCount(); count != 7 {
		t.Errorf("event count = %d, want 7", count)
	}
}
