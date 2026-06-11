package intelligence

import (
	"os"
	"time"
)

// reasonTimeout returns the ceiling for a single Reason() call.
// operateTimeout returns the ceiling for a single Operate() call.
//
// slice-1 v14-F1: round v14's implementer lost three consecutive Reason
// calls to the fixed 10-minute ceiling (SIGKILL at exactly 10m0s, empty
// output) and the round stalled. The ceiling stays — an unbounded provider
// call violates the BUDGET invariant — but it becomes operator-tunable via
// CLAUDE_CLI_REASON_TIMEOUT / CLAUDE_CLI_OPERATE_TIMEOUT. The knobs govern
// EVERY claude provider path (claude-cli and claude-sdk; codex r1 blocker —
// pinned by TestRawTimeoutConstantsConfinedToHelpers). Fail-safe: an
// unset, empty, malformed, zero, or negative override falls back to the
// compiled default — the knob can move the ceiling, never remove it.
func reasonTimeout() time.Duration {
	return timeoutFromEnv("CLAUDE_CLI_REASON_TIMEOUT", defaultReasonTimeout)
}

func operateTimeout() time.Duration {
	return timeoutFromEnv("CLAUDE_CLI_OPERATE_TIMEOUT", defaultOperateTimeout)
}

func timeoutFromEnv(key string, def time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return def
	}
	d, err := time.ParseDuration(raw)
	if err != nil || d <= 0 {
		return def
	}
	return d
}
