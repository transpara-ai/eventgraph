"""Store protocol and InMemory implementation."""

from __future__ import annotations

import threading
from dataclasses import dataclass
from typing import Protocol

from .errors import ChainIntegrityError, EventNotFoundError
from .event import Event, compute_hash, canonical_form, canonical_content_json
from .types import EventID, Hash, Option


@dataclass(frozen=True, slots=True)
class ChainVerification:
    """Result of a chain integrity check."""

    valid: bool
    length: int


class Store(Protocol):
    """Event persistence — append-only, hash-chained."""

    def append(self, event: Event) -> Event: ...
    def get(self, event_id: EventID) -> Event: ...
    def head(self) -> Option[Event]: ...
    def count(self) -> int: ...
    def verify_chain(self) -> ChainVerification: ...
    def close(self) -> None: ...


class InMemoryStore:
    """Thread-safe in-memory event store."""

    def __init__(self) -> None:
        self._lock = threading.Lock()
        self._events: list[Event] = []
        self._index: dict[str, int] = {}  # event_id -> position

    def append(self, event: Event) -> Event:
        with self._lock:
            if self._events:
                last = self._events[-1]
                if event.prev_hash != last.hash:
                    raise ChainIntegrityError(
                        len(self._events),
                        f"prev_hash {event.prev_hash.value} != head hash {last.hash.value}",
                    )
            self._events.append(event)
            self._index[event.id.value] = len(self._events) - 1
            return event

    def get(self, event_id: EventID) -> Event:
        with self._lock:
            pos = self._index.get(event_id.value)
            if pos is None:
                raise EventNotFoundError(event_id.value)
            return self._events[pos]

    def head(self) -> Option[Event]:
        with self._lock:
            if not self._events:
                return Option.none()
            return Option.some(self._events[-1])

    def count(self) -> int:
        with self._lock:
            return len(self._events)

    def verify_chain(self) -> ChainVerification:
        with self._lock:
            for i in range(1, len(self._events)):
                expected_prev = self._events[i - 1].hash
                actual_prev = self._events[i].prev_hash
                if expected_prev != actual_prev:
                    return ChainVerification(valid=False, length=i)
            return ChainVerification(valid=True, length=len(self._events))

    def recent(self, limit: int) -> list[Event]:
        """Return the most recent events, newest first."""
        with self._lock:
            return list(reversed(self._events[-limit:]))

    def close(self) -> None:
        pass
