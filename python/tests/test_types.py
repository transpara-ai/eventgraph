"""Tests for constrained types and typed IDs."""

import pytest

from eventgraph.errors import EmptyRequiredError, InvalidFormatError, OutOfRangeError
from eventgraph.types import (
    Activation,
    ActorID,
    Cadence,
    ConversationID,
    DomainScope,
    EdgeID,
    EnvelopeID,
    EventID,
    EventType,
    Hash,
    Layer,
    NonEmpty,
    Option,
    PublicKey,
    Score,
    Signature,
    SubscriptionPattern,
    SystemURI,
    TreatyID,
    Weight,
)


# ── Score ─────────────────────────────────────────────────────────────────

class TestScore:
    def test_valid_range(self):
        assert Score(0.0).value == 0.0
        assert Score(0.5).value == 0.5
        assert Score(1.0).value == 1.0

    def test_below_range(self):
        with pytest.raises(OutOfRangeError):
            Score(-0.1)

    def test_above_range(self):
        with pytest.raises(OutOfRangeError):
            Score(1.1)

    def test_nan(self):
        with pytest.raises(OutOfRangeError):
            Score(float("nan"))

    def test_frozen(self):
        s = Score(0.5)
        with pytest.raises(AttributeError):
            s._value = 0.9  # type: ignore


# ── Weight ────────────────────────────────────────────────────────────────

class TestWeight:
    def test_valid_range(self):
        assert Weight(-1.0).value == -1.0
        assert Weight(0.0).value == 0.0
        assert Weight(1.0).value == 1.0

    def test_below_range(self):
        with pytest.raises(OutOfRangeError):
            Weight(-1.1)

    def test_above_range(self):
        with pytest.raises(OutOfRangeError):
            Weight(1.1)


# ── Activation ────────────────────────────────────────────────────────────

class TestActivation:
    def test_valid_range(self):
        assert Activation(0.0).value == 0.0
        assert Activation(1.0).value == 1.0

    def test_below_range(self):
        with pytest.raises(OutOfRangeError):
            Activation(-0.01)


# ── Layer ─────────────────────────────────────────────────────────────────

class TestLayer:
    def test_valid_range(self):
        assert Layer(0).value == 0
        assert Layer(13).value == 13

    def test_below_range(self):
        with pytest.raises(OutOfRangeError):
            Layer(-1)

    def test_above_range(self):
        with pytest.raises(OutOfRangeError):
            Layer(14)


# ── Cadence ───────────────────────────────────────────────────────────────

class TestCadence:
    def test_valid(self):
        assert Cadence(1).value == 1
        assert Cadence(100).value == 100

    def test_zero(self):
        with pytest.raises(OutOfRangeError):
            Cadence(0)


# ── EventID ───────────────────────────────────────────────────────────────

class TestEventID:
    def test_valid_uuid_v7(self):
        eid = EventID("019462a0-0000-7000-8000-000000000001")
        assert eid.value == "019462a0-0000-7000-8000-000000000001"

    def test_normalizes_to_lowercase(self):
        eid = EventID("019462A0-0000-7000-8000-000000000001")
        assert eid.value == "019462a0-0000-7000-8000-000000000001"

    def test_rejects_non_v7(self):
        with pytest.raises(InvalidFormatError):
            EventID("019462a0-0000-4000-8000-000000000001")  # v4 not v7

    def test_rejects_garbage(self):
        with pytest.raises(InvalidFormatError):
            EventID("not-a-uuid")


# ── Hash ──────────────────────────────────────────────────────────────────

class TestHash:
    def test_valid_64_hex(self):
        h = Hash("a" * 64)
        assert h.value == "a" * 64

    def test_zero_hash(self):
        h = Hash.zero()
        assert h.is_zero()
        assert h.value == "0" * 64

    def test_rejects_short(self):
        with pytest.raises(InvalidFormatError):
            Hash("abc")

    def test_rejects_empty(self):
        with pytest.raises(InvalidFormatError):
            Hash("")


# ── EventType ─────────────────────────────────────────────────────────────

class TestEventType:
    def test_valid(self):
        et = EventType("trust.updated")
        assert et.value == "trust.updated"

    def test_rejects_uppercase(self):
        with pytest.raises(InvalidFormatError):
            EventType("Trust.Updated")

    def test_rejects_empty(self):
        with pytest.raises(InvalidFormatError):
            EventType("")


# ── ActorID ───────────────────────────────────────────────────────────────

class TestActorID:
    def test_valid(self):
        a = ActorID("actor_alice")
        assert a.value == "actor_alice"

    def test_rejects_empty(self):
        with pytest.raises(EmptyRequiredError):
            ActorID("")


# ── SubscriptionPattern ──────────────────────────────────────────────────

class TestSubscriptionPattern:
    def test_wildcard_matches_all(self):
        sp = SubscriptionPattern("*")
        assert sp.matches(EventType("trust.updated"))
        assert sp.matches(EventType("system.bootstrapped"))

    def test_prefix_match(self):
        sp = SubscriptionPattern("trust.*")
        assert sp.matches(EventType("trust.updated"))
        assert not sp.matches(EventType("system.bootstrapped"))

    def test_exact_match(self):
        sp = SubscriptionPattern("trust.updated")
        assert sp.matches(EventType("trust.updated"))
        assert not sp.matches(EventType("trust.decayed"))


# ── Option ────────────────────────────────────────────────────────────────

class TestOption:
    def test_some(self):
        opt = Option.some(42)
        assert opt.is_some()
        assert not opt.is_none()
        assert opt.unwrap() == 42

    def test_none(self):
        opt: Option[int] = Option.none()
        assert opt.is_none()
        assert not opt.is_some()
        with pytest.raises(ValueError):
            opt.unwrap()

    def test_unwrap_or(self):
        assert Option.some(42).unwrap_or(0) == 42
        assert Option.none().unwrap_or(0) == 0


# ── NonEmpty ──────────────────────────────────────────────────────────────

class TestNonEmpty:
    def test_valid(self):
        ne = NonEmpty.of([1, 2, 3])
        assert len(ne) == 3
        assert ne[0] == 1

    def test_rejects_empty(self):
        with pytest.raises(ValueError):
            NonEmpty.of([])

    def test_iterable(self):
        ne = NonEmpty.of([1, 2])
        assert list(ne) == [1, 2]


# ── PublicKey / Signature ─────────────────────────────────────────────────

class TestPublicKey:
    def test_valid_32_bytes(self):
        pk = PublicKey(b"\x00" * 32)
        assert len(pk.bytes_) == 32

    def test_rejects_wrong_length(self):
        with pytest.raises(InvalidFormatError):
            PublicKey(b"\x00" * 31)


class TestSignature:
    def test_valid_64_bytes(self):
        sig = Signature(b"\x00" * 64)
        assert len(sig.bytes_) == 64

    def test_rejects_wrong_length(self):
        with pytest.raises(InvalidFormatError):
            Signature(b"\x00" * 63)
