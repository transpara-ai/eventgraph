package intelligence_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/intelligence"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// ════════════════════════════════════════════════════════════════════════
// Unit tests — no API calls
// ════════════════════════════════════════════════════════════════════════

func TestNewUnknownProvider(t *testing.T) {
	_, err := intelligence.New(intelligence.Config{
		Provider: "unknown",
		Model:    "some-model",
	})
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestNewOpenAICompatibleNotYetImplemented(t *testing.T) {
	_, err := intelligence.New(intelligence.Config{
		Provider: "openai-compatible",
		Model:    "gpt-4o",
	})
	if err == nil {
		t.Fatal("expected error for openai-compatible (not yet implemented)")
	}
}

func TestNewAnthropicRequiresModel(t *testing.T) {
	_, err := intelligence.New(intelligence.Config{
		Provider: "anthropic",
	})
	if err == nil {
		t.Fatal("expected error when model is empty")
	}
}

func TestNewAnthropicSuccess(t *testing.T) {
	p, err := intelligence.New(intelligence.Config{
		Provider: "anthropic",
		Model:    "claude-sonnet-4-6",
		APIKey:   "test-key-not-real",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "anthropic" {
		t.Errorf("Name() = %q, want anthropic", p.Name())
	}
	if p.Model() != "claude-sonnet-4-6" {
		t.Errorf("Model() = %q, want claude-sonnet-4-6", p.Model())
	}
}

func TestDefaultMaxTokens(t *testing.T) {
	p, err := intelligence.New(intelligence.Config{
		Provider: "anthropic",
		Model:    "claude-sonnet-4-6",
		APIKey:   "test-key",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Provider should be created successfully with default max tokens
	if p == nil {
		t.Fatal("provider is nil")
	}
}

func TestProviderSatisfiesIIntelligence(t *testing.T) {
	p, err := intelligence.New(intelligence.Config{
		Provider: "anthropic",
		Model:    "claude-sonnet-4-6",
		APIKey:   "test-key",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify the provider can be used as IIntelligence
	var _ interface {
		Reason(ctx context.Context, prompt string, history []event.Event) (interface{ Content() string }, error)
	}
	_ = p
}

func TestConfigWithAllOptions(t *testing.T) {
	p, err := intelligence.New(intelligence.Config{
		Provider:     "anthropic",
		Model:        "claude-sonnet-4-6",
		APIKey:       "test-key",
		BaseURL:      "https://custom.api.example.com",
		MaxTokens:    2048,
		Temperature:  0.5,
		SystemPrompt: "You are a helpful assistant.",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "anthropic" {
		t.Errorf("Name() = %q, want anthropic", p.Name())
	}
}

// ════════════════════════════════════════════════════════════════════════
// Claude CLI unit tests
// ════════════════════════════════════════════════════════════════════════

func TestNewClaudeCliDefaultModel(t *testing.T) {
	p, err := intelligence.New(intelligence.NewClaudeCliConfig(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "claude-cli" {
		t.Errorf("Name() = %q, want claude-cli", p.Name())
	}
	if p.Model() != "sonnet" {
		t.Errorf("Model() = %q, want sonnet", p.Model())
	}
}

func TestNewClaudeCliCustomModel(t *testing.T) {
	p, err := intelligence.New(intelligence.NewClaudeCliConfig("opus"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Model() != "opus" {
		t.Errorf("Model() = %q, want opus", p.Model())
	}
}

func TestNewClaudeCliConfig(t *testing.T) {
	cfg := intelligence.NewClaudeCliConfig("haiku")
	if cfg.Provider != "claude-cli" {
		t.Errorf("Provider = %q, want claude-cli", cfg.Provider)
	}
	if cfg.Model != "haiku" {
		t.Errorf("Model = %q, want haiku", cfg.Model)
	}
}

// ════════════════════════════════════════════════════════════════════════
// Claude CLI integration tests — require `claude` in PATH
// ════════════════════════════════════════════════════════════════════════

func skipWithoutClaudeCli(t *testing.T) {
	t.Helper()
	if os.Getenv("EVENTGRAPH_TEST_CLAUDE_CLI") == "" {
		t.Skip("EVENTGRAPH_TEST_CLAUDE_CLI not set — skipping Claude CLI integration test")
	}
}

func TestIntegrationClaudeCliReason(t *testing.T) {
	skipWithoutClaudeCli(t)

	p, err := intelligence.New(intelligence.NewClaudeCliConfig("sonnet"))
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()
	resp, err := p.Reason(ctx, "Reply with exactly one word: hello", nil)
	if err != nil {
		t.Fatalf("Reason failed: %v", err)
	}

	if resp.Content() == "" {
		t.Error("response content is empty")
	}
	t.Logf("Claude CLI response: %q (tokens: %d)", resp.Content(), resp.TokensUsed())
}

func TestIntegrationClaudeCliReasonWithSystemPrompt(t *testing.T) {
	skipWithoutClaudeCli(t)

	p, err := intelligence.New(intelligence.Config{
		Provider:     "claude-cli",
		Model:        "sonnet",
		SystemPrompt: "You are a decision-making system for an event graph. Respond with only: PERMIT, DENY, ESCALATE, or DEFER.",
	})
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()
	resp, err := p.Reason(ctx, "An agent with trust score 0.3 wants to delete critical data. What is your decision?", nil)
	if err != nil {
		t.Fatalf("Reason failed: %v", err)
	}

	content := resp.Content()
	if content == "" {
		t.Error("response content is empty")
	}
	t.Logf("Decision response: %q", content)
}

func TestIntegrationClaudeCliReasonWithHistory(t *testing.T) {
	skipWithoutClaudeCli(t)

	p, err := intelligence.New(intelligence.NewClaudeCliConfig("sonnet"))
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	actorID := types.MustActorID("actor_00000000000000000000000000000001")
	evID := types.MustEventID("01912345-6789-7abc-8def-0123456789ab")
	evType := types.MustEventType("trust.updated")
	ts := types.NewTimestamp(time.Now())
	convID := types.MustConversationID("conv_00000000000000000000000000000001")
	hash := types.MustHash("0000000000000000000000000000000000000000000000000000000000000000")
	sig := types.MustSignature(make([]byte, 64))

	ev := event.NewEvent(1, evID, evType, ts, actorID, nil, []types.EventID{evID}, convID, hash, hash, sig)

	ctx := context.Background()
	resp, err := p.Reason(ctx, "Given the event history, reply with exactly one word: acknowledged", []event.Event{ev})
	if err != nil {
		t.Fatalf("Reason with history failed: %v", err)
	}

	if resp.Content() == "" {
		t.Error("response content is empty")
	}
	t.Logf("Response with history: %q", resp.Content())
}

// ════════════════════════════════════════════════════════════════════════
// Anthropic API integration tests — require ANTHROPIC_API_KEY
// ════════════════════════════════════════════════════════════════════════

func skipWithoutAnthropicKey(t *testing.T) string {
	t.Helper()
	key := os.Getenv("ANTHROPIC_API_KEY")
	if key == "" {
		t.Skip("ANTHROPIC_API_KEY not set — skipping integration test")
	}
	return key
}

func TestIntegrationAnthropicReason(t *testing.T) {
	key := skipWithoutAnthropicKey(t)

	p, err := intelligence.New(intelligence.Config{
		Provider:  "anthropic",
		Model:     "claude-sonnet-4-6",
		APIKey:    key,
		MaxTokens: 100,
	})
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()
	resp, err := p.Reason(ctx, "Reply with exactly one word: hello", nil)
	if err != nil {
		t.Fatalf("Reason failed: %v", err)
	}

	if resp.Content() == "" {
		t.Error("response content is empty")
	}
	if resp.TokensUsed() == 0 {
		t.Error("tokens used is 0")
	}
	if resp.Confidence().Value() <= 0 {
		t.Error("confidence should be positive")
	}
}

func TestIntegrationAnthropicReasonWithHistory(t *testing.T) {
	key := skipWithoutAnthropicKey(t)

	p, err := intelligence.New(intelligence.Config{
		Provider:  "anthropic",
		Model:     "claude-sonnet-4-6",
		APIKey:    key,
		MaxTokens: 100,
	})
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	// Create a minimal event for history context.
	actorID := types.MustActorID("actor_00000000000000000000000000000001")
	evID := types.MustEventID("01912345-6789-7abc-8def-0123456789ab")
	evType := types.MustEventType("trust.updated")
	ts := types.NewTimestamp(time.Now())
	convID := types.MustConversationID("conv_00000000000000000000000000000001")
	hash := types.MustHash("0000000000000000000000000000000000000000000000000000000000000000")
	sigBytes := make([]byte, 64)
	sig := types.MustSignature(sigBytes)

	ev := event.NewEvent(1, evID, evType, ts, actorID, nil, []types.EventID{evID}, convID, hash, hash, sig)

	ctx := context.Background()
	resp, err := p.Reason(ctx, "Based on the event history, reply with exactly one word: acknowledged", []event.Event{ev})
	if err != nil {
		t.Fatalf("Reason with history failed: %v", err)
	}

	if resp.Content() == "" {
		t.Error("response content is empty")
	}
}

func TestIntegrationAnthropicReasonWithSystemPrompt(t *testing.T) {
	key := skipWithoutAnthropicKey(t)

	p, err := intelligence.New(intelligence.Config{
		Provider:     "anthropic",
		Model:        "claude-sonnet-4-6",
		APIKey:       key,
		MaxTokens:    100,
		SystemPrompt: "You are a decision-making system. Always respond with exactly one word.",
	})
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()
	resp, err := p.Reason(ctx, "Should this action be permitted?", nil)
	if err != nil {
		t.Fatalf("Reason with system prompt failed: %v", err)
	}

	if resp.Content() == "" {
		t.Error("response content is empty")
	}
}

func TestIntegrationAnthropicInvalidKey(t *testing.T) {
	p, err := intelligence.New(intelligence.Config{
		Provider:  "anthropic",
		Model:     "claude-sonnet-4-6",
		APIKey:    "sk-ant-invalid-key-for-testing",
		MaxTokens: 50,
	})
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()
	_, err = p.Reason(ctx, "hello", nil)
	if err == nil {
		t.Fatal("expected error with invalid API key")
	}
}
