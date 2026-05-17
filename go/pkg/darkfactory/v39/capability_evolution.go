package v39

import (
	"fmt"
	"strings"
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
	return s.AppendEdge(derivedEdge(EdgeUsedCapability, sourceID, artifact.CommonNode.ID, common))
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
	path.Missing = append(path.Missing, orderPath.Missing...)

	for _, taskID := range taskIDsFromPath(s, orderPath) {
		task, ok := s.mustGetTask(taskID)
		if !ok {
			path.Missing = append(path.Missing, "Task "+taskID)
			continue
		}
		for _, sourceRef := range capabilitySourceRefs(task.SourceRefs) {
			artifact, ok := s.findCapabilityArtifactForSourceRef(sourceRef)
			if !ok {
				path.Missing = append(path.Missing, "CapabilityArtifact for "+sourceRef.raw+" used in Task "+taskID)
				continue
			}
			path.NodeIDs = appendUniqueStrings(path.NodeIDs, taskID, artifact.CommonNode.ID)
			edgeIDs := edgeIDsBetween(s, taskID, artifact.CommonNode.ID, EdgeUsedCapability)
			if len(edgeIDs) == 0 {
				path.Missing = append(path.Missing, "USED_CAPABILITY from Task "+taskID+" to CapabilityArtifact "+artifact.CommonNode.ID)
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
		if _, err := s.AppendEdge(derivedEdge(edgeSpec.typ, capabilityVersion.CommonNode.ID, edgeSpec.to, capabilityVersion.CommonNode)); err != nil {
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
	if _, err := s.AppendEdge(derivedEdge(EdgeActivatedBy, version.CommonNode.ID, activationPolicy.CommonNode.ID, activationPolicy.CommonNode)); err != nil {
		return nil, nil, err
	}
	if _, err := s.AppendEdge(derivedEdge(EdgePackagedAs, version.CommonNode.ID, storedRuntime.CommonNode.ID, storedRuntime.CommonNode)); err != nil {
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
	if _, err := s.AppendEdge(derivedEdge(EdgeRolledBackTo, rollback.CapabilityVersionID, rollback.RollbackTo, rollback.CommonNode)); err != nil {
		return nil, err
	}
	return rollback, nil
}

func (s *InMemoryStore) evaluateCapabilityPromotion(version *CapabilityVersion) (capabilityPromotionEvidence, error) {
	var evidence capabilityPromotionEvidence
	if err := version.Validate(); err != nil {
		return evidence, err
	}
	if version.CommonNode.Status == nil || *version.CommonNode.Status != "approved" {
		return evidence, fieldError(TypeCapabilityVersion, "status", "must be approved for promotion")
	}
	if version.RollbackTo == nil || *version.RollbackTo == "" {
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

	candidate, ok := s.findCandidateVariantForArtifact(artifact)
	if !ok {
		return evidence, fmt.Errorf("%w: CandidateVariant for CapabilityArtifact %s", ErrRequiredPathMissing, artifact.CommonNode.ID)
	}
	if candidate.CommonNode.Status == nil || *candidate.CommonNode.Status != "approved" {
		return evidence, fieldError(TypeCandidateVariant, "status", "approved candidate required")
	}
	evidence.candidate = candidate

	run, ok := s.mustGetOptimizationRun(candidate.OptimizationRunID)
	if !ok {
		return evidence, fmt.Errorf("%w: OptimizationRun %s", ErrNotFound, candidate.OptimizationRunID)
	}
	if run.CommonNode.Status == nil || *run.CommonNode.Status != "succeeded" {
		return evidence, fieldError(TypeOptimizationRun, "status", "succeeded optimization run required")
	}
	if run.CommonNode.CreatedBy == version.CommonNode.CreatedBy {
		return evidence, fieldError(TypeCapabilityVersion, "created_by", "optimizer actor cannot promote its own output")
	}
	evidence.run = run

	order, ok := s.mustGetEvolutionOrder(run.EvolutionOrderID)
	if !ok {
		return evidence, fmt.Errorf("%w: EvolutionOrder %s", ErrNotFound, run.EvolutionOrderID)
	}
	if order.CommonNode.Status == nil || *order.CommonNode.Status != "accepted" {
		return evidence, fieldError(TypeEvolutionOrder, "status", "accepted evolution order required")
	}
	evidence.order = order

	dataset, ok := s.mustGetEvalDataset(run.EvalDatasetID)
	if !ok {
		return evidence, fmt.Errorf("%w: EvalDataset %s", ErrNotFound, run.EvalDatasetID)
	}
	evidence.dataset = dataset

	benchmark, ok := s.findBenchmarkResultForCandidate(candidate.CommonNode.ID)
	if !ok {
		return evidence, fmt.Errorf("%w: BenchmarkResult for CandidateVariant %s", ErrRequiredPathMissing, candidate.CommonNode.ID)
	}
	if benchmark.CommonNode.Status == nil || *benchmark.CommonNode.Status != "pass" {
		return evidence, fieldError(TypeBenchmarkResult, "status", "benchmark regression blocks capability promotion")
	}
	evidence.benchmark = benchmark

	review, ok := s.findApprovedHumanReview(artifact, candidate)
	if !ok {
		return evidence, fmt.Errorf("%w: approved HumanReview for CapabilityArtifact %s", ErrRequiredPathMissing, artifact.CommonNode.ID)
	}
	evidence.review = review
	return evidence, nil
}

func (s *InMemoryStore) findCapabilityArtifact(ref string) (*CapabilityArtifact, bool) {
	if artifact, ok := s.mustGetCapabilityArtifact(ref); ok {
		return artifact, true
	}
	for _, record := range s.ByType(TypeCapabilityArtifact) {
		artifact := record.(*CapabilityArtifact)
		if artifact.ArtifactID == ref || artifact.Name == ref {
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
		if artifact.CommonNode.ID == ref.value || artifact.ArtifactID == ref.value || artifact.Name == ref.value || artifact.CommonNode.ID == ref.raw || artifact.ArtifactID == ref.raw {
			return artifact, true
		}
	}
	return nil, false
}

func (s *InMemoryStore) findCandidateVariantForArtifact(artifact *CapabilityArtifact) (*CandidateVariant, bool) {
	for _, record := range s.ByType(TypeCandidateVariant) {
		candidate := record.(*CandidateVariant)
		if candidate.CapabilityArtifactID == artifact.CommonNode.ID || candidate.CapabilityArtifactID == artifact.ArtifactID {
			return candidate, true
		}
	}
	return nil, false
}

func (s *InMemoryStore) findBenchmarkResultForCandidate(candidateID string) (*BenchmarkResult, bool) {
	for _, record := range s.ByType(TypeBenchmarkResult) {
		benchmark := record.(*BenchmarkResult)
		if benchmark.CandidateVariantID == candidateID {
			return benchmark, true
		}
	}
	return nil, false
}

func (s *InMemoryStore) findApprovedHumanReview(artifact *CapabilityArtifact, candidate *CandidateVariant) (*HumanReview, bool) {
	if review, ok := s.mustGetHumanReview(artifact.HumanReviewRef); ok && review.CommonNode.Status != nil && *review.CommonNode.Status == "approved" {
		return review, true
	}
	for _, edge := range s.outgoingEdges(candidate.CommonNode.ID, EdgeReviewedBy) {
		if review, ok := s.mustGetHumanReview(edge.ToID); ok && review.CommonNode.Status != nil && *review.CommonNode.Status == "approved" {
			return review, true
		}
	}
	for _, record := range s.ByType(TypeHumanReview) {
		review := record.(*HumanReview)
		if review.CommonNode.Status == nil || *review.CommonNode.Status != "approved" {
			continue
		}
		if containsString(review.CommonNode.SourceRefs, candidate.CommonNode.ID) || containsString(review.CommonNode.SourceRefs, artifact.CommonNode.ID) || containsString(review.CommonNode.SourceRefs, artifact.ArtifactID) {
			return review, true
		}
	}
	return nil, false
}

func capabilitySourceRefs(sourceRefs []string) []capabilitySourceRef {
	types := []string{"skill", "plugin", "prompt_section", "tool_description", "workflow_pack", "schema_instruction", "evaluation_prompt", "runtime_adapter", "policy_bundle"}
	var out []capabilitySourceRef
	for _, sourceRef := range sourceRefs {
		if strings.HasPrefix(sourceRef, "capability:") {
			out = append(out, capabilitySourceRef{raw: sourceRef, value: trimCapabilitySourceRef(sourceRef, "capability")})
			continue
		}
		for _, artifactType := range types {
			if strings.HasPrefix(sourceRef, artifactType+":") {
				out = append(out, capabilitySourceRef{raw: sourceRef, artifactType: artifactType, value: trimCapabilitySourceRef(sourceRef, artifactType)})
				break
			}
		}
	}
	return out
}

func trimCapabilitySourceRef(sourceRef, prefix string) string {
	value := strings.TrimPrefix(sourceRef, prefix+":")
	value = strings.TrimPrefix(value, "//")
	return value
}
