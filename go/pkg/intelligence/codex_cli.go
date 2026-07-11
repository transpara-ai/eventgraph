package intelligence

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/transpara-ai/eventgraph/go/pkg/decision"
	"github.com/transpara-ai/eventgraph/go/pkg/event"
)

const (
	// defaultCodexReasonTimeout is the maximum time a single Reason() call
	// can run against the Codex CLI.
	defaultCodexReasonTimeout = 10 * time.Minute

	// defaultCodexOperateTimeout is the maximum time a single Operate() call
	// can run against the Codex CLI.
	defaultCodexOperateTimeout = 15 * time.Minute
)

// codexEvent is a single line from `codex exec --json` JSONL output.
type codexEvent struct {
	Type    string          `json:"type"`
	Message string          `json:"message,omitempty"` // error events
	Item    *codexItem      `json:"item,omitempty"`    // item.completed events
	Usage   *codexUsage     `json:"usage,omitempty"`   // turn.completed events
	Error   *codexTurnError `json:"error,omitempty"`   // turn.failed events
}

type codexItem struct {
	ID   string `json:"id"`
	Type string `json:"type"` // "agent_message", "command_execution"
	Text string `json:"text,omitempty"`
}

type codexUsage struct {
	InputTokens       int `json:"input_tokens"`
	CachedInputTokens int `json:"cached_input_tokens"`
	OutputTokens      int `json:"output_tokens"`
	ReasoningTokens   int `json:"reasoning_output_tokens"`
}

type codexTurnError struct {
	Message string `json:"message"`
}

// codexCliProvider implements Provider by shelling out to the `codex` CLI.
// This uses whatever authentication Codex already has (ChatGPT Plus, API key,
// etc.) without requiring separate credentials.
type codexCliProvider struct {
	model        string
	systemPrompt string
	codexPath    string // path to codex binary, default "codex"
}

func newCodexCliProvider(cfg Config) (*codexCliProvider, error) {
	model := cfg.Model
	if model == "" {
		model = "o3" // Codex default
	}

	codexPath := "codex"
	if cfg.BaseURL != "" {
		// BaseURL repurposed as path to codex binary for testing.
		codexPath = cfg.BaseURL
	}

	if _, err := exec.LookPath(codexPath); err != nil {
		return nil, fmt.Errorf("codex CLI not found in PATH: %w", err)
	}

	return &codexCliProvider{
		model:        model,
		systemPrompt: cfg.SystemPrompt,
		codexPath:    codexPath,
	}, nil
}

func (p *codexCliProvider) Name() string  { return "codex-cli" }
func (p *codexCliProvider) Model() string { return p.model }

func (p *codexCliProvider) Reason(ctx context.Context, prompt string, history []event.Event) (decision.Response, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, defaultCodexReasonTimeout)
	defer cancel()

	// Reason is a decision call, never an execution authority. Codex CLI is an
	// agentic client even for `exec`: without explicit constraints it can honor
	// ambient user config, inspect the caller's repository, run tools, and write
	// or commit files while the caller records only agent.evaluated. Run it from
	// an empty non-repository root under the CLI's read-only sandbox so every
	// filesystem mutation remains exclusive to the governed Operate path.
	reasonRoot, err := os.MkdirTemp("", "eventgraph-codex-reason-")
	if err != nil {
		return decision.Response{}, fmt.Errorf("codex reason: create isolated root: %w", err)
	}
	defer func() { _ = os.RemoveAll(reasonRoot) }()

	var fullPrompt strings.Builder
	historyText := eventsToMessages(history)
	if historyText != "" {
		fullPrompt.WriteString(historyText)
		fullPrompt.WriteString("\n---\n\n")
	}
	fullPrompt.WriteString(prompt)

	args := []string{
		"exec",
		"--json",
		"--ephemeral",
		"--ignore-user-config",
		"--sandbox", "read-only",
		"--skip-git-repo-check",
		"-C", reasonRoot,
		"-m", p.model,
	}
	if p.systemPrompt != "" {
		args = append(args, "-c", fmt.Sprintf("system_prompt=%q", p.systemPrompt))
	}
	args = append(args, "-") // read prompt from stdin

	cmd := exec.CommandContext(ctx, p.codexPath, args...)
	cmd.Stdin = strings.NewReader(fullPrompt.String())

	// Reason uses the same credential isolation as Operate even though its
	// filesystem sandbox is stricter. Read-only files do not prevent a tool from
	// using inherited gh/SSH credentials for remote side effects.
	env, cleanupEnv, err := codexReasonSubprocessEnv(cmd.Environ(), reasonRoot)
	if cleanupEnv != nil {
		defer cleanupEnv()
	}
	if err != nil {
		return decision.Response{}, fmt.Errorf("codex reason: isolate environment: %w", err)
	}
	cmd.Env = env

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := runWithProgress(cmd, "  ⏳ codex thinking"); err != nil {
		// Try to extract useful info from JSONL even on error.
		if stdout.Len() > 0 {
			text, _, parseErr := parseCodexJSONL(stdout.Bytes())
			if parseErr == nil && text != "" {
				return decision.NewResponse(text, defaultConfidence(), decision.TokenUsage{}), nil
			}
			// Include the JSONL parse error (often contains the real reason).
			if parseErr != nil {
				return decision.Response{}, fmt.Errorf("codex CLI error: %w (%v)", err, parseErr)
			}
		}
		return decision.Response{}, fmt.Errorf("codex CLI error: %w\nstderr: %s", err, stderr.String())
	}

	text, usage, err := parseCodexJSONL(stdout.Bytes())
	if err != nil {
		return decision.Response{}, err
	}

	tokenUsage := decision.TokenUsage{
		InputTokens:     usage.InputTokens,
		OutputTokens:    usage.OutputTokens,
		CacheReadTokens: usage.CachedInputTokens,
	}

	return decision.NewResponse(text, defaultConfidence(), tokenUsage), nil
}

// codexReasonSubprocessEnv reduces Reason's environment to transport/runtime
// essentials plus explicit Git/GitHub neutralizers. Unlike Operate, Reason has
// no authority to run repository or remote mutations, so arbitrary ambient
// variables are not inherited: unknown present or future credential names fail
// closed by omission rather than relying on an ever-complete denylist.
func codexReasonSubprocessEnv(parent []string, reasonRoot string) (env []string, cleanup func(), err error) {
	isolated, cleanup, err := operateSubprocessEnv(parent)
	if err != nil {
		return nil, cleanup, err
	}

	parentValues := make(map[string]string, len(parent))
	for _, entry := range parent {
		if key, value, ok := strings.Cut(entry, "="); ok {
			parentValues[key] = value
		}
	}
	codexHome := parentValues["CODEX_HOME"]
	if codexHome == "" && parentValues["HOME"] != "" {
		codexHome = filepath.Join(parentValues["HOME"], ".codex")
	}

	allowed := map[string]bool{
		"PATH": true,
		"LANG": true, "LANGUAGE": true, "TZ": true, "TERM": true, "NO_COLOR": true,
		"SSL_CERT_FILE": true, "SSL_CERT_DIR": true, "CURL_CA_BUNDLE": true, "REQUESTS_CA_BUNDLE": true,
		"GH_CONFIG_DIR": true, "GIT_CONFIG_GLOBAL": true, "GIT_CONFIG_NOSYSTEM": true,
		"GIT_TERMINAL_PROMPT": true, "GIT_SSH_COMMAND": true,
	}
	for _, entry := range isolated {
		key, _, ok := strings.Cut(entry, "=")
		if ok && (allowed[key] || strings.HasPrefix(key, "LC_")) {
			env = append(env, entry)
		}
	}

	// A temporary HOME hides default credential stores (~/.ssh, ~/.aws,
	// ~/.config/gh). CODEX_HOME names the subscription-auth home explicitly;
	// --ignore-user-config prevents its config/MCP settings from becoming tools.
	env = append(env, "HOME="+reasonRoot, "TMPDIR="+reasonRoot)
	if codexHome != "" {
		env = append(env, "CODEX_HOME="+codexHome)
	}
	return env, cleanup, nil
}

func (p *codexCliProvider) Operate(ctx context.Context, task decision.OperateTask) (decision.OperateResult, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, defaultCodexOperateTimeout)
	defer cancel()

	if task.WorkDir == "" {
		return decision.OperateResult{}, fmt.Errorf("Operate requires WorkDir")
	}
	if task.Instruction == "" {
		return decision.OperateResult{}, fmt.Errorf("Operate requires Instruction")
	}

	args := []string{
		"exec",
		"--json",
		"--ephemeral",
		"-m", p.model,
		"-C", task.WorkDir,
		"--dangerously-bypass-approvals-and-sandbox",
	}
	if p.systemPrompt != "" {
		args = append(args, "-c", fmt.Sprintf("system_prompt=%q", p.systemPrompt))
	}
	args = append(args, "-") // read from stdin

	cmd := exec.CommandContext(ctx, p.codexPath, args...)
	cmd.Stdin = strings.NewReader(task.Instruction)

	// Credential-isolated environment (slice-1 v10-F1): codex Operate runs with
	// --dangerously-bypass-approvals-and-sandbox, so it must not inherit the
	// daemon's ambient git/gh credentials. Fail closed.
	env, cleanupEnv, err := operateSubprocessEnv(cmd.Environ())
	if err != nil {
		return decision.OperateResult{}, err
	}
	defer cleanupEnv()
	cmd.Env = env

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := runWithProgress(cmd, "  ⏳ codex working"); err != nil {
		if stdout.Len() > 0 {
			text, _, parseErr := parseCodexJSONL(stdout.Bytes())
			if parseErr == nil && text != "" {
				return decision.OperateResult{Summary: text}, nil
			}
			if parseErr != nil {
				return decision.OperateResult{}, fmt.Errorf("codex CLI operate error: %w (%v)", err, parseErr)
			}
		}
		return decision.OperateResult{}, fmt.Errorf("codex CLI operate error: %w\nstderr: %s", err, stderr.String())
	}

	text, usage, err := parseCodexJSONL(stdout.Bytes())
	if err != nil {
		return decision.OperateResult{}, err
	}

	return decision.OperateResult{
		Summary: text,
		Usage: decision.TokenUsage{
			InputTokens:     usage.InputTokens,
			OutputTokens:    usage.OutputTokens,
			CacheReadTokens: usage.CachedInputTokens,
		},
	}, nil
}

// parseCodexJSONL extracts the final agent message text and usage from Codex
// JSONL output. It collects all agent_message items and returns the last
// turn.completed usage block.
func parseCodexJSONL(data []byte) (string, codexUsage, error) {
	var messages []string
	var usage codexUsage
	var turnError string

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev codexEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			continue // skip unparseable lines
		}

		switch ev.Type {
		case "item.completed":
			if ev.Item != nil && ev.Item.Type == "agent_message" && ev.Item.Text != "" {
				messages = append(messages, ev.Item.Text)
			}
		case "turn.completed":
			if ev.Usage != nil {
				usage = *ev.Usage
			}
		case "turn.failed":
			if ev.Error != nil {
				turnError = ev.Error.Message
			}
		case "error":
			if ev.Message != "" {
				turnError = ev.Message
			}
		}
	}

	if turnError != "" && len(messages) == 0 {
		return "", usage, fmt.Errorf("codex turn failed: %s", turnError)
	}

	if len(messages) == 0 {
		return "", usage, fmt.Errorf("codex returned no agent messages")
	}

	// Return the last agent message (most recent response).
	return messages[len(messages)-1], usage, nil
}

// Compile-time checks.
var (
	_ Provider           = (*codexCliProvider)(nil)
	_ decision.IOperator = (*codexCliProvider)(nil)
)
