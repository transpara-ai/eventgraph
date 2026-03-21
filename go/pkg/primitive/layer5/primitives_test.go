package layer5_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive/layer5"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var (
	systemActor = types.MustActorID("actor_00000000000000000000000000000001")
	actor2      = types.MustActorID("actor_00000000000000000000000000000002")
	convID      = types.MustConversationID("conv_00000000000000000000000000000001")
)

type testSigner struct{}

func (testSigner) Sign(data []byte) (types.Signature, error) {
	sig := make([]byte, 64)
	copy(sig, data)
	return types.MustSignature(sig), nil
}

type headFromStore struct{ s store.Store }

func (h headFromStore) Head() (types.Option[event.Event], error) { return h.s.Head() }

func bootstrapStore(t *testing.T) (store.Store, event.Event) {
	t.Helper()
	s := store.NewInMemoryStore()
	factory := event.NewBootstrapFactory(event.DefaultRegistry())
	ev, err := factory.Init(systemActor, testSigner{})
	if err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	if _, err := s.Append(ev); err != nil {
		t.Fatalf("append: %v", err)
	}
	return s, ev
}

func TestAllPrimitivesRegister(t *testing.T) {
	reg := primitive.NewRegistry()

	prims := []primitive.Primitive{
		// Group A: Investigation
		layer5.NewMethodPrimitive(),
		layer5.NewMeasurementPrimitive(),
		layer5.NewKnowledgePrimitive(),
		layer5.NewModelPrimitive(),
		// Group B: Creation
		layer5.NewToolPrimitive(),
		layer5.NewTechniquePrimitive(),
		layer5.NewInventionPrimitive(),
		layer5.NewAbstractionPrimitive(),
		// Group C: Systems
		layer5.NewInfrastructurePrimitive(),
		layer5.NewStandardPrimitive(),
		layer5.NewEfficiencyPrimitive(),
		layer5.NewAutomationPrimitive(),
	}

	for _, p := range prims {
		if err := reg.Register(p); err != nil {
			t.Errorf("Register %q: %v", p.ID().Value(), err)
		}
		if p.Layer().Value() != 5 {
			t.Errorf("%q: Layer = %d, want 5", p.ID().Value(), p.Layer().Value())
		}
		if p.Lifecycle() != types.LifecycleActive {
			t.Errorf("%q: Lifecycle = %v, want Active", p.ID().Value(), p.Lifecycle())
		}
		if len(p.Subscriptions()) == 0 {
			t.Errorf("%q: no subscriptions", p.ID().Value())
		}
	}

	if reg.Count() != 12 {
		t.Errorf("registered %d primitives, want 12", reg.Count())
	}
}

func TestMethodProcess(t *testing.T) {
	_, bootstrap := bootstrapStore(t)
	p := layer5.NewMethodPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{bootstrap}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}

func TestToolProcess(t *testing.T) {
	_, bootstrap := bootstrapStore(t)
	p := layer5.NewToolPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{bootstrap}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}

func TestAutomationProcess(t *testing.T) {
	_, bootstrap := bootstrapStore(t)
	p := layer5.NewAutomationPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{bootstrap}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}

func TestInfrastructureProcess(t *testing.T) {
	_, bootstrap := bootstrapStore(t)
	p := layer5.NewInfrastructurePrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{bootstrap}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}
