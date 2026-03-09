import { describe, it, expect } from "vitest";
import {
  allCodeGraphEventTypes,
  allCodeGraphPrimitives,
  allCodeGraphCompositions,
  registerAllCodeGraph,
  isCodeGraphPrimitive,
  boardComposition,
  detailComposition,
  feedComposition,
  dashboardComposition,
  inboxComposition,
  wizardComposition,
  skinComposition,
} from "../src/codegraph.js";
import { Registry } from "../src/primitive.js";
import type { Snapshot } from "../src/primitive.js";
import type { Event } from "../src/event.js";
import { EventType, PrimitiveId } from "../src/types.js";

// ── Helpers ──────────────────────────────────────────────────────────────

const fakeSnapshot: Snapshot = {
  tick: 1,
  primitives: new Map(),
  pendingEvents: [],
  recentEvents: [],
};

function fakeEvent(typeStr: string): Event {
  return { type: new EventType(typeStr) } as unknown as Event;
}

// ── Code Graph Event Types ──────────────────────────────────────────────

describe("code graph event types", () => {
  it("returns 35 event types", () => {
    const types = allCodeGraphEventTypes();
    expect(types.length).toBe(35);
  });

  it("all event types are unique", () => {
    const types = allCodeGraphEventTypes();
    const values = new Set(types.map((t) => t.value));
    expect(values.size).toBe(35);
  });

  it("all event types start with codegraph.", () => {
    for (const t of allCodeGraphEventTypes()) {
      expect(t.value.startsWith("codegraph.")).toBe(true);
    }
  });
});

// ── Code Graph Primitives ───────────────────────────────────────────────

describe("code graph primitives", () => {
  const all = allCodeGraphPrimitives();

  it("returns exactly 61 primitives", () => {
    expect(all.length).toBe(61);
  });

  it("all IDs are unique", () => {
    const ids = new Set(all.map((p) => p.id().value));
    expect(ids.size).toBe(61);
  });

  it("all primitives are at Layer 5", () => {
    for (const p of all) {
      expect(p.layer().value).toBe(5);
    }
  });

  it("all primitives have cadence 1", () => {
    for (const p of all) {
      expect(p.cadence().value).toBe(1);
    }
  });

  it("all primitives have non-empty subscriptions", () => {
    for (const p of all) {
      const subs = p.subscriptions();
      expect(subs.length).toBeGreaterThan(0);
      for (const s of subs) {
        expect(s.value).toBeTruthy();
      }
    }
  });

  it("all primitive IDs start with CG", () => {
    for (const p of all) {
      expect(p.id().value.startsWith("CG")).toBe(true);
    }
  });

  it("isCodeGraphPrimitive recognises CG primitives", () => {
    for (const p of all) {
      expect(isCodeGraphPrimitive(p.id())).toBe(true);
    }
    expect(isCodeGraphPrimitive(new PrimitiveId("Event"))).toBe(false);
    expect(isCodeGraphPrimitive(new PrimitiveId("agent.Identity"))).toBe(false);
  });

  it("all 61 primitives process with empty events", () => {
    for (const p of all) {
      const mutations = p.process(1, [], fakeSnapshot);
      expect(mutations.length).toBeGreaterThan(0);
      const lastTickMutation = mutations.find(
        (m) => m.kind === "updateState" && m.key === "lastTick",
      );
      expect(lastTickMutation).toBeDefined();
      if (lastTickMutation && lastTickMutation.kind === "updateState") {
        expect(lastTickMutation.value).toBe(1);
      }
    }
  });

  it("process records tick correctly at tick 42", () => {
    for (const p of all) {
      const mutations = p.process(42, [], fakeSnapshot);
      const lastTickMutation = mutations.find(
        (m) => m.kind === "updateState" && m.key === "lastTick",
      );
      expect(lastTickMutation).toBeDefined();
      if (lastTickMutation && lastTickMutation.kind === "updateState") {
        expect(lastTickMutation.value).toBe(42);
      }
    }
  });

  it("process returns eventsProcessed mutation with correct count", () => {
    const p = all[0];
    const events = [
      fakeEvent("codegraph.entity.defined"),
      fakeEvent("codegraph.entity.modified"),
      fakeEvent("codegraph.entity.related"),
    ];
    const mutations = p.process(1, events, fakeSnapshot);
    const processed = mutations.find(
      (m) => m.kind === "updateState" && m.key === "eventsProcessed",
    );
    expect(processed).toBeDefined();
    if (processed && processed.kind === "updateState") {
      expect(processed.value).toBe(3);
    }
  });
});

// ── Registry ─────────────────────────────────────────────────────────────

describe("registerAllCodeGraph", () => {
  it("registers 61 primitives in the registry", () => {
    const reg = new Registry();
    registerAllCodeGraph(reg);
    expect(reg.count).toBe(61);
  });

  it("all registered primitives are active", () => {
    const reg = new Registry();
    registerAllCodeGraph(reg);
    for (const p of allCodeGraphPrimitives()) {
      expect(reg.getLifecycle(p.id())).toBe("active");
    }
  });

  it("double registration throws", () => {
    const reg = new Registry();
    registerAllCodeGraph(reg);
    expect(() => registerAllCodeGraph(reg)).toThrow("already registered");
  });
});

// ── Compositions ─────────────────────────────────────────────────────────

describe("code graph compositions", () => {
  it("returns 7 compositions", () => {
    const comps = allCodeGraphCompositions();
    expect(comps.length).toBe(7);
  });

  it("all compositions have unique names", () => {
    const names = new Set(allCodeGraphCompositions().map((c) => c.name));
    expect(names.size).toBe(7);
  });

  it("all compositions have non-empty primitives", () => {
    for (const c of allCodeGraphCompositions()) {
      expect(c.primitives.length).toBeGreaterThan(0);
    }
  });

  it("Board has 10 primitives", () => {
    const c = boardComposition();
    expect(c.name).toBe("Board");
    expect(c.primitives.length).toBe(10);
  });

  it("Detail has 9 primitives", () => {
    const c = detailComposition();
    expect(c.name).toBe("Detail");
    expect(c.primitives.length).toBe(9);
  });

  it("Feed has 7 primitives", () => {
    const c = feedComposition();
    expect(c.name).toBe("Feed");
    expect(c.primitives.length).toBe(7);
  });

  it("Dashboard has 5 primitives", () => {
    const c = dashboardComposition();
    expect(c.name).toBe("Dashboard");
    expect(c.primitives.length).toBe(5);
  });

  it("Inbox has 8 primitives", () => {
    const c = inboxComposition();
    expect(c.name).toBe("Inbox");
    expect(c.primitives.length).toBe(8);
  });

  it("Wizard has 8 primitives", () => {
    const c = wizardComposition();
    expect(c.name).toBe("Wizard");
    expect(c.primitives.length).toBe(8);
  });

  it("Skin has 7 primitives", () => {
    const c = skinComposition();
    expect(c.name).toBe("Skin");
    expect(c.primitives.length).toBe(7);
  });

  it("all composition primitive IDs reference valid code graph primitives", () => {
    const validIds = new Set(allCodeGraphPrimitives().map((p) => p.id().value));
    for (const c of allCodeGraphCompositions()) {
      for (const pid of c.primitives) {
        expect(validIds.has(pid)).toBe(true);
      }
    }
  });
});
