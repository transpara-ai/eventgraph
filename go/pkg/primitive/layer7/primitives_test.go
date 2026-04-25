package layer7_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/primitive"
	"github.com/transpara-ai/eventgraph/go/pkg/primitive/layer7"
	"github.com/transpara-ai/eventgraph/go/pkg/store"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
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
		// Group A: Moral Standing
		layer7.NewMoralStatusPrimitive(),
		layer7.NewDignityPrimitive(),
		layer7.NewAutonomyPrimitive(),
		layer7.NewFlourishingPrimitive(),
		// Group B: Moral Obligation
		layer7.NewDutyPrimitive(),
		layer7.NewHarmPrimitive(),
		layer7.NewCarePrimitive(),
		layer7.NewJusticePrimitive(),
		// Group C: Moral Agency
		layer7.NewConsciencePrimitive(),
		layer7.NewVirtuePrimitive(),
		layer7.NewResponsibilityPrimitive(),
		layer7.NewMotivePrimitive(),
	}

	for _, p := range prims {
		if err := reg.Register(p); err != nil {
			t.Errorf("Register %q: %v", p.ID().Value(), err)
		}
		if p.Layer().Value() != 7 {
			t.Errorf("%q: Layer = %d, want 7", p.ID().Value(), p.Layer().Value())
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

func TestMoralStatusProcess(t *testing.T) {
	_, bootstrap := bootstrapStore(t)
	p := layer7.NewMoralStatusPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{bootstrap}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}

func TestHarmProcess(t *testing.T) {
	_, bootstrap := bootstrapStore(t)
	p := layer7.NewHarmPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{bootstrap}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}

func TestCarePrimitive(t *testing.T) {
	_, bootstrap := bootstrapStore(t)
	p := layer7.NewCarePrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{bootstrap}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}

func TestResponsibilityProcess(t *testing.T) {
	_, bootstrap := bootstrapStore(t)
	p := layer7.NewResponsibilityPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{bootstrap}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}
