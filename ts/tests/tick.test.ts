import { describe, it, expect } from "vitest";
import { createBootstrap, NoopSigner, Event } from "../src/event.js";
import { Registry, Lifecycle, type Mutation, type Snapshot, type Primitive } from "../src/primitive.js";
import { InMemoryStore } from "../src/store.js";
import { TickEngine } from "../src/tick.js";
import {
  Activation, ActorId, Cadence, ConversationId, EventType, Layer, PrimitiveId, SubscriptionPattern,
} from "../src/types.js";

class CountingPrimitive implements Primitive {
  received = 0;
  private _id: PrimitiveId;
  constructor(name: string, private _layer = 0) { this._id = new PrimitiveId(name); }
  id() { return this._id; }
  layer() { return new Layer(this._layer); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(_tick: number, events: Event[], _snap: Snapshot): Mutation[] {
    this.received += events.length;
    return [{ kind: "updateState", primitiveId: this._id, key: "count", value: this.received }];
  }
}

class EmittingPrimitive implements Primitive {
  emissions = 0;
  private _id: PrimitiveId;
  constructor(name: string, private max: number) { this._id = new PrimitiveId(name); }
  id() { return this._id; }
  layer() { return new Layer(0); }
  cadence() { return new Cadence(1); }
  subscriptions() { return [new SubscriptionPattern("*")]; }
  process(_tick: number, events: Event[]): Mutation[] {
    if (events.length === 0 || this.emissions >= this.max) return [];
    this.emissions++;
    return [{ kind: "addEvent", type: new EventType("test.emitted"), source: new ActorId("emitter"), content: { w: this.emissions }, causes: [events[0].id], conversationId: new ConversationId("conv_t") }];
  }
}

function setup(prims: Primitive[] = [], config?: { maxWavesPerTick: number }) {
  const reg = new Registry();
  const store = new InMemoryStore();
  const b = createBootstrap(new ActorId("system"), new NoopSigner());
  store.append(b);
  for (const p of prims) { reg.register(p); reg.activate(p.id()); }
  const engine = new TickEngine(reg, store, config);
  return { reg, store, engine, boot: b };
}

describe("TickEngine", () => {
  it("basic tick", () => {
    const c = new CountingPrimitive("counter");
    const { engine, boot } = setup([c]);
    const r = engine.tick([boot]);
    expect(r.tick).toBe(1);
    expect(c.received).toBe(1);
  });

  it("quiescence", () => {
    const { engine, boot } = setup([new CountingPrimitive("c")]);
    expect(engine.tick([boot]).quiesced).toBe(true);
  });

  it("ripple waves", () => {
    const e = new EmittingPrimitive("e", 3);
    const c = new CountingPrimitive("c");
    const { engine, boot } = setup([e, c]);
    const r = engine.tick([boot]);
    expect(r.waves).toBeGreaterThan(1);
    expect(c.received).toBeGreaterThan(1);
  });

  it("max waves limit", () => {
    const inf: Primitive = {
      id: () => new PrimitiveId("inf"),
      layer: () => new Layer(0),
      cadence: () => new Cadence(1),
      subscriptions: () => [new SubscriptionPattern("*")],
      process: (_t, events) => events.length ? [{ kind: "addEvent", type: new EventType("test.loop"), source: new ActorId("inf"), content: {}, causes: [events[0].id], conversationId: new ConversationId("c") }] : [],
    };
    const { engine, boot } = setup([inf], { maxWavesPerTick: 3 });
    const r = engine.tick([boot]);
    expect(r.waves).toBe(3);
    expect(r.quiesced).toBe(false);
  });

  it("inactive skipped", () => {
    const c = new CountingPrimitive("c");
    const reg = new Registry();
    const store = new InMemoryStore();
    const b = createBootstrap(new ActorId("sys"), new NoopSigner());
    store.append(b);
    reg.register(c); // don't activate
    const engine = new TickEngine(reg, store);
    engine.tick([b]);
    expect(c.received).toBe(0);
  });

  it("tick counter increments", () => {
    const { engine, boot } = setup();
    expect(engine.tick([boot]).tick).toBe(1);
    expect(engine.tick().tick).toBe(2);
  });

  it("layer ordering", () => {
    const order: string[] = [];
    const mkTracker = (name: string, layer: number): Primitive => ({
      id: () => new PrimitiveId(name),
      layer: () => new Layer(layer),
      cadence: () => new Cadence(1),
      subscriptions: () => [new SubscriptionPattern("*")],
      process: () => { order.push(name); return []; },
    });
    const { engine, boot } = setup([mkTracker("high", 5), mkTracker("low", 0), mkTracker("mid", 2)]);
    engine.tick([boot]);
    expect(order).toEqual(["low", "mid", "high"]);
  });
});
