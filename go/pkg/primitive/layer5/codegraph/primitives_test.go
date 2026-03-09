package codegraph_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive/layer5/codegraph"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var systemActor = types.MustActorID("actor_00000000000000000000000000000001")

type testSigner struct{}

func (testSigner) Sign(data []byte) (types.Signature, error) {
	sig := make([]byte, 64)
	copy(sig, data)
	return types.MustSignature(sig), nil
}

func TestAllPrimitivesRegister(t *testing.T) {
	reg := primitive.NewRegistry()
	prims := codegraph.AllPrimitives()

	if len(prims) != 61 {
		t.Fatalf("AllPrimitives() returned %d, want 61", len(prims))
	}

	for _, p := range prims {
		if err := reg.Register(p); err != nil {
			t.Errorf("Register %q: %v", p.ID().Value(), err)
		}
		if p.Layer().Value() != 5 {
			t.Errorf("%q: Layer = %d, want 5", p.ID().Value(), p.Layer().Value())
		}
		if len(p.Subscriptions()) == 0 {
			t.Errorf("%q: no subscriptions", p.ID().Value())
		}
	}
	if reg.Count() != 61 {
		t.Errorf("registered %d primitives, want 61", reg.Count())
	}
}

func TestNoDuplicateIDs(t *testing.T) {
	seen := map[string]bool{}
	for _, p := range codegraph.AllPrimitives() {
		id := p.ID().Value()
		if seen[id] {
			t.Errorf("duplicate primitive ID: %q", id)
		}
		seen[id] = true
	}
}

func TestEntityProcess(t *testing.T) {
	s := store.NewInMemoryStore()
	factory := event.NewBootstrapFactory(event.DefaultRegistry())
	bootstrap, _ := factory.Init(systemActor, testSigner{})
	s.Append(bootstrap)

	p := codegraph.NewEntityPrimitive()
	mutations, err := p.Process(types.MustTick(1), []event.Event{bootstrap}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) != 2 {
		t.Errorf("got %d mutations, want 2", len(mutations))
	}
}

func TestCommandProcess(t *testing.T) {
	p := codegraph.NewCommandPrimitive()
	mutations, err := p.Process(types.MustTick(5), nil, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) != 2 {
		t.Errorf("got %d mutations, want 2", len(mutations))
	}
}

func TestAuditSubscribesAll(t *testing.T) {
	p := codegraph.NewAuditPrimitive()
	subs := p.Subscriptions()
	if len(subs) != 1 {
		t.Fatalf("Audit subscriptions = %d, want 1 (codegraph.*)", len(subs))
	}
	if subs[0].Value() != "codegraph.*" {
		t.Errorf("Audit subscription = %q, want codegraph.*", subs[0].Value())
	}
}

func TestAllCompositions(t *testing.T) {
	compositions := codegraph.AllCompositions()
	if len(compositions) != 7 {
		t.Errorf("AllCompositions() = %d, want 7", len(compositions))
	}

	names := map[string]bool{}
	for _, c := range compositions {
		if c.Name == "" {
			t.Error("composition has empty name")
		}
		if c.Purpose == "" {
			t.Errorf("composition %q has empty purpose", c.Name)
		}
		if len(c.Primitives) == 0 {
			t.Errorf("composition %q has no primitives", c.Name)
		}
		if names[c.Name] {
			t.Errorf("duplicate composition name: %q", c.Name)
		}
		names[c.Name] = true
	}
}

func TestCompositionPrimitivesExist(t *testing.T) {
	reg := primitive.NewRegistry()
	for _, p := range codegraph.AllPrimitives() {
		reg.Register(p)
	}

	for _, c := range codegraph.AllCompositions() {
		for _, pid := range c.Primitives {
			if _, ok := reg.Get(pid); !ok {
				t.Errorf("composition %q references unregistered primitive %q", c.Name, pid.Value())
			}
		}
	}
}

func TestCodeGraphEventTypesRegistered(t *testing.T) {
	etypes := event.AllCodeGraphEventTypes()
	if len(etypes) != 35 {
		t.Errorf("AllCodeGraphEventTypes() = %d, want 35", len(etypes))
	}

	reg := event.DefaultRegistry()
	for _, et := range etypes {
		if !reg.IsRegistered(et) {
			t.Errorf("event type %q not registered in DefaultRegistry", et.Value())
		}
	}
}
