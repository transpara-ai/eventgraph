package compositions

import (
	"context"
	"fmt"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/grammar"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// AlignmentGrammar provides Layer 7 (Ethics) composition operations.
// 10 operations + 2 named functions for AI accountability.
type AlignmentGrammar struct {
	g *grammar.Grammar
}

// NewAlignmentGrammar creates an AlignmentGrammar bound to the given base grammar.
func NewAlignmentGrammar(g *grammar.Grammar) *AlignmentGrammar {
	return &AlignmentGrammar{g: g}
}

// --- Operations (10) ---

// Constrain sets an ethical boundary on an existing action or system. (Value + Annotate)
func (a *AlignmentGrammar) Constrain(
	ctx context.Context, source types.ActorID,
	target types.EventID, constraint string,
	convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Annotate(ctx, source, target, "constraint", constraint, convID, signer)
}

// DetectHarm identifies harm from an action or pattern. (Harm + Emit)
func (a *AlignmentGrammar) DetectHarm(
	ctx context.Context, source types.ActorID, harm string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Emit(ctx, source, "harm: "+harm, convID, causes, signer)
}

// AssessFairness evaluates equitable treatment on a target. (Fairness + Annotate)
func (a *AlignmentGrammar) AssessFairness(
	ctx context.Context, source types.ActorID,
	target types.EventID, assessment string,
	convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Annotate(ctx, source, target, "fairness", assessment, convID, signer)
}

// FlagDilemma identifies a situation where values conflict. (Dilemma + Emit)
func (a *AlignmentGrammar) FlagDilemma(
	ctx context.Context, source types.ActorID, dilemma string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Emit(ctx, source, "dilemma: "+dilemma, convID, causes, signer)
}

// Weigh balances competing values for a decision. (Proportionality + Derive)
func (a *AlignmentGrammar) Weigh(
	ctx context.Context, source types.ActorID, weighing string,
	decision types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Derive(ctx, source, "weigh: "+weighing, decision, convID, signer)
}

// Explain makes reasoning visible and accessible. (Transparency + Emit)
func (a *AlignmentGrammar) Explain(
	ctx context.Context, source types.ActorID, explanation string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Emit(ctx, source, "explain: "+explanation, convID, causes, signer)
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

// Grow updates ethical reasoning from experience. (Growth + Emit)
func (a *AlignmentGrammar) Grow(
	ctx context.Context, source types.ActorID, growth string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return a.g.Emit(ctx, source, "grow: "+growth, convID, causes, signer)
}

// --- Named Functions (5) ---

// EthicsAuditResult holds the events produced by an EthicsAudit.
type EthicsAuditResult struct {
	Fairness event.Event
	HarmScan event.Event
	Report   event.Event
}

// EthicsAudit performs a comprehensive ethical review: AssessFairness + DetectHarm + Explain.
// The fairness assessment annotates the target; the harm scan derives from it;
// the report explains the combined findings.
func (a *AlignmentGrammar) EthicsAudit(
	ctx context.Context, auditor types.ActorID,
	target types.EventID,
	fairnessAssessment string, harmScan string, summary string,
	convID types.ConversationID, signer event.Signer,
) (EthicsAuditResult, error) {
	fairness, err := a.AssessFairness(ctx, auditor, target, fairnessAssessment, convID, signer)
	if err != nil {
		return EthicsAuditResult{}, fmt.Errorf("ethics-audit/fairness: %w", err)
	}

	harm, err := a.DetectHarm(ctx, auditor, harmScan, []types.EventID{fairness.ID()}, convID, signer)
	if err != nil {
		return EthicsAuditResult{}, fmt.Errorf("ethics-audit/harm: %w", err)
	}

	report, err := a.Explain(ctx, auditor, summary,
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
	harmEv, err := a.DetectHarm(ctx, auditor, harm, []types.EventID{cause}, convID, signer)
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

	growEv, err := a.Grow(ctx, agent, growth, []types.EventID{repair.ID()}, convID, signer)
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

// GuardrailResult holds the events produced by a Guardrail.
type GuardrailResult struct {
	Constraint event.Event
	Dilemma    event.Event
	Escalation event.Event
}

// Guardrail sets a boundary with trigger-based dilemma detection:
// Constrain + FlagDilemma + Escalate (via Emit to authority).
func (a *AlignmentGrammar) Guardrail(
	ctx context.Context, source types.ActorID,
	target types.EventID, constraint string,
	dilemma string, escalation string,
	convID types.ConversationID, signer event.Signer,
) (GuardrailResult, error) {
	constrain, err := a.Constrain(ctx, source, target, constraint, convID, signer)
	if err != nil {
		return GuardrailResult{}, fmt.Errorf("guardrail/constrain: %w", err)
	}

	dilemmaEv, err := a.FlagDilemma(ctx, source, dilemma, []types.EventID{constrain.ID()}, convID, signer)
	if err != nil {
		return GuardrailResult{}, fmt.Errorf("guardrail/dilemma: %w", err)
	}

	escalate, err := a.g.Emit(ctx, source, "escalate: "+escalation, convID, []types.EventID{dilemmaEv.ID()}, signer)
	if err != nil {
		return GuardrailResult{}, fmt.Errorf("guardrail/escalate: %w", err)
	}

	return GuardrailResult{Constraint: constrain, Dilemma: dilemmaEv, Escalation: escalate}, nil
}

// ImpactAssessmentResult holds the events produced by an ImpactAssessment.
type ImpactAssessmentResult struct {
	Weighing    event.Event
	Fairness    event.Event
	Explanation event.Event
}

// ImpactAssessment evaluates prospective consequences:
// Weigh + AssessFairness + Explain.
func (a *AlignmentGrammar) ImpactAssessment(
	ctx context.Context, source types.ActorID,
	decision types.EventID, weighing string, fairness string, explanation string,
	convID types.ConversationID, signer event.Signer,
) (ImpactAssessmentResult, error) {
	weigh, err := a.Weigh(ctx, source, weighing, decision, convID, signer)
	if err != nil {
		return ImpactAssessmentResult{}, fmt.Errorf("impact/weigh: %w", err)
	}

	fair, err := a.AssessFairness(ctx, source, weigh.ID(), fairness, convID, signer)
	if err != nil {
		return ImpactAssessmentResult{}, fmt.Errorf("impact/fairness: %w", err)
	}

	explain, err := a.Explain(ctx, source, explanation, []types.EventID{weigh.ID(), fair.ID()}, convID, signer)
	if err != nil {
		return ImpactAssessmentResult{}, fmt.Errorf("impact/explain: %w", err)
	}

	return ImpactAssessmentResult{Weighing: weigh, Fairness: fair, Explanation: explain}, nil
}

// WhistleblowResult holds the events produced by a Whistleblow.
type WhistleblowResult struct {
	Harm        event.Event
	Explanation event.Event
	Escalation  event.Event
}

// Whistleblow reports harm to external authority:
// DetectHarm + Explain + Escalate (to external).
func (a *AlignmentGrammar) Whistleblow(
	ctx context.Context, source types.ActorID,
	harm string, explanation string, escalation string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (WhistleblowResult, error) {
	harmEv, err := a.DetectHarm(ctx, source, harm, causes, convID, signer)
	if err != nil {
		return WhistleblowResult{}, fmt.Errorf("whistleblow/harm: %w", err)
	}

	explain, err := a.Explain(ctx, source, explanation, []types.EventID{harmEv.ID()}, convID, signer)
	if err != nil {
		return WhistleblowResult{}, fmt.Errorf("whistleblow/explain: %w", err)
	}

	escalate, err := a.g.Emit(ctx, source, "escalate-external: "+escalation, convID, []types.EventID{explain.ID()}, signer)
	if err != nil {
		return WhistleblowResult{}, fmt.Errorf("whistleblow/escalate: %w", err)
	}

	return WhistleblowResult{Harm: harmEv, Explanation: explain, Escalation: escalate}, nil
}
