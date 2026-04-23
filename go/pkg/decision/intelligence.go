package decision

import (
	"context"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// TokenUsage tracks per-call token consumption with full breakdown.
type TokenUsage struct {
	InputTokens       int     // prompt tokens sent
	OutputTokens      int     // completion tokens received
	CacheReadTokens   int     // tokens served from cache (cheaper)
	CacheWriteTokens  int     // tokens written to cache
	CostUSD           float64 // total cost for this call
}

// Total returns the sum of input and output tokens.
func (u TokenUsage) Total() int {
	return u.InputTokens + u.OutputTokens
}

// Response is the result of an IIntelligence.Reason call.
type Response struct {
	content    string
	confidence types.Score
	usage      TokenUsage
}

// NewResponse creates a Response with full token usage.
func NewResponse(content string, confidence types.Score, usage TokenUsage) Response {
	return Response{content: content, confidence: confidence, usage: usage}
}

func (r Response) Content() string        { return r.content }
func (r Response) Confidence() types.Score { return r.confidence }
func (r Response) TokensUsed() int         { return r.usage.Total() }
func (r Response) Usage() TokenUsage       { return r.usage }

// IIntelligence is anything that reasons. Not every primitive needs this.
type IIntelligence interface {
	Reason(ctx context.Context, prompt string, history []event.Event) (Response, error)
}

// OperateTask describes an agentic task with filesystem access.
type OperateTask struct {
	WorkDir      string   // directory the agent operates in
	Instruction  string   // what to do (natural language)
	AllowedTools []string // which tools to grant (Read, Edit, Write, Bash, etc.)
}

// OperateResult is the outcome of an agentic filesystem operation.
type OperateResult struct {
	Summary string     // text summary of what was done
	Usage   TokenUsage // token consumption
}

// IOperator extends IIntelligence with agentic filesystem operations.
// Only providers that can interact with the filesystem implement this
// (e.g., Claude CLI with tool access). The pipeline does a type assertion
// to check if a provider supports this.
type IOperator interface {
	Operate(ctx context.Context, task OperateTask) (OperateResult, error)
}

// IDecisionMaker is anything that makes decisions.
// An AI agent, a human with a UI, a committee vote, a rules engine.
type IDecisionMaker interface {
	Decide(ctx context.Context, input event.DecisionInput) (event.Decision, error)
}

// NoOpIntelligence is a mechanical-only intelligence that always returns an error.
// Used when no LLM is configured (FallbackToMechanical mode).
type NoOpIntelligence struct{}

func (NoOpIntelligence) Reason(_ context.Context, _ string, _ []event.Event) (Response, error) {
	return Response{}, &IntelligenceUnavailableError{}
}

// IntelligenceUnavailableError is returned when IIntelligence is needed but not available.
type IntelligenceUnavailableError struct{}

func (e *IntelligenceUnavailableError) Error() string {
	return "intelligence unavailable: no IIntelligence configured"
}

// DecisionError is the marker interface for decision-domain errors.
type DecisionError interface {
	error
	decisionError()
}

// ActorNotFoundError is returned when a referenced actor doesn't exist.
type ActorNotFoundError struct {
	ID types.ActorID
}

func (e *ActorNotFoundError) Error() string {
	return "actor not found: " + e.ID.Value()
}
func (e *ActorNotFoundError) decisionError() {}

// InsufficientAuthorityError is returned when an actor lacks authority for an action.
type InsufficientAuthorityError struct {
	Actor    types.ActorID
	Action   string
	Required event.AuthorityLevel
}

func (e *InsufficientAuthorityError) Error() string {
	return "insufficient authority: " + e.Actor.Value() + " cannot " + e.Action + " (requires " + string(e.Required) + ")"
}
func (e *InsufficientAuthorityError) decisionError() {}

// TrustBelowThresholdError is returned when trust is too low for an action.
type TrustBelowThresholdError struct {
	Actor    types.ActorID
	Score    types.Score
	Required types.Score
}

func (e *TrustBelowThresholdError) Error() string {
	return "trust below threshold for actor: " + e.Actor.Value()
}
func (e *TrustBelowThresholdError) decisionError() {}

// CausesRequiredError is returned when a decision must reference causing events.
type CausesRequiredError struct {
	Action string
}

func (e *CausesRequiredError) Error() string {
	return "causes required for action: " + e.Action
}
func (e *CausesRequiredError) decisionError() {}
