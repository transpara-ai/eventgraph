package event

import (
	"encoding/json"
	"fmt"
)

type authorityEvidenceContent struct{}

func (authorityEvidenceContent) Accept(EventContentVisitor) {}

type AuthorityDecisionOutcome string

const (
	AuthorityDecisionOutcomeAutonomous       AuthorityDecisionOutcome = "autonomous"
	AuthorityDecisionOutcomeNotify           AuthorityDecisionOutcome = "notify"
	AuthorityDecisionOutcomeApprovalRequired AuthorityDecisionOutcome = "approval_required"
	AuthorityDecisionOutcomeForbidden        AuthorityDecisionOutcome = "forbidden"
)

func (o AuthorityDecisionOutcome) IsValid() bool {
	switch o {
	case AuthorityDecisionOutcomeAutonomous, AuthorityDecisionOutcomeNotify,
		AuthorityDecisionOutcomeApprovalRequired, AuthorityDecisionOutcomeForbidden:
		return true
	default:
		return false
	}
}

func (o *AuthorityDecisionOutcome) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	outcome := AuthorityDecisionOutcome(value)
	if !outcome.IsValid() {
		return fmt.Errorf("invalid authority decision outcome %q", value)
	}
	*o = outcome
	return nil
}

type AuthorityBoundaryType string

const (
	AuthorityBoundaryProtectedAction     AuthorityBoundaryType = "protected_action"
	AuthorityBoundaryHumanScope          AuthorityBoundaryType = "human_scope"
	AuthorityBoundaryRuntimeExecution    AuthorityBoundaryType = "runtime_execution"
	AuthorityBoundaryEventGraphWrite     AuthorityBoundaryType = "eventgraph_write"
	AuthorityBoundaryMerge               AuthorityBoundaryType = "merge"
	AuthorityBoundaryDeployment          AuthorityBoundaryType = "deployment"
	AuthorityBoundaryAutonomyIncrease    AuthorityBoundaryType = "autonomy_increase"
	AuthorityBoundaryValueAllocation     AuthorityBoundaryType = "value_allocation"
	AuthorityBoundaryResidualRiskClosure AuthorityBoundaryType = "residual_risk_closure"
	AuthorityBoundaryProtectedSettings   AuthorityBoundaryType = "protected_settings"
)

func (t AuthorityBoundaryType) IsValid() bool {
	switch t {
	case AuthorityBoundaryProtectedAction, AuthorityBoundaryHumanScope,
		AuthorityBoundaryRuntimeExecution, AuthorityBoundaryEventGraphWrite,
		AuthorityBoundaryMerge, AuthorityBoundaryDeployment, AuthorityBoundaryAutonomyIncrease,
		AuthorityBoundaryValueAllocation, AuthorityBoundaryResidualRiskClosure,
		AuthorityBoundaryProtectedSettings:
		return true
	default:
		return false
	}
}

func (t *AuthorityBoundaryType) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	boundaryType := AuthorityBoundaryType(value)
	if !boundaryType.IsValid() {
		return fmt.Errorf("invalid authority boundary type %q", value)
	}
	*t = boundaryType
	return nil
}

type AuthorityBoundaryState string

const (
	AuthorityBoundaryStateBlocked          AuthorityBoundaryState = "blocked"
	AuthorityBoundaryStateApprovalRequired AuthorityBoundaryState = "approval_required"
	AuthorityBoundaryStateAuthorized       AuthorityBoundaryState = "authorized"
	AuthorityBoundaryStateNotApplicable    AuthorityBoundaryState = "not_applicable"
)

func (s AuthorityBoundaryState) IsValid() bool {
	switch s {
	case AuthorityBoundaryStateBlocked, AuthorityBoundaryStateApprovalRequired,
		AuthorityBoundaryStateAuthorized, AuthorityBoundaryStateNotApplicable:
		return true
	default:
		return false
	}
}

func (s *AuthorityBoundaryState) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	state := AuthorityBoundaryState(value)
	if !state.IsValid() {
		return fmt.Errorf("invalid authority boundary state %q", value)
	}
	*s = state
	return nil
}

type AuthorityResidualStatus string

const (
	AuthorityResidualStatusOpen          AuthorityResidualStatus = "open"
	AuthorityResidualStatusCarried       AuthorityResidualStatus = "carried"
	AuthorityResidualStatusWaived        AuthorityResidualStatus = "waived"
	AuthorityResidualStatusClosed        AuthorityResidualStatus = "closed"
	AuthorityResidualStatusNotApplicable AuthorityResidualStatus = "not_applicable"
)

func (s AuthorityResidualStatus) IsValid() bool {
	switch s {
	case AuthorityResidualStatusOpen, AuthorityResidualStatusCarried,
		AuthorityResidualStatusWaived, AuthorityResidualStatusClosed,
		AuthorityResidualStatusNotApplicable:
		return true
	default:
		return false
	}
}

func (s *AuthorityResidualStatus) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	status := AuthorityResidualStatus(value)
	if !status.IsValid() {
		return fmt.Errorf("invalid authority residual status %q", value)
	}
	*s = status
	return nil
}

type AuthorityResidualSeverity string

const (
	AuthorityResidualSeverityLow      AuthorityResidualSeverity = "low"
	AuthorityResidualSeverityMedium   AuthorityResidualSeverity = "medium"
	AuthorityResidualSeverityHigh     AuthorityResidualSeverity = "high"
	AuthorityResidualSeverityCritical AuthorityResidualSeverity = "critical"
)

func (s AuthorityResidualSeverity) IsValid() bool {
	switch s {
	case AuthorityResidualSeverityLow, AuthorityResidualSeverityMedium,
		AuthorityResidualSeverityHigh, AuthorityResidualSeverityCritical:
		return true
	default:
		return false
	}
}

func (s *AuthorityResidualSeverity) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	severity := AuthorityResidualSeverity(value)
	if !severity.IsValid() {
		return fmt.Errorf("invalid authority residual severity %q", value)
	}
	*s = severity
	return nil
}

type AuthorityStoreWriteStatus string

const (
	AuthorityStoreWriteStatusSchemaOnly          AuthorityStoreWriteStatus = "schema_only"
	AuthorityStoreWriteStatusProjectionOnly      AuthorityStoreWriteStatus = "projection_only"
	AuthorityStoreWriteStatusMigrationRequired   AuthorityStoreWriteStatus = "migration_required"
	AuthorityStoreWriteStatusWritePathBlocked    AuthorityStoreWriteStatus = "write_path_blocked"
	AuthorityStoreWriteStatusWritePathAuthorized AuthorityStoreWriteStatus = "write_path_authorized"
)

func (s AuthorityStoreWriteStatus) IsValid() bool {
	switch s {
	case AuthorityStoreWriteStatusSchemaOnly, AuthorityStoreWriteStatusProjectionOnly,
		AuthorityStoreWriteStatusMigrationRequired, AuthorityStoreWriteStatusWritePathBlocked,
		AuthorityStoreWriteStatusWritePathAuthorized:
		return true
	default:
		return false
	}
}

func (s *AuthorityStoreWriteStatus) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	status := AuthorityStoreWriteStatus(value)
	if !status.IsValid() {
		return fmt.Errorf("invalid authority store write status %q", value)
	}
	*s = status
	return nil
}

type AuthorityDecisionRecordedContent struct {
	authorityEvidenceContent
	DecisionID       string                   `json:"decision_id"`
	SubjectType      string                   `json:"subject_type"`
	SubjectRef       string                   `json:"subject_ref"`
	Outcome          AuthorityDecisionOutcome `json:"outcome"`
	ActorRef         string                   `json:"actor_ref,omitempty"`
	ProtectedActions []ProtectedAction        `json:"protected_actions,omitempty"`
	SourceIssueRefs  []NativeEvidenceIssueRef `json:"source_issue_refs,omitempty"`
	AuthorityRefs    []string                 `json:"authority_refs,omitempty"`
	EvidenceRefs     []string                 `json:"evidence_refs,omitempty"`
	NonClaimRefs     []string                 `json:"non_claim_refs,omitempty"`
	RecordedAt       string                   `json:"recorded_at,omitempty"`
}

func (c *AuthorityDecisionRecordedContent) UnmarshalJSON(data []byte) error {
	if err := requireAuthorityEvidenceJSONField(data, "outcome"); err != nil {
		return err
	}
	type alias AuthorityDecisionRecordedContent
	var out alias
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	if err := validateAuthorityDecisionRecorded(AuthorityDecisionRecordedContent(out)); err != nil {
		return err
	}
	*c = AuthorityDecisionRecordedContent(out)
	return nil
}

func (c AuthorityDecisionRecordedContent) EventTypeName() string {
	return EventTypeAuthorityDecisionRecorded.Value()
}

type AuthorityBoundaryRecordedContent struct {
	authorityEvidenceContent
	BoundaryID            string                   `json:"boundary_id"`
	BoundaryType          AuthorityBoundaryType    `json:"boundary_type"`
	SubjectRef            string                   `json:"subject_ref"`
	State                 AuthorityBoundaryState   `json:"state"`
	ProtectedActions      []ProtectedAction        `json:"protected_actions,omitempty"`
	RequiredAuthorityRefs []string                 `json:"required_authority_refs,omitempty"`
	SourceIssueRefs       []NativeEvidenceIssueRef `json:"source_issue_refs,omitempty"`
	EvidenceRefs          []string                 `json:"evidence_refs,omitempty"`
	RecordedAt            string                   `json:"recorded_at,omitempty"`
}

func (c *AuthorityBoundaryRecordedContent) UnmarshalJSON(data []byte) error {
	if err := requireAuthorityEvidenceJSONField(data, "boundary_type"); err != nil {
		return err
	}
	if err := requireAuthorityEvidenceJSONField(data, "state"); err != nil {
		return err
	}
	type alias AuthorityBoundaryRecordedContent
	var out alias
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	if err := validateAuthorityBoundaryRecorded(AuthorityBoundaryRecordedContent(out)); err != nil {
		return err
	}
	*c = AuthorityBoundaryRecordedContent(out)
	return nil
}

func (c AuthorityBoundaryRecordedContent) EventTypeName() string {
	return EventTypeAuthorityBoundaryRecorded.Value()
}

type AuthorityResidualRecordedContent struct {
	authorityEvidenceContent
	ResidualID      string                    `json:"residual_id"`
	SubjectRef      string                    `json:"subject_ref"`
	Status          AuthorityResidualStatus   `json:"status"`
	Severity        AuthorityResidualSeverity `json:"severity"`
	Description     string                    `json:"description"`
	RequiredAction  string                    `json:"required_action,omitempty"`
	SourceIssueRefs []NativeEvidenceIssueRef  `json:"source_issue_refs,omitempty"`
	AuthorityRefs   []string                  `json:"authority_refs,omitempty"`
	EvidenceRefs    []string                  `json:"evidence_refs,omitempty"`
	RecordedAt      string                    `json:"recorded_at,omitempty"`
}

func (c *AuthorityResidualRecordedContent) UnmarshalJSON(data []byte) error {
	if err := requireAuthorityEvidenceJSONField(data, "status"); err != nil {
		return err
	}
	if err := requireAuthorityEvidenceJSONField(data, "severity"); err != nil {
		return err
	}
	type alias AuthorityResidualRecordedContent
	var out alias
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	if err := validateAuthorityResidualRecorded(AuthorityResidualRecordedContent(out)); err != nil {
		return err
	}
	*c = AuthorityResidualRecordedContent(out)
	return nil
}

func (c AuthorityResidualRecordedContent) EventTypeName() string {
	return EventTypeAuthorityResidualRecorded.Value()
}

type AuthorityStoreGovernanceRecordedContent struct {
	authorityEvidenceContent
	GovernanceID           string                    `json:"governance_id"`
	StoreName              string                    `json:"store_name"`
	SchemaVersion          string                    `json:"schema_version"`
	WriteStatus            AuthorityStoreWriteStatus `json:"write_status"`
	MigrationRefs          []string                  `json:"migration_refs,omitempty"`
	RequiredValidationRefs []string                  `json:"required_validation_refs,omitempty"`
	AuthorityRefs          []string                  `json:"authority_refs,omitempty"`
	SourceIssueRefs        []NativeEvidenceIssueRef  `json:"source_issue_refs,omitempty"`
	EvidenceRefs           []string                  `json:"evidence_refs,omitempty"`
	RecordedAt             string                    `json:"recorded_at,omitempty"`
}

func (c *AuthorityStoreGovernanceRecordedContent) UnmarshalJSON(data []byte) error {
	if err := requireAuthorityEvidenceJSONField(data, "write_status"); err != nil {
		return err
	}
	type alias AuthorityStoreGovernanceRecordedContent
	var out alias
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	if err := validateAuthorityStoreGovernanceRecorded(AuthorityStoreGovernanceRecordedContent(out)); err != nil {
		return err
	}
	*c = AuthorityStoreGovernanceRecordedContent(out)
	return nil
}

func (c AuthorityStoreGovernanceRecordedContent) EventTypeName() string {
	return EventTypeAuthorityStoreGovernanceRecorded.Value()
}

func validateAuthorityEvidenceContent(content EventContent) error {
	switch c := content.(type) {
	case AuthorityDecisionRecordedContent:
		return validateAuthorityDecisionRecorded(c)
	case AuthorityBoundaryRecordedContent:
		return validateAuthorityBoundaryRecorded(c)
	case AuthorityResidualRecordedContent:
		return validateAuthorityResidualRecorded(c)
	case AuthorityStoreGovernanceRecordedContent:
		return validateAuthorityStoreGovernanceRecorded(c)
	default:
		return fmt.Errorf("unexpected authority evidence content %T", content)
	}
}

func validateAuthorityDecisionRecorded(c AuthorityDecisionRecordedContent) error {
	if err := requireAuthorityEvidenceString("decision_id", c.DecisionID); err != nil {
		return err
	}
	if err := requireAuthorityEvidenceString("subject_type", c.SubjectType); err != nil {
		return err
	}
	if err := requireAuthorityEvidenceString("subject_ref", c.SubjectRef); err != nil {
		return err
	}
	if !c.Outcome.IsValid() {
		return fmt.Errorf("invalid authority decision outcome %q", c.Outcome)
	}
	if (c.Outcome == AuthorityDecisionOutcomeAutonomous || c.Outcome == AuthorityDecisionOutcomeApprovalRequired) && len(c.AuthorityRefs) == 0 {
		return fmt.Errorf("%s authority decision requires at least one authority_ref", c.Outcome)
	}
	if err := validateAuthorityProtectedActions(c.ProtectedActions); err != nil {
		return err
	}
	return validateNativeEvidenceIssueRefs(c.SourceIssueRefs)
}

func validateAuthorityBoundaryRecorded(c AuthorityBoundaryRecordedContent) error {
	if err := requireAuthorityEvidenceString("boundary_id", c.BoundaryID); err != nil {
		return err
	}
	if err := requireAuthorityEvidenceString("subject_ref", c.SubjectRef); err != nil {
		return err
	}
	if !c.BoundaryType.IsValid() {
		return fmt.Errorf("invalid authority boundary type %q", c.BoundaryType)
	}
	if !c.State.IsValid() {
		return fmt.Errorf("invalid authority boundary state %q", c.State)
	}
	if c.BoundaryType == AuthorityBoundaryProtectedAction && len(c.ProtectedActions) == 0 {
		return fmt.Errorf("protected_action boundary requires at least one protected action")
	}
	if c.State == AuthorityBoundaryStateAuthorized && len(c.RequiredAuthorityRefs) == 0 {
		return fmt.Errorf("authorized boundary requires at least one required_authority_ref")
	}
	if err := validateAuthorityProtectedActions(c.ProtectedActions); err != nil {
		return err
	}
	return validateNativeEvidenceIssueRefs(c.SourceIssueRefs)
}

func validateAuthorityResidualRecorded(c AuthorityResidualRecordedContent) error {
	if err := requireAuthorityEvidenceString("residual_id", c.ResidualID); err != nil {
		return err
	}
	if err := requireAuthorityEvidenceString("subject_ref", c.SubjectRef); err != nil {
		return err
	}
	if err := requireAuthorityEvidenceString("description", c.Description); err != nil {
		return err
	}
	if !c.Status.IsValid() {
		return fmt.Errorf("invalid authority residual status %q", c.Status)
	}
	if !c.Severity.IsValid() {
		return fmt.Errorf("invalid authority residual severity %q", c.Severity)
	}
	if (c.Status == AuthorityResidualStatusOpen || c.Status == AuthorityResidualStatusCarried) && c.RequiredAction == "" {
		return fmt.Errorf("%s residual requires required_action", c.Status)
	}
	return validateNativeEvidenceIssueRefs(c.SourceIssueRefs)
}

func validateAuthorityStoreGovernanceRecorded(c AuthorityStoreGovernanceRecordedContent) error {
	if err := requireAuthorityEvidenceString("governance_id", c.GovernanceID); err != nil {
		return err
	}
	if err := requireAuthorityEvidenceString("store_name", c.StoreName); err != nil {
		return err
	}
	if err := requireAuthorityEvidenceString("schema_version", c.SchemaVersion); err != nil {
		return err
	}
	if !c.WriteStatus.IsValid() {
		return fmt.Errorf("invalid authority store write status %q", c.WriteStatus)
	}
	if c.WriteStatus == AuthorityStoreWriteStatusWritePathAuthorized && len(c.AuthorityRefs) == 0 {
		return fmt.Errorf("write_path_authorized requires at least one authority_ref")
	}
	if c.WriteStatus == AuthorityStoreWriteStatusWritePathAuthorized && len(c.RequiredValidationRefs) == 0 {
		return fmt.Errorf("write_path_authorized requires at least one required_validation_ref")
	}
	return validateNativeEvidenceIssueRefs(c.SourceIssueRefs)
}

func validateAuthorityProtectedActions(actions []ProtectedAction) error {
	for i, action := range actions {
		if !IsProtectedAction(string(action)) {
			return fmt.Errorf("protected_actions[%d] is not a known protected action: %q", i, action)
		}
	}
	return nil
}

func requireAuthorityEvidenceString(field string, value string) error {
	if value == "" {
		return fmt.Errorf("authority evidence field %q is required", field)
	}
	return nil
}

func requireAuthorityEvidenceJSONField(data []byte, field string) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if _, ok := fields[field]; !ok {
		return fmt.Errorf("authority evidence field %q is required", field)
	}
	return nil
}
