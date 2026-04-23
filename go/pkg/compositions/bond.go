package compositions

import (
	"context"
	"fmt"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/grammar"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// BondGrammar provides Layer 9 (Relationship) composition operations.
// 10 operations + 3 named functions for deep relational bonds.
type BondGrammar struct {
	g *grammar.Grammar
}

// NewBondGrammar creates a BondGrammar bound to the given base grammar.
func NewBondGrammar(g *grammar.Grammar) *BondGrammar {
	return &BondGrammar{g: g}
}

// --- Operations (10) ---

// Connect initiates a relational bond. (Attachment + Subscribe mutual)
func (b *BondGrammar) Connect(
	ctx context.Context, source types.ActorID, target types.ActorID,
	scope types.Option[types.DomainScope],
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (sub1 event.Event, sub2 event.Event, err error) {
	sub1, err = b.g.Subscribe(ctx, source, target, scope, cause, convID, signer)
	if err != nil {
		return event.Event{}, event.Event{}, fmt.Errorf("connect/subscribe-1: %w", err)
	}
	sub2, err = b.g.Subscribe(ctx, target, source, scope, sub1.ID(), convID, signer)
	if err != nil {
		return event.Event{}, event.Event{}, fmt.Errorf("connect/subscribe-2: %w", err)
	}
	return sub1, sub2, nil
}

// Balance assesses and adjusts reciprocity. (Reciprocity + Annotate)
func (b *BondGrammar) Balance(
	ctx context.Context, source types.ActorID, target types.EventID,
	assessment string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Annotate(ctx, source, target, "reciprocity", assessment, convID, signer)
}

// Deepen extends relational trust beyond transactional. (Trust + Consent)
func (b *BondGrammar) Deepen(
	ctx context.Context, source types.ActorID, other types.ActorID,
	basis string, scope types.DomainScope,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Consent(ctx, source, other, "deepen: "+basis, scope, cause, convID, signer)
}

// Open shares vulnerability with another. (Vulnerability + Channel)
func (b *BondGrammar) Open(
	ctx context.Context, source types.ActorID, target types.ActorID,
	scope types.Option[types.DomainScope],
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Channel(ctx, source, target, scope, cause, convID, signer)
}

// Attune develops accurate understanding of another. (Understanding + Emit)
func (b *BondGrammar) Attune(
	ctx context.Context, source types.ActorID, understanding string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "attune: "+understanding, convID, causes, signer)
}

// FeelWith expresses empathy for another's state. (Empathy + Respond)
func (b *BondGrammar) FeelWith(
	ctx context.Context, source types.ActorID, empathy string,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Respond(ctx, source, "empathy: "+empathy, target, convID, signer)
}

// Break acknowledges a relational rupture. (Rupture + Emit)
func (b *BondGrammar) Break(
	ctx context.Context, source types.ActorID, rupture string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "rupture: "+rupture, convID, causes, signer)
}

// Apologize acknowledges harm and takes responsibility. (Apology + Emit)
func (b *BondGrammar) Apologize(
	ctx context.Context, source types.ActorID, apology string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "apology: "+apology, convID, causes, signer)
}

// Reconcile rebuilds relationship after rupture. (Reconciliation + Consent)
func (b *BondGrammar) Reconcile(
	ctx context.Context, source types.ActorID, other types.ActorID,
	progress string, scope types.DomainScope,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Consent(ctx, source, other, "reconcile: "+progress, scope, cause, convID, signer)
}

// Mourn processes the permanent end of a relationship. (Loss + Emit)
func (b *BondGrammar) Mourn(
	ctx context.Context, source types.ActorID, loss string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "mourn: "+loss, convID, causes, signer)
}

// --- Named Functions (5) ---

// BetrayalRepairResult holds the events produced by BetrayalRepair.
type BetrayalRepairResult struct {
	Rupture       event.Event
	Apology       event.Event
	Reconciliation event.Event
	Deepened      event.Event
}

// BetrayalRepair runs a full rupture-to-growth cycle: Break + Apologize + Reconcile + Deepen.
func (b *BondGrammar) BetrayalRepair(
	ctx context.Context, injured types.ActorID, offender types.ActorID,
	rupture string, apology string, reconciliation string, newBasis string,
	scope types.DomainScope,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (BetrayalRepairResult, error) {
	ruptureEv, err := b.Break(ctx, injured, rupture, causes, convID, signer)
	if err != nil {
		return BetrayalRepairResult{}, fmt.Errorf("betrayal-repair/break: %w", err)
	}

	apologyEv, err := b.Apologize(ctx, offender, apology, []types.EventID{ruptureEv.ID()}, convID, signer)
	if err != nil {
		return BetrayalRepairResult{}, fmt.Errorf("betrayal-repair/apologize: %w", err)
	}

	reconcileEv, err := b.Reconcile(ctx, injured, offender, reconciliation, scope, apologyEv.ID(), convID, signer)
	if err != nil {
		return BetrayalRepairResult{}, fmt.Errorf("betrayal-repair/reconcile: %w", err)
	}

	deepen, err := b.Deepen(ctx, injured, offender, newBasis, scope, reconcileEv.ID(), convID, signer)
	if err != nil {
		return BetrayalRepairResult{}, fmt.Errorf("betrayal-repair/deepen: %w", err)
	}

	return BetrayalRepairResult{
		Rupture:        ruptureEv,
		Apology:        apologyEv,
		Reconciliation: reconcileEv,
		Deepened:       deepen,
	}, nil
}

// CheckInResult holds the events produced by a CheckIn.
type CheckInResult struct {
	Balance event.Event
	Attunement event.Event
	Empathy event.Event
}

// CheckIn is a regular relationship health assessment: Balance + Attune + FeelWith.
func (b *BondGrammar) CheckIn(
	ctx context.Context, source types.ActorID,
	balanceTarget types.EventID, assessment string,
	attunement string, empathy string,
	convID types.ConversationID, signer event.Signer,
) (CheckInResult, error) {
	bal, err := b.Balance(ctx, source, balanceTarget, assessment, convID, signer)
	if err != nil {
		return CheckInResult{}, fmt.Errorf("check-in/balance: %w", err)
	}

	att, err := b.Attune(ctx, source, attunement, []types.EventID{bal.ID()}, convID, signer)
	if err != nil {
		return CheckInResult{}, fmt.Errorf("check-in/attune: %w", err)
	}

	emp, err := b.FeelWith(ctx, source, empathy, att.ID(), convID, signer)
	if err != nil {
		return CheckInResult{}, fmt.Errorf("check-in/feel-with: %w", err)
	}

	return CheckInResult{Balance: bal, Attunement: att, Empathy: emp}, nil
}

// BondMentorshipResult holds the events produced by a Mentorship.
type BondMentorshipResult struct {
	Connection event.Event
	Deepening  event.Event
	Attunement event.Event
	Teaching   event.Event
}

// Mentorship establishes a mentor-mentee bond: Connect (one-way) + Deepen + Attune + Channel (teaching).
func (b *BondGrammar) Mentorship(
	ctx context.Context, mentor types.ActorID, mentee types.ActorID,
	basis string, understanding string,
	scope types.DomainScope, teachingScope types.Option[types.DomainScope],
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (BondMentorshipResult, error) {
	connect, err := b.g.Subscribe(ctx, mentee, mentor, teachingScope, cause, convID, signer)
	if err != nil {
		return BondMentorshipResult{}, fmt.Errorf("mentorship/connect: %w", err)
	}

	deepen, err := b.Deepen(ctx, mentor, mentee, basis, scope, connect.ID(), convID, signer)
	if err != nil {
		return BondMentorshipResult{}, fmt.Errorf("mentorship/deepen: %w", err)
	}

	attune, err := b.Attune(ctx, mentor, understanding, []types.EventID{deepen.ID()}, convID, signer)
	if err != nil {
		return BondMentorshipResult{}, fmt.Errorf("mentorship/attune: %w", err)
	}

	teach, err := b.g.Channel(ctx, mentor, mentee, teachingScope, attune.ID(), convID, signer)
	if err != nil {
		return BondMentorshipResult{}, fmt.Errorf("mentorship/teach: %w", err)
	}

	return BondMentorshipResult{Connection: connect, Deepening: deepen, Attunement: attune, Teaching: teach}, nil
}

// BondFarewellResult holds the events produced by a Farewell.
type BondFarewellResult struct {
	Mourning  event.Event
	Memorial  event.Event
	Gratitude event.Event
}

// Farewell processes the end of a bond: Mourn + Memorialize (via Emit) + Gratitude (via Endorse).
func (b *BondGrammar) Farewell(
	ctx context.Context, source types.ActorID, departing types.ActorID,
	loss string, memorial string, gratitudeWeight types.Weight,
	scope types.Option[types.DomainScope],
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (BondFarewellResult, error) {
	mourn, err := b.Mourn(ctx, source, loss, causes, convID, signer)
	if err != nil {
		return BondFarewellResult{}, fmt.Errorf("farewell/mourn: %w", err)
	}

	mem, err := b.g.Emit(ctx, source, "memorialize: "+memorial, convID, []types.EventID{mourn.ID()}, signer)
	if err != nil {
		return BondFarewellResult{}, fmt.Errorf("farewell/memorialize: %w", err)
	}

	gratitude, err := b.g.Endorse(ctx, source, mem.ID(), departing, gratitudeWeight, scope, convID, signer)
	if err != nil {
		return BondFarewellResult{}, fmt.Errorf("farewell/gratitude: %w", err)
	}

	return BondFarewellResult{Mourning: mourn, Memorial: mem, Gratitude: gratitude}, nil
}

// Forgive re-establishes connection after sever: Subscribe after Sever.
func (b *BondGrammar) Forgive(
	ctx context.Context, source types.ActorID,
	severEvent types.EventID, target types.ActorID,
	scope types.Option[types.DomainScope],
	convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Forgive(ctx, source, severEvent, target, scope, convID, signer)
}
