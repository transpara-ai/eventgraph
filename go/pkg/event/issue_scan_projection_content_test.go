package event

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

func TestIssueScanProjectionEventTypesRegistered(t *testing.T) {
	registry := DefaultRegistry()
	all := AllIssueScanProjectionEventTypes()
	if len(all) != 4 {
		t.Fatalf("AllIssueScanProjectionEventTypes returned %d types, want 4", len(all))
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

func TestIssueScanProjectionContentRoundTrip(t *testing.T) {
	cases := []struct {
		name      string
		eventType string
		content   EventContent
	}{
		{
			name:      "run",
			eventType: EventTypeIssueScanRunProjected.Value(),
			content: IssueScanRunProjectedContent{
				RunID:            "run_docs_172",
				FactoryOrderID:   "fo_docs_172",
				LifecycleVersion: "civilization_issue_to_human_ready_pr_v0.4",
				State:            IssueScanRunStateParked,
				TargetIssue:      IssueScanIssueRef{Repo: "transpara-ai/docs", Number: 172, URL: "https://github.com/transpara-ai/docs/issues/172", State: "open", Labels: []string{"cc:pr-ready"}},
				SelectedIssue:    IssueScanIssueRef{Repo: "transpara-ai/docs", Number: 172, URL: "https://github.com/transpara-ai/docs/issues/172", State: "open"},
				CandidateIssues:  []IssueScanIssueRef{{Repo: "transpara-ai/docs", Number: 172}, {Repo: "transpara-ai/site", Number: 115}},
				SourceRefs:       []string{"github:transpara-ai/docs#172"},
			},
		},
		{
			name:      "stage",
			eventType: EventTypeIssueScanStageProjected.Value(),
			content: IssueScanStageProjectedContent{
				RunID:             "run_docs_172",
				FactoryOrderID:    "fo_docs_172",
				StageID:           "research_issue_and_repo_context",
				StageNumber:       1,
				StageCount:        7,
				CanonicalTaskID:   "tsk_docs_172_research_issue_and_repo_context",
				TaskID:            "019c0000-0000-7000-8000-000000000172",
				CurrentState:      IssueScanStageStateBlocked,
				CompletionGate:    "required role outputs recorded",
				AuthorityBoundary: "no protected action or merge",
				AssignedAgentIDs:  []string{"agent_researcher"},
				TouchingAgentIDs:  []string{"agent_researcher", "agent_reviewer"},
				EvidenceRefs:      []string{"artifact:issue_scan_stage_role_contract"},
			},
		},
		{
			name:      "blocker",
			eventType: EventTypeIssueScanBlockerProjected.Value(),
			content: IssueScanBlockerProjectedContent{
				RunID:          "run_docs_172",
				FactoryOrderID: "fo_docs_172",
				StageID:        "research_issue_and_repo_context",
				BlockerType:    IssueScanBlockerDuplicateChain,
				Reason:         "canonical task chain duplicated",
				RequiredAction: "repair canonical lineage before runtime continues",
				EvidenceRefs:   []string{"work:tsk_docs_172_research_issue_and_repo_context"},
			},
		},
		{
			name:      "lineage",
			eventType: EventTypeIssueScanLineageProjected.Value(),
			content: IssueScanLineageProjectedContent{
				RunID:            "run_docs_172",
				FactoryOrderID:   "fo_docs_172",
				StageID:          "research_issue_and_repo_context",
				CanonicalTaskID:  "tsk_docs_172_research_issue_and_repo_context",
				PrimaryTaskID:    "019c0000-0000-7000-8000-000000000172",
				TaskIDs:          []string{"019c0000-0000-7000-8000-000000000172", "019c0000-0000-7000-8000-000000000173"},
				DuplicateTaskIDs: []string{"019c0000-0000-7000-8000-000000000173"},
				DuplicateOf:      "019c0000-0000-7000-8000-000000000172",
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

func TestIssueScanProjectionRejectsInvalidEnums(t *testing.T) {
	registry := DefaultRegistry()
	cases := []struct {
		name      string
		eventType types.EventType
		json      string
		content   EventContent
	}{
		{
			name:      "run state",
			eventType: EventTypeIssueScanRunProjected,
			json:      `{"run_id":"run_bad","lifecycle_version":"v","state":"banana","target_issue":{"repo":"transpara-ai/docs","number":172},"selected_issue":{"repo":"transpara-ai/docs","number":172}}`,
			content: IssueScanRunProjectedContent{
				RunID:            "run_bad",
				LifecycleVersion: "v",
				State:            IssueScanRunState("banana"),
			},
		},
		{
			name:      "stage state",
			eventType: EventTypeIssueScanStageProjected,
			json:      `{"run_id":"run_bad","stage_id":"research","stage_number":1,"canonical_task_id":"tsk_bad","current_state":"banana","completion_gate":"gate","authority_boundary":"none"}`,
			content: IssueScanStageProjectedContent{
				RunID:        "run_bad",
				StageID:      "research",
				StageNumber:  1,
				CurrentState: IssueScanStageState("banana"),
			},
		},
		{
			name:      "blocker type",
			eventType: EventTypeIssueScanBlockerProjected,
			json:      `{"run_id":"run_bad","blocker_type":"banana","required_action":"stop"}`,
			content: IssueScanBlockerProjectedContent{
				RunID:          "run_bad",
				BlockerType:    IssueScanBlockerType("banana"),
				RequiredAction: "stop",
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

func TestIssueScanProjectionDocs172Site115FixtureHasTypedKanbanFields(t *testing.T) {
	runDocs := IssueScanRunProjectedContent{
		RunID:            "run_docs_172",
		FactoryOrderID:   "fo_docs_172",
		LifecycleVersion: "civilization_issue_to_human_ready_pr_v0.4",
		State:            IssueScanRunStateParked,
		TargetIssue:      IssueScanIssueRef{Repo: "transpara-ai/docs", Number: 172, URL: "https://github.com/transpara-ai/docs/issues/172", State: "open"},
		SelectedIssue:    IssueScanIssueRef{Repo: "transpara-ai/docs", Number: 172, URL: "https://github.com/transpara-ai/docs/issues/172", State: "open"},
		CandidateIssues: []IssueScanIssueRef{
			{Repo: "transpara-ai/docs", Number: 172, Labels: []string{"cc:pr-ready"}},
			{Repo: "transpara-ai/site", Number: 115, Labels: []string{"cc:pr-ready", "cc:protected-action"}},
		},
		SourceRefs: []string{"github:transpara-ai/docs#172", "github:transpara-ai/site#115"},
	}
	stageDocs := IssueScanStageProjectedContent{
		RunID:             runDocs.RunID,
		FactoryOrderID:    runDocs.FactoryOrderID,
		StageID:           "run_adversarial_review",
		StageNumber:       5,
		StageCount:        7,
		CanonicalTaskID:   "tsk_docs_172_run_adversarial_review",
		TaskID:            "019c0000-0000-7000-8000-000000001172",
		CurrentState:      IssueScanStageStateBlocked,
		CompletionGate:    "exact-head adversarial review returns zero blockers",
		AuthorityBoundary: "review only; no merge or deploy",
		AssignedAgentIDs:  []string{"agent_reviewer"},
		TouchingAgentIDs:  []string{"agent_reviewer", "agent_blocker_repair"},
		EvidenceRefs:      []string{"code.review.submitted:docs-172"},
	}
	duplicateLineage := IssueScanLineageProjectedContent{
		RunID:            runDocs.RunID,
		FactoryOrderID:   runDocs.FactoryOrderID,
		StageID:          stageDocs.StageID,
		CanonicalTaskID:  stageDocs.CanonicalTaskID,
		PrimaryTaskID:    stageDocs.TaskID,
		TaskIDs:          []string{stageDocs.TaskID, "019c0000-0000-7000-8000-000000001173"},
		DuplicateTaskIDs: []string{"019c0000-0000-7000-8000-000000001173"},
		DuplicateOf:      stageDocs.TaskID,
		SourceRefs:       []string{"work:duplicate-stage-chain"},
	}
	blockers := []IssueScanBlockerProjectedContent{
		{RunID: runDocs.RunID, FactoryOrderID: runDocs.FactoryOrderID, StageID: stageDocs.StageID, BlockerType: IssueScanBlockerDuplicateChain, RequiredAction: "collapse duplicate canonical stage chain"},
		{RunID: "run_site_115", FactoryOrderID: "fo_site_115", StageID: "surface_ready_for_human_result_pr", BlockerType: IssueScanBlockerProtectedAction, RequiredAction: "human must authorize protected repo action"},
		{RunID: "run_site_115", FactoryOrderID: "fo_site_115", StageID: "research_issue_and_repo_context", BlockerType: IssueScanBlockerStaleTarget, RequiredAction: "confirm target issue is still live"},
		{RunID: "run_docs_172_scope", FactoryOrderID: "fo_docs_172_scope", StageID: "select_and_design_approach", BlockerType: IssueScanBlockerNeedsHumanScope, RequiredAction: "human must clarify issue scope before runtime continues"},
	}

	if runDocs.SelectedIssue.Repo != "transpara-ai/docs" || runDocs.SelectedIssue.Number != 172 {
		t.Fatalf("selected issue is not typed docs#172: %+v", runDocs.SelectedIssue)
	}
	if len(runDocs.CandidateIssues) != 2 || runDocs.CandidateIssues[1].Repo != "transpara-ai/site" || runDocs.CandidateIssues[1].Number != 115 {
		t.Fatalf("candidate issues do not preserve typed site#115: %+v", runDocs.CandidateIssues)
	}
	if stageDocs.CanonicalTaskID == "" || stageDocs.StageNumber == 0 || stageDocs.CompletionGate == "" || stageDocs.AuthorityBoundary == "" {
		t.Fatalf("stage is missing Kanban-critical typed fields: %+v", stageDocs)
	}
	if len(stageDocs.AssignedAgentIDs) == 0 || len(stageDocs.TouchingAgentIDs) == 0 {
		t.Fatalf("stage lacks typed agent assignment/touch fields: %+v", stageDocs)
	}
	if duplicateLineage.DuplicateOf != stageDocs.TaskID || len(duplicateLineage.DuplicateTaskIDs) != 1 {
		t.Fatalf("duplicate lineage is not structural: %+v", duplicateLineage)
	}
	wantBlockers := map[IssueScanBlockerType]bool{
		IssueScanBlockerDuplicateChain:  true,
		IssueScanBlockerNeedsHumanScope: true,
		IssueScanBlockerProtectedAction: true,
		IssueScanBlockerStaleTarget:     true,
	}
	for _, blocker := range blockers {
		delete(wantBlockers, blocker.BlockerType)
		if blocker.RequiredAction == "" {
			t.Fatalf("blocker lacks typed required action: %+v", blocker)
		}
	}
	if len(wantBlockers) != 0 {
		t.Fatalf("fixture missing blocker types: %+v", wantBlockers)
	}
}
