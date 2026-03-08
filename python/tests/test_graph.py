"""Tests for the Graph facade module."""

from __future__ import annotations

import pytest

from eventgraph import (
    ActorID,
    ConversationID,
    EventID,
    EventType,
    Graph,
    GraphConfig,
    InMemoryActorStore,
    InMemoryStore,
    NoopSigner,
    Query,
    new_event_id,
)
from eventgraph.graph import (
    AlreadyBootstrappedError,
    GraphClosedError,
    GraphNotStartedError,
)
from eventgraph.trust import DefaultTrustModel
from eventgraph.authority import DefaultAuthorityChain
from eventgraph.decision import AuthorityLevel


# ── Helpers ─────────────────────────────────────────────────────────────


def _make_graph(**kwargs) -> Graph:
    """Create a Graph with default in-memory components."""
    defaults = dict(
        store=InMemoryStore(),
        actor_store=InMemoryActorStore(),
    )
    defaults.update(kwargs)
    return Graph(**defaults)


def _bootstrapped_graph(**kwargs) -> Graph:
    """Create a started and bootstrapped Graph."""
    g = _make_graph(**kwargs)
    g.start()
    g.bootstrap(ActorID("system"))
    return g


SYSTEM = ActorID("system")
ALICE = ActorID("alice")


# ── test_start_and_close ────────────────────────────────────────────────


def test_start_and_close():
    g = _make_graph()

    # Not started yet, so record should fail
    with pytest.raises(GraphNotStartedError):
        g.record(
            event_type=EventType("test.event"),
            source=SYSTEM,
            content={"key": "value"},
            causes=[new_event_id()],
            conversation_id=ConversationID("conv-1"),
        )

    g.start()
    # start() is idempotent
    g.start()

    g.close()
    # close() is idempotent
    g.close()


# ── test_bootstrap_creates_genesis ──────────────────────────────────────


def test_bootstrap_creates_genesis():
    g = _make_graph()
    g.start()
    genesis = g.bootstrap(SYSTEM)

    assert genesis is not None
    assert genesis.type.value == "system.bootstrapped"
    assert genesis.source.value == "system"
    assert g.store.count() == 1


# ── test_bootstrap_requires_start ───────────────────────────────────────


def test_bootstrap_requires_start():
    g = _make_graph()
    with pytest.raises(GraphNotStartedError):
        g.bootstrap(SYSTEM)


# ── test_bootstrap_idempotent_error ─────────────────────────────────────


def test_bootstrap_idempotent_error():
    g = _make_graph()
    g.start()
    g.bootstrap(SYSTEM)
    with pytest.raises(AlreadyBootstrappedError):
        g.bootstrap(SYSTEM)


# ── test_record_creates_event ───────────────────────────────────────────


def test_record_creates_event():
    g = _bootstrapped_graph()
    genesis = g.store.head().unwrap()

    event = g.record(
        event_type=EventType("test.recorded"),
        source=ALICE,
        content={"message": "hello"},
        causes=[genesis.id],
        conversation_id=ConversationID("conv-1"),
    )

    assert event is not None
    assert event.type.value == "test.recorded"
    assert event.source.value == "alice"
    assert event.prev_hash == genesis.hash
    assert g.store.count() == 2


# ── test_record_requires_start ──────────────────────────────────────────


def test_record_requires_start():
    g = _make_graph()
    with pytest.raises(GraphNotStartedError):
        g.record(
            event_type=EventType("test.event"),
            source=SYSTEM,
            content={},
            causes=[new_event_id()],
            conversation_id=ConversationID("conv-1"),
        )


# ── test_record_after_close_fails ───────────────────────────────────────


def test_record_after_close_fails():
    g = _bootstrapped_graph()
    g.close()
    with pytest.raises(GraphClosedError):
        g.record(
            event_type=EventType("test.event"),
            source=SYSTEM,
            content={},
            causes=[new_event_id()],
            conversation_id=ConversationID("conv-1"),
        )


# ── test_evaluate_delegates_to_authority ────────────────────────────────


def test_evaluate_delegates_to_authority():
    g = _bootstrapped_graph()
    result = g.evaluate(ALICE, "some.action")

    assert result is not None
    assert result.level == AuthorityLevel.NOTIFICATION  # default policy
    assert result.weight.value == 1.0
    assert len(result.chain) == 1
    assert result.chain[0].actor.value == "alice"


# ── test_query_by_type ──────────────────────────────────────────────────


def test_query_by_type():
    g = _bootstrapped_graph()
    genesis = g.store.head().unwrap()

    g.record(
        event_type=EventType("test.alpha"),
        source=ALICE,
        content={"n": 1},
        causes=[genesis.id],
        conversation_id=ConversationID("conv-1"),
    )
    e2 = g.record(
        event_type=EventType("test.beta"),
        source=ALICE,
        content={"n": 2},
        causes=[genesis.id],
        conversation_id=ConversationID("conv-1"),
    )

    q = g.query()
    alphas = q.by_type(EventType("test.alpha"))
    betas = q.by_type(EventType("test.beta"))

    assert len(alphas) == 1
    assert alphas[0].type.value == "test.alpha"
    assert len(betas) == 1
    assert betas[0].type.value == "test.beta"


# ── test_query_by_source ────────────────────────────────────────────────


def test_query_by_source():
    g = _bootstrapped_graph()
    genesis = g.store.head().unwrap()

    g.record(
        event_type=EventType("test.event"),
        source=ALICE,
        content={"from": "alice"},
        causes=[genesis.id],
        conversation_id=ConversationID("conv-1"),
    )
    g.record(
        event_type=EventType("test.event"),
        source=SYSTEM,
        content={"from": "system"},
        causes=[genesis.id],
        conversation_id=ConversationID("conv-1"),
    )

    q = g.query()
    alice_events = q.by_source(ALICE)
    # genesis is from SYSTEM, plus the one we recorded
    system_events = q.by_source(SYSTEM)

    assert len(alice_events) == 1
    assert alice_events[0].source.value == "alice"
    assert len(system_events) == 2  # genesis + recorded


# ── test_query_ancestors ────────────────────────────────────────────────


def test_query_ancestors():
    g = _bootstrapped_graph()
    genesis = g.store.head().unwrap()

    e1 = g.record(
        event_type=EventType("test.event"),
        source=ALICE,
        content={"step": 1},
        causes=[genesis.id],
        conversation_id=ConversationID("conv-1"),
    )
    e2 = g.record(
        event_type=EventType("test.event"),
        source=ALICE,
        content={"step": 2},
        causes=[e1.id],
        conversation_id=ConversationID("conv-1"),
    )

    q = g.query()
    ancestors = q.ancestors(e2.id, max_depth=5)

    # e2's cause is e1, e1's cause is genesis
    ancestor_ids = {a.id.value for a in ancestors}
    assert e1.id.value in ancestor_ids
    assert genesis.id.value in ancestor_ids


# ── test_graph_defaults ─────────────────────────────────────────────────


def test_graph_defaults():
    """Graph creates default trust model and authority chain if not provided."""
    store = InMemoryStore()
    actor_store = InMemoryActorStore()
    g = Graph(store=store, actor_store=actor_store)

    # Should have created defaults internally
    assert g.store is store
    assert g.actor_store is actor_store
    assert g.bus is not None

    # Query should work with the defaults
    g.start()
    g.bootstrap(SYSTEM)
    q = g.query()
    assert q.event_count() == 1

    # Trust query should return default metrics
    metrics = q.trust_score(ALICE)
    assert metrics.overall.value == 0.0  # default initial trust

    # Authority should return notification level by default
    result = g.evaluate(ALICE, "any.action")
    assert result.level == AuthorityLevel.NOTIFICATION

    g.close()


# ── test_query_event_count ──────────────────────────────────────────────


def test_query_event_count():
    g = _bootstrapped_graph()
    genesis = g.store.head().unwrap()

    q = g.query()
    assert q.event_count() == 1

    g.record(
        event_type=EventType("test.event"),
        source=ALICE,
        content={},
        causes=[genesis.id],
        conversation_id=ConversationID("conv-1"),
    )

    assert q.event_count() == 2


# ── test_query_trust_between ────────────────────────────────────────────


def test_query_trust_between():
    g = _bootstrapped_graph()
    q = g.query()
    metrics = q.trust_between(ALICE, SYSTEM)
    assert metrics is not None
    assert metrics.overall.value == 0.0  # default, no evidence
