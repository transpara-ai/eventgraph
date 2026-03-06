# Roadmap

## How to Use This

Pick any unclaimed task. Create a branch. Submit a PR. See `CONTRIBUTING.md` for standards.

Tasks are ordered by dependency — work higher in the list before lower. Tasks within a section can often be parallelised.

**Status key:** DONE | IN PROGRESS | NEEDED | BLOCKED (by what)

---

## Phase 1: Foundation

The event graph core — the substrate everything else builds on.

### Event Graph Core — NEEDED (extract from mind-zero-five)

The reference implementation exists in [mind-zero-five](https://github.com/mattxo/mind-zero-five). These tasks are primarily extraction, cleanup, and making the code package-friendly.

- [ ] `go/pkg/event/event.go` — Event struct, hash computation, canonical form
- [ ] `go/pkg/event/event_test.go` — Hash chain tests, canonical form tests
- [ ] `go/pkg/store/store.go` — Store interface definition
- [ ] `go/pkg/store/memory.go` — InMemoryStore implementation
- [ ] `go/pkg/store/memory_test.go` — Full conformance tests
- [ ] `go/pkg/store/postgres.go` — PostgresStore implementation
- [ ] `go/pkg/store/postgres_test.go` — Full conformance tests (requires test DB)
- [ ] `go/pkg/store/conformance_test.go` — Shared conformance test suite (any Store impl must pass)
- [ ] `go/pkg/bus/bus.go` — Event bus (pub/sub fan-out)
- [ ] `go/pkg/bus/bus_test.go` — Concurrency tests, backpressure tests
- [ ] `go/pkg/actor/actor.go` — Actor struct, ActorStore interface
- [ ] `go/pkg/actor/actor_test.go`
- [ ] `go/pkg/authority/authority.go` — Request, Policy, three-tier approval
- [ ] `go/pkg/authority/authority_test.go`
- [ ] `go/cmd/eg/main.go` — CLI for interacting with any store

### Primitive Framework — NEEDED

The architecture for primitives — the 200 agents that form the cognitive layers.

- [ ] `go/pkg/primitive/primitive.go` — Primitive interface, Mutation types, Registry
- [ ] `go/pkg/primitive/lifecycle.go` — Lifecycle state machine (dormant→active→processing→emitting)
- [ ] `go/pkg/primitive/registry.go` — Primitive registry (register, get, by-layer, subscribers-for)
- [ ] `go/pkg/primitive/registry_test.go`
- [ ] `go/pkg/primitive/lifecycle_test.go`

### Tick Engine — NEEDED (depends on Primitive Framework)

The ripple-wave processor — the system's heartbeat.

- [ ] `go/pkg/tick/engine.go` — Tick engine, wave processing, quiescence detection
- [ ] `go/pkg/tick/snapshot.go` — Deep copy snapshot mechanism
- [ ] `go/pkg/tick/cadence.go` — Cadence gating logic
- [ ] `go/pkg/tick/engine_test.go` — Ripple tests, wave limit tests, quiescence tests

### Decision Tree Engine — NEEDED

The mechanical-to-intelligent continuum.

- [ ] `go/pkg/decision/tree.go` — Tree structure, internal nodes, leaf nodes, conditions, branches
- [ ] `go/pkg/decision/evaluate.go` — Tree evaluation, path tracking
- [ ] `go/pkg/decision/evolve.go` — Pattern recognition, branch extraction, cost demotion
- [ ] `go/pkg/decision/intelligence.go` — IIntelligence interface, IDecisionMaker interface
- [ ] `go/pkg/decision/tree_test.go`
- [ ] `go/pkg/decision/evaluate_test.go`
- [ ] `go/pkg/decision/evolve_test.go`

---

## Phase 2: Layer 0 Primitives

The 44 foundation primitives in 11 groups. Each primitive needs: implementation, tests, documentation.

### Group 0 — Core
- [ ] Event primitive
- [ ] EventStore primitive
- [ ] Clock primitive
- [ ] Hash primitive
- [ ] Self primitive (identity + routing)

### Group 1 — Causality
- [ ] CausalLink primitive
- [ ] Ancestry primitive
- [ ] Descendancy primitive
- [ ] FirstCause primitive

### Group 2 — Identity
- [ ] ActorID primitive
- [ ] ActorRegistry primitive
- [ ] Signature primitive
- [ ] Verify primitive

### Group 3 — Expectations
- [ ] Expectation primitive
- [ ] Timeout primitive
- [ ] Violation primitive
- [ ] Severity primitive

### Group 4 — Trust
- [ ] TrustScore primitive
- [ ] TrustUpdate primitive
- [ ] Corroboration primitive
- [ ] Contradiction primitive

### Group 5 — Confidence
- [ ] Confidence primitive
- [ ] Evidence primitive
- [ ] Revision primitive
- [ ] Uncertainty primitive

### Group 6 — Instrumentation
- [ ] InstrumentationSpec primitive
- [ ] CoverageCheck primitive
- [ ] Gap primitive
- [ ] Blind primitive

### Group 7 — Query
- [ ] PathQuery primitive
- [ ] SubgraphExtract primitive
- [ ] Annotate primitive
- [ ] Timeline primitive

### Group 8 — Integrity
- [ ] HashChain primitive
- [ ] ChainVerify primitive
- [ ] Witness primitive
- [ ] IntegrityViolation primitive

### Group 9 — Deception
- [ ] Pattern primitive
- [ ] DeceptionIndicator primitive
- [ ] Suspicion primitive
- [ ] Quarantine primitive

### Group 10 — Health
- [ ] GraphHealth primitive
- [ ] Invariant primitive
- [ ] InvariantCheck primitive
- [ ] Bootstrap primitive

---

## Phase 3: Communication Protocol

Inter-primitive communication within a single system.

- [ ] `go/pkg/protocol/message.go` — Four event types: MessageSent, MessageReceived, Decision, Action
- [ ] `go/pkg/protocol/listen_say.go` — Listen/Say interface for communicator primitives
- [ ] `go/pkg/protocol/router.go` — Semantic routing (Self primitive routes to relevant domain primitives)
- [ ] `go/pkg/protocol/knowledge.go` — Three-layer knowledge architecture (context, memory, structural change)
- [ ] Tests for all of the above

---

## Phase 4: EGIP (Inter-System Protocol)

Sovereign systems communicating across graph boundaries.

- [ ] `go/pkg/protocol/egip/identity.go` — Ed25519 keypair, System URI
- [ ] `go/pkg/protocol/egip/cger.go` — Cross-Graph Event Reference
- [ ] `go/pkg/protocol/egip/envelope.go` — Signed message envelope
- [ ] `go/pkg/protocol/egip/messages.go` — Seven message types
- [ ] `go/pkg/protocol/egip/treaty.go` — Treaty model, lifecycle, bilateral governance
- [ ] `go/pkg/protocol/egip/trust.go` — Trust accumulation model
- [ ] `go/pkg/protocol/egip/proof.go` — Integrity proofs (chain segment, event existence, chain summary)
- [ ] Tests for all of the above

---

## Phase 5: Layers 1-13

Each layer has 12 primitives in 3 groups of 4. Layer N depends on Layer N-1 being stable.

### Layer 1 — Agency (Observer → Participant)
- [ ] Specification (`docs/layers/01-agency.md`)
- [ ] 12 primitives implementation + tests

### Layer 2 — Exchange (Individual → Dyad)
- [ ] Specification (`docs/layers/02-exchange.md`)
- [ ] 12 primitives implementation + tests

### Layer 3 — Society (Dyad → Group)
- [ ] Specification
- [ ] 12 primitives implementation + tests

### Layer 4 — Legal (Informal → Formal)
- [ ] Specification
- [ ] 12 primitives implementation + tests

### Layer 5 — Technology (Governing → Building)
- [ ] Specification
- [ ] 12 primitives implementation + tests

### Layer 6 — Information (Physical → Symbolic)
- [ ] Specification
- [ ] 12 primitives implementation + tests

### Layer 7 — Ethics (Is → Ought)
- [ ] Specification
- [ ] 12 primitives implementation + tests

### Layer 8 — Identity (Doing → Being)
- [ ] Specification
- [ ] 12 primitives implementation + tests

### Layer 9 — Relationship (Self → Self-with-Other)
- [ ] Specification
- [ ] 12 primitives implementation + tests

### Layer 10 — Community (Relationship → Belonging)
- [ ] Specification
- [ ] 12 primitives implementation + tests

### Layer 11 — Culture (Living → Seeing)
- [ ] Specification
- [ ] 12 primitives implementation + tests

### Layer 12 — Emergence (Content → Architecture)
- [ ] Specification
- [ ] 12 primitives implementation + tests

### Layer 13 — Existence (Everything → The Fact of Everything)
- [ ] Specification
- [ ] 12 primitives implementation + tests

---

## Phase 6: Language Packages

Each language package must pass the language-agnostic conformance test suite.

### Rust
- [ ] Core event types + hash chain
- [ ] Store trait + InMemory implementation
- [ ] Bus
- [ ] Primitive trait + Registry
- [ ] Tick engine
- [ ] Conformance tests passing

### Python
- [ ] Core event types + hash chain
- [ ] Store protocol + InMemory implementation
- [ ] Bus
- [ ] Primitive protocol + Registry
- [ ] Tick engine
- [ ] Conformance tests passing

### .NET
- [ ] Core event types + hash chain
- [ ] IStore interface + InMemory implementation
- [ ] Bus
- [ ] IPrimitive interface + Registry
- [ ] Tick engine
- [ ] Conformance tests passing

---

## Phase 7: Documentation & Examples

- [ ] `docs/conformance/` — Language-agnostic conformance test specification
- [ ] `examples/minimal/` — Smallest possible event graph (10 lines of code)
- [ ] `examples/social/` — The 15 social grammar operations on the event graph
- [ ] `examples/multi-system/` — Two systems communicating via EGIP
- [ ] Tutorial: "Build your first primitive"
- [ ] Tutorial: "Implement a custom store"
- [ ] Tutorial: "Connect two event graphs"

---

## Future

These are on the horizon but not yet specified:

- [ ] Product layer: Social Grammar (15 operations from Post 35)
- [ ] Product layer: Governance (Post 34)
- [ ] Product layer: Exchange / Market
- [ ] WebAssembly builds for browser-based event graphs
- [ ] Mobile SDKs
- [ ] Reference UI implementations
