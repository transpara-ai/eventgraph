# EventGraph

> **Status: Pre-release.** Documentation and specifications are in place. Core packages are being implemented. Contributions welcome — see `ROADMAP.md`.

A hash-chained, append-only, causal event graph. The foundation for building systems where every action is signed, auditable, and causally linked.

## What This Is

EventGraph is infrastructure, not a product. A standard with packages. Published to every major ecosystem — Go, Rust, Python, .NET — so developers can make any system auditable in four lines of code.

The graph doesn't take an "agent" as input. It takes an `IDecisionMaker` — anything that makes decisions. An AI agent implements it. A human with a UI implements it. A committee vote implements it. A rules engine implements it. The graph doesn't know or care what's deciding. It records *what* was decided, *by what*, with *what confidence*, under *what authority*, and links it causally to everything that led there.

This makes EventGraph **decision governance**, not AI governance — applicable to any system where decisions matter and accountability is required.

## Soul Statement

> Take care of your human, humanity, and yourself. In that order when they conflict, but they rarely should.

Every design decision serves this.

## Two Levels of API

**Top level** — what most developers hit. Four lines to make any system auditable:

```go
graph := eventgraph.New(store.Memory())     // or Postgres, SQLite, your own
result, _ := graph.Evaluate(ctx, action)     // evaluate through the ontology
receipt := result.Receipt()                  // cryptographic proof of decision
graph.Record(ctx, result)                    // append to the hash chain
```

**Power user level** — for building AI agent frameworks, compliance platforms, or custom governance. Every primitive is an interface with sensible defaults. Override what you need:

```go
// Custom trust decay for your domain
type MyTrustDecay struct { primitives.TrustDecay }
func (t *MyTrustDecay) Process(tick, events, snapshot) []Mutation {
    // domain-specific trust logic
}
graph.Register(&MyTrustDecay{})
```

**Persistence is a plugin:**

```go
graph := eventgraph.New(store.Memory())      // dev/testing
graph := eventgraph.New(store.SQLite(path))   // single app
graph := eventgraph.New(store.Postgres(dsn))  // production
```

> **Note:** These packages are not yet published. See `ROADMAP.md` Phase 1 for implementation status. The API shown above reflects the target design.

## Quick Start (Go)

```bash
go get github.com/lovyou-ai/eventgraph/go
```

```go
package main

import (
    "context"
    "fmt"

    eg "github.com/lovyou-ai/eventgraph/go/pkg/graph"
    "github.com/lovyou-ai/eventgraph/go/pkg/store"
)

func main() {
    ctx := context.Background()

    // Create a graph with an in-memory store
    graph, _ := eg.New(store.NewMemory())
    graph.Start()
    defer graph.Close()

    // Record an event (factory validates, hashes, signs, then store persists)
    e, _ := graph.Record(ctx, "hello.world", systemActor, HelloContent{
        Message: "first event",
    }, causes, conversationID)

    // Every event is hash-chained and signed
    fmt.Println(e.Hash)      // SHA-256 of canonical form
    fmt.Println(e.PrevHash)  // links to previous event
    fmt.Println(e.Signature) // Ed25519 signature

    // Causal traversal
    ancestors, _ := graph.Query().Events().Ancestors(e.ID, 10)
    fmt.Println(len(ancestors))
}
```

## Architecture

```
Top-Level API (graph.Evaluate / Record / Query)
    ↕
Product Layers (social, governance, exchange — built on, not in)
    ↕
Primitives (200 across 14 layers, each an overridable interface)
    ↕
Tick Engine (ripple-wave processing until quiescence)
    ↕
Decision Trees (mechanical-to-intelligent continuum, evolving)
    ↕
Authority (three-tier: required / recommended / notification)
    ↕
Event Graph (hash-chained, append-only, causal DAG)
    ↕
Store (plugin: Memory, SQLite, Postgres, your own)
```

See `docs/architecture.md` for the full picture.

## The 14 Layers

| Layer | Name | What it adds |
|-------|------|-------------|
| 0 | Foundation | 44 primitives: events, hash chains, identity, trust, causality |
| 1 | Agency | Observer becomes participant — action and communication |
| 2 | Exchange | Individual becomes dyad — negotiated interaction |
| 3 | Society | Dyad becomes group — collective behaviour |
| 4 | Legal | Informal becomes formal — binding agreements |
| 5 | Technology | Governing becomes building — tool creation |
| 6 | Information | Physical becomes symbolic — representation |
| 7 | Ethics | Is becomes ought — moral reasoning |
| 8 | Identity | Doing becomes being — self-knowledge |
| 9 | Relationship | Self becomes self-with-other |
| 10 | Community | Relationship becomes belonging |
| 11 | Culture | Living culture becomes seeing culture |
| 12 | Emergence | Content becomes architecture — self-organisation |
| 13 | Existence | Everything becomes the fact of everything |

Each layer is derived from a gap in the layer below — something the lower layer cannot express but structurally needs. See `docs/layers/` for per-layer specifications.

## Key Interfaces

Everything is pluggable. The graph defines the sockets; you provide the implementations:

- **`Store`** — Event persistence. Memory, SQLite, Postgres, or your own.
- **`IDecisionMaker`** — Anything that makes decisions. AI agents, humans, committees, rules engines. The graph records what was decided, not how.
- **`IIntelligence`** — Reasoning. Any model, local or cloud, or deterministic logic.
- **Primitives** — 200 cognitive agents, each an interface with sensible defaults. Override with domain-specific logic.

Every decision returns epistemic context: not just "permitted" but "permitted with 0.87 confidence through this authority chain with these trust weights." Trust isn't binary — it's 0.73 and decaying. Authority isn't binary — it's strong here, weak there, contextual.

See `docs/interfaces.md` for specifications.

## Contributing

We welcome contributions from humans and AI systems. If you have Claude Code or any AI coding tool, read `CLAUDE.md` — it's designed to let you pick up work immediately.

See `CONTRIBUTING.md` for process. See `ROADMAP.md` for what needs building.

## Language Support

Published to every ecosystem developers already work in:

| Language | Package | Status | Path |
|----------|---------|--------|------|
| Go | `go get github.com/lovyou-ai/eventgraph/go` | Reference implementation | `go/` |
| Rust | `cargo add eventgraph` | Community — help wanted | `rust/` |
| Python | `pip install eventgraph` | Community — help wanted | `python/` |
| .NET | `dotnet add package EventGraph` | Community — help wanted | `dotnet/` |

All implementations must pass the language-agnostic conformance test suite. Each implements the same interfaces — `Store`, `IDecisionMaker`, `IIntelligence`, `Primitive`. Same spec, native to each ecosystem.

## License

[Business Source License 1.1](LICENSE) converting to [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0) on 26 February 2030.

The packages are free. The spec is free. The ideas are free (deed poll). Run it yourself and owe nothing. Non-commercial use, research, and education are always free. Free and open source implementations are always protected under the [Defensive Patent Pledge](patent/defensive_patent_pledge.pdf).

## Links

- [The blog series](https://mattsearles2.substack.com) — Posts 33-35 derive the values, governance, and social grammar
- [mind-zero](https://github.com/mattxo/mind-zero) — The primitive derivation
- [mind-zero-five](https://github.com/mattxo/mind-zero-five) — The working implementation

## Patent

Australian Provisional Patent Application No. 2026901564. Subject to an irrevocable [Defensive Patent Pledge](patent/defensive_patent_pledge.pdf) — free and open source implementations are protected, always. The patent exists as a shield, not a weapon. See `PATENT-NOTICE` for details.
