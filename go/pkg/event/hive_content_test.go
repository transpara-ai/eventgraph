package event

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// --- GapCategory ---

func TestGapCategoryIsValid(t *testing.T) {
	valid := []GapCategory{
		GapCategoryLeadership, GapCategoryTechnical, GapCategoryProcess,
		GapCategoryStaffing, GapCategoryCapability,
	}
	for _, c := range valid {
		if !c.IsValid() {
			t.Errorf("expected %q to be valid", c)
		}
	}
	if GapCategory("bogus").IsValid() {
		t.Error("expected bogus GapCategory to be invalid")
	}
}

// --- DirectivePriority ---

func TestDirectivePriorityIsValid(t *testing.T) {
	valid := []DirectivePriority{
		DirectivePriorityCritical, DirectivePriorityHigh,
		DirectivePriorityMedium, DirectivePriorityLow,
	}
	for _, p := range valid {
		if !p.IsValid() {
			t.Errorf("expected %q to be valid", p)
		}
	}
	if DirectivePriority("bogus").IsValid() {
		t.Error("expected bogus DirectivePriority to be invalid")
	}
}

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
	c.Accept(nil)
}

func TestGapDetectedContentRoundTrip(t *testing.T) {
	c := GapDetectedContent{
		Category:    GapCategoryLeadership,
		MissingRole: "CTO",
		Evidence:    "no technical decisions in 30 days",
		Severity:    SeverityLevelSerious,
	}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	// Verify PascalCase JSON keys.
	raw := string(data)
	for _, key := range []string{`"Category"`, `"MissingRole"`, `"Evidence"`, `"Severity"`} {
		if !strings.Contains(raw, key) {
			t.Errorf("serialized JSON missing PascalCase key %s: %s", key, raw)
		}
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
	c := NewGapDetectedContent(GapCategoryLeadership, "CTO", "no technical decisions in 30 days", SeverityLevelSerious)
	if c.Category != GapCategoryLeadership {
		t.Errorf("Category = %q, want %q", c.Category, GapCategoryLeadership)
	}
	if c.MissingRole != "CTO" {
		t.Errorf("MissingRole = %q, want %q", c.MissingRole, "CTO")
	}
	if c.Evidence != "no technical decisions in 30 days" {
		t.Errorf("Evidence = %q, want %q", c.Evidence, "no technical decisions in 30 days")
	}
	if c.Severity != SeverityLevelSerious {
		t.Errorf("Severity = %q, want %q", c.Severity, SeverityLevelSerious)
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
		Priority: DirectivePriorityCritical,
	}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	// Verify PascalCase JSON keys.
	raw := string(data)
	for _, key := range []string{`"Target"`, `"Action"`, `"Reason"`, `"Priority"`} {
		if !strings.Contains(raw, key) {
			t.Errorf("serialized JSON missing PascalCase key %s: %s", key, raw)
		}
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
	c := NewDirectiveIssuedContent("engineering-team", "hire-cto", "gap detected in technical leadership", DirectivePriorityCritical)
	if c.Target != "engineering-team" {
		t.Errorf("Target = %q, want %q", c.Target, "engineering-team")
	}
	if c.Action != "hire-cto" {
		t.Errorf("Action = %q, want %q", c.Action, "hire-cto")
	}
	if c.Reason != "gap detected in technical leadership" {
		t.Errorf("Reason = %q, want %q", c.Reason, "gap detected in technical leadership")
	}
	if c.Priority != DirectivePriorityCritical {
		t.Errorf("Priority = %q, want %q", c.Priority, DirectivePriorityCritical)
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
