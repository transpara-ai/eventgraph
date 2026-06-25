package event

import "github.com/transpara-ai/eventgraph/go/pkg/types"

var (
	EventTypeNativeTestRunRecorded     = types.MustEventType("evidence.testrun.recorded")
	EventTypeNativeGateResultRecorded  = types.MustEventType("evidence.gateresult.recorded")
	EventTypeNativeAuditReportRecorded = types.MustEventType("evidence.auditreport.recorded")
)

func AllNativeEvidenceEventTypes() []types.EventType {
	return []types.EventType{
		EventTypeNativeTestRunRecorded,
		EventTypeNativeGateResultRecorded,
		EventTypeNativeAuditReportRecorded,
	}
}
