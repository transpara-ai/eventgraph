use eventgraph::egip::*;
use eventgraph::types::*;

// ── Helper: build a test identity ───────────────────────────────────────

fn test_identity(name: &str) -> SystemIdentity {
    let uri = SystemUri::new(name).unwrap();
    SystemIdentity::generate(uri).unwrap()
}

fn test_envelope_id() -> EnvelopeId {
    // Use a valid UUID format.
    EnvelopeId::new("a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d").unwrap()
}

fn test_envelope_id_2() -> EnvelopeId {
    EnvelopeId::new("b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e").unwrap()
}

fn test_system_uri(name: &str) -> SystemUri {
    SystemUri::new(name).unwrap()
}

// ── Enum display tests ──────────────────────────────────────────────────

#[test]
fn test_message_type_as_str() {
    assert_eq!(MessageType::Hello.as_str(), "hello");
    assert_eq!(MessageType::Message.as_str(), "message");
    assert_eq!(MessageType::Receipt.as_str(), "receipt");
    assert_eq!(MessageType::Proof.as_str(), "proof");
    assert_eq!(MessageType::Treaty.as_str(), "treaty");
    assert_eq!(MessageType::AuthorityRequest.as_str(), "authorityrequest");
    assert_eq!(MessageType::Discover.as_str(), "discover");
}

#[test]
fn test_treaty_status_as_str() {
    assert_eq!(TreatyStatus::Proposed.as_str(), "Proposed");
    assert_eq!(TreatyStatus::Active.as_str(), "Active");
    assert_eq!(TreatyStatus::Suspended.as_str(), "Suspended");
    assert_eq!(TreatyStatus::Terminated.as_str(), "Terminated");
}

#[test]
fn test_treaty_action_as_str() {
    assert_eq!(TreatyAction::Propose.as_str(), "Propose");
    assert_eq!(TreatyAction::Accept.as_str(), "Accept");
    assert_eq!(TreatyAction::Modify.as_str(), "Modify");
    assert_eq!(TreatyAction::Suspend.as_str(), "Suspend");
    assert_eq!(TreatyAction::Terminate.as_str(), "Terminate");
}

#[test]
fn test_receipt_status_as_str() {
    assert_eq!(ReceiptStatus::Delivered.as_str(), "Delivered");
    assert_eq!(ReceiptStatus::Processed.as_str(), "Processed");
    assert_eq!(ReceiptStatus::Rejected.as_str(), "Rejected");
}

#[test]
fn test_proof_type_as_str() {
    assert_eq!(ProofType::ChainSegment.as_str(), "ChainSegment");
    assert_eq!(ProofType::EventExistence.as_str(), "EventExistence");
    assert_eq!(ProofType::ChainSummary.as_str(), "ChainSummary");
}

#[test]
fn test_cger_relationship_as_str() {
    assert_eq!(CGERRelationship::CausedBy.as_str(), "CausedBy");
    assert_eq!(CGERRelationship::References.as_str(), "References");
    assert_eq!(CGERRelationship::RespondsTo.as_str(), "RespondsTo");
}

#[test]
fn test_authority_level_as_str() {
    assert_eq!(AuthorityLevel::Required.as_str(), "Required");
    assert_eq!(AuthorityLevel::Recommended.as_str(), "Recommended");
    assert_eq!(AuthorityLevel::Notification.as_str(), "Notification");
}

// ── Identity tests ──────────────────────────────────────────────────────

#[test]
fn test_system_identity_generate() {
    let id = test_identity("system-alpha");
    assert_eq!(id.system_uri().value(), "system-alpha");
    assert_eq!(id.public_key().bytes().len(), 32);
}

#[test]
fn test_identity_sign_and_verify() {
    let id = test_identity("system-beta");
    let data = b"hello world";
    let sig = id.sign(data).unwrap();
    assert_eq!(sig.bytes().len(), 64);

    let valid = id.verify(id.public_key(), data, &sig).unwrap();
    assert!(valid, "signature should verify against own key");
}

#[test]
fn test_identity_verify_wrong_data() {
    let id = test_identity("system-gamma");
    let sig = id.sign(b"correct data").unwrap();
    let valid = id.verify(id.public_key(), b"wrong data", &sig).unwrap();
    assert!(!valid, "signature should not verify against wrong data");
}

#[test]
fn test_identity_from_key() {
    let id1 = test_identity("system-delta");
    let key = *id1.private_key();
    let id2 = SystemIdentity::from_key(
        SystemUri::new("system-delta").unwrap(),
        key,
    );
    assert_eq!(id2.public_key().bytes(), id1.public_key().bytes());
}

// ── Envelope canonical form and signing ─────────────────────────────────

#[test]
fn test_envelope_canonical_form() {
    let env = Envelope {
        protocol_version: 1,
        id: test_envelope_id(),
        from: test_system_uri("system-a"),
        to: test_system_uri("system-b"),
        message_type: MessageType::Hello,
        payload: MessagePayload::Hello(HelloPayload {
            system_uri: "system-a".to_string(),
            public_key: vec![0u8; 32],
            protocol_versions: vec![1],
            capabilities: vec!["treaty".to_string()],
            chain_length: 42,
        }),
        timestamp_nanos: 1000000000,
        signature: Signature::zero(),
        in_reply_to: None,
    };

    let canon = env.canonical_form().unwrap();
    assert!(canon.starts_with("1|"), "should start with version");
    assert!(canon.contains("system-a"), "should contain from URI");
    assert!(canon.contains("system-b"), "should contain to URI");
    assert!(canon.contains("hello"), "should contain message type");
    assert!(canon.contains("1000000000"), "should contain timestamp");
}

#[test]
fn test_envelope_canonical_form_with_reply_to() {
    let env = Envelope {
        protocol_version: 1,
        id: test_envelope_id(),
        from: test_system_uri("system-a"),
        to: test_system_uri("system-b"),
        message_type: MessageType::Receipt,
        payload: MessagePayload::Receipt(ReceiptPayload {
            envelope_id: "some-id".to_string(),
            status: "Delivered".to_string(),
            local_event_id: None,
            reason: None,
            signature: vec![0u8; 64],
        }),
        timestamp_nanos: 2000000000,
        signature: Signature::zero(),
        in_reply_to: Some(test_envelope_id_2()),
    };

    let canon = env.canonical_form().unwrap();
    assert!(canon.contains(test_envelope_id_2().value()), "should contain in_reply_to ID");
}

#[test]
fn test_sign_and_verify_envelope() {
    let id = test_identity("system-signer");
    let env = Envelope {
        protocol_version: 1,
        id: test_envelope_id(),
        from: test_system_uri("system-signer"),
        to: test_system_uri("system-receiver"),
        message_type: MessageType::Hello,
        payload: MessagePayload::Hello(HelloPayload {
            system_uri: "system-signer".to_string(),
            public_key: id.public_key().bytes().to_vec(),
            protocol_versions: vec![1],
            capabilities: vec![],
            chain_length: 0,
        }),
        timestamp_nanos: 3000000000,
        signature: Signature::zero(),
        in_reply_to: None,
    };

    let signed = sign_envelope(&env, &id).unwrap();
    assert_ne!(signed.signature.bytes(), Signature::zero().bytes());

    let valid = verify_envelope(&signed, &id, id.public_key()).unwrap();
    assert!(valid, "signed envelope should verify");
}

// ── Version negotiation ────────────────────────────────────────────────

#[test]
fn test_negotiate_version_common() {
    assert_eq!(negotiate_version(&[1, 2, 3], &[2, 3, 4]), Some(3));
}

#[test]
fn test_negotiate_version_single() {
    assert_eq!(negotiate_version(&[1], &[1]), Some(1));
}

#[test]
fn test_negotiate_version_none() {
    assert_eq!(negotiate_version(&[1, 2], &[3, 4]), None);
}

#[test]
fn test_negotiate_version_empty() {
    assert_eq!(negotiate_version(&[], &[1]), None);
    assert_eq!(negotiate_version(&[1], &[]), None);
}

// ── Treaty state machine ───────────────────────────────────────────────

fn test_treaty() -> Treaty {
    Treaty::new(
        TreatyId::new("a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d").unwrap(),
        test_system_uri("system-a"),
        test_system_uri("system-b"),
        vec![],
    )
}

#[test]
fn test_treaty_initial_status() {
    let treaty = test_treaty();
    assert_eq!(treaty.status, TreatyStatus::Proposed);
}

#[test]
fn test_treaty_propose_to_active() {
    let mut treaty = test_treaty();
    treaty.transition(TreatyStatus::Active).unwrap();
    assert_eq!(treaty.status, TreatyStatus::Active);
}

#[test]
fn test_treaty_propose_to_terminated() {
    let mut treaty = test_treaty();
    treaty.transition(TreatyStatus::Terminated).unwrap();
    assert_eq!(treaty.status, TreatyStatus::Terminated);
}

#[test]
fn test_treaty_propose_to_suspended_invalid() {
    let mut treaty = test_treaty();
    let result = treaty.transition(TreatyStatus::Suspended);
    assert!(result.is_err(), "Proposed -> Suspended should be invalid");
}

#[test]
fn test_treaty_active_to_suspended() {
    let mut treaty = test_treaty();
    treaty.transition(TreatyStatus::Active).unwrap();
    treaty.transition(TreatyStatus::Suspended).unwrap();
    assert_eq!(treaty.status, TreatyStatus::Suspended);
}

#[test]
fn test_treaty_suspended_to_active() {
    let mut treaty = test_treaty();
    treaty.transition(TreatyStatus::Active).unwrap();
    treaty.transition(TreatyStatus::Suspended).unwrap();
    treaty.transition(TreatyStatus::Active).unwrap();
    assert_eq!(treaty.status, TreatyStatus::Active);
}

#[test]
fn test_treaty_terminated_is_terminal() {
    let mut treaty = test_treaty();
    treaty.transition(TreatyStatus::Terminated).unwrap();
    let result = treaty.transition(TreatyStatus::Active);
    assert!(result.is_err(), "Terminated is a terminal state");
}

#[test]
fn test_treaty_apply_action_accept() {
    let mut treaty = test_treaty();
    treaty.apply_action(TreatyAction::Accept).unwrap();
    assert_eq!(treaty.status, TreatyStatus::Active);
}

#[test]
fn test_treaty_apply_action_modify_requires_active() {
    let mut treaty = test_treaty();
    let result = treaty.apply_action(TreatyAction::Modify);
    assert!(result.is_err(), "Modify requires Active status");

    treaty.apply_action(TreatyAction::Accept).unwrap();
    treaty.apply_action(TreatyAction::Modify).unwrap();
    assert_eq!(treaty.status, TreatyStatus::Active); // status unchanged
}

#[test]
fn test_treaty_apply_action_propose_on_existing() {
    let mut treaty = test_treaty();
    let result = treaty.apply_action(TreatyAction::Propose);
    assert!(result.is_err(), "Propose on existing treaty should fail");
}

// ── PeerStore ──────────────────────────────────────────────────────────

#[test]
fn test_peer_store_register_and_get() {
    let store = PeerStore::new();
    let uri = test_system_uri("peer-alpha");
    let pk = PublicKey::new([1u8; 32]);
    store.register(uri.clone(), pk.clone(), vec!["treaty".to_string()], 1);

    let peer = store.get(&uri).unwrap();
    assert_eq!(peer.system_uri.value(), "peer-alpha");
    assert_eq!(peer.public_key.bytes(), pk.bytes());
    assert_eq!(peer.negotiated_version, 1);
    assert_eq!(peer.trust.value(), 0.0);
}

#[test]
fn test_peer_store_get_unknown() {
    let store = PeerStore::new();
    let uri = test_system_uri("unknown");
    assert!(store.get(&uri).is_none());
}

#[test]
fn test_peer_store_reregister_preserves_key() {
    let store = PeerStore::new();
    let uri = test_system_uri("peer-beta");
    let pk1 = PublicKey::new([1u8; 32]);
    let pk2 = PublicKey::new([2u8; 32]);

    store.register(uri.clone(), pk1.clone(), vec![], 1);
    store.register(uri.clone(), pk2, vec!["new-cap".to_string()], 2);

    let peer = store.get(&uri).unwrap();
    // Public key should NOT be overwritten (security: prevents key-substitution).
    assert_eq!(peer.public_key.bytes(), pk1.bytes());
    // But capabilities and version should be updated.
    assert_eq!(peer.capabilities, vec!["new-cap"]);
    assert_eq!(peer.negotiated_version, 2);
}

#[test]
fn test_peer_store_update_trust_positive() {
    let store = PeerStore::new();
    let uri = test_system_uri("peer-trust");
    store.register(uri.clone(), PublicKey::new([0u8; 32]), vec![], 1);

    // Positive delta capped at INTER_SYSTEM_MAX_ADJUSTMENT (0.05).
    let score = store.update_trust(&uri, 0.10).unwrap();
    assert!((score.value() - 0.05).abs() < 1e-9, "should cap at 0.05");
}

#[test]
fn test_peer_store_update_trust_negative() {
    let store = PeerStore::new();
    let uri = test_system_uri("peer-untrust");
    store.register(uri.clone(), PublicKey::new([0u8; 32]), vec![], 1);

    // Start at 0, negative should clamp to 0.
    let score = store.update_trust(&uri, -0.20).unwrap();
    assert_eq!(score.value(), 0.0);
}

#[test]
fn test_peer_store_update_trust_unknown() {
    let store = PeerStore::new();
    let uri = test_system_uri("unknown-peer");
    assert!(store.update_trust(&uri, 0.01).is_none());
}

#[test]
fn test_peer_store_all() {
    let store = PeerStore::new();
    store.register(test_system_uri("p1"), PublicKey::new([0u8; 32]), vec![], 1);
    store.register(test_system_uri("p2"), PublicKey::new([1u8; 32]), vec![], 1);

    let all = store.all();
    assert_eq!(all.len(), 2);
}

// ── TreatyStore ────────────────────────────────────────────────────────

#[test]
fn test_treaty_store_put_and_get() {
    let store = TreatyStore::new();
    let treaty = test_treaty();
    let id = treaty.id.clone();
    store.put(treaty);

    let retrieved = store.get(&id).unwrap();
    assert_eq!(retrieved.status, TreatyStatus::Proposed);
}

#[test]
fn test_treaty_store_get_unknown() {
    let store = TreatyStore::new();
    let id = TreatyId::new("a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d").unwrap();
    assert!(store.get(&id).is_none());
}

#[test]
fn test_treaty_store_apply() {
    let store = TreatyStore::new();
    let treaty = test_treaty();
    let id = treaty.id.clone();
    store.put(treaty);

    store.apply(&id, |t| t.apply_action(TreatyAction::Accept)).unwrap();

    let retrieved = store.get(&id).unwrap();
    assert_eq!(retrieved.status, TreatyStatus::Active);
}

#[test]
fn test_treaty_store_apply_not_found() {
    let store = TreatyStore::new();
    let id = TreatyId::new("a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d").unwrap();
    let result = store.apply(&id, |t| t.apply_action(TreatyAction::Accept));
    assert!(result.is_err());
}

#[test]
fn test_treaty_store_by_system() {
    let store = TreatyStore::new();
    let t1 = Treaty::new(
        TreatyId::new("a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d").unwrap(),
        test_system_uri("sys-a"),
        test_system_uri("sys-b"),
        vec![],
    );
    let t2 = Treaty::new(
        TreatyId::new("b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e").unwrap(),
        test_system_uri("sys-c"),
        test_system_uri("sys-a"),
        vec![],
    );
    store.put(t1);
    store.put(t2);

    let results = store.by_system(&test_system_uri("sys-a"));
    assert_eq!(results.len(), 2);
}

#[test]
fn test_treaty_store_active() {
    let store = TreatyStore::new();
    let mut t1 = test_treaty();
    t1.status = TreatyStatus::Active;
    store.put(t1);

    let t2 = Treaty::new(
        TreatyId::new("b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e").unwrap(),
        test_system_uri("sys-c"),
        test_system_uri("sys-d"),
        vec![],
    );
    store.put(t2);

    let active = store.active();
    assert_eq!(active.len(), 1);
    assert_eq!(active[0].status, TreatyStatus::Active);
}

// ── EnvelopeDedup ──────────────────────────────────────────────────────

#[test]
fn test_dedup_first_check_returns_true() {
    let dedup = EnvelopeDedup::new();
    let id = test_envelope_id();
    assert!(dedup.check(&id), "first check should return true");
}

#[test]
fn test_dedup_second_check_returns_false() {
    let dedup = EnvelopeDedup::new();
    let id = test_envelope_id();
    dedup.check(&id);
    assert!(!dedup.check(&id), "second check should return false (duplicate)");
}

#[test]
fn test_dedup_different_ids() {
    let dedup = EnvelopeDedup::new();
    let id1 = test_envelope_id();
    let id2 = test_envelope_id_2();
    assert!(dedup.check(&id1));
    assert!(dedup.check(&id2));
    assert!(!dedup.check(&id1));
    assert!(!dedup.check(&id2));
}

#[test]
fn test_dedup_size() {
    let dedup = EnvelopeDedup::new();
    assert_eq!(dedup.size(), 0);
    dedup.check(&test_envelope_id());
    assert_eq!(dedup.size(), 1);
    dedup.check(&test_envelope_id_2());
    assert_eq!(dedup.size(), 2);
}

#[test]
fn test_dedup_prune_no_expired() {
    let dedup = EnvelopeDedup::new();
    dedup.check(&test_envelope_id());
    let removed = dedup.prune();
    assert_eq!(removed, 0, "nothing should be pruned when TTL hasn't expired");
}

#[test]
fn test_dedup_with_short_ttl() {
    let dedup = EnvelopeDedup::with_ttl(std::time::Duration::from_millis(1));
    dedup.check(&test_envelope_id());
    std::thread::sleep(std::time::Duration::from_millis(10));
    let removed = dedup.prune();
    assert_eq!(removed, 1, "expired entry should be pruned");
}

// ── Proof verification ─────────────────────────────────────────────────

#[test]
fn test_verify_chain_segment_valid() {
    let proof = ChainSegmentProof {
        event_hashes: vec!["hash1".to_string(), "hash2".to_string(), "hash3".to_string()],
        start_hash: "hash0".to_string(),
        end_hash: "hash3".to_string(),
    };
    assert!(verify_chain_segment(&proof));
}

#[test]
fn test_verify_chain_segment_empty() {
    let proof = ChainSegmentProof {
        event_hashes: vec![],
        start_hash: "hash0".to_string(),
        end_hash: "hash1".to_string(),
    };
    assert!(!verify_chain_segment(&proof));
}

#[test]
fn test_verify_chain_segment_end_hash_mismatch() {
    let proof = ChainSegmentProof {
        event_hashes: vec!["hash1".to_string(), "hash2".to_string()],
        start_hash: "hash0".to_string(),
        end_hash: "wrong".to_string(),
    };
    assert!(!verify_chain_segment(&proof));
}

#[test]
fn test_verify_event_existence_valid() {
    let proof = EventExistenceProof {
        event_hash: "abc123".to_string(),
        prev_hash: "prev".to_string(),
        next_hash: Some("next".to_string()),
        position: 5,
        chain_length: 100,
    };
    assert!(verify_event_existence(&proof));
}

#[test]
fn test_verify_event_existence_negative_position() {
    let proof = EventExistenceProof {
        event_hash: "abc123".to_string(),
        prev_hash: "prev".to_string(),
        next_hash: None,
        position: -1,
        chain_length: 100,
    };
    assert!(!verify_event_existence(&proof));
}

#[test]
fn test_verify_event_existence_position_exceeds_length() {
    let proof = EventExistenceProof {
        event_hash: "abc123".to_string(),
        prev_hash: "prev".to_string(),
        next_hash: None,
        position: 100,
        chain_length: 100,
    };
    assert!(!verify_event_existence(&proof));
}

#[test]
fn test_verify_event_existence_empty_hash() {
    let proof = EventExistenceProof {
        event_hash: "".to_string(),
        prev_hash: "prev".to_string(),
        next_hash: None,
        position: 0,
        chain_length: 10,
    };
    assert!(!verify_event_existence(&proof));
}

#[test]
fn test_validate_proof_chain_summary() {
    let payload = ProofPayload {
        proof_type: "ChainSummary".to_string(),
        data: ProofData::ChainSummary(ChainSummaryProof {
            length: 42,
            head_hash: "head".to_string(),
            genesis_hash: "genesis".to_string(),
            timestamp_nanos: 1000000000,
        }),
    };
    assert!(validate_proof(&payload).unwrap());
}

#[test]
fn test_validate_proof_chain_summary_zero_length() {
    let payload = ProofPayload {
        proof_type: "ChainSummary".to_string(),
        data: ProofData::ChainSummary(ChainSummaryProof {
            length: 0,
            head_hash: "head".to_string(),
            genesis_hash: "genesis".to_string(),
            timestamp_nanos: 1000000000,
        }),
    };
    assert!(!validate_proof(&payload).unwrap());
}

#[test]
fn test_proof_type_from_data() {
    assert_eq!(
        proof_type_from_data(&ProofData::ChainSegment(ChainSegmentProof {
            event_hashes: vec![], start_hash: "".into(), end_hash: "".into(),
        })),
        ProofType::ChainSegment,
    );
    assert_eq!(
        proof_type_from_data(&ProofData::EventExistence(EventExistenceProof {
            event_hash: "".into(), prev_hash: "".into(), next_hash: None, position: 0, chain_length: 0,
        })),
        ProofType::EventExistence,
    );
    assert_eq!(
        proof_type_from_data(&ProofData::ChainSummary(ChainSummaryProof {
            length: 0, head_hash: "".into(), genesis_hash: "".into(), timestamp_nanos: 0,
        })),
        ProofType::ChainSummary,
    );
}

// ── CGER ───────────────────────────────────────────────────────────────

#[test]
fn test_cger_relationship_parsing() {
    let cger = CGER {
        local_event_id: "evt1".to_string(),
        remote_system: "remote".to_string(),
        remote_event_id: "rev1".to_string(),
        remote_hash: "hash1".to_string(),
        relationship: "CausedBy".to_string(),
        verified: false,
    };
    assert_eq!(cger.cger_relationship(), Some(CGERRelationship::CausedBy));

    let cger2 = CGER {
        relationship: "Unknown".to_string(),
        ..cger
    };
    assert_eq!(cger2.cger_relationship(), None);
}

// ── TreatyPayload parsing ──────────────────────────────────────────────

#[test]
fn test_treaty_payload_action_parsing() {
    let payload = TreatyPayload {
        treaty_id: "abc".to_string(),
        action: "Accept".to_string(),
        terms: vec![],
        reason: None,
    };
    assert_eq!(payload.treaty_action(), Some(TreatyAction::Accept));

    let unknown = TreatyPayload { action: "Unknown".to_string(), ..payload };
    assert_eq!(unknown.treaty_action(), None);
}

// ── ReceiptPayload parsing ─────────────────────────────────────────────

#[test]
fn test_receipt_payload_status_parsing() {
    let r = ReceiptPayload {
        envelope_id: "id".to_string(),
        status: "Processed".to_string(),
        local_event_id: None,
        reason: None,
        signature: vec![0; 64],
    };
    assert_eq!(r.receipt_status(), Some(ReceiptStatus::Processed));
}

// ── EgipError display ──────────────────────────────────────────────────

#[test]
fn test_egip_error_display() {
    let err = EgipError::SystemNotFound { uri: "sys-x".to_string() };
    assert_eq!(err.to_string(), "system not found: sys-x");

    let err = EgipError::DuplicateEnvelope { envelope_id: "env-1".to_string() };
    assert_eq!(err.to_string(), "duplicate envelope: env-1");

    let err = EgipError::VersionIncompatible { local: vec![1], remote: vec![2] };
    assert!(err.to_string().contains("no compatible protocol version"));

    let err = EgipError::TreatyNotFound { treaty_id: "t1".to_string() };
    assert_eq!(err.to_string(), "treaty not found: t1");
}

// ── Handler integration (minimal) ──────────────────────────────────────

/// A test transport that records sent envelopes.
struct MockTransport {
    sent: std::sync::Mutex<Vec<(String, Envelope)>>,
}

impl MockTransport {
    fn new() -> Self {
        Self {
            sent: std::sync::Mutex::new(Vec::new()),
        }
    }
}

impl Transport for MockTransport {
    fn send(&self, to: &SystemUri, envelope: &Envelope) -> EgipResult<Option<ReceiptPayload>> {
        let mut sent = self.sent.lock().unwrap();
        sent.push((to.value().to_string(), envelope.clone()));
        Ok(None)
    }
}

#[test]
fn test_handler_hello() {
    let id = test_identity("handler-system");
    let transport = MockTransport::new();
    let handler = Handler::new(
        Box::new(SystemIdentity::from_key(
            SystemUri::new("handler-system").unwrap(),
            *id.private_key(),
        )),
        Box::new(transport),
    );

    let target = test_system_uri("remote-system");
    handler.hello(&target).unwrap();
}

#[test]
fn test_handler_handle_incoming_hello() {
    let id_a = test_identity("system-a");
    let id_b = test_identity("system-b");

    let handler = Handler::new(
        Box::new(SystemIdentity::from_key(
            SystemUri::new("system-b").unwrap(),
            *id_b.private_key(),
        )),
        Box::new(MockTransport::new()),
    );

    // Build a HELLO envelope from system-a.
    let env = Envelope {
        protocol_version: 1,
        id: test_envelope_id(),
        from: test_system_uri("system-a"),
        to: test_system_uri("system-b"),
        message_type: MessageType::Hello,
        payload: MessagePayload::Hello(HelloPayload {
            system_uri: "system-a".to_string(),
            public_key: id_a.public_key().bytes().to_vec(),
            protocol_versions: vec![1],
            capabilities: vec!["treaty".to_string()],
            chain_length: 10,
        }),
        timestamp_nanos: std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .unwrap()
            .as_nanos() as u64,
        signature: Signature::zero(),
        in_reply_to: None,
    };

    // Sign with system-a's identity.
    let signed = sign_envelope(&env, &id_a).unwrap();

    // Handler-b uses id_b which cannot verify id_a's signature (our test impl
    // only verifies own signatures). So this will fail with signature invalid.
    // This is expected — the test validates the flow, not real crypto.
    let result = handler.handle_incoming(&signed);
    // With our simplified crypto, cross-system verification fails as expected.
    assert!(result.is_err());
}

#[test]
fn test_handler_handle_incoming_duplicate() {
    let id = test_identity("system-x");
    let handler = Handler::new(
        Box::new(SystemIdentity::from_key(
            SystemUri::new("system-x").unwrap(),
            *id.private_key(),
        )),
        Box::new(MockTransport::new()),
    );

    let env = Envelope {
        protocol_version: 1,
        id: test_envelope_id(),
        from: test_system_uri("system-x"),
        to: test_system_uri("system-x"),
        message_type: MessageType::Hello,
        payload: MessagePayload::Hello(HelloPayload {
            system_uri: "system-x".to_string(),
            public_key: id.public_key().bytes().to_vec(),
            protocol_versions: vec![1],
            capabilities: vec![],
            chain_length: 0,
        }),
        timestamp_nanos: std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .unwrap()
            .as_nanos() as u64,
        signature: Signature::zero(),
        in_reply_to: None,
    };

    let signed = sign_envelope(&env, &id).unwrap();

    // First call should succeed (self-verify works).
    let r1 = handler.handle_incoming(&signed);
    assert!(r1.is_ok(), "first incoming should succeed: {:?}", r1);

    // Second call with same envelope ID should be rejected as duplicate.
    let r2 = handler.handle_incoming(&signed);
    assert!(r2.is_err());
    let err_msg = r2.unwrap_err().to_string();
    assert!(err_msg.contains("duplicate"), "should be duplicate error: {err_msg}");
}

#[test]
fn test_handler_handle_incoming_unknown_sender() {
    let id = test_identity("system-y");
    let handler = Handler::new(
        Box::new(SystemIdentity::from_key(
            SystemUri::new("system-y").unwrap(),
            *id.private_key(),
        )),
        Box::new(MockTransport::new()),
    );

    // Non-HELLO from unknown sender should fail.
    let env = Envelope {
        protocol_version: 1,
        id: test_envelope_id(),
        from: test_system_uri("unknown-sender"),
        to: test_system_uri("system-y"),
        message_type: MessageType::Message,
        payload: MessagePayload::Message(MessagePayloadContent {
            content: serde_json::json!({"text": "hi"}),
            content_type: "chat.message".to_string(),
            conversation_id: None,
            cgers: vec![],
        }),
        timestamp_nanos: std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .unwrap()
            .as_nanos() as u64,
        signature: Signature::zero(),
        in_reply_to: None,
    };

    let result = handler.handle_incoming(&env);
    assert!(result.is_err());
    let err_msg = result.unwrap_err().to_string();
    assert!(err_msg.contains("system not found"), "should be system not found: {err_msg}");
}

// ── Thread safety ──────────────────────────────────────────────────────

#[test]
fn test_peer_store_concurrent_access() {
    use std::sync::Arc;
    use std::thread;

    let store = Arc::new(PeerStore::new());
    let mut handles = vec![];

    for i in 0..10 {
        let store = Arc::clone(&store);
        handles.push(thread::spawn(move || {
            let uri = SystemUri::new(&format!("peer-{i}")).unwrap();
            store.register(uri.clone(), PublicKey::new([i as u8; 32]), vec![], 1);
            store.update_trust(&uri, 0.01);
            store.get(&uri);
        }));
    }

    for h in handles {
        h.join().unwrap();
    }

    assert_eq!(store.all().len(), 10);
}

#[test]
fn test_dedup_concurrent_access() {
    use std::sync::Arc;
    use std::thread;

    let dedup = Arc::new(EnvelopeDedup::new());
    let id = test_envelope_id();
    let mut handles = vec![];

    let success_count = Arc::new(std::sync::atomic::AtomicU32::new(0));

    for _ in 0..10 {
        let dedup = Arc::clone(&dedup);
        let id = id.clone();
        let sc = Arc::clone(&success_count);
        handles.push(thread::spawn(move || {
            if dedup.check(&id) {
                sc.fetch_add(1, std::sync::atomic::Ordering::Relaxed);
            }
        }));
    }

    for h in handles {
        h.join().unwrap();
    }

    // Exactly one thread should succeed.
    assert_eq!(
        success_count.load(std::sync::atomic::Ordering::Relaxed),
        1,
        "only one thread should see a new envelope ID",
    );
}

// ── Constants ──────────────────────────────────────────────────────────

#[test]
fn test_trust_impact_constants() {
    assert!(TRUST_IMPACT_VALID_PROOF > 0.0);
    assert!(TRUST_IMPACT_RECEIPT_ON_TIME > 0.0);
    assert!(TRUST_IMPACT_TREATY_HONOURED > 0.0);
    assert!(TRUST_IMPACT_TREATY_VIOLATED < 0.0);
    assert!(TRUST_IMPACT_INVALID_PROOF < 0.0);
    assert!(TRUST_IMPACT_SIGNATURE_INVALID < 0.0);
    assert!(TRUST_IMPACT_NO_HELLO_RESPONSE < 0.0);
}

#[test]
fn test_current_protocol_version() {
    assert_eq!(CURRENT_PROTOCOL_VERSION, 1);
}
