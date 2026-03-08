"""Tests for Event creation, canonical form, and hash chain."""

import pytest

from eventgraph.event import (
    Event,
    NoopSigner,
    canonical_content_json,
    canonical_form,
    compute_hash,
    create_bootstrap,
    create_event,
    new_event_id,
)
from eventgraph.types import (
    ActorID,
    ConversationID,
    EventType,
    Hash,
    NonEmpty,
)


class TestCanonicalContentJson:
    def test_sorted_keys(self):
        result = canonical_content_json({"b": 1, "a": 2})
        assert result == '{"a":2,"b":1}'

    def test_no_whitespace(self):
        result = canonical_content_json({"key": "value"})
        assert " " not in result

    def test_nested_sorted(self):
        result = canonical_content_json({"z": {"b": 1, "a": 2}, "a": 0})
        assert result == '{"a":0,"z":{"a":2,"b":1}}'


class TestCanonicalForm:
    def test_format(self):
        canon = canonical_form(
            version=1,
            prev_hash="0" * 64,
            causes=["c2", "c1"],
            event_id="eid",
            event_type="trust.updated",
            source="actor_alice",
            conversation_id="conv_1",
            timestamp_nanos=123456789,
            content_json='{"key":"val"}',
        )
        # Causes should be sorted
        assert "c1,c2" in canon
        # Pipe-separated
        parts = canon.split("|")
        assert parts[0] == "1"  # version
        assert parts[1] == "0" * 64  # prev_hash
        assert parts[2] == "c1,c2"  # sorted causes
        assert parts[3] == "eid"  # event_id
        assert parts[4] == "trust.updated"  # event_type
        assert parts[5] == "actor_alice"  # source
        assert parts[6] == "conv_1"  # conversation_id
        assert parts[7] == "123456789"  # timestamp_nanos
        assert parts[8] == '{"key":"val"}'  # content_json

    def test_empty_causes(self):
        canon = canonical_form(
            version=1, prev_hash="", causes=[], event_id="eid",
            event_type="system.bootstrapped", source="s", conversation_id="c",
            timestamp_nanos=0, content_json="{}",
        )
        parts = canon.split("|")
        assert parts[2] == ""  # empty causes


class TestComputeHash:
    def test_deterministic(self):
        h1 = compute_hash("hello")
        h2 = compute_hash("hello")
        assert h1.value == h2.value

    def test_different_input(self):
        h1 = compute_hash("hello")
        h2 = compute_hash("world")
        assert h1.value != h2.value

    def test_returns_valid_hash(self):
        h = compute_hash("test")
        assert len(h.value) == 64


class TestNewEventId:
    def test_generates_valid_uuid_v7(self):
        eid = new_event_id()
        # Should not raise — already validated by EventID constructor
        assert len(eid.value) == 36
        # Version nibble is '7'
        assert eid.value[14] == "7"


class TestConformanceVectors:
    """Tests from docs/conformance/canonical-vectors.json."""

    def test_bootstrap_conformance_hash(self):
        content_json = canonical_content_json({
            "ActorID": "actor_00000000000000000000000000000001",
            "ChainGenesis": "0000000000000000000000000000000000000000000000000000000000000000",
            "Timestamp": "2023-11-14T22:13:20Z",
        })
        canon = canonical_form(
            version=1, prev_hash="", causes=[],
            event_id="019462a0-0000-7000-8000-000000000001",
            event_type="system.bootstrapped",
            source="actor_00000000000000000000000000000001",
            conversation_id="conv_00000000000000000000000000000001",
            timestamp_nanos=1700000000000000000,
            content_json=content_json,
        )
        h = compute_hash(canon)
        assert h.value == "f7cae7ae11c1232a932c64f2302432c0e304dffce80f3935e688980dfbafeb75"

    def test_trust_updated_conformance_hash(self):
        content_json = canonical_content_json({
            "Actor": "actor_00000000000000000000000000000002",
            "Cause": "019462a0-0000-7000-8000-000000000001",
            "Current": 0.85,
            "Domain": "code_review",
            "Previous": 0.8,
        })
        canon = canonical_form(
            version=1,
            prev_hash="a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
            causes=["019462a0-0000-7000-8000-000000000001"],
            event_id="019462a0-0000-7000-8000-000000000002",
            event_type="trust.updated",
            source="actor_00000000000000000000000000000001",
            conversation_id="conv_00000000000000000000000000000001",
            timestamp_nanos=1700000001000000000,
            content_json=content_json,
        )
        h = compute_hash(canon)
        assert h.value == "b2fbcd2684868f0b0d07d2f5136b52f14b8e749da7b4b7bae2a22f67147152b7"

    def test_content_key_ordering_conformance_hash(self):
        content_json = canonical_content_json({
            "Weight": 0.5,
            "From": "actor_00000000000000000000000000000001",
            "To": "actor_00000000000000000000000000000002",
            "EdgeType": "Trust",
            "Direction": "Centripetal",
        })
        canon = canonical_form(
            version=1,
            prev_hash="b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3",
            causes=["019462a0-0000-7000-8000-000000000001"],
            event_id="019462a0-0000-7000-8000-000000000003",
            event_type="edge.created",
            source="actor_00000000000000000000000000000001",
            conversation_id="conv_00000000000000000000000000000001",
            timestamp_nanos=1700000002000000000,
            content_json=content_json,
        )
        h = compute_hash(canon)
        assert h.value == "4e5c6710ca9325676663b4a66d2e82114fcd8fb49dbe5705795051e0b0be374c"

    def test_multiple_causes_sorted_conformance(self):
        content_json = canonical_content_json({"Message": "derived from multiple causes"})
        canon = canonical_form(
            version=1,
            prev_hash="c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4",
            causes=[
                "019462a0-0000-7000-8000-000000000003",
                "019462a0-0000-7000-8000-000000000001",
                "019462a0-0000-7000-8000-000000000002",
            ],
            event_id="019462a0-0000-7000-8000-000000000004",
            event_type="grammar.derived",
            source="actor_00000000000000000000000000000001",
            conversation_id="conv_00000000000000000000000000000001",
            timestamp_nanos=1700000003000000000,
            content_json=content_json,
        )
        h = compute_hash(canon)
        assert h.value == "0c0e47ee89f8a7a21bb47f60d5f3887833297c945f946c8f3695ff2638f6cd50"

    def test_integer_float_formatting_conformance(self):
        content_json = canonical_content_json({
            "Current": 1.0,
            "Domain": "testing",
            "Previous": 0.5,
        })
        assert content_json == '{"Current":1,"Domain":"testing","Previous":0.5}'

    def test_null_omission_conformance(self):
        content_json = canonical_content_json({
            "Actor": "actor_00000000000000000000000000000001",
            "Scope": None,
            "Reason": None,
        })
        assert content_json == '{"Actor":"actor_00000000000000000000000000000001"}'

    def test_nested_objects_conformance(self):
        content_json = canonical_content_json({
            "Outer": "value",
            "Nested": {
                "Zebra": 1,
                "Alpha": "first",
                "Middle": None,
            },
        })
        assert content_json == '{"Nested":{"Alpha":"first","Zebra":1},"Outer":"value"}'

    def test_number_formatting_rules(self):
        """Integer-like floats become integers, non-integer floats keep decimals."""
        assert canonical_content_json({"v": 1.0}) == '{"v":1}'
        assert canonical_content_json({"v": 0.5}) == '{"v":0.5}'
        assert canonical_content_json({"v": 0.1}) == '{"v":0.1}'
        assert canonical_content_json({"v": 100}) == '{"v":100}'
        assert canonical_content_json({"v": 0.001}) == '{"v":0.001}'
        assert canonical_content_json({"v": 2.0}) == '{"v":2}'
        assert canonical_content_json({"v": -1.0}) == '{"v":-1}'
        assert canonical_content_json({"v": 0.0}) == '{"v":0}'


class TestCreateBootstrap:
    def test_creates_valid_bootstrap(self):
        signer = NoopSigner()
        source = ActorID("actor_alice")
        event = create_bootstrap(source=source, signer=signer)

        assert event.version == 1
        assert event.type.value == "system.bootstrapped"
        assert event.source.value == "actor_alice"
        assert event.prev_hash.is_zero()
        assert len(event.causes) == 1  # self-referencing
        assert event.causes[0].value == event.id.value
        assert len(event.hash.value) == 64
        assert event.signature.bytes_ == b"\x00" * 64

    def test_bootstrap_hash_is_deterministic_for_same_canonical(self):
        # Two bootstraps will differ (different timestamps/IDs),
        # but hash should match canonical form
        signer = NoopSigner()
        event = create_bootstrap(source=ActorID("alice"), signer=signer)
        content_json = canonical_content_json(event.content)
        canon = canonical_form(
            version=event.version,
            prev_hash="",
            causes=[],
            event_id=event.id.value,
            event_type=event.type.value,
            source=event.source.value,
            conversation_id=event.conversation_id.value,
            timestamp_nanos=event.timestamp_nanos,
            content_json=content_json,
        )
        assert event.hash == compute_hash(canon)


class TestCreateEvent:
    def test_creates_valid_event(self):
        signer = NoopSigner()
        bootstrap = create_bootstrap(source=ActorID("alice"), signer=signer)

        event = create_event(
            event_type=EventType("trust.updated"),
            source=ActorID("alice"),
            content={"score": 0.8},
            causes=[bootstrap.id],
            conversation_id=ConversationID("conv_1"),
            prev_hash=bootstrap.hash,
            signer=signer,
        )

        assert event.version == 1
        assert event.type.value == "trust.updated"
        assert event.source.value == "alice"
        assert event.prev_hash == bootstrap.hash
        assert len(event.causes) == 1
        assert event.causes[0].value == bootstrap.id.value
        assert event.content == {"score": 0.8}

    def test_event_is_frozen(self):
        signer = NoopSigner()
        event = create_bootstrap(source=ActorID("alice"), signer=signer)
        with pytest.raises(AttributeError):
            event._version = 2  # type: ignore

    def test_content_is_defensive_copy(self):
        signer = NoopSigner()
        event = create_bootstrap(source=ActorID("alice"), signer=signer)
        content = event.content
        content["injected"] = True
        assert "injected" not in event.content
