package egip

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// IIdentity represents a system's cryptographic identity.
type IIdentity interface {
	// SystemURI returns this system's address.
	SystemURI() types.SystemURI

	// PublicKey returns the Ed25519 public key (32 bytes).
	PublicKey() types.PublicKey

	// Sign produces an Ed25519 signature of the given data.
	Sign(data []byte) (types.Signature, error)

	// Verify checks an Ed25519 signature against a public key.
	Verify(publicKey types.PublicKey, data []byte, signature types.Signature) (bool, error)
}

// SystemIdentity is the Ed25519-based implementation of IIdentity.
type SystemIdentity struct {
	uri        types.SystemURI
	publicKey  types.PublicKey
	privateKey ed25519.PrivateKey
	createdAt  time.Time
}

// GenerateIdentity creates a new system identity with a fresh Ed25519 keypair.
func GenerateIdentity(uri types.SystemURI) (*SystemIdentity, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate keypair: %w", err)
	}

	publicKey, err := types.NewPublicKey([]byte(pub))
	if err != nil {
		return nil, fmt.Errorf("create public key: %w", err)
	}

	return &SystemIdentity{
		uri:        uri,
		publicKey:  publicKey,
		privateKey: priv,
		createdAt:  time.Now(),
	}, nil
}

// NewIdentityFromKey creates a system identity from an existing Ed25519 private key.
func NewIdentityFromKey(uri types.SystemURI, privateKey ed25519.PrivateKey) (*SystemIdentity, error) {
	pub := privateKey.Public().(ed25519.PublicKey)
	publicKey, err := types.NewPublicKey([]byte(pub))
	if err != nil {
		return nil, fmt.Errorf("create public key: %w", err)
	}

	return &SystemIdentity{
		uri:        uri,
		publicKey:  publicKey,
		privateKey: privateKey,
		createdAt:  time.Now(),
	}, nil
}

func (id *SystemIdentity) SystemURI() types.SystemURI { return id.uri }
func (id *SystemIdentity) PublicKey() types.PublicKey  { return id.publicKey }
func (id *SystemIdentity) CreatedAt() time.Time        { return id.createdAt }

// Sign produces an Ed25519 signature of the given data.
func (id *SystemIdentity) Sign(data []byte) (types.Signature, error) {
	sig := ed25519.Sign(id.privateKey, data)
	s, err := types.NewSignature(sig)
	if err != nil {
		return types.Signature{}, fmt.Errorf("create signature: %w", err)
	}
	return s, nil
}

// Verify checks an Ed25519 signature against a public key.
func (id *SystemIdentity) Verify(publicKey types.PublicKey, data []byte, signature types.Signature) (bool, error) {
	pubBytes := publicKey.Bytes()
	if len(pubBytes) != ed25519.PublicKeySize {
		return false, fmt.Errorf("invalid public key length: %d", len(pubBytes))
	}

	sigBytes := signature.Bytes()
	if len(sigBytes) != ed25519.SignatureSize {
		return false, fmt.Errorf("invalid signature length: %d", len(sigBytes))
	}

	return ed25519.Verify(ed25519.PublicKey(pubBytes), data, sigBytes), nil
}
