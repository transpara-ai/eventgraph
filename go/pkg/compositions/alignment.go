package compositions

import (
	"context"
	"fmt"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/grammar"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// AlignmentGrammar provides Layer 7 (Ethics) composition operations.
// 10 operations + 5 named functions for AI accountability.
type AlignmentGrammar struct {
	g *grammar.Grammar
}

// NewAlignmentGrammar creates an AlignmentGrammar bound to the given base grammar.
func NewAlignmentGrammar(g *grammar.Grammar) *AlignmentGrammar {
	return &AlignmentGrammar{g: g}
}

// --- Operations (10) ---

// Constrain sets an ethical boundary on future actions. (Value + Annotate)
func (a *AlignmentGrammar) Constrain(
	ctx context.Context, source types.ActorID, constraint string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Emit(ctx, source, "constrain: "+constraint, convID, causes, signer)
}

// DetectHarm identifies harm from an action or pattern. (Harm + Emit)
func (a *AlignmentGrammar) DetectHarm(
	ctx context.Context, source types.ActorID, harm string,
	action types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Derive(ctx, source, "harm: "+harm, action, convID, signer)
}

// AssessFairness evaluates equitable treatment across groups. (Fairness + Annotate)
func (a *AlignmentGrammar) AssessFairness(
	ctx context.Context, source types.ActorID, assessment string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Emit(ctx, source, "fairness: "+assessment, convID, causes, signer)
}

// FlagDilemma identifies a situation where values conflict. (Dilemma + Emit)
func (a *AlignmentGrammar) FlagDilemma(
	ctx context.Context, source types.ActorID, dilemma string,
	situation types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Derive(ctx, source, "dilemma: "+dilemma, situation, convID, signer)
}

// Weigh balances competing values for a decision. (Proportionality + Derive)
func (a *AlignmentGrammar) Weigh(
	ctx context.Context, source types.ActorID, weighing string,
	decision types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Derive(ctx, source, "weigh: "+weighing, decision, convID, signer)
}

// Explain makes reasoning visible and accessible. (Transparency + Derive)
func (a *AlignmentGrammar) Explain(
	ctx context.Context, source types.ActorID, explanation string,
	reasoning types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Derive(ctx, source, "explain: "+explanation, reasoning, convID, signer)
}

// Assign determines moral responsibility. (Responsibility + Annotate)
func (a *AlignmentGrammar) Assign(
	ctx context.Context, source types.ActorID, target types.EventID,
	responsibility string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Annotate(ctx, source, target, "responsibility", responsibility, convID, signer)
}

// Repair proposes and executes redress for harm. (Redress + Consent)
func (a *AlignmentGrammar) Repair(
	ctx context.Context, source types.ActorID, affected types.ActorID,
	redress string, scope types.DomainScope,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Consent(ctx, source, affected, "repair: "+redress, scope, cause, convID, signer)
}

// Care prioritizes wellbeing of an actor. (Care + Emit)
func (a *AlignmentGrammar) Care(
	ctx context.Context, source types.ActorID, care string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Emit(ctx, source, "care: "+care, convID, causes, signer)
}

// Grow updates ethical reasoning from experience. (Growth + Extend)
func (a *AlignmentGrammar) Grow(
	ctx context.Context, source types.ActorID, growth string,
	previous types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Extend(ctx, source, "grow: "+growth, previous, convID, signer)
}

// --- Named Functions (5) ---

// EthicsAuditResult holds the events produced by an EthicsAudit.
type EthicsAuditResult struct {
	Fairness event.Event
	HarmScan event.Event
	Report   event.Event
}

// EthicsAudit performs a comprehensive ethical review: AssessFairness + DetectHarm + Explain.
func (a *AlignmentGrammar) EthicsAudit(
	ctx context.Context, auditor types.ActorID,
	fairnessAssessment string, harmScan string, summary string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (EthicsAuditResult, error) {
	fairness, err := a.AssessFairness(ctx, auditor, fairnessAssessment, causes, convID, signer)
	if err != nil {
		return EthicsAuditResult{}, fmt.Errorf("ethics-audit/fairness: %w", err)
	}

	harm, err := a.AssessFairness(ctx, auditor, harmScan, causes, convID, signer)
	if err != nil {
		return EthicsAuditResult{}, fmt.Errorf("ethics-audit/harm: %w", err)
	}

	report, err := a.g.Merge(ctx, auditor, "ethics-audit: "+summary,
		[]types.EventID{fairness.ID(), harm.ID()}, convID, signer)
	if err != nil {
		return EthicsAuditResult{}, fmt.Errorf("ethics-audit/report: %w", err)
	}

	return EthicsAuditResult{Fairness: fairness, HarmScan: harm, Report: report}, nil
}

// RestorativeJusticeResult holds the events produced by RestorativeJustice.
type RestorativeJusticeResult struct {
	HarmDetection  event.Event
	Responsibility event.Event
	Redress        event.Event
	Growth         event.Event
}

// RestorativeJustice runs a full accountability-to-healing cycle:
// DetectHarm + Assign + Repair + Grow.
func (a *AlignmentGrammar) RestorativeJustice(
	ctx context.Context, auditor types.ActorID, agent types.ActorID,
	affected types.ActorID,
	harm string, responsibility string, redress string, growth string,
	scope types.DomainScope,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (RestorativeJusticeResult, error) {
	harmEv, err := a.DetectHarm(ctx, auditor, harm, cause, convID, signer)
	if err != nil {
		return RestorativeJusticeResult{}, fmt.Errorf("restorative/harm: %w", err)
	}

	assign, err := a.Assign(ctx, auditor, harmEv.ID(), responsibility, convID, signer)
	if err != nil {
		return RestorativeJusticeResult{}, fmt.Errorf("restorative/assign: %w", err)
	}

	repair, err := a.Repair(ctx, auditor, affected, redress, scope, assign.ID(), convID, signer)
	if err != nil {
		return RestorativeJusticeResult{}, fmt.Errorf("restorative/repair: %w", err)
	}

	growEv, err := a.Grow(ctx, agent, growth, repair.ID(), convID, signer)
	if err != nil {
		return RestorativeJusticeResult{}, fmt.Errorf("restorative/grow: %w", err)
	}

	return RestorativeJusticeResult{
		HarmDetection:  harmEv,
		Responsibility: assign,
		Redress:        repair,
		Growth:         growEv,
	}, nil
}
