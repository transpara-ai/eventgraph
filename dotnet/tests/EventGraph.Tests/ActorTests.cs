namespace EventGraph.Tests;

public class ActorTests
{
    private static readonly EventId TestReason = new("019462a0-0000-7000-8000-000000000001");

    private static PublicKey TestPublicKey(byte b)
    {
        var key = new byte[32];
        key[0] = b;
        return new PublicKey(key);
    }

    // ── Register ────────────────────────────────────────────────────────

    [Fact]
    public void Register()
    {
        var store = new InMemoryActorStore();
        var pk = TestPublicKey(1);

        var actor = store.Register(pk, "Alice", ActorType.Human);

        Assert.Equal("Alice", actor.DisplayName);
        Assert.Equal(ActorType.Human, actor.Type);
        Assert.Equal(ActorStatus.Active, actor.Status);
    }

    [Fact]
    public void RegisterIdempotent()
    {
        var store = new InMemoryActorStore();
        var pk = TestPublicKey(1);

        var a1 = store.Register(pk, "Alice", ActorType.Human);
        var a2 = store.Register(pk, "Alice Again", ActorType.Human);

        Assert.Equal(a1.Id, a2.Id);
        Assert.Equal(1, store.ActorCount);
    }

    // ── Get ─────────────────────────────────────────────────────────────

    [Fact]
    public void Get()
    {
        var store = new InMemoryActorStore();
        var pk = TestPublicKey(1);
        var actor = store.Register(pk, "Alice", ActorType.Human);

        var got = store.Get(actor.Id);

        Assert.Equal(actor.Id, got.Id);
    }

    [Fact]
    public void GetNotFound()
    {
        var store = new InMemoryActorStore();

        var ex = Assert.Throws<ActorNotFoundException>(() =>
            store.Get(new ActorId("actor_nonexistent")));
        Assert.Equal("actor_nonexistent", ex.ActorId.Value);
    }

    // ── GetByPublicKey ──────────────────────────────────────────────────

    [Fact]
    public void GetByPublicKey()
    {
        var store = new InMemoryActorStore();
        var pk = TestPublicKey(1);
        var actor = store.Register(pk, "Alice", ActorType.Human);

        var got = store.GetByPublicKey(pk);

        Assert.Equal(actor.Id, got.Id);
    }

    [Fact]
    public void GetByPublicKeyNotFound()
    {
        var store = new InMemoryActorStore();

        Assert.Throws<ActorKeyNotFoundException>(() =>
            store.GetByPublicKey(TestPublicKey(99)));
    }

    // ── Update ──────────────────────────────────────────────────────────

    [Fact]
    public void Update()
    {
        var store = new InMemoryActorStore();
        var pk = TestPublicKey(1);
        var actor = store.Register(pk, "Alice", ActorType.Human);

        var updated = store.Update(actor.Id, new ActorUpdate { DisplayName = "Alice Updated" });

        Assert.Equal("Alice Updated", updated.DisplayName);
    }

    [Fact]
    public void UpdateMetadataMerge()
    {
        var store = new InMemoryActorStore();
        var pk = TestPublicKey(1);
        var actor = store.Register(pk, "Alice", ActorType.Human);

        // Set initial metadata
        store.Update(actor.Id, new ActorUpdate
        {
            Metadata = new Dictionary<string, object?> { ["role"] = "builder" }
        });

        // Merge additional metadata
        var updated = store.Update(actor.Id, new ActorUpdate
        {
            Metadata = new Dictionary<string, object?> { ["team"] = "core" }
        });

        var md = updated.Metadata;
        Assert.Equal("builder", md["role"]);
        Assert.Equal("core", md["team"]);
    }

    [Fact]
    public void UpdateNotFound()
    {
        var store = new InMemoryActorStore();

        Assert.Throws<ActorNotFoundException>(() =>
            store.Update(new ActorId("actor_nonexistent"), new ActorUpdate()));
    }

    // ── Suspend ─────────────────────────────────────────────────────────

    [Fact]
    public void Suspend()
    {
        var store = new InMemoryActorStore();
        var pk = TestPublicKey(1);
        var actor = store.Register(pk, "Alice", ActorType.Human);

        var suspended = store.Suspend(actor.Id, TestReason);

        Assert.Equal(ActorStatus.Suspended, suspended.Status);
    }

    [Fact]
    public void SuspendNotFound()
    {
        var store = new InMemoryActorStore();

        Assert.Throws<ActorNotFoundException>(() =>
            store.Suspend(new ActorId("actor_nonexistent"), TestReason));
    }

    [Fact]
    public void SuspendAndReactivate()
    {
        var store = new InMemoryActorStore();
        var pk = TestPublicKey(1);
        var actor = store.Register(pk, "Alice", ActorType.Human);

        store.Suspend(actor.Id, TestReason);

        var got = store.Get(actor.Id);
        Assert.Equal(ActorStatus.Suspended, got.Status);

        var reactivated = store.Reactivate(actor.Id, TestReason);
        Assert.Equal(ActorStatus.Active, reactivated.Status);
    }

    // ── Memorial ────────────────────────────────────────────────────────

    [Fact]
    public void Memorial()
    {
        var store = new InMemoryActorStore();
        var pk = TestPublicKey(1);
        var actor = store.Register(pk, "Alice", ActorType.Human);

        var memorial = store.Memorial(actor.Id, TestReason);

        Assert.Equal(ActorStatus.Memorial, memorial.Status);
    }

    [Fact]
    public void MemorialIsTerminal()
    {
        var store = new InMemoryActorStore();
        var pk = TestPublicKey(1);
        var actor = store.Register(pk, "Alice", ActorType.Human);

        store.Memorial(actor.Id, TestReason);

        Assert.Throws<InvalidTransitionException>(() =>
            store.Suspend(actor.Id, TestReason));
    }

    [Fact]
    public void MemorialReactivateIsError()
    {
        var store = new InMemoryActorStore();
        var pk = TestPublicKey(1);
        var actor = store.Register(pk, "Alice", ActorType.Human);

        store.Memorial(actor.Id, TestReason);

        Assert.Throws<InvalidTransitionException>(() =>
            store.Reactivate(actor.Id, TestReason));
    }

    // ── Reactivate ──────────────────────────────────────────────────────

    [Fact]
    public void ReactivateFromActiveIsError()
    {
        var store = new InMemoryActorStore();
        var pk = TestPublicKey(1);
        var actor = store.Register(pk, "Alice", ActorType.Human);

        Assert.Throws<InvalidTransitionException>(() =>
            store.Reactivate(actor.Id, TestReason));
    }

    [Fact]
    public void ReactivateNotFound()
    {
        var store = new InMemoryActorStore();

        Assert.Throws<ActorNotFoundException>(() =>
            store.Reactivate(new ActorId("actor_nonexistent"), TestReason));
    }

    // ── List ────────────────────────────────────────────────────────────

    [Fact]
    public void List()
    {
        var store = new InMemoryActorStore();
        for (byte i = 1; i <= 5; i++)
            store.Register(TestPublicKey(i), "Actor", ActorType.Human);

        var page = store.List(new ActorFilter { Limit = 10 });

        Assert.Equal(5, page.Items.Count);
    }

    [Fact]
    public void ListWithStatusFilter()
    {
        var store = new InMemoryActorStore();
        store.Register(TestPublicKey(1), "Active", ActorType.Human);
        store.Register(TestPublicKey(2), "Active2", ActorType.Human);
        store.Register(TestPublicKey(3), "Suspended", ActorType.Human);

        // Suspend the third
        var all = store.List(new ActorFilter { Limit = 10 });
        if (all.Items.Count >= 3)
            store.Suspend(all.Items[2].Id, TestReason);

        var activePage = store.List(new ActorFilter
        {
            Status = ActorStatus.Active,
            Limit = 10
        });

        Assert.Equal(2, activePage.Items.Count);
    }

    [Fact]
    public void ListWithTypeFilter()
    {
        var store = new InMemoryActorStore();
        store.Register(TestPublicKey(1), "Human", ActorType.Human);
        store.Register(TestPublicKey(2), "AI", ActorType.AI);
        store.Register(TestPublicKey(3), "System", ActorType.System);

        var page = store.List(new ActorFilter
        {
            Type = ActorType.AI,
            Limit = 10
        });

        Assert.Single(page.Items);
    }

    [Fact]
    public void ListPagination()
    {
        var store = new InMemoryActorStore();
        for (byte i = 1; i <= 5; i++)
            store.Register(TestPublicKey(i), "Actor", ActorType.Human);

        // Page 1
        var page1 = store.List(new ActorFilter { Limit = 2 });
        Assert.Equal(2, page1.Items.Count);
        Assert.True(page1.HasMore);

        // Page 2
        var page2 = store.List(new ActorFilter { Limit = 2, After = page1.Cursor });
        Assert.Equal(2, page2.Items.Count);
    }

    // ── Actor getters ───────────────────────────────────────────────────

    [Fact]
    public void ActorGetters()
    {
        var store = new InMemoryActorStore();
        var pk = TestPublicKey(42);
        var actor = store.Register(pk, "Alice", ActorType.AI);

        Assert.False(string.IsNullOrEmpty(actor.Id.Value));
        _ = actor.PublicKey;
        _ = actor.CreatedAtNanos;
        _ = actor.Metadata;
    }
}
