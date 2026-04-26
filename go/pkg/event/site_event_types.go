package event

import "github.com/transpara-ai/eventgraph/go/pkg/types"

// Site event type constants. All use the "site." prefix.
// Emitted by the hive's bridge layer to anchor, translate, reject,
// and mirror operations between the site and the hive chain.
var (
	EventTypeSiteOpReceived   = types.MustEventType("site.op.received")
	EventTypeSiteOpTranslated = types.MustEventType("site.op.translated")
	EventTypeSiteOpRejected   = types.MustEventType("site.op.rejected")
	EventTypeSiteOpMirrored   = types.MustEventType("site.op.mirrored")
)

// AllSiteEventTypes returns all registered site event types.
func AllSiteEventTypes() []types.EventType {
	return []types.EventType{
		EventTypeSiteOpReceived,
		EventTypeSiteOpTranslated,
		EventTypeSiteOpRejected,
		EventTypeSiteOpMirrored,
	}
}
