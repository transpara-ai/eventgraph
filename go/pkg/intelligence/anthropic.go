package intelligence

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/transpara-ai/eventgraph/go/pkg/decision"
	"github.com/transpara-ai/eventgraph/go/pkg/event"
)

// anthropicProvider implements Provider using the official Anthropic SDK.
type anthropicProvider struct {
	client       anthropic.Client
	model        anthropic.Model
	modelStr     string
	maxTokens    int
	temperature  float64
	systemPrompt string
}

func newAnthropicProvider(cfg Config) (*anthropicProvider, error) {
	if cfg.Model == "" {
		return nil, fmt.Errorf("anthropic provider requires a model (e.g., claude-sonnet-4-6)")
	}

	// Resolve credentials: explicit config > env var.
	// OAuth tokens (sk-ant-oat*) use Authorization: Bearer.
	// API keys (sk-ant-api*) use x-api-key header.
	// The SDK auto-reads ANTHROPIC_API_KEY and ANTHROPIC_AUTH_TOKEN from env,
	// but we need to route OAuth tokens to the correct auth method.
	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}

	var opts []option.RequestOption
	if apiKey != "" {
		if strings.HasPrefix(apiKey, "sk-ant-oat") {
			// OAuth tokens use Bearer auth. We must also clear the x-api-key
			// header that the SDK may set from ANTHROPIC_API_KEY env var.
			opts = append(opts,
				option.WithAuthToken(apiKey),
				option.WithAPIKey(""),
			)
		} else {
			opts = append(opts, option.WithAPIKey(apiKey))
		}
	}
	if cfg.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(cfg.BaseURL))
	}

	client := anthropic.NewClient(opts...)

	return &anthropicProvider{
		client:       client,
		model:        anthropic.Model(cfg.Model),
		modelStr:     cfg.Model,
		maxTokens:    cfg.MaxTokens,
		temperature:  cfg.Temperature,
		systemPrompt: cfg.SystemPrompt,
	}, nil
}

func (p *anthropicProvider) Name() string  { return "anthropic" }
func (p *anthropicProvider) Model() string { return p.modelStr }

func (p *anthropicProvider) Reason(ctx context.Context, prompt string, history []event.Event) (decision.Response, error) {
	var messages []anthropic.MessageParam

	// Include event history as context if present.
	historyText := eventsToMessages(history)
	if historyText != "" {
		messages = append(messages,
			anthropic.NewUserMessage(anthropic.NewTextBlock(historyText)),
			anthropic.NewAssistantMessage(anthropic.NewTextBlock("I understand the event history. What would you like me to reason about?")),
		)
	}

	messages = append(messages,
		anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
	)

	params := anthropic.MessageNewParams{
		Model:     p.model,
		MaxTokens: int64(p.maxTokens),
		Messages:  messages,
	}

	if p.systemPrompt != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: p.systemPrompt},
		}
	}

	if p.temperature > 0 {
		params.Temperature = anthropic.Float(p.temperature)
	}

	msg, err := p.client.Messages.New(ctx, params)
	if err != nil {
		return decision.Response{}, fmt.Errorf("anthropic API error: %w", err)
	}

	// Extract text content from response.
	var content strings.Builder
	for _, block := range msg.Content {
		if tb, ok := block.AsAny().(anthropic.TextBlock); ok {
			content.WriteString(tb.Text)
		}
	}

	confidence := defaultConfidence()
	usage := decision.TokenUsage{
		InputTokens:  int(msg.Usage.InputTokens),
		OutputTokens: int(msg.Usage.OutputTokens),
	}

	return decision.NewResponse(content.String(), confidence, usage), nil
}

// Compile-time check.
var _ Provider = (*anthropicProvider)(nil)
