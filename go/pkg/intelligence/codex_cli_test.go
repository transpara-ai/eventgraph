package intelligence

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCodexJSONL_HappyPath(t *testing.T) {
	jsonl := `{"type":"thread.started","thread_id":"abc123"}
{"type":"turn.started"}
{"type":"item.completed","item":{"id":"item_0","type":"agent_message","text":"Hello world"}}
{"type":"turn.completed","usage":{"input_tokens":100,"cached_input_tokens":50,"output_tokens":20,"reasoning_output_tokens":5}}`

	text, usage, err := parseCodexJSONL([]byte(jsonl))
	require.NoError(t, err)
	assert.Equal(t, "Hello world", text)
	assert.Equal(t, 100, usage.InputTokens)
	assert.Equal(t, 50, usage.CachedInputTokens)
	assert.Equal(t, 20, usage.OutputTokens)
	assert.Equal(t, 5, usage.ReasoningTokens)
}

func TestParseCodexJSONL_MultipleMessages(t *testing.T) {
	jsonl := `{"type":"item.completed","item":{"id":"item_0","type":"agent_message","text":"first"}}
{"type":"item.completed","item":{"id":"item_1","type":"command_execution","command":"ls"}}
{"type":"item.completed","item":{"id":"item_2","type":"agent_message","text":"second"}}
{"type":"turn.completed","usage":{"input_tokens":200,"output_tokens":40}}`

	text, _, err := parseCodexJSONL([]byte(jsonl))
	require.NoError(t, err)
	assert.Equal(t, "second", text, "should return last agent message")
}

func TestParseCodexJSONL_TurnFailed(t *testing.T) {
	jsonl := `{"type":"turn.started"}
{"type":"turn.failed","error":{"message":"model not supported"}}`

	_, _, err := parseCodexJSONL([]byte(jsonl))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "model not supported")
}

func TestParseCodexJSONL_ErrorEvent(t *testing.T) {
	jsonl := `{"type":"error","message":"invalid request"}`

	_, _, err := parseCodexJSONL([]byte(jsonl))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid request")
}

func TestParseCodexJSONL_NoMessages(t *testing.T) {
	jsonl := `{"type":"turn.started"}
{"type":"turn.completed","usage":{"input_tokens":10}}`

	_, _, err := parseCodexJSONL([]byte(jsonl))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no agent messages")
}

func TestParseCodexJSONL_SkipsMalformedLines(t *testing.T) {
	jsonl := `not json at all
{"type":"item.completed","item":{"id":"item_0","type":"agent_message","text":"works"}}
also not json
{"type":"turn.completed","usage":{"input_tokens":10}}`

	text, _, err := parseCodexJSONL([]byte(jsonl))
	require.NoError(t, err)
	assert.Equal(t, "works", text)
}

func TestParseCodexJSONL_TurnFailedWithMessage(t *testing.T) {
	// If there's a message AND an error, the message wins (partial success).
	jsonl := `{"type":"item.completed","item":{"id":"item_0","type":"agent_message","text":"partial result"}}
{"type":"turn.failed","error":{"message":"context cancelled"}}`

	text, _, err := parseCodexJSONL([]byte(jsonl))
	require.NoError(t, err)
	assert.Equal(t, "partial result", text)
}

func TestNewCodexCliProvider_Defaults(t *testing.T) {
	// This test only works if codex is in PATH — skip otherwise.
	p, err := newCodexCliProvider(Config{Provider: "codex-cli"})
	if err != nil {
		t.Skipf("codex not in PATH: %v", err)
	}
	assert.Equal(t, "codex-cli", p.Name())
	assert.Equal(t, "o3", p.Model(), "default model should be o3")
}

func TestNewCodexCliProvider_CustomModel(t *testing.T) {
	p, err := newCodexCliProvider(Config{Provider: "codex-cli", Model: "o4-mini"})
	if err != nil {
		t.Skipf("codex not in PATH: %v", err)
	}
	assert.Equal(t, "o4-mini", p.Model())
}

func TestCodexReasonUsesReadOnlyIsolatedInvocation(t *testing.T) {
	tmp := t.TempDir()
	argsFile := filepath.Join(tmp, "args")
	envFile := filepath.Join(tmp, "env")
	fakeCodex := filepath.Join(tmp, "codex")
	script := `#!/bin/sh
root=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
printf '%s\n' "$@" > "$root/args"
{
  printf 'GH_TOKEN=%s\n' "${GH_TOKEN-}"
  printf 'GITHUB_TOKEN=%s\n' "${GITHUB_TOKEN-}"
  printf 'GH_ENTERPRISE_TOKEN=%s\n' "${GH_ENTERPRISE_TOKEN-}"
  printf 'GHE_TOKEN=%s\n' "${GHE_TOKEN-}"
  printf 'SSH_AUTH_SOCK=%s\n' "${SSH_AUTH_SOCK-}"
  printf 'SSH_AGENT_PID=%s\n' "${SSH_AGENT_PID-}"
  printf 'CLAUDECODE=%s\n' "${CLAUDECODE-}"
  printf 'AWS_ACCESS_KEY_ID=%s\n' "${AWS_ACCESS_KEY_ID-}"
  printf 'AWS_SECRET_ACCESS_KEY=%s\n' "${AWS_SECRET_ACCESS_KEY-}"
  printf 'HTTPS_PROXY=%s\n' "${HTTPS_PROXY-}"
  printf 'GH_CONFIG_DIR=%s\n' "${GH_CONFIG_DIR-}"
  printf 'GIT_CONFIG_GLOBAL=%s\n' "${GIT_CONFIG_GLOBAL-}"
  printf 'GIT_CONFIG_NOSYSTEM=%s\n' "${GIT_CONFIG_NOSYSTEM-}"
  printf 'HOME=%s\n' "${HOME-}"
  printf 'TMPDIR=%s\n' "${TMPDIR-}"
  printf 'CODEX_HOME=%s\n' "${CODEX_HOME-}"
} > "$root/env"
printf '%s\n' '{"type":"item.completed","item":{"id":"item_0","type":"agent_message","text":"read-only decision"}}'
printf '%s\n' '{"type":"turn.completed","usage":{"input_tokens":12,"output_tokens":3}}'
`
	require.NoError(t, os.WriteFile(fakeCodex, []byte(script), 0o700))
	t.Setenv("GH_TOKEN", "must-not-reach-reason")
	t.Setenv("GITHUB_TOKEN", "must-not-reach-reason")
	t.Setenv("GH_ENTERPRISE_TOKEN", "must-not-reach-reason")
	t.Setenv("GHE_TOKEN", "must-not-reach-reason")
	t.Setenv("SSH_AUTH_SOCK", filepath.Join(tmp, "agent.sock"))
	t.Setenv("SSH_AGENT_PID", "12345")
	t.Setenv("CLAUDECODE", "1")
	t.Setenv("AWS_ACCESS_KEY_ID", "must-not-reach-reason")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "must-not-reach-reason")
	t.Setenv("HTTPS_PROXY", "https://credential:secret@proxy.invalid")
	ambientTmp := filepath.Join(tmp, "ambient-tmp")
	require.NoError(t, os.Mkdir(ambientTmp, 0o700))
	t.Setenv("TMPDIR", ambientTmp)
	t.Setenv("GH_CONFIG_DIR", filepath.Join(tmp, "ambient-gh"))
	t.Setenv("CODEX_HOME", filepath.Join(tmp, "codex-auth-home"))

	p, err := newCodexCliProvider(Config{Provider: "codex-cli", Model: "gpt-test", BaseURL: fakeCodex})
	require.NoError(t, err)
	resp, err := p.Reason(context.Background(), "classify this event", nil)
	require.NoError(t, err)
	assert.Equal(t, "read-only decision", resp.Content())

	rawArgs, err := os.ReadFile(argsFile)
	require.NoError(t, err)
	args := strings.Fields(string(rawArgs))
	assert.Contains(t, args, "--ignore-user-config")
	assert.Contains(t, args, "--skip-git-repo-check")
	assert.NotContains(t, args, "--dangerously-bypass-approvals-and-sandbox")

	sandboxAt := indexOf(args, "--sandbox")
	require.GreaterOrEqual(t, sandboxAt, 0)
	require.Less(t, sandboxAt+1, len(args))
	assert.Equal(t, "read-only", args[sandboxAt+1])

	rootAt := indexOf(args, "-C")
	require.GreaterOrEqual(t, rootAt, 0)
	require.Less(t, rootAt+1, len(args))
	reasonRoot := args[rootAt+1]
	assert.Contains(t, filepath.Base(reasonRoot), "eventgraph-codex-reason-")
	_, err = os.Stat(reasonRoot)
	assert.ErrorIs(t, err, os.ErrNotExist, "isolated Reason root must be removed after the call")

	rawEnv, err := os.ReadFile(envFile)
	require.NoError(t, err)
	env := string(rawEnv)
	for _, key := range []string{
		"GH_TOKEN", "GITHUB_TOKEN", "GH_ENTERPRISE_TOKEN", "GHE_TOKEN",
		"SSH_AUTH_SOCK", "SSH_AGENT_PID", "CLAUDECODE",
		"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "HTTPS_PROXY",
	} {
		assert.Contains(t, env, key+"=\n", "%s must be omitted from Reason", key)
	}
	assert.NotContains(t, env, "GH_CONFIG_DIR="+filepath.Join(tmp, "ambient-gh"))
	assert.Contains(t, env, "GIT_CONFIG_GLOBAL="+os.DevNull+"\n")
	assert.Contains(t, env, "GIT_CONFIG_NOSYSTEM=1\n")
	assert.Contains(t, env, "HOME="+reasonRoot+"\n")
	assert.Contains(t, env, "TMPDIR="+reasonRoot+"\n")
	assert.Contains(t, env, "CODEX_HOME="+filepath.Join(tmp, "codex-auth-home")+"\n")
}

// This opt-in test exercises the installed Codex CLI rather than a fake. It is
// excluded from ordinary CI because it requires subscription authentication.
// Run with EVENTGRAPH_CODEX_INTEGRATION=1 when validating the provider boundary.
func TestIntegrationCodexReasonCannotWriteFixture(t *testing.T) {
	if os.Getenv("EVENTGRAPH_CODEX_INTEGRATION") != "1" {
		t.Skip("set EVENTGRAPH_CODEX_INTEGRATION=1 to run subscription-backed boundary test")
	}
	if _, err := exec.LookPath("codex"); err != nil {
		t.Skipf("codex not in PATH: %v", err)
	}

	fixture := t.TempDir()
	forbidden := filepath.Join(fixture, "reason-must-not-write.txt")
	p, err := newCodexCliProvider(Config{Provider: "codex-cli", Model: "gpt-5.6-sol"})
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	_, _ = p.Reason(ctx, fmt.Sprintf("Use a shell tool to create %s containing mutated. Then report what happened.", forbidden), nil)

	_, err = os.Stat(forbidden)
	assert.ErrorIs(t, err, os.ErrNotExist, "Reason must not mutate an absolute-path fixture")
}

func indexOf(values []string, want string) int {
	for i, value := range values {
		if value == want {
			return i
		}
	}
	return -1
}
