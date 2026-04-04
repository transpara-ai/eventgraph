package event

import "github.com/lovyou-ai/eventgraph/go/pkg/types"

// Hive event type constants. All use the "hive." prefix.
var (
	EventTypeGapDetected     = types.MustEventType("hive.gap.detected")
	EventTypeDirectiveIssued = types.MustEventType("hive.directive.issued")
)

// AllHiveEventTypes returns all registered hive event types.
func AllHiveEventTypes() []types.EventType {
	return []types.EventType{
		EventTypeGapDetected,
		EventTypeDirectiveIssued,
	}
}
