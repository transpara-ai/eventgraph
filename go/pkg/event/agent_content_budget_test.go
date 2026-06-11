package event

import (
	"encoding/json"
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// v14-F3(c): the budget-adjusted event names WHICH resource was adjusted.
// Legacy payloads (no Resource key) must read as empty — consumers treat
// empty as iterations, the only dimension that existed before.

func TestAgentBudgetAdjustedContentResourceRoundTrip(t *testing.T) {
	aid, err := types.NewActorID("actor_4398821bf5b937fc51b0fb7117ab9c5c")
	if err != nil {
		t.Fatalf("actor id: %v", err)
	}
	in := AgentBudgetAdjustedContent{
		AgentID:   aid,
		AgentName: "implementer",
		Action:    "set",
		NewBudget: 120,
		Resource:  "duration",
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out AgentBudgetAdjustedContent
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Resource != "duration" {
		t.Fatalf("Resource round-trip = %q; want \"duration\"", out.Resource)
	}
}

func TestAgentBudgetAdjustedContentLegacyJSONHasEmptyResource(t *testing.T) {
	var out AgentBudgetAdjustedContent
	if err := json.Unmarshal([]byte(`{"AgentName":"cto","Action":"set","NewBudget":100}`), &out); err != nil {
		t.Fatalf("unmarshal legacy payload: %v", err)
	}
	if out.Resource != "" {
		t.Fatalf("legacy payload Resource = %q; want empty (reads as iterations)", out.Resource)
	}
}
