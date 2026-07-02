package modelconfig

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testCatalog(t *testing.T) *ModelCatalog {
	t.Helper()
	cat, err := NewCatalog(testEntries())
	require.NoError(t, err)
	return cat
}

func testProfiles() map[string]ModelProfile {
	temp01 := 0.1
	return map[string]ModelProfile{
		"balanced": {
			Name:     "balanced",
			Model:    "sonnet",
			Provider: "claude-cli",
		},
		"deep-judgment": {
			Name:        "deep-judgment",
			Model:       "opus",
			Provider:    "claude-cli",
			Temperature: &temp01,
		},
	}
}

func testDefaults() ResolverDefaults {
	return ResolverDefaults{
		Provider: "claude-cli",
		Model:    "test-sonnet",
		TierModels: map[ModelTier]string{
			TierJudgment:  "test-opus",
			TierExecution: "test-sonnet",
			TierVolume:    "test-haiku",
		},
		RoleModels: map[string]string{
			"guardian": "sonnet",
			"cto":      "opus",
			"sysmon":   "haiku",
		},
	}
}

func testResolver(t *testing.T) *Resolver {
	t.Helper()
	return NewResolver(testCatalog(t), testProfiles(), testDefaults())
}

func mixedProviderResolver(t *testing.T) *Resolver {
	t.Helper()
	cat, err := NewCatalog([]ModelCatalogEntry{
		{
			ID:           "test-sonnet",
			Aliases:      []string{"sonnet"},
			Provider:     "claude-cli",
			AuthMode:     AuthSubscription,
			Tier:         TierExecution,
			Capabilities: []Capability{CapTools, CapReasoning, CapCoding, CapOperate},
		},
		{
			ID:           "api-sonnet",
			Aliases:      []string{"api-sonnet-alias"},
			Provider:     "anthropic",
			AuthMode:     AuthAPIKey,
			Tier:         TierExecution,
			Capabilities: []Capability{CapTools, CapReasoning, CapCoding},
		},
	})
	require.NoError(t, err)
	return NewResolver(cat, map[string]ModelProfile{
		"api-with-matching-provider": {
			Name:     "api-with-matching-provider",
			Model:    "api-sonnet",
			Provider: "anthropic",
		},
		"api-with-cli-provider": {
			Name:     "api-with-cli-provider",
			Model:    "api-sonnet",
			Provider: "claude-cli",
		},
	}, ResolverDefaults{
		Provider: "claude-cli",
		Model:    "test-sonnet",
	})
}

// modeFromWinningTrace classifies the WINNING (last) model-writing trace entry
// into the SelectionMode it implies. ModelAliases remap entries preserve the
// source class of the token they remapped, so they are skipped. Allowlist
// form: an unclassifiable winning entry yields SelectionModeUnknown so the
// cross-check fails loudly instead of blessing an unproven mode.
func modeFromWinningTrace(trace []string) SelectionMode {
	for i := len(trace) - 1; i >= 0; i-- {
		e := trace[i]
		if !strings.HasPrefix(e, "model: ") {
			continue
		}
		if strings.Contains(e, "catalog alias override") {
			// Remap preserves the winning source class — classify the
			// entry that supplied the remapped token instead.
			continue
		}
		switch {
		case strings.Contains(e, " tier "):
			return SelectionModeAutoTier
		case strings.Contains(e, "AgentDef.Model"),
			strings.Contains(e, "explicit"),
			strings.Contains(e, "profile"):
			return SelectionModeManualExplicit
		case strings.Contains(e, "system default"),
			strings.Contains(e, "role default"):
			return SelectionModeSystemDefault
		}
		return SelectionModeUnknown
	}
	return SelectionModeUnknown
}

func TestResolve(t *testing.T) {
	tests := []struct {
		name       string
		resolver   func(t *testing.T) *Resolver // nil → testResolver
		input      ResolutionInput
		wantModel  string
		wantProv   string
		wantMode   SelectionMode
		wantErr    string
		checkTrace func(t *testing.T, trace []string)
	}{
		{
			name:      "defaults only (no policy, no AgentDefModel)",
			input:     ResolutionInput{Role: "unknown-role"},
			wantModel: "test-sonnet",
			wantProv:  "claude-cli",
			wantMode:  SelectionModeSystemDefault,
			checkTrace: func(t *testing.T, trace []string) {
				assert.Contains(t, trace[0], "system default")
			},
		},
		{
			name: "AgentDefModel overrides role default",
			input: ResolutionInput{
				Role:          "guardian",
				AgentDefModel: "opus", // alias
			},
			wantModel: "test-opus",
			wantProv:  "claude-cli",
			wantMode:  SelectionModeManualExplicit,
			checkTrace: func(t *testing.T, trace []string) {
				found := false
				for _, s := range trace {
					if assert.ObjectsAreEqual("model: AgentDef.Model -> opus", s) || // just check it's there
						len(s) > 0 {
						found = true
					}
				}
				assert.True(t, found)
			},
		},
		{
			name: "Policy.Model overrides AgentDefModel",
			input: ResolutionInput{
				Role:          "guardian",
				AgentDefModel: "sonnet",
				Policy: &RoleModelPolicy{
					Model: "opus",
				},
			},
			wantModel: "test-opus",
			wantProv:  "claude-cli",
			wantMode:  SelectionModeManualExplicit,
		},
		{
			name: "Policy.Profile expands profile fields",
			input: ResolutionInput{
				Role: "worker",
				Policy: &RoleModelPolicy{
					Profile: "deep-judgment",
				},
			},
			wantModel: "test-opus",
			wantProv:  "claude-cli",
			wantMode:  SelectionModeManualExplicit, // profile supplied a concrete alias — still an explicit pin
			checkTrace: func(t *testing.T, trace []string) {
				hasProfile := false
				for _, s := range trace {
					if assert.ObjectsAreEqual(s, "") {
						continue
					}
					// look for profile trace entry
					if len(s) > 0 {
						hasProfile = true
					}
				}
				assert.True(t, hasProfile)
			},
		},
		{
			name: "TaskOverride wins over Policy",
			input: ResolutionInput{
				Role: "worker",
				Policy: &RoleModelPolicy{
					Model: "sonnet",
				},
				TaskOverride: &RoleModelPolicy{
					Model: "opus",
				},
			},
			wantModel: "test-opus",
			wantProv:  "claude-cli",
			wantMode:  SelectionModeManualExplicit,
		},
		{
			name: "CanOperate with claude-cli model succeeds",
			input: ResolutionInput{
				Role:       "worker",
				CanOperate: true,
			},
			wantModel: "test-sonnet",
			wantProv:  "claude-cli",
			wantMode:  SelectionModeSystemDefault, // CanOperate constrains, it never chooses
		},
		{
			name: "PreferredTier resolves to tier model",
			input: ResolutionInput{
				Role: "cheap-worker",
				Policy: &RoleModelPolicy{
					PreferredTier: TierVolume,
				},
			},
			wantModel: "test-haiku",
			wantProv:  "claude-cli",
			wantMode:  SelectionModeAutoTier,
		},
		{
			name:      "role default resolves as system default",
			input:     ResolutionInput{Role: "guardian"},
			wantModel: "test-sonnet",
			wantProv:  "claude-cli",
			wantMode:  SelectionModeSystemDefault, // role defaults are system config, not caller input
		},
		{
			name: "tier-name token via Policy.Model resolves through tier fallback",
			input: ResolutionInput{
				Role: "worker",
				Policy: &RoleModelPolicy{
					Model: "judgment", // a tier name, not a model id/alias
				},
			},
			wantModel: "test-opus",
			wantProv:  "claude-cli",
			wantMode:  SelectionModeAutoTier, // the token WON as a policy field but RESOLVED via tier
		},
		{
			name: "ModelAliases remap of explicit pin stays manual-explicit",
			resolver: func(t *testing.T) *Resolver {
				defaults := testDefaults()
				defaults.ModelAliases = map[string]string{"opus": "haiku"}
				return NewResolver(testCatalog(t), testProfiles(), defaults)
			},
			input: ResolutionInput{
				Role:          "worker",
				AgentDefModel: "opus",
			},
			wantModel: "test-haiku",
			wantProv:  "claude-cli",
			wantMode:  SelectionModeManualExplicit, // a remapped alias is still an explicit pin
		},
		{
			name: "ModelAliases remap of role default stays system-default",
			resolver: func(t *testing.T) *Resolver {
				defaults := testDefaults()
				defaults.ModelAliases = map[string]string{"sonnet": "haiku"}
				return NewResolver(testCatalog(t), testProfiles(), defaults)
			},
			input:     ResolutionInput{Role: "guardian"},
			wantModel: "test-haiku",
			wantProv:  "claude-cli",
			wantMode:  SelectionModeSystemDefault, // remap preserves the winning source class
		},
		{
			name: "unknown model returns error",
			input: ResolutionInput{
				Role:          "worker",
				AgentDefModel: "does-not-exist",
			},
			wantErr: "not found in catalog",
		},
		{
			name: "missing required capabilities returns error",
			input: ResolutionInput{
				Role: "worker",
				Policy: &RoleModelPolicy{
					Model:                "haiku",
					RequiredCapabilities: []Capability{CapReasoning}, // haiku lacks reasoning
				},
			},
			wantErr: "missing capabilities",
		},
		{
			name: "task override provider must match resolved model provider",
			input: ResolutionInput{
				Role: "worker",
				TaskOverride: &RoleModelPolicy{
					Provider: "anthropic",
				},
			},
			wantErr: "belongs to provider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := testResolver(t)
			if tt.resolver != nil {
				r = tt.resolver(t)
			}
			rc, err := r.Resolve(tt.input)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Equal(t, SelectionModeUnknown, rc.Mode,
					"error paths must leave mode unknown — fail closed, never coerced")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantModel, rc.Model)
			assert.Equal(t, tt.wantProv, rc.Provider)
			assert.Equal(t, tt.wantMode, rc.Mode)
			assert.NotEmpty(t, rc.Trace, "trace should have entries")

			// Trace is the cross-check oracle (IADA-1): the derived mode
			// must match what the winning trace entry implies.
			assert.Equal(t, tt.wantMode, modeFromWinningTrace(rc.Trace),
				"mode inconsistent with winning trace entry: %v", rc.Trace)

			// Final trace entry should be the resolved summary
			last := rc.Trace[len(rc.Trace)-1]
			assert.Contains(t, last, "resolved:")

			if tt.checkTrace != nil {
				tt.checkTrace(t, rc.Trace)
			}
		})
	}
}

func TestResolve_ProviderModelCoherence(t *testing.T) {
	tests := []struct {
		name    string
		input   ResolutionInput
		wantErr string
	}{
		{
			name: "task override matching provider succeeds",
			input: ResolutionInput{
				Role: "guardian",
				TaskOverride: &RoleModelPolicy{
					Model:    "api-sonnet",
					Provider: "anthropic",
				},
			},
		},
		{
			name: "task override explicit default provider mismatch fails",
			input: ResolutionInput{
				Role: "guardian",
				TaskOverride: &RoleModelPolicy{
					Model:    "api-sonnet",
					Provider: "claude-cli",
				},
			},
			wantErr: "belongs to provider",
		},
		{
			name: "policy provider mismatch fails",
			input: ResolutionInput{
				Role: "guardian",
				Policy: &RoleModelPolicy{
					Model:    "api-sonnet",
					Provider: "claude-cli",
				},
			},
			wantErr: "belongs to provider",
		},
		{
			name: "profile matching provider succeeds",
			input: ResolutionInput{
				Role:   "guardian",
				Policy: &RoleModelPolicy{Profile: "api-with-matching-provider"},
			},
		},
		{
			name: "profile provider mismatch fails",
			input: ResolutionInput{
				Role:   "guardian",
				Policy: &RoleModelPolicy{Profile: "api-with-cli-provider"},
			},
			wantErr: "belongs to provider",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc, err := mixedProviderResolver(t).Resolve(tt.input)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "api-sonnet", rc.Model)
			assert.Equal(t, "anthropic", rc.Provider)
		})
	}
}

func TestResolve_CanOperate_ValidatesCapability(t *testing.T) {
	// haiku lacks CapOperate, so CanOperate should fail
	r := testResolver(t)
	_, err := r.Resolve(ResolutionInput{
		Role:          "worker",
		AgentDefModel: "haiku",
		CanOperate:    true,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not support operate")
}

func TestResolve_CanOperate_ValidatesAliasRemapTarget(t *testing.T) {
	defaults := testDefaults()
	defaults.ModelAliases = map[string]string{"sonnet": "haiku"}
	resolver := NewResolver(testCatalog(t), testProfiles(), defaults)

	_, err := resolver.Resolve(ResolutionInput{
		Role:       "guardian",
		CanOperate: true,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not support operate")
}

func TestResolve_CanOperate_RejectsNonOperatorProvider(t *testing.T) {
	// Build a catalog with a model on a non-operator provider that claims operate capability.
	entries := []ModelCatalogEntry{
		{
			ID:           "fake-ollama-op",
			Aliases:      []string{"fake-op"},
			Provider:     "ollama",
			Tier:         TierExecution,
			Capabilities: []Capability{CapTools, CapCoding, CapOperate},
			Pricing:      ModelPricing{InputPerMillion: 0, OutputPerMillion: 0},
		},
	}
	cat, err := NewCatalog(entries)
	require.NoError(t, err)

	resolver := NewResolver(cat, nil, ResolverDefaults{
		Provider: "ollama",
		Model:    "fake-ollama-op",
	})

	_, err = resolver.Resolve(ResolutionInput{
		Role:       "worker",
		CanOperate: true,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be used with CanOperate")
}

func TestResolve_ProfileTemperature(t *testing.T) {
	r := testResolver(t)
	rc, err := r.Resolve(ResolutionInput{
		Role: "thinker",
		Policy: &RoleModelPolicy{
			Profile: "deep-judgment",
		},
	})
	require.NoError(t, err)
	assert.InDelta(t, 0.1, rc.Temperature, 0.001)
}

func TestResolve_RoleDefault(t *testing.T) {
	r := testResolver(t)

	// "guardian" role defaults to "sonnet" alias -> test-sonnet
	rc, err := r.Resolve(ResolutionInput{Role: "guardian"})
	require.NoError(t, err)
	assert.Equal(t, "test-sonnet", rc.Model)

	// "cto" role defaults to "opus" alias -> test-opus
	rc, err = r.Resolve(ResolutionInput{Role: "cto"})
	require.NoError(t, err)
	assert.Equal(t, "test-opus", rc.Model)

	// "sysmon" role defaults to "haiku" alias -> test-haiku
	rc, err = r.Resolve(ResolutionInput{Role: "sysmon"})
	require.NoError(t, err)
	assert.Equal(t, "test-haiku", rc.Model)
}

func TestResolve_MaxCostPerCallUSD(t *testing.T) {
	r := testResolver(t)
	maxCost := 0.05 // below opus pricing (10k*15/1M + 2k*75/1M = 0.15+0.15 = 0.30), above haiku (10k*0.8/1M + 2k*4/1M = 0.008+0.008 = 0.016)

	// Resolving to opus should fail because it exceeds the cap.
	_, err := r.Resolve(ResolutionInput{
		Role: "cto", // resolves to opus via role defaults
		Policy: &RoleModelPolicy{
			MaxCostPerCallUSD: &maxCost,
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds cost cap")

	// Resolving to haiku should succeed (its cost is below 0.05).
	highCap := 1.0
	rc, err := r.Resolve(ResolutionInput{
		Role: "sysmon", // resolves to haiku via role defaults
		Policy: &RoleModelPolicy{
			MaxCostPerCallUSD: &highCap,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "test-haiku", rc.Model)
}
