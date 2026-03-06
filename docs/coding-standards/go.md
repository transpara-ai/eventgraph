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
go/pkg/<package>/
    <package>.go       # types, interfaces, core logic
    <impl>.go          # specific implementations (e.g., postgres.go, memory.go)
    <package>_test.go  # tests
```

## Interfaces

- Define interfaces in the package that *uses* them, not the package that implements them (Go convention)
- Exception: Store, Primitive, IIntelligence, IDecisionMaker are defined centrally because they're core contracts
- Keep interfaces small — prefer many small interfaces over few large ones

## Error Handling

- Errors are values. Return them. Don't panic.
- Only panic on unrecoverable invariant violations (hash chain corruption, impossible state)
- Wrap errors with context: `fmt.Errorf("appending event: %w", err)`
- Don't swallow errors silently — if an error occurs, it must be visible (returned, logged to event graph, or both)

## Types

- Avoid `interface{}` / `any` except in `Event.Content` (which is `map[string]any` by design — typed at the primitive level)
- Use strong types for enums: `type Level string` not raw strings
- Use `time.Time` for timestamps, never Unix integers in public APIs

## Testing

- Table-driven tests preferred
- Test file lives next to the code it tests
- Name tests descriptively: `TestAppend_WithCauses_LinksCausally`
- Use `t.Helper()` in test helpers
- Use `t.Parallel()` where safe

### Coverage Thresholds

| Package | Minimum |
|---------|---------|
| event, store, bus, primitive, tick | 90% |
| protocol | 85% |
| Primitive implementations | 80% |
| Utilities | 70% |

### Conformance Tests

Store implementations must pass `store/conformance_test.go`. This tests:

- Append and retrieve
- Hash chain integrity across multiple appends
- Causal traversal (ancestors, descendants)
- Query operations (by type, by source, by conversation, since, search)
- Chain verification
- Concurrent append safety

## Concurrency

- Protect shared state with `sync.Mutex` or `sync.RWMutex`
- Document which methods are safe for concurrent use
- Use channels for event fan-out (Bus pattern)
- Non-blocking sends: prefer `select { case ch <- e: default: }` over blocking
- Test concurrent behaviour with `-race` flag

## Dependencies

- Minimise external dependencies
- Standard library preferred
- Any new dependency requires justification in the PR description
- No dependency on a specific AI provider — IIntelligence is the abstraction

## Documentation

- All public types, functions, and methods must have doc comments
- Package-level doc comment in the primary `.go` file
- Examples in `_test.go` files using `Example` functions where helpful

## Hash Chain Integrity

This is the most critical invariant. Every code path that creates events must:

1. Acquire the chain lock
2. Read the previous hash
3. Compute the new hash from canonical form
4. Store with both hash and prev_hash
5. Release the chain lock

No shortcuts. No "we'll fix the hash later." No batch appends that skip intermediate hashing. The chain is sacred.
