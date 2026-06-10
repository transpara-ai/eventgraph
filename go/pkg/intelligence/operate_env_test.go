package intelligence

import (
	"os"
	"os/exec"
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

// buildOperateEnv runs operateSubprocessEnv and registers its cleanup with the
// test, returning the env map.
func buildOperateEnv(t *testing.T, parent []string) ([]string, map[string]string) {
	t.Helper()
	env, cleanup, err := operateSubprocessEnv(parent)
	if err != nil {
		t.Fatalf("operateSubprocessEnv: %v", err)
	}
	t.Cleanup(cleanup)
	return env, envMap(env)
}

// TestOperateSubprocessEnv_StripsRemoteCredentials is the v10-F1 keystone: the
// Operate subprocess (headless coding agent with Bash) must not inherit any
// channel that lets `git push` or `gh` reach a remote by default. On the v10
// run the implementer pushed a branch and opened a PR with the daemon's ambient
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
		"GIT_ASKPASS=/usr/bin/askpass",
		"SSH_ASKPASS=/usr/bin/ssh-askpass",
		"SSH_AUTH_SOCK=/run/agent.sock",
		"DISPLAY=:0",
		"GIT_CONFIG_PARAMETERS='credential.helper=evil'",
		"GIT_SSH=/usr/bin/evil-ssh",
		"CLAUDECODE=1",
		"DATABASE_URL=postgres://x",
		"HIVE_AGENT_ID=actor_x",
		"HIVE_HUMAN_ID=actor_h",
		"ANTHROPIC_API_KEY=sk-ant",
		"HIVE_ANTHROPIC_API_KEY=sk-ant2",
	}

	_, m := buildOperateEnv(t, parent)

	// Every credential-bearing variable must be gone (empty counts as gone for
	// the token vars — gh treats an empty token as unauthenticated). The askpass
	// and SSH-agent channels must be gone too: GIT_TERMINAL_PROMPT=0 does not
	// stop an inherited askpass, and an agent socket re-enables SSH key auth.
	for _, k := range []string{
		"GH_TOKEN", "GITHUB_TOKEN", "GH_ENTERPRISE_TOKEN", "GITHUB_ENTERPRISE_TOKEN",
		"GIT_ASKPASS", "SSH_ASKPASS", "SSH_AUTH_SOCK", "DISPLAY",
		"GIT_CONFIG_PARAMETERS", "GIT_SSH",
		"CLAUDECODE", "DATABASE_URL", "HIVE_AGENT_ID", "HIVE_HUMAN_ID",
		"ANTHROPIC_API_KEY", "HIVE_ANTHROPIC_API_KEY",
	} {
		if v, ok := m[k]; ok && v != "" {
			t.Errorf("credential/askpass var %q survived with value %q; the default remote path is still reachable", k, v)
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
// ~/.ssh and credential.helper): the builder must neutralize each default
// store, not just env tokens.
func TestOperateSubprocessEnv_NeutralizesCredentialStores(t *testing.T) {
	_, m := buildOperateEnv(t, []string{"HOME=/home/op", "PATH=/usr/bin"})

	// gh: no token + an isolated empty config dir => unauthenticated.
	ghDir, ok := m["GH_CONFIG_DIR"]
	if !ok || ghDir == "" {
		t.Fatal("GH_CONFIG_DIR not set; gh would read ~/.config/gh/hosts.yml and stay authenticated")
	}
	info, statErr := os.Stat(ghDir)
	if statErr != nil || !info.IsDir() {
		t.Fatalf("GH_CONFIG_DIR %q is not an existing directory: %v", ghDir, statErr)
	}
	if info.Mode().Perm() != 0o700 {
		t.Errorf("GH_CONFIG_DIR %q mode = %v; want 0700 (MkdirTemp)", ghDir, info.Mode().Perm())
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
	for _, want := range []string{"IdentitiesOnly=yes", os.DevNull, "IdentityAgent=none", "BatchMode=yes"} {
		if !strings.Contains(ssh, want) {
			t.Errorf("GIT_SSH_COMMAND = %q; must contain %q so ~/.ssh keys are not offered", ssh, want)
		}
	}
}

// TestOperateSubprocessEnv_PreservesLocalCommitIdentity guards the legitimate
// path: dropping global/system git config removes user.name/user.email AND the
// usual safe.directory entry, so the builder must supply a factory identity and
// safe.directory=* via env or `git commit` (and the commit-verification gate)
// breaks.
func TestOperateSubprocessEnv_PreservesLocalCommitIdentity(t *testing.T) {
	_, m := buildOperateEnv(t, []string{"HOME=/home/op"})
	for _, k := range []string{"GIT_AUTHOR_NAME", "GIT_AUTHOR_EMAIL", "GIT_COMMITTER_NAME", "GIT_COMMITTER_EMAIL"} {
		if m[k] == "" {
			t.Errorf("%s is empty; local `git commit` would fail with no global config", k)
		}
	}
	// safe.directory=* injected via the env mechanism (independent of the
	// dropped file config).
	if m["GIT_CONFIG_COUNT"] != "1" || m["GIT_CONFIG_KEY_0"] != "safe.directory" || m["GIT_CONFIG_VALUE_0"] != "*" {
		t.Errorf("safe.directory not injected via GIT_CONFIG_COUNT (count=%q key=%q val=%q); a foreign-owned worktree commit would fail dubious-ownership",
			m["GIT_CONFIG_COUNT"], m["GIT_CONFIG_KEY_0"], m["GIT_CONFIG_VALUE_0"])
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
	env, m := buildOperateEnv(t, parent)

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
	if strings.Contains(m["GIT_SSH_COMMAND"], "id_ed25519") {
		t.Error("inherited GIT_SSH_COMMAND with a real key survived; the neutralizer must override it")
	}
	if m["GH_CONFIG_DIR"] == "/home/op/.config/gh" {
		t.Error("inherited GH_CONFIG_DIR pointing at the real gh config survived")
	}
}

// TestOperateSubprocessEnv_GhDirPerCallAndCleaned proves the gh config dir is
// per-call (no shared /tmp state a stale or planted hosts.yml could re-auth
// from) and that cleanup removes it (codex review of #50, finding 4).
func TestOperateSubprocessEnv_GhDirPerCallAndCleaned(t *testing.T) {
	env1, c1, err := operateSubprocessEnv([]string{"HOME=/home/op"})
	if err != nil {
		t.Fatalf("operateSubprocessEnv #1: %v", err)
	}
	env2, c2, err := operateSubprocessEnv([]string{"HOME=/home/op"})
	if err != nil {
		t.Fatalf("operateSubprocessEnv #2: %v", err)
	}
	dir1, dir2 := envMap(env1)["GH_CONFIG_DIR"], envMap(env2)["GH_CONFIG_DIR"]
	if dir1 == dir2 {
		t.Fatalf("two calls shared the same GH_CONFIG_DIR %q; a stale/planted hosts.yml would re-authenticate gh", dir1)
	}
	if !strings.HasPrefix(dir1, filepath.Clean(os.TempDir())) {
		t.Errorf("GH_CONFIG_DIR = %q; want it under the OS temp dir, not HOME or a workspace", dir1)
	}
	c1()
	if _, err := os.Stat(dir1); !os.IsNotExist(err) {
		t.Errorf("cleanup did not remove GH_CONFIG_DIR %q (stat err=%v)", dir1, err)
	}
	c2()
}

// TestOperateSubprocessEnv_LocalCommitSucceeds is the integration proof for
// codex finding 6: under the full Operate env (global+system config dropped),
// a real `git commit` in a fresh repo still succeeds via the env-supplied
// identity and safe.directory. Skips if git is unavailable.
func TestOperateSubprocessEnv_LocalCommitSucceeds(t *testing.T) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git not available")
	}
	repo := t.TempDir()
	env, cleanup, err := operateSubprocessEnv(os.Environ())
	if err != nil {
		t.Fatalf("operateSubprocessEnv: %v", err)
	}
	t.Cleanup(cleanup)

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command(gitPath, args...)
		cmd.Dir = repo
		cmd.Env = env
		if out, e := cmd.CombinedOutput(); e != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), e, out)
		}
	}
	run("init")
	if err := os.WriteFile(filepath.Join(repo, "f.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	run("add", "f.txt")
	run("commit", "-m", "factory commit")

	// Identity came from the env, with no global/system config.
	cmd := exec.Command(gitPath, "log", "-1", "--format=%an <%ae>")
	cmd.Dir = repo
	cmd.Env = env
	out, e := cmd.CombinedOutput()
	if e != nil {
		t.Fatalf("git log: %v\n%s", e, out)
	}
	if got := strings.TrimSpace(string(out)); got != "hive-factory <factory@hive.local>" {
		t.Errorf("commit identity = %q; want the env-supplied factory identity", got)
	}
}

// TestOperateSubprocessEnv_DefeatsGitConfigParametersInjection is the NEW-F1
// integration proof (codex re-review of #50): GIT_CONFIG_PARAMETERS is honored
// by git even with global/system config dropped and GIT_CONFIG_COUNT bounded,
// so an inherited one can reintroduce credential.helper. The Operate env must
// strip it — proven by running git under the built env and confirming the
// injected credential.helper is absent. Skips if git is unavailable.
func TestOperateSubprocessEnv_DefeatsGitConfigParametersInjection(t *testing.T) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git not available")
	}
	repo := t.TempDir()

	// Control: a parent carrying a hostile GIT_CONFIG_PARAMETERS would inject
	// credential.helper into every git invocation.
	parent := append(os.Environ(), "GIT_CONFIG_PARAMETERS='credential.helper=evilhelper'")
	env, cleanup, err := operateSubprocessEnv(parent)
	if err != nil {
		t.Fatalf("operateSubprocessEnv: %v", err)
	}
	t.Cleanup(cleanup)

	initCmd := exec.Command(gitPath, "init")
	initCmd.Dir = repo
	initCmd.Env = env
	if out, e := initCmd.CombinedOutput(); e != nil {
		t.Fatalf("git init: %v\n%s", e, out)
	}

	cmd := exec.Command(gitPath, "config", "--get", "credential.helper")
	cmd.Dir = repo
	cmd.Env = env
	out, _ := cmd.CombinedOutput() // exit 1 when unset — that is the pass case
	if strings.Contains(string(out), "evilhelper") {
		t.Fatalf("GIT_CONFIG_PARAMETERS injection survived: credential.helper = %q; the ambient config-injection channel is still open", strings.TrimSpace(string(out)))
	}
}
