package compositions

import (
	"context"
	"fmt"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/grammar"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// SocialGrammar provides Layer 3 (Society) composition operations.
// 5 society-specific extensions + 2 named functions for user-owned social platforms.
type SocialGrammar struct {
	g *grammar.Grammar
}

// NewSocialGrammar creates a SocialGrammar bound to the given base grammar.
func NewSocialGrammar(g *grammar.Grammar) *SocialGrammar {
	return &SocialGrammar{g: g}
}

// --- Society-Specific Extensions (5) ---

// Norm establishes a shared behavioural expectation. (Norm + Consent)
func (s *SocialGrammar) Norm(
	ctx context.Context, proposer types.ActorID, supporter types.ActorID,
	norm string, scope types.DomainScope,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return s.g.Consent(ctx, proposer, supporter, "norm: "+norm, scope, cause, convID, signer)
}

// Moderate enforces community norms on content. (Sanction + Retract or Annotate)
// Uses Annotate to flag content without removing it.
func (s *SocialGrammar) Moderate(
	ctx context.Context, moderator types.ActorID, target types.EventID,
	action string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return s.g.Annotate(ctx, moderator, target, "moderation", action, convID, signer)
}

// Elect assigns a community role through collective decision. (Role + Consent)
func (s *SocialGrammar) Elect(
	ctx context.Context, community types.ActorID, elected types.ActorID,
	role string, scope types.DomainScope,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return s.g.Consent(ctx, community, elected, "elect: "+role, scope, cause, convID, signer)
}

// Welcome is structured onboarding of a new member. (Inclusion + Invite)
func (s *SocialGrammar) Welcome(
	ctx context.Context, sponsor types.ActorID, newcomer types.ActorID,
	weight types.Weight, scope types.Option[types.DomainScope],
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (endorseEv event.Event, subscribeEv event.Event, err error) {
	return s.g.Invite(ctx, sponsor, newcomer, weight, scope, cause, convID, signer)
}

// ExileResult holds the events produced by an Exile.
type ExileResult struct {
	Exclusion event.Event
	Sever     event.Event
	Sanction  event.Event
}

// Exile is structured removal of a member. (Exclusion + Sever + Sanction)
func (s *SocialGrammar) Exile(
	ctx context.Context, moderator types.ActorID,
	edge types.EdgeID, reason string,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (ExileResult, error) {
	exclusion, err := s.g.Emit(ctx, moderator, "exile: "+reason, convID, []types.EventID{cause}, signer)
	if err != nil {
		return ExileResult{}, fmt.Errorf("exile/exclusion: %w", err)
	}

	sever, err := s.g.Sever(ctx, moderator, edge, exclusion.ID(), convID, signer)
	if err != nil {
		return ExileResult{}, fmt.Errorf("exile/sever: %w", err)
	}

	sanction, err := s.g.Annotate(ctx, moderator, sever.ID(), "sanction", reason, convID, signer)
	if err != nil {
		return ExileResult{}, fmt.Errorf("exile/sanction: %w", err)
	}

	return ExileResult{Exclusion: exclusion, Sever: sever, Sanction: sanction}, nil
}

// --- Named Functions (3) ---

// PollResult holds the events produced by a Poll.
type PollResult struct {
	Proposal event.Event
	Votes    []event.Event
}

// Poll runs a quick community sentiment check: Norm (proposed) + Consent (batch).
func (s *SocialGrammar) Poll(
	ctx context.Context, proposer types.ActorID, question string,
	voters []types.ActorID, scope types.DomainScope,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (PollResult, error) {
	proposal, err := s.g.Emit(ctx, proposer, "poll: "+question, convID, []types.EventID{cause}, signer)
	if err != nil {
		return PollResult{}, fmt.Errorf("poll/propose: %w", err)
	}

	result := PollResult{Proposal: proposal}
	for i, voter := range voters {
		vote, err := s.g.Consent(ctx, voter, proposer, "vote: "+question, scope, proposal.ID(), convID, signer)
		if err != nil {
			return PollResult{}, fmt.Errorf("poll/vote[%d]: %w", i, err)
		}
		result.Votes = append(result.Votes, vote)
	}
	return result, nil
}

// FederationResult holds the events produced by a Federation.
type FederationResult struct {
	Agreement  event.Event
	Delegation event.Event
}

// Federation creates cooperation between communities while maintaining autonomy.
func (s *SocialGrammar) Federation(
	ctx context.Context, communityA types.ActorID, communityB types.ActorID,
	terms string, scope types.DomainScope, weight types.Weight,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (FederationResult, error) {
	agreement, err := s.g.Consent(ctx, communityA, communityB, "federation: "+terms, scope, cause, convID, signer)
	if err != nil {
		return FederationResult{}, fmt.Errorf("federation/consent: %w", err)
	}

	delegation, err := s.g.Delegate(ctx, communityA, communityB, scope, weight, agreement.ID(), convID, signer)
	if err != nil {
		return FederationResult{}, fmt.Errorf("federation/delegate: %w", err)
	}

	return FederationResult{Agreement: agreement, Delegation: delegation}, nil
}
