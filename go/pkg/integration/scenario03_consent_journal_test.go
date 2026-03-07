package integration_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// TestScenario03_ConsentJournal exercises consent-based shared journaling.
// Alice and Bob share a journal with explicit consent, trust accumulates,
// betrayal is detected, trust crashes, channel is severed, and eventually forgiven.
func TestScenario03_ConsentJournal(t *testing.T) {
	env := newTestEnv(t)

	alice := env.registerActor("Alice", 1, event.ActorTypeHuman)
	bob := env.registerActor("Bob", 2, event.ActorTypeHuman)

	// 1. Alice invites Bob (Endorse + Subscribe)
	endorseEv, subscribeEv, err := env.grammar.Invite(env.ctx, alice.ID(), bob.ID(),
		types.MustWeight(0.5), types.Some(types.MustDomainScope("journaling")),
		env.boot.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Invite: %v", err)
	}

	// 2. Bob subscribes back (reciprocal)
	bobSub, err := env.grammar.Subscribe(env.ctx, bob.ID(), alice.ID(),
		types.Some(types.MustDomainScope("journaling")),
		subscribeEv.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Subscribe Bob: %v", err)
	}

	// 3. Both open private channel
	channel, err := env.grammar.Channel(env.ctx, alice.ID(), bob.ID(),
		types.Some(types.MustDomainScope("journaling")),
		bobSub.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Channel: %v", err)
	}

	// 4. Alice writes journal entry
	entry, err := env.grammar.Emit(env.ctx, alice.ID(),
		"journal: feeling uncertain about career change, weighing options",
		env.convID, []types.EventID{channel.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit journal: %v", err)
	}

	// 5. Alice requests consent to share with Bob
	consentReq, err := env.graph.Record(
		event.EventTypeAuthorityRequested, alice.ID(),
		event.AuthorityRequestContent{
			Actor:  alice.ID(),
			Action: "share_journal_entry",
			Level:  event.AuthorityLevelRequired,
		},
		[]types.EventID{entry.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Record consent request: %v", err)
	}

	// 6. Bob consents
	consentApproval, err := env.graph.Record(
		event.EventTypeAuthorityResolved, bob.ID(),
		event.AuthorityResolvedContent{
			RequestID: consentReq.ID(),
			Approved:  true,
			Resolver:  bob.ID(),
			Reason:    types.None[string](),
		},
		[]types.EventID{consentReq.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Record consent resolved: %v", err)
	}

	// 7. Bob responds with own journal entry (causally linked to Alice's + consent)
	bobEntry, err := env.grammar.Respond(env.ctx, bob.ID(),
		"journal: I went through something similar last year, here's what helped...",
		consentApproval.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Respond bob entry: %v", err)
	}

	// 8. Trust accumulates from reciprocal sharing
	_, err = env.graph.Record(
		event.EventTypeTrustUpdated, env.system,
		event.TrustUpdatedContent{
			Actor:    bob.ID(),
			Previous: types.MustScore(0.1),
			Current:  types.MustScore(0.52),
			Domain:   types.MustDomainScope("journaling"),
			Cause:    bobEntry.ID(),
		},
		[]types.EventID{bobEntry.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Record trust up: %v", err)
	}

	// 9. Bob betrays — shares Alice's private entry externally
	betrayal, err := env.grammar.Emit(env.ctx, bob.ID(),
		"shared externally: Alice's private journal entry about career uncertainty",
		env.convID, []types.EventID{entry.ID()}, signer)
	if err != nil {
		t.Fatalf("Emit betrayal: %v", err)
	}

	// 10. Violation detected
	violation, err := env.graph.Record(
		event.EventTypeViolationDetected, env.system,
		event.ViolationDetectedContent{
			Expectation: entry.ID(),
			Actor:       bob.ID(),
			Severity:    event.SeverityLevelCritical,
			Description: "shared private channel content externally",
			Evidence:    types.MustNonEmpty([]types.EventID{betrayal.ID()}),
		},
		[]types.EventID{betrayal.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Record violation: %v", err)
	}

	// 11. Trust drops sharply
	_, err = env.graph.Record(
		event.EventTypeTrustUpdated, env.system,
		event.TrustUpdatedContent{
			Actor:    bob.ID(),
			Previous: types.MustScore(0.52),
			Current:  types.MustScore(0.1),
			Domain:   types.MustDomainScope("journaling"),
			Cause:    violation.ID(),
		},
		[]types.EventID{violation.ID()}, env.convID, signer)
	if err != nil {
		t.Fatalf("Record trust crash: %v", err)
	}

	// 12. Alice severs the channel
	channelEdgeID, _ := types.NewEdgeID(channel.ID().Value())
	severEv, err := env.grammar.Sever(env.ctx, alice.ID(),
		channelEdgeID, violation.ID(), env.convID, signer)
	if err != nil {
		t.Fatalf("Sever: %v", err)
	}

	// 13. Much later — Alice forgives (re-subscribes)
	forgiveEv, err := env.grammar.Forgive(env.ctx, alice.ID(),
		severEv.ID(), bob.ID(),
		types.Some(types.MustDomainScope("journaling")),
		env.convID, signer)
	if err != nil {
		t.Fatalf("Forgive: %v", err)
	}

	// --- Assertions ---

	// Forgive ≠ Forget: causal chain from betrayal through sever to forgiveness
	forgiveAncestors := env.ancestors(forgiveEv.ID(), 10)
	if !containsEvent(forgiveAncestors, severEv.ID()) {
		t.Error("forgiveness should have sever in ancestors")
	}

	// Sever traces to violation
	severAncestors := env.ancestors(severEv.ID(), 10)
	if !containsEvent(severAncestors, violation.ID()) {
		t.Error("sever should have violation in ancestors")
	}

	// Consent chain: Bob's entry traces through consent to Alice's entry
	bobEntryAncestors := env.ancestors(bobEntry.ID(), 10)
	if !containsEvent(bobEntryAncestors, consentApproval.ID()) {
		t.Error("Bob's entry should trace through consent approval")
	}

	// Invitation records exist
	_ = endorseEv

	// Chain integrity
	env.verifyChain()

	// Count: bootstrap(1) + invite(endorse+subscribe=2) + bobSub(1) + channel(1) +
	// entry(1) + consentReq(1) + consentApproval(1) + bobEntry(1) + trustUp(1) +
	// betrayal(1) + violation(1) + trustCrash(1) + sever(1) + forgive(1) = 15
	if count := env.eventCount(); count != 15 {
		t.Errorf("event count = %d, want 15", count)
	}
}
