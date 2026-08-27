package v39

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
)

const historyDigestA = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
const historyDigestB = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

func historyIdentity(ordinal uint64, attempt uint32) event.TLC51EventIdentity {
	return event.TLC51EventIdentity{
		ProtocolVersion: event.TLC51ProtocolVersion,
		FactoryOrderID:  "FO-HISTORY-001",
		ChangeSeriesID:  "series-1",
		PlanDigest:      historyDigestA,
		SubjectDigest:   historyDigestB,
		EventOrdinal:    ordinal,
		AttemptOrdinal:  attempt,
	}
}

func historyPlan() event.TLC51ExactJSON {
	raw := []byte(`{"plan_digest":"` + historyDigestA + `","schema_version":"tlc-gate-plan/v1"}` + "\n")
	return event.TLC51ExactJSON{
		SchemaVersion: "tlc-gate-plan/v1",
		CanonicalJSON: string(raw),
		SHA256:        fmt.Sprintf("%x", sha256.Sum256(raw)),
	}
}

func TestTLC51HistoryAppendIdempotenceOrderAndChain(t *testing.T) {
	now := time.Date(2026, 8, 27, 17, 0, 0, 0, time.UTC)
	history := NewTLC51History()
	plan := event.TLC51PlanRecordedContent{TLC51EventIdentity: historyIdentity(1, 0), Plan: historyPlan(), RecordedAt: now}
	first, err := history.AppendTLC51(plan, now)
	if err != nil {
		t.Fatalf("Append plan: %v", err)
	}
	repeated, err := history.AppendTLC51(plan, now.Add(time.Minute))
	if err != nil || repeated.EntryDigest != first.EntryDigest || repeated.RecordedAt != first.RecordedAt {
		t.Fatalf("idempotent append = %+v, %v", repeated, err)
	}
	ready := event.TLC51ObligationReadyContent{
		TLC51EventIdentity: historyIdentity(2, 1),
		ObligationID:       "O001-source-relationship-facts",
		ReadyAt:            now.Add(time.Second),
	}
	second, err := history.AppendTLC51(ready, now.Add(time.Second))
	if err != nil {
		t.Fatalf("Append ready: %v", err)
	}
	if second.PreviousDigest != first.EntryDigest {
		t.Fatalf("PreviousDigest = %q, want %q", second.PreviousDigest, first.EntryDigest)
	}
	if err := history.Verify("FO-HISTORY-001", "series-1"); err != nil {
		t.Fatalf("Verify: %v", err)
	}
	rows := history.Entries("FO-HISTORY-001", "series-1")
	rows[0].Payload[0] ^= 1
	if err := history.Verify("FO-HISTORY-001", "series-1"); err != nil {
		t.Fatalf("caller mutation changed stored history: %v", err)
	}
}

func TestTLC51HistoryRejectsGapAndConflict(t *testing.T) {
	now := time.Date(2026, 8, 27, 17, 0, 0, 0, time.UTC)
	history := NewTLC51History()
	gap := event.TLC51ObligationReadyContent{TLC51EventIdentity: historyIdentity(2, 1), ObligationID: "O001", ReadyAt: now}
	if _, err := history.AppendTLC51(gap, now); err == nil {
		t.Fatal("ordinal gap unexpectedly appended")
	}
	plan := event.TLC51PlanRecordedContent{TLC51EventIdentity: historyIdentity(1, 0), Plan: historyPlan(), RecordedAt: now}
	if _, err := history.AppendTLC51(plan, now); err != nil {
		t.Fatalf("Append plan: %v", err)
	}
	changed := plan
	changed.RecordedAt = now.Add(time.Second)
	if _, err := history.AppendTLC51(changed, now.Add(time.Second)); err == nil {
		t.Fatal("conflicting ordinal unexpectedly appended")
	}
}

func TestTLC51WorkReconciliationRepairsOnlyFromValidEventGraph(t *testing.T) {
	now := time.Date(2026, 8, 27, 17, 0, 0, 0, time.UTC)
	history := NewTLC51History()
	entry, err := history.AppendTLC51(
		event.TLC51PlanRecordedContent{TLC51EventIdentity: historyIdentity(1, 0), Plan: historyPlan(), RecordedAt: now},
		now,
	)
	if err != nil {
		t.Fatalf("Append: %v", err)
	}
	repair, err := ReconcileTLC51Work(&entry, nil)
	if err != nil || repair.Action != TLC51ReconciliationRepairWork || repair.RepairedArtifact == nil || repair.HumanRequired {
		t.Fatalf("repair = %+v, %v", repair, err)
	}
	match, err := ReconcileTLC51Work(&entry, repair.RepairedArtifact)
	if err != nil || match.Action != TLC51ReconciliationMatch {
		t.Fatalf("match = %+v, %v", match, err)
	}
	conflictArtifact := *repair.RepairedArtifact
	conflictArtifact.Payload = append(json.RawMessage(nil), conflictArtifact.Payload...)
	conflictArtifact.Payload[0] ^= 1
	conflict, err := ReconcileTLC51Work(&entry, &conflictArtifact)
	if err != nil || conflict.Action != TLC51ReconciliationQuarantine || !conflict.HumanRequired {
		t.Fatalf("conflict = %+v, %v", conflict, err)
	}
	missingEvent, err := ReconcileTLC51Work(nil, repair.RepairedArtifact)
	if err != nil || missingEvent.Action != TLC51ReconciliationQuarantine || !missingEvent.HumanRequired {
		t.Fatalf("missing EventGraph = %+v, %v", missingEvent, err)
	}
}

func TestTLC51CutoverRejectsOnlyNewLegacyTransitions(t *testing.T) {
	effective := time.Date(2026, 8, 27, 18, 0, 0, 0, time.UTC)
	guard := TLC51CutoverGuard{
		EffectiveAt:                   effective,
		AcceptedReleaseIdentitySHA256: historyDigestA,
		AdapterIdentitySHA256:         historyDigestB,
		HiveBinarySHA256:              historyDigestA,
		HiveConfigurationSHA256:       historyDigestB,
		WorkerGroup:                   "tlc51-canary",
		ActivationReceiptDigest:       historyDigestA,
	}
	if err := guard.AdmitTransition(FactoryProtocolLegacy, effective, false); !errors.Is(err, ErrLegacyAfterCutover) {
		t.Fatalf("new legacy at cutover = %v", err)
	}
	if err := guard.AdmitTransition(FactoryProtocolLegacy, effective, true); err != nil {
		t.Fatalf("legacy replay rejected: %v", err)
	}
	if err := guard.AdmitTransition(FactoryProtocolLegacy, effective.Add(-time.Nanosecond), false); err != nil {
		t.Fatalf("pre-cutover legacy transition rejected: %v", err)
	}
	if err := guard.AdmitTransition(FactoryProtocolTLC51, effective, false); err != nil {
		t.Fatalf("TLC 5.1 transition rejected: %v", err)
	}
}
