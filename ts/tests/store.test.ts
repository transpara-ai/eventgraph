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
});
