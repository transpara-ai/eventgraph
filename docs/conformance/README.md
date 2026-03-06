# Conformance Test Suite

Language-agnostic test vectors for EventGraph implementations. Every language package must produce identical results for these inputs.

## Files

- `canonical-vectors.json` — Test vectors for canonical form, hash computation, type validation, and state machine transitions

## What Conformance Tests Cover

1. **Canonical form** — given event fields, produce the exact canonical string
2. **Hash computation** — SHA-256 of canonical form matches expected hash
3. **Signature verification** — Ed25519 signatures verify against known keypairs
4. **Content JSON rules** — sorted keys, no whitespace, number formatting, null omission
5. **Type validation** — constrained types reject invalid values at construction
6. **State machine transitions** — valid transitions succeed, invalid transitions return errors
7. **Hash chain integrity** — appending N events produces a verifiable chain

## How to Use

Each language package loads `canonical-vectors.json` and runs tests against its implementation:

```go
// Go example
func TestConformance(t *testing.T) {
    vectors := LoadVectors("../../docs/conformance/canonical-vectors.json")
    for _, tc := range vectors.CanonicalForm.Cases {
        t.Run(tc.Name, func(t *testing.T) {
            canonical := ComputeCanonical(tc.Input)
            assert.Equal(t, tc.Expected.Canonical, canonical)
        })
    }
}
```

## Extending

When adding new event types, content schemas, or type constraints, add corresponding test vectors here. All language packages must update to pass the new vectors before the change is considered complete.

## Reference

- `docs/interfaces.md` — Canonical form specification
- `docs/interfaces.md` — Type constraints and state machines
- `ROADMAP.md` Phase 7 — Documentation and conformance status
