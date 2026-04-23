package layer6_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/primitive"
	"github.com/transpara-ai/eventgraph/go/pkg/primitive/layer6"
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
		// Group A: Representation
		layer6.NewSymbolPrimitive(),
		layer6.NewLanguagePrimitive(),
		layer6.NewEncodingPrimitive(),
		layer6.NewRecordPrimitive(),
		// Group B: Dynamics
		layer6.NewChannelPrimitive(),
		layer6.NewCopyPrimitive(),
		layer6.NewNoisePrimitive(),
		layer6.NewRedundancyPrimitive(),
		// Group C: Transformation
		layer6.NewDataPrimitive(),
		layer6.NewComputationPrimitive(),
		layer6.NewAlgorithmPrimitive(),
		layer6.NewEntropyPrimitive(),
	}

	for _, p := range prims {
		if err := reg.Register(p); err != nil {
			t.Errorf("Register %q: %v", p.ID().Value(), err)
		}
		if p.Layer().Value() != 6 {
			t.Errorf("%q: Layer = %d, want 6", p.ID().Value(), p.Layer().Value())
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

func TestSymbolProcess(t *testing.T) {
	_, bootstrap := bootstrapStore(t)
	p := layer6.NewSymbolPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{bootstrap}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}

func TestChannelProcess(t *testing.T) {
	_, bootstrap := bootstrapStore(t)
	p := layer6.NewChannelPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{bootstrap}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}

func TestComputationProcess(t *testing.T) {
	_, bootstrap := bootstrapStore(t)
	p := layer6.NewComputationPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{bootstrap}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}

func TestEntropyProcess(t *testing.T) {
	_, bootstrap := bootstrapStore(t)
	p := layer6.NewEntropyPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{bootstrap}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}
