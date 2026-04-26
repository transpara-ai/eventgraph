package primitive_test

import (
	"fmt"
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/primitive"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

func TestHarnessProcess(t *testing.T) {
	primID := types.MustPrimitiveID("test_prim")
	p := &testPrimitive{
		id:        primID,
		layer:     types.MustLayer(0),
		cadence:   types.MustCadence(1),
		lifecycle: types.LifecycleActive,
		subscriptions: []types.SubscriptionPattern{
			types.MustSubscriptionPattern("*"),
		},
		processFunc: func(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
			return []primitive.Mutation{
				primitive.UpdateState{PrimitiveID: primID, Key: "count", Value: 1},
			}, nil
		},
	}

	h := primitive.NewHarness()
	mutations, err := h.Process(p, nil)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) != 1 {
		t.Errorf("expected 1 mutation, got %d", len(mutations))
	}

	changes := h.StateChanges()
	if changes["count"] != 1 {
		t.Errorf("expected count=1, got %v", changes["count"])
	}
}

func TestHarnessEmittedEvents(t *testing.T) {
	primID := types.MustPrimitiveID("emitter")
	actorID := types.MustActorID("actor_system0000000000000000001")

	p := &testPrimitive{
		id:        primID,
		layer:     types.MustLayer(0),
		cadence:   types.MustCadence(1),
		lifecycle: types.LifecycleActive,
		subscriptions: []types.SubscriptionPattern{
			types.MustSubscriptionPattern("*"),
		},
		processFunc: func(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
			return []primitive.Mutation{
				primitive.AddEvent{
					Type:    event.EventTypeTrustUpdated,
					Source:  actorID,
					Content: event.TrustUpdatedContent{},
				},
				primitive.AddEvent{
					Type:    event.EventTypeEdgeCreated,
					Source:  actorID,
					Content: event.EdgeCreatedContent{},
				},
			}, nil
		},
	}

	h := primitive.NewHarness()
	h.Process(p, nil)

	emitted := h.EmittedEvents()
	if len(emitted) != 2 {
		t.Fatalf("expected 2 emitted events, got %d", len(emitted))
	}
	if emitted[0].Type != event.EventTypeTrustUpdated {
		t.Errorf("first event type = %v, want trust.updated", emitted[0].Type)
	}
	if emitted[1].Type != event.EventTypeEdgeCreated {
		t.Errorf("second event type = %v, want edge.created", emitted[1].Type)
	}
}

func TestHarnessEdgeMutations(t *testing.T) {
	primID := types.MustPrimitiveID("edge_creator")
	from := types.MustActorID("actor_system0000000000000000001")
	to := types.MustActorID("actor_system0000000000000000002")

	p := &testPrimitive{
		id:        primID,
		layer:     types.MustLayer(0),
		cadence:   types.MustCadence(1),
		lifecycle: types.LifecycleActive,
		processFunc: func(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
			return []primitive.Mutation{
				primitive.AddEdge{
					From:     from,
					To:       to,
					EdgeType: event.EdgeTypeTrust,
					Weight:   types.MustWeight(0.5),
					Scope:    types.None[types.DomainScope](),
				},
			}, nil
		},
	}

	h := primitive.NewHarness()
	h.Process(p, nil)

	edges := h.EdgeMutations()
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge mutation, got %d", len(edges))
	}
	if edges[0].From != from {
		t.Errorf("From = %v, want %v", edges[0].From, from)
	}
	if edges[0].To != to {
		t.Errorf("To = %v, want %v", edges[0].To, to)
	}
}

func TestHarnessActivationChanges(t *testing.T) {
	primID := types.MustPrimitiveID("activator")
	p := &testPrimitive{
		id:        primID,
		layer:     types.MustLayer(0),
		cadence:   types.MustCadence(1),
		lifecycle: types.LifecycleActive,
		processFunc: func(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
			return []primitive.Mutation{
				primitive.UpdateActivation{PrimitiveID: primID, Level: types.MustActivation(0.9)},
			}, nil
		},
	}

	h := primitive.NewHarness()
	h.Process(p, nil)

	activations := h.ActivationChanges()
	if len(activations) != 1 {
		t.Fatalf("expected 1 activation change, got %d", len(activations))
	}
	if activations[0].Level.Value() != 0.9 {
		t.Errorf("Level = %v, want 0.9", activations[0].Level.Value())
	}
}

func TestHarnessLifecycleChanges(t *testing.T) {
	primID := types.MustPrimitiveID("lifecycle_changer")
	p := &testPrimitive{
		id:        primID,
		layer:     types.MustLayer(0),
		cadence:   types.MustCadence(1),
		lifecycle: types.LifecycleActive,
		processFunc: func(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
			return []primitive.Mutation{
				primitive.UpdateLifecycle{PrimitiveID: primID, State: types.LifecycleDeactivating},
			}, nil
		},
	}

	h := primitive.NewHarness()
	h.Process(p, nil)

	lifecycle := h.LifecycleChanges()
	if len(lifecycle) != 1 {
		t.Fatalf("expected 1 lifecycle change, got %d", len(lifecycle))
	}
	if lifecycle[0].State != types.LifecycleDeactivating {
		t.Errorf("State = %v, want Deactivating", lifecycle[0].State)
	}
}

func TestHarnessMixedMutations(t *testing.T) {
	primID := types.MustPrimitiveID("mixed")
	actorID := types.MustActorID("actor_system0000000000000000001")

	p := &testPrimitive{
		id:        primID,
		layer:     types.MustLayer(0),
		cadence:   types.MustCadence(1),
		lifecycle: types.LifecycleActive,
		processFunc: func(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
			return []primitive.Mutation{
				primitive.AddEvent{Type: event.EventTypeTrustUpdated, Source: actorID, Content: event.TrustUpdatedContent{}},
				primitive.UpdateState{PrimitiveID: primID, Key: "k1", Value: "v1"},
				primitive.UpdateState{PrimitiveID: primID, Key: "k2", Value: 42},
				primitive.UpdateActivation{PrimitiveID: primID, Level: types.MustActivation(0.5)},
				primitive.AddEdge{From: actorID, To: actorID, EdgeType: event.EdgeTypeTrust, Weight: types.MustWeight(0.1), Scope: types.None[types.DomainScope]()},
			}, nil
		},
	}

	h := primitive.NewHarness()
	h.Process(p, nil)

	if len(h.Mutations()) != 5 {
		t.Errorf("total mutations = %d, want 5", len(h.Mutations()))
	}
	if len(h.EmittedEvents()) != 1 {
		t.Errorf("emitted events = %d, want 1", len(h.EmittedEvents()))
	}
	if len(h.StateChanges()) != 2 {
		t.Errorf("state changes = %d, want 2", len(h.StateChanges()))
	}
	if len(h.ActivationChanges()) != 1 {
		t.Errorf("activation changes = %d, want 1", len(h.ActivationChanges()))
	}
	if len(h.EdgeMutations()) != 1 {
		t.Errorf("edge mutations = %d, want 1", len(h.EdgeMutations()))
	}
}

func TestHarnessNoMutations(t *testing.T) {
	p := newTestPrimitive("quiet", 0)

	h := primitive.NewHarness()
	mutations, err := h.Process(p, nil)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) != 0 {
		t.Errorf("expected 0 mutations, got %d", len(mutations))
	}
	if len(h.EmittedEvents()) != 0 {
		t.Errorf("expected 0 emitted events")
	}
	if len(h.StateChanges()) != 0 {
		t.Errorf("expected 0 state changes")
	}
}

func TestHarnessProcessError(t *testing.T) {
	p := &testPrimitive{
		id:        types.MustPrimitiveID("error_prim"),
		layer:     types.MustLayer(0),
		cadence:   types.MustCadence(1),
		lifecycle: types.LifecycleActive,
		processFunc: func(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
			return nil, fmt.Errorf("processing failed")
		},
	}

	h := primitive.NewHarness()
	_, err := h.Process(p, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestHarnessWithTick(t *testing.T) {
	var receivedTick types.Tick
	p := &testPrimitive{
		id:        types.MustPrimitiveID("tick_reader"),
		layer:     types.MustLayer(0),
		cadence:   types.MustCadence(1),
		lifecycle: types.LifecycleActive,
		processFunc: func(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
			receivedTick = tick
			return nil, nil
		},
	}

	h := primitive.NewHarness().WithTick(types.MustTick(42))
	h.Process(p, nil)

	if receivedTick.Value() != 42 {
		t.Errorf("received tick = %d, want 42", receivedTick.Value())
	}
}

func TestHarnessWithEvents(t *testing.T) {
	var receivedSnap primitive.Snapshot
	p := &testPrimitive{
		id:        types.MustPrimitiveID("snap_reader"),
		layer:     types.MustLayer(0),
		cadence:   types.MustCadence(1),
		lifecycle: types.LifecycleActive,
		processFunc: func(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
			receivedSnap = snap
			return nil, nil
		},
	}

	ev := event.NewEvent(1,
		types.MustEventID("019462a0-0000-7000-8000-000000000099"),
		event.EventTypeTrustUpdated, types.Now(),
		types.MustActorID("actor_system0000000000000000001"),
		event.TrustUpdatedContent{},
		[]types.EventID{types.MustEventID("019462a0-0000-7000-8000-000000000001")},
		types.MustConversationID("conv_test000000000000000000000001"),
		types.ZeroHash(), types.ZeroHash(),
		types.MustSignature(make([]byte, 64)),
	)

	h := primitive.NewHarness().WithEvents([]event.Event{ev})
	h.Process(p, []event.Event{ev})

	if len(receivedSnap.PendingEvents) != 1 {
		t.Errorf("pending events = %d, want 1", len(receivedSnap.PendingEvents))
	}
}

func TestHarnessSnapshotIncludesPrimitiveState(t *testing.T) {
	var receivedSnap primitive.Snapshot
	primID := types.MustPrimitiveID("snap_checker")
	p := &testPrimitive{
		id:        primID,
		layer:     types.MustLayer(0),
		cadence:   types.MustCadence(1),
		lifecycle: types.LifecycleActive,
		processFunc: func(tick types.Tick, events []event.Event, snap primitive.Snapshot) ([]primitive.Mutation, error) {
			receivedSnap = snap
			return nil, nil
		},
	}

	h := primitive.NewHarness()
	h.Process(p, nil)

	ps, ok := receivedSnap.Primitives[primID]
	if !ok {
		t.Fatal("primitive not in snapshot")
	}
	if ps.Lifecycle != types.LifecycleActive {
		t.Errorf("lifecycle = %v, want Active", ps.Lifecycle)
	}
}
