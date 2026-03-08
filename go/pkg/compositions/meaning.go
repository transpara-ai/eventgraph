package compositions

import (
	"context"
	"fmt"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/grammar"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// MeaningGrammar provides Layer 11 (Culture) composition operations.
// 10 operations + 2 named functions for cross-cultural communication.
type MeaningGrammar struct {
	g *grammar.Grammar
}

// NewMeaningGrammar creates a MeaningGrammar bound to the given base grammar.
func NewMeaningGrammar(g *grammar.Grammar) *MeaningGrammar {
	return &MeaningGrammar{g: g}
}

// --- Operations (10) ---

// Examine identifies blind spots and assumptions. (SelfAwareness + Emit)
func (m *MeaningGrammar) Examine(
	ctx context.Context, source types.ActorID, examination string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Emit(ctx, source, "examine: "+examination, convID, causes, signer)
}

// Reframe sees a situation from a different viewpoint. (Perspective + Derive)
func (m *MeaningGrammar) Reframe(
	ctx context.Context, source types.ActorID, reframing string,
	original types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Derive(ctx, source, "reframe: "+reframing, original, convID, signer)
}

// Question challenges what's taken for granted. (Critique + Respond)
func (m *MeaningGrammar) Question(
	ctx context.Context, source types.ActorID, question string,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Respond(ctx, source, "question: "+question, target, convID, signer)
}

// Distill extracts what truly matters from experience. (Wisdom + Derive)
func (m *MeaningGrammar) Distill(
	ctx context.Context, source types.ActorID, wisdom string,
	experience types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Derive(ctx, source, "distill: "+wisdom, experience, convID, signer)
}

// Beautify recognizes or creates beauty and elegance. (Aesthetic + Annotate or Emit)
func (m *MeaningGrammar) Beautify(
	ctx context.Context, source types.ActorID, beauty string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Emit(ctx, source, "beautify: "+beauty, convID, causes, signer)
}

// Liken explains one thing in terms of another. (Metaphor + Derive)
func (m *MeaningGrammar) Liken(
	ctx context.Context, source types.ActorID, metaphor string,
	subject types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Derive(ctx, source, "liken: "+metaphor, subject, convID, signer)
}

// Lighten finds incongruity and playfulness. (Humour + Emit)
func (m *MeaningGrammar) Lighten(
	ctx context.Context, source types.ActorID, humour string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Emit(ctx, source, "lighten: "+humour, convID, causes, signer)
}

// Teach deliberately transfers knowledge. (Teaching + Channel)
func (m *MeaningGrammar) Teach(
	ctx context.Context, source types.ActorID, student types.ActorID,
	scope types.Option[types.DomainScope],
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Channel(ctx, source, student, scope, cause, convID, signer)
}

// Translate makes meaning accessible across boundaries. (Translation + Derive)
func (m *MeaningGrammar) Translate(
	ctx context.Context, source types.ActorID, translation string,
	original types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Derive(ctx, source, "translate: "+translation, original, convID, signer)
}

// Prophesy extrapolates where current patterns lead. (Prophecy + Emit)
func (m *MeaningGrammar) Prophesy(
	ctx context.Context, source types.ActorID, prediction string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return m.g.Emit(ctx, source, "prophesy: "+prediction, convID, causes, signer)
}

// --- Named Functions (5) ---

// MeaningPostMortemResult holds the events produced by a PostMortem.
type MeaningPostMortemResult struct {
	Examination event.Event
	Questions   event.Event
	Wisdom      event.Event
}

// PostMortem learns from failure through reflection: Examine + Question + Distill.
func (m *MeaningGrammar) PostMortem(
	ctx context.Context, source types.ActorID,
	examination string, question string, wisdom string,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (MeaningPostMortemResult, error) {
	exam, err := m.Examine(ctx, source, examination, []types.EventID{cause}, convID, signer)
	if err != nil {
		return MeaningPostMortemResult{}, fmt.Errorf("post-mortem/examine: %w", err)
	}

	q, err := m.Question(ctx, source, question, exam.ID(), convID, signer)
	if err != nil {
		return MeaningPostMortemResult{}, fmt.Errorf("post-mortem/question: %w", err)
	}

	w, err := m.Distill(ctx, source, wisdom, q.ID(), convID, signer)
	if err != nil {
		return MeaningPostMortemResult{}, fmt.Errorf("post-mortem/distill: %w", err)
	}

	return MeaningPostMortemResult{Examination: exam, Questions: q, Wisdom: w}, nil
}

// MentorshipResult holds the events produced by a Mentorship.
type MentorshipResult struct {
	Channel     event.Event
	Reframing   event.Event
	Wisdom      event.Event
	Translation event.Event
}

// Mentorship is deep knowledge transfer: Teach + Reframe + Distill + Translate.
func (m *MeaningGrammar) Mentorship(
	ctx context.Context, mentor types.ActorID, student types.ActorID,
	reframing string, wisdom string, translation string,
	scope types.Option[types.DomainScope],
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (MentorshipResult, error) {
	channel, err := m.Teach(ctx, mentor, student, scope, cause, convID, signer)
	if err != nil {
		return MentorshipResult{}, fmt.Errorf("mentorship/teach: %w", err)
	}

	reframe, err := m.Reframe(ctx, mentor, reframing, channel.ID(), convID, signer)
	if err != nil {
		return MentorshipResult{}, fmt.Errorf("mentorship/reframe: %w", err)
	}

	distill, err := m.Distill(ctx, mentor, wisdom, reframe.ID(), convID, signer)
	if err != nil {
		return MentorshipResult{}, fmt.Errorf("mentorship/distill: %w", err)
	}

	translate, err := m.Translate(ctx, student, translation, distill.ID(), convID, signer)
	if err != nil {
		return MentorshipResult{}, fmt.Errorf("mentorship/translate: %w", err)
	}

	return MentorshipResult{
		Channel:     channel,
		Reframing:   reframe,
		Wisdom:      distill,
		Translation: translate,
	}, nil
}
