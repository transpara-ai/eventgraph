// Requires Microsoft.Data.SqlClient NuGet package (the modern SQL Server client).
// Install via: dotnet add package Microsoft.Data.SqlClient

using System.Text.Json;
using Microsoft.Data.SqlClient;

using EventGraph;

namespace EventGraph.SqlServer;

/// <summary>
/// SQL Server-backed event store. Satisfies IStore. Data persists across process restarts.
/// Requires Microsoft.Data.SqlClient NuGet package.
/// </summary>
public sealed class SqlServerStore : IStore, IDisposable
{
    private readonly SqlConnection _conn;
    private readonly Lock _lock = new();

    public SqlServerStore(string connectionString)
    {
        _conn = new SqlConnection(connectionString);
        _conn.Open();
        InitSchema();
    }

    private void InitSchema()
    {
        using var cmd = _conn.CreateCommand();
        cmd.CommandText = """
            IF OBJECT_ID('events', 'U') IS NULL
            CREATE TABLE events (
                position        BIGINT IDENTITY(1,1) PRIMARY KEY,
                event_id        NVARCHAR(100) NOT NULL UNIQUE,
                event_type      NVARCHAR(200) NOT NULL,
                version         INT NOT NULL,
                timestamp_nanos BIGINT NOT NULL,
                source          NVARCHAR(200) NOT NULL,
                content         NVARCHAR(MAX) NOT NULL,
                causes          NVARCHAR(MAX) NOT NULL,
                conversation_id NVARCHAR(200) NOT NULL,
                hash            NVARCHAR(100) NOT NULL,
                prev_hash       NVARCHAR(100) NOT NULL,
                signature       VARBINARY(64) NOT NULL
            );

            IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'idx_events_type')
                CREATE INDEX idx_events_type ON events(event_type);

            IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'idx_events_source')
                CREATE INDEX idx_events_source ON events(source);

            IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'idx_events_conversation')
                CREATE INDEX idx_events_conversation ON events(conversation_id);
            """;
        cmd.ExecuteNonQuery();
    }

    // Column mapping: 0=position, 1=event_id, 2=event_type, 3=version,
    // 4=timestamp_nanos, 5=source, 6=content, 7=causes, 8=conversation_id,
    // 9=hash, 10=prev_hash, 11=signature

    private Event ReadEvent(SqlDataReader r)
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

    private List<Event> QueryEvents(string sql, params SqlParameter[] parameters)
    {
        using var cmd = _conn.CreateCommand();
        cmd.CommandText = sql;
        foreach (var p in parameters) cmd.Parameters.Add(p);
        using var reader = cmd.ExecuteReader();
        var result = new List<Event>();
        while (reader.Read()) result.Add(ReadEvent(reader));
        return result;
    }

    private Event? QuerySingleEvent(string sql, params SqlParameter[] parameters)
    {
        using var cmd = _conn.CreateCommand();
        cmd.CommandText = sql;
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
                new SqlParameter("@id", ev.Id.Value));
            if (existing is not null)
            {
                if (existing.Hash != ev.Hash)
                    throw new ChainIntegrityException(0,
                        $"hash mismatch for existing event {ev.Id.Value}");
                return existing;
            }

            // Chain continuity
            var head = QuerySingleEvent(
                "SELECT TOP(1) * FROM events ORDER BY position DESC");
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

            using var cmd = _conn.CreateCommand();
            cmd.CommandText = """
                INSERT INTO events
                (event_id, event_type, version, timestamp_nanos, source,
                 content, causes, conversation_id, hash, prev_hash, signature)
                VALUES (@eid, @et, @v, @ts, @src, @cnt, @cau, @cid, @h, @ph, @sig)
                """;
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
                new SqlParameter("@id", eventId.Value))
                ?? throw new EventNotFoundException(eventId.Value);
        }
    }

    public Option<Event> Head()
    {
        lock (_lock)
        {
            var ev = QuerySingleEvent(
                "SELECT TOP(1) * FROM events ORDER BY position DESC");
            return ev is null ? Option<Event>.None() : Option<Event>.Some(ev);
        }
    }

    public int Count()
    {
        lock (_lock)
        {
            using var cmd = _conn.CreateCommand();
            cmd.CommandText = "SELECT COUNT(*) FROM events";
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
                "SELECT TOP(@lim) * FROM events ORDER BY position DESC",
                new SqlParameter("@lim", limit));
        }
    }

    public List<Event> ByType(EventType type, int limit)
    {
        lock (_lock)
        {
            return QueryEvents(
                "SELECT TOP(@lim) * FROM events WHERE event_type = @t ORDER BY position DESC",
                new SqlParameter("@t", type.Value),
                new SqlParameter("@lim", limit));
        }
    }

    public List<Event> BySource(ActorId source, int limit)
    {
        lock (_lock)
        {
            return QueryEvents(
                "SELECT TOP(@lim) * FROM events WHERE source = @s ORDER BY position DESC",
                new SqlParameter("@s", source.Value),
                new SqlParameter("@lim", limit));
        }
    }

    public List<Event> ByConversation(ConversationId id, int limit)
    {
        lock (_lock)
        {
            return QueryEvents(
                "SELECT TOP(@lim) * FROM events WHERE conversation_id = @c ORDER BY position DESC",
                new SqlParameter("@c", id.Value),
                new SqlParameter("@lim", limit));
        }
    }

    public List<Event> Ancestors(EventId id, int maxDepth)
    {
        lock (_lock)
        {
            var start = QuerySingleEvent(
                "SELECT * FROM events WHERE event_id = @id",
                new SqlParameter("@id", id.Value))
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
                        new SqlParameter("@id", eid));
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
                new SqlParameter("@id", id.Value))
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
                        new SqlParameter("@id", eid));
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
