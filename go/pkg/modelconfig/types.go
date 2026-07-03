package modelconfig

// ModelTier classifies models by capability/cost tradeoff.
type ModelTier string

const (
	TierJudgment  ModelTier = "judgment"  // expensive, high-capability (Opus-class)
	TierExecution ModelTier = "execution" // mid-cost, strong workers (Sonnet-class)
	TierVolume    ModelTier = "volume"    // cheap, fast (Haiku-class)
)

// AuthMode describes how a provider authenticates.
type AuthMode string

const (
	AuthSubscription AuthMode = "subscription" // claude-cli, codex-cli (flat rate)
	AuthAPIKey       AuthMode = "api-key"      // anthropic, openai-compatible
	AuthLocal        AuthMode = "local"        // ollama, local models (no auth)
)

// SelectionMode records HOW the winning model token resolved. It is derived
// at the RESOLUTION POINT of the winning token, never inferred from the layer
// that supplied it — a policy token may hold a model id, an alias, or a tier
// name, so "a policy supplied the token" proves nothing about how it resolved.
// Closed vocabulary; the zero value means "not derived" and consumers MUST
// treat it as not-projected, never coerce it to a real mode.
type SelectionMode string

const (
	// SelectionModeManualExplicit — the winning token resolved as a concrete
	// catalog id or alias (a ModelAliases remap of an explicit pin included).
	SelectionModeManualExplicit SelectionMode = "manual-explicit"
	// SelectionModeAutoTier — the winning token fell back through tier
	// resolution (a tier-name token) or PreferredTier selected the model.
	SelectionModeAutoTier SelectionMode = "auto-tier"
	// SelectionModeSystemDefault — no caller-supplied token won: system,
	// role, or tier defaults chose the model.
	SelectionModeSystemDefault SelectionMode = "system-default"
	// SelectionModeUnknown — absent. Consumers MUST treat it as
	// not-projected, never coerce.
	SelectionModeUnknown SelectionMode = ""
)

// Capability describes what a model can do.
type Capability string

const (
	CapTools         Capability = "tools"
	CapReasoning     Capability = "reasoning"
	CapCoding        Capability = "coding"
	CapVision        Capability = "vision"
	CapOperate       Capability = "operate"           // agentic filesystem access
	CapLargeContext  Capability = "large-context"     // >100k context window
	CapFastLatency   Capability = "fast-latency"      // optimized for speed
	CapStructuredOut Capability = "structured-output" // JSON schema output
)
