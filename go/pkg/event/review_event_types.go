package event

import "github.com/transpara-ai/eventgraph/go/pkg/types"

// Review event type constants. All use the "code.review." prefix.
var (
	EventTypeCodeReviewSubmitted = types.MustEventType("code.review.submitted")
)

// AllReviewEventTypes returns all registered review event types.
func AllReviewEventTypes() []types.EventType {
	return []types.EventType{
		EventTypeCodeReviewSubmitted,
	}
}
