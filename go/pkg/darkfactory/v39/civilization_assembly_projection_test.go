package v39

import (
	"reflect"
	"strings"
	"testing"
)

func TestCivilizationAssemblyProjectionCompleteDeterministicFixture(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)

	first := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{
		GeneratedAt:    fixedTime,
		ValidationRefs: []string{"external-review:pr-head"},
	})
	second := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{
		GeneratedAt:    fixedTime,
		ValidationRefs: []string{"external-review:pr-head"},
	})

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("projection is not deterministic\nfirst=%+v\nsecond=%+v", first, second)
	}
	if first.DerivationStatus != CivilizationAssemblyDerivationComplete {
		t.Fatalf("derivation status = %s, want complete: %+v", first.DerivationStatus, first.WithheldOrUnavailableFields)
	}
	if first.ProjectionSubject != CivilizationAssemblyProjectionSubject {
		t.Fatalf("subject = %s", first.ProjectionSubject)
	}
	if first.SourceEventGraphHeadOrStateVersion == "" || !strings.HasPrefix(first.SourceEventGraphHeadOrStateVersion, "sha256:") {
		t.Fatalf("missing deterministic source state version: %s", first.SourceEventGraphHeadOrStateVersion)
	}
	if first.AuthorityState.Status != CivilizationAssemblyFieldAvailable || len(first.AuthorityState.AuthorityDecisions) != 1 {
		t.Fatalf("authority not derived from AuthorityDecision records: %+v", first.AuthorityState)
	}
	if first.ExternalCommitteeState.Status != CivilizationAssemblyFieldAvailable || !containsString(first.ExternalCommitteeState.ApprovalRefs, "approval_civ_001") {
		t.Fatalf("external committee state missing approval evidence: %+v", first.ExternalCommitteeState)
	}
	if len(first.ActorRoster) != 2 {
		t.Fatalf("actor roster = %+v, want agent and human actors", first.ActorRoster)
	}
	if len(first.FactoryOrderSummary) != 1 || !containsString(first.FactoryOrderSummary[0].TaskRefs, "tsk_civ_001") {
		t.Fatalf("factory order summary missing work task refs: %+v", first.FactoryOrderSummary)
	}
	if first.WorkEvidenceSummary.Status != CivilizationAssemblyFieldAvailable || !containsString(first.WorkEvidenceSummary.TestRunRefs, "tr_civ_001") {
		t.Fatalf("work evidence missing test run refs: %+v", first.WorkEvidenceSummary)
	}
	if first.SiteConsumerStatus.Status != CivilizationAssemblyFieldAvailable || !containsString(first.SiteConsumerStatus.SourceRefs, "site_consumer_civ_001") {
		t.Fatalf("site consumer status missing EventGraph artifact evidence: %+v", first.SiteConsumerStatus)
	}
	for _, forbidden := range []string{"no_eventgraph_writes", "no_runtime_execution", "no_protected_actions", "no_site_replacement"} {
		if !containsString(first.BoundaryFlags, forbidden) {
			t.Fatalf("boundary flag %s missing from %+v", forbidden, first.BoundaryFlags)
		}
	}
	if len(first.WithheldOrUnavailableFields) != 0 {
		t.Fatalf("complete fixture should not have unavailable fields: %+v", first.WithheldOrUnavailableFields)
	}
}

func TestCivilizationAssemblyProjectionMissingAuthorityFailsClosed(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	deleteRecord(store, "auth_dec_civ_001")
	deleteRecord(store, "exec_civ_001")
	appendRecord(t, store, &MemoryReference{AdvisoryReference: advisory("mem_civ_context", TypeMemoryReference, "tsk_civ_001")})

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if projection.DerivationStatus != CivilizationAssemblyDerivationPartial {
		t.Fatalf("derivation status = %s, want partial", projection.DerivationStatus)
	}
	if projection.AuthorityState.Status != CivilizationAssemblyFieldUnavailable {
		t.Fatalf("authority state = %+v, want unavailable", projection.AuthorityState)
	}
	if len(projection.AuthorityState.AuthorityDecisions) != 0 {
		t.Fatalf("authority decisions should be absent after deleting authority evidence: %+v", projection.AuthorityState.AuthorityDecisions)
	}
	if !projectionHasUnavailableField(projection, "authority_state") {
		t.Fatalf("missing authority was not surfaced as unavailable: %+v", projection.WithheldOrUnavailableFields)
	}
	for _, binding := range projection.RoleBindings {
		if binding.SourceType == TypeMemoryReference {
			t.Fatalf("memory evidence must not create authority role bindings: %+v", projection.RoleBindings)
		}
	}
}

func TestCivilizationAssemblyProjectionConflictingAuthorityFails(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	appendRecord(t, store, &AuthorityDecision{
		CommonNode:         common("auth_dec_civ_conflict", TypeAuthorityDecision, "approved"),
		AuthorityRequestID: "auth_req_civ_001",
		DeciderActorID:     "act_human",
		DeciderRole:        "External Committee",
		Decision:           "Forbidden",
		Reason:             "negative conflict fixture",
		Scope:              []string{"eventgraph.read.projection"},
	})

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if projection.DerivationStatus != CivilizationAssemblyDerivationFailed {
		t.Fatalf("derivation status = %s, want failed", projection.DerivationStatus)
	}
	if len(projection.FailureReasons) == 0 || !strings.Contains(projection.FailureReasons[0], "conflicting AuthorityDecision") {
		t.Fatalf("missing authority conflict failure reason: %+v", projection.FailureReasons)
	}
}

func TestCivilizationAssemblyProjectionOpenGateAndResidualRiskArePartial(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	appendRecord(t, store, &GateResult{
		CommonNode:     common("gate_civ_fail", TypeGateResult, "fail"),
		FactoryOrderID: "fo_civ_001",
		GateName:       "gate_s_projection_residual",
		EvidenceRefs:   []string{"tr_civ_001"},
	})
	appendRecord(t, store, &Failure{
		CommonNode:     common("failure_civ_001", TypeFailure, "open"),
		FactoryOrderID: strPtr("fo_civ_001"),
		TaskID:         strPtr("tsk_civ_001"),
		GateResultID:   strPtr("gate_civ_fail"),
		FailureClass:   "traceability_gap",
		Severity:       "high",
		Summary:        "residual projection evidence is incomplete",
	})

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if projection.DerivationStatus != CivilizationAssemblyDerivationPartial {
		t.Fatalf("derivation status = %s, want partial", projection.DerivationStatus)
	}
	if len(projection.OpenGateSummary) != 1 || projection.OpenGateSummary[0].ID != "gate_civ_fail" {
		t.Fatalf("open gate not projected: %+v", projection.OpenGateSummary)
	}
	if len(projection.ResidualRiskSummary) != 1 || projection.ResidualRiskSummary[0].ID != "failure_civ_001" {
		t.Fatalf("residual risk not projected: %+v", projection.ResidualRiskSummary)
	}
}

func TestCivilizationAssemblyProjectionResolvedFailureDoesNotStayPartial(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	appendRecord(t, store, &Failure{
		CommonNode:     common("failure_civ_closed", TypeFailure, "closed"),
		FactoryOrderID: strPtr("fo_civ_001"),
		FailureClass:   "traceability_gap",
		Severity:       "high",
		Summary:        "historical failure was closed",
	})

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if projection.DerivationStatus != CivilizationAssemblyDerivationComplete {
		t.Fatalf("closed historical failure should not make projection partial: %+v", projection)
	}
	if len(projection.ResidualRiskSummary) != 0 {
		t.Fatalf("closed failure should not appear as residual risk: %+v", projection.ResidualRiskSummary)
	}
}

func TestCivilizationAssemblyProjectionEmptyStoreUnavailable(t *testing.T) {
	projection := NewInMemoryStore().ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if projection.DerivationStatus != CivilizationAssemblyDerivationUnavailable {
		t.Fatalf("derivation status = %s, want unavailable", projection.DerivationStatus)
	}
	if !projectionHasUnavailableField(projection, "authority_state") {
		t.Fatalf("empty store should surface unavailable authority: %+v", projection.WithheldOrUnavailableFields)
	}
}

func TestCivilizationAssemblyProjectionCloneFailureFailsClosed(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	store.mu.Lock()
	store.records["clone_failure"] = &projectionCloneFailureRecord{CommonNode: common("clone_failure", "ProjectionCloneFailure", "recorded")}
	store.canonicalByID["clone_failure"] = []byte(`{"id":"clone_failure"}`)
	store.mu.Unlock()

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if projection.DerivationStatus != CivilizationAssemblyDerivationFailed {
		t.Fatalf("clone failure status = %s, want failed", projection.DerivationStatus)
	}
	if len(projection.FailureReasons) == 0 || !strings.Contains(projection.FailureReasons[0], "could not be cloned") {
		t.Fatalf("missing clone failure reason: %+v", projection.FailureReasons)
	}
}

func TestCivilizationAssemblyProjectionSiteConsumerRequiresExplicitMarker(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	appendRecord(t, store, &Artifact{
		CommonNode:   common("artifact_unrelated_civilization_report", TypeArtifact, "verified"),
		ArtifactType: "report",
		Path:         strPtr("report://civilization/unrelated"),
		ContentHash:  strPtr("sha256:unrelated-civilization-report"),
	})

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if !reflect.DeepEqual(projection.SiteConsumerStatus.SourceRefs, []string{"site_consumer_civ_001"}) {
		t.Fatalf("unmarked civilization report should not count as Site consumer evidence: %+v", projection.SiteConsumerStatus)
	}
}

func TestCivilizationAssemblyProjectionOpenCriticalContradictionFails(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	appendRecord(t, store, &ContradictionLog{
		CommonNode:      common("contradiction_civ_critical", TypeContradictionLog, "open"),
		ContradictionID: "contradiction_civ_critical",
		ClaimARef:       "auth_dec_civ_001",
		ClaimBRef:       "approval_civ_001",
		Severity:        "critical",
	})

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if projection.DerivationStatus != CivilizationAssemblyDerivationFailed {
		t.Fatalf("derivation status = %s, want failed", projection.DerivationStatus)
	}
	if len(projection.FailureReasons) == 0 || !strings.Contains(projection.FailureReasons[0], "critical contradiction") {
		t.Fatalf("missing critical contradiction failure: %+v", projection.FailureReasons)
	}
}

func civilizationAssemblyProjectionStore(t *testing.T) *InMemoryStore {
	t.Helper()
	store := NewInMemoryStore()
	foID := "fo_civ_001"
	reqID := "req_civ_001"
	acID := "ac_civ_001"
	taskID := "tsk_civ_001"
	testCaseID := "tc_civ_001"
	testRunID := "tr_civ_001"
	gateID := "gate_civ_pass"
	authReqID := "auth_req_civ_001"
	authDecisionID := "auth_dec_civ_001"

	appendRecord(t, store, &FactoryOrder{
		CommonNode:          common(foID, TypeFactoryOrder, "accepted"),
		FactoryOrderVersion: 1,
		SourceIntentHash:    "sha256:civilization-intent",
		SourceIntentRef:     "eventgraph://fixture/civilization",
		RiskClass:           "high",
		ReleasePolicy:       "human_approval_required",
	})
	appendRecord(t, store, &Requirement{
		CommonNode:     common(reqID, TypeRequirement, "accepted"),
		FactoryOrderID: foID,
		Text:           "derive civilization assembly projection from EventGraph records",
		Source:         "explicit",
		RiskClass:      "high",
	})
	appendRecord(t, store, &AcceptanceCriterion{
		CommonNode:           common(acID, TypeAcceptanceCriterion, "verified"),
		RequirementID:        reqID,
		Text:                 "projection fails closed when authority evidence is absent",
		Source:               "explicit",
		VerificationMethod:   "test",
		RequiredEvidenceType: "go_test",
		OwnerRole:            "maintainer",
		RiskClass:            "high",
	})
	appendRecord(t, store, &Task{
		CommonNode:     common(taskID, TypeTask, "verified"),
		FactoryOrderID: &foID,
		Cell:           "cell_projection",
		State:          "verified",
		Priority:       1,
		RiskClass:      "high",
	})
	appendRecord(t, store, &ActorIdentity{
		CommonNode:   common("actor_identity_agent", TypeActorIdentity, "active"),
		ActorID:      "act_agent",
		ActorType:    "agent",
		IdentityMode: "fixture",
	})
	appendRecord(t, store, &ActorIdentity{
		CommonNode:   common("actor_identity_human", TypeActorIdentity, "active"),
		ActorID:      "act_human",
		ActorType:    "human",
		IdentityMode: "externally_managed",
	})
	appendRecord(t, store, &AuthorityRequest{
		CommonNode:   common(authReqID, TypeAuthorityRequest, "open"),
		ActorID:      "act_agent",
		ActorRole:    "ProjectionBuilder",
		Action:       "eventgraph.read.projection",
		TargetType:   "civilization_assembly",
		TargetID:     foID,
		RiskClass:    "high",
		Reason:       "derive read-only civilization assembly state",
		EvidenceRefs: []string{reqID, acID},
	})
	appendRecord(t, store, &AuthorityDecision{
		CommonNode:         common(authDecisionID, TypeAuthorityDecision, "approved"),
		AuthorityRequestID: authReqID,
		DeciderActorID:     "act_human",
		DeciderRole:        "External Committee",
		Decision:           "Autonomous",
		Reason:             "bounded deterministic read-model fixture",
		Scope:              []string{"eventgraph.read.projection"},
		Conditions:         []string{"read-only", "no persistent writes"},
	})
	appendRecord(t, store, &HumanApproval{
		CommonNode:      common("approval_civ_001", TypeHumanApproval, "approved"),
		RequestRef:      authReqID,
		ApproverActorID: "act_human",
		ApproverRole:    "External Committee",
		Decision:        "approved",
		Reason:          "fixture approval for complete projection evidence",
	})
	appendRecord(t, store, &ExecutionReceipt{
		CommonNode:          common("exec_civ_001", TypeExecutionReceipt, "recorded"),
		AuthorityDecisionID: authDecisionID,
		Action:              "eventgraph.read.projection",
		TargetID:            foID,
		Result:              "succeeded",
		EvidenceRefs:        []string{testRunID},
	})
	appendRecord(t, store, &LifecycleTransition{
		CommonNode:          common("lt_civ_001", TypeLifecycleTransition, "recorded"),
		ActorID:             "act_agent",
		FromState:           "trial",
		ToState:             "active",
		Reason:              "projection fixture verified",
		AuthorityDecisionID: &authDecisionID,
	})
	appendRecord(t, store, &TrustRecord{
		CommonNode:     common("trust_civ_001", TypeTrustRecord, "recorded"),
		SubjectActorID: "act_agent",
		TrustLevel:     "fixture",
		EvidenceRefs:   []string{"exec_civ_001"},
		Reason:         "deterministic projection fixture",
	})
	appendRecord(t, store, &Artifact{
		CommonNode:   common("artifact_civ_001", TypeArtifact, "verified"),
		TaskID:       &taskID,
		ArtifactType: "test",
		Path:         strPtr("go/pkg/darkfactory/v39/civilization_assembly_projection_test.go"),
		ContentHash:  strPtr("sha256:projection-fixture"),
	})
	appendRecord(t, store, &Artifact{
		CommonNode:   commonWithSourceRefs("site_consumer_civ_001", TypeArtifact, "verified", []string{civilizationAssemblySiteConsumerSourceRef}),
		TaskID:       &taskID,
		ArtifactType: "report",
		Path:         strPtr("site://ops/civilization/read-only"),
		ContentHash:  strPtr("sha256:site-consumer-status"),
	})
	appendRecord(t, store, &TestCase{
		CommonNode:            common(testCaseID, TypeTestCase, "active"),
		AcceptanceCriterionID: &acID,
		RequirementID:         &reqID,
		Name:                  "civilization assembly projection",
		TestType:              "unit",
		Path:                  strPtr("go/pkg/darkfactory/v39/civilization_assembly_projection_test.go"),
	})
	appendRecord(t, store, &TestRun{
		CommonNode: common(testRunID, TypeTestRun, "pass"),
		TestCaseID: &testCaseID,
		Command:    "go test ./pkg/darkfactory/v39",
	})
	appendRecord(t, store, &GateResult{
		CommonNode:     common(gateID, TypeGateResult, "pass"),
		FactoryOrderID: foID,
		GateName:       "gate_s_projection_fixture",
		EvidenceRefs:   []string{testRunID},
	})
	appendRecord(t, store, &AuditReport{
		CommonNode: common("audit_civ_001", TypeAuditReport, "complete"),
		TargetType: "factory_order",
		TargetID:   foID,
		TraceScore: 1,
	})

	appendEdge(t, store, edge("edge_civ_fo_req", EdgeRequires, foID, reqID))
	appendEdge(t, store, edge("edge_civ_req_ac", EdgeRequires, reqID, acID))
	appendEdge(t, store, edge("edge_civ_ac_task", EdgeDecomposedInto, acID, taskID))
	appendEdge(t, store, edge("edge_civ_actor_auth", EdgeRequestedAuthority, "actor_identity_agent", authReqID))
	appendEdge(t, store, edge("edge_civ_auth_dec", EdgeDecidedBy, authReqID, authDecisionID))
	appendEdge(t, store, edge("edge_civ_auth_exec", EdgeReceiptedBy, authDecisionID, "exec_civ_001"))
	appendEdge(t, store, edge("edge_civ_task_artifact", EdgeProduced, taskID, "artifact_civ_001"))
	appendEdge(t, store, edge("edge_civ_task_site", EdgeProduced, taskID, "site_consumer_civ_001"))
	appendEdge(t, store, edge("edge_civ_task_tc", EdgeVerifies, taskID, testCaseID))
	appendEdge(t, store, edge("edge_civ_tc_tr", EdgeVerifies, testCaseID, testRunID))
	appendEdge(t, store, edge("edge_civ_tr_gate", EdgeProduced, testRunID, gateID))
	return store
}

type projectionCloneFailureRecord struct {
	CommonNode
}

func (r *projectionCloneFailureRecord) Validate() error {
	return nil
}

func projectionHasUnavailableField(projection CivilizationAssemblyProjection, field string) bool {
	for _, unavailable := range projection.WithheldOrUnavailableFields {
		if unavailable.Field == field {
			return true
		}
	}
	return false
}
