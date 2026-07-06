package event

import (
	"encoding/json"
	"fmt"
)

type issueScanProjectionContent struct{}

func (issueScanProjectionContent) Accept(EventContentVisitor) {}

type IssueScanRunState string

const (
	IssueScanRunStateQueued         IssueScanRunState = "queued"
	IssueScanRunStateDispatched     IssueScanRunState = "dispatched"
	IssueScanRunStateRunning        IssueScanRunState = "running"
	IssueScanRunStateBlocked        IssueScanRunState = "blocked"
	IssueScanRunStateParked         IssueScanRunState = "parked"
	IssueScanRunStateHumanAction    IssueScanRunState = "human_action"
	IssueScanRunStateReadyForHuman  IssueScanRunState = "ready_for_human"
	IssueScanRunStateSuperseded     IssueScanRunState = "superseded"
	IssueScanRunStateCompleted      IssueScanRunState = "completed"
	IssueScanRunStateProjectionOnly IssueScanRunState = "projection_only"
)

func (s IssueScanRunState) IsValid() bool {
	switch s {
	case IssueScanRunStateQueued, IssueScanRunStateDispatched, IssueScanRunStateRunning,
		IssueScanRunStateBlocked, IssueScanRunStateParked, IssueScanRunStateHumanAction,
		IssueScanRunStateReadyForHuman, IssueScanRunStateSuperseded, IssueScanRunStateCompleted,
		IssueScanRunStateProjectionOnly:
		return true
	default:
		return false
	}
}

func (s *IssueScanRunState) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	state := IssueScanRunState(value)
	if !state.IsValid() {
		return fmt.Errorf("invalid issue-scan run state %q", value)
	}
	*s = state
	return nil
}

type IssueScanStageState string

const (
	IssueScanStageStateDeclared       IssueScanStageState = "declared"
	IssueScanStageStateBlocked        IssueScanStageState = "blocked"
	IssueScanStageStateReady          IssueScanStageState = "ready"
	IssueScanStageStateRunning        IssueScanStageState = "running"
	IssueScanStageStateComplete       IssueScanStageState = "complete"
	IssueScanStageStateHumanAction    IssueScanStageState = "human_action"
	IssueScanStageStateParked         IssueScanStageState = "parked"
	IssueScanStageStateSuperseded     IssueScanStageState = "superseded"
	IssueScanStageStateProjectionOnly IssueScanStageState = "projection_only"
)

func (s IssueScanStageState) IsValid() bool {
	switch s {
	case IssueScanStageStateDeclared, IssueScanStageStateBlocked, IssueScanStageStateReady,
		IssueScanStageStateRunning, IssueScanStageStateComplete, IssueScanStageStateHumanAction,
		IssueScanStageStateParked, IssueScanStageStateSuperseded, IssueScanStageStateProjectionOnly:
		return true
	default:
		return false
	}
}

func (s *IssueScanStageState) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	state := IssueScanStageState(value)
	if !state.IsValid() {
		return fmt.Errorf("invalid issue-scan stage state %q", value)
	}
	*s = state
	return nil
}

type IssueScanBlockerType string

const (
	IssueScanBlockerNeedsHumanScope     IssueScanBlockerType = "needs_human_scope"
	IssueScanBlockerProtectedAction     IssueScanBlockerType = "protected_action"
	IssueScanBlockerStaleTarget         IssueScanBlockerType = "stale_target"
	IssueScanBlockerDuplicateChain      IssueScanBlockerType = "duplicate_chain"
	IssueScanBlockerMissingGateEvidence IssueScanBlockerType = "missing_gate_evidence"
)

func (t IssueScanBlockerType) IsValid() bool {
	switch t {
	case IssueScanBlockerNeedsHumanScope, IssueScanBlockerProtectedAction, IssueScanBlockerStaleTarget,
		IssueScanBlockerDuplicateChain, IssueScanBlockerMissingGateEvidence:
		return true
	default:
		return false
	}
}

func (t *IssueScanBlockerType) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	blockerType := IssueScanBlockerType(value)
	if !blockerType.IsValid() {
		return fmt.Errorf("invalid issue-scan blocker type %q", value)
	}
	*t = blockerType
	return nil
}

type IssueScanSourceMarkerTransition string

const (
	IssueScanSourceMarkerAcquired          IssueScanSourceMarkerTransition = "acquired"
	IssueScanSourceMarkerParkedHumanAction IssueScanSourceMarkerTransition = "parked_human_action"
	IssueScanSourceMarkerReadyForHuman     IssueScanSourceMarkerTransition = "ready_for_human"
	IssueScanSourceMarkerCompleted         IssueScanSourceMarkerTransition = "completed"
	IssueScanSourceMarkerAbandoned         IssueScanSourceMarkerTransition = "abandoned"
	IssueScanSourceMarkerSuperseded        IssueScanSourceMarkerTransition = "superseded"
)

func (t IssueScanSourceMarkerTransition) IsValid() bool {
	switch t {
	case IssueScanSourceMarkerAcquired, IssueScanSourceMarkerParkedHumanAction,
		IssueScanSourceMarkerReadyForHuman, IssueScanSourceMarkerCompleted,
		IssueScanSourceMarkerAbandoned, IssueScanSourceMarkerSuperseded:
		return true
	default:
		return false
	}
}

func (t *IssueScanSourceMarkerTransition) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	transition := IssueScanSourceMarkerTransition(value)
	if !transition.IsValid() {
		return fmt.Errorf("invalid issue-scan source-marker transition %q", value)
	}
	*t = transition
	return nil
}

type IssueScanIssueRef struct {
	Repo        string   `json:"repo"`
	Number      int      `json:"number"`
	URL         string   `json:"url,omitempty"`
	Title       string   `json:"title,omitempty"`
	State       string   `json:"state,omitempty"`
	StateReason string   `json:"state_reason,omitempty"`
	Labels      []string `json:"labels,omitempty"`
}

type IssueScanRunProjectedContent struct {
	issueScanProjectionContent
	RunID            string              `json:"run_id"`
	FactoryOrderID   string              `json:"factory_order_id,omitempty"`
	LifecycleVersion string              `json:"lifecycle_version"`
	State            IssueScanRunState   `json:"state"`
	TargetIssue      IssueScanIssueRef   `json:"target_issue"`
	SelectedIssue    IssueScanIssueRef   `json:"selected_issue"`
	CandidateIssues  []IssueScanIssueRef `json:"candidate_issues,omitempty"`
	SourceRefs       []string            `json:"source_refs,omitempty"`
	EvidenceRefs     []string            `json:"evidence_refs,omitempty"`
}

func (c *IssueScanRunProjectedContent) UnmarshalJSON(data []byte) error {
	if err := requireIssueScanProjectionJSONField(data, "state"); err != nil {
		return err
	}
	type alias IssueScanRunProjectedContent
	var out alias
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	if !out.State.IsValid() {
		return fmt.Errorf("invalid issue-scan run state %q", out.State)
	}
	*c = IssueScanRunProjectedContent(out)
	return nil
}

func (c IssueScanRunProjectedContent) EventTypeName() string {
	return EventTypeIssueScanRunProjected.Value()
}

type IssueScanStageProjectedContent struct {
	issueScanProjectionContent
	RunID             string              `json:"run_id"`
	FactoryOrderID    string              `json:"factory_order_id,omitempty"`
	StageID           string              `json:"stage_id"`
	StageNumber       int                 `json:"stage_number"`
	StageCount        int                 `json:"stage_count,omitempty"`
	CanonicalTaskID   string              `json:"canonical_task_id"`
	TaskID            string              `json:"task_id,omitempty"`
	CurrentState      IssueScanStageState `json:"current_state"`
	CompletionGate    string              `json:"completion_gate"`
	AuthorityBoundary string              `json:"authority_boundary"`
	AssignedAgentIDs  []string            `json:"assigned_agent_ids,omitempty"`
	TouchingAgentIDs  []string            `json:"touching_agent_ids,omitempty"`
	EvidenceRefs      []string            `json:"evidence_refs,omitempty"`
	SourceRefs        []string            `json:"source_refs,omitempty"`
}

func (c *IssueScanStageProjectedContent) UnmarshalJSON(data []byte) error {
	if err := requireIssueScanProjectionJSONField(data, "current_state"); err != nil {
		return err
	}
	type alias IssueScanStageProjectedContent
	var out alias
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	if !out.CurrentState.IsValid() {
		return fmt.Errorf("invalid issue-scan stage state %q", out.CurrentState)
	}
	*c = IssueScanStageProjectedContent(out)
	return nil
}

func (c IssueScanStageProjectedContent) EventTypeName() string {
	return EventTypeIssueScanStageProjected.Value()
}

type IssueScanBlockerProjectedContent struct {
	issueScanProjectionContent
	RunID          string               `json:"run_id"`
	FactoryOrderID string               `json:"factory_order_id,omitempty"`
	StageID        string               `json:"stage_id,omitempty"`
	BlockerType    IssueScanBlockerType `json:"blocker_type"`
	Reason         string               `json:"reason,omitempty"`
	RequiredAction string               `json:"required_action"`
	EvidenceRefs   []string             `json:"evidence_refs,omitempty"`
	SourceRefs     []string             `json:"source_refs,omitempty"`
}

func (c *IssueScanBlockerProjectedContent) UnmarshalJSON(data []byte) error {
	if err := requireIssueScanProjectionJSONField(data, "blocker_type"); err != nil {
		return err
	}
	type alias IssueScanBlockerProjectedContent
	var out alias
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	if !out.BlockerType.IsValid() {
		return fmt.Errorf("invalid issue-scan blocker type %q", out.BlockerType)
	}
	*c = IssueScanBlockerProjectedContent(out)
	return nil
}

func (c IssueScanBlockerProjectedContent) EventTypeName() string {
	return EventTypeIssueScanBlockerProjected.Value()
}

type IssueScanLineageProjectedContent struct {
	issueScanProjectionContent
	RunID             string   `json:"run_id"`
	FactoryOrderID    string   `json:"factory_order_id,omitempty"`
	StageID           string   `json:"stage_id,omitempty"`
	CanonicalTaskID   string   `json:"canonical_task_id"`
	PrimaryTaskID     string   `json:"primary_task_id,omitempty"`
	TaskIDs           []string `json:"task_ids"`
	DuplicateTaskIDs  []string `json:"duplicate_task_ids,omitempty"`
	DuplicateOf       string   `json:"duplicate_of,omitempty"`
	SupersededTaskIDs []string `json:"superseded_task_ids,omitempty"`
	SourceRefs        []string `json:"source_refs,omitempty"`
}

func (c IssueScanLineageProjectedContent) EventTypeName() string {
	return EventTypeIssueScanLineageProjected.Value()
}

type IssueScanMarkerTargetRef struct {
	Repository  string `json:"repository"`
	IssueNumber int    `json:"issue_number"`
}

type IssueScanMarkerBlockerRef struct {
	Reason       IssueScanBlockerType `json:"reason"`
	Detail       string               `json:"detail,omitempty"`
	EvidenceRefs []string             `json:"evidence_refs,omitempty"`
}

type IssueScanMarkerGateRef struct {
	Gate         string   `json:"gate"`
	EvidenceRefs []string `json:"evidence_refs,omitempty"`
}

type IssueScanMarkerEvidenceRefs struct {
	TestCaseIDs      []string `json:"test_case_ids,omitempty"`
	TestRunIDs       []string `json:"test_run_ids,omitempty"`
	GateResultIDs    []string `json:"gate_result_ids,omitempty"`
	FailureIDs       []string `json:"failure_ids,omitempty"`
	RepairAttemptIDs []string `json:"repair_attempt_ids,omitempty"`
	WaiverIDs        []string `json:"waiver_ids,omitempty"`
}

type IssueScanMarkerWorkRef struct {
	SchemaVersion          string                      `json:"schema_version"`
	ProjectionKind         string                      `json:"projection_kind"`
	CanonicalSource        string                      `json:"canonical_source"`
	ProjectionOnly         bool                        `json:"projection_only"`
	RunID                  string                      `json:"run_id"`
	Target                 IssueScanMarkerTargetRef    `json:"target"`
	Stage                  string                      `json:"stage"`
	StageNumber            int                         `json:"stage_number"`
	Gate                   string                      `json:"gate"`
	TaskID                 string                      `json:"task_id"`
	CanonicalTaskID        string                      `json:"canonical_task_id"`
	FactoryOrderID         string                      `json:"factory_order_id"`
	RequirementIDs         []string                    `json:"requirement_ids,omitempty"`
	AcceptanceCriterionIDs []string                    `json:"acceptance_criterion_ids,omitempty"`
	LifecycleState         string                      `json:"lifecycle_state"`
	Ready                  bool                        `json:"ready"`
	Blocked                bool                        `json:"blocked"`
	MissingGates           []string                    `json:"missing_gates,omitempty"`
	MissingFacts           []string                    `json:"missing_facts,omitempty"`
	SupersededBy           string                      `json:"superseded_by,omitempty"`
	LastTransitionEvent    string                      `json:"last_transition_event,omitempty"`
	LatestBlocker          *IssueScanMarkerBlockerRef  `json:"latest_blocker,omitempty"`
	LatestGate             *IssueScanMarkerGateRef     `json:"latest_gate,omitempty"`
	VerificationRefs       IssueScanMarkerEvidenceRefs `json:"verification_refs"`
	FailureRepairRefs      IssueScanMarkerEvidenceRefs `json:"failure_repair_refs"`
	SourceIssueRefs        []string                    `json:"source_issue_refs,omitempty"`
	AuthorityExclusions    []string                    `json:"authority_exclusions"`
}

type IssueScanSourceMarkerOutputRef struct {
	System         string   `json:"system"`
	Repository     string   `json:"repository,omitempty"`
	IssueNumber    int      `json:"issue_number,omitempty"`
	CommentID      string   `json:"comment_id,omitempty"`
	CommentURL     string   `json:"comment_url,omitempty"`
	LabelNames     []string `json:"label_names,omitempty"`
	DerivedOutput  bool     `json:"derived_output"`
	ProjectionSink bool     `json:"projection_sink"`
}

type IssueScanSourceMarkerProjectedContent struct {
	issueScanProjectionContent
	SchemaVersion       string                          `json:"schema_version"`
	ProjectionKind      string                          `json:"projection_kind"`
	Transition          IssueScanSourceMarkerTransition `json:"transition"`
	RunID               string                          `json:"run_id"`
	Target              IssueScanIssueRef               `json:"target"`
	StageID             string                          `json:"stage_id"`
	StageNumber         int                             `json:"stage_number"`
	Gate                string                          `json:"gate,omitempty"`
	WorkRef             IssueScanMarkerWorkRef          `json:"work_ref"`
	ActorID             string                          `json:"actor_id"`
	ActorRole           string                          `json:"actor_role"`
	OccurredAt          string                          `json:"occurred_at"`
	IdempotencyKey      string                          `json:"idempotency_key"`
	AuthorityBoundary   string                          `json:"authority_boundary"`
	AuthorityExclusions []string                        `json:"authority_exclusions"`
	EvidenceRefs        IssueScanMarkerEvidenceRefs     `json:"evidence_refs"`
	SourceRefs          []string                        `json:"source_refs,omitempty"`
	GitHubMarker        IssueScanSourceMarkerOutputRef  `json:"github_marker,omitempty"`
	CanonicalSource     string                          `json:"canonical_source"`
	ProjectionOnly      bool                            `json:"projection_only"`
	SupersededBy        string                          `json:"superseded_by,omitempty"`
	StaleTarget         bool                            `json:"stale_target,omitempty"`
}

func (c *IssueScanSourceMarkerProjectedContent) UnmarshalJSON(data []byte) error {
	if err := requireIssueScanProjectionJSONField(data, "transition"); err != nil {
		return err
	}
	type alias IssueScanSourceMarkerProjectedContent
	var out alias
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	if err := validateIssueScanSourceMarkerProjection(IssueScanSourceMarkerProjectedContent(out)); err != nil {
		return err
	}
	*c = IssueScanSourceMarkerProjectedContent(out)
	return nil
}

func (c IssueScanSourceMarkerProjectedContent) EventTypeName() string {
	return EventTypeIssueScanSourceMarkerProjected.Value()
}

func validateIssueScanProjectionContent(content EventContent) error {
	switch c := content.(type) {
	case IssueScanRunProjectedContent:
		if !c.State.IsValid() {
			return fmt.Errorf("invalid issue-scan run state %q", c.State)
		}
	case IssueScanStageProjectedContent:
		if !c.CurrentState.IsValid() {
			return fmt.Errorf("invalid issue-scan stage state %q", c.CurrentState)
		}
	case IssueScanBlockerProjectedContent:
		if !c.BlockerType.IsValid() {
			return fmt.Errorf("invalid issue-scan blocker type %q", c.BlockerType)
		}
	case IssueScanLineageProjectedContent:
		return nil
	case IssueScanSourceMarkerProjectedContent:
		return validateIssueScanSourceMarkerProjection(c)
	default:
		return fmt.Errorf("unexpected issue-scan projection content %T", content)
	}
	return nil
}

func validateIssueScanSourceMarkerProjection(c IssueScanSourceMarkerProjectedContent) error {
	if !c.Transition.IsValid() {
		return fmt.Errorf("invalid issue-scan source-marker transition %q", c.Transition)
	}
	if c.RunID == "" {
		return fmt.Errorf("issue-scan source-marker run_id is required")
	}
	if c.Target.Repo == "" || c.Target.Number <= 0 {
		return fmt.Errorf("issue-scan source-marker target repo and number are required")
	}
	if c.StageID == "" {
		return fmt.Errorf("issue-scan source-marker stage_id is required")
	}
	if c.WorkRef.ProjectionKind != "work.issue_scan.source_marker_ref" {
		return fmt.Errorf("issue-scan source-marker work_ref projection_kind must be work.issue_scan.source_marker_ref")
	}
	if c.WorkRef.CanonicalSource != "work" {
		return fmt.Errorf("issue-scan source-marker work_ref canonical_source must be work")
	}
	if !c.WorkRef.ProjectionOnly {
		return fmt.Errorf("issue-scan source-marker work_ref must be projection_only")
	}
	if c.WorkRef.RunID != c.RunID {
		return fmt.Errorf("issue-scan source-marker work_ref run_id mismatch")
	}
	if c.WorkRef.Target.Repository != c.Target.Repo || c.WorkRef.Target.IssueNumber != c.Target.Number {
		return fmt.Errorf("issue-scan source-marker work_ref target mismatch")
	}
	if c.WorkRef.Stage != c.StageID {
		return fmt.Errorf("issue-scan source-marker work_ref stage mismatch")
	}
	if c.ActorRole == "" {
		return fmt.Errorf("issue-scan source-marker actor_role is required")
	}
	if c.OccurredAt == "" {
		return fmt.Errorf("issue-scan source-marker occurred_at is required")
	}
	if c.IdempotencyKey == "" {
		return fmt.Errorf("issue-scan source-marker idempotency_key is required")
	}
	if c.AuthorityBoundary == "" {
		return fmt.Errorf("issue-scan source-marker authority_boundary is required")
	}
	if c.CanonicalSource != "work_eventgraph_projection" {
		return fmt.Errorf("issue-scan source-marker canonical_source must be work_eventgraph_projection")
	}
	if !c.ProjectionOnly {
		return fmt.Errorf("issue-scan source-marker projection_only must be true")
	}
	if c.GitHubMarker.System == "github" && (!c.GitHubMarker.DerivedOutput || !c.GitHubMarker.ProjectionSink) {
		return fmt.Errorf("issue-scan source-marker github marker must be a derived projection sink")
	}
	if c.Transition == IssueScanSourceMarkerSuperseded && c.SupersededBy == "" && c.WorkRef.SupersededBy == "" {
		return fmt.Errorf("issue-scan source-marker superseded transition requires superseded_by")
	}
	if c.StaleTarget && c.Transition != IssueScanSourceMarkerParkedHumanAction && c.Transition != IssueScanSourceMarkerSuperseded {
		return fmt.Errorf("issue-scan source-marker stale targets must park or supersede")
	}
	return nil
}

func requireIssueScanProjectionJSONField(data []byte, field string) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if _, ok := fields[field]; !ok {
		return fmt.Errorf("issue-scan projection field %q is required", field)
	}
	return nil
}
