package store_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/store"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// appendTrustEvents appends n trust.updated events to s, reusing one factory.
// Returns the appended events in insertion order.
func appendTrustEvents(tb testing.TB, s store.Store, n int, cause types.EventID) []event.Event {
	tb.Helper()
	registry := event.DefaultRegistry()
	factory := event.NewEventFactory(registry)
	events := make([]event.Event, 0, n)
	for i := 0; i < n; i++ {
		ev, err := factory.Create(
			event.EventTypeTrustUpdated,
			types.MustActorID("actor_00000000000000000000000000000001"),
			event.TrustUpdatedContent{
				Actor:    types.MustActorID("actor_00000000000000000000000000000002"),
				Previous: types.MustScore(0.5),
				Current:  types.MustScore(0.6),
				Domain:   types.MustDomainScope("test"),
				Cause:    cause,
			},
			[]types.EventID{cause},
			types.MustConversationID("conv_00000000000000000000000000000001"),
			headFromStore{s},
			testSigner{},
		)
		if err != nil {
			tb.Fatalf("create trust event %d: %v", i, err)
		}
		if _, err := s.Append(ev); err != nil {
			tb.Fatalf("append trust event %d: %v", i, err)
		}
		events = append(events, ev)
	}
	return events
}

// walkAllPages pages through query until exhaustion, returning all items.
// Fails the test if the walk does not terminate within maxPages pages.
func walkAllPages(tb testing.TB, query func(limit int, after types.Option[types.Cursor]) (types.Page[event.Event], error), pageSize, maxPages int) []event.Event {
	tb.Helper()
	var all []event.Event
	after := types.None[types.Cursor]()
	for pages := 0; ; pages++ {
		if pages > maxPages {
			tb.Fatalf("walk did not terminate within %d pages", maxPages)
		}
		page, err := query(pageSize, after)
		if err != nil {
			tb.Fatalf("page %d (after %d items): %v", pages, len(all), err)
		}
		all = append(all, page.Items()...)
		if !page.HasMore() {
			return all
		}
		after = page.Cursor()
		if !after.IsSome() {
			tb.Fatalf("page %d reports HasMore but carries no cursor", pages)
		}
	}
}

// TestByTypeMultiPageWalkReturnsAllEventsInReverseOrder pins the multi-page
// contract: a cursor walk over a type index yields every event of that type
// exactly once, most recent first, even with other types interleaved.
func TestByTypeMultiPageWalkReturnsAllEventsInReverseOrder(t *testing.T) {
	s := store.NewInMemoryStore()
	boot := makeBootstrapEvent(t)
	if _, err := s.Append(boot); err != nil {
		t.Fatalf("append bootstrap: %v", err)
	}

	// Interleave trust.updated with decision.recorded so the type index is a
	// strict subset of the global event sequence.
	var trust []event.Event
	for i := 0; i < 25; i++ {
		tev := makeEvent(t, s, event.EventTypeTrustUpdated, []types.EventID{boot.ID()})
		if _, err := s.Append(tev); err != nil {
			t.Fatalf("append trust %d: %v", i, err)
		}
		trust = append(trust, tev)
		dev := makeDecisionRecordedEvent(t, s, []types.EventID{boot.ID()})
		if _, err := s.Append(dev); err != nil {
			t.Fatalf("append decision %d: %v", i, err)
		}
	}

	got := walkAllPages(t, func(limit int, after types.Option[types.Cursor]) (types.Page[event.Event], error) {
		return s.ByType(event.EventTypeTrustUpdated, limit, after)
	}, 4, 100)

	if len(got) != len(trust) {
		t.Fatalf("walk returned %d events, want %d", len(got), len(trust))
	}
	for i, ev := range got {
		want := trust[len(trust)-1-i]
		if ev.ID() != want.ID() {
			t.Fatalf("position %d: got %s, want %s (reverse insertion order)", i, ev.ID(), want.ID())
		}
	}
}

// TestPaginateReverseCursorEdgeSemantics pins the memory store's fail-closed
// cursor semantics: any cursor that does not exactly match an event in the
// queried index is rejected with InvalidCursorError.
func TestPaginateReverseCursorEdgeSemantics(t *testing.T) {
	s := store.NewInMemoryStore()
	boot := makeBootstrapEvent(t)
	if _, err := s.Append(boot); err != nil {
		t.Fatalf("append bootstrap: %v", err)
	}
	trust := appendTrustEvents(t, s, 3, boot.ID())
	decision := makeDecisionRecordedEvent(t, s, []types.EventID{boot.ID()})
	if _, err := s.Append(decision); err != nil {
		t.Fatalf("append decision: %v", err)
	}

	byTrustType := func(after types.Cursor) error {
		_, err := s.ByType(event.EventTypeTrustUpdated, 10, types.Some(after))
		return err
	}
	wantInvalidCursor := func(t *testing.T, err error) {
		t.Helper()
		if err == nil {
			t.Fatal("expected InvalidCursorError, got nil")
		}
		var invalid *store.InvalidCursorError
		if !errors.As(err, &invalid) {
			t.Fatalf("expected InvalidCursorError, got %T: %v", err, err)
		}
	}

	t.Run("CursorFromDifferentTypeIsRejected", func(t *testing.T) {
		// decision exists in the store but not in the trust.updated index.
		wantInvalidCursor(t, byTrustType(types.MustCursor(decision.ID().Value())))
	})

	t.Run("UnknownEventIDCursorIsRejected", func(t *testing.T) {
		wantInvalidCursor(t, byTrustType(types.MustCursor("019462a0-0000-7000-8000-0000000000ff")))
	})

	t.Run("MalformedCursorIsRejected", func(t *testing.T) {
		wantInvalidCursor(t, byTrustType(types.MustCursor("not-an-event-id")))
	})

	t.Run("CaseVariantCursorIsRejected", func(t *testing.T) {
		// Stored IDs are lowercase; a case-variant cursor has never matched
		// and must keep failing closed rather than being normalized.
		wantInvalidCursor(t, byTrustType(types.MustCursor(strings.ToUpper(trust[1].ID().Value()))))
	})

	t.Run("CursorAtOldestItemYieldsEmptyFinalPage", func(t *testing.T) {
		page, err := s.ByType(event.EventTypeTrustUpdated, 10, types.Some(types.MustCursor(trust[0].ID().Value())))
		if err != nil {
			t.Fatalf("ByType: %v", err)
		}
		if len(page.Items()) != 0 {
			t.Fatalf("expected empty page past the oldest item, got %d items", len(page.Items()))
		}
		if page.HasMore() {
			t.Fatal("expected HasMore=false past the oldest item")
		}
	})
}

// TestByTypeMultiPageWalkScalesLinearly guards against the O(N²/pageSize)
// regression where paginateReverse resolved each page's cursor by scanning
// the type index linearly. It measures a full ByType walk at N and 4N events:
// a linear walk grows ~4x, the quadratic scan ~16x. The 8x threshold sits
// between the two with 2x margin on each side; min-of-5 damps scheduler noise.
func TestByTypeMultiPageWalkScalesLinearly(t *testing.T) {
	if testing.Short() {
		t.Skip("scaling measurement skipped in -short mode")
	}
	const (
		smallN   = 2500
		largeN   = 10000
		pageSize = 10
		maxRatio = 8.0
	)

	measureWalk := func(n int) time.Duration {
		s := store.NewInMemoryStore()
		boot := makeBootstrapEvent(t)
		if _, err := s.Append(boot); err != nil {
			t.Fatalf("append bootstrap: %v", err)
		}
		appendTrustEvents(t, s, n, boot.ID())

		var best time.Duration
		for rep := 0; rep < 5; rep++ {
			start := time.Now()
			got := walkAllPages(t, func(limit int, after types.Option[types.Cursor]) (types.Page[event.Event], error) {
				return s.ByType(event.EventTypeTrustUpdated, limit, after)
			}, pageSize, n/pageSize+2)
			elapsed := time.Since(start)
			if len(got) != n {
				t.Fatalf("walk returned %d of %d events", len(got), n)
			}
			if rep == 0 || elapsed < best {
				best = elapsed
			}
		}
		return best
	}

	small := measureWalk(smallN)
	large := measureWalk(largeN)
	// Floor the denominator so clock jitter on a very fast small walk cannot
	// fabricate a failure; a genuinely quadratic large walk (tens of ms at
	// 10k events) still exceeds the threshold against this floor.
	const denominatorFloor = 100 * time.Microsecond
	denom := small
	if denom < denominatorFloor {
		denom = denominatorFloor
	}
	ratio := float64(large) / float64(denom)
	t.Logf("full walk: %d events in %v, %d events in %v (ratio %.1fx)", smallN, small, largeN, large, ratio)
	if ratio > maxRatio {
		t.Fatalf("multi-page walk grew %.1fx for a 4x input increase (want <= %.1fx): cursor resolution is superlinear", ratio, maxRatio)
	}
}

// BenchmarkByTypeMultiPageWalk measures a full cursor walk over a type index.
// Before the byID+binary-search cursor resolution this was quadratic in the
// number of events of the type; after it, linear.
func BenchmarkByTypeMultiPageWalk(b *testing.B) {
	for _, n := range []int{2500, 10000, 25000} {
		b.Run(fmt.Sprintf("events=%d", n), func(b *testing.B) {
			s := store.NewInMemoryStore()
			boot := makeBootstrapEvent(b)
			if _, err := s.Append(boot); err != nil {
				b.Fatalf("append bootstrap: %v", err)
			}
			appendTrustEvents(b, s, n, boot.ID())
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				got := walkAllPages(b, func(limit int, after types.Option[types.Cursor]) (types.Page[event.Event], error) {
					return s.ByType(event.EventTypeTrustUpdated, limit, after)
				}, 100, n/100+2)
				if len(got) != n {
					b.Fatalf("walk returned %d of %d events", len(got), n)
				}
			}
		})
	}
}
