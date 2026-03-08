import { describe, it, expect } from "vitest";
import { Graph, Query, GraphStateError } from "../src/graph.js";
import { InMemoryStore } from "../src/store.js";
import { InMemoryActorStore, Actor, ActorType, ActorStatus } from "../src/actor.js";
import { NoopSigner } from "../src/event.js";
import {
  ActorId, ConversationId, EventType, PublicKey, SubscriptionPattern,
} from "../src/types.js";
import { DefaultTrustModel } from "../src/trust.js";
import { DefaultAuthorityChain, AuthorityLevel } from "../src/authority.js";

function makeGraph(options?: Parameters<typeof Graph.prototype.constructor>[2]) {
  return new Graph(new InMemoryStore(), new InMemoryActorStore(), options);
}

function makeActor(name: string): { actorId: ActorId; publicKey: PublicKey } {
  const bytes = new Uint8Array(32);
  for (let i = 0; i < name.length && i < 32; i++) {
    bytes[i] = name.charCodeAt(i);
  }
  return { actorId: new ActorId(name), publicKey: new PublicKey(bytes) };
}

describe("Graph", () => {
  it("bootstrap creates genesis event", () => {
    const g = makeGraph();
    const actorId = new ActorId("system");
    const ev = g.bootstrap(actorId);

    expect(ev.type.value).toBe("system.bootstrapped");
    expect(ev.source.value).toBe("system");
    expect(g.store.count()).toBe(1);
  });

  it("start and close lifecycle", () => {
    const g = makeGraph();
    g.bootstrap(new ActorId("system"));
    g.start();
    // Should be able to query after start
    const q = g.query();
    expect(q.eventCount()).toBe(1);

    g.close();
    // After close, operations throw
    expect(() => g.query()).toThrow(GraphStateError);
  });

  it("record requires started graph", () => {
    const g = makeGraph();
    g.bootstrap(new ActorId("system"));
    // Not started yet
    expect(() =>
      g.record(
        new EventType("test.event"), new ActorId("alice"), {},
        [], new ConversationId("conv_1"),
      ),
    ).toThrow(GraphStateError);
  });

  it("record requires non-empty store", () => {
    const g = makeGraph();
    g.start();
    expect(() =>
      g.record(
        new EventType("test.event"), new ActorId("alice"), {},
        [], new ConversationId("conv_1"),
      ),
    ).toThrow(GraphStateError);
  });

  it("record appends and publishes event", () => {
    const g = makeGraph();
    const boot = g.bootstrap(new ActorId("system"));
    g.start();

    const published: string[] = [];
    g.bus.subscribe(new SubscriptionPattern("*"), (ev) => {
      published.push(ev.type.value);
    });

    const ev = g.record(
      new EventType("trust.updated"), new ActorId("alice"),
      { current: 0.5 }, [boot.id], new ConversationId("conv_1"),
    );

    expect(ev.type.value).toBe("trust.updated");
    expect(g.store.count()).toBe(2);
    expect(published).toContain("trust.updated");
  });

  it("bootstrap publishes on bus", () => {
    const g = makeGraph();
    const published: string[] = [];
    g.bus.subscribe(new SubscriptionPattern("system.*"), (ev) => {
      published.push(ev.type.value);
    });

    g.bootstrap(new ActorId("system"));
    expect(published).toEqual(["system.bootstrapped"]);
  });

  it("evaluate delegates to authority chain", () => {
    const trustModel = new DefaultTrustModel();
    const chain = new DefaultAuthorityChain(trustModel);
    const g = makeGraph({ trustModel, authorityChain: chain });
    g.bootstrap(new ActorId("system"));
    g.start();

    const { publicKey } = makeActor("bob");
    const actor = g.actorStore.register(publicKey, "Bob", ActorType.Human);

    const result = g.evaluate(actor, "some.action");
    expect(result.level).toBe(AuthorityLevel.Notification);
    expect(result.chain.length).toBeGreaterThan(0);
  });

  it("query.recent returns events newest first", () => {
    const g = makeGraph();
    const boot = g.bootstrap(new ActorId("system"));
    g.start();

    const e1 = g.record(
      new EventType("test.one"), new ActorId("alice"), {},
      [boot.id], new ConversationId("conv_1"),
    );
    const e2 = g.record(
      new EventType("test.two"), new ActorId("alice"), {},
      [e1.id], new ConversationId("conv_1"),
    );

    const recent = g.query().recent(10);
    expect(recent.length).toBe(3);
    expect(recent[0].id.value).toBe(e2.id.value);
    expect(recent[2].id.value).toBe(boot.id.value);
  });

  it("query.byType filters by event type", () => {
    const g = makeGraph();
    const boot = g.bootstrap(new ActorId("system"));
    g.start();

    g.record(
      new EventType("test.alpha"), new ActorId("alice"), {},
      [boot.id], new ConversationId("conv_1"),
    );
    g.record(
      new EventType("test.beta"), new ActorId("alice"), {},
      [boot.id], new ConversationId("conv_1"),
    );

    const results = g.query().byType(new EventType("test.alpha"), 10);
    expect(results.length).toBe(1);
    expect(results[0].type.value).toBe("test.alpha");
  });

  it("query.bySource filters by actor", () => {
    const g = makeGraph();
    const boot = g.bootstrap(new ActorId("system"));
    g.start();

    g.record(
      new EventType("test.event"), new ActorId("alice"), {},
      [boot.id], new ConversationId("conv_1"),
    );
    g.record(
      new EventType("test.event"), new ActorId("bob"), {},
      [boot.id], new ConversationId("conv_1"),
    );

    const results = g.query().bySource(new ActorId("alice"), 10);
    expect(results.length).toBe(1);
    expect(results[0].source.value).toBe("alice");
  });

  it("query.eventCount reflects store state", () => {
    const g = makeGraph();
    g.bootstrap(new ActorId("system"));
    g.start();
    expect(g.query().eventCount()).toBe(1);
  });

  it("query.actor retrieves registered actors", () => {
    const g = makeGraph();
    g.bootstrap(new ActorId("system"));
    g.start();

    const { publicKey } = makeActor("carol");
    const registered = g.actorStore.register(publicKey, "Carol", ActorType.Human);

    const found = g.query().actor(registered.id);
    expect(found.displayName).toBe("Carol");
  });

  it("close is idempotent", () => {
    const g = makeGraph();
    g.bootstrap(new ActorId("system"));
    g.start();
    g.close();
    g.close(); // second close should not throw
    expect(() => g.start()).toThrow(GraphStateError);
  });

  it("custom signer is used for bootstrap and record", () => {
    const calls: number[] = [];
    const spy: import("../src/event.js").Signer = {
      sign(data: Uint8Array) {
        calls.push(data.length);
        return new NoopSigner().sign(data);
      },
    };

    const g = makeGraph({ signer: spy });
    g.bootstrap(new ActorId("system"));
    expect(calls.length).toBe(1);

    g.start();
    const boot = g.store.head().unwrap();
    g.record(
      new EventType("test.event"), new ActorId("alice"), {},
      [boot.id], new ConversationId("conv_1"),
    );
    expect(calls.length).toBe(2);
  });

  it("record with per-call signer overrides default", () => {
    const defaultCalls: number[] = [];
    const overrideCalls: number[] = [];
    const defaultSigner: import("../src/event.js").Signer = {
      sign(data: Uint8Array) {
        defaultCalls.push(data.length);
        return new NoopSigner().sign(data);
      },
    };
    const overrideSigner: import("../src/event.js").Signer = {
      sign(data: Uint8Array) {
        overrideCalls.push(data.length);
        return new NoopSigner().sign(data);
      },
    };

    const g = makeGraph({ signer: defaultSigner });
    g.bootstrap(new ActorId("system"));
    g.start();

    const boot = g.store.head().unwrap();
    g.record(
      new EventType("test.event"), new ActorId("alice"), {},
      [boot.id], new ConversationId("conv_1"), overrideSigner,
    );

    // Default was used for bootstrap, override for record
    expect(defaultCalls.length).toBe(1);
    expect(overrideCalls.length).toBe(1);
  });

  it("query.trustScore and trustBetween delegate to trust model", () => {
    const g = makeGraph();
    g.bootstrap(new ActorId("system"));
    g.start();

    const { publicKey: pk1 } = makeActor("dave");
    const { publicKey: pk2 } = makeActor("eve\0");
    const a1 = g.actorStore.register(pk1, "Dave", ActorType.Human);
    const a2 = g.actorStore.register(pk2, "Eve", ActorType.AI);

    const score = g.query().trustScore(a1);
    expect(score.overall.value).toBeGreaterThanOrEqual(0);
    expect(score.overall.value).toBeLessThanOrEqual(1);

    const between = g.query().trustBetween(a1, a2);
    expect(between.overall.value).toBeGreaterThanOrEqual(0);
  });

  it("query.byConversation filters correctly", () => {
    const g = makeGraph();
    const boot = g.bootstrap(new ActorId("system"));
    g.start();

    g.record(
      new EventType("test.event"), new ActorId("alice"), {},
      [boot.id], new ConversationId("conv_a"),
    );
    g.record(
      new EventType("test.event"), new ActorId("alice"), {},
      [boot.id], new ConversationId("conv_b"),
    );

    const results = g.query().byConversation(new ConversationId("conv_a"), 10);
    expect(results.length).toBe(1);
    expect(results[0].conversationId.value).toBe("conv_a");
  });
});
