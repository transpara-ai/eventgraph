"""SQLite-backed event store implementation."""

from __future__ import annotations

import json
import sqlite3
import threading
from pathlib import Path
from typing import Any

from .errors import ChainIntegrityError, EventNotFoundError
from .event import Event, canonical_content_json, canonical_form, compute_hash
from .store import ChainVerification
from .types import (
    ActorID,
    ConversationID,
    EventID,
    EventType,
    Hash,
    NonEmpty,
    Option,
    Signature,
)

_SCHEMA = """
CREATE TABLE IF NOT EXISTS events (
    position   INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id   TEXT NOT NULL UNIQUE,
    event_type TEXT NOT NULL,
    version    INTEGER NOT NULL,
    timestamp_nanos INTEGER NOT NULL,
    source     TEXT NOT NULL,
    content    TEXT NOT NULL,
    causes     TEXT NOT NULL,
    conversation_id TEXT NOT NULL,
    hash       TEXT NOT NULL,
    prev_hash  TEXT NOT NULL,
    signature  BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_events_type ON events(event_type);
CREATE INDEX IF NOT EXISTS idx_events_source ON events(source);
CREATE INDEX IF NOT EXISTS idx_events_conversation ON events(conversation_id);
"""


def _event_to_row(ev: Event) -> tuple:
    causes_json = json.dumps([c.value for c in ev.causes])
    content_json = json.dumps(ev.content, sort_keys=True)
    return (
        ev.id.value,
        ev.type.value,
        ev.version,
        ev.timestamp_nanos,
        ev.source.value,
        content_json,
        causes_json,
        ev.conversation_id.value,
        ev.hash.value,
        ev.prev_hash.value,
        ev.signature.bytes_,
    )


def _row_to_event(row: tuple) -> Event:
    (position, event_id, event_type, version, timestamp_nanos, source,
     content_json, causes_json, conversation_id, hash_val, prev_hash, sig_bytes) = row
    cause_ids = [EventID(c) for c in json.loads(causes_json)]
    content = json.loads(content_json)
    return Event(
        _version=version,
        _id=EventID(event_id),
        _type=EventType(event_type),
        _timestamp_nanos=timestamp_nanos,
        _source=ActorID(source),
        _content=content,
        _causes=NonEmpty.of(cause_ids),
        _conversation_id=ConversationID(conversation_id),
        _hash=Hash(hash_val),
        _prev_hash=Hash(prev_hash),
        _signature=Signature(sig_bytes if isinstance(sig_bytes, bytes) else b"\x00" * 64),
    )


class SQLiteStore:
    """Thread-safe SQLite-backed event store.

    Satisfies the Store protocol. Data persists across process restarts.
    """

    def __init__(self, path: str | Path = ":memory:") -> None:
        self._path = str(path)
        self._lock = threading.Lock()
        self._conn = sqlite3.connect(self._path, check_same_thread=False)
        self._conn.execute("PRAGMA journal_mode=WAL")
        self._conn.execute("PRAGMA foreign_keys=ON")
        self._conn.executescript(_SCHEMA)
        self._conn.commit()

    def append(self, event: Event) -> Event:
        with self._lock:
            # Idempotency check
            row = self._conn.execute(
                "SELECT * FROM events WHERE event_id = ?", (event.id.value,)
            ).fetchone()
            if row is not None:
                stored = _row_to_event(row)
                if stored.hash.value != event.hash.value:
                    raise ChainIntegrityError(
                        0, f"hash mismatch for existing event {event.id.value}"
                    )
                return stored

            # Verify chain continuity
            head_row = self._conn.execute(
                "SELECT * FROM events ORDER BY position DESC LIMIT 1"
            ).fetchone()
            if head_row is not None:
                head_hash = head_row[9]  # hash column
                if event.prev_hash.value != head_hash:
                    raise ChainIntegrityError(
                        head_row[0] + 1,
                        f"prev_hash {event.prev_hash.value} != head hash {head_hash}",
                    )
            else:
                if event.prev_hash.value != Hash.zero().value:
                    raise ChainIntegrityError(
                        0,
                        f"first event prev_hash must be zero hash, got {event.prev_hash.value}",
                    )

            self._conn.execute(
                """INSERT INTO events
                   (event_id, event_type, version, timestamp_nanos, source,
                    content, causes, conversation_id, hash, prev_hash, signature)
                   VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)""",
                _event_to_row(event),
            )
            self._conn.commit()
            return event

    def get(self, event_id: EventID) -> Event:
        with self._lock:
            row = self._conn.execute(
                "SELECT * FROM events WHERE event_id = ?", (event_id.value,)
            ).fetchone()
            if row is None:
                raise EventNotFoundError(event_id.value)
            return _row_to_event(row)

    def head(self) -> Option[Event]:
        with self._lock:
            row = self._conn.execute(
                "SELECT * FROM events ORDER BY position DESC LIMIT 1"
            ).fetchone()
            if row is None:
                return Option.none()
            return Option.some(_row_to_event(row))

    def count(self) -> int:
        with self._lock:
            row = self._conn.execute("SELECT COUNT(*) FROM events").fetchone()
            return row[0]

    def verify_chain(self) -> ChainVerification:
        with self._lock:
            rows = self._conn.execute(
                "SELECT * FROM events ORDER BY position ASC"
            ).fetchall()
            for i, row in enumerate(rows):
                ev = _row_to_event(row)
                if i == 0:
                    if ev.prev_hash.value != Hash.zero().value:
                        return ChainVerification(valid=False, length=i)
                else:
                    prev_ev = _row_to_event(rows[i - 1])
                    if ev.prev_hash.value != prev_ev.hash.value:
                        return ChainVerification(valid=False, length=i)
            return ChainVerification(valid=True, length=len(rows))

    def recent(self, limit: int) -> list[Event]:
        with self._lock:
            rows = self._conn.execute(
                "SELECT * FROM events ORDER BY position DESC LIMIT ?", (limit,)
            ).fetchall()
            return [_row_to_event(r) for r in rows]

    def by_type(self, event_type: EventType, limit: int) -> list[Event]:
        with self._lock:
            rows = self._conn.execute(
                "SELECT * FROM events WHERE event_type = ? ORDER BY position DESC LIMIT ?",
                (event_type.value, limit),
            ).fetchall()
            return [_row_to_event(r) for r in rows]

    def by_source(self, source: ActorID, limit: int) -> list[Event]:
        with self._lock:
            rows = self._conn.execute(
                "SELECT * FROM events WHERE source = ? ORDER BY position DESC LIMIT ?",
                (source.value, limit),
            ).fetchall()
            return [_row_to_event(r) for r in rows]

    def by_conversation(self, conversation_id: ConversationID, limit: int) -> list[Event]:
        with self._lock:
            rows = self._conn.execute(
                "SELECT * FROM events WHERE conversation_id = ? ORDER BY position DESC LIMIT ?",
                (conversation_id.value, limit),
            ).fetchall()
            return [_row_to_event(r) for r in rows]

    def ancestors(self, event_id: EventID, max_depth: int) -> list[Event]:
        with self._lock:
            row = self._conn.execute(
                "SELECT * FROM events WHERE event_id = ?", (event_id.value,)
            ).fetchone()
            if row is None:
                raise EventNotFoundError(event_id.value)

            result: list[Event] = []
            seen: set[str] = {event_id.value}
            ev = _row_to_event(row)
            frontier = [c.value for c in ev.causes if c.value != event_id.value]

            for _ in range(max_depth):
                next_frontier: list[str] = []
                for eid in frontier:
                    if eid in seen:
                        continue
                    seen.add(eid)
                    r = self._conn.execute(
                        "SELECT * FROM events WHERE event_id = ?", (eid,)
                    ).fetchone()
                    if r is None:
                        continue
                    ancestor = _row_to_event(r)
                    result.append(ancestor)
                    for c in ancestor.causes:
                        if c.value not in seen:
                            next_frontier.append(c.value)
                frontier = next_frontier
                if not frontier:
                    break
            return result

    def descendants(self, event_id: EventID, max_depth: int) -> list[Event]:
        with self._lock:
            row = self._conn.execute(
                "SELECT * FROM events WHERE event_id = ?", (event_id.value,)
            ).fetchone()
            if row is None:
                raise EventNotFoundError(event_id.value)

            # Build reverse cause index
            all_rows = self._conn.execute(
                "SELECT * FROM events ORDER BY position ASC"
            ).fetchall()
            children: dict[str, list[str]] = {}
            for r in all_rows:
                ev = _row_to_event(r)
                for c in ev.causes:
                    if c.value != ev.id.value:
                        children.setdefault(c.value, []).append(ev.id.value)

            result: list[Event] = []
            seen: set[str] = {event_id.value}
            frontier = children.get(event_id.value, [])

            for _ in range(max_depth):
                next_frontier: list[str] = []
                for eid in frontier:
                    if eid in seen:
                        continue
                    seen.add(eid)
                    r = self._conn.execute(
                        "SELECT * FROM events WHERE event_id = ?", (eid,)
                    ).fetchone()
                    if r is None:
                        continue
                    desc = _row_to_event(r)
                    result.append(desc)
                    for child_id in children.get(eid, []):
                        if child_id not in seen:
                            next_frontier.append(child_id)
                frontier = next_frontier
                if not frontier:
                    break
            return result

    def close(self) -> None:
        with self._lock:
            self._conn.close()
