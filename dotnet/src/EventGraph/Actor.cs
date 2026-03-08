using System.Security.Cryptography;

namespace EventGraph;

// ── ActorType ───────────────────────────────────────────────────────────

/// <summary>What kind of decision-maker an actor is.</summary>
public enum ActorType
{
    Human,
    AI,
    System,
    Committee,
    RulesEngine,
}

// ── ActorStatus ─────────────────────────────────────────────────────────

/// <summary>Lifecycle status of an actor. State machine with enforced valid transitions.</summary>
public enum ActorStatus
{
    Active,
    Suspended,
    Memorial,
}

public static class ActorStatusExtensions
{
    private static readonly IReadOnlyDictionary<ActorStatus, ActorStatus[]> ValidTransitionsMap =
        new Dictionary<ActorStatus, ActorStatus[]>
        {
            [ActorStatus.Active] = [ActorStatus.Suspended, ActorStatus.Memorial],
            [ActorStatus.Suspended] = [ActorStatus.Active, ActorStatus.Memorial],
            [ActorStatus.Memorial] = [], // terminal
        };

    /// <summary>Attempt a state transition. Throws InvalidTransitionException on invalid transitions.</summary>
    public static ActorStatus TransitionTo(this ActorStatus from, ActorStatus to)
    {
        var valid = ValidTransitions(from);
        foreach (var v in valid)
        {
            if (v == to) return to;
        }
        throw new InvalidTransitionException(from.ToString(), to.ToString());
    }

    /// <summary>Returns the set of valid target states from the given status.</summary>
    public static IReadOnlyList<ActorStatus> ValidTransitions(this ActorStatus status)
    {
        return ValidTransitionsMap.TryGetValue(status, out var targets)
            ? targets
            : [];
    }
}

// ── Actor ───────────────────────────────────────────────────────────────

/// <summary>Immutable actor — validated at construction, guaranteed valid for lifetime.</summary>
public sealed class Actor
{
    public ActorId Id { get; }
    public PublicKey PublicKey { get; }
    public string DisplayName { get; }
    public ActorType Type { get; }
    private readonly Dictionary<string, object?> _metadata;
    public IReadOnlyDictionary<string, object?> Metadata => DeepCopyMetadata(_metadata);
    public long CreatedAtNanos { get; }
    public ActorStatus Status { get; }

    public Actor(
        ActorId id,
        PublicKey publicKey,
        string displayName,
        ActorType type,
        Dictionary<string, object?>? metadata,
        long createdAtNanos,
        ActorStatus status)
    {
        Id = id;
        PublicKey = publicKey;
        DisplayName = displayName;
        Type = type;
        _metadata = DeepCopyMetadata(metadata);
        CreatedAtNanos = createdAtNanos;
        Status = status;
    }

    internal Actor WithStatus(ActorStatus status) =>
        new(Id, PublicKey, DisplayName, Type, _metadata, CreatedAtNanos, status);

    internal Actor WithUpdates(ActorUpdate update)
    {
        var md = DeepCopyMetadata(_metadata);
        var name = update.DisplayName ?? DisplayName;

        if (update.Metadata is not null)
        {
            foreach (var kv in update.Metadata)
            {
                md[kv.Key] = DeepCopyValue(kv.Value);
            }
        }

        return new Actor(Id, PublicKey, name, Type, md, CreatedAtNanos, Status);
    }

    // ── Deep copy helpers ───────────────────────────────────────────────

    internal static Dictionary<string, object?> DeepCopyMetadata(IDictionary<string, object?>? src)
    {
        if (src is null || src.Count == 0)
            return new Dictionary<string, object?>();

        var cp = new Dictionary<string, object?>(src.Count);
        foreach (var kv in src)
            cp[kv.Key] = DeepCopyValue(kv.Value);
        return cp;
    }

    internal static object? DeepCopyValue(object? v)
    {
        return v switch
        {
            IDictionary<string, object?> dict => DeepCopyMetadata(dict),
            IList<object?> list => list.Select(DeepCopyValue).ToList(),
            _ => v, // primitives are immutable
        };
    }
}

// ── ActorUpdate ─────────────────────────────────────────────────────────

/// <summary>Describes updates to apply to an actor. Null fields are left unchanged.</summary>
public sealed record ActorUpdate
{
    public string? DisplayName { get; init; }
    public Dictionary<string, object?>? Metadata { get; init; }
}

// ── ActorFilter ─────────────────────────────────────────────────────────

/// <summary>Criteria for listing actors.</summary>
public sealed record ActorFilter
{
    public ActorStatus? Status { get; init; }
    public ActorType? Type { get; init; }
    public int Limit { get; init; } = 100;
    public string? After { get; init; }
}

// ── Page<T> + Cursor ────────────────────────────────────────────────────

/// <summary>Cursor-based pagination result.</summary>
public sealed class Page<T>
{
    public IReadOnlyList<T> Items { get; }
    public string? Cursor { get; }
    public bool HasMore { get; }

    public Page(IReadOnlyList<T> items, string? cursor, bool hasMore)
    {
        Items = items;
        Cursor = cursor;
        HasMore = hasMore;
    }
}

// ── IActorStore ─────────────────────────────────────────────────────────

/// <summary>Actor persistence interface. Separate from Store — single responsibility.</summary>
public interface IActorStore
{
    /// <summary>Register a new actor. Idempotent on public key.</summary>
    Actor Register(PublicKey publicKey, string displayName, ActorType type);

    /// <summary>Get an actor by ID.</summary>
    Actor Get(ActorId id);

    /// <summary>Get an actor by public key.</summary>
    Actor GetByPublicKey(PublicKey publicKey);

    /// <summary>Update an actor's display name and/or metadata.</summary>
    Actor Update(ActorId id, ActorUpdate update);

    /// <summary>List actors with optional filtering and pagination.</summary>
    Page<Actor> List(ActorFilter filter);

    /// <summary>Suspend an actor. Active → Suspended.</summary>
    Actor Suspend(ActorId id, EventId reason);

    /// <summary>Reactivate a suspended actor. Suspended → Active.</summary>
    Actor Reactivate(ActorId id, EventId reason);

    /// <summary>Memorialize an actor. Terminal state.</summary>
    Actor Memorial(ActorId id, EventId reason);
}

// ── InMemoryActorStore ──────────────────────────────────────────────────

/// <summary>Thread-safe in-memory implementation of IActorStore.</summary>
public sealed class InMemoryActorStore : IActorStore
{
    private readonly object _lock = new();
    private readonly Dictionary<string, Actor> _actors = new();      // actorId.Value → Actor
    private readonly Dictionary<string, string> _byKey = new();      // hex(publicKey) → actorId.Value
    private readonly List<string> _ordered = new();                  // insertion order for pagination

    /// <summary>Derive actor_id from SHA-256 of public key bytes: actor_{hex(sha256(pk)[0..16])}.</summary>
    private static ActorId DeriveActorId(PublicKey pk)
    {
        var hash = SHA256.HashData(pk.Bytes.ToArray());
        var hex = Convert.ToHexString(hash, 0, 16).ToLowerInvariant();
        return new ActorId($"actor_{hex}");
    }

    private static string PubKeyHex(PublicKey pk) =>
        Convert.ToHexString(pk.Bytes.ToArray()).ToLowerInvariant();

    public Actor Register(PublicKey publicKey, string displayName, ActorType type)
    {
        lock (_lock)
        {
            var keyHex = PubKeyHex(publicKey);
            if (_byKey.TryGetValue(keyHex, out var existingId))
                return _actors[existingId];

            var id = DeriveActorId(publicKey);
            var now = DateTimeOffset.UtcNow.ToUnixTimeMilliseconds() * 1_000_000;
            var actor = new Actor(id, publicKey, displayName, type, null, now, ActorStatus.Active);

            _actors[id.Value] = actor;
            _byKey[keyHex] = id.Value;
            _ordered.Add(id.Value);
            return actor;
        }
    }

    public Actor Get(ActorId id)
    {
        lock (_lock)
        {
            if (_actors.TryGetValue(id.Value, out var actor))
                return actor;
            throw new ActorNotFoundException(id);
        }
    }

    public Actor GetByPublicKey(PublicKey publicKey)
    {
        lock (_lock)
        {
            var keyHex = PubKeyHex(publicKey);
            if (_byKey.TryGetValue(keyHex, out var actorIdValue))
                return _actors[actorIdValue];
            throw new ActorKeyNotFoundException(keyHex);
        }
    }

    public Actor Update(ActorId id, ActorUpdate update)
    {
        lock (_lock)
        {
            if (!_actors.TryGetValue(id.Value, out var actor))
                throw new ActorNotFoundException(id);

            var updated = actor.WithUpdates(update);
            _actors[id.Value] = updated;
            return updated;
        }
    }

    public Page<Actor> List(ActorFilter filter)
    {
        lock (_lock)
        {
            var limit = filter.Limit <= 0 ? 100 : filter.Limit;
            var startIdx = 0;

            if (filter.After is not null)
            {
                var found = false;
                for (var i = 0; i < _ordered.Count; i++)
                {
                    if (_ordered[i] == filter.After)
                    {
                        startIdx = i + 1;
                        found = true;
                        break;
                    }
                }
                if (!found)
                    return new Page<Actor>([], null, false);
            }

            var items = new List<Actor>();
            for (var i = startIdx; i < _ordered.Count && items.Count < limit; i++)
            {
                var actor = _actors[_ordered[i]];
                if (filter.Status.HasValue && actor.Status != filter.Status.Value)
                    continue;
                if (filter.Type.HasValue && actor.Type != filter.Type.Value)
                    continue;
                items.Add(actor);
            }

            var hasMore = false;
            string? cursor = null;

            if (items.Count == limit)
            {
                // Check if there are more matching items beyond the last one returned
                var lastId = items[^1].Id.Value;
                var lastIdx = _ordered.IndexOf(lastId);
                for (var i = lastIdx + 1; i < _ordered.Count; i++)
                {
                    var actor = _actors[_ordered[i]];
                    if (filter.Status.HasValue && actor.Status != filter.Status.Value)
                        continue;
                    if (filter.Type.HasValue && actor.Type != filter.Type.Value)
                        continue;
                    hasMore = true;
                    break;
                }
                if (hasMore)
                    cursor = lastId;
            }

            return new Page<Actor>(items, cursor, hasMore);
        }
    }

    public Actor Suspend(ActorId id, EventId reason)
    {
        lock (_lock)
        {
            if (!_actors.TryGetValue(id.Value, out var actor))
                throw new ActorNotFoundException(id);

            var newStatus = actor.Status.TransitionTo(ActorStatus.Suspended);
            var updated = actor.WithStatus(newStatus);
            _actors[id.Value] = updated;
            _ = reason; // recorded on the event graph, not stored here
            return updated;
        }
    }

    public Actor Reactivate(ActorId id, EventId reason)
    {
        lock (_lock)
        {
            if (!_actors.TryGetValue(id.Value, out var actor))
                throw new ActorNotFoundException(id);

            var newStatus = actor.Status.TransitionTo(ActorStatus.Active);
            var updated = actor.WithStatus(newStatus);
            _actors[id.Value] = updated;
            _ = reason;
            return updated;
        }
    }

    public Actor Memorial(ActorId id, EventId reason)
    {
        lock (_lock)
        {
            if (!_actors.TryGetValue(id.Value, out var actor))
                throw new ActorNotFoundException(id);

            var newStatus = actor.Status.TransitionTo(ActorStatus.Memorial);
            var updated = actor.WithStatus(newStatus);
            _actors[id.Value] = updated;
            _ = reason;
            return updated;
        }
    }

    /// <summary>Returns the number of registered actors. For testing.</summary>
    public int ActorCount
    {
        get { lock (_lock) { return _actors.Count; } }
    }
}
