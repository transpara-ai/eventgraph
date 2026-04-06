package event

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// --- knowledge.insight.recorded ---

func TestKnowledgeInsightContentEventTypeName(t *testing.T) {
	c := KnowledgeInsightContent{}
	if c.EventTypeName() != "knowledge.insight.recorded" {
		t.Errorf("EventTypeName() = %q, want %q", c.EventTypeName(), "knowledge.insight.recorded")
	}
}

func TestKnowledgeInsightContentAccept(t *testing.T) {
	c := KnowledgeInsightContent{}
	c.Accept(nil)
}

func TestKnowledgeInsightContentRoundTrip(t *testing.T) {
	c := NewKnowledgeInsightContent(
		"insight-001", "infrastructure", "cache hit rates drop under load",
		[]string{"sre", "backend"},
		types.MustScore(0.85),
		42,
		"pattern-distiller",
		3600,
		types.None[string](),
	)
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	raw := string(data)
	for _, key := range []string{
		`"insight_id"`, `"domain"`, `"summary"`, `"relevant_roles"`,
		`"confidence"`, `"evidence_count"`, `"source"`, `"ttl"`,
	} {
		if !strings.Contains(raw, key) {
			t.Errorf("serialized JSON missing key %s: %s", key, raw)
		}
	}
	got, err := UnmarshalContent("knowledge.insight.recorded", data)
	if err != nil {
		t.Fatalf("UnmarshalContent: %v", err)
	}
	typed, ok := got.(KnowledgeInsightContent)
	if !ok {
		t.Fatalf("got type %T, want KnowledgeInsightContent", got)
	}
	if typed.InsightID != c.InsightID {
		t.Errorf("InsightID = %q, want %q", typed.InsightID, c.InsightID)
	}
	if typed.Domain != c.Domain {
		t.Errorf("Domain = %q, want %q", typed.Domain, c.Domain)
	}
	if typed.Summary != c.Summary {
		t.Errorf("Summary = %q, want %q", typed.Summary, c.Summary)
	}
	if typed.Confidence != c.Confidence {
		t.Errorf("Confidence = %v, want %v", typed.Confidence, c.Confidence)
	}
	if typed.EvidenceCount != c.EvidenceCount {
		t.Errorf("EvidenceCount = %d, want %d", typed.EvidenceCount, c.EvidenceCount)
	}
	if typed.Source != c.Source {
		t.Errorf("Source = %q, want %q", typed.Source, c.Source)
	}
	if typed.TTL != c.TTL {
		t.Errorf("TTL = %d, want %d", typed.TTL, c.TTL)
	}
	if len(typed.RelevantRoles) != len(c.RelevantRoles) {
		t.Errorf("RelevantRoles len = %d, want %d", len(typed.RelevantRoles), len(c.RelevantRoles))
	}
}

func TestNewKnowledgeInsightContentSortsRoles(t *testing.T) {
	c := NewKnowledgeInsightContent(
		"insight-002", "ops", "summary",
		[]string{"zulu", "alpha", "mike"},
		types.MustScore(0.5),
		10,
		"distiller",
		1800,
		types.None[string](),
	)
	if len(c.RelevantRoles) != 3 {
		t.Fatalf("RelevantRoles len = %d, want 3", len(c.RelevantRoles))
	}
	if c.RelevantRoles[0] != "alpha" {
		t.Errorf("RelevantRoles[0] = %q, want %q", c.RelevantRoles[0], "alpha")
	}
	if c.RelevantRoles[1] != "mike" {
		t.Errorf("RelevantRoles[1] = %q, want %q", c.RelevantRoles[1], "mike")
	}
	if c.RelevantRoles[2] != "zulu" {
		t.Errorf("RelevantRoles[2] = %q, want %q", c.RelevantRoles[2], "zulu")
	}
}

func TestNewKnowledgeInsightContentDoesNotMutateInput(t *testing.T) {
	input := []string{"beta", "alpha"}
	_ = NewKnowledgeInsightContent(
		"insight-003", "ops", "summary",
		input,
		types.MustScore(0.5),
		5,
		"distiller",
		900,
		types.None[string](),
	)
	if input[0] != "beta" {
		t.Errorf("input[0] = %q, want %q — constructor mutated the input slice", input[0], "beta")
	}
}

func TestKnowledgeInsightContentWithSupersedes(t *testing.T) {
	c := NewKnowledgeInsightContent(
		"insight-004", "security", "token rotation interval too long",
		[]string{"security"},
		types.MustScore(0.9),
		15,
		"distiller",
		7200,
		types.Some("insight-001"),
	)
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	raw := string(data)
	if !strings.Contains(raw, `"supersedes_id"`) {
		t.Errorf("serialized JSON missing supersedes_id key: %s", raw)
	}
	got, err := UnmarshalContent("knowledge.insight.recorded", data)
	if err != nil {
		t.Fatalf("UnmarshalContent: %v", err)
	}
	typed, ok := got.(KnowledgeInsightContent)
	if !ok {
		t.Fatalf("got type %T, want KnowledgeInsightContent", got)
	}
	if !typed.SupersedesID.IsSome() {
		t.Fatal("SupersedesID should be Some")
	}
	if typed.SupersedesID.Unwrap() != "insight-001" {
		t.Errorf("SupersedesID = %q, want %q", typed.SupersedesID.Unwrap(), "insight-001")
	}
}

// --- knowledge.insight.superseded ---

func TestKnowledgeSupersessionContentEventTypeName(t *testing.T) {
	c := KnowledgeSupersessionContent{}
	if c.EventTypeName() != "knowledge.insight.superseded" {
		t.Errorf("EventTypeName() = %q, want %q", c.EventTypeName(), "knowledge.insight.superseded")
	}
}

func TestKnowledgeSupersessionContentAccept(t *testing.T) {
	c := KnowledgeSupersessionContent{}
	c.Accept(nil)
}

func TestKnowledgeSupersessionContentRoundTrip(t *testing.T) {
	c := KnowledgeSupersessionContent{
		OldInsightID: "insight-001",
		NewInsightID: "insight-004",
		Reason:       "newer data available with higher confidence",
	}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	raw := string(data)
	for _, key := range []string{`"old_insight_id"`, `"new_insight_id"`, `"reason"`} {
		if !strings.Contains(raw, key) {
			t.Errorf("serialized JSON missing key %s: %s", key, raw)
		}
	}
	got, err := UnmarshalContent("knowledge.insight.superseded", data)
	if err != nil {
		t.Fatalf("UnmarshalContent: %v", err)
	}
	typed, ok := got.(KnowledgeSupersessionContent)
	if !ok {
		t.Fatalf("got type %T, want KnowledgeSupersessionContent", got)
	}
	if typed.OldInsightID != c.OldInsightID {
		t.Errorf("OldInsightID = %q, want %q", typed.OldInsightID, c.OldInsightID)
	}
	if typed.NewInsightID != c.NewInsightID {
		t.Errorf("NewInsightID = %q, want %q", typed.NewInsightID, c.NewInsightID)
	}
	if typed.Reason != c.Reason {
		t.Errorf("Reason = %q, want %q", typed.Reason, c.Reason)
	}
}

// --- knowledge.insight.expired ---

func TestKnowledgeExpirationContentEventTypeName(t *testing.T) {
	c := KnowledgeExpirationContent{}
	if c.EventTypeName() != "knowledge.insight.expired" {
		t.Errorf("EventTypeName() = %q, want %q", c.EventTypeName(), "knowledge.insight.expired")
	}
}

func TestKnowledgeExpirationContentAccept(t *testing.T) {
	c := KnowledgeExpirationContent{}
	c.Accept(nil)
}

func TestKnowledgeExpirationContentRoundTrip(t *testing.T) {
	c := KnowledgeExpirationContent{
		InsightID: "insight-001",
		Reason:    "TTL elapsed after 3600 seconds",
	}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	raw := string(data)
	for _, key := range []string{`"insight_id"`, `"reason"`} {
		if !strings.Contains(raw, key) {
			t.Errorf("serialized JSON missing key %s: %s", key, raw)
		}
	}
	got, err := UnmarshalContent("knowledge.insight.expired", data)
	if err != nil {
		t.Fatalf("UnmarshalContent: %v", err)
	}
	typed, ok := got.(KnowledgeExpirationContent)
	if !ok {
		t.Fatalf("got type %T, want KnowledgeExpirationContent", got)
	}
	if typed.InsightID != c.InsightID {
		t.Errorf("InsightID = %q, want %q", typed.InsightID, c.InsightID)
	}
	if typed.Reason != c.Reason {
		t.Errorf("Reason = %q, want %q", typed.Reason, c.Reason)
	}
}

// --- Event type constants ---

func TestKnowledgeEventTypeConstants(t *testing.T) {
	tests := []struct {
		name  string
		et    types.EventType
		value string
	}{
		{"InsightRecorded", EventTypeKnowledgeInsightRecorded, "knowledge.insight.recorded"},
		{"InsightSuperseded", EventTypeKnowledgeInsightSuperseded, "knowledge.insight.superseded"},
		{"InsightExpired", EventTypeKnowledgeInsightExpired, "knowledge.insight.expired"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.et.Value() != tt.value {
				t.Errorf("Value() = %q, want %q", tt.et.Value(), tt.value)
			}
		})
	}
}

// --- AllKnowledgeEventTypes ---

func TestAllKnowledgeEventTypesContainsAll(t *testing.T) {
	all := AllKnowledgeEventTypes()
	if len(all) != 3 {
		t.Fatalf("AllKnowledgeEventTypes() returned %d types, want 3", len(all))
	}
	found := map[string]bool{}
	for _, et := range all {
		found[et.Value()] = true
	}
	for _, want := range []string{
		"knowledge.insight.recorded",
		"knowledge.insight.superseded",
		"knowledge.insight.expired",
	} {
		if !found[want] {
			t.Errorf("AllKnowledgeEventTypes() missing %q", want)
		}
	}
}

// --- DefaultRegistry ---

func TestDefaultRegistryContainsKnowledgeTypes(t *testing.T) {
	r := DefaultRegistry()
	for _, et := range []types.EventType{
		EventTypeKnowledgeInsightRecorded,
		EventTypeKnowledgeInsightSuperseded,
		EventTypeKnowledgeInsightExpired,
	} {
		if !r.IsRegistered(et) {
			t.Errorf("DefaultRegistry() missing %q", et.Value())
		}
	}
}
