namespace EventGraph.Tests;

public class EmittingPrimitive : IPrimitive
{
    public PrimitiveId Id { get; }
    public Layer Layer { get; } = new(0);
    public List<SubscriptionPattern> Subscriptions { get; } = new() { new SubscriptionPattern("*") };
    public Cadence Cadence { get; } = new(1);
    public int Emissions;
    private readonly int _maxEmissions;

    public EmittingPrimitive(string name, int maxEmissions = 1)
    {
        Id = new PrimitiveId(name);
        _maxEmissions = maxEmissions;
    }

    public List<Mutation> Process(int tick, List<Event> events, Snapshot snapshot)
    {
        if (events.Count == 0 || Emissions >= _maxEmissions) return new();
        Emissions++;
        return new()
        {
            new AddEventMutation(new EventType("test.emitted"), new ActorId("emitter"),
                new Dictionary<string, object?> { ["wave"] = Emissions },
                new List<EventId> { events[0].Id },
                new ConversationId("conv_tick"))
        };
    }
}

public class TickEngineTests
{
    private static (PrimitiveRegistry, InMemoryStore, TickEngine, Event) Setup(
        IPrimitive[]? prims = null, TickConfig? config = null)
    {
        var reg = new PrimitiveRegistry();
        var store = new InMemoryStore();
        var boot = EventFactory.CreateBootstrap(new ActorId("system"), new NoopSigner());
        store.Append(boot);
        foreach (var p in prims ?? [])
        {
            reg.Register(p);
            reg.Activate(p.Id);
        }
        var engine = new TickEngine(reg, store, config);
        return (reg, store, engine, boot);
    }

    [Fact]
    public void BasicTick()
    {
        var counter = new StubPrimitive("counter");
        var (_, _, engine, boot) = Setup(new IPrimitive[] { counter });
        var result = engine.Tick(new() { boot });
        Assert.Equal(1, result.Tick);
        Assert.True(result.Mutations >= 1);
        Assert.Equal(1, counter.ReceivedCount);
    }

    [Fact]
    public void Quiescence()
    {
        var counter = new StubPrimitive("counter");
        var (_, _, engine, boot) = Setup(new IPrimitive[] { counter });
        var result = engine.Tick(new() { boot });
        Assert.True(result.Quiesced);
    }

    [Fact]
    public void RippleWaves()
    {
        var emitter = new EmittingPrimitive("emitter", 3);
        var counter = new StubPrimitive("counter");
        var (_, _, engine, boot) = Setup(new IPrimitive[] { emitter, counter });
        var result = engine.Tick(new() { boot });
        Assert.True(result.Waves > 1);
        Assert.True(counter.ReceivedCount > 1);
    }

    [Fact]
    public void MaxWavesLimit()
    {
        var infinite = new InfiniteEmitter();
        var config = new TickConfig(MaxWavesPerTick: 3);
        var (_, _, engine, boot) = Setup(new IPrimitive[] { infinite }, config);
        var result = engine.Tick(new() { boot });
        Assert.Equal(3, result.Waves);
        Assert.False(result.Quiesced);
    }

    [Fact]
    public void InactivePrimitivesSkipped()
    {
        var counter = new StubPrimitive("dormant");
        var reg = new PrimitiveRegistry();
        var store = new InMemoryStore();
        var boot = EventFactory.CreateBootstrap(new ActorId("system"), new NoopSigner());
        store.Append(boot);
        reg.Register(counter); // don't activate
        var engine = new TickEngine(reg, store);
        engine.Tick(new() { boot });
        Assert.Equal(0, counter.ReceivedCount);
    }

    [Fact]
    public void TickCounterIncrements()
    {
        var (_, _, engine, boot) = Setup();
        Assert.Equal(1, engine.Tick(new() { boot }).Tick);
        Assert.Equal(2, engine.Tick().Tick);
        Assert.Equal(3, engine.Tick().Tick);
    }

    [Fact]
    public void LayerOrdering()
    {
        var order = new List<string>();
        var high = new OrderTracker("high", 5, order);
        var low = new OrderTracker("low", 0, order);
        var mid = new OrderTracker("mid", 2, order);
        var (_, _, engine, boot) = Setup(new IPrimitive[] { high, low, mid });
        engine.Tick(new() { boot });
        Assert.Equal(new[] { "low", "mid", "high" }, order.ToArray());
    }

    private class InfiniteEmitter : IPrimitive
    {
        public PrimitiveId Id { get; } = new("infinite");
        public Layer Layer { get; } = new(0);
        public List<SubscriptionPattern> Subscriptions { get; } = new() { new SubscriptionPattern("*") };
        public Cadence Cadence { get; } = new(1);

        public List<Mutation> Process(int tick, List<Event> events, Snapshot snapshot)
        {
            if (events.Count == 0) return new();
            return new()
            {
                new AddEventMutation(new EventType("test.loop"), new ActorId("inf"),
                    new(), new List<EventId> { events[0].Id }, new ConversationId("conv_inf"))
            };
        }
    }

    private class OrderTracker : IPrimitive
    {
        public PrimitiveId Id { get; }
        public Layer Layer { get; }
        public List<SubscriptionPattern> Subscriptions { get; } = new() { new SubscriptionPattern("*") };
        public Cadence Cadence { get; } = new(1);
        private readonly List<string> _order;

        public OrderTracker(string name, int layer, List<string> order)
        {
            Id = new PrimitiveId(name);
            Layer = new Layer(layer);
            _order = order;
        }

        public List<Mutation> Process(int tick, List<Event> events, Snapshot snapshot)
        {
            _order.Add(Id.Value);
            return new();
        }
    }
}
