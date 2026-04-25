package event

import (
	"time"

	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// siteContent is embedded in all site content types to satisfy the
// EventContent interface's Accept method. Site content types use their
// own visitor rather than the base EventContentVisitor.
type siteContent struct{}

func (siteContent) Accept(EventContentVisitor) {}

// SiteOpReceivedContent is the anchor event — emitted synchronously
// when the hive's webhook handler receives a site op.
// The raw payload is NOT stored on the chain; only its SHA-256 hash.
type SiteOpReceivedContent struct {
	siteContent
	ExternalRef   ExternalRef `json:"external_ref"`
	SpaceID       string      `json:"space_id"`
	NodeID        string      `json:"node_id,omitempty"`
	NodeTitle     string      `json:"node_title,omitempty"`
	Actor         string      `json:"actor"`
	ActorID       string      `json:"actor_id"`
	ActorKind     string      `json:"actor_kind"`
	OpKind        string      `json:"op_kind"`
	PayloadHash   string      `json:"payload_hash"`
	ReceivedAt    time.Time   `json:"received_at"`
	SiteCreatedAt time.Time   `json:"site_created_at"`
}

func (c SiteOpReceivedContent) EventTypeName() string { return "site.op.received" }

// SiteOpTranslatedContent is emitted after the translation goroutine
// successfully converts an anchored op into a hive bus event.
type SiteOpTranslatedContent struct {
	siteContent
	ExternalRef  ExternalRef   `json:"external_ref"`
	BusEventID   types.EventID `json:"bus_event_id"`
	TranslatedAt time.Time     `json:"translated_at"`
}

func (c SiteOpTranslatedContent) EventTypeName() string { return "site.op.translated" }

// SiteOpRejectedContent is emitted when translation fails.
type SiteOpRejectedContent struct {
	siteContent
	ExternalRef ExternalRef `json:"external_ref"`
	Reason      string      `json:"reason"`
	RejectedAt  time.Time   `json:"rejected_at"`
}

func (c SiteOpRejectedContent) EventTypeName() string { return "site.op.rejected" }

// SiteOpMirroredContent is emitted after POST /api/hive/mirror returns 2xx.
// MirrorEventID is the hive event whose update was mirrored to site.
type SiteOpMirroredContent struct {
	siteContent
	ExternalRef   ExternalRef   `json:"external_ref"`
	MirrorEventID types.EventID `json:"mirror_event_id"`
	HTTPStatus    int           `json:"http_status"`
	MirroredAt    time.Time     `json:"mirrored_at"`
}

func (c SiteOpMirroredContent) EventTypeName() string { return "site.op.mirrored" }
