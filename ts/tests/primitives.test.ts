import { describe, it, expect } from "vitest";
import { createAllPrimitives } from "../src/primitives.js";
import type { Primitive, Snapshot } from "../src/primitive.js";
import type { Event } from "../src/event.js";

describe("primitives", () => {
  const all = createAllPrimitives();

  it("createAllPrimitives returns 201 primitives", () => {
    expect(all.length).toBe(201);
  });

  it("all 201 primitives instantiate with valid id", () => {
    const ids = new Set<string>();
    for (const p of all) {
      const id = p.id();
      expect(id.value).toBeTruthy();
      ids.add(id.value);
    }
    // All unique
    expect(ids.size).toBe(201);
  });

  it("Layer 0 has 45 primitives", () => {
    expect(all.filter((p) => p.layer().value === 0).length).toBe(45);
  });

  it("Layer 1 has 12 primitives", () => {
    expect(all.filter((p) => p.layer().value === 1).length).toBe(12);
  });

  it("Layer 2 has 12 primitives", () => {
    expect(all.filter((p) => p.layer().value === 2).length).toBe(12);
  });

  it("Layer 3 has 12 primitives", () => {
    expect(all.filter((p) => p.layer().value === 3).length).toBe(12);
  });

  it("Layer 4 has 12 primitives", () => {
    expect(all.filter((p) => p.layer().value === 4).length).toBe(12);
  });

  it("Layer 5 has 12 primitives", () => {
    expect(all.filter((p) => p.layer().value === 5).length).toBe(12);
  });

  it("Layer 6 has 12 primitives", () => {
    expect(all.filter((p) => p.layer().value === 6).length).toBe(12);
  });

  it("Layer 7 has 12 primitives", () => {
    expect(all.filter((p) => p.layer().value === 7).length).toBe(12);
  });

  it("Layer 8 has 12 primitives", () => {
    expect(all.filter((p) => p.layer().value === 8).length).toBe(12);
  });

  it("Layer 9 has 12 primitives", () => {
    expect(all.filter((p) => p.layer().value === 9).length).toBe(12);
  });

  it("Layer 10 has 12 primitives", () => {
    expect(all.filter((p) => p.layer().value === 10).length).toBe(12);
  });

  it("Layer 11 has 12 primitives", () => {
    expect(all.filter((p) => p.layer().value === 11).length).toBe(12);
  });

  it("Layer 12 has 12 primitives", () => {
    expect(all.filter((p) => p.layer().value === 12).length).toBe(12);
  });

  it("Layer 13 has 12 primitives", () => {
    expect(all.filter((p) => p.layer().value === 13).length).toBe(12);
  });

  it("every primitive has non-empty subscriptions", () => {
    for (const p of all) {
      const subs = p.subscriptions();
      expect(subs.length).toBeGreaterThan(0);
      for (const s of subs) {
        expect(s.value).toBeTruthy();
      }
    }
  });

  it("every primitive has cadence >= 1", () => {
    for (const p of all) {
      expect(p.cadence().value).toBeGreaterThanOrEqual(1);
    }
  });

  it("process returns mutations with correct structure", () => {
    const fakeSnapshot: Snapshot = {
      tick: 1,
      primitives: new Map(),
      pendingEvents: [],
      recentEvents: [],
    };

    for (const p of all) {
      const mutations = p.process(1, [], fakeSnapshot);
      expect(mutations.length).toBe(2);

      const m0 = mutations[0];
      expect(m0.kind).toBe("updateState");
      if (m0.kind === "updateState") {
        expect(m0.key).toBe("eventsProcessed");
        expect(m0.value).toBe(0);
      }

      const m1 = mutations[1];
      expect(m1.kind).toBe("updateState");
      if (m1.kind === "updateState") {
        expect(m1.key).toBe("lastTick");
        expect(m1.value).toBe(1);
      }
    }
  });

  it("process records event count correctly", () => {
    const fakeSnapshot: Snapshot = {
      tick: 5,
      primitives: new Map(),
      pendingEvents: [],
      recentEvents: [],
    };

    // Use a non-empty events array (cast to avoid needing real Event objects)
    const fakeEvents = [{}, {}, {}] as unknown as Event[];
    const p = all[0];
    const mutations = p.process(5, fakeEvents, fakeSnapshot);

    if (mutations[0].kind === "updateState") {
      expect(mutations[0].value).toBe(3);
    }
    if (mutations[1].kind === "updateState") {
      expect(mutations[1].value).toBe(5);
    }
  });
});
