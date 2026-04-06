package event

// knowledgeContent is embedded in all knowledge content types to satisfy the
// EventContent interface's Accept method. Knowledge content types use their
// own visitor rather than the base EventContentVisitor.
type knowledgeContent struct{}

func (knowledgeContent) Accept(EventContentVisitor) {}

// KnowledgeInsightContent represents a distilled insight from pattern analysis.
// Emitted by automated distillers or knowledge agents.
type KnowledgeInsightContent struct {
	knowledgeContent
	InsightID     string   `json:"insight_id"`
	Domain        string   `json:"domain"`
	Summary       string   `json:"summary"`
	RelevantRoles []string `json:"relevant_roles"`
	Confidence    float64  `json:"confidence"`
	EvidenceCount int      `json:"evidence_count"`
	Source        string   `json:"source"`
	TTL           int      `json:"ttl"`
	SupersedesID  string   `json:"supersedes_id,omitempty"`
}

func (c KnowledgeInsightContent) EventTypeName() string {
	return EventTypeKnowledgeInsightRecorded.Value()
}

// KnowledgeSupersessionContent records when an insight is replaced by a newer version.
type KnowledgeSupersessionContent struct {
	knowledgeContent
	OldInsightID string `json:"old_insight_id"`
	NewInsightID string `json:"new_insight_id"`
	Reason       string `json:"reason"`
}

func (c KnowledgeSupersessionContent) EventTypeName() string {
	return EventTypeKnowledgeInsightSuperseded.Value()
}

// KnowledgeExpirationContent records when an insight's TTL expires.
type KnowledgeExpirationContent struct {
	knowledgeContent
	InsightID string `json:"insight_id"`
	Reason    string `json:"reason"`
}

func (c KnowledgeExpirationContent) EventTypeName() string {
	return EventTypeKnowledgeInsightExpired.Value()
}
