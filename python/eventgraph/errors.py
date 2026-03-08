"""Typed domain errors for EventGraph."""

from __future__ import annotations


class EventGraphError(Exception):
    """Base error for all EventGraph errors."""


class ValidationError(EventGraphError):
    """Base for validation failures at construction time."""


class OutOfRangeError(ValidationError):
    """A numeric value is outside its valid range."""

    def __init__(self, field: str, value: float, min_val: float, max_val: float) -> None:
        self.field = field
        self.value = value
        self.min_val = min_val
        self.max_val = max_val
        super().__init__(f"{field}: {value} outside [{min_val}, {max_val}]")


class EmptyRequiredError(ValidationError):
    """A required string field is empty."""

    def __init__(self, field: str) -> None:
        self.field = field
        super().__init__(f"{field}: cannot be empty")


class InvalidFormatError(ValidationError):
    """A string field does not match the expected format."""

    def __init__(self, field: str, value: str, expected: str) -> None:
        self.field = field
        self.value = value
        self.expected = expected
        super().__init__(f"{field}: '{value}' does not match expected format: {expected}")


class InvalidTransitionError(ValidationError):
    """A state machine transition is invalid."""

    def __init__(self, field: str, from_state: str, to_state: str) -> None:
        self.field = field
        self.from_state = from_state
        self.to_state = to_state
        super().__init__(f"{field}: invalid transition {from_state} -> {to_state}")


class StoreError(EventGraphError):
    """Base for store-related errors."""


class EventNotFoundError(StoreError):
    """An event was not found in the store."""

    def __init__(self, event_id: str) -> None:
        self.event_id = event_id
        super().__init__(f"event not found: {event_id}")


class ChainIntegrityError(StoreError):
    """The hash chain is broken."""

    def __init__(self, position: int, detail: str) -> None:
        self.position = position
        self.detail = detail
        super().__init__(f"chain integrity error at position {position}: {detail}")
