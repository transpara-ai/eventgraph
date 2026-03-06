# EGIP — EventGraph Inter-system Protocol

Sovereign systems communicate without shared infrastructure. Each system has its own event graph, its own types, its own trust model. EGIP is the anti-corruption layer that allows bilateral communication without either system importing the other's domain model.

## Design Principles

1. **Self-sovereign identity** — Ed25519 keypairs. No central registry. No certificate authority. You are your key.
2. **Bilateral, not multilateral** — two systems negotiate directly. No shared infrastructure. No consensus protocol. No blockchain.
3. **Trust, not authority** — systems earn trust through behaviour, not through certificates. Trust is asymmetric and non-transitive.
4. **Causal integrity** — cross-graph references maintain causal links. A response to a remote event references that event.
5. **Signed everything** — every envelope is signed. Every receipt is signed. Repudiation requires breaking Ed25519.

---

## Identity

### IIdentity Interface

```
IIdentity {
    SystemURI() → SystemURI                                              // this system's address
    PublicKey() → PublicKey                                               // Ed25519 public key (32 bytes)
    Sign(data []byte) → Result<Signature, ValidationError>               // Ed25519 signature
    Verify(publicKey PublicKey, data []byte, signature Signature) → Result<bool, ValidationError>
}
```

### System Identity

Each system has exactly one identity — an Ed25519 keypair generated at bootstrap:

```
SystemIdentity {
    URI:        SystemURI               // e.g., "eg://my-system.example.com"
    PublicKey:  PublicKey                // 32 bytes
    PrivateKey: [private]               // never leaves the system
    CreatedAt:  time
}
```

The `SystemURI` is the system's address. It is a value object validated at construction. Systems discover each other's URIs through DISCOVER messages or out-of-band configuration.

Identity is recorded at bootstrap:
```
"system.bootstrapped" → BootstrapContent {
    ActorID:      ActorID
    ChainGenesis: Hash
    Timestamp:    time
}
```

---

## Transport

### ITransport Interface

```
ITransport {
    Send(context, to SystemURI, envelope Envelope) → Result<Receipt, EGIPError>
    Listen(context) → <-chan Result<Envelope, EGIPError>
}
```

Transport is pluggable — HTTP, WebSocket, gRPC, message queue, carrier pigeon. The protocol doesn't care how bytes move, only what they contain and that they're signed.

---

## Envelope

The signed message container. Every EGIP communication is wrapped in an envelope.

```
Envelope {
    ProtocolVersion: int                 // EGIP protocol version (starts at 1)
    ID:           EnvelopeID             // UUID
    From:         SystemURI              // sender
    To:           SystemURI              // recipient
    Type:         MessageType            // Hello, Message, Receipt, Proof, Treaty, AuthorityRequest, Discover
    Payload:      MessagePayload         // interface — concrete type determined by MessageType
    Timestamp:    time
    Signature:    Signature              // Ed25519 signature of canonical envelope form
    InReplyTo:    Option<EnvelopeID>     // None = not a reply. Some = response to this envelope.
}
```

### Envelope Canonical Form

```
canonical(envelope) → string:
    protocol_version|id|from|to|type|timestamp_nanos|payload_json

Where:
    protocol_version → int, decimal string
    id               → EnvelopeID string
    from             → SystemURI string
    to               → SystemURI string
    type             → MessageType string (lowercase: "hello", "message", etc.)
    timestamp_nanos  → int64 nanoseconds since Unix epoch
    payload_json     → JSON with sorted keys, no whitespace, UTF-8 NFC (same rules as event content)
    |                → literal pipe character (U+007C)
```

Signature covers the canonical form. Tampering with any field invalidates the signature.

---

## Message Types

### 1. HELLO

Introduction and capability exchange. The first message between two systems.

```
HelloPayload {
    SystemURI:      SystemURI
    PublicKey:       PublicKey
    ProtocolVersions: []int             // supported protocol versions
    Capabilities:   []string            // what this system supports (e.g., "treaty", "proof", "authority")
    ChainLength:    int                 // how many events in our chain (size signal, not trust signal)
}
```

**Flow:**
1. System A sends HELLO to System B
2. System B verifies A's signature using A's provided PublicKey
3. System B responds with its own HELLO
4. Both systems store the other's identity for future communication
5. Protocol version is negotiated to the highest version both support

HELLO does not establish trust — it establishes identity. Trust starts at `Score(0.0)`.

### 2. MESSAGE

Content delivery between systems.

```
MessagePayload {
    Content:        EventContent        // typed content — same as intra-graph events
    ContentType:    EventType           // what kind of content
    ConversationID: Option<ConversationID>  // thread grouping across systems
    CGERs:          []CGER              // cross-graph event references
}
```

Messages carry typed content and optional causal links to remote events via CGERs.

### 3. RECEIPT

Delivery and processing confirmation.

```
ReceiptPayload {
    EnvelopeID:     EnvelopeID          // which envelope this receipts
    Status:         ReceiptStatus       // Delivered, Processed, Rejected
    LocalEventID:   Option<EventID>     // if processed, the local event created from this message
    Reason:         Option<string>      // if rejected, why
    Signature:      Signature           // signed by recipient — proof of delivery
}

ReceiptStatus { Delivered, Processed, Rejected }
```

Receipts are signed — they are non-repudiable proof that a system received and processed (or rejected) a message.

### 4. PROOF

Integrity proof. Demonstrates that a system's chain is valid without revealing the full chain.

```
ProofPayload {
    ProofType:      ProofType
    Data:           ProofData           // type-specific proof data
}

ProofType { ChainSegment, EventExistence, ChainSummary }

// Chain segment — contiguous portion of the hash chain
ChainSegmentProof {
    Events:     []Event                 // contiguous events with hashes
    StartHash:  Hash                    // hash before the first event
    EndHash:    Hash                    // hash after the last event
}

// Event existence — proof that a specific event exists in the chain
EventExistenceProof {
    Event:      Event                   // the event itself
    PrevHash:   Hash                    // hash before it
    NextHash:   Option<Hash>            // hash after it (None if head)
    Position:   int                     // position in chain
    ChainLength: int                    // total chain length
}

// Chain summary — high-level integrity attestation
ChainSummaryProof {
    Length:     int
    HeadHash:  Hash
    GenesisHash: Hash
    Timestamp: time                     // when this summary was computed
}
```

Proof messages build inter-system trust. A system that can produce valid chain proofs is demonstrating integrity.

### 5. TREATY

Bilateral governance agreement.

```
TreatyPayload {
    TreatyID:       TreatyID
    Action:         TreatyAction        // Propose, Accept, Modify, Suspend, Terminate
    Terms:          []TreatyTerm        // for Propose and Modify
    Reason:         Option<string>      // for Suspend and Terminate
}

TreatyAction { Propose, Accept, Modify, Suspend, Terminate }
```

**Treaty lifecycle:**

```
           Propose
    (none) ───────→ Proposed
                      │
              Accept  │  Modify (counter-proposal)
                      ↓
                    Active ←──── Modify (renegotiation)
                      │
           Suspend    │  Terminate
                      ↓
                   Suspended ──→ Active (resume) or Terminated
                      │
                      └──→ Terminated (final)
```

**Treaty terms:**

```
TreatyTerm {
    Scope:     DomainScope              // what domain this term covers
    Policy:    string                   // freeform policy text
    Symmetric: bool                     // true = applies to both systems equally
}
```

Treaties require dual signatures — both systems must sign for a treaty to become Active. Treaty violations are detectable and trigger trust updates.

### 6. AUTHORITY_REQUEST

Cross-system approval request. System A asks System B for authority to take an action that affects B.

```
AuthorityRequestPayload {
    Action:         string
    Actor:          ActorID
    Level:          AuthorityLevel
    Justification:  string
    Context:        map[string]any      // domain-specific context
    TreatyID:       Option<TreatyID>    // under which treaty this request is made
}
```

Cross-system authority follows the same three-tier model as intra-system authority. The remote system evaluates the request against its own policies and trust scores.

### 7. DISCOVER

System discovery. How systems find each other.

```
DiscoverPayload {
    Query:          DiscoverQuery
    Results:        []DiscoverResult    // for responses
}

DiscoverQuery {
    Capabilities:   []string            // looking for systems with these capabilities
    MinTrust:       Option<Score>       // only suggest systems above this trust threshold
}

DiscoverResult {
    SystemURI:      SystemURI
    PublicKey:       PublicKey
    Capabilities:   []string
    TrustScore:     Score               // the responder's trust in this system
}
```

Discovery is peer-to-peer. System A asks System B "do you know any systems that can do X?" System B responds with systems it trusts. Trust scores in discovery results are the responder's scores — they don't transfer (trust is non-transitive).

---

## Cross-Graph Event References (CGERs)

Causal links across graph boundaries. When a local event is caused by or references a remote event, the CGER records that relationship.

```
CGER {
    LocalEventID:   EventID             // our event
    RemoteSystem:   SystemURI           // their system
    RemoteEventID:  string              // their event ID (opaque to us)
    RemoteHash:     Hash                // their event's hash (for verification)
    Relationship:   CGERRelationship    // CausedBy, References, RespondsTo
    Verified:       bool                // have we verified this via PROOF?
}

CGERRelationship { CausedBy, References, RespondsTo }
```

### CGER Lifecycle

1. **Creation** — when processing a remote MESSAGE, the local system creates a local event and a CGER linking it to the remote event
2. **Verification** — the local system can request a PROOF from the remote system to verify the referenced event exists and has the claimed hash
3. **Trust impact** — verified CGERs increase inter-system trust; unverifiable CGERs decrease it

CGERs maintain causality across graph boundaries. A response to a remote message references that message via CGER, just as a response to a local event references it via `Causes`.

---

## Inter-System Trust

Trust between systems follows the `ITrustModel` interface, with system-specific defaults:

| Property | Intra-system | Inter-system |
|---|---|---|
| Initial trust | `Score(0.0)` | `Score(0.0)` |
| Decay rate | `Score(0.01)` per day | `Score(0.02)` per day (faster) |
| Max single adjustment | `Weight(0.1)` | `Weight(0.05)` (more conservative) |
| Evidence types | Events on the graph | EGIP interactions (receipts, proofs, treaty compliance) |

### Trust-Building Actions

| Action | Trust Impact |
|---|---|
| Valid PROOF provided | `+0.02` |
| Receipt delivered on time | `+0.01` |
| Treaty term honoured | `+0.03` |
| Treaty term violated | `-0.15` |
| Invalid PROOF | `-0.10` |
| Message signature invalid | `-0.20` |
| No response to HELLO | `-0.05` |

---

## Protocol Version Negotiation

HELLO messages include supported protocol versions. The agreed version is the highest version both systems support.

```
fn NegotiateVersion(local []int, remote []int) → Option<int>:
    common = intersection(local, remote)
    if common is empty:
        return None     // incompatible
    return Some(max(common))
```

If no common version exists, systems cannot communicate. This is reported as `EGIPError.TransportFailure`.

---

## Error Handling

```
EGIPError =
    | SystemNotFound { uri: SystemURI }
    | EnvelopeSignatureInvalid { envelopeID: EnvelopeID }
    | TreatyViolation { treatyID: TreatyID, term: string }
    | TrustInsufficient { system: SystemURI, score: Score, required: Score }
    | TransportFailure { to: SystemURI, reason: string }
```

All errors are typed. Transport failures are retryable. Signature failures and treaty violations are not.

---

## Security

- **Replay protection** — EnvelopeID is unique. Systems reject duplicate EnvelopeIDs.
- **Tampering detection** — every envelope is signed. Any modification invalidates the signature.
- **Impersonation prevention** — PublicKey is verified during HELLO. Subsequent messages must be signed by the established key.
- **Trust decay** — inactive systems lose trust over time. A compromised system that goes silent loses influence.
- **No shared state** — systems don't share databases, keys, or infrastructure. Compromise of one system doesn't compromise another.

---

## EGIP Events

All EGIP activity is recorded on the local graph:

```
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
"egip.trust.updated" → EGIPTrustUpdatedContent { System: SystemURI, Previous: Score, Current: Score }
```

---

## Reference

- `docs/interfaces.md` — `IIdentity`, `ITransport`, `Envelope`, `CGER`, `Treaty`, EGIP errors
- `docs/trust.md` — Trust model (applies to inter-system trust)
- `docs/authority.md` — Authority model (applies to cross-system authority requests)
- `ROADMAP.md` Phase 4 — Implementation status
