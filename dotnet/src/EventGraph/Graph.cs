namespace EventGraph;

/// <summary>Configuration for the Graph facade.</summary>
public sealed record GraphConfig
{
    public int SubscriberBufferSize { get; init; } = 256;
    public bool FallbackToMechanical { get; init; } = true;
}

/// <summary>
/// Top-level facade — IGraph.
/// Provides Evaluate(), Record(), Query(), Start(), Close().
/// Thread-safe with lock; tracks started/closed state.
/// </summary>
public sealed class Graph : IDisposable
{
    private readonly object _lock = new();
    private readonly InMemoryStore _store;
    private readonly IActorStore _actorStore;
    private readonly EventBus _bus;
    private readonly ISigner _signer;
    private readonly ITrustModel _trustModel;
    private readonly IAuthorityChain _authorityChain;
    private readonly GraphConfig _config;
    private bool _started;
    private bool _closed;

    /// <summary>The underlying event store.</summary>
    public InMemoryStore Store => _store;

    /// <summary>The underlying actor store.</summary>
    public IActorStore ActorStore => _actorStore;

    /// <summary>The event bus.</summary>
    public EventBus Bus => _bus;

    public Graph(
        InMemoryStore store,
        IActorStore actorStore,
        ITrustModel? trustModel = null,
        IAuthorityChain? authorityChain = null,
        ISigner? signer = null,
        GraphConfig? config = null)
    {
        _store = store;
        _actorStore = actorStore;
        _config = config ?? new GraphConfig();
        _signer = signer ?? new NoopSigner();
        _trustModel = trustModel ?? new DefaultTrustModel();
        _authorityChain = authorityChain ?? new DefaultAuthorityChain(_trustModel);
        _bus = new EventBus(_store, _config.SubscriberBufferSize);
    }

    /// <summary>Initialize the graph. Must be called before Record/Evaluate/Query.</summary>
    public void Start()
    {
        lock (_lock)
        {
            if (_started) return;
            _started = true;
        }
    }

    /// <summary>Graceful shutdown. Disposes the bus and closes the store.</summary>
    public void Close()
    {
        lock (_lock)
        {
            if (_closed) return;
            _closed = true;
        }
        // Release lock before bus/store shutdown to prevent deadlock
        // if a subscriber calls back into the Graph.
        _bus.Dispose();
        _store.Close();
    }

    /// <summary>Implements IDisposable by delegating to Close.</summary>
    public void Dispose() => Close();

    /// <summary>
    /// Initialize the graph with a genesis event.
    /// Must be called on a started, empty graph.
    /// </summary>
    public Event Bootstrap(ActorId systemActor, ISigner? signer = null)
    {
        Event stored;
        lock (_lock)
        {
            EnsureRunning();

            if (_store.Count() > 0)
                throw new InvalidOperationException("Graph already bootstrapped");

            var s = signer ?? _signer;
            var ev = EventFactory.CreateBootstrap(systemActor, s);
            stored = _store.Append(ev);
        }
        // Publish outside lock to prevent deadlock
        _bus.Publish(stored);
        return stored;
    }

    /// <summary>
    /// Create and persist an event, then notify the bus.
    /// Uses exclusive lock to serialize event creation + append, ensuring hash chain integrity.
    /// </summary>
    public Event Record(
        EventType eventType,
        ActorId source,
        Dictionary<string, object?> content,
        List<EventId> causes,
        ConversationId conversationId,
        ISigner? signer = null)
    {
        Event stored;
        lock (_lock)
        {
            EnsureRunning();

            var s = signer ?? _signer;
            var prevHash = _store.Head().IsSome
                ? _store.Head().Unwrap().Hash
                : Hash.Zero();

            var ev = EventFactory.CreateEvent(
                eventType, source, content, causes,
                conversationId, prevHash, s);
            stored = _store.Append(ev);
        }
        // Publish outside lock to prevent deadlock
        _bus.Publish(stored);
        return stored;
    }

    /// <summary>Evaluate authority for an actor performing an action.</summary>
    public AuthorityResult Evaluate(Actor actor, string action)
    {
        lock (_lock)
        {
            EnsureRunning();
        }
        return _authorityChain.Evaluate(actor.Id, action);
    }

    /// <summary>Return a query builder for the graph.</summary>
    public GraphQuery Query()
    {
        lock (_lock)
        {
            EnsureRunning();
        }
        return new GraphQuery(_store, _actorStore, _trustModel);
    }

    /// <summary>Throw if the graph is not in a usable state.</summary>
    private void EnsureRunning()
    {
        if (_closed)
            throw new InvalidOperationException("Graph is closed");
        if (!_started)
            throw new InvalidOperationException("Graph is not started (call Start first)");
    }
}

/// <summary>Query builder for the graph — wraps store, actor store, and trust model.</summary>
public sealed class GraphQuery
{
    private readonly InMemoryStore _store;
    private readonly IActorStore _actorStore;
    private readonly ITrustModel _trustModel;

    internal GraphQuery(InMemoryStore store, IActorStore actorStore, ITrustModel trustModel)
    {
        _store = store;
        _actorStore = actorStore;
        _trustModel = trustModel;
    }

    /// <summary>Return the most recent events.</summary>
    public List<Event> Recent(int limit) => _store.Recent(limit);

    /// <summary>Return events of a given type.</summary>
    public List<Event> ByType(EventType type, int limit) => _store.ByType(type, limit);

    /// <summary>Return events from a given source.</summary>
    public List<Event> BySource(ActorId source, int limit) => _store.BySource(source, limit);

    /// <summary>Return events in a conversation.</summary>
    public List<Event> ByConversation(ConversationId id, int limit) => _store.ByConversation(id, limit);

    /// <summary>Return causal ancestors of an event.</summary>
    public List<Event> Ancestors(EventId id, int maxDepth) => _store.Ancestors(id, maxDepth);

    /// <summary>Return causal descendants of an event.</summary>
    public List<Event> Descendants(EventId id, int maxDepth) => _store.Descendants(id, maxDepth);

    /// <summary>Return trust metrics for an actor.</summary>
    public TrustMetrics TrustScore(Actor actor) => _trustModel.Score(actor);

    /// <summary>Return directional trust from one actor toward another.</summary>
    public TrustMetrics TrustBetween(Actor from, Actor to) => _trustModel.Between(from, to);

    /// <summary>Return an actor by ID.</summary>
    public Actor GetActor(ActorId id) => _actorStore.Get(id);

    /// <summary>Return the total number of events.</summary>
    public int EventCount() => _store.Count();
}
