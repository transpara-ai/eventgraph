// Optional dependency — requires: npm install pg
//
// PostgresStore is an async PostgreSQL-backed event store.
// Because the Store interface defines synchronous methods, PostgresStore does NOT
// implement Store directly. All methods return Promises instead. Consumer code must
// await each call.

import { ChainIntegrityError, EventNotFoundError } from "./errors.js";
import { Event } from "./event.js";
import type { ChainVerification } from "./store.js";
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

// Minimal type declarations for pg to avoid requiring @types/pg
interface PgPool {
  query(text: string, values?: unknown[]): Promise<{ rows: unknown[] }>;
  end(): Promise<void>;
}

const SCHEMA = `
CREATE TABLE IF NOT EXISTS events (
  position        BIGSERIAL PRIMARY KEY,
  event_id        TEXT NOT NULL UNIQUE,
  event_type      TEXT NOT NULL,
  version         INTEGER NOT NULL,
  timestamp_nanos BIGINT NOT NULL,
  source          TEXT NOT NULL,
  content         TEXT NOT NULL,
  causes          TEXT NOT NULL,
  conversation_id TEXT NOT NULL,
  hash            TEXT NOT NULL,
  prev_hash       TEXT NOT NULL,
  signature       BYTEA NOT NULL
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
  timestamp_nanos: string; // BIGINT comes back as string from pg
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
    Number(row.timestamp_nanos),
    new ActorId(row.source),
    content,
    NonEmpty.of(causeIds),
    new ConversationId(row.conversation_id),
    new Hash(row.hash),
    new Hash(row.prev_hash),
    new Signature(new Uint8Array(row.signature)),
  );
}

export class PostgresStore {
  private readonly pool: PgPool;
  private readonly initPromise: Promise<void>;

  constructor(connectionString: string) {
    // Dynamic require to avoid hard dependency
    // eslint-disable-next-line @typescript-eslint/no-require-imports
    const { Pool } = require("pg");
    this.pool = new Pool({ connectionString }) as PgPool;
    this.initPromise = this.init();
  }

  private async init(): Promise<void> {
    await this.pool.query(SCHEMA);
  }

  async ensureReady(): Promise<void> {
    await this.initPromise;
  }

  async append(event: Event): Promise<Event> {
    await this.ensureReady();

    // Idempotency
    const existingResult = await this.pool.query(
      "SELECT * FROM events WHERE event_id = $1",
      [event.id.value],
    );
    if (existingResult.rows.length > 0) {
      const existing = existingResult.rows[0] as EventRow;
      if (existing.hash !== event.hash.value) {
        throw new ChainIntegrityError(0, `hash mismatch for existing event ${event.id.value}`);
      }
      return rowToEvent(existing);
    }

    // Verify chain
    const headResult = await this.pool.query(
      "SELECT * FROM events ORDER BY position DESC LIMIT 1",
    );
    if (headResult.rows.length > 0) {
      const headRow = headResult.rows[0] as EventRow;
      if (event.prevHash.value !== headRow.hash) {
        throw new ChainIntegrityError(
          Number(headRow.position) + 1,
          `prev_hash ${event.prevHash.value} != head hash ${headRow.hash}`,
        );
      }
    }

    const causesJson = JSON.stringify(event.causes.toArray().map((c) => c.value));
    const contentJson = JSON.stringify(event.content, Object.keys(event.content).sort());

    await this.pool.query(
      `INSERT INTO events
       (event_id, event_type, version, timestamp_nanos, source,
        content, causes, conversation_id, hash, prev_hash, signature)
       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
      [
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
      ],
    );

    return event;
  }

  async get(eventId: EventId): Promise<Event> {
    await this.ensureReady();
    const result = await this.pool.query(
      "SELECT * FROM events WHERE event_id = $1",
      [eventId.value],
    );
    if (result.rows.length === 0) throw new EventNotFoundError(eventId.value);
    return rowToEvent(result.rows[0] as EventRow);
  }

  async head(): Promise<Option<Event>> {
    await this.ensureReady();
    const result = await this.pool.query(
      "SELECT * FROM events ORDER BY position DESC LIMIT 1",
    );
    if (result.rows.length === 0) return Option.none();
    return Option.some(rowToEvent(result.rows[0] as EventRow));
  }

  async count(): Promise<number> {
    await this.ensureReady();
    const result = await this.pool.query("SELECT COUNT(*) as cnt FROM events");
    return Number((result.rows[0] as { cnt: string }).cnt);
  }

  async verifyChain(): Promise<ChainVerification> {
    await this.ensureReady();
    const result = await this.pool.query(
      "SELECT * FROM events ORDER BY position ASC",
    );
    const rows = result.rows as EventRow[];
    for (let i = 0; i < rows.length; i++) {
      if (i > 0 && rows[i].prev_hash !== rows[i - 1].hash) {
        return { valid: false, length: i };
      }
    }
    return { valid: true, length: rows.length };
  }

  async recent(limit: number): Promise<Event[]> {
    await this.ensureReady();
    const result = await this.pool.query(
      "SELECT * FROM events ORDER BY position DESC LIMIT $1",
      [limit],
    );
    return (result.rows as EventRow[]).map(rowToEvent);
  }

  async byType(eventType: EventType, limit: number): Promise<Event[]> {
    await this.ensureReady();
    const result = await this.pool.query(
      "SELECT * FROM events WHERE event_type = $1 ORDER BY position DESC LIMIT $2",
      [eventType.value, limit],
    );
    return (result.rows as EventRow[]).map(rowToEvent);
  }

  async bySource(source: ActorId, limit: number): Promise<Event[]> {
    await this.ensureReady();
    const result = await this.pool.query(
      "SELECT * FROM events WHERE source = $1 ORDER BY position DESC LIMIT $2",
      [source.value, limit],
    );
    return (result.rows as EventRow[]).map(rowToEvent);
  }

  async byConversation(conversationId: ConversationId, limit: number): Promise<Event[]> {
    await this.ensureReady();
    const result = await this.pool.query(
      "SELECT * FROM events WHERE conversation_id = $1 ORDER BY position DESC LIMIT $2",
      [conversationId.value, limit],
    );
    return (result.rows as EventRow[]).map(rowToEvent);
  }

  async ancestors(eventId: EventId, maxDepth: number): Promise<Event[]> {
    await this.ensureReady();

    const startResult = await this.pool.query(
      "SELECT * FROM events WHERE event_id = $1",
      [eventId.value],
    );
    if (startResult.rows.length === 0) throw new EventNotFoundError(eventId.value);

    const result: Event[] = [];
    const visited = new Set<string>([eventId.value]);
    const startEv = rowToEvent(startResult.rows[0] as EventRow);
    let frontier = startEv.causes.toArray()
      .filter((c) => c.value !== eventId.value)
      .map((c) => c.value);

    for (let d = 0; d < maxDepth && frontier.length > 0; d++) {
      const nextFrontier: string[] = [];
      for (const eid of frontier) {
        if (visited.has(eid)) continue;
        visited.add(eid);
        const rowResult = await this.pool.query(
          "SELECT * FROM events WHERE event_id = $1",
          [eid],
        );
        if (rowResult.rows.length === 0) continue;
        const ev = rowToEvent(rowResult.rows[0] as EventRow);
        result.push(ev);
        for (const c of ev.causes) {
          if (!visited.has(c.value)) nextFrontier.push(c.value);
        }
      }
      frontier = nextFrontier;
    }
    return result;
  }

  async descendants(eventId: EventId, maxDepth: number): Promise<Event[]> {
    await this.ensureReady();

    const existsResult = await this.pool.query(
      "SELECT 1 FROM events WHERE event_id = $1",
      [eventId.value],
    );
    if (existsResult.rows.length === 0) throw new EventNotFoundError(eventId.value);

    // Build reverse index
    const allResult = await this.pool.query(
      "SELECT * FROM events ORDER BY position ASC",
    );
    const allRows = allResult.rows as EventRow[];
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
        const rowResult = await this.pool.query(
          "SELECT * FROM events WHERE event_id = $1",
          [eid],
        );
        if (rowResult.rows.length === 0) continue;
        result.push(rowToEvent(rowResult.rows[0] as EventRow));
        for (const child of children.get(eid) ?? []) {
          if (!visited.has(child)) nextFrontier.push(child);
        }
      }
      frontier = nextFrontier;
    }
    return result;
  }

  async close(): Promise<void> {
    await this.ensureReady();
    await this.pool.end();
  }
}
