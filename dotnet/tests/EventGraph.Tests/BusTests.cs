namespace EventGraph.Tests;

public class EventBusTests
{
    private static Event Bootstrap() => EventFactory.CreateBootstrap(new ActorId("alice"), new NoopSigner());

    [Fact]
    public void SubscribeAndPublish()
    {
        using var bus = new EventBus(new InMemoryStore());
        var received = new List<Event>();
        var latch = new ManualResetEventSlim();

        bus.Subscribe(new SubscriptionPattern("*"), ev => { received.Add(ev); latch.Set(); });
        var boot = Bootstrap();
        bus.Publish(boot);
        latch.Wait(TimeSpan.FromSeconds(2));

        Assert.Single(received);
        Assert.Equal(boot.Id, received[0].Id);
    }

    [Fact]
    public void PatternFiltering()
    {
        var store = new InMemoryStore();
        using var bus = new EventBus(store);
        var trustEvents = new List<Event>();
        var allEvents = new List<Event>();
        var latch = new CountdownEvent(2);

        bus.Subscribe(new SubscriptionPattern("trust.*"), ev => trustEvents.Add(ev));
        bus.Subscribe(new SubscriptionPattern("*"), ev => { allEvents.Add(ev); if (allEvents.Count >= 2) latch.Signal(latch.CurrentCount); });

        var boot = Bootstrap();
        store.Append(boot);
        var e1 = EventFactory.CreateEvent(new EventType("trust.updated"), new ActorId("alice"), new(), new() { boot.Id }, new ConversationId("c"), boot.Hash, new NoopSigner());

        bus.Publish(boot);
        bus.Publish(e1);
        latch.Wait(TimeSpan.FromSeconds(2));
        Thread.Sleep(100);

        Assert.Equal(2, allEvents.Count);
        Assert.Single(trustEvents);
    }

    [Fact]
    public void ClosePreventsSubscribe()
    {
        var bus = new EventBus(new InMemoryStore());
        bus.Dispose();
        Assert.Equal(-1, bus.Subscribe(new SubscriptionPattern("*"), _ => { }));
    }

    [Fact]
    public void StoreProperty()
    {
        var store = new InMemoryStore();
        using var bus = new EventBus(store);
        Assert.Same(store, bus.Store);
    }
}
