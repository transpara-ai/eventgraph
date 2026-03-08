namespace EventGraph.Tests;

public class StubPrimitive : IPrimitive
{
    public PrimitiveId Id { get; }
    public Layer Layer { get; }
    public List<SubscriptionPattern> Subscriptions { get; } = new() { new SubscriptionPattern("*") };
    public Cadence Cadence { get; } = new(1);
    public int ReceivedCount;

    public StubPrimitive(string name, int layer = 0)
    {
        Id = new PrimitiveId(name);
        Layer = new Layer(layer);
    }

    public virtual List<Mutation> Process(int tick, List<Event> events, Snapshot snapshot)
    {
        ReceivedCount += events.Count;
        return new() { new UpdateStateMutation(Id, "count", ReceivedCount) };
    }
}

public class LifecycleTests
{
    [Fact] public void ValidTransitions() { Assert.True(Lifecycle.IsValidTransition(Lifecycle.Dormant, Lifecycle.Activating)); Assert.True(Lifecycle.IsValidTransition(Lifecycle.Active, Lifecycle.Processing)); }
    [Fact] public void InvalidTransitions() { Assert.False(Lifecycle.IsValidTransition(Lifecycle.Dormant, Lifecycle.Active)); Assert.False(Lifecycle.IsValidTransition(Lifecycle.Memorial, Lifecycle.Dormant)); }
    [Fact] public void MemorialIsTerminal() { Assert.False(Lifecycle.IsValidTransition(Lifecycle.Memorial, Lifecycle.Dormant)); Assert.False(Lifecycle.IsValidTransition(Lifecycle.Memorial, Lifecycle.Active)); }
}

public class PrimitiveRegistryTests
{
    [Fact]
    public void RegisterAndGet()
    {
        var reg = new PrimitiveRegistry();
        var p = new StubPrimitive("test");
        reg.Register(p);
        Assert.Same(p, reg.Get(new PrimitiveId("test")));
    }

    [Fact]
    public void RegisterDuplicateThrows()
    {
        var reg = new PrimitiveRegistry();
        reg.Register(new StubPrimitive("dup"));
        Assert.Throws<InvalidOperationException>(() => reg.Register(new StubPrimitive("dup")));
    }

    [Fact]
    public void Count()
    {
        var reg = new PrimitiveRegistry();
        Assert.Equal(0, reg.Count);
        reg.Register(new StubPrimitive("a"));
        Assert.Equal(1, reg.Count);
    }

    [Fact]
    public void AllOrderedByLayer()
    {
        var reg = new PrimitiveRegistry();
        reg.Register(new StubPrimitive("high", 5));
        reg.Register(new StubPrimitive("low", 0));
        reg.Register(new StubPrimitive("mid", 2));
        var all = reg.All();
        Assert.Equal("low", all[0].Id.Value);
        Assert.Equal("mid", all[1].Id.Value);
        Assert.Equal("high", all[2].Id.Value);
    }

    [Fact]
    public void LifecycleStartsDormant()
    {
        var reg = new PrimitiveRegistry();
        reg.Register(new StubPrimitive("p"));
        Assert.Equal(Lifecycle.Dormant, reg.GetLifecycle(new PrimitiveId("p")));
    }

    [Fact]
    public void Activate()
    {
        var reg = new PrimitiveRegistry();
        reg.Register(new StubPrimitive("p"));
        reg.Activate(new PrimitiveId("p"));
        Assert.Equal(Lifecycle.Active, reg.GetLifecycle(new PrimitiveId("p")));
    }

    [Fact]
    public void SetLifecycleInvalid()
    {
        var reg = new PrimitiveRegistry();
        reg.Register(new StubPrimitive("p"));
        Assert.Throws<InvalidTransitionException>(() => reg.SetLifecycle(new PrimitiveId("p"), Lifecycle.Active));
    }

    [Fact]
    public void UpdateAndReadState()
    {
        var reg = new PrimitiveRegistry();
        reg.Register(new StubPrimitive("p"));
        reg.UpdateState(new PrimitiveId("p"), "count", 42);
        var states = reg.AllStates();
        Assert.Equal(42, states["p"].State["count"]);
    }

    [Fact]
    public void GetNonexistent()
    {
        var reg = new PrimitiveRegistry();
        Assert.Null(reg.Get(new PrimitiveId("nope")));
    }
}
