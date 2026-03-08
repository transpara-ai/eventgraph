using System.Collections.Concurrent;

namespace EventGraph;

/// <summary>Thread-safe event bus with per-subscriber buffering and non-blocking delivery.</summary>
public sealed class EventBus : IDisposable
{
    private readonly InMemoryStore _store;
    private readonly int _bufferSize;
    private readonly Lock _lock = new();
    private readonly Dictionary<int, Subscription> _subs = new();
    private int _nextId;
    private bool _closed;

    public EventBus(InMemoryStore store, int bufferSize = 256)
    {
        _store = store;
        _bufferSize = Math.Max(bufferSize, 1);
    }

    public InMemoryStore Store => _store;

    public int Subscribe(SubscriptionPattern pattern, Action<Event> handler)
    {
        lock (_lock)
        {
            if (_closed) return -1;
            var id = ++_nextId;
            var sub = new Subscription(id, pattern, handler, _bufferSize);
            _subs[id] = sub;
            sub.Start();
            return id;
        }
    }

    public void Unsubscribe(int subId)
    {
        Subscription? sub;
        lock (_lock) { _subs.Remove(subId, out sub); }
        sub?.Stop();
    }

    public void Publish(Event ev)
    {
        List<Subscription> snapshot;
        lock (_lock)
        {
            if (_closed) return;
            snapshot = new(_subs.Values);
        }
        foreach (var sub in snapshot)
            sub.Deliver(ev);
    }

    public void Dispose()
    {
        List<Subscription> subs;
        lock (_lock)
        {
            if (_closed) return;
            _closed = true;
            subs = new(_subs.Values);
            _subs.Clear();
        }
        foreach (var sub in subs)
            sub.Stop();
        foreach (var sub in subs)
            sub.Join(TimeSpan.FromSeconds(30));
    }

    private sealed class Subscription
    {
        private readonly int _id;
        private readonly SubscriptionPattern _pattern;
        private readonly Action<Event> _handler;
        private readonly BlockingCollection<Event> _buffer;
        private Thread? _thread;
        public Exception? LastError;

        public Subscription(int id, SubscriptionPattern pattern, Action<Event> handler, int bufferSize)
        {
            _id = id;
            _pattern = pattern;
            _handler = handler;
            _buffer = new BlockingCollection<Event>(bufferSize);
        }

        public void Start()
        {
            _thread = new Thread(Run) { IsBackground = true };
            _thread.Start();
        }

        public void Deliver(Event ev)
        {
            if (!_pattern.Matches(ev.Type)) return;
            _buffer.TryAdd(ev); // drops if full
        }

        public void Stop() => _buffer.CompleteAdding();

        public void Join(TimeSpan timeout) => _thread?.Join(timeout);

        private void Run()
        {
            try
            {
                foreach (var ev in _buffer.GetConsumingEnumerable())
                {
                    try { _handler(ev); }
                    catch (Exception ex) { LastError = ex; }
                }
            }
            catch (OperationCanceledException) { }
        }
    }
}
