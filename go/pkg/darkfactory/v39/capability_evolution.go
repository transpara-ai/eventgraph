package v39

import (
	"fmt"
	"strings"
)

const (
	CapabilityPromotionAction = "capability.promote"
	CapabilityReleaseRole     = "CapabilityRelease"
)

type capabilityPromotionEvidence struct {
	artifact  *CapabilityArtifact
	candidate *CandidateVariant
	run       *OptimizationRun
	order     *EvolutionOrder
	dataset   *EvalDataset
	benchmark *BenchmarkResult
	review    *HumanReview
	rollback  *CapabilityVersion
	authority capabilityPromotionAuthority
}

type capabilityPromotionAuthority struct {
	identity *ActorIdentity
	request  *AuthorityRequest
	decision *AuthorityDecision
	receipt  *ExecutionReceipt
}

type CapabilityArtifactUsageLoggingFinding struct {
	CapabilityArtifactRecordID string `json:"capability_artifact_record_id"`
	ArtifactID                 string `json:"artifact_id"`
	ArtifactType               string `json:"artifact_type"`
	Name                       string `json:"name"`
	Reason                     string `json:"reason"`
}

type capabilitySourceRef struct {
	raw          string
	artifactType string
	value        string
}

func (s *InMemoryStore) RecordCapabilityUsage(sourceID, capabilityArtifactID string, common CommonNode) (CommonEdge, error) {
	if _, err := s.Get(sourceID); err != nil {
		return CommonEdge{}, err
	}
	artifact, ok := s.findCapabilityArtifact(capabilityArtifactID)
	if !ok {
		return CommonEdge{}, fmt.Errorf("%w: CapabilityArtifact %s", ErrNotFound, capabilityArtifactID)
	}
	if !artifact.UsageLoggingRequired {
		return CommonEdge{}, fieldError(TypeCapabilityArtifact, "usage_logging_required", "must be true before capability use")
	}
	return s.appendTrustedEdge(derivedEdge(EdgeUsedCapability, sourceID, artifact.CommonNode.ID, common))
}

func (s *InMemoryStore) CapabilityUsageEvidencePath(releaseCandidateID string) (RequiredPath, error) {
	path := RequiredPath{Name: "Material capability influence -> CapabilityArtifact / USED_CAPABILITY", NodeIDs: []string{releaseCandidateID}}
	rc, ok := s.mustGetReleaseCandidate(releaseCandidateID)
	if !ok {
		path.Missing = append(path.Missing, "ReleaseCandidate "+releaseCandidateID)
		return path, path.Err()
	}
	orderPath, _ := s.FactoryOrderRequirementAcceptanceTask(rc.FactoryOrderID)
	path.EdgeIDs = append(path.EdgeIDs, orderPath.EdgeIDs...)
	path.Missing = appendUniqueStrings(path.Missing, orderPath.Missing...)

	for _, taskID := range taskIDsFromPath(s, orderPath) {
		task, ok := s.mustGetTask(taskID)
		if !ok {
			path.Missing = appendUniqueStrings(path.Missing, "Task "+taskID)
			continue
		}
		for _, sourceRef := range s.capabilitySourceRefs(task.SourceRefs) {
			artifact, ok := s.findCapabilityArtifactForSourceRef(sourceRef)
			if !ok {
				path.Missing = appendUniqueStrings(path.Missing, "CapabilityArtifact for "+sourceRef.raw+" used in Task "+taskID)
				continue
			}
			path.NodeIDs = appendUniqueStrings(path.NodeIDs, taskID, artifact.CommonNode.ID)
			edgeIDs := edgeIDsBetween(s, taskID, artifact.CommonNode.ID, EdgeUsedCapability)
			if len(edgeIDs) == 0 {
				path.Missing = appendUniqueStrings(path.Missing, "USED_CAPABILITY from Task "+taskID+" to CapabilityArtifact "+artifact.CommonNode.ID)
				continue
			}
			path.EdgeIDs = appendUniqueStrings(path.EdgeIDs, edgeIDs...)
		}
	}

	path.Completed = len(path.Missing) == 0
	return path, path.Err()
}

func (s *InMemoryStore) PromoteCapabilityVersion(version *CapabilityVersion) (*CapabilityVersion, error) {
	if version == nil {
		return nil, fmt.Errorf("%w: nil CapabilityVersion", ErrInvalidRecord)
	}
	evidence, err := s.evaluateCapabilityPromotion(version)
	if err != nil {
		return nil, err
	}
	stored, err := s.AppendRecord(version)
	if err != nil {
		return nil, err
	}
	capabilityVersion, ok := stored.(*CapabilityVersion)
	if !ok {
		return nil, fmt.Errorf("%w: CapabilityVersion append returned %T", ErrInvalidRecord, stored)
	}
	edges := []struct {
		typ string
		to  string
	}{
		{EdgePromotedTo, evidence.artifact.CommonNode.ID},
		{EdgeOptimizedBy, evidence.order.CommonNode.ID},
		{EdgeEvaluatedBy, evidence.dataset.CommonNode.ID},
		{EdgeEvaluatedBy, evidence.benchmark.CommonNode.ID},
		{EdgeReviewedBy, evidence.review.CommonNode.ID},
		{EdgeRolledBackTo, evidence.rollback.CommonNode.ID},
	}
	for _, edgeSpec := range edges {
		edge := derivedEdge(edgeSpec.typ, capabilityVersion.CommonNode.ID, edgeSpec.to, capabilityVersion.CommonNode)
		if edgeSpec.typ == EdgePromotedTo {
			edge.EvidenceRefs = appendUniqueStrings(edge.EvidenceRefs, evidence.authority.identity.CommonNode.ID, evidence.authority.request.CommonNode.ID, evidence.authority.decision.CommonNode.ID, evidence.authority.receipt.CommonNode.ID)
		}
		if _, err := s.appendTrustedEdge(edge); err != nil {
			return nil, err
		}
	}
	return capabilityVersion, nil
}

func (s *InMemoryStore) ActivateCapabilityVersion(policy *ActivationPolicy, runtimeVersion *FactoryRuntimeVersion) (*ActivationPolicy, *FactoryRuntimeVersion, error) {
	if policy == nil {
		return nil, nil, fmt.Errorf("%w: nil ActivationPolicy", ErrInvalidRecord)
	}
	if runtimeVersion == nil {
		return nil, nil, fmt.Errorf("%w: nil FactoryRuntimeVersion", ErrInvalidRecord)
	}
	if policy.Scope == "global" {
		return nil, nil, fieldError(TypeActivationPolicy, "scope", "global activation disabled for MVP")
	}
	version, ok := s.mustGetCapabilityVersion(policy.CapabilityVersionID)
	if !ok {
		return nil, nil, fmt.Errorf("%w: CapabilityVersion %s", ErrNotFound, policy.CapabilityVersionID)
	}
	if version.RollbackTo == nil || *version.RollbackTo == "" {
		return nil, nil, fieldError(TypeCapabilityVersion, "rollback_to", "required before activation")
	}
	if version.CommonNode.Status == nil || *version.CommonNode.Status != "approved" {
		return nil, nil, fieldError(TypeCapabilityVersion, "status", "must be approved before activation")
	}
	if s.capabilityVersionHasActiveRollback(version.CommonNode.ID) {
		return nil, nil, fieldError(TypeCapabilityVersion, "status", "rolled back capability version cannot be activated")
	}
	if !s.capabilityVersionHasPromotionEvidence(version) {
		return nil, nil, fmt.Errorf("%w: promotion evidence for CapabilityVersion %s", ErrRequiredPathMissing, version.CommonNode.ID)
	}
	if len(policy.ApprovedBy) == 0 {
		return nil, nil, fieldError(TypeActivationPolicy, "approved_by", "required before activation")
	}
	if !containsString(runtimeVersion.CapabilityVersionRefs, version.CommonNode.ID) {
		return nil, nil, fieldError(TypeFactoryRuntimeVersion, "capability_version_refs", "must include activated capability version")
	}

	storedPolicy, err := s.AppendRecord(policy)
	if err != nil {
		return nil, nil, err
	}
	activationPolicy, ok := storedPolicy.(*ActivationPolicy)
	if !ok {
		return nil, nil, fmt.Errorf("%w: ActivationPolicy append returned %T", ErrInvalidRecord, storedPolicy)
	}
	storedRuntime, err := s.RecordFactoryRuntimeVersionBOM(runtimeVersion)
	if err != nil {
		return nil, nil, err
	}
	if _, err := s.appendTrustedEdge(derivedEdge(EdgeActivatedBy, version.CommonNode.ID, activationPolicy.CommonNode.ID, activationPolicy.CommonNode)); err != nil {
		return nil, nil, err
	}
	if _, err := s.appendTrustedEdge(derivedEdge(EdgePackagedAs, version.CommonNode.ID, storedRuntime.CommonNode.ID, storedRuntime.CommonNode)); err != nil {
		return nil, nil, err
	}
	return activationPolicy, storedRuntime, nil
}

func (s *InMemoryStore) RecordRollbackRecord(record *RollbackRecord) (*RollbackRecord, error) {
	if record == nil {
		return nil, fmt.Errorf("%w: nil RollbackRecord", ErrInvalidRecord)
	}
	stored, err := s.AppendRecord(record)
	if err != nil {
		return nil, err
	}
	rollback, ok := stored.(*RollbackRecord)
	if !ok {
		return nil, fmt.Errorf("%w: RollbackRecord append returned %T", ErrInvalidRecord, stored)
	}
	edge := derivedEdge(EdgeRolledBackTo, rollback.CapabilityVersionID, rollback.RollbackTo, rollback.CommonNode)
	edge.ID = "edge:" + rollback.CommonNode.ID + ":" + EdgeRolledBackTo + ":" + rollback.RollbackTo
	edge.IdempotencyKey = edge.ID
	edge.EvidenceRefs = []string{rollback.CommonNode.ID}
	if _, err := s.appendTrustedEdge(edge); err != nil {
		return nil, err
	}
	return rollback, nil
}

func (s *InMemoryStore) evaluateCapabilityPromotion(version *CapabilityVersion) (capabilityPromotionEvidence, error) {
	if err := version.Validate(); err != nil {
		return capabilityPromotionEvidence{}, err
	}
	if version.CommonNode.Status == nil || *version.CommonNode.Status != "approved" {
		return capabilityPromotionEvidence{}, fieldError(TypeCapabilityVersion, "status", "must be approved for promotion")
	}
	return s.resolveCapabilityPromotionEvidence(version)
}

func (s *InMemoryStore) resolveCapabilityPromotionEvidence(version *CapabilityVersion) (capabilityPromotionEvidence, error) {
	var evidence capabilityPromotionEvidence
	for field, value := range map[string]string{
		"evolution_order_id":   version.EvolutionOrderID,
		"optimization_run_id":  version.OptimizationRunID,
		"candidate_variant_id": version.CandidateVariantID,
		"eval_dataset_id":      version.EvalDatasetID,
		"benchmark_result_id":  version.BenchmarkResultID,
		"human_review_id":      version.HumanReviewID,
		"promoter_actor_id":    version.PromoterActorID,
		"promoter_role":        version.PromoterRole,
	} {
		if err := requireNonEmpty(TypeCapabilityVersion, field, value); err != nil {
			return evidence, err
		}
	}
	if version.PromoterRole != CapabilityReleaseRole {
		return evidence, fieldError(TypeCapabilityVersion, "promoter_role", "must be CapabilityRelease")
	}
	authority, err := s.resolveCapabilityPromotionAuthority(version)
	if err != nil {
		return evidence, err
	}
	evidence.authority = authority
	if !s.capabilityVersionHasRollbackTarget(version) {
		return evidence, fieldError(TypeCapabilityVersion, "rollback_to", "required before promotion")
	}
	rollback, ok := s.mustGetCapabilityVersion(*version.RollbackTo)
	if !ok {
		return evidence, fmt.Errorf("%w: rollback CapabilityVersion %s", ErrNotFound, *version.RollbackTo)
	}
	evidence.rollback = rollback

	artifact, ok := s.findCapabilityArtifact(version.CapabilityArtifactID)
	if !ok {
		return evidence, fmt.Errorf("%w: CapabilityArtifact %s", ErrNotFound, version.CapabilityArtifactID)
	}
	if artifact.ArtifactType != "skill" && artifact.ArtifactType != "tool_description" {
		return evidence, fieldError(TypeCapabilityArtifact, "artifact_type", "MVP promotion supports only skill or tool_description")
	}
	if artifact.RiskClass != "low" && artifact.RiskClass != "medium" {
		return evidence, fieldError(TypeCapabilityArtifact, "risk_class", "MVP promotion supports only low or medium risk")
	}
	evidence.artifact = artifact

	candidate, ok := s.mustGetCandidateVariant(version.CandidateVariantID)
	if !ok {
		return evidence, fmt.Errorf("%w: CandidateVariant %s", ErrNotFound, version.CandidateVariantID)
	}
	if candidate.CapabilityArtifactID != artifact.CommonNode.ID && candidate.CapabilityArtifactID != artifact.ArtifactID {
		return evidence, fieldError(TypeCandidateVariant, "capability_artifact_id", "must match promoted CapabilityArtifact")
	}
	if candidate.CommonNode.Status == nil || *candidate.CommonNode.Status != "approved" {
		return evidence, fieldError(TypeCandidateVariant, "status", "approved candidate required")
	}
	evidence.candidate = candidate

	run, ok := s.mustGetOptimizationRun(version.OptimizationRunID)
	if !ok {
		return evidence, fmt.Errorf("%w: OptimizationRun %s", ErrNotFound, version.OptimizationRunID)
	}
	if candidate.OptimizationRunID != run.CommonNode.ID {
		return evidence, fieldError(TypeCandidateVariant, "optimization_run_id", "must match promoted OptimizationRun")
	}
	if run.CommonNode.Status == nil || *run.CommonNode.Status != "succeeded" {
		return evidence, fieldError(TypeOptimizationRun, "status", "succeeded optimization run required")
	}
	if run.CommonNode.CreatedBy == version.PromoterActorID || candidate.CommonNode.CreatedBy == version.PromoterActorID {
		return evidence, fieldError(TypeCapabilityVersion, "promoter_actor_id", "optimizer actor cannot promote its own output")
	}
	evidence.run = run

	order, ok := s.mustGetEvolutionOrder(version.EvolutionOrderID)
	if !ok {
		return evidence, fmt.Errorf("%w: EvolutionOrder %s", ErrNotFound, version.EvolutionOrderID)
	}
	if run.EvolutionOrderID != order.CommonNode.ID {
		return evidence, fieldError(TypeOptimizationRun, "evolution_order_id", "must match promoted EvolutionOrder")
	}
	if order.CommonNode.Status == nil || *order.CommonNode.Status != "accepted" {
		return evidence, fieldError(TypeEvolutionOrder, "status", "accepted evolution order required")
	}
	evidence.order = order

	dataset, ok := s.mustGetEvalDataset(version.EvalDatasetID)
	if !ok {
		return evidence, fmt.Errorf("%w: EvalDataset %s", ErrNotFound, version.EvalDatasetID)
	}
	if run.EvalDatasetID != dataset.CommonNode.ID {
		return evidence, fieldError(TypeOptimizationRun, "eval_dataset_id", "must match promoted EvalDataset")
	}
	evidence.dataset = dataset

	benchmark, ok := s.mustGetBenchmarkResult(version.BenchmarkResultID)
	if !ok {
		return evidence, fmt.Errorf("%w: BenchmarkResult %s", ErrNotFound, version.BenchmarkResultID)
	}
	if benchmark.CandidateVariantID != candidate.CommonNode.ID {
		return evidence, fieldError(TypeBenchmarkResult, "candidate_variant_id", "must match promoted CandidateVariant")
	}
	if benchmark.CommonNode.Status == nil || *benchmark.CommonNode.Status != "pass" {
		return evidence, fieldError(TypeBenchmarkResult, "status", "benchmark regression blocks capability promotion")
	}
	if stale, ok := s.laterBlockingBenchmarkForCandidate(benchmark); ok {
		return evidence, fieldError(TypeBenchmarkResult, "created_at", "stale benchmark result "+benchmark.CommonNode.ID+" is superseded by later "+*stale.CommonNode.Status+" result "+stale.CommonNode.ID)
	}
	evidence.benchmark = benchmark

	review, ok := s.mustGetHumanReview(version.HumanReviewID)
	if !ok {
		return evidence, fmt.Errorf("%w: HumanReview %s", ErrRequiredPathMissing, version.HumanReviewID)
	}
	if review.CommonNode.Status == nil || *review.CommonNode.Status != "approved" {
		return evidence, fieldError(TypeHumanReview, "status", "approved HumanReview required")
	}
	if review.ReviewerRole != "CapabilityReviewer" {
		return evidence, fieldError(TypeHumanReview, "reviewer_role", "CapabilityReviewer required")
	}
	if !containsString(review.CommonNode.SourceRefs, candidate.CommonNode.ID) || !containsString(review.CommonNode.SourceRefs, benchmark.CommonNode.ID) {
		return evidence, fieldError(TypeHumanReview, "source_refs", "must include promoted CandidateVariant and BenchmarkResult")
	}
	evidence.review = review
	return evidence, nil
}

func (s *InMemoryStore) resolveCapabilityPromotionAuthority(version *CapabilityVersion) (capabilityPromotionAuthority, error) {
	identity, ok := s.actorIdentityForActor(version.PromoterActorID)
	if !ok {
		return capabilityPromotionAuthority{}, fmt.Errorf("%w: ActorIdentity for promoter actor %s", ErrRequiredPathMissing, version.PromoterActorID)
	}
	if identity.CommonNode.Status == nil || *identity.CommonNode.Status != "active" {
		return capabilityPromotionAuthority{}, fieldError(TypeActorIdentity, "status", "active promoter identity required")
	}

	for _, record := range s.ByType(TypeAuthorityRequest) {
		request := record.(*AuthorityRequest)
		if !capabilityPromotionRequestMatches(request, version) {
			continue
		}
		authorityPath, err := s.ActorAuthorityRequestDecisionReceipt(request.CommonNode.ID)
		if err != nil {
			return capabilityPromotionAuthority{}, err
		}
		if !authorityPath.Completed {
			return capabilityPromotionAuthority{}, fmt.Errorf("%w: capability promotion authority path for %s", ErrRequiredPathMissing, request.CommonNode.ID)
		}
		decisionID := authorityPath.NodeIDs[len(authorityPath.NodeIDs)-2]
		receiptID := authorityPath.NodeIDs[len(authorityPath.NodeIDs)-1]
		decision, ok := s.mustGetAuthorityDecision(decisionID)
		if !ok {
			return capabilityPromotionAuthority{}, fmt.Errorf("%w: AuthorityDecision %s", ErrNotFound, decisionID)
		}
		receipt, ok := s.mustGetExecutionReceipt(receiptID)
		if !ok {
			return capabilityPromotionAuthority{}, fmt.Errorf("%w: ExecutionReceipt %s", ErrNotFound, receiptID)
		}
		if !capabilityPromotionDecisionAuthorizes(decision, version) {
			return capabilityPromotionAuthority{}, fieldError(TypeAuthorityDecision, "decision", "must approve capability.promote for CapabilityRelease")
		}
		if !capabilityPromotionReceiptMatches(receipt, version) {
			return capabilityPromotionAuthority{}, fieldError(TypeExecutionReceipt, "result", "succeeded capability.promote receipt required")
		}
		return capabilityPromotionAuthority{identity: identity, request: request, decision: decision, receipt: receipt}, nil
	}
	return capabilityPromotionAuthority{}, fmt.Errorf("%w: AuthorityRequest for %s %s on CapabilityVersion %s", ErrRequiredPathMissing, CapabilityReleaseRole, CapabilityPromotionAction, version.CommonNode.ID)
}

func capabilityPromotionRequestMatches(request *AuthorityRequest, version *CapabilityVersion) bool {
	return request.ActorID == version.PromoterActorID &&
		request.ActorRole == CapabilityReleaseRole &&
		request.Action == CapabilityPromotionAction &&
		request.TargetType == TypeCapabilityVersion &&
		request.TargetID == version.CommonNode.ID
}

func capabilityPromotionDecisionAuthorizes(decision *AuthorityDecision, version *CapabilityVersion) bool {
	if decision == nil || decision.CommonNode.Status == nil || *decision.CommonNode.Status != "approved" {
		return false
	}
	if decision.Decision != "Autonomous" && decision.Decision != "ApprovalRequired" {
		return false
	}
	return containsString(decision.Scope, CapabilityPromotionAction) || containsString(decision.Scope, version.CommonNode.ID)
}

func capabilityPromotionReceiptMatches(receipt *ExecutionReceipt, version *CapabilityVersion) bool {
	return receipt != nil &&
		receipt.Action == CapabilityPromotionAction &&
		receipt.TargetID == version.CommonNode.ID &&
		receipt.Result == "succeeded"
}

// Benchmark freshness policy: CapabilityVersion.BenchmarkResultID is the
// reviewed benchmark, but it stops being authoritative if a later fail/error
// BenchmarkResult exists for the same CandidateVariant before promotion.
func (s *InMemoryStore) laterBlockingBenchmarkForCandidate(benchmark *BenchmarkResult) (*BenchmarkResult, bool) {
	for _, record := range s.ByType(TypeBenchmarkResult) {
		candidate := record.(*BenchmarkResult)
		if candidate.CommonNode.ID == benchmark.CommonNode.ID || candidate.CandidateVariantID != benchmark.CandidateVariantID {
			continue
		}
		if !candidate.CommonNode.CreatedAt.After(benchmark.CommonNode.CreatedAt) || candidate.CommonNode.Status == nil {
			continue
		}
		if *candidate.CommonNode.Status == "fail" || *candidate.CommonNode.Status == "error" {
			return candidate, true
		}
	}
	return nil, false
}

func (s *InMemoryStore) capabilityVersionHasPromotionEvidence(version *CapabilityVersion) bool {
	evidence, err := s.resolveCapabilityPromotionEvidence(version)
	if err != nil {
		return false
	}
	for _, edgeSpec := range []struct {
		typ string
		to  string
	}{
		{EdgePromotedTo, evidence.artifact.CommonNode.ID},
		{EdgeOptimizedBy, evidence.order.CommonNode.ID},
		{EdgeEvaluatedBy, evidence.dataset.CommonNode.ID},
		{EdgeEvaluatedBy, evidence.benchmark.CommonNode.ID},
		{EdgeReviewedBy, evidence.review.CommonNode.ID},
		{EdgeRolledBackTo, evidence.rollback.CommonNode.ID},
	} {
		if len(edgeIDsBetween(s, version.CommonNode.ID, edgeSpec.to, edgeSpec.typ)) == 0 {
			return false
		}
	}
	return true
}

// CapabilityArtifactUsageLoggingFindings inventories legacy or externally
// loaded CapabilityArtifact records that predate the v3.9 usage logging
// invariant. The append-only v3.9 store does not mutate those records in
// place; callers must backfill by re-emitting compliant artifact records before
// material capability use.
func (s *InMemoryStore) CapabilityArtifactUsageLoggingFindings() []CapabilityArtifactUsageLoggingFinding {
	var findings []CapabilityArtifactUsageLoggingFinding
	for _, record := range s.ByType(TypeCapabilityArtifact) {
		artifact := record.(*CapabilityArtifact)
		if artifact.UsageLoggingRequired {
			continue
		}
		findings = append(findings, CapabilityArtifactUsageLoggingFinding{
			CapabilityArtifactRecordID: artifact.CommonNode.ID,
			ArtifactID:                 artifact.ArtifactID,
			ArtifactType:               artifact.ArtifactType,
			Name:                       artifact.Name,
			Reason:                     "usage_logging_required is false or missing; material capability artifacts must set usage_logging_required=true before use",
		})
	}
	return findings
}

func (s *InMemoryStore) capabilityVersionHasRollbackTarget(version *CapabilityVersion) bool {
	return version != nil && version.RollbackTo != nil && *version.RollbackTo != ""
}

func (s *InMemoryStore) capabilityVersionHasActiveRollback(versionID string) bool {
	for _, record := range s.ByType(TypeRollbackRecord) {
		rollback := record.(*RollbackRecord)
		if rollback.CapabilityVersionID != versionID {
			continue
		}
		if rollback.CommonNode.Status != nil && (*rollback.CommonNode.Status == "planned" || *rollback.CommonNode.Status == "completed") {
			return true
		}
	}
	return false
}

func (s *InMemoryStore) findCapabilityArtifact(ref string) (*CapabilityArtifact, bool) {
	if artifact, ok := s.mustGetCapabilityArtifact(ref); ok {
		return artifact, true
	}
	for _, record := range s.ByType(TypeCapabilityArtifact) {
		artifact := record.(*CapabilityArtifact)
		if artifact.ArtifactID == ref {
			return artifact, true
		}
	}
	return nil, false
}

func (s *InMemoryStore) findCapabilityArtifactForSourceRef(ref capabilitySourceRef) (*CapabilityArtifact, bool) {
	for _, record := range s.ByType(TypeCapabilityArtifact) {
		artifact := record.(*CapabilityArtifact)
		if ref.artifactType != "" && artifact.ArtifactType != ref.artifactType {
			continue
		}
		if artifact.CommonNode.ID == ref.value || artifact.ArtifactID == ref.value || artifact.CommonNode.ID == ref.raw || artifact.ArtifactID == ref.raw {
			return artifact, true
		}
	}
	return nil, false
}

func (s *InMemoryStore) capabilitySourceRefs(sourceRefs []string) []capabilitySourceRef {
	types := []string{"skill", "plugin", "prompt_section", "tool_description", "workflow_pack", "schema_instruction", "evaluation_prompt", "runtime_adapter", "policy_bundle"}
	var out []capabilitySourceRef
	seen := map[string]bool{}
	for _, sourceRef := range sourceRefs {
		if artifact, ok := s.findCapabilityArtifact(sourceRef); ok {
			key := "bare:" + artifact.CommonNode.ID
			if !seen[key] {
				out = append(out, capabilitySourceRef{raw: sourceRef, artifactType: artifact.ArtifactType, value: artifact.CommonNode.ID})
				seen[key] = true
			}
			continue
		}
		artifactType, value, ok := parseCapabilitySourceRef(sourceRef, types)
		if !ok {
			continue
		}
		key := artifactType + ":" + value
		if seen[key] {
			continue
		}
		out = append(out, capabilitySourceRef{raw: sourceRef, artifactType: artifactType, value: value})
		seen[key] = true
	}
	return out
}

func parseCapabilitySourceRef(sourceRef string, types []string) (string, string, bool) {
	sourceRef = strings.TrimSpace(sourceRef)
	if sourceRef == "" {
		return "", "", false
	}
	lower := strings.ToLower(sourceRef)
	if strings.HasPrefix(lower, "capability:") {
		return "", trimCapabilitySourceRef(sourceRef[len("capability:"):]), true
	}
	if strings.HasPrefix(lower, "capability//") {
		return "", trimCapabilitySourceRef(sourceRef[len("capability//"):]), true
	}
	for _, artifactType := range types {
		for _, delimiter := range []string{":", "//"} {
			prefix := artifactType + delimiter
			if strings.HasPrefix(lower, prefix) {
				return artifactType, trimCapabilitySourceRef(sourceRef[len(prefix):]), true
			}
		}
	}
	return "", "", false
}

func trimCapabilitySourceRef(value string) string {
	return strings.TrimPrefix(value, "//")
}
