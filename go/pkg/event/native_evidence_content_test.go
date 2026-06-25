package event

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

func TestNativeEvidenceEventTypesRegistered(t *testing.T) {
	registry := DefaultRegistry()
	all := AllNativeEvidenceEventTypes()
	if len(all) != 3 {
		t.Fatalf("AllNativeEvidenceEventTypes returned %d types, want 3", len(all))
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

func TestNativeEvidenceContentRoundTrip(t *testing.T) {
	cases := []struct {
		name      string
		eventType string
		content   EventContent
	}{
		{
			name:      "test run",
			eventType: EventTypeNativeTestRunRecorded.Value(),
			content: NativeTestRunRecordedContent{
				TestRunID:         "testrun_work_63",
				TestCaseID:        "testcase_audit_report_traceability",
				ActorInvocationID: "actor_invocation_work_63",
				Command:           "go test ./pkg/event -run TestNativeEvidence",
				Outcome:           NativeEvidenceOutcomePassed,
				StartedAt:         "2026-06-25T19:52:00Z",
				CompletedAt:       "2026-06-25T19:53:00Z",
				SourceIssueRefs: []NativeEvidenceIssueRef{
					{Repo: "transpara-ai/eventgraph", Number: 63, URL: "https://github.com/transpara-ai/eventgraph/issues/63"},
				},
				PRRefs:         []string{"transpara-ai/eventgraph#999"},
				EvidenceRefs:   []string{"artifact:native_evidence_content_test"},
				ValidationRefs: []string{"validation:go-test-native-evidence"},
				SourceRefs:     []string{"github:transpara-ai/eventgraph#63"},
			},
		},
		{
			name:      "gate result",
			eventType: EventTypeNativeGateResultRecorded.Value(),
			content: NativeGateResultRecordedContent{
				GateResultID:       "gate_eventgraph_63_cfar",
				FactoryOrderID:     "fo_eventgraph_63",
				ReleaseCandidateID: "rc_eventgraph_63",
				GateName:           "cross-family-adversarial-review",
				Outcome:            NativeEvidenceOutcomePassed,
				EvidenceRefs:       []string{"cfar:eventgraph-pr-999"},
				SourceIssueRefs: []NativeEvidenceIssueRef{
					{Repo: "transpara-ai/eventgraph", Number: 63},
				},
				PRRefs:         []string{"transpara-ai/eventgraph#999"},
				ValidationRefs: []string{"validation:go-test-native-evidence"},
				SourceRefs:     []string{"github:transpara-ai/eventgraph#63"},
			},
		},
		{
			name:      "audit report",
			eventType: EventTypeNativeAuditReportRecorded.Value(),
			content: NativeAuditReportRecordedContent{
				AuditReportID:         "audit_eventgraph_63",
				TargetType:            "pull_request",
				TargetID:              "transpara-ai/eventgraph#999",
				Outcome:               NativeEvidenceOutcomePassed,
				TraceScoreBasisPoints: nativeEvidenceIntPtr(10000),
				SourceIssueRefs: []NativeEvidenceIssueRef{
					{Repo: "transpara-ai/eventgraph", Number: 63},
				},
				PRRefs:                []string{"transpara-ai/eventgraph#999"},
				ValidationRefs:        []string{"validation:go-test-native-evidence"},
				CFARRefs:              []string{"cfar:eventgraph-pr-999"},
				AuthorityBoundaryRefs: []string{"issue:eventgraph#63-pr-ready-boundary"},
				ResidualRiskRefs:      []string{"none"},
				EvidenceRefs:          []string{"testrun_work_63", "gate_eventgraph_63_cfar"},
				SourceRefs:            []string{"github:transpara-ai/eventgraph#63"},
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

func TestNativeEvidenceRejectsInvalidOutcomesAndTraceScores(t *testing.T) {
	registry := DefaultRegistry()
	cases := []struct {
		name      string
		eventType types.EventType
		json      string
		content   EventContent
	}{
		{
			name:      "test run outcome",
			eventType: EventTypeNativeTestRunRecorded,
			json:      `{"test_run_id":"tr_bad","command":"go test","outcome":"banana"}`,
			content: NativeTestRunRecordedContent{
				TestRunID: "tr_bad",
				Command:   "go test",
				Outcome:   NativeEvidenceOutcome("banana"),
			},
		},
		{
			name:      "missing test run outcome",
			eventType: EventTypeNativeTestRunRecorded,
			json:      `{"test_run_id":"tr_bad","command":"go test"}`,
			content: NativeTestRunRecordedContent{
				TestRunID: "tr_bad",
				Command:   "go test",
			},
		},
		{
			name:      "gate result outcome",
			eventType: EventTypeNativeGateResultRecorded,
			json:      `{"gate_result_id":"gate_bad","factory_order_id":"fo","gate_name":"cfar","outcome":"banana"}`,
			content: NativeGateResultRecordedContent{
				GateResultID:   "gate_bad",
				FactoryOrderID: "fo",
				GateName:       "cfar",
				Outcome:        NativeEvidenceOutcome("banana"),
			},
		},
		{
			name:      "audit report outcome",
			eventType: EventTypeNativeAuditReportRecorded,
			json:      `{"audit_report_id":"audit_bad","target_type":"pull_request","target_id":"pr","outcome":"banana","trace_score_basis_points":10000}`,
			content: NativeAuditReportRecordedContent{
				AuditReportID:         "audit_bad",
				TargetType:            "pull_request",
				TargetID:              "pr",
				Outcome:               NativeEvidenceOutcome("banana"),
				TraceScoreBasisPoints: nativeEvidenceIntPtr(10000),
			},
		},
		{
			name:      "audit report trace score",
			eventType: EventTypeNativeAuditReportRecorded,
			json:      `{"audit_report_id":"audit_bad","target_type":"pull_request","target_id":"pr","outcome":"passed","trace_score_basis_points":10001}`,
			content: NativeAuditReportRecordedContent{
				AuditReportID:         "audit_bad",
				TargetType:            "pull_request",
				TargetID:              "pr",
				Outcome:               NativeEvidenceOutcomePassed,
				TraceScoreBasisPoints: nativeEvidenceIntPtr(10001),
			},
		},
		{
			name:      "missing audit report trace score",
			eventType: EventTypeNativeAuditReportRecorded,
			json:      `{"audit_report_id":"audit_bad","target_type":"pull_request","target_id":"pr","outcome":"passed"}`,
			content: NativeAuditReportRecordedContent{
				AuditReportID: "audit_bad",
				TargetType:    "pull_request",
				TargetID:      "pr",
				Outcome:       NativeEvidenceOutcomePassed,
			},
		},
		{
			name:      "empty gate identity",
			eventType: EventTypeNativeGateResultRecorded,
			json:      `{"gate_result_id":"","factory_order_id":"fo","gate_name":"cfar","outcome":"passed"}`,
			content: NativeGateResultRecordedContent{
				FactoryOrderID: "fo",
				GateName:       "cfar",
				Outcome:        NativeEvidenceOutcomePassed,
			},
		},
		{
			name:      "invalid source issue ref",
			eventType: EventTypeNativeTestRunRecorded,
			json:      `{"test_run_id":"tr_bad","command":"go test","outcome":"passed","source_issue_refs":[{"repo":"transpara-ai/eventgraph","number":0}]}`,
			content: NativeTestRunRecordedContent{
				TestRunID:       "tr_bad",
				Command:         "go test",
				Outcome:         NativeEvidenceOutcomePassed,
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

func TestNativeEvidenceFixturePreservesIssueToAuditTraceability(t *testing.T) {
	issueRef := NativeEvidenceIssueRef{
		Repo:   "transpara-ai/eventgraph",
		Number: 63,
		URL:    "https://github.com/transpara-ai/eventgraph/issues/63",
		Title:  "Native TestRun GateResult and AuditReport persistence contract",
	}
	testRun := NativeTestRunRecordedContent{
		TestRunID:       "tr_eventgraph_63_native_evidence",
		Command:         "go test ./pkg/event -run TestNativeEvidence",
		Outcome:         NativeEvidenceOutcomePassed,
		SourceIssueRefs: []NativeEvidenceIssueRef{issueRef},
		PRRefs:          []string{"transpara-ai/eventgraph#999"},
		ValidationRefs:  []string{"validation:go-test-native-evidence"},
	}
	gate := NativeGateResultRecordedContent{
		GateResultID:    "gate_eventgraph_63_cfar",
		FactoryOrderID:  "fo_eventgraph_63",
		GateName:        "cross-family-adversarial-review",
		Outcome:         NativeEvidenceOutcomePassed,
		EvidenceRefs:    []string{testRun.TestRunID},
		SourceIssueRefs: []NativeEvidenceIssueRef{issueRef},
		PRRefs:          []string{"transpara-ai/eventgraph#999"},
	}
	audit := NativeAuditReportRecordedContent{
		AuditReportID:         "audit_eventgraph_63",
		TargetType:            "pull_request",
		TargetID:              "transpara-ai/eventgraph#999",
		Outcome:               NativeEvidenceOutcomePassed,
		TraceScoreBasisPoints: nativeEvidenceIntPtr(10000),
		SourceIssueRefs:       []NativeEvidenceIssueRef{issueRef},
		PRRefs:                []string{"transpara-ai/eventgraph#999"},
		ValidationRefs:        []string{"validation:go-test-native-evidence"},
		CFARRefs:              []string{"cfar:eventgraph-pr-999"},
		AuthorityBoundaryRefs: []string{"issue:eventgraph#63-pr-ready-boundary"},
		ResidualRiskRefs:      []string{"none"},
		EvidenceRefs:          []string{testRun.TestRunID, gate.GateResultID},
	}

	if audit.TraceScoreBasisPoints == nil || *audit.TraceScoreBasisPoints != 10000 || len(audit.MissingLinks) != 0 {
		t.Fatalf("audit fixture should be fully traced: %+v", audit)
	}
	if len(audit.SourceIssueRefs) != 1 || audit.SourceIssueRefs[0].Repo != issueRef.Repo || audit.SourceIssueRefs[0].Number != issueRef.Number {
		t.Fatalf("audit does not preserve source issue ref: %+v", audit.SourceIssueRefs)
	}
	if len(audit.PRRefs) != 1 || audit.PRRefs[0] != "transpara-ai/eventgraph#999" {
		t.Fatalf("audit does not preserve PR ref: %+v", audit.PRRefs)
	}
	if len(audit.ValidationRefs) == 0 || len(audit.CFARRefs) == 0 || len(audit.AuthorityBoundaryRefs) == 0 || len(audit.ResidualRiskRefs) == 0 {
		t.Fatalf("audit is missing closeout evidence refs: %+v", audit)
	}
	if len(audit.EvidenceRefs) != 2 || audit.EvidenceRefs[0] != testRun.TestRunID || audit.EvidenceRefs[1] != gate.GateResultID {
		t.Fatalf("audit does not link test run and gate result evidence: %+v", audit.EvidenceRefs)
	}
}

func nativeEvidenceIntPtr(value int) *int {
	return &value
}
