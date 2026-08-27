package event

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

const testTLC51DigestA = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
const testTLC51DigestB = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

func testTLC51Identity(attempt uint32) TLC51EventIdentity {
	return TLC51EventIdentity{
		ProtocolVersion: TLC51ProtocolVersion,
		FactoryOrderID:  "FO-TLC51-TEST-001",
		ChangeSeriesID:  "series-1",
		PlanDigest:      testTLC51DigestA,
		SubjectDigest:   testTLC51DigestB,
		EventOrdinal:    1,
		AttemptOrdinal:  attempt,
	}
}

func testTLC51JSON(schema string, fields string) TLC51ExactJSON {
	input := []byte(fmt.Sprintf("{\"schema_version\":\"%s\"%s}", schema, fields))
	var value any
	if err := json.Unmarshal(input, &value); err != nil {
		panic(err)
	}
	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		panic(err)
	}
	raw := output.Bytes()
	return TLC51ExactJSON{SchemaVersion: schema, CanonicalJSON: string(raw), SHA256: fmt.Sprintf("%x", sha256.Sum256(raw))}
}

func TestAllTLC51EventTypesAreClosedAndRegistered(t *testing.T) {
	all := AllTLC51EventTypes()
	if len(all) != 16 {
		t.Fatalf("AllTLC51EventTypes length = %d, want 16", len(all))
	}
	seen := map[string]bool{}
	registry := DefaultRegistry()
	for _, eventType := range all {
		if seen[eventType.Value()] {
			t.Fatalf("duplicate TLC 5.1 event type %q", eventType.Value())
		}
		seen[eventType.Value()] = true
		if !registry.IsRegistered(eventType) || !IsKnownEventType(eventType.Value()) {
			t.Errorf("TLC 5.1 event type %q not registered", eventType.Value())
		}
	}
}

func TestTLC51PlanRecordedValidatesExactPayloadAndRoundTrips(t *testing.T) {
	content := TLC51PlanRecordedContent{
		TLC51EventIdentity: testTLC51Identity(0),
		Plan:               testTLC51JSON("tlc-gate-plan/v1", `,"plan_digest":"`+testTLC51DigestA+`"`),
		RecordedAt:         time.Date(2026, 8, 27, 17, 0, 0, 0, time.UTC),
	}
	if err := DefaultRegistry().Validate(EventTypeTLC51PlanRecorded, content); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	raw, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	decoded, err := UnmarshalContent(EventTypeTLC51PlanRecorded.Value(), raw)
	if err != nil {
		t.Fatalf("UnmarshalContent: %v", err)
	}
	round, ok := decoded.(TLC51PlanRecordedContent)
	if !ok || round.Plan.CanonicalJSON != content.Plan.CanonicalJSON || round.Plan.SHA256 != content.Plan.SHA256 {
		t.Fatalf("roundtrip = %#v", decoded)
	}
}

func TestTLC51ValidationRejectsAttemptAndPayloadForgery(t *testing.T) {
	ready := TLC51ObligationReadyContent{
		TLC51EventIdentity: testTLC51Identity(0),
		ObligationID:       "O001-source-relationship-facts",
		ReadyAt:            time.Date(2026, 8, 27, 17, 0, 0, 0, time.UTC),
	}
	if err := DefaultRegistry().Validate(EventTypeTLC51ObligationReady, ready); err == nil {
		t.Fatal("zero attempt ordinal unexpectedly validated")
	}
	plan := TLC51PlanRecordedContent{
		TLC51EventIdentity: testTLC51Identity(0),
		Plan:               testTLC51JSON("tlc-gate-plan/v1", `,"plan_digest":"`+testTLC51DigestA+`"`),
		RecordedAt:         time.Date(2026, 8, 27, 17, 0, 0, 0, time.UTC),
	}
	raw := []byte(plan.Plan.CanonicalJSON)
	raw[5] ^= 1
	plan.Plan.CanonicalJSON = string(raw)
	if err := DefaultRegistry().Validate(EventTypeTLC51PlanRecorded, plan); err == nil {
		t.Fatal("forged exact payload unexpectedly validated")
	}
}

func TestTLC51SuccessfulEffectRequiresReceipt(t *testing.T) {
	content := TLC51ProtectedEffectTerminalContent{
		TLC51EventIdentity: testTLC51Identity(1),
		Effect:             "merge",
		OperationID:        "operation-1",
		Outcome:            TLC51TerminalSucceeded,
		Reason:             "provider state exact",
		TerminalAt:         time.Date(2026, 8, 27, 17, 0, 0, 0, time.UTC),
	}
	if err := DefaultRegistry().Validate(EventTypeTLC51ProtectedEffectTerminal, content); err == nil {
		t.Fatal("successful protected effect without receipt unexpectedly validated")
	}
}
