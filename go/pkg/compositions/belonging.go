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

// FestivalResult holds the events produced by a Festival.
type FestivalResult struct {
	Celebration event.Event
	Practice    event.Event
	Story       event.Event
	Gift        event.Event
}

// Festival is a collective celebration: Celebrate + Practice + Tell + Gift.
func (b *BelongingGrammar) Festival(
	ctx context.Context, source types.ActorID,
	celebration string, tradition string, story string, gift string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (FestivalResult, error) {
	celebrate, err := b.Celebrate(ctx, source, celebration, causes, convID, signer)
	if err != nil {
		return FestivalResult{}, fmt.Errorf("festival/celebrate: %w", err)
	}

	practice, err := b.Practice(ctx, source, tradition, []types.EventID{celebrate.ID()}, convID, signer)
	if err != nil {
		return FestivalResult{}, fmt.Errorf("festival/practice: %w", err)
	}

	tell, err := b.Tell(ctx, source, story, []types.EventID{practice.ID()}, convID, signer)
	if err != nil {
		return FestivalResult{}, fmt.Errorf("festival/tell: %w", err)
	}

	giftEv, err := b.Gift(ctx, source, gift, []types.EventID{tell.ID()}, convID, signer)
	if err != nil {
		return FestivalResult{}, fmt.Errorf("festival/gift: %w", err)
	}

	return FestivalResult{Celebration: celebrate, Practice: practice, Story: tell, Gift: giftEv}, nil
}

// CommonsGovernanceResult holds the events produced by CommonsGovernance.
type CommonsGovernanceResult struct {
	Stewardship event.Event
	Assessment  event.Event
	Legislation event.Event
	Audit       event.Event
}

// CommonsGovernance manages shared resources: Steward + Sustain + Legislate (via Emit) + Audit (via Annotate).
func (b *BelongingGrammar) CommonsGovernance(
	ctx context.Context, source types.ActorID, steward types.ActorID,
	scope types.DomainScope, weight types.Weight,
	assessment string, rule string, findings string,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (CommonsGovernanceResult, error) {
	stewardship, err := b.Steward(ctx, source, steward, scope, weight, cause, convID, signer)
	if err != nil {
		return CommonsGovernanceResult{}, fmt.Errorf("commons/steward: %w", err)
	}

	sustain, err := b.Sustain(ctx, steward, assessment, []types.EventID{stewardship.ID()}, convID, signer)
	if err != nil {
		return CommonsGovernanceResult{}, fmt.Errorf("commons/sustain: %w", err)
	}

	legislate, err := b.g.Emit(ctx, source, "legislate: "+rule, convID, []types.EventID{sustain.ID()}, signer)
	if err != nil {
		return CommonsGovernanceResult{}, fmt.Errorf("commons/legislate: %w", err)
	}

	audit, err := b.g.Annotate(ctx, steward, legislate.ID(), "audit", findings, convID, signer)
	if err != nil {
		return CommonsGovernanceResult{}, fmt.Errorf("commons/audit: %w", err)
	}

	return CommonsGovernanceResult{Stewardship: stewardship, Assessment: sustain, Legislation: legislate, Audit: audit}, nil
}

// RenewalResult holds the events produced by a Renewal.
type RenewalResult struct {
	Assessment event.Event
	Practice   event.Event
	Story      event.Event
}

// Renewal refreshes a community after crisis: Sustain + Practice (evolved) + Tell (new chapter).
func (b *BelongingGrammar) Renewal(
	ctx context.Context, source types.ActorID,
	assessment string, evolvedPractice string, newStory string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (RenewalResult, error) {
	sustain, err := b.Sustain(ctx, source, assessment, causes, convID, signer)
	if err != nil {
		return RenewalResult{}, fmt.Errorf("renewal/sustain: %w", err)
	}

	practice, err := b.Practice(ctx, source, evolvedPractice, []types.EventID{sustain.ID()}, convID, signer)
	if err != nil {
		return RenewalResult{}, fmt.Errorf("renewal/practice: %w", err)
	}

	story, err := b.Tell(ctx, source, newStory, []types.EventID{practice.ID()}, convID, signer)
	if err != nil {
		return RenewalResult{}, fmt.Errorf("renewal/tell: %w", err)
	}

	return RenewalResult{Assessment: sustain, Practice: practice, Story: story}, nil
}

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
