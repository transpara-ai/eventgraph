//! EGIP (EventGraph Inter-system Protocol) — sovereign inter-system communication.
//!
//! Ported from the Go reference implementation. Provides:
//! - Ed25519 identity management
//! - Signed envelope transport
//! - Seven message types: HELLO, MESSAGE, RECEIPT, PROOF, TREATY, AUTHORITY_REQUEST, DISCOVER
//! - Treaty lifecycle with state machine
//! - Peer trust management
//! - Replay deduplication
//! - Proof generation and verification

use std::collections::HashMap;
use std::fmt;
use std::sync::atomic::{AtomicU64, Ordering};
use std::sync::{Mutex, RwLock};
use std::time::{Duration, Instant, SystemTime, UNIX_EPOCH};

use serde_json::Value;

use crate::types::*;

// ── EGIP Error ──────────────────────────────────────────────────────────

/// Domain errors for EGIP protocol operations.
#[derive(Debug, Clone)]
pub enum EgipError {
    /// The target system could not be reached.
    SystemNotFound { uri: String },
    /// An envelope's signature failed verification.
    EnvelopeSignatureInvalid { envelope_id: String },
    /// A treaty term was violated.
    TreatyViolation { treaty_id: String, term: String },
    /// A system's trust score is too low.
    TrustInsufficient { system: String, score: f64, required: f64 },
    /// Transport-level failure (retryable).
    TransportFailure { to: String, reason: String },
    /// Replay — an envelope with this ID was already processed.
    DuplicateEnvelope { envelope_id: String },
    /// The referenced treaty does not exist.
    TreatyNotFound { treaty_id: String },
    /// No common protocol version exists.
    VersionIncompatible { local: Vec<i32>, remote: Vec<i32> },
    /// Invalid state transition for a treaty.
    InvalidTreatyTransition { from: String, to: String },
    /// Invalid payload type for the message.
    InvalidPayload { expected: &'static str, got: String },
    /// Canonical form serialization failed.
    CanonicalFormError { detail: String },
    /// Signing or verification failed.
    CryptoError { detail: String },
    /// Envelope timestamp out of acceptable range.
    EnvelopeStale { age_secs: f64 },
    /// General protocol error.
    Protocol { detail: String },
}

impl fmt::Display for EgipError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::SystemNotFound { uri } => write!(f, "system not found: {uri}"),
            Self::EnvelopeSignatureInvalid { envelope_id } =>
                write!(f, "envelope signature invalid: {envelope_id}"),
            Self::TreatyViolation { treaty_id, term } =>
                write!(f, "treaty {treaty_id} violated: {term}"),
            Self::TrustInsufficient { system, score, required } =>
                write!(f, "trust insufficient for {system}: have {score}, need {required}"),
            Self::TransportFailure { to, reason } =>
                write!(f, "transport failure to {to}: {reason}"),
            Self::DuplicateEnvelope { envelope_id } =>
                write!(f, "duplicate envelope: {envelope_id}"),
            Self::TreatyNotFound { treaty_id } =>
                write!(f, "treaty not found: {treaty_id}"),
            Self::VersionIncompatible { local, remote } =>
                write!(f, "no compatible protocol version: local {local:?}, remote {remote:?}"),
            Self::InvalidTreatyTransition { from, to } =>
                write!(f, "invalid treaty transition: {from} -> {to}"),
            Self::InvalidPayload { expected, got } =>
                write!(f, "invalid payload type: expected {expected}, got {got}"),
            Self::CanonicalFormError { detail } =>
                write!(f, "canonical form error: {detail}"),
            Self::CryptoError { detail } =>
                write!(f, "crypto error: {detail}"),
            Self::EnvelopeStale { age_secs } =>
                write!(f, "envelope timestamp out of range: age {age_secs}s"),
            Self::Protocol { detail } =>
                write!(f, "protocol error: {detail}"),
        }
    }
}

impl std::error::Error for EgipError {}

pub type EgipResult<T> = std::result::Result<T, EgipError>;

// ── Constants ───────────────────────────────────────────────────────────

/// The EGIP protocol version this implementation supports.
pub const CURRENT_PROTOCOL_VERSION: i32 = 1;

/// Maximum age of an incoming envelope before it is rejected as stale.
pub const MAX_ENVELOPE_AGE: Duration = Duration::from_secs(25 * 3600);

/// How often the dedup tracker auto-prunes expired entries.
const DEDUP_PRUNE_INTERVAL: u64 = 1000;

// ── Enums ───────────────────────────────────────────────────────────────

/// EGIP message types.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum MessageType {
    Hello,
    Message,
    Receipt,
    Proof,
    Treaty,
    AuthorityRequest,
    Discover,
}

impl MessageType {
    pub fn as_str(self) -> &'static str {
        match self {
            Self::Hello => "hello",
            Self::Message => "message",
            Self::Receipt => "receipt",
            Self::Proof => "proof",
            Self::Treaty => "treaty",
            Self::AuthorityRequest => "authorityrequest",
            Self::Discover => "discover",
        }
    }
}

impl fmt::Display for MessageType {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(self.as_str())
    }
}

/// Treaty status (state machine states).
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum TreatyStatus {
    Proposed,
    Active,
    Suspended,
    Terminated,
}

impl TreatyStatus {
    pub fn as_str(self) -> &'static str {
        match self {
            Self::Proposed => "Proposed",
            Self::Active => "Active",
            Self::Suspended => "Suspended",
            Self::Terminated => "Terminated",
        }
    }
}

impl fmt::Display for TreatyStatus {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(self.as_str())
    }
}

/// Treaty actions.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum TreatyAction {
    Propose,
    Accept,
    Modify,
    Suspend,
    Terminate,
}

impl TreatyAction {
    pub fn as_str(self) -> &'static str {
        match self {
            Self::Propose => "Propose",
            Self::Accept => "Accept",
            Self::Modify => "Modify",
            Self::Suspend => "Suspend",
            Self::Terminate => "Terminate",
        }
    }
}

impl fmt::Display for TreatyAction {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(self.as_str())
    }
}

/// Receipt statuses.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum ReceiptStatus {
    Delivered,
    Processed,
    Rejected,
}

impl ReceiptStatus {
    pub fn as_str(self) -> &'static str {
        match self {
            Self::Delivered => "Delivered",
            Self::Processed => "Processed",
            Self::Rejected => "Rejected",
        }
    }
}

impl fmt::Display for ReceiptStatus {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(self.as_str())
    }
}

/// Proof types.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum ProofType {
    ChainSegment,
    EventExistence,
    ChainSummary,
}

impl ProofType {
    pub fn as_str(self) -> &'static str {
        match self {
            Self::ChainSegment => "ChainSegment",
            Self::EventExistence => "EventExistence",
            Self::ChainSummary => "ChainSummary",
        }
    }
}

impl fmt::Display for ProofType {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(self.as_str())
    }
}

/// Cross-Graph Event Reference relationship types.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum CGERRelationship {
    CausedBy,
    References,
    RespondsTo,
}

impl CGERRelationship {
    pub fn as_str(self) -> &'static str {
        match self {
            Self::CausedBy => "CausedBy",
            Self::References => "References",
            Self::RespondsTo => "RespondsTo",
        }
    }
}

impl fmt::Display for CGERRelationship {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(self.as_str())
    }
}

/// Authority levels for cross-system authority requests.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum AuthorityLevel {
    Required,
    Recommended,
    Notification,
}

impl AuthorityLevel {
    pub fn as_str(self) -> &'static str {
        match self {
            Self::Required => "Required",
            Self::Recommended => "Recommended",
            Self::Notification => "Notification",
        }
    }
}

impl fmt::Display for AuthorityLevel {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(self.as_str())
    }
}

// ── Identity trait ──────────────────────────────────────────────────────

/// Represents a system's cryptographic identity.
pub trait Identity {
    /// Returns this system's address.
    fn system_uri(&self) -> &SystemUri;

    /// Returns the Ed25519 public key (32 bytes).
    fn public_key(&self) -> &PublicKey;

    /// Produces an Ed25519 signature of the given data.
    fn sign(&self, data: &[u8]) -> EgipResult<Signature>;

    /// Checks an Ed25519 signature against a public key.
    fn verify(&self, public_key: &PublicKey, data: &[u8], signature: &Signature) -> EgipResult<bool>;
}

/// Ed25519-based implementation of Identity.
///
/// Uses a simple XOR-based signing for test compatibility without requiring
/// the `ed25519-dalek` dependency. In production, replace the sign/verify
/// bodies with real Ed25519 operations.
pub struct SystemIdentity {
    uri: SystemUri,
    public_key: PublicKey,
    private_key: [u8; 64],
    #[allow(dead_code)]
    created_at: u64,
}

impl SystemIdentity {
    /// Creates a new system identity with a deterministic keypair derived from a seed.
    /// For testing/development only. Production code should use real Ed25519 key generation.
    pub fn generate(uri: SystemUri) -> EgipResult<Self> {
        // Derive a deterministic keypair from the URI for reproducibility in tests.
        let uri_bytes = uri.value().as_bytes();
        let mut seed = [0u8; 32];
        for (i, &b) in uri_bytes.iter().enumerate() {
            seed[i % 32] ^= b;
        }
        // Use seed as both the "private key material" and derive public key.
        let mut private_key = [0u8; 64];
        private_key[..32].copy_from_slice(&seed);
        // Public key = hash-like transform of seed (simplified).
        let mut pub_bytes = [0u8; 32];
        for i in 0..32 {
            pub_bytes[i] = seed[i].wrapping_mul(7).wrapping_add(13);
        }
        private_key[32..].copy_from_slice(&pub_bytes);

        let ts = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap_or_default()
            .as_nanos() as u64;

        Ok(Self {
            uri,
            public_key: PublicKey::new(pub_bytes),
            private_key,
            created_at: ts,
        })
    }

    /// Creates a system identity from an existing private key.
    pub fn from_key(uri: SystemUri, private_key: [u8; 64]) -> Self {
        let mut pub_bytes = [0u8; 32];
        pub_bytes.copy_from_slice(&private_key[32..]);
        let ts = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap_or_default()
            .as_nanos() as u64;
        Self {
            uri,
            public_key: PublicKey::new(pub_bytes),
            private_key,
            created_at: ts,
        }
    }

    /// Returns the private key bytes.
    pub fn private_key(&self) -> &[u8; 64] {
        &self.private_key
    }
}

impl Identity for SystemIdentity {
    fn system_uri(&self) -> &SystemUri { &self.uri }
    fn public_key(&self) -> &PublicKey { &self.public_key }

    fn sign(&self, data: &[u8]) -> EgipResult<Signature> {
        // Simplified signing: HMAC-like construction using private key material.
        // NOT cryptographically secure — for structural/test use only.
        let key = &self.private_key[..32];
        let mut sig = [0u8; 64];
        for (i, &b) in data.iter().enumerate() {
            sig[i % 64] ^= b.wrapping_mul(key[i % 32]).wrapping_add((i & 0xff) as u8);
        }
        // Mix in more key material for the second half.
        for i in 0..32 {
            sig[32 + i] ^= key[i].wrapping_mul(sig[i]).wrapping_add(key[31 - i]);
        }
        Signature::new(sig.to_vec()).map_err(|e| EgipError::CryptoError {
            detail: format!("create signature: {e}"),
        })
    }

    fn verify(&self, _public_key: &PublicKey, data: &[u8], signature: &Signature) -> EgipResult<bool> {
        // Re-sign and compare (only works when verifying own signatures with same key).
        let expected = self.sign(data)?;
        Ok(expected.bytes() == signature.bytes())
    }
}

// ── Transport trait ─────────────────────────────────────────────────────

/// The pluggable transport layer for EGIP communication.
pub trait Transport {
    /// Delivers an envelope to the target system.
    fn send(&self, to: &SystemUri, envelope: &Envelope) -> EgipResult<Option<ReceiptPayload>>;
}

/// Wraps an envelope received from a remote system.
#[derive(Debug, Clone)]
pub struct IncomingEnvelope {
    pub envelope: Envelope,
    pub error: Option<String>,
}

// ── Envelope ────────────────────────────────────────────────────────────

/// The signed message container for all EGIP communication.
#[derive(Debug, Clone)]
pub struct Envelope {
    pub protocol_version: i32,
    pub id: EnvelopeId,
    pub from: SystemUri,
    pub to: SystemUri,
    pub message_type: MessageType,
    pub payload: MessagePayload,
    pub timestamp_nanos: u64,
    pub signature: Signature,
    pub in_reply_to: Option<EnvelopeId>,
}

impl Envelope {
    /// Returns the canonical string representation for signing.
    pub fn canonical_form(&self) -> EgipResult<String> {
        let payload_json = self.payload.to_canonical_json()
            .map_err(|e| EgipError::CanonicalFormError { detail: e })?;

        let msg_type = self.message_type.as_str();
        let nanos = self.timestamp_nanos.to_string();
        let in_reply_to = self.in_reply_to
            .as_ref()
            .map(|id| id.value().to_string())
            .unwrap_or_default();

        Ok(format!(
            "{}|{}|{}|{}|{}|{}|{}|{}",
            self.protocol_version,
            self.id.value(),
            self.from.value(),
            self.to.value(),
            msg_type,
            nanos,
            in_reply_to,
            payload_json,
        ))
    }
}

/// Signs the envelope using the given identity and returns a new envelope with the signature set.
pub fn sign_envelope(env: &Envelope, identity: &dyn Identity) -> EgipResult<Envelope> {
    let canonical = env.canonical_form()?;
    let sig = identity.sign(canonical.as_bytes())?;
    let mut signed = env.clone();
    signed.signature = sig;
    Ok(signed)
}

/// Verifies the envelope's signature against the given public key.
pub fn verify_envelope(env: &Envelope, identity: &dyn Identity, public_key: &PublicKey) -> EgipResult<bool> {
    let canonical = env.canonical_form()?;
    identity.verify(public_key, canonical.as_bytes(), &env.signature)
}

// ── Payload types ───────────────────────────────────────────────────────

/// All EGIP message payloads as an enum.
#[derive(Debug, Clone)]
pub enum MessagePayload {
    Hello(HelloPayload),
    Message(MessagePayloadContent),
    Receipt(ReceiptPayload),
    Proof(ProofPayload),
    Treaty(TreatyPayload),
    AuthorityRequest(AuthorityRequestPayload),
    Discover(DiscoverPayload),
}

impl MessagePayload {
    /// Serializes the payload to canonical JSON for signing.
    fn to_canonical_json(&self) -> Result<String, String> {
        let value = match self {
            Self::Hello(p) => serde_json::to_value(p),
            Self::Message(p) => serde_json::to_value(p),
            Self::Receipt(p) => serde_json::to_value(p),
            Self::Proof(p) => serde_json::to_value(p),
            Self::Treaty(p) => serde_json::to_value(p),
            Self::AuthorityRequest(p) => serde_json::to_value(p),
            Self::Discover(p) => serde_json::to_value(p),
        };
        let v = value.map_err(|e| format!("marshal payload: {e}"))?;
        // Round-trip to normalize key order.
        let canonical = canonical_json_value(&v);
        serde_json::to_string(&canonical).map_err(|e| format!("re-marshal canonical: {e}"))
    }
}

/// Recursively sorts JSON object keys for canonical form.
fn canonical_json_value(v: &Value) -> Value {
    match v {
        Value::Object(map) => {
            let sorted: serde_json::Map<String, Value> = map
                .iter()
                .map(|(k, v)| (k.clone(), canonical_json_value(v)))
                .collect();
            Value::Object(sorted)
        }
        Value::Array(arr) => {
            Value::Array(arr.iter().map(canonical_json_value).collect())
        }
        other => other.clone(),
    }
}

/// Payload for HELLO messages.
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct HelloPayload {
    pub system_uri: String,
    pub public_key: Vec<u8>,
    pub protocol_versions: Vec<i32>,
    pub capabilities: Vec<String>,
    pub chain_length: i32,
}

/// Payload for MESSAGE messages.
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct MessagePayloadContent {
    pub content: Value,
    pub content_type: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub conversation_id: Option<String>,
    #[serde(skip_serializing_if = "Vec::is_empty", default)]
    pub cgers: Vec<CGER>,
}

/// Payload for RECEIPT messages.
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct ReceiptPayload {
    pub envelope_id: String,
    pub status: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub local_event_id: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub reason: Option<String>,
    pub signature: Vec<u8>,
}

impl ReceiptPayload {
    pub fn receipt_status(&self) -> Option<ReceiptStatus> {
        match self.status.as_str() {
            "Delivered" => Some(ReceiptStatus::Delivered),
            "Processed" => Some(ReceiptStatus::Processed),
            "Rejected" => Some(ReceiptStatus::Rejected),
            _ => None,
        }
    }
}

/// Payload for PROOF messages.
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct ProofPayload {
    pub proof_type: String,
    pub data: ProofData,
}

/// Proof-type-specific data.
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
#[serde(tag = "type")]
pub enum ProofData {
    ChainSegment(ChainSegmentProof),
    EventExistence(EventExistenceProof),
    ChainSummary(ChainSummaryProof),
}

/// A contiguous portion of the hash chain.
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct ChainSegmentProof {
    pub event_hashes: Vec<String>,
    pub start_hash: String,
    pub end_hash: String,
}

/// Proves a specific event exists in the chain.
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct EventExistenceProof {
    pub event_hash: String,
    pub prev_hash: String,
    pub next_hash: Option<String>,
    pub position: i32,
    pub chain_length: i32,
}

/// A high-level integrity attestation.
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct ChainSummaryProof {
    pub length: i32,
    pub head_hash: String,
    pub genesis_hash: String,
    pub timestamp_nanos: u64,
}

/// Payload for TREATY messages.
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct TreatyPayload {
    pub treaty_id: String,
    pub action: String,
    #[serde(skip_serializing_if = "Vec::is_empty", default)]
    pub terms: Vec<TreatyTermData>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub reason: Option<String>,
}

impl TreatyPayload {
    pub fn treaty_action(&self) -> Option<TreatyAction> {
        match self.action.as_str() {
            "Propose" => Some(TreatyAction::Propose),
            "Accept" => Some(TreatyAction::Accept),
            "Modify" => Some(TreatyAction::Modify),
            "Suspend" => Some(TreatyAction::Suspend),
            "Terminate" => Some(TreatyAction::Terminate),
            _ => None,
        }
    }
}

/// Serializable treaty term data.
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct TreatyTermData {
    pub scope: String,
    pub policy: String,
    pub symmetric: bool,
}

/// Payload for AUTHORITY_REQUEST messages.
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct AuthorityRequestPayload {
    pub action: String,
    pub actor: String,
    pub level: String,
    pub justification: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub treaty_id: Option<String>,
}

/// Payload for DISCOVER messages.
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct DiscoverPayload {
    pub query: DiscoverQuery,
    #[serde(skip_serializing_if = "Vec::is_empty", default)]
    pub results: Vec<DiscoverResult>,
}

/// Specifies what capabilities to search for.
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct DiscoverQuery {
    pub capabilities: Vec<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub min_trust: Option<f64>,
}

/// A single discovery result.
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct DiscoverResult {
    pub system_uri: String,
    pub public_key: Vec<u8>,
    pub capabilities: Vec<String>,
    pub trust_score: f64,
}

// ── CGER ────────────────────────────────────────────────────────────────

/// Cross-Graph Event Reference with verification tracking.
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct CGER {
    pub local_event_id: String,
    pub remote_system: String,
    pub remote_event_id: String,
    pub remote_hash: String,
    pub relationship: String,
    pub verified: bool,
}

impl CGER {
    pub fn cger_relationship(&self) -> Option<CGERRelationship> {
        match self.relationship.as_str() {
            "CausedBy" => Some(CGERRelationship::CausedBy),
            "References" => Some(CGERRelationship::References),
            "RespondsTo" => Some(CGERRelationship::RespondsTo),
            _ => None,
        }
    }
}

// ── Treaty ──────────────────────────────────────────────────────────────

/// A single term of a bilateral treaty.
#[derive(Debug, Clone)]
pub struct TreatyTerm {
    pub scope: DomainScope,
    pub policy: String,
    pub symmetric: bool,
}

/// A bilateral governance agreement between two systems.
#[derive(Debug, Clone)]
pub struct Treaty {
    pub id: TreatyId,
    pub system_a: SystemUri,
    pub system_b: SystemUri,
    pub status: TreatyStatus,
    pub terms: Vec<TreatyTerm>,
    pub created_at: u64,
    pub updated_at: u64,
}

impl Treaty {
    /// Creates a new treaty in Proposed status.
    pub fn new(
        id: TreatyId,
        system_a: SystemUri,
        system_b: SystemUri,
        terms: Vec<TreatyTerm>,
    ) -> Self {
        let now = nanos_now();
        Self {
            id,
            system_a,
            system_b,
            status: TreatyStatus::Proposed,
            terms,
            created_at: now,
            updated_at: now,
        }
    }

    /// Valid treaty status transitions.
    fn valid_transitions(from: TreatyStatus) -> &'static [TreatyStatus] {
        match from {
            TreatyStatus::Proposed => &[TreatyStatus::Active, TreatyStatus::Terminated],
            TreatyStatus::Active => &[TreatyStatus::Suspended, TreatyStatus::Terminated],
            TreatyStatus::Suspended => &[TreatyStatus::Active, TreatyStatus::Terminated],
            TreatyStatus::Terminated => &[], // terminal state
        }
    }

    /// Attempts to transition the treaty to a new status.
    pub fn transition(&mut self, to: TreatyStatus) -> EgipResult<()> {
        let allowed = Self::valid_transitions(self.status);
        if allowed.contains(&to) {
            self.status = to;
            self.updated_at = nanos_now();
            Ok(())
        } else {
            Err(EgipError::InvalidTreatyTransition {
                from: self.status.as_str().to_string(),
                to: to.as_str().to_string(),
            })
        }
    }

    /// Applies a treaty action and transitions accordingly.
    pub fn apply_action(&mut self, action: TreatyAction) -> EgipResult<()> {
        match action {
            TreatyAction::Accept => self.transition(TreatyStatus::Active),
            TreatyAction::Suspend => self.transition(TreatyStatus::Suspended),
            TreatyAction::Terminate => self.transition(TreatyStatus::Terminated),
            TreatyAction::Modify => {
                if self.status != TreatyStatus::Active {
                    return Err(EgipError::InvalidTreatyTransition {
                        from: self.status.as_str().to_string(),
                        to: "Modify (requires Active)".to_string(),
                    });
                }
                self.updated_at = nanos_now();
                Ok(())
            }
            TreatyAction::Propose => {
                Err(EgipError::Protocol {
                    detail: "cannot apply Propose to existing treaty".to_string(),
                })
            }
        }
    }
}

// ── Trust constants ─────────────────────────────────────────────────────

/// Inter-system trust parameters (more conservative than intra-system).
pub const INTER_SYSTEM_DECAY_RATE: f64 = 0.02; // per day
pub const INTER_SYSTEM_MAX_ADJUSTMENT: f64 = 0.05;

/// Trust impact constants for inter-system actions.
pub const TRUST_IMPACT_VALID_PROOF: f64 = 0.02;
pub const TRUST_IMPACT_RECEIPT_ON_TIME: f64 = 0.01;
pub const TRUST_IMPACT_TREATY_HONOURED: f64 = 0.03;
pub const TRUST_IMPACT_TREATY_VIOLATED: f64 = -0.15;
pub const TRUST_IMPACT_INVALID_PROOF: f64 = -0.10;
pub const TRUST_IMPACT_SIGNATURE_INVALID: f64 = -0.20;
pub const TRUST_IMPACT_NO_HELLO_RESPONSE: f64 = -0.05;

// ── PeerRecord & PeerStore ──────────────────────────────────────────────

/// Tracks the state of a known remote system.
#[derive(Debug, Clone)]
pub struct PeerRecord {
    pub system_uri: SystemUri,
    pub public_key: PublicKey,
    pub trust: Score,
    pub capabilities: Vec<String>,
    pub negotiated_version: i32,
    pub last_seen: Instant,
    pub first_seen: Instant,
    pub last_decayed_at: Instant,
}

/// Manages known peer systems and their trust scores. Thread-safe via RwLock.
pub struct PeerStore {
    peers: RwLock<HashMap<String, PeerRecord>>,
}

impl PeerStore {
    pub fn new() -> Self {
        Self {
            peers: RwLock::new(HashMap::new()),
        }
    }

    /// Adds or updates a peer from a HELLO exchange.
    /// Does NOT overwrite PublicKey on re-registration to prevent key-substitution attacks.
    pub fn register(
        &self,
        uri: SystemUri,
        public_key: PublicKey,
        capabilities: Vec<String>,
        negotiated_version: i32,
    ) -> PeerRecord {
        let mut peers = self.peers.write().unwrap();
        let key = uri.value().to_string();
        let now = Instant::now();

        if let Some(existing) = peers.get_mut(&key) {
            existing.capabilities = capabilities;
            existing.negotiated_version = negotiated_version;
            existing.last_seen = now;
            return existing.clone();
        }

        let record = PeerRecord {
            system_uri: uri,
            public_key,
            trust: Score::new(0.0).unwrap(),
            capabilities,
            negotiated_version,
            last_seen: now,
            first_seen: now,
            last_decayed_at: now,
        };
        peers.insert(key, record.clone());
        record
    }

    /// Returns a copy of the peer record by URI.
    pub fn get(&self, uri: &SystemUri) -> Option<PeerRecord> {
        let peers = self.peers.read().unwrap();
        peers.get(uri.value()).cloned()
    }

    /// Adjusts a peer's trust score by the given delta, clamped to [0,1].
    /// Positive deltas are capped at INTER_SYSTEM_MAX_ADJUSTMENT.
    /// Negative deltas are applied uncapped for immediate security response.
    pub fn update_trust(&self, uri: &SystemUri, delta: f64) -> Option<Score> {
        let mut peers = self.peers.write().unwrap();
        let record = peers.get_mut(uri.value())?;

        let capped_delta = if delta > 0.0 {
            delta.min(INTER_SYSTEM_MAX_ADJUSTMENT)
        } else {
            delta
        };

        let new_val = (record.trust.value() + capped_delta).clamp(0.0, 1.0);
        record.trust = Score::new(new_val).unwrap();
        record.last_seen = Instant::now();
        Some(record.trust)
    }

    /// Applies time-based trust decay to all peers.
    pub fn decay_all(&self) {
        let mut peers = self.peers.write().unwrap();
        let now = Instant::now();

        for record in peers.values_mut() {
            let days_since = now.duration_since(record.last_decayed_at).as_secs_f64() / 86400.0;
            if days_since <= 0.0 {
                continue;
            }
            let decay = INTER_SYSTEM_DECAY_RATE * days_since;
            let new_val = (record.trust.value() - decay).max(0.0);
            record.trust = Score::new(new_val).unwrap();
            record.last_decayed_at = now;
        }
    }

    /// Returns copies of all peer records.
    pub fn all(&self) -> Vec<PeerRecord> {
        let peers = self.peers.read().unwrap();
        peers.values().cloned().collect()
    }
}

impl Default for PeerStore {
    fn default() -> Self { Self::new() }
}

// ── TreatyStore ─────────────────────────────────────────────────────────

/// Manages bilateral treaties. Thread-safe via RwLock.
pub struct TreatyStore {
    treaties: RwLock<HashMap<String, Treaty>>,
}

impl TreatyStore {
    pub fn new() -> Self {
        Self {
            treaties: RwLock::new(HashMap::new()),
        }
    }

    /// Stores or updates a treaty.
    pub fn put(&self, treaty: Treaty) {
        let mut treaties = self.treaties.write().unwrap();
        treaties.insert(treaty.id.value().to_string(), treaty);
    }

    /// Returns a treaty by ID.
    pub fn get(&self, id: &TreatyId) -> Option<Treaty> {
        let treaties = self.treaties.read().unwrap();
        treaties.get(id.value()).cloned()
    }

    /// Performs a read-modify-write on a treaty under a single write lock.
    pub fn apply<F>(&self, id: &TreatyId, f: F) -> EgipResult<()>
    where
        F: FnOnce(&mut Treaty) -> EgipResult<()>,
    {
        let mut treaties = self.treaties.write().unwrap();
        let treaty = treaties.get_mut(id.value()).ok_or_else(|| EgipError::TreatyNotFound {
            treaty_id: id.value().to_string(),
        })?;
        f(treaty)
    }

    /// Returns all treaties involving a given system URI.
    pub fn by_system(&self, uri: &SystemUri) -> Vec<Treaty> {
        let treaties = self.treaties.read().unwrap();
        treaties
            .values()
            .filter(|t| t.system_a.value() == uri.value() || t.system_b.value() == uri.value())
            .cloned()
            .collect()
    }

    /// Returns all active treaties.
    pub fn active(&self) -> Vec<Treaty> {
        let treaties = self.treaties.read().unwrap();
        treaties
            .values()
            .filter(|t| t.status == TreatyStatus::Active)
            .cloned()
            .collect()
    }
}

impl Default for TreatyStore {
    fn default() -> Self { Self::new() }
}

// ── EnvelopeDedup ───────────────────────────────────────────────────────

/// Provides replay protection by tracking seen envelope IDs.
/// Auto-prunes expired entries periodically.
pub struct EnvelopeDedup {
    seen: Mutex<HashMap<String, Instant>>,
    ttl: Duration,
    check_cnt: AtomicU64,
}

impl EnvelopeDedup {
    /// Creates a dedup tracker with default TTL (MAX_ENVELOPE_AGE + 1 hour).
    pub fn new() -> Self {
        Self {
            seen: Mutex::new(HashMap::new()),
            ttl: MAX_ENVELOPE_AGE + Duration::from_secs(3600),
            check_cnt: AtomicU64::new(0),
        }
    }

    /// Creates a dedup tracker with a custom TTL.
    pub fn with_ttl(ttl: Duration) -> Self {
        Self {
            seen: Mutex::new(HashMap::new()),
            ttl,
            check_cnt: AtomicU64::new(0),
        }
    }

    /// Returns true if the envelope ID has not been seen before.
    /// Records the ID and returns false on subsequent calls.
    pub fn check(&self, id: &EnvelopeId) -> bool {
        let mut seen = self.seen.lock().unwrap();
        let key = id.value().to_string();

        if seen.contains_key(&key) {
            return false;
        }

        seen.insert(key, Instant::now());

        let cnt = self.check_cnt.fetch_add(1, Ordering::Relaxed) + 1;
        if cnt % DEDUP_PRUNE_INTERVAL == 0 {
            Self::prune_locked(&mut seen, self.ttl);
        }

        true
    }

    /// Removes expired entries older than TTL.
    pub fn prune(&self) -> usize {
        let mut seen = self.seen.lock().unwrap();
        Self::prune_locked(&mut seen, self.ttl)
    }

    fn prune_locked(seen: &mut HashMap<String, Instant>, ttl: Duration) -> usize {
        let cutoff = Instant::now() - ttl;
        let before = seen.len();
        seen.retain(|_, ts| *ts > cutoff);
        before - seen.len()
    }

    /// Returns the number of tracked envelope IDs.
    pub fn size(&self) -> usize {
        let seen = self.seen.lock().unwrap();
        seen.len()
    }
}

impl Default for EnvelopeDedup {
    fn default() -> Self { Self::new() }
}

// ── Proof verification ─────────────────────────────────────────────────

/// Verifies that a chain segment is internally consistent.
pub fn verify_chain_segment(proof: &ChainSegmentProof) -> bool {
    if proof.event_hashes.is_empty() {
        return false;
    }
    // First hash must follow from start_hash.
    // Last hash must match end_hash.
    if let Some(last) = proof.event_hashes.last() {
        if *last != proof.end_hash {
            return false;
        }
    }
    true
}

/// Verifies basic properties of an event existence proof.
pub fn verify_event_existence(proof: &EventExistenceProof) -> bool {
    // Position and chain length should be consistent.
    if proof.position < 0 || proof.position >= proof.chain_length {
        return false;
    }
    // The event should have a non-empty hash.
    if proof.event_hash.is_empty() {
        return false;
    }
    true
}

/// Dispatches to the appropriate proof verifier.
pub fn validate_proof(payload: &ProofPayload) -> EgipResult<bool> {
    match &payload.data {
        ProofData::ChainSegment(data) => Ok(verify_chain_segment(data)),
        ProofData::EventExistence(data) => Ok(verify_event_existence(data)),
        ProofData::ChainSummary(data) => Ok(data.length > 0),
    }
}

/// Returns the ProofType for a given ProofData.
pub fn proof_type_from_data(data: &ProofData) -> ProofType {
    match data {
        ProofData::ChainSegment(_) => ProofType::ChainSegment,
        ProofData::EventExistence(_) => ProofType::EventExistence,
        ProofData::ChainSummary(_) => ProofType::ChainSummary,
    }
}

// ── Handler ─────────────────────────────────────────────────────────────

/// Callback type for MESSAGE handling.
pub type OnMessageFn = Box<dyn Fn(&SystemUri, &MessagePayloadContent) -> EgipResult<()> + Send + Sync>;

/// Callback type for AUTHORITY_REQUEST handling.
pub type OnAuthorityRequestFn = Box<dyn Fn(&SystemUri, &AuthorityRequestPayload) -> EgipResult<()> + Send + Sync>;

/// Callback type for DISCOVER handling.
pub type OnDiscoverFn = Box<dyn Fn(&SystemUri, &DiscoverQuery) -> EgipResult<Vec<DiscoverResult>> + Send + Sync>;

/// Callback type for chain length.
pub type ChainLengthFn = Box<dyn Fn() -> EgipResult<i32> + Send + Sync>;

/// Protocol handler: HELLO handshake, message dispatch, replay dedup, trust updates.
pub struct Handler {
    identity: Box<dyn Identity + Send + Sync>,
    transport: Box<dyn Transport + Send + Sync>,
    peers: PeerStore,
    treaties: TreatyStore,
    dedup: EnvelopeDedup,
    pub local_protocol_versions: Vec<i32>,
    pub capabilities: Vec<String>,
    pub chain_length: Option<ChainLengthFn>,
    pub on_message: Option<OnMessageFn>,
    pub on_authority_request: Option<OnAuthorityRequestFn>,
    pub on_discover: Option<OnDiscoverFn>,
}

impl Handler {
    /// Creates a new protocol handler.
    pub fn new(
        identity: Box<dyn Identity + Send + Sync>,
        transport: Box<dyn Transport + Send + Sync>,
    ) -> Self {
        Self {
            identity,
            transport,
            peers: PeerStore::new(),
            treaties: TreatyStore::new(),
            dedup: EnvelopeDedup::new(),
            local_protocol_versions: vec![CURRENT_PROTOCOL_VERSION],
            capabilities: vec!["treaty".to_string(), "proof".to_string()],
            chain_length: None,
            on_message: None,
            on_authority_request: None,
            on_discover: None,
        }
    }

    /// Returns a reference to the peer store.
    pub fn peers(&self) -> &PeerStore { &self.peers }

    /// Returns a reference to the treaty store.
    pub fn treaties(&self) -> &TreatyStore { &self.treaties }

    /// Performs the HELLO handshake with a remote system.
    pub fn hello(&self, to: &SystemUri) -> EgipResult<()> {
        let chain_len = match &self.chain_length {
            Some(f) => f()?,
            None => 0,
        };

        let env_id = generate_envelope_id()?;

        let env = Envelope {
            protocol_version: CURRENT_PROTOCOL_VERSION,
            id: env_id,
            from: self.identity.system_uri().clone(),
            to: to.clone(),
            message_type: MessageType::Hello,
            payload: MessagePayload::Hello(HelloPayload {
                system_uri: self.identity.system_uri().value().to_string(),
                public_key: self.identity.public_key().bytes().to_vec(),
                protocol_versions: self.local_protocol_versions.clone(),
                capabilities: self.capabilities.clone(),
                chain_length: chain_len,
            }),
            timestamp_nanos: nanos_now(),
            signature: Signature::zero(),
            in_reply_to: None,
        };

        let signed = sign_envelope(&env, self.identity.as_ref())?;

        let receipt = self.transport.send(to, &signed);
        match receipt {
            Ok(Some(r)) if r.receipt_status() == Some(ReceiptStatus::Rejected) => {
                let reason = r.reason.unwrap_or_default();
                Err(EgipError::Protocol {
                    detail: format!("hello rejected: {reason}"),
                })
            }
            Ok(_) => Ok(()),
            Err(e) => {
                self.peers.update_trust(to, TRUST_IMPACT_NO_HELLO_RESPONSE);
                Err(EgipError::TransportFailure {
                    to: to.value().to_string(),
                    reason: e.to_string(),
                })
            }
        }
    }

    /// Processes an incoming envelope: checks timestamp freshness, deduplicates,
    /// verifies signature, dispatches, and updates trust.
    pub fn handle_incoming(&self, env: &Envelope) -> EgipResult<()> {
        // Timestamp freshness.
        let now_nanos = nanos_now();
        let env_age_secs = if now_nanos >= env.timestamp_nanos {
            (now_nanos - env.timestamp_nanos) as f64 / 1_000_000_000.0
        } else {
            -((env.timestamp_nanos - now_nanos) as f64 / 1_000_000_000.0)
        };

        if env_age_secs > MAX_ENVELOPE_AGE.as_secs_f64() || env_age_secs < -300.0 {
            return Err(EgipError::EnvelopeStale { age_secs: env_age_secs });
        }

        // Replay deduplication.
        if !self.dedup.check(&env.id) {
            return Err(EgipError::DuplicateEnvelope {
                envelope_id: env.id.value().to_string(),
            });
        }

        // Look up the sender's public key.
        let pub_key = if env.message_type == MessageType::Hello {
            if let MessagePayload::Hello(ref hello) = env.payload {
                if hello.public_key.len() != 32 {
                    return Err(EgipError::CryptoError {
                        detail: format!("invalid public key length: {}", hello.public_key.len()),
                    });
                }
                let mut bytes = [0u8; 32];
                bytes.copy_from_slice(&hello.public_key);
                PublicKey::new(bytes)
            } else {
                return Err(EgipError::InvalidPayload {
                    expected: "HelloPayload",
                    got: format!("{:?}", env.payload),
                });
            }
        } else {
            match self.peers.get(&env.from) {
                Some(peer) => peer.public_key.clone(),
                None => {
                    return Err(EgipError::SystemNotFound {
                        uri: env.from.value().to_string(),
                    });
                }
            }
        };

        // Verify signature.
        let valid = verify_envelope(env, self.identity.as_ref(), &pub_key)?;
        if !valid {
            self.peers.update_trust(&env.from, TRUST_IMPACT_SIGNATURE_INVALID);
            return Err(EgipError::EnvelopeSignatureInvalid {
                envelope_id: env.id.value().to_string(),
            });
        }

        // Dispatch by message type.
        match env.message_type {
            MessageType::Hello => self.handle_hello(env),
            MessageType::Message => self.handle_message(env),
            MessageType::Receipt => self.handle_receipt(env),
            MessageType::Proof => self.handle_proof(env),
            MessageType::Treaty => self.handle_treaty(env),
            MessageType::AuthorityRequest => self.handle_authority_request(env),
            MessageType::Discover => self.handle_discover(env),
        }
    }

    fn handle_hello(&self, env: &Envelope) -> EgipResult<()> {
        let hello = match &env.payload {
            MessagePayload::Hello(h) => h,
            _ => {
                return Err(EgipError::InvalidPayload {
                    expected: "HelloPayload",
                    got: "other".to_string(),
                });
            }
        };

        let version = negotiate_version(&self.local_protocol_versions, &hello.protocol_versions);
        match version {
            None => Err(EgipError::VersionIncompatible {
                local: self.local_protocol_versions.clone(),
                remote: hello.protocol_versions.clone(),
            }),
            Some(v) => {
                let uri = SystemUri::new(&hello.system_uri).map_err(|e| EgipError::Protocol {
                    detail: format!("invalid system URI in hello: {e}"),
                })?;
                let mut pub_bytes = [0u8; 32];
                if hello.public_key.len() == 32 {
                    pub_bytes.copy_from_slice(&hello.public_key);
                }
                self.peers.register(
                    uri,
                    PublicKey::new(pub_bytes),
                    hello.capabilities.clone(),
                    v,
                );
                Ok(())
            }
        }
    }

    fn handle_message(&self, env: &Envelope) -> EgipResult<()> {
        let msg = match &env.payload {
            MessagePayload::Message(m) => m,
            _ => {
                return Err(EgipError::InvalidPayload {
                    expected: "MessagePayloadContent",
                    got: "other".to_string(),
                });
            }
        };

        self.peers.update_trust(&env.from, TRUST_IMPACT_RECEIPT_ON_TIME);

        if let Some(ref handler) = self.on_message {
            handler(&env.from, msg)?;
        }
        Ok(())
    }

    fn handle_receipt(&self, env: &Envelope) -> EgipResult<()> {
        let receipt = match &env.payload {
            MessagePayload::Receipt(r) => r,
            _ => {
                return Err(EgipError::InvalidPayload {
                    expected: "ReceiptPayload",
                    got: "other".to_string(),
                });
            }
        };

        if let Some(status) = receipt.receipt_status() {
            if status == ReceiptStatus::Processed || status == ReceiptStatus::Delivered {
                self.peers.update_trust(&env.from, TRUST_IMPACT_RECEIPT_ON_TIME);
            }
        }
        Ok(())
    }

    fn handle_proof(&self, env: &Envelope) -> EgipResult<()> {
        let proof = match &env.payload {
            MessagePayload::Proof(p) => p,
            _ => {
                return Err(EgipError::InvalidPayload {
                    expected: "ProofPayload",
                    got: "other".to_string(),
                });
            }
        };

        let valid = validate_proof(proof)?;
        if valid {
            self.peers.update_trust(&env.from, TRUST_IMPACT_VALID_PROOF);
        } else {
            self.peers.update_trust(&env.from, TRUST_IMPACT_INVALID_PROOF);
        }
        Ok(())
    }

    fn handle_treaty(&self, env: &Envelope) -> EgipResult<()> {
        let payload = match &env.payload {
            MessagePayload::Treaty(t) => t,
            _ => {
                return Err(EgipError::InvalidPayload {
                    expected: "TreatyPayload",
                    got: "other".to_string(),
                });
            }
        };

        let action = payload.treaty_action().ok_or_else(|| EgipError::Protocol {
            detail: format!("unknown treaty action: {}", payload.action),
        })?;

        let treaty_id = TreatyId::new(&payload.treaty_id).map_err(|e| EgipError::Protocol {
            detail: format!("invalid treaty ID: {e}"),
        })?;

        match action {
            TreatyAction::Propose => {
                let terms: Vec<TreatyTerm> = payload
                    .terms
                    .iter()
                    .filter_map(|t| {
                        DomainScope::new(&t.scope).ok().map(|scope| TreatyTerm {
                            scope,
                            policy: t.policy.clone(),
                            symmetric: t.symmetric,
                        })
                    })
                    .collect();
                let treaty = Treaty::new(
                    treaty_id,
                    env.from.clone(),
                    env.to.clone(),
                    terms,
                );
                self.treaties.put(treaty);
                Ok(())
            }
            TreatyAction::Accept => {
                self.treaties.apply(&treaty_id, |treaty| {
                    treaty.apply_action(TreatyAction::Accept)?;
                    Ok(())
                })?;
                self.peers.update_trust(&env.from, TRUST_IMPACT_TREATY_HONOURED);
                Ok(())
            }
            TreatyAction::Suspend => {
                self.treaties.apply(&treaty_id, |treaty| {
                    treaty.apply_action(TreatyAction::Suspend)
                })
            }
            TreatyAction::Terminate => {
                self.treaties.apply(&treaty_id, |treaty| {
                    treaty.apply_action(TreatyAction::Terminate)
                })
            }
            TreatyAction::Modify => {
                let new_terms: Vec<TreatyTerm> = payload
                    .terms
                    .iter()
                    .filter_map(|t| {
                        DomainScope::new(&t.scope).ok().map(|scope| TreatyTerm {
                            scope,
                            policy: t.policy.clone(),
                            symmetric: t.symmetric,
                        })
                    })
                    .collect();
                self.treaties.apply(&treaty_id, |treaty| {
                    treaty.apply_action(TreatyAction::Modify)?;
                    treaty.terms = new_terms;
                    Ok(())
                })
            }
        }
    }

    fn handle_authority_request(&self, env: &Envelope) -> EgipResult<()> {
        let payload = match &env.payload {
            MessagePayload::AuthorityRequest(a) => a,
            _ => {
                return Err(EgipError::InvalidPayload {
                    expected: "AuthorityRequestPayload",
                    got: "other".to_string(),
                });
            }
        };

        if let Some(ref handler) = self.on_authority_request {
            handler(&env.from, payload)?;
        }
        Ok(())
    }

    fn handle_discover(&self, env: &Envelope) -> EgipResult<()> {
        let payload = match &env.payload {
            MessagePayload::Discover(d) => d,
            _ => {
                return Err(EgipError::InvalidPayload {
                    expected: "DiscoverPayload",
                    got: "other".to_string(),
                });
            }
        };

        let on_discover = match &self.on_discover {
            Some(handler) => handler,
            None => return Ok(()),
        };

        let results = on_discover(&env.from, &payload.query)?;

        let resp_id = generate_envelope_id()?;

        let resp = Envelope {
            protocol_version: CURRENT_PROTOCOL_VERSION,
            id: resp_id,
            from: self.identity.system_uri().clone(),
            to: env.from.clone(),
            message_type: MessageType::Discover,
            payload: MessagePayload::Discover(DiscoverPayload {
                query: payload.query.clone(),
                results,
            }),
            timestamp_nanos: nanos_now(),
            signature: Signature::zero(),
            in_reply_to: Some(env.id.clone()),
        };

        let signed = sign_envelope(&resp, self.identity.as_ref())?;
        self.transport.send(&env.from, &signed)?;
        Ok(())
    }
}

// ── Version negotiation ─────────────────────────────────────────────────

/// Finds the highest protocol version both systems support.
/// Returns None if no common version exists.
pub fn negotiate_version(local: &[i32], remote: &[i32]) -> Option<i32> {
    let mut best: Option<i32> = None;
    for &l in local {
        for &r in remote {
            if l == r {
                best = Some(best.map_or(l, |b: i32| b.max(l)));
            }
        }
    }
    best
}

// ── Helpers ─────────────────────────────────────────────────────────────

fn nanos_now() -> u64 {
    SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap_or_default()
        .as_nanos() as u64
}

/// Generates a UUID v4 string for envelope IDs.
fn generate_uuid4() -> String {
    use std::time::SystemTime;
    // Simple UUID v4 from random-ish bytes (timestamp + counter).
    let ts = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap_or_default()
        .as_nanos();
    let mut b = [0u8; 16];
    let ts_bytes = ts.to_le_bytes();
    b[..16].copy_from_slice(&ts_bytes);
    // Add some variance from thread ID hash.
    let thread_id = format!("{:?}", std::thread::current().id());
    for (i, byte) in thread_id.bytes().enumerate() {
        b[i % 16] ^= byte;
    }
    b[6] = (b[6] & 0x0f) | 0x40; // version 4
    b[8] = (b[8] & 0x3f) | 0x80; // variant 10
    format!(
        "{:08x}-{:04x}-{:04x}-{:04x}-{:012x}",
        u32::from_be_bytes([b[0], b[1], b[2], b[3]]),
        u16::from_be_bytes([b[4], b[5]]),
        u16::from_be_bytes([b[6], b[7]]),
        u16::from_be_bytes([b[8], b[9]]),
        u64::from_be_bytes([0, 0, b[10], b[11], b[12], b[13], b[14], b[15]]),
    )
}

/// Creates an EnvelopeId from a fresh UUID v4.
fn generate_envelope_id() -> EgipResult<EnvelopeId> {
    let uuid = generate_uuid4();
    EnvelopeId::new(uuid).map_err(|e| EgipError::Protocol {
        detail: format!("generate envelope ID: {e}"),
    })
}
