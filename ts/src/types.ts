import { OutOfRangeError, EmptyRequiredError, InvalidFormatError } from "./errors.js";

// ── Option<T> ───────────────────────────────────────────────────────────

export class Option<T> {
  private constructor(
    private readonly _value: T | undefined,
    private readonly _present: boolean,
  ) {}

  static some<T>(value: T): Option<T> { return new Option(value, true); }
  static none<T>(): Option<T> { return new Option<T>(undefined, false); }

  get isSome(): boolean { return this._present; }
  get isNone(): boolean { return !this._present; }

  unwrap(): T {
    if (!this._present) throw new Error("Unwrap called on None Option");
    return this._value!;
  }

  unwrapOr(defaultValue: T): T {
    return this._present ? this._value! : defaultValue;
  }
}

// ── NonEmpty<T> ─────────────────────────────────────────────────────────

export class NonEmpty<T> {
  private readonly _items: readonly T[];
  private constructor(items: readonly T[]) { this._items = items; }

  static of<T>(items: readonly T[]): NonEmpty<T> {
    if (items.length === 0) throw new Error("NonEmpty requires at least one element");
    return new NonEmpty([...items]);
  }

  get length(): number { return this._items.length; }
  get(index: number): T { return this._items[index]; }
  [Symbol.iterator](): Iterator<T> { return this._items[Symbol.iterator](); }
  toArray(): readonly T[] { return [...this._items]; }
}

// ── Constrained numerics ────────────────────────────────────────────────

export class Score {
  readonly value: number;
  constructor(value: number) {
    if (Number.isNaN(value) || value < 0.0 || value > 1.0)
      throw new OutOfRangeError("Score", value, 0.0, 1.0);
    this.value = Object.is(value, -0) ? 0.0 : value;
    Object.freeze(this);
  }
}

export class Weight {
  readonly value: number;
  constructor(value: number) {
    if (Number.isNaN(value) || value < -1.0 || value > 1.0)
      throw new OutOfRangeError("Weight", value, -1.0, 1.0);
    this.value = Object.is(value, -0) ? 0.0 : value;
    Object.freeze(this);
  }
}

export class Activation {
  readonly value: number;
  constructor(value: number) {
    if (Number.isNaN(value) || value < 0.0 || value > 1.0)
      throw new OutOfRangeError("Activation", value, 0.0, 1.0);
    this.value = Object.is(value, -0) ? 0.0 : value;
    Object.freeze(this);
  }
}

export class Layer {
  readonly value: number;
  constructor(value: number) {
    if (value < 0 || value > 13)
      throw new OutOfRangeError("Layer", value, 0, 13);
    this.value = value;
    Object.freeze(this);
  }
}

export class Cadence {
  readonly value: number;
  constructor(value: number) {
    if (value < 1)
      throw new OutOfRangeError("Cadence", value, 1, Infinity);
    this.value = value;
    Object.freeze(this);
  }
}

// ── Regex patterns ──────────────────────────────────────────────────────

const UUID_V7_RE = /^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/;
const UUID_RE = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/;
const EVENT_TYPE_RE = /^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$/;
const DOMAIN_SCOPE_RE = /^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)*$/;
const SUBSCRIPTION_RE = /^(\*|[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*(\.\*)?)$/;

// ── Typed IDs ───────────────────────────────────────────────────────────

export class EventId {
  readonly value: string;
  constructor(value: string) {
    value = value.toLowerCase();
    if (!UUID_V7_RE.test(value))
      throw new InvalidFormatError("EventId", value, "UUID v7");
    this.value = value;
    Object.freeze(this);
  }
  toString(): string { return this.value; }
}

export class EdgeId {
  readonly value: string;
  constructor(value: string) {
    value = value.toLowerCase();
    if (!UUID_V7_RE.test(value))
      throw new InvalidFormatError("EdgeId", value, "UUID v7");
    this.value = value;
    Object.freeze(this);
  }
  toString(): string { return this.value; }
}

export class Hash {
  readonly value: string;
  constructor(value: string) {
    value = value.toLowerCase();
    if (value.length !== 64 || !/^[0-9a-f]{64}$/.test(value))
      throw new InvalidFormatError("Hash", value, "64 hex characters (SHA-256)");
    this.value = value;
    Object.freeze(this);
  }

  static zero(): Hash { return new Hash("0".repeat(64)); }
  get isZero(): boolean { return this.value === "0".repeat(64); }
  toString(): string { return this.value; }
}

export class EventType {
  readonly value: string;
  constructor(value: string) {
    if (!EVENT_TYPE_RE.test(value))
      throw new InvalidFormatError("EventType", value, "dot-separated lowercase segments");
    this.value = value;
    Object.freeze(this);
  }
  toString(): string { return this.value; }
}

function makeStringId(name: string) {
  return class {
    readonly value: string;
    constructor(value: string) {
      if (!value) throw new EmptyRequiredError(name);
      this.value = value;
      Object.freeze(this);
    }
    toString(): string { return this.value; }
  };
}

export const ActorId = makeStringId("ActorId");
export type ActorId = InstanceType<typeof ActorId>;

export const ConversationId = makeStringId("ConversationId");
export type ConversationId = InstanceType<typeof ConversationId>;

export const SystemUri = makeStringId("SystemUri");
export type SystemUri = InstanceType<typeof SystemUri>;

export const PrimitiveId = makeStringId("PrimitiveId");
export type PrimitiveId = InstanceType<typeof PrimitiveId>;

export class DomainScope {
  readonly value: string;
  constructor(value: string) {
    if (!DOMAIN_SCOPE_RE.test(value))
      throw new InvalidFormatError("DomainScope", value, "lowercase dot/underscore-separated namespace");
    this.value = value;
    Object.freeze(this);
  }
  toString(): string { return this.value; }
}

export class SubscriptionPattern {
  readonly value: string;
  constructor(value: string) {
    if (!SUBSCRIPTION_RE.test(value))
      throw new InvalidFormatError("SubscriptionPattern", value, "dot-separated with optional trailing .* or bare *");
    this.value = value;
    Object.freeze(this);
  }

  matches(et: EventType): boolean {
    if (this.value === "*") return true;
    if (this.value.endsWith(".*")) {
      const prefix = this.value.slice(0, -2);
      return et.value === prefix || et.value.startsWith(prefix + ".");
    }
    return this.value === et.value;
  }
}

export class EnvelopeId {
  readonly value: string;
  constructor(value: string) {
    value = value.toLowerCase();
    if (!UUID_RE.test(value))
      throw new InvalidFormatError("EnvelopeId", value, "UUID");
    this.value = value;
    Object.freeze(this);
  }
  toString(): string { return this.value; }
}

export class TreatyId {
  readonly value: string;
  constructor(value: string) {
    value = value.toLowerCase();
    if (!UUID_RE.test(value))
      throw new InvalidFormatError("TreatyId", value, "UUID");
    this.value = value;
    Object.freeze(this);
  }
  toString(): string { return this.value; }
}

export class PublicKey {
  readonly bytes: Uint8Array;
  constructor(value: Uint8Array) {
    if (value.length !== 32)
      throw new InvalidFormatError("PublicKey", Buffer.from(value).toString("hex"), "32 bytes (Ed25519 public key)");
    this.bytes = new Uint8Array(value);
    Object.freeze(this);
  }
  toString(): string { return Buffer.from(this.bytes).toString("hex"); }
}

export class Signature {
  readonly bytes: Uint8Array;
  constructor(value: Uint8Array) {
    if (value.length !== 64)
      throw new InvalidFormatError("Signature", Buffer.from(value).toString("hex"), "64 bytes (Ed25519 signature)");
    this.bytes = new Uint8Array(value);
    Object.freeze(this);
  }
  toString(): string { return Buffer.from(this.bytes).toString("hex"); }
}
