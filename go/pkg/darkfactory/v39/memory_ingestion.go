package v39

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

const (
	MemoryClassificationPublic       = "public"
	MemoryClassificationInternal     = "internal"
	MemoryClassificationConfidential = "confidential"
	MemoryClassificationSecret       = "secret"

	MemoryRedactionNone          = "none"
	MemoryRedactionRedacted      = "redacted"
	MemoryRedactionQuarantined   = "quarantined"
	MemoryRedactionReferenceOnly = "reference_only"

	MemoryInfluencePlanning            = "planning"
	MemoryInfluenceRepair              = "repair"
	MemoryInfluenceReview              = "review"
	MemoryInfluenceAudit               = "audit"
	MemoryInfluenceCapabilityEvolution = "capability_evolution"
)

type MemoryIngestionInput struct {
	MemoryID                     string
	SourceSystem                 string
	SourceRef                    string
	SourceHashOrImmutableLocator string
	Content                      string
	ContentHash                  string
	Classification               string
	Scope                        string
	Summary                      string
	CreatedAt                    time.Time
	CreatedBy                    string
	CorrelationID                string
	SecretIndicators             []string
}

type MemoryIngestionResult struct {
	MemoryID       string
	Classification string
	RedactionState string
	Scope          string
	Ingested       *MemoryIngested
	Indexed        *MemoryIndexed
	ScopeAssigned  *MemoryScopeAssigned
	Redaction      *MemoryRedactionApplied
}

type MemoryRetrieval struct {
	MemoryID                     string
	SourceSystem                 string
	SourceRef                    string
	SourceHashOrImmutableLocator string
	Classification               string
	Scope                        string
	Summary                      string
	RetrievedAt                  time.Time
	TrustLevel                   string
	FreshnessStatus              string
	RedactionState               string
}

func (s *InMemoryStore) IngestMemory(input MemoryIngestionInput) (MemoryIngestionResult, error) {
	if input.MemoryID == "" {
		return MemoryIngestionResult{}, fieldError(TypeMemoryIngested, "memory_id", "required")
	}
	if input.SourceSystem == "" {
		return MemoryIngestionResult{}, fieldError(TypeMemoryIngested, "source_system", "required")
	}
	if input.SourceRef == "" {
		return MemoryIngestionResult{}, fieldError(TypeMemoryIngested, "source_ref", "required")
	}
	if input.Scope == "" {
		return MemoryIngestionResult{}, fieldError(TypeMemoryScopeAssigned, "scope", "required")
	}
	if input.Content == "" && input.ContentHash == "" {
		return MemoryIngestionResult{}, fieldError(TypeMemoryIngested, "content_hash", "required when content is not provided")
	}
	createdAt := UTC(input.CreatedAt)
	if createdAt.IsZero() {
		createdAt = UTC(time.Now())
	}
	createdBy := input.CreatedBy
	if createdBy == "" {
		createdBy = "memory_ingestion_adapter"
	}
	correlationID := input.CorrelationID
	if correlationID == "" {
		correlationID = "memory:" + input.MemoryID
	}

	classification, indicators := classifyMemorySource(input)
	contentHash := input.ContentHash
	if contentHash == "" {
		contentHash = hashMemoryContent(input.Content)
	}
	sourceHash := input.SourceHashOrImmutableLocator
	if sourceHash == "" {
		sourceHash = contentHash
	}
	redactionState := redactionStateForClassification(classification)
	summary := input.Summary
	if summary == "" {
		summary = "memory indexed from " + input.SourceRef
	}
	if redactionState == MemoryRedactionQuarantined {
		summary = "quarantined reference only for " + input.SourceRef
	}

	ingested := &MemoryIngested{
		CommonNode:                   memoryEventCommon(input.MemoryID, TypeMemoryIngested, createdAt, createdBy, correlationID),
		MemoryID:                     input.MemoryID,
		SourceSystem:                 input.SourceSystem,
		SourceRef:                    input.SourceRef,
		SourceHashOrImmutableLocator: sourceHash,
		Classification:               classification,
		ContentHash:                  contentHash,
		SecretIndicators:             indicators,
		IngestedAt:                   createdAt,
	}
	indexed := &MemoryIndexed{
		CommonNode:                    memoryEventCommon(input.MemoryID, TypeMemoryIndexed, createdAt, createdBy, correlationID),
		MemoryID:                      input.MemoryID,
		IndexID:                       "idx:" + input.MemoryID,
		IndexRef:                      "memory-index://" + input.MemoryID,
		ContentHashOrImmutableLocator: sourceHash,
		Summary:                       summary,
		IndexedAt:                     createdAt,
	}
	scopeAssigned := &MemoryScopeAssigned{
		CommonNode: memoryEventCommon(input.MemoryID, TypeMemoryScopeAssigned, createdAt, createdBy, correlationID),
		MemoryID:   input.MemoryID,
		Scope:      input.Scope,
		AssignedBy: createdBy,
		Reason:     "assigned during memory ingestion",
	}
	redaction := &MemoryRedactionApplied{
		CommonNode:       memoryEventCommon(input.MemoryID, TypeMemoryRedactionApplied, createdAt, createdBy, correlationID),
		MemoryID:         input.MemoryID,
		RedactionState:   redactionState,
		SecretIndicators: indicators,
	}
	if redactionState == MemoryRedactionRedacted {
		ref := "memory-redacted://" + input.MemoryID
		redaction.RedactedContentRef = &ref
	}
	if redactionState == MemoryRedactionQuarantined {
		reason := "secret-bearing memory is quarantined"
		redaction.QuarantineReason = &reason
	}

	for _, record := range []Record{ingested, indexed, scopeAssigned, redaction} {
		if _, err := s.AppendRecord(record); err != nil {
			return MemoryIngestionResult{}, err
		}
	}
	for _, toID := range []string{indexed.CommonNode.ID, scopeAssigned.CommonNode.ID, redaction.CommonNode.ID} {
		if _, err := s.AppendEdge(derivedEdge(EdgeDerivedFrom, ingested.CommonNode.ID, toID, ingested.CommonNode)); err != nil {
			return MemoryIngestionResult{}, err
		}
	}

	return MemoryIngestionResult{
		MemoryID:       input.MemoryID,
		Classification: classification,
		RedactionState: redactionState,
		Scope:          input.Scope,
		Ingested:       ingested,
		Indexed:        indexed,
		ScopeAssigned:  scopeAssigned,
		Redaction:      redaction,
	}, nil
}

func (s *InMemoryStore) RetrieveMemory(scope string, retrievedAt time.Time) []MemoryRetrieval {
	if scope == "" {
		return nil
	}
	retrievedAt = UTC(retrievedAt)
	if retrievedAt.IsZero() {
		retrievedAt = UTC(time.Now())
	}

	scopes := map[string]*MemoryScopeAssigned{}
	for _, record := range s.ByType(TypeMemoryScopeAssigned) {
		scopeRecord := record.(*MemoryScopeAssigned)
		scopes[scopeRecord.MemoryID] = scopeRecord
	}
	redactions := map[string]*MemoryRedactionApplied{}
	for _, record := range s.ByType(TypeMemoryRedactionApplied) {
		redaction := record.(*MemoryRedactionApplied)
		redactions[redaction.MemoryID] = redaction
	}
	indexed := map[string]*MemoryIndexed{}
	for _, record := range s.ByType(TypeMemoryIndexed) {
		indexRecord := record.(*MemoryIndexed)
		indexed[indexRecord.MemoryID] = indexRecord
	}

	var out []MemoryRetrieval
	for _, record := range s.ByType(TypeMemoryIngested) {
		ingested := record.(*MemoryIngested)
		scopeRecord, ok := scopes[ingested.MemoryID]
		if !ok || scopeRecord.Scope != scope {
			continue
		}
		redaction, ok := redactions[ingested.MemoryID]
		if !ok || redaction.RedactionState == MemoryRedactionQuarantined {
			continue
		}
		indexRecord, ok := indexed[ingested.MemoryID]
		if !ok {
			continue
		}
		out = append(out, MemoryRetrieval{
			MemoryID:                     ingested.MemoryID,
			SourceSystem:                 ingested.SourceSystem,
			SourceRef:                    ingested.SourceRef,
			SourceHashOrImmutableLocator: ingested.SourceHashOrImmutableLocator,
			Classification:               ingested.Classification,
			Scope:                        scopeRecord.Scope,
			Summary:                      indexRecord.Summary,
			RetrievedAt:                  retrievedAt,
			TrustLevel:                   "system_observed",
			FreshnessStatus:              "current",
			RedactionState:               redaction.RedactionState,
		})
	}
	return out
}

func (s *InMemoryStore) RecordRetrievedMemoryInfluence(retrieval MemoryRetrieval, taskID, actorID, influenceType, influenceSummary, riskScope string) (*MemoryReference, error) {
	if !isMaterialMemoryInfluence(influenceType) {
		return nil, fieldError(TypeMemoryReference, "influence_summary", "memory reference requires material influence in planning, repair, review, audit, or capability_evolution")
	}
	if retrieval.MemoryID == "" {
		return nil, fieldError(TypeMemoryReference, "source_ref", "retrieved memory required")
	}
	if retrieval.RedactionState == MemoryRedactionQuarantined {
		return nil, fieldError(TypeMemoryReference, "redaction_state", "quarantined memory cannot be used")
	}
	if influenceSummary == "" {
		return nil, fieldError(TypeMemoryReference, "influence_summary", "required")
	}
	if actorID == "" {
		return nil, fieldError(TypeMemoryReference, "used_by_actor", "required")
	}
	if riskScope == "" {
		riskScope = "medium"
	}

	ref := &MemoryReference{AdvisoryReference: AdvisoryReference{
		CommonNode:                   memoryInfluenceCommon(retrieval.MemoryID, taskID, influenceType, actorID, retrieval.RetrievedAt),
		ReferenceCreatedAt:           retrieval.RetrievedAt,
		SourceSystem:                 retrieval.SourceSystem,
		SourceRef:                    retrieval.SourceRef,
		SourceHashOrImmutableLocator: retrieval.SourceHashOrImmutableLocator,
		RetrievedAt:                  retrieval.RetrievedAt,
		UsedByActor:                  actorID,
		UsedInTask:                   taskID,
		InfluenceSummary:             influenceType + ": " + influenceSummary,
		RiskScope:                    riskScope,
		TrustLevel:                   retrieval.TrustLevel,
		FreshnessStatus:              retrieval.FreshnessStatus,
		RedactionState:               retrieval.RedactionState,
		ContradictionRefs:            []string{},
	}}
	return s.RecordMemoryReference(ref)
}

func (s *InMemoryStore) validateRecordRelations(r Record) error {
	switch typed := r.(type) {
	case *MemoryReference:
		return s.validateAdvisoryReferenceUse(&typed.AdvisoryReference, TypeMemoryReference)
	case *KnowledgeReference:
		return s.validateAdvisoryReferenceUse(&typed.AdvisoryReference, TypeKnowledgeReference)
	case *MemoryIndexed:
		if _, ok := s.memoryIngestedByMemoryID(typed.MemoryID); !ok {
			return fmt.Errorf("%w: MemoryIngested %s", ErrNotFound, typed.MemoryID)
		}
	case *MemoryScopeAssigned:
		if _, ok := s.memoryIngestedByMemoryID(typed.MemoryID); !ok {
			return fmt.Errorf("%w: MemoryIngested %s", ErrNotFound, typed.MemoryID)
		}
	case *MemoryRedactionApplied:
		if _, ok := s.memoryIngestedByMemoryID(typed.MemoryID); !ok {
			return fmt.Errorf("%w: MemoryIngested %s", ErrNotFound, typed.MemoryID)
		}
	case *CivilizationAssemblyProjectionStoreRecord:
		decision, ok := s.mustGetAuthorityDecision(typed.AuthorityDecisionRef)
		if !ok {
			return fmt.Errorf("%w: AuthorityDecision %s", ErrNotFound, typed.AuthorityDecisionRef)
		}
		if !civilizationAssemblyProjectionStoreDecisionAuthorizes(decision) {
			return fieldError(TypeAuthorityDecision, "scope", "must include "+CivilizationAssemblyProjectionStoreAction)
		}
	case *OptimizationRun:
		if _, ok := s.mustGetEvolutionOrder(typed.EvolutionOrderID); !ok {
			return fmt.Errorf("%w: EvolutionOrder %s", ErrNotFound, typed.EvolutionOrderID)
		}
		if _, ok := s.mustGetEvalDataset(typed.EvalDatasetID); !ok {
			return fmt.Errorf("%w: EvalDataset %s", ErrNotFound, typed.EvalDatasetID)
		}
	case *CandidateVariant:
		if _, ok := s.mustGetOptimizationRun(typed.OptimizationRunID); !ok {
			return fmt.Errorf("%w: OptimizationRun %s", ErrNotFound, typed.OptimizationRunID)
		}
		if _, ok := s.findCapabilityArtifact(typed.CapabilityArtifactID); !ok {
			return fmt.Errorf("%w: CapabilityArtifact %s", ErrNotFound, typed.CapabilityArtifactID)
		}
	case *BenchmarkResult:
		if _, ok := s.mustGetCandidateVariant(typed.CandidateVariantID); !ok {
			return fmt.Errorf("%w: CandidateVariant %s", ErrNotFound, typed.CandidateVariantID)
		}
	case *CapabilityVersion:
		if _, ok := s.findCapabilityArtifact(typed.CapabilityArtifactID); !ok {
			return fmt.Errorf("%w: CapabilityArtifact %s", ErrNotFound, typed.CapabilityArtifactID)
		}
		for field, ref := range map[string]string{
			"evolution_order_id":   typed.EvolutionOrderID,
			"optimization_run_id":  typed.OptimizationRunID,
			"candidate_variant_id": typed.CandidateVariantID,
			"eval_dataset_id":      typed.EvalDatasetID,
			"benchmark_result_id":  typed.BenchmarkResultID,
			"human_review_id":      typed.HumanReviewID,
		} {
			if ref == "" {
				continue
			}
			if !s.recordExistsForCapabilityVersionEvidence(field, ref) {
				return fmt.Errorf("%w: %s %s", ErrNotFound, field, ref)
			}
		}
		if typed.RollbackTo != nil && *typed.RollbackTo != "" {
			if _, ok := s.mustGetCapabilityVersion(*typed.RollbackTo); !ok {
				return fmt.Errorf("%w: rollback CapabilityVersion %s", ErrNotFound, *typed.RollbackTo)
			}
		}
	case *ActivationPolicy:
		version, ok := s.mustGetCapabilityVersion(typed.CapabilityVersionID)
		if !ok {
			return fmt.Errorf("%w: CapabilityVersion %s", ErrNotFound, typed.CapabilityVersionID)
		}
		if typed.Scope != "disabled" && !s.capabilityVersionHasRollbackTarget(version) {
			return fieldError(TypeCapabilityVersion, "rollback_to", "required before activation policy")
		}
		if s.capabilityVersionHasActiveRollback(version.CommonNode.ID) {
			return fieldError(TypeCapabilityVersion, "status", "rolled back capability version cannot be activated")
		}
	case *RollbackRecord:
		if _, ok := s.mustGetCapabilityVersion(typed.CapabilityVersionID); !ok {
			return fmt.Errorf("%w: CapabilityVersion %s", ErrNotFound, typed.CapabilityVersionID)
		}
		if _, ok := s.mustGetCapabilityVersion(typed.RollbackTo); !ok {
			return fmt.Errorf("%w: rollback CapabilityVersion %s", ErrNotFound, typed.RollbackTo)
		}
		if _, ok := s.mustGetFactoryRuntimeVersion(typed.FactoryRuntimeVersionID); !ok {
			return fmt.Errorf("%w: FactoryRuntimeVersion %s", ErrNotFound, typed.FactoryRuntimeVersionID)
		}
	case *FactoryRuntimeVersion:
		for _, capabilityVersionRef := range typed.CapabilityVersionRefs {
			version, ok := s.mustGetCapabilityVersion(capabilityVersionRef)
			if !ok {
				return fmt.Errorf("%w: CapabilityVersion %s", ErrNotFound, capabilityVersionRef)
			}
			if version.CommonNode.Status != nil {
				switch *version.CommonNode.Status {
				case "rolled_back", "rejected", "superseded":
					return fieldError(TypeFactoryRuntimeVersion, "capability_version_refs", "must not reference inactive capability version "+capabilityVersionRef)
				}
			}
			if !s.capabilityVersionHasRollbackTarget(version) {
				return fieldError(TypeCapabilityVersion, "rollback_to", "required before runtime packaging")
			}
			if s.capabilityVersionHasActiveRollback(version.CommonNode.ID) {
				return fieldError(TypeFactoryRuntimeVersion, "capability_version_refs", "rolled back capability version cannot be packaged")
			}
			if !s.capabilityVersionHasPromotionEvidence(version) {
				return fmt.Errorf("%w: promotion evidence for CapabilityVersion %s", ErrRequiredPathMissing, capabilityVersionRef)
			}
		}
	}
	return nil
}

func (s *InMemoryStore) recordExistsForCapabilityVersionEvidence(field, ref string) bool {
	switch field {
	case "evolution_order_id":
		_, ok := s.mustGetEvolutionOrder(ref)
		return ok
	case "optimization_run_id":
		_, ok := s.mustGetOptimizationRun(ref)
		return ok
	case "candidate_variant_id":
		_, ok := s.mustGetCandidateVariant(ref)
		return ok
	case "eval_dataset_id":
		_, ok := s.mustGetEvalDataset(ref)
		return ok
	case "benchmark_result_id":
		_, ok := s.mustGetBenchmarkResult(ref)
		return ok
	case "human_review_id":
		_, ok := s.mustGetHumanReview(ref)
		return ok
	default:
		return false
	}
}

func (s *InMemoryStore) memoryIngestedByMemoryID(memoryID string) (*MemoryIngested, bool) {
	for _, record := range s.ByType(TypeMemoryIngested) {
		ingested := record.(*MemoryIngested)
		if ingested.MemoryID == memoryID {
			return ingested, true
		}
	}
	return nil, false
}

func classifyMemorySource(input MemoryIngestionInput) (string, []string) {
	indicators := append([]string{}, input.SecretIndicators...)
	indicators = append(indicators, detectedSecretIndicators(input.Content)...)
	classification := strings.ToLower(strings.TrimSpace(input.Classification))
	if classification == "" {
		classification = MemoryClassificationInternal
	}
	if len(indicators) > 0 || classification == MemoryClassificationSecret {
		classification = MemoryClassificationSecret
	}
	return classification, appendUniqueStrings(nil, indicators...)
}

func detectedSecretIndicators(content string) []string {
	lower := strings.ToLower(content)
	checks := []struct {
		needle    string
		indicator string
	}{
		{needle: "password", indicator: "password"},
		{needle: "api_key", indicator: "api_key"},
		{needle: "secret", indicator: "secret"},
		{needle: "private key", indicator: "private_key"},
		{needle: "token=", indicator: "token"},
	}
	var indicators []string
	for _, check := range checks {
		if strings.Contains(lower, check.needle) {
			indicators = append(indicators, check.indicator)
		}
	}
	return indicators
}

func redactionStateForClassification(classification string) string {
	switch classification {
	case MemoryClassificationSecret:
		return MemoryRedactionQuarantined
	case MemoryClassificationConfidential:
		return MemoryRedactionRedacted
	default:
		return MemoryRedactionNone
	}
}

func isMaterialMemoryInfluence(influenceType string) bool {
	switch influenceType {
	case MemoryInfluencePlanning, MemoryInfluenceRepair, MemoryInfluenceReview, MemoryInfluenceAudit, MemoryInfluenceCapabilityEvolution:
		return true
	default:
		return false
	}
}

func materialInfluenceTypeFromSummary(summary string) string {
	influenceType, _, ok := strings.Cut(summary, ":")
	if !ok {
		return ""
	}
	return strings.TrimSpace(influenceType)
}

func memoryEventCommon(memoryID, typ string, createdAt time.Time, createdBy, correlationID string) CommonNode {
	id := "memory:" + memoryID + ":" + typ
	return CommonNode{
		ID:             id,
		Type:           typ,
		CreatedAt:      createdAt,
		CreatedBy:      createdBy,
		Status:         stringPtr("recorded"),
		IdempotencyKey: id,
		CorrelationID:  correlationID,
	}
}

func memoryInfluenceCommon(memoryID, taskID, influenceType, actorID string, createdAt time.Time) CommonNode {
	id := "memory:" + memoryID + ":reference:" + taskID + ":" + influenceType
	return CommonNode{
		ID:             id,
		Type:           TypeMemoryReference,
		CreatedAt:      createdAt,
		CreatedBy:      actorID,
		Status:         stringPtr("recorded"),
		IdempotencyKey: id,
		CorrelationID:  "memory:" + memoryID,
	}
}

func stringPtr(s string) *string {
	return &s
}

func hashMemoryContent(content string) string {
	sum := sha256.Sum256([]byte(content))
	return "sha256:" + hex.EncodeToString(sum[:])
}
