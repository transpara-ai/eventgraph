import { describe, it, expect } from "vitest";
import {
  Actor,
  ActorType,
  ActorStatus,
  InMemoryActorStore,
  transitionTo,
  validTransitions,
} from "../src/actor.js";
import { ActorId, EventId, PublicKey } from "../src/types.js";
import {
  ActorNotFoundError,
  ActorKeyNotFoundError,
  InvalidTransitionError,
} from "../src/errors.js";

function testPublicKey(b: number): PublicKey {
  const key = new Uint8Array(32);
  key[0] = b;
  return new PublicKey(key);
}

const reason = new EventId("019462a0-0000-7000-8000-000000000001");

describe("InMemoryActorStore", () => {
  it("Register", () => {
    const s = new InMemoryActorStore();
    const pk = testPublicKey(1);
    const a = s.register(pk, "Alice", ActorType.Human);
    expect(a.displayName).toBe("Alice");
    expect(a.actorType).toBe(ActorType.Human);
    expect(a.status).toBe(ActorStatus.Active);
  });

  it("RegisterIdempotent", () => {
    const s = new InMemoryActorStore();
    const pk = testPublicKey(1);
    const a1 = s.register(pk, "Alice", ActorType.Human);
    const a2 = s.register(pk, "Alice Again", ActorType.Human);
    expect(a1.id.value).toBe(a2.id.value);
    expect(s.actorCount()).toBe(1);
  });

  it("Get", () => {
    const s = new InMemoryActorStore();
    const pk = testPublicKey(1);
    const a = s.register(pk, "Alice", ActorType.Human);
    const got = s.get(a.id);
    expect(got.id.value).toBe(a.id.value);
  });

  it("GetNotFound", () => {
    const s = new InMemoryActorStore();
    expect(() => s.get(new ActorId("actor_nonexistent"))).toThrow(ActorNotFoundError);
  });

  it("GetByPublicKey", () => {
    const s = new InMemoryActorStore();
    const pk = testPublicKey(1);
    const a = s.register(pk, "Alice", ActorType.Human);
    const got = s.getByPublicKey(pk);
    expect(got.id.value).toBe(a.id.value);
  });

  it("GetByPublicKeyNotFound", () => {
    const s = new InMemoryActorStore();
    expect(() => s.getByPublicKey(testPublicKey(99))).toThrow(ActorKeyNotFoundError);
  });

  it("Update", () => {
    const s = new InMemoryActorStore();
    const pk = testPublicKey(1);
    const a = s.register(pk, "Alice", ActorType.Human);
    const updated = s.update(a.id, { displayName: "Alice Updated" });
    expect(updated.displayName).toBe("Alice Updated");
  });

  it("UpdateMetadataMerge", () => {
    const s = new InMemoryActorStore();
    const pk = testPublicKey(1);
    const a = s.register(pk, "Alice", ActorType.Human);

    // Set initial metadata
    s.update(a.id, { metadata: { role: "builder" } });

    // Merge additional metadata
    const updated = s.update(a.id, { metadata: { team: "core" } });
    const md = updated.metadata;
    expect(md["role"]).toBe("builder");
    expect(md["team"]).toBe("core");
  });

  it("UpdateNotFound", () => {
    const s = new InMemoryActorStore();
    expect(() => s.update(new ActorId("actor_nonexistent"), {})).toThrow(ActorNotFoundError);
  });

  it("Suspend", () => {
    const s = new InMemoryActorStore();
    const pk = testPublicKey(1);
    const a = s.register(pk, "Alice", ActorType.Human);
    const suspended = s.suspend(a.id, reason);
    expect(suspended.status).toBe(ActorStatus.Suspended);
  });

  it("SuspendNotFound", () => {
    const s = new InMemoryActorStore();
    expect(() => s.suspend(new ActorId("actor_nonexistent"), reason)).toThrow(ActorNotFoundError);
  });

  it("SuspendAndReactivate", () => {
    const s = new InMemoryActorStore();
    const pk = testPublicKey(1);
    const a = s.register(pk, "Alice", ActorType.Human);
    s.suspend(a.id, reason);

    const got = s.get(a.id);
    expect(got.status).toBe(ActorStatus.Suspended);

    const reactivated = s.reactivate(a.id, reason);
    expect(reactivated.status).toBe(ActorStatus.Active);
  });

  it("Memorial", () => {
    const s = new InMemoryActorStore();
    const pk = testPublicKey(1);
    const a = s.register(pk, "Alice", ActorType.Human);
    const memorial = s.memorial(a.id, reason);
    expect(memorial.status).toBe(ActorStatus.Memorial);
  });

  it("MemorialIsTerminal", () => {
    const s = new InMemoryActorStore();
    const pk = testPublicKey(1);
    const a = s.register(pk, "Alice", ActorType.Human);
    s.memorial(a.id, reason);

    // Try to suspend — should fail
    expect(() => s.suspend(a.id, reason)).toThrow(InvalidTransitionError);
  });

  it("MemorialReactivateIsError", () => {
    const s = new InMemoryActorStore();
    const pk = testPublicKey(1);
    const a = s.register(pk, "Alice", ActorType.Human);
    s.memorial(a.id, reason);

    expect(() => s.reactivate(a.id, reason)).toThrow(InvalidTransitionError);
  });

  it("ReactivateFromActiveIsError", () => {
    const s = new InMemoryActorStore();
    const pk = testPublicKey(1);
    const a = s.register(pk, "Alice", ActorType.Human);
    expect(() => s.reactivate(a.id, reason)).toThrow(InvalidTransitionError);
  });

  it("ReactivateNotFound", () => {
    const s = new InMemoryActorStore();
    expect(() => s.reactivate(new ActorId("actor_nonexistent"), reason)).toThrow(ActorNotFoundError);
  });

  it("List", () => {
    const s = new InMemoryActorStore();
    for (let i = 1; i <= 5; i++) {
      s.register(testPublicKey(i), "Actor", ActorType.Human);
    }
    const page = s.list({ limit: 10 });
    expect(page.items.length).toBe(5);
  });

  it("ListWithStatusFilter", () => {
    const s = new InMemoryActorStore();
    s.register(testPublicKey(1), "Active1", ActorType.Human);
    s.register(testPublicKey(2), "Active2", ActorType.Human);
    const a3 = s.register(testPublicKey(3), "ToBeSuspended", ActorType.Human);

    s.suspend(a3.id, reason);

    const activePage = s.list({ status: ActorStatus.Active, limit: 10 });
    expect(activePage.items.length).toBe(2);
  });

  it("ListWithTypeFilter", () => {
    const s = new InMemoryActorStore();
    s.register(testPublicKey(1), "Human", ActorType.Human);
    s.register(testPublicKey(2), "AI", ActorType.AI);
    s.register(testPublicKey(3), "System", ActorType.System);

    const page = s.list({ actorType: ActorType.AI, limit: 10 });
    expect(page.items.length).toBe(1);
  });

  it("ListPagination", () => {
    const s = new InMemoryActorStore();
    for (let i = 1; i <= 5; i++) {
      s.register(testPublicKey(i), "Actor", ActorType.Human);
    }

    // Page 1
    const page1 = s.list({ limit: 2 });
    expect(page1.items.length).toBe(2);
    expect(page1.hasMore).toBe(true);
    expect(page1.cursor).toBeDefined();

    // Page 2
    const page2 = s.list({ limit: 2, after: page1.cursor });
    expect(page2.items.length).toBe(2);
    expect(page2.hasMore).toBe(true);

    // Page 3
    const page3 = s.list({ limit: 2, after: page2.cursor });
    expect(page3.items.length).toBe(1);
    expect(page3.hasMore).toBe(false);
  });

  it("ActorGetters", () => {
    const s = new InMemoryActorStore();
    const pk = testPublicKey(42);
    const a = s.register(pk, "Alice", ActorType.AI);

    expect(a.id.value).toBeTruthy();
    expect(a.publicKey).toBeDefined();
    expect(a.createdAt).toBeGreaterThan(0);
    expect(a.metadata).toBeDefined();
    expect(typeof a.metadata).toBe("object");
  });
});

describe("ActorStatus transitions", () => {
  it("active to suspended", () => {
    expect(transitionTo(ActorStatus.Active, ActorStatus.Suspended)).toBe(ActorStatus.Suspended);
  });

  it("active to memorial", () => {
    expect(transitionTo(ActorStatus.Active, ActorStatus.Memorial)).toBe(ActorStatus.Memorial);
  });

  it("suspended to active", () => {
    expect(transitionTo(ActorStatus.Suspended, ActorStatus.Active)).toBe(ActorStatus.Active);
  });

  it("suspended to memorial", () => {
    expect(transitionTo(ActorStatus.Suspended, ActorStatus.Memorial)).toBe(ActorStatus.Memorial);
  });

  it("memorial to anything throws", () => {
    expect(() => transitionTo(ActorStatus.Memorial, ActorStatus.Active)).toThrow(InvalidTransitionError);
    expect(() => transitionTo(ActorStatus.Memorial, ActorStatus.Suspended)).toThrow(InvalidTransitionError);
  });

  it("active to active throws", () => {
    expect(() => transitionTo(ActorStatus.Active, ActorStatus.Active)).toThrow(InvalidTransitionError);
  });

  it("validTransitions returns correct targets", () => {
    expect(validTransitions(ActorStatus.Active)).toEqual([ActorStatus.Suspended, ActorStatus.Memorial]);
    expect(validTransitions(ActorStatus.Suspended)).toEqual([ActorStatus.Active, ActorStatus.Memorial]);
    expect(validTransitions(ActorStatus.Memorial)).toEqual([]);
  });
});

describe("Actor immutability", () => {
  it("metadata is deep copied on construction and access", () => {
    const md = { key: "value", nested: { a: 1 } };
    const actor = new Actor(
      new ActorId("actor_test"),
      testPublicKey(1),
      "Test",
      ActorType.Human,
      md,
      Date.now() * 1_000_000,
      ActorStatus.Active,
    );

    // Mutating original should not affect actor
    md.key = "changed";
    expect(actor.metadata["key"]).toBe("value");

    // Mutating returned metadata should not affect actor
    const got = actor.metadata;
    got["key"] = "hacked";
    expect(actor.metadata["key"]).toBe("value");
  });
});
