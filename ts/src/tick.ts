import { Event, createEvent, NoopSigner } from "./event.js";
import { Lifecycle, type Mutation, type Primitive, Registry, type Snapshot } from "./primitive.js";
import { InMemoryStore } from "./store.js";
import { Hash, PrimitiveId, type SubscriptionPattern } from "./types.js";

export interface TickConfig {
  maxWavesPerTick: number;
}

export interface TickResult {
  tick: number;
  waves: number;
  mutations: number;
  quiesced: boolean;
  durationMs: number;
  errors: string[];
}

export class TickEngine {
  private readonly registry: Registry;
  private readonly store: InMemoryStore;
  private readonly config: TickConfig;
  private readonly publisher?: (ev: Event) => void;
  private readonly signer = new NoopSigner();
  private currentTick = 0;

  constructor(
    registry: Registry,
    store: InMemoryStore,
    config?: Partial<TickConfig>,
    publisher?: (ev: Event) => void,
  ) {
    this.registry = registry;
    this.store = store;
    this.config = { maxWavesPerTick: 10, ...config };
    this.publisher = publisher;
  }

  tick(pendingEvents?: Event[]): TickResult {
    const start = performance.now();
    this.currentTick++;
    const tickNum = this.currentTick;

    let waveEvents = [...(pendingEvents ?? [])];
    let totalMutations = 0;
    const errors: string[] = [];
    let quiesced = false;
    const invokedThisTick = new Set<string>();
    let wavesRun = 0;

    for (let wave = 0; wave < this.config.maxWavesPerTick; wave++) {
      if (waveEvents.length === 0 && wave > 0) {
        quiesced = true;
        break;
      }
      wavesRun = wave + 1;

      const snapshot: Snapshot = {
        tick: tickNum,
        primitives: this.registry.allStates(),
        pendingEvents: waveEvents,
        recentEvents: this.store.recent(50),
      };

      const newEvents: Event[] = [];
      const allMutations: Mutation[] = [];

      for (const prim of this.registry.all()) {
        const pid = prim.id();
        const lifecycle = this.registry.getLifecycle(pid);
        if (lifecycle !== Lifecycle.Active) continue;

        if (!invokedThisTick.has(pid.value)) {
          const last = this.registry.getLastTick(pid);
          if (tickNum - last < prim.cadence().value) continue;
        }

        const matched = filterEvents(waveEvents, prim.subscriptions());
        if (matched.length === 0 && wave > 0) continue;

        try { this.registry.setLifecycle(pid, Lifecycle.Processing); }
        catch { continue; }

        try {
          const mutations = prim.process(tickNum, matched, snapshot);
          allMutations.push(...mutations);
        } catch (e) { errors.push(`${pid.value}: ${e}`); }

        try {
          this.registry.setLifecycle(pid, Lifecycle.Emitting);
          this.registry.setLifecycle(pid, Lifecycle.Active);
        } catch (e) { errors.push(`${pid.value} lifecycle: ${e}`); }

        invokedThisTick.add(pid.value);
      }

      for (const m of allMutations) {
        totalMutations++;
        try {
          const ev = this.applyMutation(m);
          if (ev) newEvents.push(ev);
        } catch (e) { errors.push(`mutation: ${e}`); }
      }

      waveEvents = newEvents;
    }

    for (const pidVal of invokedThisTick) {
      this.registry.setLastTick(new PrimitiveId(pidVal), tickNum);
    }

    if (!quiesced && wavesRun >= this.config.maxWavesPerTick) {
      quiesced = false;
    }

    return {
      tick: tickNum,
      waves: wavesRun,
      mutations: totalMutations,
      quiesced,
      durationMs: performance.now() - start,
      errors,
    };
  }

  private applyMutation(m: Mutation): Event | null {
    switch (m.kind) {
      case "addEvent": {
        const head = this.store.head();
        const prevHash = head.isSome ? head.unwrap().hash : Hash.zero();
        const ev = createEvent(m.type, m.source, m.content, m.causes, m.conversationId, prevHash, this.signer);
        this.store.append(ev);
        this.publisher?.(ev);
        return ev;
      }
      case "updateState":
        this.registry.updateState(m.primitiveId, m.key, m.value);
        return null;
      case "updateActivation":
        this.registry.setActivation(m.primitiveId, m.level);
        return null;
      case "updateLifecycle":
        this.registry.setLifecycle(m.primitiveId, m.state);
        return null;
    }
  }
}

function filterEvents(events: Event[], patterns: SubscriptionPattern[]): Event[] {
  return events.filter((ev) => patterns.some((p) => p.matches(ev.type)));
}
