namespace EventGraph.Tests;

public class GraphTests
{
    private static readonly ActorId SystemActor = new("system_bootstrap");

    private static Graph CreateGraph(
        InMemoryStore? store = null,
        IActorStore? actorStore = null,
        ITrustModel? trustModel = null,
        IAuthorityChain? authorityChain = null,
        ISigner? signer = null,
        GraphConfig? config = null)
    {
        return new Graph(
            store ?? new InMemoryStore(),
            actorStore ?? new InMemoryActorStore(),
            trustModel,
            authorityChain,
            signer,
            config);
    }

    private static Graph StartedGraph(
        InMemoryStore? store = null,
        IActorStore? actorStore = null)
    {
        var g = CreateGraph(store: store, actorStore: actorStore);
        g.Start();
        return g;
    }

    private static Graph BootstrappedGraph(
        InMemoryStore? store = null,
        IActorStore? actorStore = null)
    {
        var g = StartedGraph(store: store, actorStore: actorStore);
        g.Bootstrap(SystemActor);
        return g;
    }

    // ── 1. Start / Close lifecycle ─────────────────────────────────────

    [Fact]
    public void StartAndClose()
    {
        var g = CreateGraph();
        g.Start();
        g.Close();
        // Double close is safe
        g.Close();
    }

    [Fact]
    public void StartIsIdempotent()
    {
        var g = CreateGraph();
        g.Start();
        g.Start(); // no exception
        g.Close();
    }

    // ── 2. Record before Start throws ──────────────────────────────────

    [Fact]
    public void RecordBeforeStartThrows()
    {
        var g = CreateGraph();
        Assert.Throws<InvalidOperationException>(() =>
            g.Record(
                new EventType("test.event"),
                new ActorId("alice"),
                new Dictionary<string, object?> { ["key"] = "val" },
                new List<EventId> { EventFactory.NewEventId() },
                new ConversationId("conv_1")));
    }

    // ── 3. Record after Close throws ───────────────────────────────────

    [Fact]
    public void RecordAfterCloseThrows()
    {
        var g = BootstrappedGraph();
        g.Close();
        Assert.Throws<InvalidOperationException>(() =>
            g.Record(
                new EventType("test.event"),
                new ActorId("alice"),
                new Dictionary<string, object?> { ["key"] = "val" },
                new List<EventId> { EventFactory.NewEventId() },
                new ConversationId("conv_1")));
    }

    // ── 4. Bootstrap creates genesis event ─────────────────────────────

    [Fact]
    public void BootstrapCreatesGenesisEvent()
    {
        var g = StartedGraph();
        var genesis = g.Bootstrap(SystemActor);

        Assert.Equal(new EventType("system.bootstrapped"), genesis.Type);
        Assert.Equal(SystemActor, genesis.Source);
        Assert.True(genesis.PrevHash.IsZero);
        Assert.Equal(1, g.Store.Count());
    }

    // ── 5. Double bootstrap throws ─────────────────────────────────────

    [Fact]
    public void DoubleBootstrapThrows()
    {
        var g = BootstrappedGraph();
        Assert.Throws<InvalidOperationException>(() => g.Bootstrap(SystemActor));
    }

    // ── 6. Record creates event with correct chain ─────────────────────

    [Fact]
    public void RecordCreatesChainedEvent()
    {
        var g = BootstrappedGraph();
        var genesis = g.Store.Head().Unwrap();

        var ev = g.Record(
            new EventType("trust.updated"),
            new ActorId("alice"),
            new Dictionary<string, object?> { ["score"] = 0.5 },
            new List<EventId> { genesis.Id },
            new ConversationId("conv_1"));

        Assert.Equal(new EventType("trust.updated"), ev.Type);
        Assert.Equal(genesis.Hash, ev.PrevHash);
        Assert.Equal(2, g.Store.Count());
    }

    // ── 7. Query Recent ────────────────────────────────────────────────

    [Fact]
    public void QueryRecent()
    {
        var g = BootstrappedGraph();
        var genesis = g.Store.Head().Unwrap();

        g.Record(
            new EventType("test.a"), new ActorId("alice"),
            new Dictionary<string, object?>(), new List<EventId> { genesis.Id },
            new ConversationId("conv_1"));

        var q = g.Query();
        var recent = q.Recent(10);
        Assert.Equal(2, recent.Count);
    }

    // ── 8. Query ByType ────────────────────────────────────────────────

    [Fact]
    public void QueryByType()
    {
        var g = BootstrappedGraph();
        var genesis = g.Store.Head().Unwrap();

        g.Record(
            new EventType("test.alpha"), new ActorId("alice"),
            new Dictionary<string, object?>(), new List<EventId> { genesis.Id },
            new ConversationId("conv_1"));

        var q = g.Query();
        var byType = q.ByType(new EventType("test.alpha"), 10);
        Assert.Single(byType);
        Assert.Equal(new EventType("test.alpha"), byType[0].Type);

        var byBoot = q.ByType(new EventType("system.bootstrapped"), 10);
        Assert.Single(byBoot);
    }

    // ── 9. Query BySource ──────────────────────────────────────────────

    [Fact]
    public void QueryBySource()
    {
        var g = BootstrappedGraph();
        var genesis = g.Store.Head().Unwrap();

        g.Record(
            new EventType("test.alpha"), new ActorId("bob"),
            new Dictionary<string, object?>(), new List<EventId> { genesis.Id },
            new ConversationId("conv_1"));

        var q = g.Query();
        var byBob = q.BySource(new ActorId("bob"), 10);
        Assert.Single(byBob);

        var bySystem = q.BySource(SystemActor, 10);
        Assert.Single(bySystem);
    }

    // ── 10. Query EventCount ───────────────────────────────────────────

    [Fact]
    public void QueryEventCount()
    {
        var g = BootstrappedGraph();
        var q = g.Query();
        Assert.Equal(1, q.EventCount());
    }

    // ── 11. Evaluate returns authority result ───────────────────────────

    [Fact]
    public void EvaluateReturnsAuthorityResult()
    {
        var g = BootstrappedGraph();
        var actor = new Actor(
            new ActorId("alice"),
            new PublicKey(new byte[32]),
            "Alice",
            ActorType.Human,
            null,
            0,
            ActorStatus.Active);

        var result = g.Evaluate(actor, "test.action");

        // Default authority chain returns Notification for unmatched actions
        Assert.Equal(AuthorityLevel.Notification, result.Level);
    }

    // ── 12. Dispose delegates to Close ─────────────────────────────────

    [Fact]
    public void DisposeIsClose()
    {
        var g = BootstrappedGraph();
        g.Dispose();

        // After dispose, operations should throw as closed
        Assert.Throws<InvalidOperationException>(() =>
            g.Record(
                new EventType("test.event"),
                new ActorId("alice"),
                new Dictionary<string, object?>(),
                new List<EventId> { EventFactory.NewEventId() },
                new ConversationId("conv_1")));
    }
}
