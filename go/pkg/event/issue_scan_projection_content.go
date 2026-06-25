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
	default:
		return fmt.Errorf("unexpected issue-scan projection content %T", content)
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
