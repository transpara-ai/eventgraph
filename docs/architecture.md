# Architecture

## Overview

EventGraph is a hash-chained, append-only, causal event graph. Decision governance infrastructure — not AI governance. The graph takes an `IDecisionMaker` (anything that makes decisions: AI agent, human, committee, rules engine) and records what was decided, by what, with what confidence, under what authority, causally linked to everything that led there.

```
┌─────────────────────────────────────────────┐
│            Top-Level API                     │
│    graph.Evaluate() / Record() / Query()    │
├─────────────────────────────────────────────┤
│            Product Layers                    │
│    (Social Grammar, Governance, Exchange,   │
│     Task Management)                        │
├─────────────────────────────────────────────┤
│         Cognitive Primitives                 │
│    (201 across 14 layers, tick engine)       │
├─────────────────────────────────────────────┤
│           Communication                      │
│    (listen/say, routing, EGIP)              │
├─────────────────────────────────────────────┤
│            Event Graph                       │
│    (events, hash chain, causal DAG, bus)    │
├─────────────────────────────────────────────┤
│              Store (plugin)                  │
│    (Postgres, SQLite, Memory, your own)     │
└─────────────────────────────────────────────┘
```

## Two API Levels

**Top level** — four lines to make any system auditable. Developers describe what their system wants to do, the graph evaluates it through the ontology, and returns a decision with a cryptographic receipt. They don't need to understand the 201 primitives.

**Primitive level** — for power users building AI agent frameworks, compliance platforms, or their own implementations. Every primitive is an interface with sensible defaults. Override with domain-specific logic. Custom trust decay, custom authority chains, custom confidence models.

## The Event

The fundamental unit. Every significant action is an event.

```
Event {
    Version:        int             // schema version for this event type (starts at 1)
    ID:             EventID         // UUID v7 (time-ordered)
    Type:           EventType       // validated against EventTypeRegistry (e.g., "trust.updated")
    Timestamp:      time            // nanosecond precision
    Source:         ActorID         // who/what emitted this
    Content:        EventContent    // typed per EventType via EventTypeRegistry
    Causes:         NonEmpty<EventID> // causal DAG — at least one cause (except Bootstrap)
    ConversationID: ConversationID  // thread grouping
    Hash:           Hash            // SHA-256 of canonical form
    PrevHash:       Hash            // hash of previous event — linear chain
    Signature:      Signature       // Ed25519 signature
}
```

All IDs are typed — `EventID`, `ActorID`, `ConversationID`, `Hash` are distinct types, not bare strings. See `interfaces.md` for the full type system.

Two overlapping structures:
- **Linear hash chain** (PrevHash) — ordering and tamper detection
- **Causal DAG** (Causes) — reasoning about why events occurred

Both are maintained simultaneously. The chain provides integrity. The DAG provides meaning.

## Hash Chain

On each append:
1. Acquire exclusive lock on most recent event
2. Retrieve previous hash
3. Construct canonical string: `prev_hash|id|type|source|conversation_id|timestamp_nanos|content_json`
4. Compute SHA-256
5. Store event with hash and prev_hash

Verification walks the entire chain, recomputing each hash. Any tampering breaks the chain at the modification point.

## Store Interface

All methods return `Result<T, StoreError>` with typed domain errors. Queries use cursor-based pagination via `Page<T>`. See `interfaces.md` for the full type system.

```
Store {
    // Persistence — takes a pre-built Event from EventFactory
    Append(event Event) → Result<Event, StoreError>
    Get(id EventID) → Result<Event, StoreError>
    Head() → Result<Option<Event>, StoreError>

    // Queries — all return paginated results
    Recent(limit int, after Option<Cursor>) → Result<Page<Event>, StoreError>
    ByType(type EventType, limit int, after Option<Cursor>) → Result<Page<Event>, StoreError>
    BySource(source ActorID, limit int, after Option<Cursor>) → Result<Page<Event>, StoreError>
    ByConversation(id ConversationID, limit int, after Option<Cursor>) → Result<Page<Event>, StoreError>
    Since(afterID EventID, limit int) → Result<Page<Event>, StoreError>
    Search(query string, limit int, after Option<Cursor>) → Result<Page<Event>, StoreError>

    // Causal traversal
    Ancestors(id EventID, maxDepth int) → Result<[]Event, StoreError>
    Descendants(id EventID, maxDepth int) → Result<[]Event, StoreError>

    // Edge queries (typed, weighted relationships)
    EdgesFrom(id ActorID, edgeType EdgeType) → Result<[]Edge, StoreError>
    EdgesTo(id ActorID, edgeType EdgeType) → Result<[]Edge, StoreError>
    EdgeBetween(from ActorID, to ActorID, edgeType EdgeType) → Result<Option<Edge>, StoreError>

    // Integrity
    Count() → Result<int, StoreError>
    VerifyChain() → Result<ChainVerifiedContent, StoreError>
}

IActorStore {
    Register(publicKey PublicKey, displayName string, actorType ActorType) → Result<IActor, StoreError>
    Get(id ActorID) → Result<IActor, StoreError>
    GetByPublicKey(publicKey PublicKey) → Result<IActor, StoreError>
    Update(id ActorID, updates ActorUpdate) → Result<IActor, StoreError>
    List(filter ActorFilter) → Result<Page<IActor>, StoreError>
    Suspend(id ActorID, reason EventID) → Result<IActor, StoreError>
    Memorial(id ActorID, reason EventID) → Result<IActor, StoreError>
}
```

Store handles events and edges. IActorStore handles actors. Both share the backing database but are separate interfaces (single responsibility). Edges are events — creating a trust relationship emits an event. The Store indexes edge-creating events for efficient traversal.

Key type safety patterns:
- `NonEmpty<EventID>` for causes — every event (except Bootstrap) must have at least one cause
- `Option<Cursor>` for pagination — `None` means start from the beginning
- `Option<Edge>` for EdgeBetween — no edge between two actors is valid, not an error
- `Page<T>` wraps results with cursor and `HasMore` flag
- Events created by `EventFactory`, persisted by `Store` (separation of creation and persistence)
- Content validated against `EventTypeRegistry` by the factory before reaching the Store

See `interfaces.md` for the full type system (typed IDs, constrained numerics, state machines, domain errors, Edge data structure, EdgeType enum, EventTypeRegistry).

Implement these interfaces for any backing store. Reference implementations: InMemoryStore, PostgresStore.

## Bus

Wraps a Store with pub/sub fan-out. When an event is appended, all subscribers receive it. Non-blocking — slow subscribers get dropped events, not blocked writers.

## Primitives

A cognitive primitive is a software agent for a specific domain. See `primitives.md` for the full specification.

Key properties: name, `Layer` [0-13], `Activation` [0,1], state (key-value), `LifecycleState` (state machine with enforced transitions), subscriptions (event type prefixes), `Cadence` [1,∞).

Interface: `Process(tick Tick, events []Event, snapshot Frozen<Snapshot>) → Result<[]Mutation, StoreError>`

The snapshot is deeply immutable (`Frozen<Snapshot>`) — no primitive can mutate another's state. Mutations are declarative values: `AddEvent`, `AddEdge`, `UpdateState`, `UpdateActivation`, `UpdateLifecycle`.

## Tick Engine

The ripple-wave processor. Each tick:
1. Snapshot all primitive states
2. Distribute events to subscribers
3. Invoke primitives (subject to cadence + lifecycle)
4. Collect mutations
5. Apply atomically — new events feed next wave
6. Repeat until quiescence or wave limit (10)
7. Persist states

See `tick-engine.md` for details.

## Decision Trees

Each primitive can have an evolving decision tree. Internal nodes have conditions and branches. Leaf nodes return deterministic outcomes or flag `NeedsLLM = true`.

Over time, patterns are recognised and extracted into deterministic branches, migrating expensive AI calls to cheap rules. See `decision-trees.md`.

## Authority

Three tiers:
- **Required** — blocks until human approves
- **Recommended** — auto-approves after timeout (15 min default)
- **Notification** — auto-approves immediately, logged

Requests and resolutions are events on the graph. See `authority.md`.

## Layer Ontology

14 layers, each derived from a gap in the layer below:

- Layer 0: Foundation (45 primitives, 11 groups)
- Layers 1-13: 12 primitives each (3 groups of 4)

Layer N activates only when Layer N-1 is stable. See `docs/layers/` for per-layer specifications.

## EGIP (Inter-System Protocol)

Sovereign systems communicate without shared infrastructure:
- Self-sovereign identity (Ed25519 keypairs, no registry)
- Cross-Graph Event References (causal links across graph boundaries)
- Signed envelopes
- Treaties for bilateral governance
- Trust accumulation (asymmetric, non-transitive)

See `protocol.md`.

## Weighted Trust and Confidence

Trust and confidence are not binary — they are continuous, contextual, and decaying.

- **Trust** — 0.0 to 1.0, asymmetric, non-transitive. Custom trust mechanics are trivial: override `TrustDecay` with domain logic (+0.05 for completing a task, -0.1 for missing a deadline). The graph handles the rest.
- **Confidence** — Every decision returns epistemic context: not just "permitted" but "permitted with 0.87 confidence through this authority chain with these trust weights."
- **Authority** — Strong here, weak there, contextual. Not a binary gate but a weighted chain.

## Interfaces and Extensibility

The graph defines trait boundaries — the sockets that plugins plug into. The tick engine calls the trait, not the implementation. Swappable, composable, testable.

Core extension points:
- **Store** — plug in any database (memory, SQLite, Postgres, hosted)
- **IDecisionMaker** — plug in anything that makes decisions (AI, human, committee, rules)
- **IIntelligence** — plug in any reasoning engine
- **Primitives** — override defaults, add domain-specific ones

Every component is designed to be replaced. The event graph is the constant. Everything else is pluggable.

## Product Layers (not in this package)

Product layers are built *on* the event graph, not part of it:
- **Social Grammar** — 15 operations (see `grammar.md`)
- **Governance** — Authority chains, voting, delegation
- **Exchange / Market** — Bilateral negotiation, consent, escrow
- **Task Management** — Hierarchical decomposition, model-tier routing (see `docs/task-management.md`)
