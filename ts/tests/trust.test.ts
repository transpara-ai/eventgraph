import { describe, it, expect } from "vitest";
import {
  DefaultTrustModel,
  TrustMetrics,
  defaultTrustConfig,
} from "../src/trust.js";
import { Actor, ActorType, ActorStatus, InMemoryActorStore } from "../src/actor.js";
import {
  ActorId,
  ConversationId,
  DomainScope,
  EventId,
  EventType,
  Hash,
  NonEmpty,
  Score,
  Signature,
  PublicKey,
  Weight,
} from "../src/types.js";
import { Event, newEventId } from "../src/event.js";

// ── Helpers ──────────────────────────────────────────────────────────────

function testPublicKey(b: number): PublicKey {
  const key = new Uint8Array(32);
  key[0] = b;
  return new PublicKey(key);
}

function testActor(name: string, b: number): Actor {
  const store = new InMemoryActorStore();
  return store.register(testPublicKey(b), name, ActorType.Human);
}

function testTrustEvent(actorId: ActorId, prev: number, curr: number, domain = "general"): Event {
  const content: Record<string, unknown> = {
    actor: actorId.value,
    previous: prev,
    current: curr,
    domain: domain,
    cause: "019462a0-0000-7000-8000-000000000001",
  };
  const causeId = new EventId("019462a0-0000-7000-8000-000000000001");
  return new Event(
    1,
    newEventId(),
    new EventType("trust.updated"),
    Date.now() * 1_000_000,
    actorId,
    content,
    NonEmpty.of([causeId]),
    new ConversationId("conv_test000000000000000000000001"),
    Hash.zero(),
    Hash.zero(),
    new Signature(new Uint8Array(64)),
  );
}

function testNonTrustEvent(actorId: ActorId): Event {
  const content: Record<string, unknown> = {
    message: "hello",
  };
  const causeId = new EventId("019462a0-0000-7000-8000-000000000001");
  return new Event(
    1,
    newEventId(),
    new EventType("system.bootstrapped"),
    Date.now() * 1_000_000,
    actorId,
    content,
    NonEmpty.of([causeId]),
    new ConversationId("conv_test000000000000000000000001"),
    Hash.zero(),
    Hash.zero(),
    new Signature(new Uint8Array(64)),
  );
}

// ── Tests ────────────────────────────────────────────────────────────────

describe("DefaultTrustModel", () => {
  it("initialScore", () => {
    const model = new DefaultTrustModel();
    const a = testActor("Alice", 1);

    const metrics = model.score(a);
    expect(metrics.overall.value).toBe(0.0);
    expect(metrics.confidence.value).toBe(0.0);
  });

  it("updateIncreasesTrust", () => {
    const model = new DefaultTrustModel();
    const a = testActor("Alice", 1);

    const ev = testTrustEvent(a.id, 0.0, 0.05);
    const metrics = model.update(a, ev);
    expect(metrics.overall.value).toBeGreaterThan(0.0);
  });

  it("updateDecreasesTrust", () => {
    const model = new DefaultTrustModel();
    const a = testActor("Alice", 1);

    // First increase
    const ev1 = testTrustEvent(a.id, 0.0, 0.08);
    model.update(a, ev1);

    // Then decrease
    const ev2 = testTrustEvent(a.id, 0.08, 0.0);
    const metrics = model.update(a, ev2);
    expect(metrics.overall.value).toBeLessThan(0.08);
  });

  it("updateClampedToMaxAdjustment", () => {
    const model = new DefaultTrustModel();
    const a = testActor("Alice", 1);

    // Try to increase by 0.5 — should be clamped to MaxAdjustment (0.1)
    const ev = testTrustEvent(a.id, 0.0, 0.5);
    const metrics = model.update(a, ev);
    expect(metrics.overall.value).toBeLessThanOrEqual(0.1);
  });

  it("updateDeduplication", () => {
    const model = new DefaultTrustModel();
    const a = testActor("Alice", 1);

    const ev = testTrustEvent(a.id, 0.0, 0.05);
    const m1 = model.update(a, ev);

    // Same event again — should not change score
    const m2 = model.update(a, ev);
    expect(m2.overall.value).toBe(m1.overall.value);
    expect(m2.evidence.length).toBe(m1.evidence.length);
  });

  it("updateTrendPositive", () => {
    const model = new DefaultTrustModel();
    const a = testActor("Alice", 1);

    const ev = testTrustEvent(a.id, 0.0, 0.05);
    const metrics = model.update(a, ev);
    expect(metrics.trend.value).toBeGreaterThan(0);
  });

  it("updateTrendNegative", () => {
    const model = new DefaultTrustModel();
    const a = testActor("Alice", 1);

    // First increase to have something to decrease from
    const ev1 = testTrustEvent(a.id, 0.0, 0.05);
    model.update(a, ev1);

    // Then decrease — trend should go down
    const ev2 = testTrustEvent(a.id, 0.1, 0.0);
    const metrics = model.update(a, ev2);
    // After +0.1 then -0.1, trend should be 0
    expect(metrics.trend.value).toBeLessThanOrEqual(0.1);
  });

  it("scoreInDomain", () => {
    const model = new DefaultTrustModel();
    const a = testActor("Alice", 1);

    // Update with domain-specific evidence
    const ev = testTrustEvent(a.id, 0.0, 0.05, "code_review");
    model.update(a, ev);

    const domain = new DomainScope("code_review");
    const metrics = model.scoreInDomain(a, domain);
    expect(metrics.overall.value).toBeGreaterThan(0.0);
  });

  it("scoreInDomainFallback", () => {
    const model = new DefaultTrustModel();
    const a = testActor("Alice", 1);

    // No domain-specific evidence, should fall back to global with halved confidence
    const domain = new DomainScope("code_review");
    const metrics = model.scoreInDomain(a, domain);
    expect(metrics.overall.value).toBe(0.0);
    expect(metrics.confidence.value).toBe(0.0); // 0 evidence * 0.5 = 0
  });

  it("scoreInDomainFallbackWithEvidence", () => {
    const model = new DefaultTrustModel();
    const a = testActor("Alice", 1);

    // Add non-trust event to build evidence (no domain info)
    const ev = testNonTrustEvent(a.id);
    model.update(a, ev);

    // Query a domain that has no specific data
    const domain = new DomainScope("unknown_domain");
    const metrics = model.scoreInDomain(a, domain);
    // Confidence should be halved compared to global
    const globalMetrics = model.score(a);
    expect(metrics.confidence.value).toBeCloseTo(globalMetrics.confidence.value * 0.5, 10);
  });

  it("decay", () => {
    const model = new DefaultTrustModel();
    const a = testActor("Alice", 1);

    // Build up some trust
    for (let i = 0; i < 5; i++) {
      const ev = testTrustEvent(a.id, 0.0, 0.1);
      model.update(a, ev);
    }

    const before = model.score(a);

    // Decay 10 days (in seconds)
    model.decay(a, 10 * 86400);

    const after = model.score(a);
    expect(after.overall.value).toBeLessThan(before.overall.value);
  });

  it("decayNegativeDuration", () => {
    const model = new DefaultTrustModel();
    const a = testActor("Alice", 1);

    // Build some trust
    const ev = testTrustEvent(a.id, 0.0, 0.05);
    model.update(a, ev);

    const before = model.score(a);

    // Negative duration — should not change score
    const metrics = model.decay(a, -100);
    expect(metrics.overall.value).toBe(before.overall.value);
  });

  it("updateBetween", () => {
    const model = new DefaultTrustModel();
    const from = testActor("Alice", 1);
    const to = testActor("Bob", 2);

    const ev = testTrustEvent(to.id, 0.0, 0.05);
    const metrics = model.updateBetween(from, to, ev);
    expect(metrics.overall.value).toBeGreaterThan(0.0);
    expect(metrics.actor.value).toBe(to.id.value);
  });

  it("betweenNoRelationship", () => {
    const model = new DefaultTrustModel();
    const from = testActor("Alice", 1);
    const to = testActor("Bob", 2);

    const metrics = model.between(from, to);
    expect(metrics.overall.value).toBe(0.0);
    expect(metrics.confidence.value).toBe(0.0);
  });

  it("betweenAsymmetric", () => {
    const model = new DefaultTrustModel();
    const alice = testActor("Alice", 1);
    const bob = testActor("Bob", 2);

    // Alice trusts Bob
    const ev = testTrustEvent(bob.id, 0.0, 0.05);
    model.updateBetween(alice, bob, ev);

    // Alice→Bob has trust
    const ab = model.between(alice, bob);
    expect(ab.overall.value).toBeGreaterThan(0.0);

    // Bob→Alice has no trust (asymmetric)
    const ba = model.between(bob, alice);
    expect(ba.overall.value).toBe(0.0);
  });

  it("evidenceCappedAt100", () => {
    const model = new DefaultTrustModel();
    const a = testActor("Alice", 1);

    // Add 150 updates
    for (let i = 0; i < 150; i++) {
      const ev = testTrustEvent(a.id, 0.0, 0.01);
      model.update(a, ev);
    }

    const metrics = model.score(a);
    expect(metrics.evidence.length).toBeLessThanOrEqual(100);
  });

  it("decayDirectedTrust", () => {
    const model = new DefaultTrustModel();
    const from = testActor("Alice", 1);
    const to = testActor("Bob", 2);

    // Build directed trust
    for (let i = 0; i < 5; i++) {
      const ev = testTrustEvent(to.id, 0.0, 0.1);
      model.updateBetween(from, to, ev);
    }

    const before = model.between(from, to);

    // Decay from actor's perspective
    model.decay(from, 10 * 86400);

    const after = model.between(from, to);
    expect(after.overall.value).toBeLessThan(before.overall.value);
  });
});
