//! PostgreSQL-backed event store using the `postgres` crate (sync).
//!
//! Optional — only available when the `postgres` feature is enabled.
//! Satisfies the Store trait with persistent storage.

use std::collections::{HashMap, HashSet};
use std::sync::Mutex;

use postgres::{Client, NoTls};

use crate::errors::{EventGraphError, Result};
use crate::event::Event;
use crate::store::{ChainVerification, Store};
use crate::types::*;

const SCHEMA: &str = "
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
";

fn row_to_event(row: &postgres::Row) -> Event {
    let event_id: String = row.get("event_id");
    let event_type: String = row.get("event_type");
    let version: i32 = row.get("version");
    let timestamp_nanos: i64 = row.get("timestamp_nanos");
    let source: String = row.get("source");
    let content_json: String = row.get("content");
    let causes_json: String = row.get("causes");
    let conversation_id: String = row.get("conversation_id");
    let hash: String = row.get("hash");
    let prev_hash: String = row.get("prev_hash");
    let signature: Vec<u8> = row.get("signature");

    let causes: Vec<String> = serde_json::from_str(&causes_json).unwrap_or_default();
    let cause_ids: Vec<EventId> = causes
        .into_iter()
        .map(|c| EventId::new(&c).expect("stored EventId must be valid"))
        .collect();
    let content: std::collections::BTreeMap<String, serde_json::Value> =
        serde_json::from_str(&content_json).unwrap_or_default();

    Event {
        version: version as u32,
        id: EventId::new(&event_id).expect("stored EventId must be valid"),
        event_type: EventType::new(&event_type).expect("stored EventType must be valid"),
        timestamp_nanos: timestamp_nanos as u64,
        source: ActorId::new(&source).expect("stored ActorId must be valid"),
        content,
        causes: NonEmpty::of(cause_ids).expect("causes must be non-empty"),
        conversation_id: ConversationId::new(&conversation_id)
            .expect("stored ConversationId must be valid"),
        hash: Hash::new(&hash).expect("stored Hash must be valid"),
        prev_hash: Hash::new(&prev_hash).expect("stored prev_hash must be valid"),
        signature: Signature::new(signature).expect("stored Signature must be 64 bytes"),
    }
}

/// PostgreSQL-backed event store. Thread-safe via Mutex.
pub struct PostgresStore {
    client: Mutex<Client>,
}

impl PostgresStore {
    /// Connects to a PostgreSQL database and ensures the schema exists.
    ///
    /// `connection_string` is a libpq-style connection string, e.g.
    /// `"host=localhost user=postgres dbname=eventgraph"`.
    pub fn new(connection_string: &str) -> Result<Self> {
        let mut client =
            Client::connect(connection_string, NoTls).map_err(|e| {
                EventGraphError::StoreUnavailable {
                    detail: e.to_string(),
                }
            })?;

        // Create schema — execute each statement individually since
        // batch_execute doesn't support parameterised queries but works
        // fine for DDL.
        client
            .batch_execute(SCHEMA)
            .map_err(|e| EventGraphError::StoreUnavailable {
                detail: e.to_string(),
            })?;

        Ok(Self {
            client: Mutex::new(client),
        })
    }

    // ── Query helpers ───────────────────────────────────────────────────

    fn query_events(client: &mut Client, sql: &str, params: &[&(dyn postgres::types::ToSql + Sync)]) -> Vec<Event> {
        client
            .query(sql, params)
            .unwrap_or_default()
            .iter()
            .map(row_to_event)
            .collect()
    }

    fn query_single(client: &mut Client, sql: &str, params: &[&(dyn postgres::types::ToSql + Sync)]) -> Option<Event> {
        client
            .query_opt(sql, params)
            .ok()
            .flatten()
            .map(|row| row_to_event(&row))
    }

    // ── Owned query methods ─────────────────────────────────────────────

    /// Get an event by ID, returning an owned Event.
    pub fn get_owned(&self, event_id: &EventId) -> Result<Event> {
        let mut client = self.client.lock().unwrap();
        Self::query_single(
            &mut client,
            "SELECT * FROM events WHERE event_id = $1",
            &[&event_id.value()],
        )
        .ok_or_else(|| EventGraphError::EventNotFound {
            event_id: event_id.value().to_string(),
        })
    }

    /// Get the most recent event as an owned value.
    pub fn head_owned(&self) -> Option<Event> {
        let mut client = self.client.lock().unwrap();
        Self::query_single(
            &mut client,
            "SELECT * FROM events ORDER BY position DESC LIMIT 1",
            &[],
        )
    }

    pub fn recent(&self, limit: usize) -> Vec<Event> {
        let mut client = self.client.lock().unwrap();
        Self::query_events(
            &mut client,
            "SELECT * FROM events ORDER BY position DESC LIMIT $1",
            &[&(limit as i64)],
        )
    }

    pub fn by_type(&self, event_type: &EventType, limit: usize) -> Vec<Event> {
        let mut client = self.client.lock().unwrap();
        Self::query_events(
            &mut client,
            "SELECT * FROM events WHERE event_type = $1 ORDER BY position DESC LIMIT $2",
            &[&event_type.value(), &(limit as i64)],
        )
    }

    pub fn by_source(&self, source: &ActorId, limit: usize) -> Vec<Event> {
        let mut client = self.client.lock().unwrap();
        Self::query_events(
            &mut client,
            "SELECT * FROM events WHERE source = $1 ORDER BY position DESC LIMIT $2",
            &[&source.value(), &(limit as i64)],
        )
    }

    pub fn by_conversation(&self, conversation_id: &ConversationId, limit: usize) -> Vec<Event> {
        let mut client = self.client.lock().unwrap();
        Self::query_events(
            &mut client,
            "SELECT * FROM events WHERE conversation_id = $1 ORDER BY position DESC LIMIT $2",
            &[&conversation_id.value(), &(limit as i64)],
        )
    }

    pub fn ancestors(&self, event_id: &EventId, max_depth: usize) -> Result<Vec<Event>> {
        let mut client = self.client.lock().unwrap();
        let start = Self::query_single(
            &mut client,
            "SELECT * FROM events WHERE event_id = $1",
            &[&event_id.value()],
        )
        .ok_or_else(|| EventGraphError::EventNotFound {
            event_id: event_id.value().to_string(),
        })?;

        let mut result = Vec::new();
        let mut visited = HashSet::new();
        visited.insert(event_id.value().to_string());
        let mut frontier: Vec<String> = start
            .causes
            .iter()
            .filter(|c| c.value() != event_id.value())
            .map(|c| c.value().to_string())
            .collect();

        for _ in 0..max_depth {
            if frontier.is_empty() {
                break;
            }
            let mut next = Vec::new();
            for eid in &frontier {
                if !visited.insert(eid.clone()) {
                    continue;
                }
                if let Some(ev) = Self::query_single(
                    &mut client,
                    "SELECT * FROM events WHERE event_id = $1",
                    &[&eid.as_str() as &(dyn postgres::types::ToSql + Sync)],
                ) {
                    for c in ev.causes.iter() {
                        if !visited.contains(c.value()) {
                            next.push(c.value().to_string());
                        }
                    }
                    result.push(ev);
                }
            }
            frontier = next;
        }
        Ok(result)
    }

    pub fn descendants(&self, event_id: &EventId, max_depth: usize) -> Result<Vec<Event>> {
        let mut client = self.client.lock().unwrap();
        Self::query_single(
            &mut client,
            "SELECT * FROM events WHERE event_id = $1",
            &[&event_id.value()],
        )
        .ok_or_else(|| EventGraphError::EventNotFound {
            event_id: event_id.value().to_string(),
        })?;

        // Build reverse index (children map)
        let all_events = Self::query_events(
            &mut client,
            "SELECT * FROM events ORDER BY position ASC",
            &[],
        );
        let mut children: HashMap<String, Vec<String>> = HashMap::new();
        for ev in &all_events {
            for c in ev.causes.iter() {
                if c.value() != ev.id.value() {
                    children
                        .entry(c.value().to_string())
                        .or_default()
                        .push(ev.id.value().to_string());
                }
            }
        }

        let mut result = Vec::new();
        let mut visited = HashSet::new();
        visited.insert(event_id.value().to_string());
        let mut frontier = children
            .get(event_id.value())
            .cloned()
            .unwrap_or_default();

        for _ in 0..max_depth {
            if frontier.is_empty() {
                break;
            }
            let mut next = Vec::new();
            for eid in &frontier {
                if !visited.insert(eid.clone()) {
                    continue;
                }
                if let Some(ev) = Self::query_single(
                    &mut client,
                    "SELECT * FROM events WHERE event_id = $1",
                    &[&eid.as_str() as &(dyn postgres::types::ToSql + Sync)],
                ) {
                    result.push(ev);
                    if let Some(ch) = children.get(eid) {
                        for child in ch {
                            if !visited.contains(child.as_str()) {
                                next.push(child.clone());
                            }
                        }
                    }
                }
            }
            frontier = next;
        }
        Ok(result)
    }
}

impl Store for PostgresStore {
    fn append(&mut self, event: Event) -> Result<Event> {
        let mut client = self.client.lock().unwrap();

        // Idempotency — if event already stored with same hash, return it
        if let Some(existing) = Self::query_single(
            &mut client,
            "SELECT * FROM events WHERE event_id = $1",
            &[&event.id.value()],
        ) {
            if existing.hash.value() != event.hash.value() {
                return Err(EventGraphError::ChainIntegrity {
                    position: 0,
                    detail: format!("hash mismatch for existing event {}", event.id.value()),
                });
            }
            return Ok(existing);
        }

        // Chain continuity — prev_hash must match head
        if let Some(head) = Self::query_single(
            &mut client,
            "SELECT * FROM events ORDER BY position DESC LIMIT 1",
            &[],
        ) {
            if event.prev_hash.value() != head.hash.value() {
                return Err(EventGraphError::ChainIntegrity {
                    position: 0,
                    detail: format!(
                        "prev_hash {} != head hash {}",
                        event.prev_hash.value(),
                        head.hash.value()
                    ),
                });
            }
        }

        let causes_json =
            serde_json::to_string(&event.causes.iter().map(|c| c.value()).collect::<Vec<_>>())
                .unwrap();
        let content_json = serde_json::to_string(&event.content()).unwrap();

        client.execute(
            "INSERT INTO events (event_id, event_type, version, timestamp_nanos, source, content, causes, conversation_id, hash, prev_hash, signature) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)",
            &[
                &event.id.value(),
                &event.event_type.value(),
                &(event.version as i32),
                &(event.timestamp_nanos as i64),
                &event.source.value(),
                &content_json,
                &causes_json,
                &event.conversation_id.value(),
                &event.hash.value(),
                &event.prev_hash.value(),
                &event.signature.bytes(),
            ],
        ).map_err(|e| EventGraphError::StoreUnavailable { detail: e.to_string() })?;

        Ok(event)
    }

    fn get(&self, _event_id: &EventId) -> Result<&Event> {
        // PostgresStore cannot return references — data lives in the database, not in memory.
        // Use get_owned() instead.
        unimplemented!("PostgresStore cannot return references; use get_owned() instead")
    }

    fn head(&self) -> Option<&Event> {
        // PostgresStore cannot return references — data lives in the database, not in memory.
        // Use head_owned() instead.
        unimplemented!("PostgresStore cannot return references; use head_owned() instead")
    }

    fn count(&self) -> usize {
        let mut client = self.client.lock().unwrap();
        let row = client
            .query_one("SELECT COUNT(*) FROM events", &[])
            .ok();
        row.map(|r| {
            let c: i64 = r.get(0);
            c as usize
        })
        .unwrap_or(0)
    }

    fn verify_chain(&self) -> ChainVerification {
        let mut client = self.client.lock().unwrap();
        let events = Self::query_events(
            &mut client,
            "SELECT * FROM events ORDER BY position ASC",
            &[],
        );
        for i in 1..events.len() {
            if events[i - 1].hash.value() != events[i].prev_hash.value() {
                return ChainVerification {
                    valid: false,
                    length: i,
                };
            }
        }
        ChainVerification {
            valid: true,
            length: events.len(),
        }
    }

    fn close(&mut self) {
        // Connection is closed when PostgresStore is dropped.
    }
}
