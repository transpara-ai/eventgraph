// Social Grammar example — demonstrates the 15 social grammar operations on an event graph.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/grammar"
	"github.com/lovyou-ai/eventgraph/go/pkg/graph"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

type noopSigner struct{}

func (noopSigner) Sign(data []byte) (types.Signature, error) {
	return types.MustSignature(make([]byte, 64)), nil
}

var sign = noopSigner{}

func main() {
	s := store.NewInMemoryStore()
	g := graph.New(s, nil)
	if err := g.Start(); err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	systemActor := types.MustActorID("actor_system")
	boot, err := g.Bootstrap(systemActor, sign)
	if err != nil {
		log.Fatal(err)
	}

	gr := grammar.New(g)
	ctx := context.Background()
	conv := types.MustConversationID("conv_social_demo")

	alice := types.MustActorID("actor_alice")
	bob := types.MustActorID("actor_bob")

	// 1. Emit — Alice starts a conversation.
	emit, err := gr.Emit(ctx, alice, "I have an idea for a new project", conv, []types.EventID{boot.ID()}, sign)
	must(err, "Emit")
	show("Emit", emit)

	// 2. Respond — Bob replies.
	respond, err := gr.Respond(ctx, bob, "That sounds great, tell me more", emit.ID(), conv, sign)
	must(err, "Respond")
	show("Respond", respond)

	// 3. Derive — Alice builds on her idea.
	derive, err := gr.Derive(ctx, alice, "Here's the detailed spec", emit.ID(), conv, sign)
	must(err, "Derive")
	show("Derive", derive)

	// 4. Extend — Bob adds to Alice's spec.
	extend, err := gr.Extend(ctx, bob, "Adding performance requirements", derive.ID(), conv, sign)
	must(err, "Extend")
	show("Extend", extend)

	// 5. Retract — Alice retracts her initial emit (only author can retract).
	retract, err := gr.Retract(ctx, alice, emit.ID(), "Superseded by the spec", conv, sign)
	must(err, "Retract")
	show("Retract", retract)

	// 6. Annotate — Alice adds metadata to the spec.
	annotate, err := gr.Annotate(ctx, alice, derive.ID(), "priority", "high", conv, sign)
	must(err, "Annotate")
	show("Annotate", annotate)

	// 7. Acknowledge — Bob acknowledges Alice's spec.
	ack, err := gr.Acknowledge(ctx, bob, derive.ID(), alice, conv, sign)
	must(err, "Acknowledge")
	show("Acknowledge", ack)

	// 8. Propagate — Alice propagates the spec to Bob.
	propagate, err := gr.Propagate(ctx, alice, derive.ID(), bob, conv, sign)
	must(err, "Propagate")
	show("Propagate", propagate)

	// 9. Endorse — Bob endorses Alice's spec (reputation-staked).
	endorse, err := gr.Endorse(ctx, bob, derive.ID(), alice,
		types.MustWeight(0.9), types.Some(types.MustDomainScope("innovation")), conv, sign)
	must(err, "Endorse")
	show("Endorse", endorse)

	// 10. Subscribe — Bob subscribes to Alice's updates.
	subscribe, err := gr.Subscribe(ctx, bob, alice,
		types.Some(types.MustDomainScope("projects")), boot.ID(), conv, sign)
	must(err, "Subscribe")
	show("Subscribe", subscribe)

	// 11. Channel — Alice creates a private channel with Bob.
	channel, err := gr.Channel(ctx, alice, bob,
		types.Some(types.MustDomainScope("project_alpha")), emit.ID(), conv, sign)
	must(err, "Channel")
	show("Channel", channel)

	// 12. Delegate — Alice delegates a task to Bob.
	delegate, err := gr.Delegate(ctx, alice, bob,
		types.MustDomainScope("engineering"), types.MustWeight(0.8), derive.ID(), conv, sign)
	must(err, "Delegate")
	show("Delegate", delegate)

	// 13. Consent — Alice and Bob establish mutual consent.
	consent, err := gr.Consent(ctx, alice, bob,
		"Both parties agree to collaborate on project alpha",
		types.MustDomainScope("collaboration"), derive.ID(), conv, sign)
	must(err, "Consent")
	show("Consent", consent)

	// 14. Sever — Alice severs the delegation (only parties can sever).
	delegateEdgeID := types.MustEdgeID(delegate.ID().Value())
	sever, err := gr.Sever(ctx, alice, delegateEdgeID, consent.ID(), conv, sign)
	must(err, "Sever")
	show("Sever", sever)

	// 15. Merge — Alice merges her spec with Bob's extension.
	merge, err := gr.Merge(ctx, alice, "Final spec with all requirements",
		[]types.EventID{derive.ID(), extend.ID()}, conv, sign)
	must(err, "Merge")
	show("Merge", merge)

	// Verify integrity.
	result, _ := s.VerifyChain()
	count, _ := s.Count()
	fmt.Printf("\nChain: %d events, valid=%v\n", count, result.Valid)
}

func show(op string, ev event.Event) {
	fmt.Printf("%-12s → %s (%s…)\n", op, ev.Type().Value(), ev.ID().Value()[:13])
}

func must(err error, op string) {
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}
}
