package v39

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

var (
	ErrTLC51HistoryConflict = errors.New("factory-tlc51/v1 immutable history conflict")
	ErrTLC51HistoryGap      = errors.New("factory-tlc51/v1 event ordinal gap")
	ErrLegacyAfterCutover   = errors.New("new tlc-v1 transition rejected after TLC 5.1 cutover")
)

const (
	FactoryProtocolLegacy = "tlc-v1"
	FactoryProtocolTLC51  = event.TLC51ProtocolVersion
)

// TLC51HistoryEntry is an append-only, hash-chained projection of a typed
// EventGraph content payload. Payload is preserved verbatim for deterministic
// Work-store repair; it never conveys authority.
type TLC51HistoryEntry struct {
	FactoryOrderID string          `json:"factory_order_id"`
	ChangeSeriesID string          `json:"change_series_id"`
	EventType      types.EventType `json:"event_type"`
	EventOrdinal   uint64          `json:"event_ordinal"`
	PlanDigest     string          `json:"plan_digest"`
	SubjectDigest  string          `json:"subject_digest"`
	AttemptOrdinal uint32          `json:"attempt_ordinal"`
	Payload        json.RawMessage `json:"payload"`
	PayloadSHA256  string          `json:"payload_sha256"`
	PreviousDigest string          `json:"previous_digest"`
	EntryDigest    string          `json:"entry_digest"`
	RecordedAt     time.Time       `json:"recorded_at"`
}

type TLC51History struct {
	mu      sync.RWMutex
	entries map[string][]TLC51HistoryEntry
}

func NewTLC51History() *TLC51History {
	return &TLC51History{entries: map[string][]TLC51HistoryEntry{}}
}

func historyKey(factoryOrderID, changeSeriesID string) string {
	return factoryOrderID + "\x00" + changeSeriesID
}

func cloneHistoryEntry(value TLC51HistoryEntry) TLC51HistoryEntry {
	value.Payload = append(json.RawMessage(nil), value.Payload...)
	return value
}

func digestHistoryEntry(value TLC51HistoryEntry) (string, error) {
	value.EntryDigest = ""
	payload, err := CanonicalJSON(value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(payload)), nil
}

// AppendTLC51 validates a typed content value through EventGraph's registry,
// serializes it as data, and appends exactly one ordinal to its history.
func (h *TLC51History) AppendTLC51(content event.EventContent, recordedAt time.Time) (TLC51HistoryEntry, error) {
	identity, ok := event.TLC51IdentityOf(content)
	if !ok {
		return TLC51HistoryEntry{}, fmt.Errorf("%w: content %T is not TLC 5.1", ErrInvalidRecord, content)
	}
	eventType := types.MustEventType(content.EventTypeName())
	if err := event.DefaultRegistry().Validate(eventType, content); err != nil {
		return TLC51HistoryEntry{}, fmt.Errorf("%w: %v", ErrInvalidRecord, err)
	}
	if recordedAt.IsZero() || recordedAt.Location() != time.UTC {
		return TLC51HistoryEntry{}, fmt.Errorf("%w: recorded_at must be explicit UTC", ErrInvalidRecord)
	}
	payload, err := json.Marshal(content)
	if err != nil {
		return TLC51HistoryEntry{}, err
	}
	payloadSHA := fmt.Sprintf("%x", sha256.Sum256(payload))
	key := historyKey(identity.FactoryOrderID, identity.ChangeSeriesID)

	h.mu.Lock()
	defer h.mu.Unlock()
	rows := h.entries[key]
	if identity.EventOrdinal <= uint64(len(rows)) {
		existing := rows[identity.EventOrdinal-1]
		if existing.EventType == eventType && existing.PayloadSHA256 == payloadSHA && bytes.Equal(existing.Payload, payload) {
			return cloneHistoryEntry(existing), nil
		}
		return TLC51HistoryEntry{}, fmt.Errorf("%w: %s ordinal %d", ErrTLC51HistoryConflict, key, identity.EventOrdinal)
	}
	if identity.EventOrdinal != uint64(len(rows)+1) {
		return TLC51HistoryEntry{}, fmt.Errorf("%w: got %d want %d", ErrTLC51HistoryGap, identity.EventOrdinal, len(rows)+1)
	}
	previous := ""
	if len(rows) > 0 {
		previous = rows[len(rows)-1].EntryDigest
	}
	entry := TLC51HistoryEntry{
		FactoryOrderID: identity.FactoryOrderID,
		ChangeSeriesID: identity.ChangeSeriesID,
		EventType:      eventType,
		EventOrdinal:   identity.EventOrdinal,
		PlanDigest:     identity.PlanDigest,
		SubjectDigest:  identity.SubjectDigest,
		AttemptOrdinal: identity.AttemptOrdinal,
		Payload:        append(json.RawMessage(nil), payload...),
		PayloadSHA256:  payloadSHA,
		PreviousDigest: previous,
		RecordedAt:     recordedAt,
	}
	entry.EntryDigest, err = digestHistoryEntry(entry)
	if err != nil {
		return TLC51HistoryEntry{}, err
	}
	h.entries[key] = append(rows, entry)
	return cloneHistoryEntry(entry), nil
}

func (h *TLC51History) Entries(factoryOrderID, changeSeriesID string) []TLC51HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	rows := h.entries[historyKey(factoryOrderID, changeSeriesID)]
	result := make([]TLC51HistoryEntry, len(rows))
	for index, row := range rows {
		result[index] = cloneHistoryEntry(row)
	}
	return result
}

func (h *TLC51History) Verify(factoryOrderID, changeSeriesID string) error {
	rows := h.Entries(factoryOrderID, changeSeriesID)
	previous := ""
	for index, row := range rows {
		if row.EventOrdinal != uint64(index+1) || row.PreviousDigest != previous {
			return fmt.Errorf("%w: chain position %d", ErrTLC51HistoryConflict, index+1)
		}
		if !validHistorySHA(row.PayloadSHA256) || fmt.Sprintf("%x", sha256.Sum256(row.Payload)) != row.PayloadSHA256 {
			return fmt.Errorf("%w: payload digest %d", ErrTLC51HistoryConflict, index+1)
		}
		digest, err := digestHistoryEntry(row)
		if err != nil || digest != row.EntryDigest {
			return fmt.Errorf("%w: entry digest %d", ErrTLC51HistoryConflict, index+1)
		}
		previous = row.EntryDigest
	}
	return nil
}

func validHistorySHA(value string) bool {
	if len(value) != sha256.Size*2 {
		return false
	}
	decoded, err := hex.DecodeString(value)
	return err == nil && hex.EncodeToString(decoded) == value
}

type TLC51WorkArtifact struct {
	FactoryOrderID string          `json:"factory_order_id"`
	ChangeSeriesID string          `json:"change_series_id"`
	EventOrdinal   uint64          `json:"event_ordinal"`
	EventType      types.EventType `json:"event_type"`
	Payload        json.RawMessage `json:"payload"`
	PayloadSHA256  string          `json:"payload_sha256"`
}

type TLC51ReconciliationAction string

const (
	TLC51ReconciliationMatch       TLC51ReconciliationAction = "match"
	TLC51ReconciliationRepairWork  TLC51ReconciliationAction = "repair_work_from_eventgraph"
	TLC51ReconciliationQuarantine  TLC51ReconciliationAction = "quarantine_conflict_human_required"
	TLC51ReconciliationMissingBoth TLC51ReconciliationAction = "missing_both_human_required"
)

type TLC51Reconciliation struct {
	Action            TLC51ReconciliationAction `json:"action"`
	EventEntryDigest  *string                   `json:"event_entry_digest,omitempty"`
	WorkPayloadSHA256 *string                   `json:"work_payload_sha256,omitempty"`
	RepairedArtifact  *TLC51WorkArtifact        `json:"repaired_artifact,omitempty"`
	HumanRequired     bool                      `json:"human_required"`
}

// ReconcileTLC51Work repairs only a missing Work twin from valid EventGraph
// history. A conflict or missing EventGraph source is quarantined; Work bytes
// never overwrite or invent EventGraph history.
func ReconcileTLC51Work(entry *TLC51HistoryEntry, work *TLC51WorkArtifact) (TLC51Reconciliation, error) {
	if entry == nil && work == nil {
		return TLC51Reconciliation{Action: TLC51ReconciliationMissingBoth, HumanRequired: true}, nil
	}
	if entry == nil {
		workDigest := work.PayloadSHA256
		return TLC51Reconciliation{Action: TLC51ReconciliationQuarantine, WorkPayloadSHA256: &workDigest, HumanRequired: true}, nil
	}
	if !validHistorySHA(entry.PayloadSHA256) || fmt.Sprintf("%x", sha256.Sum256(entry.Payload)) != entry.PayloadSHA256 {
		return TLC51Reconciliation{}, fmt.Errorf("%w: invalid EventGraph payload", ErrTLC51HistoryConflict)
	}
	entryDigest := entry.EntryDigest
	if work == nil {
		artifact := &TLC51WorkArtifact{
			FactoryOrderID: entry.FactoryOrderID,
			ChangeSeriesID: entry.ChangeSeriesID,
			EventOrdinal:   entry.EventOrdinal,
			EventType:      entry.EventType,
			Payload:        append(json.RawMessage(nil), entry.Payload...),
			PayloadSHA256:  entry.PayloadSHA256,
		}
		return TLC51Reconciliation{Action: TLC51ReconciliationRepairWork, EventEntryDigest: &entryDigest, RepairedArtifact: artifact}, nil
	}
	workDigest := work.PayloadSHA256
	identityMatches := work.FactoryOrderID == entry.FactoryOrderID && work.ChangeSeriesID == entry.ChangeSeriesID && work.EventOrdinal == entry.EventOrdinal && work.EventType == entry.EventType
	payloadMatches := validHistorySHA(work.PayloadSHA256) && fmt.Sprintf("%x", sha256.Sum256(work.Payload)) == work.PayloadSHA256 && work.PayloadSHA256 == entry.PayloadSHA256 && bytes.Equal(work.Payload, entry.Payload)
	if identityMatches && payloadMatches {
		return TLC51Reconciliation{Action: TLC51ReconciliationMatch, EventEntryDigest: &entryDigest, WorkPayloadSHA256: &workDigest}, nil
	}
	return TLC51Reconciliation{Action: TLC51ReconciliationQuarantine, EventEntryDigest: &entryDigest, WorkPayloadSHA256: &workDigest, HumanRequired: true}, nil
}

type TLC51CutoverGuard struct {
	EffectiveAt                   time.Time `json:"effective_at"`
	AcceptedReleaseIdentitySHA256 string    `json:"accepted_release_identity_sha256"`
	AdapterIdentitySHA256         string    `json:"adapter_identity_sha256"`
	HiveBinarySHA256              string    `json:"hive_binary_sha256"`
	HiveConfigurationSHA256       string    `json:"hive_configuration_sha256"`
	WorkerGroup                   string    `json:"worker_group"`
	ActivationReceiptDigest       string    `json:"activation_receipt_digest"`
}

func (g TLC51CutoverGuard) Validate() error {
	if g.EffectiveAt.IsZero() || g.EffectiveAt.Location() != time.UTC || g.WorkerGroup == "" {
		return fmt.Errorf("%w: cutover time and worker group required", ErrInvalidRecord)
	}
	for _, value := range []string{g.AcceptedReleaseIdentitySHA256, g.AdapterIdentitySHA256, g.HiveBinarySHA256, g.HiveConfigurationSHA256, g.ActivationReceiptDigest} {
		if !validHistorySHA(value) {
			return fmt.Errorf("%w: cutover digest required", ErrInvalidRecord)
		}
	}
	return nil
}

// AdmitTransition keeps immutable legacy replay readable while rejecting new
// tlc-v1 transition emission at or after the separately authorized cutover.
func (g TLC51CutoverGuard) AdmitTransition(protocol string, transitionAt time.Time, replay bool) error {
	if err := g.Validate(); err != nil {
		return err
	}
	if transitionAt.IsZero() || transitionAt.Location() != time.UTC {
		return fmt.Errorf("%w: transition time must be explicit UTC", ErrInvalidRecord)
	}
	if protocol != FactoryProtocolLegacy && protocol != FactoryProtocolTLC51 {
		return fmt.Errorf("%w: unknown Factory protocol", ErrInvalidRecord)
	}
	if protocol == FactoryProtocolLegacy && !replay && !transitionAt.Before(g.EffectiveAt) {
		return ErrLegacyAfterCutover
	}
	return nil
}

// SortTLC51History gives deterministic projection order without changing the
// append-only source history.
func SortTLC51History(rows []TLC51HistoryEntry) []TLC51HistoryEntry {
	result := make([]TLC51HistoryEntry, len(rows))
	for index, row := range rows {
		result[index] = cloneHistoryEntry(row)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].FactoryOrderID != result[j].FactoryOrderID {
			return result[i].FactoryOrderID < result[j].FactoryOrderID
		}
		if result[i].ChangeSeriesID != result[j].ChangeSeriesID {
			return result[i].ChangeSeriesID < result[j].ChangeSeriesID
		}
		return result[i].EventOrdinal < result[j].EventOrdinal
	})
	return result
}
