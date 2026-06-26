package event

import "github.com/transpara-ai/eventgraph/go/pkg/types"

var (
	EventTypeAuthorityDecisionRecorded        = types.MustEventType("authority.decision.recorded")
	EventTypeAuthorityBoundaryRecorded        = types.MustEventType("authority.boundary.recorded")
	EventTypeAuthorityResidualRecorded        = types.MustEventType("authority.residual.recorded")
	EventTypeAuthorityStoreGovernanceRecorded = types.MustEventType("authority.storegovernance.recorded")
)

func AllAuthorityEvidenceEventTypes() []types.EventType {
	return []types.EventType{
		EventTypeAuthorityDecisionRecorded,
		EventTypeAuthorityBoundaryRecorded,
		EventTypeAuthorityResidualRecorded,
		EventTypeAuthorityStoreGovernanceRecorded,
	}
}
