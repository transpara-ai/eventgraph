"""Tests for the tick engine — ripple-wave processing."""

from eventgraph.event import NoopSigner, create_bootstrap
from eventgraph.primitive import (
    LIFECYCLE_ACTIVE,
    AddEvent,
    Mutation,
    Registry,
    Snapshot,
    UpdateState,
)
from eventgraph.store import InMemoryStore
from eventgraph.tick import TickConfig, TickEngine, TickResult
from eventgraph.types import (
    Activation,
    ActorID,
    Cadence,
    ConversationID,
    EventType,
    Layer,
    PrimitiveID,
    SubscriptionPattern,
)


class CountingPrimitive:
    """Counts events it receives."""

    def __init__(self, name: str, layer: int = 0) -> None:
        self._id = PrimitiveID(name)
        self._layer = Layer(layer)
        self.received_count = 0

    def id(self) -> PrimitiveID:
        return self._id

    def layer(self) -> Layer:
        return self._layer

    def process(self, tick, events, snapshot) -> list[Mutation]:
        self.received_count += len(events)
        return [
            UpdateState(primitive_id=self._id, key="count", value=self.received_count)
        ]

    def subscriptions(self) -> list[SubscriptionPattern]:
        return [SubscriptionPattern("*")]

    def cadence(self) -> Cadence:
        return Cadence(1)


class EmittingPrimitive:
    """Emits a new event on each process call (causes ripple waves)."""

    def __init__(self, name: str, max_emissions: int = 1) -> None:
        self._id = PrimitiveID(name)
        self._layer = Layer(0)
        self.emissions = 0
        self.max_emissions = max_emissions

    def id(self) -> PrimitiveID:
        return self._id

    def layer(self) -> Layer:
        return self._layer

    def process(self, tick, events, snapshot) -> list[Mutation]:
        if not events or self.emissions >= self.max_emissions:
            return []
        self.emissions += 1
        return [
            AddEvent(
                type=EventType("test.emitted"),
                source=ActorID("emitter"),
                content={"wave": self.emissions},
                causes=[events[0].id],
                conversation_id=ConversationID("conv_tick"),
            )
        ]

    def subscriptions(self) -> list[SubscriptionPattern]:
        return [SubscriptionPattern("*")]

    def cadence(self) -> Cadence:
        return Cadence(1)


def _setup(prims=None, config=None):
    """Create a registry, store, and engine with bootstrap event."""
    reg = Registry()
    store = InMemoryStore()
    boot = create_bootstrap(source=ActorID("system"), signer=NoopSigner())
    store.append(boot)

    for p in (prims or []):
        reg.register(p)
        reg.activate(p.id())

    engine = TickEngine(reg, store, config)
    return reg, store, engine, boot


class TestTickEngine:
    def test_basic_tick(self):
        counter = CountingPrimitive("counter")
        reg, store, engine, boot = _setup([counter])

        result = engine.tick([boot])
        assert result.tick == 1
        assert result.mutations >= 1
        assert counter.received_count == 1

    def test_quiescence(self):
        counter = CountingPrimitive("counter")
        _, _, engine, boot = _setup([counter])

        result = engine.tick([boot])
        assert result.quiesced is True  # no new events emitted

    def test_ripple_waves(self):
        emitter = EmittingPrimitive("emitter", max_emissions=3)
        counter = CountingPrimitive("counter")
        reg, store, engine, boot = _setup([emitter, counter])

        result = engine.tick([boot])
        assert result.waves > 1  # ripple happened
        assert counter.received_count > 1  # received original + emitted

    def test_max_waves_limit(self):
        # Emitter that always emits (infinite ripple)
        class InfiniteEmitter:
            def __init__(self):
                self._id = PrimitiveID("infinite")

            def id(self): return self._id
            def layer(self): return Layer(0)
            def subscriptions(self): return [SubscriptionPattern("*")]
            def cadence(self): return Cadence(1)

            def process(self, tick, events, snapshot):
                if events:
                    return [AddEvent(
                        type=EventType("test.loop"),
                        source=ActorID("inf"),
                        content={},
                        causes=[events[0].id],
                        conversation_id=ConversationID("conv_inf"),
                    )]
                return []

        config = TickConfig(max_waves_per_tick=3)
        _, _, engine, boot = _setup([InfiniteEmitter()], config)

        result = engine.tick([boot])
        assert result.waves == 3
        assert result.quiesced is False

    def test_inactive_primitives_skipped(self):
        counter = CountingPrimitive("dormant_counter")
        reg = Registry()
        store = InMemoryStore()
        boot = create_bootstrap(source=ActorID("system"), signer=NoopSigner())
        store.append(boot)
        reg.register(counter)
        # Don't activate — stays dormant

        engine = TickEngine(reg, store)
        result = engine.tick([boot])
        assert counter.received_count == 0

    def test_cadence_respected(self):
        counter = CountingPrimitive("slow")

        class SlowPrimitive(CountingPrimitive):
            def cadence(self): return Cadence(3)

        slow = SlowPrimitive("slow_prim")
        reg, store, engine, boot = _setup([slow])

        r1 = engine.tick([boot])  # tick 1: should process (last=0, cadence=3, 1-0>=3? no wait...)
        # Actually cadence check: tick - last_tick < cadence → skip
        # tick=1, last=0, 1-0=1 < 3 → skip... but we need first invocation
        # The Go engine invokes on first tick regardless. Let me check:
        # Our check is: tick_num - last < cadence → skip
        # 1 - 0 = 1 < 3 → skip. This means cadence=3 skips ticks 1,2 and runs on 3.
        # That's correct — cadence means "minimum ticks between invocations"
        assert slow.received_count == 0

        engine.tick([])  # tick 2
        assert slow.received_count == 0

        engine.tick([boot])  # tick 3: 3-0=3, not < 3, should process
        assert slow.received_count == 1

    def test_tick_counter_increments(self):
        _, _, engine, boot = _setup()
        r1 = engine.tick([boot])
        r2 = engine.tick([])
        r3 = engine.tick([])
        assert r1.tick == 1
        assert r2.tick == 2
        assert r3.tick == 3

    def test_published_events(self):
        emitter = EmittingPrimitive("pub_emitter", max_emissions=1)
        published: list = []

        reg = Registry()
        store = InMemoryStore()
        boot = create_bootstrap(source=ActorID("system"), signer=NoopSigner())
        store.append(boot)
        reg.register(emitter)
        reg.activate(emitter.id())

        engine = TickEngine(reg, store, publisher=lambda ev: published.append(ev))
        engine.tick([boot])
        assert len(published) >= 1

    def test_layer_ordering(self):
        """Lower layers process before higher layers."""
        order = []

        class OrderTracker:
            def __init__(self, name, layer_val):
                self._id = PrimitiveID(name)
                self._layer = Layer(layer_val)

            def id(self): return self._id
            def layer(self): return self._layer
            def subscriptions(self): return [SubscriptionPattern("*")]
            def cadence(self): return Cadence(1)

            def process(self, tick, events, snapshot):
                order.append(self._id.value)
                return []

        p_high = OrderTracker("high", 5)
        p_low = OrderTracker("low", 0)
        p_mid = OrderTracker("mid", 2)

        reg, store, engine, boot = _setup([p_high, p_low, p_mid])
        engine.tick([boot])

        assert order == ["low", "mid", "high"]
