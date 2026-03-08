"""Always-valid domain types — constrained numerics, typed IDs, Option, NonEmpty, Page, Cursor."""

from __future__ import annotations

import math
import re
from dataclasses import dataclass
from typing import Generic, Iterator, Sequence, TypeVar

from .errors import EmptyRequiredError, InvalidFormatError, OutOfRangeError

T = TypeVar("T")

# ── Option ────────────────────────────────────────────────────────────────

@dataclass(frozen=True, slots=True)
class Option(Generic[T]):
    """Explicit optionality — Some(value) or NONE. No None/null ambiguity."""

    _value: T | None
    _present: bool

    @staticmethod
    def some(value: T) -> Option[T]:
        return Option(_value=value, _present=True)

    @staticmethod
    def none() -> Option[T]:
        return Option(_value=None, _present=False)

    def is_some(self) -> bool:
        return self._present

    def is_none(self) -> bool:
        return not self._present

    def unwrap(self) -> T:
        if not self._present:
            raise ValueError("unwrap called on None Option")
        return self._value  # type: ignore[return-value]

    def unwrap_or(self, default: T) -> T:
        return self._value if self._present else default  # type: ignore[return-value]


# ── NonEmpty ──────────────────────────────────────────────────────────────

@dataclass(frozen=True, slots=True)
class NonEmpty(Generic[T]):
    """A collection with at least one element."""

    _items: tuple[T, ...]

    @staticmethod
    def of(items: Sequence[T]) -> NonEmpty[T]:
        if not items:
            raise ValueError("NonEmpty requires at least one element")
        return NonEmpty(_items=tuple(items))

    def __iter__(self) -> Iterator[T]:
        return iter(self._items)

    def __len__(self) -> int:
        return len(self._items)

    def __getitem__(self, index: int) -> T:
        return self._items[index]

    def items(self) -> tuple[T, ...]:
        return self._items


# ── Constrained numerics ──────────────────────────────────────────────────

@dataclass(frozen=True, slots=True, order=True)
class Score:
    """Float constrained to [0.0, 1.0]."""

    _value: float

    def __init__(self, value: float) -> None:
        if math.isnan(value) or value < 0.0 or value > 1.0:
            raise OutOfRangeError("Score", value, 0.0, 1.0)
        if value == 0.0 and math.copysign(1.0, value) < 0:
            value = 0.0  # normalize -0.0
        object.__setattr__(self, "_value", value)

    @property
    def value(self) -> float:
        return self._value


@dataclass(frozen=True, slots=True, order=True)
class Weight:
    """Float constrained to [-1.0, 1.0]."""

    _value: float

    def __init__(self, value: float) -> None:
        if math.isnan(value) or value < -1.0 or value > 1.0:
            raise OutOfRangeError("Weight", value, -1.0, 1.0)
        if value == 0.0 and math.copysign(1.0, value) < 0:
            value = 0.0
        object.__setattr__(self, "_value", value)

    @property
    def value(self) -> float:
        return self._value


@dataclass(frozen=True, slots=True, order=True)
class Activation:
    """Float constrained to [0.0, 1.0]."""

    _value: float

    def __init__(self, value: float) -> None:
        if math.isnan(value) or value < 0.0 or value > 1.0:
            raise OutOfRangeError("Activation", value, 0.0, 1.0)
        if value == 0.0 and math.copysign(1.0, value) < 0:
            value = 0.0
        object.__setattr__(self, "_value", value)

    @property
    def value(self) -> float:
        return self._value


@dataclass(frozen=True, slots=True, order=True)
class Layer:
    """Int constrained to [0, 13]."""

    _value: int

    def __init__(self, value: int) -> None:
        if value < 0 or value > 13:
            raise OutOfRangeError("Layer", float(value), 0, 13)
        object.__setattr__(self, "_value", value)

    @property
    def value(self) -> int:
        return self._value


@dataclass(frozen=True, slots=True, order=True)
class Cadence:
    """Int constrained to [1, ∞)."""

    _value: int

    def __init__(self, value: int) -> None:
        if value < 1:
            raise OutOfRangeError("Cadence", float(value), 1, float("inf"))
        object.__setattr__(self, "_value", value)

    @property
    def value(self) -> int:
        return self._value


# ── Typed IDs ─────────────────────────────────────────────────────────────

_UUID_V7_RE = re.compile(
    r"^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
)
_UUID_RE = re.compile(
    r"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
)
_EVENT_TYPE_RE = re.compile(r"^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$")
_DOMAIN_SCOPE_RE = re.compile(r"^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)*$")
_SUBSCRIPTION_RE = re.compile(
    r"^(\*|[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*(\.\*)?)$"
)


def _str_id(name: str) -> type:
    """Create a frozen dataclass for a simple non-empty string ID."""

    @dataclass(frozen=True, slots=True)
    class _ID:
        _value: str

        def __init__(self, value: str) -> None:
            if not value:
                raise EmptyRequiredError(name)
            object.__setattr__(self, "_value", value)

        @property
        def value(self) -> str:
            return self._value

        def __str__(self) -> str:
            return self._value

    _ID.__name__ = name
    _ID.__qualname__ = name
    return _ID


ActorID = _str_id("ActorID")
ConversationID = _str_id("ConversationID")
SystemURI = _str_id("SystemURI")
PrimitiveID = _str_id("PrimitiveID")


@dataclass(frozen=True, slots=True)
class EventID:
    """UUID v7 event identifier."""

    _value: str

    def __init__(self, value: str) -> None:
        value = value.lower()
        if not _UUID_V7_RE.match(value):
            raise InvalidFormatError("EventID", value, "UUID v7")
        object.__setattr__(self, "_value", value)

    @property
    def value(self) -> str:
        return self._value

    def __str__(self) -> str:
        return self._value


@dataclass(frozen=True, slots=True)
class EdgeID:
    """UUID v7 edge identifier (the EventID of the edge-creating event)."""

    _value: str

    def __init__(self, value: str) -> None:
        value = value.lower()
        if not _UUID_V7_RE.match(value):
            raise InvalidFormatError("EdgeID", value, "UUID v7")
        object.__setattr__(self, "_value", value)

    @property
    def value(self) -> str:
        return self._value

    def __str__(self) -> str:
        return self._value


@dataclass(frozen=True, slots=True)
class Hash:
    """SHA-256 hex string (64 characters)."""

    _value: str

    def __init__(self, value: str) -> None:
        value = value.lower()
        if len(value) != 64 or not all(c in "0123456789abcdef" for c in value):
            raise InvalidFormatError("Hash", value, "64 hex characters (SHA-256)")
        object.__setattr__(self, "_value", value)

    @staticmethod
    def zero() -> Hash:
        return Hash("0" * 64)

    @property
    def value(self) -> str:
        return self._value

    def is_zero(self) -> bool:
        return self._value == "0" * 64

    def __str__(self) -> str:
        return self._value


@dataclass(frozen=True, slots=True)
class EventType:
    """Dot-separated lowercase event type (e.g., 'trust.updated')."""

    _value: str

    def __init__(self, value: str) -> None:
        if not _EVENT_TYPE_RE.match(value):
            raise InvalidFormatError(
                "EventType", value, "dot-separated lowercase segments"
            )
        object.__setattr__(self, "_value", value)

    @property
    def value(self) -> str:
        return self._value

    def __str__(self) -> str:
        return self._value


@dataclass(frozen=True, slots=True)
class DomainScope:
    """Trust/authority domain (e.g., 'code_review')."""

    _value: str

    def __init__(self, value: str) -> None:
        if not _DOMAIN_SCOPE_RE.match(value):
            raise InvalidFormatError(
                "DomainScope", value, "lowercase dot/underscore-separated namespace"
            )
        object.__setattr__(self, "_value", value)

    @property
    def value(self) -> str:
        return self._value

    def __str__(self) -> str:
        return self._value


@dataclass(frozen=True, slots=True)
class SubscriptionPattern:
    """Glob pattern for event type matching (e.g., 'trust.*')."""

    _value: str

    def __init__(self, value: str) -> None:
        if not _SUBSCRIPTION_RE.match(value):
            raise InvalidFormatError(
                "SubscriptionPattern", value,
                "dot-separated segments with optional trailing .* or bare *"
            )
        object.__setattr__(self, "_value", value)

    def matches(self, et: EventType) -> bool:
        if self._value == "*":
            return True
        if self._value.endswith(".*"):
            prefix = self._value[:-2]
            return et.value == prefix or et.value.startswith(prefix + ".")
        return self._value == et.value

    @property
    def value(self) -> str:
        return self._value


@dataclass(frozen=True, slots=True)
class EnvelopeID:
    """UUID envelope identifier."""

    _value: str

    def __init__(self, value: str) -> None:
        value = value.lower()
        if not _UUID_RE.match(value):
            raise InvalidFormatError("EnvelopeID", value, "UUID")
        object.__setattr__(self, "_value", value)

    @property
    def value(self) -> str:
        return self._value

    def __str__(self) -> str:
        return self._value


@dataclass(frozen=True, slots=True)
class TreatyID:
    """UUID treaty identifier."""

    _value: str

    def __init__(self, value: str) -> None:
        value = value.lower()
        if not _UUID_RE.match(value):
            raise InvalidFormatError("TreatyID", value, "UUID")
        object.__setattr__(self, "_value", value)

    @property
    def value(self) -> str:
        return self._value

    def __str__(self) -> str:
        return self._value


@dataclass(frozen=True, slots=True)
class PublicKey:
    """Ed25519 public key (32 bytes)."""

    _value: bytes

    def __init__(self, value: bytes) -> None:
        if len(value) != 32:
            raise InvalidFormatError(
                "PublicKey", value.hex(), "32 bytes (Ed25519 public key)"
            )
        object.__setattr__(self, "_value", bytes(value))

    @property
    def bytes_(self) -> bytes:
        return self._value

    def __str__(self) -> str:
        return self._value.hex()


# ── Cursor & Page ─────────────────────────────────────────────────────────

@dataclass(frozen=True, slots=True)
class Cursor:
    """Opaque pagination token."""

    _value: str

    def __init__(self, value: str) -> None:
        object.__setattr__(self, "_value", value)

    @property
    def value(self) -> str:
        return self._value

    def __str__(self) -> str:
        return self._value


@dataclass(frozen=True, slots=True)
class Page(Generic[T]):
    """Paginated result set with cursor-based navigation."""

    _items: tuple[T, ...]
    _cursor: Option[Cursor]
    _has_more: bool

    def __init__(
        self,
        items: Sequence[T],
        cursor: Option[Cursor],
        has_more: bool,
    ) -> None:
        object.__setattr__(self, "_items", tuple(items))
        object.__setattr__(self, "_cursor", cursor)
        object.__setattr__(self, "_has_more", has_more)

    def items(self) -> tuple[T, ...]:
        return self._items

    def cursor(self) -> Option[Cursor]:
        return self._cursor

    def has_more(self) -> bool:
        return self._has_more

    def __len__(self) -> int:
        return len(self._items)


@dataclass(frozen=True, slots=True)
class Signature:
    """Ed25519 signature (64 bytes)."""

    _value: bytes

    def __init__(self, value: bytes) -> None:
        if len(value) != 64:
            raise InvalidFormatError(
                "Signature", value.hex(), "64 bytes (Ed25519 signature)"
            )
        object.__setattr__(self, "_value", bytes(value))

    @property
    def bytes_(self) -> bytes:
        return self._value

    def __str__(self) -> str:
        return self._value.hex()
