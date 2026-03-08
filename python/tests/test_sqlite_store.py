"""Tests for SQLiteStore — verifies it satisfies the Store protocol."""

import pytest
from eventgraph.event import create_bootstrap, create_event, NoopSigner
from eventgraph.sqlite_store import SQLiteStore
from eventgraph.types import (
    ActorID, ConversationID, EventID, EventType, Hash,
)


@pytest.fixture
def store():
    s = SQLiteStore(":memory:")
    yield s
    s.close()


@pytest.fixture
def signer():
    return NoopSigner()


@pytest.fixture
def actor():
    return ActorID("actor_00000000000000000000000000000001")


@pytest.fixture
def conv():
    return ConversationID("conv_00000000000000000000000000000001")


def bootstrap(store, signer, actor):
    ev = create_bootstrap(actor, signer)
    return store.append(ev)


def chained_event(store, signer, actor, conv, causes):
    head = store.head()
    prev_hash = head.unwrap().hash if head.is_some() else Hash.zero()  # type: ignore
    return create_event(
        EventType("trust.updated"), actor,
        {"actor": "test", "score": 0.5},
        causes, conv, prev_hash, signer,
    )


class TestSQLiteStoreBasic:
    def test_append_and_get(self, store, signer, actor):
        ev = bootstrap(store, signer, actor)
        got = store.get(ev.id)
        assert got.id == ev.id

    def test_get_not_found(self, store):
        with pytest.raises(Exception):
            store.get(EventID("019462a0-0000-7000-8000-000000000099"))

    def test_idempotent(self, store, signer, actor):
        ev = create_bootstrap(actor, signer)
        store.append(ev)
        stored = store.append(ev)
        assert stored.id == ev.id
        assert store.count() == 1

    def test_chain_integrity(self, store, signer, actor, conv):
        ev = bootstrap(store, signer, actor)
        ev2 = chained_event(store, signer, actor, conv, [ev.id])
        store.append(ev2)
        assert store.count() == 2

    def test_head_empty(self, store):
        head = store.head()
        assert not head.is_some()

    def test_head(self, store, signer, actor):
        ev = bootstrap(store, signer, actor)
        head = store.head()
        assert head.is_some()
        assert head.unwrap().id == ev.id

    def test_count(self, store, signer, actor):
        assert store.count() == 0
        bootstrap(store, signer, actor)
        assert store.count() == 1


class TestSQLiteStoreQueries:
    def test_recent(self, store, signer, actor, conv):
        ev = bootstrap(store, signer, actor)
        ev2 = chained_event(store, signer, actor, conv, [ev.id])
        store.append(ev2)
        recent = store.recent(10)
        assert len(recent) == 2
        assert recent[0].id == ev2.id

    def test_by_type(self, store, signer, actor, conv):
        ev = bootstrap(store, signer, actor)
        ev2 = chained_event(store, signer, actor, conv, [ev.id])
        store.append(ev2)
        results = store.by_type(EventType("trust.updated"), 10)
        assert len(results) == 1
        assert results[0].id == ev2.id

    def test_by_source(self, store, signer, actor):
        ev = bootstrap(store, signer, actor)
        results = store.by_source(actor, 10)
        assert len(results) == 1

    def test_by_conversation(self, store, signer, actor, conv):
        ev = bootstrap(store, signer, actor)
        ev2 = chained_event(store, signer, actor, conv, [ev.id])
        store.append(ev2)
        results = store.by_conversation(conv, 10)
        assert len(results) == 1


class TestSQLiteStoreCausality:
    def test_ancestors(self, store, signer, actor, conv):
        ev0 = bootstrap(store, signer, actor)
        ev1 = chained_event(store, signer, actor, conv, [ev0.id])
        store.append(ev1)
        ev2 = chained_event(store, signer, actor, conv, [ev1.id])
        store.append(ev2)
        ancestors = store.ancestors(ev2.id, 10)
        assert len(ancestors) >= 1

    def test_ancestors_not_found(self, store):
        with pytest.raises(Exception):
            store.ancestors(EventID("019462a0-0000-7000-8000-000000000099"), 10)

    def test_descendants(self, store, signer, actor, conv):
        ev0 = bootstrap(store, signer, actor)
        ev1 = chained_event(store, signer, actor, conv, [ev0.id])
        store.append(ev1)
        descendants = store.descendants(ev0.id, 10)
        assert len(descendants) == 1

    def test_descendants_not_found(self, store):
        with pytest.raises(Exception):
            store.descendants(EventID("019462a0-0000-7000-8000-000000000099"), 10)


class TestSQLiteStoreChain:
    def test_verify_chain(self, store, signer, actor, conv):
        ev = bootstrap(store, signer, actor)
        ev2 = chained_event(store, signer, actor, conv, [ev.id])
        store.append(ev2)
        result = store.verify_chain()
        assert result.valid
        assert result.length == 2

    def test_verify_chain_empty(self, store):
        result = store.verify_chain()
        assert result.valid
        assert result.length == 0

    def test_close(self, store):
        store.close()
