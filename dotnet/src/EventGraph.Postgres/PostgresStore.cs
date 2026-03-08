using System.Text.Json;
using Npgsql;

using EventGraph;

namespace EventGraph.Postgres;

/// <summary>
/// PostgreSQL-backed event store. Satisfies IStore. Data persists across process restarts.
/// Requires Npgsql NuGet package: dotnet add package Npgsql --version 9.0.0
/// </summary>
public sealed class PostgresStore : IStore, IDisposable
{
    private readonly NpgsqlConnection _conn;
    private readonly Lock _lock = new();

    public PostgresStore(string connectionString)
    {
        _conn = new NpgsqlConnection(connectionString);
        _conn.Open();
        InitSchema();
    }

    private void InitSchema()
    {
        using var cmd = new NpgsqlCommand("""
            CREATE TABLE IF NOT EXISTS events (
                position        BIGSERIAL PRIMARY KEY,
                event_id        TEXT NOT NULL UNIQUE,
                event_type      TEXT NOT NULL,
                version         INTEGER NOT NULL,
                timestamp_nanos BIGINT NOT NULL,
                source          TEXT NOT NULL,
                content         TEXT NOT NULL,
                causes          TEXT NOT NULL,
                conversation_id TEXT NOT NULL,
                hash            TEXT NOT NULL,
                prev_hash       TEXT NOT NULL,
                signature       BYTEA NOT NULL
            );
            CREATE INDEX IF NOT EXISTS idx_events_type ON events(event_type);
            CREATE INDEX IF NOT EXISTS idx_events_source ON events(source);
            CREATE INDEX IF NOT EXISTS idx_events_conversation ON events(conversation_id);
            """, _conn);
        cmd.ExecuteNonQuery();
    }

    // Column mapping: 0=position, 1=event_id, 2=event_type, 3=version,
    // 4=timestamp_nanos, 5=source, 6=content, 7=causes, 8=conversation_id,
    // 9=hash, 10=prev_hash, 11=signature

    private static Event ReadEvent(NpgsqlDataReader r)
    {
        var causeIds = JsonSerializer.Deserialize<List<string>>(r.GetString(7))!
            .Select(c => new EventId(c)).ToList();
        var contentJson = r.GetString(6);
        var content = JsonSerializer.Deserialize<Dictionary<string, object?>>(contentJson)
            ?? new Dictionary<string, object?>();
        var sigBytes = (byte[])r["signature"];

        return new Event(
            r.GetInt32(3),
            new EventId(r.GetString(1)),
            new EventType(r.GetString(2)),
            r.GetInt64(4),
            new ActorId(r.GetString(5)),
            content,
            NonEmpty<EventId>.Of(causeIds),
            new ConversationId(r.GetString(8)),
            new Hash(r.GetString(9)),
            new Hash(r.GetString(10)),
            new Signature(sigBytes)
        );
    }

    private List<Event> QueryEvents(string sql, params NpgsqlParameter[] parameters)
    {
        using var cmd = new NpgsqlCommand(sql, _conn);
        foreach (var p in parameters) cmd.Parameters.Add(p);
        using var reader = cmd.ExecuteReader();
        var result = new List<Event>();
        while (reader.Read()) result.Add(ReadEvent(reader));
        return result;
    }

    private Event? QuerySingleEvent(string sql, params NpgsqlParameter[] parameters)
    {
        using var cmd = new NpgsqlCommand(sql, _conn);
        foreach (var p in parameters) cmd.Parameters.Add(p);
        using var reader = cmd.ExecuteReader();
        return reader.Read() ? ReadEvent(reader) : null;
    }

    public Event Append(Event ev)
    {
        lock (_lock)
        {
            // Idempotency
            var existing = QuerySingleEvent(
                "SELECT * FROM events WHERE event_id = @id",
                new NpgsqlParameter("@id", ev.Id.Value));
            if (existing is not null)
            {
                if (existing.Hash != ev.Hash)
                    throw new ChainIntegrityException(0,
                        $"hash mismatch for existing event {ev.Id.Value}");
                return existing;
            }

            // Chain continuity
            var head = QuerySingleEvent(
                "SELECT * FROM events ORDER BY position DESC LIMIT 1");
            if (head is not null)
            {
                if (ev.PrevHash != head.Hash)
                    throw new ChainIntegrityException(0,
                        $"prev_hash {ev.PrevHash.Value} != head hash {head.Hash.Value}");
            }

            var causesJson = JsonSerializer.Serialize(
                ev.Causes.Select(c => c.Value).ToList());
            var contentJson = JsonSerializer.Serialize(
                ev.Content, new JsonSerializerOptions { WriteIndented = false });

            using var cmd = new NpgsqlCommand("""
                INSERT INTO events
                (event_id, event_type, version, timestamp_nanos, source,
                 content, causes, conversation_id, hash, prev_hash, signature)
                VALUES (@eid, @et, @v, @ts, @src, @cnt, @cau, @cid, @h, @ph, @sig)
                """, _conn);
            cmd.Parameters.AddWithValue("@eid", ev.Id.Value);
            cmd.Parameters.AddWithValue("@et", ev.Type.Value);
            cmd.Parameters.AddWithValue("@v", ev.Version);
            cmd.Parameters.AddWithValue("@ts", ev.TimestampNanos);
            cmd.Parameters.AddWithValue("@src", ev.Source.Value);
            cmd.Parameters.AddWithValue("@cnt", contentJson);
            cmd.Parameters.AddWithValue("@cau", causesJson);
            cmd.Parameters.AddWithValue("@cid", ev.ConversationId.Value);
            cmd.Parameters.AddWithValue("@h", ev.Hash.Value);
            cmd.Parameters.AddWithValue("@ph", ev.PrevHash.Value);
            cmd.Parameters.AddWithValue("@sig", ev.Signature.Bytes.ToArray());
            cmd.ExecuteNonQuery();

            return ev;
        }
    }

    public Event Get(EventId eventId)
    {
        lock (_lock)
        {
            return QuerySingleEvent(
                "SELECT * FROM events WHERE event_id = @id",
                new NpgsqlParameter("@id", eventId.Value))
                ?? throw new EventNotFoundException(eventId.Value);
        }
    }

    public Option<Event> Head()
    {
        lock (_lock)
        {
            var ev = QuerySingleEvent(
                "SELECT * FROM events ORDER BY position DESC LIMIT 1");
            return ev is null ? Option<Event>.None() : Option<Event>.Some(ev);
        }
    }

    public int Count()
    {
        lock (_lock)
        {
            using var cmd = new NpgsqlCommand("SELECT COUNT(*) FROM events", _conn);
            return Convert.ToInt32(cmd.ExecuteScalar());
        }
    }

    public ChainVerification VerifyChain()
    {
        lock (_lock)
        {
            var events = QueryEvents("SELECT * FROM events ORDER BY position ASC");
            for (int i = 0; i < events.Count; i++)
            {
                if (i > 0 && events[i].PrevHash != events[i - 1].Hash)
                    return new ChainVerification(false, i);
            }
            return new ChainVerification(true, events.Count);
        }
    }

    public List<Event> Recent(int limit)
    {
        lock (_lock)
        {
            return QueryEvents(
                "SELECT * FROM events ORDER BY position DESC LIMIT @lim",
                new NpgsqlParameter("@lim", limit));
        }
    }

    public List<Event> ByType(EventType type, int limit)
    {
        lock (_lock)
        {
            return QueryEvents(
                "SELECT * FROM events WHERE event_type = @t ORDER BY position DESC LIMIT @lim",
                new NpgsqlParameter("@t", type.Value),
                new NpgsqlParameter("@lim", limit));
        }
    }

    public List<Event> BySource(ActorId source, int limit)
    {
        lock (_lock)
        {
            return QueryEvents(
                "SELECT * FROM events WHERE source = @s ORDER BY position DESC LIMIT @lim",
                new NpgsqlParameter("@s", source.Value),
                new NpgsqlParameter("@lim", limit));
        }
    }

    public List<Event> ByConversation(ConversationId id, int limit)
    {
        lock (_lock)
        {
            return QueryEvents(
                "SELECT * FROM events WHERE conversation_id = @c ORDER BY position DESC LIMIT @lim",
                new NpgsqlParameter("@c", id.Value),
                new NpgsqlParameter("@lim", limit));
        }
    }

    public List<Event> Ancestors(EventId id, int maxDepth)
    {
        lock (_lock)
        {
            var start = QuerySingleEvent(
                "SELECT * FROM events WHERE event_id = @id",
                new NpgsqlParameter("@id", id.Value))
                ?? throw new EventNotFoundException(id.Value);

            var result = new List<Event>();
            var seen = new HashSet<string> { id.Value };
            var frontier = start.Causes
                .Where(c => c.Value != id.Value)
                .Select(c => c.Value).ToList();

            for (int d = 0; d < maxDepth && frontier.Count > 0; d++)
            {
                var next = new List<string>();
                foreach (var eid in frontier)
                {
                    if (!seen.Add(eid)) continue;
                    var ev = QuerySingleEvent(
                        "SELECT * FROM events WHERE event_id = @id",
                        new NpgsqlParameter("@id", eid));
                    if (ev is null) continue;
                    result.Add(ev);
                    foreach (var c in ev.Causes)
                        if (!seen.Contains(c.Value)) next.Add(c.Value);
                }
                frontier = next;
            }
            return result;
        }
    }

    public List<Event> Descendants(EventId id, int maxDepth)
    {
        lock (_lock)
        {
            var exists = QuerySingleEvent(
                "SELECT * FROM events WHERE event_id = @id",
                new NpgsqlParameter("@id", id.Value))
                ?? throw new EventNotFoundException(id.Value);

            // Build reverse index
            var allEvents = QueryEvents("SELECT * FROM events ORDER BY position ASC");
            var children = new Dictionary<string, List<string>>();
            foreach (var ev in allEvents)
            {
                foreach (var c in ev.Causes)
                {
                    if (c.Value == ev.Id.Value) continue;
                    if (!children.ContainsKey(c.Value))
                        children[c.Value] = new List<string>();
                    children[c.Value].Add(ev.Id.Value);
                }
            }

            var result = new List<Event>();
            var seen = new HashSet<string> { id.Value };
            var frontier = children.GetValueOrDefault(id.Value, new List<string>());

            for (int d = 0; d < maxDepth && frontier.Count > 0; d++)
            {
                var next = new List<string>();
                foreach (var eid in frontier)
                {
                    if (!seen.Add(eid)) continue;
                    var ev = QuerySingleEvent(
                        "SELECT * FROM events WHERE event_id = @id",
                        new NpgsqlParameter("@id", eid));
                    if (ev is null) continue;
                    result.Add(ev);
                    foreach (var child in children.GetValueOrDefault(eid, new List<string>()))
                        if (!seen.Contains(child)) next.Add(child);
                }
                frontier = next;
            }
            return result;
        }
    }

    public void Close() => Dispose();

    public void Dispose()
    {
        _conn.Close();
        _conn.Dispose();
    }
}
