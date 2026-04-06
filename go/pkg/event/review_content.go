package event

// reviewContent is embedded in all review content types to satisfy the
// EventContent interface's Accept method.
type reviewContent struct{}

func (reviewContent) Accept(EventContentVisitor) {}

// CodeReviewContent represents a code review verdict from the Reviewer agent.
// Emitted when the Reviewer evaluates a completed task's code changes.
type CodeReviewContent struct {
	reviewContent
	TaskID     string   `json:"task_id"`
	Verdict    string   `json:"verdict"`    // "approve", "request_changes", "reject"
	Summary    string   `json:"summary"`
	Issues     []string `json:"issues"`
	Confidence float64  `json:"confidence"` // 0.0–1.0
}

func (c CodeReviewContent) EventTypeName() string {
	return EventTypeCodeReviewSubmitted.Value()
}
