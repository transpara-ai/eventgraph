package intelligence

import (
	"fmt"
	"os"
	"path/filepath"
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

// operateSubprocessEnv builds the environment for an Operate coding subprocess.
// It strips the daemon's ambient remote credentials and neutralizes every
// credential STORE so an Operate that runs `git push` or `gh` inside Bash
// cannot reach a remote — while keeping local `git commit` working.
//
// Defense in depth, because an env-var denylist alone leaves stores reachable
// (gh reads ~/.config/gh/hosts.yml; git falls back to ~/.ssh and a configured
// credential.helper):
//   - strip the credential-bearing env vars (operateCredentialVars);
//   - GH_CONFIG_DIR -> an isolated empty dir, so ~/.config/gh/hosts.yml is not
//     read and gh is unauthenticated;
//   - GIT_CONFIG_GLOBAL=/dev/null + GIT_CONFIG_NOSYSTEM=1 drop credential.helper
//     and the user identity from ~/.gitconfig and /etc/gitconfig...
//   - ...so a factory identity is supplied via GIT_AUTHOR_*/GIT_COMMITTER_* to
//     keep local commits working;
//   - GIT_SSH_COMMAND offers no identity (IdentitiesOnly + IdentityFile=/dev/null
//     + IdentityAgent=none) so SSH remotes cannot authenticate;
//   - GIT_TERMINAL_PROMPT=0 so a missing credential fails fast instead of
//     hanging on a prompt.
//
// HOME and ~/.claude are deliberately untouched: the claude CLI's own Max
// subscription auth lives there and must keep working.
//
// Fails closed: if the isolated gh config dir cannot be created, it returns an
// error rather than letting the subprocess run with the ambient environment.
func operateSubprocessEnv(parent []string) ([]string, error) {
	ghConfigDir := filepath.Join(os.TempDir(), "hive-operate-gh-noauth")
	if err := os.MkdirAll(ghConfigDir, 0o700); err != nil {
		return nil, fmt.Errorf("operate env: isolate gh config: %w", err)
	}

	// The managed keys we set ourselves — strip every inherited occurrence so
	// the neutralizer wins over any hostile inherited value and appears once.
	managed := map[string]string{
		"GH_CONFIG_DIR":       ghConfigDir,
		"GIT_CONFIG_GLOBAL":   os.DevNull,
		"GIT_CONFIG_NOSYSTEM": "1",
		"GIT_TERMINAL_PROMPT": "0",
		"GIT_SSH_COMMAND":     "ssh -o IdentitiesOnly=yes -o IdentityFile=/dev/null -o IdentityAgent=none -o BatchMode=yes",
	}
	for k, v := range operateGitIdentity {
		managed[k] = v
	}

	strip := append([]string{}, operateCredentialVars...)
	for k := range managed {
		strip = append(strip, k)
	}

	env := scrubEnv(parent, strip...)
	for k, v := range managed {
		env = append(env, k+"="+v)
	}
	return env, nil
}
