package event

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// --- ExternalRef ---

func TestExternalRefRoundTrip(t *testing.T) {
	ref := ExternalRef{System: "site", ID: "op_12345"}
	data, err := json.Marshal(ref)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if got := string(data); got != `{"system":"site","id":"op_12345"}` {
		t.Errorf("Marshal = %s, want canonical snake_case form", got)
	}
	var round ExternalRef
	if err := json.Unmarshal(data, &round); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if round != ref {
		t.Errorf("roundtrip = %+v, want %+v", round, ref)
	}
}

// --- site.op.received ---

func TestSiteOpReceivedContentEventTypeName(t *testing.T) {
	c := SiteOpReceivedContent{}
	if c.EventTypeName() != "site.op.received" {
		t.Errorf("EventTypeName() = %q, want %q", c.EventTypeName(), "site.op.received")
	}
}

func TestSiteOpReceivedContentAccept(t *testing.T) {
	c := SiteOpReceivedContent{}
	c.Accept(nil)
}

func TestSiteOpReceivedContentRoundTrip(t *testing.T) {
	now := time.Date(2026, 4, 17, 12, 0, 0, 0, time.UTC)
	siteCreated := now.Add(-time.Second)
	c := SiteOpReceivedContent{
		ExternalRef:   ExternalRef{System: "site", ID: "op_abc"},
		SpaceID:       "space_123",
		NodeID:        "node_456",
		NodeTitle:     "Design Review",
		Actor:         "alice",
		ActorID:       "actor_789",
		ActorKind:     "human",
		OpKind:        "node.updated",
		PayloadHash:   "sha256:deadbeef",
		ReceivedAt:    now,
		SiteCreatedAt: siteCreated,
	}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	raw := string(data)
	for _, key := range []string{
		`"external_ref"`, `"space_id"`, `"node_id"`, `"node_title"`,
		`"actor"`, `"actor_id"`, `"actor_kind"`, `"op_kind"`,
		`"payload_hash"`, `"received_at"`, `"site_created_at"`,
	} {
		if !strings.Contains(raw, key) {
			t.Errorf("serialized JSON missing snake_case key %s: %s", key, raw)
		}
	}
	got, err := UnmarshalContent("site.op.received", data)
	if err != nil {
		t.Fatalf("UnmarshalContent: %v", err)
	}
	typed, ok := got.(SiteOpReceivedContent)
	if !ok {
		t.Fatalf("got type %T, want SiteOpReceivedContent", got)
	}
	if typed.ExternalRef != c.ExternalRef {
		t.Errorf("ExternalRef = %+v, want %+v", typed.ExternalRef, c.ExternalRef)
	}
	if typed.SpaceID != c.SpaceID {
		t.Errorf("SpaceID = %q, want %q", typed.SpaceID, c.SpaceID)
	}
	if typed.NodeID != c.NodeID {
		t.Errorf("NodeID = %q, want %q", typed.NodeID, c.NodeID)
	}
	if typed.NodeTitle != c.NodeTitle {
		t.Errorf("NodeTitle = %q, want %q", typed.NodeTitle, c.NodeTitle)
	}
	if typed.Actor != c.Actor {
		t.Errorf("Actor = %q, want %q", typed.Actor, c.Actor)
	}
	if typed.ActorID != c.ActorID {
		t.Errorf("ActorID = %q, want %q", typed.ActorID, c.ActorID)
	}
	if typed.ActorKind != c.ActorKind {
		t.Errorf("ActorKind = %q, want %q", typed.ActorKind, c.ActorKind)
	}
	if typed.OpKind != c.OpKind {
		t.Errorf("OpKind = %q, want %q", typed.OpKind, c.OpKind)
	}
	if typed.PayloadHash != c.PayloadHash {
		t.Errorf("PayloadHash = %q, want %q", typed.PayloadHash, c.PayloadHash)
	}
	if !typed.ReceivedAt.Equal(c.ReceivedAt) {
		t.Errorf("ReceivedAt = %v, want %v", typed.ReceivedAt, c.ReceivedAt)
	}
	if !typed.SiteCreatedAt.Equal(c.SiteCreatedAt) {
		t.Errorf("SiteCreatedAt = %v, want %v", typed.SiteCreatedAt, c.SiteCreatedAt)
	}
}

// --- site.op.translated ---

func TestSiteOpTranslatedContentEventTypeName(t *testing.T) {
	c := SiteOpTranslatedContent{}
	if c.EventTypeName() != "site.op.translated" {
		t.Errorf("EventTypeName() = %q, want %q", c.EventTypeName(), "site.op.translated")
	}
}

func TestSiteOpTranslatedContentAccept(t *testing.T) {
	c := SiteOpTranslatedContent{}
	c.Accept(nil)
}

func TestSiteOpTranslatedContentRoundTrip(t *testing.T) {
	now := time.Date(2026, 4, 17, 12, 0, 0, 0, time.UTC)
	busEventID := types.MustEventID("01912345-6789-7abc-8def-0123456789ab")
	c := SiteOpTranslatedContent{
		ExternalRef:  ExternalRef{System: "site", ID: "op_abc"},
		BusEventID:   busEventID,
		TranslatedAt: now,
	}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	raw := string(data)
	for _, key := range []string{`"external_ref"`, `"bus_event_id"`, `"translated_at"`} {
		if !strings.Contains(raw, key) {
			t.Errorf("serialized JSON missing snake_case key %s: %s", key, raw)
		}
	}
	got, err := UnmarshalContent("site.op.translated", data)
	if err != nil {
		t.Fatalf("UnmarshalContent: %v", err)
	}
	typed, ok := got.(SiteOpTranslatedContent)
	if !ok {
		t.Fatalf("got type %T, want SiteOpTranslatedContent", got)
	}
	if typed.ExternalRef != c.ExternalRef {
		t.Errorf("ExternalRef = %+v, want %+v", typed.ExternalRef, c.ExternalRef)
	}
	if typed.BusEventID != c.BusEventID {
		t.Errorf("BusEventID = %v, want %v", typed.BusEventID, c.BusEventID)
	}
	if !typed.TranslatedAt.Equal(c.TranslatedAt) {
		t.Errorf("TranslatedAt = %v, want %v", typed.TranslatedAt, c.TranslatedAt)
	}
}

// --- site.op.rejected ---

func TestSiteOpRejectedContentEventTypeName(t *testing.T) {
	c := SiteOpRejectedContent{}
	if c.EventTypeName() != "site.op.rejected" {
		t.Errorf("EventTypeName() = %q, want %q", c.EventTypeName(), "site.op.rejected")
	}
}

func TestSiteOpRejectedContentAccept(t *testing.T) {
	c := SiteOpRejectedContent{}
	c.Accept(nil)
}

func TestSiteOpRejectedContentRoundTrip(t *testing.T) {
	now := time.Date(2026, 4, 17, 12, 0, 0, 0, time.UTC)
	c := SiteOpRejectedContent{
		ExternalRef: ExternalRef{System: "site", ID: "op_abc"},
		Reason:      "unknown op_kind",
		RejectedAt:  now,
	}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	raw := string(data)
	for _, key := range []string{`"external_ref"`, `"reason"`, `"rejected_at"`} {
		if !strings.Contains(raw, key) {
			t.Errorf("serialized JSON missing snake_case key %s: %s", key, raw)
		}
	}
	got, err := UnmarshalContent("site.op.rejected", data)
	if err != nil {
		t.Fatalf("UnmarshalContent: %v", err)
	}
	typed, ok := got.(SiteOpRejectedContent)
	if !ok {
		t.Fatalf("got type %T, want SiteOpRejectedContent", got)
	}
	if typed.ExternalRef != c.ExternalRef {
		t.Errorf("ExternalRef = %+v, want %+v", typed.ExternalRef, c.ExternalRef)
	}
	if typed.Reason != c.Reason {
		t.Errorf("Reason = %q, want %q", typed.Reason, c.Reason)
	}
	if !typed.RejectedAt.Equal(c.RejectedAt) {
		t.Errorf("RejectedAt = %v, want %v", typed.RejectedAt, c.RejectedAt)
	}
}

// --- site.op.mirrored ---

func TestSiteOpMirroredContentEventTypeName(t *testing.T) {
	c := SiteOpMirroredContent{}
	if c.EventTypeName() != "site.op.mirrored" {
		t.Errorf("EventTypeName() = %q, want %q", c.EventTypeName(), "site.op.mirrored")
	}
}

func TestSiteOpMirroredContentAccept(t *testing.T) {
	c := SiteOpMirroredContent{}
	c.Accept(nil)
}

func TestSiteOpMirroredContentRoundTrip(t *testing.T) {
	now := time.Date(2026, 4, 17, 12, 0, 0, 0, time.UTC)
	mirrorEventID := types.MustEventID("01912345-6789-7abc-8def-0123456789ab")
	c := SiteOpMirroredContent{
		ExternalRef:   ExternalRef{System: "site", ID: "op_abc"},
		MirrorEventID: mirrorEventID,
		HTTPStatus:    200,
		MirroredAt:    now,
	}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	raw := string(data)
	for _, key := range []string{
		`"external_ref"`, `"mirror_event_id"`, `"http_status"`, `"mirrored_at"`,
	} {
		if !strings.Contains(raw, key) {
			t.Errorf("serialized JSON missing snake_case key %s: %s", key, raw)
		}
	}
	got, err := UnmarshalContent("site.op.mirrored", data)
	if err != nil {
		t.Fatalf("UnmarshalContent: %v", err)
	}
	typed, ok := got.(SiteOpMirroredContent)
	if !ok {
		t.Fatalf("got type %T, want SiteOpMirroredContent", got)
	}
	if typed.ExternalRef != c.ExternalRef {
		t.Errorf("ExternalRef = %+v, want %+v", typed.ExternalRef, c.ExternalRef)
	}
	if typed.MirrorEventID != c.MirrorEventID {
		t.Errorf("MirrorEventID = %v, want %v", typed.MirrorEventID, c.MirrorEventID)
	}
	if typed.HTTPStatus != c.HTTPStatus {
		t.Errorf("HTTPStatus = %d, want %d", typed.HTTPStatus, c.HTTPStatus)
	}
	if !typed.MirroredAt.Equal(c.MirroredAt) {
		t.Errorf("MirroredAt = %v, want %v", typed.MirroredAt, c.MirroredAt)
	}
}

// --- Event type constants ---

func TestSiteEventTypeConstants(t *testing.T) {
	tests := []struct {
		name  string
		et    types.EventType
		value string
	}{
		{"SiteOpReceived", EventTypeSiteOpReceived, "site.op.received"},
		{"SiteOpTranslated", EventTypeSiteOpTranslated, "site.op.translated"},
		{"SiteOpRejected", EventTypeSiteOpRejected, "site.op.rejected"},
		{"SiteOpMirrored", EventTypeSiteOpMirrored, "site.op.mirrored"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.et.Value() != tt.value {
				t.Errorf("Value() = %q, want %q", tt.et.Value(), tt.value)
			}
		})
	}
}

// --- AllSiteEventTypes ---

func TestAllSiteEventTypes(t *testing.T) {
	all := AllSiteEventTypes()
	if len(all) != 4 {
		t.Fatalf("AllSiteEventTypes() returned %d types, want 4", len(all))
	}
	found := map[string]bool{}
	for _, et := range all {
		found[et.Value()] = true
	}
	for _, want := range []string{
		"site.op.received", "site.op.translated",
		"site.op.rejected", "site.op.mirrored",
	} {
		if !found[want] {
			t.Errorf("AllSiteEventTypes() missing %q", want)
		}
	}
}

// --- DefaultRegistry ---

func TestDefaultRegistryContainsSiteTypes(t *testing.T) {
	r := DefaultRegistry()
	for _, et := range []types.EventType{
		EventTypeSiteOpReceived, EventTypeSiteOpTranslated,
		EventTypeSiteOpRejected, EventTypeSiteOpMirrored,
	} {
		if !r.IsRegistered(et) {
			t.Errorf("DefaultRegistry() missing %q", et.Value())
		}
	}
}

// --- Unmarshaler completeness ---

func TestSiteUnmarshalersRegistered(t *testing.T) {
	for _, name := range []string{
		"site.op.received", "site.op.translated",
		"site.op.rejected", "site.op.mirrored",
	} {
		if !IsKnownEventType(name) {
			t.Errorf("IsKnownEventType(%q) = false, want true", name)
		}
	}
}
