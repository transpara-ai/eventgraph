"""Tests for the trust module."""

from __future__ import annotations

import time

from eventgraph.event import (
    Event,
    NoopSigner,
    create_event,
    new_event_id,
)
from eventgraph.trust import DefaultTrustModel, TrustConfig, TrustMetrics
from eventgraph.types import (
    ActorID,
    ConversationID,
    DomainScope,
    EventType,
    Hash,
    NonEmpty,
    Score,
    Weight,
)


def _make_trust_event(actor_id: ActorID, current: float) -> Event:
    """Create a trust.updated evidence event with a 'current' value."""
    return create_event(
        event_type=EventType("trust.updated"),
        source=actor_id,
        content={"current": current, "domain": "general"},
        causes=[new_event_id()],
        conversation_id=ConversationID("conv_test"),
        prev_hash=Hash.zero(),
        signer=NoopSigner(),
    )


def _make_non_trust_event(actor_id: ActorID) -> Event:
    """Create a non-trust evidence event (no 'current' key)."""
    return create_event(
        event_type=EventType("task.completed"),
        source=actor_id,
        content={"result": "success"},
        causes=[new_event_id()],
        conversation_id=ConversationID("conv_test"),
        prev_hash=Hash.zero(),
        signer=NoopSigner(),
    )


# ── Tests ────────────────────────────────────────────────────────────────


def test_initial_score():
    """New actor gets initial trust (default 0.0)."""
    model = DefaultTrustModel()
    actor = ActorID("alice")

    metrics = model.score(actor)
    assert metrics.overall.value == 0.0
    assert metrics.confidence.value == 0.0


def test_update_increases_trust():
    """Evidence with higher current raises score."""
    model = DefaultTrustModel()
    actor = ActorID("alice")

    ev = _make_trust_event(actor, 0.05)
    metrics = model.update(actor, ev)
    assert metrics.overall.value > 0.0


def test_update_decreases_trust():
    """Evidence with lower current lowers score."""
    model = DefaultTrustModel()
    actor = ActorID("alice")

    # First increase
    ev1 = _make_trust_event(actor, 0.08)
    model.update(actor, ev1)

    # Then decrease
    ev2 = _make_trust_event(actor, 0.0)
    metrics = model.update(actor, ev2)
    assert metrics.overall.value < 0.08


def test_update_clamped_to_max_adjustment():
    """Delta clamped to max_adjustment (default 0.1)."""
    model = DefaultTrustModel()
    actor = ActorID("alice")

    # Try to jump to 0.5 from 0.0 — delta would be 0.5, clamped to 0.1
    ev = _make_trust_event(actor, 0.5)
    metrics = model.update(actor, ev)
    assert metrics.overall.value <= 0.1


def test_update_deduplication():
    """Same evidence applied twice has no additional effect."""
    model = DefaultTrustModel()
    actor = ActorID("alice")

    ev = _make_trust_event(actor, 0.05)
    metrics1 = model.update(actor, ev)

    # Apply same event again
    metrics2 = model.update(actor, ev)
    assert metrics1.overall.value == metrics2.overall.value
    assert len(metrics1.evidence) == len(metrics2.evidence)


def test_update_trend_positive():
    """Positive delta increases trend."""
    model = DefaultTrustModel()
    actor = ActorID("alice")

    ev = _make_trust_event(actor, 0.05)
    metrics = model.update(actor, ev)
    assert metrics.trend.value > 0


def test_update_trend_negative():
    """Negative delta decreases trend."""
    model = DefaultTrustModel()
    actor = ActorID("alice")

    # First increase
    ev1 = _make_trust_event(actor, 0.05)
    model.update(actor, ev1)

    # Then decrease — delta will be negative (target 0.0, current ~0.05)
    ev2 = _make_trust_event(actor, 0.0)
    metrics = model.update(actor, ev2)
    # Trend started at 0, went to +0.1, then back to 0.0 with -0.1
    assert metrics.trend.value <= 0.0


def test_score_in_domain():
    """Domain-specific score returned when domain data exists."""
    model = DefaultTrustModel()
    actor = ActorID("alice")

    # Update with domain-carrying event
    ev = _make_trust_event(actor, 0.05)
    model.update(actor, ev)

    domain = DomainScope("general")
    metrics = model.score_in_domain(actor, domain)
    # Domain score should exist and be used as overall in the result
    assert metrics.overall.value > 0.0


def test_score_in_domain_fallback():
    """No domain data falls back to global score with halved confidence."""
    model = DefaultTrustModel()
    actor = ActorID("alice")

    # Add some evidence (non-domain)
    ev = _make_non_trust_event(actor)
    model.update(actor, ev)

    # Query a domain that has no data
    domain = DomainScope("code_review")
    metrics = model.score_in_domain(actor, domain)

    # Confidence should be halved: 1 evidence / 50 * 0.5 = 0.01
    assert metrics.confidence.value == (1 / 50.0) * 0.5


def test_decay():
    """Trust decays over time."""
    model = DefaultTrustModel()
    actor = ActorID("alice")

    # Build some trust
    for _ in range(5):
        ev = _make_trust_event(actor, 0.1)
        model.update(actor, ev)

    before = model.score(actor)

    # Decay 10 days (in seconds)
    model.decay(actor, 10 * 86400.0)

    after = model.score(actor)
    assert after.overall.value < before.overall.value


def test_decay_negative_duration():
    """Negative elapsed time causes no change."""
    model = DefaultTrustModel()
    actor = ActorID("alice")

    ev = _make_trust_event(actor, 0.05)
    model.update(actor, ev)

    before = model.score(actor)
    model.decay(actor, -100.0)
    after = model.score(actor)

    assert after.overall.value == before.overall.value


def test_update_between():
    """Directional trust update works."""
    model = DefaultTrustModel()
    alice = ActorID("alice")
    bob = ActorID("bob")

    ev = _make_trust_event(bob, 0.05)
    metrics = model.update_between(alice, bob, ev)
    assert metrics.overall.value > 0.0
    assert metrics.actor.value == bob.value


def test_between_no_relationship():
    """No directed trust returns initial trust with zero confidence."""
    model = DefaultTrustModel()
    alice = ActorID("alice")
    bob = ActorID("bob")

    metrics = model.between(alice, bob)
    assert metrics.overall.value == 0.0
    assert metrics.confidence.value == 0.0


def test_between_asymmetric():
    """A->B trust is different from B->A trust."""
    model = DefaultTrustModel()
    alice = ActorID("alice")
    bob = ActorID("bob")

    # Alice trusts Bob more
    ev1 = _make_trust_event(bob, 0.08)
    model.update_between(alice, bob, ev1)

    # Bob gives Alice small trust
    ev2 = _make_trust_event(alice, 0.02)
    model.update_between(bob, alice, ev2)

    ab = model.between(alice, bob)
    ba = model.between(bob, alice)

    assert ab.overall.value != ba.overall.value
    assert ab.overall.value > ba.overall.value


def test_evidence_capped_at_100():
    """Evidence list never exceeds 100 entries."""
    model = DefaultTrustModel()
    actor = ActorID("alice")

    for _ in range(150):
        ev = _make_trust_event(actor, 0.01)
        model.update(actor, ev)

    metrics = model.score(actor)
    assert len(metrics.evidence) <= 100


def test_decay_directed_trust():
    """Decay also affects directed trust where actor is the from side."""
    model = DefaultTrustModel()
    alice = ActorID("alice")
    bob = ActorID("bob")

    # Build directed trust
    for _ in range(5):
        ev = _make_trust_event(bob, 0.1)
        model.update_between(alice, bob, ev)

    before = model.between(alice, bob)

    # Decay alice (the from side) — should affect alice->bob
    model.decay(alice, 10 * 86400.0)

    after = model.between(alice, bob)
    assert after.overall.value < before.overall.value
