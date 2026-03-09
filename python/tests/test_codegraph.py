"""Tests for the 61 code graph primitives, 35 event types, and 7 compositions."""

from __future__ import annotations

import pytest

from eventgraph.codegraph import (
    # Primitives
    ALL_CODEGRAPH_PRIMITIVE_CLASSES,
    _CodeGraphBase,
    all_codegraph_primitives,
    register_all_codegraph,
    is_codegraph_primitive,
    # Compositions
    CodeGraphComposition,
    all_codegraph_compositions,
    board,
    detail,
    feed,
    dashboard,
    inbox,
    wizard,
    skin,
    # Event types
    all_codegraph_event_types,
)
from eventgraph.event import Event, NoopSigner, create_bootstrap, create_event
from eventgraph.primitive import (
    Registry,
    Snapshot,
    UpdateState,
)
from eventgraph.types import (
    ActorID,
    ConversationID,
    EventType,
    PrimitiveID,
    SubscriptionPattern,
)


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def _empty_snapshot(tick: int = 1) -> Snapshot:
    return Snapshot(tick=tick, primitives={}, pending_events=[], recent_events=[])


def _make_event(event_type: str) -> Event:
    """Create a minimal Event with the given type for testing process()."""
    signer = NoopSigner()
    bootstrap = create_bootstrap(source=ActorID("system"), signer=signer)
    return create_event(
        event_type=EventType(event_type),
        source=ActorID("test-cg"),
        content={"test": True},
        causes=[bootstrap.id],
        conversation_id=ConversationID("conv_test"),
        prev_hash=bootstrap.hash,
        signer=signer,
    )


# ===========================================================================
# 1. Code Graph Event Types
# ===========================================================================

class TestCodeGraphEventTypes:
    def test_all_codegraph_event_types(self) -> None:
        types = all_codegraph_event_types()
        assert len(types) == 35

    def test_all_start_with_codegraph(self) -> None:
        for et in all_codegraph_event_types():
            assert et.value.startswith("codegraph.")

    def test_all_unique(self) -> None:
        values = [et.value for et in all_codegraph_event_types()]
        assert len(values) == len(set(values))


# ===========================================================================
# 2. All 61 Code Graph Primitives
# ===========================================================================

class TestAllCodeGraphPrimitives:
    def test_all_primitives_count(self) -> None:
        prims = all_codegraph_primitives()
        assert len(prims) == 61

    def test_no_duplicate_ids(self) -> None:
        prims = all_codegraph_primitives()
        ids = [p.id().value for p in prims]
        assert len(ids) == len(set(ids)), f"duplicate IDs found: {ids}"

    def test_all_layer_5(self) -> None:
        for p in all_codegraph_primitives():
            assert p.layer().value == 5, f"{p.id().value} not at layer 5"

    def test_all_cadence_1(self) -> None:
        for p in all_codegraph_primitives():
            assert p.cadence().value == 1

    def test_all_have_subscriptions(self) -> None:
        for p in all_codegraph_primitives():
            subs = p.subscriptions()
            assert len(subs) > 0, f"{p.id().value} has no subscriptions"
            for s in subs:
                assert isinstance(s, SubscriptionPattern)

    def test_all_ids_start_with_cg(self) -> None:
        for p in all_codegraph_primitives():
            assert p.id().value.startswith("CG"), f"{p.id().value} missing CG prefix"

    def test_is_codegraph_primitive(self) -> None:
        for p in all_codegraph_primitives():
            assert is_codegraph_primitive(p.id()) is True
        assert is_codegraph_primitive(PrimitiveID("notacodegraph")) is False
        assert is_codegraph_primitive(PrimitiveID("agent.Identity")) is False

    def test_isinstance_codegraph_base(self) -> None:
        for p in all_codegraph_primitives():
            assert isinstance(p, _CodeGraphBase)


# ===========================================================================
# 3. Registration
# ===========================================================================

class TestRegistration:
    def test_register_all(self) -> None:
        reg = Registry()
        register_all_codegraph(reg)
        assert reg.count() == 61

    def test_register_all_unique(self) -> None:
        reg = Registry()
        register_all_codegraph(reg)
        prims = reg.all()
        ids = [p.id().value for p in prims]
        assert len(ids) == len(set(ids))

    def test_all_active_after_register(self) -> None:
        from eventgraph.primitive import LIFECYCLE_ACTIVE
        reg = Registry()
        register_all_codegraph(reg)
        for p in reg.all():
            assert reg.lifecycle(p.id()) == LIFECYCLE_ACTIVE


# ===========================================================================
# 4. Process Returns Mutations
# ===========================================================================

class TestProcessReturnsMutations:
    def test_process_returns_mutations(self) -> None:
        snap = _empty_snapshot()
        for p in all_codegraph_primitives():
            mutations = p.process(tick=1, events=[], snapshot=snap)
            assert isinstance(mutations, list)
            assert len(mutations) > 0, f"{p.id().value} returned no mutations"
            # All should include a lastTick mutation
            last_ticks = [
                m for m in mutations
                if isinstance(m, UpdateState) and m.key == "lastTick"
            ]
            assert len(last_ticks) == 1, f"{p.id().value} missing lastTick mutation"
            assert last_ticks[0].value == 1

    def test_process_with_events(self) -> None:
        snap = _empty_snapshot(tick=5)
        ev = _make_event("codegraph.entity.defined")
        for p in all_codegraph_primitives():
            mutations = p.process(tick=5, events=[ev], snapshot=snap)
            state = {m.key: m.value for m in mutations if isinstance(m, UpdateState)}
            assert state["lastTick"] == 5
            assert state["eventsProcessed"] == 1

    def test_process_counts_multiple_events(self) -> None:
        snap = _empty_snapshot(tick=3)
        evs = [
            _make_event("codegraph.entity.defined"),
            _make_event("codegraph.ui.view.rendered"),
            _make_event("codegraph.io.query.executed"),
        ]
        p = all_codegraph_primitives()[0]
        mutations = p.process(tick=3, events=evs, snapshot=snap)
        state = {m.key: m.value for m in mutations if isinstance(m, UpdateState)}
        assert state["eventsProcessed"] == 3


# ===========================================================================
# 5. Compositions
# ===========================================================================

class TestCompositions:
    def test_all_compositions(self) -> None:
        comps = all_codegraph_compositions()
        assert len(comps) == 7

    def test_all_unique_names(self) -> None:
        comps = all_codegraph_compositions()
        names = [c.name for c in comps]
        assert len(names) == len(set(names))

    def test_composition_primitives_exist(self) -> None:
        """Every primitive ID referenced in compositions exists in the 61."""
        valid_ids = {p.id().value for p in all_codegraph_primitives()}
        for c in all_codegraph_compositions():
            for pid in c.primitives:
                assert pid in valid_ids, (
                    f"composition {c.name!r} references unknown primitive {pid!r}"
                )

    def test_board(self) -> None:
        c = board()
        assert c.name == "Board"
        assert len(c.primitives) == 10
        assert "CGLayout" in c.primitives
        assert "CGDrag" in c.primitives

    def test_detail(self) -> None:
        c = detail()
        assert c.name == "Detail"
        assert len(c.primitives) == 9
        assert "CGForm" in c.primitives
        assert "CGThread" in c.primitives

    def test_feed(self) -> None:
        c = feed()
        assert c.name == "Feed"
        assert len(c.primitives) == 7
        assert "CGPagination" in c.primitives
        assert "CGRecency" in c.primitives

    def test_dashboard(self) -> None:
        c = dashboard()
        assert c.name == "Dashboard"
        assert len(c.primitives) == 5
        assert "CGTransform" in c.primitives
        assert "CGSalience" in c.primitives

    def test_inbox(self) -> None:
        c = inbox()
        assert c.name == "Inbox"
        assert len(c.primitives) == 8
        assert "CGSalience" in c.primitives
        assert "CGSelection" in c.primitives

    def test_wizard(self) -> None:
        c = wizard()
        assert c.name == "Wizard"
        assert len(c.primitives) == 8
        assert "CGSequence" in c.primitives
        assert "CGConsequencePreview" in c.primitives

    def test_skin(self) -> None:
        c = skin()
        assert c.name == "Skin"
        assert len(c.primitives) == 7
        assert "CGPalette" in c.primitives
        assert "CGShape" in c.primitives

    def test_composition_is_frozen(self) -> None:
        c = board()
        with pytest.raises(AttributeError):
            c.name = "Modified"  # type: ignore[misc]
