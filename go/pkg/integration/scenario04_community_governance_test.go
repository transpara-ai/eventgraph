package integration_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// TestScenario04_CommunityGovernance exercises collective decision making.
// Alice proposes a budget, discussion and amendment follow, members vote,
// outcome is tallied, and the full chain from enactment to proposal is traversable.
func TestScenario04_CommunityGovernance(t *testing.T) {
	env := newTestEnv(t)

	alice := env.registerActor("Alice", 1, event.ActorTypeHuman)
	bob := env.registerActor("Bob", 2, event.ActorTypeHuman)
	carol := env.registerActor("Carol", 3, event.ActorTypeHuman)
	dave := env.registerActor("Dave", 4, event.ActorTypeHuman)  // elder, high trust
	tallyBot := env.registerActor("TallyBot", 5, event.ActorTypeAI)

	// 1. Alice proposes budget for community garden
	proposal, err := env.grammar.Emit(env.ctx, alice.ID(),
		"proposal: allocate $2000 for community garden supplies and maintenance",
		env.convID, []types.EventID{env.boot.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit proposal: %v", err)
	}

	// 2. Bob raises concern
	concern, err := env.grammar.Respond(env.ctx, bob.ID(),
		"concern: $2000 is steep, could we do it for $1500 and use volunteers?",
		proposal.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Respond concern: %v", err)
	}

	// 3. Carol supports Alice
	support, err := env.grammar.Respond(env.ctx, carol.ID(),
		"support: the garden benefits everyone, $2000 is reasonable for quality materials",
		proposal.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Respond support: %v", err)
	}

	// 4. Bob proposes amendment
	amendment, err := env.grammar.Annotate(env.ctx, bob.ID(),
		proposal.ID(), "amendment",
		"reduce budget to $1500, recruit volunteer labour for installation",
		env.convID, signer)
	if err != nil {
		t.Fatalf("Annotate amendment: %v", err)
	}

	// 5. Dave (elder) endorses amendment
	_, err = env.grammar.Endorse(env.ctx, dave.ID(),
		amendment.ID(), bob.ID(), types.MustWeight(0.9),
		types.Some(types.MustDomainScope("governance")),
		env.convID, signer)
	if err != nil {
		t.Fatalf("Endorse amendment: %v", err)
	}

	// 6. Vote opens (system event after discussion period)
	voteOpen, err := env.grammar.Derive(env.ctx, tallyBot.ID(),
		"vote open: original ($2000) vs amended ($1500 + volunteers)",
		proposal.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive vote open: %v", err)
	}

	// 7-8. Members vote using Consent
	aliceVote, err := env.grammar.Consent(env.ctx, alice.ID(), tallyBot.ID(),
		"vote: original ($2000)",
		types.MustDomainScope("governance"),
		voteOpen.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Alice vote: %v", err)
	}

	bobVote, err := env.grammar.Consent(env.ctx, bob.ID(), tallyBot.ID(),
		"vote: amended ($1500)",
		types.MustDomainScope("governance"),
		voteOpen.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Bob vote: %v", err)
	}

	carolVote, err := env.grammar.Consent(env.ctx, carol.ID(), tallyBot.ID(),
		"vote: amended ($1500)",
		types.MustDomainScope("governance"),
		voteOpen.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Carol vote: %v", err)
	}

	daveVote, err := env.grammar.Consent(env.ctx, dave.ID(), tallyBot.ID(),
		"vote: amended ($1500)",
		types.MustDomainScope("governance"),
		voteOpen.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Dave vote: %v", err)
	}

	// 9. Bot tallies outcome (all votes in causes)
	outcome, err := env.grammar.Merge(env.ctx, tallyBot.ID(),
		"outcome: amended budget ($1500) passes 3-1",
		[]types.EventID{aliceVote.ID(), bobVote.ID(), carolVote.ID(), daveVote.ID()},
		env.convID, signer)
	if err != nil {
		t.Fatalf("Merge outcome: %v", err)
	}

	// 10. Budget enacted
	enacted, err := env.grammar.Derive(env.ctx, tallyBot.ID(),
		"enacted: community garden budget $1500 with volunteer labour",
		outcome.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive enacted: %v", err)
	}

	// --- Assertions ---

	// Full transparency: enacted → outcome → all votes → discussions → proposal
	enactedAncestors := env.ancestors(enacted.ID(), 10)
	if !containsEvent(enactedAncestors, outcome.ID()) {
		t.Error("enacted should trace to outcome")
	}

	outcomeAncestors := env.ancestors(outcome.ID(), 10)
	if !containsEvent(outcomeAncestors, aliceVote.ID()) {
		t.Error("outcome should include Alice's vote")
	}
	if !containsEvent(outcomeAncestors, bobVote.ID()) {
		t.Error("outcome should include Bob's vote")
	}
	if !containsEvent(outcomeAncestors, carolVote.ID()) {
		t.Error("outcome should include Carol's vote")
	}
	if !containsEvent(outcomeAncestors, daveVote.ID()) {
		t.Error("outcome should include Dave's vote")
	}

	// Amendment addresses concern (causal link through proposal)
	amendmentAncestors := env.ancestors(amendment.ID(), 10)
	if !containsEvent(amendmentAncestors, proposal.ID()) {
		t.Error("amendment should trace to proposal")
	}

	// Discussion events exist
	_ = concern
	_ = support

	// Chain integrity
	env.verifyChain()

	// bootstrap(1) + proposal(1) + concern(1) + support(1) + amendment(1) +
	// endorse(1) + voteOpen(1) + 4votes(4) + outcome(1) + enacted(1) = 13
	if count := env.eventCount(); count != 13 {
		t.Errorf("event count = %d, want 13", count)
	}
}
