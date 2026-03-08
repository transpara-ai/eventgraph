import { Event } from "./event.js";
import { EventType, SubscriptionPattern } from "./types.js";
import { InMemoryStore } from "./store.js";

interface Sub {
  id: number;
  pattern: SubscriptionPattern;
  handler: (event: Event) => void;
}

export class EventBus {
  private readonly _store: InMemoryStore;
  private readonly subs = new Map<number, Sub>();
  private nextId = 0;
  private closed = false;

  constructor(store: InMemoryStore) {
    this._store = store;
  }

  get store(): InMemoryStore { return this._store; }

  subscribe(pattern: SubscriptionPattern, handler: (event: Event) => void): number {
    if (this.closed) return -1;
    const id = ++this.nextId;
    this.subs.set(id, { id, pattern, handler });
    return id;
  }

  unsubscribe(subId: number): void {
    this.subs.delete(subId);
  }

  publish(event: Event): void {
    if (this.closed) return;
    for (const sub of this.subs.values()) {
      if (sub.pattern.matches(event.type)) {
        try { sub.handler(event); } catch { /* swallow */ }
      }
    }
  }

  close(): void {
    this.closed = true;
    this.subs.clear();
  }
}
