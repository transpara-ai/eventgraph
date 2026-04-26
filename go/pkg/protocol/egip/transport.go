package egip

import (
	"context"

	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// ITransport is the pluggable transport layer for EGIP communication.
// Implementations may use HTTP, WebSocket, gRPC, message queues, or any other mechanism.
type ITransport interface {
	// Send delivers an envelope to the target system.
	Send(ctx context.Context, to types.SystemURI, envelope *Envelope) (*ReceiptPayload, error)

	// Listen returns a channel of incoming envelopes.
	Listen(ctx context.Context) <-chan IncomingEnvelope
}

// IncomingEnvelope wraps an envelope received from a remote system.
type IncomingEnvelope struct {
	Envelope *Envelope
	Err      error
}
