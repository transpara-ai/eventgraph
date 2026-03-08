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

// --- Named Functions (2) ---

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
