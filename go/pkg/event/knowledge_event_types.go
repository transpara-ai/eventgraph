package event

import "github.com/lovyou-ai/eventgraph/go/pkg/types"

// Knowledge event type constants. All use the "knowledge." prefix.
var (
	EventTypeKnowledgeInsightRecorded   = types.MustEventType("knowledge.insight.recorded")
	EventTypeKnowledgeInsightSuperseded = types.MustEventType("knowledge.insight.superseded")
	EventTypeKnowledgeInsightExpired    = types.MustEventType("knowledge.insight.expired")
)

// AllKnowledgeEventTypes returns all registered knowledge event types.
func AllKnowledgeEventTypes() []types.EventType {
	return []types.EventType{
		EventTypeKnowledgeInsightRecorded,
		EventTypeKnowledgeInsightSuperseded,
		EventTypeKnowledgeInsightExpired,
	}
}
