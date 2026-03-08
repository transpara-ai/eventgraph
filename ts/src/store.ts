import { ChainIntegrityError, EventNotFoundError } from "./errors.js";
import { Event } from "./event.js";
import { EventId, Hash, Option } from "./types.js";

export interface ChainVerification {
  valid: boolean;
  length: number;
}

export interface Store {
  append(event: Event): Event;
  get(eventId: EventId): Event;
  head(): Option<Event>;
  count(): number;
  verifyChain(): ChainVerification;
  close(): void;
}

export class InMemoryStore implements Store {
  private readonly events: Event[] = [];
  private readonly index = new Map<string, number>();

  append(event: Event): Event {
    if (this.events.length > 0) {
      const last = this.events[this.events.length - 1];
      if (event.prevHash.value !== last.hash.value) {
        throw new ChainIntegrityError(
          this.events.length,
          `prev_hash ${event.prevHash.value} != head hash ${last.hash.value}`,
        );
      }
    }
    this.events.push(event);
    this.index.set(event.id.value, this.events.length - 1);
    return event;
  }

  get(eventId: EventId): Event {
    const pos = this.index.get(eventId.value);
    if (pos === undefined) throw new EventNotFoundError(eventId.value);
    return this.events[pos];
  }

  head(): Option<Event> {
    if (this.events.length === 0) return Option.none();
    return Option.some(this.events[this.events.length - 1]);
  }

  count(): number {
    return this.events.length;
  }

  verifyChain(): ChainVerification {
    for (let i = 1; i < this.events.length; i++) {
      if (this.events[i - 1].hash.value !== this.events[i].prevHash.value) {
        return { valid: false, length: i };
      }
    }
    return { valid: true, length: this.events.length };
  }

  recent(limit: number): Event[] {
    return this.events.slice(-limit).reverse();
  }

  close(): void {}
}
