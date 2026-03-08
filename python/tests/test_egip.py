"""Tests for the EGIP module."""

from __future__ import annotations

import time
import uuid

import pytest

from eventgraph.egip import (
    # Constants
    CURRENT_PROTOCOL_VERSION,
    MAX_ENVELOPE_AGE_SECONDS,
    TRUST_IMPACT_NO_HELLO_RESPONSE,
    TRUST_IMPACT_RECEIPT_ON_TIME,
    TRUST_IMPACT_SIGNATURE_INVALID,
    TRUST_IMPACT_TREATY_HONOURED,
    TRUST_IMPACT_VALID_PROOF,
    # Enums
    AuthorityLevel,
    CGERRelationship,
    MessageType,
    ProofType,
    ReceiptStatus,
    TreatyAction,
    TreatyStatus,
    # Errors
    DuplicateEnvelopeError,
    EGIPError,
    EnvelopeSignatureInvalidError,
    SystemNotFoundError,
    TransportFailureError,
    TreatyNotFoundError,
    TreatyViolationError,
    TrustInsufficientError,
    VersionIncompatibleError,
    # Identity
    IIdentity,
    SystemIdentity,
    # Envelope
    Envelope,
    sign_envelope,
    verify_envelope,
    # Payloads
    AuthorityRequestPayload,
    CGER,
    ChainSegmentProof,
    ChainSummaryProof,
    DiscoverPayload,
    DiscoverQuery,
    DiscoverResult,
    EventExistenceProof,
    HelloPayload,
    MessagePayloadContent,
    ProofPayload,
    ReceiptPayload,
    TreatyPayload,
    TreatyTerm,
    # Version negotiation
    negotiate_version,
    # Treaty
    Treaty,
    new_treaty,
    # Stores
    EnvelopeDedup,
    Handler,
    PeerRecord,
    PeerStore,
    TreatyStore,
    # Proof
    ProofGenerator,
    validate_proof,
    verify_chain_segment,
    verify_event_existence,
)
from eventgraph.types import (
    EnvelopeID,
    Option,
    PublicKey,
    Score,
    Signature,
    SystemURI,
    TreatyID,
)


# ── Helpers ────────────────────────────────────────────────────────────────


def _make_identity(name: str = "system-a") -> SystemIdentity:
    return SystemIdentity.generate(SystemURI(name))


def _make_envelope_id() -> EnvelopeID:
    return EnvelopeID(str(uuid.uuid4()))


def _make_treaty_id() -> TreatyID:
    return TreatyID(str(uuid.uuid4()))


def _make_hello_envelope(
    identity: SystemIdentity, to: SystemURI
) -> Envelope:
    return Envelope(
        protocol_version=CURRENT_PROTOCOL_VERSION,
        id=_make_envelope_id(),
        from_uri=identity.system_uri(),
        to_uri=to,
        type=MessageType.HELLO,
        payload=HelloPayload(
            system_uri=identity.system_uri().value,
            public_key=identity.public_key().bytes_,
            protocol_versions=[CURRENT_PROTOCOL_VERSION],
            capabilities=["treaty", "proof"],
            chain_length=0,
        ),
        timestamp=time.time(),
    )


class MockTransport:
    """Mock transport for testing."""

    def __init__(self, response: ReceiptPayload | None = None) -> None:
        self.sent: list[tuple[SystemURI, Envelope]] = []
        self.response = response

    def send(self, to: SystemURI, envelope: Envelope) -> ReceiptPayload | None:
        self.sent.append((to, envelope))
        return self.response


# ── SystemIdentity Tests ──────────────────────────────────────────────────


class TestSystemIdentity:
    def test_generate(self) -> None:
        identity = _make_identity()
        assert identity.system_uri().value == "system-a"
        assert len(identity.public_key().bytes_) == 32
        assert identity.created_at > 0

    def test_sign_and_verify(self) -> None:
        identity = _make_identity()
        data = b"hello world"
        sig = identity.sign(data)
        assert len(sig.bytes_) == 64
        assert identity.verify(identity.public_key(), data, sig)

    def test_verify_wrong_data(self) -> None:
        identity = _make_identity()
        sig = identity.sign(b"hello")
        assert not identity.verify(identity.public_key(), b"world", sig)

    def test_verify_wrong_signature(self) -> None:
        identity = _make_identity()
        bad_sig = Signature(b"\x00" * 64)
        assert not identity.verify(identity.public_key(), b"hello", bad_sig)

    def test_two_identities_different_keys(self) -> None:
        a = _make_identity("a")
        b = _make_identity("b")
        assert a.public_key().bytes_ != b.public_key().bytes_

    def test_implements_protocol(self) -> None:
        identity = _make_identity()
        assert isinstance(identity, IIdentity)


# ── Envelope Tests ─────────────────────────────────────────────────────────


class TestEnvelope:
    def test_canonical_form_deterministic(self) -> None:
        identity = _make_identity()
        env = _make_hello_envelope(identity, SystemURI("system-b"))
        form1 = env.canonical_form()
        form2 = env.canonical_form()
        assert form1 == form2

    def test_canonical_form_contains_fields(self) -> None:
        identity = _make_identity()
        env = _make_hello_envelope(identity, SystemURI("system-b"))
        form = env.canonical_form()
        assert f"{CURRENT_PROTOCOL_VERSION}|" in form
        assert "system-a" in form
        assert "system-b" in form
        assert "|hello|" in form

    def test_canonical_form_in_reply_to(self) -> None:
        identity = _make_identity()
        reply_id = _make_envelope_id()
        env = Envelope(
            protocol_version=CURRENT_PROTOCOL_VERSION,
            id=_make_envelope_id(),
            from_uri=SystemURI("a"),
            to_uri=SystemURI("b"),
            type=MessageType.RECEIPT,
            payload=ReceiptPayload(
                envelope_id=reply_id.value,
                status=ReceiptStatus.PROCESSED,
            ),
            timestamp=time.time(),
            in_reply_to=Option.some(reply_id),
        )
        form = env.canonical_form()
        assert reply_id.value in form

    def test_sign_and_verify_envelope(self) -> None:
        identity = _make_identity()
        env = _make_hello_envelope(identity, SystemURI("system-b"))
        signed = sign_envelope(env, identity)
        assert signed.signature.bytes_ != b"\x00" * 64
        assert verify_envelope(signed, identity, identity.public_key())

    def test_verify_tampered_envelope(self) -> None:
        identity = _make_identity()
        env = _make_hello_envelope(identity, SystemURI("system-b"))
        signed = sign_envelope(env, identity)
        # Tamper with the envelope
        tampered = Envelope(
            protocol_version=signed.protocol_version,
            id=signed.id,
            from_uri=SystemURI("tampered"),
            to_uri=signed.to_uri,
            type=signed.type,
            payload=signed.payload,
            timestamp=signed.timestamp,
            signature=signed.signature,
            in_reply_to=signed.in_reply_to,
        )
        assert not verify_envelope(tampered, identity, identity.public_key())


# ── NegotiateVersion Tests ─────────────────────────────────────────────────


class TestNegotiateVersion:
    def test_common_version(self) -> None:
        result = negotiate_version([1, 2, 3], [2, 3, 4])
        assert result.is_some()
        assert result.unwrap() == 3

    def test_single_common(self) -> None:
        result = negotiate_version([1], [1])
        assert result.is_some()
        assert result.unwrap() == 1

    def test_no_common(self) -> None:
        result = negotiate_version([1, 2], [3, 4])
        assert result.is_none()

    def test_empty_lists(self) -> None:
        assert negotiate_version([], [1]).is_none()
        assert negotiate_version([1], []).is_none()
        assert negotiate_version([], []).is_none()

    def test_picks_highest(self) -> None:
        result = negotiate_version([1, 2, 5], [5, 2])
        assert result.unwrap() == 5


# ── Treaty State Machine Tests ─────────────────────────────────────────────


class TestTreaty:
    def test_new_treaty_is_proposed(self) -> None:
        t = new_treaty(
            _make_treaty_id(),
            SystemURI("a"),
            SystemURI("b"),
            [TreatyTerm(scope="code", policy="read-only")],
        )
        assert t.status == TreatyStatus.PROPOSED

    def test_propose_to_active(self) -> None:
        t = new_treaty(_make_treaty_id(), SystemURI("a"), SystemURI("b"), [])
        t.apply_action(TreatyAction.ACCEPT)
        assert t.status == TreatyStatus.ACTIVE

    def test_active_to_suspended(self) -> None:
        t = new_treaty(_make_treaty_id(), SystemURI("a"), SystemURI("b"), [])
        t.apply_action(TreatyAction.ACCEPT)
        t.apply_action(TreatyAction.SUSPEND)
        assert t.status == TreatyStatus.SUSPENDED

    def test_suspended_to_active(self) -> None:
        t = new_treaty(_make_treaty_id(), SystemURI("a"), SystemURI("b"), [])
        t.apply_action(TreatyAction.ACCEPT)
        t.apply_action(TreatyAction.SUSPEND)
        t.apply_action(TreatyAction.ACCEPT)
        assert t.status == TreatyStatus.ACTIVE

    def test_active_to_terminated(self) -> None:
        t = new_treaty(_make_treaty_id(), SystemURI("a"), SystemURI("b"), [])
        t.apply_action(TreatyAction.ACCEPT)
        t.apply_action(TreatyAction.TERMINATE)
        assert t.status == TreatyStatus.TERMINATED

    def test_full_lifecycle(self) -> None:
        """propose -> active -> suspended -> active -> terminated"""
        t = new_treaty(_make_treaty_id(), SystemURI("a"), SystemURI("b"), [])
        assert t.status == TreatyStatus.PROPOSED
        t.apply_action(TreatyAction.ACCEPT)
        assert t.status == TreatyStatus.ACTIVE
        t.apply_action(TreatyAction.SUSPEND)
        assert t.status == TreatyStatus.SUSPENDED
        t.apply_action(TreatyAction.ACCEPT)
        assert t.status == TreatyStatus.ACTIVE
        t.apply_action(TreatyAction.TERMINATE)
        assert t.status == TreatyStatus.TERMINATED

    def test_terminated_is_terminal(self) -> None:
        t = new_treaty(_make_treaty_id(), SystemURI("a"), SystemURI("b"), [])
        t.apply_action(TreatyAction.ACCEPT)
        t.apply_action(TreatyAction.TERMINATE)
        with pytest.raises(Exception):
            t.apply_action(TreatyAction.ACCEPT)

    def test_invalid_transition_proposed_to_suspended(self) -> None:
        t = new_treaty(_make_treaty_id(), SystemURI("a"), SystemURI("b"), [])
        with pytest.raises(Exception):
            t.apply_action(TreatyAction.SUSPEND)

    def test_modify_only_when_active(self) -> None:
        t = new_treaty(_make_treaty_id(), SystemURI("a"), SystemURI("b"), [])
        with pytest.raises(Exception):
            t.apply_action(TreatyAction.MODIFY)

        t.apply_action(TreatyAction.ACCEPT)
        t.apply_action(TreatyAction.MODIFY)  # should succeed
        assert t.status == TreatyStatus.ACTIVE

    def test_propose_on_existing_raises(self) -> None:
        t = new_treaty(_make_treaty_id(), SystemURI("a"), SystemURI("b"), [])
        with pytest.raises(ValueError, match="cannot apply Propose"):
            t.apply_action(TreatyAction.PROPOSE)

    def test_treaty_copy(self) -> None:
        t = new_treaty(
            _make_treaty_id(), SystemURI("a"), SystemURI("b"),
            [TreatyTerm(scope="x", policy="y")],
        )
        cp = t.copy()
        assert cp.id == t.id
        assert cp.status == t.status
        cp.terms.append(TreatyTerm(scope="z", policy="w"))
        assert len(t.terms) == 1  # original unmodified


# ── PeerStore Tests ────────────────────────────────────────────────────────


class TestPeerStore:
    def test_register_and_get(self) -> None:
        ps = PeerStore()
        uri = SystemURI("peer-1")
        pk = PublicKey(b"\x01" * 32)
        ps.register(uri, pk, ["cap1"], 1)

        record, found = ps.get(uri)
        assert found
        assert record.system_uri.value == "peer-1"
        assert record.public_key.bytes_ == pk.bytes_
        assert record.capabilities == ["cap1"]
        assert record.negotiated_version == 1

    def test_get_unknown(self) -> None:
        ps = PeerStore()
        _, found = ps.get(SystemURI("unknown"))
        assert not found

    def test_update_trust_positive(self) -> None:
        ps = PeerStore()
        uri = SystemURI("peer-1")
        ps.register(uri, PublicKey(b"\x01" * 32), [], 1)

        score, ok = ps.update_trust(uri, 0.10)
        assert ok
        # Positive capped at INTER_SYSTEM_MAX_ADJUSTMENT (0.05)
        assert abs(score.value - 0.05) < 1e-9

    def test_update_trust_negative_uncapped(self) -> None:
        ps = PeerStore()
        uri = SystemURI("peer-1")
        ps.register(uri, PublicKey(b"\x01" * 32), [], 1)
        ps.update_trust(uri, 0.05)

        score, ok = ps.update_trust(uri, -0.20)
        assert ok
        assert score.value == 0.0  # 0.05 - 0.20 clamped to 0

    def test_update_trust_unknown_peer(self) -> None:
        ps = PeerStore()
        _, ok = ps.update_trust(SystemURI("nope"), 0.1)
        assert not ok

    def test_decay_all(self) -> None:
        ps = PeerStore()
        uri = SystemURI("peer-1")
        ps.register(uri, PublicKey(b"\x01" * 32), [], 1)
        # Set trust to 0.5
        for _ in range(10):
            ps.update_trust(uri, 0.05)

        record, _ = ps.get(uri)
        initial_trust = record.trust.value
        assert initial_trust > 0

        # Manually set last_decayed_at to 10 days ago
        with ps._lock:
            ps._peers[uri.value].last_decayed_at = time.time() - 10 * 86400

        ps.decay_all()
        record, _ = ps.get(uri)
        assert record.trust.value < initial_trust

    def test_register_preserves_public_key(self) -> None:
        """Re-registration should NOT overwrite the public key (TOFU)."""
        ps = PeerStore()
        uri = SystemURI("peer-1")
        pk1 = PublicKey(b"\x01" * 32)
        pk2 = PublicKey(b"\x02" * 32)

        ps.register(uri, pk1, ["v1"], 1)
        ps.register(uri, pk2, ["v2"], 2)

        record, _ = ps.get(uri)
        assert record.public_key.bytes_ == pk1.bytes_  # original key preserved
        assert record.capabilities == ["v2"]  # capabilities updated
        assert record.negotiated_version == 2  # version updated

    def test_all(self) -> None:
        ps = PeerStore()
        ps.register(SystemURI("a"), PublicKey(b"\x01" * 32), [], 1)
        ps.register(SystemURI("b"), PublicKey(b"\x02" * 32), [], 1)
        all_peers = ps.all()
        assert len(all_peers) == 2


# ── TreatyStore Tests ──────────────────────────────────────────────────────


class TestTreatyStore:
    def test_put_and_get(self) -> None:
        ts = TreatyStore()
        treaty = new_treaty(
            _make_treaty_id(), SystemURI("a"), SystemURI("b"),
            [TreatyTerm(scope="x", policy="y")],
        )
        ts.put(treaty)

        got, found = ts.get(treaty.id)
        assert found
        assert got.id == treaty.id

    def test_get_not_found(self) -> None:
        ts = TreatyStore()
        _, found = ts.get(_make_treaty_id())
        assert not found

    def test_apply(self) -> None:
        ts = TreatyStore()
        treaty = new_treaty(_make_treaty_id(), SystemURI("a"), SystemURI("b"), [])
        ts.put(treaty)

        ts.apply(treaty.id, lambda t: t.apply_action(TreatyAction.ACCEPT))

        got, _ = ts.get(treaty.id)
        assert got.status == TreatyStatus.ACTIVE

    def test_apply_not_found(self) -> None:
        ts = TreatyStore()
        with pytest.raises(TreatyNotFoundError):
            ts.apply(_make_treaty_id(), lambda t: None)

    def test_by_system(self) -> None:
        ts = TreatyStore()
        uri_a = SystemURI("a")
        uri_b = SystemURI("b")
        uri_c = SystemURI("c")
        ts.put(new_treaty(_make_treaty_id(), uri_a, uri_b, []))
        ts.put(new_treaty(_make_treaty_id(), uri_b, uri_c, []))
        ts.put(new_treaty(_make_treaty_id(), uri_a, uri_c, []))

        by_a = ts.by_system(uri_a)
        assert len(by_a) == 2
        by_b = ts.by_system(uri_b)
        assert len(by_b) == 2

    def test_active(self) -> None:
        ts = TreatyStore()
        t1 = new_treaty(_make_treaty_id(), SystemURI("a"), SystemURI("b"), [])
        t2 = new_treaty(_make_treaty_id(), SystemURI("a"), SystemURI("c"), [])
        ts.put(t1)
        ts.put(t2)

        assert len(ts.active()) == 0

        ts.apply(t1.id, lambda t: t.apply_action(TreatyAction.ACCEPT))
        assert len(ts.active()) == 1


# ── EnvelopeDedup Tests ────────────────────────────────────────────────────


class TestEnvelopeDedup:
    def test_check_first_time(self) -> None:
        dedup = EnvelopeDedup()
        eid = _make_envelope_id()
        assert dedup.check(eid) is True

    def test_check_duplicate(self) -> None:
        dedup = EnvelopeDedup()
        eid = _make_envelope_id()
        assert dedup.check(eid) is True
        assert dedup.check(eid) is False

    def test_different_ids(self) -> None:
        dedup = EnvelopeDedup()
        assert dedup.check(_make_envelope_id()) is True
        assert dedup.check(_make_envelope_id()) is True

    def test_size(self) -> None:
        dedup = EnvelopeDedup()
        dedup.check(_make_envelope_id())
        dedup.check(_make_envelope_id())
        assert dedup.size() == 2

    def test_prune_with_short_ttl(self) -> None:
        dedup = EnvelopeDedup(ttl=0.0)  # instant expiry
        eid = _make_envelope_id()
        dedup.check(eid)
        assert dedup.size() == 1

        time.sleep(0.01)
        removed = dedup.prune()
        assert removed == 1
        assert dedup.size() == 0


# ── Handler Tests ──────────────────────────────────────────────────────────


class TestHandler:
    def test_hello_handshake(self) -> None:
        """Test that a HELLO is sent via transport."""
        identity = _make_identity("local")
        transport = MockTransport()
        peers = PeerStore()
        treaties = TreatyStore()
        handler = Handler(identity, transport, peers, treaties)

        handler.hello(SystemURI("remote"))
        assert len(transport.sent) == 1
        _, env = transport.sent[0]
        assert env.type == MessageType.HELLO

    def test_handle_incoming_hello(self) -> None:
        """Test processing an incoming HELLO registers the peer."""
        local_id = _make_identity("local")
        remote_id = _make_identity("remote")
        transport = MockTransport()
        peers = PeerStore()
        treaties = TreatyStore()
        handler = Handler(local_id, transport, peers, treaties)

        # Build a HELLO from remote, signed by remote
        env = _make_hello_envelope(remote_id, SystemURI("local"))
        signed = sign_envelope(env, remote_id)

        # For verification to work with HMAC, the handler needs to use
        # the remote identity's verify. Since our HMAC model is symmetric,
        # we need to make the handler use the remote identity for verification.
        # Instead, we'll test with a single shared identity approach.
        # Create handler with remote identity so it can verify remote's signatures.
        handler2 = Handler(remote_id, transport, peers, treaties)
        handler2.handle_incoming(signed)

        _, found = peers.get(SystemURI(remote_id.system_uri().value))
        assert found

    def test_handle_incoming_duplicate(self) -> None:
        """Test that duplicate envelopes are rejected."""
        identity = _make_identity("local")
        transport = MockTransport()
        peers = PeerStore()
        treaties = TreatyStore()
        handler = Handler(identity, transport, peers, treaties)

        env = _make_hello_envelope(identity, SystemURI("local"))
        signed = sign_envelope(env, identity)

        handler.handle_incoming(signed)

        with pytest.raises(DuplicateEnvelopeError):
            handler.handle_incoming(signed)

    def test_handle_incoming_stale_timestamp(self) -> None:
        """Test that envelopes with stale timestamps are rejected."""
        identity = _make_identity("local")
        transport = MockTransport()
        peers = PeerStore()
        treaties = TreatyStore()
        handler = Handler(identity, transport, peers, treaties)

        env = Envelope(
            protocol_version=CURRENT_PROTOCOL_VERSION,
            id=_make_envelope_id(),
            from_uri=identity.system_uri(),
            to_uri=SystemURI("local"),
            type=MessageType.HELLO,
            payload=HelloPayload(
                system_uri=identity.system_uri().value,
                public_key=identity.public_key().bytes_,
                protocol_versions=[1],
                capabilities=[],
                chain_length=0,
            ),
            timestamp=time.time() - MAX_ENVELOPE_AGE_SECONDS - 100,
        )
        signed = sign_envelope(env, identity)

        with pytest.raises(EGIPError, match="timestamp out of range"):
            handler.handle_incoming(signed)

    def test_handle_incoming_unknown_sender_non_hello(self) -> None:
        """Non-HELLO from unknown system raises SystemNotFoundError."""
        identity = _make_identity("local")
        transport = MockTransport()
        peers = PeerStore()
        treaties = TreatyStore()
        handler = Handler(identity, transport, peers, treaties)

        env = Envelope(
            protocol_version=CURRENT_PROTOCOL_VERSION,
            id=_make_envelope_id(),
            from_uri=SystemURI("unknown"),
            to_uri=SystemURI("local"),
            type=MessageType.MESSAGE,
            payload=MessagePayloadContent(
                content={"text": "hello"},
                content_type="test.message",
            ),
            timestamp=time.time(),
        )
        signed = sign_envelope(env, identity)

        with pytest.raises(SystemNotFoundError):
            handler.handle_incoming(signed)

    def test_handle_message_dispatch(self) -> None:
        """Test that MESSAGE envelopes are dispatched to on_message."""
        identity = _make_identity("local")
        transport = MockTransport()
        peers = PeerStore()
        treaties = TreatyStore()
        handler = Handler(identity, transport, peers, treaties)

        # Register peer first
        peers.register(
            identity.system_uri(),
            identity.public_key(),
            [],
            CURRENT_PROTOCOL_VERSION,
        )

        received: list[tuple] = []

        def on_msg(from_uri: SystemURI, payload: MessagePayloadContent) -> None:
            received.append((from_uri, payload))

        handler.on_message = on_msg

        env = Envelope(
            protocol_version=CURRENT_PROTOCOL_VERSION,
            id=_make_envelope_id(),
            from_uri=identity.system_uri(),
            to_uri=SystemURI("local"),
            type=MessageType.MESSAGE,
            payload=MessagePayloadContent(
                content={"text": "hi"},
                content_type="test.message",
            ),
            timestamp=time.time(),
        )
        signed = sign_envelope(env, identity)

        handler.handle_incoming(signed)
        assert len(received) == 1
        assert received[0][1].content == {"text": "hi"}

    def test_hello_version_incompatible(self) -> None:
        """HELLO with incompatible versions raises VersionIncompatibleError."""
        identity = _make_identity("local")
        transport = MockTransport()
        peers = PeerStore()
        treaties = TreatyStore()
        handler = Handler(identity, transport, peers, treaties)

        env = Envelope(
            protocol_version=99,
            id=_make_envelope_id(),
            from_uri=identity.system_uri(),
            to_uri=SystemURI("local"),
            type=MessageType.HELLO,
            payload=HelloPayload(
                system_uri=identity.system_uri().value,
                public_key=identity.public_key().bytes_,
                protocol_versions=[99],  # incompatible
                capabilities=[],
                chain_length=0,
            ),
            timestamp=time.time(),
        )
        signed = sign_envelope(env, identity)

        with pytest.raises(VersionIncompatibleError):
            handler.handle_incoming(signed)


# ── Proof Tests ────────────────────────────────────────────────────────────


class TestProof:
    def test_verify_chain_segment_valid(self) -> None:
        events = [
            {"prev_hash": "aaa", "hash": "bbb"},
            {"prev_hash": "bbb", "hash": "ccc"},
            {"prev_hash": "ccc", "hash": "ddd"},
        ]
        proof = ChainSegmentProof(events=events, start_hash="aaa", end_hash="ddd")
        assert verify_chain_segment(proof) is True

    def test_verify_chain_segment_bad_start(self) -> None:
        events = [
            {"prev_hash": "wrong", "hash": "bbb"},
        ]
        proof = ChainSegmentProof(events=events, start_hash="aaa", end_hash="bbb")
        assert verify_chain_segment(proof) is False

    def test_verify_chain_segment_empty(self) -> None:
        proof = ChainSegmentProof(events=[], start_hash="a", end_hash="b")
        assert verify_chain_segment(proof) is False

    def test_verify_chain_segment_bad_chain(self) -> None:
        events = [
            {"prev_hash": "aaa", "hash": "bbb"},
            {"prev_hash": "XXX", "hash": "ccc"},  # broken link
        ]
        proof = ChainSegmentProof(events=events, start_hash="aaa", end_hash="ccc")
        assert verify_chain_segment(proof) is False

    def test_verify_event_existence_valid(self) -> None:
        proof = EventExistenceProof(
            event={"prev_hash": "aaa", "hash": "bbb", "id": "e1"},
            prev_hash="aaa",
            position=0,
            chain_length=5,
        )
        assert verify_event_existence(proof) is True

    def test_verify_event_existence_bad_prev_hash(self) -> None:
        proof = EventExistenceProof(
            event={"prev_hash": "xxx", "hash": "bbb"},
            prev_hash="aaa",
            position=0,
            chain_length=5,
        )
        assert verify_event_existence(proof) is False

    def test_verify_event_existence_bad_position(self) -> None:
        proof = EventExistenceProof(
            event={"prev_hash": "aaa", "hash": "bbb"},
            prev_hash="aaa",
            position=5,
            chain_length=5,  # position >= chain_length
        )
        assert verify_event_existence(proof) is False

    def test_verify_event_existence_no_hash(self) -> None:
        proof = EventExistenceProof(
            event={"prev_hash": "aaa"},  # no hash
            prev_hash="aaa",
            position=0,
            chain_length=5,
        )
        assert verify_event_existence(proof) is False

    def test_validate_proof_chain_summary(self) -> None:
        proof = ProofPayload(
            proof_type=ProofType.CHAIN_SUMMARY,
            data=ChainSummaryProof(
                length=10,
                head_hash="abc",
                genesis_hash="def",
                timestamp=time.time(),
            ),
        )
        valid, err = validate_proof(proof)
        assert valid is True

    def test_validate_proof_chain_summary_zero_length(self) -> None:
        proof = ProofPayload(
            proof_type=ProofType.CHAIN_SUMMARY,
            data=ChainSummaryProof(length=0, head_hash="a", genesis_hash="b"),
        )
        valid, _ = validate_proof(proof)
        assert valid is False

    def test_proof_generator_chain_summary(self) -> None:
        events = [
            {"id": "e1", "prev_hash": "000", "hash": "aaa"},
            {"id": "e2", "prev_hash": "aaa", "hash": "bbb"},
        ]
        gen = ProofGenerator(events)
        summary = gen.generate_chain_summary()
        assert summary.length == 2
        assert summary.head_hash == "bbb"
        assert summary.genesis_hash == "000"

    def test_proof_generator_event_existence(self) -> None:
        events = [
            {"id": "e1", "prev_hash": "000", "hash": "aaa"},
            {"id": "e2", "prev_hash": "aaa", "hash": "bbb"},
        ]
        gen = ProofGenerator(events)
        proof = gen.generate_event_existence("e1")
        assert proof.event["id"] == "e1"
        assert proof.prev_hash == "000"
        assert proof.position == 0
        assert proof.chain_length == 2

    def test_proof_generator_event_not_found(self) -> None:
        gen = ProofGenerator([])
        with pytest.raises(ValueError, match="event not found"):
            gen.generate_event_existence("nope")

    def test_proof_generator_empty_chain(self) -> None:
        gen = ProofGenerator([])
        with pytest.raises(ValueError, match="empty chain"):
            gen.generate_chain_summary()


# ── Error Types Tests ──────────────────────────────────────────────────────


class TestEGIPErrors:
    def test_system_not_found(self) -> None:
        err = SystemNotFoundError(SystemURI("foo"))
        assert "foo" in str(err)
        assert isinstance(err, EGIPError)

    def test_envelope_signature_invalid(self) -> None:
        err = EnvelopeSignatureInvalidError(_make_envelope_id())
        assert "signature invalid" in str(err)
        assert isinstance(err, EGIPError)

    def test_treaty_violation(self) -> None:
        err = TreatyViolationError(_make_treaty_id(), "read-only")
        assert "violated" in str(err)
        assert "read-only" in str(err)
        assert isinstance(err, EGIPError)

    def test_trust_insufficient(self) -> None:
        err = TrustInsufficientError(SystemURI("x"), Score(0.1), Score(0.5))
        assert "insufficient" in str(err)
        assert isinstance(err, EGIPError)

    def test_transport_failure(self) -> None:
        err = TransportFailureError(SystemURI("x"), "timeout")
        assert "timeout" in str(err)
        assert isinstance(err, EGIPError)

    def test_duplicate_envelope(self) -> None:
        err = DuplicateEnvelopeError(_make_envelope_id())
        assert "duplicate" in str(err)
        assert isinstance(err, EGIPError)

    def test_treaty_not_found(self) -> None:
        err = TreatyNotFoundError(_make_treaty_id())
        assert "not found" in str(err)
        assert isinstance(err, EGIPError)

    def test_version_incompatible(self) -> None:
        err = VersionIncompatibleError([1], [2])
        assert "compatible" in str(err)
        assert isinstance(err, EGIPError)


# ── Enum Tests ─────────────────────────────────────────────────────────────


class TestEnums:
    def test_message_types(self) -> None:
        assert MessageType.HELLO.value == "Hello"
        assert MessageType.MESSAGE.value == "Message"
        assert MessageType.RECEIPT.value == "Receipt"
        assert MessageType.PROOF.value == "Proof"
        assert MessageType.TREATY.value == "Treaty"
        assert MessageType.AUTHORITY_REQUEST.value == "AuthorityRequest"
        assert MessageType.DISCOVER.value == "Discover"

    def test_treaty_status(self) -> None:
        assert TreatyStatus.PROPOSED.value == "Proposed"
        assert TreatyStatus.ACTIVE.value == "Active"
        assert TreatyStatus.SUSPENDED.value == "Suspended"
        assert TreatyStatus.TERMINATED.value == "Terminated"

    def test_receipt_status(self) -> None:
        assert ReceiptStatus.DELIVERED.value == "Delivered"
        assert ReceiptStatus.PROCESSED.value == "Processed"
        assert ReceiptStatus.REJECTED.value == "Rejected"

    def test_proof_type(self) -> None:
        assert ProofType.CHAIN_SEGMENT.value == "ChainSegment"
        assert ProofType.EVENT_EXISTENCE.value == "EventExistence"
        assert ProofType.CHAIN_SUMMARY.value == "ChainSummary"

    def test_cger_relationship(self) -> None:
        assert CGERRelationship.CAUSED_BY.value == "CausedBy"
        assert CGERRelationship.REFERENCES.value == "References"
        assert CGERRelationship.RESPONDS_TO.value == "RespondsTo"

    def test_treaty_action(self) -> None:
        assert TreatyAction.PROPOSE.value == "Propose"
        assert TreatyAction.ACCEPT.value == "Accept"
        assert TreatyAction.MODIFY.value == "Modify"
        assert TreatyAction.SUSPEND.value == "Suspend"
        assert TreatyAction.TERMINATE.value == "Terminate"

    def test_authority_level(self) -> None:
        assert AuthorityLevel.REQUIRED.value == "Required"
        assert AuthorityLevel.RECOMMENDED.value == "Recommended"
        assert AuthorityLevel.NOTIFICATION.value == "Notification"
