package v39

import (
	"errors"
	"strings"
	"testing"
)

func TestPromoteCapabilityVersionRecordsRequiredEvidenceEdges(t *testing.T) {
	store, version := capabilityEvolutionStore(t, capabilityEvolutionOptions{})

	recorded, err := store.PromoteCapabilityVersion(version)
	if err != nil {
		t.Fatalf("promote capability version: %v", err)
	}

	edges := store.EdgesFrom(recorded.CommonNode.ID)
	for _, want := range []string{EdgePromotedTo, EdgeOptimizedBy, EdgeEvaluatedBy, EdgeReviewedBy, EdgeRolledBackTo} {
		if !hasEdgeType(edges, want) {
			t.Fatalf("promotion missing %s edge: %+v", want, edges)
		}
	}
}

func TestPromoteCapabilityVersionRequiresHumanReview(t *testing.T) {
	store, version := capabilityEvolutionStore(t, capabilityEvolutionOptions{omitReview: true})

	_, err := store.PromoteCapabilityVersion(version)
	if !errors.Is(err, ErrRequiredPathMissing) {
		t.Fatalf("expected missing human review to block promotion, got %v", err)
	}
	if !strings.Contains(err.Error(), "HumanReview") {
		t.Fatalf("expected HumanReview in error, got %v", err)
	}
}

func TestPromoteCapabilityVersionBlocksBenchmarkRegression(t *testing.T) {
	store, version := capabilityEvolutionStore(t, capabilityEvolutionOptions{benchmarkStatus: "fail"})

	_, err := store.PromoteCapabilityVersion(version)
	if !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expected benchmark regression to block promotion, got %v", err)
	}
	if !strings.Contains(err.Error(), "benchmark regression") {
		t.Fatalf("expected benchmark regression in error, got %v", err)
	}
}

func TestPromoteCapabilityVersionRequiresRollbackTarget(t *testing.T) {
	store, version := capabilityEvolutionStore(t, capabilityEvolutionOptions{omitRollback: true})

	_, err := store.PromoteCapabilityVersion(version)
	if !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expected missing rollback target to block promotion, got %v", err)
	}
	if !strings.Contains(err.Error(), "rollback_to") {
		t.Fatalf("expected rollback_to in error, got %v", err)
	}
}

func TestPromoteCapabilityVersionBlocksOptimizerActor(t *testing.T) {
	store, version := capabilityEvolutionStore(t, capabilityEvolutionOptions{promoterIsOptimizer: true})

	_, err := store.PromoteCapabilityVersion(version)
	if !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expected optimizer actor to be blocked from promotion, got %v", err)
	}
	if !strings.Contains(err.Error(), "optimizer actor") {
		t.Fatalf("expected optimizer actor in error, got %v", err)
	}
}

func TestActivateCapabilityVersionRejectsGlobalScopeAndMissingRollback(t *testing.T) {
	store, version := capabilityEvolutionStore(t, capabilityEvolutionOptions{})
	if _, err := store.PromoteCapabilityVersion(version); err != nil {
		t.Fatalf("promote capability version: %v", err)
	}

	globalPolicy := activationPolicy("activation_global", version.CommonNode.ID, "global")
	runtime := capabilityRuntimeVersion("frv_cap_global", version.CommonNode.ID)
	if _, _, err := store.ActivateCapabilityVersion(globalPolicy, runtime); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expected global activation to be rejected, got %v", err)
	}

	baselinePolicy := activationPolicy("activation_missing_rollback", "cap_version_base", "project")
	baselineRuntime := capabilityRuntimeVersion("frv_cap_missing_rollback", "cap_version_base")
	if _, _, err := store.ActivateCapabilityVersion(baselinePolicy, baselineRuntime); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expected missing rollback target to block activation, got %v", err)
	}
}

func TestCertificationFailsWhenCapabilityInfluenceLacksUsageEvidence(t *testing.T) {
	store := stage6BaseStoreWithTaskSourceRefs(t, []string{"tool_description:cap_art_tool_001"})
	recordStage6ReleaseCandidate(t, store, "rc_missing_capability_usage")

	_, err := store.CertifyReleaseCandidate(stage6Certification("cert_missing_capability_usage", "rc_missing_capability_usage"))
	if !errors.Is(err, ErrRequiredPathMissing) {
		t.Fatalf("expected missing capability usage evidence to block certification, got %v", err)
	}
	if !strings.Contains(err.Error(), "USED_CAPABILITY") {
		t.Fatalf("expected USED_CAPABILITY in error, got %v", err)
	}
}

func TestCertificationPassesWithCapabilityArtifactUsageEvidence(t *testing.T) {
	store := stage6BaseStoreWithTaskSourceRefs(t, []string{"tool_description:cap_art_tool_001"})
	if _, err := store.RecordCapabilityUsage("tsk_001", "cap_art_tool_001", common("capability_usage_001", TypeCapabilityArtifact, "recorded")); err != nil {
		t.Fatalf("record capability usage: %v", err)
	}
	recordStage6ReleaseCandidate(t, store, "rc_with_capability_usage")

	cert, err := store.CertifyReleaseCandidate(stage6Certification("cert_with_capability_usage", "rc_with_capability_usage"))
	if err != nil {
		t.Fatalf("certify with capability usage: %v", err)
	}
	if !containsString(cert.EvidenceRefs, "cap_art_tool_001") {
		t.Fatalf("certification missing capability artifact evidence: %+v", cert.EvidenceRefs)
	}
}

type capabilityEvolutionOptions struct {
	omitReview          bool
	benchmarkStatus     string
	omitRollback        bool
	promoterIsOptimizer bool
}

func capabilityEvolutionStore(t *testing.T, opts capabilityEvolutionOptions) (*InMemoryStore, *CapabilityVersion) {
	t.Helper()
	store := NewInMemoryStore()
	appendRecord(t, store, &CapabilityArtifact{CommonNode: common("cap_art_e2", TypeCapabilityArtifact, "recorded"), ArtifactID: "capability_artifact_e2", ArtifactType: "tool_description", Name: "E2 planner tool description", ArtifactVersion: "1.1.0", SourceRepoOrOrigin: "eventgraph", ContentHash: "sha256:e2-tool-description", Owner: "eventgraph", RiskClass: "medium", ActivationScope: "project", EvalRefs: []string{"bench_e2"}, HumanReviewRef: "review_e2", RollbackRef: "cap_version_base", UsageLoggingRequired: true})
	appendRecord(t, store, &EvolutionOrder{CommonNode: common("evo_e2", TypeEvolutionOrder, "accepted"), EvolutionOrderVersion: 1, TargetCapabilityType: "tool_description", TargetRepo: "eventgraph", TargetPath: "go/pkg/darkfactory/v39", RiskClass: "medium", Motivation: "improve deterministic tool-description selection", EvalSource: "golden benchmark", Constraints: []string{"text_only"}, ReviewRequirements: []string{"CapabilityReviewer approval"}})
	appendRecord(t, store, &EvalDataset{CommonNode: common("eval_e2", TypeEvalDataset, "active"), SourceType: "golden", TrustLevel: "reviewed", TrainCount: 20, ValidationCount: 10, HoldoutCount: 30})
	runCommon := common("opt_e2", TypeOptimizationRun, "succeeded")
	runCommon.CreatedBy = "act_optimizer"
	appendRecord(t, store, &OptimizationRun{CommonNode: runCommon, EvolutionOrderID: "evo_e2", EvalDatasetID: "eval_e2", Engine: "manual"})
	appendRecord(t, store, &CandidateVariant{CommonNode: common("cand_e2", TypeCandidateVariant, "approved"), OptimizationRunID: "opt_e2", CapabilityArtifactID: "cap_art_e2"})
	benchmarkStatus := opts.benchmarkStatus
	if benchmarkStatus == "" {
		benchmarkStatus = "pass"
	}
	appendRecord(t, store, &BenchmarkResult{CommonNode: common("bench_e2", TypeBenchmarkResult, benchmarkStatus), CandidateVariantID: "cand_e2", BaselineRef: "cap_version_base", MetricDeltas: map[string]float64{"success_rate": 0.05, "regression_count": 0}})
	if !opts.omitReview {
		appendRecord(t, store, &HumanReview{CommonNode: common("review_e2", TypeHumanReview, "approved"), ReviewerActorID: "act_human", ReviewerRole: "CapabilityReviewer", Rationale: "approved text-only non-regressing candidate"})
	}
	appendRecord(t, store, &CapabilityVersion{CommonNode: common("cap_version_base", TypeCapabilityVersion, "active"), CapabilityArtifactID: "cap_art_e2", CapabilitySemver: "1.0.0"})

	versionCommon := common("cap_version_e2", TypeCapabilityVersion, "approved")
	versionCommon.CreatedBy = "act_capability_release"
	if opts.promoterIsOptimizer {
		versionCommon.CreatedBy = "act_optimizer"
	}
	version := &CapabilityVersion{CommonNode: versionCommon, CapabilityArtifactID: "cap_art_e2", CapabilitySemver: "1.1.0"}
	if !opts.omitRollback {
		version.RollbackTo = strPtr("cap_version_base")
	}
	return store, version
}

func activationPolicy(id, capabilityVersionID, scope string) *ActivationPolicy {
	return &ActivationPolicy{CommonNode: common(id, TypeActivationPolicy, "approved"), ActivationPolicyID: id, CapabilityVersionID: capabilityVersionID, Scope: scope, AllowedProjects: []string{"dark-factory"}, MonitoringWindowRuns: 5, RollbackTriggers: []string{"benchmark_regression"}, ApprovedBy: []string{"act_human"}}
}

func capabilityRuntimeVersion(id, capabilityVersionID string) *FactoryRuntimeVersion {
	return &FactoryRuntimeVersion{CommonNode: common(id, TypeFactoryRuntimeVersion, "active"), RuntimeVersion: "3.9.0", CapabilityVersionRefs: []string{capabilityVersionID}, RuntimeRefs: []string{"local@0.1.0"}}
}

func hasEdgeType(edges []CommonEdge, edgeType string) bool {
	for _, edge := range edges {
		if edge.Type == edgeType {
			return true
		}
	}
	return false
}
