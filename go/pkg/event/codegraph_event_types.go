package event

import "github.com/transpara-ai/eventgraph/go/pkg/types"

// Code graph event type constants. All use the "codegraph." prefix.
var (
	// Data events
	EventTypeCodeGraphEntityDefined  = types.MustEventType("codegraph.entity.defined")
	EventTypeCodeGraphEntityModified = types.MustEventType("codegraph.entity.modified")
	EventTypeCodeGraphEntityRelated  = types.MustEventType("codegraph.entity.related")
	EventTypeCodeGraphStateTransitioned = types.MustEventType("codegraph.state.transitioned")

	// Logic events
	EventTypeCodeGraphTransformApplied     = types.MustEventType("codegraph.logic.transform.applied")
	EventTypeCodeGraphConditionEvaluated   = types.MustEventType("codegraph.logic.condition.evaluated")
	EventTypeCodeGraphSequenceExecuted     = types.MustEventType("codegraph.logic.sequence.executed")
	EventTypeCodeGraphLoopIterated         = types.MustEventType("codegraph.logic.loop.iterated")
	EventTypeCodeGraphTriggerFired         = types.MustEventType("codegraph.logic.trigger.fired")
	EventTypeCodeGraphConstraintViolated   = types.MustEventType("codegraph.logic.constraint.violated")

	// IO events
	EventTypeCodeGraphQueryExecuted       = types.MustEventType("codegraph.io.query.executed")
	EventTypeCodeGraphCommandExecuted     = types.MustEventType("codegraph.io.command.executed")
	EventTypeCodeGraphCommandRejected     = types.MustEventType("codegraph.io.command.rejected")
	EventTypeCodeGraphSubscribeRegistered = types.MustEventType("codegraph.io.subscribe.registered")
	EventTypeCodeGraphAuthorizeGranted    = types.MustEventType("codegraph.io.authorize.granted")
	EventTypeCodeGraphAuthorizeDenied     = types.MustEventType("codegraph.io.authorize.denied")
	EventTypeCodeGraphSearchExecuted      = types.MustEventType("codegraph.io.search.executed")
	EventTypeCodeGraphInteropSent         = types.MustEventType("codegraph.io.interop.sent")
	EventTypeCodeGraphInteropReceived     = types.MustEventType("codegraph.io.interop.received")

	// UI events
	EventTypeCodeGraphViewRendered        = types.MustEventType("codegraph.ui.view.rendered")
	EventTypeCodeGraphActionInvoked       = types.MustEventType("codegraph.ui.action.invoked")
	EventTypeCodeGraphNavigationTriggered = types.MustEventType("codegraph.ui.navigation.triggered")
	EventTypeCodeGraphFeedbackEmitted     = types.MustEventType("codegraph.ui.feedback.emitted")
	EventTypeCodeGraphAlertDispatched     = types.MustEventType("codegraph.ui.alert.dispatched")
	EventTypeCodeGraphDragCompleted       = types.MustEventType("codegraph.ui.drag.completed")
	EventTypeCodeGraphSelectionChanged    = types.MustEventType("codegraph.ui.selection.changed")
	EventTypeCodeGraphConfirmationResolved = types.MustEventType("codegraph.ui.confirmation.resolved")

	// Aesthetic events
	EventTypeCodeGraphSkinApplied = types.MustEventType("codegraph.aesthetic.skin.applied")

	// Temporal events
	EventTypeCodeGraphUndoExecuted    = types.MustEventType("codegraph.temporal.undo.executed")
	EventTypeCodeGraphRetryAttempted  = types.MustEventType("codegraph.temporal.retry.attempted")

	// Resilience events
	EventTypeCodeGraphOfflineEntered = types.MustEventType("codegraph.resilience.offline.entered")
	EventTypeCodeGraphOfflineSynced  = types.MustEventType("codegraph.resilience.offline.synced")

	// Structural events
	EventTypeCodeGraphScopeDefined = types.MustEventType("codegraph.structural.scope.defined")

	// Social events
	EventTypeCodeGraphPresenceUpdated    = types.MustEventType("codegraph.social.presence.updated")
	EventTypeCodeGraphSalienceTriggered  = types.MustEventType("codegraph.social.salience.triggered")
)

// AllCodeGraphEventTypes returns all registered code graph event types.
func AllCodeGraphEventTypes() []types.EventType {
	return []types.EventType{
		// Data
		EventTypeCodeGraphEntityDefined, EventTypeCodeGraphEntityModified,
		EventTypeCodeGraphEntityRelated, EventTypeCodeGraphStateTransitioned,
		// Logic
		EventTypeCodeGraphTransformApplied, EventTypeCodeGraphConditionEvaluated,
		EventTypeCodeGraphSequenceExecuted, EventTypeCodeGraphLoopIterated,
		EventTypeCodeGraphTriggerFired, EventTypeCodeGraphConstraintViolated,
		// IO
		EventTypeCodeGraphQueryExecuted, EventTypeCodeGraphCommandExecuted,
		EventTypeCodeGraphCommandRejected, EventTypeCodeGraphSubscribeRegistered,
		EventTypeCodeGraphAuthorizeGranted, EventTypeCodeGraphAuthorizeDenied,
		EventTypeCodeGraphSearchExecuted, EventTypeCodeGraphInteropSent,
		EventTypeCodeGraphInteropReceived,
		// UI
		EventTypeCodeGraphViewRendered, EventTypeCodeGraphActionInvoked,
		EventTypeCodeGraphNavigationTriggered, EventTypeCodeGraphFeedbackEmitted,
		EventTypeCodeGraphAlertDispatched, EventTypeCodeGraphDragCompleted,
		EventTypeCodeGraphSelectionChanged, EventTypeCodeGraphConfirmationResolved,
		// Aesthetic
		EventTypeCodeGraphSkinApplied,
		// Temporal
		EventTypeCodeGraphUndoExecuted, EventTypeCodeGraphRetryAttempted,
		// Resilience
		EventTypeCodeGraphOfflineEntered, EventTypeCodeGraphOfflineSynced,
		// Structural
		EventTypeCodeGraphScopeDefined,
		// Social
		EventTypeCodeGraphPresenceUpdated, EventTypeCodeGraphSalienceTriggered,
	}
}
