# Architecture

## Overview

EventGraph is a hash-chained, append-only, causal event graph with cognitive primitives.

```
┌─────────────────────────────────────────────┐
│                 Interfaces                   │
│         (social, governance, market)         │
├─────────────────────────────────────────────┤
│            Product Layers                    │
│    (Social Grammar, Governance, Exchange)    │
├─────────────────────────────────────────────┤
│         Cognitive Primitives                 │
│    (200 across 14 layers, tick engine)       │
├─────────────────────────────────────────────┤
│           Communication                      │
│    (listen/say, routing, EGIP)              │
├─────────────────────────────────────────────┤
│            Event Graph                       │
│    (events, hash chain, causal DAG, bus)    │
├─────────────────────────────────────────────┤
│              Store                           │
│    (Postgres, SQLite, Memory, your own)     │
└─────────────────────────────────────────────┘
```

## The Event

The fundamental unit. Every significant action is an event.

```
Event {
    ID:             UUID v7 (time-ordered)
    Type:           hierarchical string (e.g., "trust.updated")
    Timestamp:      nanosecond precision
    Source:         who/what emitted this
    Content:        structured payload (JSON)
    Causes:         [event IDs] — causal DAG
    ConversationID: thread grouping
    Hash:           SHA-256 of canonical form
    PrevHash:       hash of previous event — linear chain
    Signature:      Ed25519 signature
}
```

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

```
Store {
    Append(type, source, content, causes, conversationID) → Event
    Get(id) → Event
    Recent(limit) → []Event
    ByType(type, limit) → []Event
    BySource(source, limit) → []Event
    ByConversation(conversationID, limit) → []Event
    Since(afterID, limit) → []Event
    Search(query, limit) → []Event
    Ancestors(id, maxDepth) → []Event
    Descendants(id, maxDepth) → []Event
    Count() → int
    VerifyChain() → error
}
```

Implement this interface for any backing store. Reference implementations: PostgresStore, InMemoryStore.

## Bus

Wraps a Store with pub/sub fan-out. When an event is appended, all subscribers receive it. Non-blocking — slow subscribers get dropped events, not blocked writers.

## Primitives

A cognitive primitive is a software agent for a specific domain. See `primitives.md` for the full specification.

Key properties: name, layer, activation level, state (key-value), lifecycle state machine, subscriptions (event type prefixes), cadence (minimum ticks between invocations).

Interface: `Process(tick, events, snapshot) → []Mutation`

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

- Layer 0: Foundation (44 primitives, 11 groups)
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

## Interfaces and Extensibility

Core extension points:
- **Store** — plug in any database
- **IIntelligence** — plug in any reasoning engine
- **IDecisionMaker** — plug in any decision routing
- **Primitives** — override defaults, add new ones

Every component is designed to be replaced. The event graph is the constant. Everything else is pluggable.
