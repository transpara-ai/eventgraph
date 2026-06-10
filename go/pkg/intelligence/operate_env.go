package intelligence

import (
	"fmt"
	"os"
)

// operateCredentialVars are the environment variables stripped from every
// Operate subprocess. A headless coding agent runs with the Bash tool and
// --dangerously-skip-permissions, so any credential it inherits is a path to an
// ungoverned remote mutation (slice-1 v10-F1: the implementer pushed a branch
// and opened a PR with the daemon's ambient credentials). The governed
// create-PR path runs in the daemon and keeps these; the subprocess must not.
var operateCredentialVars = []string{
	// GitHub tokens (gh + git-over-https).
	"GH_TOKEN", "GITHUB_TOKEN", "GH_ENTERPRISE_TOKEN", "GITHUB_ENTERPRISE_TOKEN",
	// Credential-prompt and SSH-agent channels: GIT_TERMINAL_PROMPT=0 does not
	// stop git from calling an inherited askpass, and an inherited agent socket
	// re-enables SSH key auth despite IdentityAgent=none (codex review of
	// eventgraph#50, findings 1 + 3).
	"GIT_ASKPASS", "SSH_ASKPASS", "SSH_AUTH_SOCK", "DISPLAY",
	// Hive/daemon internals that must never leak into a subprocess.
	"CLAUDECODE", "DATABASE_URL", "HIVE_AGENT_ID", "HIVE_HUMAN_ID",
	// Anthropic keys poison the claude CLI's Max-subscription auth.
	"ANTHROPIC_API_KEY", "HIVE_ANTHROPIC_API_KEY",
}

// operateGitIdentity is the git author/committer supplied to the Operate
// subprocess. Dropping global+system git config (below) removes user.name and
// user.email, so a local `git commit` would fail — and the commit-verification
// gate would never pass — without this. A recognizable factory identity also
// aids forensics on any commit a run produces.
var operateGitIdentity = map[string]string{
	"GIT_AUTHOR_NAME":     "hive-factory",
	"GIT_AUTHOR_EMAIL":    "factory@hive.local",
	"GIT_COMMITTER_NAME":  "hive-factory",
	"GIT_COMMITTER_EMAIL": "factory@hive.local",
}

// noopCleanup is returned alongside every error so callers can unconditionally
// `defer cleanup()`.
func noopCleanup() {}

// operateSubprocessEnv builds the environment for an Operate coding subprocess
// and returns a cleanup func the caller must defer (it removes the per-call gh
// config dir). It strips the daemon's ambient remote credentials and
// neutralizes the credential STORES an LLM coding agent reaches by default, so
// the obvious `git push` / `gh pr create` an Operate would run fails — while a
// local `git commit` still works.
//
// SCOPE / THREAT MODEL — this is defense in depth, NOT an OS sandbox. It closes
// the DEFAULT credential path (the v10 breach: an aligned-but-mistaken agent
// following planner gates that echoed "push a branch, open a PR" used the
// ambient credentials). It does NOT contain a deliberately adversarial Bash
// child: HOME is preserved (the claude CLI's Max-subscription auth lives in
// ~/.claude), so a child that constructs absolute paths could still read
// ~/.ssh or point GH_CONFIG_DIR back at the real config, and a run workspace's
// local .git/config is trusted. A hard boundary (separate uid / container /
// namespace, credential-less run workspaces) is the slice's sandbox layer,
// routed to G-2.x.
//
// Defense in depth, because an env-var denylist alone leaves stores reachable
// (gh reads ~/.config/gh/hosts.yml; git falls back to ~/.ssh, an askpass, an
// agent socket, and a configured credential.helper):
//   - strip credential-bearing env vars (operateCredentialVars), including the
//     askpass + SSH-agent channels;
//   - GH_CONFIG_DIR -> a per-call empty dir (MkdirTemp: unique, 0700,
//     owned-by-us, no symlink/stale-state race), so gh is unauthenticated;
//   - GIT_CONFIG_GLOBAL=/dev/null + GIT_CONFIG_NOSYSTEM=1 drop credential.helper
//     and the user identity from ~/.gitconfig and /etc/gitconfig...
//   - ...so a factory identity is supplied via GIT_AUTHOR_*/GIT_COMMITTER_*, and
//     safe.directory=* is injected via GIT_CONFIG_COUNT (dropping global config
//     also drops the usual safe.directory entry, which would otherwise break a
//     commit in a foreign-owned worktree with "dubious ownership");
//   - GIT_SSH_COMMAND offers no identity (IdentitiesOnly + IdentityFile=/dev/null
//     + IdentityAgent=none) so SSH remotes cannot authenticate;
//   - GIT_TERMINAL_PROMPT=0 so a missing credential fails fast.
//
// Fails closed: if the isolated gh config dir cannot be created, it returns an
// error (and a no-op cleanup) rather than letting the subprocess run ambient.
func operateSubprocessEnv(parent []string) (env []string, cleanup func(), err error) {
	ghConfigDir, err := os.MkdirTemp("", "hive-operate-gh-")
	if err != nil {
		return nil, noopCleanup, fmt.Errorf("operate env: isolate gh config: %w", err)
	}
	cleanup = func() { _ = os.RemoveAll(ghConfigDir) }

	// The managed keys we set ourselves — strip every inherited occurrence so
	// the neutralizer wins over any hostile inherited value and appears once.
	managed := map[string]string{
		"GH_CONFIG_DIR":       ghConfigDir,
		"GIT_CONFIG_GLOBAL":   os.DevNull,
		"GIT_CONFIG_NOSYSTEM": "1",
		"GIT_TERMINAL_PROMPT": "0",
		"GIT_SSH_COMMAND":     "ssh -o IdentitiesOnly=yes -o IdentityFile=/dev/null -o IdentityAgent=none -o BatchMode=yes",
		// safe.directory=* via the env-injection mechanism (independent of the
		// dropped file config) so `git commit`/`git status` work in a worktree
		// git would otherwise reject for dubious ownership.
		"GIT_CONFIG_COUNT":   "1",
		"GIT_CONFIG_KEY_0":   "safe.directory",
		"GIT_CONFIG_VALUE_0": "*",
	}
	for k, v := range operateGitIdentity {
		managed[k] = v
	}

	strip := append([]string{}, operateCredentialVars...)
	for k := range managed {
		strip = append(strip, k)
	}

	env = scrubEnv(parent, strip...)
	for k, v := range managed {
		env = append(env, k+"="+v)
	}
	return env, cleanup, nil
}
