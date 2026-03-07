package compositions

import (
	"context"
	"fmt"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/grammar"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// EvolutionGrammar provides Layer 12 (Emergence) composition operations.
// 10 operations + 4 named functions for system self-awareness and evolution.
type EvolutionGrammar struct {
	g *grammar.Grammar
}

// NewEvolutionGrammar creates an EvolutionGrammar bound to the given base grammar.
func NewEvolutionGrammar(g *grammar.Grammar) *EvolutionGrammar {
	return &EvolutionGrammar{g: g}
}

// --- Operations (10) ---

// DetectPattern finds a pattern in how patterns form. (MetaPattern + Emit)
func (e *EvolutionGrammar) DetectPattern(
	ctx context.Context, source types.ActorID, pattern string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return e.g.Emit(ctx, source, "pattern: "+pattern, convID, causes, signer)
}

// Model maps how components interact to produce behaviour. (SystemDynamic + Emit)
func (e *EvolutionGrammar) Model(
	ctx context.Context, source types.ActorID, model string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return e.g.Emit(ctx, source, "model: "+model, convID, causes, signer)
}

// TraceLoop identifies a feedback cycle. (FeedbackLoop + Emit)
func (e *EvolutionGrammar) TraceLoop(
	ctx context.Context, source types.ActorID, loop string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return e.g.Emit(ctx, source, "loop: "+loop, convID, causes, signer)
}

// WatchThreshold tracks approach to a qualitative transition. (Threshold + Annotate)
func (e *EvolutionGrammar) WatchThreshold(
	ctx context.Context, source types.ActorID, target types.EventID,
	threshold string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return e.g.Annotate(ctx, source, target, "threshold", threshold, convID, signer)
}

// Adapt proposes a structural change. (Adaptation + Emit)
func (e *EvolutionGrammar) Adapt(
	ctx context.Context, source types.ActorID, proposal string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return e.g.Emit(ctx, source, "adapt: "+proposal, convID, causes, signer)
}

// Select tests and keeps or discards an adaptation. (Selection + Emit)
func (e *EvolutionGrammar) Select(
	ctx context.Context, source types.ActorID, result string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return e.g.Emit(ctx, source, "select: "+result, convID, causes, signer)
}

// Simplify removes unnecessary complexity. (Simplification + Emit)
func (e *EvolutionGrammar) Simplify(
	ctx context.Context, source types.ActorID, simplification string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return e.g.Emit(ctx, source, "simplify: "+simplification, convID, causes, signer)
}

// CheckIntegrity assesses structural soundness. (SystemicIntegrity + Emit)
func (e *EvolutionGrammar) CheckIntegrity(
	ctx context.Context, source types.ActorID, assessment string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return e.g.Emit(ctx, source, "integrity: "+assessment, convID, causes, signer)
}

// AssessResilience evaluates ability to absorb shocks. (Resilience + Emit)
func (e *EvolutionGrammar) AssessResilience(
	ctx context.Context, source types.ActorID, assessment string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return e.g.Emit(ctx, source, "resilience: "+assessment, convID, causes, signer)
}

// AlignPurpose verifies alignment with the system's purpose. (Purpose + Emit)
func (e *EvolutionGrammar) AlignPurpose(
	ctx context.Context, source types.ActorID, alignment string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return e.g.Emit(ctx, source, "purpose: "+alignment, convID, causes, signer)
}

// --- Named Functions (4) ---

// SelfEvolveResult holds the events produced by SelfEvolve.
type SelfEvolveResult struct {
	Pattern       event.Event
	Adaptation    event.Event
	Selection     event.Event
	Simplification event.Event
}

// SelfEvolve is a full mechanical-to-intelligent migration:
// DetectPattern + Adapt + Select + Simplify.
func (e *EvolutionGrammar) SelfEvolve(
	ctx context.Context, source types.ActorID,
	pattern string, adaptation string, selection string, simplification string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (SelfEvolveResult, error) {
	pat, err := e.DetectPattern(ctx, source, pattern, causes, convID, signer)
	if err != nil {
		return SelfEvolveResult{}, fmt.Errorf("self-evolve/detect: %w", err)
	}

	adapt, err := e.Adapt(ctx, source, adaptation, []types.EventID{pat.ID()}, convID, signer)
	if err != nil {
		return SelfEvolveResult{}, fmt.Errorf("self-evolve/adapt: %w", err)
	}

	sel, err := e.Select(ctx, source, selection, []types.EventID{adapt.ID()}, convID, signer)
	if err != nil {
		return SelfEvolveResult{}, fmt.Errorf("self-evolve/select: %w", err)
	}

	simp, err := e.Simplify(ctx, source, simplification, []types.EventID{sel.ID()}, convID, signer)
	if err != nil {
		return SelfEvolveResult{}, fmt.Errorf("self-evolve/simplify: %w", err)
	}

	return SelfEvolveResult{
		Pattern:        pat,
		Adaptation:     adapt,
		Selection:      sel,
		Simplification: simp,
	}, nil
}

// HealthCheckResult holds the events produced by a HealthCheck.
type HealthCheckResult struct {
	Integrity  event.Event
	Resilience event.Event
	Model      event.Event
	Purpose    event.Event
}

// HealthCheck is a comprehensive system assessment:
// CheckIntegrity + AssessResilience + Model + AlignPurpose.
func (e *EvolutionGrammar) HealthCheck(
	ctx context.Context, source types.ActorID,
	integrity string, resilience string, model string, purpose string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (HealthCheckResult, error) {
	integ, err := e.CheckIntegrity(ctx, source, integrity, causes, convID, signer)
	if err != nil {
		return HealthCheckResult{}, fmt.Errorf("health-check/integrity: %w", err)
	}

	resil, err := e.AssessResilience(ctx, source, resilience, []types.EventID{integ.ID()}, convID, signer)
	if err != nil {
		return HealthCheckResult{}, fmt.Errorf("health-check/resilience: %w", err)
	}

	mod, err := e.Model(ctx, source, model, []types.EventID{resil.ID()}, convID, signer)
	if err != nil {
		return HealthCheckResult{}, fmt.Errorf("health-check/model: %w", err)
	}

	purp, err := e.AlignPurpose(ctx, source, purpose, []types.EventID{mod.ID()}, convID, signer)
	if err != nil {
		return HealthCheckResult{}, fmt.Errorf("health-check/purpose: %w", err)
	}

	return HealthCheckResult{
		Integrity:  integ,
		Resilience: resil,
		Model:      mod,
		Purpose:    purp,
	}, nil
}
