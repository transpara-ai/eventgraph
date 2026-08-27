package event

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type tlc51Content struct{}

func (tlc51Content) Accept(EventContentVisitor) {}

// TLC51EventIdentity binds every event to one order/change-series/plan/subject
// history and a strict per-history ordinal. AttemptOrdinal is zero only for
// event kinds that are not attempts.
type TLC51EventIdentity struct {
	ProtocolVersion string `json:"protocol_version"`
	FactoryOrderID  string `json:"factory_order_id"`
	ChangeSeriesID  string `json:"change_series_id"`
	PlanDigest      string `json:"plan_digest"`
	SubjectDigest   string `json:"subject_digest"`
	EventOrdinal    uint64 `json:"event_ordinal"`
	AttemptOrdinal  uint32 `json:"attempt_ordinal"`
}

// TLC51ExactJSON preserves the exact canonical TLC artifact bytes. SHA256 is
// over JSON including its required terminal LF; the embedded schema identity
// is checked without rewriting the payload.
type TLC51ExactJSON struct {
	SchemaVersion string `json:"schema_version"`
	CanonicalJSON string `json:"canonical_json"`
	SHA256        string `json:"sha256"`
}

type TLC51ObligationOutcome string

const (
	TLC51ObligationPassed    TLC51ObligationOutcome = "passed"
	TLC51ObligationFailed    TLC51ObligationOutcome = "failed"
	TLC51ObligationBlocked   TLC51ObligationOutcome = "blocked"
	TLC51ObligationCancelled TLC51ObligationOutcome = "cancelled"
	TLC51ObligationUnknown   TLC51ObligationOutcome = "unknown"
)

func (v TLC51ObligationOutcome) IsValid() bool {
	switch v {
	case TLC51ObligationPassed, TLC51ObligationFailed, TLC51ObligationBlocked, TLC51ObligationCancelled, TLC51ObligationUnknown:
		return true
	default:
		return false
	}
}

type TLC51ExternalEffectState string

const (
	TLC51EffectAbsent   TLC51ExternalEffectState = "absent"
	TLC51EffectExact    TLC51ExternalEffectState = "exact"
	TLC51EffectConflict TLC51ExternalEffectState = "conflict"
	TLC51EffectUnknown  TLC51ExternalEffectState = "unknown"
)

func (v TLC51ExternalEffectState) IsValid() bool {
	switch v {
	case TLC51EffectAbsent, TLC51EffectExact, TLC51EffectConflict, TLC51EffectUnknown:
		return true
	default:
		return false
	}
}

type TLC51TerminalOutcome string

const (
	TLC51TerminalSucceeded TLC51TerminalOutcome = "succeeded"
	TLC51TerminalFailed    TLC51TerminalOutcome = "failed"
	TLC51TerminalBlocked   TLC51TerminalOutcome = "blocked"
	TLC51TerminalUnknown   TLC51TerminalOutcome = "unknown"
)

func (v TLC51TerminalOutcome) IsValid() bool {
	switch v {
	case TLC51TerminalSucceeded, TLC51TerminalFailed, TLC51TerminalBlocked, TLC51TerminalUnknown:
		return true
	default:
		return false
	}
}

type TLC51PlanRecordedContent struct {
	tlc51Content
	TLC51EventIdentity
	Plan       TLC51ExactJSON `json:"plan"`
	RecordedAt time.Time      `json:"recorded_at"`
}

func (TLC51PlanRecordedContent) EventTypeName() string { return EventTypeTLC51PlanRecorded.Value() }

type TLC51PlanSupersededContent struct {
	tlc51Content
	TLC51EventIdentity
	SupersedingPlanDigest string    `json:"superseding_plan_digest"`
	Reason                string    `json:"reason"`
	SupersededAt          time.Time `json:"superseded_at"`
}

func (TLC51PlanSupersededContent) EventTypeName() string { return EventTypeTLC51PlanSuperseded.Value() }

type TLC51ObligationReadyContent struct {
	tlc51Content
	TLC51EventIdentity
	ObligationID string    `json:"obligation_id"`
	ReadyAt      time.Time `json:"ready_at"`
}

func (TLC51ObligationReadyContent) EventTypeName() string {
	return EventTypeTLC51ObligationReady.Value()
}

type TLC51ObligationClaimedContent struct {
	tlc51Content
	TLC51EventIdentity
	ObligationID          string    `json:"obligation_id"`
	ProviderBindingID     string    `json:"provider_binding_id"`
	ProviderBindingSHA256 string    `json:"provider_binding_sha256"`
	ClaimedAt             time.Time `json:"claimed_at"`
}

func (TLC51ObligationClaimedContent) EventTypeName() string {
	return EventTypeTLC51ObligationClaimed.Value()
}

type TLC51ObligationRunningContent struct {
	tlc51Content
	TLC51EventIdentity
	ObligationID string    `json:"obligation_id"`
	InvocationID string    `json:"invocation_id"`
	RunningAt    time.Time `json:"running_at"`
}

func (TLC51ObligationRunningContent) EventTypeName() string {
	return EventTypeTLC51ObligationRunning.Value()
}

type TLC51ObligationTerminalContent struct {
	tlc51Content
	TLC51EventIdentity
	ObligationID string                 `json:"obligation_id"`
	Outcome      TLC51ObligationOutcome `json:"outcome"`
	Reason       string                 `json:"reason"`
	TerminalAt   time.Time              `json:"terminal_at"`
}

func (TLC51ObligationTerminalContent) EventTypeName() string {
	return EventTypeTLC51ObligationTerminal.Value()
}

type TLC51EvidenceLinkedContent struct {
	tlc51Content
	TLC51EventIdentity
	ObligationID     string         `json:"obligation_id"`
	EvidenceRecordID string         `json:"evidence_record_id"`
	EvidenceRecord   TLC51ExactJSON `json:"evidence_record"`
	LinkedAt         time.Time      `json:"linked_at"`
}

func (TLC51EvidenceLinkedContent) EventTypeName() string { return EventTypeTLC51EvidenceLinked.Value() }

type TLC51DecisionRecordedContent struct {
	tlc51Content
	TLC51EventIdentity
	Receipt    TLC51ExactJSON `json:"receipt"`
	Decision   string         `json:"decision"`
	RecordedAt time.Time      `json:"recorded_at"`
}

func (TLC51DecisionRecordedContent) EventTypeName() string {
	return EventTypeTLC51DecisionRecorded.Value()
}

type TLC51DecisionInvalidatedContent struct {
	tlc51Content
	TLC51EventIdentity
	ReceiptDigest string    `json:"receipt_digest"`
	Reason        string    `json:"reason"`
	InvalidatedAt time.Time `json:"invalidated_at"`
}

func (TLC51DecisionInvalidatedContent) EventTypeName() string {
	return EventTypeTLC51DecisionInvalidated.Value()
}

type TLC51ProtectedEffectProposedContent struct {
	tlc51Content
	TLC51EventIdentity
	Effect         string    `json:"effect"`
	OperationID    string    `json:"operation_id"`
	IdempotencyKey string    `json:"idempotency_key"`
	ReceiptDigest  string    `json:"receipt_digest"`
	ProposedAt     time.Time `json:"proposed_at"`
}

func (TLC51ProtectedEffectProposedContent) EventTypeName() string {
	return EventTypeTLC51ProtectedEffectProposed.Value()
}

type TLC51ProtectedEffectObservedContent struct {
	tlc51Content
	TLC51EventIdentity
	Effect                    string                   `json:"effect"`
	OperationID               string                   `json:"operation_id"`
	ExternalState             TLC51ExternalEffectState `json:"external_state"`
	ProviderObservationID     string                   `json:"provider_observation_id"`
	ProviderObservationDigest string                   `json:"provider_observation_digest"`
	ObservedAt                time.Time                `json:"observed_at"`
}

func (TLC51ProtectedEffectObservedContent) EventTypeName() string {
	return EventTypeTLC51ProtectedEffectObserved.Value()
}

type TLC51ProtectedEffectReconciledContent struct {
	tlc51Content
	TLC51EventIdentity
	Effect        string                   `json:"effect"`
	OperationID   string                   `json:"operation_id"`
	ExternalState TLC51ExternalEffectState `json:"external_state"`
	Action        string                   `json:"action"`
	ReconciledAt  time.Time                `json:"reconciled_at"`
}

func (TLC51ProtectedEffectReconciledContent) EventTypeName() string {
	return EventTypeTLC51ProtectedEffectReconciled.Value()
}

type TLC51ProtectedEffectTerminalContent struct {
	tlc51Content
	TLC51EventIdentity
	Effect        string               `json:"effect"`
	OperationID   string               `json:"operation_id"`
	Outcome       TLC51TerminalOutcome `json:"outcome"`
	EffectReceipt *TLC51ExactJSON      `json:"effect_receipt,omitempty"`
	Reason        string               `json:"reason"`
	TerminalAt    time.Time            `json:"terminal_at"`
}

func (TLC51ProtectedEffectTerminalContent) EventTypeName() string {
	return EventTypeTLC51ProtectedEffectTerminal.Value()
}

type TLC51HumanInterventionRequestedContent struct {
	tlc51Content
	TLC51EventIdentity
	RequestID   string    `json:"request_id"`
	Boundary    string    `json:"boundary"`
	Reason      string    `json:"reason"`
	RequestedAt time.Time `json:"requested_at"`
}

func (TLC51HumanInterventionRequestedContent) EventTypeName() string {
	return EventTypeTLC51HumanInterventionRequested.Value()
}

type TLC51HumanInterventionResolvedContent struct {
	tlc51Content
	TLC51EventIdentity
	RequestID   string         `json:"request_id"`
	HumanRecord TLC51ExactJSON `json:"human_record"`
	Resolution  string         `json:"resolution"`
	ResolvedAt  time.Time      `json:"resolved_at"`
}

func (TLC51HumanInterventionResolvedContent) EventTypeName() string {
	return EventTypeTLC51HumanInterventionResolved.Value()
}

type TLC51CutoverRecordedContent struct {
	tlc51Content
	TLC51EventIdentity
	AcceptedReleaseIdentitySHA256 string    `json:"accepted_release_identity_sha256"`
	AdapterIdentitySHA256         string    `json:"adapter_identity_sha256"`
	HiveBinarySHA256              string    `json:"hive_binary_sha256"`
	HiveConfigurationSHA256       string    `json:"hive_configuration_sha256"`
	WorkerGroup                   string    `json:"worker_group"`
	ActivationReceiptDigest       string    `json:"activation_receipt_digest"`
	EffectiveAt                   time.Time `json:"effective_at"`
}

func (TLC51CutoverRecordedContent) EventTypeName() string {
	return EventTypeTLC51CutoverRecorded.Value()
}

func validSHA256(value string) bool {
	if len(value) != sha256.Size*2 {
		return false
	}
	decoded, err := hex.DecodeString(value)
	return err == nil && hex.EncodeToString(decoded) == value
}

func validateTLC51Identity(value TLC51EventIdentity, attemptRequired bool) error {
	if value.ProtocolVersion != TLC51ProtocolVersion {
		return fmt.Errorf("protocol_version must be %q", TLC51ProtocolVersion)
	}
	if value.FactoryOrderID == "" || value.ChangeSeriesID == "" || value.EventOrdinal == 0 {
		return fmt.Errorf("factory_order_id, change_series_id, and event_ordinal are required")
	}
	if !validSHA256(value.PlanDigest) || !validSHA256(value.SubjectDigest) {
		return fmt.Errorf("plan_digest and subject_digest must be lowercase SHA-256")
	}
	if attemptRequired && value.AttemptOrdinal == 0 {
		return fmt.Errorf("attempt_ordinal must be positive for attempt events")
	}
	if !attemptRequired && value.AttemptOrdinal != 0 {
		return fmt.Errorf("attempt_ordinal must be zero for non-attempt events")
	}
	return nil
}

func validateTLC51ExactJSON(value TLC51ExactJSON, schema string) error {
	raw := []byte(value.CanonicalJSON)
	if value.SchemaVersion != schema || len(raw) == 0 || !json.Valid(raw) {
		return fmt.Errorf("exact JSON must carry schema %q and valid JSON", schema)
	}
	if raw[len(raw)-1] != '\n' || bytes.TrimSpace(raw)[0] != '{' {
		return fmt.Errorf("exact JSON must be an object with one canonical terminal LF")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var decoded any
	if err := decoder.Decode(&decoded); err != nil {
		return fmt.Errorf("exact JSON decode failed: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return fmt.Errorf("exact JSON must contain one value")
	}
	var canonical bytes.Buffer
	encoder := json.NewEncoder(&canonical)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(decoded); err != nil || !bytes.Equal(canonical.Bytes(), raw) {
		return fmt.Errorf("exact JSON is not TLC canonical JSON")
	}
	if !validSHA256(value.SHA256) || fmt.Sprintf("%x", sha256.Sum256(raw)) != value.SHA256 {
		return fmt.Errorf("exact JSON SHA-256 mismatch")
	}
	var identity struct {
		SchemaVersion string `json:"schema_version"`
		PlanDigest    string `json:"plan_digest"`
		ReceiptDigest string `json:"receipt_digest"`
		RecordDigest  string `json:"record_digest"`
	}
	if err := json.Unmarshal(raw, &identity); err != nil || identity.SchemaVersion != schema {
		return fmt.Errorf("embedded schema_version mismatch")
	}
	return nil
}

func validateTLC51Content(content EventContent) error {
	require := func(identity TLC51EventIdentity, attempt bool, fields ...string) error {
		if err := validateTLC51Identity(identity, attempt); err != nil {
			return err
		}
		for _, field := range fields {
			if field == "" {
				return fmt.Errorf("required TLC 5.1 identity field is empty")
			}
		}
		return nil
	}
	nonzero := func(value time.Time) error {
		if value.IsZero() || value.Location() != time.UTC {
			return fmt.Errorf("event observation time must be explicit UTC")
		}
		return nil
	}
	switch value := content.(type) {
	case TLC51PlanRecordedContent:
		if err := require(value.TLC51EventIdentity, false); err != nil {
			return err
		}
		if err := validateTLC51ExactJSON(value.Plan, "tlc-gate-plan/v1"); err != nil {
			return err
		}
		var plan struct {
			PlanDigest string `json:"plan_digest"`
		}
		if err := json.Unmarshal([]byte(value.Plan.CanonicalJSON), &plan); err != nil || plan.PlanDigest != value.PlanDigest {
			return fmt.Errorf("plan payload digest identity mismatch")
		}
		return nonzero(value.RecordedAt)
	case TLC51PlanSupersededContent:
		if err := require(value.TLC51EventIdentity, false, value.SupersedingPlanDigest, value.Reason); err != nil {
			return err
		}
		if !validSHA256(value.SupersedingPlanDigest) || value.SupersedingPlanDigest == value.PlanDigest {
			return fmt.Errorf("superseding plan digest must be a different SHA-256")
		}
		return nonzero(value.SupersededAt)
	case TLC51ObligationReadyContent:
		if err := require(value.TLC51EventIdentity, true, value.ObligationID); err != nil {
			return err
		}
		return nonzero(value.ReadyAt)
	case TLC51ObligationClaimedContent:
		if err := require(value.TLC51EventIdentity, true, value.ObligationID, value.ProviderBindingID); err != nil {
			return err
		}
		if !validSHA256(value.ProviderBindingSHA256) {
			return fmt.Errorf("provider binding digest required")
		}
		return nonzero(value.ClaimedAt)
	case TLC51ObligationRunningContent:
		if err := require(value.TLC51EventIdentity, true, value.ObligationID, value.InvocationID); err != nil {
			return err
		}
		return nonzero(value.RunningAt)
	case TLC51ObligationTerminalContent:
		if err := require(value.TLC51EventIdentity, true, value.ObligationID, value.Reason); err != nil {
			return err
		}
		if !value.Outcome.IsValid() {
			return fmt.Errorf("invalid obligation outcome")
		}
		return nonzero(value.TerminalAt)
	case TLC51EvidenceLinkedContent:
		if err := require(value.TLC51EventIdentity, true, value.ObligationID, value.EvidenceRecordID); err != nil {
			return err
		}
		if err := validateTLC51ExactJSON(value.EvidenceRecord, "tlc-gate-record/v1"); err != nil {
			return err
		}
		var record struct {
			RecordID string `json:"record_id"`
		}
		if err := json.Unmarshal([]byte(value.EvidenceRecord.CanonicalJSON), &record); err != nil || record.RecordID != value.EvidenceRecordID {
			return fmt.Errorf("evidence record identity mismatch")
		}
		return nonzero(value.LinkedAt)
	case TLC51DecisionRecordedContent:
		if err := require(value.TLC51EventIdentity, false, value.Decision); err != nil {
			return err
		}
		if value.Decision != "pass" && value.Decision != "fail" && value.Decision != "unknown" {
			return fmt.Errorf("invalid gate decision")
		}
		if err := validateTLC51ExactJSON(value.Receipt, "tlc-gate-receipt/v1"); err != nil {
			return err
		}
		var receipt struct {
			PlanDigest    string `json:"plan_digest"`
			SubjectDigest string `json:"subject_digest"`
			Decision      string `json:"decision"`
		}
		if err := json.Unmarshal([]byte(value.Receipt.CanonicalJSON), &receipt); err != nil || receipt.PlanDigest != value.PlanDigest || receipt.SubjectDigest != value.SubjectDigest || receipt.Decision != value.Decision {
			return fmt.Errorf("receipt event identity mismatch")
		}
		return nonzero(value.RecordedAt)
	case TLC51DecisionInvalidatedContent:
		if err := require(value.TLC51EventIdentity, false, value.ReceiptDigest, value.Reason); err != nil {
			return err
		}
		if !validSHA256(value.ReceiptDigest) {
			return fmt.Errorf("receipt digest required")
		}
		return nonzero(value.InvalidatedAt)
	case TLC51ProtectedEffectProposedContent:
		if err := require(value.TLC51EventIdentity, true, value.Effect, value.OperationID, value.IdempotencyKey, value.ReceiptDigest); err != nil {
			return err
		}
		if !validSHA256(value.ReceiptDigest) {
			return fmt.Errorf("receipt digest required")
		}
		return nonzero(value.ProposedAt)
	case TLC51ProtectedEffectObservedContent:
		if err := require(value.TLC51EventIdentity, true, value.Effect, value.OperationID, value.ProviderObservationID); err != nil {
			return err
		}
		if !value.ExternalState.IsValid() || !validSHA256(value.ProviderObservationDigest) {
			return fmt.Errorf("invalid effect observation")
		}
		return nonzero(value.ObservedAt)
	case TLC51ProtectedEffectReconciledContent:
		if err := require(value.TLC51EventIdentity, true, value.Effect, value.OperationID, value.Action); err != nil {
			return err
		}
		if !value.ExternalState.IsValid() || (value.Action != "settle" && value.Action != "retry" && value.Action != "block") {
			return fmt.Errorf("invalid reconciliation")
		}
		return nonzero(value.ReconciledAt)
	case TLC51ProtectedEffectTerminalContent:
		if err := require(value.TLC51EventIdentity, true, value.Effect, value.OperationID, value.Reason); err != nil {
			return err
		}
		if !value.Outcome.IsValid() {
			return fmt.Errorf("invalid effect outcome")
		}
		if value.Outcome == TLC51TerminalSucceeded && value.EffectReceipt == nil {
			return fmt.Errorf("successful effect requires exact effect receipt")
		}
		if value.EffectReceipt != nil {
			if err := validateTLC51ExactJSON(*value.EffectReceipt, "factory-tlc51-effect-receipt/v1"); err != nil {
				return err
			}
		}
		return nonzero(value.TerminalAt)
	case TLC51HumanInterventionRequestedContent:
		if err := require(value.TLC51EventIdentity, false, value.RequestID, value.Boundary, value.Reason); err != nil {
			return err
		}
		return nonzero(value.RequestedAt)
	case TLC51HumanInterventionResolvedContent:
		if err := require(value.TLC51EventIdentity, false, value.RequestID, value.Resolution); err != nil {
			return err
		}
		if err := validateTLC51ExactJSON(value.HumanRecord, "tlc-gate-record/v1"); err != nil {
			return err
		}
		return nonzero(value.ResolvedAt)
	case TLC51CutoverRecordedContent:
		if err := require(value.TLC51EventIdentity, false, value.WorkerGroup); err != nil {
			return err
		}
		for _, digest := range []string{value.AcceptedReleaseIdentitySHA256, value.AdapterIdentitySHA256, value.HiveBinarySHA256, value.HiveConfigurationSHA256, value.ActivationReceiptDigest} {
			if !validSHA256(digest) {
				return fmt.Errorf("cutover identity digest required")
			}
		}
		return nonzero(value.EffectiveAt)
	default:
		return fmt.Errorf("unexpected TLC 5.1 content type %T", content)
	}
}

// TLC51IdentityOf returns the immutable protocol identity embedded in a typed
// TLC 5.1 content value. The returned boolean is false for every other event
// family.
func TLC51IdentityOf(content EventContent) (TLC51EventIdentity, bool) {
	switch value := content.(type) {
	case TLC51PlanRecordedContent:
		return value.TLC51EventIdentity, true
	case TLC51PlanSupersededContent:
		return value.TLC51EventIdentity, true
	case TLC51ObligationReadyContent:
		return value.TLC51EventIdentity, true
	case TLC51ObligationClaimedContent:
		return value.TLC51EventIdentity, true
	case TLC51ObligationRunningContent:
		return value.TLC51EventIdentity, true
	case TLC51ObligationTerminalContent:
		return value.TLC51EventIdentity, true
	case TLC51EvidenceLinkedContent:
		return value.TLC51EventIdentity, true
	case TLC51DecisionRecordedContent:
		return value.TLC51EventIdentity, true
	case TLC51DecisionInvalidatedContent:
		return value.TLC51EventIdentity, true
	case TLC51ProtectedEffectProposedContent:
		return value.TLC51EventIdentity, true
	case TLC51ProtectedEffectObservedContent:
		return value.TLC51EventIdentity, true
	case TLC51ProtectedEffectReconciledContent:
		return value.TLC51EventIdentity, true
	case TLC51ProtectedEffectTerminalContent:
		return value.TLC51EventIdentity, true
	case TLC51HumanInterventionRequestedContent:
		return value.TLC51EventIdentity, true
	case TLC51HumanInterventionResolvedContent:
		return value.TLC51EventIdentity, true
	case TLC51CutoverRecordedContent:
		return value.TLC51EventIdentity, true
	default:
		return TLC51EventIdentity{}, false
	}
}
