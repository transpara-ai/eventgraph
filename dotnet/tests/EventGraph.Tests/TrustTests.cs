namespace EventGraph.Tests;

public class TrustTests
{
    private static readonly EventId CauseId = new("019462a0-0000-7000-8000-000000000001");

    private static PublicKey TestPublicKey(byte b)
    {
        var key = new byte[32];
        key[0] = b;
        return new PublicKey(key);
    }

    private static Actor TestActor(string name, byte b)
    {
        var store = new InMemoryActorStore();
        return store.Register(TestPublicKey(b), name, ActorType.Human);
    }

    /// <summary>Create a trust evidence event with "current" and "domain" in content.</summary>
    private static Event MakeTrustEvent(ActorId source, double current, string domain = "general")
    {
        var content = new Dictionary<string, object?>
        {
            ["current"] = current,
            ["domain"] = domain,
        };
        return EventFactory.CreateEvent(
            new EventType("trust.updated"),
            source,
            content,
            new List<EventId> { CauseId },
            new ConversationId("conv_test"),
            Hash.Zero(),
            new NoopSigner());
    }

    /// <summary>Create a non-trust evidence event (no "current" key).</summary>
    private static Event MakeGenericEvent(ActorId source)
    {
        var content = new Dictionary<string, object?> { ["action"] = "observed" };
        return EventFactory.CreateEvent(
            new EventType("system.observed"),
            source,
            content,
            new List<EventId> { CauseId },
            new ConversationId("conv_test"),
            Hash.Zero(),
            new NoopSigner());
    }

    // ── InitialScore ────────────────────────────────────────────────────

    [Fact]
    public void InitialScore()
    {
        var model = new DefaultTrustModel();
        var actor = TestActor("Alice", 1);

        var metrics = model.Score(actor);

        Assert.Equal(0.0, metrics.Overall.Value);
        Assert.Equal(0.0, metrics.Confidence.Value);
    }

    // ── UpdateIncreasesTrust ────────────────────────────────────────────

    [Fact]
    public void UpdateIncreasesTrust()
    {
        var model = new DefaultTrustModel();
        var actor = TestActor("Alice", 1);

        var ev = MakeTrustEvent(actor.Id, 0.05);
        var metrics = model.Update(actor, ev);

        Assert.True(metrics.Overall.Value > 0.0, $"Expected trust > 0, got {metrics.Overall.Value}");
    }

    // ── UpdateDecreasesTrust ────────────────────────────────────────────

    [Fact]
    public void UpdateDecreasesTrust()
    {
        var model = new DefaultTrustModel();
        var actor = TestActor("Alice", 1);

        // First increase
        var ev1 = MakeTrustEvent(actor.Id, 0.08);
        model.Update(actor, ev1);

        // Then decrease
        var ev2 = MakeTrustEvent(actor.Id, 0.0);
        var metrics = model.Update(actor, ev2);

        Assert.True(metrics.Overall.Value < 0.08, $"Expected trust < 0.08, got {metrics.Overall.Value}");
    }

    // ── UpdateClampedToMaxAdjustment ────────────────────────────────────

    [Fact]
    public void UpdateClampedToMaxAdjustment()
    {
        var model = new DefaultTrustModel();
        var actor = TestActor("Alice", 1);

        // Try to increase by 0.5 — should be clamped to MaxAdjustment (0.1)
        var ev = MakeTrustEvent(actor.Id, 0.5);
        var metrics = model.Update(actor, ev);

        Assert.True(metrics.Overall.Value <= 0.1, $"Expected trust <= 0.1, got {metrics.Overall.Value}");
    }

    // ── UpdateDeduplication ─────────────────────────────────────────────

    [Fact]
    public void UpdateDeduplication()
    {
        var model = new DefaultTrustModel();
        var actor = TestActor("Alice", 1);

        var ev = MakeTrustEvent(actor.Id, 0.05);
        var m1 = model.Update(actor, ev);
        var m2 = model.Update(actor, ev); // same event again

        Assert.Equal(m1.Overall.Value, m2.Overall.Value);
        Assert.Equal(1, m2.Evidence.Count); // only one evidence entry
    }

    // ── UpdateTrendPositive ─────────────────────────────────────────────

    [Fact]
    public void UpdateTrendPositive()
    {
        var model = new DefaultTrustModel();
        var actor = TestActor("Alice", 1);

        var ev = MakeTrustEvent(actor.Id, 0.05);
        var metrics = model.Update(actor, ev);

        Assert.True(metrics.Trend.Value > 0, $"Expected positive trend, got {metrics.Trend.Value}");
    }

    // ── UpdateTrendNegative ─────────────────────────────────────────────

    [Fact]
    public void UpdateTrendNegative()
    {
        var model = new DefaultTrustModel();
        var actor = TestActor("Alice", 1);

        // First increase
        var ev1 = MakeTrustEvent(actor.Id, 0.05);
        model.Update(actor, ev1);

        // Then decrease — the model's current score is ~0.05, target is 0.0
        var ev2 = MakeTrustEvent(actor.Id, 0.0);
        var metrics = model.Update(actor, ev2);

        // Trend should not still be at +0.1 — it should have moved toward negative
        Assert.True(metrics.Trend.Value < 0.1, $"Expected trend < 0.1 after negative update, got {metrics.Trend.Value}");
    }

    // ── ScoreInDomain ───────────────────────────────────────────────────

    [Fact]
    public void ScoreInDomain()
    {
        var model = new DefaultTrustModel();
        var actor = TestActor("Alice", 1);

        // Update with domain "code_review"
        var ev = MakeTrustEvent(actor.Id, 0.05, "code_review");
        model.Update(actor, ev);

        var metrics = model.ScoreInDomain(actor, new DomainScope("code_review"));

        // Domain-specific score should exist and reflect the update
        Assert.True(metrics.Overall.Value > 0.0, $"Expected domain score > 0, got {metrics.Overall.Value}");
    }

    // ── ScoreInDomainFallback ───────────────────────────────────────────

    [Fact]
    public void ScoreInDomainFallback()
    {
        var model = new DefaultTrustModel();
        var actor = TestActor("Alice", 1);

        // Update in "general" domain but query a different domain
        var ev = MakeTrustEvent(actor.Id, 0.05);
        model.Update(actor, ev);

        var metrics = model.ScoreInDomain(actor, new DomainScope("unknown_domain"));

        // Falls back to global score with halved confidence
        var globalMetrics = model.Score(actor);
        Assert.Equal(globalMetrics.Overall.Value, metrics.Overall.Value);
        Assert.True(metrics.Confidence.Value <= globalMetrics.Confidence.Value * 0.5 + 0.001,
            $"Fallback confidence ({metrics.Confidence.Value}) should be <= half of global ({globalMetrics.Confidence.Value})");
    }

    // ── Decay ───────────────────────────────────────────────────────────

    [Fact]
    public void Decay()
    {
        var model = new DefaultTrustModel();
        var actor = TestActor("Alice", 1);

        // Build up trust
        for (var i = 0; i < 5; i++)
        {
            var ev = MakeTrustEvent(actor.Id, 0.1);
            model.Update(actor, ev);
        }

        var before = model.Score(actor);

        // Decay 10 days
        model.Decay(actor, TimeSpan.FromDays(10));

        var after = model.Score(actor);
        Assert.True(after.Overall.Value < before.Overall.Value,
            $"Trust should decrease after decay: before={before.Overall.Value}, after={after.Overall.Value}");
    }

    // ── DecayNegativeDuration ───────────────────────────────────────────

    [Fact]
    public void DecayNegativeDuration()
    {
        var model = new DefaultTrustModel();
        var actor = TestActor("Alice", 1);

        // Build some trust
        var ev = MakeTrustEvent(actor.Id, 0.05);
        model.Update(actor, ev);

        var before = model.Score(actor);

        // Negative duration should be a no-op
        var after = model.Decay(actor, TimeSpan.FromDays(-5));
        Assert.Equal(before.Overall.Value, after.Overall.Value);
    }

    // ── UpdateBetween ───────────────────────────────────────────────────

    [Fact]
    public void UpdateBetween()
    {
        var model = new DefaultTrustModel();
        var alice = TestActor("Alice", 1);
        var bob = TestActor("Bob", 2);

        var ev = MakeTrustEvent(alice.Id, 0.05);
        var metrics = model.UpdateBetween(alice, bob, ev);

        Assert.True(metrics.Overall.Value > 0.0,
            $"Directed trust should increase, got {metrics.Overall.Value}");
        Assert.Equal(bob.Id, metrics.Actor);
    }

    // ── BetweenNoRelationship ───────────────────────────────────────────

    [Fact]
    public void BetweenNoRelationship()
    {
        var model = new DefaultTrustModel();
        var alice = TestActor("Alice", 1);
        var bob = TestActor("Bob", 2);

        var metrics = model.Between(alice, bob);

        Assert.Equal(0.0, metrics.Overall.Value);
        Assert.Equal(0.0, metrics.Confidence.Value);
    }

    // ── BetweenAsymmetric ───────────────────────────────────────────────

    [Fact]
    public void BetweenAsymmetric()
    {
        var model = new DefaultTrustModel();
        var alice = TestActor("Alice", 1);
        var bob = TestActor("Bob", 2);

        // Alice trusts Bob
        var ev = MakeTrustEvent(alice.Id, 0.05);
        model.UpdateBetween(alice, bob, ev);

        var aliceToBob = model.Between(alice, bob);
        var bobToAlice = model.Between(bob, alice);

        // Trust is asymmetric — Alice→Bob has evidence, Bob→Alice does not
        Assert.True(aliceToBob.Overall.Value > 0.0, "Alice→Bob should have trust");
        Assert.Equal(0.0, bobToAlice.Overall.Value);
    }

    // ── EvidenceCappedAt100 ─────────────────────────────────────────────

    [Fact]
    public void EvidenceCappedAt100()
    {
        var model = new DefaultTrustModel();
        var actor = TestActor("Alice", 1);

        // Add 150 updates (each with unique EventId)
        for (var i = 0; i < 150; i++)
        {
            var ev = MakeGenericEvent(actor.Id);
            model.Update(actor, ev);
        }

        var metrics = model.Score(actor);
        Assert.True(metrics.Evidence.Count <= 100,
            $"Evidence should be capped at 100, got {metrics.Evidence.Count}");
    }

    // ── DecayDirectedTrust ──────────────────────────────────────────────

    [Fact]
    public void DecayDirectedTrust()
    {
        var model = new DefaultTrustModel();
        var alice = TestActor("Alice", 1);
        var bob = TestActor("Bob", 2);

        // Build directed trust Alice→Bob
        var ev = MakeTrustEvent(alice.Id, 0.1);
        model.UpdateBetween(alice, bob, ev);

        var before = model.Between(alice, bob);

        // Decay Alice — should also decay directed trust where Alice is "from"
        model.Decay(alice, TimeSpan.FromDays(10));

        var after = model.Between(alice, bob);
        Assert.True(after.Overall.Value < before.Overall.Value,
            $"Directed trust should decay: before={before.Overall.Value}, after={after.Overall.Value}");
    }
}
