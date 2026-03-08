package compositions

import (
	"context"
	"fmt"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/grammar"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// IdentityGrammar provides Layer 8 (Identity) composition operations.
// 10 operations + 2 named functions for self-sovereign identity.
type IdentityGrammar struct {
	g *grammar.Grammar
}

// NewIdentityGrammar creates an IdentityGrammar bound to the given base grammar.
func NewIdentityGrammar(g *grammar.Grammar) *IdentityGrammar {
	return &IdentityGrammar{g: g}
}

// --- Operations (10) ---

// Introspect forms or updates a self-model. (SelfModel + Emit)
func (i *IdentityGrammar) Introspect(
	ctx context.Context, source types.ActorID, selfModel string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return i.g.Emit(ctx, source, "introspect: "+selfModel, convID, causes, signer)
}

// Narrate constructs an identity narrative from history. (NarrativeIdentity + Derive)
func (i *IdentityGrammar) Narrate(
	ctx context.Context, source types.ActorID, narrative string,
	basis types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return i.g.Derive(ctx, source, "narrate: "+narrative, basis, convID, signer)
}

// Align checks behaviour against self-model. (Authenticity + Annotate)
func (i *IdentityGrammar) Align(
	ctx context.Context, source types.ActorID, target types.EventID,
	alignment string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return i.g.Annotate(ctx, source, target, "alignment", alignment, convID, signer)
}

// Bound defines or enforces personal boundaries. (Boundary + Emit)
func (i *IdentityGrammar) Bound(
	ctx context.Context, source types.ActorID, boundary string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return i.g.Emit(ctx, source, "bound: "+boundary, convID, causes, signer)
}

// Aspire declares who you want to become. (Aspiration + Emit)
func (i *IdentityGrammar) Aspire(
	ctx context.Context, source types.ActorID, aspiration string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return i.g.Emit(ctx, source, "aspire: "+aspiration, convID, causes, signer)
}

// Transform acknowledges fundamental identity change. (Transformation + Derive)
func (i *IdentityGrammar) Transform(
	ctx context.Context, source types.ActorID, transformation string,
	basis types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return i.g.Derive(ctx, source, "transform: "+transformation, basis, convID, signer)
}

// Disclose selectively reveals aspects of identity. (SelfModel + Derive)
func (i *IdentityGrammar) Disclose(
	ctx context.Context, source types.ActorID, disclosure string,
	selfModel types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return i.g.Derive(ctx, source, "disclose: "+disclosure, selfModel, convID, signer)
}

// Recognize acknowledges another's unique identity. (Dignity + Emit)
func (i *IdentityGrammar) Recognize(
	ctx context.Context, source types.ActorID, recognition string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return i.g.Emit(ctx, source, "recognize: "+recognition, convID, causes, signer)
}

// Distinguish identifies what makes an actor unique. (Uniqueness + Annotate)
func (i *IdentityGrammar) Distinguish(
	ctx context.Context, source types.ActorID, target types.EventID,
	uniqueness string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return i.g.Annotate(ctx, source, target, "uniqueness", uniqueness, convID, signer)
}

// Memorialize preserves identity of a departed actor. (Memorial + Emit)
func (i *IdentityGrammar) Memorialize(
	ctx context.Context, source types.ActorID, memorial string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return i.g.Emit(ctx, source, "memorialize: "+memorial, convID, causes, signer)
}

// --- Named Functions (5) ---

// IdentityAuditResult holds the events produced by an IdentityAudit.
type IdentityAuditResult struct {
	SelfModel event.Event
	Alignment event.Event
	Narrative event.Event
}

// IdentityAudit performs a comprehensive self-assessment: Introspect + Align + Narrate.
func (i *IdentityGrammar) IdentityAudit(
	ctx context.Context, source types.ActorID,
	selfModel string, alignment string, narrative string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (IdentityAuditResult, error) {
	intro, err := i.Introspect(ctx, source, selfModel, causes, convID, signer)
	if err != nil {
		return IdentityAuditResult{}, fmt.Errorf("identity-audit/introspect: %w", err)
	}

	align, err := i.Align(ctx, source, intro.ID(), alignment, convID, signer)
	if err != nil {
		return IdentityAuditResult{}, fmt.Errorf("identity-audit/align: %w", err)
	}

	narr, err := i.Narrate(ctx, source, narrative, align.ID(), convID, signer)
	if err != nil {
		return IdentityAuditResult{}, fmt.Errorf("identity-audit/narrate: %w", err)
	}

	return IdentityAuditResult{SelfModel: intro, Alignment: align, Narrative: narr}, nil
}

// RetirementResult holds the events produced by a Retirement.
type RetirementResult struct {
	Memorial   event.Event
	Transfer   event.Event
}

// Retirement gracefully departs: Memorialize + Transfer (authority).
func (i *IdentityGrammar) Retirement(
	ctx context.Context, system types.ActorID, departing types.ActorID,
	successor types.ActorID, memorial string,
	scope types.DomainScope, weight types.Weight,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (RetirementResult, error) {
	mem, err := i.Memorialize(ctx, system,
		fmt.Sprintf("retirement of %s: %s", departing.Value(), memorial),
		causes, convID, signer)
	if err != nil {
		return RetirementResult{}, fmt.Errorf("retirement/memorialize: %w", err)
	}

	transfer, err := i.g.Delegate(ctx, system, successor, scope, weight, mem.ID(), convID, signer)
	if err != nil {
		return RetirementResult{}, fmt.Errorf("retirement/transfer: %w", err)
	}

	return RetirementResult{Memorial: mem, Transfer: transfer}, nil
}
