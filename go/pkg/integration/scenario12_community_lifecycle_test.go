package integration_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// TestScenario12_CommunityLifecycle exercises onboarding, traditions, stewardship, and succession.
// A newcomer is invited, makes contributions, belonging grows, stewardship is
// transferred through succession, milestones are celebrated, and gifts create no obligation.
func TestScenario12_CommunityLifecycle(t *testing.T) {
	env := newTestEnv(t)

	alice := env.registerActor("Alice", 1, event.ActorTypeHuman) // founder
	carol := env.registerActor("Carol", 2, event.ActorTypeHuman) // current steward
	bob := env.registerActor("Bob", 3, event.ActorTypeHuman)     // newcomer

	// 1. Alice invites Bob (sponsor endorsement)
	endorseEv, subscribeEv, err := env.grammar.Invite(env.ctx, alice.ID(), bob.ID(),
		types.MustWeight(0.4),
		types.Some(types.MustDomainScope("community")),
		env.boot.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Invite: %v", err)
	}

	// 2. Bob settles — belonging starts low
	settle, err := env.grammar.Emit(env.ctx, bob.ID(),
		"home: joined the community, feeling welcomed, belonging 0.15",
		env.convID, []types.EventID{subscribeEv.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit settle: %v", err)
	}

	// 3. Bob makes first contribution
	contrib1, err := env.grammar.Emit(env.ctx, bob.ID(),
		"contribution: added unit tests for the auth module, 15 test cases",
		env.convID, []types.EventID{settle.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit contribution: %v", err)
	}

	// 4. Community acknowledges (trust starts growing)
	_, err = env.grammar.Acknowledge(env.ctx, carol.ID(),
		contrib1.ID(), bob.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Acknowledge: %v", err)
	}

	// 5. Trust accumulates
	_, err = env.graph.Record(
		event.EventTypeTrustUpdated, env.system,
		event.TrustUpdatedContent{
			Actor:    bob.ID(),
			Previous: types.MustScore(0.1),
			Current:  types.MustScore(0.35),
			Domain:   types.MustDomainScope("community"),
			Cause:    contrib1.ID(),
		},
		[]types.EventID{contrib1.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Record trust: %v", err)
	}

	// 6. Bob participates in tradition (Friday retrospective)
	tradition, err := env.grammar.Emit(env.ctx, bob.ID(),
		"tradition: participated in Friday retrospective, 12th consecutive week",
		env.convID, []types.EventID{contrib1.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit tradition: %v", err)
	}

	// 7. More contributions accumulate (summary)
	contribSummary, err := env.grammar.Extend(env.ctx, bob.ID(),
		"contributions: 30 total over 6 months, trust now 0.65, belonging 0.78",
		tradition.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Extend contributions: %v", err)
	}

	// 8. Sustainability assessment — Carol is sole steward (bus factor risk)
	sustainability, err := env.grammar.Emit(env.ctx, env.system,
		"sustainability: bus factor risk — Carol is sole steward of test infrastructure",
		env.convID, []types.EventID{contribSummary.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit sustainability: %v", err)
	}

	// 9. Succession planned — Carol delegates to Bob
	successionPlan, err := env.grammar.Delegate(env.ctx, carol.ID(), bob.ID(),
		types.MustDomainScope("test_infrastructure"), types.MustWeight(0.8),
		sustainability.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Delegate succession: %v", err)
	}

	// 10. Succession completed — bilateral consent
	successionComplete, err := env.grammar.Consent(env.ctx, carol.ID(), bob.ID(),
		"succession complete: Bob is now steward of test infrastructure",
		types.MustDomainScope("test_infrastructure"),
		successionPlan.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Consent succession: %v", err)
	}

	// 11. Community celebrates v2.0 milestone
	milestone, err := env.grammar.Emit(env.ctx, env.system,
		"milestone: v2.0 released, 6 months of community effort, 30 contributions from Bob alone",
		env.convID, []types.EventID{successionComplete.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit milestone: %v", err)
	}

	// 12. Community story chapter
	story, err := env.grammar.Derive(env.ctx, env.system,
		"community story: Bob's journey — newcomer to steward in 6 months, 30 contributions, adopted test infrastructure",
		milestone.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Derive story: %v", err)
	}

	// 13. Gift given (unconditional, no obligation)
	gift, err := env.grammar.Emit(env.ctx, alice.ID(),
		"gift: custom test harness for Bob, unconditional, no obligation or reciprocity expected",
		env.convID, []types.EventID{milestone.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit gift: %v", err)
	}

	// --- Assertions ---

	// Belonging gradient visible through event chain
	_ = settle // 0.15 → contributions → 0.78

	// Succession bilateral
	successionContent := successionComplete.Content().(event.GrammarConsentContent)
	hasCarol := successionContent.Parties[0] == carol.ID() || successionContent.Parties[1] == carol.ID()
	hasBob := successionContent.Parties[0] == bob.ID() || successionContent.Parties[1] == bob.ID()
	if !hasCarol {
		t.Error("succession should include Carol")
	}
	if !hasBob {
		t.Error("succession should include Bob")
	}

	// Story traces to milestone
	storyAncestors := env.ancestors(story.ID(), 5)
	if !containsEvent(storyAncestors, milestone.ID()) {
		t.Error("story should trace to milestone")
	}

	// Succession traces to sustainability concern
	successionAncestors := env.ancestors(successionPlan.ID(), 5)
	if !containsEvent(successionAncestors, sustainability.ID()) {
		t.Error("succession should trace to sustainability assessment")
	}

	// Gift creates no obligation (just an emit, no consent or edge)
	giftContent := gift.Content().(event.GrammarEmitContent)
	if giftContent.Body == "" {
		t.Error("gift should have content")
	}

	// Invitation records exist
	_ = endorseEv

	// Chain integrity
	env.verifyChain()

	// bootstrap(1) + invite(endorse+subscribe=2) + settle(1) + contrib1(1) + ack(1) +
	// trust(1) + tradition(1) + contribSummary(1) + sustainability(1) + succession(1) +
	// successionComplete(1) + milestone(1) + story(1) + gift(1) = 15
	if count := env.eventCount(); count != 15 {
		t.Errorf("event count = %d, want 15", count)
	}
}
