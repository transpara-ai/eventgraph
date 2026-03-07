package layer11_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive/layer11"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var systemActor = types.MustActorID("actor_00000000000000000000000000000001")

type testSigner struct{}
func (testSigner) Sign(data []byte) (types.Signature, error) {
	sig := make([]byte, 64); copy(sig, data); return types.MustSignature(sig), nil
}

func TestAllPrimitivesRegister(t *testing.T) {
	reg := primitive.NewRegistry()
	prims := []primitive.Primitive{
		layer11.NewSelfAwarenessPrimitive(), layer11.NewPerspectivePrimitive(),
		layer11.NewCritiquePrimitive(), layer11.NewWisdomPrimitive(),
		layer11.NewAestheticPrimitive(), layer11.NewMetaphorPrimitive(),
		layer11.NewHumourPrimitive(), layer11.NewSilencePrimitive(),
		layer11.NewTeachingPrimitive(), layer11.NewTranslationPrimitive(),
		layer11.NewArchivePrimitive(), layer11.NewProphecyPrimitive(),
	}
	for _, p := range prims {
		if err := reg.Register(p); err != nil {
			t.Errorf("Register %q: %v", p.ID().Value(), err)
		}
		if p.Layer().Value() != 11 {
			t.Errorf("%q: Layer = %d, want 11", p.ID().Value(), p.Layer().Value())
		}
		if len(p.Subscriptions()) == 0 {
			t.Errorf("%q: no subscriptions", p.ID().Value())
		}
	}
	if reg.Count() != 12 {
		t.Errorf("registered %d primitives, want 12", reg.Count())
	}
}

func TestWisdomProcess(t *testing.T) {
	s := store.NewInMemoryStore()
	factory := event.NewBootstrapFactory(event.DefaultRegistry())
	bootstrap, _ := factory.Init(systemActor, testSigner{})
	s.Append(bootstrap)
	p := layer11.NewWisdomPrimitive()
	mutations, err := p.Process(types.MustTick(1), []event.Event{bootstrap}, primitive.Snapshot{})
	if err != nil { t.Fatalf("Process: %v", err) }
	if len(mutations) == 0 { t.Fatal("expected mutations") }
}
