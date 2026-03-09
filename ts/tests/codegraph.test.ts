import { describe, it, expect } from "vitest";
import {
  allCodeGraphEventTypes,
  allCodeGraphCompositions,
  boardComposition,
  detailComposition,
  feedComposition,
  dashboardComposition,
  inboxComposition,
  wizardComposition,
  skinComposition,
} from "../src/codegraph.js";

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
});
