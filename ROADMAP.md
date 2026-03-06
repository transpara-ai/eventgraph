# Roadmap

## How to Use This

Pick any unclaimed task. Create a branch. Submit a PR. See `CONTRIBUTING.md` for standards.

Tasks are ordered by dependency — work higher in the list before lower. Tasks within a section can often be parallelised.

**For automated implementers:** See `docs/implementation-order.md` for the strict dependency DAG with compilation-order tasks, acceptance criteria, and the implementer loop.

**Status key:** DONE | IN PROGRESS | NEEDED | BLOCKED (by what)

---

## Phase 1: Foundation

The event graph core — the substrate everything else builds on.

### Event Graph Core — NEEDED

The reference implementation exists in [mind-zero-five](https://github.com/mattxo/mind-zero-five). These tasks are primarily extraction, cleanup, and making the code package-friendly.

**See `docs/implementation-order.md` for the strict dependency-ordered implementation DAG with acceptance criteria.**

#### Tier 0-1: Foundation Types
- [ ] `go/pkg/types/option.go` — `Option[T]` generic type (Some, None, Unwrap, JSON)
- [ ] `go/pkg/types/nonempty.go` — `NonEmpty[T]` generic type (rejects empty)
- [ ] `go/pkg/types/page.go` — `Page[T]` pagination, `Cursor`
- [ ] `go/pkg/types/errors.go` — All `ValidationError` types
- [ ] `go/pkg/types/ids.go` — Value objects: EventID, ActorID, Hash, ConversationID, SystemURI, PublicKey, Signature, etc.
- [ ] `go/pkg/types/constrained.go` — Constrained numerics: Score [0,1], Weight [-1,1], Activation [0,1], Layer [0,13], Cadence [1,∞), Tick [0,∞)
- [ ] `go/pkg/types/statemachine.go` — LifecycleState and ActorStatus state machines (enforced valid transitions)
- [ ] `go/pkg/types/types_test.go` — Construction validation, rejection of invalid values, equality, state transitions, conformance vectors

#### Tier 2-3: Events and Content
- [ ] `go/pkg/event/constants.go` — All enums with `IsValid()` and Visitor interfaces
- [ ] `go/pkg/event/content.go` — EventContent interface, all content structs, EventTypeRegistry, EventContentVisitor
- [ ] `go/pkg/event/edge.go` — Edge struct, EdgeMetadata interface, EdgeTypeRegistry, EdgeMetadataVisitor, all metadata types
- [ ] `go/pkg/event/event.go` — Event struct (immutable), canonical form, hash computation
- [ ] `go/pkg/event/decision.go` — Decision, DecisionInput, Receipt, TrustMetrics, AuthorityLink, TrustWeight, Expectation, ViolationRecord
- [ ] `go/pkg/event/event_test.go` — Canonical form vectors, hash chain tests, content validation

#### Tier 4-5: Store and Actor
- [ ] `go/pkg/store/errors.go` — All `StoreError` types, StoreErrorVisitor
- [ ] `go/pkg/store/store.go` — Store interface definition
- [ ] `go/pkg/store/memory.go` — InMemoryStore implementation (chain locking, edge indexing, concurrent-safe)
- [ ] `go/pkg/store/conformance_test.go` — Shared conformance test suite
- [ ] `go/pkg/store/memory_test.go` — Runs conformance suite + memory-specific tests
- [ ] `go/pkg/store/postgres.go` — PostgresStore implementation
- [ ] `go/pkg/store/postgres_test.go` — Runs conformance suite (requires test DB)
- [ ] `go/pkg/actor/actor.go` — IActor, Actor, IActorStore, ActorUpdate, ActorFilter
- [ ] `go/pkg/actor/memory.go` — InMemoryActorStore
- [ ] `go/pkg/actor/actor_test.go` — Registration, lookup, lifecycle, pagination
- [ ] `go/pkg/event/factory.go` — EventFactory, BootstrapFactory, EdgeFactory

#### Tier 6-7: Bus, Decision, Trust, Authority
- [ ] `go/pkg/bus/bus.go` — IBus, EventBus (non-blocking, overflow handling)
- [ ] `go/pkg/bus/bus_test.go` — Concurrency, backpressure, overflow
- [ ] `go/pkg/decision/tree.go` — DecisionTree, nodes, conditions, stats
- [ ] `go/pkg/decision/evaluate.go` — Tree evaluation, path tracking, Semantic conditions
- [ ] `go/pkg/decision/intelligence.go` — IIntelligence, IDecisionMaker interfaces
- [ ] `go/pkg/trust/model.go` — ITrustModel, DefaultTrustModel (decay, recovery)
- [ ] `go/pkg/trust/model_test.go` — Scoring, decay, domain-specific, boundary values
- [ ] `go/pkg/authority/authority.go` — IAuthorityChain, AuthorityResult, policies
- [ ] `go/pkg/authority/chain.go` — Delegation chain walk, weight propagation
- [ ] `go/pkg/authority/authority_test.go` — Evaluation, delegation, trust-based demotion

#### Tier 8-9: Primitives and Tick Engine
- [ ] `go/pkg/primitive/primitive.go` — Primitive interface, Mutation types, MutationVisitor
- [ ] `go/pkg/primitive/registry.go` — PrimitiveRegistry
- [ ] `go/pkg/primitive/lifecycle.go` — Lifecycle integration with tick
- [ ] `go/pkg/primitive/harness.go` — PrimitiveTestHarness
- [ ] `go/pkg/primitive/registry_test.go`
- [ ] `go/pkg/primitive/lifecycle_test.go`
- [ ] `go/pkg/tick/snapshot.go` — FrozenSnapshot, PrimitiveState (deep copy, immutable)
- [ ] `go/pkg/tick/cadence.go` — Cadence gating logic
- [ ] `go/pkg/tick/engine.go` — Tick engine, wave processing, quiescence, layer ordering
- [ ] `go/pkg/tick/engine_test.go` — Ripple, wave limit, quiescence, layer constraint, concurrency

#### Tier 10: Top-Level API
- [ ] `go/pkg/graph/graph.go` — IGraph (Evaluate, Record, Query), IGraphQuery, GraphConfig
- [ ] `go/pkg/graph/graph_test.go` — End-to-end integration tests
- [ ] `go/cmd/eg/main.go` — CLI for interacting with any store

### Primitive Framework — NEEDED

The architecture for primitives — the 201 agents that form the cognitive layers.

- [ ] `go/pkg/primitive/primitive.go` — Primitive interface, Mutation types, Registry
- [ ] `go/pkg/primitive/lifecycle.go` — Lifecycle state machine (dormant→active→processing→emitting)
- [ ] `go/pkg/primitive/registry.go` — Primitive registry (register, get, by-layer, subscribers-for)
- [ ] `go/pkg/primitive/registry_test.go`
- [ ] `go/pkg/primitive/lifecycle_test.go`

### Tick Engine — NEEDED (depends on Primitive Framework)

The ripple-wave processor — the system's heartbeat.

- [ ] `go/pkg/tick/engine.go` — Tick engine, wave processing, quiescence detection
- [ ] `go/pkg/tick/snapshot.go` — Snapshot, PrimitiveState, Frozen<Snapshot> deep copy mechanism
- [ ] `go/pkg/tick/cadence.go` — Cadence gating logic
- [ ] `go/pkg/tick/engine_test.go` — Ripple tests, wave limit tests, quiescence tests

### Decision Tree Engine — NEEDED

The mechanical-to-intelligent continuum.

- [ ] `go/pkg/decision/tree.go` — Tree structure, internal nodes, leaf nodes, conditions, branches
- [ ] `go/pkg/decision/evaluate.go` — Tree evaluation, path tracking
- [ ] `go/pkg/decision/evolve.go` — Pattern recognition, branch extraction, cost demotion
- [ ] `go/pkg/decision/intelligence.go` — IIntelligence interface, IDecisionMaker interface, DecisionInput/Decision types
- [ ] `go/pkg/trust/model.go` — ITrustModel interface, default implementation
- [ ] `go/pkg/trust/model_test.go` — Trust scoring, decay, domain-specific trust tests
- [ ] `go/pkg/authority/chain.go` — IAuthorityChain interface, default implementation
- [ ] `go/pkg/authority/chain_test.go` — Authority evaluation, delegation, weighted chains
- [ ] `go/pkg/decision/tree_test.go`
- [ ] `go/pkg/decision/evaluate_test.go`
- [ ] `go/pkg/decision/evolve_test.go`

---

## Phase 2: Layer 0 Primitives

The 45 foundation primitives in 11 groups. Each primitive needs: implementation, tests, documentation.

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
- [ ] `docs/conformance/canonical-vectors.json` — Test vectors for canonical form, hash, and signature verification
- [ ] `examples/minimal/` — Smallest possible event graph (10 lines of code)
- [ ] `examples/social/` — The 15 social grammar operations on the event graph
- [ ] `examples/multi-system/` — Two systems communicating via EGIP
- [ ] Tutorial: "Build your first primitive"
- [ ] Tutorial: "Implement a custom store"
- [ ] Tutorial: "Connect two event graphs"

---

## Future

These are product layers — built *on* the event graph, not part of the infrastructure package:

- [ ] Product layer: Social Grammar (15 operations from Post 35)
- [ ] Product layer: Governance (Post 34)
- [ ] Product layer: Exchange / Market
- [ ] Product layer: Task Management (hierarchical decomposition, model-tier routing)
- [ ] Top-level API: `graph.Evaluate()` / `graph.Record()` / `graph.Query()`
- [ ] WebAssembly builds for browser-based event graphs
- [ ] Mobile SDKs
- [ ] Reference UI implementations
- [ ] Hosted persistence service
