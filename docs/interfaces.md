# Core Interfaces

## Design Philosophy

SOLID throughout. Always-valid domain models. No magic values. No bare strings. No illegal states. Typed errors. Explicit optionality. Exhaustive matching. Immutable domain objects.

- **S** — Each primitive has a single domain. Each interface has a single responsibility.
- **O** — Every primitive ships with sensible defaults. Override with domain-specific logic without modifying the base.
- **L** — Any implementation of an interface is substitutable. The conformance test suite guarantees this.
- **I** — Small, focused interfaces. `IActor` doesn't know about trust. `ITrustModel` doesn't know about persistence.
- **D** — Primitives depend on abstractions (`IActor`, `ITrustModel`), never concrete types. The tick engine calls interfaces, not implementations.

### No Magic Values

Every type, every weight, every edge kind, every authority level is a defined constant or enum — never a bare string or number with implicit meaning. If you're comparing against `"trust"` instead of `EdgeType.Trust`, something is wrong. If `0.5` means something, it has a named constant.

### Always-Valid Domain Models

Domain objects guard themselves from ever becoming invalid. Validation happens at construction — if construction succeeds, the object is valid for its entire lifetime. There is no "partially constructed" or "needs validation later" state.

- **Construction fails, not runtime.** A `TrustScore` of 5.0 is a construction error, not a runtime surprise three layers deep.
- **No defensive checks downstream.** If you have an `Event`, it has a valid `Hash`, valid `Signature`, non-empty `Causes` (except Bootstrap). You don't re-check.
- **Boundary between valid and invalid is the constructor.** Inside the domain, everything is valid. Outside (user input, external APIs, deserialisation), everything is suspect.

### Make Illegal States Unrepresentable

The type system prevents invalid states, not runtime checks. If a state shouldn't exist, the types don't allow it.

- A `float64` that "should be 0-1" is a bug waiting to happen → use a constrained `Score` type
- An `Event` without a `Hash` should not compile → `Hash` is required, not optional
- A lifecycle transition from `Dormant` → `Processing` is invalid → the state machine enforces valid transitions
- An `Edge` with `Weight > 1.0` is impossible → construction rejects it
- An `Event` with empty `Causes` (except Bootstrap) is impossible → use `NonEmpty<EventID>`

### Immutability

All domain objects are immutable after construction. Events, Edges, Decisions, Receipts, TrustMetrics — once created, frozen. No setters. No mutation. If you need a different value, create a new object.

The `Snapshot` passed to `Primitive.Process()` is a deep-frozen read-only view. No primitive can mutate another primitive's state through the snapshot. Mutations are returned as declarative values, not applied by side effect.

The only mutable state is primitive internal state (key-value store), which is only modifiable via `UpdateState` mutations applied by the tick engine.

### Typed Errors

Errors are domain types, not strings. Every interface method returns a `Result<T, Error>` where `Error` is a discriminated union of possible failure modes. You match on the error type, you don't parse a string.

### Explicit Optionality

No null. No zero-value-means-absent. Optional fields use `Option<T>` — the caller must handle both `Some(value)` and `None`. If a field is always present, it's not `Option`.

### Exhaustive Matching

All enums must be matched exhaustively. When you switch on `DecisionOutcome`, the compiler forces you to handle `Permit`, `Deny`, `Defer`, AND `Escalate`. No silent `default:` that swallows new variants. Adding a new variant to an enum is a breaking change that the compiler surfaces everywhere.

---

## Design Patterns

Patterns used throughout. Named explicitly so implementations follow the same structure.

### Event Sourcing

The entire architecture. All state is derived from events. The event graph is the source of truth — if it isn't recorded, it didn't happen. Edges are derived from events. Actor status is derived from events. Trust scores are derived from events. There is no mutable state that isn't backed by an event trail.

### Factory

Domain objects with complex construction use factories, not bare constructors.

```
EventFactory {
    Create(type EventType, source ActorID, content EventContent, causes NonEmpty<EventID>,
           conversationID ConversationID, store Store, signer IIdentity) → Result<Event, StoreError>
}
```

The factory:
1. Validates `type` against `EventTypeRegistry`
2. Validates `content` matches the schema for the type
3. Validates `source` is an Active actor
4. Validates all `causes` exist in the Store
5. Retrieves `PrevHash` from the Store (last event in chain)
6. Computes `Hash` from canonical form
7. Signs with `signer` to produce `Signature`
8. Returns an immutable, fully-valid `Event`

No other path creates an Event. The constructor is private. This is the only way in.

```
BootstrapFactory {
    Init(systemActor ActorID, signer IIdentity) → Result<Event, StoreError>
}
```

Bootstrap is separate because it has no causes and no PrevHash — the only event where that's valid.

```
EdgeFactory {
    Create(from ActorID, to ActorID, edgeType EdgeType, weight Weight,
           direction EdgeDirection, scope Option<DomainScope>, store Store, signer IIdentity) → Result<Event, StoreError>
}
```

EdgeFactory emits an `edge.created` event and returns it. The Store indexes the event for edge queries. Edges are events — there is no separate edge creation path.

### Command (Mutations)

Mutations are commands — declarative descriptions of intent, not direct side effects. Primitives return mutations. The tick engine collects and applies them atomically.

This separation means:
- Primitives are pure functions (input → mutations). No hidden side effects.
- The tick engine can validate, reorder, or reject mutations before applying.
- Failed mutations don't leave the system in a partial state.
- Mutations are inspectable, loggable, auditable.

### Mediator (Tick Engine)

The tick engine mediates all communication between primitives. Primitives never call each other. They emit events, and the tick engine routes those events to subscribing primitives in the next wave.

This prevents:
- Circular dependencies between primitives
- One primitive blocking another
- Ordering surprises from direct calls
- State corruption from interleaved access

### Strategy

All `I`-prefixed interfaces are strategies with swappable implementations:

| Strategy | Default | Override for |
|---|---|---|
| `ITrustModel` | Linear decay, equal weighting | Domain-specific trust (e.g., +0.05 for task completion) |
| `IAuthorityChain` | Flat authority, no delegation | Hierarchical orgs, role-based access |
| `IIntelligence` | No-op (mechanical) | Claude, GPT, local models, deterministic logic |
| `IDecisionMaker` | Tree → intelligence fallthrough | Custom decision pipelines |
| `Store` | InMemoryStore | Postgres, SQLite, hosted service |
| `IActorStore` | InMemoryActorStore | Same backing as Store in production |
| `ITransport` | None | HTTP, WebSocket, gRPC, message queue |

### Repository (Store, IActorStore)

`Store` and `IActorStore` are repositories — they abstract persistence behind domain-oriented interfaces. The consumer doesn't know whether events are in memory, SQLite, or Postgres. The conformance test suite guarantees identical behaviour across implementations.

### Observer (Bus)

The bus wraps a Store with pub/sub. When an event is appended, all subscribers receive it. Primitives declare subscriptions via `SubscriptionPattern`. The bus matches events to subscribers.

Non-blocking — slow subscribers get dropped events, not blocked writers. Back-pressure is handled by the bus, not by the publisher.

### Chain of Responsibility (Authority)

`IAuthorityChain.Evaluate()` walks a delegation chain. Each link has a weight. The chain produces a weighted authority result, not a binary yes/no.

### Unit of Work (Tick Wave)

Each wave in a tick is a unit of work. All mutations from all primitives in a wave are collected and applied atomically. If any mutation fails validation, the entire wave can be rolled back.

### Visitor (Exhaustive Dispatch)

The mechanism that enforces exhaustive matching in languages without native sum types.

Every discriminated union in the spec has a corresponding visitor interface. Adding a new variant to the union forces a compile error in every visitor implementation — no silent `default:` branches.

```
// Mutation visitor — tick engine implements this
MutationVisitor {
    VisitAddEvent(AddEvent)
    VisitAddEdge(AddEdge)
    VisitUpdateState(UpdateState)
    VisitUpdateActivation(UpdateActivation)
    VisitUpdateLifecycle(UpdateLifecycle)
}

Mutation {
    Accept(visitor MutationVisitor)
}

// Error visitors
StoreErrorVisitor {
    VisitEventNotFound(EventNotFound)
    VisitActorNotFound(ActorNotFound)
    VisitEdgeNotFound(EdgeNotFound)
    VisitDuplicateEvent(DuplicateEvent)
    VisitCausalLinkMissing(CausalLinkMissing)
    VisitChainIntegrityViolation(ChainIntegrityViolation)
    VisitHashMismatch(HashMismatch)
    VisitSignatureInvalid(SignatureInvalid)
    VisitActorSuspended(ActorSuspended)
    VisitActorMemorial(ActorMemorial)
    VisitRateLimitExceeded(RateLimitExceeded)
    VisitStoreUnavailable(StoreUnavailable)
}

// Event content visitor — for processing events by type
EventContentVisitor {
    // Trust
    VisitTrustUpdated(TrustUpdatedContent)
    VisitTrustScore(TrustScoreContent)
    VisitTrustDecayed(TrustDecayedContent)
    // Authority
    VisitAuthorityRequested(AuthorityRequestContent)
    VisitAuthorityResolved(AuthorityResolvedContent)
    VisitAuthorityDelegated(AuthorityDelegatedContent)
    VisitAuthorityRevoked(AuthorityRevokedContent)
    VisitAuthorityTimeout(AuthorityTimeoutContent)
    // Actor
    VisitActorRegistered(ActorRegisteredContent)
    VisitActorSuspended(ActorSuspendedContent)
    VisitActorMemorial(ActorMemorialContent)
    // Edge
    VisitEdgeCreated(EdgeCreatedContent)
    VisitEdgeSuperseded(EdgeSupersededContent)
    // Integrity
    VisitViolationDetected(ViolationDetectedContent)
    VisitChainVerified(ChainVerifiedContent)
    VisitChainBroken(ChainBrokenContent)
    // System
    VisitBootstrap(BootstrapContent)
    VisitClockTick(ClockTickContent)
    VisitHealthReport(HealthReportContent)
    // Decision tree evolution
    VisitBranchProposed(BranchProposedContent)
    VisitBranchInserted(BranchInsertedContent)
    VisitCostReport(CostReportContent)
    // EGIP
    VisitEGIPHelloSent(EGIPHelloSentContent)
    VisitEGIPHelloReceived(EGIPHelloReceivedContent)
    VisitEGIPMessageSent(EGIPMessageSentContent)
    VisitEGIPMessageReceived(EGIPMessageReceivedContent)
    VisitEGIPReceiptSent(EGIPReceiptSentContent)
    VisitEGIPReceiptReceived(EGIPReceiptReceivedContent)
    VisitEGIPProofRequested(EGIPProofRequestedContent)
    VisitEGIPProofReceived(EGIPProofReceivedContent)
    VisitEGIPTreatyProposed(EGIPTreatyProposedContent)
    VisitEGIPTreatyActive(EGIPTreatyActiveContent)
    VisitEGIPTrustUpdated(EGIPTrustUpdatedContent)
}

EventContent {
    Accept(visitor EventContentVisitor)
}
```

**Language mapping:**
- **Rust** — `match` on enums is exhaustive natively. Visitor unnecessary — use `match`.
- **Go** — No exhaustive type switch. Visitor interface is the enforcement mechanism. Adding a method to the visitor interface breaks all implementations at compile time.
- **C#** — Pattern matching with exhaustiveness warnings. Visitor optional but recommended for complex hierarchies.
- **Python** — `@singledispatch` or `match` (3.10+). Visitor provides structural enforcement.

```
// Decision tree node visitor
DecisionNodeVisitor {
    VisitInternal(InternalNode)
    VisitLeaf(LeafNode)
}

// EGIP proof data visitor
ProofDataVisitor {
    VisitChainSegment(ChainSegmentProof)
    VisitEventExistence(EventExistenceProof)
    VisitChainSummary(ChainSummaryProof)
}

// EGIP message payload visitor
MessagePayloadVisitor {
    VisitHello(HelloPayload)
    VisitMessage(MessagePayload)
    VisitReceipt(ReceiptPayload)
    VisitProof(ProofPayload)
    VisitTreaty(TreatyPayload)
    VisitAuthorityRequest(AuthorityRequestPayload)
    VisitDiscover(DiscoverPayload)
}
```

The visitor pattern applies to: `Mutation`, `StoreError`, `DecisionError`, `ValidationError`, `EGIPError`, `EventContent`, `EdgeMetadata`, `MessagePayload`, `DecisionNode`, `ProofData`, and decision tree nodes.

### State (Lifecycle)

`LifecycleState` and `ActorStatus` are the state pattern — behaviour changes based on current state. A `Dormant` primitive ignores events. An `Active` primitive processes them. A `Suspended` actor can't emit. The state machine governs which operations are valid, and the transition function enforces it.

### Memento (Snapshot)

`Frozen<Snapshot>` is a memento — a captured, read-only snapshot of system state at a point in time. Primitives receive it during tick processing. It can't be mutated, and it doesn't change while the primitive is processing.

### Domain Service (Trust, Authority)

`ITrustModel` and `IAuthorityChain` are domain services — business logic that doesn't belong to a single entity. Trust is computed across events, actors, and domains. Authority evaluation walks delegation chains. These are stateless services that operate on the graph.

### Anti-Corruption Layer (EGIP)

EGIP is an anti-corruption layer between sovereign systems. Each system has its own event graph, its own types, its own trust model. EGIP translates between them via signed envelopes, CGERs, and treaties — without either system importing the other's domain model.

### Specification (Filters and Matching)

- `ActorFilter` — specifies which actors to query
- `SubscriptionPattern` — specifies which events a primitive receives
- `EventTypeRegistry.Validate()` — specifies valid content for an event type

### Facade (Top-Level API)

`graph.Evaluate()`, `graph.Record()`, `graph.Query()` are facades — simple entry points that hide the 200 primitives, tick engine, trust model, and authority chain behind a four-line API. See the Top-Level API section below.

### Decorator (Store Wrapping)

Stores can be wrapped for cross-cutting concerns without changing the interface:

- **Logging Store** — logs all operations for debugging
- **Metrics Store** — tracks operation counts, latencies
- **Caching Store** — caches frequent queries (Recent, ByType)
- **Validating Store** — additional integrity checks (paranoid mode)

Each wraps a `Store` and delegates, adding behaviour before/after. The conformance test suite must pass for any decorator chain.

---

## Generic Types

These are language-agnostic type constructors used throughout the spec.

```
Result<T, E>        // Either a success value T or an error E. Never both. Never neither.
                    // Go: (T, error) with typed error. Rust: Result<T, E>. C#: OneOf<T, E>.

Option<T>           // Either Some(T) or None. Explicit absence — no null, no zero values.
                    // Go: pointer + nil check, or custom Option type. Rust: Option<T>.

NonEmpty<T>         // A collection guaranteed to have at least one element.
                    // Construction: rejects empty input.
                    // .First() always succeeds — no need to check length.

Page<T>             // A paginated result set with cursor-based navigation.
                    // { Items: []T, Cursor: Option<Cursor>, HasMore: bool }

Frozen<T>           // A deeply immutable view of T. No references to mutable state.
                    // Used for Snapshot passed to primitives.
```

---

## Value Objects

Immutable. Compared by value. Validated at construction. These are the building blocks everything else is made of.

### Typed IDs

Bare strings are not IDs. These are distinct types — you cannot pass an `ActorID` where an `EventID` is expected.

```
EventID         // UUID v7 (time-ordered). Identifies an event.
                // Construction: validates UUID v7 format. Rejects malformed input.
                // Equality: by value. Two EventIDs with the same string are equal.

ActorID         // Identifies an actor. Derived from public key hash.
                // Construction: validated format. Cannot be empty.

EdgeID          // The EventID of the event that created this edge.
                // Construction: same as EventID.

ConversationID  // Groups related events into threads.
                // Construction: validated format. Cannot be empty.

Hash            // SHA-256 hex string (64 characters). Used for chain integrity.
                // Construction: validates length (64) and hex characters. Rejects "abc" or "".

SystemURI       // Identifies a remote system in EGIP.
                // Construction: validates URI format.

EnvelopeID      // UUID. Identifies an EGIP message.
                // Construction: validates UUID format.

TreatyID        // UUID. Identifies a bilateral treaty.
                // Construction: validates UUID format.

PrimitiveID     // Identifies a primitive instance.
                // Construction: validated name + layer. Layer must be 0-13.

EventType       // A registered event type string (e.g., "trust.updated").
                // Construction: validates against EventTypeRegistry. Rejects unregistered types.
                // Hierarchical — dot-separated namespace (e.g., "trust.updated", "actor.registered").

SubscriptionPattern // A glob pattern for matching event types (e.g., "trust.*", "*").
                    // Construction: validates dot-separated segments with optional * wildcard.
                    // Used by primitives to declare which events they subscribe to.

DomainScope     // A trust/authority domain (e.g., "code_review", "financial", "medical").
                // Construction: validated format (lowercase, dot-separated namespace).
                // Used for: domain-specific trust, scoped authority, edge scoping.
                // Open vocabulary — product layers define their own domains.

PublicKey       // Ed25519 public key (32 bytes).
                // Construction: validates length. Rejects wrong size.

Signature       // Ed25519 signature (64 bytes).
                // Construction: validates length. Rejects wrong size.
```

All are strings underneath in serialisation, but typed in code. The compiler catches misuse, not a runtime crash.

### Constrained Numerics

Bare `float64` is not a score, weight, or confidence. These are domain types with invariants enforced at construction.

```
Score           // float64 constrained to [0.0, 1.0].
                // Construction: rejects values outside range.
                // Used for: trust scores, confidence levels, activation levels, authority weights.

Weight          // float64 constrained to [-1.0, 1.0].
                // Construction: rejects values outside range.
                // Negative values represent penalties.
                // Used for: edge weights, trust adjustments, trends.

Activation      // float64 constrained to [0.0, 1.0].
                // Construction: rejects values outside range.
                // Used for: primitive activation levels.

Cadence         // int constrained to [1, ∞).
                // Construction: rejects zero or negative.
                // Used for: minimum ticks between primitive invocations.

Tick            // int constrained to [0, ∞).
                // Construction: rejects negative.

Layer           // int constrained to [0, 13].
                // Construction: rejects values outside range.
```

### Lifecycle State Machine

Not an enum you can set to anything — a state machine that enforces valid transitions. Invalid transitions return an error, not a runtime surprise.

```
LifecycleState {
    Dormant      → { Activating }
    Activating   → { Active, Dormant }       // activation can fail
    Active       → { Processing, Deactivating }
    Processing   → { Emitting, Active }       // processing can produce nothing
    Emitting     → { Active }
    Deactivating → { Dormant }
}

transition(from, to) → Result<LifecycleState, InvalidTransition>
validTransitions(from) → []LifecycleState
```

### Actor Status Transitions

Similarly constrained:

```
ActorStatus {
    Active     → { Suspended, Memorial }
    Suspended  → { Active, Memorial }       // can be reactivated
    Memorial   → { }                         // terminal — graph preserved forever
}
```

`Memorial` is irreversible. Once an actor is memorialised, their graph is preserved but they can never emit new events.

---

## Domain Errors

Errors are typed, not strings. Every failure mode is a distinct type. You match on the type, you don't parse a message.

```
// Store errors
StoreError =
    | EventNotFound { id: EventID }
    | ActorNotFound { id: ActorID }
    | EdgeNotFound { from: ActorID, to: ActorID, edgeType: EdgeType }
    | DuplicateEvent { id: EventID }                    // idempotency — same event submitted twice
    | CausalLinkMissing { eventID: EventID, missingCause: EventID }
    | ChainIntegrityViolation { position: int, expected: Hash, actual: Hash }
    | HashMismatch { eventID: EventID, computed: Hash, stored: Hash }
    | SignatureInvalid { eventID: EventID, signer: ActorID }
    | ActorSuspended { id: ActorID }                    // suspended actor tried to emit
    | ActorMemorial { id: ActorID }                     // memorialised actor tried to emit
    | RateLimitExceeded { actor: ActorID, limit: int, window: string }
    | StoreUnavailable { reason: string }               // backing database down

// Decision errors
DecisionError =
    | ActorNotFound { id: ActorID }
    | InsufficientAuthority { actor: ActorID, action: string, required: AuthorityLevel }
    | TrustBelowThreshold { actor: ActorID, score: Score, required: Score }
    | CausesRequired { action: string }                 // decision must reference causing events

// Validation errors (construction failures)
ValidationError =
    | OutOfRange { field: string, value: float64, min: float64, max: float64 }
    | InvalidFormat { field: string, value: string, expected: string }
    | EmptyRequired { field: string }
    | InvalidLifecycleTransition { from: LifecycleState, to: LifecycleState, validTargets: []LifecycleState }
    | InvalidActorTransition { from: ActorStatus, to: ActorStatus, validTargets: []ActorStatus }

// EGIP errors
EGIPError =
    | SystemNotFound { uri: SystemURI }
    | EnvelopeSignatureInvalid { envelopeID: EnvelopeID }
    | TreatyViolation { treatyID: TreatyID, term: string }
    | TrustInsufficient { system: SystemURI, score: Score, required: Score }
    | TransportFailure { to: SystemURI, reason: string }
```

Every interface method in this spec returns `Result<T, Error>` using the appropriate error type. For example:

- `Store.Get(id EventID) → Result<Event, StoreError>`
- `IDecisionMaker.Decide(ctx, input) → Result<Decision, DecisionError>`
- `IActorStore.Get(id ActorID) → Result<IActor, StoreError>`

---

## Constants and Enums

No magic values. These are the defined vocabularies. All enums must be matched exhaustively — no silent `default:` branches.

```
// Edge types — every relationship in the graph
EdgeType {
    Trust           // actor trusts actor (weighted)
    Authority       // actor has authority over actor/scope (weighted)
    Subscription    // actor follows actor
    Endorsement     // actor stakes reputation on event/actor (weighted)
    Delegation      // actor grants authority to actor (scoped)
    Causation       // event caused event (also in Event.Causes)
    Reference       // event references event (cross-graph via CGER)
    Channel         // private bidirectional link between actors
    Annotation      // metadata attached to event
}

// Authority levels
AuthorityLevel {
    Required        // blocks until human approves
    Recommended     // auto-approves after timeout
    Notification    // auto-approves immediately, logged
}

// Decision outcomes
DecisionOutcome {
    Permit          // action allowed
    Deny            // action denied
    Defer           // insufficient information — try again later
    Escalate        // confidence too low — needs human
}

// Actor types — what kind of decision-maker this is
ActorType {
    Human           // a person
    AI              // an AI agent or model
    System          // the graph infrastructure itself
    Committee       // a group that votes
    RulesEngine     // deterministic logic
}

// Actor status (transitions enforced by state machine — see above)
ActorStatus { Active, Suspended, Memorial }

// Lifecycle states (transitions enforced by state machine — see above)
LifecycleState { Dormant, Activating, Active, Processing, Emitting, Deactivating }

// Edge direction (from social grammar)
EdgeDirection {
    Centripetal      // toward the content/target
    Centrifugal      // into the actor's subgraph
}

// Expectation status
ExpectationStatus {
    Pending         // waiting to be fulfilled
    Met             // fulfilled before deadline
    Violated        // not fulfilled, violation detected
    Expired         // deadline passed, violation assessment pending
}

// Severity levels (for violations)
SeverityLevel {
    Info            // notable but not concerning
    Warning         // potentially problematic
    Serious         // requires attention
    Critical        // requires immediate action
}

// EGIP message types
MessageType { Hello, Message, Receipt, Proof, Treaty, AuthorityRequest, Discover }

// Treaty status
TreatyStatus { Proposed, Active, Suspended, Terminated }

// Integrity violation types
IntegrityViolationType {
    ChainBreak          // hash chain is broken at a specific position
    HashMismatch        // stored hash doesn't match computed hash
    MissingCause        // event references a cause that doesn't exist
    SignatureInvalid    // signature doesn't verify against public key
    OrphanEvent         // event exists outside the chain
}

// System invariants — the 10 non-negotiable rules
InvariantName {
    Causality       // every event declares its causes
    Integrity       // all events are hash-chained
    Observable      // all significant operations emit events
    SelfEvolve      // the system improves itself
    Dignity         // agents have identity, state, and lifecycle
    Transparent     // humans know when interacting with automation
    Consent         // no significant action without approval
    Authority       // significant actions require appropriate authority
    Verify          // all code changes are built and tested
    Record          // if it isn't recorded, it didn't happen
}

// Cross-graph event relationship
CGERRelationship { CausedBy, References, RespondsTo }

// EGIP receipt status
ReceiptStatus { Delivered, Processed, Rejected }

// EGIP proof types
ProofType { ChainSegment, EventExistence, ChainSummary }

// EGIP treaty actions
TreatyAction { Propose, Accept, Modify, Suspend, Terminate }

// Decision tree condition operators
ConditionOperator {
    Equals              // exact match
    GreaterThan         // numeric comparison
    LessThan            // numeric comparison
    InRange             // between two values
    Matches             // pattern match (SubscriptionPattern)
    Exists              // field is present (not None)
    Semantic            // delegates to IIntelligence — returns Score, branches on threshold
}
```

---

## Event Type Registry

Event types are not bare strings — they are registered in a compile-time registry that maps each type string to its content schema. `Store.Append()` validates content against the registered schema before the event enters the graph.

```
EventTypeRegistry {
    Register<C>(type EventType, schema Schema<C>)
    Validate(type EventType, content EventContent) → Result<EventContent, ValidationError>
    Schema(type EventType) → Option<Schema>
    AllTypes() → []EventType
}
```

### Typed Event Content

Every event type has a corresponding content struct. No `map[string]any` — the content is typed at the event type level.

```
// Trust events
"trust.updated" → TrustUpdatedContent {
    Actor:    ActorID
    Previous: Score
    Current:  Score
    Domain:   DomainScope
    Cause:    EventID
}

"trust.score" → TrustScoreContent {
    Actor:   ActorID
    Metrics: TrustMetrics
}

// Authority events
"authority.requested" → AuthorityRequestContent {
    Action:      string
    Actor:       ActorID
    Level:       AuthorityLevel
    Justification: string
    Causes:      NonEmpty<EventID>
}

"authority.resolved" → AuthorityResolvedContent {
    RequestID:  EventID
    Approved:   bool
    Resolver:   ActorID
    Reason:     Option<string>
}

// Actor events
"actor.registered" → ActorRegisteredContent {
    ActorID:    ActorID
    PublicKey:  PublicKey
    Type:       ActorType
}

"actor.suspended" → ActorSuspendedContent {
    ActorID:    ActorID
    Reason:     EventID
}

"actor.memorial" → ActorMemorialContent {
    ActorID:    ActorID
    Reason:     EventID
}

// Edge events
"edge.created" → EdgeCreatedContent {
    From:      ActorID
    To:        ActorID
    EdgeType:  EdgeType
    Weight:    Weight
    Direction: EdgeDirection
    Scope:     Option<DomainScope>
    ExpiresAt: Option<time>
}

"edge.superseded" → EdgeSupersededContent {
    PreviousEdge: EdgeID
    NewEdge:      EdgeID
    Reason:       EventID
}

// Violation events
"violation.detected" → ViolationDetectedContent {
    Expectation: EventID
    Actor:       ActorID
    Severity:    SeverityLevel
    Description: string
    Evidence:    NonEmpty<EventID>
}

// Integrity events
"chain.verified" → ChainVerifiedContent {
    Valid:     bool
    Length:    int
    Duration:  duration
}

"chain.broken" → ChainBrokenContent {
    Position: int
    Expected: Hash
    Actual:   Hash
}

// System events
"system.bootstrapped" → BootstrapContent {
    ActorID:      ActorID
    ChainGenesis: Hash
    Timestamp:    time
}

// Clock events
"clock.tick" → ClockTickContent {
    Tick:      Tick
    Timestamp: time
    Elapsed:   duration
}

// Health events
"health.report" → HealthReportContent {
    Overall:          Score
    ChainIntegrity:   bool
    PrimitiveHealth:  map[PrimitiveID]Score
    ActiveActors:     int
    EventRate:        float64         // events per second — not a Score, genuinely unbounded
}

// Trust decay events
"trust.decayed" → TrustDecayedContent {
    Actor:    ActorID
    Previous: Score
    Current:  Score
    Elapsed:  duration
    Rate:     Score
}

// Authority events (additional)
"authority.delegated" → AuthorityDelegatedContent {
    From:       ActorID
    To:         ActorID
    Scope:      DomainScope
    Weight:     Score
    ExpiresAt:  Option<time>
}

"authority.revoked" → AuthorityRevokedContent {
    From:       ActorID
    To:         ActorID
    Scope:      DomainScope
    Reason:     EventID
}

"authority.timeout" → AuthorityTimeoutContent {
    RequestID:  EventID
    Level:      AuthorityLevel           // Recommended — the only level that times out
    Duration:   duration
}

// Decision tree evolution events
"decision.branch.proposed" → BranchProposedContent {
    PrimitiveID:  PrimitiveID
    TreeVersion:  int
    Condition:    Condition
    Outcome:      DecisionOutcome
    Accuracy:     Score
    SampleSize:   int
}

"decision.branch.inserted" → BranchInsertedContent {
    PrimitiveID:  PrimitiveID
    TreeVersion:  int
    Path:         []PathStep
    Outcome:      DecisionOutcome
    Confidence:   Score
}

"decision.cost.report" → CostReportContent {
    PrimitiveID:    PrimitiveID
    TreeVersion:    int
    TotalLeaves:    int
    LLMLeaves:      int
    MechanicalRate: Score               // proportion handled without LLM
    TotalTokens:    int
}

// EGIP events (local graph records of inter-system activity)
"egip.hello.sent" → EGIPHelloSentContent { To: SystemURI }
"egip.hello.received" → EGIPHelloReceivedContent { From: SystemURI, PublicKey: PublicKey }
"egip.message.sent" → EGIPMessageSentContent { To: SystemURI, EnvelopeID: EnvelopeID }
"egip.message.received" → EGIPMessageReceivedContent { From: SystemURI, EnvelopeID: EnvelopeID }
"egip.receipt.sent" → EGIPReceiptSentContent { EnvelopeID: EnvelopeID, Status: ReceiptStatus }
"egip.receipt.received" → EGIPReceiptReceivedContent { EnvelopeID: EnvelopeID, Status: ReceiptStatus }
"egip.proof.requested" → EGIPProofRequestedContent { System: SystemURI, ProofType: ProofType }
"egip.proof.received" → EGIPProofReceivedContent { System: SystemURI, Valid: bool }
"egip.treaty.proposed" → EGIPTreatyProposedContent { TreatyID: TreatyID, To: SystemURI }
"egip.treaty.active" → EGIPTreatyActiveContent { TreatyID: TreatyID, With: SystemURI }
"egip.trust.updated" → EGIPTrustUpdatedContent { System: SystemURI, Previous: Score, Current: Score, Evidence: EnvelopeID }
```

When implementing a new primitive that emits events, you must:
1. Define the content struct for each event type
2. Register it in the EventTypeRegistry
3. The Store will reject events whose content doesn't match the registered schema

This eliminates the last `map[string]any` hole. Event content is typed at the domain level.

---

## Edge Type Registry

Parallel to `EventTypeRegistry` — maps each `EdgeType` to its metadata schema. Edges carry typed metadata, not generic maps.

```
EdgeTypeRegistry {
    Register<M>(edgeType EdgeType, schema Schema<M>)
    Validate(edgeType EdgeType, metadata EdgeMetadata) → Result<EdgeMetadata, ValidationError>
    Schema(edgeType EdgeType) → Option<Schema>
    AllTypes() → []EdgeType
}
```

### Typed Edge Metadata

Every edge type has a corresponding metadata struct. `EdgeMetadata` is the interface — concrete type is determined by `EdgeType`.

```
// Trust edge metadata
EdgeType.Trust → TrustEdgeMetadata {
    Domain:         DomainScope
    Evidence:       []EventID           // events that contributed to this trust level
    DecayRate:      Score               // [0.0, 1.0]
    LastUpdated:    time
}

// Authority edge metadata
EdgeType.Authority → AuthorityEdgeMetadata {
    Scope:          DomainScope
    Delegated:      bool                // true if this is a delegation, not direct authority
    DelegatedFrom:  Option<ActorID>     // who delegated (if Delegated)
    Constraints:    []string            // freeform policy constraints on this delegation
}

// Subscription edge metadata
EdgeType.Subscription → SubscriptionEdgeMetadata {
    Patterns:       []SubscriptionPattern   // what event patterns the subscriber wants
    Muted:          bool                    // subscriber has silenced notifications
}

// Endorsement edge metadata
EdgeType.Endorsement → EndorsementEdgeMetadata {
    Stake:          Score               // [0.0, 1.0] — how much reputation is staked
    Target:         EventID             // which event is being endorsed
    Domain:         Option<DomainScope>
}

// Delegation edge metadata
EdgeType.Delegation → DelegationEdgeMetadata {
    Scope:          DomainScope
    Constraints:    []string            // what the delegate may do
    RevokedBy:      Option<EventID>     // None = active. Some = revoked by this event.
}

// Causation edge metadata
EdgeType.Causation → CausationEdgeMetadata {
    Relationship:   string              // "caused", "contributed", "triggered"
}

// Reference edge metadata (cross-graph)
EdgeType.Reference → ReferenceEdgeMetadata {
    CGER:           Option<CGER>        // cross-graph reference if remote
    Verified:       bool
}

// Channel edge metadata
EdgeType.Channel → ChannelEdgeMetadata {
    Encrypted:      bool
    CreatedBy:      ActorID             // who initiated the channel
}

// Annotation edge metadata
EdgeType.Annotation → AnnotationEdgeMetadata {
    Key:            string              // annotation key
    Value:          EventContent        // typed annotation value
    Annotator:      ActorID
}
```

The `EdgeMetadataVisitor` enforces exhaustive handling:

```
EdgeMetadataVisitor {
    VisitTrust(TrustEdgeMetadata)
    VisitAuthority(AuthorityEdgeMetadata)
    VisitSubscription(SubscriptionEdgeMetadata)
    VisitEndorsement(EndorsementEdgeMetadata)
    VisitDelegation(DelegationEdgeMetadata)
    VisitCausation(CausationEdgeMetadata)
    VisitReference(ReferenceEdgeMetadata)
    VisitChannel(ChannelEdgeMetadata)
    VisitAnnotation(AnnotationEdgeMetadata)
}
```

---

## Data Structures

All data structures are **immutable after construction**. No setters. No mutation. If you need a different value, construct a new object. Construction validates all fields — if the constructor returns, the object is valid forever.

### Event

The fundamental unit. Every significant action is an event. Immutable after construction.

```
Event {
    Version:        int                  // schema version for this event type (starts at 1)
    ID:             EventID              // UUID v7 (time-ordered)
    Type:           EventType            // validated against EventTypeRegistry at construction
    Timestamp:      time                 // nanosecond precision
    Source:         ActorID              // who/what emitted this — must be Active actor
    Content:        EventContent         // interface — concrete type determined by EventType via EventTypeRegistry
    Causes:         NonEmpty<EventID>    // causal DAG — at least one cause (see Bootstrap)
    ConversationID: ConversationID       // thread grouping
    Hash:           Hash                 // SHA-256 of canonical form
    PrevHash:       Hash                 // hash chain link
    Signature:      Signature            // Ed25519 signature of canonical form
}
```

**Bootstrap exception:** The genesis event is the only event with empty causes. It is constructed via `BootstrapFactory.Init()`, not `EventFactory.Create()`. The normal factory enforces `NonEmpty<EventID>` — only the bootstrap path bypasses this.

**Construction validates:**
- ID is valid UUID v7
- Type is registered in EventTypeRegistry
- Source is an Active actor (not Suspended, not Memorial)
- Content matches the schema for the event type
- All Causes exist in the Store
- Hash is correctly computed from canonical form
- Signature verifies against Source's public key

### Edge

A typed, weighted, directional relationship. Immutable after construction.

```
Edge {
    ID:        EdgeID              // the EventID of the event that created this edge
    From:      ActorID             // who/what this edge originates from
    To:        ActorID             // who/what this edge points to
    Type:      EdgeType            // Trust, Authority, Subscription, etc.
    Weight:    Weight              // [-1.0, 1.0] — constrained at construction
    Direction: EdgeDirection       // Centripetal or Centrifugal
    Scope:     Option<DomainScope>  // what domain this edge applies to (e.g., "code_review")
    Metadata:  EdgeMetadata        // interface — concrete type determined by EdgeType via EdgeTypeRegistry (see Edge Type Registry section)
    CreatedAt: time
    ExpiresAt: Option<time>        // None = permanent. Some(time) = transient edge.
}
```

Edges are never deleted — they are superseded by new edge events. Querying returns the current effective state.

### Decision

What `IDecisionMaker.Decide()` returns. Immutable. Full epistemic context.

```
Decision {
    Outcome:        DecisionOutcome      // Permit, Deny, Defer, Escalate
    Confidence:     Score                // [0.0, 1.0]
    AuthorityChain: NonEmpty<AuthorityLink> // the chain of authority — at least one link
    TrustWeights:   []TrustWeight        // trust scores of relevant actors (may be empty for mechanical decisions)
    Evidence:       NonEmpty<EventID>    // events supporting the decision — a decision without evidence is not a decision
    Receipt:        Receipt              // cryptographic proof
    NeedsHuman:     bool                 // true if confidence too low for autonomous action
}

AuthorityLink {
    Actor:    ActorID
    Level:    AuthorityLevel
    Weight:   Score                      // [0.0, 1.0]
}

TrustWeight {
    Actor:    ActorID
    Score:    Score                       // [0.0, 1.0]
    Domain:   DomainScope
}

Receipt {
    Hash:       Hash                     // SHA-256 of the decision + inputs
    Timestamp:  time
    SignedBy:   ActorID
    Signature:  Signature                // Ed25519 signature
    InputHash:  Hash                     // hash of what was evaluated
    ChainPos:   EventID                  // position in the hash chain
}
```

### TrustMetrics

What trust queries return. Immutable.

```
TrustMetrics {
    Actor:       ActorID
    Overall:     Score                    // [0.0, 1.0]
    ByDomain:    map[DomainScope]Score    // per-domain trust
    Confidence:  Score                    // [0.0, 1.0]
    Trend:       Weight                   // [-1.0, 1.0]
    Evidence:    []EventID
    LastUpdated: time
    DecayRate:   Score                    // [0.0, 1.0]
}
```

### Expectation

What should happen after an event. Immutable.

```
Expectation {
    ID:          EventID
    Trigger:     EventID
    Description: string
    Deadline:    time
    Severity:    SeverityLevel
    Status:      ExpectationStatus       // Pending, Met, Violated, Expired
}
```

### Snapshot

The read-only view passed to primitives during tick processing. Deeply frozen via `Frozen<Snapshot>`.

```
Snapshot {
    Tick:            Tick                               // current tick number
    Primitives:      map[PrimitiveID]PrimitiveState     // all primitive states (read-only)
    PendingEvents:   []Event                            // events waiting to be processed
    RecentEvents:    []Event                            // last N events for context
    ActiveActors:    []IActor                           // currently active actors
}

PrimitiveState {
    ID:          PrimitiveID
    Layer:       Layer
    Lifecycle:   LifecycleState
    Activation:  Activation
    Cadence:     Cadence
    State:       map[string]any                        // primitive-internal state (read-only view)
    LastTick:    Tick                                   // last tick this primitive was invoked
}
```

Primitives use the snapshot to read other primitives' states, check recent events, and see active actors — without being able to mutate anything.

### Schema

Defines the expected structure for typed content (events, edges, message payloads).

```
Schema<C> {
    Validate(content C) → Result<C, ValidationError>
    ContentType() → string                             // Go: reflect.Type, Rust: TypeId
}
```

Schemas are registered in `EventTypeRegistry` at startup. The Store validates all content against the registered schema on `Append`. Unregistered event types are rejected at `EventType` construction.

### ViolationRecord

When an expectation is not met. Immutable.

```
ViolationRecord {
    ID:           EventID
    Expectation:  EventID
    Severity:     SeverityLevel
    Actor:        ActorID
    Description:  string
    Evidence:     NonEmpty<EventID>       // at least one piece of evidence
}
```

### Decision Tree

The mechanical-to-intelligent continuum. Optional per primitive. See `decision-trees.md` for the full evolution algorithm.

```
DecisionTree {
    Root:       DecisionNode
    Version:    int                     // incremented on each evolution step
    Stats:      TreeStats
}

DecisionNode = InternalNode | LeafNode

InternalNode {
    Condition:  Condition               // what to test
    Branches:   NonEmpty<Branch>        // at least one branch
    Default:    DecisionNode            // fallthrough if no branch matches
}

Branch {
    Match:      MatchValue              // what the condition result must equal
    Child:      DecisionNode            // where to go on match
}

LeafNode {
    Outcome:      Option<DecisionOutcome>  // Some = deterministic. None = needs intelligence.
    NeedsLLM:     bool                     // true if this leaf requires IIntelligence
    Confidence:   Score                    // [0.0, 1.0]
    Stats:        LeafStats
}

Condition {
    Field:      string                  // dot-path into DecisionInput context
    Operator:   ConditionOperator
    Threshold:  Option<Score>           // for numeric and Semantic comparisons
    Prompt:     Option<string>          // only for Semantic — what to ask IIntelligence
}

MatchValue {
    String:     Option<string>
    Number:     Option<float64>
    Boolean:    Option<bool>
    EventType:  Option<EventType>
}

PathStep {
    Condition:  Condition
    Branch:     MatchValue              // which branch was taken (or "default")
}

TreeStats {
    TotalHits:       int
    MechanicalHits:  int                // resolved without IIntelligence
    LLMHits:         int                // required IIntelligence
    TotalTokens:     int                // cumulative token usage
    LastEvolution:   Option<time>       // last time a branch was extracted
}

LeafStats {
    HitCount:        int
    LLMCallCount:    int
    LastLLMResponse: Option<Response>
    ResponseHistory: []ResponseRecord   // last N responses for pattern detection
    PatternScore:    Score              // [0.0, 1.0] — how patterned the responses are
}

ResponseRecord {
    Input:      DecisionInput
    Output:     DecisionOutcome
    Confidence: Score
    Timestamp:  time
}

SemanticEvalRecord {
    Input:      DecisionInput
    Score:      Score                   // what IIntelligence returned
    Branch:     MatchValue              // which branch was taken
    Prompt:     string                  // the Semantic condition's prompt
    TokensUsed: int                    // for budget tracking
}
```

### Authority Request and Resolution

Authority workflow data structures. See `authority.md` for the full flow.

```
AuthorityRequest {
    ID:             EventID
    Action:         string
    Actor:          ActorID
    Level:          AuthorityLevel
    Justification:  string
    Causes:         NonEmpty<EventID>
    Context:        map[string]any      // domain-specific context for the approver
    CreatedAt:      time
    ExpiresAt:      Option<time>        // only for Recommended (timeout)
}

AuthorityResolution {
    RequestID:      EventID
    Approved:       bool
    Resolver:       ActorID
    Reason:         Option<string>
    ResolvedAt:     time
    AutoApproved:   bool                // true if timeout/notification auto-approved
}

AuthorityPolicy {
    Action:         string              // or pattern — "actor.suspend", "trust.*", "*"
    Level:          AuthorityLevel
    Approvers:      []ActorID           // who can approve (empty = any human)
    MinTrust:       Option<Score>       // minimum trust to bypass Required → Recommended
    Scope:          Option<DomainScope>
}
```

### Tick Result

What the tick engine returns after processing. See `tick-engine.md`.

```
TickResult {
    Tick:       Tick
    Waves:      int                     // how many waves before quiescence
    Mutations:  int                     // total mutations applied
    Duration:   duration                // wall-clock time
    Quiesced:   bool                    // true if terminated naturally, false if hit wave limit
}
```

### Subgraph Criteria

Filter for subgraph extraction queries.

```
SubgraphCriteria {
    MaxDepth:   int                     // max causal distance from root
    EventTypes: []EventType             // filter by type (empty = all)
    Sources:    []ActorID               // filter by source (empty = all)
    Since:      Option<time>            // only events after this time
    Until:      Option<time>            // only events before this time
}
```

### EGIP Message Payloads

Concrete types for `Envelope.Payload`. See `protocol.md` for the full protocol specification.

```
HelloPayload {
    SystemURI:          SystemURI
    PublicKey:           PublicKey
    ProtocolVersions:   []int
    Capabilities:       []string        // what this system supports
    ChainLength:        int
}

MessagePayload {                        // the MessagePayload interface — this is one concrete type
    Content:            EventContent
    ContentType:        EventType
    ConversationID:     Option<ConversationID>
    CGERs:              []CGER
}

ReceiptPayload {
    EnvelopeID:         EnvelopeID
    Status:             ReceiptStatus
    LocalEventID:       Option<EventID>
    Reason:             Option<string>
    Signature:          Signature
}

ProofPayload {
    ProofType:          ProofType
    Data:               ProofData       // discriminated union — see below
}

ChainSegmentProof {
    Events:             []Event
    StartHash:          Hash
    EndHash:            Hash
}

EventExistenceProof {
    Event:              Event
    PrevHash:           Hash
    NextHash:           Option<Hash>
    Position:           int
    ChainLength:        int
}

ChainSummaryProof {
    Length:              int
    HeadHash:           Hash
    GenesisHash:        Hash
    Timestamp:          time
}

TreatyPayload {
    TreatyID:           TreatyID
    Action:             TreatyAction
    Terms:              []TreatyTerm    // for Propose and Modify
    Reason:             Option<string>  // for Suspend and Terminate
}

AuthorityRequestPayload {
    Action:             string
    Actor:              ActorID
    Level:              AuthorityLevel
    Justification:      string
    Context:            map[string]any
    TreatyID:           Option<TreatyID>
}

DiscoverPayload {
    Query:              Option<DiscoverQuery>    // Some for requests, None for responses
    Results:            []DiscoverResult         // populated in responses
}

DiscoverQuery {
    Capabilities:       []string
    MinTrust:           Option<Score>
}

DiscoverResult {
    SystemURI:          SystemURI
    PublicKey:           PublicKey
    Capabilities:       []string
    TrustScore:         Score           // the responder's trust in this system
}
```

`ProofData` is a discriminated union of `ChainSegmentProof`, `EventExistenceProof`, and `ChainSummaryProof`. Use the Visitor pattern for exhaustive dispatch:

```
ProofDataVisitor {
    VisitChainSegment(ChainSegmentProof)
    VisitEventExistence(EventExistenceProof)
    VisitChainSummary(ChainSummaryProof)
}
```

---

## Core Interfaces

All interface methods return `Result<T, Error>` with typed domain errors. No untyped exceptions. No string error messages you have to parse.

### IActor

Identity. What everything else references. Immutable — actor state changes create new events and new state, not mutation.

```
IActor {
    ID() → ActorID
    PublicKey() → PublicKey
    DisplayName() → string
    Type() → ActorType             // Human, AI, System, Committee, RulesEngine — for TRANSPARENT invariant
    Metadata() → map[string]any    // extension point — product layers attach domain-specific data
    CreatedAt() → time
    Status() → ActorStatus         // Active, Suspended, Memorial
}
```

Implementations: human actor (keypair from wallet/device), AI actor (keypair from system), system actor (the graph itself). The graph doesn't distinguish behaviourally — an actor is an actor.

### IActorStore

Actor persistence. Separate from event Store — single responsibility.

```
IActorStore {
    Register(publicKey PublicKey, displayName string, actorType ActorType) → Result<IActor, StoreError>
    Get(id ActorID) → Result<IActor, StoreError>
    GetByPublicKey(publicKey PublicKey) → Result<IActor, StoreError>
    Update(id ActorID, updates ActorUpdate) → Result<IActor, StoreError>
    List(filter ActorFilter) → Result<Page<IActor>, StoreError>
    Suspend(id ActorID, reason EventID) → Result<IActor, StoreError>
    Memorial(id ActorID, reason EventID) → Result<IActor, StoreError>
}

ActorUpdate {
    DisplayName: Option<string>
    Metadata:    Option<map[string]any>      // extension point — merged with existing if Some
}

ActorFilter {
    Status:    Option<ActorStatus>           // None = all statuses
    Type:      Option<ActorType>             // None = all types
    Limit:     int
    After:     Option<Cursor>                // for pagination
}
```

Actor lifecycle changes (Register, Suspend, Memorial) emit events on the graph. The IActorStore handles persistence; the events provide the audit trail.

### Store

Event and edge persistence. A repository — it persists and queries, it doesn't create.

Events are created by `EventFactory`, then persisted by Store. This separates creation logic (validation, hashing, signing) from persistence logic (storage, indexing, retrieval).

**Concurrency contract:** Store implementations must be safe for concurrent access from multiple goroutines/threads. `Append` serialises hash chain writes internally (exclusive lock on chain head). All query methods are safe for concurrent reads.

```
Store {
    // Persistence — takes a pre-built Event from EventFactory
    Append(event Event) → Result<Event, StoreError>
    Get(id EventID) → Result<Event, StoreError>

    // Chain state — used by EventFactory to build new events
    Head() → Result<Option<Event>, StoreError>   // most recent event, None if empty

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

    // Lifecycle
    Close() → Result<void, StoreError>           // release connections, flush buffers
}

Cursor              // Opaque pagination token. Obtained from Page<T>.Cursor.
                    // Value object — validated at construction.

Page<T> {
    Items:    []T
    Cursor:   Option<Cursor>             // None = no more pages
    HasMore:  bool
}
```

**Append validates:**
- `PrevHash` matches current chain head (rejects stale events from concurrent writers)
- `Hash` is correctly computed (recomputes and compares — doesn't trust the caller)
- No duplicate `EventID` (idempotency: returns existing event if same ID already exists)
- Indexes edge-creating events (event types like `edge.created`) for edge queries

**EdgeBetween** returns `Option<Edge>` — no edge between two actors is a valid state, not an error.

Edge queries are derived from events — the Store indexes edge-creating events for efficient traversal.

### IDecisionMaker

The key abstraction. Anything that makes decisions.

```
IDecisionMaker {
    Decide(context, input DecisionInput) → Result<Decision, DecisionError>
}

DecisionInput {
    Action:     string
    Actor:      IActor
    Context:    map[string]any              // extension point — domain-specific decision context from product layer
    Causes:     NonEmpty<EventID>         // decisions must reference what caused them
}
```

An AI agent implements it. A human with a UI implements it. A committee vote implements it. A rules engine implements it. The graph records what was decided, not how.

### IIntelligence

Anything that reasons. Not every primitive needs this — most start mechanical.

```
IIntelligence {
    Reason(context, prompt string, history []Event) → Result<Response, StoreError>
}

Response {
    Content:    string
    Confidence: Score                    // [0.0, 1.0]
    TokensUsed: int                     // for cost tracking / decision tree evolution
}
```

### ITrustModel

How trust is calculated, updated, and decayed. Override for domain-specific trust mechanics.

```
ITrustModel {
    Score(context, actor IActor) → Result<TrustMetrics, StoreError>
    ScoreInDomain(context, actor IActor, domain DomainScope) → Result<TrustMetrics, StoreError>
    Update(context, actor IActor, evidence Event) → Result<TrustMetrics, StoreError>
    Decay(context, actor IActor, elapsed duration) → Result<TrustMetrics, StoreError>
    Between(context, from IActor, to IActor) → Result<TrustMetrics, StoreError>
}
```

Example override:
```
type ProjectTrust struct { trust.DefaultModel }

var (
    TaskCompleted  = MustEventType("task.completed")
    DeadlineMissed = MustEventType("deadline.missed")
    ReviewApproved = MustEventType("review.approved")
)

func (t *ProjectTrust) Update(ctx, actor, evidence) Result<TrustMetrics, StoreError> {
    switch evidence.Type {
    case TaskCompleted:
        return t.Adjust(actor, Weight(+0.05))
    case DeadlineMissed:
        return t.Adjust(actor, Weight(-0.10))
    case ReviewApproved:
        return t.Adjust(actor, Weight(+0.03))
    }
    return t.DefaultModel.Update(ctx, actor, evidence)
}
```

### IAuthorityChain

How authority is evaluated. Returns weighted authority, not binary permission.

```
IAuthorityChain {
    Evaluate(context, actor IActor, action string) → Result<AuthorityResult, DecisionError>
    Chain(context, actor IActor, action string) → Result<[]AuthorityLink, DecisionError>
    Grant(context, from IActor, to IActor, scope DomainScope, weight Score) → Result<Edge, StoreError>
    Revoke(context, from IActor, to IActor, scope DomainScope) → Result<Edge, StoreError>
}

AuthorityResult {
    Level:      AuthorityLevel
    Weight:     Score                    // [0.0, 1.0]
    Chain:      []AuthorityLink
    Delegated:  bool
    ExpiresAt:  Option<time>
}
```

### Primitive

A cognitive agent for a specific domain.

```
Primitive {
    ID() → PrimitiveID                       // encodes name + layer
    Layer() → Layer                          // [0, 13] — also derivable from ID, but explicit for clarity
    Process(tick Tick, events []Event, snapshot Frozen<Snapshot>) → Result<[]Mutation, StoreError>
    Subscriptions() → []SubscriptionPattern  // event type globs (e.g., "trust.*", "clock.*")
    Cadence() → Cadence                     // [1, ∞)
    Lifecycle() → LifecycleState
}
```

**Dependency injection:** Primitives receive dependencies at construction, not through the interface. A trust primitive is constructed with `NewTrustScorePrimitive(trustModel ITrustModel, actorStore IActorStore)`. The Primitive interface is what the tick engine sees — the constructor is what the wiring layer sees.

**Invocation rules:**
- `Process()` is called when matching events exist AND cadence allows AND lifecycle is `Active`
- `events` may be empty if the primitive is invoked on cadence alone (heartbeat-style)
- Primitives must handle empty events gracefully — return empty mutations, not an error

`snapshot` is `Frozen<Snapshot>` — a deeply immutable view. No primitive can mutate another primitive's state through the snapshot.

Mutations are declarative:
```
Mutation = AddEvent { type EventType, source ActorID, content EventContent, causes NonEmpty<EventID> }
         | AddEdge  { from ActorID, to ActorID, edgeType EdgeType, weight Weight, scope Option<DomainScope> }
         | UpdateState { key string, value any }             // primitive-internal state — typed per primitive implementation
         | UpdateActivation { level Activation }
         | UpdateLifecycle { state LifecycleState }  // must be a valid transition
```

`AddEvent` content must match the EventTypeRegistry schema for the given type. The tick engine validates this before applying.

---

## EGIP Interfaces

For inter-system communication. These are the interfaces someone building EGIP would implement.

### IIdentity

Self-sovereign identity management.

```
IIdentity {
    SystemURI() → SystemURI
    PublicKey() → PublicKey
    Sign(data []byte) → Result<Signature, ValidationError>
    Verify(publicKey PublicKey, data []byte, signature Signature) → Result<bool, ValidationError>
}
```

### ITransport

How messages move between systems. Transport-agnostic.

```
ITransport {
    Send(context, to SystemURI, envelope Envelope) → Result<Receipt, EGIPError>
    Listen(context) → <-chan Result<Envelope, EGIPError>
}
```

### Envelope

Signed message container. Immutable after construction.

```
Envelope {
    ProtocolVersion: int                 // EGIP protocol version
    ID:           EnvelopeID
    From:         SystemURI
    To:           SystemURI
    Type:         MessageType
    Payload:      MessagePayload       // interface — concrete type determined by MessageType
    Timestamp:    time
    Signature:    Signature            // Ed25519 signature of canonical form
    InReplyTo:    Option<EnvelopeID>   // None = not a reply. Some = response to this envelope.
}
```

### CGER (Cross-Graph Event Reference)

Causal links across graph boundaries. Immutable.

```
CGER {
    LocalEventID:   EventID
    RemoteSystem:   SystemURI
    RemoteEventID:  string               // opaque to us — their ID format
    RemoteHash:     Hash
    Relationship:   CGERRelationship     // CausedBy, References, RespondsTo
    Verified:       bool
}
```

### Treaty

Bilateral governance agreement. Immutable (new versions supersede old ones).

```
Treaty {
    ID:            TreatyID
    SystemA:       SystemURI
    SystemB:       SystemURI
    Terms:         NonEmpty<TreatyTerm>  // a treaty with no terms is not a treaty
    Status:        TreatyStatus
    TrustRequired: Score                 // [0.0, 1.0]
    CreatedAt:     time
    ExpiresAt:     Option<time>          // None = no expiry
    Signatures:    [2]Signature
}

TreatyTerm {
    Scope:     DomainScope
    Policy:    string               // treaty policy text — genuinely freeform legal/governance language
    Symmetric: bool
}
```

---

## Bus

Wraps a Store with pub/sub fan-out. The observer pattern — when an event is appended, all matching subscribers receive it.

```
IBus {
    // Delegates to underlying Store
    Store() → Store

    // Pub/sub
    Subscribe(pattern SubscriptionPattern, handler func(Event)) → SubscriptionID
    Unsubscribe(id SubscriptionID)

    // Lifecycle
    Close() → Result<void, StoreError>
}

SubscriptionID      // Opaque ID for an active subscription. Used to unsubscribe.
```

**Non-blocking:** Slow subscribers get dropped events, not blocked writers. The bus maintains a buffer per subscriber. If the buffer fills, oldest events are dropped and a `bus.overflow` event is emitted.

**Ordering:** Events are delivered to subscribers in hash chain order (the order they were appended to the Store). Within a single tick wave, the order is deterministic.

The tick engine uses the bus internally to distribute events to primitives. External consumers (UIs, monitoring, integrations) can also subscribe.

---

## Top-Level API

The facade — four lines to make any system auditable. Hides the 200 primitives, tick engine, trust model, and authority chain.

```
IGraph {
    // Evaluate an action through the full ontology
    Evaluate(actor IActor, action string, context map[string]any, causes NonEmpty<EventID>) → Result<Decision, DecisionError>

    // Record an event on the graph (goes through factory, validation, persistence)
    Record(type EventType, source ActorID, content EventContent, causes NonEmpty<EventID>,
           conversationID ConversationID) → Result<Event, StoreError>

    // Query the graph
    Query() → IGraphQuery

    // Access underlying components (power-user level)
    Store() → Store
    ActorStore() → IActorStore
    Bus() → IBus
    Registry() → EventTypeRegistry

    // Lifecycle
    Start() → Result<void, StoreError>           // starts tick engine, bus, primitives
    Close() → Result<void, StoreError>           // graceful shutdown
}

IGraphQuery {
    Events() → EventQuery
    Actors() → ActorQuery
    Edges() → EdgeQuery
    Trust() → TrustQuery
}

EventQuery {
    Recent(limit int) → Result<Page<Event>, StoreError>
    ByType(type EventType, limit int) → Result<Page<Event>, StoreError>
    BySource(source ActorID, limit int) → Result<Page<Event>, StoreError>
    ByConversation(id ConversationID, limit int) → Result<Page<Event>, StoreError>
    Ancestors(id EventID, maxDepth int) → Result<[]Event, StoreError>
    Descendants(id EventID, maxDepth int) → Result<[]Event, StoreError>
}

TrustQuery {
    Score(actor IActor) → Result<TrustMetrics, StoreError>
    ScoreInDomain(actor IActor, domain DomainScope) → Result<TrustMetrics, StoreError>
    Between(from IActor, to IActor) → Result<TrustMetrics, StoreError>
}
```

`IGraph.Evaluate()` is the key entry point:
1. Constructs a `DecisionInput` from the arguments
2. Runs it through `IDecisionMaker` (which uses the decision tree → intelligence fallthrough)
3. The decision tree may invoke primitives via the tick engine
4. Returns a `Decision` with full epistemic context and a cryptographic `Receipt`
5. Records the decision as an event on the graph

`IGraph.Record()` is the convenience method — it calls `EventFactory.Create()` then `Store.Append()` then notifies the bus.

---

## Canonical Form

The canonical form is the deterministic string representation used for hashing and signing. All implementations must produce the exact same bytes for the same event — any deviation breaks cross-language hash chain verification.

```
canonical(event) → string:
    version|prev_hash|id|type|source|conversation_id|timestamp_nanos|content_json

Where:
    version         → int, decimal string (e.g., "1")
    prev_hash       → Hash hex string (64 chars), or "" for Bootstrap
    id              → EventID string representation
    type            → EventType string representation
    source          → ActorID string representation
    conversation_id → ConversationID string representation
    timestamp_nanos → int64 nanoseconds since Unix epoch, decimal string
    content_json    → JSON with sorted keys, no whitespace, UTF-8 normalised (NFC)
    |               → literal pipe character (U+007C)
```

**Content JSON rules:**
- Keys sorted lexicographically (Unicode code point order)
- No whitespace between tokens
- Numbers: no trailing zeros, no leading zeros (except `0.x`), no `+` prefix
- Strings: minimal escaping (only `"`, `\`, control characters)
- UTF-8 normalised to NFC
- Null fields omitted entirely (not `"field": null`)

**Hash computation:**
```
hash = SHA-256(canonical(event))
```

**Signature computation:**
```
signature = Ed25519.Sign(private_key, canonical(event))
```

The canonical form is the **conformance-critical** part of the spec. The language-agnostic conformance test suite includes test vectors — known events with pre-computed canonical forms, hashes, and signatures that all implementations must match exactly.

---

## Versioning

### Event Version

Events carry a version for forward compatibility:

```
Event {
    Version:        int                  // schema version (starts at 1, monotonically increasing)
    // ... all other fields
}
```

When the schema for an event type changes:
1. Register the new version in `EventTypeRegistry` with a migration function
2. Old events retain their version — they are never modified (append-only)
3. Readers handle both versions (the registry provides version-aware deserialisation)
4. New events use the latest version

### Protocol Version

EGIP envelopes carry a protocol version:

```
Envelope {
    ProtocolVersion: int                 // EGIP protocol version
    // ... all other fields
}
```

Systems negotiate compatible versions during the HELLO handshake.

---

## Operational Constraints

### Rate Limiting

Actors have configurable rate limits to prevent graph flooding:

```
RateLimit {
    Actor:          ActorID
    MaxPerSecond:   int                  // max events per second
    MaxPerMinute:   int                  // max events per minute
    MaxContentSize: int                  // max content size in bytes
    MaxCauses:      int                  // max causes per event (default: 100)
}
```

Rate limit violations return `StoreError.RateLimitExceeded { actor, limit, window }`.

### Resource Lifecycle

All interfaces that hold resources implement `Close()`:

| Interface | Resources held |
|---|---|
| `Store` | Database connections, file handles |
| `IActorStore` | Database connections |
| `IBus` | Subscriber goroutines, buffers |
| `IGraph` | All of the above + tick engine |
| `ITransport` | Network connections |

`Close()` is idempotent — calling it twice is safe. After `Close()`, all other methods return `StoreError.StoreUnavailable`.

---

## Configuration

System behaviour is controlled by a configuration type — no magic numbers buried in code.

```
GraphConfig {
    // Tick engine
    MaxWavesPerTick:     int              // max ripple waves before forced quiescence (default: 10)
    TickInterval:        duration         // time between ticks (default: 100ms)

    // Authority
    RecommendedTimeout:  duration         // auto-approve timeout for Recommended level (default: 15min)

    // Bus
    SubscriberBufferSize: int             // per-subscriber event buffer (default: 1000)

    // Rate limiting
    DefaultRateLimit:    RateLimit        // applied to new actors unless overridden

    // Integrity
    VerifyOnAppend:      bool             // recompute hash on every Append (default: true, disable for perf)

    // Intelligence
    MaxTokensPerTick:    int              // budget for LLM calls per tick (cost control)
    FallbackToMechanical: bool            // if true, never escalate to IIntelligence (fully mechanical mode)
}
```

`IGraph` takes a `GraphConfig` at construction. Every magic number in the architecture has a name here.

### Token Budget and Semantic Conditions

`MaxTokensPerTick` is a shared budget across all intelligence calls in a tick — both `LeafNode.NeedsLLM` leaves and `Semantic` condition evaluations. The tick engine tracks cumulative token usage:

- Before invoking a `Semantic` condition or LLM leaf, check remaining budget
- If budget is exhausted: `Semantic` conditions fall through to their `Default` branch; LLM leaves return `DecisionOutcome.Defer`
- Token usage is recorded per-primitive for cost attribution and evolution prioritisation
- Primitives with the highest token cost are prioritised for pattern extraction (decision tree evolution targets them first)

When `FallbackToMechanical` is true, all `Semantic` conditions immediately fall through to `Default` and all `NeedsLLM` leaves return `Defer`. The system runs entirely on mechanical branches.

---

## Testing Infrastructure

### Primitive Test Harness

Testing a primitive in isolation requires a controlled environment — a Store with known events, a tick clock you control, and the ability to inspect mutations without applying them.

```
PrimitiveTestHarness {
    // Setup
    WithStore(store Store) → PrimitiveTestHarness
    WithEvents(events []Event) → PrimitiveTestHarness
    WithActors(actors []IActor) → PrimitiveTestHarness
    WithConfig(config GraphConfig) → PrimitiveTestHarness

    // Execute
    Process(primitive Primitive, events []Event) → Result<[]Mutation, StoreError>
    Tick(primitive Primitive) → Result<[]Mutation, StoreError>         // full tick cycle

    // Inspect
    Mutations() → []Mutation
    EmittedEvents() → []Event                                         // AddEvent mutations as Events
    StateChanges() → map[string]any                                   // UpdateState mutations
}
```

The harness uses `InMemoryStore` by default. It builds a `Frozen<Snapshot>` from the provided state and invokes the primitive's `Process()`. Mutations are captured, not applied — so you can assert on them without side effects.

### Conformance Test Suite

Every `Store` implementation must pass the shared conformance suite. The suite tests:

1. **Append and retrieve** — round-trip fidelity
2. **Hash chain integrity** — VerifyChain passes after N appends
3. **Causal traversal** — Ancestors/Descendants follow causes correctly
4. **Edge indexing** — edge-creating events are queryable via EdgesFrom/EdgesTo/EdgeBetween
5. **Pagination** — cursor-based pagination returns correct pages
6. **Concurrency** — concurrent Append from multiple goroutines
7. **Idempotency** — duplicate EventID returns existing event
8. **Chain head conflict** — stale PrevHash is rejected
9. **Canonical form** — pre-computed test vectors match

The suite is parameterised — pass in a `Store` factory, get a full test run.

```
func RunConformanceSuite(t *testing.T, factory func() Store)
```

---

## Interface Dependency Graph

```
Store ← (everything reads/writes events)
  ↑
IBus ← (wraps Store, delivers events to subscribers)
  ↑
IActorStore ← (everything that needs actors)
  ↑
EventTypeRegistry ← EventFactory, Store (validates content)
  ↑
EventFactory ← IGraph.Record(), tick engine (creates events)
  ↑
IActor ← ITrustModel, IAuthorityChain, IDecisionMaker
  ↑
ITrustModel ← Primitive (trust primitives call this)
  ↑
IAuthorityChain ← Primitive (authority primitives call this)
  ↑
IDecisionMaker ← IGraph.Evaluate() (calls this)
  ↑
IIntelligence ← IDecisionMaker (fallthrough when tree can't decide)
  ↑
IGraph ← Application code (the four-line entry point)
```

For EGIP:
```
IIdentity ← ITransport (signs envelopes)
  ↑
ITransport ← EGIP primitives (send/receive)
```

No circular dependencies. Each interface depends only on interfaces below it or on data structures.
