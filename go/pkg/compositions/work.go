package compositions

import (
	"context"
	"fmt"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/grammar"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// WorkGrammar provides Layer 1 (Agency) composition operations.
// 12 operations + 3 named functions for task management where AI agents
// and humans operate on the same graph.
type WorkGrammar struct {
	g *grammar.Grammar
}

// NewWorkGrammar creates a WorkGrammar bound to the given base grammar.
func NewWorkGrammar(g *grammar.Grammar) *WorkGrammar {
	return &WorkGrammar{g: g}
}

// --- Operations (12) ---

// Intend declares a goal or desired outcome. (Goal + Emit)
func (w *WorkGrammar) Intend(
	ctx context.Context, source types.ActorID, goal string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Emit(ctx, source, "intend: "+goal, convID, causes, signer)
}

// Decompose breaks a goal into actionable steps. (Plan + Derive)
func (w *WorkGrammar) Decompose(
	ctx context.Context, source types.ActorID, subtask string,
	goal types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Derive(ctx, source, "decompose: "+subtask, goal, convID, signer)
}

// Assign gives work to a specific actor. (Delegation + Delegate)
func (w *WorkGrammar) Assign(
	ctx context.Context, source types.ActorID, assignee types.ActorID,
	scope types.DomainScope, weight types.Weight,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Delegate(ctx, source, assignee, scope, weight, cause, convID, signer)
}

// Claim takes on unassigned work. (Initiative + Emit)
func (w *WorkGrammar) Claim(
	ctx context.Context, source types.ActorID, work string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Emit(ctx, source, "claim: "+work, convID, causes, signer)
}

// Prioritize ranks work by importance. (Focus + Annotate)
func (w *WorkGrammar) Prioritize(
	ctx context.Context, source types.ActorID, target types.EventID,
	priority string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Annotate(ctx, source, target, "priority", priority, convID, signer)
}

// Block flags work that cannot proceed. (Salience + Annotate)
func (w *WorkGrammar) Block(
	ctx context.Context, source types.ActorID, target types.EventID,
	blocker string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Annotate(ctx, source, target, "blocked", blocker, convID, signer)
}

// Unblock removes an impediment to work. (Salience + Emit)
func (w *WorkGrammar) Unblock(
	ctx context.Context, source types.ActorID, resolution string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Emit(ctx, source, "unblock: "+resolution, convID, causes, signer)
}

// Progress reports incremental advancement. (Commitment + Extend)
func (w *WorkGrammar) Progress(
	ctx context.Context, source types.ActorID, update string,
	previous types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Extend(ctx, source, "progress: "+update, previous, convID, signer)
}

// Complete marks work as done with evidence. (Commitment + Emit)
func (w *WorkGrammar) Complete(
	ctx context.Context, source types.ActorID, summary string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Emit(ctx, source, "complete: "+summary, convID, causes, signer)
}

// Handoff transfers work between actors. (Delegation + Consent)
func (w *WorkGrammar) Handoff(
	ctx context.Context, from types.ActorID, to types.ActorID,
	description string, scope types.DomainScope,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Consent(ctx, from, to, "handoff: "+description, scope, cause, convID, signer)
}

// Scope defines what an actor may do autonomously. (Permission + Delegate)
func (w *WorkGrammar) Scope(
	ctx context.Context, source types.ActorID, target types.ActorID,
	scope types.DomainScope, weight types.Weight,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Delegate(ctx, source, target, scope, weight, cause, convID, signer)
}

// Review evaluates completed work. (Accountability + Respond)
func (w *WorkGrammar) Review(
	ctx context.Context, source types.ActorID, assessment string,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Respond(ctx, source, "review: "+assessment, target, convID, signer)
}

// --- Named Functions (6) ---

// SprintResult holds the events produced by a Sprint.
type SprintResult struct {
	Intent        event.Event
	Subtasks      []event.Event
	Assignments   []event.Event
}

// Sprint plans a work cycle: Intend + Decompose + Assign (batch).
func (w *WorkGrammar) Sprint(
	ctx context.Context, source types.ActorID, goal string,
	subtasks []string, assignees []types.ActorID, scopes []types.DomainScope,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (SprintResult, error) {
	if len(subtasks) != len(assignees) || len(subtasks) != len(scopes) {
		return SprintResult{}, fmt.Errorf("sprint: subtasks, assignees, and scopes must have equal length")
	}

	intent, err := w.Intend(ctx, source, goal, causes, convID, signer)
	if err != nil {
		return SprintResult{}, fmt.Errorf("sprint/intend: %w", err)
	}

	result := SprintResult{Intent: intent}
	for i, st := range subtasks {
		task, err := w.Decompose(ctx, source, st, intent.ID(), convID, signer)
		if err != nil {
			return SprintResult{}, fmt.Errorf("sprint/decompose[%d]: %w", i, err)
		}
		result.Subtasks = append(result.Subtasks, task)

		assign, err := w.Assign(ctx, source, assignees[i], scopes[i], types.MustWeight(0.5), task.ID(), convID, signer)
		if err != nil {
			return SprintResult{}, fmt.Errorf("sprint/assign[%d]: %w", i, err)
		}
		result.Assignments = append(result.Assignments, assign)
	}
	return result, nil
}

// EscalateResult holds the events produced by an Escalate.
type EscalateResult struct {
	BlockEvent   event.Event
	HandoffEvent event.Event
}

// Escalate moves stuck work up: Block + Handoff (to higher authority).
func (w *WorkGrammar) Escalate(
	ctx context.Context, source types.ActorID, blocker string,
	task types.EventID, authority types.ActorID, scope types.DomainScope,
	convID types.ConversationID, signer event.Signer,
) (EscalateResult, error) {
	block, err := w.Block(ctx, source, task, blocker, convID, signer)
	if err != nil {
		return EscalateResult{}, fmt.Errorf("escalate/block: %w", err)
	}

	handoff, err := w.Handoff(ctx, source, authority, blocker, scope, block.ID(), convID, signer)
	if err != nil {
		return EscalateResult{}, fmt.Errorf("escalate/handoff: %w", err)
	}

	return EscalateResult{BlockEvent: block, HandoffEvent: handoff}, nil
}

// DelegateAndVerifyResult holds the events produced by DelegateAndVerify.
type DelegateAndVerifyResult struct {
	AssignEvent event.Event
	ScopeEvent  event.Event
}

// DelegateAndVerify is a full delegation cycle: Assign + Scope.
// The Review step happens later when work is complete.
func (w *WorkGrammar) DelegateAndVerify(
	ctx context.Context, source types.ActorID, assignee types.ActorID,
	scope types.DomainScope, weight types.Weight,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (DelegateAndVerifyResult, error) {
	assign, err := w.Assign(ctx, source, assignee, scope, weight, cause, convID, signer)
	if err != nil {
		return DelegateAndVerifyResult{}, fmt.Errorf("delegate-and-verify/assign: %w", err)
	}

	scopeEv, err := w.Scope(ctx, source, assignee, scope, weight, assign.ID(), convID, signer)
	if err != nil {
		return DelegateAndVerifyResult{}, fmt.Errorf("delegate-and-verify/scope: %w", err)
	}

	return DelegateAndVerifyResult{AssignEvent: assign, ScopeEvent: scopeEv}, nil
}
