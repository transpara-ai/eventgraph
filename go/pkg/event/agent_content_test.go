package event

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// --- agent.vital.reported ---

func TestAgentVitalReportedContentEventTypeName(t *testing.T) {
	c := AgentVitalReportedContent{}
	if c.EventTypeName() != "agent.vital.reported" {
		t.Errorf("EventTypeName() = %q, want %q", c.EventTypeName(), "agent.vital.reported")
	}
}

func TestAgentVitalReportedContentAccept(t *testing.T) {
	c := AgentVitalReportedContent{}
	c.Accept(nil)
}

func TestAgentVitalReportedContentRoundTrip(t *testing.T) {
	c := AgentVitalReportedContent{
		AgentID:               types.MustActorID("actor_13c4fb4af6dc7d3def47139827601e53"),
		IterationsPct:         0.73,
		TrustScore:            0.91,
		BudgetBurnRatePerHour: 4.5,
		LastHeartbeatTicks:    12,
		Severity:              "OK",
		HealthReportCycleID:   "01HV1A2B3C4D5E6F7G8H9JKMNP",
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	raw := string(data)
	for _, key := range []string{
		`"AgentID"`, `"IterationsPct"`, `"TrustScore"`,
		`"BudgetBurnRatePerHour"`, `"LastHeartbeatTicks"`,
		`"Severity"`, `"HealthReportCycleID"`,
	} {
		if !strings.Contains(raw, key) {
			t.Errorf("serialized JSON missing PascalCase key %s: %s", key, raw)
		}
	}

	got, err := UnmarshalContent("agent.vital.reported", data)
	if err != nil {
		t.Fatalf("UnmarshalContent: %v", err)
	}
	gotC, ok := got.(AgentVitalReportedContent)
	if !ok {
		t.Fatalf("UnmarshalContent returned %T, want AgentVitalReportedContent", got)
	}
	if gotC != c {
		t.Errorf("round-trip mismatch:\n got: %+v\nwant: %+v", gotC, c)
	}
}

func TestAgentVitalReportedContentRegistered(t *testing.T) {
	if !IsKnownEventType("agent.vital.reported") {
		t.Error(`IsKnownEventType("agent.vital.reported") = false, want true`)
	}
	r := DefaultRegistry()
	if !r.IsRegistered(EventTypeAgentVitalReported) {
		t.Errorf("DefaultRegistry() missing %q", EventTypeAgentVitalReported.Value())
	}
}

func TestAllAgentEventTypesContainsVitalReported(t *testing.T) {
	all := AllAgentEventTypes()
	found := false
	for _, et := range all {
		if et == EventTypeAgentVitalReported {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("AllAgentEventTypes() does not contain %q", EventTypeAgentVitalReported.Value())
	}
}
