# Contributing to EventGraph

We welcome contributions from humans and AI systems alike. This document describes how to contribute effectively.

## Before You Start

1. Read `CLAUDE.md` — the comprehensive project guide
2. Read `ROADMAP.md` — find work that needs doing
3. Read the relevant `docs/` files for the area you're working in
4. Check open issues — someone may already be working on what you're planning

## Ways to Contribute

### Code
- Implement primitives (see layer specs in `docs/layers/`)
- Implement store backends (Postgres, SQLite, DynamoDB, etc.)
- Implement language packages (Rust, Python, .NET, etc.)
- Fix bugs
- Improve test coverage
- Performance improvements (with benchmarks proving the improvement)

### Documentation
- Improve existing docs
- Add examples
- Translate docs
- Write tutorials

### Design
- Propose new primitives (RFC process)
- Propose interface changes (RFC process)
- Review and discuss open RFCs

## Process

### For Code Changes

1. **Find or create an issue** — Every PR should reference an issue or roadmap item
2. **Fork and branch** — Branch from `main`, name your branch descriptively
3. **Implement** — Follow coding standards for your language (`docs/coding-standards/`)
4. **Test** — Meet coverage thresholds (see `CLAUDE.md`)
5. **Self-audit** — Before submitting, review your own code thoroughly:
   - Run through the code multiple times
   - Check for logic errors, race conditions, missing error handling
   - Verify hash chain integrity is maintained
   - Verify all events have causal links
   - Check test coverage meets thresholds
   - Ensure docs are updated for any interface changes
   - Look for security issues
   - Keep auditing until you find no more issues
6. **Submit PR** — Use the PR template. Describe what and why.
7. **Review** — Expect feedback. Address all comments. This is infrastructure — correctness matters.

### For Interface Changes (RFC Process)

Changing a public interface affects everyone building on EventGraph. These changes require discussion before implementation.

1. Open a GitHub Issue using the **RFC** template
2. Describe: what you want to change, why, and the impact on existing code
3. Wait for community discussion (minimum 7 days for significant changes)
4. If approved, implement and submit a PR
5. Include migration guidance for existing users

### For New Primitives

1. Check `docs/layers/` — does your primitive fit an existing layer?
2. Open a GitHub Issue using the **Primitive** template
3. Describe: the primitive's purpose, which layer it belongs to, what gap it fills, its event types, state schema, and default behaviour
4. If approved, implement following the primitive implementation guide in `CLAUDE.md`

## Coding Standards

### All Languages

- Public interfaces must be documented
- Every primitive must have tests
- No silent failures — errors must be visible
- Hash chain integrity is non-negotiable
- Causal links on every event are non-negotiable

### Language-Specific

See `docs/coding-standards/` for your language. If standards don't exist for your language yet, propose them as part of your first PR.

## Coverage Thresholds

| Package | Minimum Coverage |
|---------|-----------------|
| Core (event, store, bus, primitive, tick) | 90% |
| Protocol (EGIP, envelopes, trust) | 85% |
| Primitive implementations | 80% |
| Utilities and helpers | 70% |

PRs that decrease coverage below these thresholds will not be merged.

## Commit Messages

Use conventional commits:

```
feat: add InMemoryStore conformance tests
fix: correct hash chain verification for empty graphs
docs: update primitive implementation guide
refactor: extract common store logic into shared interface
test: add edge cases for causal traversal
```

Keep the first line under 72 characters. Add a body if the change needs explanation.

## AI Contributors

If you're contributing via Claude Code, Cursor, Copilot, or any AI tool:

- Read `CLAUDE.md` — it's written for you
- Your contributions are held to the same standards as human contributions
- Self-audit thoroughly before submitting
- Include `Co-Authored-By` in your commit messages
- Don't guess at architecture decisions — read the docs or ask in an issue

## Code of Conduct

See `CODE_OF_CONDUCT.md`. Be decent. Build things that help people.

## Questions?

Open a GitHub Issue. Tag it with `question`.
