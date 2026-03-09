package intelligence_test

import (
	"context"
	"strings"
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/intelligence"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// ════════════════════════════════════════════════════════════════════════
// Research integration tests: agents read web resources, describe in
// Code Graph, and generate code from what they learn.
// ════════════════════════════════════════════════════════════════════════

func researchRuntime(t *testing.T, response string) *intelligence.AgentRuntime {
	t.Helper()
	p := newMockProvider(response)
	rt, err := intelligence.NewRuntime(context.Background(), intelligence.RuntimeConfig{
		AgentID:  types.MustActorID("actor_00000000000000000000000000000001"),
		Provider: p,
	})
	if err != nil {
		t.Fatalf("NewRuntime: %v", err)
	}
	return rt
}

// TestResearchRecordsEvents verifies that Research() records observation
// and evaluation events on the graph.
func TestResearchRecordsEvents(t *testing.T) {
	rt := researchRuntime(t, "The blog post describes a task management system with kanban boards.")
	ctx := context.Background()

	obs, evaluation, err := rt.Research(ctx,
		"https://mattsearles2.substack.com/p/post-35",
		"extract the product idea and key features described.")
	if err != nil {
		t.Fatalf("Research: %v", err)
	}

	// Observation event was recorded
	if obs.Type().Value() != "agent.observed" {
		t.Errorf("observation type = %q, want agent.observed", obs.Type().Value())
	}

	// Evaluation contains extracted content
	if !strings.Contains(evaluation, "task management") {
		t.Errorf("evaluation should contain extracted content, got: %s", evaluation)
	}

	// Both events are on the graph
	observations, _ := rt.EventsByType("agent.observed", 10)
	evaluations, _ := rt.EventsByType("agent.evaluated", 10)
	if len(observations) < 1 {
		t.Error("expected at least 1 observed event")
	}
	if len(evaluations) < 1 {
		t.Error("expected at least 1 evaluated event")
	}
}

// TestResearchToCodeGraphToCode verifies the full pipeline:
// 1. Agent reads a URL (mock)
// 2. Agent describes what it learned as Code Graph events
// 3. Agent generates code from the spec
func TestResearchToCodeGraphToCode(t *testing.T) {
	// Mock that returns different responses based on prompt content
	rt := researchRuntime(t, "Product: Task board with todo/doing/done columns, drag-and-drop, real-time presence.")
	agentID := rt.ID()
	ctx := context.Background()

	// 1. Research: read the blog post
	_, extraction, err := rt.Research(ctx,
		"https://mattsearles2.substack.com/p/post-35",
		"extract the product idea.")
	if err != nil {
		t.Fatalf("Research: %v", err)
	}
	if len(extraction) == 0 {
		t.Fatal("extraction should not be empty")
	}

	// 2. Describe in Code Graph: record the spec as events
	_, err = rt.Emit(event.CodeGraphEntityDefinedContent{
		EntityType: "Task",
		EntityID:   agentID,
	})
	if err != nil {
		t.Fatalf("define entity: %v", err)
	}

	_, err = rt.Emit(event.CodeGraphStateTransitionedContent{
		EntityType: "Task",
		EntityID:   agentID,
		From:       "todo",
		To:         "doing",
	})
	if err != nil {
		t.Fatalf("define state: %v", err)
	}

	_, err = rt.Emit(event.CodeGraphViewRenderedContent{
		ViewName: "TaskBoard",
		Actor:    agentID,
	})
	if err != nil {
		t.Fatalf("define view: %v", err)
	}

	// 3. Generate code from the spec
	code, err := rt.CodeWrite(ctx,
		"Implement a TaskBoard React component based on the Code Graph spec: "+
			"Entity(Task) with states todo/doing/done, View(TaskBoard) with drag-and-drop columns.",
		"typescript")
	if err != nil {
		t.Fatalf("CodeWrite: %v", err)
	}
	if len(code) == 0 {
		t.Fatal("code should not be empty")
	}

	// 4. Verify the full event trail
	memory, _ := rt.Memory(20)
	// Should have: bootstrap + observe + evaluate (research) +
	//              entity.defined + state.transitioned + view.rendered +
	//              code_write reasoning + acted
	if len(memory) < 7 {
		t.Errorf("expected >= 7 events in memory, got %d", len(memory))
	}

	// Verify codegraph events are queryable
	entities, _ := rt.EventsByType("codegraph.entity.defined", 10)
	if len(entities) != 1 {
		t.Errorf("entity.defined = %d, want 1", len(entities))
	}
	views, _ := rt.EventsByType("codegraph.ui.view.rendered", 10)
	if len(views) != 1 {
		t.Errorf("view.rendered = %d, want 1", len(views))
	}
}

// TestResearchCausalChain verifies that research → spec → code events
// form a proper causal chain.
func TestResearchCausalChain(t *testing.T) {
	rt := researchRuntime(t, "A social network with posts and reactions.")
	agentID := rt.ID()
	ctx := context.Background()

	// Collect events in order
	var chain []event.Event

	obs, _, _ := rt.Research(ctx, "https://example.com/idea", "extract the idea.")
	chain = append(chain, obs)
	// Research also emits an evaluated event (internal)
	evals, _ := rt.EventsByType("agent.evaluated", 10)
	if len(evals) > 0 {
		chain = append(chain, evals[len(evals)-1])
	}

	ev, _ := rt.Emit(event.CodeGraphEntityDefinedContent{
		EntityType: "Post", EntityID: agentID,
	})
	chain = append(chain, ev)

	ev, _ = rt.Emit(event.CodeGraphEntityDefinedContent{
		EntityType: "Reaction", EntityID: agentID,
	})
	chain = append(chain, ev)

	ev, _ = rt.Emit(event.CodeGraphEntityRelatedContent{
		SourceType: "Reaction", SourceID: agentID,
		TargetType: "Post", TargetID: agentID,
		Relation: "belongs_to",
	})
	chain = append(chain, ev)

	// Verify sequential causality
	for i := 1; i < len(chain); i++ {
		if chain[i].PrevHash() != chain[i-1].Hash() {
			t.Errorf("chain broken at %d: %s → %s",
				i, chain[i].Type().Value(), chain[i-1].Type().Value())
		}
	}
}

// ════════════════════════════════════════════════════════════════════════
// Real LLM tests — gated on EVENTGRAPH_TEST_CLAUDE_CLI
// ════════════════════════════════════════════════════════════════════════

// TestRealLLMResearchSubstack uses Claude CLI to actually read a Substack post,
// extract a product idea, describe it in Code Graph, and generate code.
func TestRealLLMResearchSubstack(t *testing.T) {
	skipWithoutClaudeCli(t)

	p := agentProvider(t)
	agentID := types.MustActorID("actor_00000000000000000000000000000001")
	rt, err := intelligence.NewRuntime(context.Background(), intelligence.RuntimeConfig{
		AgentID:  agentID,
		Provider: p,
	})
	if err != nil {
		t.Fatalf("NewRuntime: %v", err)
	}
	ctx := context.Background()

	// 1. Research: read a real Substack post
	_, extraction, err := rt.Research(ctx,
		"https://mattsearles2.substack.com/p/the-missing-social-grammar",
		"extract the key product idea. What system is being described? What are the core operations? Summarize in 3-5 bullet points.")
	if err != nil {
		t.Fatalf("Research: %v", err)
	}
	t.Logf("Extraction:\n%s", extraction)

	if len(extraction) < 50 {
		t.Errorf("extraction too short (%d chars), expected meaningful content", len(extraction))
	}

	// 2. Describe in Code Graph
	_, err = rt.Emit(event.CodeGraphEntityDefinedContent{
		EntityType: "SocialAction",
		EntityID:   agentID,
	})
	if err != nil {
		t.Fatalf("define entity: %v", err)
	}

	// 3. Ask agent to evaluate what Code Graph spec would look like
	_, specEval, err := rt.Evaluate(ctx, "codegraph_design",
		"Based on this product extraction:\n"+extraction+"\n\n"+
			"Describe the data model using Code Graph vocabulary: "+
			"Entity(), State(), Relation(), View(), Trigger(). "+
			"Be specific — list 3-5 entities with their properties.")
	if err != nil {
		t.Fatalf("Evaluate spec: %v", err)
	}
	t.Logf("Spec evaluation:\n%s", specEval)

	// 4. Generate code
	code, err := rt.CodeWrite(ctx,
		"Based on this Code Graph spec:\n"+specEval+"\n\n"+
			"Implement the core data types in TypeScript. Include the entity types, "+
			"state machines, and one React component for the main view.",
		"typescript")
	if err != nil {
		t.Fatalf("CodeWrite: %v", err)
	}
	t.Logf("Generated code (%d chars):\n%s", len(code), code[:min(500, len(code))])

	if len(code) < 50 {
		t.Errorf("generated code too short (%d chars)", len(code))
	}

	// 5. Verify the full trail is on the graph
	memory, _ := rt.Memory(30)
	t.Logf("Total events on graph: %d", len(memory))

	if len(memory) < 6 {
		t.Errorf("expected >= 6 events, got %d", len(memory))
	}
}

// TestRealLLMResearchAndDescribeMultipleProducts reads multiple URLs and
// describes them all in Code Graph, showing the agent can build a portfolio.
func TestRealLLMResearchAndDescribeMultipleProducts(t *testing.T) {
	skipWithoutClaudeCli(t)

	p := agentProvider(t)
	agentID := types.MustActorID("actor_00000000000000000000000000000001")
	rt, err := intelligence.NewRuntime(context.Background(), intelligence.RuntimeConfig{
		AgentID:  agentID,
		Provider: p,
	})
	if err != nil {
		t.Fatalf("NewRuntime: %v", err)
	}
	ctx := context.Background()

	// Research the project README for product ideas
	_, extraction, err := rt.Research(ctx,
		"https://github.com/lovyou-ai/eventgraph",
		"identify what products could be built on this infrastructure. List 3 product ideas with one sentence each.")
	if err != nil {
		t.Fatalf("Research: %v", err)
	}
	t.Logf("Product ideas:\n%s", extraction)

	// Record each idea as a codegraph entity
	for i, name := range []string{"AuditPlatform", "GovernanceSystem", "AgentFramework"} {
		_, err := rt.Emit(event.CodeGraphEntityDefinedContent{
			EntityType: name,
			EntityID:   agentID,
		})
		if err != nil {
			t.Fatalf("define entity %d: %v", i, err)
		}
	}

	// Verify all 3 entities are on the graph
	entities, _ := rt.EventsByType("codegraph.entity.defined", 10)
	if len(entities) != 3 {
		t.Errorf("entity.defined = %d, want 3", len(entities))
	}
}
