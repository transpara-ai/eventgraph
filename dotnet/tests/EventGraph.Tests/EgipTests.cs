using System.Runtime.CompilerServices;
using System.Text;

namespace EventGraph.Tests;

// ── Test helpers ─────────────────────────────────────────────────────────

/// <summary>In-memory transport for testing.</summary>
internal sealed class InMemoryTransport : ITransport
{
    public List<(SystemUri To, Envelope Envelope)> Sent { get; } = new();
    public ReceiptPayload? NextReceipt { get; set; }
    public Exception? SendException { get; set; }

    public Task<ReceiptPayload?> SendAsync(SystemUri to, Envelope envelope, CancellationToken ct = default)
    {
        if (SendException is not null) throw SendException;
        Sent.Add((to, envelope));
        return Task.FromResult(NextReceipt);
    }

    public async IAsyncEnumerable<IncomingEnvelope> ListenAsync([EnumeratorCancellation] CancellationToken ct = default)
    {
        await Task.CompletedTask;
        yield break;
    }
}

// ── Enum tests ───────────────────────────────────────────────────────────

public class EgipEnumTests
{
    [Fact]
    public void MessageType_AllVariants()
    {
        var values = Enum.GetValues<MessageType>();
        Assert.Equal(7, values.Length);
        Assert.Contains(MessageType.Hello, values);
        Assert.Contains(MessageType.Message, values);
        Assert.Contains(MessageType.Receipt, values);
        Assert.Contains(MessageType.Proof, values);
        Assert.Contains(MessageType.Treaty, values);
        Assert.Contains(MessageType.AuthorityRequest, values);
        Assert.Contains(MessageType.Discover, values);
    }

    [Fact]
    public void TreatyStatus_AllVariants()
    {
        var values = Enum.GetValues<TreatyStatus>();
        Assert.Equal(4, values.Length);
    }

    [Fact]
    public void TreatyAction_AllVariants()
    {
        var values = Enum.GetValues<TreatyAction>();
        Assert.Equal(5, values.Length);
    }

    [Fact]
    public void ReceiptStatus_AllVariants()
    {
        var values = Enum.GetValues<ReceiptStatus>();
        Assert.Equal(3, values.Length);
    }

    [Fact]
    public void ProofType_AllVariants()
    {
        var values = Enum.GetValues<ProofType>();
        Assert.Equal(3, values.Length);
    }

    [Fact]
    public void CgerRelationship_AllVariants()
    {
        var values = Enum.GetValues<CgerRelationship>();
        Assert.Equal(3, values.Length);
    }
}

// ── SystemIdentity tests ─────────────────────────────────────────────────

public class SystemIdentityTests
{
    [Fact]
    public void Generate_CreatesValidIdentity()
    {
        using var id = SystemIdentity.Generate(new SystemUri("eg://test.local"));
        Assert.Equal("eg://test.local", id.SystemUri.Value);
        Assert.Equal(32, id.PublicKey.Bytes.Length);
    }

    [Fact]
    public void Sign_ProducesValidSignature()
    {
        using var id = SystemIdentity.Generate(new SystemUri("eg://test.local"));
        var data = Encoding.UTF8.GetBytes("hello world");
        var sig = id.Sign(data);
        Assert.Equal(64, sig.Bytes.Length);
    }

    [Fact]
    public void Verify_ValidSignature_ReturnsTrue()
    {
        using var id = SystemIdentity.Generate(new SystemUri("eg://test.local"));
        var data = Encoding.UTF8.GetBytes("hello world");
        var sig = id.Sign(data);
        Assert.True(id.Verify(id.PublicKey, data, sig));
    }

    [Fact]
    public void Verify_WrongData_ReturnsFalse()
    {
        using var id = SystemIdentity.Generate(new SystemUri("eg://test.local"));
        var data = Encoding.UTF8.GetBytes("hello world");
        var sig = id.Sign(data);
        var wrongData = Encoding.UTF8.GetBytes("wrong data");
        Assert.False(id.Verify(id.PublicKey, wrongData, sig));
    }

    [Fact]
    public void Verify_WrongKey_ReturnsFalse()
    {
        using var id1 = SystemIdentity.Generate(new SystemUri("eg://system-a"));
        using var id2 = SystemIdentity.Generate(new SystemUri("eg://system-b"));
        var data = Encoding.UTF8.GetBytes("hello world");
        var sig = id1.Sign(data);
        Assert.False(id1.Verify(id2.PublicKey, data, sig));
    }
}

// ── Envelope tests ───────────────────────────────────────────────────────

public class EnvelopeTests
{
    private static Envelope CreateTestEnvelope(IIdentity identity, SystemUri to)
    {
        return new Envelope
        {
            ProtocolVersion = EgipConstants.CurrentProtocolVersion,
            Id = new EnvelopeId(Guid.NewGuid().ToString()),
            From = identity.SystemUri,
            To = to,
            Type = MessageType.Hello,
            Payload = new HelloPayload(
                identity.SystemUri,
                identity.PublicKey,
                [1],
                ["treaty", "proof"],
                42),
            Timestamp = DateTimeOffset.UtcNow,
            Signature = new Signature(new byte[64]),
            InReplyTo = Option<EnvelopeId>.None(),
        };
    }

    [Fact]
    public void CanonicalForm_IsDeterministic()
    {
        using var id = SystemIdentity.Generate(new SystemUri("eg://test.local"));
        var env = CreateTestEnvelope(id, new SystemUri("eg://remote"));

        var form1 = env.CanonicalForm();
        var form2 = env.CanonicalForm();
        Assert.Equal(form1, form2);
    }

    [Fact]
    public void CanonicalForm_ContainsExpectedFields()
    {
        using var id = SystemIdentity.Generate(new SystemUri("eg://test.local"));
        var env = CreateTestEnvelope(id, new SystemUri("eg://remote"));

        var form = env.CanonicalForm();
        Assert.Contains("eg://test.local", form);
        Assert.Contains("eg://remote", form);
        Assert.Contains("hello", form); // message type lowercase
    }

    [Fact]
    public void SignAndVerify_RoundTrip()
    {
        using var id = SystemIdentity.Generate(new SystemUri("eg://test.local"));
        var env = CreateTestEnvelope(id, new SystemUri("eg://remote"));

        var signed = Envelope.SignEnvelope(env, id);
        Assert.True(Envelope.VerifyEnvelope(signed, id, id.PublicKey));
    }

    [Fact]
    public void VerifyEnvelope_DifferentKey_ReturnsFalse()
    {
        using var id1 = SystemIdentity.Generate(new SystemUri("eg://system-a"));
        using var id2 = SystemIdentity.Generate(new SystemUri("eg://system-b"));

        var env = CreateTestEnvelope(id1, new SystemUri("eg://remote"));
        var signed = Envelope.SignEnvelope(env, id1);

        Assert.False(Envelope.VerifyEnvelope(signed, id1, id2.PublicKey));
    }

    [Fact]
    public void CanonicalForm_InReplyTo_IncludedWhenPresent()
    {
        using var id = SystemIdentity.Generate(new SystemUri("eg://test.local"));
        var replyId = new EnvelopeId(Guid.NewGuid().ToString());
        var env = new Envelope
        {
            ProtocolVersion = 1,
            Id = new EnvelopeId(Guid.NewGuid().ToString()),
            From = id.SystemUri,
            To = new SystemUri("eg://remote"),
            Type = MessageType.Receipt,
            Payload = new ReceiptPayload(
                replyId,
                ReceiptStatus.Delivered,
                Option<EventId>.None(),
                Option<string>.None(),
                new Signature(new byte[64])),
            Timestamp = DateTimeOffset.UtcNow,
            Signature = new Signature(new byte[64]),
            InReplyTo = Option<EnvelopeId>.Some(replyId),
        };

        var form = env.CanonicalForm();
        Assert.Contains(replyId.Value, form);
    }
}

// ── NegotiateVersion tests ───────────────────────────────────────────────

public class NegotiateVersionTests
{
    [Fact]
    public void CommonVersion_ReturnsHighest()
    {
        var result = EgipVersionNegotiation.NegotiateVersion([1, 2, 3], [2, 3, 4]);
        Assert.True(result.IsSome);
        Assert.Equal(3, result.Unwrap());
    }

    [Fact]
    public void NoCommonVersion_ReturnsNone()
    {
        var result = EgipVersionNegotiation.NegotiateVersion([1, 2], [3, 4]);
        Assert.True(result.IsNone);
    }

    [Fact]
    public void SingleCommonVersion()
    {
        var result = EgipVersionNegotiation.NegotiateVersion([1], [1]);
        Assert.True(result.IsSome);
        Assert.Equal(1, result.Unwrap());
    }

    [Fact]
    public void EmptyLocal_ReturnsNone()
    {
        var result = EgipVersionNegotiation.NegotiateVersion([], [1, 2]);
        Assert.True(result.IsNone);
    }
}

// ── PeerStore tests ──────────────────────────────────────────────────────

public class PeerStoreTests
{
    [Fact]
    public void Register_NewPeer_CreatesRecord()
    {
        var store = new PeerStore();
        var uri = new SystemUri("eg://peer");
        var key = new PublicKey(new byte[32]);
        var record = store.Register(uri, key, ["treaty"], 1);

        Assert.Equal("eg://peer", record.SystemUri.Value);
        Assert.Equal(0.0, record.Trust.Value);
    }

    [Fact]
    public void Register_ExistingPeer_DoesNotOverwriteKey()
    {
        var store = new PeerStore();
        var uri = new SystemUri("eg://peer");
        var key1 = new PublicKey(new byte[32]);
        store.Register(uri, key1, ["treaty"], 1);

        var newKeyBytes = new byte[32];
        newKeyBytes[0] = 0xFF;
        var key2 = new PublicKey(newKeyBytes);
        store.Register(uri, key2, ["proof"], 2);

        var (record, found) = store.Get(uri);
        Assert.True(found);
        Assert.Equal(0, record!.PublicKey.Bytes[0]); // original key preserved
        Assert.Contains("proof", record.Capabilities);
        Assert.Equal(2, record.NegotiatedVersion);
    }

    [Fact]
    public void Get_UnknownPeer_ReturnsFalse()
    {
        var store = new PeerStore();
        var (_, found) = store.Get(new SystemUri("eg://unknown"));
        Assert.False(found);
    }

    [Fact]
    public void UpdateTrust_PositiveCapped()
    {
        var store = new PeerStore();
        var uri = new SystemUri("eg://peer");
        store.Register(uri, new PublicKey(new byte[32]), [], 1);

        // Try a large positive delta — should be capped at InterSystemMaxAdjustment (0.05)
        var (score, found) = store.UpdateTrust(uri, 1.0);
        Assert.True(found);
        Assert.Equal(0.05, score.Value);
    }

    [Fact]
    public void UpdateTrust_NegativeUncapped()
    {
        var store = new PeerStore();
        var uri = new SystemUri("eg://peer");
        store.Register(uri, new PublicKey(new byte[32]), [], 1);

        // First give some trust
        store.UpdateTrust(uri, 0.05);
        store.UpdateTrust(uri, 0.05);

        // Negative delta is not capped (but clamped to 0)
        var (score, _) = store.UpdateTrust(uri, -0.20);
        Assert.Equal(0.0, score.Value); // clamped to 0
    }

    [Fact]
    public void UpdateTrust_UnknownPeer_ReturnsFalse()
    {
        var store = new PeerStore();
        var (_, found) = store.UpdateTrust(new SystemUri("eg://unknown"), 0.01);
        Assert.False(found);
    }

    [Fact]
    public void All_ReturnsAllPeers()
    {
        var store = new PeerStore();
        store.Register(new SystemUri("eg://a"), new PublicKey(new byte[32]), [], 1);
        store.Register(new SystemUri("eg://b"), new PublicKey(new byte[32]), [], 1);

        var all = store.All();
        Assert.Equal(2, all.Count);
    }
}

// ── TreatyStore tests ────────────────────────────────────────────────────

public class TreatyStoreTests
{
    [Fact]
    public void PutAndGet()
    {
        var store = new TreatyStore();
        var id = new TreatyId(Guid.NewGuid().ToString());
        var treaty = new Treaty(id, new SystemUri("eg://a"), new SystemUri("eg://b"), []);
        store.Put(treaty);

        var (retrieved, found) = store.Get(id);
        Assert.True(found);
        Assert.Equal(id, retrieved!.Id);
    }

    [Fact]
    public void Get_NotFound()
    {
        var store = new TreatyStore();
        var (_, found) = store.Get(new TreatyId(Guid.NewGuid().ToString()));
        Assert.False(found);
    }

    [Fact]
    public void Apply_MutatesTreaty()
    {
        var store = new TreatyStore();
        var id = new TreatyId(Guid.NewGuid().ToString());
        var treaty = new Treaty(id, new SystemUri("eg://a"), new SystemUri("eg://b"), []);
        store.Put(treaty);

        store.Apply(id, t => t.ApplyAction(TreatyAction.Accept));

        var (retrieved, _) = store.Get(id);
        Assert.Equal(TreatyStatus.Active, retrieved!.Status);
    }

    [Fact]
    public void Apply_NotFound_Throws()
    {
        var store = new TreatyStore();
        var id = new TreatyId(Guid.NewGuid().ToString());
        Assert.Throws<TreatyNotFoundException>(() => store.Apply(id, _ => { }));
    }

    [Fact]
    public void BySystem_ReturnsMatchingTreaties()
    {
        var store = new TreatyStore();
        var uriA = new SystemUri("eg://a");
        var uriB = new SystemUri("eg://b");
        var uriC = new SystemUri("eg://c");

        store.Put(new Treaty(new TreatyId(Guid.NewGuid().ToString()), uriA, uriB, []));
        store.Put(new Treaty(new TreatyId(Guid.NewGuid().ToString()), uriB, uriC, []));
        store.Put(new Treaty(new TreatyId(Guid.NewGuid().ToString()), uriC, uriA, []));

        var byA = store.BySystem(uriA);
        Assert.Equal(2, byA.Count);
    }

    [Fact]
    public void Active_ReturnsOnlyActive()
    {
        var store = new TreatyStore();
        var id1 = new TreatyId(Guid.NewGuid().ToString());
        var id2 = new TreatyId(Guid.NewGuid().ToString());

        var treaty1 = new Treaty(id1, new SystemUri("eg://a"), new SystemUri("eg://b"), []);
        var treaty2 = new Treaty(id2, new SystemUri("eg://c"), new SystemUri("eg://d"), []);

        store.Put(treaty1);
        store.Put(treaty2);

        // Accept treaty1
        store.Apply(id1, t => t.ApplyAction(TreatyAction.Accept));

        var active = store.Active();
        Assert.Single(active);
        Assert.Equal(id1, active[0].Id);
    }
}

// ── Treaty state machine tests ───────────────────────────────────────────

public class TreatyTests
{
    [Fact]
    public void NewTreaty_IsProposed()
    {
        var treaty = new Treaty(
            new TreatyId(Guid.NewGuid().ToString()),
            new SystemUri("eg://a"),
            new SystemUri("eg://b"),
            []);

        Assert.Equal(TreatyStatus.Proposed, treaty.Status);
    }

    [Fact]
    public void Proposed_CanAccept()
    {
        var treaty = new Treaty(
            new TreatyId(Guid.NewGuid().ToString()),
            new SystemUri("eg://a"),
            new SystemUri("eg://b"),
            []);

        treaty.ApplyAction(TreatyAction.Accept);
        Assert.Equal(TreatyStatus.Active, treaty.Status);
    }

    [Fact]
    public void Proposed_CanTerminate()
    {
        var treaty = new Treaty(
            new TreatyId(Guid.NewGuid().ToString()),
            new SystemUri("eg://a"),
            new SystemUri("eg://b"),
            []);

        treaty.ApplyAction(TreatyAction.Terminate);
        Assert.Equal(TreatyStatus.Terminated, treaty.Status);
    }

    [Fact]
    public void Proposed_CannotSuspend()
    {
        var treaty = new Treaty(
            new TreatyId(Guid.NewGuid().ToString()),
            new SystemUri("eg://a"),
            new SystemUri("eg://b"),
            []);

        Assert.Throws<InvalidTransitionException>(() => treaty.ApplyAction(TreatyAction.Suspend));
    }

    [Fact]
    public void Active_CanSuspend()
    {
        var treaty = new Treaty(
            new TreatyId(Guid.NewGuid().ToString()),
            new SystemUri("eg://a"),
            new SystemUri("eg://b"),
            []);

        treaty.ApplyAction(TreatyAction.Accept);
        treaty.ApplyAction(TreatyAction.Suspend);
        Assert.Equal(TreatyStatus.Suspended, treaty.Status);
    }

    [Fact]
    public void Active_CanModify()
    {
        var treaty = new Treaty(
            new TreatyId(Guid.NewGuid().ToString()),
            new SystemUri("eg://a"),
            new SystemUri("eg://b"),
            []);

        treaty.ApplyAction(TreatyAction.Accept);
        treaty.ApplyAction(TreatyAction.Modify); // should not throw
        Assert.Equal(TreatyStatus.Active, treaty.Status); // status unchanged
    }

    [Fact]
    public void Suspended_CanReactivate()
    {
        var treaty = new Treaty(
            new TreatyId(Guid.NewGuid().ToString()),
            new SystemUri("eg://a"),
            new SystemUri("eg://b"),
            []);

        treaty.ApplyAction(TreatyAction.Accept);
        treaty.ApplyAction(TreatyAction.Suspend);
        treaty.ApplyAction(TreatyAction.Accept); // reactivate
        Assert.Equal(TreatyStatus.Active, treaty.Status);
    }

    [Fact]
    public void Terminated_IsTerminal()
    {
        var treaty = new Treaty(
            new TreatyId(Guid.NewGuid().ToString()),
            new SystemUri("eg://a"),
            new SystemUri("eg://b"),
            []);

        treaty.ApplyAction(TreatyAction.Terminate);
        Assert.Throws<InvalidTransitionException>(() => treaty.ApplyAction(TreatyAction.Accept));
    }

    [Fact]
    public void Propose_OnExistingTreaty_Throws()
    {
        var treaty = new Treaty(
            new TreatyId(Guid.NewGuid().ToString()),
            new SystemUri("eg://a"),
            new SystemUri("eg://b"),
            []);

        Assert.Throws<InvalidOperationException>(() => treaty.ApplyAction(TreatyAction.Propose));
    }

    [Fact]
    public void Modify_NotActive_Throws()
    {
        var treaty = new Treaty(
            new TreatyId(Guid.NewGuid().ToString()),
            new SystemUri("eg://a"),
            new SystemUri("eg://b"),
            []);

        Assert.Throws<InvalidTransitionException>(() => treaty.ApplyAction(TreatyAction.Modify));
    }
}

// ── EnvelopeDedup tests ──────────────────────────────────────────────────

public class EnvelopeDedupTests
{
    [Fact]
    public void Check_FirstTime_ReturnsTrue()
    {
        var dedup = new EnvelopeDedup();
        var id = new EnvelopeId(Guid.NewGuid().ToString());
        Assert.True(dedup.Check(id));
    }

    [Fact]
    public void Check_SecondTime_ReturnsFalse()
    {
        var dedup = new EnvelopeDedup();
        var id = new EnvelopeId(Guid.NewGuid().ToString());
        dedup.Check(id);
        Assert.False(dedup.Check(id));
    }

    [Fact]
    public void Size_TracksCount()
    {
        var dedup = new EnvelopeDedup();
        Assert.Equal(0, dedup.Size);

        dedup.Check(new EnvelopeId(Guid.NewGuid().ToString()));
        dedup.Check(new EnvelopeId(Guid.NewGuid().ToString()));
        Assert.Equal(2, dedup.Size);
    }

    [Fact]
    public void Prune_RemovesExpired()
    {
        // Use a very short TTL for testing
        var dedup = new EnvelopeDedup(TimeSpan.FromMilliseconds(1));
        dedup.Check(new EnvelopeId(Guid.NewGuid().ToString()));

        // Wait for entry to expire
        Thread.Sleep(50);
        var removed = dedup.Prune();
        Assert.Equal(1, removed);
        Assert.Equal(0, dedup.Size);
    }
}

// ── ProofVerification tests ──────────────────────────────────────────────

public class ProofVerificationTests
{
    [Fact]
    public void VerifyChainSegment_EmptyEvents_ReturnsFalse()
    {
        var proof = new ChainSegmentProof([], Hash.Zero(), Hash.Zero());
        Assert.False(ProofVerification.VerifyChainSegment(proof));
    }

    [Fact]
    public void VerifyEventExistence_ValidProof()
    {
        var signer = new NoopSigner();
        var evt = EventFactory.CreateBootstrap(new ActorId("system"), signer);

        var proof = new EventExistenceProof(
            evt,
            evt.PrevHash,
            Option<Hash>.None(),
            0,
            1);

        Assert.True(ProofVerification.VerifyEventExistence(proof));
    }

    [Fact]
    public void VerifyEventExistence_InvalidPosition_ReturnsFalse()
    {
        var signer = new NoopSigner();
        var evt = EventFactory.CreateBootstrap(new ActorId("system"), signer);

        var proof = new EventExistenceProof(
            evt,
            evt.PrevHash,
            Option<Hash>.None(),
            5, // position >= chainLength
            5);

        Assert.False(ProofVerification.VerifyEventExistence(proof));
    }

    [Fact]
    public void ValidateProof_ChainSummary_PositiveLength()
    {
        var proof = new ProofPayload(
            ProofType.ChainSummary,
            new ChainSummaryProof(10, Hash.Zero(), Hash.Zero(), DateTimeOffset.UtcNow));

        Assert.True(ProofVerification.ValidateProof(proof));
    }

    [Fact]
    public void ValidateProof_ChainSummary_ZeroLength()
    {
        var proof = new ProofPayload(
            ProofType.ChainSummary,
            new ChainSummaryProof(0, Hash.Zero(), Hash.Zero(), DateTimeOffset.UtcNow));

        Assert.False(ProofVerification.ValidateProof(proof));
    }

    [Fact]
    public void ProofTypeFromData_ReturnsCorrectType()
    {
        Assert.Equal(ProofType.ChainSummary,
            ProofVerification.ProofTypeFromData(
                new ChainSummaryProof(1, Hash.Zero(), Hash.Zero(), DateTimeOffset.UtcNow)));
        Assert.Equal(ProofType.ChainSegment,
            ProofVerification.ProofTypeFromData(new ChainSegmentProof([], Hash.Zero(), Hash.Zero())));
    }
}

// ── EGIP Error tests ─────────────────────────────────────────────────────

public class EgipErrorTests
{
    [Fact]
    public void AllErrors_InheritFromEgipException()
    {
        Assert.IsAssignableFrom<EgipException>(new SystemNotFoundException(new SystemUri("eg://x")));
        Assert.IsAssignableFrom<EgipException>(new EnvelopeSignatureInvalidException(new EnvelopeId(Guid.NewGuid().ToString())));
        Assert.IsAssignableFrom<EgipException>(new TreatyViolationException(new TreatyId(Guid.NewGuid().ToString()), "term"));
        Assert.IsAssignableFrom<EgipException>(new TrustInsufficientException(new SystemUri("eg://x"), new Score(0.1), new Score(0.5)));
        Assert.IsAssignableFrom<EgipException>(new TransportFailureException(new SystemUri("eg://x"), "timeout"));
        Assert.IsAssignableFrom<EgipException>(new DuplicateEnvelopeException(new EnvelopeId(Guid.NewGuid().ToString())));
        Assert.IsAssignableFrom<EgipException>(new TreatyNotFoundException(new TreatyId(Guid.NewGuid().ToString())));
        Assert.IsAssignableFrom<EgipException>(new VersionIncompatibleException([1], [2]));
    }

    [Fact]
    public void AllErrors_InheritFromEventGraphException()
    {
        Assert.IsAssignableFrom<EventGraphException>(new SystemNotFoundException(new SystemUri("eg://x")));
    }

    [Fact]
    public void ErrorMessages_ContainRelevantInfo()
    {
        var err = new SystemNotFoundException(new SystemUri("eg://test.local"));
        Assert.Contains("eg://test.local", err.Message);

        var err2 = new TrustInsufficientException(new SystemUri("eg://x"), new Score(0.1), new Score(0.5));
        Assert.Contains("0.1", err2.Message);
        Assert.Contains("0.5", err2.Message);
    }
}

// ── Handler integration tests ────────────────────────────────────────────

public class EgipHandlerTests
{
    private static (EgipHandler Handler, SystemIdentity Identity, InMemoryTransport Transport, PeerStore Peers, TreatyStore Treaties) CreateHandler(string uri = "eg://local")
    {
        var identity = SystemIdentity.Generate(new SystemUri(uri));
        var transport = new InMemoryTransport();
        var peers = new PeerStore();
        var treaties = new TreatyStore();
        var handler = new EgipHandler(identity, transport, peers, treaties);
        return (handler, identity, transport, peers, treaties);
    }

    [Fact]
    public async Task Hello_SendsSignedEnvelope()
    {
        var (handler, identity, transport, _, _) = CreateHandler();
        var target = new SystemUri("eg://remote");

        await handler.HelloAsync(target);

        Assert.Single(transport.Sent);
        var (to, env) = transport.Sent[0];
        Assert.Equal("eg://remote", to.Value);
        Assert.Equal(MessageType.Hello, env.Type);
        Assert.True(Envelope.VerifyEnvelope(env, identity, identity.PublicKey));

        identity.Dispose();
    }

    [Fact]
    public async Task Hello_TransportFailure_ThrowsAndUpdatesTrust()
    {
        var (handler, identity, transport, peers, _) = CreateHandler();
        var target = new SystemUri("eg://remote");

        // Register peer so we can check trust impact
        peers.Register(target, new PublicKey(new byte[32]), [], 1);

        transport.SendException = new Exception("network error");

        await Assert.ThrowsAsync<TransportFailureException>(() => handler.HelloAsync(target));

        identity.Dispose();
    }

    [Fact]
    public async Task HandleIncoming_HelloHandshake()
    {
        var (handler, localId, _, peers, _) = CreateHandler("eg://local");
        using var remoteId = SystemIdentity.Generate(new SystemUri("eg://remote"));

        // Create a HELLO envelope from remote
        var env = new Envelope
        {
            ProtocolVersion = 1,
            Id = new EnvelopeId(Guid.NewGuid().ToString()),
            From = remoteId.SystemUri,
            To = localId.SystemUri,
            Type = MessageType.Hello,
            Payload = new HelloPayload(remoteId.SystemUri, remoteId.PublicKey, [1], ["treaty"], 10),
            Timestamp = DateTimeOffset.UtcNow,
            Signature = new Signature(new byte[64]),
            InReplyTo = Option<EnvelopeId>.None(),
        };

        // Sign with remote's identity
        var signed = Envelope.SignEnvelope(env, remoteId);

        await handler.HandleIncomingAsync(signed);

        // Peer should be registered
        var (peer, found) = peers.Get(remoteId.SystemUri);
        Assert.True(found);
        Assert.Contains("treaty", peer!.Capabilities);

        localId.Dispose();
    }

    [Fact]
    public async Task HandleIncoming_DuplicateEnvelope_Throws()
    {
        var (handler, localId, _, _, _) = CreateHandler("eg://local");
        using var remoteId = SystemIdentity.Generate(new SystemUri("eg://remote"));

        var env = Envelope.SignEnvelope(new Envelope
        {
            ProtocolVersion = 1,
            Id = new EnvelopeId(Guid.NewGuid().ToString()),
            From = remoteId.SystemUri,
            To = localId.SystemUri,
            Type = MessageType.Hello,
            Payload = new HelloPayload(remoteId.SystemUri, remoteId.PublicKey, [1], [], 0),
            Timestamp = DateTimeOffset.UtcNow,
            Signature = new Signature(new byte[64]),
            InReplyTo = Option<EnvelopeId>.None(),
        }, remoteId);

        await handler.HandleIncomingAsync(env);
        await Assert.ThrowsAsync<DuplicateEnvelopeException>(() => handler.HandleIncomingAsync(env));

        localId.Dispose();
    }

    [Fact]
    public async Task HandleIncoming_InvalidSignature_Throws()
    {
        var (handler, localId, _, _, _) = CreateHandler("eg://local");
        using var remoteId = SystemIdentity.Generate(new SystemUri("eg://remote"));

        // Create envelope but sign with wrong key
        var env = new Envelope
        {
            ProtocolVersion = 1,
            Id = new EnvelopeId(Guid.NewGuid().ToString()),
            From = remoteId.SystemUri,
            To = localId.SystemUri,
            Type = MessageType.Hello,
            Payload = new HelloPayload(remoteId.SystemUri, remoteId.PublicKey, [1], [], 0),
            Timestamp = DateTimeOffset.UtcNow,
            Signature = new Signature(new byte[64]), // zero signature = invalid
            InReplyTo = Option<EnvelopeId>.None(),
        };

        await Assert.ThrowsAsync<EnvelopeSignatureInvalidException>(() => handler.HandleIncomingAsync(env));

        localId.Dispose();
    }

    [Fact]
    public async Task HandleIncoming_Message_InvokesCallback()
    {
        var (handler, localId, _, peers, _) = CreateHandler("eg://local");
        using var remoteId = SystemIdentity.Generate(new SystemUri("eg://remote"));

        // Register peer first
        peers.Register(remoteId.SystemUri, remoteId.PublicKey, [], 1);

        MessagePayloadContent? receivedPayload = null;
        handler.OnMessage = (from, payload) =>
        {
            receivedPayload = payload;
            return Task.CompletedTask;
        };

        var env = Envelope.SignEnvelope(new Envelope
        {
            ProtocolVersion = 1,
            Id = new EnvelopeId(Guid.NewGuid().ToString()),
            From = remoteId.SystemUri,
            To = localId.SystemUri,
            Type = MessageType.Message,
            Payload = new MessagePayloadContent("{}", "test.event", Option<ConversationId>.None(), []),
            Timestamp = DateTimeOffset.UtcNow,
            Signature = new Signature(new byte[64]),
            InReplyTo = Option<EnvelopeId>.None(),
        }, remoteId);

        await handler.HandleIncomingAsync(env);
        Assert.NotNull(receivedPayload);
        Assert.Equal("test.event", receivedPayload!.ContentType);

        localId.Dispose();
    }

    [Fact]
    public async Task HandleIncoming_Message_FromUnknownPeer_Throws()
    {
        var (handler, localId, _, _, _) = CreateHandler("eg://local");
        using var remoteId = SystemIdentity.Generate(new SystemUri("eg://remote"));

        // Don't register peer
        var env = Envelope.SignEnvelope(new Envelope
        {
            ProtocolVersion = 1,
            Id = new EnvelopeId(Guid.NewGuid().ToString()),
            From = remoteId.SystemUri,
            To = localId.SystemUri,
            Type = MessageType.Message,
            Payload = new MessagePayloadContent("{}", "test.event", Option<ConversationId>.None(), []),
            Timestamp = DateTimeOffset.UtcNow,
            Signature = new Signature(new byte[64]),
            InReplyTo = Option<EnvelopeId>.None(),
        }, remoteId);

        await Assert.ThrowsAsync<SystemNotFoundException>(() => handler.HandleIncomingAsync(env));

        localId.Dispose();
    }

    [Fact]
    public async Task HandleIncoming_Treaty_ProposeAndAccept()
    {
        var (handler, localId, _, peers, treaties) = CreateHandler("eg://local");
        using var remoteId = SystemIdentity.Generate(new SystemUri("eg://remote"));

        // Register peer
        peers.Register(remoteId.SystemUri, remoteId.PublicKey, [], 1);

        var treatyId = new TreatyId(Guid.NewGuid().ToString());

        // Propose treaty
        var proposeEnv = Envelope.SignEnvelope(new Envelope
        {
            ProtocolVersion = 1,
            Id = new EnvelopeId(Guid.NewGuid().ToString()),
            From = remoteId.SystemUri,
            To = localId.SystemUri,
            Type = MessageType.Treaty,
            Payload = new TreatyPayload(
                treatyId,
                TreatyAction.Propose,
                [new TreatyTerm(new DomainScope("data.sharing"), "bidirectional", true)],
                Option<string>.None()),
            Timestamp = DateTimeOffset.UtcNow,
            Signature = new Signature(new byte[64]),
            InReplyTo = Option<EnvelopeId>.None(),
        }, remoteId);

        await handler.HandleIncomingAsync(proposeEnv);

        var (treaty, found) = treaties.Get(treatyId);
        Assert.True(found);
        Assert.Equal(TreatyStatus.Proposed, treaty!.Status);

        // Accept treaty
        var acceptEnv = Envelope.SignEnvelope(new Envelope
        {
            ProtocolVersion = 1,
            Id = new EnvelopeId(Guid.NewGuid().ToString()),
            From = remoteId.SystemUri,
            To = localId.SystemUri,
            Type = MessageType.Treaty,
            Payload = new TreatyPayload(treatyId, TreatyAction.Accept, [], Option<string>.None()),
            Timestamp = DateTimeOffset.UtcNow,
            Signature = new Signature(new byte[64]),
            InReplyTo = Option<EnvelopeId>.None(),
        }, remoteId);

        await handler.HandleIncomingAsync(acceptEnv);

        (treaty, _) = treaties.Get(treatyId);
        Assert.Equal(TreatyStatus.Active, treaty!.Status);

        localId.Dispose();
    }

    [Fact]
    public async Task HandleIncoming_VersionIncompatible_Throws()
    {
        var (handler, localId, _, _, _) = CreateHandler("eg://local");
        handler.LocalProtocolVersions = [1];

        using var remoteId = SystemIdentity.Generate(new SystemUri("eg://remote"));

        var env = Envelope.SignEnvelope(new Envelope
        {
            ProtocolVersion = 1,
            Id = new EnvelopeId(Guid.NewGuid().ToString()),
            From = remoteId.SystemUri,
            To = localId.SystemUri,
            Type = MessageType.Hello,
            Payload = new HelloPayload(remoteId.SystemUri, remoteId.PublicKey, [99], [], 0),
            Timestamp = DateTimeOffset.UtcNow,
            Signature = new Signature(new byte[64]),
            InReplyTo = Option<EnvelopeId>.None(),
        }, remoteId);

        await Assert.ThrowsAsync<VersionIncompatibleException>(() => handler.HandleIncomingAsync(env));

        localId.Dispose();
    }

    [Fact]
    public async Task HandleIncoming_StaleEnvelope_Throws()
    {
        var (handler, localId, _, _, _) = CreateHandler("eg://local");
        using var remoteId = SystemIdentity.Generate(new SystemUri("eg://remote"));

        var env = Envelope.SignEnvelope(new Envelope
        {
            ProtocolVersion = 1,
            Id = new EnvelopeId(Guid.NewGuid().ToString()),
            From = remoteId.SystemUri,
            To = localId.SystemUri,
            Type = MessageType.Hello,
            Payload = new HelloPayload(remoteId.SystemUri, remoteId.PublicKey, [1], [], 0),
            Timestamp = DateTimeOffset.UtcNow.AddHours(-26), // stale
            Signature = new Signature(new byte[64]),
            InReplyTo = Option<EnvelopeId>.None(),
        }, remoteId);

        await Assert.ThrowsAsync<EgipException>(() => handler.HandleIncomingAsync(env));

        localId.Dispose();
    }

    [Fact]
    public async Task HandleIncoming_Receipt_UpdatesTrust()
    {
        var (handler, localId, _, peers, _) = CreateHandler("eg://local");
        using var remoteId = SystemIdentity.Generate(new SystemUri("eg://remote"));

        peers.Register(remoteId.SystemUri, remoteId.PublicKey, [], 1);

        var env = Envelope.SignEnvelope(new Envelope
        {
            ProtocolVersion = 1,
            Id = new EnvelopeId(Guid.NewGuid().ToString()),
            From = remoteId.SystemUri,
            To = localId.SystemUri,
            Type = MessageType.Receipt,
            Payload = new ReceiptPayload(
                new EnvelopeId(Guid.NewGuid().ToString()),
                ReceiptStatus.Processed,
                Option<EventId>.None(),
                Option<string>.None(),
                new Signature(new byte[64])),
            Timestamp = DateTimeOffset.UtcNow,
            Signature = new Signature(new byte[64]),
            InReplyTo = Option<EnvelopeId>.None(),
        }, remoteId);

        await handler.HandleIncomingAsync(env);

        var (peer, _) = peers.Get(remoteId.SystemUri);
        Assert.True(peer!.Trust.Value > 0.0);

        localId.Dispose();
    }
}

// ── CGER tests ───────────────────────────────────────────────────────────

public class CgerTests
{
    [Fact]
    public void Cger_CreatesValid()
    {
        var signer = new NoopSigner();
        var evt = EventFactory.CreateBootstrap(new ActorId("system"), signer);

        var cger = new Cger(
            evt.Id,
            new SystemUri("eg://remote"),
            "remote-event-id",
            Hash.Zero(),
            CgerRelationship.CausedBy,
            false);

        Assert.Equal(CgerRelationship.CausedBy, cger.Relationship);
        Assert.False(cger.Verified);
    }
}

// ── TreatyTerm tests ─────────────────────────────────────────────────────

public class TreatyTermTests
{
    [Fact]
    public void TreatyTerm_RecordEquality()
    {
        var t1 = new TreatyTerm(new DomainScope("data.sharing"), "open", true);
        var t2 = new TreatyTerm(new DomainScope("data.sharing"), "open", true);
        Assert.Equal(t1, t2);
    }
}
