using System.Text.RegularExpressions;

namespace EventGraph;

// ── Option<T> ───────────────────────────────────────────────────────────

/// <summary>Explicit optionality — Some(value) or None. No null ambiguity.</summary>
public readonly struct Option<T> where T : notnull
{
    private readonly T? _value;
    private readonly bool _present;

    private Option(T? value, bool present) { _value = value; _present = present; }

    public static Option<T> Some(T value) => new(value, true);
    public static Option<T> None() => new(default, false);

    public bool IsSome => _present;
    public bool IsNone => !_present;

    public T Unwrap() => _present ? _value! : throw new InvalidOperationException("Unwrap called on None Option");
    public T UnwrapOr(T defaultValue) => _present ? _value! : defaultValue;
}

// ── NonEmpty<T> ─────────────────────────────────────────────────────────

/// <summary>A collection with at least one element.</summary>
public sealed class NonEmpty<T> : IReadOnlyList<T>
{
    private readonly IReadOnlyList<T> _items;

    private NonEmpty(IReadOnlyList<T> items) => _items = items;

    public static NonEmpty<T> Of(IReadOnlyList<T> items)
    {
        if (items.Count == 0) throw new ArgumentException("NonEmpty requires at least one element");
        return new NonEmpty<T>(items.ToArray());
    }

    public T this[int index] => _items[index];
    public int Count => _items.Count;
    public IEnumerator<T> GetEnumerator() => _items.GetEnumerator();
    System.Collections.IEnumerator System.Collections.IEnumerable.GetEnumerator() => GetEnumerator();
}

// ── Constrained numerics ────────────────────────────────────────────────

/// <summary>Float constrained to [0.0, 1.0].</summary>
public readonly record struct Score
{
    public double Value { get; }
    public Score(double value)
    {
        if (double.IsNaN(value) || value < 0.0 || value > 1.0)
            throw new OutOfRangeException("Score", value, 0.0, 1.0);
        Value = value == 0.0 ? 0.0 : value; // normalize -0.0
    }
}

/// <summary>Float constrained to [-1.0, 1.0].</summary>
public readonly record struct Weight
{
    public double Value { get; }
    public Weight(double value)
    {
        if (double.IsNaN(value) || value < -1.0 || value > 1.0)
            throw new OutOfRangeException("Weight", value, -1.0, 1.0);
        Value = value == 0.0 ? 0.0 : value;
    }
}

/// <summary>Float constrained to [0.0, 1.0].</summary>
public readonly record struct Activation
{
    public double Value { get; }
    public Activation(double value)
    {
        if (double.IsNaN(value) || value < 0.0 || value > 1.0)
            throw new OutOfRangeException("Activation", value, 0.0, 1.0);
        Value = value == 0.0 ? 0.0 : value;
    }
}

/// <summary>Int constrained to [0, 13].</summary>
public readonly record struct Layer : IComparable<Layer>
{
    public int Value { get; }
    public Layer(int value)
    {
        if (value < 0 || value > 13)
            throw new OutOfRangeException("Layer", value, 0, 13);
        Value = value;
    }
    public int CompareTo(Layer other) => Value.CompareTo(other.Value);
}

/// <summary>Int constrained to [1, ∞).</summary>
public readonly record struct Cadence
{
    public int Value { get; }
    public Cadence(int value)
    {
        if (value < 1)
            throw new OutOfRangeException("Cadence", value, 1, double.PositiveInfinity);
        Value = value;
    }
}

// ── Typed IDs ───────────────────────────────────────────────────────────

internal static partial class Patterns
{
    [GeneratedRegex(@"^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")]
    internal static partial Regex UuidV7();

    [GeneratedRegex(@"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")]
    internal static partial Regex Uuid();

    [GeneratedRegex(@"^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$")]
    internal static partial Regex EventTypePattern();

    [GeneratedRegex(@"^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)*$")]
    internal static partial Regex DomainScopePattern();

    [GeneratedRegex(@"^(\*|[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*(\.\*)?)$")]
    internal static partial Regex SubscriptionPatternRe();
}

/// <summary>UUID v7 event identifier.</summary>
public readonly record struct EventId
{
    public string Value { get; }
    public EventId(string value)
    {
        value = value.ToLowerInvariant();
        if (!Patterns.UuidV7().IsMatch(value))
            throw new InvalidFormatException("EventId", value, "UUID v7");
        Value = value;
    }
    public override string ToString() => Value;
}

/// <summary>UUID v7 edge identifier.</summary>
public readonly record struct EdgeId
{
    public string Value { get; }
    public EdgeId(string value)
    {
        value = value.ToLowerInvariant();
        if (!Patterns.UuidV7().IsMatch(value))
            throw new InvalidFormatException("EdgeId", value, "UUID v7");
        Value = value;
    }
    public override string ToString() => Value;
}

/// <summary>SHA-256 hex string (64 characters).</summary>
public readonly record struct Hash
{
    public string Value { get; }
    public Hash(string value)
    {
        value = value.ToLowerInvariant();
        if (value.Length != 64 || !value.All(c => "0123456789abcdef".Contains(c)))
            throw new InvalidFormatException("Hash", value, "64 hex characters (SHA-256)");
        Value = value;
    }
    public static Hash Zero() => new(new string('0', 64));
    public bool IsZero => Value == new string('0', 64);
    public override string ToString() => Value;
}

/// <summary>Dot-separated lowercase event type.</summary>
public readonly record struct EventType
{
    public string Value { get; }
    public EventType(string value)
    {
        if (!Patterns.EventTypePattern().IsMatch(value))
            throw new InvalidFormatException("EventType", value, "dot-separated lowercase segments");
        Value = value;
    }
    public override string ToString() => Value;
}

/// <summary>Non-empty string actor identifier.</summary>
public readonly record struct ActorId
{
    public string Value { get; }
    public ActorId(string value)
    {
        if (string.IsNullOrEmpty(value))
            throw new EmptyRequiredException("ActorId");
        Value = value;
    }
    public override string ToString() => Value;
}

/// <summary>Non-empty string conversation identifier.</summary>
public readonly record struct ConversationId
{
    public string Value { get; }
    public ConversationId(string value)
    {
        if (string.IsNullOrEmpty(value))
            throw new EmptyRequiredException("ConversationId");
        Value = value;
    }
    public override string ToString() => Value;
}

/// <summary>Non-empty system URI.</summary>
public readonly record struct SystemUri
{
    public string Value { get; }
    public SystemUri(string value)
    {
        if (string.IsNullOrEmpty(value))
            throw new EmptyRequiredException("SystemUri");
        Value = value;
    }
    public override string ToString() => Value;
}

/// <summary>Non-empty primitive identifier.</summary>
public readonly record struct PrimitiveId
{
    public string Value { get; }
    public PrimitiveId(string value)
    {
        if (string.IsNullOrEmpty(value))
            throw new EmptyRequiredException("PrimitiveId");
        Value = value;
    }
    public override string ToString() => Value;
}

/// <summary>Trust/authority domain scope.</summary>
public readonly record struct DomainScope
{
    public string Value { get; }
    public DomainScope(string value)
    {
        if (!Patterns.DomainScopePattern().IsMatch(value))
            throw new InvalidFormatException("DomainScope", value, "lowercase dot/underscore-separated namespace");
        Value = value;
    }
    public override string ToString() => Value;
}

/// <summary>Glob pattern for event type matching.</summary>
public readonly record struct SubscriptionPattern
{
    public string Value { get; }
    public SubscriptionPattern(string value)
    {
        if (!Patterns.SubscriptionPatternRe().IsMatch(value))
            throw new InvalidFormatException("SubscriptionPattern", value, "dot-separated with optional trailing .* or bare *");
        Value = value;
    }

    public bool Matches(EventType et)
    {
        if (Value == "*") return true;
        if (Value.EndsWith(".*"))
        {
            var prefix = Value[..^2];
            return et.Value == prefix || et.Value.StartsWith(prefix + ".");
        }
        return Value == et.Value;
    }
}

/// <summary>UUID envelope identifier.</summary>
public readonly record struct EnvelopeId
{
    public string Value { get; }
    public EnvelopeId(string value)
    {
        value = value.ToLowerInvariant();
        if (!Patterns.Uuid().IsMatch(value))
            throw new InvalidFormatException("EnvelopeId", value, "UUID");
        Value = value;
    }
    public override string ToString() => Value;
}

/// <summary>UUID treaty identifier.</summary>
public readonly record struct TreatyId
{
    public string Value { get; }
    public TreatyId(string value)
    {
        value = value.ToLowerInvariant();
        if (!Patterns.Uuid().IsMatch(value))
            throw new InvalidFormatException("TreatyId", value, "UUID");
        Value = value;
    }
    public override string ToString() => Value;
}

/// <summary>Ed25519 public key (32 bytes).</summary>
public readonly record struct PublicKey
{
    private readonly byte[] _value;
    public ReadOnlySpan<byte> Bytes => _value;

    public PublicKey(byte[] value)
    {
        if (value.Length != 32)
            throw new InvalidFormatException("PublicKey", Convert.ToHexString(value), "32 bytes (Ed25519 public key)");
        _value = (byte[])value.Clone();
    }
    public override string ToString() => Convert.ToHexString(_value).ToLowerInvariant();
}

/// <summary>Ed25519 signature (64 bytes).</summary>
public readonly record struct Signature
{
    private readonly byte[] _value;
    public ReadOnlySpan<byte> Bytes => _value;

    public Signature(byte[] value)
    {
        if (value.Length != 64)
            throw new InvalidFormatException("Signature", Convert.ToHexString(value), "64 bytes (Ed25519 signature)");
        _value = (byte[])value.Clone();
    }
    public override string ToString() => Convert.ToHexString(_value).ToLowerInvariant();
}
