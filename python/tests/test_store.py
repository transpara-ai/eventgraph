"""Tests for InMemoryStore — append, get, head, chain verification."""

import pytest

from eventgraph.errors import ChainIntegrityError, EventNotFoundError
from eventgraph.event import (
    NoopSigner,
    create_bootstrap,
    create_event,
)
from eventgraph.store import InMemoryStore
from eventgraph.types import (
    ActorID,
    ConversationID,
    EventID,
    EventType,
)


def _bootstrap():
    return create_bootstrap(source=ActorID("alice"), signer=NoopSigner())


def _next_event(prev_event):
    return create_event(
        event_type=EventType("trust.updated"),
        source=ActorID("alice"),
        content={"score": 0.5},
        causes=[prev_event.id],
        conversation_id=ConversationID("conv_1"),
        prev_hash=prev_event.hash,
        signer=NoopSigner(),
    )


class TestInMemoryStore:
    def test_append_and_get(self):
        store = InMemoryStore()
        boot = _bootstrap()
        store.append(boot)
        retrieved = store.get(boot.id)
        assert retrieved.id.value == boot.id.value

    def test_head_empty(self):
        store = InMemoryStore()
        assert store.head().is_none()

    def test_head_after_append(self):
        store = InMemoryStore()
        boot = _bootstrap()
        store.append(boot)
        head = store.head().unwrap()
        assert head.id.value == boot.id.value

    def test_count(self):
        store = InMemoryStore()
        assert store.count() == 0
        store.append(_bootstrap())
        assert store.count() == 1

    def test_chain_of_events(self):
        store = InMemoryStore()
        boot = _bootstrap()
        store.append(boot)

        e1 = _next_event(boot)
        store.append(e1)

        e2 = _next_event(e1)
        store.append(e2)

        assert store.count() == 3
        assert store.head().unwrap().id.value == e2.id.value

    def test_rejects_broken_chain(self):
        store = InMemoryStore()
        boot = _bootstrap()
        store.append(boot)

        # Create an event with wrong prev_hash
        bad = create_event(
            event_type=EventType("trust.updated"),
            source=ActorID("alice"),
            content={},
            causes=[boot.id],
            conversation_id=ConversationID("conv_1"),
            prev_hash=boot.prev_hash,  # wrong — should be boot.hash
            signer=NoopSigner(),
        )
        with pytest.raises(ChainIntegrityError):
            store.append(bad)

    def test_get_nonexistent(self):
        store = InMemoryStore()
        with pytest.raises(EventNotFoundError):
            store.get(EventID("019462a0-0000-7000-8000-000000000099"))

    def test_verify_chain_empty(self):
        store = InMemoryStore()
        v = store.verify_chain()
        assert v.valid is True
        assert v.length == 0

    def test_verify_chain_valid(self):
        store = InMemoryStore()
        boot = _bootstrap()
        store.append(boot)
        e1 = _next_event(boot)
        store.append(e1)

        v = store.verify_chain()
        assert v.valid is True
        assert v.length == 2

    def test_recent(self):
        store = InMemoryStore()
        boot = _bootstrap()
        store.append(boot)
        e1 = _next_event(boot)
        store.append(e1)
        e2 = _next_event(e1)
        store.append(e2)

        recent = store.recent(2)
        assert len(recent) == 2
        # Newest first
        assert recent[0].id.value == e2.id.value
        assert recent[1].id.value == e1.id.value

    def test_close_is_noop(self):
        store = InMemoryStore()
        store.close()  # should not raise
