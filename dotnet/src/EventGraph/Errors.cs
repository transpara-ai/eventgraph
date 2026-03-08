namespace EventGraph;

/// <summary>Base class for all EventGraph errors.</summary>
public class EventGraphException : Exception
{
    public EventGraphException(string message) : base(message) { }
    public EventGraphException(string message, Exception inner) : base(message, inner) { }
}

/// <summary>A validated value is out of its allowed range.</summary>
public class OutOfRangeException : EventGraphException
{
    public string TypeName { get; }
    public double Value { get; }
    public double Min { get; }
    public double Max { get; }

    public OutOfRangeException(string typeName, double value, double min, double max)
        : base($"{typeName} value {value} is out of range [{min}, {max}]")
    {
        TypeName = typeName;
        Value = value;
        Min = min;
        Max = max;
    }
}

/// <summary>A required value was empty.</summary>
public class EmptyRequiredException : EventGraphException
{
    public string TypeName { get; }
    public EmptyRequiredException(string typeName) : base($"{typeName} cannot be empty") => TypeName = typeName;
}

/// <summary>A value did not match the required format.</summary>
public class InvalidFormatException : EventGraphException
{
    public string TypeName { get; }
    public string ProvidedValue { get; }
    public string ExpectedFormat { get; }

    public InvalidFormatException(string typeName, string value, string expected)
        : base($"{typeName} value \"{value}\" does not match expected format: {expected}")
    {
        TypeName = typeName;
        ProvidedValue = value;
        ExpectedFormat = expected;
    }
}

/// <summary>An invalid state machine transition was attempted.</summary>
public class InvalidTransitionException : EventGraphException
{
    public string From { get; }
    public string To { get; }

    public InvalidTransitionException(string from, string to)
        : base($"Invalid transition from {from} to {to}")
    {
        From = from;
        To = to;
    }
}

/// <summary>An event was not found in the store.</summary>
public class EventNotFoundException : EventGraphException
{
    public string EventId { get; }
    public EventNotFoundException(string eventId) : base($"Event {eventId} not found") => EventId = eventId;
}

/// <summary>The hash chain integrity was violated.</summary>
public class ChainIntegrityException : EventGraphException
{
    public int Position { get; }
    public ChainIntegrityException(int position, string detail)
        : base($"Chain integrity violation at position {position}: {detail}") => Position = position;
}
