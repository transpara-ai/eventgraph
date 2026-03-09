// Package codegraph implements the 65 Code Graph primitives — semantic atoms
// for describing any application. Data, Logic, IO, UI, Aesthetic, Accessibility,
// Temporal, Resilience, Structural, and Social categories.
// Spec: docs/codegraph-spec.md
package codegraph

import (
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

var layer5 = types.MustLayer(5)
var cadence1 = types.MustCadence(1)

// cgPrimitive is the common implementation for all code graph primitives.
type cgPrimitive struct {
	id   types.PrimitiveID
	subs []types.SubscriptionPattern
}

func (p *cgPrimitive) ID() types.PrimitiveID                   { return p.id }
func (p *cgPrimitive) Layer() types.Layer                       { return layer5 }
func (p *cgPrimitive) Lifecycle() types.LifecycleState          { return types.LifecycleActive }
func (p *cgPrimitive) Cadence() types.Cadence                   { return cadence1 }
func (p *cgPrimitive) Subscriptions() []types.SubscriptionPattern { return p.subs }

func (p *cgPrimitive) Process(tick types.Tick, events []event.Event, _ primitive.Snapshot) ([]primitive.Mutation, error) {
	return []primitive.Mutation{
		primitive.UpdateState{PrimitiveID: p.id, Key: "eventsProcessed", Value: len(events)},
		primitive.UpdateState{PrimitiveID: p.id, Key: "lastTick", Value: tick.Value()},
	}, nil
}

func newCG(id string, subs ...string) *cgPrimitive {
	patterns := make([]types.SubscriptionPattern, len(subs))
	for i, s := range subs {
		patterns[i] = types.MustSubscriptionPattern(s)
	}
	return &cgPrimitive{id: types.MustPrimitiveID(id), subs: patterns}
}

// ════════════════════════════════════════════════════════════════════════
// Category 1: Data (6)
// ════════════════════════════════════════════════════════════════════════

func NewEntityPrimitive() primitive.Primitive     { return newCG("CGEntity", "codegraph.entity.*", "codegraph.io.command.*") }
func NewPropertyPrimitive() primitive.Primitive   { return newCG("CGProperty", "codegraph.entity.*") }
func NewRelationPrimitive() primitive.Primitive   { return newCG("CGRelation", "codegraph.entity.*") }
func NewCollectionPrimitive() primitive.Primitive { return newCG("CGCollection", "codegraph.entity.*", "codegraph.io.query.*") }
func NewStatePrimitive() primitive.Primitive      { return newCG("CGState", "codegraph.state.*", "codegraph.io.command.*") }
func NewEventPrimitive() primitive.Primitive      { return newCG("CGEvent", "codegraph.*") }

// ════════════════════════════════════════════════════════════════════════
// Category 2: Logic (6)
// ════════════════════════════════════════════════════════════════════════

func NewTransformPrimitive() primitive.Primitive  { return newCG("CGTransform", "codegraph.logic.transform.*") }
func NewConditionPrimitive() primitive.Primitive  { return newCG("CGCondition", "codegraph.logic.condition.*") }
func NewSequencePrimitive() primitive.Primitive   { return newCG("CGSequence", "codegraph.logic.sequence.*") }
func NewLoopPrimitive() primitive.Primitive       { return newCG("CGLoop", "codegraph.logic.loop.*") }
func NewTriggerPrimitive() primitive.Primitive    { return newCG("CGTrigger", "codegraph.logic.trigger.*", "codegraph.entity.*") }
func NewConstraintPrimitive() primitive.Primitive { return newCG("CGConstraint", "codegraph.logic.constraint.*", "codegraph.io.command.*") }

// ════════════════════════════════════════════════════════════════════════
// Category 3: IO (6)
// ════════════════════════════════════════════════════════════════════════

func NewQueryPrimitive() primitive.Primitive     { return newCG("CGQuery", "codegraph.io.query.*") }
func NewCommandPrimitive() primitive.Primitive   { return newCG("CGCommand", "codegraph.io.command.*") }
func NewSubscribePrimitive() primitive.Primitive { return newCG("CGSubscribe", "codegraph.io.subscribe.*") }
func NewAuthorizePrimitive() primitive.Primitive { return newCG("CGAuthorize", "codegraph.io.authorize.*", "authority.*") }
func NewSearchPrimitive() primitive.Primitive    { return newCG("CGSearch", "codegraph.io.search.*") }
func NewInteropPrimitive() primitive.Primitive   { return newCG("CGInterop", "codegraph.io.interop.*") }

// ════════════════════════════════════════════════════════════════════════
// Category 4: UI (19)
// ════════════════════════════════════════════════════════════════════════

func NewDisplayPrimitive() primitive.Primitive      { return newCG("CGDisplay", "codegraph.ui.*") }
func NewInputPrimitive() primitive.Primitive        { return newCG("CGInput", "codegraph.ui.*") }
func NewLayoutPrimitive() primitive.Primitive       { return newCG("CGLayout", "codegraph.ui.*") }
func NewListPrimitive() primitive.Primitive         { return newCG("CGList", "codegraph.ui.*", "codegraph.io.query.*") }
func NewFormPrimitive() primitive.Primitive         { return newCG("CGForm", "codegraph.ui.*", "codegraph.io.command.*") }
func NewActionPrimitive() primitive.Primitive       { return newCG("CGAction", "codegraph.ui.action.*") }
func NewNavigationPrimitive() primitive.Primitive   { return newCG("CGNavigation", "codegraph.ui.navigation.*") }
func NewViewPrimitive() primitive.Primitive         { return newCG("CGView", "codegraph.ui.view.*") }
func NewFeedbackPrimitive() primitive.Primitive     { return newCG("CGFeedback", "codegraph.ui.feedback.*") }
func NewAlertPrimitive() primitive.Primitive        { return newCG("CGAlert", "codegraph.ui.alert.*") }
func NewThreadPrimitive() primitive.Primitive       { return newCG("CGThread", "codegraph.ui.*", "codegraph.entity.*") }
func NewAvatarPrimitive() primitive.Primitive       { return newCG("CGAvatar", "codegraph.ui.*") }
func NewAuditPrimitive() primitive.Primitive        { return newCG("CGAudit", "codegraph.*") }
func NewDragPrimitive() primitive.Primitive         { return newCG("CGDrag", "codegraph.ui.drag.*") }
func NewSelectionPrimitive() primitive.Primitive    { return newCG("CGSelection", "codegraph.ui.selection.*") }
func NewConfirmationPrimitive() primitive.Primitive { return newCG("CGConfirmation", "codegraph.ui.confirmation.*") }
func NewEmptyPrimitive() primitive.Primitive        { return newCG("CGEmpty", "codegraph.ui.*") }
func NewLoadingPrimitive() primitive.Primitive      { return newCG("CGLoading", "codegraph.ui.*") }
func NewPaginationPrimitive() primitive.Primitive   { return newCG("CGPagination", "codegraph.ui.*", "codegraph.io.query.*") }

// ════════════════════════════════════════════════════════════════════════
// Category 5: Aesthetic (7)
// ════════════════════════════════════════════════════════════════════════

func NewPalettePrimitive() primitive.Primitive    { return newCG("CGPalette", "codegraph.aesthetic.*") }
func NewTypographyPrimitive() primitive.Primitive { return newCG("CGTypography", "codegraph.aesthetic.*") }
func NewSpacingPrimitive() primitive.Primitive    { return newCG("CGSpacing", "codegraph.aesthetic.*") }
func NewElevationPrimitive() primitive.Primitive  { return newCG("CGElevation", "codegraph.aesthetic.*") }
func NewMotionPrimitive() primitive.Primitive     { return newCG("CGMotion", "codegraph.aesthetic.*") }
func NewDensityPrimitive() primitive.Primitive    { return newCG("CGDensity", "codegraph.aesthetic.*") }
func NewShapePrimitive() primitive.Primitive      { return newCG("CGShape", "codegraph.aesthetic.*") }

// ════════════════════════════════════════════════════════════════════════
// Category 6: Accessibility (4)
// ════════════════════════════════════════════════════════════════════════

func NewAnnouncePrimitive() primitive.Primitive { return newCG("CGAnnounce", "codegraph.ui.*", "codegraph.aesthetic.*") }
func NewFocusPrimitive() primitive.Primitive    { return newCG("CGFocus", "codegraph.ui.*") }
func NewContrastPrimitive() primitive.Primitive { return newCG("CGContrast", "codegraph.aesthetic.*") }
func NewSimplifyPrimitive() primitive.Primitive { return newCG("CGSimplify", "codegraph.ui.*", "codegraph.aesthetic.*") }

// ════════════════════════════════════════════════════════════════════════
// Category 7: Temporal (3)
// ════════════════════════════════════════════════════════════════════════

func NewRecencyPrimitive() primitive.Primitive  { return newCG("CGRecency", "codegraph.entity.*", "codegraph.temporal.*") }
func NewHistoryPrimitive() primitive.Primitive  { return newCG("CGHistory", "codegraph.entity.*", "codegraph.temporal.*") }
func NewLivenessPrimitive() primitive.Primitive { return newCG("CGLiveness", "codegraph.io.subscribe.*", "codegraph.social.presence.*") }

// ════════════════════════════════════════════════════════════════════════
// Category 8: Resilience (4)
// ════════════════════════════════════════════════════════════════════════

func NewUndoPrimitive() primitive.Primitive    { return newCG("CGUndo", "codegraph.temporal.undo.*", "codegraph.io.command.*") }
func NewRetryPrimitive() primitive.Primitive   { return newCG("CGRetry", "codegraph.temporal.retry.*", "codegraph.io.command.*") }
func NewFallbackPrimitive() primitive.Primitive { return newCG("CGFallback", "codegraph.resilience.*", "codegraph.ui.*") }
func NewOfflinePrimitive() primitive.Primitive  { return newCG("CGOffline", "codegraph.resilience.offline.*") }

// ════════════════════════════════════════════════════════════════════════
// Category 9: Structural (3)
// ════════════════════════════════════════════════════════════════════════

func NewScopePrimitive() primitive.Primitive   { return newCG("CGScope", "codegraph.structural.*") }
func NewFormatPrimitive() primitive.Primitive  { return newCG("CGFormat", "codegraph.io.*", "codegraph.structural.*") }
func NewGesturePrimitive() primitive.Primitive { return newCG("CGGesture", "codegraph.ui.*") }

// ════════════════════════════════════════════════════════════════════════
// Category 10: Social (3)
// ════════════════════════════════════════════════════════════════════════

func NewPresencePrimitive() primitive.Primitive           { return newCG("CGPresence", "codegraph.social.presence.*") }
func NewSaliencePrimitive() primitive.Primitive           { return newCG("CGSalience", "codegraph.social.salience.*", "codegraph.entity.*") }
func NewConsequencePreviewPrimitive() primitive.Primitive { return newCG("CGConsequencePreview", "codegraph.ui.confirmation.*", "codegraph.io.command.*") }

// AllPrimitives returns all 65 code graph primitives.
func AllPrimitives() []primitive.Primitive {
	return []primitive.Primitive{
		// Data
		NewEntityPrimitive(), NewPropertyPrimitive(), NewRelationPrimitive(),
		NewCollectionPrimitive(), NewStatePrimitive(), NewEventPrimitive(),
		// Logic
		NewTransformPrimitive(), NewConditionPrimitive(), NewSequencePrimitive(),
		NewLoopPrimitive(), NewTriggerPrimitive(), NewConstraintPrimitive(),
		// IO
		NewQueryPrimitive(), NewCommandPrimitive(), NewSubscribePrimitive(),
		NewAuthorizePrimitive(), NewSearchPrimitive(), NewInteropPrimitive(),
		// UI
		NewDisplayPrimitive(), NewInputPrimitive(), NewLayoutPrimitive(),
		NewListPrimitive(), NewFormPrimitive(), NewActionPrimitive(),
		NewNavigationPrimitive(), NewViewPrimitive(), NewFeedbackPrimitive(),
		NewAlertPrimitive(), NewThreadPrimitive(), NewAvatarPrimitive(),
		NewAuditPrimitive(), NewDragPrimitive(), NewSelectionPrimitive(),
		NewConfirmationPrimitive(), NewEmptyPrimitive(), NewLoadingPrimitive(),
		NewPaginationPrimitive(),
		// Aesthetic
		NewPalettePrimitive(), NewTypographyPrimitive(), NewSpacingPrimitive(),
		NewElevationPrimitive(), NewMotionPrimitive(), NewDensityPrimitive(),
		NewShapePrimitive(),
		// Accessibility
		NewAnnouncePrimitive(), NewFocusPrimitive(), NewContrastPrimitive(),
		NewSimplifyPrimitive(),
		// Temporal
		NewRecencyPrimitive(), NewHistoryPrimitive(), NewLivenessPrimitive(),
		// Resilience
		NewUndoPrimitive(), NewRetryPrimitive(), NewFallbackPrimitive(),
		NewOfflinePrimitive(),
		// Structural
		NewScopePrimitive(), NewFormatPrimitive(), NewGesturePrimitive(),
		// Social
		NewPresencePrimitive(), NewSaliencePrimitive(), NewConsequencePreviewPrimitive(),
	}
}
