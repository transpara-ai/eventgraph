"""Tests for the 201 EventGraph primitives."""

from __future__ import annotations

import hashlib
import time
import uuid

import pytest

from eventgraph.event import Event, create_bootstrap, NoopSigner
from eventgraph.primitive import (
    Mutation,
    PrimitiveState,
    Snapshot,
    UpdateState,
)
from eventgraph.primitives import (
    ALL_PRIMITIVE_CLASSES,
    create_all_primitives,
    _Base,
)
from eventgraph.types import (
    Activation,
    ActorID,
    Cadence,
    ConversationID,
    EventID,
    Hash,
    Layer,
    PrimitiveID,
    SubscriptionPattern,
)


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def _empty_snapshot(tick: int = 1) -> Snapshot:
    return Snapshot(tick=tick, primitives={}, pending_events=[], recent_events=[])


def _make_bootstrap() -> Event:
    signer = NoopSigner()
    return create_bootstrap(
        source=ActorID("system"),
        signer=signer,
    )


# ---------------------------------------------------------------------------
# 1. All 201 primitives instantiate
# ---------------------------------------------------------------------------

class TestAllPrimitivesInstantiate:
    def test_instantiate_all(self) -> None:
        primitives = create_all_primitives()
        assert len(primitives) == 201
        for p in primitives:
            assert isinstance(p, _Base)

    def test_each_class_instantiates(self) -> None:
        for cls in ALL_PRIMITIVE_CLASSES:
            p = cls()
            assert p.id().value  # non-empty name
            assert isinstance(p.layer(), Layer)
            assert isinstance(p.cadence(), Cadence)


# ---------------------------------------------------------------------------
# 2. create_all_primitives returns 201
# ---------------------------------------------------------------------------

class TestCreateAllPrimitives:
    def test_returns_201(self) -> None:
        assert len(create_all_primitives()) == 201

    def test_all_unique_ids(self) -> None:
        primitives = create_all_primitives()
        ids = [p.id().value for p in primitives]
        assert len(set(ids)) == 201, f"Duplicate IDs found: {[x for x in ids if ids.count(x) > 1]}"


# ---------------------------------------------------------------------------
# 3. Layer counts correct
# ---------------------------------------------------------------------------

EXPECTED_LAYER_COUNTS = {
    0: 45,
    1: 12,
    2: 12,
    3: 12,
    4: 12,
    5: 12,
    6: 12,
    7: 12,
    8: 12,
    9: 12,
    10: 12,
    11: 12,
    12: 12,
    13: 12,
}


class TestLayerCounts:
    def test_layer_counts(self) -> None:
        primitives = create_all_primitives()
        counts: dict[int, int] = {}
        for p in primitives:
            layer = p.layer().value
            counts[layer] = counts.get(layer, 0) + 1
        assert counts == EXPECTED_LAYER_COUNTS

    def test_all_layers_present(self) -> None:
        primitives = create_all_primitives()
        layers = {p.layer().value for p in primitives}
        assert layers == set(range(14))

    def test_total_from_layer_counts(self) -> None:
        assert sum(EXPECTED_LAYER_COUNTS.values()) == 201


# ---------------------------------------------------------------------------
# 4. Process returns mutations
# ---------------------------------------------------------------------------

class TestProcessReturnsMutations:
    def test_process_returns_list(self) -> None:
        for cls in ALL_PRIMITIVE_CLASSES:
            p = cls()
            result = p.process(1, [], _empty_snapshot())
            assert isinstance(result, list)

    def test_process_returns_update_state_mutations(self) -> None:
        primitives = create_all_primitives()
        for p in primitives:
            result = p.process(42, [], _empty_snapshot(42))
            assert len(result) == 2
            # eventsProcessed
            assert isinstance(result[0], UpdateState)
            assert result[0].key == "eventsProcessed"
            assert result[0].value == 0
            assert result[0].primitive_id == p.id()
            # lastTick
            assert isinstance(result[1], UpdateState)
            assert result[1].key == "lastTick"
            assert result[1].value == 42

    def test_process_counts_events(self) -> None:
        p = create_all_primitives()[0]  # Event primitive
        bootstrap = _make_bootstrap()
        result = p.process(5, [bootstrap, bootstrap, bootstrap], _empty_snapshot(5))
        assert result[0].value == 3  # eventsProcessed == 3
        assert result[1].value == 5  # lastTick == 5


# ---------------------------------------------------------------------------
# 5. Subscriptions non-empty
# ---------------------------------------------------------------------------

class TestSubscriptionsNonEmpty:
    def test_all_have_subscriptions(self) -> None:
        for cls in ALL_PRIMITIVE_CLASSES:
            p = cls()
            subs = p.subscriptions()
            assert len(subs) > 0, f"{p.id().value} has no subscriptions"
            for s in subs:
                assert isinstance(s, SubscriptionPattern)

    def test_subscription_values_non_empty(self) -> None:
        for cls in ALL_PRIMITIVE_CLASSES:
            p = cls()
            for s in p.subscriptions():
                assert s.value, f"{p.id().value} has empty subscription pattern"


# ---------------------------------------------------------------------------
# Cadence
# ---------------------------------------------------------------------------

class TestCadence:
    def test_all_cadence_at_least_one(self) -> None:
        for cls in ALL_PRIMITIVE_CLASSES:
            p = cls()
            assert p.cadence().value >= 1


# ---------------------------------------------------------------------------
# Specific primitives spot-checks
# ---------------------------------------------------------------------------

class TestSpecificPrimitives:
    def test_event_primitive(self) -> None:
        from eventgraph.primitives import EventPrimitive
        p = EventPrimitive()
        assert p.id().value == "Event"
        assert p.layer().value == 0
        assert SubscriptionPattern("*") in p.subscriptions()

    def test_being_primitive(self) -> None:
        from eventgraph.primitives import BeingPrimitive
        p = BeingPrimitive()
        assert p.id().value == "Being"
        assert p.layer().value == 13
        subs = [s.value for s in p.subscriptions()]
        assert "clock.tick" in subs

    def test_wonder_primitive(self) -> None:
        from eventgraph.primitives import WonderPrimitive
        p = WonderPrimitive()
        assert p.id().value == "Wonder"
        assert p.layer().value == 13
        assert SubscriptionPattern("*") in p.subscriptions()

    def test_deception_indicator_multi_sub(self) -> None:
        from eventgraph.primitives import DeceptionIndicatorPrimitive
        p = DeceptionIndicatorPrimitive()
        subs = [s.value for s in p.subscriptions()]
        assert "trust.*" in subs
        assert "violation.*" in subs
