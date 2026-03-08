import { describe, it, expect } from "vitest";
import { createBootstrap, createEvent, NoopSigner } from "../src/event.js";
import { InMemoryStore } from "../src/store.js";
import { ActorId, ConversationId, EventId, EventType } from "../src/types.js";
import { ChainIntegrityError, EventNotFoundError } from "../src/errors.js";

const boot = () => createBootstrap(new ActorId("alice"), new NoopSigner());
const next = (prev: ReturnType<typeof boot>) => createEvent(
  new EventType("trust.updated"), new ActorId("alice"), {},
  [prev.id], new ConversationId("conv_1"), prev.hash, new NoopSigner(),
);

describe("InMemoryStore", () => {
  it("append and get", () => {
    const s = new InMemoryStore();
    const b = boot();
    s.append(b);
    expect(s.get(b.id).id.value).toBe(b.id.value);
  });

  it("head empty", () => expect(new InMemoryStore().head().isNone).toBe(true));

  it("head after append", () => {
    const s = new InMemoryStore();
    const b = boot();
    s.append(b);
    expect(s.head().unwrap().id.value).toBe(b.id.value);
  });

  it("count", () => {
    const s = new InMemoryStore();
    expect(s.count()).toBe(0);
    s.append(boot());
    expect(s.count()).toBe(1);
  });

  it("chain of events", () => {
    const s = new InMemoryStore();
    const b = boot(); s.append(b);
    const e1 = next(b); s.append(e1);
    const e2 = next(e1); s.append(e2);
    expect(s.count()).toBe(3);
  });

  it("rejects broken chain", () => {
    const s = new InMemoryStore();
    const b = boot(); s.append(b);
    const bad = createEvent(new EventType("trust.updated"), new ActorId("alice"), {}, [b.id], new ConversationId("c"), b.prevHash, new NoopSigner());
    expect(() => s.append(bad)).toThrow(ChainIntegrityError);
  });

  it("get nonexistent", () => {
    expect(() => new InMemoryStore().get(new EventId("019462a0-0000-7000-8000-000000000099"))).toThrow(EventNotFoundError);
  });

  it("verify chain valid", () => {
    const s = new InMemoryStore();
    const b = boot(); s.append(b); s.append(next(b));
    expect(s.verifyChain()).toEqual({ valid: true, length: 2 });
  });

  it("recent", () => {
    const s = new InMemoryStore();
    const b = boot(); s.append(b);
    const e1 = next(b); s.append(e1);
    const r = s.recent(2);
    expect(r).toHaveLength(2);
    expect(r[0].id.value).toBe(e1.id.value);
  });

  it("byType returns matching events newest first", () => {
    const s = new InMemoryStore();
    const b = boot(); s.append(b);
    const e1 = next(b); s.append(e1);
    const e2 = next(e1); s.append(e2);

    const results = s.byType(new EventType("trust.updated"), 10);
    expect(results).toHaveLength(2);
    expect(results[0].id.value).toBe(e2.id.value);
    expect(results[1].id.value).toBe(e1.id.value);
  });

  it("byType respects limit", () => {
    const s = new InMemoryStore();
    const b = boot(); s.append(b);
    const e1 = next(b); s.append(e1);
    const e2 = next(e1); s.append(e2);

    const results = s.byType(new EventType("trust.updated"), 1);
    expect(results).toHaveLength(1);
    expect(results[0].id.value).toBe(e2.id.value);
  });

  it("byType returns empty for no matches", () => {
    const s = new InMemoryStore();
    const b = boot(); s.append(b);
    expect(s.byType(new EventType("trust.updated"), 10)).toHaveLength(0);
  });

  it("bySource returns matching events newest first", () => {
    const s = new InMemoryStore();
    const b = boot(); s.append(b);
    const e1 = next(b); s.append(e1);

    const results = s.bySource(new ActorId("alice"), 10);
    expect(results).toHaveLength(2);
    expect(results[0].id.value).toBe(e1.id.value);
    expect(results[1].id.value).toBe(b.id.value);
  });

  it("bySource returns empty for unknown actor", () => {
    const s = new InMemoryStore();
    const b = boot(); s.append(b);
    expect(s.bySource(new ActorId("bob"), 10)).toHaveLength(0);
  });

  it("byConversation returns matching events newest first", () => {
    const s = new InMemoryStore();
    const b = boot(); s.append(b);
    const e1 = next(b); s.append(e1);

    const results = s.byConversation(new ConversationId("conv_1"), 10);
    expect(results).toHaveLength(1);
    expect(results[0].id.value).toBe(e1.id.value);
  });

  it("byConversation returns empty for unknown conversation", () => {
    const s = new InMemoryStore();
    const b = boot(); s.append(b);
    expect(s.byConversation(new ConversationId("conv_unknown"), 10)).toHaveLength(0);
  });

  it("ancestors returns causal ancestors", () => {
    const s = new InMemoryStore();
    const b = boot(); s.append(b);
    const e1 = next(b); s.append(e1);
    const e2 = next(e1); s.append(e2);

    const anc = s.ancestors(e2.id, 10);
    expect(anc).toHaveLength(2);
    // e1 is direct cause, b is cause of e1
    const ids = anc.map((e) => e.id.value);
    expect(ids).toContain(e1.id.value);
    expect(ids).toContain(b.id.value);
  });

  it("ancestors respects maxDepth", () => {
    const s = new InMemoryStore();
    const b = boot(); s.append(b);
    const e1 = next(b); s.append(e1);
    const e2 = next(e1); s.append(e2);

    const anc = s.ancestors(e2.id, 1);
    expect(anc).toHaveLength(1);
    expect(anc[0].id.value).toBe(e1.id.value);
  });

  it("ancestors throws for unknown event", () => {
    const s = new InMemoryStore();
    expect(() => s.ancestors(new EventId("019462a0-0000-7000-8000-000000000099"), 10)).toThrow(EventNotFoundError);
  });

  it("descendants returns causal descendants", () => {
    const s = new InMemoryStore();
    const b = boot(); s.append(b);
    const e1 = next(b); s.append(e1);
    const e2 = next(e1); s.append(e2);

    const desc = s.descendants(b.id, 10);
    // e1 has b in causes, e2 has e1 in causes
    const ids = desc.map((e) => e.id.value);
    expect(ids).toContain(e1.id.value);
    expect(ids).toContain(e2.id.value);
  });

  it("descendants respects maxDepth", () => {
    const s = new InMemoryStore();
    const b = boot(); s.append(b);
    const e1 = next(b); s.append(e1);
    const e2 = next(e1); s.append(e2);

    const desc = s.descendants(b.id, 1);
    expect(desc).toHaveLength(1);
    expect(desc[0].id.value).toBe(e1.id.value);
  });

  it("descendants throws for unknown event", () => {
    const s = new InMemoryStore();
    expect(() => s.descendants(new EventId("019462a0-0000-7000-8000-000000000099"), 10)).toThrow(EventNotFoundError);
  });
});
