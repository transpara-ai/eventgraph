package compositions

import (
	"context"
	"fmt"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/grammar"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// BuildGrammar provides Layer 5 (Technology) composition operations.
// 12 operations + 5 named functions for development, CI/CD, and artefact lifecycle.
type BuildGrammar struct {
	g *grammar.Grammar
}

// NewBuildGrammar creates a BuildGrammar bound to the given base grammar.
func NewBuildGrammar(g *grammar.Grammar) *BuildGrammar {
	return &BuildGrammar{g: g}
}

// --- Operations (12) ---

// Build creates a new artefact with provenance. (Create + Emit)
func (b *BuildGrammar) Build(
	ctx context.Context, source types.ActorID, artefact string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "build: "+artefact, convID, causes, signer)
}

// Version releases a new version of an artefact. (Create version + Derive)
func (b *BuildGrammar) Version(
	ctx context.Context, source types.ActorID, version string,
	previous types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Derive(ctx, source, "version: "+version, previous, convID, signer)
}

// Ship deploys an artefact for use. (Tool registered + Emit)
func (b *BuildGrammar) Ship(
	ctx context.Context, source types.ActorID, deployment string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "ship: "+deployment, convID, causes, signer)
}

// Sunset deprecates an artefact with migration path. (Deprecation + Annotate)
func (b *BuildGrammar) Sunset(
	ctx context.Context, source types.ActorID, target types.EventID,
	migration string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Annotate(ctx, source, target, "deprecated", migration, convID, signer)
}

// Define establishes a repeatable workflow. (Workflow + Emit)
func (b *BuildGrammar) Define(
	ctx context.Context, source types.ActorID, workflow string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "define: "+workflow, convID, causes, signer)
}

// Automate converts a manual step to mechanical. (Automation + Derive)
func (b *BuildGrammar) Automate(
	ctx context.Context, source types.ActorID, automation string,
	workflow types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Derive(ctx, source, "automate: "+automation, workflow, convID, signer)
}

// Test runs verification against an artefact. (Testing + Emit)
func (b *BuildGrammar) Test(
	ctx context.Context, source types.ActorID, results string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "test: "+results, convID, causes, signer)
}

// Review performs peer assessment. (Review + Respond)
func (b *BuildGrammar) Review(
	ctx context.Context, source types.ActorID, assessment string,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Respond(ctx, source, "review: "+assessment, target, convID, signer)
}

// Measure assesses quality against criteria. (Quality + Annotate)
func (b *BuildGrammar) Measure(
	ctx context.Context, source types.ActorID, target types.EventID,
	scores string, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Annotate(ctx, source, target, "quality", scores, convID, signer)
}

// Feedback provides structured input on outcomes. (Feedback + Respond)
func (b *BuildGrammar) Feedback(
	ctx context.Context, source types.ActorID, feedback string,
	target types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Respond(ctx, source, "feedback: "+feedback, target, convID, signer)
}

// Iterate improves through repeated refinement. (Iteration + Derive)
func (b *BuildGrammar) Iterate(
	ctx context.Context, source types.ActorID, improvement string,
	previous types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Derive(ctx, source, "iterate: "+improvement, previous, convID, signer)
}

// Innovate creates something genuinely new. (Innovation + Emit)
func (b *BuildGrammar) Innovate(
	ctx context.Context, source types.ActorID, innovation string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (event.Event, error) {
	return b.g.Emit(ctx, source, "innovate: "+innovation, convID, causes, signer)
}

// --- Named Functions (5) ---

// PipelineResult holds the events produced by a Pipeline.
type PipelineResult struct {
	Definition event.Event
	TestResult event.Event
	Metrics    event.Event
	Deployment event.Event
}

// Pipeline runs a full CI/CD flow: Define + Test + Measure + Ship.
func (b *BuildGrammar) Pipeline(
	ctx context.Context, source types.ActorID,
	workflow string, testResults string, metrics string, deployment string,
	causes []types.EventID, convID types.ConversationID, signer event.Signer,
) (PipelineResult, error) {
	def, err := b.Define(ctx, source, workflow, causes, convID, signer)
	if err != nil {
		return PipelineResult{}, fmt.Errorf("pipeline/define: %w", err)
	}

	test, err := b.Test(ctx, source, testResults, []types.EventID{def.ID()}, convID, signer)
	if err != nil {
		return PipelineResult{}, fmt.Errorf("pipeline/test: %w", err)
	}

	measure, err := b.Measure(ctx, source, test.ID(), metrics, convID, signer)
	if err != nil {
		return PipelineResult{}, fmt.Errorf("pipeline/measure: %w", err)
	}

	ship, err := b.Ship(ctx, source, deployment, []types.EventID{measure.ID()}, convID, signer)
	if err != nil {
		return PipelineResult{}, fmt.Errorf("pipeline/ship: %w", err)
	}

	return PipelineResult{Definition: def, TestResult: test, Metrics: measure, Deployment: ship}, nil
}

// PostMortemResult holds the events produced by a PostMortem.
type PostMortemResult struct {
	Feedback    []event.Event
	Analysis    event.Event
	Improvements event.Event
}

// PostMortem learns from failure: Feedback (batch) + Measure + Define (improvements).
func (b *BuildGrammar) PostMortem(
	ctx context.Context, lead types.ActorID,
	contributors []types.ActorID, feedbacks []string,
	analysis string, improvements string,
	incident types.EventID, convID types.ConversationID, signer event.Signer,
) (PostMortemResult, error) {
	if len(contributors) != len(feedbacks) {
		return PostMortemResult{}, fmt.Errorf("post-mortem: contributors and feedbacks must have equal length")
	}

	result := PostMortemResult{}
	fbIDs := make([]types.EventID, 0, len(contributors))
	for i, contrib := range contributors {
		fb, err := b.Feedback(ctx, contrib, feedbacks[i], incident, convID, signer)
		if err != nil {
			return PostMortemResult{}, fmt.Errorf("post-mortem/feedback[%d]: %w", i, err)
		}
		result.Feedback = append(result.Feedback, fb)
		fbIDs = append(fbIDs, fb.ID())
	}

	analysisEv, err := b.g.Merge(ctx, lead, "post-mortem: "+analysis, fbIDs, convID, signer)
	if err != nil {
		return PostMortemResult{}, fmt.Errorf("post-mortem/analysis: %w", err)
	}
	result.Analysis = analysisEv

	improve, err := b.g.Derive(ctx, lead, "improvements: "+improvements, analysisEv.ID(), convID, signer)
	if err != nil {
		return PostMortemResult{}, fmt.Errorf("post-mortem/improvements: %w", err)
	}
	result.Improvements = improve

	return result, nil
}
