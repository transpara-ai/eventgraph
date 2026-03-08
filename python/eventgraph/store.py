"""Store protocol and InMemory implementation."""

from __future__ import annotations

import threading
from dataclasses import dataclass
from typing import Protocol

from .errors import ChainIntegrityError, EventNotFoundError
from .event import Event, compute_hash, canonical_form, canonical_content_json
from .types import ActorID, ConversationID, EventID, EventType, Hash, Option


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
    def by_type(self, event_type: EventType, limit: int) -> list[Event]: ...
    def by_source(self, source: ActorID, limit: int) -> list[Event]: ...
    def by_conversation(self, conversation_id: ConversationID, limit: int) -> list[Event]: ...
    def ancestors(self, event_id: EventID, max_depth: int) -> list[Event]: ...
    def descendants(self, event_id: EventID, max_depth: int) -> list[Event]: ...
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

    def by_type(self, event_type: EventType, limit: int) -> list[Event]:
        """Return events of the given type, newest first, up to limit."""
        with self._lock:
            matches: list[Event] = []
            for event in reversed(self._events):
                if event.type.value == event_type.value:
                    matches.append(event)
                    if len(matches) >= limit:
                        break
            return matches

    def by_source(self, source: ActorID, limit: int) -> list[Event]:
        """Return events from the given source, newest first, up to limit."""
        with self._lock:
            matches: list[Event] = []
            for event in reversed(self._events):
                if event.source.value == source.value:
                    matches.append(event)
                    if len(matches) >= limit:
                        break
            return matches

    def by_conversation(self, conversation_id: ConversationID, limit: int) -> list[Event]:
        """Return events in the given conversation, newest first, up to limit."""
        with self._lock:
            matches: list[Event] = []
            for event in reversed(self._events):
                if event.conversation_id.value == conversation_id.value:
                    matches.append(event)
                    if len(matches) >= limit:
                        break
            return matches

    def ancestors(self, event_id: EventID, max_depth: int) -> list[Event]:
        """Return causal ancestors via BFS on causes, up to max_depth levels."""
        with self._lock:
            result: list[Event] = []
            seen: set[str] = {event_id.value}
            # Current frontier: causes of the starting event
            pos = self._index.get(event_id.value)
            if pos is None:
                raise EventNotFoundError(event_id.value)
            frontier: list[str] = [
                c.value for c in self._events[pos].causes
                if c.value != event_id.value  # skip bootstrap self-ref
            ]
            for _ in range(max_depth):
                next_frontier: list[str] = []
                for eid in frontier:
                    if eid in seen:
                        continue
                    seen.add(eid)
                    p = self._index.get(eid)
                    if p is None:
                        continue
                    evt = self._events[p]
                    result.append(evt)
                    for c in evt.causes:
                        if c.value not in seen:
                            next_frontier.append(c.value)
                frontier = next_frontier
                if not frontier:
                    break
            return result

    def descendants(self, event_id: EventID, max_depth: int) -> list[Event]:
        """Return causal descendants via BFS on reverse-cause index, up to max_depth levels."""
        with self._lock:
            if self._index.get(event_id.value) is None:
                raise EventNotFoundError(event_id.value)
            # Build reverse index lazily
            children: dict[str, list[str]] = {}
            for evt in self._events:
                for c in evt.causes:
                    if c.value != evt.id.value:  # skip bootstrap self-ref
                        children.setdefault(c.value, []).append(evt.id.value)

            result: list[Event] = []
            seen: set[str] = {event_id.value}
            frontier = children.get(event_id.value, [])
            for _ in range(max_depth):
                next_frontier: list[str] = []
                for eid in frontier:
                    if eid in seen:
                        continue
                    seen.add(eid)
                    p = self._index.get(eid)
                    if p is None:
                        continue
                    evt = self._events[p]
                    result.append(evt)
                    for child_id in children.get(eid, []):
                        if child_id not in seen:
                            next_frontier.append(child_id)
                frontier = next_frontier
                if not frontier:
                    break
            return result

    def close(self) -> None:
        pass
