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
    List<Event> ByType(EventType type, int limit);
    List<Event> BySource(ActorId source, int limit);
    List<Event> ByConversation(ConversationId id, int limit);
    List<Event> Ancestors(EventId id, int maxDepth);
    List<Event> Descendants(EventId id, int maxDepth);
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

    public List<Event> ByType(EventType type, int limit)
    {
        lock (_lock)
        {
            var result = new List<Event>();
            for (int i = _events.Count - 1; i >= 0 && result.Count < limit; i--)
            {
                if (_events[i].Type == type)
                    result.Add(_events[i]);
            }
            return result;
        }
    }

    public List<Event> BySource(ActorId source, int limit)
    {
        lock (_lock)
        {
            var result = new List<Event>();
            for (int i = _events.Count - 1; i >= 0 && result.Count < limit; i--)
            {
                if (_events[i].Source == source)
                    result.Add(_events[i]);
            }
            return result;
        }
    }

    public List<Event> ByConversation(ConversationId id, int limit)
    {
        lock (_lock)
        {
            var result = new List<Event>();
            for (int i = _events.Count - 1; i >= 0 && result.Count < limit; i--)
            {
                if (_events[i].ConversationId == id)
                    result.Add(_events[i]);
            }
            return result;
        }
    }

    public List<Event> Ancestors(EventId id, int maxDepth)
    {
        lock (_lock)
        {
            var result = new List<Event>();
            var visited = new HashSet<string>();
            var queue = new Queue<(EventId Id, int Depth)>();
            queue.Enqueue((id, 0));
            visited.Add(id.Value);

            while (queue.Count > 0)
            {
                var (currentId, depth) = queue.Dequeue();
                if (depth > maxDepth) continue;

                if (!_index.TryGetValue(currentId.Value, out var pos))
                    continue;

                var ev = _events[pos];
                if (depth > 0) // don't include the starting event
                    result.Add(ev);

                if (depth < maxDepth)
                {
                    foreach (var causeId in ev.Causes)
                    {
                        if (visited.Add(causeId.Value))
                            queue.Enqueue((causeId, depth + 1));
                    }
                }
            }
            return result;
        }
    }

    public List<Event> Descendants(EventId id, int maxDepth)
    {
        lock (_lock)
        {
            // Build reverse index lazily
            var childIndex = new Dictionary<string, List<int>>();
            for (int i = 0; i < _events.Count; i++)
            {
                foreach (var causeId in _events[i].Causes)
                {
                    if (!childIndex.TryGetValue(causeId.Value, out var list))
                    {
                        list = new List<int>();
                        childIndex[causeId.Value] = list;
                    }
                    list.Add(i);
                }
            }

            var result = new List<Event>();
            var visited = new HashSet<string>();
            var queue = new Queue<(string Id, int Depth)>();
            queue.Enqueue((id.Value, 0));
            visited.Add(id.Value);

            while (queue.Count > 0)
            {
                var (currentId, depth) = queue.Dequeue();
                if (depth > maxDepth) continue;

                if (depth > 0 && _index.TryGetValue(currentId, out var pos))
                    result.Add(_events[pos]);

                if (depth < maxDepth && childIndex.TryGetValue(currentId, out var children))
                {
                    foreach (var childPos in children)
                    {
                        var childId = _events[childPos].Id.Value;
                        if (visited.Add(childId))
                            queue.Enqueue((childId, depth + 1));
                    }
                }
            }
            return result;
        }
    }

    public void Close() { }
}
