package event

import (
	"sort"

	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// knowledgeContent is embedded in all knowledge content types to satisfy the
// EventContent interface's Accept method. Knowledge content types use their
// own visitor rather than the base EventContentVisitor.
type knowledgeContent struct{}

func (knowledgeContent) Accept(EventContentVisitor) {}

// KnowledgeInsightContent represents a distilled insight from pattern analysis.
// Emitted by automated distillers or knowledge agents.
// Use NewKnowledgeInsightContent to ensure RelevantRoles is sorted for
// deterministic canonical-form hashing.
type KnowledgeInsightContent struct {
	knowledgeContent
	InsightID     string                `json:"insight_id"`
	Domain        string                `json:"domain"`
	Summary       string                `json:"summary"`
	RelevantRoles []string              `json:"relevant_roles"`
	Confidence    types.Score           `json:"confidence"`
	EvidenceCount int                   `json:"evidence_count"`
	Source        string                `json:"source"`
	TTL           int                   `json:"ttl"`
	SupersedesID  types.Option[string]  `json:"supersedes_id"`
}

func (c KnowledgeInsightContent) EventTypeName() string {
	return EventTypeKnowledgeInsightRecorded.Value()
}

// NewKnowledgeInsightContent creates a KnowledgeInsightContent with
// RelevantRoles sorted lexicographically for deterministic hashing.
func NewKnowledgeInsightContent(
	insightID, domain, summary string,
	relevantRoles []string,
	confidence types.Score,
	evidenceCount int,
	source string,
	ttl int,
	supersedesID types.Option[string],
) KnowledgeInsightContent {
	sorted := make([]string, len(relevantRoles))
	copy(sorted, relevantRoles)
	sort.Strings(sorted)
	return KnowledgeInsightContent{
		InsightID:     insightID,
		Domain:        domain,
		Summary:       summary,
		RelevantRoles: sorted,
		Confidence:    confidence,
		EvidenceCount: evidenceCount,
		Source:        source,
		TTL:           ttl,
		SupersedesID:  supersedesID,
	}
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
