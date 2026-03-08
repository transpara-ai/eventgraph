"""Actor module — identity, lifecycle, and in-memory actor store."""

from __future__ import annotations

import copy
import hashlib
import threading
import time
from dataclasses import dataclass, field
from enum import Enum
from typing import Any, Protocol, runtime_checkable

from .errors import (
    ActorKeyNotFoundError,
    ActorNotFoundError,
    InvalidTransitionError,
)
from .types import ActorID, Cursor, EventID, Option, Page, PublicKey


# ── ActorType ────────────────────────────────────────────────────────────


class ActorType(str, Enum):
    """What kind of decision-maker an actor is."""

    HUMAN = "Human"
    AI = "AI"
    SYSTEM = "System"
    COMMITTEE = "Committee"
    RULES_ENGINE = "RulesEngine"


# ── ActorStatus ──────────────────────────────────────────────────────────

_VALID_ACTOR_TRANSITIONS: dict[str, list[str]] = {
    "active": ["suspended", "memorial"],
    "suspended": ["active", "memorial"],
    "memorial": [],  # terminal
}


class ActorStatus(str, Enum):
    """Actor lifecycle status. Memorial is terminal."""

    ACTIVE = "active"
    SUSPENDED = "suspended"
    MEMORIAL = "memorial"

    def transition_to(self, target: ActorStatus) -> ActorStatus:
        """Validate and perform a status transition.

        Raises InvalidTransitionError if the transition is not allowed.
        """
        valid = _VALID_ACTOR_TRANSITIONS[self.value]
        if target.value in valid:
            return target
        raise InvalidTransitionError("ActorStatus", self.value, target.value)

    def valid_transitions(self) -> list[ActorStatus]:
        """Return the list of valid target statuses from this status."""
        return [ActorStatus(v) for v in _VALID_ACTOR_TRANSITIONS[self.value]]


# ── Actor ────────────────────────────────────────────────────────────────


def _deep_copy_metadata(src: dict[str, Any] | None) -> dict[str, Any]:
    if src is None:
        return {}
    return copy.deepcopy(src)


@dataclass(frozen=True, slots=True)
class Actor:
    """Immutable actor identity."""

    _id: ActorID
    _public_key: PublicKey
    _display_name: str
    _actor_type: ActorType
    _metadata: dict[str, Any]
    _created_at: int
    _status: ActorStatus

    def __init__(
        self,
        id: ActorID,
        public_key: PublicKey,
        display_name: str,
        actor_type: ActorType,
        metadata: dict[str, Any] | None,
        created_at: int,
        status: ActorStatus,
    ) -> None:
        object.__setattr__(self, "_id", id)
        object.__setattr__(self, "_public_key", public_key)
        object.__setattr__(self, "_display_name", display_name)
        object.__setattr__(self, "_actor_type", actor_type)
        object.__setattr__(self, "_metadata", _deep_copy_metadata(metadata))
        object.__setattr__(self, "_created_at", created_at)
        object.__setattr__(self, "_status", status)

    @property
    def id(self) -> ActorID:
        return self._id

    @property
    def public_key(self) -> PublicKey:
        return self._public_key

    @property
    def display_name(self) -> str:
        return self._display_name

    @property
    def actor_type(self) -> ActorType:
        return self._actor_type

    @property
    def metadata(self) -> dict[str, Any]:
        return _deep_copy_metadata(self._metadata)

    @property
    def created_at(self) -> int:
        return self._created_at

    @property
    def status(self) -> ActorStatus:
        return self._status

    def _with_status(self, status: ActorStatus) -> Actor:
        return Actor(
            id=self._id,
            public_key=self._public_key,
            display_name=self._display_name,
            actor_type=self._actor_type,
            metadata=self._metadata,
            created_at=self._created_at,
            status=status,
        )

    def _with_updates(self, updates: ActorUpdate) -> Actor:
        new_name = updates.display_name.unwrap() if updates.display_name.is_some() else self._display_name
        md = _deep_copy_metadata(self._metadata)
        if updates.metadata.is_some():
            for k, v in updates.metadata.unwrap().items():
                md[k] = copy.deepcopy(v)
        return Actor(
            id=self._id,
            public_key=self._public_key,
            display_name=new_name,
            actor_type=self._actor_type,
            metadata=md,
            created_at=self._created_at,
            status=self._status,
        )


# ── ActorUpdate ──────────────────────────────────────────────────────────


@dataclass(frozen=True, slots=True)
class ActorUpdate:
    """Describes updates to apply to an actor."""

    display_name: Option[str] = field(default_factory=lambda: Option.none())
    metadata: Option[dict[str, Any]] = field(default_factory=lambda: Option.none())


# ── ActorFilter ──────────────────────────────────────────────────────────


@dataclass(frozen=True, slots=True)
class ActorFilter:
    """Criteria for listing actors."""

    status: Option[ActorStatus] = field(default_factory=lambda: Option.none())
    actor_type: Option[ActorType] = field(default_factory=lambda: Option.none())
    limit: int = 100
    after: Option[str] = field(default_factory=lambda: Option.none())


# ── ActorStore Protocol ──────────────────────────────────────────────────


@runtime_checkable
class ActorStore(Protocol):
    """Actor persistence interface."""

    def register(
        self, public_key: PublicKey, display_name: str, actor_type: ActorType
    ) -> Actor: ...

    def get(self, actor_id: ActorID) -> Actor: ...

    def get_by_public_key(self, public_key: PublicKey) -> Actor: ...

    def update(self, actor_id: ActorID, updates: ActorUpdate) -> Actor: ...

    def list(self, filter: ActorFilter) -> Page[Actor]: ...

    def suspend(self, actor_id: ActorID, reason_event_id: EventID) -> Actor: ...

    def reactivate(self, actor_id: ActorID, reason_event_id: EventID) -> Actor: ...

    def memorial(self, actor_id: ActorID, reason_event_id: EventID) -> Actor: ...


# ── InMemoryActorStore ───────────────────────────────────────────────────


def _derive_actor_id(pk: PublicKey) -> ActorID:
    """Derive a deterministic ActorID from a public key via SHA-256."""
    h = hashlib.sha256(pk.bytes_).digest()
    return ActorID(f"actor_{h[:16].hex()}")


def _now_nanos() -> int:
    return int(time.time() * 1_000_000_000)


class InMemoryActorStore:
    """Thread-safe in-memory implementation of ActorStore."""

    def __init__(self) -> None:
        self._lock = threading.Lock()
        self._actors: dict[str, Actor] = {}  # actor_id.value -> Actor
        self._by_key: dict[str, str] = {}  # hex(pk) -> actor_id.value
        self._ordered: list[str] = []  # insertion order of actor_id values

    def register(
        self, public_key: PublicKey, display_name: str, actor_type: ActorType
    ) -> Actor:
        with self._lock:
            key_hex = public_key.bytes_.hex()
            if key_hex in self._by_key:
                return self._actors[self._by_key[key_hex]]

            aid = _derive_actor_id(public_key)
            actor = Actor(
                id=aid,
                public_key=public_key,
                display_name=display_name,
                actor_type=actor_type,
                metadata=None,
                created_at=_now_nanos(),
                status=ActorStatus.ACTIVE,
            )
            self._actors[aid.value] = actor
            self._by_key[key_hex] = aid.value
            self._ordered.append(aid.value)
            return actor

    def get(self, actor_id: ActorID) -> Actor:
        with self._lock:
            a = self._actors.get(actor_id.value)
            if a is None:
                raise ActorNotFoundError(actor_id.value)
            return a

    def get_by_public_key(self, public_key: PublicKey) -> Actor:
        with self._lock:
            key_hex = public_key.bytes_.hex()
            aid_value = self._by_key.get(key_hex)
            if aid_value is None:
                raise ActorKeyNotFoundError(key_hex)
            return self._actors[aid_value]

    def update(self, actor_id: ActorID, updates: ActorUpdate) -> Actor:
        with self._lock:
            a = self._actors.get(actor_id.value)
            if a is None:
                raise ActorNotFoundError(actor_id.value)
            updated = a._with_updates(updates)
            self._actors[actor_id.value] = updated
            return updated

    def list(self, filter: ActorFilter) -> Page[Actor]:
        with self._lock:
            limit = filter.limit if filter.limit > 0 else 100

            start_idx = 0
            if filter.after.is_some():
                cursor_val = filter.after.unwrap()
                found = False
                for i, aid_val in enumerate(self._ordered):
                    if aid_val == cursor_val:
                        start_idx = i + 1
                        found = True
                        break
                if not found:
                    # Return empty page for invalid cursor
                    return Page((), Option.none(), False)

            items: list[Actor] = []
            for i in range(start_idx, len(self._ordered)):
                if len(items) >= limit:
                    break
                a = self._actors[self._ordered[i]]
                if filter.status.is_some() and a.status != filter.status.unwrap():
                    continue
                if filter.actor_type.is_some() and a.actor_type != filter.actor_type.unwrap():
                    continue
                items.append(a)

            has_more = False
            cursor: Option[Cursor] = Option.none()
            if len(items) == limit:
                last_id = items[-1].id.value
                last_idx = 0
                for i, aid_val in enumerate(self._ordered):
                    if aid_val == last_id:
                        last_idx = i
                        break
                for i in range(last_idx + 1, len(self._ordered)):
                    a = self._actors[self._ordered[i]]
                    if filter.status.is_some() and a.status != filter.status.unwrap():
                        continue
                    if filter.actor_type.is_some() and a.actor_type != filter.actor_type.unwrap():
                        continue
                    has_more = True
                    break
                if has_more:
                    cursor = Option.some(Cursor(last_id))

            return Page(items, cursor, has_more)

    def suspend(self, actor_id: ActorID, reason_event_id: EventID) -> Actor:
        with self._lock:
            a = self._actors.get(actor_id.value)
            if a is None:
                raise ActorNotFoundError(actor_id.value)
            new_status = a.status.transition_to(ActorStatus.SUSPENDED)
            updated = a._with_status(new_status)
            self._actors[actor_id.value] = updated
            _ = reason_event_id  # recorded on the event graph, not stored here
            return updated

    def reactivate(self, actor_id: ActorID, reason_event_id: EventID) -> Actor:
        with self._lock:
            a = self._actors.get(actor_id.value)
            if a is None:
                raise ActorNotFoundError(actor_id.value)
            new_status = a.status.transition_to(ActorStatus.ACTIVE)
            updated = a._with_status(new_status)
            self._actors[actor_id.value] = updated
            _ = reason_event_id
            return updated

    def memorial(self, actor_id: ActorID, reason_event_id: EventID) -> Actor:
        with self._lock:
            a = self._actors.get(actor_id.value)
            if a is None:
                raise ActorNotFoundError(actor_id.value)
            new_status = a.status.transition_to(ActorStatus.MEMORIAL)
            updated = a._with_status(new_status)
            self._actors[actor_id.value] = updated
            _ = reason_event_id
            return updated

    def actor_count(self) -> int:
        """Return the number of registered actors. For testing."""
        with self._lock:
            return len(self._actors)
