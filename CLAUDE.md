# EventGraph

## Soul

> Take care of your human, humanity, and yourself.

This is infrastructure for accountable AI and human systems. Every design decision serves that soul statement.

## What This Is

A hash-chained, append-only, causal event graph — the foundation for building systems where every action is signed, auditable, and causally linked. Built as a set of packages anyone can import, extend, and build on.

The event graph is not a product. It is the infrastructure that products are built on. Social networks, governance systems, marketplaces, identity systems — these are interfaces on the same substrate.

## Architecture

### The Event

The fundamental unit. Every significant action is an event:

```
Event {
    ID:             string      // UUID v7 (time-ordered)
    Type:           string      // hierarchical (e.g., "trust.updated")
    Timestamp:      time        // nanosecond precision
    Source:         string      // who/what emitted this
    Content:        object      // structured payload
    Causes:         []string    // IDs of causing events (causal DAG)
    ConversationID: string      // thread grouping
    Hash:           string      // SHA-256 of canonical form
    PrevHash:       string      // hash chain link
    Signature:      string      // Ed25519 signature
}
```

Every event has declared causes (except Bootstrap). All events are hash-chained. All operations emit events. The event graph is the source of truth.

### Core Interfaces

These are the extension points. Implement them to plug in your own infrastructure:

- **`Store`** — Event persistence. Reference: PostgresStore. Implement for your database.
- **`IIntelligence`** — Anything that reasons. Implement with Claude, GPT, local models, or deterministic logic. Not every primitive needs intelligence — most start mechanical and grow toward it only when needed.
- **`IDecisionMaker`** — The decision tree engine. Routes decisions through deterministic branches first, falls through to IIntelligence only when the tree can't handle it. Evolves over time — expensive model calls become cheap deterministic rules as patterns emerge.

### Primitives

A primitive is a software agent that embodies a specific domain of intelligence. 200 primitives across 14 layers, from foundational (Event, Hash, Clock) to existential (Being, Wonder, Mystery).

Each primitive has:
- **Name** and **Layer** — position in the ontological hierarchy
- **Activation** — 0.0 to 1.0, current engagement level
- **State** — mutable key-value store
- **Lifecycle** — dormant → activating → active → processing → emitting → active (or → deactivating → dormant)
- **Subscriptions** — event type prefixes this primitive listens to
- **Cadence** — minimum ticks between invocations
- **Decision Tree** — optional, evolving, mechanical-to-intelligent

The primitive interface:

```
Process(tick, events, snapshot) → []Mutation
```

Mutations are declarative: AddEvent, UpdateState, UpdateActivation, UpdateLifecycle. The tick engine collects and applies them atomically.

### Tick Engine (Ripple-Wave Processor)

The system's heartbeat. Each tick:
1. Snapshot all primitive states (read-only, shared)
2. Distribute pending events to subscribing primitives
3. Invoke each primitive's Process function (subject to cadence + lifecycle)
4. Collect mutations
5. Apply atomically — new events become input for next wave
6. Repeat until quiescence or max waves (10)
7. Persist all primitive states

Layer constraint: Layer N primitives activate only when Layer N-1 is stable.

### Authority

Three-tier approval for significant actions:
- **Required** — blocks until human approves/rejects
- **Recommended** — auto-approves after 15 min timeout
- **Notification** — auto-approves immediately, logged for audit

Authority requests and resolutions are events on the graph.

### Inter-System Protocol (EGIP)

Sovereign systems communicate without shared infrastructure:
- Ed25519 identity (no central registry)
- Cross-Graph Event References (CGERs) — causal links across graph boundaries
- Signed message envelopes
- Seven message types: HELLO, MESSAGE, RECEIPT, PROOF, TREATY, AUTHORITY_REQUEST, DISCOVER
- Treaties for bilateral governance
- Trust accumulation (0.0-1.0, asymmetric, non-transitive)

## Current State

See `ROADMAP.md` for what's built and what needs building.

## How to Contribute

### Finding Work

1. Check `ROADMAP.md` for unclaimed tasks
2. Check GitHub Issues for open tasks
3. If you see something that needs doing that isn't listed, open an issue first

### Before You Code

1. Read this file completely
2. Read `CONTRIBUTING.md` for process
3. Read the relevant `docs/` files for the area you're working in
4. Read the relevant `docs/coding-standards/` file for your language

### Coding Standards

**All languages:**
- Every public interface must have documentation
- Every primitive must have tests
- Tests must cover: happy path, error cases, edge cases
- No silent failures — errors must be visible
- Hash chain integrity must be maintained at all times
- Causal links must be declared on every event
- Every store implementation must pass the shared conformance test suite

**Coverage requirements:**
- Core packages (event, store, bus, primitive, tick): **90%** minimum
- Primitive implementations: **80%** minimum
- Protocol packages: **85%** minimum
- Utility packages: **70%** minimum

**Go (reference implementation):**
- See `docs/coding-standards/go.md`
- `go vet`, `staticcheck`, and `golangci-lint` must pass
- No `interface{}` / `any` except in Event.Content (which is typed at the primitive level)
- Errors are values, not panics. Only panic on unrecoverable invariant violations.
- Table-driven tests preferred

### PR Requirements

Your PR must:

1. **Pass CI** — build, test, lint, coverage thresholds
2. **Be self-audited** — Before submitting, audit your own code thoroughly. Run through it multiple times looking for:
   - Logic errors
   - Missing error handling
   - Race conditions
   - Hash chain integrity violations
   - Missing causal links on events
   - Interface contract violations
   - Test coverage gaps
   - Documentation gaps on public interfaces
   - Security issues (injection, overflow, timing attacks)
3. **Include tests** — No code without tests. No exceptions.
4. **Update docs** — If you changed an interface, update the relevant docs/ file
5. **One concern per PR** — Don't bundle unrelated changes. Each PR should do one thing.
6. **Describe what and why** — PR description explains what changed and why. Link to the issue or roadmap item.

### After PR Submission

Reviewers (human or AI) will audit the code. Common review criteria:
- Does this maintain hash chain integrity?
- Are all events causally linked?
- Does this respect the primitive lifecycle?
- Does this break any existing interfaces?
- Are the tests actually testing the right things?
- Is this the simplest solution that works?

Expect multiple rounds of review. This is infrastructure — correctness matters more than speed.

### Implementing a New Primitive

1. Check `docs/layers/` for the layer specification
2. Create the primitive in the appropriate package
3. Implement the Primitive interface (Process, lifecycle, subscriptions)
4. Define sensible defaults for all behaviours
5. Define the events this primitive emits and subscribes to
6. Write tests (80% minimum coverage)
7. Document the primitive's purpose, state schema, event types, and default behaviour
8. Submit PR referencing the layer spec

### Implementing a New Store

1. Implement the `Store` interface
2. Pass the **full conformance test suite** (`go/pkg/store/conformance_test.go`)
3. The conformance suite tests: append, get, query, causal traversal, hash chain verification, concurrent access
4. Document any limitations or configuration requirements
5. Submit PR

### Implementing a New Language Package

1. Read `docs/interfaces.md` for the interface specifications
2. Read `docs/coding-standards/` for your language — if none exists, propose one
3. Implement all core interfaces (Store, Primitive, Bus, etc.)
4. Pass the **language-agnostic conformance test suite** (defined in `docs/conformance/`)
5. Match the Go reference implementation's behaviour exactly
6. Submit PR with full test coverage

### Changing Existing Interfaces

Interface changes require an RFC:
1. Open a GitHub Issue using the RFC template
2. Describe what you want to change and why
3. Show the impact on existing implementations
4. Wait for discussion and approval before implementing

This prevents people building on shifting sand.

## What NOT to Do

- Don't break existing interfaces without an RFC
- Don't add primitives without checking the layer specification
- Don't skip tests
- Don't commit secrets, keys, or credentials
- Don't add dependencies without justification
- Don't optimise prematurely — correctness first, then performance
- Don't ignore hash chain integrity — ever
- Don't emit events without causal links — ever

## Invariants

These are non-negotiable. Every contribution must maintain them:

1. Every event has declared causes
2. All events are hash-chained
3. All operations emit events
4. The system improves itself (decision trees evolve, primitives learn)
5. Significant actions require authority
6. Build and test before done
7. The event graph is the source of truth

## Reference

- `docs/architecture.md` — System architecture
- `docs/primitives.md` — All 200 primitives across 14 layers
- `docs/grammar.md` — The 15 social grammar operations
- `docs/protocol.md` — EGIP inter-system protocol
- `docs/interfaces.md` — Core interfaces specification
- `docs/tick-engine.md` — Ripple-wave processing
- `docs/decision-trees.md` — Mechanical-to-intelligent continuum
- `docs/authority.md` — Three-tier approval system
- `docs/trust.md` — Trust model and dynamics
- `docs/layers/` — Per-layer primitive specifications
- `ROADMAP.md` — What needs building
- `CONTRIBUTING.md` — How to contribute
