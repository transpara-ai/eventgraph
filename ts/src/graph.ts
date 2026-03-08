import { EventGraphError } from "./errors.js";
import {
  ActorId, ConversationId, EventId, EventType, Hash, Option,
} from "./types.js";
import {
  Event, Signer, NoopSigner, createBootstrap, createEvent,
} from "./event.js";
import { InMemoryStore } from "./store.js";
import { EventBus } from "./bus.js";
import { Actor, InMemoryActorStore } from "./actor.js";
import { TrustModel, DefaultTrustModel, TrustMetrics } from "./trust.js";
import {
  AuthorityChain, DefaultAuthorityChain, AuthorityResult,
} from "./authority.js";

// ── Errors ──────────────────────────────────────────────────────────────

/** Thrown when the graph is used before start() or after close(). */
export class GraphStateError extends EventGraphError {
  constructor(message: string) {
    super(message);
    this.name = "GraphStateError";
  }
}

// ── Config ──────────────────────────────────────────────────────────────

/** Tuning knobs for Graph behaviour. */
export interface GraphConfig {
  /** Buffer size for bus subscribers. Default: 256. */
  subscriberBufferSize?: number;
  /** Whether to fall back to mechanical decision when intelligence is unavailable. Default: true. */
  fallbackToMechanical?: boolean;
}

/** Optional collaborators injected at construction time. */
export interface GraphOptions {
  trustModel?: TrustModel;
  authorityChain?: AuthorityChain;
  signer?: Signer;
  config?: GraphConfig;
}

// ── Graph ───────────────────────────────────────────────────────────────

/**
 * Top-level facade for the event graph.
 *
 * Provides: bootstrap, record, evaluate, query.
 * Owns: Store, ActorStore, Bus, TrustModel, AuthorityChain, Signer.
 */
export class Graph {
  private readonly _store: InMemoryStore;
  private readonly _actorStore: InMemoryActorStore;
  private readonly _bus: EventBus;
  private readonly _trustModel: TrustModel;
  private readonly _authorityChain: AuthorityChain;
  private readonly _signer: Signer;
  private readonly _config: Required<GraphConfig>;

  private _started = false;
  private _closed = false;

  constructor(
    store: InMemoryStore,
    actorStore: InMemoryActorStore,
    options?: GraphOptions,
  ) {
    this._store = store;
    this._actorStore = actorStore;
    this._bus = new EventBus(store);

    this._signer = options?.signer ?? new NoopSigner();

    this._trustModel = options?.trustModel ?? new DefaultTrustModel();
    this._authorityChain = options?.authorityChain
      ?? new DefaultAuthorityChain(this._trustModel);

    this._config = {
      subscriberBufferSize: options?.config?.subscriberBufferSize ?? 256,
      fallbackToMechanical: options?.config?.fallbackToMechanical ?? true,
    };
  }

  // ── Lifecycle ───────────────────────────────────────────────────────

  /** Marks the graph as started. Must be called before record/evaluate/query. */
  start(): void {
    if (this._closed) throw new GraphStateError("graph is closed");
    this._started = true;
  }

  /** Closes the graph. Closes the bus and underlying store. */
  close(): void {
    if (this._closed) return;
    this._closed = true;
    this._started = false;
    this._bus.close();
    this._store.close();
  }

  // ── Getters ─────────────────────────────────────────────────────────

  get store(): InMemoryStore { return this._store; }
  get actorStore(): InMemoryActorStore { return this._actorStore; }
  get bus(): EventBus { return this._bus; }

  // ── Bootstrap ───────────────────────────────────────────────────────

  /**
   * Creates and appends a genesis (bootstrap) event for the given system actor.
   * May be called before start() — bootstrapping is a prerequisite to starting.
   */
  bootstrap(systemActor: ActorId, signer?: Signer): Event {
    if (this._closed) throw new GraphStateError("graph is closed");

    const sig = signer ?? this._signer;
    const event = createBootstrap(systemActor, sig);
    this._store.append(event);
    this._bus.publish(event);
    return event;
  }

  // ── Record ──────────────────────────────────────────────────────────

  /**
   * Creates an event, appends it to the store, and publishes it on the bus.
   * Requires that the graph has been started and that at least one event
   * (typically bootstrap) exists in the store.
   */
  record(
    eventType: EventType,
    source: ActorId,
    content: Record<string, unknown>,
    causes: EventId[],
    conversationId: ConversationId,
    signer?: Signer,
  ): Event {
    this.requireStarted();

    const head = this._store.head();
    if (head.isNone) {
      throw new GraphStateError("cannot record: store is empty (bootstrap first)");
    }
    const prevHash = head.unwrap().hash;
    const sig = signer ?? this._signer;

    const event = createEvent(
      eventType, source, content, causes,
      conversationId, prevHash, sig,
    );

    this._store.append(event);
    this._bus.publish(event);
    return event;
  }

  // ── Evaluate ────────────────────────────────────────────────────────

  /** Delegates authority evaluation to the configured AuthorityChain. */
  evaluate(actor: Actor, action: string): AuthorityResult {
    this.requireStarted();
    return this._authorityChain.evaluate(actor, action);
  }

  // ── Query ───────────────────────────────────────────────────────────

  /** Returns a Query helper bound to this graph's stores and trust model. */
  query(): Query {
    this.requireStarted();
    return new Query(this._store, this._actorStore, this._trustModel);
  }

  // ── Private ─────────────────────────────────────────────────────────

  private requireStarted(): void {
    if (this._closed) throw new GraphStateError("graph is closed");
    if (!this._started) throw new GraphStateError("graph is not started");
  }
}

// ── Query ───────────────────────────────────────────────────────────────

/**
 * Read-only query facade over the store, actor store, and trust model.
 * Returned by Graph.query().
 */
export class Query {
  private readonly store: InMemoryStore;
  private readonly actorStore: InMemoryActorStore;
  private readonly trustModel: TrustModel;

  constructor(
    store: InMemoryStore,
    actorStore: InMemoryActorStore,
    trustModel: TrustModel,
  ) {
    this.store = store;
    this.actorStore = actorStore;
    this.trustModel = trustModel;
  }

  /** Returns the most recent events, newest first. */
  recent(limit: number): Event[] {
    return this.store.recent(limit);
  }

  /** Returns events matching the given type, newest first. */
  byType(eventType: EventType, limit: number): Event[] {
    return this.store.byType(eventType, limit);
  }

  /** Returns events emitted by the given source, newest first. */
  bySource(source: ActorId, limit: number): Event[] {
    return this.store.bySource(source, limit);
  }

  /** Returns events in the given conversation, newest first. */
  byConversation(conversationId: ConversationId, limit: number): Event[] {
    return this.store.byConversation(conversationId, limit);
  }

  /** Returns causal ancestors of an event up to maxDepth. */
  ancestors(eventId: EventId, maxDepth: number): Event[] {
    return this.store.ancestors(eventId, maxDepth);
  }

  /** Returns causal descendants of an event up to maxDepth. */
  descendants(eventId: EventId, maxDepth: number): Event[] {
    return this.store.descendants(eventId, maxDepth);
  }

  /** Returns the trust score for an actor. */
  trustScore(actor: Actor): TrustMetrics {
    return this.trustModel.score(actor);
  }

  /** Returns the directed trust from one actor to another. */
  trustBetween(from: Actor, to: Actor): TrustMetrics {
    return this.trustModel.between(from, to);
  }

  /** Returns an actor by ID. */
  actor(id: ActorId): Actor {
    return this.actorStore.get(id);
  }

  /** Returns the number of events in the store. */
  eventCount(): number {
    return this.store.count();
  }
}
