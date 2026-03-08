"""Tick engine — ripple-wave processor."""

from __future__ import annotations

import threading
import time
from dataclasses import dataclass, field
from typing import Any, Callable

from .event import Event, NoopSigner, create_event
from .primitive import (
    LIFECYCLE_ACTIVE,
    LIFECYCLE_EMITTING,
    LIFECYCLE_PROCESSING,
    AddEvent,
    Mutation,
    Primitive,
    Registry,
    Snapshot,
    UpdateActivation,
    UpdateLifecycle,
    UpdateState,
)
from .store import InMemoryStore
from .types import (
    ActorID,
    ConversationID,
    EventType,
    Hash,
    NonEmpty,
    PrimitiveID,
    SubscriptionPattern,
)


@dataclass(frozen=True)
class TickConfig:
    max_waves_per_tick: int = 10


@dataclass(frozen=True)
class TickResult:
    tick: int
    waves: int
    mutations: int
    quiesced: bool
    duration_ms: float
    errors: list[str] = field(default_factory=list)


class TickEngine:
    """Ripple-wave tick processor.

    Each tick:
    1. Snapshot all primitive states
    2. Distribute pending events to subscribing primitives
    3. Invoke each primitive's process() (subject to cadence + lifecycle)
    4. Collect mutations
    5. Apply atomically — new events become input for next wave
    6. Repeat until quiescence or max waves
    """

    def __init__(
        self,
        registry: Registry,
        store: InMemoryStore,
        config: TickConfig | None = None,
        publisher: Callable[[Event], None] | None = None,
    ) -> None:
        self._lock = threading.Lock()
        self._registry = registry
        self._store = store
        self._config = config or TickConfig()
        self._publisher = publisher
        self._current_tick = 0
        self._signer = NoopSigner()

    def tick(self, pending_events: list[Event] | None = None) -> TickResult:
        """Run a single tick. Returns the result."""
        with self._lock:
            start = time.monotonic()
            self._current_tick += 1
            tick_num = self._current_tick

            wave_events = list(pending_events or [])
            total_mutations = 0
            errors: list[str] = []
            quiesced = False

            # Track which primitives were invoked this tick (for cadence)
            invoked_this_tick: set[str] = set()
            waves_run = 0

            for wave in range(self._config.max_waves_per_tick):
                if not wave_events and wave > 0:
                    quiesced = True
                    break

                waves_run = wave + 1

                # Build snapshot
                snapshot = Snapshot(
                    tick=tick_num,
                    primitives=self._registry.all_states(),
                    pending_events=wave_events,
                    recent_events=self._store.recent(50),
                )

                # Process primitives in layer order
                new_events: list[Event] = []
                all_mutations: list[Mutation] = []

                for prim in self._registry.all():
                    pid = prim.id()
                    lifecycle = self._registry.lifecycle(pid)
                    if lifecycle != LIFECYCLE_ACTIVE:
                        continue

                    # Cadence check — only between ticks, not within waves
                    if pid.value not in invoked_this_tick:
                        last = self._registry.last_tick(pid)
                        if tick_num - last < prim.cadence().value:
                            continue

                    # Filter events by subscription
                    matched = _filter_events(wave_events, prim.subscriptions())
                    if not matched and wave > 0:
                        continue

                    # Transition to processing
                    try:
                        self._registry.set_lifecycle(pid, LIFECYCLE_PROCESSING)
                    except ValueError:
                        continue

                    try:
                        mutations = prim.process(tick_num, matched, snapshot)
                        all_mutations.extend(mutations)
                    except Exception as e:
                        errors.append(f"{pid.value}: {e}")

                    # Transition back to active (processing -> emitting -> active)
                    try:
                        self._registry.set_lifecycle(pid, LIFECYCLE_EMITTING)
                        self._registry.set_lifecycle(pid, LIFECYCLE_ACTIVE)
                    except ValueError as e:
                        errors.append(f"{pid.value} lifecycle: {e}")

                    invoked_this_tick.add(pid.value)

                # Apply mutations
                for m in all_mutations:
                    total_mutations += 1
                    try:
                        ev = self._apply_mutation(m)
                        if ev is not None:
                            new_events.append(ev)
                    except Exception as e:
                        errors.append(f"mutation: {e}")

                wave_events = new_events

            # Set last_tick for all invoked primitives
            for pid_val in invoked_this_tick:
                self._registry.set_last_tick(PrimitiveID(pid_val), tick_num)

            # If we hit max waves without quiescing
            if not quiesced and waves_run >= self._config.max_waves_per_tick:
                quiesced = False

            elapsed = (time.monotonic() - start) * 1000
            return TickResult(
                tick=tick_num,
                waves=waves_run,
                mutations=total_mutations,
                quiesced=quiesced,
                duration_ms=elapsed,
                errors=errors,
            )

    def _apply_mutation(self, m: Mutation) -> Event | None:
        if isinstance(m, AddEvent):
            head = self._store.head()
            prev_hash = head.unwrap().hash if head.is_some() else Hash.zero()
            event = create_event(
                event_type=m.type,
                source=m.source,
                content=m.content,
                causes=m.causes,
                conversation_id=m.conversation_id,
                prev_hash=prev_hash,
                signer=self._signer,
            )
            self._store.append(event)
            if self._publisher:
                self._publisher(event)
            return event
        elif isinstance(m, UpdateState):
            self._registry.update_state(m.primitive_id, m.key, m.value)
        elif isinstance(m, UpdateActivation):
            self._registry.set_activation(m.primitive_id, m.level)
        elif isinstance(m, UpdateLifecycle):
            self._registry.set_lifecycle(m.primitive_id, m.state)
        return None


def _filter_events(
    events: list[Event], patterns: list[SubscriptionPattern]
) -> list[Event]:
    """Return events matching any of the subscription patterns."""
    result = []
    for ev in events:
        for pat in patterns:
            if pat.matches(ev.type):
                result.append(ev)
                break
    return result
