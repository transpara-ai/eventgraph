"""Tests for EventBus — subscribe, publish, unsubscribe, close."""

import threading
import time

from eventgraph.bus import EventBus
from eventgraph.event import NoopSigner, create_bootstrap, create_event
from eventgraph.store import InMemoryStore
from eventgraph.types import (
    ActorID,
    ConversationID,
    EventType,
    SubscriptionPattern,
)


def _bootstrap():
    return create_bootstrap(source=ActorID("alice"), signer=NoopSigner())


def _event(prev, event_type="trust.updated"):
    return create_event(
        event_type=EventType(event_type),
        source=ActorID("alice"),
        content={},
        causes=[prev.id],
        conversation_id=ConversationID("conv_1"),
        prev_hash=prev.hash,
        signer=NoopSigner(),
    )


class TestEventBus:
    def test_subscribe_and_publish(self):
        store = InMemoryStore()
        bus = EventBus(store)
        received = []
        latch = threading.Event()

        def handler(ev):
            received.append(ev)
            latch.set()

        bus.subscribe(SubscriptionPattern("*"), handler)
        boot = _bootstrap()
        bus.publish(boot)
        latch.wait(timeout=2.0)
        bus.close()

        assert len(received) == 1
        assert received[0].id.value == boot.id.value

    def test_pattern_filtering(self):
        store = InMemoryStore()
        bus = EventBus(store)
        trust_events = []
        all_events = []
        latch = threading.Event()

        def trust_handler(ev):
            trust_events.append(ev)

        def all_handler(ev):
            all_events.append(ev)
            if len(all_events) >= 2:
                latch.set()

        bus.subscribe(SubscriptionPattern("trust.*"), trust_handler)
        bus.subscribe(SubscriptionPattern("*"), all_handler)

        boot = _bootstrap()
        store.append(boot)
        e1 = _event(boot, "trust.updated")
        store.append(e1)

        bus.publish(boot)  # system.bootstrapped — only wildcard
        bus.publish(e1)  # trust.updated — both

        latch.wait(timeout=2.0)
        time.sleep(0.1)  # let trust handler process
        bus.close()

        assert len(all_events) == 2
        assert len(trust_events) == 1
        assert trust_events[0].type.value == "trust.updated"

    def test_unsubscribe(self):
        store = InMemoryStore()
        bus = EventBus(store)
        received = []

        sub_id = bus.subscribe(SubscriptionPattern("*"), lambda ev: received.append(ev))
        bus.unsubscribe(sub_id)
        time.sleep(0.05)  # let thread exit

        boot = _bootstrap()
        bus.publish(boot)
        time.sleep(0.1)
        bus.close()

        assert len(received) == 0

    def test_close_prevents_subscribe(self):
        store = InMemoryStore()
        bus = EventBus(store)
        bus.close()

        sub_id = bus.subscribe(SubscriptionPattern("*"), lambda ev: None)
        assert sub_id == -1

    def test_close_prevents_publish(self):
        store = InMemoryStore()
        bus = EventBus(store)
        received = []
        bus.subscribe(SubscriptionPattern("*"), lambda ev: received.append(ev))
        bus.close()

        boot = _bootstrap()
        bus.publish(boot)
        time.sleep(0.05)

        assert len(received) == 0

    def test_overflow_drops(self):
        store = InMemoryStore()
        bus = EventBus(store, buffer_size=2)
        blocker = threading.Event()
        received = []

        def slow_handler(ev):
            blocker.wait(timeout=5.0)
            received.append(ev)

        bus.subscribe(SubscriptionPattern("*"), slow_handler)

        boot = _bootstrap()
        store.append(boot)
        # Fill buffer (2) + overflow
        for _ in range(5):
            bus.publish(boot)

        # Unblock
        blocker.set()
        time.sleep(0.2)
        bus.close()

        # Should have received at most buffer_size events
        assert len(received) <= 2

    def test_handler_exception_doesnt_kill_delivery(self):
        store = InMemoryStore()
        bus = EventBus(store)
        received = []
        latch = threading.Event()
        call_count = 0

        def bad_handler(ev):
            nonlocal call_count
            call_count += 1
            if call_count == 1:
                raise ValueError("boom")
            received.append(ev)
            latch.set()

        bus.subscribe(SubscriptionPattern("*"), bad_handler)

        boot = _bootstrap()
        store.append(boot)
        e1 = _event(boot)
        store.append(e1)

        bus.publish(boot)  # will raise
        bus.publish(e1)  # should still deliver

        latch.wait(timeout=2.0)
        bus.close()

        assert len(received) == 1

    def test_store_property(self):
        store = InMemoryStore()
        bus = EventBus(store)
        assert bus.store is store
        bus.close()
