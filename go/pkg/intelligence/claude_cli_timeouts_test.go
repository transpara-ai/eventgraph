package intelligence

import (
	"os"
	"strings"
	"testing"
	"time"
)

// ════════════════════════════════════════════════════════════════════════
// slice-1 v14-F1: reason/operate timeout knobs
//
// The ceilings must hold by default, move only under a valid override, and
// fail SAFE (back to the default) on any unset, empty, malformed, zero, or
// negative value — a knob that can silently remove the ceiling would trade
// a visible stall for an unbounded hang.
// ════════════════════════════════════════════════════════════════════════

func TestReasonTimeoutDefault(t *testing.T) {
	t.Setenv("CLAUDE_CLI_REASON_TIMEOUT", "")
	if got := reasonTimeout(); got != defaultReasonTimeout {
		t.Fatalf("reasonTimeout() = %v; want compiled default %v", got, defaultReasonTimeout)
	}
}

func TestReasonTimeoutEnvOverride(t *testing.T) {
	t.Setenv("CLAUDE_CLI_REASON_TIMEOUT", "25m")
	if got := reasonTimeout(); got != 25*time.Minute {
		t.Fatalf("reasonTimeout() with 25m override = %v; want 25m", got)
	}
}

func TestReasonTimeoutInvalidFallsBack(t *testing.T) {
	for _, bad := range []string{"bogus", "10", "0", "0s", "-5m"} {
		t.Setenv("CLAUDE_CLI_REASON_TIMEOUT", bad)
		if got := reasonTimeout(); got != defaultReasonTimeout {
			t.Fatalf("reasonTimeout() with override %q = %v; want fail-safe default %v", bad, got, defaultReasonTimeout)
		}
	}
}

// TestRawTimeoutConstantsConfinedToHelpers pins the CLASS, not the instance
// (codex r1 blocker: claude_sdk bypassed the knobs while claude_cli got
// them). The raw defaults may appear only in their const declaration and
// the helper functions — every provider call site must route through
// reasonTimeout()/operateTimeout(), or an operator's override silently
// fails to govern that provider and the fixed failure class returns.
func TestRawTimeoutConstantsConfinedToHelpers(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read package dir: %v", err)
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") || name == "claude_cli_timeouts.go" {
			continue
		}
		src, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		for lineNo, line := range strings.Split(string(src), "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "//") {
				continue // prose may name the constants
			}
			for _, c := range []string{"defaultReasonTimeout", "defaultOperateTimeout"} {
				if !strings.Contains(line, c) {
					continue
				}
				if strings.Contains(line, "time.Minute") && strings.Contains(line, "=") {
					continue // the const declaration itself
				}
				t.Errorf("%s:%d uses raw %s — route through the timeout helpers so the env knobs govern every provider", name, lineNo+1, c)
			}
		}
	}
}

func TestOperateTimeoutDefaultOverrideAndFallback(t *testing.T) {
	t.Setenv("CLAUDE_CLI_OPERATE_TIMEOUT", "")
	if got := operateTimeout(); got != defaultOperateTimeout {
		t.Fatalf("operateTimeout() = %v; want compiled default %v", got, defaultOperateTimeout)
	}
	t.Setenv("CLAUDE_CLI_OPERATE_TIMEOUT", "45m")
	if got := operateTimeout(); got != 45*time.Minute {
		t.Fatalf("operateTimeout() with 45m override = %v; want 45m", got)
	}
	t.Setenv("CLAUDE_CLI_OPERATE_TIMEOUT", "nope")
	if got := operateTimeout(); got != defaultOperateTimeout {
		t.Fatalf("operateTimeout() with invalid override = %v; want fail-safe default %v", got, defaultOperateTimeout)
	}
}
