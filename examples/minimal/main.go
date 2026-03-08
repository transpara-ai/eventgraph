// Minimal EventGraph example — bootstrap a graph, record two events, verify the chain.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lovyou-ai/eventgraph/go/pkg/grammar"
	"github.com/lovyou-ai/eventgraph/go/pkg/graph"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

type noopSigner struct{}

func (noopSigner) Sign(data []byte) (types.Signature, error) {
	return types.MustSignature(make([]byte, 64)), nil
}

func main() {
	// 1. Create in-memory store and graph.
	s := store.NewInMemoryStore()
	g := graph.New(s, nil)
	if err := g.Start(); err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	// 2. Bootstrap the hash chain.
	systemActor := types.MustActorID("actor_system")
	boot, err := g.Bootstrap(systemActor, noopSigner{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bootstrap: %s (hash: %.16s…)\n", boot.ID().Value(), boot.Hash().Value())

	// 3. Record events using social grammar.
	gr := grammar.New(g)
	ctx := context.Background()
	convID := types.MustConversationID("conv_example")
	actor := types.MustActorID("actor_alice")

	ev1, err := gr.Emit(ctx, actor, "Hello, EventGraph!", convID, []types.EventID{boot.ID()}, noopSigner{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Event 1:   %s (type: %s)\n", ev1.ID().Value(), ev1.Type().Value())

	ev2, err := gr.Derive(ctx, actor, "A derived thought", ev1.ID(), convID, noopSigner{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Event 2:   %s (type: %s)\n", ev2.ID().Value(), ev2.Type().Value())

	// 4. Verify chain integrity.
	result, err := s.VerifyChain()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Chain:     %d events, valid=%v\n", result.Length, result.Valid)
}
