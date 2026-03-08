namespace EventGraph;

public sealed record TickConfig(int MaxWavesPerTick = 10);

public sealed record TickResult(int Tick, int Waves, int Mutations, bool Quiesced, double DurationMs, List<string> Errors);

public sealed class TickEngine
{
    private readonly Lock _lock = new();
    private readonly PrimitiveRegistry _registry;
    private readonly InMemoryStore _store;
    private readonly TickConfig _config;
    private readonly Action<Event>? _publisher;
    private readonly ISigner _signer = new NoopSigner();
    private int _currentTick;

    public TickEngine(PrimitiveRegistry registry, InMemoryStore store, TickConfig? config = null, Action<Event>? publisher = null)
    {
        _registry = registry;
        _store = store;
        _config = config ?? new TickConfig();
        _publisher = publisher;
    }

    public TickResult Tick(List<Event>? pendingEvents = null)
    {
        lock (_lock)
        {
            var start = DateTime.UtcNow;
            _currentTick++;
            var tickNum = _currentTick;

            var waveEvents = new List<Event>(pendingEvents ?? []);
            var totalMutations = 0;
            var errors = new List<string>();
            var quiesced = false;
            var invokedThisTick = new HashSet<string>();
            var wavesRun = 0;

            for (var wave = 0; wave < _config.MaxWavesPerTick; wave++)
            {
                if (waveEvents.Count == 0 && wave > 0)
                {
                    quiesced = true;
                    break;
                }
                wavesRun = wave + 1;

                var snapshot = new Snapshot(
                    tickNum,
                    _registry.AllStates(),
                    waveEvents,
                    _store.Recent(50));

                var newEvents = new List<Event>();
                var allMutations = new List<Mutation>();

                foreach (var prim in _registry.All())
                {
                    var pid = prim.Id;
                    var lifecycle = _registry.GetLifecycle(pid);
                    if (lifecycle != Lifecycle.Active) continue;

                    // Cadence check — only between ticks
                    if (!invokedThisTick.Contains(pid.Value))
                    {
                        var last = _registry.GetLastTick(pid);
                        if (tickNum - last < prim.Cadence.Value) continue;
                    }

                    var matched = FilterEvents(waveEvents, prim.Subscriptions);
                    if (matched.Count == 0 && wave > 0) continue;

                    try { _registry.SetLifecycle(pid, Lifecycle.Processing); }
                    catch { continue; }

                    try
                    {
                        var mutations = prim.Process(tickNum, matched, snapshot);
                        allMutations.AddRange(mutations);
                    }
                    catch (Exception ex) { errors.Add($"{pid.Value}: {ex.Message}"); }

                    try
                    {
                        _registry.SetLifecycle(pid, Lifecycle.Emitting);
                        _registry.SetLifecycle(pid, Lifecycle.Active);
                    }
                    catch (Exception ex) { errors.Add($"{pid.Value} lifecycle: {ex.Message}"); }

                    invokedThisTick.Add(pid.Value);
                }

                foreach (var m in allMutations)
                {
                    totalMutations++;
                    try
                    {
                        var ev = ApplyMutation(m);
                        if (ev != null) newEvents.Add(ev);
                    }
                    catch (Exception ex) { errors.Add($"mutation: {ex.Message}"); }
                }

                waveEvents = newEvents;
            }

            foreach (var pidVal in invokedThisTick)
                _registry.SetLastTick(new PrimitiveId(pidVal), tickNum);

            if (!quiesced && wavesRun >= _config.MaxWavesPerTick)
                quiesced = false;

            var elapsed = (DateTime.UtcNow - start).TotalMilliseconds;
            return new TickResult(tickNum, wavesRun, totalMutations, quiesced, elapsed, errors);
        }
    }

    private Event? ApplyMutation(Mutation m)
    {
        switch (m)
        {
            case AddEventMutation ae:
                var head = _store.Head();
                var prevHash = head.IsSome ? head.Unwrap().Hash : Hash.Zero();
                var ev = EventFactory.CreateEvent(ae.Type, ae.Source, ae.Content, ae.Causes, ae.ConversationId, prevHash, _signer);
                _store.Append(ev);
                _publisher?.Invoke(ev);
                return ev;
            case UpdateStateMutation us:
                _registry.UpdateState(us.PrimitiveId, us.Key, us.Value);
                return null;
            case UpdateActivationMutation ua:
                _registry.SetActivation(ua.PrimitiveId, ua.Level);
                return null;
            case UpdateLifecycleMutation ul:
                _registry.SetLifecycle(ul.PrimitiveId, ul.State);
                return null;
            default:
                return null;
        }
    }

    private static List<Event> FilterEvents(List<Event> events, List<SubscriptionPattern> patterns)
    {
        var result = new List<Event>();
        foreach (var ev in events)
        {
            foreach (var pat in patterns)
            {
                if (pat.Matches(ev.Type))
                {
                    result.Add(ev);
                    break;
                }
            }
        }
        return result;
    }
}
