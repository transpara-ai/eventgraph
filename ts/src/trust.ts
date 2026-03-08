import { ActorId, DomainScope, EventId, Score, Weight } from "./types.js";
import { Actor } from "./actor.js";
import { Event } from "./event.js";

// ── TrustMetrics ─────────────────────────────────────────────────────────

/** Immutable trust query result for an actor. */
export class TrustMetrics {
  readonly actor: ActorId;
  readonly overall: Score;
  readonly byDomain: ReadonlyMap<string, Score>;
  readonly confidence: Score;
  readonly trend: Weight;
  readonly evidence: readonly EventId[];
  readonly lastUpdated: number; // nanoseconds
  readonly decayRate: Score;

  constructor(
    actor: ActorId,
    overall: Score,
    byDomain: Map<string, Score>,
    confidence: Score,
    trend: Weight,
    evidence: EventId[],
    lastUpdated: number,
    decayRate: Score,
  ) {
    this.actor = actor;
    this.overall = overall;
    this.byDomain = new Map(byDomain);
    this.confidence = confidence;
    this.trend = trend;
    this.evidence = [...evidence];
    this.lastUpdated = lastUpdated;
    this.decayRate = decayRate;
    Object.freeze(this);
  }
}

// ── TrustConfig ──────────────────────────────────────────────────────────

/** Configuration for the default trust model. */
export interface TrustConfig {
  /** Initial trust score for new actors. Default: 0.0 */
  initialTrust: Score;
  /** Trust decay per day. Default: 0.01 */
  decayRate: Score;
  /** Maximum trust adjustment from a single event. Default: 0.1 */
  maxAdjustment: Weight;
  /** Small trust delta for non-trust events. Default: 0.01 */
  observedEventDelta: number;
  /** Trend decay per day toward zero. Default: 0.01 */
  trendDecayRate: number;
}

/** Returns a TrustConfig with sensible defaults. */
export function defaultTrustConfig(): TrustConfig {
  return {
    initialTrust: new Score(0.0),
    decayRate: new Score(0.01),
    maxAdjustment: new Weight(0.1),
    observedEventDelta: 0.01,
    trendDecayRate: 0.01,
  };
}

// ── TrustModel interface ─────────────────────────────────────────────────

/** Trust model interface — calculates, updates, and decays trust. */
export interface TrustModel {
  score(actor: Actor): TrustMetrics;
  scoreInDomain(actor: Actor, domain: DomainScope): TrustMetrics;
  update(actor: Actor, evidence: Event): TrustMetrics;
  updateBetween(from: Actor, to: Actor, evidence: Event): TrustMetrics;
  decay(actor: Actor, elapsedSeconds: number): TrustMetrics;
  between(from: Actor, to: Actor): TrustMetrics;
}

// ── Internal trust state ─────────────────────────────────────────────────

interface TrustState {
  score: Score;
  byDomain: Map<string, Score>;
  evidence: EventId[];
  lastUpdated: number;
  trend: Weight;
}

function newTrustState(initialTrust: Score): TrustState {
  return {
    score: initialTrust,
    byDomain: new Map(),
    evidence: [],
    lastUpdated: Date.now() * 1_000_000,
    trend: new Weight(0.0),
  };
}

// ── DefaultTrustModel ────────────────────────────────────────────────────

/** Default trust model with linear decay and equal weighting. */
export class DefaultTrustModel implements TrustModel {
  private readonly config: TrustConfig;
  private readonly scores = new Map<string, TrustState>();
  private readonly directed = new Map<string, TrustState>();

  constructor(config?: Partial<TrustConfig>) {
    const defaults = defaultTrustConfig();
    this.config = config
      ? { ...defaults, ...config }
      : defaults;
  }

  score(actor: Actor): TrustMetrics {
    const state = this.getOrDefault(actor.id);
    return this.buildMetrics(actor.id, state);
  }

  scoreInDomain(actor: Actor, domain: DomainScope): TrustMetrics {
    const state = this.getOrDefault(actor.id);

    // If domain-specific score exists, return metrics with that score
    const domainScore = state.byDomain.get(domain.value);
    if (domainScore !== undefined) {
      return this.buildDomainMetrics(actor.id, state, domainScore);
    }

    // Fall back to global score with halved confidence (no domain-specific data)
    const evidenceCount = state.evidence.length;
    const globalConfidence = Math.min(1.0, evidenceCount / 50.0);
    return new TrustMetrics(
      actor.id,
      state.score,
      state.byDomain,
      new Score(globalConfidence * 0.5),
      state.trend,
      state.evidence,
      state.lastUpdated,
      this.config.decayRate,
    );
  }

  update(actor: Actor, evidence: Event): TrustMetrics {
    const state = this.getOrCreate(actor.id);

    // Deduplicate: if this evidence was already applied, return current metrics unchanged
    for (const id of state.evidence) {
      if (id.value === evidence.id.value) {
        return this.buildMetrics(actor.id, state);
      }
    }

    // Extract trust delta from the evidence event
    let delta = this.extractDeltaForScore(evidence, state.score.value);

    // Clamp to MaxAdjustment
    const maxAdj = this.config.maxAdjustment.value;
    delta = Math.max(-maxAdj, Math.min(maxAdj, delta));

    // Apply delta, clamp to [0, 1]
    const newScore = Math.max(0, Math.min(1, state.score.value + delta));
    state.score = new Score(newScore);

    // Update domain-specific score if evidence carries domain info
    this.updateDomainScore(state, evidence, delta);

    // Update trend
    this.updateTrend(state, delta);

    // Track evidence (capped at 100)
    state.evidence.push(evidence.id);
    if (state.evidence.length > 100) {
      state.evidence.splice(0, state.evidence.length - 100);
    }

    state.lastUpdated = Date.now() * 1_000_000;

    return this.buildMetrics(actor.id, state);
  }

  updateBetween(from: Actor, to: Actor, evidence: Event): TrustMetrics {
    const key = `${from.id.value}\0${to.id.value}`;
    let state = this.directed.get(key);
    if (!state) {
      state = newTrustState(this.config.initialTrust);
      this.directed.set(key, state);
    }

    // Deduplicate
    for (const id of state.evidence) {
      if (id.value === evidence.id.value) {
        return this.buildMetrics(to.id, state);
      }
    }

    let delta = this.extractDeltaForScore(evidence, state.score.value);

    const maxAdj = this.config.maxAdjustment.value;
    delta = Math.max(-maxAdj, Math.min(maxAdj, delta));

    const newScore = Math.max(0, Math.min(1, state.score.value + delta));
    state.score = new Score(newScore);

    // Update domain-specific score if evidence carries domain info
    this.updateDomainScore(state, evidence, delta);

    // Update trend
    this.updateTrend(state, delta);

    // Track evidence (capped at 100)
    state.evidence.push(evidence.id);
    if (state.evidence.length > 100) {
      state.evidence.splice(0, state.evidence.length - 100);
    }

    state.lastUpdated = Date.now() * 1_000_000;

    return this.buildMetrics(to.id, state);
  }

  decay(actor: Actor, elapsedSeconds: number): TrustMetrics {
    // Guard against negative durations
    if (elapsedSeconds <= 0) {
      const state = this.getOrDefault(actor.id);
      return this.buildMetrics(actor.id, state);
    }

    const days = elapsedSeconds / 86400;
    const decayAmount = this.config.decayRate.value * days;
    const trendDecay = this.config.trendDecayRate * days;

    // Decay undirected trust
    const key = actor.id.value;
    const state = this.scores.get(key);
    if (state) {
      this.decayState(state, decayAmount, trendDecay);
    }

    // Decay directed trust where actor is the trust holder (from)
    for (const [dkey, dstate] of this.directed) {
      if (dkey.startsWith(actor.id.value + "\0")) {
        this.decayState(dstate, decayAmount, trendDecay);
      }
    }

    if (!state) {
      return this.buildMetrics(actor.id, newTrustState(this.config.initialTrust));
    }

    return this.buildMetrics(actor.id, state);
  }

  between(from: Actor, to: Actor): TrustMetrics {
    const key = `${from.id.value}\0${to.id.value}`;
    const state = this.directed.get(key);
    if (!state) {
      // No direct trust relationship
      return new TrustMetrics(
        to.id,
        this.config.initialTrust,
        new Map(),
        new Score(0.0),  // low confidence
        new Weight(0.0),
        [],
        Date.now() * 1_000_000,
        this.config.decayRate,
      );
    }
    return this.buildMetrics(to.id, state);
  }

  // ── Private helpers ──────────────────────────────────────────────────

  private getOrCreate(actorId: ActorId): TrustState {
    const key = actorId.value;
    let state = this.scores.get(key);
    if (!state) {
      state = newTrustState(this.config.initialTrust);
      this.scores.set(key, state);
    }
    return state;
  }

  private getOrDefault(actorId: ActorId): TrustState {
    const key = actorId.value;
    const state = this.scores.get(key);
    if (state) {
      return state;
    }
    return newTrustState(this.config.initialTrust);
  }

  private extractDeltaForScore(ev: Event, currentScore: number): number {
    const content = ev.content;
    if (typeof content.current === "number") {
      return content.current - currentScore;
    }
    // Small positive trust boost for any observed (non-trust) event
    return this.config.observedEventDelta;
  }

  private updateDomainScore(state: TrustState, evidence: Event, delta: number): void {
    const content = evidence.content;
    if (typeof content.domain === "string" && content.domain.length > 0) {
      const domainKey = content.domain;
      const existing = state.byDomain.get(domainKey);
      let domainScore: number;
      if (existing !== undefined) {
        domainScore = Math.max(0, Math.min(1, existing.value + delta));
      } else {
        domainScore = Math.max(0, Math.min(1, this.config.initialTrust.value + delta));
      }
      state.byDomain.set(domainKey, new Score(domainScore));
    }
  }

  private updateTrend(state: TrustState, delta: number): void {
    if (delta > 0) {
      state.trend = new Weight(Math.min(1, state.trend.value + 0.1));
    } else if (delta < 0) {
      state.trend = new Weight(Math.max(-1, state.trend.value - 0.1));
    }
  }

  private decayState(state: TrustState, decayAmount: number, trendDecay: number): void {
    state.score = new Score(Math.max(0, state.score.value - decayAmount));

    for (const [domain, ds] of state.byDomain) {
      state.byDomain.set(domain, new Score(Math.max(0, ds.value - decayAmount)));
    }

    if (state.trend.value > 0) {
      state.trend = new Weight(Math.max(0, state.trend.value - trendDecay));
    } else if (state.trend.value < 0) {
      state.trend = new Weight(Math.min(0, state.trend.value + trendDecay));
    }

    state.lastUpdated = Date.now() * 1_000_000;
  }

  private buildMetrics(actorId: ActorId, state: TrustState): TrustMetrics {
    const evidenceCount = state.evidence.length;
    const confidence = Math.min(1.0, evidenceCount / 50.0);
    return new TrustMetrics(
      actorId,
      state.score,
      state.byDomain,
      new Score(confidence),
      state.trend,
      state.evidence,
      state.lastUpdated,
      this.config.decayRate,
    );
  }

  private buildDomainMetrics(actorId: ActorId, state: TrustState, domainScore: Score): TrustMetrics {
    const evidenceCount = state.evidence.length;
    const confidence = Math.min(1.0, evidenceCount / 50.0);
    return new TrustMetrics(
      actorId,
      domainScore,
      state.byDomain,
      new Score(confidence),
      state.trend,
      state.evidence,
      state.lastUpdated,
      this.config.decayRate,
    );
  }
}
