package event

import "github.com/transpara-ai/eventgraph/go/pkg/types"

const TLC51ProtocolVersion = "factory-tlc51/v1"

var (
	EventTypeTLC51PlanRecorded               = types.MustEventType("factory.tlc51.plan.recorded")
	EventTypeTLC51PlanSuperseded             = types.MustEventType("factory.tlc51.plan.superseded")
	EventTypeTLC51ObligationReady            = types.MustEventType("factory.tlc51.obligation.ready")
	EventTypeTLC51ObligationClaimed          = types.MustEventType("factory.tlc51.obligation.claimed")
	EventTypeTLC51ObligationRunning          = types.MustEventType("factory.tlc51.obligation.running")
	EventTypeTLC51ObligationTerminal         = types.MustEventType("factory.tlc51.obligation.terminal")
	EventTypeTLC51EvidenceLinked             = types.MustEventType("factory.tlc51.evidence.linked")
	EventTypeTLC51DecisionRecorded           = types.MustEventType("factory.tlc51.decision.recorded")
	EventTypeTLC51DecisionInvalidated        = types.MustEventType("factory.tlc51.decision.invalidated")
	EventTypeTLC51ProtectedEffectProposed    = types.MustEventType("factory.tlc51.effect.proposed")
	EventTypeTLC51ProtectedEffectObserved    = types.MustEventType("factory.tlc51.effect.observed")
	EventTypeTLC51ProtectedEffectReconciled  = types.MustEventType("factory.tlc51.effect.reconciled")
	EventTypeTLC51ProtectedEffectTerminal    = types.MustEventType("factory.tlc51.effect.terminal")
	EventTypeTLC51HumanInterventionRequested = types.MustEventType("factory.tlc51.human.requested")
	EventTypeTLC51HumanInterventionResolved  = types.MustEventType("factory.tlc51.human.resolved")
	EventTypeTLC51CutoverRecorded            = types.MustEventType("factory.tlc51.cutover.recorded")
)

// AllTLC51EventTypes returns the complete closed factory-tlc51/v1 event set.
func AllTLC51EventTypes() []types.EventType {
	return []types.EventType{
		EventTypeTLC51PlanRecorded,
		EventTypeTLC51PlanSuperseded,
		EventTypeTLC51ObligationReady,
		EventTypeTLC51ObligationClaimed,
		EventTypeTLC51ObligationRunning,
		EventTypeTLC51ObligationTerminal,
		EventTypeTLC51EvidenceLinked,
		EventTypeTLC51DecisionRecorded,
		EventTypeTLC51DecisionInvalidated,
		EventTypeTLC51ProtectedEffectProposed,
		EventTypeTLC51ProtectedEffectObserved,
		EventTypeTLC51ProtectedEffectReconciled,
		EventTypeTLC51ProtectedEffectTerminal,
		EventTypeTLC51HumanInterventionRequested,
		EventTypeTLC51HumanInterventionResolved,
		EventTypeTLC51CutoverRecorded,
	}
}
