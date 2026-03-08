import { createPrivateKey, createPublicKey, sign, verify, generateKeyPairSync, randomBytes } from "node:crypto";
import {
  Option, Score, Weight, DomainScope, EnvelopeId, TreatyId, SystemUri,
  PublicKey, Signature, Hash, EventId, EventType, ActorId, ConversationId,
} from "./types.js";
import { EventGraphError, InvalidTransitionError } from "./errors.js";

// ── Constants / Enums ───────────────────────────────────────────────────

export const MessageType = {
  Hello: "Hello",
  Message: "Message",
  Receipt: "Receipt",
  Proof: "Proof",
  Treaty: "Treaty",
  AuthorityRequest: "AuthorityRequest",
  Discover: "Discover",
} as const;
export type MessageType = (typeof MessageType)[keyof typeof MessageType];

const validMessageTypes = new Set<string>(Object.values(MessageType));
export function isValidMessageType(v: string): v is MessageType {
  return validMessageTypes.has(v);
}

export const TreatyStatus = {
  Proposed: "Proposed",
  Active: "Active",
  Suspended: "Suspended",
  Terminated: "Terminated",
} as const;
export type TreatyStatus = (typeof TreatyStatus)[keyof typeof TreatyStatus];

export const TreatyAction = {
  Propose: "Propose",
  Accept: "Accept",
  Modify: "Modify",
  Suspend: "Suspend",
  Terminate: "Terminate",
} as const;
export type TreatyAction = (typeof TreatyAction)[keyof typeof TreatyAction];

export const ReceiptStatus = {
  Delivered: "Delivered",
  Processed: "Processed",
  Rejected: "Rejected",
} as const;
export type ReceiptStatus = (typeof ReceiptStatus)[keyof typeof ReceiptStatus];

export const ProofType = {
  ChainSegment: "ChainSegment",
  EventExistence: "EventExistence",
  ChainSummary: "ChainSummary",
} as const;
export type ProofType = (typeof ProofType)[keyof typeof ProofType];

export const CGERRelationship = {
  CausedBy: "CausedBy",
  References: "References",
  RespondsTo: "RespondsTo",
} as const;
export type CGERRelationship = (typeof CGERRelationship)[keyof typeof CGERRelationship];

export const AuthorityLevel = {
  Required: "Required",
  Recommended: "Recommended",
  Notification: "Notification",
} as const;
export type AuthorityLevel = (typeof AuthorityLevel)[keyof typeof AuthorityLevel];

// ── Protocol constants ──────────────────────────────────────────────────

export const CurrentProtocolVersion = 1;
export const MaxEnvelopeAgeMs = 25 * 60 * 60 * 1000; // 25 hours in ms
const DedupPruneInterval = 1000;

// ── Trust impact constants ──────────────────────────────────────────────

export const TrustImpactValidProof = 0.02;
export const TrustImpactReceiptOnTime = 0.01;
export const TrustImpactTreatyHonoured = 0.03;
export const TrustImpactTreatyViolated = -0.15;
export const TrustImpactInvalidProof = -0.10;
export const TrustImpactSignatureInvalid = -0.20;
export const TrustImpactNoHelloResponse = -0.05;
export const InterSystemDecayRate = new Score(0.02); // per day
export const InterSystemMaxAdjustment = new Weight(0.05);

// ── EGIP Errors ─────────────────────────────────────────────────────────

export class EGIPError extends EventGraphError {
  constructor(message: string) {
    super(message);
    this.name = "EGIPError";
  }
}

export class SystemNotFoundError extends EGIPError {
  constructor(public readonly uri: SystemUri) {
    super(`system not found: ${uri.value}`);
    this.name = "SystemNotFoundError";
  }
}

export class EnvelopeSignatureInvalidError extends EGIPError {
  constructor(public readonly envelopeId: EnvelopeId) {
    super(`envelope signature invalid: ${envelopeId.value}`);
    this.name = "EnvelopeSignatureInvalidError";
  }
}

export class TreatyViolationError extends EGIPError {
  constructor(
    public readonly treatyId: TreatyId,
    public readonly term: string,
  ) {
    super(`treaty ${treatyId.value} violated: ${term}`);
    this.name = "TreatyViolationError";
  }
}

export class TrustInsufficientError extends EGIPError {
  constructor(
    public readonly system: SystemUri,
    public readonly score: Score,
    public readonly required: Score,
  ) {
    super(`trust insufficient for ${system.value}: have ${score.value}, need ${required.value}`);
    this.name = "TrustInsufficientError";
  }
}

export class TransportFailureError extends EGIPError {
  constructor(
    public readonly to: SystemUri,
    public readonly reason: string,
  ) {
    super(`transport failure to ${to.value}: ${reason}`);
    this.name = "TransportFailureError";
  }
}

export class DuplicateEnvelopeError extends EGIPError {
  constructor(public readonly envelopeId: EnvelopeId) {
    super(`duplicate envelope: ${envelopeId.value}`);
    this.name = "DuplicateEnvelopeError";
  }
}

export class TreatyNotFoundError extends EGIPError {
  constructor(public readonly treatyId: TreatyId) {
    super(`treaty not found: ${treatyId.value}`);
    this.name = "TreatyNotFoundError";
  }
}

export class VersionIncompatibleError extends EGIPError {
  constructor(
    public readonly local: readonly number[],
    public readonly remote: readonly number[],
  ) {
    super(`no compatible protocol version: local ${JSON.stringify(local)}, remote ${JSON.stringify(remote)}`);
    this.name = "VersionIncompatibleError";
  }
}

// ── IIdentity ───────────────────────────────────────────────────────────

export interface IIdentity {
  systemUri(): SystemUri;
  publicKey(): PublicKey;
  sign(data: Uint8Array): Signature;
  verify(publicKey: PublicKey, data: Uint8Array, signature: Signature): boolean;
}

// ── SystemIdentity (Ed25519 via node:crypto) ────────────────────────────

export class SystemIdentity implements IIdentity {
  private readonly _uri: SystemUri;
  private readonly _publicKey: PublicKey;
  private readonly _privateKeyDer: Buffer;
  private readonly _publicKeyDer: Buffer;
  readonly createdAt: Date;

  private constructor(
    uri: SystemUri,
    publicKey: PublicKey,
    privateKeyDer: Buffer,
    publicKeyDer: Buffer,
  ) {
    this._uri = uri;
    this._publicKey = publicKey;
    this._privateKeyDer = privateKeyDer;
    this._publicKeyDer = publicKeyDer;
    this.createdAt = new Date();
  }

  static generate(uri: SystemUri): SystemIdentity {
    const { publicKey: pub, privateKey: priv } = generateKeyPairSync("ed25519");
    const privDer = priv.export({ type: "pkcs8", format: "der" });
    const pubDer = pub.export({ type: "spki", format: "der" });
    // Ed25519 raw public key is the last 32 bytes of the SPKI DER encoding
    const rawPub = pubDer.subarray(pubDer.length - 32);
    return new SystemIdentity(uri, new PublicKey(rawPub), privDer as Buffer, pubDer as Buffer);
  }

  systemUri(): SystemUri { return this._uri; }
  publicKey(): PublicKey { return this._publicKey; }

  sign(data: Uint8Array): Signature {
    const privKey = createPrivateKey({ key: this._privateKeyDer, format: "der", type: "pkcs8" });
    const sig = sign(null, Buffer.from(data), privKey);
    return new Signature(new Uint8Array(sig));
  }

  verify(publicKey: PublicKey, data: Uint8Array, signature: Signature): boolean {
    try {
      // Build SPKI DER from raw 32-byte Ed25519 public key
      const spkiPrefix = Buffer.from("302a300506032b6570032100", "hex");
      const spkiDer = Buffer.concat([spkiPrefix, Buffer.from(publicKey.bytes)]);
      const pubKey = createPublicKey({ key: spkiDer, format: "der", type: "spki" });
      return verify(null, Buffer.from(data), pubKey, Buffer.from(signature.bytes));
    } catch {
      return false;
    }
  }
}

// ── Payload types ───────────────────────────────────────────────────────

export interface HelloPayload {
  readonly kind: "hello";
  readonly systemUri: SystemUri;
  readonly publicKey: PublicKey;
  readonly protocolVersions: readonly number[];
  readonly capabilities: readonly string[];
  readonly chainLength: number;
}

export function helloPayload(p: Omit<HelloPayload, "kind">): HelloPayload {
  return { kind: "hello", ...p };
}

export interface MessagePayloadContent {
  readonly kind: "message";
  readonly content: Record<string, unknown>;
  readonly contentType: EventType;
  readonly conversationId: Option<ConversationId>;
  readonly cgers: readonly CGER[];
}

export function messagePayload(p: Omit<MessagePayloadContent, "kind">): MessagePayloadContent {
  return { kind: "message", ...p };
}

export interface ReceiptPayload {
  readonly kind: "receipt";
  readonly envelopeId: EnvelopeId;
  readonly status: ReceiptStatus;
  readonly localEventId: Option<EventId>;
  readonly reason: Option<string>;
  readonly signature: Signature;
}

export function receiptPayload(p: Omit<ReceiptPayload, "kind">): ReceiptPayload {
  return { kind: "receipt", ...p };
}

export interface ProofPayload {
  readonly kind: "proof";
  readonly proofType: ProofType;
  readonly data: ProofData;
}

export function proofPayload(p: Omit<ProofPayload, "kind">): ProofPayload {
  return { kind: "proof", ...p };
}

export type ProofData = ChainSegmentProof | EventExistenceProof | ChainSummaryProof;

export interface ChainSegmentProof {
  readonly proofKind: "chain_segment";
  readonly events: readonly ProofEvent[];
  readonly startHash: Hash;
  readonly endHash: Hash;
}

/** Minimal event representation for proofs (not the full Event class). */
export interface ProofEvent {
  readonly hash: Hash;
  readonly prevHash: Hash;
}

export interface EventExistenceProof {
  readonly proofKind: "event_existence";
  readonly event: ProofEvent;
  readonly prevHash: Hash;
  readonly nextHash: Option<Hash>;
  readonly position: number;
  readonly chainLength: number;
}

export interface ChainSummaryProof {
  readonly proofKind: "chain_summary";
  readonly length: number;
  readonly headHash: Hash;
  readonly genesisHash: Hash;
  readonly timestamp: Date;
}

export interface TreatyPayload {
  readonly kind: "treaty";
  readonly treatyId: TreatyId;
  readonly action: TreatyAction;
  readonly terms: readonly TreatyTerm[];
  readonly reason: Option<string>;
}

export function treatyPayload(p: Omit<TreatyPayload, "kind">): TreatyPayload {
  return { kind: "treaty", ...p };
}

export interface AuthorityRequestPayload {
  readonly kind: "authority_request";
  readonly action: DomainScope;
  readonly actor: ActorId;
  readonly level: AuthorityLevel;
  readonly justification: string;
  readonly treatyId: Option<TreatyId>;
}

export function authorityRequestPayload(p: Omit<AuthorityRequestPayload, "kind">): AuthorityRequestPayload {
  return { kind: "authority_request", ...p };
}

export interface DiscoverPayload {
  readonly kind: "discover";
  readonly query: DiscoverQuery;
  readonly results: readonly DiscoverResult[];
}

export function discoverPayload(p: Omit<DiscoverPayload, "kind">): DiscoverPayload {
  return { kind: "discover", ...p };
}

export interface DiscoverQuery {
  readonly capabilities: readonly string[];
  readonly minTrust: Option<Score>;
}

export interface DiscoverResult {
  readonly systemUri: SystemUri;
  readonly publicKey: PublicKey;
  readonly capabilities: readonly string[];
  readonly trustScore: Score;
}

export type MessagePayload =
  | HelloPayload
  | MessagePayloadContent
  | ReceiptPayload
  | ProofPayload
  | TreatyPayload
  | AuthorityRequestPayload
  | DiscoverPayload;

// ── CGER ────────────────────────────────────────────────────────────────

export interface CGER {
  readonly localEventId: EventId;
  readonly remoteSystem: SystemUri;
  readonly remoteEventId: string;
  readonly remoteHash: Hash;
  readonly relationship: CGERRelationship;
  readonly verified: boolean;
}

// ── TreatyTerm ──────────────────────────────────────────────────────────

export interface TreatyTerm {
  readonly scope: DomainScope;
  readonly policy: string;
  readonly symmetric: boolean;
}

// ── Envelope ────────────────────────────────────────────────────────────

export class Envelope {
  constructor(
    public readonly protocolVersion: number,
    public readonly id: EnvelopeId,
    public readonly from: SystemUri,
    public readonly to: SystemUri,
    public readonly type: MessageType,
    public readonly payload: MessagePayload,
    public readonly timestamp: Date,
    public readonly signature: Signature,
    public readonly inReplyTo: Option<EnvelopeId>,
  ) {}

  /** Returns the canonical string representation for signing. */
  canonicalForm(): string {
    const payloadJson = canonicalPayloadJson(this.payload);
    const msgType = this.type.toLowerCase();
    const nanos = BigInt(this.timestamp.getTime()) * 1_000_000n;
    const inReplyTo = this.inReplyTo.isSome ? this.inReplyTo.unwrap().value : "";

    return `${this.protocolVersion}|${this.id.value}|${this.from.value}|${this.to.value}|${msgType}|${nanos.toString()}|${inReplyTo}|${payloadJson}`;
  }

  /** Returns a new Envelope with the given signature. */
  withSignature(sig: Signature): Envelope {
    return new Envelope(
      this.protocolVersion, this.id, this.from, this.to,
      this.type, this.payload, this.timestamp, sig, this.inReplyTo,
    );
  }
}

/** Deterministic JSON for payload (sorted keys, no undefined). */
function canonicalPayloadJson(payload: MessagePayload): string {
  return JSON.stringify(payloadToPlain(payload), Object.keys(payloadToPlain(payload)).sort());
}

function payloadToPlain(payload: MessagePayload): Record<string, unknown> {
  // Build a plain object from the typed payload, then sort keys via JSON round-trip
  const obj: Record<string, unknown> = {};
  switch (payload.kind) {
    case "hello":
      obj["system_uri"] = payload.systemUri.value;
      obj["public_key"] = payload.publicKey.toString();
      obj["protocol_versions"] = [...payload.protocolVersions];
      obj["capabilities"] = [...payload.capabilities];
      obj["chain_length"] = payload.chainLength;
      break;
    case "message":
      obj["content"] = payload.content;
      obj["content_type"] = payload.contentType.value;
      obj["conversation_id"] = payload.conversationId.isSome ? payload.conversationId.unwrap().value : null;
      obj["cgers"] = payload.cgers.map(cgerToPlain);
      break;
    case "receipt":
      obj["envelope_id"] = payload.envelopeId.value;
      obj["status"] = payload.status;
      obj["local_event_id"] = payload.localEventId.isSome ? payload.localEventId.unwrap().value : null;
      obj["reason"] = payload.reason.isSome ? payload.reason.unwrap() : null;
      obj["signature"] = payload.signature.toString();
      break;
    case "proof":
      obj["proof_type"] = payload.proofType;
      obj["data"] = proofDataToPlain(payload.data);
      break;
    case "treaty":
      obj["treaty_id"] = payload.treatyId.value;
      obj["action"] = payload.action;
      obj["terms"] = payload.terms.map(t => ({
        scope: t.scope.value,
        policy: t.policy,
        symmetric: t.symmetric,
      }));
      obj["reason"] = payload.reason.isSome ? payload.reason.unwrap() : null;
      break;
    case "authority_request":
      obj["action"] = payload.action.value;
      obj["actor"] = payload.actor.value;
      obj["level"] = payload.level;
      obj["justification"] = payload.justification;
      obj["treaty_id"] = payload.treatyId.isSome ? payload.treatyId.unwrap().value : null;
      break;
    case "discover":
      obj["query"] = {
        capabilities: [...payload.query.capabilities],
        min_trust: payload.query.minTrust.isSome ? payload.query.minTrust.unwrap().value : null,
      };
      obj["results"] = payload.results.map(r => ({
        system_uri: r.systemUri.value,
        public_key: r.publicKey.toString(),
        capabilities: [...r.capabilities],
        trust_score: r.trustScore.value,
      }));
      break;
  }
  // Round-trip through JSON.parse to normalize and sort keys recursively
  const raw = JSON.parse(JSON.stringify(obj));
  return sortObjectKeys(raw) as Record<string, unknown>;
}

function sortObjectKeys(obj: unknown): unknown {
  if (obj === null || obj === undefined) return obj;
  if (Array.isArray(obj)) return obj.map(sortObjectKeys);
  if (typeof obj === "object") {
    const sorted: Record<string, unknown> = {};
    for (const key of Object.keys(obj as Record<string, unknown>).sort()) {
      sorted[key] = sortObjectKeys((obj as Record<string, unknown>)[key]);
    }
    return sorted;
  }
  return obj;
}

function cgerToPlain(c: CGER): Record<string, unknown> {
  return {
    local_event_id: c.localEventId.value,
    remote_system: c.remoteSystem.value,
    remote_event_id: c.remoteEventId,
    remote_hash: c.remoteHash.value,
    relationship: c.relationship,
    verified: c.verified,
  };
}

function proofDataToPlain(data: ProofData): Record<string, unknown> {
  switch (data.proofKind) {
    case "chain_segment":
      return {
        events: data.events.map(e => ({ hash: e.hash.value, prev_hash: e.prevHash.value })),
        start_hash: data.startHash.value,
        end_hash: data.endHash.value,
      };
    case "event_existence":
      return {
        event: { hash: data.event.hash.value, prev_hash: data.event.prevHash.value },
        prev_hash: data.prevHash.value,
        next_hash: data.nextHash.isSome ? data.nextHash.unwrap().value : null,
        position: data.position,
        chain_length: data.chainLength,
      };
    case "chain_summary":
      return {
        length: data.length,
        head_hash: data.headHash.value,
        genesis_hash: data.genesisHash.value,
        timestamp: data.timestamp.toISOString(),
      };
  }
}

// ── Sign / Verify Envelope ──────────────────────────────────────────────

export function signEnvelope(env: Envelope, identity: IIdentity): Envelope {
  const canonical = env.canonicalForm();
  const sig = identity.sign(new TextEncoder().encode(canonical));
  return env.withSignature(sig);
}

export function verifyEnvelope(env: Envelope, identity: IIdentity, publicKey: PublicKey): boolean {
  const canonical = env.canonicalForm();
  return identity.verify(publicKey, new TextEncoder().encode(canonical), env.signature);
}

// ── Version negotiation ─────────────────────────────────────────────────

export function negotiateVersion(local: readonly number[], remote: readonly number[]): Option<number> {
  let best = -1;
  for (const l of local) {
    for (const r of remote) {
      if (l === r && l > best) best = l;
    }
  }
  return best < 0 ? Option.none<number>() : Option.some(best);
}

// ── Treaty (state machine) ──────────────────────────────────────────────

const validTreatyTransitions: Record<TreatyStatus, readonly TreatyStatus[]> = {
  [TreatyStatus.Proposed]: [TreatyStatus.Active, TreatyStatus.Terminated],
  [TreatyStatus.Active]: [TreatyStatus.Suspended, TreatyStatus.Terminated],
  [TreatyStatus.Suspended]: [TreatyStatus.Active, TreatyStatus.Terminated],
  [TreatyStatus.Terminated]: [],
};

export class Treaty {
  id: TreatyId;
  systemA: SystemUri;
  systemB: SystemUri;
  status: TreatyStatus;
  terms: TreatyTerm[];
  createdAt: Date;
  updatedAt: Date;

  constructor(id: TreatyId, systemA: SystemUri, systemB: SystemUri, terms: TreatyTerm[]) {
    this.id = id;
    this.systemA = systemA;
    this.systemB = systemB;
    this.status = TreatyStatus.Proposed;
    this.terms = [...terms];
    const now = new Date();
    this.createdAt = now;
    this.updatedAt = now;
  }

  transition(to: TreatyStatus): void {
    const allowed = validTreatyTransitions[this.status];
    if (!allowed.includes(to)) {
      throw new InvalidTransitionError(this.status, to);
    }
    this.status = to;
    this.updatedAt = new Date();
  }

  applyAction(action: TreatyAction): void {
    switch (action) {
      case TreatyAction.Accept:
        this.transition(TreatyStatus.Active);
        break;
      case TreatyAction.Suspend:
        this.transition(TreatyStatus.Suspended);
        break;
      case TreatyAction.Terminate:
        this.transition(TreatyStatus.Terminated);
        break;
      case TreatyAction.Modify:
        if (this.status !== TreatyStatus.Active) {
          throw new EGIPError(`can only modify active treaties, current status: ${this.status}`);
        }
        this.updatedAt = new Date();
        break;
      case TreatyAction.Propose:
        throw new EGIPError("cannot apply Propose to existing treaty");
      default:
        throw new EGIPError(`unknown treaty action: ${action as string}`);
    }
  }
}

// ── PeerRecord / PeerStore ──────────────────────────────────────────────

export interface PeerRecord {
  systemUri: SystemUri;
  publicKey: PublicKey;
  trust: Score;
  capabilities: string[];
  negotiatedVersion: number;
  lastSeen: Date;
  firstSeen: Date;
  lastDecayedAt: Date;
}

export class PeerStore {
  private readonly peers = new Map<string, PeerRecord>();

  register(uri: SystemUri, publicKey: PublicKey, capabilities: string[], negotiatedVersion: number): PeerRecord {
    const key = uri.value;
    const now = new Date();
    const existing = this.peers.get(key);
    if (existing) {
      // Do NOT overwrite PublicKey on re-registration (key-substitution protection)
      existing.capabilities = [...capabilities];
      existing.negotiatedVersion = negotiatedVersion;
      existing.lastSeen = now;
      return existing;
    }

    const record: PeerRecord = {
      systemUri: uri,
      publicKey,
      trust: new Score(0.0),
      capabilities: [...capabilities],
      negotiatedVersion,
      lastSeen: now,
      firstSeen: now,
      lastDecayedAt: now,
    };
    this.peers.set(key, record);
    return record;
  }

  get(uri: SystemUri): [PeerRecord, boolean] {
    const record = this.peers.get(uri.value);
    if (!record) return [undefined as unknown as PeerRecord, false];
    return [{ ...record, capabilities: [...record.capabilities] }, true];
  }

  updateTrust(uri: SystemUri, delta: number): [Score, boolean] {
    const record = this.peers.get(uri.value);
    if (!record) return [new Score(0.0), false];

    // Positive trust accumulates gradually; negative trust hits immediately.
    let adjustedDelta = delta;
    if (adjustedDelta > 0) {
      const maxAdj = InterSystemMaxAdjustment.value;
      if (adjustedDelta > maxAdj) adjustedDelta = maxAdj;
    }

    let newVal = record.trust.value + adjustedDelta;
    if (newVal < 0.0) newVal = 0.0;
    if (newVal > 1.0) newVal = 1.0;

    record.trust = new Score(newVal);
    record.lastSeen = new Date();
    return [record.trust, true];
  }

  decayAll(): void {
    const now = new Date();
    const decayPerDay = InterSystemDecayRate.value;
    for (const record of this.peers.values()) {
      const daysSince = (now.getTime() - record.lastDecayedAt.getTime()) / (24 * 60 * 60 * 1000);
      if (daysSince <= 0) continue;
      const decay = decayPerDay * daysSince;
      let newVal = record.trust.value - decay;
      if (newVal < 0.0) newVal = 0.0;
      record.trust = new Score(newVal);
      record.lastDecayedAt = now;
    }
  }

  all(): PeerRecord[] {
    return [...this.peers.values()].map(r => ({ ...r, capabilities: [...r.capabilities] }));
  }
}

// ── TreatyStore ─────────────────────────────────────────────────────────

export class TreatyStore {
  private readonly treaties = new Map<string, Treaty>();

  put(treaty: Treaty): void {
    this.treaties.set(treaty.id.value, treaty);
  }

  get(id: TreatyId): [Treaty | undefined, boolean] {
    const t = this.treaties.get(id.value);
    if (!t) return [undefined, false];
    // Return a shallow copy with copied terms
    const copy = Object.create(Treaty.prototype) as Treaty;
    Object.assign(copy, t, { terms: [...t.terms] });
    return [copy, true];
  }

  apply(id: TreatyId, fn: (treaty: Treaty) => void): void {
    const t = this.treaties.get(id.value);
    if (!t) throw new TreatyNotFoundError(id);
    fn(t);
  }

  bySystem(uri: SystemUri): Treaty[] {
    const result: Treaty[] = [];
    for (const t of this.treaties.values()) {
      if (t.systemA.value === uri.value || t.systemB.value === uri.value) {
        const copy = Object.create(Treaty.prototype) as Treaty;
        Object.assign(copy, t, { terms: [...t.terms] });
        result.push(copy);
      }
    }
    return result;
  }

  active(): Treaty[] {
    const result: Treaty[] = [];
    for (const t of this.treaties.values()) {
      if (t.status === TreatyStatus.Active) {
        const copy = Object.create(Treaty.prototype) as Treaty;
        Object.assign(copy, t, { terms: [...t.terms] });
        result.push(copy);
      }
    }
    return result;
  }
}

// ── EnvelopeDedup ───────────────────────────────────────────────────────

export class EnvelopeDedup {
  private readonly seen = new Map<string, number>(); // key -> timestamp ms
  private readonly ttlMs: number;
  private checkCount = 0;

  constructor(ttlMs?: number) {
    this.ttlMs = ttlMs ?? MaxEnvelopeAgeMs + 60 * 60 * 1000; // MaxEnvelopeAge + 1 hour
  }

  /** Returns true if the envelope ID has not been seen before. */
  check(id: EnvelopeId): boolean {
    const key = id.value;
    if (this.seen.has(key)) return false;
    this.seen.set(key, Date.now());
    this.checkCount++;
    if (this.checkCount % DedupPruneInterval === 0) {
      this.pruneInternal();
    }
    return true;
  }

  prune(): number {
    return this.pruneInternal();
  }

  private pruneInternal(): number {
    const cutoff = Date.now() - this.ttlMs;
    let removed = 0;
    for (const [key, ts] of this.seen) {
      if (ts <= cutoff) {
        this.seen.delete(key);
        removed++;
      }
    }
    return removed;
  }

  get size(): number { return this.seen.size; }
}

// ── ITransport ──────────────────────────────────────────────────────────

export interface IncomingEnvelope {
  readonly envelope?: Envelope;
  readonly error?: Error;
}

export interface ITransport {
  send(to: SystemUri, envelope: Envelope): Promise<ReceiptPayload | undefined>;
  listen(): AsyncIterable<IncomingEnvelope>;
}

// ── Proof verification ──────────────────────────────────────────────────

export function verifyChainSegment(proof: ChainSegmentProof): boolean {
  if (proof.events.length === 0) return false;

  // First event's prevHash must match startHash
  if (proof.events[0].prevHash.value !== proof.startHash.value) return false;

  // Internal hash chain continuity
  for (let i = 1; i < proof.events.length; i++) {
    if (proof.events[i].prevHash.value !== proof.events[i - 1].hash.value) return false;
  }

  // Last event's hash must match endHash
  if (proof.events[proof.events.length - 1].hash.value !== proof.endHash.value) return false;

  return true;
}

export function verifyEventExistence(proof: EventExistenceProof): boolean {
  // Event's prevHash should match proof's prevHash
  if (proof.event.prevHash.value !== proof.prevHash.value) return false;

  // Position and chain length should be consistent
  if (proof.position < 0 || proof.position >= proof.chainLength) return false;

  // Event should have a non-zero hash
  if (proof.event.hash.isZero) return false;

  return true;
}

export function validateProof(payload: ProofPayload): boolean {
  switch (payload.data.proofKind) {
    case "chain_segment":
      return verifyChainSegment(payload.data);
    case "event_existence":
      return verifyEventExistence(payload.data);
    case "chain_summary":
      return payload.data.length > 0;
  }
}

export function proofTypeFromData(data: ProofData): ProofType {
  switch (data.proofKind) {
    case "chain_segment": return ProofType.ChainSegment;
    case "event_existence": return ProofType.EventExistence;
    case "chain_summary": return ProofType.ChainSummary;
  }
}

// ── UUID v4 generator ───────────────────────────────────────────────────

function generateUUID4(): string {
  const b = randomBytes(16);
  b[6] = (b[6] & 0x0f) | 0x40; // version 4
  b[8] = (b[8] & 0x3f) | 0x80; // variant 10
  const hex = b.toString("hex");
  return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20, 32)}`;
}

// ── Handler ─────────────────────────────────────────────────────────────

export class Handler {
  readonly identity: IIdentity;
  readonly transport: ITransport;
  readonly peers: PeerStore;
  readonly treaties: TreatyStore;
  readonly dedup: EnvelopeDedup;

  localProtocolVersions: number[];
  capabilities: string[];
  chainLength?: () => number;

  onMessage?: (from: SystemUri, payload: MessagePayloadContent) => void;
  onAuthorityRequest?: (from: SystemUri, payload: AuthorityRequestPayload) => void;
  onDiscover?: (from: SystemUri, query: DiscoverQuery) => DiscoverResult[];

  constructor(identity: IIdentity, transport: ITransport, peers: PeerStore, treaties: TreatyStore) {
    this.identity = identity;
    this.transport = transport;
    this.peers = peers;
    this.treaties = treaties;
    this.dedup = new EnvelopeDedup();
    this.localProtocolVersions = [CurrentProtocolVersion];
    this.capabilities = ["treaty", "proof"];
  }

  async hello(to: SystemUri): Promise<void> {
    const chainLen = this.chainLength ? this.chainLength() : 0;
    const envId = new EnvelopeId(generateUUID4());

    const env = new Envelope(
      CurrentProtocolVersion,
      envId,
      this.identity.systemUri(),
      to,
      MessageType.Hello,
      helloPayload({
        systemUri: this.identity.systemUri(),
        publicKey: this.identity.publicKey(),
        protocolVersions: this.localProtocolVersions,
        capabilities: this.capabilities,
        chainLength: chainLen,
      }),
      new Date(),
      new Signature(new Uint8Array(64)), // placeholder
      Option.none<EnvelopeId>(),
    );

    const signed = signEnvelope(env, this.identity);

    let receipt: ReceiptPayload | undefined;
    try {
      receipt = await this.transport.send(to, signed);
    } catch (err) {
      this.peers.updateTrust(to, TrustImpactNoHelloResponse);
      throw new TransportFailureError(to, err instanceof Error ? err.message : String(err));
    }

    if (receipt && receipt.status === ReceiptStatus.Rejected) {
      const reason = receipt.reason.isSome ? receipt.reason.unwrap() : "";
      throw new EGIPError(`hello rejected: ${reason}`);
    }
  }

  handleIncoming(env: Envelope): void {
    // Timestamp freshness check
    const ageMs = Date.now() - env.timestamp.getTime();
    if (ageMs > MaxEnvelopeAgeMs || ageMs < -5 * 60 * 1000) {
      throw new EGIPError(`envelope timestamp out of range: age ${ageMs}ms`);
    }

    // Replay deduplication
    if (!this.dedup.check(env.id)) {
      throw new DuplicateEnvelopeError(env.id);
    }

    // Look up sender's public key
    const [peer, known] = this.peers.get(env.from);

    let pubKey: PublicKey;
    if (env.type === MessageType.Hello) {
      const hello = env.payload as HelloPayload;
      if (hello.kind !== "hello") {
        throw new EGIPError(`invalid hello payload type`);
      }
      pubKey = hello.publicKey;
    } else {
      if (!known) throw new SystemNotFoundError(env.from);
      pubKey = peer.publicKey;
    }

    // Verify signature
    const valid = verifyEnvelope(env, this.identity, pubKey);
    if (!valid) {
      this.peers.updateTrust(env.from, TrustImpactSignatureInvalid);
      throw new EnvelopeSignatureInvalidError(env.id);
    }

    // Dispatch by message type
    switch (env.type) {
      case MessageType.Hello: this.handleHello(env); break;
      case MessageType.Message: this.handleMessage(env); break;
      case MessageType.Receipt: this.handleReceipt(env); break;
      case MessageType.Proof: this.handleProof(env); break;
      case MessageType.Treaty: this.handleTreaty(env); break;
      case MessageType.AuthorityRequest: this.handleAuthorityRequest(env); break;
      case MessageType.Discover: this.handleDiscover(env); break;
      default:
        throw new EGIPError(`unknown message type: ${env.type as string}`);
    }
  }

  private handleHello(env: Envelope): void {
    const hello = env.payload as HelloPayload;
    const version = negotiateVersion(this.localProtocolVersions, [...hello.protocolVersions]);
    if (version.isNone) {
      throw new VersionIncompatibleError(this.localProtocolVersions, [...hello.protocolVersions]);
    }
    this.peers.register(hello.systemUri, hello.publicKey, [...hello.capabilities], version.unwrap());
  }

  private handleMessage(env: Envelope): void {
    const msg = env.payload as MessagePayloadContent;
    this.peers.updateTrust(env.from, TrustImpactReceiptOnTime);
    if (this.onMessage) this.onMessage(env.from, msg);
  }

  private handleReceipt(env: Envelope): void {
    const receipt = env.payload as ReceiptPayload;
    if (receipt.status === ReceiptStatus.Processed || receipt.status === ReceiptStatus.Delivered) {
      this.peers.updateTrust(env.from, TrustImpactReceiptOnTime);
    }
  }

  private handleProof(env: Envelope): void {
    const proof = env.payload as ProofPayload;
    const valid = validateProof(proof);
    if (valid) {
      this.peers.updateTrust(env.from, TrustImpactValidProof);
    } else {
      this.peers.updateTrust(env.from, TrustImpactInvalidProof);
    }
  }

  private handleTreaty(env: Envelope): void {
    const payload = env.payload as TreatyPayload;
    switch (payload.action) {
      case TreatyAction.Propose: {
        const treaty = new Treaty(payload.treatyId, env.from, env.to, [...payload.terms]);
        this.treaties.put(treaty);
        break;
      }
      case TreatyAction.Accept:
        this.treaties.apply(payload.treatyId, (treaty) => {
          treaty.applyAction(TreatyAction.Accept);
          this.peers.updateTrust(env.from, TrustImpactTreatyHonoured);
        });
        break;
      case TreatyAction.Suspend:
        this.treaties.apply(payload.treatyId, (treaty) => {
          treaty.applyAction(TreatyAction.Suspend);
        });
        break;
      case TreatyAction.Terminate:
        this.treaties.apply(payload.treatyId, (treaty) => {
          treaty.applyAction(TreatyAction.Terminate);
        });
        break;
      case TreatyAction.Modify:
        this.treaties.apply(payload.treatyId, (treaty) => {
          treaty.applyAction(TreatyAction.Modify);
          treaty.terms = [...payload.terms];
        });
        break;
      default:
        throw new EGIPError(`unknown treaty action: ${payload.action as string}`);
    }
  }

  private handleAuthorityRequest(env: Envelope): void {
    const payload = env.payload as AuthorityRequestPayload;
    if (this.onAuthorityRequest) this.onAuthorityRequest(env.from, payload);
  }

  private handleDiscover(env: Envelope): void {
    const payload = env.payload as DiscoverPayload;
    if (!this.onDiscover) return;
    // Note: in a full implementation this would send results back via transport.
    // Keeping synchronous for the TS port since transport.send is async.
    this.onDiscover(env.from, payload.query);
  }
}
