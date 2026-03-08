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

    private static Event NextEventWithType(Event prev, EventType type, ActorId source, ConversationId convId) =>
        EventFactory.CreateEvent(
            type, source,
            new Dictionary<string, object?> { ["score"] = 0.5 },
            new List<EventId> { prev.Id },
            convId, prev.Hash, new NoopSigner());

    [Fact]
    public void ByTypeFiltersCorrectly()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap();
        store.Append(boot);
        var trustType = new EventType("trust.updated");
        var otherType = new EventType("edge.created");

        var e1 = NextEventWithType(boot, trustType, new ActorId("alice"), new ConversationId("conv_1"));
        store.Append(e1);
        var e2 = NextEventWithType(e1, otherType, new ActorId("alice"), new ConversationId("conv_1"));
        store.Append(e2);
        var e3 = NextEventWithType(e2, trustType, new ActorId("alice"), new ConversationId("conv_1"));
        store.Append(e3);

        var results = store.ByType(trustType, 10);
        Assert.Equal(2, results.Count);
        Assert.Equal(e3.Id, results[0].Id); // newest first
        Assert.Equal(e1.Id, results[1].Id);
    }

    [Fact]
    public void ByTypeRespectsLimit()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap();
        store.Append(boot);
        var trustType = new EventType("trust.updated");

        var e1 = NextEventWithType(boot, trustType, new ActorId("alice"), new ConversationId("conv_1"));
        store.Append(e1);
        var e2 = NextEventWithType(e1, trustType, new ActorId("alice"), new ConversationId("conv_1"));
        store.Append(e2);

        var results = store.ByType(trustType, 1);
        Assert.Single(results);
        Assert.Equal(e2.Id, results[0].Id); // newest first
    }

    [Fact]
    public void BySourceFiltersCorrectly()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap(); // source: alice
        store.Append(boot);

        var e1 = NextEventWithType(boot, new EventType("trust.updated"), new ActorId("alice"), new ConversationId("conv_1"));
        store.Append(e1);
        var e2 = NextEventWithType(e1, new EventType("trust.updated"), new ActorId("bob"), new ConversationId("conv_1"));
        store.Append(e2);

        var results = store.BySource(new ActorId("bob"), 10);
        Assert.Single(results);
        Assert.Equal(e2.Id, results[0].Id);
    }

    [Fact]
    public void ByConversationFiltersCorrectly()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap(); // conv: conv_alice
        store.Append(boot);

        var conv1 = new ConversationId("conv_1");
        var conv2 = new ConversationId("conv_2");

        var e1 = NextEventWithType(boot, new EventType("trust.updated"), new ActorId("alice"), conv1);
        store.Append(e1);
        var e2 = NextEventWithType(e1, new EventType("trust.updated"), new ActorId("alice"), conv2);
        store.Append(e2);
        var e3 = NextEventWithType(e2, new EventType("trust.updated"), new ActorId("alice"), conv1);
        store.Append(e3);

        var results = store.ByConversation(conv1, 10);
        Assert.Equal(2, results.Count);
        Assert.Equal(e3.Id, results[0].Id); // newest first
        Assert.Equal(e1.Id, results[1].Id);
    }

    [Fact]
    public void AncestorsTraversesCausalChain()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap();
        store.Append(boot);
        var e1 = NextEvent(boot);
        store.Append(e1);
        var e2 = NextEvent(e1);
        store.Append(e2);

        var ancestors = store.Ancestors(e2.Id, 10);
        Assert.Equal(2, ancestors.Count); // e1 and boot
        Assert.Contains(ancestors, e => e.Id == e1.Id);
        Assert.Contains(ancestors, e => e.Id == boot.Id);
    }

    [Fact]
    public void AncestorsRespectsMaxDepth()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap();
        store.Append(boot);
        var e1 = NextEvent(boot);
        store.Append(e1);
        var e2 = NextEvent(e1);
        store.Append(e2);

        var ancestors = store.Ancestors(e2.Id, 1);
        Assert.Single(ancestors); // only e1, depth=1 doesn't reach boot
        Assert.Equal(e1.Id, ancestors[0].Id);
    }

    [Fact]
    public void AncestorsDoesNotIncludeStartEvent()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap();
        store.Append(boot);

        var ancestors = store.Ancestors(boot.Id, 10);
        Assert.DoesNotContain(ancestors, e => e.Id == boot.Id);
    }

    [Fact]
    public void DescendantsTraversesCausalChain()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap();
        store.Append(boot);
        var e1 = NextEvent(boot);
        store.Append(e1);
        var e2 = NextEvent(e1);
        store.Append(e2);

        var descendants = store.Descendants(boot.Id, 10);
        Assert.Equal(2, descendants.Count); // e1 and e2
        Assert.Contains(descendants, e => e.Id == e1.Id);
        Assert.Contains(descendants, e => e.Id == e2.Id);
    }

    [Fact]
    public void DescendantsRespectsMaxDepth()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap();
        store.Append(boot);
        var e1 = NextEvent(boot);
        store.Append(e1);
        var e2 = NextEvent(e1);
        store.Append(e2);

        var descendants = store.Descendants(boot.Id, 1);
        Assert.Single(descendants); // only e1, depth=1 doesn't reach e2
        Assert.Equal(e1.Id, descendants[0].Id);
    }

    [Fact]
    public void DescendantsDoesNotIncludeStartEvent()
    {
        var store = new InMemoryStore();
        var boot = Bootstrap();
        store.Append(boot);
        var e1 = NextEvent(boot);
        store.Append(e1);

        var descendants = store.Descendants(boot.Id, 10);
        Assert.DoesNotContain(descendants, e => e.Id == boot.Id);
    }
}
