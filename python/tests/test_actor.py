"""Tests for the actor module — ports Go actor_test.go to pytest."""

import pytest

from eventgraph.actor import (
    Actor,
    ActorFilter,
    ActorStatus,
    ActorStore,
    ActorType,
    ActorUpdate,
    InMemoryActorStore,
)
from eventgraph.errors import (
    ActorKeyNotFoundError,
    ActorNotFoundError,
    InvalidTransitionError,
)
from eventgraph.types import ActorID, EventID, Option, PublicKey


def _test_public_key(b: int) -> PublicKey:
    """Create a test public key with the given first byte."""
    key = bytes([b]) + bytes(31)
    return PublicKey(key)


REASON = EventID("019462a0-0000-7000-8000-000000000001")


# ── Registration ─────────────────────────────────────────────────────────


class TestRegister:
    def test_register(self):
        s = InMemoryActorStore()
        pk = _test_public_key(1)
        a = s.register(pk, "Alice", ActorType.HUMAN)

        assert a.display_name == "Alice"
        assert a.actor_type == ActorType.HUMAN
        assert a.status == ActorStatus.ACTIVE

    def test_register_idempotent(self):
        s = InMemoryActorStore()
        pk = _test_public_key(1)

        a1 = s.register(pk, "Alice", ActorType.HUMAN)
        a2 = s.register(pk, "Alice Again", ActorType.HUMAN)

        assert a1.id.value == a2.id.value
        assert s.actor_count() == 1


# ── Get ──────────────────────────────────────────────────────────────────


class TestGet:
    def test_get(self):
        s = InMemoryActorStore()
        pk = _test_public_key(1)
        a = s.register(pk, "Alice", ActorType.HUMAN)

        got = s.get(a.id)
        assert got.id.value == a.id.value

    def test_get_not_found(self):
        s = InMemoryActorStore()
        with pytest.raises(ActorNotFoundError):
            s.get(ActorID("actor_nonexistent"))


# ── GetByPublicKey ───────────────────────────────────────────────────────


class TestGetByPublicKey:
    def test_get_by_public_key(self):
        s = InMemoryActorStore()
        pk = _test_public_key(1)
        a = s.register(pk, "Alice", ActorType.HUMAN)

        got = s.get_by_public_key(pk)
        assert got.id.value == a.id.value

    def test_get_by_public_key_not_found(self):
        s = InMemoryActorStore()
        with pytest.raises(ActorKeyNotFoundError):
            s.get_by_public_key(_test_public_key(99))


# ── Update ───────────────────────────────────────────────────────────────


class TestUpdate:
    def test_update(self):
        s = InMemoryActorStore()
        pk = _test_public_key(1)
        a = s.register(pk, "Alice", ActorType.HUMAN)

        updated = s.update(a.id, ActorUpdate(display_name=Option.some("Alice Updated")))
        assert updated.display_name == "Alice Updated"

    def test_update_metadata_merge(self):
        s = InMemoryActorStore()
        pk = _test_public_key(1)
        a = s.register(pk, "Alice", ActorType.HUMAN)

        # Set initial metadata
        s.update(a.id, ActorUpdate(metadata=Option.some({"role": "builder"})))

        # Merge additional metadata
        updated = s.update(a.id, ActorUpdate(metadata=Option.some({"team": "core"})))
        md = updated.metadata
        assert md["role"] == "builder"
        assert md["team"] == "core"

    def test_update_not_found(self):
        s = InMemoryActorStore()
        with pytest.raises(ActorNotFoundError):
            s.update(ActorID("actor_nonexistent"), ActorUpdate())


# ── Suspend ──────────────────────────────────────────────────────────────


class TestSuspend:
    def test_suspend(self):
        s = InMemoryActorStore()
        pk = _test_public_key(1)
        a = s.register(pk, "Alice", ActorType.HUMAN)

        suspended = s.suspend(a.id, REASON)
        assert suspended.status == ActorStatus.SUSPENDED

    def test_suspend_not_found(self):
        s = InMemoryActorStore()
        with pytest.raises(ActorNotFoundError):
            s.suspend(ActorID("actor_nonexistent"), REASON)

    def test_suspend_and_reactivate(self):
        s = InMemoryActorStore()
        pk = _test_public_key(1)
        a = s.register(pk, "Alice", ActorType.HUMAN)

        s.suspend(a.id, REASON)
        got = s.get(a.id)
        assert got.status == ActorStatus.SUSPENDED

        reactivated = s.reactivate(a.id, REASON)
        assert reactivated.status == ActorStatus.ACTIVE


# ── Memorial ─────────────────────────────────────────────────────────────


class TestMemorial:
    def test_memorial(self):
        s = InMemoryActorStore()
        pk = _test_public_key(1)
        a = s.register(pk, "Alice", ActorType.HUMAN)

        memorial = s.memorial(a.id, REASON)
        assert memorial.status == ActorStatus.MEMORIAL

    def test_memorial_is_terminal(self):
        s = InMemoryActorStore()
        pk = _test_public_key(1)
        a = s.register(pk, "Alice", ActorType.HUMAN)

        s.memorial(a.id, REASON)

        # Try to suspend — should fail
        with pytest.raises(InvalidTransitionError):
            s.suspend(a.id, REASON)

    def test_memorial_reactivate_is_error(self):
        s = InMemoryActorStore()
        pk = _test_public_key(1)
        a = s.register(pk, "Alice", ActorType.HUMAN)

        s.memorial(a.id, REASON)

        with pytest.raises(InvalidTransitionError):
            s.reactivate(a.id, REASON)


# ── Reactivate ───────────────────────────────────────────────────────────


class TestReactivate:
    def test_reactivate_from_active_is_error(self):
        s = InMemoryActorStore()
        pk = _test_public_key(1)
        a = s.register(pk, "Alice", ActorType.HUMAN)

        with pytest.raises(InvalidTransitionError):
            s.reactivate(a.id, REASON)

    def test_reactivate_not_found(self):
        s = InMemoryActorStore()
        with pytest.raises(ActorNotFoundError):
            s.reactivate(ActorID("actor_nonexistent"), REASON)


# ── List ─────────────────────────────────────────────────────────────────


class TestList:
    def test_list(self):
        s = InMemoryActorStore()
        for i in range(1, 6):
            s.register(_test_public_key(i), "Actor", ActorType.HUMAN)

        page = s.list(ActorFilter(limit=10))
        assert len(page.items()) == 5

    def test_list_with_status_filter(self):
        s = InMemoryActorStore()
        s.register(_test_public_key(1), "Active1", ActorType.HUMAN)
        s.register(_test_public_key(2), "Active2", ActorType.HUMAN)
        a3 = s.register(_test_public_key(3), "ToBeSuspended", ActorType.HUMAN)

        s.suspend(a3.id, REASON)

        active_page = s.list(ActorFilter(
            status=Option.some(ActorStatus.ACTIVE),
            limit=10,
        ))
        assert len(active_page.items()) == 2

    def test_list_with_type_filter(self):
        s = InMemoryActorStore()
        s.register(_test_public_key(1), "Human", ActorType.HUMAN)
        s.register(_test_public_key(2), "AI", ActorType.AI)
        s.register(_test_public_key(3), "System", ActorType.SYSTEM)

        page = s.list(ActorFilter(
            actor_type=Option.some(ActorType.AI),
            limit=10,
        ))
        assert len(page.items()) == 1

    def test_list_pagination(self):
        s = InMemoryActorStore()
        for i in range(1, 6):
            s.register(_test_public_key(i), "Actor", ActorType.HUMAN)

        # Page 1
        page1 = s.list(ActorFilter(limit=2))
        assert len(page1.items()) == 2
        assert page1.has_more()

        # Page 2
        cursor = page1.cursor()
        assert cursor.is_some()
        page2 = s.list(ActorFilter(limit=2, after=Option.some(cursor.unwrap().value)))
        assert len(page2.items()) == 2


# ── Actor Getters ────────────────────────────────────────────────────────


class TestActorGetters:
    def test_actor_getters(self):
        s = InMemoryActorStore()
        pk = _test_public_key(42)
        a = s.register(pk, "Alice", ActorType.AI)

        assert a.id.value != ""
        assert a.public_key is not None
        assert a.created_at > 0
        assert isinstance(a.metadata, dict)


# ── ActorStatus transitions ─────────────────────────────────────────────


class TestActorStatusTransitions:
    def test_active_to_suspended(self):
        result = ActorStatus.ACTIVE.transition_to(ActorStatus.SUSPENDED)
        assert result == ActorStatus.SUSPENDED

    def test_active_to_memorial(self):
        result = ActorStatus.ACTIVE.transition_to(ActorStatus.MEMORIAL)
        assert result == ActorStatus.MEMORIAL

    def test_suspended_to_active(self):
        result = ActorStatus.SUSPENDED.transition_to(ActorStatus.ACTIVE)
        assert result == ActorStatus.ACTIVE

    def test_suspended_to_memorial(self):
        result = ActorStatus.SUSPENDED.transition_to(ActorStatus.MEMORIAL)
        assert result == ActorStatus.MEMORIAL

    def test_memorial_is_terminal(self):
        with pytest.raises(InvalidTransitionError):
            ActorStatus.MEMORIAL.transition_to(ActorStatus.ACTIVE)
        with pytest.raises(InvalidTransitionError):
            ActorStatus.MEMORIAL.transition_to(ActorStatus.SUSPENDED)

    def test_valid_transitions(self):
        assert set(ActorStatus.ACTIVE.valid_transitions()) == {
            ActorStatus.SUSPENDED,
            ActorStatus.MEMORIAL,
        }
        assert set(ActorStatus.SUSPENDED.valid_transitions()) == {
            ActorStatus.ACTIVE,
            ActorStatus.MEMORIAL,
        }
        assert ActorStatus.MEMORIAL.valid_transitions() == []


# ── Protocol conformance ────────────────────────────────────────────────


class TestProtocol:
    def test_in_memory_store_implements_protocol(self):
        s = InMemoryActorStore()
        assert isinstance(s, ActorStore)
