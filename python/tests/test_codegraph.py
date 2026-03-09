"""Tests for the 35 code graph event types and 7 compositions."""

from __future__ import annotations

import pytest

from eventgraph.codegraph import (
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
# 2. Compositions
# ===========================================================================

class TestCompositions:
    def test_all_compositions(self) -> None:
        comps = all_codegraph_compositions()
        assert len(comps) == 7

    def test_composition_names_unique(self) -> None:
        comps = all_codegraph_compositions()
        names = [c.name for c in comps]
        assert len(names) == len(set(names))

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
