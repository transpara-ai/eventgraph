package compositions

import (
	"context"
	"fmt"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/grammar"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
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

// --- Named Functions (6) ---

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

// ConstitutionalAmendmentResult holds the events produced by a ConstitutionalAmendment.
type ConstitutionalAmendmentResult struct {
	Reform      event.Event
	Legislation event.Event
	RightsCheck event.Event
}

// ConstitutionalAmendment proposes systemic change with supermajority:
// Reform + Legislate + Audit (rights check).
func (j *JusticeGrammar) ConstitutionalAmendment(
	ctx context.Context, proposer types.ActorID,
	proposal string, legislation string, rightsAssessment string,
	precedent types.EventID, convID types.ConversationID, signer event.Signer,
) (ConstitutionalAmendmentResult, error) {
	reform, err := j.Reform(ctx, proposer, proposal, precedent, convID, signer)
	if err != nil {
		return ConstitutionalAmendmentResult{}, fmt.Errorf("amendment/reform: %w", err)
	}

	legislate, err := j.Legislate(ctx, proposer, legislation, []types.EventID{reform.ID()}, convID, signer)
	if err != nil {
		return ConstitutionalAmendmentResult{}, fmt.Errorf("amendment/legislate: %w", err)
	}

	rights, err := j.Audit(ctx, proposer, legislate.ID(), rightsAssessment, convID, signer)
	if err != nil {
		return ConstitutionalAmendmentResult{}, fmt.Errorf("amendment/rights-check: %w", err)
	}

	return ConstitutionalAmendmentResult{Reform: reform, Legislation: legislate, RightsCheck: rights}, nil
}

// InjunctionResult holds the events produced by an Injunction.
type InjunctionResult struct {
	Filing      event.Event
	Ruling      event.Event
	Enforcement event.Event
}

// Injunction issues an emergency order: File + Judge (emergency) + Enforce (temporary).
func (j *JusticeGrammar) Injunction(
	ctx context.Context, petitioner types.ActorID, judge types.ActorID,
	executor types.ActorID, complaint string, ruling string,
	scope types.DomainScope, weight types.Weight,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (InjunctionResult, error) {
	filing, err := j.File(ctx, petitioner, complaint, target, convID, signer)
	if err != nil {
		return InjunctionResult{}, fmt.Errorf("injunction/file: %w", err)
	}

	verdict, err := j.Judge(ctx, judge, "emergency: "+ruling, []types.EventID{filing.ID()}, convID, signer)
	if err != nil {
		return InjunctionResult{}, fmt.Errorf("injunction/judge: %w", err)
	}

	enforce, err := j.Enforce(ctx, judge, executor, scope, weight, verdict.ID(), convID, signer)
	if err != nil {
		return InjunctionResult{}, fmt.Errorf("injunction/enforce: %w", err)
	}

	return InjunctionResult{Filing: filing, Ruling: verdict, Enforcement: enforce}, nil
}

// PleaResult holds the events produced by a Plea.
type PleaResult struct {
	Filing      event.Event
	Acceptance  event.Event
	Enforcement event.Event
}

// Plea accepts a reduced penalty: File + Consent (plea deal) + Enforce.
func (j *JusticeGrammar) Plea(
	ctx context.Context, defendant types.ActorID, prosecutor types.ActorID,
	executor types.ActorID, complaint string, deal string,
	scope types.DomainScope, weight types.Weight,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (PleaResult, error) {
	filing, err := j.File(ctx, prosecutor, complaint, target, convID, signer)
	if err != nil {
		return PleaResult{}, fmt.Errorf("plea/file: %w", err)
	}

	acceptance, err := j.Pardon(ctx, prosecutor, defendant, deal, scope, filing.ID(), convID, signer)
	if err != nil {
		return PleaResult{}, fmt.Errorf("plea/accept: %w", err)
	}

	enforce, err := j.Enforce(ctx, prosecutor, executor, scope, weight, acceptance.ID(), convID, signer)
	if err != nil {
		return PleaResult{}, fmt.Errorf("plea/enforce: %w", err)
	}

	return PleaResult{Filing: filing, Acceptance: acceptance, Enforcement: enforce}, nil
}

// ClassActionResult holds the events produced by a ClassAction.
type ClassActionResult struct {
	Filings []event.Event
	Merged  event.Event
	Trial   TrialResult
}

// ClassAction files from multiple parties then runs trial:
// File (multiple, Merge) + Trial.
func (j *JusticeGrammar) ClassAction(
	ctx context.Context, plaintiffs []types.ActorID,
	defendant types.ActorID, judge types.ActorID,
	complaints []string, evidence string, argument string,
	defenseEvidence string, defenseArgument string, ruling string,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (ClassActionResult, error) {
	if len(plaintiffs) != len(complaints) {
		return ClassActionResult{}, fmt.Errorf("class-action: plaintiffs and complaints must have equal length")
	}

	result := ClassActionResult{}
	filingIDs := make([]types.EventID, 0, len(plaintiffs))
	for i, plaintiff := range plaintiffs {
		filing, err := j.File(ctx, plaintiff, complaints[i], target, convID, signer)
		if err != nil {
			return ClassActionResult{}, fmt.Errorf("class-action/file[%d]: %w", i, err)
		}
		result.Filings = append(result.Filings, filing)
		filingIDs = append(filingIDs, filing.ID())
	}

	merged, err := j.g.Merge(ctx, plaintiffs[0], "class-action: merged complaints", filingIDs, convID, signer)
	if err != nil {
		return ClassActionResult{}, fmt.Errorf("class-action/merge: %w", err)
	}
	result.Merged = merged

	trial, err := j.Trial(ctx, plaintiffs[0], defendant, judge,
		"class-action", evidence, defenseEvidence, argument, defenseArgument, ruling,
		merged.ID(), convID, signer)
	if err != nil {
		return ClassActionResult{}, fmt.Errorf("class-action/trial: %w", err)
	}
	result.Trial = trial

	return result, nil
}

// RecallResult holds the events produced by a Recall.
type RecallResult struct {
	Audit      event.Event
	Filing     event.Event
	Consent    event.Event
	Revocation event.Event
}

// Recall removes an authority figure: Audit + File + Consent (community) + revocation.
func (j *JusticeGrammar) Recall(
	ctx context.Context, auditor types.ActorID, community types.ActorID,
	official types.ActorID, findings string, complaint string,
	scope types.DomainScope,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (RecallResult, error) {
	audit, err := j.Audit(ctx, auditor, target, findings, convID, signer)
	if err != nil {
		return RecallResult{}, fmt.Errorf("recall/audit: %w", err)
	}

	filing, err := j.File(ctx, auditor, complaint, audit.ID(), convID, signer)
	if err != nil {
		return RecallResult{}, fmt.Errorf("recall/file: %w", err)
	}

	consent, err := j.g.Consent(ctx, community, official, "recall: "+complaint, scope, filing.ID(), convID, signer)
	if err != nil {
		return RecallResult{}, fmt.Errorf("recall/consent: %w", err)
	}

	revocation, err := j.g.Emit(ctx, community, "role-revoked: "+complaint, convID, []types.EventID{consent.ID()}, signer)
	if err != nil {
		return RecallResult{}, fmt.Errorf("recall/revocation: %w", err)
	}

	return RecallResult{Audit: audit, Filing: filing, Consent: consent, Revocation: revocation}, nil
}
