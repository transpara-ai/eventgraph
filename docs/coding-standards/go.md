# Go Coding Standards

EventGraph's Go packages are the reference implementation. All other language packages must match its behaviour.

## Version

Go 1.22 or later. Use the standard library where possible.

## Formatting and Linting

All code must pass:

```bash
go vet ./...
staticcheck ./...
golangci-lint run
```

Use `gofmt` (or `goimports`) for formatting. No exceptions.

## Package Structure

```
go/
    pkg/
        types/          # Value objects, constrained numerics, state machines
        event/          # Event, Edge, Decision, factories, content, constants
        store/          # Store interface, InMemoryStore, conformance suite
        actor/          # IActor, IActorStore, InMemoryActorStore
        bus/            # IBus interface, event bus
        graph/          # IGraph top-level API, IGraphQuery
        primitive/      # Primitive interface, registry, lifecycle
        tick/           # Tick engine, snapshot, cadence
        decision/       # Decision tree, evaluation, evolution
        trust/          # ITrustModel, default implementation
        authority/      # IAuthorityChain, policies, delegation
        protocol/       # Intra-system communication
            egip/       # Inter-system protocol
    cmd/
        eg/             # CLI
```

Within each package:

```
<package>.go            # types, interfaces, core logic
<impl>.go               # specific implementations (e.g., memory.go)
<package>_test.go       # tests
```

## Interfaces

- Define interfaces in the package that *uses* them, not the package that implements them (Go convention)
- Exception: `Store`, `Primitive`, `IIntelligence`, `IDecisionMaker` are defined centrally because they're core contracts
- Keep interfaces small — prefer many small interfaces over few large ones

---

## Generic Type Mapping

The spec uses language-agnostic generic types. This section defines their Go implementations. **All implementers must follow these patterns** — consistency across packages is non-negotiable.

### Result\<T, E\>

Go uses `(T, error)`. Our typed errors implement the `error` interface. Use `errors.As` for type-safe dispatch.

```go
// Spec: Result<Event, StoreError>
// Go:
func (s *MemoryStore) Get(id EventID) (Event, error) {
    event, ok := s.events[id]
    if !ok {
        return Event{}, &EventNotFoundError{ID: id}
    }
    return event, nil
}

// Caller uses errors.As for typed dispatch:
var notFound *EventNotFoundError
if errors.As(err, &notFound) {
    // handle specifically
}
```

**Error types** are structs implementing `error`, not bare strings:

```go
// Each domain error is a concrete struct
type EventNotFoundError struct{ ID EventID }
func (e *EventNotFoundError) Error() string { return "event not found: " + string(e.ID) }

type ChainIntegrityViolationError struct {
    Position int
    Expected Hash
    Actual   Hash
}
func (e *ChainIntegrityViolationError) Error() string { ... }

// Group error types with a sentinel interface for is-a checks:
type StoreError interface {
    error
    storeError() // unexported marker method
}

// All store errors implement StoreError:
func (e *EventNotFoundError) storeError()         {}
func (e *ChainIntegrityViolationError) storeError() {}
```

**Exhaustive dispatch** — use the Visitor pattern via interface. Adding a new error type forces compile errors in all visitors:

```go
type StoreErrorVisitor interface {
    VisitEventNotFound(*EventNotFoundError)
    VisitChainIntegrityViolation(*ChainIntegrityViolationError)
    VisitActorNotFound(*ActorNotFoundError)
    // ... one method per error type
}

// Each error type accepts the visitor
type VisitableStoreError interface {
    StoreError
    Accept(StoreErrorVisitor)
}

func (e *EventNotFoundError) Accept(v StoreErrorVisitor) { v.VisitEventNotFound(e) }
```

### Option\<T\>

Use Go generics. A custom `Option[T]` type that is explicit about presence vs absence.

```go
// Option represents an explicitly optional value.
// Zero value is None.
type Option[T any] struct {
    value T
    valid bool
}

func Some[T any](v T) Option[T] { return Option[T]{value: v, valid: true} }
func None[T any]() Option[T]    { return Option[T]{} }

func (o Option[T]) IsSome() bool  { return o.valid }
func (o Option[T]) IsNone() bool  { return !o.valid }
func (o Option[T]) Unwrap() T     { if !o.valid { panic("unwrap on None") }; return o.value }
func (o Option[T]) UnwrapOr(def T) T { if o.valid { return o.value }; return def }

// JSON marshalling: Some(v) → v, None → omitted (not null)
```

**Do NOT use pointers for optionality** on value types. `*Score` is ambiguous — is it nil because it's optional, or because something went wrong? `Option[Score]` is unambiguous.

Pointer types (like `*IActor`) naturally express optionality via `nil` — use `Option` only for value types where nil isn't available.

### NonEmpty\<T\>

A slice guaranteed to have at least one element. Construction rejects empty input.

```go
// NonEmpty is a slice guaranteed to have at least one element.
type NonEmpty[T any] struct {
    head T
    tail []T
}

func NewNonEmpty[T any](items []T) (NonEmpty[T], error) {
    if len(items) == 0 {
        return NonEmpty[T]{}, &EmptyRequiredError{Field: "NonEmpty"}
    }
    return NonEmpty[T]{head: items[0], tail: items[1:]}, nil
}

func MustNonEmpty[T any](items []T) NonEmpty[T] {
    ne, err := NewNonEmpty(items)
    if err != nil { panic(err) }
    return ne
}

func (ne NonEmpty[T]) First() T      { return ne.head }
func (ne NonEmpty[T]) All() []T      { return append([]T{ne.head}, ne.tail...) }
func (ne NonEmpty[T]) Len() int      { return 1 + len(ne.tail) }
```

### Frozen\<T\>

Go has no deep immutability. We enforce it by:
1. Deep-copying on construction
2. Exposing only getter methods (no exported fields)
3. Returning copies of slices/maps, never references

```go
// FrozenSnapshot is the read-only view passed to primitives during tick processing.
// It is deeply copied at construction — mutations to the source do not affect the snapshot.
type FrozenSnapshot struct {
    tick           Tick
    primitives     map[PrimitiveID]PrimitiveState   // deep-copied
    pendingEvents  []Event                           // deep-copied
    recentEvents   []Event                           // deep-copied
    activeActors   []IActor                          // deep-copied
}

func NewFrozenSnapshot(tick Tick, primitives map[PrimitiveID]PrimitiveState, ...) FrozenSnapshot {
    // Deep copy all slices and maps
    ...
}

// Only getter methods — no exported fields, no setters
func (s FrozenSnapshot) Tick() Tick                                { return s.tick }
func (s FrozenSnapshot) Primitives() map[PrimitiveID]PrimitiveState { return deepCopyMap(s.primitives) }
func (s FrozenSnapshot) PendingEvents() []Event                     { return copySlice(s.pendingEvents) }
```

### Page\<T\>

Cursor-based pagination. Straightforward generic struct.

```go
type Page[T any] struct {
    Items   []T
    Cursor  Option[Cursor]
    HasMore bool
}

// Cursor is opaque — implementation-specific (could be an offset, a timestamp, or an event ID).
type Cursor struct {
    value string
}

func NewCursor(v string) Cursor         { return Cursor{value: v} }
func (c Cursor) String() string         { return c.value }
```

---

## Constrained Type Pattern

All constrained types follow the same pattern. The zero value is NOT valid — you must use the constructor.

```go
// Score is constrained to [0.0, 1.0].
type Score struct{ value float64 }

func NewScore(v float64) (Score, error) {
    if v < 0.0 || v > 1.0 {
        return Score{}, &OutOfRangeError{Field: "Score", Value: v, Min: 0.0, Max: 1.0}
    }
    return Score{value: v}, nil
}

func MustScore(v float64) Score {
    s, err := NewScore(v)
    if err != nil { panic(err) }
    return s
}

func (s Score) Value() float64 { return s.value }

// Implement fmt.Stringer, encoding.TextMarshaler, json.Marshaler
```

**Pattern applies to:** `Score`, `Weight`, `Activation`, `Layer`, `Cadence`, `Tick`.

**`Must*` constructors** are for tests and known-valid literals only. Production code uses the error-returning constructors.

---

## Typed ID Pattern

All IDs are distinct types wrapping strings. You cannot pass an `ActorID` where an `EventID` is expected.

```go
type EventID struct{ value string }

func NewEventID(v string) (EventID, error) {
    if !isValidUUIDv7(v) {
        return EventID{}, &InvalidFormatError{Field: "EventID", Value: v, Expected: "UUID v7"}
    }
    return EventID{value: v}, nil
}

func (id EventID) String() string { return id.value }

// Comparable (used as map keys):
// Go structs with comparable fields are automatically comparable.
```

---

## Enum Pattern

Go has no enums. Use typed string constants with a closed set and an `IsValid()` method.

```go
type EdgeType string

const (
    EdgeTypeTrust        EdgeType = "trust"
    EdgeTypeAuthority    EdgeType = "authority"
    EdgeTypeSubscription EdgeType = "subscription"
    EdgeTypeEndorsement  EdgeType = "endorsement"
    EdgeTypeDelegation   EdgeType = "delegation"
    EdgeTypeCausation    EdgeType = "causation"
    EdgeTypeReference    EdgeType = "reference"
    EdgeTypeChannel      EdgeType = "channel"
    EdgeTypeAnnotation   EdgeType = "annotation"
)

var validEdgeTypes = map[EdgeType]bool{
    EdgeTypeTrust: true, EdgeTypeAuthority: true, EdgeTypeSubscription: true,
    EdgeTypeEndorsement: true, EdgeTypeDelegation: true, EdgeTypeCausation: true,
    EdgeTypeReference: true, EdgeTypeChannel: true, EdgeTypeAnnotation: true,
}

func (e EdgeType) IsValid() bool { return validEdgeTypes[e] }

// Exhaustive dispatch: use Visitor interface (see Error Handling above).
// DO NOT use switch with default — it silently swallows new variants.
```

For enums where exhaustive matching matters most (`DecisionOutcome`, `AuthorityLevel`, `LifecycleState`), always provide a Visitor interface. For simpler enums used mostly in data (`SeverityLevel`, `ExpectationStatus`), the `IsValid()` method is sufficient.

---

## State Machine Pattern

State machines enforce valid transitions at the type level.

```go
type LifecycleState string

const (
    LifecycleDormant      LifecycleState = "dormant"
    LifecycleActivating   LifecycleState = "activating"
    LifecycleActive       LifecycleState = "active"
    LifecycleProcessing   LifecycleState = "processing"
    LifecycleEmitting     LifecycleState = "emitting"
    LifecycleDeactivating LifecycleState = "deactivating"
)

var validLifecycleTransitions = map[LifecycleState][]LifecycleState{
    LifecycleDormant:      {LifecycleActivating},
    LifecycleActivating:   {LifecycleActive, LifecycleDormant},
    LifecycleActive:       {LifecycleProcessing, LifecycleDeactivating},
    LifecycleProcessing:   {LifecycleEmitting, LifecycleActive},
    LifecycleEmitting:     {LifecycleActive},
    LifecycleDeactivating: {LifecycleDormant},
}

func (s LifecycleState) TransitionTo(target LifecycleState) (LifecycleState, error) {
    valid := validLifecycleTransitions[s]
    for _, v := range valid {
        if v == target {
            return target, nil
        }
    }
    return s, &InvalidLifecycleTransitionError{
        From:         s,
        To:           target,
        ValidTargets: valid,
    }
}

func (s LifecycleState) ValidTransitions() []LifecycleState {
    return validLifecycleTransitions[s]
}
```

---

## Immutable Struct Pattern

Domain objects are immutable after construction. Use unexported fields + getter methods.

```go
type Event struct {
    version        int
    id             EventID
    eventType      EventType
    timestamp      time.Time
    source         ActorID
    content        EventContent
    causes         NonEmpty[EventID]
    conversationID ConversationID
    hash           Hash
    prevHash       Hash
    signature      Signature
}

// Only getter methods — no setters, no exported fields
func (e Event) Version() int                     { return e.version }
func (e Event) ID() EventID                      { return e.id }
func (e Event) Type() EventType                  { return e.eventType }
func (e Event) Timestamp() time.Time             { return e.timestamp }
func (e Event) Source() ActorID                   { return e.source }
func (e Event) Content() EventContent            { return e.content }
func (e Event) Causes() NonEmpty[EventID]        { return e.causes }
func (e Event) ConversationID() ConversationID   { return e.conversationID }
func (e Event) Hash() Hash                       { return e.hash }
func (e Event) PrevHash() Hash                   { return e.prevHash }
func (e Event) Signature() Signature             { return e.signature }

// Constructed only via EventFactory — no public constructor
```

---

## EventContent Interface

`EventContent` is a Go interface. Each event type has a concrete struct implementing it. The `EventTypeRegistry` maps `EventType` to the expected concrete type.

```go
// EventContent is the interface all typed event content implements.
type EventContent interface {
    EventType() EventType       // which event type this content is for
    Accept(EventContentVisitor) // visitor for exhaustive dispatch
}

// Concrete content types:
type TrustUpdatedContent struct {
    actor    ActorID
    previous Score
    current  Score
    domain   DomainScope
    cause    EventID
}

func (c TrustUpdatedContent) EventType() EventType          { return MustEventType("trust.updated") }
func (c TrustUpdatedContent) Accept(v EventContentVisitor)  { v.VisitTrustUpdated(c) }

// Getters for each field...
```

---

## Error Handling

- Errors are values. Return them. Don't panic.
- Only panic on unrecoverable invariant violations (hash chain corruption, impossible state)
- Wrap errors with context: `fmt.Errorf("appending event: %w", err)`
- Don't swallow errors silently — if an error occurs, it must be visible (returned, logged to event graph, or both)
- No `interface{}` / `any` — `Event.Content` is typed per event type via `EventTypeRegistry`. No `map[string]any` for event content.
- Use strong types for enums: `type EdgeType string` not raw strings
- Use `time.Time` for timestamps, never Unix integers in public APIs

---

## JSON Serialisation

Events must serialise to JSON for canonical form computation and storage. Follow these rules:

```go
// Constrained types marshal to their underlying value:
func (s Score) MarshalJSON() ([]byte, error)  { return json.Marshal(s.value) }
func (s *Score) UnmarshalJSON(b []byte) error {
    var v float64
    if err := json.Unmarshal(b, &v); err != nil { return err }
    score, err := NewScore(v)
    if err != nil { return err }
    *s = score
    return nil
}

// Typed IDs marshal to their string value:
func (id EventID) MarshalJSON() ([]byte, error)  { return json.Marshal(id.value) }

// Option marshals to value or is omitted:
func (o Option[T]) MarshalJSON() ([]byte, error) {
    if o.IsNone() { return []byte("null"), nil }
    return json.Marshal(o.value)
}
// Use `omitempty` tags on struct fields with Option types.

// EventContent uses a type discriminator for deserialisation:
// { "type": "trust.updated", "data": { ... } }
```

**Canonical form JSON** (for hashing) uses sorted keys, no whitespace, and NFC normalisation. This is separate from storage JSON — use a dedicated `CanonicalJSON(v any) []byte` function.

---

## Testing

- Table-driven tests preferred
- Test file lives next to the code it tests
- Name tests descriptively: `TestAppend_WithCauses_LinksCausally`
- Use `t.Helper()` in test helpers
- Use `t.Parallel()` where safe
- Run with `-race` flag always

### Coverage Thresholds

| Package | Minimum |
|---------|---------|
| types, event, store, bus, primitive, tick | 90% |
| protocol, egip | 85% |
| Primitive implementations | 80% |
| Utilities and helpers | 70% |

### Conformance Tests

Store implementations must pass `store/conformance_test.go`. Run with:

```go
func TestMemoryStoreConformance(t *testing.T) {
    RunConformanceSuite(t, func() Store { return NewMemoryStore() })
}
```

The suite tests: append/retrieve, hash chain integrity, causal traversal, queries, edge indexing, pagination, concurrency, idempotency, chain head conflicts, canonical form vectors.

---

## Concurrency

- Protect shared state with `sync.Mutex` or `sync.RWMutex`
- Document which methods are safe for concurrent use
- Use channels for event fan-out (Bus pattern)
- Non-blocking sends: prefer `select { case ch <- e: default: }` over blocking
- Test concurrent behaviour with `-race` flag
- Store.Append serialises via chain head lock — document this

---

## Dependencies

- Minimise external dependencies
- Standard library preferred
- Any new dependency requires justification in the PR description
- No dependency on a specific AI provider — `IIntelligence` is the abstraction
- Allowed: `github.com/google/uuid` for UUID v7 generation

---

## Documentation

- All public types, functions, and methods must have doc comments
- Package-level doc comment in the primary `.go` file
- Examples in `_test.go` files using `Example` functions where helpful

---

## Hash Chain Integrity

This is the most critical invariant. Every code path that creates events must:

1. Acquire the chain lock
2. Read the previous hash
3. Compute the new hash from canonical form
4. Store with both hash and prev_hash
5. Release the chain lock

No shortcuts. No "we'll fix the hash later." No batch appends that skip intermediate hashing. The chain is sacred.
