namespace EventGraph.Tests;

public class PrimitivesTests
{
    [Fact]
    public void CreateAll_Returns201Primitives()
    {
        var all = PrimitiveFactory.CreateAll();
        Assert.Equal(201, all.Count);
    }

    [Fact]
    public void AllPrimitives_Instantiate()
    {
        var all = PrimitiveFactory.CreateAll();
        foreach (var p in all)
        {
            Assert.NotNull(p);
            Assert.False(string.IsNullOrEmpty(p.Id.Value));
        }
    }

    [Fact]
    public void AllPrimitives_HaveNonEmptySubscriptions()
    {
        var all = PrimitiveFactory.CreateAll();
        foreach (var p in all)
        {
            Assert.True(p.Subscriptions.Count > 0, $"Primitive '{p.Id.Value}' has no subscriptions");
        }
    }

    [Fact]
    public void AllPrimitives_HaveUniqueNames()
    {
        var all = PrimitiveFactory.CreateAll();
        var names = all.Select(p => p.Id.Value).ToList();
        Assert.Equal(names.Count, names.Distinct().Count());
    }

    [Fact]
    public void Layer0_Has45Primitives()
    {
        var all = PrimitiveFactory.CreateAll();
        var count = all.Count(p => p.Layer.Value == 0);
        Assert.Equal(45, count);
    }

    [Fact]
    public void Layer1_Has12Primitives()
    {
        var all = PrimitiveFactory.CreateAll();
        Assert.Equal(12, all.Count(p => p.Layer.Value == 1));
    }

    [Fact]
    public void Layer2_Has12Primitives()
    {
        var all = PrimitiveFactory.CreateAll();
        Assert.Equal(12, all.Count(p => p.Layer.Value == 2));
    }

    [Fact]
    public void Layer3_Has12Primitives()
    {
        var all = PrimitiveFactory.CreateAll();
        Assert.Equal(12, all.Count(p => p.Layer.Value == 3));
    }

    [Fact]
    public void Layer4_Has12Primitives()
    {
        var all = PrimitiveFactory.CreateAll();
        Assert.Equal(12, all.Count(p => p.Layer.Value == 4));
    }

    [Fact]
    public void Layer5_Has12Primitives()
    {
        var all = PrimitiveFactory.CreateAll();
        Assert.Equal(12, all.Count(p => p.Layer.Value == 5));
    }

    [Fact]
    public void Layer6_Has12Primitives()
    {
        var all = PrimitiveFactory.CreateAll();
        Assert.Equal(12, all.Count(p => p.Layer.Value == 6));
    }

    [Fact]
    public void Layer7_Has12Primitives()
    {
        var all = PrimitiveFactory.CreateAll();
        Assert.Equal(12, all.Count(p => p.Layer.Value == 7));
    }

    [Fact]
    public void Layer8_Has12Primitives()
    {
        var all = PrimitiveFactory.CreateAll();
        Assert.Equal(12, all.Count(p => p.Layer.Value == 8));
    }

    [Fact]
    public void Layer9_Has12Primitives()
    {
        var all = PrimitiveFactory.CreateAll();
        Assert.Equal(12, all.Count(p => p.Layer.Value == 9));
    }

    [Fact]
    public void Layer10_Has12Primitives()
    {
        var all = PrimitiveFactory.CreateAll();
        Assert.Equal(12, all.Count(p => p.Layer.Value == 10));
    }

    [Fact]
    public void Layer11_Has12Primitives()
    {
        var all = PrimitiveFactory.CreateAll();
        Assert.Equal(12, all.Count(p => p.Layer.Value == 11));
    }

    [Fact]
    public void Layer12_Has12Primitives()
    {
        var all = PrimitiveFactory.CreateAll();
        Assert.Equal(12, all.Count(p => p.Layer.Value == 12));
    }

    [Fact]
    public void Layer13_Has12Primitives()
    {
        var all = PrimitiveFactory.CreateAll();
        Assert.Equal(12, all.Count(p => p.Layer.Value == 13));
    }

    [Fact]
    public void Process_ReturnsMutations()
    {
        var all = PrimitiveFactory.CreateAll();
        var snapshot = new Snapshot(1, new Dictionary<string, PrimitiveState>(), new List<Event>(), new List<Event>());

        foreach (var p in all)
        {
            var mutations = p.Process(1, new List<Event>(), snapshot);
            Assert.NotNull(mutations);
            Assert.True(mutations.Count >= 2, $"Primitive '{p.Id.Value}' should return at least 2 mutations (eventsProcessed, lastTick)");

            // Verify the mutations are UpdateStateMutation with expected keys
            var stateMutations = mutations.OfType<UpdateStateMutation>().ToList();
            Assert.Contains(stateMutations, m => m.Key == "eventsProcessed");
            Assert.Contains(stateMutations, m => m.Key == "lastTick");
        }
    }

    [Fact]
    public void Process_TracksEventCount()
    {
        var p = new ConcretePrimitive("TestPrimitive", 0, new List<SubscriptionPattern> { new("*") });
        var snapshot = new Snapshot(5, new Dictionary<string, PrimitiveState>(), new List<Event>(), new List<Event>());

        // Create some events to process
        var signer = new NoopSigner();
        var bootstrap = EventFactory.CreateBootstrap(new ActorId("test-actor"), signer);
        var events = new List<Event> { bootstrap, bootstrap, bootstrap };

        var mutations = p.Process(5, events, snapshot);
        var eventsProcessed = mutations.OfType<UpdateStateMutation>().First(m => m.Key == "eventsProcessed");
        var lastTick = mutations.OfType<UpdateStateMutation>().First(m => m.Key == "lastTick");

        Assert.Equal(3, eventsProcessed.Value);
        Assert.Equal(5, lastTick.Value);
    }

    [Fact]
    public void CreateRegistry_RegistersAll201()
    {
        var registry = PrimitiveFactory.CreateRegistry();
        Assert.Equal(201, registry.Count);
    }

    [Fact]
    public void CreateRegistry_AllRetrievable()
    {
        var registry = PrimitiveFactory.CreateRegistry();
        var all = registry.All();
        Assert.Equal(201, all.Count);

        // Spot-check some primitives from different layers
        Assert.NotNull(registry.Get(new PrimitiveId("Event")));
        Assert.NotNull(registry.Get(new PrimitiveId("Goal")));
        Assert.NotNull(registry.Get(new PrimitiveId("Message")));
        Assert.NotNull(registry.Get(new PrimitiveId("Group")));
        Assert.NotNull(registry.Get(new PrimitiveId("Rule")));
        Assert.NotNull(registry.Get(new PrimitiveId("Create")));
        Assert.NotNull(registry.Get(new PrimitiveId("Symbol")));
        Assert.NotNull(registry.Get(new PrimitiveId("Value")));
        Assert.NotNull(registry.Get(new PrimitiveId("SelfModel")));
        Assert.NotNull(registry.Get(new PrimitiveId("Attachment")));
        Assert.NotNull(registry.Get(new PrimitiveId("Home")));
        Assert.NotNull(registry.Get(new PrimitiveId("SelfAwareness")));
        Assert.NotNull(registry.Get(new PrimitiveId("MetaPattern")));
        Assert.NotNull(registry.Get(new PrimitiveId("Being")));
        Assert.NotNull(registry.Get(new PrimitiveId("Wonder")));
    }

    [Fact]
    public void AllPrimitives_CadenceIsValid()
    {
        var all = PrimitiveFactory.CreateAll();
        foreach (var p in all)
        {
            Assert.True(p.Cadence.Value >= 1, $"Primitive '{p.Id.Value}' has invalid cadence {p.Cadence.Value}");
        }
    }

    [Fact]
    public void AllPrimitives_LayerInRange()
    {
        var all = PrimitiveFactory.CreateAll();
        foreach (var p in all)
        {
            Assert.InRange(p.Layer.Value, 0, 13);
        }
    }

    [Fact]
    public void Registry_OrderedByLayerThenName()
    {
        var registry = PrimitiveFactory.CreateRegistry();
        var all = registry.All();

        for (int i = 1; i < all.Count; i++)
        {
            var prev = all[i - 1];
            var curr = all[i];
            if (prev.Layer.Value == curr.Layer.Value)
                Assert.True(string.Compare(prev.Id.Value, curr.Id.Value, StringComparison.Ordinal) <= 0,
                    $"'{prev.Id.Value}' should come before '{curr.Id.Value}' within layer {curr.Layer.Value}");
            else
                Assert.True(prev.Layer.Value < curr.Layer.Value,
                    $"Layer {prev.Layer.Value} should come before layer {curr.Layer.Value}");
        }
    }
}
