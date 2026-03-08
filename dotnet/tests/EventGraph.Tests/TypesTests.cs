namespace EventGraph.Tests;

public class ScoreTests
{
    [Fact] public void ValidRange() { Assert.Equal(0.5, new Score(0.5).Value); }
    [Fact] public void BoundaryValues() { Assert.Equal(0.0, new Score(0.0).Value); Assert.Equal(1.0, new Score(1.0).Value); }
    [Fact] public void BelowRange() => Assert.Throws<OutOfRangeException>(() => new Score(-0.1));
    [Fact] public void AboveRange() => Assert.Throws<OutOfRangeException>(() => new Score(1.1));
    [Fact] public void NaN() => Assert.Throws<OutOfRangeException>(() => new Score(double.NaN));
}

public class WeightTests
{
    [Fact] public void ValidRange() { Assert.Equal(-1.0, new Weight(-1.0).Value); Assert.Equal(1.0, new Weight(1.0).Value); }
    [Fact] public void BelowRange() => Assert.Throws<OutOfRangeException>(() => new Weight(-1.1));
    [Fact] public void AboveRange() => Assert.Throws<OutOfRangeException>(() => new Weight(1.1));
}

public class ActivationTests
{
    [Fact] public void ValidRange() { Assert.Equal(0.0, new Activation(0.0).Value); Assert.Equal(1.0, new Activation(1.0).Value); }
    [Fact] public void BelowRange() => Assert.Throws<OutOfRangeException>(() => new Activation(-0.01));
}

public class LayerTests
{
    [Fact] public void ValidRange() { Assert.Equal(0, new Layer(0).Value); Assert.Equal(13, new Layer(13).Value); }
    [Fact] public void BelowRange() => Assert.Throws<OutOfRangeException>(() => new Layer(-1));
    [Fact] public void AboveRange() => Assert.Throws<OutOfRangeException>(() => new Layer(14));
}

public class CadenceTests
{
    [Fact] public void Valid() { Assert.Equal(1, new Cadence(1).Value); Assert.Equal(100, new Cadence(100).Value); }
    [Fact] public void Zero() => Assert.Throws<OutOfRangeException>(() => new Cadence(0));
}

public class EventIdTests
{
    [Fact] public void ValidUuidV7() => Assert.Equal("019462a0-0000-7000-8000-000000000001", new EventId("019462a0-0000-7000-8000-000000000001").Value);
    [Fact] public void NormalizesToLowercase() => Assert.Equal("019462a0-0000-7000-8000-000000000001", new EventId("019462A0-0000-7000-8000-000000000001").Value);
    [Fact] public void RejectsNonV7() => Assert.Throws<InvalidFormatException>(() => new EventId("019462a0-0000-4000-8000-000000000001"));
    [Fact] public void RejectsGarbage() => Assert.Throws<InvalidFormatException>(() => new EventId("not-a-uuid"));
}

public class HashTests
{
    [Fact] public void Valid64Hex() => Assert.Equal(new string('a', 64), new Hash(new string('a', 64)).Value);
    [Fact] public void ZeroHash() { var h = Hash.Zero(); Assert.True(h.IsZero); Assert.Equal(new string('0', 64), h.Value); }
    [Fact] public void RejectsShort() => Assert.Throws<InvalidFormatException>(() => new Hash("abc"));
    [Fact] public void RejectsEmpty() => Assert.Throws<InvalidFormatException>(() => new Hash(""));
}

public class EventTypeTests
{
    [Fact] public void Valid() => Assert.Equal("trust.updated", new EventType("trust.updated").Value);
    [Fact] public void RejectsUppercase() => Assert.Throws<InvalidFormatException>(() => new EventType("Trust.Updated"));
    [Fact] public void RejectsEmpty() => Assert.Throws<InvalidFormatException>(() => new EventType(""));
}

public class ActorIdTests
{
    [Fact] public void Valid() => Assert.Equal("actor_alice", new ActorId("actor_alice").Value);
    [Fact] public void RejectsEmpty() => Assert.Throws<EmptyRequiredException>(() => new ActorId(""));
}

public class SubscriptionPatternTests
{
    [Fact]
    public void WildcardMatchesAll()
    {
        var sp = new SubscriptionPattern("*");
        Assert.True(sp.Matches(new EventType("trust.updated")));
        Assert.True(sp.Matches(new EventType("system.bootstrapped")));
    }

    [Fact]
    public void PrefixMatch()
    {
        var sp = new SubscriptionPattern("trust.*");
        Assert.True(sp.Matches(new EventType("trust.updated")));
        Assert.False(sp.Matches(new EventType("system.bootstrapped")));
    }

    [Fact]
    public void ExactMatch()
    {
        var sp = new SubscriptionPattern("trust.updated");
        Assert.True(sp.Matches(new EventType("trust.updated")));
        Assert.False(sp.Matches(new EventType("trust.decayed")));
    }
}

public class OptionTests
{
    [Fact]
    public void Some()
    {
        var opt = Option<int>.Some(42);
        Assert.True(opt.IsSome);
        Assert.False(opt.IsNone);
        Assert.Equal(42, opt.Unwrap());
    }

    [Fact]
    public void None()
    {
        var opt = Option<int>.None();
        Assert.True(opt.IsNone);
        Assert.Throws<InvalidOperationException>(() => opt.Unwrap());
    }

    [Fact]
    public void UnwrapOr()
    {
        Assert.Equal(42, Option<int>.Some(42).UnwrapOr(0));
        Assert.Equal(0, Option<int>.None().UnwrapOr(0));
    }
}

public class NonEmptyTests
{
    [Fact] public void Valid() { var ne = NonEmpty<int>.Of(new[] { 1, 2, 3 }); Assert.Equal(3, ne.Count); Assert.Equal(1, ne[0]); }
    [Fact] public void RejectsEmpty() => Assert.Throws<ArgumentException>(() => NonEmpty<int>.Of(Array.Empty<int>()));
    [Fact] public void Iterable() { var ne = NonEmpty<int>.Of(new[] { 1, 2 }); Assert.Equal(new[] { 1, 2 }, ne.ToArray()); }
}

public class PublicKeyTests
{
    [Fact] public void Valid32Bytes() { var pk = new PublicKey(new byte[32]); Assert.Equal(32, pk.Bytes.Length); }
    [Fact] public void RejectsWrongLength() => Assert.Throws<InvalidFormatException>(() => new PublicKey(new byte[31]));
}

public class SignatureTests
{
    [Fact] public void Valid64Bytes() { var sig = new Signature(new byte[64]); Assert.Equal(64, sig.Bytes.Length); }
    [Fact] public void RejectsWrongLength() => Assert.Throws<InvalidFormatException>(() => new Signature(new byte[63]));
}
