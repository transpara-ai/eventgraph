package layer13_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive/layer13"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var systemActor = types.MustActorID("actor_00000000000000000000000000000001")

type testSigner struct{}
func (testSigner) Sign(data []byte) (types.Signature, error) {
	sig := make([]byte, 64); copy(sig, data); return types.MustSignature(sig), nil
}

func TestAllPrimitivesRegister(t *testing.T) {
	s := store.NewInMemoryStore()
	reg := primitive.NewRegistry()
	prims := []primitive.Primitive{
		layer13.NewBeingPrimitive(), layer13.NewFinitudePrimitive(),
		layer13.NewChangePrimitive(s), layer13.NewInterdependencePrimitive(s),
		layer13.NewMysteryPrimitive(), layer13.NewParadoxPrimitive(),
		layer13.NewInfinityPrimitive(), layer13.NewVoidPrimitive(),
		layer13.NewAwePrimitive(), layer13.NewExistentialGratitudePrimitive(),
		layer13.NewPlayPrimitive(), layer13.NewWonderPrimitive(),
	}
	for _, p := range prims {
		if err := reg.Register(p); err != nil {
			t.Errorf("Register %q: %v", p.ID().Value(), err)
		}
		if p.Layer().Value() != 13 {
			t.Errorf("%q: Layer = %d, want 13", p.ID().Value(), p.Layer().Value())
		}
		if len(p.Subscriptions()) == 0 {
			t.Errorf("%q: no subscriptions", p.ID().Value())
		}
	}
	if reg.Count() != 12 {
		t.Errorf("registered %d primitives, want 12", reg.Count())
	}
}

func TestBeingProcess(t *testing.T) {
	p := layer13.NewBeingPrimitive()
	mutations, err := p.Process(types.MustTick(1), nil, primitive.Snapshot{})
	if err != nil { t.Fatalf("Process: %v", err) }
	if len(mutations) != 2 { t.Fatalf("expected 2 mutations, got %d", len(mutations)) }
}

func TestWonderProcess(t *testing.T) {
	s := store.NewInMemoryStore()
	factory := event.NewBootstrapFactory(event.DefaultRegistry())
	bootstrap, _ := factory.Init(systemActor, testSigner{})
	s.Append(bootstrap)
	p := layer13.NewWonderPrimitive()
	mutations, err := p.Process(types.MustTick(1), []event.Event{bootstrap}, primitive.Snapshot{})
	if err != nil { t.Fatalf("Process: %v", err) }
	if len(mutations) == 0 { t.Fatal("expected mutations") }
}
