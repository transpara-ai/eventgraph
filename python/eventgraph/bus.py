"""Event bus — pub/sub fan-out with non-blocking delivery."""

from __future__ import annotations

import threading
from typing import Callable

from .event import Event
from .store import InMemoryStore
from .types import EventType, SubscriptionPattern


class EventBus:
    """Thread-safe event bus with per-subscriber buffering.

    Slow subscribers get dropped events, not blocked publishers.
    """

    def __init__(self, store: InMemoryStore, buffer_size: int = 256) -> None:
        self._store = store
        self._buffer_size = max(buffer_size, 1)
        self._lock = threading.Lock()
        self._subs: dict[int, _Subscription] = {}
        self._next_id = 0
        self._closed = False

    @property
    def store(self) -> InMemoryStore:
        return self._store

    def subscribe(
        self, pattern: SubscriptionPattern, handler: Callable[[Event], None]
    ) -> int:
        """Register a handler for events matching the pattern.

        Returns subscription ID, or -1 if the bus is closed.
        """
        with self._lock:
            if self._closed:
                return -1
            self._next_id += 1
            sub_id = self._next_id
            sub = _Subscription(sub_id, pattern, handler, self._buffer_size)
            self._subs[sub_id] = sub
            sub.start()
            return sub_id

    def unsubscribe(self, sub_id: int) -> None:
        with self._lock:
            sub = self._subs.pop(sub_id, None)
        if sub is not None:
            sub.close()

    def publish(self, event: Event) -> None:
        """Deliver event to all matching subscribers. Non-blocking."""
        with self._lock:
            if self._closed:
                return
            subs = list(self._subs.values())

        for sub in subs:
            sub.deliver(event)

    def close(self, timeout: float = 30.0) -> None:
        """Stop all subscriber threads."""
        with self._lock:
            if self._closed:
                return
            self._closed = True
            subs = list(self._subs.values())
            self._subs.clear()

        for sub in subs:
            sub.close()
        for sub in subs:
            sub.join(timeout=timeout)


class _Subscription:
    """Per-subscriber delivery thread with bounded buffer."""

    def __init__(
        self,
        sub_id: int,
        pattern: SubscriptionPattern,
        handler: Callable[[Event], None],
        buffer_size: int,
    ) -> None:
        self.id = sub_id
        self.pattern = pattern
        self.handler = handler
        self._buffer: list[Event] = []
        self._buffer_size = buffer_size
        self._cond = threading.Condition()
        self._closed = False
        self._thread: threading.Thread | None = None
        self.last_error: Exception | None = None

    def start(self) -> None:
        self._thread = threading.Thread(target=self._run, daemon=True)
        self._thread.start()

    def deliver(self, event: Event) -> None:
        """Non-blocking delivery. Drops if buffer full or pattern doesn't match."""
        if not self.pattern.matches(event.type):
            return
        with self._cond:
            if self._closed:
                return
            if len(self._buffer) < self._buffer_size:
                self._buffer.append(event)
                self._cond.notify()

    def close(self) -> None:
        with self._cond:
            self._closed = True
            self._cond.notify()

    def join(self, timeout: float = 30.0) -> None:
        if self._thread is not None:
            self._thread.join(timeout=timeout)

    def _run(self) -> None:
        while True:
            with self._cond:
                while not self._buffer and not self._closed:
                    self._cond.wait()
                if self._closed and not self._buffer:
                    return
                event = self._buffer.pop(0)

            try:
                self.handler(event)
            except Exception as e:
                self.last_error = e
