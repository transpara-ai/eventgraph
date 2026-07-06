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
	if len(all) != 5 {
		t.Fatalf("AllIssueScanProjectionEventTypes returned %d types, want 5", len(all))
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
		{
			name:      "source marker",
			eventType: EventTypeIssueScanSourceMarkerProjected.Value(),
			content:   sourceMarkerProjectionFixture(IssueScanSourceMarkerAcquired),
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

func TestIssueScanSourceMarkerProjectionTransitions(t *testing.T) {
	transitions := []IssueScanSourceMarkerTransition{
		IssueScanSourceMarkerAcquired,
		IssueScanSourceMarkerParkedHumanAction,
		IssueScanSourceMarkerReadyForHuman,
		IssueScanSourceMarkerCompleted,
		IssueScanSourceMarkerAbandoned,
		IssueScanSourceMarkerSuperseded,
	}

	for _, transition := range transitions {
		t.Run(string(transition), func(t *testing.T) {
			content := sourceMarkerProjectionFixture(transition)
			if transition == IssueScanSourceMarkerParkedHumanAction {
				content.StaleTarget = true
				content.WorkRef.LatestBlocker = &IssueScanMarkerBlockerRef{
					Reason:       IssueScanBlockerStaleTarget,
					Detail:       "source issue was closed after acquisition",
					EvidenceRefs: []string{"github:transpara-ai/docs#256"},
				}
			}
			if transition == IssueScanSourceMarkerSuperseded {
				content.SupersededBy = "task:successor"
				content.WorkRef.SupersededBy = "task:successor"
			}
			data, err := json.Marshal(content)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			got, err := UnmarshalContent(EventTypeIssueScanSourceMarkerProjected.Value(), data)
			if err != nil {
				t.Fatalf("UnmarshalContent: %v", err)
			}
			if !reflect.DeepEqual(got, content) {
				t.Fatalf("round trip mismatch\n got: %#v\nwant: %#v", got, content)
			}
			if err := DefaultRegistry().Validate(EventTypeIssueScanSourceMarkerProjected, content); err != nil {
				t.Fatalf("Validate: %v", err)
			}
		})
	}
}

func TestIssueScanSourceMarkerProjectionRejectsCanonicalGitHubMarkers(t *testing.T) {
	content := sourceMarkerProjectionFixture(IssueScanSourceMarkerAcquired)
	content.GitHubMarker.DerivedOutput = false
	if err := DefaultRegistry().Validate(EventTypeIssueScanSourceMarkerProjected, content); err == nil {
		t.Fatal("registry accepted a GitHub marker that was not a derived output")
	}

	content = sourceMarkerProjectionFixture(IssueScanSourceMarkerAcquired)
	content.GitHubMarker.ProjectionSink = false
	if err := DefaultRegistry().Validate(EventTypeIssueScanSourceMarkerProjected, content); err == nil {
		t.Fatal("registry accepted a GitHub marker that was not a projection sink")
	}

	content = sourceMarkerProjectionFixture(IssueScanSourceMarkerAcquired)
	content.GitHubMarker.System = "GitHub"
	if err := DefaultRegistry().Validate(EventTypeIssueScanSourceMarkerProjected, content); err == nil {
		t.Fatal("registry accepted a GitHub marker with a non-canonical system")
	}

	content = sourceMarkerProjectionFixture(IssueScanSourceMarkerAcquired)
	content.GitHubMarker = nil
	if err := DefaultRegistry().Validate(EventTypeIssueScanSourceMarkerProjected, content); err != nil {
		t.Fatalf("registry rejected absent GitHub marker: %v", err)
	}
}

func TestIssueScanSourceMarkerProjectionRejectsMismatchedWorkRef(t *testing.T) {
	cases := []struct {
		name string
		edit func(*IssueScanSourceMarkerProjectedContent)
	}{
		{
			name: "invalid transition",
			edit: func(c *IssueScanSourceMarkerProjectedContent) {
				c.Transition = IssueScanSourceMarkerTransition("parsed_from_github_comment")
			},
		},
		{
			name: "target mismatch",
			edit: func(c *IssueScanSourceMarkerProjectedContent) {
				c.WorkRef.Target.IssueNumber = 999
			},
		},
		{
			name: "stage mismatch",
			edit: func(c *IssueScanSourceMarkerProjectedContent) {
				c.WorkRef.Stage = "github_comment_body"
			},
		},
		{
			name: "stage number mismatch",
			edit: func(c *IssueScanSourceMarkerProjectedContent) {
				c.WorkRef.StageNumber = 2
			},
		},
		{
			name: "gate mismatch",
			edit: func(c *IssueScanSourceMarkerProjectedContent) {
				c.WorkRef.Gate = "different_gate"
			},
		},
		{
			name: "top-level projection kind mismatch",
			edit: func(c *IssueScanSourceMarkerProjectedContent) {
				c.ProjectionKind = "work.issue_scan.source_marker_ref"
			},
		},
		{
			name: "top-level projection_only false",
			edit: func(c *IssueScanSourceMarkerProjectedContent) {
				c.ProjectionOnly = false
			},
		},
		{
			name: "work projection_only false",
			edit: func(c *IssueScanSourceMarkerProjectedContent) {
				c.WorkRef.ProjectionOnly = false
			},
		},
		{
			name: "work canonical source mismatch",
			edit: func(c *IssueScanSourceMarkerProjectedContent) {
				c.WorkRef.CanonicalSource = "github"
			},
		},
		{
			name: "empty authority exclusions",
			edit: func(c *IssueScanSourceMarkerProjectedContent) {
				c.AuthorityExclusions = nil
			},
		},
		{
			name: "empty work authority exclusions",
			edit: func(c *IssueScanSourceMarkerProjectedContent) {
				c.WorkRef.AuthorityExclusions = nil
			},
		},
		{
			name: "empty actor id",
			edit: func(c *IssueScanSourceMarkerProjectedContent) {
				c.ActorID = ""
			},
		},
		{
			name: "invalid nested blocker reason",
			edit: func(c *IssueScanSourceMarkerProjectedContent) {
				c.WorkRef.LatestBlocker = &IssueScanMarkerBlockerRef{Reason: IssueScanBlockerType("github_comment")}
			},
		},
		{
			name: "superseded without successor",
			edit: func(c *IssueScanSourceMarkerProjectedContent) {
				c.Transition = IssueScanSourceMarkerSuperseded
			},
		},
		{
			name: "stale target advanced",
			edit: func(c *IssueScanSourceMarkerProjectedContent) {
				c.Transition = IssueScanSourceMarkerCompleted
				c.StaleTarget = true
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			content := sourceMarkerProjectionFixture(IssueScanSourceMarkerAcquired)
			tc.edit(&content)
			if err := DefaultRegistry().Validate(EventTypeIssueScanSourceMarkerProjected, content); err == nil {
				t.Fatalf("registry accepted invalid source marker projection: %+v", content)
			}
		})
	}
}

func TestIssueScanSourceMarkerProjectionRejectsInvalidJSON(t *testing.T) {
	content := sourceMarkerProjectionFixture(IssueScanSourceMarkerAcquired)
	content.WorkRef.RunID = "different-run"
	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if _, err := UnmarshalContent(EventTypeIssueScanSourceMarkerProjected.Value(), data); err == nil {
		t.Fatal("UnmarshalContent accepted mismatched work_ref run_id")
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
			name:      "missing run state",
			eventType: EventTypeIssueScanRunProjected,
			json:      `{"run_id":"run_bad","lifecycle_version":"v","target_issue":{"repo":"transpara-ai/docs","number":172},"selected_issue":{"repo":"transpara-ai/docs","number":172}}`,
			content: IssueScanRunProjectedContent{
				RunID:            "run_bad",
				LifecycleVersion: "v",
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
			name:      "missing stage state",
			eventType: EventTypeIssueScanStageProjected,
			json:      `{"run_id":"run_bad","stage_id":"research","stage_number":1,"canonical_task_id":"tsk_bad","completion_gate":"gate","authority_boundary":"none"}`,
			content: IssueScanStageProjectedContent{
				RunID:       "run_bad",
				StageID:     "research",
				StageNumber: 1,
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
		{
			name:      "missing blocker type",
			eventType: EventTypeIssueScanBlockerProjected,
			json:      `{"run_id":"run_bad","required_action":"stop"}`,
			content: IssueScanBlockerProjectedContent{
				RunID:          "run_bad",
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

func sourceMarkerProjectionFixture(transition IssueScanSourceMarkerTransition) IssueScanSourceMarkerProjectedContent {
	workRef := IssueScanMarkerWorkRef{
		SchemaVersion:          "1",
		ProjectionKind:         "work.issue_scan.source_marker_ref",
		CanonicalSource:        "work",
		ProjectionOnly:         true,
		RunID:                  "2026-07-06-docs-256",
		Target:                 IssueScanMarkerTargetRef{Repository: "transpara-ai/docs", IssueNumber: 256},
		Stage:                  "research_issue_and_repo_context",
		StageNumber:            1,
		Gate:                   "research_packet_posted",
		TaskID:                 "019f5000-0000-7000-8000-000000000256",
		CanonicalTaskID:        "tsk_issue_scan_2026_07_06_docs_256_transpara_ai_docs_256_research_issue_and_repo_context_cc56864a0bb3",
		FactoryOrderID:         "fo_issue_scan_2026_07_06_docs_256",
		RequirementIDs:         []string{"req_issue_scan_2026_07_06_docs_256_transpara_ai_docs_256_research_issue_and_repo_context_cc56864a0bb3"},
		AcceptanceCriterionIDs: []string{"ac_issue_scan_2026_07_06_docs_256_transpara_ai_docs_256_research_issue_and_repo_context_cc56864a0bb3"},
		LifecycleState:         "created",
		Ready:                  false,
		Blocked:                false,
		MissingGates:           []string{"definition_of_done", "acceptance_criteria", "test_plan"},
		VerificationRefs:       IssueScanMarkerEvidenceRefs{TestCaseIDs: []string{"eventgraph:issuescan-source-marker-projection"}},
		FailureRepairRefs:      IssueScanMarkerEvidenceRefs{},
		SourceIssueRefs:        []string{"github:transpara-ai/docs#256"},
		AuthorityExclusions: []string{
			"github_issue_markers_are_projection_only",
			"github_comments_are_not_work_lifecycle_truth",
			"github_labels_are_not_work_lifecycle_truth",
			"no_live_github_mutation_authority",
			"no_eventgraph_production_write",
			"no_hive_write_action_or_authority_api",
			"no_deployment",
			"no_test_001_green",
			"no_merge_authority",
			"no_issue_closure",
			"no_autonomy_increase",
			"no_value_allocation",
		},
	}
	return IssueScanSourceMarkerProjectedContent{
		SchemaVersion:       "1",
		ProjectionKind:      "eventgraph.issue_scan.source_marker_projection",
		Transition:          transition,
		RunID:               workRef.RunID,
		Target:              IssueScanIssueRef{Repo: workRef.Target.Repository, Number: workRef.Target.IssueNumber, URL: "https://github.com/transpara-ai/docs/issues/256", State: "open"},
		StageID:             workRef.Stage,
		StageNumber:         workRef.StageNumber,
		Gate:                workRef.Gate,
		WorkRef:             workRef,
		ActorID:             "agent:eventgraph-projection",
		ActorRole:           "projection_recorder",
		OccurredAt:          "2026-07-06T14:00:00Z",
		IdempotencyKey:      "issuescan-source-marker:2026-07-06-docs-256:research_issue_and_repo_context:" + string(transition),
		AuthorityBoundary:   "projection only; no production EventGraph write or GitHub mutation",
		AuthorityExclusions: append([]string(nil), workRef.AuthorityExclusions...),
		EvidenceRefs:        IssueScanMarkerEvidenceRefs{TestCaseIDs: []string{"eventgraph:issuescan-source-marker-projection"}},
		SourceRefs:          []string{"work:fo_issue_scan_2026_07_06_docs_256", "github:transpara-ai/docs#256"},
		GitHubMarker: &IssueScanSourceMarkerOutputRef{
			System:         "github",
			Repository:     "transpara-ai/docs",
			IssueNumber:    256,
			CommentID:      "planned-marker-comment",
			LabelNames:     []string{"cc:civilization-presence"},
			DerivedOutput:  true,
			ProjectionSink: true,
		},
		CanonicalSource: "work_eventgraph_projection",
		ProjectionOnly:  true,
	}
}
