"""Graph facade — top-level API integrating store, actor store, bus, trust, authority, and decision."""

from __future__ import annotations

import threading
from dataclasses import dataclass
from typing import Any

from .actor import ActorStore, InMemoryActorStore
from .authority import AuthorityChain, AuthorityResult, DefaultAuthorityChain
from .bus import EventBus
from .errors import EventGraphError
from .event import Event, NoopSigner, Signer, create_bootstrap, create_event
from .store import InMemoryStore, Store
from .trust import DefaultTrustModel, TrustMetrics, TrustModel
from .types import (
    ActorID,
    ConversationID,
    EventID,
    EventType,
    Hash,
    Option,
)


# ── Errors ──────────────────────────────────────────────────────────────


class GraphNotStartedError(EventGraphError):
    """The graph has not been started."""

    def __init__(self) -> None:
        super().__init__("graph not started")


class GraphClosedError(EventGraphError):
    """The graph has been closed."""

    def __init__(self) -> None:
        super().__init__("graph is closed")


class AlreadyBootstrappedError(EventGraphError):
    """The graph has already been bootstrapped."""

    def __init__(self) -> None:
        super().__init__("graph already bootstrapped")


# ── GraphConfig ─────────────────────────────────────────────────────────


@dataclass(frozen=True, slots=True)
class GraphConfig:
    """Configuration for the Graph facade."""

    subscriber_buffer_size: int = 256
    fallback_to_mechanical: bool = True


# ── Query ───────────────────────────────────────────────────────────────


class Query:
    """Query interface over the event store, actor store, and trust model."""

    def __init__(
        self,
        store: Store | InMemoryStore,
        actor_store: ActorStore | InMemoryActorStore,
        trust_model: TrustModel | DefaultTrustModel,
    ) -> None:
        self._store = store
        self._actor_store = actor_store
        self._trust_model = trust_model

    def recent(self, limit: int = 50) -> list[Event]:
        """Return the most recent events, newest first."""
        if hasattr(self._store, "recent"):
            return self._store.recent(limit)
        # Fallback: not available on protocol, but InMemoryStore has it
        raise NotImplementedError("store does not support recent()")

    def by_type(self, event_type: EventType, limit: int = 50) -> list[Event]:
        """Return events of the given type, newest first, up to limit."""
        return self._store.by_type(event_type, limit)

    def by_source(self, source: ActorID, limit: int = 50) -> list[Event]:
        """Return events from the given source, newest first, up to limit."""
        return self._store.by_source(source, limit)

    def by_conversation(self, conversation_id: ConversationID, limit: int = 50) -> list[Event]:
        """Return events in the given conversation, newest first, up to limit."""
        return self._store.by_conversation(conversation_id, limit)

    def ancestors(self, event_id: EventID, max_depth: int = 10) -> list[Event]:
        """Return causal ancestors via BFS on causes, up to max_depth levels."""
        return self._store.ancestors(event_id, max_depth)

    def descendants(self, event_id: EventID, max_depth: int = 10) -> list[Event]:
        """Return causal descendants via BFS, up to max_depth levels."""
        return self._store.descendants(event_id, max_depth)

    def trust_score(self, actor: ActorID) -> TrustMetrics:
        """Return current trust metrics for an actor."""
        return self._trust_model.score(actor)

    def trust_between(self, from_actor: ActorID, to_actor: ActorID) -> TrustMetrics:
        """Return directed trust from one actor to another."""
        return self._trust_model.between(from_actor, to_actor)

    def actor(self, actor_id: ActorID) -> Any:
        """Return an actor by ID."""
        return self._actor_store.get(actor_id)

    def event_count(self) -> int:
        """Return the total number of events in the store."""
        return self._store.count()


# ── Graph ───────────────────────────────────────────────────────────────


class Graph:
    """Top-level facade: Evaluate(), Record(), Query(), Start(), Close().

    Integrates store, actor store, bus, trust, authority, and decision
    into a single cohesive API.
    """

    def __init__(
        self,
        store: Store | InMemoryStore,
        actor_store: ActorStore | InMemoryActorStore,
        trust_model: TrustModel | DefaultTrustModel | None = None,
        authority_chain: AuthorityChain | DefaultAuthorityChain | None = None,
        signer: Signer | None = None,
        config: GraphConfig | None = None,
    ) -> None:
        self._config = config or GraphConfig()
        self._store = store
        self._actor_store = actor_store
        self._signer: Signer = signer or NoopSigner()
        self._trust_model: TrustModel | DefaultTrustModel = trust_model or DefaultTrustModel()
        self._authority_chain: AuthorityChain | DefaultAuthorityChain = (
            authority_chain or DefaultAuthorityChain(self._trust_model)
        )
        self._bus = EventBus(store, buffer_size=self._config.subscriber_buffer_size)
        self._lock = threading.Lock()
        self._started = False
        self._closed = False

    # ── Properties ──────────────────────────────────────────────────────

    @property
    def store(self) -> Store | InMemoryStore:
        """The underlying event store."""
        return self._store

    @property
    def actor_store(self) -> ActorStore | InMemoryActorStore:
        """The underlying actor store."""
        return self._actor_store

    @property
    def bus(self) -> EventBus:
        """The event bus for pub/sub."""
        return self._bus

    # ── Lifecycle ───────────────────────────────────────────────────────

    def start(self) -> None:
        """Start the graph. Idempotent."""
        with self._lock:
            if self._closed:
                raise GraphClosedError()
            self._started = True

    def close(self) -> None:
        """Close the graph. Idempotent. Closes the bus and store."""
        with self._lock:
            if self._closed:
                return
            self._closed = True

        self._bus.close()
        self._store.close()

    # ── Guard helpers ───────────────────────────────────────────────────

    def _check_running(self) -> None:
        """Raise if the graph is not started or is closed. Caller must hold self._lock."""
        if not self._started:
            raise GraphNotStartedError()
        if self._closed:
            raise GraphClosedError()

    # ── Bootstrap ───────────────────────────────────────────────────────

    def bootstrap(
        self,
        system_actor: ActorID,
        signer: Signer | None = None,
    ) -> Event:
        """Create the genesis event if the store is empty.

        Raises GraphNotStartedError if start() has not been called.
        Raises GraphClosedError if close() has been called.
        Raises AlreadyBootstrappedError if the store already contains events.
        """
        with self._lock:
            self._check_running()

            if self._store.count() > 0:
                raise AlreadyBootstrappedError()

            s = signer or self._signer
            genesis = create_bootstrap(source=system_actor, signer=s)
            self._store.append(genesis)
            self._bus.publish(genesis)
            return genesis

    # ── Record ──────────────────────────────────────────────────────────

    def record(
        self,
        event_type: EventType,
        source: ActorID,
        content: dict[str, Any],
        causes: list[EventID],
        conversation_id: ConversationID,
        signer: Signer | None = None,
    ) -> Event:
        """Create an event, append to store, and publish to bus.

        Uses the store's head event hash as prev_hash.

        Raises GraphNotStartedError if start() has not been called.
        Raises GraphClosedError if close() has been called.
        """
        with self._lock:
            self._check_running()

            head = self._store.head()
            if head.is_some():
                prev_hash = head.unwrap().hash
            else:
                prev_hash = Hash.zero()

            s = signer or self._signer
            event = create_event(
                event_type=event_type,
                source=source,
                content=content,
                causes=causes,
                conversation_id=conversation_id,
                prev_hash=prev_hash,
                signer=s,
            )
            self._store.append(event)
            self._bus.publish(event)
            return event

    # ── Evaluate ────────────────────────────────────────────────────────

    def evaluate(self, actor: ActorID, action: str) -> AuthorityResult:
        """Evaluate authority for an actor performing an action.

        Delegates to the authority chain.
        """
        return self._authority_chain.evaluate(actor, action)

    # ── Query ───────────────────────────────────────────────────────────

    def query(self) -> Query:
        """Return a Query object for reading events, actors, and trust."""
        return Query(self._store, self._actor_store, self._trust_model)
