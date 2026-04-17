package event

import "time"

// SpecIngestedContent records that a spec has been ingested into the hive
// from an upstream site op. SourceOpID is the bus event ID of the
// translated op that produced the spec.
type SpecIngestedContent struct {
	hiveContent
	SpecRef    string    `json:"spec_ref"`
	SourceOpID string    `json:"source_op_id"`
	IngestedAt time.Time `json:"ingested_at"`
}

func (c SpecIngestedContent) EventTypeName() string { return "hive.spec.ingested" }

// SpecParsedContent records that a spec has been parsed into
// individual tasks.
type SpecParsedContent struct {
	hiveContent
	SpecRef   string    `json:"spec_ref"`
	TaskCount int       `json:"task_count"`
	ParsedAt  time.Time `json:"parsed_at"`
}

func (c SpecParsedContent) EventTypeName() string { return "hive.spec.parsed" }

// SpecAssignedContent records that parsed tasks have been assigned
// to specific agents. Assignments maps task_id → agent_name.
type SpecAssignedContent struct {
	hiveContent
	SpecRef     string            `json:"spec_ref"`
	Assignments map[string]string `json:"assignments"`
	AssignedAt  time.Time         `json:"assigned_at"`
}

func (c SpecAssignedContent) EventTypeName() string { return "hive.spec.assigned" }

// SpecCompletedContent records the terminal state of a spec.
// Outcome is one of "success", "partial", or "failed".
type SpecCompletedContent struct {
	hiveContent
	SpecRef     string    `json:"spec_ref"`
	Outcome     string    `json:"outcome"`
	CompletedAt time.Time `json:"completed_at"`
}

func (c SpecCompletedContent) EventTypeName() string { return "hive.spec.completed" }
