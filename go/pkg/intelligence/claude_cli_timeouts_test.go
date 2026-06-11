package intelligence

import (
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
