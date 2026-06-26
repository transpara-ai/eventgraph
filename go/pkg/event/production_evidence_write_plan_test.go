package event

import (
	"strings"
	"testing"
)

func TestPlanProductionEvidenceWriteAuthorizesLocalAppend(t *testing.T) {
	candidate := baseProductionEvidenceWriteCandidate()

	plan, err := PlanProductionEvidenceWrite(candidate)
	if err != nil {
		t.Fatalf("PlanProductionEvidenceWrite: %v", err)
	}
	if plan.CandidateID != candidate.CandidateID {
		t.Fatalf("CandidateID = %q, want %q", plan.CandidateID, candidate.CandidateID)
	}
	if len(plan.RuntimeEvidenceRefs) != 1 || plan.RuntimeEvidenceRefs[0] != candidate.RuntimeEvidenceRefs[0] {
		t.Fatalf("RuntimeEvidenceRefs = %#v, want %#v", plan.RuntimeEvidenceRefs, candidate.RuntimeEvidenceRefs)
	}
	if len(plan.IssueEvidenceRefs) != 1 || plan.IssueEvidenceRefs[0] != candidate.IssueEvidenceRefs[0] {
		t.Fatalf("IssueEvidenceRefs = %#v, want %#v", plan.IssueEvidenceRefs, candidate.IssueEvidenceRefs)
	}
	if len(plan.Entries) != 5 {
		t.Fatalf("Entries length = %d, want 5", len(plan.Entries))
	}

	wantOrder := []string{
		EventTypeAuthorityDecisionRecorded.Value(),
		EventTypeAuthorityStoreGovernanceRecorded.Value(),
		EventTypeNativeTestRunRecorded.Value(),
		EventTypeNativeGateResultRecorded.Value(),
		EventTypeNativeAuditReportRecorded.Value(),
	}
	for i, want := range wantOrder {
		if plan.Entries[i].EventType != want {
			t.Fatalf("Entries[%d].EventType = %q, want %q", i, plan.Entries[i].EventType, want)
		}
	}

	applied, err := ApplyProductionEvidenceWritePlanToFixture(nil, plan)
	if err != nil {
		t.Fatalf("ApplyProductionEvidenceWritePlanToFixture: %v", err)
	}
	if len(applied) != len(plan.Entries) {
		t.Fatalf("applied entries = %d, want %d", len(applied), len(plan.Entries))
	}

	replayed, err := ApplyProductionEvidenceWritePlanToFixture(applied, plan)
	if err != nil {
		t.Fatalf("ApplyProductionEvidenceWritePlanToFixture replay: %v", err)
	}
	if len(replayed) != len(applied) {
		t.Fatalf("replay appended duplicates: got %d entries, want %d", len(replayed), len(applied))
	}
}

func TestPlanProductionEvidenceWriteFailsClosed(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*ProductionEvidenceWriteCandidate)
		wantErr string
	}{
		{
			name: "missing runtime evidence",
			mutate: func(candidate *ProductionEvidenceWriteCandidate) {
				candidate.RuntimeEvidenceRefs = nil
			},
			wantErr: "runtime_evidence_ref",
		},
		{
			name: "missing issue evidence",
			mutate: func(candidate *ProductionEvidenceWriteCandidate) {
				candidate.IssueEvidenceRefs = nil
			},
			wantErr: "issue_evidence_ref",
		},
		{
			name: "missing authority decision",
			mutate: func(candidate *ProductionEvidenceWriteCandidate) {
				candidate.AuthorityDecision = AuthorityDecisionRecordedContent{}
			},
			wantErr: "authority decision",
		},
		{
			name: "approval required is not write authority",
			mutate: func(candidate *ProductionEvidenceWriteCandidate) {
				candidate.AuthorityDecision.Outcome = AuthorityDecisionOutcomeApprovalRequired
			},
			wantErr: "outcome must be",
		},
		{
			name: "wrong actor",
			mutate: func(candidate *ProductionEvidenceWriteCandidate) {
				candidate.AuthorityDecision.ActorRef = "agent:hive"
			},
			wantErr: "actor_ref",
		},
		{
			name: "wrong repo",
			mutate: func(candidate *ProductionEvidenceWriteCandidate) {
				candidate.TestRuns[0].SourceIssueRefs = []NativeEvidenceIssueRef{{Repo: "transpara-ai/work", Number: 59}}
			},
			wantErr: "source_issue_refs do not match candidate",
		},
		{
			name: "wrong action class",
			mutate: func(candidate *ProductionEvidenceWriteCandidate) {
				candidate.ActionClass = "runtime.external_adapter.invoke"
			},
			wantErr: "wrong action class",
		},
		{
			name: "invalid schema",
			mutate: func(candidate *ProductionEvidenceWriteCandidate) {
				candidate.SchemaVersion = "production_evidence_write_v99"
			},
			wantErr: "invalid schema",
		},
		{
			name: "stale source",
			mutate: func(candidate *ProductionEvidenceWriteCandidate) {
				candidate.CurrentStateRef = "github:transpara-ai/eventgraph#61@newer"
			},
			wantErr: "stale source",
		},
		{
			name: "gate lacks test-run evidence",
			mutate: func(candidate *ProductionEvidenceWriteCandidate) {
				candidate.GateResults[0].EvidenceRefs = []string{"artifact:cfar"}
			},
			wantErr: "must include a test run",
		},
		{
			name: "audit lacks cfar evidence",
			mutate: func(candidate *ProductionEvidenceWriteCandidate) {
				candidate.AuditReports[0].CFARRefs = nil
			},
			wantErr: "cfar_refs are required",
		},
		{
			name: "passed audit has missing links",
			mutate: func(candidate *ProductionEvidenceWriteCandidate) {
				candidate.AuditReports[0].MissingLinks = []string{"external_committee_approval"}
			},
			wantErr: "missing_links",
		},
		{
			name: "duplicate evidence in candidate",
			mutate: func(candidate *ProductionEvidenceWriteCandidate) {
				candidate.TestRuns = append(candidate.TestRuns, candidate.TestRuns[0])
			},
			wantErr: "duplicate evidence",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			candidate := baseProductionEvidenceWriteCandidate()
			tc.mutate(&candidate)
			_, err := PlanProductionEvidenceWrite(candidate)
			if err == nil {
				t.Fatalf("PlanProductionEvidenceWrite accepted invalid candidate")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tc.wantErr)
			}
		})
	}
}

func TestApplyProductionEvidenceWritePlanToFixtureRejectsConflictingExistingEvidence(t *testing.T) {
	plan, err := PlanProductionEvidenceWrite(baseProductionEvidenceWriteCandidate())
	if err != nil {
		t.Fatalf("PlanProductionEvidenceWrite: %v", err)
	}

	conflicting := plan.Entries[0]
	decision := conflicting.Content.(AuthorityDecisionRecordedContent)
	decision.EvidenceRefs = append(decision.EvidenceRefs, "artifact:conflicting-pr-visible-evidence")
	conflicting.Content = decision

	if _, err := ApplyProductionEvidenceWritePlanToFixture([]ProductionEvidenceWriteEntry{conflicting}, plan); err == nil {
		t.Fatalf("ApplyProductionEvidenceWritePlanToFixture accepted conflicting existing evidence")
	} else if !strings.Contains(err.Error(), "conflicting evidence") {
		t.Fatalf("error = %q, want conflicting evidence", err.Error())
	}
}

func baseProductionEvidenceWriteCandidate() ProductionEvidenceWriteCandidate {
	issueRef := NativeEvidenceIssueRef{
		Repo:   "transpara-ai/eventgraph",
		Number: 61,
		URL:    "https://github.com/transpara-ai/eventgraph/issues/61",
		Title:  "Production EventGraph write path for runtime and issue evidence",
	}
	decisionID := "decision_eventgraph_61_production_write"
	testRunID := "testrun_eventgraph_61_make_verify"
	gateID := "gate_eventgraph_61_cfar"

	return ProductionEvidenceWriteCandidate{
		CandidateID:     "candidate_eventgraph_61_production_write",
		Repo:            "transpara-ai/eventgraph",
		ActorRef:        "human:michael",
		ActionClass:     ProductionEvidenceWriteActionClass,
		SchemaVersion:   ProductionEvidenceWriteSchemaVersion,
		SourceStateRef:  "github:transpara-ai/eventgraph#61@4811902951",
		CurrentStateRef: "github:transpara-ai/eventgraph#61@4811902951",
		SourceIssueRefs: []NativeEvidenceIssueRef{issueRef},
		RuntimeEvidenceRefs: []string{
			"runtime:governed-observation-envelope:not-executed",
		},
		IssueEvidenceRefs: []string{
			"github:transpara-ai/eventgraph#61@4811902951",
		},
		AuthorityDecision: AuthorityDecisionRecordedContent{
			DecisionID:      decisionID,
			SubjectType:     ProductionEvidenceWriteSubjectType,
			SubjectRef:      "transpara-ai/eventgraph",
			Outcome:         AuthorityDecisionOutcomeAutonomous,
			ActorRef:        "human:michael",
			SourceIssueRefs: []NativeEvidenceIssueRef{issueRef},
			AuthorityRefs: []string{
				"docs:DF-V4.0-EPIC-017-AUTHORITY-DECISION",
				"github:transpara-ai/docs@ad7ecdf",
			},
			EvidenceRefs: []string{
				"github:transpara-ai/eventgraph#61",
				"docs:dark-factory/v4.0/implementation/epics/epic-17-production-eventgraph-runtime-wiring/03-production-eventgraph-runtime-wiring-authority-decision-v4.0.md",
			},
			NonClaimRefs: []string{
				"no_live_production_write",
				"no_runtime_execution",
				"no_hive_start",
			},
			RecordedAt: "2026-06-26T18:00:00Z",
		},
		StoreGovernance: AuthorityStoreGovernanceRecordedContent{
			GovernanceID:           "storegov_eventgraph_61_production_evidence",
			StoreName:              ProductionEvidenceWriteStoreName,
			SchemaVersion:          ProductionEvidenceWriteSchemaVersion,
			WriteStatus:            AuthorityStoreWriteStatusWritePathAuthorized,
			RequiredValidationRefs: []string{"validation:make-verify", "validation:go-test-pkg-event"},
			AuthorityRefs: []string{
				decisionID,
				"docs:DF-V4.0-EPIC-017-AUTHORITY-DECISION",
			},
			SourceIssueRefs: []NativeEvidenceIssueRef{issueRef},
			EvidenceRefs:    []string{"github:transpara-ai/eventgraph#61"},
			RecordedAt:      "2026-06-26T18:01:00Z",
		},
		TestRuns: []NativeTestRunRecordedContent{
			{
				TestRunID:       testRunID,
				TestCaseID:      "eventgraph_61_production_write_semantics",
				Command:         "make verify",
				Outcome:         NativeEvidenceOutcomePassed,
				SourceIssueRefs: []NativeEvidenceIssueRef{issueRef},
				PRRefs:          []string{"transpara-ai/eventgraph#draft"},
				EvidenceRefs:    []string{"artifact:make-verify"},
				ValidationRefs:  []string{"validation:make-verify"},
				SourceRefs:      []string{"github:transpara-ai/eventgraph#61"},
			},
		},
		GateResults: []NativeGateResultRecordedContent{
			{
				GateResultID:    gateID,
				FactoryOrderID:  "fo_eventgraph_61",
				GateName:        "cross-family-adversarial-review",
				Outcome:         NativeEvidenceOutcomePassed,
				EvidenceRefs:    []string{testRunID, "cfar:eventgraph-pr-draft"},
				SourceIssueRefs: []NativeEvidenceIssueRef{issueRef},
				PRRefs:          []string{"transpara-ai/eventgraph#draft"},
				ValidationRefs:  []string{"validation:cfar-exact-head"},
				SourceRefs:      []string{"github:transpara-ai/eventgraph#61"},
			},
		},
		AuditReports: []NativeAuditReportRecordedContent{
			{
				AuditReportID:         "audit_eventgraph_61",
				TargetType:            "pull_request",
				TargetID:              "transpara-ai/eventgraph#draft",
				Outcome:               NativeEvidenceOutcomePassed,
				TraceScoreBasisPoints: productionEvidenceIntPtr(10000),
				SourceIssueRefs:       []NativeEvidenceIssueRef{issueRef},
				PRRefs:                []string{"transpara-ai/eventgraph#draft"},
				ValidationRefs:        []string{"validation:make-verify", "validation:go-test-pkg-event"},
				CFARRefs:              []string{"cfar:eventgraph-pr-draft"},
				AuthorityBoundaryRefs: []string{"docs:DF-V4.0-EPIC-017-AUTHORITY-DECISION"},
				ResidualRiskRefs:      []string{"none"},
				EvidenceRefs:          []string{testRunID, gateID},
				SourceRefs:            []string{"github:transpara-ai/eventgraph#61"},
			},
		},
	}
}

func productionEvidenceIntPtr(value int) *int {
	return &value
}
