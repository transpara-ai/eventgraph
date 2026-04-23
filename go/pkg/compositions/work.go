package compositions

import (
	"context"
	"fmt"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/grammar"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
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

// Intend declares a desired outcome. (Intent + Emit)
func (w *WorkGrammar) Intend(
	ctx context.Context, source types.ActorID, goal string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Emit(ctx, source, "intend: "+goal, convID, causes, signer)
}

// Decompose breaks intent into actionable steps. (Choice + Derive)
func (w *WorkGrammar) Decompose(
	ctx context.Context, source types.ActorID, subtask string,
	goal types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Derive(ctx, source, "decompose: "+subtask, goal, convID, signer)
}

// Assign gives work to a specific actor. (Commitment + Delegate)
func (w *WorkGrammar) Assign(
	ctx context.Context, source types.ActorID, assignee types.ActorID,
	scope types.DomainScope, weight types.Weight,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Delegate(ctx, source, assignee, scope, weight, cause, convID, signer)
}

// Claim takes on unassigned work. (Intent + Emit)
func (w *WorkGrammar) Claim(
	ctx context.Context, source types.ActorID, work string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Emit(ctx, source, "claim: "+work, convID, causes, signer)
}

// Prioritize ranks work by importance. (Value + Annotate)
func (w *WorkGrammar) Prioritize(
	ctx context.Context, source types.ActorID, target types.EventID,
	priority string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Annotate(ctx, source, target, "priority", priority, convID, signer)
}

// Block flags work that cannot proceed. (Risk + Annotate)
func (w *WorkGrammar) Block(
	ctx context.Context, source types.ActorID, target types.EventID,
	blocker string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Annotate(ctx, source, target, "blocked", blocker, convID, signer)
}

// Unblock removes an impediment to work. (Consequence + Emit)
func (w *WorkGrammar) Unblock(
	ctx context.Context, source types.ActorID, resolution string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Emit(ctx, source, "unblock: "+resolution, convID, causes, signer)
}

// Progress reports incremental advancement. (Act + Extend)
func (w *WorkGrammar) Progress(
	ctx context.Context, source types.ActorID, update string,
	previous types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Extend(ctx, source, "progress: "+update, previous, convID, signer)
}

// Complete marks work as done with evidence. (Consequence + Emit)
func (w *WorkGrammar) Complete(
	ctx context.Context, source types.ActorID, summary string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Emit(ctx, source, "complete: "+summary, convID, causes, signer)
}

// Handoff transfers work between actors. (Signal + Consent)
func (w *WorkGrammar) Handoff(
	ctx context.Context, from types.ActorID, to types.ActorID,
	description string, scope types.DomainScope,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Consent(ctx, from, to, "handoff: "+description, scope, cause, convID, signer)
}

// Scope defines what an actor may do autonomously. (Capacity + Delegate)
func (w *WorkGrammar) Scope(
	ctx context.Context, source types.ActorID, target types.ActorID,
	scope types.DomainScope, weight types.Weight,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Delegate(ctx, source, target, scope, weight, cause, convID, signer)
}

// Review evaluates completed work. (Consequence + Respond)
func (w *WorkGrammar) Review(
	ctx context.Context, source types.ActorID, assessment string,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return w.g.Respond(ctx, source, "review: "+assessment, target, convID, signer)
}

// --- Named Functions (6) ---

// StandupResult holds the events produced by a Standup.
type StandupResult struct {
	Updates   []event.Event
	Priority  event.Event
}

// Standup gathers status from all participants and sets priority:
// Progress (batch) + Prioritize.
func (w *WorkGrammar) Standup(
	ctx context.Context, participants []types.ActorID, updates []string,
	lead types.ActorID, priority string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (StandupResult, error) {
	if len(participants) != len(updates) {
		return StandupResult{}, fmt.Errorf("standup: participants and updates must have equal length")
	}

	result := StandupResult{}
	var lastID types.EventID
	for i, actor := range participants {
		var prev types.EventID
		if i == 0 {
			if len(causes) > 0 {
				prev = causes[0]
			}
		} else {
			prev = lastID
		}
		progress, err := w.Progress(ctx, actor, updates[i], prev, convID, signer)
		if err != nil {
			return StandupResult{}, fmt.Errorf("standup/progress[%d]: %w", i, err)
		}
		result.Updates = append(result.Updates, progress)
		lastID = progress.ID()
	}

	prio, err := w.Prioritize(ctx, lead, lastID, priority, convID, signer)
	if err != nil {
		return StandupResult{}, fmt.Errorf("standup/prioritize: %w", err)
	}
	result.Priority = prio

	return result, nil
}

// RetrospectiveResult holds the events produced by a Retrospective.
type RetrospectiveResult struct {
	Reviews     []event.Event
	Improvement event.Event
}

// Retrospective reviews work and identifies improvements: Review (batch) + Intend.
func (w *WorkGrammar) Retrospective(
	ctx context.Context, reviewers []types.ActorID, assessments []string,
	lead types.ActorID, improvement string,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (RetrospectiveResult, error) {
	if len(reviewers) != len(assessments) {
		return RetrospectiveResult{}, fmt.Errorf("retrospective: reviewers and assessments must have equal length")
	}

	result := RetrospectiveResult{}
	reviewIDs := make([]types.EventID, 0, len(reviewers))
	for i, reviewer := range reviewers {
		rev, err := w.Review(ctx, reviewer, assessments[i], target, convID, signer)
		if err != nil {
			return RetrospectiveResult{}, fmt.Errorf("retrospective/review[%d]: %w", i, err)
		}
		result.Reviews = append(result.Reviews, rev)
		reviewIDs = append(reviewIDs, rev.ID())
	}

	improve, err := w.Intend(ctx, lead, improvement, reviewIDs, convID, signer)
	if err != nil {
		return RetrospectiveResult{}, fmt.Errorf("retrospective/intend: %w", err)
	}
	result.Improvement = improve

	return result, nil
}

// TriageResult holds the events produced by a Triage.
type TriageResult struct {
	Priorities  []event.Event
	Assignments []event.Event
	Scopes      []event.Event
}

// Triage prioritises and assigns a batch of items: Prioritize + Assign + Scope (batch).
func (w *WorkGrammar) Triage(
	ctx context.Context, lead types.ActorID,
	items []types.EventID, priorities []string,
	assignees []types.ActorID, scopes []types.DomainScope, weights []types.Weight,
	convID types.ConversationID, signer event.Signer,
) (TriageResult, error) {
	n := len(items)
	if len(priorities) != n || len(assignees) != n || len(scopes) != n || len(weights) != n {
		return TriageResult{}, fmt.Errorf("triage: all slices must have equal length")
	}

	result := TriageResult{}
	for i := range items {
		prio, err := w.Prioritize(ctx, lead, items[i], priorities[i], convID, signer)
		if err != nil {
			return TriageResult{}, fmt.Errorf("triage/prioritize[%d]: %w", i, err)
		}
		result.Priorities = append(result.Priorities, prio)

		assign, err := w.Assign(ctx, lead, assignees[i], scopes[i], weights[i], prio.ID(), convID, signer)
		if err != nil {
			return TriageResult{}, fmt.Errorf("triage/assign[%d]: %w", i, err)
		}
		result.Assignments = append(result.Assignments, assign)

		scope, err := w.Scope(ctx, lead, assignees[i], scopes[i], weights[i], assign.ID(), convID, signer)
		if err != nil {
			return TriageResult{}, fmt.Errorf("triage/scope[%d]: %w", i, err)
		}
		result.Scopes = append(result.Scopes, scope)
	}

	return result, nil
}

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
