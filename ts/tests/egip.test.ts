import { describe, it, expect } from "vitest";
import {
  // Constants
  MessageType, TreatyStatus, TreatyAction, ReceiptStatus, ProofType,
  CGERRelationship, AuthorityLevel, CurrentProtocolVersion,
  MaxEnvelopeAgeMs, isValidMessageType,
  // Trust constants
  TrustImpactValidProof, TrustImpactReceiptOnTime, TrustImpactTreatyHonoured,
  TrustImpactSignatureInvalid, TrustImpactNoHelloResponse,
  InterSystemDecayRate, InterSystemMaxAdjustment,
  // Errors
  EGIPError, SystemNotFoundError, EnvelopeSignatureInvalidError,
  TreatyViolationError, TrustInsufficientError, TransportFailureError,
  DuplicateEnvelopeError, TreatyNotFoundError, VersionIncompatibleError,
  // Identity
  SystemIdentity,
  // Envelope
  Envelope, signEnvelope, verifyEnvelope,
  // Payloads
  helloPayload, messagePayload, receiptPayload, proofPayload,
  treatyPayload, authorityRequestPayload, discoverPayload,
  // Treaty
  Treaty,
  // Stores
  PeerStore, TreatyStore, EnvelopeDedup,
  // Proof
  verifyChainSegment, verifyEventExistence, validateProof, proofTypeFromData,
  // Version negotiation
  negotiateVersion,
  // Handler
  Handler,
  // Types re-exported
  type HelloPayload, type MessagePayloadContent, type ReceiptPayload,
  type ProofPayload, type TreatyPayload, type AuthorityRequestPayload,
  type DiscoverPayload, type CGER, type TreatyTerm, type ProofEvent,
  type ChainSegmentProof, type EventExistenceProof, type ChainSummaryProof,
  type PeerRecord, type IIdentity, type ITransport, type IncomingEnvelope,
} from "../src/egip.js";
import {
  Option, Score, Weight, DomainScope, EnvelopeId, TreatyId, SystemUri,
  PublicKey, Signature, Hash, EventId, EventType, ActorId, ConversationId,
} from "../src/types.js";
import { InvalidTransitionError } from "../src/errors.js";

// ── Test helpers ────────────────────────────────────────────────────────

function makeIdentity(name: string): SystemIdentity {
  return SystemIdentity.generate(new SystemUri(`eg://${name}`));
}

function makeEnvelopeId(): EnvelopeId {
  // Generate a random UUID v4
  const hex = Array.from({ length: 32 }, () => Math.floor(Math.random() * 16).toString(16)).join("");
  const uuid = `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20, 32)}`;
  return new EnvelopeId(uuid);
}

function makeTreatyId(): TreatyId {
  const hex = Array.from({ length: 32 }, () => Math.floor(Math.random() * 16).toString(16)).join("");
  const uuid = `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20, 32)}`;
  return new TreatyId(uuid);
}

function makeHash(seed: string): Hash {
  return new Hash(seed.padEnd(64, "0").slice(0, 64));
}

// ── Constants ───────────────────────────────────────────────────────────

describe("MessageType", () => {
  it("has all seven types", () => {
    expect(MessageType.Hello).toBe("Hello");
    expect(MessageType.Message).toBe("Message");
    expect(MessageType.Receipt).toBe("Receipt");
    expect(MessageType.Proof).toBe("Proof");
    expect(MessageType.Treaty).toBe("Treaty");
    expect(MessageType.AuthorityRequest).toBe("AuthorityRequest");
    expect(MessageType.Discover).toBe("Discover");
  });

  it("validates message types", () => {
    expect(isValidMessageType("Hello")).toBe(true);
    expect(isValidMessageType("bogus")).toBe(false);
  });
});

describe("TreatyStatus", () => {
  it("has all four statuses", () => {
    expect(TreatyStatus.Proposed).toBe("Proposed");
    expect(TreatyStatus.Active).toBe("Active");
    expect(TreatyStatus.Suspended).toBe("Suspended");
    expect(TreatyStatus.Terminated).toBe("Terminated");
  });
});

describe("TreatyAction", () => {
  it("has all five actions", () => {
    expect(TreatyAction.Propose).toBe("Propose");
    expect(TreatyAction.Accept).toBe("Accept");
    expect(TreatyAction.Modify).toBe("Modify");
    expect(TreatyAction.Suspend).toBe("Suspend");
    expect(TreatyAction.Terminate).toBe("Terminate");
  });
});

describe("ReceiptStatus", () => {
  it("has all three statuses", () => {
    expect(ReceiptStatus.Delivered).toBe("Delivered");
    expect(ReceiptStatus.Processed).toBe("Processed");
    expect(ReceiptStatus.Rejected).toBe("Rejected");
  });
});

describe("ProofType", () => {
  it("has all three types", () => {
    expect(ProofType.ChainSegment).toBe("ChainSegment");
    expect(ProofType.EventExistence).toBe("EventExistence");
    expect(ProofType.ChainSummary).toBe("ChainSummary");
  });
});

describe("CGERRelationship", () => {
  it("has all three relationships", () => {
    expect(CGERRelationship.CausedBy).toBe("CausedBy");
    expect(CGERRelationship.References).toBe("References");
    expect(CGERRelationship.RespondsTo).toBe("RespondsTo");
  });
});

describe("AuthorityLevel", () => {
  it("has all three levels", () => {
    expect(AuthorityLevel.Required).toBe("Required");
    expect(AuthorityLevel.Recommended).toBe("Recommended");
    expect(AuthorityLevel.Notification).toBe("Notification");
  });
});

// ── EGIP Errors ─────────────────────────────────────────────────────────

describe("EGIP Errors", () => {
  it("SystemNotFoundError", () => {
    const uri = new SystemUri("eg://test");
    const err = new SystemNotFoundError(uri);
    expect(err).toBeInstanceOf(EGIPError);
    expect(err.message).toContain("system not found");
    expect(err.uri).toBe(uri);
  });

  it("EnvelopeSignatureInvalidError", () => {
    const id = makeEnvelopeId();
    const err = new EnvelopeSignatureInvalidError(id);
    expect(err).toBeInstanceOf(EGIPError);
    expect(err.message).toContain("envelope signature invalid");
  });

  it("TreatyViolationError", () => {
    const id = makeTreatyId();
    const err = new TreatyViolationError(id, "rate_limit");
    expect(err).toBeInstanceOf(EGIPError);
    expect(err.message).toContain("treaty");
    expect(err.message).toContain("violated");
  });

  it("TrustInsufficientError", () => {
    const uri = new SystemUri("eg://test");
    const err = new TrustInsufficientError(uri, new Score(0.2), new Score(0.5));
    expect(err).toBeInstanceOf(EGIPError);
    expect(err.message).toContain("trust insufficient");
  });

  it("TransportFailureError", () => {
    const uri = new SystemUri("eg://test");
    const err = new TransportFailureError(uri, "connection refused");
    expect(err).toBeInstanceOf(EGIPError);
    expect(err.message).toContain("transport failure");
  });

  it("DuplicateEnvelopeError", () => {
    const id = makeEnvelopeId();
    const err = new DuplicateEnvelopeError(id);
    expect(err).toBeInstanceOf(EGIPError);
    expect(err.message).toContain("duplicate envelope");
  });

  it("TreatyNotFoundError", () => {
    const id = makeTreatyId();
    const err = new TreatyNotFoundError(id);
    expect(err).toBeInstanceOf(EGIPError);
    expect(err.message).toContain("treaty not found");
  });

  it("VersionIncompatibleError", () => {
    const err = new VersionIncompatibleError([1], [2]);
    expect(err).toBeInstanceOf(EGIPError);
    expect(err.message).toContain("no compatible protocol version");
  });
});

// ── SystemIdentity ──────────────────────────────────────────────────────

describe("SystemIdentity", () => {
  it("generates a valid identity", () => {
    const id = makeIdentity("system-a");
    expect(id.systemUri().value).toBe("eg://system-a");
    expect(id.publicKey().bytes.length).toBe(32);
  });

  it("signs and verifies data", () => {
    const id = makeIdentity("system-a");
    const data = new TextEncoder().encode("hello world");
    const sig = id.sign(data);
    expect(sig.bytes.length).toBe(64);
    expect(id.verify(id.publicKey(), data, sig)).toBe(true);
  });

  it("rejects invalid signatures", () => {
    const id = makeIdentity("system-a");
    const data = new TextEncoder().encode("hello world");
    const sig = id.sign(data);
    const wrongData = new TextEncoder().encode("wrong data");
    expect(id.verify(id.publicKey(), wrongData, sig)).toBe(false);
  });

  it("rejects signatures from different keys", () => {
    const idA = makeIdentity("system-a");
    const idB = makeIdentity("system-b");
    const data = new TextEncoder().encode("hello");
    const sig = idA.sign(data);
    expect(idA.verify(idB.publicKey(), data, sig)).toBe(false);
  });
});

// ── Envelope ────────────────────────────────────────────────────────────

describe("Envelope", () => {
  it("computes canonical form deterministically", () => {
    const id = makeIdentity("system-a");
    const envId = makeEnvelopeId();
    const payload = helloPayload({
      systemUri: id.systemUri(),
      publicKey: id.publicKey(),
      protocolVersions: [1],
      capabilities: ["treaty"],
      chainLength: 42,
    });

    const ts = new Date("2025-01-01T00:00:00Z");
    const env = new Envelope(
      1, envId, id.systemUri(), new SystemUri("eg://system-b"),
      MessageType.Hello, payload, ts,
      new Signature(new Uint8Array(64)), Option.none<EnvelopeId>(),
    );

    const form1 = env.canonicalForm();
    const form2 = env.canonicalForm();
    expect(form1).toBe(form2);
    expect(form1).toContain("hello"); // message type lowercased
    expect(form1).toContain(envId.value);
  });

  it("sign and verify round-trip", () => {
    const id = makeIdentity("system-a");
    const envId = makeEnvelopeId();
    const payload = helloPayload({
      systemUri: id.systemUri(),
      publicKey: id.publicKey(),
      protocolVersions: [1],
      capabilities: [],
      chainLength: 0,
    });

    const env = new Envelope(
      1, envId, id.systemUri(), new SystemUri("eg://system-b"),
      MessageType.Hello, payload, new Date(),
      new Signature(new Uint8Array(64)), Option.none<EnvelopeId>(),
    );

    const signed = signEnvelope(env, id);
    expect(verifyEnvelope(signed, id, id.publicKey())).toBe(true);
  });

  it("verification fails with wrong key", () => {
    const idA = makeIdentity("system-a");
    const idB = makeIdentity("system-b");
    const envId = makeEnvelopeId();
    const payload = helloPayload({
      systemUri: idA.systemUri(),
      publicKey: idA.publicKey(),
      protocolVersions: [1],
      capabilities: [],
      chainLength: 0,
    });

    const env = new Envelope(
      1, envId, idA.systemUri(), new SystemUri("eg://system-b"),
      MessageType.Hello, payload, new Date(),
      new Signature(new Uint8Array(64)), Option.none<EnvelopeId>(),
    );

    const signed = signEnvelope(env, idA);
    expect(verifyEnvelope(signed, idA, idB.publicKey())).toBe(false);
  });

  it("includes inReplyTo in canonical form", () => {
    const id = makeIdentity("system-a");
    const envId = makeEnvelopeId();
    const replyToId = makeEnvelopeId();
    const payload = helloPayload({
      systemUri: id.systemUri(),
      publicKey: id.publicKey(),
      protocolVersions: [1],
      capabilities: [],
      chainLength: 0,
    });

    const env = new Envelope(
      1, envId, id.systemUri(), new SystemUri("eg://system-b"),
      MessageType.Hello, payload, new Date(),
      new Signature(new Uint8Array(64)), Option.some(replyToId),
    );

    expect(env.canonicalForm()).toContain(replyToId.value);
  });
});

// ── Version negotiation ─────────────────────────────────────────────────

describe("negotiateVersion", () => {
  it("finds highest common version", () => {
    const result = negotiateVersion([1, 2, 3], [2, 3, 4]);
    expect(result.isSome).toBe(true);
    expect(result.unwrap()).toBe(3);
  });

  it("returns none when no common version", () => {
    const result = negotiateVersion([1, 2], [3, 4]);
    expect(result.isNone).toBe(true);
  });

  it("handles single common version", () => {
    const result = negotiateVersion([1], [1]);
    expect(result.isSome).toBe(true);
    expect(result.unwrap()).toBe(1);
  });

  it("handles empty arrays", () => {
    expect(negotiateVersion([], [1]).isNone).toBe(true);
    expect(negotiateVersion([1], []).isNone).toBe(true);
  });
});

// ── Treaty state machine ────────────────────────────────────────────────

describe("Treaty", () => {
  it("starts in Proposed status", () => {
    const treaty = new Treaty(makeTreatyId(), new SystemUri("eg://a"), new SystemUri("eg://b"), []);
    expect(treaty.status).toBe(TreatyStatus.Proposed);
  });

  it("transitions Proposed -> Active", () => {
    const treaty = new Treaty(makeTreatyId(), new SystemUri("eg://a"), new SystemUri("eg://b"), []);
    treaty.transition(TreatyStatus.Active);
    expect(treaty.status).toBe(TreatyStatus.Active);
  });

  it("transitions Active -> Suspended", () => {
    const treaty = new Treaty(makeTreatyId(), new SystemUri("eg://a"), new SystemUri("eg://b"), []);
    treaty.transition(TreatyStatus.Active);
    treaty.transition(TreatyStatus.Suspended);
    expect(treaty.status).toBe(TreatyStatus.Suspended);
  });

  it("transitions Suspended -> Active", () => {
    const treaty = new Treaty(makeTreatyId(), new SystemUri("eg://a"), new SystemUri("eg://b"), []);
    treaty.transition(TreatyStatus.Active);
    treaty.transition(TreatyStatus.Suspended);
    treaty.transition(TreatyStatus.Active);
    expect(treaty.status).toBe(TreatyStatus.Active);
  });

  it("transitions to Terminated from any non-terminal state", () => {
    for (const from of [TreatyStatus.Proposed, TreatyStatus.Active, TreatyStatus.Suspended] as TreatyStatus[]) {
      const treaty = new Treaty(makeTreatyId(), new SystemUri("eg://a"), new SystemUri("eg://b"), []);
      if (from === TreatyStatus.Active) treaty.transition(TreatyStatus.Active);
      if (from === TreatyStatus.Suspended) {
        treaty.transition(TreatyStatus.Active);
        treaty.transition(TreatyStatus.Suspended);
      }
      treaty.transition(TreatyStatus.Terminated);
      expect(treaty.status).toBe(TreatyStatus.Terminated);
    }
  });

  it("rejects invalid transitions", () => {
    const treaty = new Treaty(makeTreatyId(), new SystemUri("eg://a"), new SystemUri("eg://b"), []);
    // Proposed -> Suspended is invalid
    expect(() => treaty.transition(TreatyStatus.Suspended)).toThrow(InvalidTransitionError);
  });

  it("rejects transitions from Terminated", () => {
    const treaty = new Treaty(makeTreatyId(), new SystemUri("eg://a"), new SystemUri("eg://b"), []);
    treaty.transition(TreatyStatus.Terminated);
    expect(() => treaty.transition(TreatyStatus.Active)).toThrow(InvalidTransitionError);
  });

  it("applyAction Accept transitions to Active", () => {
    const treaty = new Treaty(makeTreatyId(), new SystemUri("eg://a"), new SystemUri("eg://b"), []);
    treaty.applyAction(TreatyAction.Accept);
    expect(treaty.status).toBe(TreatyStatus.Active);
  });

  it("applyAction Modify only works on Active", () => {
    const treaty = new Treaty(makeTreatyId(), new SystemUri("eg://a"), new SystemUri("eg://b"), []);
    expect(() => treaty.applyAction(TreatyAction.Modify)).toThrow(EGIPError);
    treaty.applyAction(TreatyAction.Accept);
    treaty.applyAction(TreatyAction.Modify); // should not throw
    expect(treaty.status).toBe(TreatyStatus.Active);
  });

  it("applyAction Propose throws on existing treaty", () => {
    const treaty = new Treaty(makeTreatyId(), new SystemUri("eg://a"), new SystemUri("eg://b"), []);
    expect(() => treaty.applyAction(TreatyAction.Propose)).toThrow(EGIPError);
  });

  it("stores terms", () => {
    const terms: TreatyTerm[] = [
      { scope: new DomainScope("data.sharing"), policy: "allow_read", symmetric: true },
    ];
    const treaty = new Treaty(makeTreatyId(), new SystemUri("eg://a"), new SystemUri("eg://b"), terms);
    expect(treaty.terms).toHaveLength(1);
    expect(treaty.terms[0].scope.value).toBe("data.sharing");
  });
});

// ── PeerStore ───────────────────────────────────────────────────────────

describe("PeerStore", () => {
  it("registers and retrieves peers", () => {
    const store = new PeerStore();
    const uri = new SystemUri("eg://peer-a");
    const key = new PublicKey(new Uint8Array(32));

    store.register(uri, key, ["treaty"], 1);
    const [record, found] = store.get(uri);
    expect(found).toBe(true);
    expect(record.systemUri.value).toBe("eg://peer-a");
    expect(record.trust.value).toBe(0.0);
  });

  it("does not overwrite public key on re-registration", () => {
    const store = new PeerStore();
    const uri = new SystemUri("eg://peer-a");
    const key1 = new PublicKey(new Uint8Array(32).fill(1));
    const key2 = new PublicKey(new Uint8Array(32).fill(2));

    store.register(uri, key1, ["treaty"], 1);
    store.register(uri, key2, ["proof"], 1);
    const [record] = store.get(uri);
    expect(Buffer.from(record.publicKey.bytes).equals(Buffer.from(key1.bytes))).toBe(true);
    expect(record.capabilities).toEqual(["proof"]); // capabilities updated
  });

  it("updates trust within bounds", () => {
    const store = new PeerStore();
    const uri = new SystemUri("eg://peer-a");
    store.register(uri, new PublicKey(new Uint8Array(32)), [], 1);

    // Positive delta capped at InterSystemMaxAdjustment
    const [score1] = store.updateTrust(uri, 0.5);
    expect(score1.value).toBe(InterSystemMaxAdjustment.value);

    // Negative delta applied directly
    const [score2] = store.updateTrust(uri, -0.20);
    expect(score2.value).toBe(0.0); // clamped at 0
  });

  it("returns false for unknown peers", () => {
    const store = new PeerStore();
    const [, found] = store.get(new SystemUri("eg://unknown"));
    expect(found).toBe(false);
  });

  it("decays all peers", () => {
    const store = new PeerStore();
    const uri = new SystemUri("eg://peer-a");
    store.register(uri, new PublicKey(new Uint8Array(32)), [], 1);
    store.updateTrust(uri, 0.05);

    // Manually set lastDecayedAt to 10 days ago
    const [record] = store.get(uri);
    // Access internal state for testing
    const internalRecord = (store as any).peers.get(uri.value) as PeerRecord;
    internalRecord.lastDecayedAt = new Date(Date.now() - 10 * 24 * 60 * 60 * 1000);

    store.decayAll();
    const [decayed] = store.get(uri);
    expect(decayed.trust.value).toBeLessThan(record.trust.value);
  });

  it("returns all peers", () => {
    const store = new PeerStore();
    store.register(new SystemUri("eg://a"), new PublicKey(new Uint8Array(32)), [], 1);
    store.register(new SystemUri("eg://b"), new PublicKey(new Uint8Array(32)), [], 1);
    expect(store.all()).toHaveLength(2);
  });
});

// ── TreatyStore ─────────────────────────────────────────────────────────

describe("TreatyStore", () => {
  it("stores and retrieves treaties", () => {
    const store = new TreatyStore();
    const id = makeTreatyId();
    const treaty = new Treaty(id, new SystemUri("eg://a"), new SystemUri("eg://b"), []);
    store.put(treaty);

    const [retrieved, found] = store.get(id);
    expect(found).toBe(true);
    expect(retrieved!.id.value).toBe(id.value);
  });

  it("applies mutations", () => {
    const store = new TreatyStore();
    const id = makeTreatyId();
    const treaty = new Treaty(id, new SystemUri("eg://a"), new SystemUri("eg://b"), []);
    store.put(treaty);

    store.apply(id, (t) => {
      t.applyAction(TreatyAction.Accept);
    });

    const [updated] = store.get(id);
    expect(updated!.status).toBe(TreatyStatus.Active);
  });

  it("throws TreatyNotFoundError for missing treaties", () => {
    const store = new TreatyStore();
    expect(() => store.apply(makeTreatyId(), () => {})).toThrow(TreatyNotFoundError);
  });

  it("finds treaties by system", () => {
    const store = new TreatyStore();
    const uriA = new SystemUri("eg://a");
    const uriB = new SystemUri("eg://b");
    const uriC = new SystemUri("eg://c");

    store.put(new Treaty(makeTreatyId(), uriA, uriB, []));
    store.put(new Treaty(makeTreatyId(), uriB, uriC, []));

    expect(store.bySystem(uriA)).toHaveLength(1);
    expect(store.bySystem(uriB)).toHaveLength(2);
    expect(store.bySystem(uriC)).toHaveLength(1);
  });

  it("returns active treaties", () => {
    const store = new TreatyStore();
    const t1 = new Treaty(makeTreatyId(), new SystemUri("eg://a"), new SystemUri("eg://b"), []);
    const t2 = new Treaty(makeTreatyId(), new SystemUri("eg://a"), new SystemUri("eg://c"), []);
    t1.applyAction(TreatyAction.Accept);
    store.put(t1);
    store.put(t2);

    const active = store.active();
    expect(active).toHaveLength(1);
    expect(active[0].status).toBe(TreatyStatus.Active);
  });
});

// ── EnvelopeDedup ───────────────────────────────────────────────────────

describe("EnvelopeDedup", () => {
  it("accepts first occurrence", () => {
    const dedup = new EnvelopeDedup();
    const id = makeEnvelopeId();
    expect(dedup.check(id)).toBe(true);
  });

  it("rejects duplicates", () => {
    const dedup = new EnvelopeDedup();
    const id = makeEnvelopeId();
    expect(dedup.check(id)).toBe(true);
    expect(dedup.check(id)).toBe(false);
  });

  it("tracks size", () => {
    const dedup = new EnvelopeDedup();
    dedup.check(makeEnvelopeId());
    dedup.check(makeEnvelopeId());
    expect(dedup.size).toBe(2);
  });

  it("prunes expired entries", () => {
    const dedup = new EnvelopeDedup(0); // 0ms TTL
    const id = makeEnvelopeId();
    dedup.check(id);
    // Wait a moment for expiry
    const removed = dedup.prune();
    expect(removed).toBe(1);
    expect(dedup.size).toBe(0);
    // Can now accept the same ID again
    expect(dedup.check(id)).toBe(true);
  });
});

// ── Proof verification ──────────────────────────────────────────────────

describe("Proof verification", () => {
  const h1 = makeHash("a".repeat(64));
  const h2 = makeHash("b".repeat(64));
  const h3 = makeHash("c".repeat(64));

  describe("verifyChainSegment", () => {
    it("returns false for empty events", () => {
      const proof: ChainSegmentProof = {
        proofKind: "chain_segment",
        events: [],
        startHash: h1,
        endHash: h2,
      };
      expect(verifyChainSegment(proof)).toBe(false);
    });

    it("verifies valid single-event segment", () => {
      const proof: ChainSegmentProof = {
        proofKind: "chain_segment",
        events: [{ hash: h2, prevHash: h1 }],
        startHash: h1,
        endHash: h2,
      };
      expect(verifyChainSegment(proof)).toBe(true);
    });

    it("verifies valid multi-event segment", () => {
      const proof: ChainSegmentProof = {
        proofKind: "chain_segment",
        events: [
          { hash: h2, prevHash: h1 },
          { hash: h3, prevHash: h2 },
        ],
        startHash: h1,
        endHash: h3,
      };
      expect(verifyChainSegment(proof)).toBe(true);
    });

    it("rejects broken chain", () => {
      const proof: ChainSegmentProof = {
        proofKind: "chain_segment",
        events: [
          { hash: h2, prevHash: h1 },
          { hash: h3, prevHash: h1 }, // should be h2
        ],
        startHash: h1,
        endHash: h3,
      };
      expect(verifyChainSegment(proof)).toBe(false);
    });

    it("rejects mismatched startHash", () => {
      const proof: ChainSegmentProof = {
        proofKind: "chain_segment",
        events: [{ hash: h2, prevHash: h1 }],
        startHash: h3, // wrong
        endHash: h2,
      };
      expect(verifyChainSegment(proof)).toBe(false);
    });

    it("rejects mismatched endHash", () => {
      const proof: ChainSegmentProof = {
        proofKind: "chain_segment",
        events: [{ hash: h2, prevHash: h1 }],
        startHash: h1,
        endHash: h3, // wrong
      };
      expect(verifyChainSegment(proof)).toBe(false);
    });
  });

  describe("verifyEventExistence", () => {
    it("verifies valid proof", () => {
      const proof: EventExistenceProof = {
        proofKind: "event_existence",
        event: { hash: h2, prevHash: h1 },
        prevHash: h1,
        nextHash: Option.none<Hash>(),
        position: 0,
        chainLength: 1,
      };
      expect(verifyEventExistence(proof)).toBe(true);
    });

    it("rejects mismatched prevHash", () => {
      const proof: EventExistenceProof = {
        proofKind: "event_existence",
        event: { hash: h2, prevHash: h1 },
        prevHash: h3, // doesn't match event.prevHash
        nextHash: Option.none<Hash>(),
        position: 0,
        chainLength: 1,
      };
      expect(verifyEventExistence(proof)).toBe(false);
    });

    it("rejects position >= chainLength", () => {
      const proof: EventExistenceProof = {
        proofKind: "event_existence",
        event: { hash: h2, prevHash: h1 },
        prevHash: h1,
        nextHash: Option.none<Hash>(),
        position: 5,
        chainLength: 5,
      };
      expect(verifyEventExistence(proof)).toBe(false);
    });

    it("rejects zero hash", () => {
      const proof: EventExistenceProof = {
        proofKind: "event_existence",
        event: { hash: Hash.zero(), prevHash: h1 },
        prevHash: h1,
        nextHash: Option.none<Hash>(),
        position: 0,
        chainLength: 1,
      };
      expect(verifyEventExistence(proof)).toBe(false);
    });
  });

  describe("validateProof", () => {
    it("validates chain summary with positive length", () => {
      const payload: ProofPayload = {
        kind: "proof",
        proofType: ProofType.ChainSummary,
        data: {
          proofKind: "chain_summary",
          length: 42,
          headHash: h1,
          genesisHash: h2,
          timestamp: new Date(),
        },
      };
      expect(validateProof(payload)).toBe(true);
    });

    it("rejects chain summary with zero length", () => {
      const payload: ProofPayload = {
        kind: "proof",
        proofType: ProofType.ChainSummary,
        data: {
          proofKind: "chain_summary",
          length: 0,
          headHash: h1,
          genesisHash: h2,
          timestamp: new Date(),
        },
      };
      expect(validateProof(payload)).toBe(false);
    });
  });

  describe("proofTypeFromData", () => {
    it("maps data kinds to proof types", () => {
      expect(proofTypeFromData({ proofKind: "chain_segment", events: [], startHash: h1, endHash: h2 })).toBe(ProofType.ChainSegment);
      expect(proofTypeFromData({
        proofKind: "event_existence", event: { hash: h1, prevHash: h2 },
        prevHash: h2, nextHash: Option.none(), position: 0, chainLength: 1,
      })).toBe(ProofType.EventExistence);
      expect(proofTypeFromData({
        proofKind: "chain_summary", length: 1, headHash: h1, genesisHash: h2, timestamp: new Date(),
      })).toBe(ProofType.ChainSummary);
    });
  });
});

// ── Handler ─────────────────────────────────────────────────────────────

describe("Handler", () => {
  function makeTestTransport(): ITransport {
    return {
      send: async () => undefined,
      listen: async function* () {},
    };
  }

  function makeSignedHelloEnvelope(fromId: SystemIdentity, toUri: SystemUri): Envelope {
    const envId = makeEnvelopeId();
    const payload = helloPayload({
      systemUri: fromId.systemUri(),
      publicKey: fromId.publicKey(),
      protocolVersions: [CurrentProtocolVersion],
      capabilities: ["treaty"],
      chainLength: 0,
    });

    const env = new Envelope(
      CurrentProtocolVersion, envId, fromId.systemUri(), toUri,
      MessageType.Hello, payload, new Date(),
      new Signature(new Uint8Array(64)), Option.none<EnvelopeId>(),
    );

    return signEnvelope(env, fromId);
  }

  it("processes HELLO and registers peer", () => {
    const idA = makeIdentity("system-a");
    const idB = makeIdentity("system-b");
    const peers = new PeerStore();
    const treaties = new TreatyStore();
    const handler = new Handler(idB, makeTestTransport(), peers, treaties);

    const env = makeSignedHelloEnvelope(idA, idB.systemUri());
    handler.handleIncoming(env);

    const [peer, found] = peers.get(idA.systemUri());
    expect(found).toBe(true);
    expect(peer.systemUri.value).toBe(idA.systemUri().value);
  });

  it("rejects duplicate envelopes", () => {
    const idA = makeIdentity("system-a");
    const idB = makeIdentity("system-b");
    const handler = new Handler(idB, makeTestTransport(), new PeerStore(), new TreatyStore());

    const env = makeSignedHelloEnvelope(idA, idB.systemUri());
    handler.handleIncoming(env);
    expect(() => handler.handleIncoming(env)).toThrow(DuplicateEnvelopeError);
  });

  it("rejects stale envelopes", () => {
    const idA = makeIdentity("system-a");
    const idB = makeIdentity("system-b");
    const handler = new Handler(idB, makeTestTransport(), new PeerStore(), new TreatyStore());

    const envId = makeEnvelopeId();
    const payload = helloPayload({
      systemUri: idA.systemUri(),
      publicKey: idA.publicKey(),
      protocolVersions: [1],
      capabilities: [],
      chainLength: 0,
    });

    const staleDate = new Date(Date.now() - 30 * 60 * 60 * 1000); // 30 hours ago
    const env = new Envelope(
      1, envId, idA.systemUri(), idB.systemUri(),
      MessageType.Hello, payload, staleDate,
      new Signature(new Uint8Array(64)), Option.none<EnvelopeId>(),
    );

    const signed = signEnvelope(env, idA);
    expect(() => handler.handleIncoming(signed)).toThrow("envelope timestamp out of range");
  });

  it("rejects invalid signatures", () => {
    const idA = makeIdentity("system-a");
    const idB = makeIdentity("system-b");
    const handler = new Handler(idB, makeTestTransport(), new PeerStore(), new TreatyStore());

    const envId = makeEnvelopeId();
    const payload = helloPayload({
      systemUri: idA.systemUri(),
      publicKey: idA.publicKey(),
      protocolVersions: [1],
      capabilities: [],
      chainLength: 0,
    });

    // Create envelope with wrong signature (sign with idB but claim to be from idA)
    const env = new Envelope(
      1, envId, idA.systemUri(), idB.systemUri(),
      MessageType.Hello, payload, new Date(),
      new Signature(new Uint8Array(64)), Option.none<EnvelopeId>(),
    );

    const wrongSigned = signEnvelope(env, idB); // signed by wrong identity
    expect(() => handler.handleIncoming(wrongSigned)).toThrow(EnvelopeSignatureInvalidError);
  });

  it("rejects unknown system for non-HELLO messages", () => {
    const idA = makeIdentity("system-a");
    const idB = makeIdentity("system-b");
    const handler = new Handler(idB, makeTestTransport(), new PeerStore(), new TreatyStore());

    const envId = makeEnvelopeId();
    const payload = receiptPayload({
      envelopeId: makeEnvelopeId(),
      status: ReceiptStatus.Delivered,
      localEventId: Option.none<EventId>(),
      reason: Option.none<string>(),
      signature: new Signature(new Uint8Array(64)),
    });

    const env = new Envelope(
      1, envId, idA.systemUri(), idB.systemUri(),
      MessageType.Receipt, payload, new Date(),
      new Signature(new Uint8Array(64)), Option.none<EnvelopeId>(),
    );

    const signed = signEnvelope(env, idA);
    expect(() => handler.handleIncoming(signed)).toThrow(SystemNotFoundError);
  });

  it("rejects incompatible protocol versions", () => {
    const idA = makeIdentity("system-a");
    const idB = makeIdentity("system-b");
    const handler = new Handler(idB, makeTestTransport(), new PeerStore(), new TreatyStore());
    handler.localProtocolVersions = [2]; // only support v2

    const envId = makeEnvelopeId();
    const payload = helloPayload({
      systemUri: idA.systemUri(),
      publicKey: idA.publicKey(),
      protocolVersions: [1], // only support v1
      capabilities: [],
      chainLength: 0,
    });

    const env = new Envelope(
      1, envId, idA.systemUri(), idB.systemUri(),
      MessageType.Hello, payload, new Date(),
      new Signature(new Uint8Array(64)), Option.none<EnvelopeId>(),
    );

    const signed = signEnvelope(env, idA);
    expect(() => handler.handleIncoming(signed)).toThrow(VersionIncompatibleError);
  });

  it("processes MESSAGE and calls onMessage", () => {
    const idA = makeIdentity("system-a");
    const idB = makeIdentity("system-b");
    const peers = new PeerStore();
    const handler = new Handler(idB, makeTestTransport(), peers, new TreatyStore());

    // First register peer via HELLO
    const hello = makeSignedHelloEnvelope(idA, idB.systemUri());
    handler.handleIncoming(hello);

    // Now send a MESSAGE
    let received: MessagePayloadContent | undefined;
    handler.onMessage = (from, payload) => { received = payload; };

    const envId = makeEnvelopeId();
    const msgPayload = messagePayload({
      content: { text: "hello" },
      contentType: new EventType("chat.message"),
      conversationId: Option.none<ConversationId>(),
      cgers: [],
    });

    const env = new Envelope(
      1, envId, idA.systemUri(), idB.systemUri(),
      MessageType.Message, msgPayload, new Date(),
      new Signature(new Uint8Array(64)), Option.none<EnvelopeId>(),
    );

    const signed = signEnvelope(env, idA);
    handler.handleIncoming(signed);

    expect(received).toBeDefined();
    expect(received!.content).toEqual({ text: "hello" });
  });

  it("processes TREATY propose and accept", () => {
    const idA = makeIdentity("system-a");
    const idB = makeIdentity("system-b");
    const peers = new PeerStore();
    const treaties = new TreatyStore();
    const handler = new Handler(idB, makeTestTransport(), peers, treaties);

    // Register peer
    handler.handleIncoming(makeSignedHelloEnvelope(idA, idB.systemUri()));

    const treatyId = makeTreatyId();
    const terms: TreatyTerm[] = [
      { scope: new DomainScope("data.access"), policy: "read_only", symmetric: true },
    ];

    // Propose
    const proposeEnv = new Envelope(
      1, makeEnvelopeId(), idA.systemUri(), idB.systemUri(),
      MessageType.Treaty, treatyPayload({
        treatyId,
        action: TreatyAction.Propose,
        terms,
        reason: Option.none<string>(),
      }), new Date(),
      new Signature(new Uint8Array(64)), Option.none<EnvelopeId>(),
    );
    handler.handleIncoming(signEnvelope(proposeEnv, idA));

    const [treaty1] = treaties.get(treatyId);
    expect(treaty1!.status).toBe(TreatyStatus.Proposed);

    // Accept
    const acceptEnv = new Envelope(
      1, makeEnvelopeId(), idA.systemUri(), idB.systemUri(),
      MessageType.Treaty, treatyPayload({
        treatyId,
        action: TreatyAction.Accept,
        terms: [],
        reason: Option.none<string>(),
      }), new Date(),
      new Signature(new Uint8Array(64)), Option.none<EnvelopeId>(),
    );
    handler.handleIncoming(signEnvelope(acceptEnv, idA));

    const [treaty2] = treaties.get(treatyId);
    expect(treaty2!.status).toBe(TreatyStatus.Active);
  });

  it("processes PROOF and updates trust", () => {
    const idA = makeIdentity("system-a");
    const idB = makeIdentity("system-b");
    const peers = new PeerStore();
    const handler = new Handler(idB, makeTestTransport(), peers, new TreatyStore());

    // Register peer
    handler.handleIncoming(makeSignedHelloEnvelope(idA, idB.systemUri()));

    const h1 = makeHash("a".repeat(64));
    const h2 = makeHash("b".repeat(64));

    const proofEnv = new Envelope(
      1, makeEnvelopeId(), idA.systemUri(), idB.systemUri(),
      MessageType.Proof, proofPayload({
        proofType: ProofType.ChainSummary,
        data: {
          proofKind: "chain_summary",
          length: 100,
          headHash: h1,
          genesisHash: h2,
          timestamp: new Date(),
        },
      }), new Date(),
      new Signature(new Uint8Array(64)), Option.none<EnvelopeId>(),
    );

    handler.handleIncoming(signEnvelope(proofEnv, idA));

    const [peer] = peers.get(idA.systemUri());
    // Trust should have increased from valid proof + receipt on time (from HELLO message handling)
    expect(peer.trust.value).toBeGreaterThan(0);
  });

  it("sends HELLO via transport", async () => {
    const idA = makeIdentity("system-a");
    let sentTo: SystemUri | undefined;
    let sentEnvelope: Envelope | undefined;

    const transport: ITransport = {
      send: async (to, env) => {
        sentTo = to;
        sentEnvelope = env;
        return undefined;
      },
      listen: async function* () {},
    };

    const handler = new Handler(idA, transport, new PeerStore(), new TreatyStore());
    await handler.hello(new SystemUri("eg://system-b"));

    expect(sentTo!.value).toBe("eg://system-b");
    expect(sentEnvelope!.type).toBe(MessageType.Hello);
    // Verify the sent envelope has a real signature
    const valid = verifyEnvelope(sentEnvelope!, idA, idA.publicKey());
    expect(valid).toBe(true);
  });

  it("handles transport failure in hello", async () => {
    const idA = makeIdentity("system-a");
    const transport: ITransport = {
      send: async () => { throw new Error("connection refused"); },
      listen: async function* () {},
    };

    const handler = new Handler(idA, transport, new PeerStore(), new TreatyStore());
    await expect(handler.hello(new SystemUri("eg://system-b"))).rejects.toThrow(TransportFailureError);
  });
});
