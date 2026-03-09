using EventGraph.CodeGraph;

namespace EventGraph.Tests;

// ── Code Graph Event Types Tests ─────────────────────────────────────────

public class CodeGraphEventTypeTests
{
    [Fact]
    public void AllCodeGraphEventTypes_Returns35()
    {
        var all = CodeGraphEventTypes.AllCodeGraphEventTypes();
        Assert.Equal(35, all.Count);
    }

    [Fact]
    public void AllEventTypesStartWithCodeGraph()
    {
        foreach (var et in CodeGraphEventTypes.AllCodeGraphEventTypes())
            Assert.StartsWith("codegraph.", et.Value);
    }

    [Fact]
    public void AllEventTypesAreUnique()
    {
        var all = CodeGraphEventTypes.AllCodeGraphEventTypes();
        var unique = new HashSet<string>(all.Select(e => e.Value));
        Assert.Equal(all.Count, unique.Count);
    }
}

// ── Code Graph Primitives Tests ──────────────────────────────────────────

public class CodeGraphPrimitivesTests
{
    [Fact]
    public void AllPrimitives_Returns61()
    {
        var all = CodeGraphPrimitiveFactory.All();
        Assert.Equal(61, all.Count);
    }

    [Fact]
    public void NoDuplicateIds()
    {
        var all = CodeGraphPrimitiveFactory.All();
        var ids = new HashSet<string>(all.Select(p => p.Id.Value));
        Assert.Equal(61, ids.Count);
    }

    [Fact]
    public void AllLayer5()
    {
        foreach (var p in CodeGraphPrimitiveFactory.All())
            Assert.Equal(5, p.Layer.Value);
    }

    [Fact]
    public void AllHaveSubscriptions()
    {
        foreach (var p in CodeGraphPrimitiveFactory.All())
            Assert.NotEmpty(p.Subscriptions);
    }

    [Fact]
    public void AllHaveCadence1()
    {
        foreach (var p in CodeGraphPrimitiveFactory.All())
            Assert.Equal(1, p.Cadence.Value);
    }

    [Fact]
    public void RegisterAll()
    {
        var registry = new PrimitiveRegistry();
        CodeGraphPrimitiveFactory.RegisterAll(registry);
        Assert.Equal(61, registry.Count);
    }

    [Fact]
    public void RegisterAllActivatesPrimitives()
    {
        var registry = new PrimitiveRegistry();
        CodeGraphPrimitiveFactory.RegisterAll(registry);
        foreach (var p in CodeGraphPrimitiveFactory.All())
            Assert.Equal(Lifecycle.Active, registry.GetLifecycle(p.Id));
    }

    [Fact]
    public void ProcessReturnsMutations()
    {
        var emptyEvents = new List<Event>();
        var snap = AgentTestHelpers.EmptySnapshot();
        foreach (var p in CodeGraphPrimitiveFactory.All())
        {
            var mutations = p.Process(1, emptyEvents, snap);
            Assert.NotNull(mutations);
            Assert.Equal(2, mutations.Count);
            Assert.Contains(mutations, m => m is UpdateStateMutation u && u.Key == "eventsProcessed" && (int)u.Value! == 0);
            Assert.Contains(mutations, m => m is UpdateStateMutation u && u.Key == "lastTick" && (int)u.Value! == 1);
        }
    }

    [Fact]
    public void ProcessCountsEvents()
    {
        var events = new List<Event>
        {
            AgentTestHelpers.MakeEvent("codegraph.entity.created"),
            AgentTestHelpers.MakeEvent("codegraph.entity.updated"),
            AgentTestHelpers.MakeEvent("codegraph.entity.deleted"),
        };
        var snap = AgentTestHelpers.EmptySnapshot(7);
        var p = new EntityPrimitive();
        var mutations = p.Process(7, events, snap);
        Assert.Contains(mutations, m => m is UpdateStateMutation u && u.Key == "eventsProcessed" && (int)u.Value! == 3);
        Assert.Contains(mutations, m => m is UpdateStateMutation u && u.Key == "lastTick" && (int)u.Value! == 7);
    }

    [Fact]
    public void IsCodeGraphPrimitive()
    {
        foreach (var p in CodeGraphPrimitiveFactory.All())
            Assert.True(CodeGraphPrimitiveFactory.IsCodeGraph(p.Id));
    }

    [Fact]
    public void IsCodeGraphReturnsFalseForNonCodeGraph()
    {
        Assert.False(CodeGraphPrimitiveFactory.IsCodeGraph(new PrimitiveId("agent.Identity")));
        Assert.False(CodeGraphPrimitiveFactory.IsCodeGraph(new PrimitiveId("system.Event")));
        Assert.False(CodeGraphPrimitiveFactory.IsCodeGraph(new PrimitiveId("Hash")));
    }
}

// ── Code Graph Compositions Tests ────────────────────────────────────────

public class CodeGraphCompositionsTests
{
    [Fact]
    public void AllCompositions_Returns7()
    {
        var all = CodeGraphCompositions.All();
        Assert.Equal(7, all.Count);
    }

    [Fact]
    public void AllCompositionsHaveUniqueNames()
    {
        var all = CodeGraphCompositions.All();
        var names = new HashSet<string>(all.Select(c => c.Name));
        Assert.Equal(7, names.Count);
    }

    [Fact]
    public void AllCompositionsHaveNonEmptyPrimitives()
    {
        foreach (var c in CodeGraphCompositions.All())
        {
            Assert.NotEmpty(c.Primitives);
            Assert.NotEmpty(c.Events);
        }
    }

    [Fact]
    public void CompositionPrimitivesExist()
    {
        var allIds = new HashSet<string>(CodeGraphPrimitiveFactory.All().Select(p => p.Id.Value));
        foreach (var c in CodeGraphCompositions.All())
            foreach (var pId in c.Primitives)
                Assert.Contains(pId, allIds);
    }

    [Fact]
    public void AllCompositionPrimitivesStartWithCG()
    {
        foreach (var c in CodeGraphCompositions.All())
            foreach (var pId in c.Primitives)
                Assert.StartsWith("CG", pId);
    }

    [Fact]
    public void Board_Has10Primitives()
    {
        var board = CodeGraphCompositions.Board();
        Assert.Equal("Board", board.Name);
        Assert.Equal(10, board.Primitives.Count);
    }

    [Fact]
    public void Detail_Has9Primitives()
    {
        var detail = CodeGraphCompositions.Detail();
        Assert.Equal("Detail", detail.Name);
        Assert.Equal(9, detail.Primitives.Count);
    }

    [Fact]
    public void Feed_Has7Primitives()
    {
        var feed = CodeGraphCompositions.Feed();
        Assert.Equal("Feed", feed.Name);
        Assert.Equal(7, feed.Primitives.Count);
    }

    [Fact]
    public void Dashboard_Has5Primitives()
    {
        var dashboard = CodeGraphCompositions.Dashboard();
        Assert.Equal("Dashboard", dashboard.Name);
        Assert.Equal(5, dashboard.Primitives.Count);
    }

    [Fact]
    public void Inbox_Has8Primitives()
    {
        var inbox = CodeGraphCompositions.Inbox();
        Assert.Equal("Inbox", inbox.Name);
        Assert.Equal(8, inbox.Primitives.Count);
    }

    [Fact]
    public void Wizard_Has8Primitives()
    {
        var wizard = CodeGraphCompositions.Wizard();
        Assert.Equal("Wizard", wizard.Name);
        Assert.Equal(8, wizard.Primitives.Count);
    }

    [Fact]
    public void Skin_Has7Primitives()
    {
        var skin = CodeGraphCompositions.Skin();
        Assert.Equal("Skin", skin.Name);
        Assert.Equal(7, skin.Primitives.Count);
    }

    [Fact]
    public void CompositionNamesAreCorrect()
    {
        var expected = new[] { "Board", "Detail", "Feed", "Dashboard", "Inbox", "Wizard", "Skin" };
        var actual = CodeGraphCompositions.All().Select(c => c.Name).ToArray();
        Assert.Equal(expected, actual);
    }
}
