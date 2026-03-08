package compositions

import (
	"context"
	"fmt"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/grammar"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// BelongingGrammar provides Layer 10 (Community) composition operations.
// 10 operations + 2 named functions for communities with shared resources.
type BelongingGrammar struct {
	g *grammar.Grammar
}

// NewBelongingGrammar creates a BelongingGrammar bound to the given base grammar.
func NewBelongingGrammar(g *grammar.Grammar) *BelongingGrammar {
	return &BelongingGrammar{g: g}
}

// --- Operations (10) ---

// Settle develops a sense of home in a community. (Home + Subscribe)
func (b *BelongingGrammar) Settle(
	ctx context.Context, source types.ActorID, community types.ActorID,
	scope types.Option[types.DomainScope],
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Subscribe(ctx, source, community, scope, cause, convID, signer)
}

// Contribute adds value to the community. (Contribution + Emit)
func (b *BelongingGrammar) Contribute(
	ctx context.Context, source types.ActorID, contribution string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "contribute: "+contribution, convID, causes, signer)
}

// Include removes barriers to participation. (Inclusion + Emit)
func (b *BelongingGrammar) Include(
	ctx context.Context, source types.ActorID, action string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "include: "+action, convID, causes, signer)
}

// Practice participates in a community tradition. (Tradition + Emit)
func (b *BelongingGrammar) Practice(
	ctx context.Context, source types.ActorID, tradition string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "practice: "+tradition, convID, causes, signer)
}

// Steward takes responsibility for shared resources. (Commons + Delegate)
func (b *BelongingGrammar) Steward(
	ctx context.Context, source types.ActorID, steward types.ActorID,
	scope types.DomainScope, weight types.Weight,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Delegate(ctx, source, steward, scope, weight, cause, convID, signer)
}

// Sustain evaluates long-term viability. (Sustainability + Emit)
func (b *BelongingGrammar) Sustain(
	ctx context.Context, source types.ActorID, assessment string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "sustain: "+assessment, convID, causes, signer)
}

// PassOn transfers stewardship to the next generation. (Succession + Consent)
func (b *BelongingGrammar) PassOn(
	ctx context.Context, from types.ActorID, to types.ActorID,
	scope types.DomainScope, description string,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Consent(ctx, from, to, "pass-on: "+description, scope, cause, convID, signer)
}

// Celebrate formally recognizes an achievement. (Milestone + Ceremony + Emit)
func (b *BelongingGrammar) Celebrate(
	ctx context.Context, source types.ActorID, celebration string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "celebrate: "+celebration, convID, causes, signer)
}

// Tell adds a chapter to the community's story. (Story + Emit)
func (b *BelongingGrammar) Tell(
	ctx context.Context, source types.ActorID, story string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "tell: "+story, convID, causes, signer)
}

// Gift gives without expectation of return. (Gift + Emit)
func (b *BelongingGrammar) Gift(
	ctx context.Context, source types.ActorID, gift string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "gift: "+gift, convID, causes, signer)
}

// --- Named Functions (5) ---

// OnboardResult holds the events produced by an Onboard.
type OnboardResult struct {
	Inclusion    event.Event
	Settlement   event.Event
	FirstPractice event.Event
	Contribution event.Event
}

// Onboard is a full newcomer welcome: Include + Settle + Practice + Contribute.
func (b *BelongingGrammar) Onboard(
	ctx context.Context, sponsor types.ActorID, newcomer types.ActorID,
	community types.ActorID, scope types.Option[types.DomainScope],
	inclusionAction string, tradition string, firstContribution string,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (OnboardResult, error) {
	inclusion, err := b.Include(ctx, sponsor, inclusionAction, []types.EventID{cause}, convID, signer)
	if err != nil {
		return OnboardResult{}, fmt.Errorf("onboard/include: %w", err)
	}

	settle, err := b.Settle(ctx, newcomer, community, scope, inclusion.ID(), convID, signer)
	if err != nil {
		return OnboardResult{}, fmt.Errorf("onboard/settle: %w", err)
	}

	practice, err := b.Practice(ctx, newcomer, tradition, []types.EventID{settle.ID()}, convID, signer)
	if err != nil {
		return OnboardResult{}, fmt.Errorf("onboard/practice: %w", err)
	}

	contrib, err := b.Contribute(ctx, newcomer, firstContribution, []types.EventID{practice.ID()}, convID, signer)
	if err != nil {
		return OnboardResult{}, fmt.Errorf("onboard/contribute: %w", err)
	}

	return OnboardResult{
		Inclusion:     inclusion,
		Settlement:    settle,
		FirstPractice: practice,
		Contribution:  contrib,
	}, nil
}

// SuccessionResult holds the events produced by a Succession.
type SuccessionResult struct {
	Assessment event.Event
	Transfer   event.Event
	Celebration event.Event
	Story      event.Event
}

// Succession is a full generational transfer: Sustain + PassOn + Celebrate + Tell.
func (b *BelongingGrammar) Succession(
	ctx context.Context, outgoing types.ActorID, incoming types.ActorID,
	assessment string, scope types.DomainScope, celebration string, story string,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (SuccessionResult, error) {
	sustain, err := b.Sustain(ctx, outgoing, assessment, []types.EventID{cause}, convID, signer)
	if err != nil {
		return SuccessionResult{}, fmt.Errorf("succession/sustain: %w", err)
	}

	transfer, err := b.PassOn(ctx, outgoing, incoming, scope, "stewardship transfer", sustain.ID(), convID, signer)
	if err != nil {
		return SuccessionResult{}, fmt.Errorf("succession/pass-on: %w", err)
	}

	celebrate, err := b.Celebrate(ctx, outgoing, celebration, []types.EventID{transfer.ID()}, convID, signer)
	if err != nil {
		return SuccessionResult{}, fmt.Errorf("succession/celebrate: %w", err)
	}

	tell, err := b.Tell(ctx, outgoing, story, []types.EventID{celebrate.ID()}, convID, signer)
	if err != nil {
		return SuccessionResult{}, fmt.Errorf("succession/tell: %w", err)
	}

	return SuccessionResult{
		Assessment:  sustain,
		Transfer:    transfer,
		Celebration: celebrate,
		Story:       tell,
	}, nil
}
