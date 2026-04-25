package integration_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// TestScenario02_FreelancerReputation exercises portable reputation across platforms.
// Carol posts a job, Bob proposes and delivers, Carol endorses Bob,
// Dave queries Bob's reputation and hires based on verifiable history.
func TestScenario02_FreelancerReputation(t *testing.T) {
	env := newTestEnv(t)

	carol := env.registerActor("Carol", 1, event.ActorTypeHuman)
	bob := env.registerActor("Bob", 2, event.ActorTypeHuman)
	dave := env.registerActor("Dave", 3, event.ActorTypeHuman)

	// 1. Carol posts a job listing
	listing, err := env.grammar.Emit(env.ctx, carol.ID(),
		"job listing: build REST API for inventory management, budget $3000",
		env.convID, []types.EventID{env.boot.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit listing: %v", err)
	}

	// 2. Bob proposes work
	proposal, err := env.grammar.Respond(env.ctx, bob.ID(),
		"proposal: can deliver in 2 weeks, $2800, Go + PostgreSQL",
		listing.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Respond proposal: %v", err)
	}

	// 3. Carol and Bob open a channel
	channel, err := env.grammar.Channel(env.ctx, carol.ID(), bob.ID(),
		types.Some(types.MustDomainScope("software_development")),
		proposal.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Channel: %v", err)
	}

	// 4. Both consent to bilateral contract
	contract, err := env.grammar.Consent(env.ctx, carol.ID(), bob.ID(),
		"REST API for inventory management, $2800, 2 week deadline",
		types.MustDomainScope("software_development"),
		channel.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Consent contract: %v", err)
	}

	// 5. Bob delivers work
	delivery, err := env.grammar.Derive(env.ctx, bob.ID(),
		"work delivered: REST API complete, 47 endpoints, 92% test coverage",
		contract.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive delivery: %v", err)
	}

	// 6. Carol acknowledges receipt
	ack, err := env.grammar.Acknowledge(env.ctx, carol.ID(),
		delivery.ID(), bob.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Acknowledge: %v", err)
	}

	// 7. Carol endorses Bob's work (reputation-staked)
	endorsement, err := env.grammar.Endorse(env.ctx, carol.ID(),
		delivery.ID(), bob.ID(), types.MustWeight(0.8),
		types.Some(types.MustDomainScope("software_development")),
		env.convID, signer)
	if err != nil {
		t.Fatalf("Endorse: %v", err)
	}

	// 8. Trust updated for Bob
	_, err = env.graph.Record(
		event.EventTypeTrustUpdated, env.system,
		event.TrustUpdatedContent{
			Actor:    bob.ID(),
			Previous: types.MustScore(0.1),
			Current:  types.MustScore(0.4),
			Domain:   types.MustDomainScope("software_development"),
			Cause:    endorsement.ID(),
		},
		[]types.EventID{endorsement.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Record trust: %v", err)
	}

	// 9. Dave queries Bob's reputation (read-only traversal)
	// Dave can see the full chain: endorsement → delivery → contract → proposal → listing
	endorseAncestors := env.ancestors(endorsement.ID(), 10)
	if !containsEvent(endorseAncestors, delivery.ID()) {
		t.Error("endorsement should trace to delivery")
	}
	if !containsEvent(endorseAncestors, contract.ID()) {
		t.Error("endorsement should trace to contract")
	}

	// 10. Dave hires Bob based on verifiable history
	daveListing, err := env.grammar.Emit(env.ctx, dave.ID(),
		"job listing: mobile app backend",
		env.convID, []types.EventID{env.boot.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit dave listing: %v", err)
	}

	daveContract, err := env.grammar.Consent(env.ctx, dave.ID(), bob.ID(),
		"mobile app backend, $4000",
		types.MustDomainScope("software_development"),
		daveListing.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Consent dave contract: %v", err)
	}

	// --- Assertions ---

	// Portable reputation: endorsement is queryable without Bob's involvement
	_ = ack
	_ = daveContract

	// Endorsement content has weight
	ec := endorsement.Content().(event.EdgeCreatedContent)
	if ec.Weight.Value() != 0.8 {
		t.Errorf("endorsement weight = %v, want 0.8", ec.Weight.Value())
	}

	// Endorsement is domain-scoped
	if !ec.Scope.IsSome() {
		t.Fatal("endorsement should have domain scope")
	}
	scope := ec.Scope.Unwrap()
	if scope.Value() != "software_development" {
		t.Errorf("endorsement scope = %v, want software_development", scope.Value())
	}

	// Chain integrity
	env.verifyChain()

	// bootstrap(1) + listing(1) + proposal(1) + channel(1) + contract(1) +
	// delivery(1) + ack(1) + endorsement(1) + trust(1) + daveListing(1) + daveContract(1) = 11
	if count := env.eventCount(); count != 11 {
		t.Errorf("event count = %d, want 11", count)
	}
}
