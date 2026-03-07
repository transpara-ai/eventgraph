package decision_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/decision"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// mockIntelligence implements IIntelligence for tests.
type mockIntelligence struct {
	content    string
	confidence types.Score
	tokens     int
	err        error
}

func (m *mockIntelligence) Reason(_ context.Context, _ string, _ []event.Event) (decision.Response, error) {
	if m.err != nil {
		return decision.Response{}, m.err
	}
	return decision.NewResponse(m.content, m.confidence, m.tokens), nil
}

func testInput(action string) decision.EvaluateInput {
	return decision.EvaluateInput{
		Action: action,
		Actor:  types.MustActorID("actor_test00000000000000000000001"),
		Context: map[string]any{
			"trust_score": 0.8,
			"event_type":  "code.reviewed",
		},
	}
}

func TestEvaluateSimpleLeaf(t *testing.T) {
	tree := decision.NewDecisionTree(
		decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.95)),
	)

	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("Outcome = %v, want Permit", result.Outcome)
	}
	if result.Confidence.Value() != 0.95 {
		t.Errorf("Confidence = %v, want 0.95", result.Confidence.Value())
	}
	if result.UsedLLM {
		t.Error("should not have used LLM")
	}
	if len(result.Path) != 0 {
		t.Errorf("Path should be empty, got %d steps", len(result.Path))
	}
}

func TestEvaluateMechanicalBranching(t *testing.T) {
	// Tree: if action == "deploy" -> Deny, else -> Permit
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("action"),
			Operator: event.ConditionOperatorEquals,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{String: types.Some("deploy")},
				Child: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(1.0)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.9)),
	})

	// Test matching branch
	result, err := decision.Evaluate(context.Background(), tree, testInput("deploy"), types.None[decision.IIntelligence]())
	if err != nil {
		t.Fatalf("Evaluate deploy: %v", err)
	}
	if result.Outcome != event.DecisionOutcomeDeny {
		t.Errorf("deploy: Outcome = %v, want Deny", result.Outcome)
	}
	if len(result.Path) != 1 {
		t.Fatalf("deploy: Path should have 1 step, got %d", len(result.Path))
	}

	// Test default branch
	result, err = decision.Evaluate(context.Background(), tree, testInput("review"), types.None[decision.IIntelligence]())
	if err != nil {
		t.Fatalf("Evaluate review: %v", err)
	}
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("review: Outcome = %v, want Permit", result.Outcome)
	}
}

func TestEvaluateNumericCondition(t *testing.T) {
	// Tree: if context.trust_score > 0.5 -> Permit, else -> Deny
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.trust_score"),
			Operator: event.ConditionOperatorGreaterThan,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{Number: types.Some(0.5)},
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.9)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.9)),
	})

	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("Outcome = %v, want Permit (trust_score 0.8 > 0.5)", result.Outcome)
	}
}

func TestEvaluateLessThan(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.trust_score"),
			Operator: event.ConditionOperatorLessThan,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{Number: types.Some(0.5)},
				Child: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.9)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.9)),
	})

	// trust_score is 0.8, not < 0.5, so should take default
	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("Outcome = %v, want Permit", result.Outcome)
	}
}

func TestEvaluateExists(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.trust_score"),
			Operator: event.ConditionOperatorExists,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{Boolean: types.Some(true)},
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.9)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.5)),
	})

	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("Outcome = %v, want Permit (field exists)", result.Outcome)
	}
}

func TestEvaluateExistsFieldMissing(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.nonexistent"),
			Operator: event.ConditionOperatorExists,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{Boolean: types.Some(true)},
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.9)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.5)),
	})

	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Outcome != event.DecisionOutcomeDeny {
		t.Errorf("Outcome = %v, want Deny (field missing)", result.Outcome)
	}
}

func TestEvaluatePatternMatch(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.event_type"),
			Operator: event.ConditionOperatorMatches,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{String: types.Some("code.*")},
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.9)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDefer, types.MustScore(0.5)),
	})

	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("Outcome = %v, want Permit (code.reviewed matches code.*)", result.Outcome)
	}
}

func TestEvaluateDeepTree(t *testing.T) {
	// Two-level tree: action == "deploy" -> trust > 0.5 -> Permit / Deny
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("action"),
			Operator: event.ConditionOperatorEquals,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{String: types.Some("deploy")},
				Child: &decision.InternalNode{
					Condition: event.Condition{
						Field:    types.MustFieldPath("context.trust_score"),
						Operator: event.ConditionOperatorGreaterThan,
					},
					Branches: []decision.Branch{
						{
							Match: event.MatchValue{Number: types.Some(0.5)},
							Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.95)),
						},
					},
					Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.9)),
				},
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDefer, types.MustScore(0.5)),
	})

	result, err := decision.Evaluate(context.Background(), tree, testInput("deploy"), types.None[decision.IIntelligence]())
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("Outcome = %v, want Permit", result.Outcome)
	}
	if len(result.Path) != 2 {
		t.Errorf("Path should have 2 steps, got %d", len(result.Path))
	}
}

func TestEvaluateLLMLeafNoIntelligence(t *testing.T) {
	tree := decision.NewDecisionTree(decision.NewLLMLeaf(types.MustScore(0.5)))

	_, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if err == nil {
		t.Fatal("expected error when LLM leaf reached without intelligence")
	}
	var insuffAuth *decision.InsufficientAuthorityError
	if !errors.As(err, &insuffAuth) {
		t.Errorf("expected InsufficientAuthorityError, got %T: %v", err, err)
	}
}

func TestEvaluateLLMLeafWithIntelligence(t *testing.T) {
	tree := decision.NewDecisionTree(decision.NewLLMLeaf(types.MustScore(0.5)))

	intel := &mockIntelligence{
		content:    "permit this action",
		confidence: types.MustScore(0.9),
		tokens:     50,
	}

	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.Some[decision.IIntelligence](intel))
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("Outcome = %v, want Permit", result.Outcome)
	}
	if !result.UsedLLM {
		t.Error("should have used LLM")
	}
	if result.Confidence.Value() != 0.9 {
		t.Errorf("Confidence = %v, want 0.9", result.Confidence.Value())
	}
}

func TestEvaluateLLMLeafDeny(t *testing.T) {
	tree := decision.NewDecisionTree(decision.NewLLMLeaf(types.MustScore(0.5)))

	intel := &mockIntelligence{
		content:    "deny access",
		confidence: types.MustScore(0.85),
		tokens:     30,
	}

	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.Some[decision.IIntelligence](intel))
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Outcome != event.DecisionOutcomeDeny {
		t.Errorf("Outcome = %v, want Deny", result.Outcome)
	}
}

func TestEvaluateLLMLeafEscalate(t *testing.T) {
	tree := decision.NewDecisionTree(decision.NewLLMLeaf(types.MustScore(0.5)))

	intel := &mockIntelligence{
		content:    "escalate to human",
		confidence: types.MustScore(0.7),
		tokens:     20,
	}

	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.Some[decision.IIntelligence](intel))
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Outcome != event.DecisionOutcomeEscalate {
		t.Errorf("Outcome = %v, want Escalate", result.Outcome)
	}
}

func TestEvaluateLLMLeafDefer(t *testing.T) {
	tree := decision.NewDecisionTree(decision.NewLLMLeaf(types.MustScore(0.5)))

	intel := &mockIntelligence{
		content:    "I'm not sure what to do",
		confidence: types.MustScore(0.3),
		tokens:     40,
	}

	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.Some[decision.IIntelligence](intel))
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Outcome != event.DecisionOutcomeDefer {
		t.Errorf("Outcome = %v, want Defer", result.Outcome)
	}
}

func TestEvaluateLLMError(t *testing.T) {
	tree := decision.NewDecisionTree(decision.NewLLMLeaf(types.MustScore(0.5)))

	intel := &mockIntelligence{err: errors.New("model unavailable")}

	_, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.Some[decision.IIntelligence](intel))
	if err == nil {
		t.Fatal("expected error from intelligence failure")
	}
}

func TestEvaluateSemanticNoIntelligence(t *testing.T) {
	// Semantic condition falls through to default when no intelligence
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:     types.MustFieldPath("context.tone"),
			Operator:  event.ConditionOperatorSemantic,
			Threshold: types.Some(types.MustScore(0.7)),
			Prompt:    types.Some("Is this message hostile?"),
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{Boolean: types.Some(true)},
				Child: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.9)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.5)),
	})

	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("Outcome = %v, want Permit (default when no intelligence)", result.Outcome)
	}
}

func TestEvaluateSemanticWithIntelligence(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:     types.MustFieldPath("context.tone"),
			Operator:  event.ConditionOperatorSemantic,
			Threshold: types.Some(types.MustScore(0.7)),
			Prompt:    types.Some("Is this message hostile?"),
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{Boolean: types.Some(true)},
				Child: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.9)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.5)),
	})

	intel := &mockIntelligence{
		content:    "yes, hostile",
		confidence: types.MustScore(0.85), // above threshold of 0.7
		tokens:     20,
	}

	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.Some[decision.IIntelligence](intel))
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Outcome != event.DecisionOutcomeDeny {
		t.Errorf("Outcome = %v, want Deny (hostile detected)", result.Outcome)
	}
}

func TestEvaluateSemanticBelowThreshold(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:     types.MustFieldPath("context.tone"),
			Operator:  event.ConditionOperatorSemantic,
			Threshold: types.Some(types.MustScore(0.7)),
			Prompt:    types.Some("Is this message hostile?"),
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{Boolean: types.Some(true)},
				Child: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.9)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.5)),
	})

	intel := &mockIntelligence{
		content:    "maybe slightly rude",
		confidence: types.MustScore(0.4), // below threshold
		tokens:     15,
	}

	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.Some[decision.IIntelligence](intel))
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("Outcome = %v, want Permit (below threshold)", result.Outcome)
	}
}

func TestEvaluateSemanticIntelligenceError(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:     types.MustFieldPath("context.tone"),
			Operator:  event.ConditionOperatorSemantic,
			Threshold: types.Some(types.MustScore(0.7)),
			Prompt:    types.Some("Is this hostile?"),
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{Boolean: types.Some(true)},
				Child: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.9)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.5)),
	})

	intel := &mockIntelligence{err: errors.New("timeout")}

	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.Some[decision.IIntelligence](intel))
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	// Should fall through to default on error
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("Outcome = %v, want Permit (fallthrough on error)", result.Outcome)
	}
}

func TestTreeStatsTracking(t *testing.T) {
	tree := decision.NewDecisionTree(
		decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.9)),
	)

	for i := 0; i < 5; i++ {
		decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	}

	if tree.Stats.TotalHits != 5 {
		t.Errorf("TotalHits = %d, want 5", tree.Stats.TotalHits)
	}
	if tree.Stats.MechanicalHits != 5 {
		t.Errorf("MechanicalHits = %d, want 5", tree.Stats.MechanicalHits)
	}
	if tree.Stats.LLMHits != 0 {
		t.Errorf("LLMHits = %d, want 0", tree.Stats.LLMHits)
	}
}

func TestTreeStatsLLMTracking(t *testing.T) {
	tree := decision.NewDecisionTree(decision.NewLLMLeaf(types.MustScore(0.5)))

	intel := &mockIntelligence{
		content:    "permit",
		confidence: types.MustScore(0.9),
		tokens:     100,
	}

	for i := 0; i < 3; i++ {
		decision.Evaluate(context.Background(), tree, testInput("test"), types.Some[decision.IIntelligence](intel))
	}

	if tree.Stats.LLMHits != 3 {
		t.Errorf("LLMHits = %d, want 3", tree.Stats.LLMHits)
	}
	if tree.Stats.TotalTokens != 300 {
		t.Errorf("TotalTokens = %d, want 300", tree.Stats.TotalTokens)
	}
}

func TestNoOpIntelligence(t *testing.T) {
	noop := decision.NoOpIntelligence{}
	_, err := noop.Reason(context.Background(), "test", nil)
	if err == nil {
		t.Fatal("expected error from NoOpIntelligence")
	}
	var unavailable *decision.IntelligenceUnavailableError
	if !errors.As(err, &unavailable) {
		t.Errorf("expected IntelligenceUnavailableError, got %T", err)
	}
}

func TestDecisionErrors(t *testing.T) {
	tests := []struct {
		name string
		err  decision.DecisionError
	}{
		{"ActorNotFound", &decision.ActorNotFoundError{ID: types.MustActorID("actor_test00000000000000000000001")}},
		{"InsufficientAuthority", &decision.InsufficientAuthorityError{
			Actor:    types.MustActorID("actor_test00000000000000000000001"),
			Action:   "deploy",
			Required: event.AuthorityLevelRequired,
		}},
		{"TrustBelowThreshold", &decision.TrustBelowThresholdError{
			Actor:    types.MustActorID("actor_test00000000000000000000001"),
			Score:    types.MustScore(0.3),
			Required: types.MustScore(0.5),
		}},
		{"CausesRequired", &decision.CausesRequiredError{Action: "deploy"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() == "" {
				t.Error("error message should not be empty")
			}
		})
	}
}

func TestResponseGetters(t *testing.T) {
	r := decision.NewResponse("test content", types.MustScore(0.8), 42)
	if r.Content() != "test content" {
		t.Errorf("Content = %q, want test content", r.Content())
	}
	if r.Confidence().Value() != 0.8 {
		t.Errorf("Confidence = %v, want 0.8", r.Confidence().Value())
	}
	if r.TokensUsed() != 42 {
		t.Errorf("TokensUsed = %d, want 42", r.TokensUsed())
	}
}

func TestNewDecisionTreeVersion(t *testing.T) {
	tree := decision.NewDecisionTree(decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.9)))
	if tree.Version != 1 {
		t.Errorf("Version = %d, want 1", tree.Version)
	}
}

func TestExtractFieldAction(t *testing.T) {
	// Tested indirectly through tree evaluation on "action" field
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("action"),
			Operator: event.ConditionOperatorEquals,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{String: types.Some("deploy")},
				Child: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(1.0)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.5)),
	})

	result, _ := decision.Evaluate(context.Background(), tree, testInput("deploy"), types.None[decision.IIntelligence]())
	if result.Outcome != event.DecisionOutcomeDeny {
		t.Errorf("action field extraction failed: got %v, want Deny", result.Outcome)
	}
}

func TestExtractFieldActor(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("actor"),
			Operator: event.ConditionOperatorEquals,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{String: types.Some("actor_test00000000000000000000001")},
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(1.0)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.5)),
	})

	result, _ := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("actor field extraction failed: got %v, want Permit", result.Outcome)
	}
}

func TestExtractFieldContextNil(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.missing"),
			Operator: event.ConditionOperatorExists,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{Boolean: types.Some(true)},
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(1.0)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.5)),
	})

	input := decision.EvaluateInput{
		Action:  "test",
		Actor:   types.MustActorID("actor_test00000000000000000000001"),
		Context: nil,
	}

	result, _ := decision.Evaluate(context.Background(), tree, input, types.None[decision.IIntelligence]())
	if result.Outcome != event.DecisionOutcomeDeny {
		t.Errorf("nil context should take default: got %v, want Deny", result.Outcome)
	}
}

func TestEqualsMatchNumber(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.trust_score"),
			Operator: event.ConditionOperatorEquals,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{Number: types.Some(0.8)},
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(1.0)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.5)),
	})

	result, _ := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("number equals failed: got %v, want Permit", result.Outcome)
	}
}

func TestEqualsMatchBoolean(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.approved"),
			Operator: event.ConditionOperatorEquals,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{Boolean: types.Some(true)},
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(1.0)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.5)),
	})

	input := decision.EvaluateInput{
		Action:  "test",
		Actor:   types.MustActorID("actor_test00000000000000000000001"),
		Context: map[string]any{"approved": true},
	}

	result, _ := decision.Evaluate(context.Background(), tree, input, types.None[decision.IIntelligence]())
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("boolean equals failed: got %v, want Permit", result.Outcome)
	}
}

func TestEqualsMatchEventType(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.event_type"),
			Operator: event.ConditionOperatorEquals,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{EventType: types.Some(types.MustEventType("code.reviewed"))},
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(1.0)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.5)),
	})

	result, _ := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("event type equals failed: got %v, want Permit", result.Outcome)
	}
}

func TestEqualsNoMatch(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.trust_score"),
			Operator: event.ConditionOperatorEquals,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{}, // no match values set
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(1.0)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.5)),
	})

	result, _ := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if result.Outcome != event.DecisionOutcomeDeny {
		t.Errorf("empty match should take default: got %v, want Deny", result.Outcome)
	}
}

func TestPatternMatchWildcard(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.event_type"),
			Operator: event.ConditionOperatorMatches,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{String: types.Some("*")},
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(1.0)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.5)),
	})

	result, _ := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("wildcard match failed: got %v, want Permit", result.Outcome)
	}
}

func TestPatternMatchExact(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.event_type"),
			Operator: event.ConditionOperatorMatches,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{String: types.Some("code.reviewed")},
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(1.0)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.5)),
	})

	result, _ := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("exact pattern match failed: got %v, want Permit", result.Outcome)
	}
}

func TestPatternMatchNoStringMatch(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.event_type"),
			Operator: event.ConditionOperatorMatches,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{Number: types.Some(42.0)}, // wrong type for pattern match
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(1.0)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.5)),
	})

	result, _ := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if result.Outcome != event.DecisionOutcomeDeny {
		t.Errorf("non-string pattern match should take default: got %v, want Deny", result.Outcome)
	}
}

func TestPatternMatchNonStringValue(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.trust_score"),
			Operator: event.ConditionOperatorMatches,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{String: types.Some("*")},
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(1.0)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.5)),
	})

	result, _ := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if result.Outcome != event.DecisionOutcomeDeny {
		t.Errorf("pattern match on non-string value should take default: got %v, want Deny", result.Outcome)
	}
}

func TestNumericCompareNoNumber(t *testing.T) {
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.trust_score"),
			Operator: event.ConditionOperatorGreaterThan,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{String: types.Some("not a number")},
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(1.0)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.5)),
	})

	result, _ := decision.Evaluate(context.Background(), tree, testInput("test"), types.None[decision.IIntelligence]())
	if result.Outcome != event.DecisionOutcomeDeny {
		t.Errorf("numeric compare with non-number match should take default: got %v, want Deny", result.Outcome)
	}
}

func TestToFloat64IntTypes(t *testing.T) {
	// Test int and int64 conversions via numeric comparison
	tests := []struct {
		name  string
		value any
	}{
		{"int", 5},
		{"int64", int64(5)},
		{"float32", float32(5.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := decision.NewDecisionTree(&decision.InternalNode{
				Condition: event.Condition{
					Field:    types.MustFieldPath("context.count"),
					Operator: event.ConditionOperatorGreaterThan,
				},
				Branches: []decision.Branch{
					{
						Match: event.MatchValue{Number: types.Some(3.0)},
						Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(1.0)),
					},
				},
				Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.5)),
			})

			input := decision.EvaluateInput{
				Action:  "test",
				Actor:   types.MustActorID("actor_test00000000000000000000001"),
				Context: map[string]any{"count": tt.value},
			}

			result, _ := decision.Evaluate(context.Background(), tree, input, types.None[decision.IIntelligence]())
			if result.Outcome != event.DecisionOutcomePermit {
				t.Errorf("%s conversion failed: got %v, want Permit", tt.name, result.Outcome)
			}
		})
	}
}

func TestExtractFieldTopLevelFromContext(t *testing.T) {
	// Field path without "context." prefix should also check context map
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("status"),
			Operator: event.ConditionOperatorEquals,
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{String: types.Some("active")},
				Child: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(1.0)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.5)),
	})

	input := decision.EvaluateInput{
		Action:  "test",
		Actor:   types.MustActorID("actor_test00000000000000000000001"),
		Context: map[string]any{"status": "active"},
	}

	result, _ := decision.Evaluate(context.Background(), tree, input, types.None[decision.IIntelligence]())
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("top-level context field failed: got %v, want Permit", result.Outcome)
	}
}

func TestIntelligenceUnavailableError(t *testing.T) {
	err := &decision.IntelligenceUnavailableError{}
	if err.Error() == "" {
		t.Error("error message should not be empty")
	}
}

func TestSemanticNoPrompt(t *testing.T) {
	// Semantic condition with no prompt set
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:     types.MustFieldPath("context.tone"),
			Operator:  event.ConditionOperatorSemantic,
			Threshold: types.Some(types.MustScore(0.7)),
			// No Prompt set
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{Boolean: types.Some(true)},
				Child: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.9)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.5)),
	})

	intel := &mockIntelligence{
		content:    "yes",
		confidence: types.MustScore(0.85),
		tokens:     10,
	}

	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.Some[decision.IIntelligence](intel))
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Outcome != event.DecisionOutcomeDeny {
		t.Errorf("Outcome = %v, want Deny", result.Outcome)
	}
}

func TestSemanticNoThreshold(t *testing.T) {
	// Semantic condition with no threshold — should take default
	tree := decision.NewDecisionTree(&decision.InternalNode{
		Condition: event.Condition{
			Field:    types.MustFieldPath("context.tone"),
			Operator: event.ConditionOperatorSemantic,
			Prompt:   types.Some("Is this hostile?"),
		},
		Branches: []decision.Branch{
			{
				Match: event.MatchValue{Boolean: types.Some(true)},
				Child: decision.NewLeaf(event.DecisionOutcomeDeny, types.MustScore(0.9)),
			},
		},
		Default: decision.NewLeaf(event.DecisionOutcomePermit, types.MustScore(0.5)),
	})

	intel := &mockIntelligence{
		content:    "yes",
		confidence: types.MustScore(0.85),
		tokens:     10,
	}

	result, err := decision.Evaluate(context.Background(), tree, testInput("test"), types.Some[decision.IIntelligence](intel))
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	// No threshold set, so threshold check is skipped, goes to default
	if result.Outcome != event.DecisionOutcomePermit {
		t.Errorf("Outcome = %v, want Permit (no threshold)", result.Outcome)
	}
}
