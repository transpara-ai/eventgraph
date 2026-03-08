/**
 * SQLite-backed event store using better-sqlite3.
 *
 * Optional dependency — only imported when SQLiteStore is used.
 * Install: npm install better-sqlite3 @types/better-sqlite3
 */

import { ChainIntegrityError, EventNotFoundError } from "./errors.js";
import { Event } from "./event.js";
import type { ChainVerification, Store } from "./store.js";
import {
  ActorId,
  ConversationId,
  EventId,
  EventType,
  Hash,
  NonEmpty,
  Option,
  Signature,
} from "./types.js";

// Minimal type declarations for better-sqlite3 to avoid requiring @types/better-sqlite3
interface BetterSqlite3Database {
  pragma(pragma: string): unknown;
  exec(sql: string): void;
  prepare(sql: string): BetterSqlite3Statement;
  close(): void;
}
interface BetterSqlite3Statement {
  run(...params: unknown[]): unknown;
  get(...params: unknown[]): unknown;
  all(...params: unknown[]): unknown[];
}

const SCHEMA = `
CREATE TABLE IF NOT EXISTS events (
  position        INTEGER PRIMARY KEY AUTOINCREMENT,
  event_id        TEXT NOT NULL UNIQUE,
  event_type      TEXT NOT NULL,
  version         INTEGER NOT NULL,
  timestamp_nanos INTEGER NOT NULL,
  source          TEXT NOT NULL,
  content         TEXT NOT NULL,
  causes          TEXT NOT NULL,
  conversation_id TEXT NOT NULL,
  hash            TEXT NOT NULL,
  prev_hash       TEXT NOT NULL,
  signature       BLOB NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_events_type ON events(event_type);
CREATE INDEX IF NOT EXISTS idx_events_source ON events(source);
CREATE INDEX IF NOT EXISTS idx_events_conversation ON events(conversation_id);
`;

interface EventRow {
  position: number;
  event_id: string;
  event_type: string;
  version: number;
  timestamp_nanos: number;
  source: string;
  content: string;
  causes: string;
  conversation_id: string;
  hash: string;
  prev_hash: string;
  signature: Buffer;
}

function rowToEvent(row: EventRow): Event {
  const causeIds = (JSON.parse(row.causes) as string[]).map((c) => new EventId(c));
  const content = JSON.parse(row.content) as Record<string, unknown>;
  return new Event(
    row.version,
    new EventId(row.event_id),
    new EventType(row.event_type),
    row.timestamp_nanos,
    new ActorId(row.source),
    content,
    NonEmpty.of(causeIds),
    new ConversationId(row.conversation_id),
    new Hash(row.hash),
    new Hash(row.prev_hash),
    new Signature(new Uint8Array(row.signature)),
  );
}

export class SQLiteStore implements Store {
  private readonly db: BetterSqlite3Database;

  constructor(path: string = ":memory:") {
    // Dynamic import to avoid hard dependency
    // eslint-disable-next-line @typescript-eslint/no-require-imports
    const BetterSqlite3 = require("better-sqlite3");
    this.db = new BetterSqlite3(path) as BetterSqlite3Database;
    this.db.pragma("journal_mode = WAL");
    this.db.pragma("foreign_keys = ON");
    this.db.exec(SCHEMA);
  }

  append(event: Event): Event {
    // Idempotency
    const existing = this.db
      .prepare("SELECT * FROM events WHERE event_id = ?")
      .get(event.id.value) as EventRow | undefined;
    if (existing) {
      if (existing.hash !== event.hash.value) {
        throw new ChainIntegrityError(0, `hash mismatch for existing event ${event.id.value}`);
      }
      return rowToEvent(existing);
    }

    // Verify chain
    const headRow = this.db
      .prepare("SELECT * FROM events ORDER BY position DESC LIMIT 1")
      .get() as EventRow | undefined;
    if (headRow) {
      if (event.prevHash.value !== headRow.hash) {
        throw new ChainIntegrityError(
          headRow.position + 1,
          `prev_hash ${event.prevHash.value} != head hash ${headRow.hash}`,
        );
      }
    }

    const causesJson = JSON.stringify(event.causes.toArray().map((c) => c.value));
    const contentJson = JSON.stringify(event.content, Object.keys(event.content).sort());

    this.db
      .prepare(
        `INSERT INTO events
         (event_id, event_type, version, timestamp_nanos, source,
          content, causes, conversation_id, hash, prev_hash, signature)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
      )
      .run(
        event.id.value,
        event.type.value,
        event.version,
        event.timestampNanos,
        event.source.value,
        contentJson,
        causesJson,
        event.conversationId.value,
        event.hash.value,
        event.prevHash.value,
        Buffer.from(event.signature.bytes),
      );

    return event;
  }

  get(eventId: EventId): Event {
    const row = this.db
      .prepare("SELECT * FROM events WHERE event_id = ?")
      .get(eventId.value) as EventRow | undefined;
    if (!row) throw new EventNotFoundError(eventId.value);
    return rowToEvent(row);
  }

  head(): Option<Event> {
    const row = this.db
      .prepare("SELECT * FROM events ORDER BY position DESC LIMIT 1")
      .get() as EventRow | undefined;
    if (!row) return Option.none();
    return Option.some(rowToEvent(row));
  }

  count(): number {
    const row = this.db.prepare("SELECT COUNT(*) as cnt FROM events").get() as { cnt: number };
    return row.cnt;
  }

  verifyChain(): ChainVerification {
    const rows = this.db
      .prepare("SELECT * FROM events ORDER BY position ASC")
      .all() as EventRow[];
    for (let i = 0; i < rows.length; i++) {
      if (i > 0 && rows[i].prev_hash !== rows[i - 1].hash) {
        return { valid: false, length: i };
      }
    }
    return { valid: true, length: rows.length };
  }

  recent(limit: number): Event[] {
    const rows = this.db
      .prepare("SELECT * FROM events ORDER BY position DESC LIMIT ?")
      .all(limit) as EventRow[];
    return rows.map(rowToEvent);
  }

  byType(eventType: EventType, limit: number): Event[] {
    const rows = this.db
      .prepare("SELECT * FROM events WHERE event_type = ? ORDER BY position DESC LIMIT ?")
      .all(eventType.value, limit) as EventRow[];
    return rows.map(rowToEvent);
  }

  bySource(source: ActorId, limit: number): Event[] {
    const rows = this.db
      .prepare("SELECT * FROM events WHERE source = ? ORDER BY position DESC LIMIT ?")
      .all(source.value, limit) as EventRow[];
    return rows.map(rowToEvent);
  }

  byConversation(conversationId: ConversationId, limit: number): Event[] {
    const rows = this.db
      .prepare("SELECT * FROM events WHERE conversation_id = ? ORDER BY position DESC LIMIT ?")
      .all(conversationId.value, limit) as EventRow[];
    return rows.map(rowToEvent);
  }

  ancestors(eventId: EventId, maxDepth: number): Event[] {
    const startRow = this.db
      .prepare("SELECT * FROM events WHERE event_id = ?")
      .get(eventId.value) as EventRow | undefined;
    if (!startRow) throw new EventNotFoundError(eventId.value);

    const result: Event[] = [];
    const visited = new Set<string>([eventId.value]);
    const startEv = rowToEvent(startRow);
    let frontier = startEv.causes.toArray()
      .filter((c) => c.value !== eventId.value)
      .map((c) => c.value);

    for (let d = 0; d < maxDepth && frontier.length > 0; d++) {
      const nextFrontier: string[] = [];
      for (const eid of frontier) {
        if (visited.has(eid)) continue;
        visited.add(eid);
        const row = this.db
          .prepare("SELECT * FROM events WHERE event_id = ?")
          .get(eid) as EventRow | undefined;
        if (!row) continue;
        const ev = rowToEvent(row);
        result.push(ev);
        for (const c of ev.causes) {
          if (!visited.has(c.value)) nextFrontier.push(c.value);
        }
      }
      frontier = nextFrontier;
    }
    return result;
  }

  descendants(eventId: EventId, maxDepth: number): Event[] {
    const exists = this.db
      .prepare("SELECT 1 FROM events WHERE event_id = ?")
      .get(eventId.value);
    if (!exists) throw new EventNotFoundError(eventId.value);

    // Build reverse index
    const allRows = this.db
      .prepare("SELECT * FROM events ORDER BY position ASC")
      .all() as EventRow[];
    const children = new Map<string, string[]>();
    for (const row of allRows) {
      const causes = JSON.parse(row.causes) as string[];
      for (const c of causes) {
        if (c !== row.event_id) {
          const list = children.get(c) ?? [];
          list.push(row.event_id);
          children.set(c, list);
        }
      }
    }

    const result: Event[] = [];
    const visited = new Set<string>([eventId.value]);
    let frontier = children.get(eventId.value) ?? [];

    for (let d = 0; d < maxDepth && frontier.length > 0; d++) {
      const nextFrontier: string[] = [];
      for (const eid of frontier) {
        if (visited.has(eid)) continue;
        visited.add(eid);
        const row = this.db
          .prepare("SELECT * FROM events WHERE event_id = ?")
          .get(eid) as EventRow | undefined;
        if (!row) continue;
        result.push(rowToEvent(row));
        for (const child of children.get(eid) ?? []) {
          if (!visited.has(child)) nextFrontier.push(child);
        }
      }
      frontier = nextFrontier;
    }
    return result;
  }

  close(): void {
    this.db.close();
  }
}
