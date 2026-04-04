package event

import (
	"encoding/json"
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// --- hive.gap.detected ---

func TestGapDetectedContentEventTypeName(t *testing.T) {
	c := GapDetectedContent{}
	if c.EventTypeName() != "hive.gap.detected" {
		t.Errorf("EventTypeName() = %q, want %q", c.EventTypeName(), "hive.gap.detected")
	}
}

func TestGapDetectedContentAccept(t *testing.T) {
	// Accept on hive content is a no-op (like agentContent) — just ensure it compiles and does not panic.
	c := GapDetectedContent{}
	c.Accept(nil) // no-op; hive content does not dispatch to the base visitor
}

func TestGapDetectedContentRoundTrip(t *testing.T) {
	c := GapDetectedContent{
		Category:    "leadership",
		MissingRole: "CTO",
		Evidence:    "no technical decisions in 30 days",
		Severity:    "high",
	}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got, err := UnmarshalContent("hive.gap.detected", data)
	if err != nil {
		t.Fatalf("UnmarshalContent: %v", err)
	}
	typed, ok := got.(GapDetectedContent)
	if !ok {
		t.Fatalf("got type %T, want GapDetectedContent", got)
	}
	if typed.Category != c.Category {
		t.Errorf("Category = %q, want %q", typed.Category, c.Category)
	}
	if typed.MissingRole != c.MissingRole {
		t.Errorf("MissingRole = %q, want %q", typed.MissingRole, c.MissingRole)
	}
	if typed.Evidence != c.Evidence {
		t.Errorf("Evidence = %q, want %q", typed.Evidence, c.Evidence)
	}
	if typed.Severity != c.Severity {
		t.Errorf("Severity = %q, want %q", typed.Severity, c.Severity)
	}
}

func TestNewGapDetectedContent(t *testing.T) {
	c := NewGapDetectedContent("leadership", "CTO", "no technical decisions in 30 days", "high")
	if c.Category != "leadership" {
		t.Errorf("Category = %q, want %q", c.Category, "leadership")
	}
	if c.MissingRole != "CTO" {
		t.Errorf("MissingRole = %q, want %q", c.MissingRole, "CTO")
	}
	if c.Evidence != "no technical decisions in 30 days" {
		t.Errorf("Evidence = %q, want %q", c.Evidence, "no technical decisions in 30 days")
	}
	if c.Severity != "high" {
		t.Errorf("Severity = %q, want %q", c.Severity, "high")
	}
}

// --- hive.directive.issued ---

func TestDirectiveIssuedContentEventTypeName(t *testing.T) {
	c := DirectiveIssuedContent{}
	if c.EventTypeName() != "hive.directive.issued" {
		t.Errorf("EventTypeName() = %q, want %q", c.EventTypeName(), "hive.directive.issued")
	}
}

func TestDirectiveIssuedContentAccept(t *testing.T) {
	c := DirectiveIssuedContent{}
	c.Accept(nil)
}

func TestDirectiveIssuedContentRoundTrip(t *testing.T) {
	c := DirectiveIssuedContent{
		Target:   "engineering-team",
		Action:   "hire-cto",
		Reason:   "gap detected in technical leadership",
		Priority: "critical",
	}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got, err := UnmarshalContent("hive.directive.issued", data)
	if err != nil {
		t.Fatalf("UnmarshalContent: %v", err)
	}
	typed, ok := got.(DirectiveIssuedContent)
	if !ok {
		t.Fatalf("got type %T, want DirectiveIssuedContent", got)
	}
	if typed.Target != c.Target {
		t.Errorf("Target = %q, want %q", typed.Target, c.Target)
	}
	if typed.Action != c.Action {
		t.Errorf("Action = %q, want %q", typed.Action, c.Action)
	}
	if typed.Reason != c.Reason {
		t.Errorf("Reason = %q, want %q", typed.Reason, c.Reason)
	}
	if typed.Priority != c.Priority {
		t.Errorf("Priority = %q, want %q", typed.Priority, c.Priority)
	}
}

func TestNewDirectiveIssuedContent(t *testing.T) {
	c := NewDirectiveIssuedContent("engineering-team", "hire-cto", "gap detected in technical leadership", "critical")
	if c.Target != "engineering-team" {
		t.Errorf("Target = %q, want %q", c.Target, "engineering-team")
	}
	if c.Action != "hire-cto" {
		t.Errorf("Action = %q, want %q", c.Action, "hire-cto")
	}
	if c.Reason != "gap detected in technical leadership" {
		t.Errorf("Reason = %q, want %q", c.Reason, "gap detected in technical leadership")
	}
	if c.Priority != "critical" {
		t.Errorf("Priority = %q, want %q", c.Priority, "critical")
	}
}

// --- Event type constants ---

func TestHiveEventTypeConstants(t *testing.T) {
	tests := []struct {
		name  string
		et    types.EventType
		value string
	}{
		{"GapDetected", EventTypeGapDetected, "hive.gap.detected"},
		{"DirectiveIssued", EventTypeDirectiveIssued, "hive.directive.issued"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.et.Value() != tt.value {
				t.Errorf("Value() = %q, want %q", tt.et.Value(), tt.value)
			}
		})
	}
}

// --- AllHiveEventTypes ---

func TestAllHiveEventTypesContainsBoth(t *testing.T) {
	all := AllHiveEventTypes()
	if len(all) != 2 {
		t.Fatalf("AllHiveEventTypes() returned %d types, want 2", len(all))
	}
	found := map[string]bool{}
	for _, et := range all {
		found[et.Value()] = true
	}
	for _, want := range []string{"hive.gap.detected", "hive.directive.issued"} {
		if !found[want] {
			t.Errorf("AllHiveEventTypes() missing %q", want)
		}
	}
}

// --- DefaultRegistry ---

func TestDefaultRegistryContainsHiveTypes(t *testing.T) {
	r := DefaultRegistry()
	for _, et := range []types.EventType{EventTypeGapDetected, EventTypeDirectiveIssued} {
		if !r.IsRegistered(et) {
			t.Errorf("DefaultRegistry() missing %q", et.Value())
		}
	}
}
