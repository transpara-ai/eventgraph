package v39

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
	"time"
)

const (
	CivilizationAssemblyProjectionSchemaVersion = "1.0.0"
	CivilizationAssemblyProjectionSubject       = "civilization_assembly"

	CivilizationAssemblyDerivationComplete    CivilizationAssemblyDerivationStatus = "complete"
	CivilizationAssemblyDerivationPartial     CivilizationAssemblyDerivationStatus = "partial"
	CivilizationAssemblyDerivationUnavailable CivilizationAssemblyDerivationStatus = "unavailable"
	CivilizationAssemblyDerivationFailed      CivilizationAssemblyDerivationStatus = "failed"

	CivilizationAssemblyFieldAvailable   CivilizationAssemblyFieldAvailability = "available"
	CivilizationAssemblyFieldUnavailable CivilizationAssemblyFieldAvailability = "unavailable"

	civilizationAssemblySiteConsumerSourceRef  = "civilization_assembly:site_consumer"
	civilizationAssemblySiteConsumerPathPrefix = "site://ops/civilization/"
)

type CivilizationAssemblyDerivationStatus string
type CivilizationAssemblyFieldAvailability string

type CivilizationAssemblyProjectionOptions struct {
	ProjectionID   string
	GeneratedAt    time.Time
	ValidationRefs []string
	BoundaryFlags  []string
}

type CivilizationAssemblyProjection struct {
	ProjectionID                       string                                 `json:"projection_id"`
	ProjectionSchemaVersion            string                                 `json:"projection_schema_version"`
	ProjectionSubject                  string                                 `json:"projection_subject"`
	GeneratedAt                        time.Time                              `json:"generated_at"`
	SourceEventGraphHeadOrStateVersion string                                 `json:"source_eventgraph_head_or_state_version"`
	SourceEventIDsOrQueryWindow        []string                               `json:"source_event_ids_or_query_window"`
	DerivationStatus                   CivilizationAssemblyDerivationStatus   `json:"derivation_status"`
	AuthorityState                     CivilizationAssemblyAuthorityState     `json:"authority_state"`
	ExternalCommitteeState             CivilizationAssemblyCommitteeState     `json:"external_committee_state"`
	ActorRoster                        []CivilizationAssemblyActorSummary     `json:"actor_roster"`
	RoleBindings                       []CivilizationAssemblyRoleBinding      `json:"role_bindings"`
	AgentLifecycleSummary              []CivilizationAssemblyLifecycleSummary `json:"agent_lifecycle_summary"`
	FactoryOrderSummary                []CivilizationAssemblyFactoryOrder     `json:"factory_order_summary"`
	WorkEvidenceSummary                CivilizationAssemblyWorkEvidence       `json:"work_evidence_summary"`
	SiteConsumerStatus                 CivilizationAssemblyFieldStatus        `json:"site_consumer_status"`
	OpenGateSummary                    []CivilizationAssemblyGateSummary      `json:"open_gate_summary"`
	ResidualRiskSummary                []CivilizationAssemblyResidualRisk     `json:"residual_risk_summary"`
	WithheldOrUnavailableFields        []CivilizationAssemblyUnavailableField `json:"withheld_or_unavailable_fields"`
	BoundaryFlags                      []string                               `json:"boundary_flags"`
	ProvenanceRefs                     []string                               `json:"provenance_refs"`
	ValidationRefs                     []string                               `json:"validation_refs"`
	FailureReasons                     []string                               `json:"failure_reasons,omitempty"`
}

type CivilizationAssemblyFieldStatus struct {
	Status     CivilizationAssemblyFieldAvailability `json:"status"`
	Summary    string                                `json:"summary"`
	SourceRefs []string                              `json:"source_refs,omitempty"`
}

type CivilizationAssemblyUnavailableField struct {
	Field  string                                `json:"field"`
	Status CivilizationAssemblyFieldAvailability `json:"status"`
	Reason string                                `json:"reason"`
}

type CivilizationAssemblyAuthorityState struct {
	Status             CivilizationAssemblyFieldAvailability   `json:"status"`
	Summary            string                                  `json:"summary"`
	AuthorityRequests  []CivilizationAssemblyAuthorityRequest  `json:"authority_requests,omitempty"`
	AuthorityDecisions []CivilizationAssemblyAuthorityDecision `json:"authority_decisions,omitempty"`
	ExecutionReceipts  []CivilizationAssemblyExecutionReceipt  `json:"execution_receipts,omitempty"`
	SourceRefs         []string                                `json:"source_refs,omitempty"`
}

type CivilizationAssemblyAuthorityRequest struct {
	ID         string `json:"id"`
	ActorID    string `json:"actor_id"`
	ActorRole  string `json:"actor_role"`
	Action     string `json:"action"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
	RiskClass  string `json:"risk_class"`
	Status     string `json:"status,omitempty"`
}

type CivilizationAssemblyAuthorityDecision struct {
	ID                 string   `json:"id"`
	AuthorityRequestID string   `json:"authority_request_id"`
	DeciderActorID     string   `json:"decider_actor_id"`
	DeciderRole        string   `json:"decider_role"`
	Decision           string   `json:"decision"`
	Status             string   `json:"status,omitempty"`
	Scope              []string `json:"scope,omitempty"`
}

type CivilizationAssemblyExecutionReceipt struct {
	ID                  string `json:"id"`
	AuthorityDecisionID string `json:"authority_decision_id"`
	Action              string `json:"action"`
	TargetID            string `json:"target_id"`
	Result              string `json:"result"`
	Status              string `json:"status,omitempty"`
}

type CivilizationAssemblyCommitteeState struct {
	Status         CivilizationAssemblyFieldAvailability `json:"status"`
	Summary        string                                `json:"summary"`
	DecisionRefs   []string                              `json:"decision_refs,omitempty"`
	ApprovalRefs   []string                              `json:"approval_refs,omitempty"`
	CommitteeRoles []string                              `json:"committee_roles,omitempty"`
}

type CivilizationAssemblyActorSummary struct {
	ID           string `json:"id"`
	ActorID      string `json:"actor_id"`
	ActorType    string `json:"actor_type"`
	IdentityMode string `json:"identity_mode"`
	Status       string `json:"status,omitempty"`
}

type CivilizationAssemblyRoleBinding struct {
	ActorID    string `json:"actor_id"`
	Role       string `json:"role"`
	SourceRef  string `json:"source_ref"`
	SourceType string `json:"source_type"`
}

type CivilizationAssemblyLifecycleSummary struct {
	ID                  string  `json:"id"`
	ActorID             string  `json:"actor_id"`
	FromState           string  `json:"from_state,omitempty"`
	ToState             string  `json:"to_state,omitempty"`
	TrustLevel          string  `json:"trust_level,omitempty"`
	AuthorityDecisionID *string `json:"authority_decision_id,omitempty"`
	Status              string  `json:"status,omitempty"`
}

type CivilizationAssemblyFactoryOrder struct {
	ID                      string   `json:"id"`
	Status                  string   `json:"status,omitempty"`
	RiskClass               string   `json:"risk_class"`
	ReleasePolicy           string   `json:"release_policy"`
	RequirementRefs         []string `json:"requirement_refs,omitempty"`
	AcceptanceCriterionRefs []string `json:"acceptance_criterion_refs,omitempty"`
	TaskRefs                []string `json:"task_refs,omitempty"`
	ReleaseCandidateRefs    []string `json:"release_candidate_refs,omitempty"`
}

type CivilizationAssemblyWorkEvidence struct {
	Status          CivilizationAssemblyFieldAvailability `json:"status"`
	Summary         string                                `json:"summary"`
	TaskRefs        []string                              `json:"task_refs,omitempty"`
	ArtifactRefs    []string                              `json:"artifact_refs,omitempty"`
	TestRunRefs     []string                              `json:"test_run_refs,omitempty"`
	GateResultRefs  []string                              `json:"gate_result_refs,omitempty"`
	AuditReportRefs []string                              `json:"audit_report_refs,omitempty"`
	SourceRefs      []string                              `json:"source_refs,omitempty"`
}

type CivilizationAssemblyGateSummary struct {
	ID                 string   `json:"id"`
	GateName           string   `json:"gate_name"`
	Status             string   `json:"status,omitempty"`
	FactoryOrderID     string   `json:"factory_order_id,omitempty"`
	ReleaseCandidateID *string  `json:"release_candidate_id,omitempty"`
	EvidenceRefs       []string `json:"evidence_refs,omitempty"`
}

type CivilizationAssemblyResidualRisk struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Severity string `json:"severity,omitempty"`
	Status   string `json:"status,omitempty"`
	Summary  string `json:"summary"`
}

type civilizationAssemblySnapshot struct {
	records      []Record
	sourceIDs    []string
	stateVersion string
	failures     []string
}

func (s *InMemoryStore) ProjectCivilizationAssembly(options CivilizationAssemblyProjectionOptions) CivilizationAssemblyProjection {
	snapshot := s.civilizationAssemblySnapshot()
	generatedAt := options.GeneratedAt
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC()
	}

	projectionID := options.ProjectionID
	if projectionID == "" {
		projectionID = "civilization_assembly:" + strings.TrimPrefix(snapshot.stateVersion, "sha256:")[:16]
	}

	projection := CivilizationAssemblyProjection{
		ProjectionID:                       projectionID,
		ProjectionSchemaVersion:            CivilizationAssemblyProjectionSchemaVersion,
		ProjectionSubject:                  CivilizationAssemblyProjectionSubject,
		GeneratedAt:                        generatedAt.UTC(),
		SourceEventGraphHeadOrStateVersion: snapshot.stateVersion,
		SourceEventIDsOrQueryWindow:        append([]string(nil), snapshot.sourceIDs...),
		BoundaryFlags:                      civilizationAssemblyBoundaryFlags(options.BoundaryFlags),
		ValidationRefs:                     appendSortedUnique(nil, options.ValidationRefs...),
	}

	projection.AuthorityState = civilizationAssemblyAuthorityState(snapshot.records)
	projection.ExternalCommitteeState = civilizationAssemblyCommitteeState(snapshot.records)
	projection.ActorRoster = civilizationAssemblyActorRoster(snapshot.records)
	projection.RoleBindings = civilizationAssemblyRoleBindings(snapshot.records)
	projection.AgentLifecycleSummary = civilizationAssemblyLifecycleSummary(snapshot.records)
	projection.FactoryOrderSummary = civilizationAssemblyFactoryOrders(snapshot.records)
	projection.WorkEvidenceSummary = civilizationAssemblyWorkEvidence(snapshot.records)
	projection.SiteConsumerStatus = civilizationAssemblySiteConsumerStatus(snapshot.records)
	projection.OpenGateSummary = civilizationAssemblyOpenGates(snapshot.records)
	projection.ResidualRiskSummary = civilizationAssemblyResidualRisks(snapshot.records)
	projection.ValidationRefs = appendSortedUnique(projection.ValidationRefs, projection.WorkEvidenceSummary.TestRunRefs...)
	projection.ValidationRefs = appendSortedUnique(projection.ValidationRefs, projection.WorkEvidenceSummary.GateResultRefs...)
	projection.ProvenanceRefs = civilizationAssemblyProvenanceRefs(projection)
	projection.WithheldOrUnavailableFields = civilizationAssemblyUnavailableFields(projection)
	projection.FailureReasons = appendSortedUnique(snapshot.failures, civilizationAssemblyFailureReasons(snapshot.records)...)
	projection.DerivationStatus = civilizationAssemblyDerivationStatus(snapshot, projection)
	return projection
}

func (s *InMemoryStore) civilizationAssemblySnapshot() civilizationAssemblySnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	recordIDs := sortedMapKeysRecord(s.records)
	edgeIDs := sortedMapKeysEdge(s.edges)
	sourceIDs := make([]string, 0, len(recordIDs)+len(edgeIDs))
	records := make([]Record, 0, len(recordIDs))
	failures := []string{}
	hasher := sha256.New()

	for _, id := range recordIDs {
		sourceIDs = append(sourceIDs, id)
		hasher.Write([]byte("record:" + id + "\n"))
		hasher.Write(s.canonicalByID[id])
		hasher.Write([]byte("\n"))
		record := s.records[id]
		if record == nil {
			failures = append(failures, "source record "+id+" could not be cloned: nil record")
			continue
		}
		if cloned, err := cloneRecord(record); err == nil {
			records = append(records, cloned)
		} else {
			failures = append(failures, "source record "+id+" could not be cloned: "+err.Error())
		}
	}
	for _, id := range edgeIDs {
		sourceIDs = append(sourceIDs, id)
		hasher.Write([]byte("edge:" + id + "\n"))
		if canonical, err := CanonicalJSON(s.edges[id]); err == nil {
			hasher.Write(canonical)
		} else {
			failures = append(failures, "source edge "+id+" could not be canonicalized: "+err.Error())
		}
		hasher.Write([]byte("\n"))
	}

	return civilizationAssemblySnapshot{
		records:      records,
		sourceIDs:    sourceIDs,
		stateVersion: "sha256:" + hex.EncodeToString(hasher.Sum(nil)),
		failures:     appendSortedUnique(nil, failures...),
	}
}

func civilizationAssemblyAuthorityState(records []Record) CivilizationAssemblyAuthorityState {
	var state CivilizationAssemblyAuthorityState
	for _, record := range records {
		switch typed := record.(type) {
		case *AuthorityRequest:
			state.AuthorityRequests = append(state.AuthorityRequests, CivilizationAssemblyAuthorityRequest{
				ID:         typed.CommonNode.ID,
				ActorID:    typed.ActorID,
				ActorRole:  typed.ActorRole,
				Action:     typed.Action,
				TargetType: typed.TargetType,
				TargetID:   typed.TargetID,
				RiskClass:  typed.RiskClass,
				Status:     commonStatus(typed.CommonNode),
			})
			state.SourceRefs = append(state.SourceRefs, typed.CommonNode.ID)
		case *AuthorityDecision:
			state.AuthorityDecisions = append(state.AuthorityDecisions, CivilizationAssemblyAuthorityDecision{
				ID:                 typed.CommonNode.ID,
				AuthorityRequestID: typed.AuthorityRequestID,
				DeciderActorID:     typed.DeciderActorID,
				DeciderRole:        typed.DeciderRole,
				Decision:           typed.Decision,
				Status:             commonStatus(typed.CommonNode),
				Scope:              appendSortedUnique(nil, typed.Scope...),
			})
			state.SourceRefs = append(state.SourceRefs, typed.CommonNode.ID)
		case *ExecutionReceipt:
			state.ExecutionReceipts = append(state.ExecutionReceipts, CivilizationAssemblyExecutionReceipt{
				ID:                  typed.CommonNode.ID,
				AuthorityDecisionID: typed.AuthorityDecisionID,
				Action:              typed.Action,
				TargetID:            typed.TargetID,
				Result:              typed.Result,
				Status:              commonStatus(typed.CommonNode),
			})
			state.SourceRefs = append(state.SourceRefs, typed.CommonNode.ID)
		}
	}
	sort.Slice(state.AuthorityRequests, func(i, j int) bool { return state.AuthorityRequests[i].ID < state.AuthorityRequests[j].ID })
	sort.Slice(state.AuthorityDecisions, func(i, j int) bool { return state.AuthorityDecisions[i].ID < state.AuthorityDecisions[j].ID })
	sort.Slice(state.ExecutionReceipts, func(i, j int) bool { return state.ExecutionReceipts[i].ID < state.ExecutionReceipts[j].ID })
	state.SourceRefs = appendSortedUnique(nil, state.SourceRefs...)
	if len(state.AuthorityDecisions) == 0 {
		state.Status = CivilizationAssemblyFieldUnavailable
		state.Summary = "no AuthorityDecision records are present in the EventGraph source snapshot"
		return state
	}
	if len(authorityReferenceIntegrityReasons(state.AuthorityRequests, state.AuthorityDecisions, state.ExecutionReceipts)) > 0 {
		state.Status = CivilizationAssemblyFieldUnavailable
		state.Summary = "authority chain contains dangling AuthorityRequest, AuthorityDecision, or ExecutionReceipt references"
		return state
	}
	if hasConflictingAuthorityDecisions(state.AuthorityDecisions) {
		state.Status = CivilizationAssemblyFieldUnavailable
		state.Summary = "conflicting AuthorityDecision records are present for the same AuthorityRequest"
		return state
	}
	state.Status = CivilizationAssemblyFieldAvailable
	state.Summary = "authority state derived from AuthorityRequest, AuthorityDecision, and ExecutionReceipt records"
	return state
}

func civilizationAssemblyCommitteeState(records []Record) CivilizationAssemblyCommitteeState {
	state := CivilizationAssemblyCommitteeState{}
	roles := []string{}
	for _, record := range records {
		switch typed := record.(type) {
		case *AuthorityDecision:
			if typed.DeciderActorID != "" {
				state.DecisionRefs = append(state.DecisionRefs, typed.CommonNode.ID)
				roles = append(roles, typed.DeciderRole)
			}
		case *HumanApproval:
			state.ApprovalRefs = append(state.ApprovalRefs, typed.CommonNode.ID)
			roles = append(roles, typed.ApproverRole)
		}
	}
	state.DecisionRefs = appendSortedUnique(nil, state.DecisionRefs...)
	state.ApprovalRefs = appendSortedUnique(nil, state.ApprovalRefs...)
	state.CommitteeRoles = appendSortedUnique(nil, roles...)
	if len(state.DecisionRefs) == 0 && len(state.ApprovalRefs) == 0 {
		state.Status = CivilizationAssemblyFieldUnavailable
		state.Summary = "no External Committee decision or human approval records are present"
		return state
	}
	state.Status = CivilizationAssemblyFieldAvailable
	state.Summary = "External Committee state derived from human decision and approval records"
	return state
}

func civilizationAssemblyActorRoster(records []Record) []CivilizationAssemblyActorSummary {
	var out []CivilizationAssemblyActorSummary
	for _, record := range records {
		if typed, ok := record.(*ActorIdentity); ok {
			out = append(out, CivilizationAssemblyActorSummary{
				ID:           typed.CommonNode.ID,
				ActorID:      typed.ActorID,
				ActorType:    typed.ActorType,
				IdentityMode: typed.IdentityMode,
				Status:       commonStatus(typed.CommonNode),
			})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func civilizationAssemblyRoleBindings(records []Record) []CivilizationAssemblyRoleBinding {
	var out []CivilizationAssemblyRoleBinding
	for _, record := range records {
		switch typed := record.(type) {
		case *AuthorityRequest:
			if typed.ActorID != "" && typed.ActorRole != "" {
				out = append(out, CivilizationAssemblyRoleBinding{ActorID: typed.ActorID, Role: typed.ActorRole, SourceRef: typed.CommonNode.ID, SourceType: TypeAuthorityRequest})
			}
		case *AuthorityDecision:
			if typed.DeciderActorID != "" && typed.DeciderRole != "" {
				out = append(out, CivilizationAssemblyRoleBinding{ActorID: typed.DeciderActorID, Role: typed.DeciderRole, SourceRef: typed.CommonNode.ID, SourceType: TypeAuthorityDecision})
			}
		case *HumanApproval:
			if typed.ApproverActorID != "" && typed.ApproverRole != "" {
				out = append(out, CivilizationAssemblyRoleBinding{ActorID: typed.ApproverActorID, Role: typed.ApproverRole, SourceRef: typed.CommonNode.ID, SourceType: TypeHumanApproval})
			}
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ActorID == out[j].ActorID {
			if out[i].Role == out[j].Role {
				return out[i].SourceRef < out[j].SourceRef
			}
			return out[i].Role < out[j].Role
		}
		return out[i].ActorID < out[j].ActorID
	})
	return out
}

func civilizationAssemblyLifecycleSummary(records []Record) []CivilizationAssemblyLifecycleSummary {
	var out []CivilizationAssemblyLifecycleSummary
	for _, record := range records {
		switch typed := record.(type) {
		case *LifecycleTransition:
			out = append(out, CivilizationAssemblyLifecycleSummary{
				ID:                  typed.CommonNode.ID,
				ActorID:             typed.ActorID,
				FromState:           typed.FromState,
				ToState:             typed.ToState,
				AuthorityDecisionID: typed.AuthorityDecisionID,
				Status:              commonStatus(typed.CommonNode),
			})
		case *TrustRecord:
			out = append(out, CivilizationAssemblyLifecycleSummary{
				ID:         typed.CommonNode.ID,
				ActorID:    typed.SubjectActorID,
				TrustLevel: typed.TrustLevel,
				Status:     commonStatus(typed.CommonNode),
			})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func civilizationAssemblyFactoryOrders(records []Record) []CivilizationAssemblyFactoryOrder {
	// v39 records carry validated relationship fields; edges are still included
	// in the source state hash/provenance, but this read model derives joins
	// from the typed record fields that the store validates on append.
	requirements := map[string][]string{}
	requirementFactoryOrders := map[string]string{}
	acceptance := map[string][]string{}
	tasks := map[string][]string{}
	releases := map[string][]string{}
	var orders []CivilizationAssemblyFactoryOrder

	for _, record := range records {
		if typed, ok := record.(*Requirement); ok {
			requirementFactoryOrders[typed.CommonNode.ID] = typed.FactoryOrderID
		}
	}

	for _, record := range records {
		switch typed := record.(type) {
		case *Requirement:
			requirements[typed.FactoryOrderID] = append(requirements[typed.FactoryOrderID], typed.CommonNode.ID)
		case *AcceptanceCriterion:
			if factoryOrderID := requirementFactoryOrders[typed.RequirementID]; factoryOrderID != "" {
				acceptance[factoryOrderID] = append(acceptance[factoryOrderID], typed.CommonNode.ID)
			}
		case *Task:
			if typed.FactoryOrderID != nil {
				tasks[*typed.FactoryOrderID] = append(tasks[*typed.FactoryOrderID], typed.CommonNode.ID)
			}
		case *ReleaseCandidate:
			releases[typed.FactoryOrderID] = append(releases[typed.FactoryOrderID], typed.CommonNode.ID)
		}
	}

	for _, record := range records {
		if typed, ok := record.(*FactoryOrder); ok {
			orders = append(orders, CivilizationAssemblyFactoryOrder{
				ID:                      typed.CommonNode.ID,
				Status:                  commonStatus(typed.CommonNode),
				RiskClass:               typed.RiskClass,
				ReleasePolicy:           typed.ReleasePolicy,
				RequirementRefs:         appendSortedUnique(nil, requirements[typed.CommonNode.ID]...),
				AcceptanceCriterionRefs: appendSortedUnique(nil, acceptance[typed.CommonNode.ID]...),
				TaskRefs:                appendSortedUnique(nil, tasks[typed.CommonNode.ID]...),
				ReleaseCandidateRefs:    appendSortedUnique(nil, releases[typed.CommonNode.ID]...),
			})
		}
	}
	sort.Slice(orders, func(i, j int) bool { return orders[i].ID < orders[j].ID })
	return orders
}

func civilizationAssemblyWorkEvidence(records []Record) CivilizationAssemblyWorkEvidence {
	evidence := CivilizationAssemblyWorkEvidence{}
	for _, record := range records {
		switch typed := record.(type) {
		case *Task:
			evidence.TaskRefs = append(evidence.TaskRefs, typed.CommonNode.ID)
		case *Artifact:
			evidence.ArtifactRefs = append(evidence.ArtifactRefs, typed.CommonNode.ID)
		case *TestRun:
			evidence.TestRunRefs = append(evidence.TestRunRefs, typed.CommonNode.ID)
		case *GateResult:
			evidence.GateResultRefs = append(evidence.GateResultRefs, typed.CommonNode.ID)
		case *AuditReport:
			evidence.AuditReportRefs = append(evidence.AuditReportRefs, typed.CommonNode.ID)
		}
	}
	evidence.TaskRefs = appendSortedUnique(nil, evidence.TaskRefs...)
	evidence.ArtifactRefs = appendSortedUnique(nil, evidence.ArtifactRefs...)
	evidence.TestRunRefs = appendSortedUnique(nil, evidence.TestRunRefs...)
	evidence.GateResultRefs = appendSortedUnique(nil, evidence.GateResultRefs...)
	evidence.AuditReportRefs = appendSortedUnique(nil, evidence.AuditReportRefs...)
	evidence.SourceRefs = appendSortedUnique(nil, evidence.TaskRefs...)
	evidence.SourceRefs = appendSortedUnique(evidence.SourceRefs, evidence.ArtifactRefs...)
	evidence.SourceRefs = appendSortedUnique(evidence.SourceRefs, evidence.TestRunRefs...)
	evidence.SourceRefs = appendSortedUnique(evidence.SourceRefs, evidence.GateResultRefs...)
	evidence.SourceRefs = appendSortedUnique(evidence.SourceRefs, evidence.AuditReportRefs...)
	if len(evidence.SourceRefs) == 0 {
		evidence.Status = CivilizationAssemblyFieldUnavailable
		evidence.Summary = "no task, artifact, test, gate, or audit records are present"
		return evidence
	}
	evidence.Status = CivilizationAssemblyFieldAvailable
	evidence.Summary = "work evidence derived from task, artifact, test, gate, and audit records"
	return evidence
}

func civilizationAssemblySiteConsumerStatus(records []Record) CivilizationAssemblyFieldStatus {
	var refs []string
	for _, record := range records {
		artifact, ok := record.(*Artifact)
		if !ok {
			continue
		}
		artifactPath := ""
		if artifact.Path != nil {
			artifactPath = strings.ToLower(strings.TrimSpace(*artifact.Path))
		}
		if strings.HasPrefix(artifactPath, civilizationAssemblySiteConsumerPathPrefix) ||
			containsString(artifact.SourceRefs, civilizationAssemblySiteConsumerSourceRef) {
			refs = append(refs, artifact.CommonNode.ID)
		}
	}
	refs = appendSortedUnique(nil, refs...)
	if len(refs) == 0 {
		return CivilizationAssemblyFieldStatus{
			Status:  CivilizationAssemblyFieldUnavailable,
			Summary: "no EventGraph artifact records declare a read-only Civilization Assembly Site consumer",
		}
	}
	return CivilizationAssemblyFieldStatus{
		Status:     CivilizationAssemblyFieldAvailable,
		Summary:    "read-only Civilization Assembly Site consumer status derived from EventGraph artifact records",
		SourceRefs: refs,
	}
}

func civilizationAssemblyOpenGates(records []Record) []CivilizationAssemblyGateSummary {
	var out []CivilizationAssemblyGateSummary
	for _, record := range records {
		gate, ok := record.(*GateResult)
		if !ok || normalizedStatus(gate.CommonNode) == "pass" {
			continue
		}
		out = append(out, CivilizationAssemblyGateSummary{
			ID:                 gate.CommonNode.ID,
			GateName:           gate.GateName,
			Status:             commonStatus(gate.CommonNode),
			FactoryOrderID:     gate.FactoryOrderID,
			ReleaseCandidateID: gate.ReleaseCandidateID,
			EvidenceRefs:       appendSortedUnique(nil, gate.EvidenceRefs...),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func civilizationAssemblyResidualRisks(records []Record) []CivilizationAssemblyResidualRisk {
	var out []CivilizationAssemblyResidualRisk
	for _, record := range records {
		switch typed := record.(type) {
		case *Failure:
			if !isUnresolvedFailureStatus(normalizedStatus(typed.CommonNode)) {
				continue
			}
			out = append(out, CivilizationAssemblyResidualRisk{
				ID:       typed.CommonNode.ID,
				Kind:     TypeFailure,
				Severity: typed.Severity,
				Status:   commonStatus(typed.CommonNode),
				Summary:  typed.Summary,
			})
		case *ContradictionLog:
			if normalizedStatus(typed.CommonNode) == "resolved" {
				continue
			}
			out = append(out, CivilizationAssemblyResidualRisk{
				ID:       typed.CommonNode.ID,
				Kind:     TypeContradictionLog,
				Severity: typed.Severity,
				Status:   commonStatus(typed.CommonNode),
				Summary:  typed.ClaimARef + " conflicts with " + typed.ClaimBRef,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func civilizationAssemblyUnavailableFields(projection CivilizationAssemblyProjection) []CivilizationAssemblyUnavailableField {
	var fields []CivilizationAssemblyUnavailableField
	add := func(field string, status CivilizationAssemblyFieldAvailability, reason string) {
		fields = append(fields, CivilizationAssemblyUnavailableField{Field: field, Status: status, Reason: reason})
	}
	if projection.AuthorityState.Status != CivilizationAssemblyFieldAvailable {
		add("authority_state", projection.AuthorityState.Status, projection.AuthorityState.Summary)
	}
	if projection.ExternalCommitteeState.Status != CivilizationAssemblyFieldAvailable {
		add("external_committee_state", projection.ExternalCommitteeState.Status, projection.ExternalCommitteeState.Summary)
	}
	if len(projection.ActorRoster) == 0 {
		add("actor_roster", CivilizationAssemblyFieldUnavailable, "no ActorIdentity records are present")
	}
	if len(projection.RoleBindings) == 0 {
		add("role_bindings", CivilizationAssemblyFieldUnavailable, "no authority or approval records bind actors to roles")
	}
	if len(projection.AgentLifecycleSummary) == 0 {
		add("agent_lifecycle_summary", CivilizationAssemblyFieldUnavailable, "no LifecycleTransition or TrustRecord records are present")
	}
	if len(projection.FactoryOrderSummary) == 0 {
		add("factory_order_summary", CivilizationAssemblyFieldUnavailable, "no FactoryOrder records are present")
	}
	if projection.WorkEvidenceSummary.Status != CivilizationAssemblyFieldAvailable {
		add("work_evidence_summary", projection.WorkEvidenceSummary.Status, projection.WorkEvidenceSummary.Summary)
	}
	if projection.SiteConsumerStatus.Status != CivilizationAssemblyFieldAvailable {
		add("site_consumer_status", projection.SiteConsumerStatus.Status, projection.SiteConsumerStatus.Summary)
	}
	if len(projection.ValidationRefs) == 0 {
		add("validation_refs", CivilizationAssemblyFieldUnavailable, "no TestRun, GateResult, or caller-provided validation refs are present")
	}
	sort.Slice(fields, func(i, j int) bool { return fields[i].Field < fields[j].Field })
	return fields
}

func civilizationAssemblyFailureReasons(records []Record) []string {
	var reasons []string
	if conflicts := conflictingAuthorityDecisionReasons(records); len(conflicts) > 0 {
		reasons = append(reasons, conflicts...)
	}
	if authorityIntegrity := authorityReferenceIntegrityReasonsFromRecords(records); len(authorityIntegrity) > 0 {
		reasons = append(reasons, authorityIntegrity...)
	}
	for _, record := range records {
		contradiction, ok := record.(*ContradictionLog)
		if !ok || normalizedStatus(contradiction.CommonNode) == "resolved" {
			continue
		}
		severity := strings.ToLower(strings.TrimSpace(contradiction.Severity))
		if severity == "high" || severity == "critical" {
			reasons = append(reasons, "unresolved "+severity+" contradiction "+contradiction.CommonNode.ID+" blocks trusted projection")
		}
	}
	return appendSortedUnique(nil, reasons...)
}

func civilizationAssemblyDerivationStatus(snapshot civilizationAssemblySnapshot, projection CivilizationAssemblyProjection) CivilizationAssemblyDerivationStatus {
	if len(projection.FailureReasons) > 0 {
		return CivilizationAssemblyDerivationFailed
	}
	if len(snapshot.records) == 0 {
		return CivilizationAssemblyDerivationUnavailable
	}
	if len(projection.WithheldOrUnavailableFields) > 0 || len(projection.OpenGateSummary) > 0 || len(projection.ResidualRiskSummary) > 0 {
		return CivilizationAssemblyDerivationPartial
	}
	return CivilizationAssemblyDerivationComplete
}

func civilizationAssemblyBoundaryFlags(extra []string) []string {
	flags := []string{
		"read_only_projection",
		"no_eventgraph_writes",
		"no_materialized_projection_store",
		"no_runtime_execution",
		"no_protected_actions",
		"no_site_replacement",
		"no_work_mutation",
		"no_hive_action_api",
		"no_deploy",
		"no_production_data",
	}
	return appendSortedUnique(flags, extra...)
}

func civilizationAssemblyProvenanceRefs(projection CivilizationAssemblyProjection) []string {
	refs := []string{}
	refs = appendSortedUnique(refs, projection.SourceEventIDsOrQueryWindow...)
	refs = appendSortedUnique(refs, projection.AuthorityState.SourceRefs...)
	refs = appendSortedUnique(refs, projection.ExternalCommitteeState.DecisionRefs...)
	refs = appendSortedUnique(refs, projection.ExternalCommitteeState.ApprovalRefs...)
	refs = appendSortedUnique(refs, projection.WorkEvidenceSummary.SourceRefs...)
	refs = appendSortedUnique(refs, projection.SiteConsumerStatus.SourceRefs...)
	return refs
}

func hasConflictingAuthorityDecisions(decisions []CivilizationAssemblyAuthorityDecision) bool {
	seen := map[string]string{}
	for _, decision := range decisions {
		authorityRequestID := strings.TrimSpace(decision.AuthorityRequestID)
		if authorityRequestID == "" {
			continue
		}
		if existing, ok := seen[authorityRequestID]; ok && existing != decision.Decision {
			return true
		}
		seen[authorityRequestID] = decision.Decision
	}
	return false
}

func conflictingAuthorityDecisionReasons(records []Record) []string {
	seen := map[string]AuthorityDecision{}
	var reasons []string
	for _, record := range records {
		decision, ok := record.(*AuthorityDecision)
		if !ok {
			continue
		}
		authorityRequestID := strings.TrimSpace(decision.AuthorityRequestID)
		if authorityRequestID == "" {
			continue
		}
		if existing, ok := seen[authorityRequestID]; ok && existing.Decision != decision.Decision {
			reasons = append(reasons, "conflicting AuthorityDecision records for "+authorityRequestID+": "+existing.CommonNode.ID+"="+existing.Decision+" "+decision.CommonNode.ID+"="+decision.Decision)
			continue
		}
		seen[authorityRequestID] = *decision
	}
	return reasons
}

func authorityReferenceIntegrityReasonsFromRecords(records []Record) []string {
	var requests []CivilizationAssemblyAuthorityRequest
	var decisions []CivilizationAssemblyAuthorityDecision
	var receipts []CivilizationAssemblyExecutionReceipt
	for _, record := range records {
		switch typed := record.(type) {
		case *AuthorityRequest:
			requests = append(requests, CivilizationAssemblyAuthorityRequest{ID: typed.CommonNode.ID})
		case *AuthorityDecision:
			decisions = append(decisions, CivilizationAssemblyAuthorityDecision{
				ID:                 typed.CommonNode.ID,
				AuthorityRequestID: typed.AuthorityRequestID,
			})
		case *ExecutionReceipt:
			receipts = append(receipts, CivilizationAssemblyExecutionReceipt{
				ID:                  typed.CommonNode.ID,
				AuthorityDecisionID: typed.AuthorityDecisionID,
			})
		}
	}
	return authorityReferenceIntegrityReasons(requests, decisions, receipts)
}

func authorityReferenceIntegrityReasons(requests []CivilizationAssemblyAuthorityRequest, decisions []CivilizationAssemblyAuthorityDecision, receipts []CivilizationAssemblyExecutionReceipt) []string {
	requestIDs := map[string]struct{}{}
	decisionIDs := map[string]struct{}{}
	for _, request := range requests {
		if id := strings.TrimSpace(request.ID); id != "" {
			requestIDs[id] = struct{}{}
		}
	}
	for _, decision := range decisions {
		if id := strings.TrimSpace(decision.ID); id != "" {
			decisionIDs[id] = struct{}{}
		}
	}

	var reasons []string
	for _, decision := range decisions {
		requestID := strings.TrimSpace(decision.AuthorityRequestID)
		if requestID == "" {
			reasons = append(reasons, "AuthorityDecision "+decision.ID+" is missing authority_request_id")
			continue
		}
		if _, ok := requestIDs[requestID]; !ok {
			reasons = append(reasons, "AuthorityDecision "+decision.ID+" references missing AuthorityRequest "+requestID)
		}
	}
	for _, receipt := range receipts {
		decisionID := strings.TrimSpace(receipt.AuthorityDecisionID)
		if decisionID == "" {
			reasons = append(reasons, "ExecutionReceipt "+receipt.ID+" is missing authority_decision_id")
			continue
		}
		if _, ok := decisionIDs[decisionID]; !ok {
			reasons = append(reasons, "ExecutionReceipt "+receipt.ID+" references missing AuthorityDecision "+decisionID)
		}
	}
	return appendSortedUnique(nil, reasons...)
}

func sortedMapKeysRecord(m map[string]Record) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedMapKeysEdge(m map[string]CommonEdge) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func commonStatus(common CommonNode) string {
	if common.Status == nil {
		return ""
	}
	return *common.Status
}

func normalizedStatus(common CommonNode) string {
	return strings.ToLower(strings.TrimSpace(commonStatus(common)))
}

func appendSortedUnique(base []string, values ...string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(base)+len(values))
	add := func(value string) {
		if value == "" {
			return
		}
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	for _, value := range base {
		add(value)
	}
	for _, value := range values {
		add(value)
	}
	sort.Strings(out)
	return out
}

func isUnresolvedFailureStatus(status string) bool {
	switch status {
	case "repaired", "waived", "closed":
		return false
	default:
		return true
	}
}
