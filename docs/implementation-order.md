# Implementation Order

Strict dependency-ordered task list for the Go reference implementation. An automated implementer reads this to know what's unblocked.

## How to Read This

Each task has:
- **File** — what to create
- **Depends on** — tasks that must be complete (tests passing) before this one starts
- **Spec** — which doc sections to read before implementing
- **Acceptance** — how to know it's done

Tasks at the same tier can be parallelised. Tasks at tier N depend on all lower tiers being complete.

---

## Tier 0: Generic Foundation Types

No dependencies. These are leaf packages that everything else imports.

### T0.1 — `go/pkg/types/option.go`

| | |
|---|---|
| **File** | `go/pkg/types/option.go` |
| **Depends on** | nothing |
| **Spec** | `docs/coding-standards/go.md` → Option\<T\> section |
| **Acceptance** | `Option[T]` generic type with `Some`, `None`, `IsSome`, `IsNone`, `Unwrap`, `UnwrapOr`. JSON marshal/unmarshal. Tests in `option_test.go`. |

### T0.2 — `go/pkg/types/nonempty.go`

| | |
|---|---|
| **File** | `go/pkg/types/nonempty.go` |
| **Depends on** | nothing |
| **Spec** | `docs/coding-standards/go.md` → NonEmpty\<T\> section |
| **Acceptance** | `NonEmpty[T]` generic type. Constructor rejects empty. `First`, `All`, `Len`. Tests: empty rejection, single element, multiple elements. |

### T0.3 — `go/pkg/types/page.go`

| | |
|---|---|
| **File** | `go/pkg/types/page.go` |
| **Depends on** | T0.1 (Option) |
| **Spec** | `docs/coding-standards/go.md` → Page\<T\> section, `docs/interfaces.md` → Page\<T\> |
| **Acceptance** | `Page[T]` struct, `Cursor` type. Tests: empty page, page with cursor, page without cursor. |

### T0.4 — `go/pkg/types/errors.go`

| | |
|---|---|
| **File** | `go/pkg/types/errors.go` |
| **Depends on** | nothing |
| **Spec** | `docs/interfaces.md` → Domain Errors, `docs/coding-standards/go.md` → Error types |
| **Acceptance** | All `ValidationError` types: `OutOfRangeError`, `InvalidFormatError`, `EmptyRequiredError`, `InvalidLifecycleTransitionError`, `InvalidActorTransitionError`. Each implements `error` and `ValidationError` marker interface. Tests: error messages, `errors.As` dispatch. |

---

## Tier 1: Value Objects

Depends on: Tier 0 (errors for construction validation).

### T1.1 — `go/pkg/types/constrained.go`

| | |
|---|---|
| **File** | `go/pkg/types/constrained.go` |
| **Depends on** | T0.4 (errors) |
| **Spec** | `docs/interfaces.md` → Constrained Numerics |
| **Acceptance** | `Score` [0,1], `Weight` [-1,1], `Activation` [0,1], `Layer` [0,13], `Cadence` [1,∞), `Tick` [0,∞). Each: `New*`, `Must*`, `Value()`, JSON marshal/unmarshal. Tests: valid construction, boundary values, rejection of invalid values, JSON round-trip. |

### T1.2 — `go/pkg/types/ids.go`

| | |
|---|---|
| **File** | `go/pkg/types/ids.go` |
| **Depends on** | T0.4 (errors) |
| **Spec** | `docs/interfaces.md` → Typed IDs |
| **Acceptance** | `EventID` (UUID v7 validated), `ActorID`, `EdgeID`, `ConversationID`, `Hash` (64 hex chars), `SystemURI`, `EnvelopeID`, `TreatyID`, `PrimitiveID`, `EventType` (validated against registry — initially just format validation, registry integration later), `SubscriptionPattern`, `DomainScope`, `PublicKey` (32 bytes), `Signature` (64 bytes). Each: constructor with validation, `String()`, JSON marshal/unmarshal, comparable. Tests: valid/invalid construction, equality, format validation. |

### T1.3 — `go/pkg/types/statemachine.go`

| | |
|---|---|
| **File** | `go/pkg/types/statemachine.go` |
| **Depends on** | T0.4 (errors) |
| **Spec** | `docs/interfaces.md` → Lifecycle State Machine, Actor Status Transitions |
| **Acceptance** | `LifecycleState` with `TransitionTo()` and `ValidTransitions()`. `ActorStatus` same. Tests: all valid transitions succeed, all invalid transitions return typed error, terminal state (Memorial) has no valid transitions. Use conformance vectors from `docs/conformance/canonical-vectors.json`. |

### T1.4 — `go/pkg/types/types_test.go`

| | |
|---|---|
| **File** | `go/pkg/types/types_test.go` |
| **Depends on** | T1.1, T1.2, T1.3 |
| **Spec** | `docs/conformance/canonical-vectors.json` → type_validation, lifecycle_transitions |
| **Acceptance** | Integration tests covering: cross-type interactions, conformance vector validation, all boundary cases. Coverage ≥ 90%. |

---

## Tier 2: Constants and Event Content

Depends on: Tier 1 (typed IDs, constrained numerics).

### T2.1 — `go/pkg/event/constants.go`

| | |
|---|---|
| **File** | `go/pkg/event/constants.go` |
| **Depends on** | T1.1, T1.2 |
| **Spec** | `docs/interfaces.md` → Constants and Enums |
| **Acceptance** | All enums: `EdgeType`, `AuthorityLevel`, `DecisionOutcome`, `ActorType`, `EdgeDirection`, `ExpectationStatus`, `SeverityLevel`, `MessageType`, `TreatyStatus`, `IntegrityViolationType`, `InvariantName`, `CGERRelationship`, `ReceiptStatus`, `ProofType`, `TreatyAction`, `ConditionOperator`. Each with `IsValid()`. Visitor interfaces for `DecisionOutcome`, `AuthorityLevel`. Tests: all variants valid, unknown variants invalid. |

### T2.2 — `go/pkg/event/content.go`

| | |
|---|---|
| **File** | `go/pkg/event/content.go` |
| **Depends on** | T1.1, T1.2, T2.1 |
| **Spec** | `docs/interfaces.md` → Event Type Registry, Typed Event Content |
| **Acceptance** | `EventContent` interface, `EventContentVisitor` interface, `EventTypeRegistry` (Register, Validate, Schema, AllTypes). All content structs: `TrustUpdatedContent`, `TrustScoreContent`, `TrustDecayedContent`, `AuthorityRequestContent`, `AuthorityResolvedContent`, `AuthorityDelegatedContent`, `AuthorityRevokedContent`, `AuthorityTimeoutContent`, `ActorRegisteredContent`, `ActorSuspendedContent`, `ActorMemorialContent`, `EdgeCreatedContent`, `EdgeSupersededContent`, `ViolationDetectedContent`, `ChainVerifiedContent`, `ChainBrokenContent`, `BootstrapContent`, `ClockTickContent`, `HealthReportContent`, `BranchProposedContent`, `BranchInsertedContent`, `CostReportContent`, all EGIP content types. Tests: registry validation, visitor dispatch, JSON serialisation. |

### T2.3 — `go/pkg/event/edge.go`

| | |
|---|---|
| **File** | `go/pkg/event/edge.go` |
| **Depends on** | T1.1, T1.2, T2.1 |
| **Spec** | `docs/interfaces.md` → Edge, Edge Type Registry, Typed Edge Metadata |
| **Acceptance** | `Edge` struct (immutable, unexported fields + getters), `EdgeMetadata` interface, `EdgeMetadataVisitor`, `EdgeTypeRegistry`, all 9 concrete metadata types. Tests: construction validation, visitor dispatch, JSON serialisation. |

---

## Tier 3: Event and Decision Structures

Depends on: Tier 2 (content, constants, edge).

### T3.1 — `go/pkg/event/event.go`

| | |
|---|---|
| **File** | `go/pkg/event/event.go` |
| **Depends on** | T1.1, T1.2, T2.2 |
| **Spec** | `docs/interfaces.md` → Event, Canonical Form |
| **Acceptance** | `Event` struct (immutable). `CanonicalForm(event) → string` matching the spec exactly. `ComputeHash(canonical) → Hash` using SHA-256. Tests: canonical form matches conformance vectors, hash computation, JSON round-trip. |

### T3.2 — `go/pkg/event/decision.go`

| | |
|---|---|
| **File** | `go/pkg/event/decision.go` |
| **Depends on** | T1.1, T1.2, T2.1 |
| **Spec** | `docs/interfaces.md` → Decision, DecisionInput, Receipt, TrustMetrics, AuthorityLink, TrustWeight, Expectation, ViolationRecord |
| **Acceptance** | All structs with constrained fields, immutable. Tests: construction with valid/invalid fields, constrained field rejection. |

### T3.3 — `go/pkg/event/event_test.go`

| | |
|---|---|
| **File** | `go/pkg/event/event_test.go` |
| **Depends on** | T3.1, T3.2, T2.2, T2.3 |
| **Spec** | `docs/conformance/canonical-vectors.json` |
| **Acceptance** | Canonical form vector tests, hash chain tests, content validation tests. Loads and runs all canonical-vectors.json test cases. Coverage ≥ 90%. |

---

## Tier 4: Store Errors and Interface

Depends on: Tier 3 (Event struct for Store interface signature).

### T4.1 — `go/pkg/store/errors.go`

| | |
|---|---|
| **File** | `go/pkg/store/errors.go` |
| **Depends on** | T1.2 (typed IDs) |
| **Spec** | `docs/interfaces.md` → Domain Errors (StoreError) |
| **Acceptance** | All `StoreError` types: `EventNotFoundError`, `ActorNotFoundError`, `EdgeNotFoundError`, `DuplicateEventError`, `CausalLinkMissingError`, `ChainIntegrityViolationError`, `HashMismatchError`, `SignatureInvalidError`, `ActorSuspendedError`, `ActorMemorialError`, `RateLimitExceededError`, `StoreUnavailableError`. `StoreError` marker interface. `StoreErrorVisitor`. Tests: all error types, `errors.As` dispatch, visitor exhaustiveness. |

### T4.2 — `go/pkg/store/store.go`

| | |
|---|---|
| **File** | `go/pkg/store/store.go` |
| **Depends on** | T3.1 (Event), T2.3 (Edge), T1.2 (IDs), T0.1 (Option), T0.3 (Page), T4.1 (errors) |
| **Spec** | `docs/interfaces.md` → Store interface |
| **Acceptance** | `Store` interface with all methods. Just the interface definition — no implementation yet. |

---

## Tier 5: Actor and Store Implementations

Depends on: Tier 4 (Store interface).

### T5.1 — `go/pkg/actor/actor.go`

| | |
|---|---|
| **File** | `go/pkg/actor/actor.go` |
| **Depends on** | T1.2, T1.3, T0.1, T0.3, T4.1 |
| **Spec** | `docs/interfaces.md` → IActor, IActorStore, ActorUpdate, ActorFilter |
| **Acceptance** | `IActor` interface, `Actor` struct (immutable), `IActorStore` interface, `ActorUpdate`, `ActorFilter`. Tests: actor construction, status transitions. |

### T5.2 — `go/pkg/actor/memory.go`

| | |
|---|---|
| **File** | `go/pkg/actor/memory.go` |
| **Depends on** | T5.1 |
| **Spec** | `docs/interfaces.md` → IActorStore |
| **Acceptance** | `InMemoryActorStore` implementing `IActorStore`. Tests: register, get, get-by-key, update, list with filters, suspend, memorial, pagination, concurrent access. |

### T5.3 — `go/pkg/event/factory.go`

| | |
|---|---|
| **File** | `go/pkg/event/factory.go` |
| **Depends on** | T3.1 (Event), T4.2 (Store interface), T5.1 (IActor), T2.2 (content/registry) |
| **Spec** | `docs/interfaces.md` → Factory (EventFactory, BootstrapFactory, EdgeFactory) |
| **Acceptance** | `EventFactory.Create()`, `BootstrapFactory.Init()`, `EdgeFactory.Create()`. Each validates all inputs, computes hash, signs. Tests: valid creation, each validation failure case, bootstrap has no causes, factory is the only path to create Events. |

### T5.4 — `go/pkg/store/memory.go`

| | |
|---|---|
| **File** | `go/pkg/store/memory.go` |
| **Depends on** | T4.2, T5.3 (factory — needed for test helpers) |
| **Spec** | `docs/interfaces.md` → Store, `docs/architecture.md` → Hash Chain |
| **Acceptance** | `InMemoryStore` implementing `Store`. Chain head locking. Edge indexing. Concurrent-safe. All query methods. |

### T5.5 — `go/pkg/store/conformance_test.go`

| | |
|---|---|
| **File** | `go/pkg/store/conformance_test.go` |
| **Depends on** | T5.4 |
| **Spec** | `docs/interfaces.md` → Testing Infrastructure (Conformance), `docs/conformance/canonical-vectors.json` |
| **Acceptance** | `RunConformanceSuite(t, factory)`. Tests: append/retrieve, hash chain integrity, causal traversal, edge indexing, pagination, concurrency, idempotency, chain head conflict, canonical form vectors. `InMemoryStore` passes the full suite. |

### T5.6 — `go/pkg/store/memory_test.go`

| | |
|---|---|
| **File** | `go/pkg/store/memory_test.go` |
| **Depends on** | T5.4, T5.5 |
| **Spec** | — |
| **Acceptance** | Runs `RunConformanceSuite` with `InMemoryStore` factory. Additional memory-specific tests if needed. Coverage ≥ 90%. |

---

## Tier 6: Bus and Decision Infrastructure

Depends on: Tier 5 (Store implementation, factory).

### T6.1 — `go/pkg/bus/bus.go`

| | |
|---|---|
| **File** | `go/pkg/bus/bus.go` |
| **Depends on** | T4.2 (Store), T1.2 (SubscriptionPattern) |
| **Spec** | `docs/interfaces.md` → Bus |
| **Acceptance** | `IBus` interface, `EventBus` implementation wrapping Store. Subscribe, Unsubscribe, pattern matching, non-blocking delivery, overflow handling (`bus.overflow` event). Tests: subscription matching, fan-out, slow subscriber doesn't block, overflow detection, concurrent subscribe/unsubscribe. Coverage ≥ 90%. |

### T6.2 — `go/pkg/decision/tree.go`

| | |
|---|---|
| **File** | `go/pkg/decision/tree.go` |
| **Depends on** | T1.1, T1.2, T2.1 |
| **Spec** | `docs/interfaces.md` → Decision Tree, `docs/decision-trees.md` → Tree Structure |
| **Acceptance** | `DecisionTree`, `DecisionNode` (interface), `InternalNode`, `LeafNode`, `Branch`, `Condition`, `MatchValue`, `PathStep`, `TreeStats`, `LeafStats`, `ResponseRecord`, `SemanticEvalRecord`. Tests: tree construction, node types. |

### T6.3 — `go/pkg/decision/evaluate.go`

| | |
|---|---|
| **File** | `go/pkg/decision/evaluate.go` |
| **Depends on** | T6.2 |
| **Spec** | `docs/decision-trees.md` → Evaluation |
| **Acceptance** | `Evaluate(tree, input, intelligence) → Decision`. Path tracking. Mechanical conditions. `Semantic` condition delegation. Token budget check. Tests: all condition operators, Semantic fallthrough when no intelligence, path recording, budget exhaustion. |

### T6.4 — `go/pkg/decision/intelligence.go`

| | |
|---|---|
| **File** | `go/pkg/decision/intelligence.go` |
| **Depends on** | T1.1, T3.1, T3.2 |
| **Spec** | `docs/interfaces.md` → IIntelligence, IDecisionMaker |
| **Acceptance** | `IIntelligence` interface, `IDecisionMaker` interface, `Response` struct. Tests: interface compliance (mock implementations). |

---

## Tier 7: Trust and Authority

Depends on: Tier 5 (Actor), Tier 4 (Store).

### T7.1 — `go/pkg/trust/model.go`

| | |
|---|---|
| **File** | `go/pkg/trust/model.go` |
| **Depends on** | T5.1 (IActor), T1.1 (Score, Weight), T1.2 (DomainScope), T3.2 (TrustMetrics) |
| **Spec** | `docs/interfaces.md` → ITrustModel, `docs/trust.md` |
| **Acceptance** | `ITrustModel` interface, `DefaultTrustModel` implementation (linear decay, equal weighting, recovery penalty). Tests: score computation, domain-specific scoring, decay over time, recovery penalty at low trust, boundary values. |

### T7.2 — `go/pkg/authority/authority.go`

| | |
|---|---|
| **File** | `go/pkg/authority/authority.go` |
| **Depends on** | T5.1, T1.1, T1.2, T2.1 |
| **Spec** | `docs/interfaces.md` → IAuthorityChain, `docs/authority.md` |
| **Acceptance** | `IAuthorityChain` interface, `AuthorityResult`, default implementation. `AuthorityPolicy` matching. Trust-authority interaction (demotion). Tests: policy matching, chain evaluation, delegation, trust-based demotion, weighted authority. |

### T7.3 — `go/pkg/authority/chain.go`

| | |
|---|---|
| **File** | `go/pkg/authority/chain.go` |
| **Depends on** | T7.2, T4.2 (Store for edge queries) |
| **Spec** | `docs/authority.md` → Delegation, Delegation Chain Walk |
| **Acceptance** | Delegation chain walk using Authority edges. Weight propagation. Expiry handling. Tests: direct authority, single delegation, delegation weight capping, expired delegation. |

---

## Tier 8: Primitive Framework

Depends on: Tier 6 (Bus), Tier 5 (Store, Actor).

### T8.1 — `go/pkg/primitive/primitive.go`

| | |
|---|---|
| **File** | `go/pkg/primitive/primitive.go` |
| **Depends on** | T1.1, T1.2, T1.3, T3.1 |
| **Spec** | `docs/interfaces.md` → Primitive |
| **Acceptance** | `Primitive` interface, `Mutation` types (AddEvent, AddEdge, UpdateState, UpdateActivation, UpdateLifecycle), `MutationVisitor`. Tests: mutation construction, visitor dispatch. |

### T8.2 — `go/pkg/primitive/registry.go`

| | |
|---|---|
| **File** | `go/pkg/primitive/registry.go` |
| **Depends on** | T8.1 |
| **Spec** | `docs/interfaces.md` → Primitive (registry mentions) |
| **Acceptance** | `PrimitiveRegistry` — register, get by ID, list by layer, find subscribers for event type. Tests: registration, lookup, subscription matching, layer grouping. |

### T8.3 — `go/pkg/primitive/lifecycle.go`

| | |
|---|---|
| **File** | `go/pkg/primitive/lifecycle.go` |
| **Depends on** | T8.1, T1.3 (state machine) |
| **Spec** | `docs/interfaces.md` → Primitive invocation rules |
| **Acceptance** | Lifecycle integration — `Process()` only called when Active + cadence allows. Tests: dormant not invoked, cadence gating, lifecycle transitions during tick. |

### T8.4 — `go/pkg/primitive/harness.go`

| | |
|---|---|
| **File** | `go/pkg/primitive/harness.go` |
| **Depends on** | T8.1, T5.4 (InMemoryStore), T5.1 (Actor) |
| **Spec** | `docs/interfaces.md` → Testing Infrastructure (PrimitiveTestHarness) |
| **Acceptance** | `PrimitiveTestHarness` — WithStore, WithEvents, WithActors, Process, Tick, Mutations, EmittedEvents, StateChanges. Tests: harness processes a mock primitive, captures mutations without applying. |

---

## Tier 9: Tick Engine

Depends on: Tier 8 (Primitive framework), Tier 6 (Bus).

### T9.1 — `go/pkg/tick/snapshot.go`

| | |
|---|---|
| **File** | `go/pkg/tick/snapshot.go` |
| **Depends on** | T8.1, T1.1, T1.2 |
| **Spec** | `docs/interfaces.md` → Snapshot, PrimitiveState, `docs/coding-standards/go.md` → Frozen\<T\> |
| **Acceptance** | `FrozenSnapshot` — deep copy on construction, only getter methods. `PrimitiveState`. Tests: mutations to source don't affect snapshot, getters return copies. |

### T9.2 — `go/pkg/tick/cadence.go`

| | |
|---|---|
| **File** | `go/pkg/tick/cadence.go` |
| **Depends on** | T1.1 (Tick, Cadence) |
| **Spec** | `docs/tick-engine.md` → Cadence Gating |
| **Acceptance** | `CadenceAllows(primitive, currentTick) → bool`. Tests: cadence 1 allows every tick, cadence 5 allows every 5th, boundary cases. |

### T9.3 — `go/pkg/tick/engine.go`

| | |
|---|---|
| **File** | `go/pkg/tick/engine.go` |
| **Depends on** | T9.1, T9.2, T8.2 (registry), T6.1 (bus), T5.4 (store), T5.3 (factory) |
| **Spec** | `docs/tick-engine.md` — Full algorithm |
| **Acceptance** | `TickEngine` — RunTick, wave processing, layer constraint, quiescence detection, mutation application. `TickResult`. Tests: single wave quiescence, multi-wave ripple, wave limit (10), layer ordering, cadence gating, parallel within layer, mutation validation rejection. Coverage ≥ 90%. |

---

## Tier 10: Top-Level API

Depends on: Tier 9 (Tick engine), Tier 7 (Trust, Authority), Tier 6 (Decision).

### T10.1 — `go/pkg/graph/graph.go`

| | |
|---|---|
| **File** | `go/pkg/graph/graph.go` |
| **Depends on** | T9.3, T7.1, T7.2, T6.3, T6.1, T5.4, T5.2, T5.3 |
| **Spec** | `docs/interfaces.md` → Top-Level API (IGraph, IGraphQuery) |
| **Acceptance** | `IGraph` — Evaluate, Record, Query, Store, ActorStore, Bus, Registry, Start, Close. `IGraphQuery` — Events, Actors, Edges, Trust. `GraphConfig`. Tests: end-to-end Evaluate flow, Record creates event on graph, Query returns correct results, Start/Close lifecycle. |

### T10.2 — `go/pkg/graph/graph_test.go`

| | |
|---|---|
| **File** | `go/pkg/graph/graph_test.go` |
| **Depends on** | T10.1 |
| **Spec** | — |
| **Acceptance** | Integration tests: create graph with InMemoryStore, bootstrap, record events, evaluate decisions, query results, verify hash chain end-to-end. This is the "four lines of code" test. |

---

## Milestone: Phase 1 Complete

After Tier 10, the core is done. You can:
- Create a graph
- Record events (hash-chained, signed, causally linked)
- Evaluate decisions through trust and authority
- Query events, actors, edges, trust
- All with typed errors, constrained types, immutable domain objects

### What's NOT in Phase 1

- Individual primitive implementations (Phase 2)
- EGIP protocol (Phase 4)
- PostgresStore (can be added by running conformance suite against a new implementation)
- CLI (`go/cmd/eg/main.go`)

---

## Implementer Loop

For an automated implementer (Claude Code or similar):

```
for each tier in 0..10:
    for each task in tier (parallelisable within tier):
        1. READ the spec sections listed in task.Spec
        2. READ go.md for type mapping conventions
        3. READ any dependency files already implemented
        4. WRITE test file first (TDD)
        5. WRITE implementation to pass tests
        6. RUN: go build ./...
        7. RUN: go test -race -cover ./pkg/<package>/...
        8. VERIFY: coverage meets threshold
        9. RUN: go vet ./... && staticcheck ./...
        10. SELF-AUDIT: check against spec, coding standards, invariants
        11. If any step fails: fix and retry from step 6
        12. COMMIT with conventional commit message
```

### Self-Audit Checklist

After each task, verify:

- [ ] All types use constrained constructors (no bare `float64`, `string`, etc.)
- [ ] All errors are typed (no `fmt.Errorf` without wrapping a domain error)
- [ ] No `interface{}`/`any` except `PrimitiveState.State` (which is internal)
- [ ] All public types/functions have doc comments
- [ ] All struct fields are unexported (immutable pattern)
- [ ] Tests cover: happy path, boundary values, invalid input rejection, error types
- [ ] No magic values — all constants are named
- [ ] JSON serialisation round-trips correctly
- [ ] Concurrent access is safe (if applicable)

### Documentation Updates

After implementing a task, check if any docs need updating:

- If the implementation reveals a spec ambiguity → note it, continue, flag for spec update
- If a type name or signature differs from spec → update the spec to match (implementation is truth)
- If a new error case is discovered → add to `docs/interfaces.md` Domain Errors
- If conformance vectors need updating → update `docs/conformance/canonical-vectors.json`
- Update `ROADMAP.md` — check off completed tasks
