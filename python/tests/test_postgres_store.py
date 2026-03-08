"""Tests for PostgresStore — verifies it satisfies the Store protocol."""

import os
import pytest

try:
    import psycopg2
    HAS_PSYCOPG2 = True
except ImportError:
    HAS_PSYCOPG2 = False

POSTGRES_URL = os.environ.get("POSTGRES_URL", "")

skip_condition = pytest.mark.skipif(
    not HAS_PSYCOPG2 or not POSTGRES_URL,
    reason="psycopg2 not installed or POSTGRES_URL env var not set",
)

from eventgraph.event import create_bootstrap, create_event, NoopSigner
from eventgraph.types import (
    ActorID, ConversationID, EventID, EventType, Hash,
)


@pytest.fixture
def store():
    if not HAS_PSYCOPG2 or not POSTGRES_URL:
        pytest.skip("psycopg2 not installed or POSTGRES_URL not set")

    from eventgraph.postgres_store import PostgresStore

    conn = psycopg2.connect(POSTGRES_URL)
    conn.autocommit = True
    with conn.cursor() as cur:
        cur.execute("DROP TABLE IF EXISTS events")
    conn.close()

    s = PostgresStore(POSTGRES_URL)
    yield s
    s.close()

    # Teardown: drop the table
    conn = psycopg2.connect(POSTGRES_URL)
    conn.autocommit = True
    with conn.cursor() as cur:
        cur.execute("DROP TABLE IF EXISTS events")
    conn.close()


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
    prev_hash = head.unwrap().hash if head.is_some() else Hash.zero()
    return create_event(
        EventType("trust.updated"), actor,
        {"actor": "test", "score": 0.5},
        causes, conv, prev_hash, signer,
    )


@skip_condition
class TestPostgresStoreBasic:
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


@skip_condition
class TestPostgresStoreQueries:
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


@skip_condition
class TestPostgresStoreCausality:
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


@skip_condition
class TestPostgresStoreChain:
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
