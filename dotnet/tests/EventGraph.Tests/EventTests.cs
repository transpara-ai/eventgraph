namespace EventGraph.Tests;

public class CanonicalFormTests
{
    [Fact]
    public void SortedKeysNoWhitespace()
    {
        var content = new Dictionary<string, object?> { ["b"] = 1, ["a"] = 2 };
        var json = CanonicalForm.CanonicalContentJson(content);
        Assert.DoesNotContain(" ", json);
        Assert.StartsWith("{\"a\"", json);
    }

    [Fact]
    public void PipeSeparatedFormat()
    {
        var canon = CanonicalForm.Build(
            1, new string('0', 64),
            new[] { "c2", "c1" },
            "eid", "trust.updated", "actor_alice", "conv_1",
            123456789, "{\"key\":\"val\"}");

        var parts = canon.Split('|');
        Assert.Equal("1", parts[0]);
        Assert.Equal(new string('0', 64), parts[1]);
        Assert.Equal("c1,c2", parts[2]); // sorted
        Assert.Equal("eid", parts[3]);
        Assert.Equal("trust.updated", parts[4]);
    }

    [Fact]
    public void EmptyCauses()
    {
        var canon = CanonicalForm.Build(1, "", Array.Empty<string>(), "eid", "system.bootstrapped", "s", "c", 0, "{}");
        var parts = canon.Split('|');
        Assert.Equal("", parts[2]);
    }
}

public class ComputeHashTests
{
    [Fact]
    public void Deterministic()
    {
        var h1 = CanonicalForm.ComputeHash("hello");
        var h2 = CanonicalForm.ComputeHash("hello");
        Assert.Equal(h1, h2);
    }

    [Fact]
    public void DifferentInput()
    {
        var h1 = CanonicalForm.ComputeHash("hello");
        var h2 = CanonicalForm.ComputeHash("world");
        Assert.NotEqual(h1, h2);
    }

    [Fact]
    public void Returns64Hex()
    {
        var h = CanonicalForm.ComputeHash("test");
        Assert.Equal(64, h.Value.Length);
    }
}

public class EventFactoryTests
{
    [Fact]
    public void NewEventIdGeneratesValidV7()
    {
        var eid = EventFactory.NewEventId();
        Assert.Equal(36, eid.Value.Length);
        Assert.Equal('7', eid.Value[14]);
    }

    [Fact]
    public void CreateBootstrapValid()
    {
        var signer = new NoopSigner();
        var source = new ActorId("actor_alice");
        var ev = EventFactory.CreateBootstrap(source, signer);

        Assert.Equal(1, ev.Version);
        Assert.Equal("system.bootstrapped", ev.Type.Value);
        Assert.Equal("actor_alice", ev.Source.Value);
        Assert.True(ev.PrevHash.IsZero);
        Assert.Equal(1, ev.Causes.Count);
        Assert.Equal(ev.Id, ev.Causes[0]); // self-referencing
    }

    [Fact]
    public void CreateEventValid()
    {
        var signer = new NoopSigner();
        var boot = EventFactory.CreateBootstrap(new ActorId("alice"), signer);

        var ev = EventFactory.CreateEvent(
            new EventType("trust.updated"),
            new ActorId("alice"),
            new Dictionary<string, object?> { ["score"] = 0.8 },
            new List<EventId> { boot.Id },
            new ConversationId("conv_1"),
            boot.Hash, signer);

        Assert.Equal("trust.updated", ev.Type.Value);
        Assert.Equal(boot.Hash, ev.PrevHash);
        Assert.Equal(1, ev.Causes.Count);
    }

    [Fact]
    public void ContentIsDefensiveCopy()
    {
        var signer = new NoopSigner();
        var ev = EventFactory.CreateBootstrap(new ActorId("alice"), signer);
        var content = ev.Content;
        Assert.False(ReferenceEquals(content, ev.Content));
    }
}
