// Package compositions implements per-layer grammar operations as compositions
// of the 15 base social grammar operations. Each layer gets its own grammar type
// with domain-specific operations and named multi-step functions.
//
// Composition operations are the vocabulary developers use when building on a
// product graph. List, Bid, Negotiate are more useful than raw Emit, Respond,
// Channel — they carry domain intent and validate domain rules.
package compositions
