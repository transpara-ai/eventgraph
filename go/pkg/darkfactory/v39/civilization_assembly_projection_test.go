package v39

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
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
	if !strings.HasPrefix(first.ProjectionID, "civilization_assembly:") || len(strings.TrimPrefix(first.ProjectionID, "civilization_assembly:")) != 16 {
		t.Fatalf("derived projection id = %s, want stable 16-character state prefix", first.ProjectionID)
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
	if first.IssueIntakeProjection.Status != CivilizationAssemblyFieldAvailable || len(first.IssueIntakeProjection.Issues) != 2 {
		t.Fatalf("issue intake projection missing typed GitHub issue refs: %+v", first.IssueIntakeProjection)
	}
	if !containsCivilizationIssue(first.IssueIntakeProjection.Issues, "transpara-ai/docs", 172) ||
		!containsCivilizationIssue(first.IssueIntakeProjection.Issues, "transpara-ai/site", 115) {
		t.Fatalf("issue intake projection does not preserve docs#172 and site#115: %+v", first.IssueIntakeProjection.Issues)
	}
	if len(first.IssueIntakeProjection.Groups) != 2 {
		t.Fatalf("issue intake groups = %+v, want two repo/substrate groups", first.IssueIntakeProjection.Groups)
	}
	if first.IssueIntakeProjection.Groups[0].GroupID != "repo-transpara-ai-docs-substrate-factoryorder-risk-high-readiness-accepted" ||
		first.IssueIntakeProjection.Groups[0].PrimaryRepo != "transpara-ai/docs" ||
		first.IssueIntakeProjection.Groups[0].TouchedSubstrate != TypeFactoryOrder ||
		first.IssueIntakeProjection.Groups[0].RiskClass != "high" ||
		first.IssueIntakeProjection.Groups[0].Readiness != "accepted" ||
		!containsCivilizationIssueRef(first.IssueIntakeProjection.Groups[0].IssueRefs, "transpara-ai/docs", 172) ||
		!containsString(first.IssueIntakeProjection.Groups[0].SourceRefs, "github:transpara-ai/docs#172") ||
		!strings.Contains(first.IssueIntakeProjection.Groups[0].Recommendation, "read-only source-intent projection") {
		t.Fatalf("docs issue intake group does not preserve expected source-intent shape: %+v", first.IssueIntakeProjection.Groups[0])
	}
	if first.IssueIntakeProjection.Groups[1].GroupID != "repo-transpara-ai-site-substrate-requirement-risk-high-readiness-accepted" ||
		first.IssueIntakeProjection.Groups[1].PrimaryRepo != "transpara-ai/site" ||
		first.IssueIntakeProjection.Groups[1].TouchedSubstrate != TypeRequirement ||
		first.IssueIntakeProjection.Groups[1].RiskClass != "high" ||
		first.IssueIntakeProjection.Groups[1].Readiness != "accepted" ||
		!containsCivilizationIssueRef(first.IssueIntakeProjection.Groups[1].IssueRefs, "transpara-ai/site", 115) ||
		!containsString(first.IssueIntakeProjection.Groups[1].SourceRefs, "https://github.com/transpara-ai/site/issues/115") ||
		first.IssueIntakeProjection.Groups[1].Summary != "1 issue(s) share repo/substrate/risk/readiness key." {
		t.Fatalf("site issue intake group does not preserve expected source-intent shape: %+v", first.IssueIntakeProjection.Groups[1])
	}
	for _, boundary := range []string{"scanner_read_only", "no_eventgraph_writes", "no_github_issue_mutation", "no_pr_creation", "no_merge", "no_protected_action_approval"} {
		if !containsString(first.IssueIntakeProjection.ScannerBoundaries, boundary) {
			t.Fatalf("issue intake scanner boundary %s missing from %+v", boundary, first.IssueIntakeProjection.ScannerBoundaries)
		}
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

func TestCivilizationAssemblyProjectionDoesNotMutateStore(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	before := civilizationAssemblyStoreFootprint(t, store)

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	after := civilizationAssemblyStoreFootprint(t, store)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("projection mutated store\nbefore=%+v\nafter=%+v\nprojection=%+v", before, after, projection)
	}
}

func TestRecordCivilizationAssemblyProjectionPersistsProjectionStoreRecord(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	authorityDecisionRef := appendCivilizationProjectionStoreAuthority(t, store)
	before := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	record, err := store.RecordCivilizationAssemblyProjection(CivilizationAssemblyProjectionStoreRecordOptions{
		CreatedAt:            fixedTime,
		CreatedBy:            "act_projection_store",
		AuthorityDecisionRef: authorityDecisionRef,
		ProjectionOptions: CivilizationAssemblyProjectionOptions{
			GeneratedAt:    fixedTime,
			ValidationRefs: []string{"cfar:pr-head", "docs:207-authority-packet"},
		},
		AdditionalSourceRefs: []string{"github:transpara-ai/eventgraph#59"},
	})
	if err != nil {
		t.Fatalf("record projection: %v", err)
	}

	if record.CommonNode.Type != TypeCivilizationAssemblyProjectionStoreRecord {
		t.Fatalf("record type = %s", record.CommonNode.Type)
	}
	if record.DerivationStatus != CivilizationAssemblyDerivationComplete {
		t.Fatalf("derivation status = %s, want complete", record.DerivationStatus)
	}
	if record.ProjectionID != before.ProjectionID ||
		record.SourceEventGraphHeadOrStateVersion != before.SourceEventGraphHeadOrStateVersion {
		t.Fatalf("record does not preserve source projection identity: before=%+v record=%+v", before, record)
	}
	for _, want := range []string{"projection_store_local_only", "no_production_eventgraph_write", "no_runtime_execution", "no_protected_actions", "no_deploy"} {
		if !containsString(record.BoundaryFlags, want) || !containsString(record.Projection.BoundaryFlags, want) {
			t.Fatalf("projection-store boundary %s missing: record=%+v projection=%+v", want, record.BoundaryFlags, record.Projection.BoundaryFlags)
		}
	}
	for _, forbidden := range []string{"no_eventgraph_writes", "no_materialized_projection_store"} {
		if containsString(record.BoundaryFlags, forbidden) || containsString(record.Projection.BoundaryFlags, forbidden) {
			t.Fatalf("recorded projection should not retain read-only projection-store blocker %s: record=%+v projection=%+v", forbidden, record.BoundaryFlags, record.Projection.BoundaryFlags)
		}
	}
	for _, want := range []string{authorityDecisionRef, "github:transpara-ai/eventgraph#59"} {
		if !containsString(record.ProvenanceRefs, want) || !containsString(record.CommonNode.SourceRefs, want) {
			t.Fatalf("record provenance missing %s: %+v source=%+v", want, record.ProvenanceRefs, record.CommonNode.SourceRefs)
		}
	}
	if !containsString(record.ValidationRefs, "cfar:pr-head") ||
		!containsString(record.ValidationRefs, "tr_civ_001") ||
		!containsString(record.ValidationRefs, "gate_civ_pass") {
		t.Fatalf("record validation refs do not include caller and derived evidence: %+v", record.ValidationRefs)
	}
	stored := store.ByType(TypeCivilizationAssemblyProjectionStoreRecord)
	if len(stored) != 1 {
		t.Fatalf("projection store record count = %d, want 1", len(stored))
	}
	after := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})
	if after.SourceEventGraphHeadOrStateVersion != before.SourceEventGraphHeadOrStateVersion ||
		containsString(after.SourceEventIDsOrQueryWindow, record.CommonNode.ID) {
		t.Fatalf("projection store records must not recursively alter source projection state: before=%+v after=%+v record=%s", before, after, record.CommonNode.ID)
	}
}

func TestRecordCivilizationAssemblyProjectionIdempotentReplay(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	authorityDecisionRef := appendCivilizationProjectionStoreAuthority(t, store)
	options := CivilizationAssemblyProjectionStoreRecordOptions{
		CreatedAt:            fixedTime,
		CreatedBy:            "act_projection_store",
		AuthorityDecisionRef: authorityDecisionRef,
		ProjectionOptions: CivilizationAssemblyProjectionOptions{
			GeneratedAt:    fixedTime,
			ValidationRefs: []string{"cfar:pr-head"},
		},
	}

	first, err := store.RecordCivilizationAssemblyProjection(options)
	if err != nil {
		t.Fatalf("first record projection: %v", err)
	}
	second, err := store.RecordCivilizationAssemblyProjection(options)
	if err != nil {
		t.Fatalf("second record projection: %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("idempotent replay changed record\nfirst=%+v\nsecond=%+v", first, second)
	}
	if got := len(store.ByType(TypeCivilizationAssemblyProjectionStoreRecord)); got != 1 {
		t.Fatalf("projection store record count = %d, want 1", got)
	}
}

func TestRecordCivilizationAssemblyProjectionRequiresScopedAuthorityDecision(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	before := civilizationAssemblyStoreFootprint(t, store)

	_, err := store.RecordCivilizationAssemblyProjection(CivilizationAssemblyProjectionStoreRecordOptions{
		CreatedAt:            fixedTime,
		CreatedBy:            "act_projection_store",
		AuthorityDecisionRef: "auth_dec_civ_001",
		ProjectionOptions: CivilizationAssemblyProjectionOptions{
			GeneratedAt:    fixedTime,
			ValidationRefs: []string{"cfar:pr-head"},
		},
	})
	if !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("unscoped authority decision error = %v, want ErrInvalidRecord", err)
	}
	after := civilizationAssemblyStoreFootprint(t, store)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("unscoped authority decision mutated store\nbefore=%+v\nafter=%+v", before, after)
	}
}

func TestRecordCivilizationAssemblyProjectionRejectsApprovalRequiredAuthorityDecision(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	authorityDecisionRef := appendCivilizationProjectionStoreAuthorityDecision(t, store, "auth_dec_civ_projection_store_approval_required", "ApprovalRequired", nil)
	before := civilizationAssemblyStoreFootprint(t, store)

	_, err := store.RecordCivilizationAssemblyProjection(CivilizationAssemblyProjectionStoreRecordOptions{
		CreatedAt:            fixedTime,
		CreatedBy:            "act_projection_store",
		AuthorityDecisionRef: authorityDecisionRef,
		ProjectionOptions: CivilizationAssemblyProjectionOptions{
			GeneratedAt:    fixedTime,
			ValidationRefs: []string{"cfar:pr-head"},
		},
	})
	if !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("approval-required authority decision error = %v, want ErrInvalidRecord", err)
	}
	after := civilizationAssemblyStoreFootprint(t, store)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("approval-required authority decision mutated store\nbefore=%+v\nafter=%+v", before, after)
	}
}

func TestRecordCivilizationAssemblyProjectionRejectsExpiredAuthorityDecision(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	expiresAt := fixedTime.Add(-1)
	authorityDecisionRef := appendCivilizationProjectionStoreAuthorityDecision(t, store, "auth_dec_civ_projection_store_expired", "Autonomous", &expiresAt)
	before := civilizationAssemblyStoreFootprint(t, store)

	_, err := store.RecordCivilizationAssemblyProjection(CivilizationAssemblyProjectionStoreRecordOptions{
		CreatedAt:            fixedTime,
		CreatedBy:            "act_projection_store",
		AuthorityDecisionRef: authorityDecisionRef,
		ProjectionOptions: CivilizationAssemblyProjectionOptions{
			GeneratedAt:    fixedTime,
			ValidationRefs: []string{"cfar:pr-head"},
		},
	})
	if !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expired authority decision error = %v, want ErrInvalidRecord", err)
	}
	after := civilizationAssemblyStoreFootprint(t, store)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("expired authority decision mutated store\nbefore=%+v\nafter=%+v", before, after)
	}
}

func TestRecordCivilizationAssemblyProjectionRejectsRevokedAuthorityDecision(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	authorityDecisionRef := appendCivilizationProjectionStoreAuthorityDecision(t, store, "auth_dec_civ_projection_store_revoked", "Autonomous", nil)
	store.mu.Lock()
	*store.records[authorityDecisionRef].(*AuthorityDecision).CommonNode.Status = "revoked"
	store.mu.Unlock()
	before := civilizationAssemblyStoreFootprint(t, store)

	_, err := store.RecordCivilizationAssemblyProjection(CivilizationAssemblyProjectionStoreRecordOptions{
		CreatedAt:            fixedTime,
		CreatedBy:            "act_projection_store",
		AuthorityDecisionRef: authorityDecisionRef,
		ProjectionOptions: CivilizationAssemblyProjectionOptions{
			GeneratedAt:    fixedTime,
			ValidationRefs: []string{"cfar:pr-head"},
		},
	})
	if !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("revoked authority decision error = %v, want ErrInvalidRecord", err)
	}
	after := civilizationAssemblyStoreFootprint(t, store)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("revoked authority decision mutated store\nbefore=%+v\nafter=%+v", before, after)
	}
}

func TestRecordCivilizationAssemblyProjectionRequiresAuthorityDecision(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	before := civilizationAssemblyStoreFootprint(t, store)

	_, err := store.RecordCivilizationAssemblyProjection(CivilizationAssemblyProjectionStoreRecordOptions{
		CreatedAt:            fixedTime,
		CreatedBy:            "act_projection_store",
		AuthorityDecisionRef: "auth_dec_civ_missing",
		ProjectionOptions: CivilizationAssemblyProjectionOptions{
			GeneratedAt:    fixedTime,
			ValidationRefs: []string{"cfar:pr-head"},
		},
	})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing authority decision error = %v, want ErrNotFound", err)
	}
	after := civilizationAssemblyStoreFootprint(t, store)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("missing authority decision mutated store\nbefore=%+v\nafter=%+v", before, after)
	}
}

func TestRecordCivilizationAssemblyProjectionRejectsFailedProjection(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	authorityDecisionRef := appendCivilizationProjectionStoreAuthority(t, store)
	appendRecord(t, store, &AuthorityDecision{
		CommonNode:         common("auth_dec_civ_projection_conflict", TypeAuthorityDecision, "approved"),
		AuthorityRequestID: "auth_req_civ_001",
		DeciderActorID:     "act_human",
		DeciderRole:        "External Committee",
		Decision:           "Forbidden",
		Reason:             "negative projection-store conflict fixture",
		Scope:              []string{CivilizationAssemblyProjectionStoreAction},
	})
	before := civilizationAssemblyStoreFootprint(t, store)

	_, err := store.RecordCivilizationAssemblyProjection(CivilizationAssemblyProjectionStoreRecordOptions{
		CreatedAt:            fixedTime,
		CreatedBy:            "act_projection_store",
		AuthorityDecisionRef: authorityDecisionRef,
		ProjectionOptions: CivilizationAssemblyProjectionOptions{
			GeneratedAt:    fixedTime,
			ValidationRefs: []string{"cfar:pr-head"},
		},
	})
	if !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("failed projection error = %v, want ErrInvalidRecord", err)
	}
	after := civilizationAssemblyStoreFootprint(t, store)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("failed projection append mutated store\nbefore=%+v\nafter=%+v", before, after)
	}
}

func TestCivilizationAssemblyIssueIntakeProjectionUnavailableWithoutIssueRefs(t *testing.T) {
	projection := civilizationAssemblyIssueIntake([]Record{
		&Artifact{
			CommonNode:   common("artifact_without_issue_ref", TypeArtifact, "verified"),
			ArtifactType: "report",
			Path:         strPtr("site://ops/civilization/no-issue-ref"),
			ContentHash:  strPtr("sha256:no-issue-ref"),
		},
	})

	if projection.Status != CivilizationAssemblyFieldUnavailable {
		t.Fatalf("status = %s, want unavailable", projection.Status)
	}
	if len(projection.Issues) != 0 || len(projection.Groups) != 0 {
		t.Fatalf("unavailable issue intake should not project issues/groups: %+v", projection)
	}
	for _, boundary := range []string{"scanner_read_only", "no_eventgraph_writes", "no_github_issue_mutation"} {
		if !containsString(projection.ScannerBoundaries, boundary) {
			t.Fatalf("unavailable projection missing scanner boundary %s: %+v", boundary, projection.ScannerBoundaries)
		}
	}
}

func TestCivilizationAssemblyIssueIntakeSkipsNilRecords(t *testing.T) {
	var typedNilArtifact *Artifact
	projection := civilizationAssemblyIssueIntake([]Record{
		nil,
		typedNilArtifact,
		&Requirement{
			CommonNode: commonWithSourceRefs("req_with_issue_ref", TypeRequirement, "accepted", []string{
				"github:transpara-ai/site#115",
			}),
			RiskClass: "high",
		},
	})

	if len(projection.Issues) != 1 || !containsCivilizationIssue(projection.Issues, "transpara-ai/site", 115) {
		t.Fatalf("nil records should be skipped without dropping valid issue refs: %+v", projection)
	}
}

func TestCivilizationAssemblyIssueIntakeAggregatesDuplicateIssueRefs(t *testing.T) {
	projection := civilizationAssemblyIssueIntake([]Record{
		&Requirement{
			CommonNode: commonWithSourceRefs("a_requirement_issue_ref", TypeRequirement, "accepted", []string{
				"github:transpara-ai/site#115",
			}),
			RiskClass: "low",
		},
		&FactoryOrder{
			CommonNode: commonWithSourceRefs("z_factory_order_issue_ref", TypeFactoryOrder, "draft", []string{
				"https://github.com/transpara-ai/site/issues/115",
			}),
			SourceIntentRef: "github:transpara-ai/site#115",
			RiskClass:       "critical",
		},
	})

	if len(projection.Issues) != 1 {
		t.Fatalf("deduped issue count = %d, want 1: %+v", len(projection.Issues), projection.Issues)
	}
	issue := projection.Issues[0]
	if issue.Repo != "transpara-ai/site" ||
		issue.Number != 115 ||
		issue.TouchedSubstrate != "multiple" ||
		issue.RiskClass != "critical" ||
		issue.Readiness != "mixed" ||
		!containsString(issue.TouchedSubstrates, TypeFactoryOrder) ||
		!containsString(issue.TouchedSubstrates, TypeRequirement) ||
		!containsString(issue.RiskClasses, "critical") ||
		!containsString(issue.RiskClasses, "low") ||
		!containsString(issue.ReadinessStates, "accepted") ||
		!containsString(issue.ReadinessStates, "draft") ||
		!containsString(issue.SourceRefs, "github:transpara-ai/site#115") ||
		!containsString(issue.SourceRefs, "https://github.com/transpara-ai/site/issues/115") {
		t.Fatalf("duplicate issue refs were not aggregated conservatively: %+v", issue)
	}
	if len(projection.Groups) != 1 ||
		projection.Groups[0].GroupID != "repo-transpara-ai-site-substrate-multiple-risk-critical-readiness-mixed" ||
		projection.Groups[0].TouchedSubstrate != "multiple" ||
		projection.Groups[0].RiskClass != "critical" ||
		projection.Groups[0].Readiness != "mixed" ||
		!containsCivilizationIssueRef(projection.Groups[0].IssueRefs, "transpara-ai/site", 115) {
		t.Fatalf("deduped issue group does not preserve aggregated fields: %+v", projection.Groups)
	}
}

func TestCivilizationAssemblyIssueIntakeSurfacesUnknownRiskSeparately(t *testing.T) {
	projection := civilizationAssemblyIssueIntake([]Record{
		&Requirement{
			CommonNode: commonWithSourceRefs("known_low_issue_ref", TypeRequirement, "accepted", []string{
				"github:transpara-ai/site#115",
			}),
			RiskClass: "low",
		},
		&Requirement{
			CommonNode: commonWithSourceRefs("future_risk_issue_ref", TypeRequirement, "accepted", []string{
				"github:transpara-ai/site#115",
			}),
			RiskClass: "blocker",
		},
	})

	if len(projection.Issues) != 1 {
		t.Fatalf("deduped issue count = %d, want 1: %+v", len(projection.Issues), projection.Issues)
	}
	issue := projection.Issues[0]
	if issue.RiskClass != "low" ||
		!containsString(issue.RiskClasses, "blocker") ||
		!containsString(issue.RiskClasses, "low") ||
		!containsString(issue.UnrecognizedRisk, "blocker") {
		t.Fatalf("unknown risk should be retained separately without replacing canonical headline: %+v", issue)
	}
}

func TestCivilizationAssemblyIssueIntakeUnknownRiskWithoutKnownClass(t *testing.T) {
	projection := civilizationAssemblyIssueIntake([]Record{
		&Requirement{
			CommonNode: commonWithSourceRefs("future_only_risk_issue_ref", TypeRequirement, "accepted", []string{
				"github:transpara-ai/site#115",
			}),
			RiskClass: "blocker",
		},
	})

	if len(projection.Issues) != 1 {
		t.Fatalf("issue count = %d, want 1: %+v", len(projection.Issues), projection.Issues)
	}
	issue := projection.Issues[0]
	if issue.RiskClass != "unknown" || !containsString(issue.UnrecognizedRisk, "blocker") {
		t.Fatalf("unknown-only risk should use canonical unknown headline and retain raw term: %+v", issue)
	}
}

func TestCivilizationAssemblyIssueIntakeGroupsMultipleIssues(t *testing.T) {
	projection := civilizationAssemblyIssueIntake([]Record{
		&Requirement{
			CommonNode: commonWithSourceRefs("req_site_issue_115", TypeRequirement, "accepted", []string{
				"github:transpara-ai/site#115",
			}),
			RiskClass: "high",
		},
		&Requirement{
			CommonNode: commonWithSourceRefs("req_site_issue_116", TypeRequirement, "accepted", []string{
				"https://github.com/transpara-ai/site/issues/116",
			}),
			RiskClass: "high",
		},
	})

	if len(projection.Issues) != 2 || len(projection.Groups) != 1 {
		t.Fatalf("projection = %+v, want two issues in one group", projection)
	}
	group := projection.Groups[0]
	if group.GroupID != "repo-transpara-ai-site-substrate-requirement-risk-high-readiness-accepted" ||
		group.Summary != "2 issue(s) share repo/substrate/risk/readiness key." ||
		len(group.IssueRefs) != 2 ||
		!containsCivilizationIssueRef(group.IssueRefs, "transpara-ai/site", 115) ||
		!containsCivilizationIssueRef(group.IssueRefs, "transpara-ai/site", 116) ||
		!containsString(group.SourceRefs, "req_site_issue_115") ||
		!containsString(group.SourceRefs, "req_site_issue_116") {
		t.Fatalf("multi-issue group missing expected aggregation: %+v", group)
	}
}

func TestCivilizationAssemblyIssueIntakeGroupsEmptyRisk(t *testing.T) {
	projection := civilizationAssemblyIssueIntake([]Record{
		&Artifact{
			CommonNode: commonWithSourceRefs("artifact_issue_ref", TypeArtifact, "verified", []string{
				"github:transpara-ai/site#115",
			}),
			ArtifactType: "report",
		},
	})

	if len(projection.Issues) != 1 || len(projection.Groups) != 1 {
		t.Fatalf("projection = %+v, want one issue and one group", projection)
	}
	if projection.Issues[0].RiskClass != "" || len(projection.Issues[0].RiskClasses) != 0 {
		t.Fatalf("non-risk record should not synthesize risk: %+v", projection.Issues[0])
	}
	if projection.Groups[0].GroupID != "repo-transpara-ai-site-substrate-artifact-risk-readiness-verified" ||
		projection.Groups[0].RiskClass != "" ||
		projection.Groups[0].Readiness != "verified" {
		t.Fatalf("empty-risk group should preserve labeled empty risk axis: %+v", projection.Groups[0])
	}
}

func TestParseCivilizationGitHubIssueRefAcceptsCanonicalRefs(t *testing.T) {
	for _, ref := range []string{
		"github:transpara-ai/site#115",
		"https://github.com/transpara-ai/site/issues/115",
	} {
		repo, number, issueURL, ok := parseCivilizationGitHubIssueRef(ref)
		if !ok || repo != "transpara-ai/site" || number != 115 || issueURL != "https://github.com/transpara-ai/site/issues/115" {
			t.Fatalf("parseCivilizationGitHubIssueRef(%q) = %s #%d %s %t, want canonical site#115", ref, repo, number, issueURL, ok)
		}
	}
}

func TestParseCivilizationGitHubIssueRefRejectsNonIssueRefs(t *testing.T) {
	for _, ref := range []string{
		"",
		"github:transpara-ai/site",
		"github:transpara-ai/site#0",
		"github:transpara-ai/site#+5",
		"github:Transpara-AI/site#115",
		"github:transpara-ai/site#0115",
		"github:transpara-ai/docs#172#5",
		"github:transpara-ai/site name#5",
		"github:transpara-ai/..#5",
		"github:../eventgraph#5",
		"https://attacker@github.com/transpara-ai/site/issues/5",
		"https://github.com/transpara-ai/site/pull/125",
		"https://github.com/transpara-ai/site/issues/+5",
		"https://github.com/transpara-ai/site/issues/115?source=tracker",
		"https://github.com/transpara-ai/site/issues/115#issuecomment-1",
		"https://github.com/transpara-ai/site/issues/115/",
		"https://github.com//transpara-ai/site/issues/115",
		"https://github.com/transpara-ai/site/issues/%31%31%35",
		"https://github.com/transpara-ai/../issues/5",
		"https://example.com/transpara-ai/site/issues/117",
		"issue://117",
	} {
		if repo, number, issueURL, ok := parseCivilizationGitHubIssueRef(ref); ok {
			t.Fatalf("parseCivilizationGitHubIssueRef(%q) = %s #%d %s, want rejected", ref, repo, number, issueURL)
		}
	}
}

func TestCivilizationAssemblyRecordRiskClassSources(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   Record
		want string
	}{
		{name: "factory order", in: &FactoryOrder{CommonNode: common("fo_risk", TypeFactoryOrder, "accepted"), RiskClass: "critical"}, want: "critical"},
		{name: "requirement", in: &Requirement{CommonNode: common("req_risk", TypeRequirement, "accepted"), RiskClass: "high"}, want: "high"},
		{name: "acceptance criterion", in: &AcceptanceCriterion{CommonNode: common("ac_risk", TypeAcceptanceCriterion, "accepted"), RiskClass: "medium"}, want: "medium"},
		{name: "assumption", in: &Assumption{CommonNode: common("asm_risk", TypeAssumption, "accepted"), RiskClass: "low"}, want: "low"},
		{name: "task", in: &Task{CommonNode: common("tsk_risk", TypeTask, "accepted"), RiskClass: "high"}, want: "high"},
		{name: "authority request", in: &AuthorityRequest{CommonNode: common("auth_req_risk", TypeAuthorityRequest, "requested"), RiskClass: "critical"}, want: "critical"},
		{name: "failure severity", in: &Failure{CommonNode: common("failure_risk", TypeFailure, "open"), Severity: "high"}, want: "high"},
		{name: "non-risk record", in: &Artifact{CommonNode: common("artifact_risk", TypeArtifact, "verified")}, want: ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := civilizationAssemblyRecordRiskClass(tc.in); got != tc.want {
				t.Fatalf("risk class = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestCivilizationAssemblyProjectionOptionsOverrideProjectionIDAndMergeBoundaryFlags(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{
		ProjectionID:  "custom_projection_id",
		GeneratedAt:   fixedTime,
		BoundaryFlags: []string{"custom_boundary", "no_deploy"},
	})

	if projection.ProjectionID != "custom_projection_id" {
		t.Fatalf("projection id = %s", projection.ProjectionID)
	}
	if !containsString(projection.BoundaryFlags, "custom_boundary") || !containsString(projection.BoundaryFlags, "no_deploy") {
		t.Fatalf("boundary flags were not merged/deduped: %+v", projection.BoundaryFlags)
	}
}

func TestCivilizationAssemblyProjectionMissingAuthorityFailsClosed(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	deleteRecord(store, "auth_dec_civ_001")
	deleteRecord(store, "exec_civ_001")
	store.records["lt_civ_001"].(*LifecycleTransition).AuthorityDecisionID = nil
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

func TestCivilizationAssemblyProjectionDanglingAuthorityDecisionReferenceFails(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	store.mu.Lock()
	deleteRecord(store, "auth_req_civ_001")
	store.mu.Unlock()

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if projection.DerivationStatus != CivilizationAssemblyDerivationFailed {
		t.Fatalf("derivation status = %s, want failed", projection.DerivationStatus)
	}
	if projection.AuthorityState.Status != CivilizationAssemblyFieldUnavailable {
		t.Fatalf("authority state = %+v, want unavailable", projection.AuthorityState)
	}
	if !containsFailureReason(projection.FailureReasons, "AuthorityDecision auth_dec_civ_001 references missing AuthorityRequest auth_req_civ_001") {
		t.Fatalf("missing dangling AuthorityRequest failure reason: %+v", projection.FailureReasons)
	}
}

func TestCivilizationAssemblyProjectionDanglingExecutionReceiptReferenceFails(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	store.mu.Lock()
	deleteRecord(store, "auth_dec_civ_001")
	store.mu.Unlock()

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if projection.DerivationStatus != CivilizationAssemblyDerivationFailed {
		t.Fatalf("derivation status = %s, want failed", projection.DerivationStatus)
	}
	if !containsFailureReason(projection.FailureReasons, "ExecutionReceipt exec_civ_001 references missing AuthorityDecision auth_dec_civ_001") {
		t.Fatalf("missing dangling AuthorityDecision failure reason: %+v", projection.FailureReasons)
	}
}

func TestCivilizationAssemblyProjectionDanglingLifecycleAuthorityReferenceFails(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	store.mu.Lock()
	*store.records["lt_civ_001"].(*LifecycleTransition).AuthorityDecisionID = "auth_dec_civ_missing"
	store.mu.Unlock()

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if projection.DerivationStatus != CivilizationAssemblyDerivationFailed {
		t.Fatalf("derivation status = %s, want failed", projection.DerivationStatus)
	}
	if projection.AuthorityState.Status != CivilizationAssemblyFieldUnavailable {
		t.Fatalf("authority state = %+v, want unavailable", projection.AuthorityState)
	}
	if !projectionHasUnavailableField(projection, "authority_state") {
		t.Fatalf("dangling lifecycle authority ref should mark authority_state unavailable: %+v", projection.WithheldOrUnavailableFields)
	}
	if !containsFailureReason(projection.FailureReasons, "LifecycleTransition lt_civ_001 references missing AuthorityDecision auth_dec_civ_missing") {
		t.Fatalf("missing dangling lifecycle authority failure reason: %+v", projection.FailureReasons)
	}
}

func TestCivilizationAssemblyProjectionDanglingHumanApprovalReferenceFails(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	store.mu.Lock()
	store.records["approval_civ_001"].(*HumanApproval).RequestRef = "auth_req_civ_missing"
	store.mu.Unlock()

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if projection.DerivationStatus != CivilizationAssemblyDerivationFailed {
		t.Fatalf("derivation status = %s, want failed", projection.DerivationStatus)
	}
	if projection.AuthorityState.Status != CivilizationAssemblyFieldUnavailable {
		t.Fatalf("authority state = %+v, want unavailable", projection.AuthorityState)
	}
	if !projectionHasUnavailableField(projection, "authority_state") {
		t.Fatalf("dangling approval authority ref should mark authority_state unavailable: %+v", projection.WithheldOrUnavailableFields)
	}
	if !containsFailureReason(projection.FailureReasons, "HumanApproval approval_civ_001 references missing AuthorityRequest auth_req_civ_missing") {
		t.Fatalf("missing dangling approval authority failure reason: %+v", projection.FailureReasons)
	}
}

func TestCivilizationAssemblyProjectionMissingAuthorityRequestIDDoesNotCreateEmptyConflict(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	appendRecord(t, store, &AuthorityDecision{
		CommonNode:         common("auth_dec_civ_missing_request_002", TypeAuthorityDecision, "approved"),
		AuthorityRequestID: "auth_req_civ_001",
		DeciderActorID:     "act_human",
		DeciderRole:        "External Committee",
		Decision:           "Forbidden",
		Reason:             "negative missing request fixture",
		Scope:              []string{"eventgraph.read.projection"},
	})
	store.mu.Lock()
	store.records["auth_dec_civ_001"].(*AuthorityDecision).AuthorityRequestID = ""
	store.records["auth_dec_civ_missing_request_002"].(*AuthorityDecision).AuthorityRequestID = ""
	store.mu.Unlock()

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if projection.DerivationStatus != CivilizationAssemblyDerivationFailed {
		t.Fatalf("derivation status = %s, want failed", projection.DerivationStatus)
	}
	if !containsFailureReason(projection.FailureReasons, "AuthorityDecision auth_dec_civ_001 is missing authority_request_id") {
		t.Fatalf("missing empty authority_request_id failure reason: %+v", projection.FailureReasons)
	}
	for _, reason := range projection.FailureReasons {
		if strings.Contains(reason, "conflicting AuthorityDecision records for :") {
			t.Fatalf("empty authority_request_id should not create blank conflict reason: %+v", projection.FailureReasons)
		}
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

func TestCivilizationAssemblyProjectionUnresolvedCriticalContradictionFails(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	appendRecord(t, store, &ContradictionLog{
		CommonNode:      common("contradiction_civ_accepted_conflict", TypeContradictionLog, "accepted_conflict"),
		ContradictionID: "contradiction_civ_accepted_conflict",
		ClaimARef:       "auth_dec_civ_001",
		ClaimBRef:       "approval_civ_001",
		Severity:        "critical",
	})

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if projection.DerivationStatus != CivilizationAssemblyDerivationFailed {
		t.Fatalf("derivation status = %s, want failed", projection.DerivationStatus)
	}
	if !containsFailureReason(projection.FailureReasons, "unresolved critical contradiction contradiction_civ_accepted_conflict blocks trusted projection") {
		t.Fatalf("missing unresolved critical contradiction failure reason: %+v", projection.FailureReasons)
	}
}

func TestCivilizationAssemblyProjectionUnresolvedHighContradictionFails(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	appendRecord(t, store, &ContradictionLog{
		CommonNode:      common("contradiction_civ_high", TypeContradictionLog, "open"),
		ContradictionID: "contradiction_civ_high",
		ClaimARef:       "auth_dec_civ_001",
		ClaimBRef:       "approval_civ_001",
		Severity:        "high",
	})

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if projection.DerivationStatus != CivilizationAssemblyDerivationFailed {
		t.Fatalf("derivation status = %s, want failed", projection.DerivationStatus)
	}
	if !containsFailureReason(projection.FailureReasons, "unresolved high contradiction contradiction_civ_high blocks trusted projection") {
		t.Fatalf("missing unresolved high contradiction failure reason: %+v", projection.FailureReasons)
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

func TestCivilizationAssemblyProjectionNormalizesStatusComparisons(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	appendRecord(t, store, &ContradictionLog{
		CommonNode:      common("contradiction_civ_mixed_case", TypeContradictionLog, "open"),
		ContradictionID: "contradiction_civ_mixed_case",
		ClaimARef:       "auth_dec_civ_001",
		ClaimBRef:       "approval_civ_001",
		Severity:        "critical",
	})
	store.mu.Lock()
	*store.records["gate_civ_pass"].(*GateResult).CommonNode.Status = "PASS"
	*store.records["contradiction_civ_mixed_case"].(*ContradictionLog).CommonNode.Status = "OPEN"
	store.records["contradiction_civ_mixed_case"].(*ContradictionLog).Severity = "Critical"
	store.mu.Unlock()

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	if len(projection.OpenGateSummary) != 0 {
		t.Fatalf("mixed-case PASS gate should not be open: %+v", projection.OpenGateSummary)
	}
	if projection.DerivationStatus != CivilizationAssemblyDerivationFailed {
		t.Fatalf("mixed-case critical contradiction should fail projection: %+v", projection)
	}
	if len(projection.FailureReasons) == 0 || !strings.Contains(projection.FailureReasons[0], "critical contradiction") {
		t.Fatalf("missing normalized contradiction failure: %+v", projection.FailureReasons)
	}
}

func TestCivilizationAssemblyProjectionSkipsEmptyActorRoleBindings(t *testing.T) {
	store := civilizationAssemblyProjectionStore(t)
	store.mu.Lock()
	store.records["auth_dec_civ_001"].(*AuthorityDecision).DeciderActorID = ""
	store.mu.Unlock()

	projection := store.ProjectCivilizationAssembly(CivilizationAssemblyProjectionOptions{GeneratedAt: fixedTime})

	for _, binding := range projection.RoleBindings {
		if binding.SourceRef == "auth_dec_civ_001" {
			t.Fatalf("empty decider actor should not produce role binding: %+v", projection.RoleBindings)
		}
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
		SourceIntentRef:     "github:transpara-ai/docs#172",
		RiskClass:           "high",
		ReleasePolicy:       "human_approval_required",
	})
	appendRecord(t, store, &Requirement{
		CommonNode:     commonWithSourceRefs(reqID, TypeRequirement, "accepted", []string{"https://github.com/transpara-ai/site/issues/115"}),
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

func appendCivilizationProjectionStoreAuthority(t *testing.T, store *InMemoryStore) string {
	t.Helper()
	return appendCivilizationProjectionStoreAuthorityDecision(t, store, "auth_dec_civ_projection_store", "Autonomous", nil)
}

func appendCivilizationProjectionStoreAuthorityDecision(t *testing.T, store *InMemoryStore, decisionID, decision string, expiresAt *time.Time) string {
	t.Helper()
	requestID := "auth_req_civ_projection_store"
	appendRecord(t, store, &AuthorityRequest{
		CommonNode:   common(requestID, TypeAuthorityRequest, "open"),
		ActorID:      "act_agent",
		ActorRole:    "ProjectionStoreBuilder",
		Action:       CivilizationAssemblyProjectionStoreAction,
		TargetType:   CivilizationAssemblyProjectionSubject,
		TargetID:     "fo_civ_001",
		RiskClass:    "high",
		Reason:       "record local Civilization Assembly projection-store truth",
		EvidenceRefs: []string{"auth_dec_civ_001", "tr_civ_001", "gate_civ_pass"},
	})
	appendRecord(t, store, &AuthorityDecision{
		CommonNode:         common(decisionID, TypeAuthorityDecision, "approved"),
		AuthorityRequestID: requestID,
		DeciderActorID:     "act_human",
		DeciderRole:        "External Committee",
		Decision:           decision,
		Reason:             "bounded local projection-store fixture only",
		Scope:              []string{CivilizationAssemblyProjectionStoreAction},
		Conditions:         []string{"projection_store_local_only", "no_production_eventgraph_write", "no_runtime_execution"},
		ExpiresAt:          expiresAt,
	})
	return decisionID
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

func containsFailureReason(reasons []string, want string) bool {
	for _, reason := range reasons {
		if strings.Contains(reason, want) {
			return true
		}
	}
	return false
}

func containsCivilizationIssue(issues []CivilizationAssemblyIssueIntakeIssue, repo string, number int) bool {
	for _, issue := range issues {
		if issue.Repo == repo && issue.Number == number {
			return true
		}
	}
	return false
}

func containsCivilizationIssueRef(issues []CivilizationAssemblyIssueRef, repo string, number int) bool {
	for _, issue := range issues {
		if issue.Repo == repo && issue.Number == number {
			return true
		}
	}
	return false
}

func civilizationAssemblyStoreFootprint(t *testing.T, store *InMemoryStore) []string {
	t.Helper()
	store.mu.RLock()
	defer store.mu.RUnlock()

	var footprint []string
	for _, id := range sortedMapKeysRecord(store.records) {
		canonical, err := CanonicalJSON(store.records[id])
		if err != nil {
			t.Fatalf("canonical record %s: %v", id, err)
		}
		footprint = append(footprint, "record:"+id+":"+string(canonical))
	}
	for _, id := range sortedMapKeysEdge(store.edges) {
		canonical, err := CanonicalJSON(store.edges[id])
		if err != nil {
			t.Fatalf("canonical edge %s: %v", id, err)
		}
		footprint = append(footprint, "edge:"+id+":"+string(canonical))
	}
	return footprint
}
