package compositions

import (
	"context"
	"fmt"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/grammar"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// KnowledgeGrammar provides Layer 6 (Information) composition operations.
// 12 operations + 2 named functions for verified, provenanced knowledge.
type KnowledgeGrammar struct {
	g *grammar.Grammar
}

// NewKnowledgeGrammar creates a KnowledgeGrammar bound to the given base grammar.
func NewKnowledgeGrammar(g *grammar.Grammar) *KnowledgeGrammar {
	return &KnowledgeGrammar{g: g}
}

// --- Operations (12) ---

// Claim makes a knowledge claim with evidence. (Fact + Emit)
func (k *KnowledgeGrammar) Claim(
	ctx context.Context, source types.ActorID, claim string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return k.g.Emit(ctx, source, "claim: "+claim, convID, causes, signer)
}

// Categorize assigns a claim to a taxonomy. (Classification + Annotate)
func (k *KnowledgeGrammar) Categorize(
	ctx context.Context, source types.ActorID, target types.EventID,
	taxonomy string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return k.g.Annotate(ctx, source, target, "classification", taxonomy, convID, signer)
}

// Abstract generalizes from specific instances. (Abstraction + Merge + Derive)
func (k *KnowledgeGrammar) Abstract(
	ctx context.Context, source types.ActorID, generalization string,
	instances []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	if len(instances) < 2 {
		return event.Event{}, fmt.Errorf("abstract: requires at least two instances")
	}
	return k.g.Merge(ctx, source, "abstract: "+generalization, instances, convID, signer)
}

// Encode transforms between representations. (Encoding + Derive)
func (k *KnowledgeGrammar) Encode(
	ctx context.Context, source types.ActorID, encoding string,
	original types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return k.g.Derive(ctx, source, "encode: "+encoding, original, convID, signer)
}

// Infer draws a new conclusion from premises. (Inference + Derive)
func (k *KnowledgeGrammar) Infer(
	ctx context.Context, source types.ActorID, conclusion string,
	premise types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return k.g.Derive(ctx, source, "infer: "+conclusion, premise, convID, signer)
}

// Remember stores knowledge for long-term retrieval. (Memory + Emit)
func (k *KnowledgeGrammar) Remember(
	ctx context.Context, source types.ActorID, knowledge string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return k.g.Emit(ctx, source, "remember: "+knowledge, convID, causes, signer)
}

// Challenge presents counter-evidence to a claim. (Narrative + Challenge)
func (k *KnowledgeGrammar) Challenge(
	ctx context.Context, source types.ActorID, counterEvidence string,
	claim types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	_, flag, err := k.g.Challenge(ctx, source, "challenge: "+counterEvidence, claim, convID, signer)
	if err != nil {
		return event.Event{}, err
	}
	return flag, nil
}

// DetectBias identifies systematic distortion. (Bias + Annotate)
func (k *KnowledgeGrammar) DetectBias(
	ctx context.Context, source types.ActorID, target types.EventID,
	bias string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return k.g.Annotate(ctx, source, target, "bias", bias, convID, signer)
}

// Correct fixes an error and propagates correction. (Correction + Derive)
func (k *KnowledgeGrammar) Correct(
	ctx context.Context, source types.ActorID, correction string,
	original types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return k.g.Derive(ctx, source, "correct: "+correction, original, convID, signer)
}

// Trace tracks a claim to its original source. (Provenance + Annotate)
func (k *KnowledgeGrammar) Trace(
	ctx context.Context, source types.ActorID, target types.EventID,
	provenance string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return k.g.Annotate(ctx, source, target, "provenance", provenance, convID, signer)
}

// Recall retrieves stored knowledge. (Memory recalled + Emit)
func (k *KnowledgeGrammar) Recall(
	ctx context.Context, source types.ActorID, query string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return k.g.Emit(ctx, source, "recall: "+query, convID, causes, signer)
}

// Learn updates behaviour based on new knowledge. (Learning + Emit)
func (k *KnowledgeGrammar) Learn(
	ctx context.Context, source types.ActorID, learning string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return k.g.Emit(ctx, source, "learn: "+learning, convID, causes, signer)
}

// --- Named Functions (6) ---

// Retract withdraws a claim with chain repair. (Challenge self + Correct)
func (k *KnowledgeGrammar) Retract(
	ctx context.Context, source types.ActorID, claim types.EventID,
	reason string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return k.g.Retract(ctx, source, claim, reason, convID, signer)
}

// FactCheckResult holds the events produced by a FactCheck.
type FactCheckResult struct {
	Provenance event.Event
	BiasCheck  event.Event
	Verdict    event.Event
}

// FactCheck performs full provenance and bias check: Trace + DetectBias + verdict.
func (k *KnowledgeGrammar) FactCheck(
	ctx context.Context, checker types.ActorID,
	claim types.EventID, provenance string, biasAnalysis string, verdict string,
	convID types.ConversationID, signer event.Signer,
) (FactCheckResult, error) {
	trace, err := k.Trace(ctx, checker, claim, provenance, convID, signer)
	if err != nil {
		return FactCheckResult{}, fmt.Errorf("fact-check/trace: %w", err)
	}

	bias, err := k.DetectBias(ctx, checker, claim, biasAnalysis, convID, signer)
	if err != nil {
		return FactCheckResult{}, fmt.Errorf("fact-check/bias: %w", err)
	}

	verdictEv, err := k.g.Merge(ctx, checker, "fact-check: "+verdict,
		[]types.EventID{trace.ID(), bias.ID()}, convID, signer)
	if err != nil {
		return FactCheckResult{}, fmt.Errorf("fact-check/verdict: %w", err)
	}

	return FactCheckResult{Provenance: trace, BiasCheck: bias, Verdict: verdictEv}, nil
}

// VerifyResult holds the events produced by a Verify.
type VerifyResult struct {
	Claim         event.Event
	Provenance    event.Event
	Corroboration event.Event
}

// Verify validates a claim with provenance and corroboration: Claim + Trace + Claim.
func (k *KnowledgeGrammar) Verify(
	ctx context.Context, source types.ActorID,
	claim string, provenance string, corroboration string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (VerifyResult, error) {
	claimEv, err := k.Claim(ctx, source, claim, causes, convID, signer)
	if err != nil {
		return VerifyResult{}, fmt.Errorf("verify/claim: %w", err)
	}

	trace, err := k.Trace(ctx, source, claimEv.ID(), provenance, convID, signer)
	if err != nil {
		return VerifyResult{}, fmt.Errorf("verify/trace: %w", err)
	}

	corroborate, err := k.Claim(ctx, source, "corroborate: "+corroboration, []types.EventID{trace.ID()}, convID, signer)
	if err != nil {
		return VerifyResult{}, fmt.Errorf("verify/corroborate: %w", err)
	}

	return VerifyResult{Claim: claimEv, Provenance: trace, Corroboration: corroborate}, nil
}

// SurveyResult holds the events produced by a Survey.
type SurveyResult struct {
	Recalls     []event.Event
	Abstraction event.Event
	Synthesis   event.Event
}

// Survey aggregates knowledge: Recall (batch) + Abstract + Claim (synthesis).
func (k *KnowledgeGrammar) Survey(
	ctx context.Context, source types.ActorID,
	queries []string, generalization string, synthesis string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (SurveyResult, error) {
	if len(queries) < 2 {
		return SurveyResult{}, fmt.Errorf("survey: requires at least two queries")
	}

	result := SurveyResult{}
	recallIDs := make([]types.EventID, 0, len(queries))
	for i, query := range queries {
		recall, err := k.Recall(ctx, source, query, causes, convID, signer)
		if err != nil {
			return SurveyResult{}, fmt.Errorf("survey/recall[%d]: %w", i, err)
		}
		result.Recalls = append(result.Recalls, recall)
		recallIDs = append(recallIDs, recall.ID())
	}

	abstract, err := k.Abstract(ctx, source, generalization, recallIDs, convID, signer)
	if err != nil {
		return SurveyResult{}, fmt.Errorf("survey/abstract: %w", err)
	}
	result.Abstraction = abstract

	synthesisClaim, err := k.Claim(ctx, source, "synthesis: "+synthesis, []types.EventID{abstract.ID()}, convID, signer)
	if err != nil {
		return SurveyResult{}, fmt.Errorf("survey/synthesis: %w", err)
	}
	result.Synthesis = synthesisClaim

	return result, nil
}

// KnowledgeBaseResult holds the events produced by a KnowledgeBase.
type KnowledgeBaseResult struct {
	Claims     []event.Event
	Categories []event.Event
	Memory     event.Event
}

// KnowledgeBase creates and organises knowledge: Claim + Categorize + Remember (batch).
func (k *KnowledgeGrammar) KnowledgeBase(
	ctx context.Context, source types.ActorID,
	claims []string, taxonomies []string, memoryLabel string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (KnowledgeBaseResult, error) {
	if len(claims) != len(taxonomies) {
		return KnowledgeBaseResult{}, fmt.Errorf("knowledge-base: claims and taxonomies must have equal length")
	}

	result := KnowledgeBaseResult{}
	claimIDs := make([]types.EventID, 0, len(claims))
	for i, c := range claims {
		claimEv, err := k.Claim(ctx, source, c, causes, convID, signer)
		if err != nil {
			return KnowledgeBaseResult{}, fmt.Errorf("knowledge-base/claim[%d]: %w", i, err)
		}
		result.Claims = append(result.Claims, claimEv)

		cat, err := k.Categorize(ctx, source, claimEv.ID(), taxonomies[i], convID, signer)
		if err != nil {
			return KnowledgeBaseResult{}, fmt.Errorf("knowledge-base/categorize[%d]: %w", i, err)
		}
		result.Categories = append(result.Categories, cat)
		claimIDs = append(claimIDs, cat.ID())
	}

	memory, err := k.Remember(ctx, source, memoryLabel, claimIDs, convID, signer)
	if err != nil {
		return KnowledgeBaseResult{}, fmt.Errorf("knowledge-base/remember: %w", err)
	}
	result.Memory = memory

	return result, nil
}

// TransferResult holds the events produced by a Transfer.
type TransferResult struct {
	Recall event.Event
	Encode event.Event
	Learn  event.Event
}

// Transfer moves knowledge to a new context: Recall + Encode + Learn.
func (k *KnowledgeGrammar) Transfer(
	ctx context.Context, source types.ActorID,
	query string, encoding string, learning string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (TransferResult, error) {
	recall, err := k.Recall(ctx, source, query, causes, convID, signer)
	if err != nil {
		return TransferResult{}, fmt.Errorf("transfer/recall: %w", err)
	}

	encode, err := k.Encode(ctx, source, encoding, recall.ID(), convID, signer)
	if err != nil {
		return TransferResult{}, fmt.Errorf("transfer/encode: %w", err)
	}

	learn, err := k.Learn(ctx, source, learning, []types.EventID{encode.ID()}, convID, signer)
	if err != nil {
		return TransferResult{}, fmt.Errorf("transfer/learn: %w", err)
	}

	return TransferResult{Recall: recall, Encode: encode, Learn: learn}, nil
}
