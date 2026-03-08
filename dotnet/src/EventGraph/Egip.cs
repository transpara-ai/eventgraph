using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;
using NSec.Cryptography;

namespace EventGraph;

// ── EGIP Enums ───────────────────────────────────────────────────────────

/// <summary>EGIP message types.</summary>
public enum MessageType
{
    Hello,
    Message,
    Receipt,
    Proof,
    Treaty,
    AuthorityRequest,
    Discover,
}

/// <summary>Treaty lifecycle status.</summary>
public enum TreatyStatus
{
    Proposed,
    Active,
    Suspended,
    Terminated,
}

/// <summary>Treaty actions that can be applied.</summary>
public enum TreatyAction
{
    Propose,
    Accept,
    Modify,
    Suspend,
    Terminate,
}

/// <summary>Receipt status for EGIP messages.</summary>
public enum ReceiptStatus
{
    Delivered,
    Processed,
    Rejected,
}

/// <summary>Proof types for chain integrity verification.</summary>
public enum ProofType
{
    ChainSegment,
    EventExistence,
    ChainSummary,
}

/// <summary>Cross-graph event reference relationship types.</summary>
public enum CgerRelationship
{
    CausedBy,
    References,
    RespondsTo,
}

/// <summary>Authority level for inter-system requests.</summary>
public enum EgipAuthorityLevel
{
    Required,
    Recommended,
    Notification,
}

// ── EGIP Constants ───────────────────────────────────────────────────────

/// <summary>EGIP protocol constants.</summary>
public static class EgipConstants
{
    /// <summary>Current protocol version this implementation supports.</summary>
    public const int CurrentProtocolVersion = 1;

    /// <summary>Maximum age of an incoming envelope before rejection (25 hours).</summary>
    public static readonly TimeSpan MaxEnvelopeAge = TimeSpan.FromHours(25);

    /// <summary>Auto-prune interval for dedup tracker.</summary>
    public const int DedupPruneInterval = 1000;

    // ── Trust impact constants (inter-system, more conservative) ─────

    public const double TrustImpactValidProof = 0.02;
    public const double TrustImpactReceiptOnTime = 0.01;
    public const double TrustImpactTreatyHonoured = 0.03;
    public const double TrustImpactTreatyViolated = -0.15;
    public const double TrustImpactInvalidProof = -0.10;
    public const double TrustImpactSignatureInvalid = -0.20;
    public const double TrustImpactNoHelloResponse = -0.05;

    /// <summary>Inter-system trust decay rate per day.</summary>
    public static readonly Score InterSystemDecayRate = new(0.02);

    /// <summary>Maximum positive trust adjustment per interaction.</summary>
    public static readonly Weight InterSystemMaxAdjustment = new(0.05);
}

// ── EGIP Errors ──────────────────────────────────────────────────────────

/// <summary>Base class for all EGIP protocol errors.</summary>
public class EgipException : EventGraphException
{
    public EgipException(string message) : base(message) { }
    public EgipException(string message, Exception inner) : base(message, inner) { }
}

/// <summary>Target system could not be reached.</summary>
public class SystemNotFoundException : EgipException
{
    public SystemUri Uri { get; }
    public SystemNotFoundException(SystemUri uri) : base($"System not found: {uri.Value}") => Uri = uri;
}

/// <summary>Envelope signature failed verification.</summary>
public class EnvelopeSignatureInvalidException : EgipException
{
    public EnvelopeId EnvelopeId { get; }
    public EnvelopeSignatureInvalidException(EnvelopeId envelopeId)
        : base($"Envelope signature invalid: {envelopeId.Value}") => EnvelopeId = envelopeId;
}

/// <summary>A treaty term was violated.</summary>
public class TreatyViolationException : EgipException
{
    public TreatyId TreatyId { get; }
    public string Term { get; }
    public TreatyViolationException(TreatyId treatyId, string term)
        : base($"Treaty {treatyId.Value} violated: {term}")
    {
        TreatyId = treatyId;
        Term = term;
    }
}

/// <summary>System's trust score is too low.</summary>
public class TrustInsufficientException : EgipException
{
    public SystemUri System { get; }
    public Score CurrentScore { get; }
    public Score RequiredScore { get; }
    public TrustInsufficientException(SystemUri system, Score current, Score required)
        : base($"Trust insufficient for {system.Value}: have {current.Value}, need {required.Value}")
    {
        System = system;
        CurrentScore = current;
        RequiredScore = required;
    }
}

/// <summary>Transport-level failure (retryable).</summary>
public class TransportFailureException : EgipException
{
    public SystemUri To { get; }
    public string Reason { get; }
    public TransportFailureException(SystemUri to, string reason)
        : base($"Transport failure to {to.Value}: {reason}")
    {
        To = to;
        Reason = reason;
    }
}

/// <summary>Replay detected — envelope already processed.</summary>
public class DuplicateEnvelopeException : EgipException
{
    public EnvelopeId EnvelopeId { get; }
    public DuplicateEnvelopeException(EnvelopeId envelopeId)
        : base($"Duplicate envelope: {envelopeId.Value}") => EnvelopeId = envelopeId;
}

/// <summary>Referenced treaty does not exist.</summary>
public class TreatyNotFoundException : EgipException
{
    public TreatyId TreatyId { get; }
    public TreatyNotFoundException(TreatyId treatyId)
        : base($"Treaty not found: {treatyId.Value}") => TreatyId = treatyId;
}

/// <summary>No common protocol version exists.</summary>
public class VersionIncompatibleException : EgipException
{
    public int[] Local { get; }
    public int[] Remote { get; }
    public VersionIncompatibleException(int[] local, int[] remote)
        : base($"No compatible protocol version: local [{string.Join(",", local)}], remote [{string.Join(",", remote)}]")
    {
        Local = local;
        Remote = remote;
    }
}

// ── IIdentity ────────────────────────────────────────────────────────────

/// <summary>Cryptographic identity for an EGIP system.</summary>
public interface IIdentity
{
    /// <summary>This system's address.</summary>
    SystemUri SystemUri { get; }

    /// <summary>Ed25519 public key (32 bytes).</summary>
    PublicKey PublicKey { get; }

    /// <summary>Produce an Ed25519 signature of the given data.</summary>
    Signature Sign(byte[] data);

    /// <summary>Check an Ed25519 signature against a public key.</summary>
    bool Verify(PublicKey publicKey, byte[] data, Signature signature);
}

// ── SystemIdentity ───────────────────────────────────────────────────────

/// <summary>Ed25519-based implementation of IIdentity using NSec.Cryptography.</summary>
public sealed class SystemIdentity : IIdentity, IDisposable
{
    private static readonly SignatureAlgorithm Ed25519Algorithm = SignatureAlgorithm.Ed25519;

    private readonly Key _key;
    private bool _disposed;

    public SystemUri SystemUri { get; }
    public PublicKey PublicKey { get; }
    public DateTimeOffset CreatedAt { get; }

    private SystemIdentity(SystemUri uri, Key key, PublicKey publicKey)
    {
        SystemUri = uri;
        _key = key;
        PublicKey = publicKey;
        CreatedAt = DateTimeOffset.UtcNow;
    }

    /// <summary>Create a new identity with a fresh Ed25519 keypair.</summary>
    public static SystemIdentity Generate(SystemUri uri)
    {
        var key = Key.Create(Ed25519Algorithm, new KeyCreationParameters { ExportPolicy = KeyExportPolicies.AllowPlaintextExport });
        var pubBytes = key.PublicKey.Export(KeyBlobFormat.RawPublicKey);
        var publicKey = new PublicKey(pubBytes);
        return new SystemIdentity(uri, key, publicKey);
    }

    /// <summary>Create an identity from an existing Ed25519 private key seed (raw 32 bytes).</summary>
    public static SystemIdentity FromPrivateKey(SystemUri uri, byte[] privateKeySeed)
    {
        var key = Key.Import(Ed25519Algorithm, privateKeySeed, KeyBlobFormat.RawPrivateKey,
            new KeyCreationParameters { ExportPolicy = KeyExportPolicies.AllowPlaintextExport });
        var pubBytes = key.PublicKey.Export(KeyBlobFormat.RawPublicKey);
        var publicKey = new PublicKey(pubBytes);
        return new SystemIdentity(uri, key, publicKey);
    }

    public Signature Sign(byte[] data)
    {
        ObjectDisposedException.ThrowIf(_disposed, this);
        var sig = Ed25519Algorithm.Sign(_key, data);
        return new Signature(sig);
    }

    public bool Verify(PublicKey publicKey, byte[] data, Signature signature)
    {
        ObjectDisposedException.ThrowIf(_disposed, this);
        var pubKeyBytes = publicKey.Bytes.ToArray();
        var nsecPubKey = NSec.Cryptography.PublicKey.Import(Ed25519Algorithm, pubKeyBytes, KeyBlobFormat.RawPublicKey);
        return Ed25519Algorithm.Verify(nsecPubKey, data, signature.Bytes);
    }

    public void Dispose()
    {
        if (!_disposed)
        {
            _key.Dispose();
            _disposed = true;
        }
    }
}

// ── ITransport ───────────────────────────────────────────────────────────

/// <summary>Pluggable transport layer for EGIP communication.</summary>
public interface ITransport
{
    /// <summary>Deliver an envelope to the target system.</summary>
    Task<ReceiptPayload?> SendAsync(SystemUri to, Envelope envelope, CancellationToken ct = default);

    /// <summary>Returns an async enumerable of incoming envelopes.</summary>
    IAsyncEnumerable<IncomingEnvelope> ListenAsync(CancellationToken ct = default);
}

/// <summary>Wraps an envelope received from a remote system.</summary>
public sealed record IncomingEnvelope(Envelope? Envelope, Exception? Error);

// ── Envelope ─────────────────────────────────────────────────────────────

/// <summary>Signed message container for all EGIP communication.</summary>
public sealed class Envelope
{
    public int ProtocolVersion { get; init; }
    public EnvelopeId Id { get; init; }
    public SystemUri From { get; init; }
    public SystemUri To { get; init; }
    public MessageType Type { get; init; }
    public IMessagePayload Payload { get; init; } = null!;
    public DateTimeOffset Timestamp { get; init; }
    public Signature Signature { get; init; }
    public Option<EnvelopeId> InReplyTo { get; init; }

    /// <summary>Returns the canonical string representation for signing.</summary>
    public string CanonicalForm()
    {
        var payloadJson = PayloadToCanonicalJson(Payload);

        var msgType = Type.ToString().ToLowerInvariant();
        var nanos = Timestamp.ToUnixTimeMilliseconds() * 1_000_000;

        var inReplyTo = InReplyTo.IsSome ? InReplyTo.Unwrap().Value : "";

        return $"{ProtocolVersion}|{Id.Value}|{From.Value}|{To.Value}|{msgType}|{nanos}|{inReplyTo}|{payloadJson}";
    }

    /// <summary>Serialize a payload to canonical JSON with sorted keys.</summary>
    private static string PayloadToCanonicalJson(IMessagePayload payload)
    {
        // Manually serialize to avoid ReadOnlySpan<byte> issues with System.Text.Json
        var dict = PayloadToDictionary(payload);
        var json = JsonSerializer.Serialize(dict);
        // Round-trip to normalize key order
        using var doc = JsonDocument.Parse(json);
        return SerializeCanonical(doc.RootElement);
    }

    private static Dictionary<string, object?> PayloadToDictionary(IMessagePayload payload)
    {
        return payload switch
        {
            HelloPayload h => new Dictionary<string, object?>
            {
                ["system_uri"] = h.SystemUri.Value,
                ["public_key"] = h.PublicKey.ToString(),
                ["protocol_versions"] = h.ProtocolVersions,
                ["capabilities"] = h.Capabilities,
                ["chain_length"] = h.ChainLength,
            },
            MessagePayloadContent m => new Dictionary<string, object?>
            {
                ["content"] = m.ContentJson,
                ["content_type"] = m.ContentType,
                ["conversation_id"] = m.ConversationId.IsSome ? m.ConversationId.Unwrap().Value : null,
                ["cgers"] = m.Cgers.Select(c => new Dictionary<string, object?>
                {
                    ["local_event_id"] = c.LocalEventId.Value,
                    ["remote_system"] = c.RemoteSystem.Value,
                    ["remote_event_id"] = c.RemoteEventId,
                    ["remote_hash"] = c.RemoteHash.Value,
                    ["relationship"] = c.Relationship.ToString(),
                    ["verified"] = c.Verified,
                }).ToArray(),
            },
            ReceiptPayload r => new Dictionary<string, object?>
            {
                ["envelope_id"] = r.EnvelopeId.Value,
                ["status"] = r.Status.ToString(),
                ["local_event_id"] = r.LocalEventId.IsSome ? r.LocalEventId.Unwrap().Value : null,
                ["reason"] = r.Reason.IsSome ? r.Reason.Unwrap() : null,
                ["signature"] = r.Signature.ToString(),
            },
            ProofPayload p => new Dictionary<string, object?>
            {
                ["proof_type"] = p.ProofType.ToString(),
                ["data"] = new Dictionary<string, object?> { ["type"] = p.Data.GetType().Name },
            },
            TreatyPayload t => new Dictionary<string, object?>
            {
                ["treaty_id"] = t.TreatyId.Value,
                ["action"] = t.Action.ToString(),
                ["terms"] = t.Terms.Select(term => new Dictionary<string, object?>
                {
                    ["scope"] = term.Scope.Value,
                    ["policy"] = term.Policy,
                    ["symmetric"] = term.Symmetric,
                }).ToArray(),
                ["reason"] = t.Reason.IsSome ? t.Reason.Unwrap() : null,
            },
            AuthorityRequestPayload a => new Dictionary<string, object?>
            {
                ["action"] = a.Action.Value,
                ["actor"] = a.Actor.Value,
                ["level"] = a.Level.ToString(),
                ["justification"] = a.Justification,
                ["treaty_id"] = a.TreatyId.IsSome ? a.TreatyId.Unwrap().Value : null,
            },
            DiscoverPayload d => new Dictionary<string, object?>
            {
                ["query"] = new Dictionary<string, object?>
                {
                    ["capabilities"] = d.Query.Capabilities,
                    ["min_trust"] = d.Query.MinTrust.IsSome ? d.Query.MinTrust.Unwrap().Value : null,
                },
                ["results"] = d.Results.Select(r => new Dictionary<string, object?>
                {
                    ["system_uri"] = r.SystemUri.Value,
                    ["public_key"] = r.PublicKey.ToString(),
                    ["capabilities"] = r.Capabilities,
                    ["trust_score"] = r.TrustScore.Value,
                }).ToArray(),
            },
            _ => throw new ArgumentException($"Unknown payload type: {payload.GetType().Name}"),
        };
    }

    /// <summary>Sign this envelope using the given identity. Returns a new envelope with the signature set.</summary>
    public static Envelope SignEnvelope(Envelope env, IIdentity identity)
    {
        var canonical = env.CanonicalForm();
        var sig = identity.Sign(Encoding.UTF8.GetBytes(canonical));

        return new Envelope
        {
            ProtocolVersion = env.ProtocolVersion,
            Id = env.Id,
            From = env.From,
            To = env.To,
            Type = env.Type,
            Payload = env.Payload,
            Timestamp = env.Timestamp,
            Signature = sig,
            InReplyTo = env.InReplyTo,
        };
    }

    /// <summary>Verify the envelope's signature against the given public key.</summary>
    public static bool VerifyEnvelope(Envelope env, IIdentity identity, PublicKey publicKey)
    {
        var canonical = env.CanonicalForm();
        return identity.Verify(publicKey, Encoding.UTF8.GetBytes(canonical), env.Signature);
    }

    /// <summary>Serialize a JsonElement to canonical (sorted-key) JSON.</summary>
    private static string SerializeCanonical(JsonElement element)
    {
        var sb = new StringBuilder();
        WriteCanonicalElement(element, sb);
        return sb.ToString();
    }

    private static void WriteCanonicalElement(JsonElement element, StringBuilder sb)
    {
        switch (element.ValueKind)
        {
            case JsonValueKind.Object:
                var props = new SortedDictionary<string, JsonElement>(StringComparer.Ordinal);
                foreach (var prop in element.EnumerateObject())
                    props[prop.Name] = prop.Value;

                sb.Append('{');
                var first = true;
                foreach (var kvp in props)
                {
                    if (!first) sb.Append(',');
                    first = false;
                    sb.Append('"');
                    sb.Append(kvp.Key);
                    sb.Append("\":");
                    WriteCanonicalElement(kvp.Value, sb);
                }
                sb.Append('}');
                break;

            case JsonValueKind.Array:
                sb.Append('[');
                var firstArr = true;
                foreach (var item in element.EnumerateArray())
                {
                    if (!firstArr) sb.Append(',');
                    firstArr = false;
                    WriteCanonicalElement(item, sb);
                }
                sb.Append(']');
                break;

            case JsonValueKind.String:
                sb.Append('"');
                sb.Append(element.GetString()!.Replace("\\", "\\\\").Replace("\"", "\\\""));
                sb.Append('"');
                break;

            case JsonValueKind.Number:
                sb.Append(element.GetRawText());
                break;

            case JsonValueKind.True:
                sb.Append("true");
                break;

            case JsonValueKind.False:
                sb.Append("false");
                break;

            case JsonValueKind.Null:
                sb.Append("null");
                break;
        }
    }
}

// ── JSON serialization context ───────────────────────────────────────────

internal static class EgipJsonContext
{
    public static JsonSerializerOptions Options { get; } = new()
    {
        PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower,
        DefaultIgnoreCondition = JsonIgnoreCondition.WhenWritingNull,
        Converters = { new JsonStringEnumConverter(JsonNamingPolicy.CamelCase) },
    };
}

// ── Payload marker interface ─────────────────────────────────────────────

/// <summary>Marker interface for all EGIP message payloads.</summary>
public interface IMessagePayload { }

// ── Payload Types ────────────────────────────────────────────────────────

/// <summary>Payload for HELLO messages.</summary>
public sealed record HelloPayload(
    SystemUri SystemUri,
    PublicKey PublicKey,
    int[] ProtocolVersions,
    string[] Capabilities,
    int ChainLength
) : IMessagePayload;

/// <summary>Payload for MESSAGE messages.</summary>
public sealed record MessagePayloadContent(
    string ContentJson,
    string ContentType,
    Option<ConversationId> ConversationId,
    Cger[] Cgers
) : IMessagePayload;

/// <summary>Payload for RECEIPT messages.</summary>
public sealed record ReceiptPayload(
    EnvelopeId EnvelopeId,
    ReceiptStatus Status,
    Option<EventId> LocalEventId,
    Option<string> Reason,
    Signature Signature
) : IMessagePayload;

/// <summary>Payload for PROOF messages.</summary>
public sealed record ProofPayload(
    ProofType ProofType,
    IProofData Data
) : IMessagePayload;

/// <summary>Payload for TREATY messages.</summary>
public sealed record TreatyPayload(
    TreatyId TreatyId,
    TreatyAction Action,
    TreatyTerm[] Terms,
    Option<string> Reason
) : IMessagePayload;

/// <summary>Payload for AUTHORITY_REQUEST messages.</summary>
public sealed record AuthorityRequestPayload(
    DomainScope Action,
    ActorId Actor,
    EgipAuthorityLevel Level,
    string Justification,
    Option<TreatyId> TreatyId
) : IMessagePayload;

/// <summary>Payload for DISCOVER messages.</summary>
public sealed record DiscoverPayload(
    DiscoverQuery Query,
    DiscoverResult[] Results
) : IMessagePayload;

/// <summary>Discovery query specifying desired capabilities.</summary>
public sealed record DiscoverQuery(
    string[] Capabilities,
    Option<Score> MinTrust
);

/// <summary>A single discovery result.</summary>
public sealed record DiscoverResult(
    SystemUri SystemUri,
    PublicKey PublicKey,
    string[] Capabilities,
    Score TrustScore
);

// ── CGER ─────────────────────────────────────────────────────────────────

/// <summary>Cross-Graph Event Reference with verification tracking.</summary>
public sealed record Cger(
    EventId LocalEventId,
    SystemUri RemoteSystem,
    string RemoteEventId,
    Hash RemoteHash,
    CgerRelationship Relationship,
    bool Verified
);

// ── Proof Data ───────────────────────────────────────────────────────────

/// <summary>Marker interface for proof-type-specific data.</summary>
public interface IProofData { }

/// <summary>Chain segment proof — contiguous portion of the hash chain.</summary>
public sealed record ChainSegmentProof(
    Event[] Events,
    Hash StartHash,
    Hash EndHash
) : IProofData;

/// <summary>Event existence proof — proves a specific event exists in the chain.</summary>
public sealed record EventExistenceProof(
    Event Event,
    Hash PrevHash,
    Option<Hash> NextHash,
    int Position,
    int ChainLength
) : IProofData;

/// <summary>Chain summary proof — high-level integrity attestation.</summary>
public sealed record ChainSummaryProof(
    int Length,
    Hash HeadHash,
    Hash GenesisHash,
    DateTimeOffset Timestamp
) : IProofData;

// ── Treaty ───────────────────────────────────────────────────────────────

/// <summary>A bilateral governance agreement between two systems.</summary>
public sealed class Treaty
{
    public TreatyId Id { get; }
    public SystemUri SystemA { get; }
    public SystemUri SystemB { get; }
    public TreatyStatus Status { get; private set; }
    public TreatyTerm[] Terms { get; internal set; }
    public DateTimeOffset CreatedAt { get; }
    public DateTimeOffset UpdatedAt { get; private set; }

    private static readonly Dictionary<TreatyStatus, TreatyStatus[]> ValidTransitions = new()
    {
        [TreatyStatus.Proposed] = [TreatyStatus.Active, TreatyStatus.Terminated],
        [TreatyStatus.Active] = [TreatyStatus.Suspended, TreatyStatus.Terminated],
        [TreatyStatus.Suspended] = [TreatyStatus.Active, TreatyStatus.Terminated],
        [TreatyStatus.Terminated] = [], // terminal
    };

    public Treaty(TreatyId id, SystemUri systemA, SystemUri systemB, TreatyTerm[] terms)
    {
        Id = id;
        SystemA = systemA;
        SystemB = systemB;
        Status = TreatyStatus.Proposed;
        Terms = terms;
        CreatedAt = DateTimeOffset.UtcNow;
        UpdatedAt = CreatedAt;
    }

    /// <summary>Attempt to transition to a new status.</summary>
    public void Transition(TreatyStatus to)
    {
        var allowed = ValidTransitions[Status];
        if (!allowed.Contains(to))
            throw new InvalidTransitionException(Status.ToString(), to.ToString());
        Status = to;
        UpdatedAt = DateTimeOffset.UtcNow;
    }

    /// <summary>Apply a treaty action and transition accordingly.</summary>
    public void ApplyAction(TreatyAction action)
    {
        switch (action)
        {
            case TreatyAction.Accept:
                Transition(TreatyStatus.Active);
                break;
            case TreatyAction.Suspend:
                Transition(TreatyStatus.Suspended);
                break;
            case TreatyAction.Terminate:
                Transition(TreatyStatus.Terminated);
                break;
            case TreatyAction.Modify:
                if (Status != TreatyStatus.Active)
                    throw new InvalidTransitionException(Status.ToString(), "Modify (requires Active)");
                UpdatedAt = DateTimeOffset.UtcNow;
                break;
            case TreatyAction.Propose:
                throw new InvalidOperationException("Cannot apply Propose to existing treaty");
            default:
                throw new ArgumentOutOfRangeException(nameof(action), action, "Unknown treaty action");
        }
    }
}

/// <summary>A single term of a bilateral treaty.</summary>
public sealed record TreatyTerm(
    DomainScope Scope,
    string Policy,
    bool Symmetric
);

// ── TreatyTerm JSON converter ────────────────────────────────────────────

// Treaty and TreatyTerm need custom serialization for the canonical form
// since DomainScope is a value object. The default JSON context handles this
// via JsonNamingPolicy.SnakeCaseLower and JsonStringEnumConverter.

// ── PeerRecord & PeerStore ───────────────────────────────────────────────

/// <summary>Tracks the state of a known remote system.</summary>
public sealed class PeerRecord
{
    public SystemUri SystemUri { get; }
    public PublicKey PublicKey { get; }
    public Score Trust { get; internal set; }
    public string[] Capabilities { get; internal set; }
    public int NegotiatedVersion { get; internal set; }
    public DateTimeOffset LastSeen { get; internal set; }
    public DateTimeOffset FirstSeen { get; }
    public DateTimeOffset LastDecayedAt { get; internal set; }

    public PeerRecord(SystemUri systemUri, PublicKey publicKey, string[] capabilities, int negotiatedVersion)
    {
        var now = DateTimeOffset.UtcNow;
        SystemUri = systemUri;
        PublicKey = publicKey;
        Trust = new Score(0.0);
        Capabilities = capabilities;
        NegotiatedVersion = negotiatedVersion;
        LastSeen = now;
        FirstSeen = now;
        LastDecayedAt = now;
    }
}

/// <summary>Thread-safe store for known peer systems and their trust scores.</summary>
public sealed class PeerStore
{
    private readonly ReaderWriterLockSlim _lock = new();
    private readonly Dictionary<string, PeerRecord> _peers = new();

    /// <summary>Add or update a peer from a HELLO exchange.</summary>
    public PeerRecord Register(SystemUri uri, PublicKey publicKey, string[] capabilities, int negotiatedVersion)
    {
        _lock.EnterWriteLock();
        try
        {
            var key = uri.Value;
            var now = DateTimeOffset.UtcNow;

            if (_peers.TryGetValue(key, out var existing))
            {
                // Do NOT overwrite PublicKey — prevents key-substitution attacks.
                existing.Capabilities = capabilities;
                existing.NegotiatedVersion = negotiatedVersion;
                existing.LastSeen = now;
                return existing;
            }

            var record = new PeerRecord(uri, publicKey, capabilities, negotiatedVersion);
            _peers[key] = record;
            return record;
        }
        finally
        {
            _lock.ExitWriteLock();
        }
    }

    /// <summary>Get a peer record by URI.</summary>
    public (PeerRecord? Record, bool Found) Get(SystemUri uri)
    {
        _lock.EnterReadLock();
        try
        {
            if (_peers.TryGetValue(uri.Value, out var record))
                return (record, true);
            return (null, false);
        }
        finally
        {
            _lock.ExitReadLock();
        }
    }

    /// <summary>Adjust a peer's trust score by the given delta, clamped to [0,1].
    /// Positive deltas are capped at InterSystemMaxAdjustment.
    /// Negative deltas are applied uncapped.</summary>
    public (Score Score, bool Found) UpdateTrust(SystemUri uri, double delta)
    {
        _lock.EnterWriteLock();
        try
        {
            if (!_peers.TryGetValue(uri.Value, out var record))
                return (new Score(0.0), false);

            // Positive trust accumulates gradually; negative hits immediately.
            if (delta > 0)
            {
                var maxAdj = EgipConstants.InterSystemMaxAdjustment.Value;
                if (delta > maxAdj)
                    delta = maxAdj;
            }

            var newVal = Math.Clamp(record.Trust.Value + delta, 0.0, 1.0);
            record.Trust = new Score(newVal);
            record.LastSeen = DateTimeOffset.UtcNow;
            return (record.Trust, true);
        }
        finally
        {
            _lock.ExitWriteLock();
        }
    }

    /// <summary>Apply time-based trust decay to all peers.</summary>
    public void DecayAll()
    {
        _lock.EnterWriteLock();
        try
        {
            var now = DateTimeOffset.UtcNow;
            var decayPerDay = EgipConstants.InterSystemDecayRate.Value;

            foreach (var record in _peers.Values)
            {
                var daysSince = (now - record.LastDecayedAt).TotalDays;
                if (daysSince <= 0) continue;

                var decay = decayPerDay * daysSince;
                var newVal = Math.Max(0.0, record.Trust.Value - decay);
                record.Trust = new Score(newVal);
                record.LastDecayedAt = now;
            }
        }
        finally
        {
            _lock.ExitWriteLock();
        }
    }

    /// <summary>Get all peer records.</summary>
    public List<PeerRecord> All()
    {
        _lock.EnterReadLock();
        try
        {
            return _peers.Values.ToList();
        }
        finally
        {
            _lock.ExitReadLock();
        }
    }
}

// ── TreatyStore ──────────────────────────────────────────────────────────

/// <summary>Thread-safe store for bilateral treaties.</summary>
public sealed class TreatyStore
{
    private readonly ReaderWriterLockSlim _lock = new();
    private readonly Dictionary<string, Treaty> _treaties = new();

    /// <summary>Store or update a treaty.</summary>
    public void Put(Treaty treaty)
    {
        _lock.EnterWriteLock();
        try
        {
            _treaties[treaty.Id.Value] = treaty;
        }
        finally
        {
            _lock.ExitWriteLock();
        }
    }

    /// <summary>Get a treaty by ID.</summary>
    public (Treaty? Treaty, bool Found) Get(TreatyId id)
    {
        _lock.EnterReadLock();
        try
        {
            if (_treaties.TryGetValue(id.Value, out var treaty))
                return (treaty, true);
            return (null, false);
        }
        finally
        {
            _lock.ExitReadLock();
        }
    }

    /// <summary>Perform a read-modify-write on a treaty under a single write lock.</summary>
    public void Apply(TreatyId id, Action<Treaty> action)
    {
        _lock.EnterWriteLock();
        try
        {
            if (!_treaties.TryGetValue(id.Value, out var treaty))
                throw new TreatyNotFoundException(id);
            action(treaty);
        }
        finally
        {
            _lock.ExitWriteLock();
        }
    }

    /// <summary>Get all treaties involving a given system URI.</summary>
    public List<Treaty> BySystem(SystemUri uri)
    {
        _lock.EnterReadLock();
        try
        {
            return _treaties.Values
                .Where(t => t.SystemA.Value == uri.Value || t.SystemB.Value == uri.Value)
                .ToList();
        }
        finally
        {
            _lock.ExitReadLock();
        }
    }

    /// <summary>Get all active treaties.</summary>
    public List<Treaty> Active()
    {
        _lock.EnterReadLock();
        try
        {
            return _treaties.Values
                .Where(t => t.Status == TreatyStatus.Active)
                .ToList();
        }
        finally
        {
            _lock.ExitReadLock();
        }
    }
}

// ── EnvelopeDedup ────────────────────────────────────────────────────────

/// <summary>Replay protection by tracking seen envelope IDs. Auto-prunes expired entries.</summary>
public sealed class EnvelopeDedup
{
    private readonly object _lock = new();
    private readonly Dictionary<string, DateTimeOffset> _seen = new();
    private readonly TimeSpan _ttl;
    private long _checkCount;

    public EnvelopeDedup()
        : this(EgipConstants.MaxEnvelopeAge + TimeSpan.FromHours(1)) { }

    public EnvelopeDedup(TimeSpan ttl) => _ttl = ttl;

    /// <summary>Returns true if the envelope ID has not been seen before.
    /// Records the ID and returns false on subsequent calls.
    /// Periodically prunes expired entries.</summary>
    public bool Check(EnvelopeId id)
    {
        lock (_lock)
        {
            var key = id.Value;
            if (_seen.ContainsKey(key))
                return false;

            _seen[key] = DateTimeOffset.UtcNow;

            _checkCount++;
            if (_checkCount % EgipConstants.DedupPruneInterval == 0)
                PruneLocked();

            return true;
        }
    }

    /// <summary>Remove expired entries older than TTL.</summary>
    public int Prune()
    {
        lock (_lock)
        {
            return PruneLocked();
        }
    }

    /// <summary>Number of tracked envelope IDs.</summary>
    public int Size
    {
        get
        {
            lock (_lock)
            {
                return _seen.Count;
            }
        }
    }

    private int PruneLocked()
    {
        var cutoff = DateTimeOffset.UtcNow - _ttl;
        var toRemove = _seen.Where(kvp => kvp.Value < cutoff).Select(kvp => kvp.Key).ToList();
        foreach (var key in toRemove)
            _seen.Remove(key);
        return toRemove.Count;
    }
}

// ── Proof verification ───────────────────────────────────────────────────

/// <summary>Proof verification and generation utilities.</summary>
public static class ProofVerification
{
    /// <summary>Verify a chain segment is internally consistent.</summary>
    public static bool VerifyChainSegment(ChainSegmentProof proof)
    {
        if (proof.Events.Length == 0) return false;

        // First event's PrevHash must match StartHash
        if (proof.Events[0].PrevHash != proof.StartHash) return false;

        // Internal hash chain continuity
        for (var i = 1; i < proof.Events.Length; i++)
        {
            if (proof.Events[i].PrevHash != proof.Events[i - 1].Hash) return false;
        }

        // Last event's hash must match EndHash
        if (proof.Events[^1].Hash != proof.EndHash) return false;

        return true;
    }

    /// <summary>Verify basic properties of an event existence proof.</summary>
    public static bool VerifyEventExistence(EventExistenceProof proof)
    {
        if (proof.Event.PrevHash != proof.PrevHash) return false;
        if (proof.Position < 0 || proof.Position >= proof.ChainLength) return false;
        if (proof.Event.Hash == Hash.Zero()) return false;
        return true;
    }

    /// <summary>Validate any proof payload by dispatching to the appropriate verifier.</summary>
    public static bool ValidateProof(ProofPayload payload)
    {
        return payload.Data switch
        {
            ChainSegmentProof csp => VerifyChainSegment(csp),
            EventExistenceProof eep => VerifyEventExistence(eep),
            ChainSummaryProof csp => csp.Length > 0,
            _ => throw new ArgumentException($"Unknown proof data type: {payload.Data.GetType().Name}"),
        };
    }

    /// <summary>Get the ProofType for a given proof data object.</summary>
    public static ProofType ProofTypeFromData(IProofData data)
    {
        return data switch
        {
            ChainSegmentProof => ProofType.ChainSegment,
            EventExistenceProof => ProofType.EventExistence,
            ChainSummaryProof => ProofType.ChainSummary,
            _ => throw new ArgumentException($"Unknown proof data type: {data.GetType().Name}"),
        };
    }
}

// ── NegotiateVersion ─────────────────────────────────────────────────────

/// <summary>EGIP version negotiation.</summary>
public static class EgipVersionNegotiation
{
    /// <summary>Find the highest protocol version both systems support.
    /// Returns None if no common version exists.</summary>
    public static Option<int> NegotiateVersion(int[] local, int[] remote)
    {
        var best = -1;
        foreach (var l in local)
        {
            foreach (var r in remote)
            {
                if (l == r && l > best)
                    best = l;
            }
        }
        return best < 0 ? Option<int>.None() : Option<int>.Some(best);
    }
}

// ── Handler ──────────────────────────────────────────────────────────────

/// <summary>EGIP protocol handler: HELLO handshake, message dispatch,
/// replay deduplication, and trust updates.</summary>
public sealed class EgipHandler
{
    private readonly IIdentity _identity;
    private readonly ITransport _transport;
    private readonly PeerStore _peers;
    private readonly TreatyStore _treaties;
    private readonly EnvelopeDedup _dedup;

    /// <summary>Protocol versions this system supports.</summary>
    public int[] LocalProtocolVersions { get; set; } = [EgipConstants.CurrentProtocolVersion];

    /// <summary>Capabilities this system advertises.</summary>
    public string[] Capabilities { get; set; } = ["treaty", "proof"];

    /// <summary>Returns the current chain length (for HELLO). May be null if not applicable.</summary>
    public Func<int>? ChainLength { get; set; }

    /// <summary>Called when a verified MESSAGE envelope arrives.</summary>
    public Func<SystemUri, MessagePayloadContent, Task>? OnMessage { get; set; }

    /// <summary>Called when a verified AUTHORITY_REQUEST arrives.</summary>
    public Func<SystemUri, AuthorityRequestPayload, Task>? OnAuthorityRequest { get; set; }

    /// <summary>Called when a verified DISCOVER query arrives. Returns results.</summary>
    public Func<SystemUri, DiscoverQuery, Task<DiscoverResult[]>>? OnDiscover { get; set; }

    public EgipHandler(IIdentity identity, ITransport transport, PeerStore peers, TreatyStore treaties)
    {
        _identity = identity;
        _transport = transport;
        _peers = peers;
        _treaties = treaties;
        _dedup = new EnvelopeDedup();
    }

    /// <summary>Perform the HELLO handshake with a remote system.</summary>
    public async Task HelloAsync(SystemUri to, CancellationToken ct = default)
    {
        var chainLen = ChainLength?.Invoke() ?? 0;

        var envId = new EnvelopeId(Guid.NewGuid().ToString());

        var env = new Envelope
        {
            ProtocolVersion = EgipConstants.CurrentProtocolVersion,
            Id = envId,
            From = _identity.SystemUri,
            To = to,
            Type = MessageType.Hello,
            Payload = new HelloPayload(
                _identity.SystemUri,
                _identity.PublicKey,
                LocalProtocolVersions,
                Capabilities,
                chainLen),
            Timestamp = DateTimeOffset.UtcNow,
            Signature = new Signature(new byte[64]), // placeholder
            InReplyTo = Option<EnvelopeId>.None(),
        };

        var signed = Envelope.SignEnvelope(env, _identity);

        ReceiptPayload? receipt;
        try
        {
            receipt = await _transport.SendAsync(to, signed, ct);
        }
        catch (Exception ex)
        {
            _peers.UpdateTrust(to, EgipConstants.TrustImpactNoHelloResponse);
            throw new TransportFailureException(to, ex.Message);
        }

        if (receipt is not null && receipt.Status == ReceiptStatus.Rejected)
        {
            var reason = receipt.Reason.IsSome ? receipt.Reason.Unwrap() : "";
            throw new EgipException($"Hello rejected: {reason}");
        }
    }

    /// <summary>Process an incoming envelope: timestamp check, dedup, signature verification, dispatch, trust update.</summary>
    public async Task HandleIncomingAsync(Envelope env, CancellationToken ct = default)
    {
        // Timestamp freshness
        var age = DateTimeOffset.UtcNow - env.Timestamp;
        if (age > EgipConstants.MaxEnvelopeAge || age < TimeSpan.FromMinutes(-5))
            throw new EgipException($"Envelope timestamp out of range: age {age}");

        // Replay deduplication
        if (!_dedup.Check(env.Id))
            throw new DuplicateEnvelopeException(env.Id);

        // Look up sender's public key
        var (peer, known) = _peers.Get(env.From);

        // For HELLO messages, use the public key from the payload (TOFU model).
        PublicKey pubKey;
        if (env.Type == MessageType.Hello)
        {
            if (env.Payload is not HelloPayload hello)
                throw new EgipException($"Invalid hello payload type: {env.Payload.GetType().Name}");
            pubKey = hello.PublicKey;
        }
        else
        {
            if (!known)
                throw new SystemNotFoundException(env.From);
            pubKey = peer!.PublicKey;
        }

        // Verify signature
        if (!Envelope.VerifyEnvelope(env, _identity, pubKey))
        {
            _peers.UpdateTrust(env.From, EgipConstants.TrustImpactSignatureInvalid);
            throw new EnvelopeSignatureInvalidException(env.Id);
        }

        // Dispatch by message type
        switch (env.Type)
        {
            case MessageType.Hello:
                HandleHello(env);
                break;
            case MessageType.Message:
                await HandleMessageAsync(env);
                break;
            case MessageType.Receipt:
                HandleReceipt(env);
                break;
            case MessageType.Proof:
                HandleProof(env);
                break;
            case MessageType.Treaty:
                HandleTreaty(env);
                break;
            case MessageType.AuthorityRequest:
                await HandleAuthorityRequestAsync(env);
                break;
            case MessageType.Discover:
                await HandleDiscoverAsync(env, ct);
                break;
            default:
                throw new EgipException($"Unknown message type: {env.Type}");
        }
    }

    private void HandleHello(Envelope env)
    {
        var hello = (HelloPayload)env.Payload;

        var version = EgipVersionNegotiation.NegotiateVersion(LocalProtocolVersions, hello.ProtocolVersions);
        if (!version.IsSome)
            throw new VersionIncompatibleException(LocalProtocolVersions, hello.ProtocolVersions);

        _peers.Register(hello.SystemUri, hello.PublicKey, hello.Capabilities, version.Unwrap());
    }

    private async Task HandleMessageAsync(Envelope env)
    {
        var msg = (MessagePayloadContent)env.Payload;
        _peers.UpdateTrust(env.From, EgipConstants.TrustImpactReceiptOnTime);

        if (OnMessage is not null)
            await OnMessage(env.From, msg);
    }

    private void HandleReceipt(Envelope env)
    {
        var receipt = (ReceiptPayload)env.Payload;
        if (receipt.Status is ReceiptStatus.Processed or ReceiptStatus.Delivered)
            _peers.UpdateTrust(env.From, EgipConstants.TrustImpactReceiptOnTime);
    }

    private void HandleProof(Envelope env)
    {
        var proof = (ProofPayload)env.Payload;
        var valid = ProofVerification.ValidateProof(proof);

        _peers.UpdateTrust(env.From, valid
            ? EgipConstants.TrustImpactValidProof
            : EgipConstants.TrustImpactInvalidProof);
    }

    private void HandleTreaty(Envelope env)
    {
        var payload = (TreatyPayload)env.Payload;

        switch (payload.Action)
        {
            case TreatyAction.Propose:
                var treaty = new Treaty(payload.TreatyId, env.From, env.To, payload.Terms);
                _treaties.Put(treaty);
                break;

            case TreatyAction.Accept:
                _treaties.Apply(payload.TreatyId, t =>
                {
                    t.ApplyAction(TreatyAction.Accept);
                    _peers.UpdateTrust(env.From, EgipConstants.TrustImpactTreatyHonoured);
                });
                break;

            case TreatyAction.Suspend:
                _treaties.Apply(payload.TreatyId, t => t.ApplyAction(TreatyAction.Suspend));
                break;

            case TreatyAction.Terminate:
                _treaties.Apply(payload.TreatyId, t => t.ApplyAction(TreatyAction.Terminate));
                break;

            case TreatyAction.Modify:
                _treaties.Apply(payload.TreatyId, t =>
                {
                    t.ApplyAction(TreatyAction.Modify);
                    t.Terms = payload.Terms;
                });
                break;

            default:
                throw new EgipException($"Unknown treaty action: {payload.Action}");
        }
    }

    private async Task HandleAuthorityRequestAsync(Envelope env)
    {
        var payload = (AuthorityRequestPayload)env.Payload;
        if (OnAuthorityRequest is not null)
            await OnAuthorityRequest(env.From, payload);
    }

    private async Task HandleDiscoverAsync(Envelope env, CancellationToken ct)
    {
        var payload = (DiscoverPayload)env.Payload;
        if (OnDiscover is null) return;

        var results = await OnDiscover(env.From, payload.Query);

        var respId = new EnvelopeId(Guid.NewGuid().ToString());
        var resp = new Envelope
        {
            ProtocolVersion = EgipConstants.CurrentProtocolVersion,
            Id = respId,
            From = _identity.SystemUri,
            To = env.From,
            Type = MessageType.Discover,
            Payload = new DiscoverPayload(payload.Query, results),
            Timestamp = DateTimeOffset.UtcNow,
            Signature = new Signature(new byte[64]),
            InReplyTo = Option<EnvelopeId>.Some(env.Id),
        };

        var signed = Envelope.SignEnvelope(resp, _identity);
        await _transport.SendAsync(env.From, signed, ct);
    }
}
