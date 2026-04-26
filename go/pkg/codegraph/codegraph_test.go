package codegraph_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/codegraph"
	"github.com/transpara-ai/eventgraph/go/pkg/event"
)

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

func TestCompositionNamesUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, c := range codegraph.AllCompositions() {
		if seen[c.Name] {
			t.Errorf("duplicate composition name: %q", c.Name)
		}
		seen[c.Name] = true
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
