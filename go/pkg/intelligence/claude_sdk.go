package intelligence

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lovyou-ai/eventgraph/go/pkg/decision"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
)

// errStaleSession is returned when the CLI reports that a resumed session
// no longer exists. Callers retry without --resume.
var errStaleSession = errors.New("stale session")

// claudeSDKProvider implements Provider using Claude CLI with session management.
// Uses --output-format json (single result) for reliable output parsing.
// Uses --resume for warm session context across calls.
// Captures session_id from results for persistence.
type claudeSDKProvider struct {
	model         string
	maxBudget     float64
	systemPrompt  string
	cliPath       string
	mcpConfigPath string
	sessionID     string
	mu            sync.Mutex
}

func newClaudeSDKProvider(cfg Config) (*claudeSDKProvider, error) {
	model := cfg.Model
	if model == "" {
		model = "sonnet"
	}

	cliPath := "claude"
	if cfg.BaseURL != "" {
		cliPath = cfg.BaseURL
	}
	if _, err := exec.LookPath(cliPath); err != nil {
		return nil, fmt.Errorf("claude CLI not found in PATH: %w", err)
	}

	maxBudget := 1.0
	if cfg.MaxBudgetUSD > 0 {
		maxBudget = cfg.MaxBudgetUSD
	}

	mcpConfig := cfg.MCPConfigPath
	if mcpConfig != "" {
		if !filepath.IsAbs(mcpConfig) {
			return nil, fmt.Errorf("MCPConfigPath must be absolute: %s", mcpConfig)
		}
		info, err := os.Stat(mcpConfig)
		if err != nil {
			return nil, fmt.Errorf("MCPConfigPath: %w", err)
		}
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("MCPConfigPath is not a regular file: %s", mcpConfig)
		}
	}

	return &claudeSDKProvider{
		model:         model,
		maxBudget:     maxBudget,
		systemPrompt:  cfg.SystemPrompt,
		cliPath:       cliPath,
		mcpConfigPath: mcpConfig,
		sessionID:     cfg.SessionID,
	}, nil
}

func (p *claudeSDKProvider) Name() string  { return "claude-sdk" }
func (p *claudeSDKProvider) Model() string { return p.model }

func (p *claudeSDKProvider) SessionID() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.sessionID
}

func (p *claudeSDKProvider) Reason(ctx context.Context, prompt string, history []event.Event) (decision.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultReasonTimeout)
	defer cancel()

	var fullPrompt strings.Builder
	historyText := eventsToMessages(history)
	if historyText != "" {
		fullPrompt.WriteString(historyText)
		fullPrompt.WriteString("\n---\n\n")
	}
	fullPrompt.WriteString(prompt)

	args := p.baseArgs()
	result, err := p.run(ctx, fullPrompt.String(), args)
	if err != nil && p.shouldRetryCold(err) {
		log.Printf("[claude-sdk] stale session for reason, retrying cold")
		p.clearSession()
		result, err = p.run(ctx, fullPrompt.String(), p.baseArgs())
	}
	if err != nil {
		return decision.Response{}, fmt.Errorf("claude-sdk reason: %w", err)
	}

	usage := decision.TokenUsage{
		InputTokens:      result.Usage.InputTokens,
		OutputTokens:     result.Usage.OutputTokens,
		CacheReadTokens:  result.Usage.CacheReadInputTokens,
		CacheWriteTokens: result.Usage.CacheCreationInputTokens,
		CostUSD:          result.TotalCostUSD,
	}
	return decision.NewResponse(result.Result, defaultConfidence(), usage), nil
}

func (p *claudeSDKProvider) Operate(ctx context.Context, task decision.OperateTask) (decision.OperateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultOperateTimeout)
	defer cancel()

	if task.WorkDir == "" {
		return decision.OperateResult{}, fmt.Errorf("Operate requires WorkDir")
	}
	if task.Instruction == "" {
		return decision.OperateResult{}, fmt.Errorf("Operate requires Instruction")
	}

	args := p.baseArgs()
	args = append(args, "--permission-mode", "bypassPermissions")
	if len(task.AllowedTools) > 0 {
		args = append(args, "--allowed-tools", strings.Join(task.AllowedTools, ","))
	} else {
		args = append(args, "--allowed-tools", "Read,Edit,Write,Bash,Glob,Grep")
	}

	buildArgs := func() []string {
		a := p.baseArgs()
		a = append(a, "--permission-mode", "bypassPermissions")
		if len(task.AllowedTools) > 0 {
			a = append(a, "--allowed-tools", strings.Join(task.AllowedTools, ","))
		} else {
			a = append(a, "--allowed-tools", "Read,Edit,Write,Bash,Glob,Grep")
		}
		return a
	}

	result, err := p.runInDir(ctx, task.Instruction, buildArgs(), task.WorkDir)
	if err != nil && p.shouldRetryCold(err) {
		log.Printf("[claude-sdk] stale session for operate, retrying cold")
		p.clearSession()
		result, err = p.runInDir(ctx, task.Instruction, buildArgs(), task.WorkDir)
	}
	if err != nil {
		return decision.OperateResult{}, fmt.Errorf("claude-sdk operate: %w", err)
	}

	usage := decision.TokenUsage{
		InputTokens:      result.Usage.InputTokens,
		OutputTokens:     result.Usage.OutputTokens,
		CacheReadTokens:  result.Usage.CacheReadInputTokens,
		CacheWriteTokens: result.Usage.CacheCreationInputTokens,
		CostUSD:          result.TotalCostUSD,
	}
	return decision.OperateResult{Summary: result.Result, Usage: usage}, nil
}

// cliResult is the JSON output from `claude -p --output-format json`.
type cliResult struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	IsError bool   `json:"is_error"`
	Result  string `json:"result"`
	Usage   struct {
		InputTokens              int `json:"input_tokens"`
		OutputTokens             int `json:"output_tokens"`
		CacheReadInputTokens     int `json:"cache_read_input_tokens"`
		CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	} `json:"usage"`
	TotalCostUSD float64 `json:"total_cost_usd"`
	SessionID    string  `json:"session_id"`
	NumTurns     int     `json:"num_turns"`
	DurationMs   int     `json:"duration_ms"`
}

func (p *claudeSDKProvider) baseArgs() []string {
	args := []string{
		"-p",
		"--output-format", "json",
		"--model", p.model,
		"--max-budget-usd", fmt.Sprintf("%.2f", p.maxBudget),
	}
	p.mu.Lock()
	sid := p.sessionID
	p.mu.Unlock()
	if sid != "" {
		args = append(args, "--resume", sid)
	}
	if p.systemPrompt != "" {
		args = append(args, "--system-prompt", p.systemPrompt)
	}
	if p.mcpConfigPath != "" {
		args = append(args, "--mcp-config", p.mcpConfigPath)
	}
	return args
}

func (p *claudeSDKProvider) run(ctx context.Context, prompt string, args []string) (*cliResult, error) {
	return p.runInDir(ctx, prompt, args, "")
}

func (p *claudeSDKProvider) runInDir(ctx context.Context, prompt string, args []string, workDir string) (*cliResult, error) {
	cmd := exec.CommandContext(ctx, p.cliPath, args...)
	cmd.Stdin = strings.NewReader(prompt)
	if workDir != "" {
		cmd.Dir = workDir
	}

	// Scrub sensitive env vars.
	env := scrubEnv(cmd.Environ(), "CLAUDECODE", "DATABASE_URL", "HIVE_AGENT_ID", "HIVE_HUMAN_ID")
	cmd.Env = env

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := runWithProgress(cmd, "  ⏳ thinking"); err != nil {
		// Check if we got JSON despite non-zero exit.
		if stdout.Len() > 0 {
			var result cliResult
			if jsonErr := json.Unmarshal(stdout.Bytes(), &result); jsonErr == nil {
				if result.IsError && result.NumTurns == 0 && result.DurationMs == 0 {
					return nil, errStaleSession
				}
				if result.Result != "" {
					p.captureSession(result.SessionID)
					return &result, nil
				}
			}
		}
		return nil, fmt.Errorf("claude CLI error: %w\nstderr: %s", err, stderr.String())
	}

	var result cliResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse claude CLI JSON output: %w\nraw: %s", err, stdout.String())
	}

	if result.IsError {
		if result.NumTurns == 0 && result.DurationMs == 0 {
			return nil, errStaleSession
		}
		return nil, fmt.Errorf("claude error (subtype=%s, turns=%d): %s", result.Subtype, result.NumTurns, result.Result)
	}

	p.captureSession(result.SessionID)
	return &result, nil
}

func (p *claudeSDKProvider) captureSession(sid string) {
	if sid == "" {
		return
	}
	p.mu.Lock()
	p.sessionID = sid
	p.mu.Unlock()
}

func (p *claudeSDKProvider) shouldRetryCold(err error) bool {
	p.mu.Lock()
	has := p.sessionID != ""
	p.mu.Unlock()
	return has && errors.Is(err, errStaleSession)
}

func (p *claudeSDKProvider) clearSession() {
	p.mu.Lock()
	p.sessionID = ""
	p.mu.Unlock()
}

func scrubEnv(env []string, keys ...string) []string {
	out := make([]string, 0, len(env))
	for _, e := range env {
		skip := false
		for _, k := range keys {
			if strings.HasPrefix(e, k+"=") {
				skip = true
				break
			}
		}
		if !skip {
			out = append(out, e)
		}
	}
	return out
}

var (
	_ Provider           = (*claudeSDKProvider)(nil)
	_ decision.IOperator = (*claudeSDKProvider)(nil)
)
