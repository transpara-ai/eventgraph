import { describe, it, expect } from "vitest";
import {
  Score, Weight, Activation, Layer, Cadence,
  EventId, Hash, EventType, ActorId, SubscriptionPattern,
  Option, NonEmpty, PublicKey, Signature,
} from "../src/types.js";
import { OutOfRangeError, InvalidFormatError, EmptyRequiredError } from "../src/errors.js";

describe("Score", () => {
  it("accepts valid range", () => { expect(new Score(0).value).toBe(0); expect(new Score(0.5).value).toBe(0.5); expect(new Score(1).value).toBe(1); });
  it("rejects below range", () => expect(() => new Score(-0.1)).toThrow(OutOfRangeError));
  it("rejects above range", () => expect(() => new Score(1.1)).toThrow(OutOfRangeError));
  it("rejects NaN", () => expect(() => new Score(NaN)).toThrow(OutOfRangeError));
});

describe("Weight", () => {
  it("accepts valid range", () => { expect(new Weight(-1).value).toBe(-1); expect(new Weight(1).value).toBe(1); });
  it("rejects below range", () => expect(() => new Weight(-1.1)).toThrow(OutOfRangeError));
  it("rejects above range", () => expect(() => new Weight(1.1)).toThrow(OutOfRangeError));
});

describe("Activation", () => {
  it("accepts valid", () => { expect(new Activation(0).value).toBe(0); expect(new Activation(1).value).toBe(1); });
  it("rejects below", () => expect(() => new Activation(-0.01)).toThrow(OutOfRangeError));
});

describe("Layer", () => {
  it("accepts valid", () => { expect(new Layer(0).value).toBe(0); expect(new Layer(13).value).toBe(13); });
  it("rejects below", () => expect(() => new Layer(-1)).toThrow(OutOfRangeError));
  it("rejects above", () => expect(() => new Layer(14)).toThrow(OutOfRangeError));
});

describe("Cadence", () => {
  it("accepts valid", () => { expect(new Cadence(1).value).toBe(1); expect(new Cadence(100).value).toBe(100); });
  it("rejects zero", () => expect(() => new Cadence(0)).toThrow(OutOfRangeError));
});

describe("EventId", () => {
  it("accepts valid UUID v7", () => expect(new EventId("019462a0-0000-7000-8000-000000000001").value).toBe("019462a0-0000-7000-8000-000000000001"));
  it("normalizes to lowercase", () => expect(new EventId("019462A0-0000-7000-8000-000000000001").value).toBe("019462a0-0000-7000-8000-000000000001"));
  it("rejects non-v7", () => expect(() => new EventId("019462a0-0000-4000-8000-000000000001")).toThrow(InvalidFormatError));
  it("rejects garbage", () => expect(() => new EventId("not-a-uuid")).toThrow(InvalidFormatError));
});

describe("Hash", () => {
  it("accepts valid 64 hex", () => expect(new Hash("a".repeat(64)).value).toBe("a".repeat(64)));
  it("zero hash", () => { const h = Hash.zero(); expect(h.isZero).toBe(true); expect(h.value).toBe("0".repeat(64)); });
  it("rejects short", () => expect(() => new Hash("abc")).toThrow(InvalidFormatError));
  it("rejects empty", () => expect(() => new Hash("")).toThrow(InvalidFormatError));
});

describe("EventType", () => {
  it("accepts valid", () => expect(new EventType("trust.updated").value).toBe("trust.updated"));
  it("rejects uppercase", () => expect(() => new EventType("Trust.Updated")).toThrow(InvalidFormatError));
  it("rejects empty", () => expect(() => new EventType("")).toThrow(InvalidFormatError));
});

describe("ActorId", () => {
  it("accepts valid", () => expect(new ActorId("actor_alice").value).toBe("actor_alice"));
  it("rejects empty", () => expect(() => new ActorId("")).toThrow(EmptyRequiredError));
});

describe("SubscriptionPattern", () => {
  it("wildcard matches all", () => {
    const sp = new SubscriptionPattern("*");
    expect(sp.matches(new EventType("trust.updated"))).toBe(true);
    expect(sp.matches(new EventType("system.bootstrapped"))).toBe(true);
  });
  it("prefix match", () => {
    const sp = new SubscriptionPattern("trust.*");
    expect(sp.matches(new EventType("trust.updated"))).toBe(true);
    expect(sp.matches(new EventType("system.bootstrapped"))).toBe(false);
  });
  it("exact match", () => {
    const sp = new SubscriptionPattern("trust.updated");
    expect(sp.matches(new EventType("trust.updated"))).toBe(true);
    expect(sp.matches(new EventType("trust.decayed"))).toBe(false);
  });
});

describe("Option", () => {
  it("some", () => { const o = Option.some(42); expect(o.isSome).toBe(true); expect(o.unwrap()).toBe(42); });
  it("none", () => { const o = Option.none<number>(); expect(o.isNone).toBe(true); expect(() => o.unwrap()).toThrow(); });
  it("unwrapOr", () => { expect(Option.some(42).unwrapOr(0)).toBe(42); expect(Option.none<number>().unwrapOr(0)).toBe(0); });
});

describe("NonEmpty", () => {
  it("valid", () => { const ne = NonEmpty.of([1, 2, 3]); expect(ne.length).toBe(3); expect(ne.get(0)).toBe(1); });
  it("rejects empty", () => expect(() => NonEmpty.of([])).toThrow());
});

describe("PublicKey", () => {
  it("accepts 32 bytes", () => expect(new PublicKey(new Uint8Array(32)).bytes.length).toBe(32));
  it("rejects wrong length", () => expect(() => new PublicKey(new Uint8Array(31))).toThrow(InvalidFormatError));
});

describe("Signature", () => {
  it("accepts 64 bytes", () => expect(new Signature(new Uint8Array(64)).bytes.length).toBe(64));
  it("rejects wrong length", () => expect(() => new Signature(new Uint8Array(63))).toThrow(InvalidFormatError));
});
