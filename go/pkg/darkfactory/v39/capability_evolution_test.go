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
	for _, want := range []struct {
		typ string
		to  string
	}{
		{EdgePromotedTo, "cap_art_e2"},
		{EdgeOptimizedBy, "evo_e2"},
		{EdgeEvaluatedBy, "eval_e2"},
		{EdgeEvaluatedBy, "bench_e2"},
		{EdgeReviewedBy, "review_e2"},
		{EdgeRolledBackTo, "cap_version_base"},
	} {
		if !hasEdge(edges, want.typ, want.to) {
			t.Fatalf("promotion missing %s -> %s edge: %+v", want.typ, want.to, edges)
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

func TestPromoteCapabilityVersionUsesExplicitCandidateAndBenchmarkRefs(t *testing.T) {
	store, _ := capabilityEvolutionStore(t, capabilityEvolutionOptions{})
	appendRecord(t, store, &CandidateVariant{CommonNode: common("cand_bad", TypeCandidateVariant, "approved"), OptimizationRunID: "opt_e2", CapabilityArtifactID: "cap_art_e2"})
	appendRecord(t, store, &BenchmarkResult{CommonNode: common("bench_bad", TypeBenchmarkResult, "fail"), CandidateVariantID: "cand_bad", BaselineRef: "cap_version_base", MetricDeltas: map[string]float64{"regression_count": 1}})
	reviewCommon := common("review_bad", TypeHumanReview, "approved")
	reviewCommon.SourceRefs = []string{"cand_bad", "bench_bad"}
	appendRecord(t, store, &HumanReview{CommonNode: reviewCommon, ReviewerActorID: "act_human", ReviewerRole: "CapabilityReviewer", Rationale: "approved only for the bad candidate to test exact refs"})

	versionCommon := common("cap_version_bad", TypeCapabilityVersion, "approved")
	version := &CapabilityVersion{CommonNode: versionCommon, CapabilityArtifactID: "cap_art_e2", EvolutionOrderID: "evo_e2", OptimizationRunID: "opt_e2", CandidateVariantID: "cand_bad", EvalDatasetID: "eval_e2", BenchmarkResultID: "bench_bad", HumanReviewID: "review_bad", PromoterActorID: "act_capability_release", PromoterRole: "CapabilityRelease", CapabilitySemver: "1.2.0", RollbackTo: strPtr("cap_version_base")}

	_, err := store.PromoteCapabilityVersion(version)
	if !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expected explicitly referenced failing benchmark to block promotion, got %v", err)
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

	if _, err := store.AppendRecord(activationPolicy("activation_policy_direct_missing_rollback", "cap_version_base", "project")); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expected direct activation policy append without rollback target to fail, got %v", err)
	}
}

func TestActivateCapabilityVersionRecordsActivationEdges(t *testing.T) {
	store, version := capabilityEvolutionStore(t, capabilityEvolutionOptions{})
	if _, err := store.PromoteCapabilityVersion(version); err != nil {
		t.Fatalf("promote capability version: %v", err)
	}

	policy := activationPolicy("activation_success", version.CommonNode.ID, "project")
	runtime := capabilityRuntimeVersion("frv_cap_success", version.CommonNode.ID)
	_, recordedRuntime, err := store.ActivateCapabilityVersion(policy, runtime)
	if err != nil {
		t.Fatalf("activate capability version: %v", err)
	}

	edges := store.EdgesFrom(version.CommonNode.ID)
	if !hasEdge(edges, EdgeActivatedBy, policy.CommonNode.ID) {
		t.Fatalf("activation missing %s edge: %+v", EdgeActivatedBy, edges)
	}
	if !hasEdge(edges, EdgePackagedAs, recordedRuntime.CommonNode.ID) {
		t.Fatalf("activation missing %s edge: %+v", EdgePackagedAs, edges)
	}
}

func TestActivateCapabilityVersionRequiresPromotionEvidence(t *testing.T) {
	store, version := capabilityEvolutionStore(t, capabilityEvolutionOptions{})
	appendRecord(t, store, version)

	policy := activationPolicy("activation_unpromoted", version.CommonNode.ID, "project")
	runtime := capabilityRuntimeVersion("frv_cap_unpromoted", version.CommonNode.ID)
	if _, _, err := store.ActivateCapabilityVersion(policy, runtime); !errors.Is(err, ErrRequiredPathMissing) {
		t.Fatalf("expected promotion evidence to be required before activation, got %v", err)
	}
}

func TestAppendRecordRejectsDirectActiveCapabilityVersion(t *testing.T) {
	store, _ := capabilityEvolutionStore(t, capabilityEvolutionOptions{})

	_, err := store.AppendRecord(&CapabilityVersion{CommonNode: common("cap_version_direct_active", TypeCapabilityVersion, "active"), CapabilityArtifactID: "cap_art_e2", CapabilitySemver: "2.0.0"})
	if !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expected direct active CapabilityVersion append to fail, got %v", err)
	}
}

func TestAppendEdgeRejectsDirectCapabilityGovernanceEdge(t *testing.T) {
	store := stage6BaseStoreWithTaskSourceRefs(t, []string{"tool_description:cap_art_tool_001"})

	_, err := store.AppendEdge(derivedEdge(EdgeUsedCapability, "tsk_001", "cap_art_tool_001", common("direct_usage_edge", TypeCapabilityArtifact, "recorded")))
	if !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expected direct USED_CAPABILITY edge append to fail, got %v", err)
	}
}

func TestPromoteCapabilityVersionRejectsFabricatedReviewSourceRefs(t *testing.T) {
	store, version := capabilityEvolutionStore(t, capabilityEvolutionOptions{omitReview: true})
	reviewCommon := common("review_e2", TypeHumanReview, "approved")
	reviewCommon.SourceRefs = []string{"cap_art_e2"}
	appendRecord(t, store, &HumanReview{CommonNode: reviewCommon, ReviewerActorID: "act_human", ReviewerRole: "CapabilityReviewer", Rationale: "does not name the candidate or benchmark"})

	_, err := store.PromoteCapabilityVersion(version)
	if !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expected fabricated review source refs to be rejected, got %v", err)
	}
	if !strings.Contains(err.Error(), "source_refs") {
		t.Fatalf("expected source_refs in error, got %v", err)
	}
}

func TestActivationPolicyRejectsBroadScopes(t *testing.T) {
	for name, policy := range map[string]*ActivationPolicy{
		"global":         activationPolicy("activation_invalid_global", "cap_version_e2", "global"),
		"project_empty":  activationPolicy("activation_invalid_project", "cap_version_e2", "project"),
		"order_empty":    activationPolicy("activation_invalid_order", "cap_version_e2", "order"),
		"canary_full":    activationPolicy("activation_invalid_canary", "cap_version_e2", "canary"),
		"canary_missing": activationPolicy("activation_invalid_canary_missing", "cap_version_e2", "canary"),
	} {
		switch name {
		case "project_empty":
			policy.AllowedProjects = nil
		case "order_empty":
			policy.AllowedProjects = nil
			policy.AllowedFactoryOrders = nil
		case "canary_full":
			canary := 100.0
			policy.CanaryPercent = &canary
		}
		if err := policy.Validate(); !errors.Is(err, ErrInvalidRecord) {
			t.Fatalf("%s: expected invalid activation policy, got %v", name, err)
		}
	}
}

func TestRecordRollbackRecordBlocksReactivation(t *testing.T) {
	store, version := capabilityEvolutionStore(t, capabilityEvolutionOptions{})
	if _, err := store.PromoteCapabilityVersion(version); err != nil {
		t.Fatalf("promote capability version: %v", err)
	}
	policy := activationPolicy("activation_before_rollback", version.CommonNode.ID, "project")
	runtime := capabilityRuntimeVersion("frv_cap_before_rollback", version.CommonNode.ID)
	if _, _, err := store.ActivateCapabilityVersion(policy, runtime); err != nil {
		t.Fatalf("activate capability version before rollback: %v", err)
	}
	_, err := store.RecordRollbackRecord(&RollbackRecord{CommonNode: common("rollback_completed", TypeRollbackRecord, "completed"), CapabilityVersionID: version.CommonNode.ID, RollbackTo: "cap_version_base", Trigger: "benchmark_regression", ActorID: "act_human", FactoryRuntimeVersionID: runtime.CommonNode.ID})
	if err != nil {
		t.Fatalf("record rollback: %v", err)
	}

	nextPolicy := activationPolicy("activation_after_rollback", version.CommonNode.ID, "project")
	nextRuntime := capabilityRuntimeVersion("frv_cap_after_rollback", version.CommonNode.ID)
	if _, _, err := store.ActivateCapabilityVersion(nextPolicy, nextRuntime); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expected rollback record to block reactivation, got %v", err)
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

func TestCertificationCapabilityInfluenceRecognizesBareAndCaseInsensitiveSourceRefs(t *testing.T) {
	for name, sourceRef := range map[string]string{
		"bare_id":          "cap_art_tool_001",
		"case_insensitive": "Tool_Description://cap_art_tool_001",
	} {
		store := stage6BaseStoreWithTaskSourceRefs(t, []string{sourceRef})
		if _, err := store.RecordCapabilityUsage("tsk_001", "cap_art_tool_001", common("capability_usage_"+name, TypeCapabilityArtifact, "recorded")); err != nil {
			t.Fatalf("%s: record capability usage: %v", name, err)
		}
		recordStage6ReleaseCandidate(t, store, "rc_with_capability_usage_"+name)

		if _, err := store.CertifyReleaseCandidate(stage6Certification("cert_with_capability_usage_"+name, "rc_with_capability_usage_"+name)); err != nil {
			t.Fatalf("%s: certify with capability usage: %v", name, err)
		}
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
		reviewCommon := common("review_e2", TypeHumanReview, "approved")
		reviewCommon.SourceRefs = []string{"cand_e2", "bench_e2"}
		appendRecord(t, store, &HumanReview{CommonNode: reviewCommon, ReviewerActorID: "act_human", ReviewerRole: "CapabilityReviewer", Rationale: "approved text-only non-regressing candidate"})
	}
	appendRecord(t, store, &CapabilityVersion{CommonNode: common("cap_version_base", TypeCapabilityVersion, "approved"), CapabilityArtifactID: "cap_art_e2", CapabilitySemver: "1.0.0"})

	versionCommon := common("cap_version_e2", TypeCapabilityVersion, "approved")
	versionCommon.CreatedBy = "act_capability_release"
	if opts.promoterIsOptimizer {
		versionCommon.CreatedBy = "act_optimizer"
	}
	version := &CapabilityVersion{CommonNode: versionCommon, CapabilityArtifactID: "cap_art_e2", EvolutionOrderID: "evo_e2", OptimizationRunID: "opt_e2", CandidateVariantID: "cand_e2", EvalDatasetID: "eval_e2", BenchmarkResultID: "bench_e2", HumanReviewID: "review_e2", PromoterActorID: "act_capability_release", PromoterRole: "CapabilityRelease", CapabilitySemver: "1.1.0"}
	if opts.promoterIsOptimizer {
		version.PromoterActorID = "act_optimizer"
	}
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

func hasEdge(edges []CommonEdge, edgeType, toID string) bool {
	for _, edge := range edges {
		if edge.Type == edgeType && edge.ToID == toID {
			return true
		}
	}
	return false
}
