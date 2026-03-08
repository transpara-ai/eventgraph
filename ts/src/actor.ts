import { createHash } from "node:crypto";
import {
  ActorNotFoundError,
  ActorKeyNotFoundError,
  InvalidTransitionError,
} from "./errors.js";
import { ActorId, EventId, PublicKey } from "./types.js";

// ── ActorType ────────────────────────────────────────────────────────────

/** What kind of decision-maker an actor is. */
export enum ActorType {
  Human = "Human",
  AI = "AI",
  System = "System",
  Committee = "Committee",
  RulesEngine = "RulesEngine",
}

// ── ActorStatus ──────────────────────────────────────────────────────────

/** An actor's lifecycle status. Memorial is terminal. */
export enum ActorStatus {
  Active = "active",
  Suspended = "suspended",
  Memorial = "memorial",
}

const validActorTransitions: Record<ActorStatus, ActorStatus[]> = {
  [ActorStatus.Active]: [ActorStatus.Suspended, ActorStatus.Memorial],
  [ActorStatus.Suspended]: [ActorStatus.Active, ActorStatus.Memorial],
  [ActorStatus.Memorial]: [], // terminal
};

/**
 * Attempt to transition from one ActorStatus to another.
 * Returns the target status on success, throws InvalidTransitionError on failure.
 */
export function transitionTo(current: ActorStatus, target: ActorStatus): ActorStatus {
  const valid = validActorTransitions[current];
  if (valid.includes(target)) {
    return target;
  }
  throw new InvalidTransitionError(current, target);
}

/** Returns the list of valid target statuses from the given status. */
export function validTransitions(status: ActorStatus): ActorStatus[] {
  return [...validActorTransitions[status]];
}

// ── Deep copy ────────────────────────────────────────────────────────────

function deepCopyMetadata(src: Record<string, unknown>): Record<string, unknown> {
  return JSON.parse(JSON.stringify(src));
}

// ── Actor ────────────────────────────────────────────────────────────────

/** Immutable actor identity in the system. */
export class Actor {
  private readonly _id: ActorId;
  private readonly _publicKey: PublicKey;
  private readonly _displayName: string;
  private readonly _actorType: ActorType;
  private readonly _metadata: Record<string, unknown>;
  private readonly _createdAt: number;
  private readonly _status: ActorStatus;

  constructor(
    id: ActorId,
    publicKey: PublicKey,
    displayName: string,
    actorType: ActorType,
    metadata: Record<string, unknown>,
    createdAt: number,
    status: ActorStatus,
  ) {
    this._id = id;
    this._publicKey = publicKey;
    this._displayName = displayName;
    this._actorType = actorType;
    this._metadata = deepCopyMetadata(metadata);
    this._createdAt = createdAt;
    this._status = status;
    Object.freeze(this);
  }

  get id(): ActorId { return this._id; }
  get publicKey(): PublicKey { return this._publicKey; }
  get displayName(): string { return this._displayName; }
  get actorType(): ActorType { return this._actorType; }
  get metadata(): Record<string, unknown> { return deepCopyMetadata(this._metadata); }
  get createdAt(): number { return this._createdAt; }
  get status(): ActorStatus { return this._status; }

  /** Returns a copy with a new status. */
  withStatus(status: ActorStatus): Actor {
    return new Actor(
      this._id, this._publicKey, this._displayName,
      this._actorType, this._metadata, this._createdAt, status,
    );
  }

  /** Returns a copy with updates applied. Metadata is merged, not replaced. */
  withUpdates(u: ActorUpdate): Actor {
    const md = deepCopyMetadata(this._metadata);
    if (u.metadata !== undefined) {
      for (const [k, v] of Object.entries(u.metadata)) {
        md[k] = JSON.parse(JSON.stringify(v));
      }
    }
    return new Actor(
      this._id, this._publicKey,
      u.displayName !== undefined ? u.displayName : this._displayName,
      this._actorType, md, this._createdAt, this._status,
    );
  }
}

// ── ActorUpdate ──────────────────────────────────────────────────────────

/** Describes updates to apply to an actor. */
export interface ActorUpdate {
  displayName?: string;
  metadata?: Record<string, unknown>;
}

// ── ActorFilter ──────────────────────────────────────────────────────────

/** Describes criteria for listing actors. */
export interface ActorFilter {
  status?: ActorStatus;
  actorType?: ActorType;
  limit?: number;
  after?: string;
}

// ── Page<T> ──────────────────────────────────────────────────────────────

/** A page of results with optional cursor for pagination. */
export interface Page<T> {
  items: T[];
  cursor?: string;
  hasMore: boolean;
}

// ── ActorStore interface ─────────────────────────────────────────────────

/** Actor persistence interface. */
export interface ActorStore {
  register(publicKey: PublicKey, displayName: string, actorType: ActorType): Actor;
  get(id: ActorId): Actor;
  getByPublicKey(publicKey: PublicKey): Actor;
  update(id: ActorId, updates: ActorUpdate): Actor;
  list(filter: ActorFilter): Page<Actor>;
  suspend(id: ActorId, reason: EventId): Actor;
  reactivate(id: ActorId, reason: EventId): Actor;
  memorial(id: ActorId, reason: EventId): Actor;
}

// ── InMemoryActorStore ───────────────────────────────────────────────────

/** Derives actor_id from SHA-256 of public key bytes: actor_{hex(sha256(pk).slice(0,16))} */
function deriveActorId(pk: PublicKey): ActorId {
  const hash = createHash("sha256").update(pk.bytes).digest();
  const id = `actor_${hash.subarray(0, 16).toString("hex")}`;
  return new ActorId(id);
}

function pubKeyHex(pk: PublicKey): string {
  return Buffer.from(pk.bytes).toString("hex");
}

/** In-memory implementation of ActorStore. */
export class InMemoryActorStore implements ActorStore {
  private readonly actors = new Map<string, Actor>();
  private readonly byKey = new Map<string, string>(); // hex(pk) → actorId value
  private readonly ordered: string[] = []; // insertion order for pagination

  register(publicKey: PublicKey, displayName: string, actorType: ActorType): Actor {
    const keyHex = pubKeyHex(publicKey);
    const existingId = this.byKey.get(keyHex);
    if (existingId !== undefined) {
      return this.actors.get(existingId)!;
    }

    const id = deriveActorId(publicKey);
    const actor = new Actor(
      id, publicKey, displayName, actorType,
      {}, Date.now() * 1_000_000, ActorStatus.Active,
    );
    this.actors.set(id.value, actor);
    this.byKey.set(keyHex, id.value);
    this.ordered.push(id.value);
    return actor;
  }

  get(id: ActorId): Actor {
    const actor = this.actors.get(id.value);
    if (actor === undefined) {
      throw new ActorNotFoundError(id.value);
    }
    return actor;
  }

  getByPublicKey(publicKey: PublicKey): Actor {
    const keyHex = pubKeyHex(publicKey);
    const id = this.byKey.get(keyHex);
    if (id === undefined) {
      throw new ActorKeyNotFoundError(keyHex);
    }
    return this.actors.get(id)!;
  }

  update(id: ActorId, updates: ActorUpdate): Actor {
    const actor = this.actors.get(id.value);
    if (actor === undefined) {
      throw new ActorNotFoundError(id.value);
    }
    const updated = actor.withUpdates(updates);
    this.actors.set(id.value, updated);
    return updated;
  }

  list(filter: ActorFilter): Page<Actor> {
    const limit = filter.limit !== undefined && filter.limit > 0 ? filter.limit : 100;

    // Find start position
    let startIdx = 0;
    if (filter.after !== undefined) {
      const idx = this.ordered.indexOf(filter.after);
      if (idx === -1) {
        throw new ActorNotFoundError(filter.after);
      }
      startIdx = idx + 1;
    }

    const items: Actor[] = [];
    for (let i = startIdx; i < this.ordered.length && items.length < limit; i++) {
      const actor = this.actors.get(this.ordered[i])!;
      if (filter.status !== undefined && actor.status !== filter.status) continue;
      if (filter.actorType !== undefined && actor.actorType !== filter.actorType) continue;
      items.push(actor);
    }

    let hasMore = false;
    let cursor: string | undefined;
    if (items.length === limit) {
      // Check if there are more matching items after the last one
      const lastActor = items[items.length - 1];
      const lastIdx = this.ordered.indexOf(lastActor.id.value);
      for (let i = lastIdx + 1; i < this.ordered.length; i++) {
        const actor = this.actors.get(this.ordered[i])!;
        if (filter.status !== undefined && actor.status !== filter.status) continue;
        if (filter.actorType !== undefined && actor.actorType !== filter.actorType) continue;
        hasMore = true;
        break;
      }
      if (hasMore) {
        cursor = lastActor.id.value;
      }
    }

    return { items, cursor, hasMore };
  }

  suspend(id: ActorId, reason: EventId): Actor {
    const actor = this.actors.get(id.value);
    if (actor === undefined) {
      throw new ActorNotFoundError(id.value);
    }
    const newStatus = transitionTo(actor.status, ActorStatus.Suspended);
    const updated = actor.withStatus(newStatus);
    this.actors.set(id.value, updated);
    void reason; // recorded on the event graph, not stored here
    return updated;
  }

  reactivate(id: ActorId, reason: EventId): Actor {
    const actor = this.actors.get(id.value);
    if (actor === undefined) {
      throw new ActorNotFoundError(id.value);
    }
    const newStatus = transitionTo(actor.status, ActorStatus.Active);
    const updated = actor.withStatus(newStatus);
    this.actors.set(id.value, updated);
    void reason;
    return updated;
  }

  memorial(id: ActorId, reason: EventId): Actor {
    const actor = this.actors.get(id.value);
    if (actor === undefined) {
      throw new ActorNotFoundError(id.value);
    }
    const newStatus = transitionTo(actor.status, ActorStatus.Memorial);
    const updated = actor.withStatus(newStatus);
    this.actors.set(id.value, updated);
    void reason;
    return updated;
  }

  /** Returns the number of registered actors. For testing. */
  actorCount(): number {
    return this.actors.size;
  }
}
