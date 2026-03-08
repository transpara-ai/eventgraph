import { ChainIntegrityError, EventNotFoundError } from "./errors.js";
import { Event } from "./event.js";
import { ActorId, ConversationId, EventId, EventType, Hash, Option } from "./types.js";

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
  byType(eventType: EventType, limit: number): Event[];
  bySource(source: ActorId, limit: number): Event[];
  byConversation(conversationId: ConversationId, limit: number): Event[];
  ancestors(eventId: EventId, maxDepth: number): Event[];
  descendants(eventId: EventId, maxDepth: number): Event[];
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
    this.causeIndex = null; // invalidate reverse index
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

  byType(eventType: EventType, limit: number): Event[] {
    const result: Event[] = [];
    for (let i = this.events.length - 1; i >= 0 && result.length < limit; i--) {
      if (this.events[i].type.value === eventType.value) {
        result.push(this.events[i]);
      }
    }
    return result;
  }

  bySource(source: ActorId, limit: number): Event[] {
    const result: Event[] = [];
    for (let i = this.events.length - 1; i >= 0 && result.length < limit; i--) {
      if (this.events[i].source.value === source.value) {
        result.push(this.events[i]);
      }
    }
    return result;
  }

  byConversation(conversationId: ConversationId, limit: number): Event[] {
    const result: Event[] = [];
    for (let i = this.events.length - 1; i >= 0 && result.length < limit; i--) {
      if (this.events[i].conversationId.value === conversationId.value) {
        result.push(this.events[i]);
      }
    }
    return result;
  }

  ancestors(eventId: EventId, maxDepth: number): Event[] {
    const pos = this.index.get(eventId.value);
    if (pos === undefined) throw new EventNotFoundError(eventId.value);

    const visited = new Set<string>();
    const result: Event[] = [];
    this.collectAncestors(pos, maxDepth, 0, visited, result);
    return result;
  }

  private collectAncestors(
    idx: number, maxDepth: number, depth: number,
    visited: Set<string>, result: Event[],
  ): void {
    if (depth >= maxDepth) return;
    const ev = this.events[idx];
    for (const causeId of ev.causes) {
      if (visited.has(causeId.value)) continue;
      const causeIdx = this.index.get(causeId.value);
      if (causeIdx === undefined) continue;
      visited.add(causeId.value);
      result.push(this.events[causeIdx]);
      this.collectAncestors(causeIdx, maxDepth, depth + 1, visited, result);
    }
  }

  descendants(eventId: EventId, maxDepth: number): Event[] {
    if (!this.index.has(eventId.value)) throw new EventNotFoundError(eventId.value);

    this.buildCauseIndex();
    const visited = new Set<string>();
    visited.add(eventId.value); // exclude the starting event itself
    const result: Event[] = [];
    this.collectDescendants(eventId.value, maxDepth, 0, visited, result);
    return result;
  }

  private causeIndex: Map<string, number[]> | null = null;

  private buildCauseIndex(): void {
    if (this.causeIndex !== null) return;
    this.causeIndex = new Map();
    for (let i = 0; i < this.events.length; i++) {
      for (const causeId of this.events[i].causes) {
        let children = this.causeIndex.get(causeId.value);
        if (!children) {
          children = [];
          this.causeIndex.set(causeId.value, children);
        }
        children.push(i);
      }
    }
  }

  private collectDescendants(
    id: string, maxDepth: number, depth: number,
    visited: Set<string>, result: Event[],
  ): void {
    if (depth >= maxDepth) return;
    const childIndices = this.causeIndex!.get(id) ?? [];
    for (const childIdx of childIndices) {
      const ev = this.events[childIdx];
      if (visited.has(ev.id.value)) continue;
      visited.add(ev.id.value);
      result.push(ev);
      this.collectDescendants(ev.id.value, maxDepth, depth + 1, visited, result);
    }
  }

  close(): void {}
}
