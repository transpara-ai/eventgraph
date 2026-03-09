package event

import "github.com/lovyou-ai/eventgraph/go/pkg/types"

// codegraphContent is embedded in all code graph content types to satisfy
// EventContent.Accept. Code graph events use their own domain, not the base visitor.
type codegraphContent struct{}

func (codegraphContent) Accept(EventContentVisitor) {}

// --- Data events ---

type CodeGraphEntityDefinedContent struct {
	codegraphContent
	EntityType string        `json:"EntityType"`
	EntityID   types.ActorID `json:"EntityID"`
}

func (c CodeGraphEntityDefinedContent) EventTypeName() string { return "codegraph.entity.defined" }

type CodeGraphEntityModifiedContent struct {
	codegraphContent
	EntityType string               `json:"EntityType"`
	EntityID   types.ActorID        `json:"EntityID"`
	Property   string               `json:"Property"`
	Previous   types.Option[string] `json:"Previous"`
	Current    string               `json:"Current"`
}

func (c CodeGraphEntityModifiedContent) EventTypeName() string { return "codegraph.entity.modified" }

type CodeGraphEntityRelatedContent struct {
	codegraphContent
	SourceType string        `json:"SourceType"`
	SourceID   types.ActorID `json:"SourceID"`
	TargetType string        `json:"TargetType"`
	TargetID   types.ActorID `json:"TargetID"`
	Relation   string        `json:"Relation"`
}

func (c CodeGraphEntityRelatedContent) EventTypeName() string { return "codegraph.entity.related" }

type CodeGraphStateTransitionedContent struct {
	codegraphContent
	EntityType string        `json:"EntityType"`
	EntityID   types.ActorID `json:"EntityID"`
	From       string        `json:"From"`
	To         string        `json:"To"`
}

func (c CodeGraphStateTransitionedContent) EventTypeName() string {
	return "codegraph.state.transitioned"
}

// --- Logic events ---

type CodeGraphTransformAppliedContent struct {
	codegraphContent
	TransformType string `json:"TransformType"`
	InputCount    int    `json:"InputCount"`
	OutputCount   int    `json:"OutputCount"`
}

func (c CodeGraphTransformAppliedContent) EventTypeName() string {
	return "codegraph.logic.transform.applied"
}

type CodeGraphConditionEvaluatedContent struct {
	codegraphContent
	Condition string `json:"Condition"`
	Result    bool   `json:"Result"`
}

func (c CodeGraphConditionEvaluatedContent) EventTypeName() string {
	return "codegraph.logic.condition.evaluated"
}

type CodeGraphSequenceExecutedContent struct {
	codegraphContent
	StepCount     int  `json:"StepCount"`
	StepsComplete int  `json:"StepsComplete"`
	Success       bool `json:"Success"`
}

func (c CodeGraphSequenceExecutedContent) EventTypeName() string {
	return "codegraph.logic.sequence.executed"
}

type CodeGraphLoopIteratedContent struct {
	codegraphContent
	CollectionType string `json:"CollectionType"`
	ItemCount      int    `json:"ItemCount"`
}

func (c CodeGraphLoopIteratedContent) EventTypeName() string {
	return "codegraph.logic.loop.iterated"
}

type CodeGraphTriggerFiredContent struct {
	codegraphContent
	TriggerID   string `json:"TriggerID"`
	CauseEvent  string `json:"CauseEvent"`
}

func (c CodeGraphTriggerFiredContent) EventTypeName() string {
	return "codegraph.logic.trigger.fired"
}

type CodeGraphConstraintViolatedContent struct {
	codegraphContent
	Constraint string `json:"Constraint"`
	EntityType string `json:"EntityType"`
	Field      string `json:"Field"`
	Message    string `json:"Message"`
}

func (c CodeGraphConstraintViolatedContent) EventTypeName() string {
	return "codegraph.logic.constraint.violated"
}

// --- IO events ---

type CodeGraphQueryExecutedContent struct {
	codegraphContent
	EntityType  string `json:"EntityType"`
	ResultCount int    `json:"ResultCount"`
}

func (c CodeGraphQueryExecutedContent) EventTypeName() string {
	return "codegraph.io.query.executed"
}

type CodeGraphCommandExecutedContent struct {
	codegraphContent
	Command    string        `json:"Command"`
	EntityType string        `json:"EntityType"`
	EntityID   types.ActorID `json:"EntityID"`
	Actor      types.ActorID `json:"Actor"`
}

func (c CodeGraphCommandExecutedContent) EventTypeName() string {
	return "codegraph.io.command.executed"
}

type CodeGraphCommandRejectedContent struct {
	codegraphContent
	Command string        `json:"Command"`
	Actor   types.ActorID `json:"Actor"`
	Reason  string        `json:"Reason"`
}

func (c CodeGraphCommandRejectedContent) EventTypeName() string {
	return "codegraph.io.command.rejected"
}

type CodeGraphSubscribeRegisteredContent struct {
	codegraphContent
	EntityType string `json:"EntityType"`
	Filter     string `json:"Filter"`
}

func (c CodeGraphSubscribeRegisteredContent) EventTypeName() string {
	return "codegraph.io.subscribe.registered"
}

type CodeGraphAuthorizeGrantedContent struct {
	codegraphContent
	Actor  types.ActorID `json:"Actor"`
	Action string        `json:"Action"`
	Scope  string        `json:"Scope"`
}

func (c CodeGraphAuthorizeGrantedContent) EventTypeName() string {
	return "codegraph.io.authorize.granted"
}

type CodeGraphAuthorizeDeniedContent struct {
	codegraphContent
	Actor  types.ActorID `json:"Actor"`
	Action string        `json:"Action"`
	Reason string        `json:"Reason"`
}

func (c CodeGraphAuthorizeDeniedContent) EventTypeName() string {
	return "codegraph.io.authorize.denied"
}

type CodeGraphSearchExecutedContent struct {
	codegraphContent
	Query       string `json:"Query"`
	ResultCount int    `json:"ResultCount"`
}

func (c CodeGraphSearchExecutedContent) EventTypeName() string {
	return "codegraph.io.search.executed"
}

type CodeGraphInteropSentContent struct {
	codegraphContent
	Target  string        `json:"Target"`
	Channel string        `json:"Channel"`
	Actor   types.ActorID `json:"Actor"`
}

func (c CodeGraphInteropSentContent) EventTypeName() string { return "codegraph.io.interop.sent" }

type CodeGraphInteropReceivedContent struct {
	codegraphContent
	Source  string `json:"Source"`
	Channel string `json:"Channel"`
}

func (c CodeGraphInteropReceivedContent) EventTypeName() string {
	return "codegraph.io.interop.received"
}

// --- UI events ---

type CodeGraphViewRenderedContent struct {
	codegraphContent
	ViewName string        `json:"ViewName"`
	Actor    types.ActorID `json:"Actor"`
}

func (c CodeGraphViewRenderedContent) EventTypeName() string { return "codegraph.ui.view.rendered" }

type CodeGraphActionInvokedContent struct {
	codegraphContent
	Label    string        `json:"Label"`
	Actor    types.ActorID `json:"Actor"`
	ViewName string        `json:"ViewName"`
}

func (c CodeGraphActionInvokedContent) EventTypeName() string { return "codegraph.ui.action.invoked" }

type CodeGraphNavigationTriggeredContent struct {
	codegraphContent
	Route  string        `json:"Route"`
	Actor  types.ActorID `json:"Actor"`
}

func (c CodeGraphNavigationTriggeredContent) EventTypeName() string {
	return "codegraph.ui.navigation.triggered"
}

type CodeGraphFeedbackEmittedContent struct {
	codegraphContent
	FeedbackType string `json:"FeedbackType"`
	Message      string `json:"Message"`
}

func (c CodeGraphFeedbackEmittedContent) EventTypeName() string {
	return "codegraph.ui.feedback.emitted"
}

type CodeGraphAlertDispatchedContent struct {
	codegraphContent
	Target  types.ActorID `json:"Target"`
	Urgency string        `json:"Urgency"`
	Message string        `json:"Message"`
}

func (c CodeGraphAlertDispatchedContent) EventTypeName() string {
	return "codegraph.ui.alert.dispatched"
}

type CodeGraphDragCompletedContent struct {
	codegraphContent
	SourceType string `json:"SourceType"`
	TargetType string `json:"TargetType"`
	Actor      types.ActorID `json:"Actor"`
}

func (c CodeGraphDragCompletedContent) EventTypeName() string { return "codegraph.ui.drag.completed" }

type CodeGraphSelectionChangedContent struct {
	codegraphContent
	EntityType string `json:"EntityType"`
	Count      int    `json:"Count"`
}

func (c CodeGraphSelectionChangedContent) EventTypeName() string {
	return "codegraph.ui.selection.changed"
}

type CodeGraphConfirmationResolvedContent struct {
	codegraphContent
	Approved bool          `json:"Approved"`
	Actor    types.ActorID `json:"Actor"`
	Message  string        `json:"Message"`
}

func (c CodeGraphConfirmationResolvedContent) EventTypeName() string {
	return "codegraph.ui.confirmation.resolved"
}

// --- Aesthetic events ---

type CodeGraphSkinAppliedContent struct {
	codegraphContent
	SkinName string        `json:"SkinName"`
	Actor    types.ActorID `json:"Actor"`
}

func (c CodeGraphSkinAppliedContent) EventTypeName() string {
	return "codegraph.aesthetic.skin.applied"
}

// --- Temporal events ---

type CodeGraphUndoExecutedContent struct {
	codegraphContent
	TargetEvent types.EventID `json:"TargetEvent"`
	Actor       types.ActorID `json:"Actor"`
}

func (c CodeGraphUndoExecutedContent) EventTypeName() string {
	return "codegraph.temporal.undo.executed"
}

type CodeGraphRetryAttemptedContent struct {
	codegraphContent
	Command     string `json:"Command"`
	Attempt     int    `json:"Attempt"`
	MaxAttempts int    `json:"MaxAttempts"`
}

func (c CodeGraphRetryAttemptedContent) EventTypeName() string {
	return "codegraph.temporal.retry.attempted"
}

// --- Resilience events ---

type CodeGraphOfflineEnteredContent struct {
	codegraphContent
	Actor types.ActorID `json:"Actor"`
}

func (c CodeGraphOfflineEnteredContent) EventTypeName() string {
	return "codegraph.resilience.offline.entered"
}

type CodeGraphOfflineSyncedContent struct {
	codegraphContent
	Actor       types.ActorID `json:"Actor"`
	QueueLength int           `json:"QueueLength"`
}

func (c CodeGraphOfflineSyncedContent) EventTypeName() string {
	return "codegraph.resilience.offline.synced"
}

// --- Structural events ---

type CodeGraphScopeDefinedContent struct {
	codegraphContent
	ScopeType string `json:"ScopeType"`
	Boundary  string `json:"Boundary"`
}

func (c CodeGraphScopeDefinedContent) EventTypeName() string {
	return "codegraph.structural.scope.defined"
}

// --- Social events ---

type CodeGraphPresenceUpdatedContent struct {
	codegraphContent
	Actor    types.ActorID `json:"Actor"`
	ViewName string        `json:"ViewName"`
	Active   bool          `json:"Active"`
}

func (c CodeGraphPresenceUpdatedContent) EventTypeName() string {
	return "codegraph.social.presence.updated"
}

type CodeGraphSalienceTriggeredContent struct {
	codegraphContent
	EntityType string `json:"EntityType"`
	Condition  string `json:"Condition"`
}

func (c CodeGraphSalienceTriggeredContent) EventTypeName() string {
	return "codegraph.social.salience.triggered"
}
