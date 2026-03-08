package compositions

import (
	"context"
	"fmt"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/grammar"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// JusticeGrammar provides Layer 4 (Legal) composition operations.
// 12 operations + 1 named function for transparent dispute resolution.
type JusticeGrammar struct {
	g *grammar.Grammar
}

// NewJusticeGrammar creates a JusticeGrammar bound to the given base grammar.
func NewJusticeGrammar(g *grammar.Grammar) *JusticeGrammar {
	return &JusticeGrammar{g: g}
}

// --- Operations (12) ---

// Legislate enacts a formal rule. (Rule + Emit)
func (j *JusticeGrammar) Legislate(
	ctx context.Context, source types.ActorID, rule string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return j.g.Emit(ctx, source, "legislate: "+rule, convID, causes, signer)
}

// Amend changes an existing rule. (Rule amended + Derive)
func (j *JusticeGrammar) Amend(
	ctx context.Context, source types.ActorID, amendment string,
	rule types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return j.g.Derive(ctx, source, "amend: "+amendment, rule, convID, signer)
}

// Repeal revokes an existing rule. (Rule repealed + Retract)
func (j *JusticeGrammar) Repeal(
	ctx context.Context, source types.ActorID, rule types.EventID,
	reason string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return j.g.Retract(ctx, source, rule, reason, convID, signer)
}

// File brings a formal complaint. (DueProcess + Challenge)
func (j *JusticeGrammar) File(
	ctx context.Context, source types.ActorID, complaint string,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	_, flag, err := j.g.Challenge(ctx, source, "file: "+complaint, target, convID, signer)
	if err != nil {
		return event.Event{}, err
	}
	return flag, nil
}

// Submit presents evidence for a case. (Precedent + Annotate)
func (j *JusticeGrammar) Submit(
	ctx context.Context, source types.ActorID, target types.EventID,
	evidence string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return j.g.Annotate(ctx, source, target, "evidence", evidence, convID, signer)
}

// Argue makes a legal argument. (Interpretation + Respond)
func (j *JusticeGrammar) Argue(
	ctx context.Context, source types.ActorID, argument string,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return j.g.Respond(ctx, source, "argue: "+argument, target, convID, signer)
}

// Judge renders a formal ruling. (Adjudication + Emit)
func (j *JusticeGrammar) Judge(
	ctx context.Context, source types.ActorID, ruling string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return j.g.Emit(ctx, source, "judge: "+ruling, convID, causes, signer)
}

// Appeal challenges a ruling to higher authority. (Appeal + Challenge)
func (j *JusticeGrammar) Appeal(
	ctx context.Context, source types.ActorID, grounds string,
	ruling types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	_, flag, err := j.g.Challenge(ctx, source, "appeal: "+grounds, ruling, convID, signer)
	if err != nil {
		return event.Event{}, err
	}
	return flag, nil
}

// Enforce executes consequences of a ruling. (Enforcement + Delegate)
func (j *JusticeGrammar) Enforce(
	ctx context.Context, source types.ActorID, executor types.ActorID,
	scope types.DomainScope, weight types.Weight,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return j.g.Delegate(ctx, source, executor, scope, weight, cause, convID, signer)
}

// Audit performs systematic review against rules. (Audit + Annotate)
func (j *JusticeGrammar) Audit(
	ctx context.Context, source types.ActorID, target types.EventID,
	findings string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return j.g.Annotate(ctx, source, target, "audit", findings, convID, signer)
}

// Pardon formally forgives a violation. (Amnesty + Consent)
func (j *JusticeGrammar) Pardon(
	ctx context.Context, authority types.ActorID, pardoned types.ActorID,
	terms string, scope types.DomainScope,
	cause types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return j.g.Consent(ctx, authority, pardoned, "pardon: "+terms, scope, cause, convID, signer)
}

// Reform proposes systemic rule change based on experience. (Reform + Derive)
func (j *JusticeGrammar) Reform(
	ctx context.Context, source types.ActorID, proposal string,
	precedent types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return j.g.Derive(ctx, source, "reform: "+proposal, precedent, convID, signer)
}

// --- Named Functions (1) ---

// TrialResult holds the events produced by a Trial.
type TrialResult struct {
	Filing     event.Event
	Submissions []event.Event
	Arguments  []event.Event
	Ruling     event.Event
}

// Trial runs a full adjudication: File + Submit + Argue + Judge.
func (j *JusticeGrammar) Trial(
	ctx context.Context, plaintiff types.ActorID, defendant types.ActorID,
	judge types.ActorID, complaint string,
	plaintiffEvidence string, defendantEvidence string,
	plaintiffArgument string, defendantArgument string,
	ruling string,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (TrialResult, error) {
	filing, err := j.File(ctx, plaintiff, complaint, target, convID, signer)
	if err != nil {
		return TrialResult{}, fmt.Errorf("trial/file: %w", err)
	}

	sub1, err := j.Submit(ctx, plaintiff, filing.ID(), plaintiffEvidence, convID, signer)
	if err != nil {
		return TrialResult{}, fmt.Errorf("trial/submit-plaintiff: %w", err)
	}

	sub2, err := j.Submit(ctx, defendant, filing.ID(), defendantEvidence, convID, signer)
	if err != nil {
		return TrialResult{}, fmt.Errorf("trial/submit-defendant: %w", err)
	}

	arg1, err := j.Argue(ctx, plaintiff, plaintiffArgument, sub1.ID(), convID, signer)
	if err != nil {
		return TrialResult{}, fmt.Errorf("trial/argue-plaintiff: %w", err)
	}

	arg2, err := j.Argue(ctx, defendant, defendantArgument, sub2.ID(), convID, signer)
	if err != nil {
		return TrialResult{}, fmt.Errorf("trial/argue-defendant: %w", err)
	}

	verdict, err := j.Judge(ctx, judge, ruling,
		[]types.EventID{arg1.ID(), arg2.ID()}, convID, signer)
	if err != nil {
		return TrialResult{}, fmt.Errorf("trial/judge: %w", err)
	}

	return TrialResult{
		Filing:      filing,
		Submissions: []event.Event{sub1, sub2},
		Arguments:   []event.Event{arg1, arg2},
		Ruling:      verdict,
	}, nil
}
