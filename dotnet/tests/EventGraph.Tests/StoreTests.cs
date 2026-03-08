namespace EventGraph.Tests;

public class InMemoryStoreTests
{
    private static Event Bootstrap() => EventFactory.CreateBootstrap(new ActorId("alice"), new NoopSigner());

    private static Event NextEvent(Event prev) => EventFactory.CreateEvent(
        new EventType("trust.updated"), new ActorId("alice"),
        new Dictionary<string, object?> { ["score"] = 0.5 },
        new List<EventId> { prev.Id },
        new ConversationId("conv_1"), prev.Hash, new NoopSigner());

    [Fact]
    public void AppendAndGet()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap();
        store.Append(boot);
        var retrieved = store.Get(boot.Id);
        Assert.Equal(boot.Id, retrieved.Id);
    }

    [Fact]
    public void HeadEmpty()
    {
        var store = new InMemoryStore();
        Assert.True(store.Head().IsNone);
    }

    [Fact]
    public void HeadAfterAppend()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap();
        store.Append(boot);
        Assert.Equal(boot.Id, store.Head().Unwrap().Id);
    }

    [Fact]
    public void Count()
    {
        var store = new InMemoryStore();
        Assert.Equal(0, store.Count());
        store.Append(Bootstrap());
        Assert.Equal(1, store.Count());
    }

    [Fact]
    public void ChainOfEvents()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap();
        store.Append(boot);
        var e1 = NextEvent(boot);
        store.Append(e1);
        var e2 = NextEvent(e1);
        store.Append(e2);
        Assert.Equal(3, store.Count());
        Assert.Equal(e2.Id, store.Head().Unwrap().Id);
    }

    [Fact]
    public void RejectsBrokenChain()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap();
        store.Append(boot);

        var bad = EventFactory.CreateEvent(
            new EventType("trust.updated"), new ActorId("alice"),
            new Dictionary<string, object?>(), new List<EventId> { boot.Id },
            new ConversationId("conv_1"), boot.PrevHash, new NoopSigner()); // wrong prev_hash

        Assert.Throws<ChainIntegrityException>(() => store.Append(bad));
    }

    [Fact]
    public void GetNonexistent()
    {
        var store = new InMemoryStore();
        Assert.Throws<EventNotFoundException>(() =>
            store.Get(new EventId("019462a0-0000-7000-8000-000000000099")));
    }

    [Fact]
    public void VerifyChainEmpty()
    {
        var store = new InMemoryStore();
        var v = store.VerifyChain();
        Assert.True(v.Valid);
        Assert.Equal(0, v.Length);
    }

    [Fact]
    public void VerifyChainValid()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap();
        store.Append(boot);
        store.Append(NextEvent(boot));
        var v = store.VerifyChain();
        Assert.True(v.Valid);
        Assert.Equal(2, v.Length);
    }

    [Fact]
    public void Recent()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap();
        store.Append(boot);
        var e1 = NextEvent(boot);
        store.Append(e1);
        var e2 = NextEvent(e1);
        store.Append(e2);

        var recent = store.Recent(2);
        Assert.Equal(2, recent.Count);
        Assert.Equal(e2.Id, recent[0].Id); // newest first
        Assert.Equal(e1.Id, recent[1].Id);
    }
}
