package intelligence

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// envMap collapses an environ slice to a map (last occurrence wins, matching
// exec semantics) so tests can assert on final effective values.
func envMap(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, e := range env {
		if i := strings.IndexByte(e, '='); i >= 0 {
			m[e[:i]] = e[i+1:]
		}
	}
	return m
}

// TestOperateSubprocessEnv_StripsRemoteCredentials is the v10-F1 keystone: the
// Operate subprocess (headless coding agent with Bash) must not inherit any
// channel that lets `git push` or `gh` reach a remote. On the v10 run the
// implementer pushed a branch and opened a PR with the daemon's ambient
// credentials; the governed create-PR path in the daemon keeps its credentials,
// the subprocess must not.
func TestOperateSubprocessEnv_StripsRemoteCredentials(t *testing.T) {
	parent := []string{
		"PATH=/usr/bin",
		"HOME=/home/op",
		"GH_TOKEN=ghp_secret",
		"GITHUB_TOKEN=gho_secret",
		"GH_ENTERPRISE_TOKEN=ghe_secret",
		"GITHUB_ENTERPRISE_TOKEN=ghe2_secret",
		"CLAUDECODE=1",
		"DATABASE_URL=postgres://x",
		"HIVE_AGENT_ID=actor_x",
		"HIVE_HUMAN_ID=actor_h",
		"ANTHROPIC_API_KEY=sk-ant",
		"HIVE_ANTHROPIC_API_KEY=sk-ant2",
	}

	env, err := operateSubprocessEnv(parent)
	if err != nil {
		t.Fatalf("operateSubprocessEnv: %v", err)
	}
	m := envMap(env)

	// Every credential-bearing variable must be gone (empty counts as gone for
	// the token vars — gh treats an empty token as unauthenticated).
	for _, k := range []string{
		"GH_TOKEN", "GITHUB_TOKEN", "GH_ENTERPRISE_TOKEN", "GITHUB_ENTERPRISE_TOKEN",
		"CLAUDECODE", "DATABASE_URL", "HIVE_AGENT_ID", "HIVE_HUMAN_ID",
		"ANTHROPIC_API_KEY", "HIVE_ANTHROPIC_API_KEY",
	} {
		if v, ok := m[k]; ok && v != "" {
			t.Errorf("credential var %q survived with value %q; remote access is still reachable", k, v)
		}
	}

	// HOME is preserved — the claude CLI's Max-subscription auth lives in
	// ~/.claude and must keep working; PATH too.
	if m["HOME"] != "/home/op" {
		t.Errorf("HOME = %q; want it preserved (claude CLI auth lives under it)", m["HOME"])
	}
	if m["PATH"] != "/usr/bin" {
		t.Errorf("PATH = %q; want it preserved", m["PATH"])
	}
}

// TestOperateSubprocessEnv_NeutralizesCredentialStores proves env-scrubbing
// alone is not enough (gh reads ~/.config/gh/hosts.yml; git falls back to
// ~/.ssh and credential.helper): the builder must neutralize every store, not
// just env tokens.
func TestOperateSubprocessEnv_NeutralizesCredentialStores(t *testing.T) {
	env, err := operateSubprocessEnv([]string{"HOME=/home/op", "PATH=/usr/bin"})
	if err != nil {
		t.Fatalf("operateSubprocessEnv: %v", err)
	}
	m := envMap(env)

	// gh: no token + an isolated empty config dir => unauthenticated.
	ghDir, ok := m["GH_CONFIG_DIR"]
	if !ok || ghDir == "" {
		t.Fatal("GH_CONFIG_DIR not set; gh would read ~/.config/gh/hosts.yml and stay authenticated")
	}
	info, statErr := os.Stat(ghDir)
	if statErr != nil || !info.IsDir() {
		t.Fatalf("GH_CONFIG_DIR %q is not an existing directory: %v", ghDir, statErr)
	}
	if entries, _ := os.ReadDir(ghDir); len(entries) != 0 {
		t.Fatalf("GH_CONFIG_DIR %q is not empty; it must carry no hosts.yml", ghDir)
	}

	// git: no global/system config (drops credential.helper and any cached
	// identity) and no terminal prompt.
	if m["GIT_CONFIG_GLOBAL"] != os.DevNull {
		t.Errorf("GIT_CONFIG_GLOBAL = %q; want %q to drop ~/.gitconfig credential helpers", m["GIT_CONFIG_GLOBAL"], os.DevNull)
	}
	if m["GIT_CONFIG_NOSYSTEM"] != "1" {
		t.Errorf("GIT_CONFIG_NOSYSTEM = %q; want \"1\" to drop /etc/gitconfig", m["GIT_CONFIG_NOSYSTEM"])
	}
	if m["GIT_TERMINAL_PROMPT"] != "0" {
		t.Errorf("GIT_TERMINAL_PROMPT = %q; want \"0\" so a missing credential fails instead of prompting", m["GIT_TERMINAL_PROMPT"])
	}

	// git over SSH: no identity offered, so SSH remotes cannot authenticate.
	ssh := m["GIT_SSH_COMMAND"]
	for _, want := range []string{"IdentitiesOnly=yes", os.DevNull, "BatchMode=yes"} {
		if !strings.Contains(ssh, want) {
			t.Errorf("GIT_SSH_COMMAND = %q; must contain %q so ~/.ssh keys are not offered", ssh, want)
		}
	}
}

// TestOperateSubprocessEnv_PreservesLocalCommitIdentity guards the legitimate
// path: dropping global/system git config removes user.name/user.email, so the
// builder must supply a factory identity via env or `git commit` (and therefore
// the commit-verification gate) breaks.
func TestOperateSubprocessEnv_PreservesLocalCommitIdentity(t *testing.T) {
	env, err := operateSubprocessEnv([]string{"HOME=/home/op"})
	if err != nil {
		t.Fatalf("operateSubprocessEnv: %v", err)
	}
	m := envMap(env)
	for _, k := range []string{"GIT_AUTHOR_NAME", "GIT_AUTHOR_EMAIL", "GIT_COMMITTER_NAME", "GIT_COMMITTER_EMAIL"} {
		if m[k] == "" {
			t.Errorf("%s is empty; local `git commit` would fail with no global config and the commit-verification gate would never pass", k)
		}
	}
}

// TestOperateSubprocessEnv_SingleOccurrenceOverridesInherited proves the
// neutralizers WIN over an inherited hostile value: a parent that already set
// GIT_SSH_COMMAND or GH_CONFIG_DIR to a credentialed value must not survive,
// and each managed key must appear exactly once (no ambiguity for exec).
func TestOperateSubprocessEnv_SingleOccurrenceOverridesInherited(t *testing.T) {
	parent := []string{
		"HOME=/home/op",
		"GIT_SSH_COMMAND=ssh -i /home/op/.ssh/id_ed25519",
		"GH_CONFIG_DIR=/home/op/.config/gh",
		"GIT_CONFIG_GLOBAL=/home/op/.gitconfig",
	}
	env, err := operateSubprocessEnv(parent)
	if err != nil {
		t.Fatalf("operateSubprocessEnv: %v", err)
	}

	counts := map[string]int{}
	for _, e := range env {
		if i := strings.IndexByte(e, '='); i >= 0 {
			counts[e[:i]]++
		}
	}
	for _, k := range []string{"GIT_SSH_COMMAND", "GH_CONFIG_DIR", "GIT_CONFIG_GLOBAL"} {
		if counts[k] != 1 {
			t.Errorf("%s appears %d times; managed keys must appear exactly once so the neutralizer is unambiguous", k, counts[k])
		}
	}
	m := envMap(env)
	if strings.Contains(m["GIT_SSH_COMMAND"], "id_ed25519") {
		t.Error("inherited GIT_SSH_COMMAND with a real key survived; the neutralizer must override it")
	}
	if m["GH_CONFIG_DIR"] == "/home/op/.config/gh" {
		t.Error("inherited GH_CONFIG_DIR pointing at the real gh config survived")
	}
}

// TestOperateSubprocessEnv_GhDirIsolatedPerHome keeps the throwaway gh config
// dir off any path a user's real config occupies, and confirms it is created
// under the OS temp area (not inside HOME or a run workspace that might be
// committed).
func TestOperateSubprocessEnv_GhDirIsolated(t *testing.T) {
	env, err := operateSubprocessEnv([]string{"HOME=/home/op"})
	if err != nil {
		t.Fatalf("operateSubprocessEnv: %v", err)
	}
	ghDir := envMap(env)["GH_CONFIG_DIR"]
	if !strings.HasPrefix(ghDir, filepath.Clean(os.TempDir())) {
		t.Errorf("GH_CONFIG_DIR = %q; want it under the OS temp dir, not HOME or a workspace", ghDir)
	}
}
