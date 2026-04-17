package event

// ExternalRef identifies a record in a non-eventgraph system
// (e.g. site, GitHub, Slack). Used by site.op.* content types to
// anchor hive-chain events to external audit logs.
type ExternalRef struct {
	System string `json:"system"` // "site", "github", "slack", future-proof
	ID     string `json:"id"`
}
