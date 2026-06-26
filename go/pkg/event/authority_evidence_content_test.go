package event

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

func TestAuthorityEvidenceEventTypesRegistered(t *testing.T) {
	registry := DefaultRegistry()
	all := AllAuthorityEvidenceEventTypes()
	if len(all) != 4 {
		t.Fatalf("AllAuthorityEvidenceEventTypes returned %d types, want 4", len(all))
	}
	for _, eventType := range all {
		if !registry.IsRegistered(eventType) {
			t.Fatalf("DefaultRegistry missing %s", eventType.Value())
		}
		if !IsKnownEventType(eventType.Value()) {
			t.Fatalf("UnmarshalContent does not know %s", eventType.Value())
		}
	}
}

func TestAuthorityEvidenceContentRoundTrip(t *testing.T) {
	issueRef := NativeEvidenceIssueRef{
		Repo:   "transpara-ai/eventgraph",
		Number: 62,
		URL:    "https://github.com/transpara-ai/eventgraph/issues/62",
		Title:  "Authority evidence schema and store governance",
	}
	cases := []struct {
		name      string
		eventType string
		content   EventContent
	}{
		{
			name:      "authority decision",
			eventType: EventTypeAuthorityDecisionRecorded.Value(),
			content: AuthorityDecisionRecordedContent{
				DecisionID:       "decision_eventgraph_62_schema_only",
				SubjectType:      "github_issue",
				SubjectRef:       "transpara-ai/eventgraph#62",
				Outcome:          AuthorityDecisionOutcomeApprovalRequired,
				ActorRef:         "human:michael",
				ProtectedActions: []ProtectedAction{ProtectedActionRepoMergeMain},
				SourceIssueRefs:  []NativeEvidenceIssueRef{issueRef},
				AuthorityRefs:    []string{"github:transpara-ai/eventgraph#62"},
				EvidenceRefs:     []string{"docs:conformance/authority-evidence-governance.md"},
				NonClaimRefs:     []string{"no_runtime", "no_eventgraph_write"},
				RecordedAt:       "2026-06-26T00:00:00Z",
			},
		},
		{
			name:      "authority boundary",
			eventType: EventTypeAuthorityBoundaryRecorded.Value(),
			content: AuthorityBoundaryRecordedContent{
				BoundaryID:            "boundary_eventgraph_62_protected_action",
				BoundaryType:          AuthorityBoundaryProtectedAction,
				SubjectRef:            "transpara-ai/eventgraph#62",
				State:                 AuthorityBoundaryStateApprovalRequired,
				ProtectedActions:      []ProtectedAction{ProtectedActionRepoMergeMain, ProtectedActionProductionDeploy},
				RequiredAuthorityRefs: []string{"github:transpara-ai/eventgraph#62"},
				SourceIssueRefs:       []NativeEvidenceIssueRef{issueRef},
				EvidenceRefs:          []string{"docs:dark-factory/authority-vocabulary.md"},
				RecordedAt:            "2026-06-26T00:00:00Z",
			},
		},
		{
			name:      "authority residual",
			eventType: EventTypeAuthorityResidualRecorded.Value(),
			content: AuthorityResidualRecordedContent{
				ResidualID:      "residual_eventgraph_62_missing_runtime_authority",
				SubjectRef:      "transpara-ai/eventgraph#62",
				Status:          AuthorityResidualStatusCarried,
				Severity:        AuthorityResidualSeverityHigh,
				Description:     "production write path remains out of scope",
				RequiredAction:  "obtain governed authority before enabling runtime writes",
				SourceIssueRefs: []NativeEvidenceIssueRef{issueRef},
				AuthorityRefs:   []string{"github:transpara-ai/docs#197"},
				EvidenceRefs:    []string{"github:transpara-ai/eventgraph#62"},
				RecordedAt:      "2026-06-26T00:00:00Z",
			},
		},
		{
			name:      "store governance",
			eventType: EventTypeAuthorityStoreGovernanceRecorded.Value(),
			content: AuthorityStoreGovernanceRecordedContent{
				GovernanceID:           "storegov_eventgraph_62_projection",
				StoreName:              "civilization_authority_projection",
				SchemaVersion:          "authority_evidence_v0",
				WriteStatus:            AuthorityStoreWriteStatusProjectionOnly,
				MigrationRefs:          []string{"migration:not_authorized"},
				RequiredValidationRefs: []string{"go test ./pkg/event -run TestAuthorityEvidence"},
				AuthorityRefs:          []string{"github:transpara-ai/eventgraph#62"},
				SourceIssueRefs:        []NativeEvidenceIssueRef{issueRef},
				EvidenceRefs:           []string{"docs:conformance/authority-evidence-governance.md"},
				RecordedAt:             "2026-06-26T00:00:00Z",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.content)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			got, err := UnmarshalContent(tc.eventType, data)
			if err != nil {
				t.Fatalf("UnmarshalContent: %v", err)
			}
			if !reflect.DeepEqual(got, tc.content) {
				t.Fatalf("round trip mismatch\n got: %#v\nwant: %#v", got, tc.content)
			}
			if got.EventTypeName() != tc.eventType {
				t.Fatalf("EventTypeName = %q, want %q", got.EventTypeName(), tc.eventType)
			}
		})
	}
}

func TestAuthorityEvidenceRejectsInvalidEnumsAndUnguardedWrites(t *testing.T) {
	registry := DefaultRegistry()
	cases := []struct {
		name      string
		eventType types.EventType
		json      string
		content   EventContent
	}{
		{
			name:      "decision outcome",
			eventType: EventTypeAuthorityDecisionRecorded,
			json:      `{"decision_id":"decision_bad","subject_type":"issue","subject_ref":"eventgraph#62","outcome":"banana"}`,
			content: AuthorityDecisionRecordedContent{
				DecisionID:  "decision_bad",
				SubjectType: "issue",
				SubjectRef:  "eventgraph#62",
				Outcome:     AuthorityDecisionOutcome("banana"),
			},
		},
		{
			name:      "missing decision outcome",
			eventType: EventTypeAuthorityDecisionRecorded,
			json:      `{"decision_id":"decision_bad","subject_type":"issue","subject_ref":"eventgraph#62"}`,
			content: AuthorityDecisionRecordedContent{
				DecisionID:  "decision_bad",
				SubjectType: "issue",
				SubjectRef:  "eventgraph#62",
			},
		},
		{
			name:      "autonomous decision missing authority ref",
			eventType: EventTypeAuthorityDecisionRecorded,
			json:      `{"decision_id":"decision_bad","subject_type":"issue","subject_ref":"eventgraph#62","outcome":"autonomous"}`,
			content: AuthorityDecisionRecordedContent{
				DecisionID:  "decision_bad",
				SubjectType: "issue",
				SubjectRef:  "eventgraph#62",
				Outcome:     AuthorityDecisionOutcomeAutonomous,
			},
		},
		{
			name:      "approval-required decision missing authority ref",
			eventType: EventTypeAuthorityDecisionRecorded,
			json:      `{"decision_id":"decision_bad","subject_type":"issue","subject_ref":"eventgraph#62","outcome":"approval_required"}`,
			content: AuthorityDecisionRecordedContent{
				DecisionID:  "decision_bad",
				SubjectType: "issue",
				SubjectRef:  "eventgraph#62",
				Outcome:     AuthorityDecisionOutcomeApprovalRequired,
			},
		},
		{
			name:      "unknown protected action",
			eventType: EventTypeAuthorityBoundaryRecorded,
			json:      `{"boundary_id":"boundary_bad","boundary_type":"protected_action","subject_ref":"eventgraph#62","state":"approval_required","protected_actions":["deploy.production"]}`,
			content: AuthorityBoundaryRecordedContent{
				BoundaryID:       "boundary_bad",
				BoundaryType:     AuthorityBoundaryProtectedAction,
				SubjectRef:       "eventgraph#62",
				State:            AuthorityBoundaryStateApprovalRequired,
				ProtectedActions: []ProtectedAction{"deploy.production"},
			},
		},
		{
			name:      "protected boundary missing actions",
			eventType: EventTypeAuthorityBoundaryRecorded,
			json:      `{"boundary_id":"boundary_bad","boundary_type":"protected_action","subject_ref":"eventgraph#62","state":"blocked"}`,
			content: AuthorityBoundaryRecordedContent{
				BoundaryID:   "boundary_bad",
				BoundaryType: AuthorityBoundaryProtectedAction,
				SubjectRef:   "eventgraph#62",
				State:        AuthorityBoundaryStateBlocked,
			},
		},
		{
			name:      "authorized boundary missing authority ref",
			eventType: EventTypeAuthorityBoundaryRecorded,
			json:      `{"boundary_id":"boundary_bad","boundary_type":"merge","subject_ref":"eventgraph#62","state":"authorized"}`,
			content: AuthorityBoundaryRecordedContent{
				BoundaryID:   "boundary_bad",
				BoundaryType: AuthorityBoundaryMerge,
				SubjectRef:   "eventgraph#62",
				State:        AuthorityBoundaryStateAuthorized,
			},
		},
		{
			name:      "open residual missing action",
			eventType: EventTypeAuthorityResidualRecorded,
			json:      `{"residual_id":"residual_bad","subject_ref":"eventgraph#62","status":"open","severity":"high","description":"needs follow-up"}`,
			content: AuthorityResidualRecordedContent{
				ResidualID:  "residual_bad",
				SubjectRef:  "eventgraph#62",
				Status:      AuthorityResidualStatusOpen,
				Severity:    AuthorityResidualSeverityHigh,
				Description: "needs follow-up",
			},
		},
		{
			name:      "write path authorized missing authority and validation",
			eventType: EventTypeAuthorityStoreGovernanceRecorded,
			json:      `{"governance_id":"store_bad","store_name":"authority_projection","schema_version":"v0","write_status":"write_path_authorized"}`,
			content: AuthorityStoreGovernanceRecordedContent{
				GovernanceID:  "store_bad",
				StoreName:     "authority_projection",
				SchemaVersion: "v0",
				WriteStatus:   AuthorityStoreWriteStatusWritePathAuthorized,
			},
		},
		{
			name:      "invalid source issue ref",
			eventType: EventTypeAuthorityDecisionRecorded,
			json:      `{"decision_id":"decision_bad","subject_type":"issue","subject_ref":"eventgraph#62","outcome":"forbidden","source_issue_refs":[{"repo":"transpara-ai/eventgraph","number":0}]}`,
			content: AuthorityDecisionRecordedContent{
				DecisionID:      "decision_bad",
				SubjectType:     "issue",
				SubjectRef:      "eventgraph#62",
				Outcome:         AuthorityDecisionOutcomeForbidden,
				SourceIssueRefs: []NativeEvidenceIssueRef{{Repo: "transpara-ai/eventgraph"}},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := UnmarshalContent(tc.eventType.Value(), []byte(tc.json)); err == nil {
				t.Fatalf("UnmarshalContent accepted invalid %s", tc.name)
			}
			if err := registry.Validate(tc.eventType, tc.content); err == nil {
				t.Fatalf("registry accepted invalid %s", tc.name)
			}
		})
	}
}

func TestAuthorityResidualClosedWaivedAndNotApplicableAllowNoRequiredAction(t *testing.T) {
	cases := []AuthorityResidualRecordedContent{
		{
			ResidualID:  "residual_closed",
			SubjectRef:  "eventgraph#62",
			Status:      AuthorityResidualStatusClosed,
			Severity:    AuthorityResidualSeverityLow,
			Description: "closed residual is complete",
		},
		{
			ResidualID:  "residual_waived",
			SubjectRef:  "eventgraph#62",
			Status:      AuthorityResidualStatusWaived,
			Severity:    AuthorityResidualSeverityMedium,
			Description: "waived residual carries waiver evidence elsewhere",
		},
		{
			ResidualID:  "residual_not_applicable",
			SubjectRef:  "eventgraph#62",
			Status:      AuthorityResidualStatusNotApplicable,
			Severity:    AuthorityResidualSeverityLow,
			Description: "residual is not applicable to this subject",
		},
	}
	for _, tc := range cases {
		if err := DefaultRegistry().Validate(EventTypeAuthorityResidualRecorded, tc); err != nil {
			t.Fatalf("residual %s rejected without required_action: %v", tc.Status, err)
		}
	}
}

func TestAuthorityBoundaryBlockedAndNotApplicableAllowNoAuthorityRefs(t *testing.T) {
	cases := []AuthorityBoundaryRecordedContent{
		{
			BoundaryID:   "boundary_blocked",
			BoundaryType: AuthorityBoundaryEventGraphWrite,
			SubjectRef:   "eventgraph#62",
			State:        AuthorityBoundaryStateBlocked,
		},
		{
			BoundaryID:   "boundary_not_applicable",
			BoundaryType: AuthorityBoundaryDeployment,
			SubjectRef:   "eventgraph#62",
			State:        AuthorityBoundaryStateNotApplicable,
		},
	}
	for _, tc := range cases {
		if err := DefaultRegistry().Validate(EventTypeAuthorityBoundaryRecorded, tc); err != nil {
			t.Fatalf("boundary %s rejected without authority refs: %v", tc.State, err)
		}
	}
}

func TestAuthorityStoreGovernanceRequiresExplicitWriteAuthority(t *testing.T) {
	content := AuthorityStoreGovernanceRecordedContent{
		GovernanceID:           "storegov_eventgraph_62_write_path",
		StoreName:              "civilization_authority_projection",
		SchemaVersion:          "authority_evidence_v0",
		WriteStatus:            AuthorityStoreWriteStatusWritePathAuthorized,
		RequiredValidationRefs: []string{"validation:production-write-replay"},
		AuthorityRefs:          []string{"authority:human-approved-eventgraph-write-path"},
		SourceIssueRefs: []NativeEvidenceIssueRef{
			{Repo: "transpara-ai/eventgraph", Number: 62},
		},
	}
	if err := DefaultRegistry().Validate(EventTypeAuthorityStoreGovernanceRecorded, content); err != nil {
		t.Fatalf("valid write governance rejected: %v", err)
	}
}
