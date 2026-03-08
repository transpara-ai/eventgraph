namespace EventGraph;

/// <summary>Result of a chain integrity check.</summary>
public readonly record struct ChainVerification(bool Valid, int Length);

/// <summary>Event persistence protocol — append-only, hash-chained.</summary>
public interface IStore
{
    Event Append(Event ev);
    Event Get(EventId eventId);
    Option<Event> Head();
    int Count();
    ChainVerification VerifyChain();
    void Close();
}

/// <summary>Thread-safe in-memory event store.</summary>
public sealed class InMemoryStore : IStore
{
    private readonly Lock _lock = new();
    private readonly List<Event> _events = new();
    private readonly Dictionary<string, int> _index = new();

    public Event Append(Event ev)
    {
        lock (_lock)
        {
            if (_events.Count > 0)
            {
                var last = _events[^1];
                if (ev.PrevHash != last.Hash)
                    throw new ChainIntegrityException(
                        _events.Count,
                        $"prev_hash {ev.PrevHash.Value} != head hash {last.Hash.Value}");
            }
            _events.Add(ev);
            _index[ev.Id.Value] = _events.Count - 1;
            return ev;
        }
    }

    public Event Get(EventId eventId)
    {
        lock (_lock)
        {
            if (_index.TryGetValue(eventId.Value, out var pos))
                return _events[pos];
            throw new EventNotFoundException(eventId.Value);
        }
    }

    public Option<Event> Head()
    {
        lock (_lock)
        {
            return _events.Count == 0
                ? Option<Event>.None()
                : Option<Event>.Some(_events[^1]);
        }
    }

    public int Count()
    {
        lock (_lock) { return _events.Count; }
    }

    public ChainVerification VerifyChain()
    {
        lock (_lock)
        {
            for (int i = 1; i < _events.Count; i++)
            {
                if (_events[i - 1].Hash != _events[i].PrevHash)
                    return new ChainVerification(false, i);
            }
            return new ChainVerification(true, _events.Count);
        }
    }

    public List<Event> Recent(int limit)
    {
        lock (_lock)
        {
            var start = Math.Max(0, _events.Count - limit);
            var result = _events.GetRange(start, _events.Count - start);
            result.Reverse();
            return result;
        }
    }

    public void Close() { }
}
