package event

import (
	"encoding/json"
	"fmt"
)

type nativeEvidenceContent struct{}

func (nativeEvidenceContent) Accept(EventContentVisitor) {}

type NativeEvidenceOutcome string

const (
	NativeEvidenceOutcomePassed  NativeEvidenceOutcome = "passed"
	NativeEvidenceOutcomeFailed  NativeEvidenceOutcome = "failed"
	NativeEvidenceOutcomeBlocked NativeEvidenceOutcome = "blocked"
	NativeEvidenceOutcomeSkipped NativeEvidenceOutcome = "skipped"
	NativeEvidenceOutcomeErrored NativeEvidenceOutcome = "errored"
	NativeEvidenceOutcomeWaived  NativeEvidenceOutcome = "waived"
	NativeEvidenceOutcomePartial NativeEvidenceOutcome = "partial"
)

func (o NativeEvidenceOutcome) IsValid() bool {
	switch o {
	case NativeEvidenceOutcomePassed, NativeEvidenceOutcomeFailed, NativeEvidenceOutcomeBlocked,
		NativeEvidenceOutcomeSkipped, NativeEvidenceOutcomeErrored, NativeEvidenceOutcomeWaived,
		NativeEvidenceOutcomePartial:
		return true
	default:
		return false
	}
}

func (o *NativeEvidenceOutcome) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	outcome := NativeEvidenceOutcome(value)
	if !outcome.IsValid() {
		return fmt.Errorf("invalid native evidence outcome %q", value)
	}
	*o = outcome
	return nil
}

type NativeEvidenceIssueRef struct {
	Repo   string `json:"repo"`
	Number int    `json:"number"`
	URL    string `json:"url,omitempty"`
	Title  string `json:"title,omitempty"`
}

type NativeTestRunRecordedContent struct {
	nativeEvidenceContent
	TestRunID         string                   `json:"test_run_id"`
	TestCaseID        string                   `json:"test_case_id,omitempty"`
	ActorInvocationID string                   `json:"actor_invocation_id,omitempty"`
	Command           string                   `json:"command"`
	Outcome           NativeEvidenceOutcome    `json:"outcome"`
	StartedAt         string                   `json:"started_at,omitempty"`
	CompletedAt       string                   `json:"completed_at,omitempty"`
	SourceIssueRefs   []NativeEvidenceIssueRef `json:"source_issue_refs,omitempty"`
	PRRefs            []string                 `json:"pr_refs,omitempty"`
	EvidenceRefs      []string                 `json:"evidence_refs,omitempty"`
	ValidationRefs    []string                 `json:"validation_refs,omitempty"`
	SourceRefs        []string                 `json:"source_refs,omitempty"`
}

func (c *NativeTestRunRecordedContent) UnmarshalJSON(data []byte) error {
	if err := requireNativeEvidenceJSONField(data, "outcome"); err != nil {
		return err
	}
	type alias NativeTestRunRecordedContent
	var out alias
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	if err := validateNativeTestRunRecorded(NativeTestRunRecordedContent(out)); err != nil {
		return err
	}
	*c = NativeTestRunRecordedContent(out)
	return nil
}

func (c NativeTestRunRecordedContent) EventTypeName() string {
	return EventTypeNativeTestRunRecorded.Value()
}

type NativeGateResultRecordedContent struct {
	nativeEvidenceContent
	GateResultID       string                   `json:"gate_result_id"`
	FactoryOrderID     string                   `json:"factory_order_id"`
	ReleaseCandidateID string                   `json:"release_candidate_id,omitempty"`
	GateName           string                   `json:"gate_name"`
	Outcome            NativeEvidenceOutcome    `json:"outcome"`
	EvidenceRefs       []string                 `json:"evidence_refs,omitempty"`
	WaiverRef          string                   `json:"waiver_ref,omitempty"`
	SourceIssueRefs    []NativeEvidenceIssueRef `json:"source_issue_refs,omitempty"`
	PRRefs             []string                 `json:"pr_refs,omitempty"`
	ValidationRefs     []string                 `json:"validation_refs,omitempty"`
	SourceRefs         []string                 `json:"source_refs,omitempty"`
}

func (c *NativeGateResultRecordedContent) UnmarshalJSON(data []byte) error {
	if err := requireNativeEvidenceJSONField(data, "outcome"); err != nil {
		return err
	}
	type alias NativeGateResultRecordedContent
	var out alias
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	if err := validateNativeGateResultRecorded(NativeGateResultRecordedContent(out)); err != nil {
		return err
	}
	*c = NativeGateResultRecordedContent(out)
	return nil
}

func (c NativeGateResultRecordedContent) EventTypeName() string {
	return EventTypeNativeGateResultRecorded.Value()
}

type NativeAuditReportRecordedContent struct {
	nativeEvidenceContent
	AuditReportID         string                   `json:"audit_report_id"`
	TargetType            string                   `json:"target_type"`
	TargetID              string                   `json:"target_id"`
	Outcome               NativeEvidenceOutcome    `json:"outcome"`
	MissingLinks          []string                 `json:"missing_links,omitempty"`
	TraceScoreBasisPoints *int                     `json:"trace_score_basis_points"`
	SourceIssueRefs       []NativeEvidenceIssueRef `json:"source_issue_refs,omitempty"`
	PRRefs                []string                 `json:"pr_refs,omitempty"`
	ValidationRefs        []string                 `json:"validation_refs,omitempty"`
	CFARRefs              []string                 `json:"cfar_refs,omitempty"`
	AuthorityBoundaryRefs []string                 `json:"authority_boundary_refs,omitempty"`
	ResidualRiskRefs      []string                 `json:"residual_risk_refs,omitempty"`
	EvidenceRefs          []string                 `json:"evidence_refs,omitempty"`
	SourceRefs            []string                 `json:"source_refs,omitempty"`
}

func (c *NativeAuditReportRecordedContent) UnmarshalJSON(data []byte) error {
	if err := requireNativeEvidenceJSONField(data, "outcome"); err != nil {
		return err
	}
	if err := requireNativeEvidenceJSONField(data, "trace_score_basis_points"); err != nil {
		return err
	}
	type alias NativeAuditReportRecordedContent
	var out alias
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	if err := validateNativeAuditReportRecorded(NativeAuditReportRecordedContent(out)); err != nil {
		return err
	}
	*c = NativeAuditReportRecordedContent(out)
	return nil
}

func (c NativeAuditReportRecordedContent) EventTypeName() string {
	return EventTypeNativeAuditReportRecorded.Value()
}

func validateNativeEvidenceContent(content EventContent) error {
	switch c := content.(type) {
	case NativeTestRunRecordedContent:
		return validateNativeTestRunRecorded(c)
	case NativeGateResultRecordedContent:
		return validateNativeGateResultRecorded(c)
	case NativeAuditReportRecordedContent:
		return validateNativeAuditReportRecorded(c)
	default:
		return fmt.Errorf("unexpected native evidence content %T", content)
	}
}

func validateNativeTestRunRecorded(c NativeTestRunRecordedContent) error {
	if err := requireNativeEvidenceString("test_run_id", c.TestRunID); err != nil {
		return err
	}
	if err := requireNativeEvidenceString("command", c.Command); err != nil {
		return err
	}
	if !c.Outcome.IsValid() {
		return fmt.Errorf("invalid native test-run outcome %q", c.Outcome)
	}
	return validateNativeEvidenceIssueRefs(c.SourceIssueRefs)
}

func validateNativeGateResultRecorded(c NativeGateResultRecordedContent) error {
	if err := requireNativeEvidenceString("gate_result_id", c.GateResultID); err != nil {
		return err
	}
	if err := requireNativeEvidenceString("factory_order_id", c.FactoryOrderID); err != nil {
		return err
	}
	if err := requireNativeEvidenceString("gate_name", c.GateName); err != nil {
		return err
	}
	if !c.Outcome.IsValid() {
		return fmt.Errorf("invalid native gate-result outcome %q", c.Outcome)
	}
	return validateNativeEvidenceIssueRefs(c.SourceIssueRefs)
}

func validateNativeAuditReportRecorded(c NativeAuditReportRecordedContent) error {
	if err := requireNativeEvidenceString("audit_report_id", c.AuditReportID); err != nil {
		return err
	}
	if err := requireNativeEvidenceString("target_type", c.TargetType); err != nil {
		return err
	}
	if err := requireNativeEvidenceString("target_id", c.TargetID); err != nil {
		return err
	}
	if !c.Outcome.IsValid() {
		return fmt.Errorf("invalid native audit-report outcome %q", c.Outcome)
	}
	if err := validateTraceScoreBasisPoints(c.TraceScoreBasisPoints); err != nil {
		return err
	}
	return validateNativeEvidenceIssueRefs(c.SourceIssueRefs)
}

func requireNativeEvidenceString(field string, value string) error {
	if value == "" {
		return fmt.Errorf("native evidence field %q is required", field)
	}
	return nil
}

func validateNativeEvidenceIssueRefs(refs []NativeEvidenceIssueRef) error {
	for i, ref := range refs {
		if ref.Repo == "" {
			return fmt.Errorf("source_issue_refs[%d].repo is required", i)
		}
		if ref.Number <= 0 {
			return fmt.Errorf("source_issue_refs[%d].number must be positive", i)
		}
	}
	return nil
}

func validateTraceScoreBasisPoints(score *int) error {
	if score == nil {
		return fmt.Errorf("trace_score_basis_points is required")
	}
	if *score < 0 || *score > 10000 {
		return fmt.Errorf("trace_score_basis_points must be between 0 and 10000, got %d", *score)
	}
	return nil
}

func requireNativeEvidenceJSONField(data []byte, field string) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if _, ok := fields[field]; !ok {
		return fmt.Errorf("native evidence field %q is required", field)
	}
	return nil
}
