package event

import "github.com/transpara-ai/eventgraph/go/pkg/types"

var (
	EventTypeIssueScanRunProjected          = types.MustEventType("issuescan.run.projected")
	EventTypeIssueScanStageProjected        = types.MustEventType("issuescan.stage.projected")
	EventTypeIssueScanBlockerProjected      = types.MustEventType("issuescan.blocker.projected")
	EventTypeIssueScanLineageProjected      = types.MustEventType("issuescan.lineage.projected")
	EventTypeIssueScanSourceMarkerProjected = types.MustEventType("issuescan.source.marker.projected")
)

func AllIssueScanProjectionEventTypes() []types.EventType {
	return []types.EventType{
		EventTypeIssueScanRunProjected,
		EventTypeIssueScanStageProjected,
		EventTypeIssueScanBlockerProjected,
		EventTypeIssueScanLineageProjected,
		EventTypeIssueScanSourceMarkerProjected,
	}
}
