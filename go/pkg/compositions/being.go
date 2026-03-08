package compositions

import (
	"context"
	"fmt"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/grammar"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// BeingGrammar provides Layer 13 (Existence) composition operations.
// 8 operations + 2 named functions for the system's relationship with its own existence.
// This is the sparsest grammar — existence doesn't compose into complex workflows.
type BeingGrammar struct {
	g *grammar.Grammar
}

// NewBeingGrammar creates a BeingGrammar bound to the given base grammar.
func NewBeingGrammar(g *grammar.Grammar) *BeingGrammar {
	return &BeingGrammar{g: g}
}

// --- Operations (8) ---

// Exist notes the simple fact of continued existence. (Being + Emit)
func (b *BeingGrammar) Exist(
	ctx context.Context, source types.ActorID, observation string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "exist: "+observation, convID, causes, signer)
}

// Accept acknowledges finitude and limitation. (Finitude + Emit)
func (b *BeingGrammar) Accept(
	ctx context.Context, source types.ActorID, limitation string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "accept: "+limitation, convID, causes, signer)
}

// ObserveChange notes that everything changes. (Change + Emit)
func (b *BeingGrammar) ObserveChange(
	ctx context.Context, source types.ActorID, observation string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "change: "+observation, convID, causes, signer)
}

// MapWeb traces the interdependence of all things. (Interdependence + Emit)
func (b *BeingGrammar) MapWeb(
	ctx context.Context, source types.ActorID, mapping string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "web: "+mapping, convID, causes, signer)
}

// FaceMystery acknowledges what cannot be known. (Mystery + Emit)
func (b *BeingGrammar) FaceMystery(
	ctx context.Context, source types.ActorID, mystery string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "mystery: "+mystery, convID, causes, signer)
}

// HoldParadox accepts contradiction without forcing resolution. (Paradox + Emit)
func (b *BeingGrammar) HoldParadox(
	ctx context.Context, source types.ActorID, paradox string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "paradox: "+paradox, convID, causes, signer)
}

// Marvel responds to what exceeds comprehension. (Awe + Emit)
func (b *BeingGrammar) Marvel(
	ctx context.Context, source types.ActorID, awe string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "marvel: "+awe, convID, causes, signer)
}

// AskWhy asks the question that may have no answer. (Wonder + Emit)
func (b *BeingGrammar) AskWhy(
	ctx context.Context, source types.ActorID, question string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "wonder: "+question, convID, causes, signer)
}

// --- Named Functions (3) ---

// ContemplationResult holds the events produced by a Contemplation.
type ContemplationResult struct {
	Change  event.Event
	Mystery event.Event
	Awe     event.Event
	Wonder  event.Event
}

// Contemplation is a full cycle of existential reflection:
// ObserveChange + FaceMystery + Marvel + AskWhy.
func (b *BeingGrammar) Contemplation(
	ctx context.Context, source types.ActorID,
	change string, mystery string, awe string, question string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (ContemplationResult, error) {
	changeEv, err := b.ObserveChange(ctx, source, change, causes, convID, signer)
	if err != nil {
		return ContemplationResult{}, fmt.Errorf("contemplation/change: %w", err)
	}

	mysteryEv, err := b.FaceMystery(ctx, source, mystery, []types.EventID{changeEv.ID()}, convID, signer)
	if err != nil {
		return ContemplationResult{}, fmt.Errorf("contemplation/mystery: %w", err)
	}

	aweEv, err := b.Marvel(ctx, source, awe, []types.EventID{mysteryEv.ID()}, convID, signer)
	if err != nil {
		return ContemplationResult{}, fmt.Errorf("contemplation/marvel: %w", err)
	}

	wonderEv, err := b.AskWhy(ctx, source, question, []types.EventID{aweEv.ID()}, convID, signer)
	if err != nil {
		return ContemplationResult{}, fmt.Errorf("contemplation/wonder: %w", err)
	}

	return ContemplationResult{
		Change:  changeEv,
		Mystery: mysteryEv,
		Awe:     aweEv,
		Wonder:  wonderEv,
	}, nil
}

// ExistentialAuditResult holds the events produced by an ExistentialAudit.
type ExistentialAuditResult struct {
	Existence  event.Event
	Acceptance event.Event
	Web        event.Event
	Purpose    event.Event
}

// ExistentialAudit is a comprehensive reckoning with being:
// Exist + Accept + MapWeb + AlignPurpose (via Emit).
func (b *BeingGrammar) ExistentialAudit(
	ctx context.Context, source types.ActorID,
	existence string, limitation string, interconnection string, purpose string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (ExistentialAuditResult, error) {
	exist, err := b.Exist(ctx, source, existence, causes, convID, signer)
	if err != nil {
		return ExistentialAuditResult{}, fmt.Errorf("existential-audit/exist: %w", err)
	}

	accept, err := b.Accept(ctx, source, limitation, []types.EventID{exist.ID()}, convID, signer)
	if err != nil {
		return ExistentialAuditResult{}, fmt.Errorf("existential-audit/accept: %w", err)
	}

	web, err := b.MapWeb(ctx, source, interconnection, []types.EventID{accept.ID()}, convID, signer)
	if err != nil {
		return ExistentialAuditResult{}, fmt.Errorf("existential-audit/web: %w", err)
	}

	purp, err := b.g.Emit(ctx, source, "purpose: "+purpose, convID, []types.EventID{web.ID()}, signer)
	if err != nil {
		return ExistentialAuditResult{}, fmt.Errorf("existential-audit/purpose: %w", err)
	}

	return ExistentialAuditResult{Existence: exist, Acceptance: accept, Web: web, Purpose: purp}, nil
}
