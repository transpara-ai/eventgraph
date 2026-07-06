# Issue-Scan Projection Events

Civilization issue-hunting dashboards must consume typed issue-scan projection
events instead of deriving state from prose comments, logs, or agent narration.
These events are read-model events: they make the current run, stage, blocker,
and lineage state durable without authorizing Hive runtime work.

## Event Types

| Event type | Purpose |
|---|---|
| `issuescan.run.projected` | Current lifecycle state for one issue-scan run and its selected/candidate GitHub issues. |
| `issuescan.stage.projected` | Current state for one canonical lifecycle stage, including gate, authority boundary, and agent assignment/touch data. |
| `issuescan.blocker.projected` | Structured blocker reason and required human or system action. |
| `issuescan.lineage.projected` | Canonical task lineage, primary task, duplicates, and superseded task IDs. |
| `issuescan.source.marker.projected` | Work/EventGraph source-marker projection state for the derived GitHub source issue marker. |

The `issuescan` prefix is intentionally unseparated. It is the valid EventGraph
event-type token for the issue-scan domain.

## Dashboard Contract

Kanban and monitor consumers should render from these typed fields:

- Run cards use `run_id`, `factory_order_id`, `lifecycle_version`, `state`,
  `target_issue`, `selected_issue`, `candidate_issues`, `source_refs`, and
  `evidence_refs`.
- Stage cards use `run_id`, `stage_id`, `stage_number`, `stage_count`,
  `canonical_task_id`, `task_id`, `current_state`, `completion_gate`,
  `authority_boundary`, `assigned_agent_ids`, `touching_agent_ids`,
  `evidence_refs`, and `source_refs`.
- Blocked columns use `blocker_type`, `reason`, and `required_action` from
  `issuescan.blocker.projected`.
- Duplicate-chain views use `canonical_task_id`, `primary_task_id`, `task_ids`,
  `duplicate_task_ids`, `duplicate_of`, and `superseded_task_ids` from
  `issuescan.lineage.projected`.
- Source-marker views use `transition`, `work_ref`, `actor_id`, `actor_role`,
  `schema_version`, `projection_kind`, `run_id`, `target`, `stage_id`,
  `stage_number`, `gate`, `occurred_at`, `idempotency_key`,
  `authority_boundary`, `authority_exclusions`, `evidence_refs`,
  `source_refs`, `github_marker`, `canonical_source`, `projection_only`,
  `superseded_by`, and `stale_target` from
  `issuescan.source.marker.projected`.

Consumers must treat `blocked`, `parked`, `human_action`, `ready_for_human`,
`superseded`, and `projection_only` as terminal or non-running dashboard states
unless a later projection event changes the state. A monitor or dashboard must
not wake Hive runtime work from these projection events.

## Source-Issue Marker Boundary

`issuescan.source.marker.projected` is the EventGraph side of the
`work.issue_scan.source_marker_ref` contract. It records the projection state
for the compact marker that Hive may later write to a source GitHub issue and
that Site may render as operator evidence.

The valid source-marker transitions are:

| Transition | Meaning |
|---|---|
| `acquired` | A Work issue-scan stage acquired a source issue and emitted canonical refs. |
| `parked_human_action` | The marker is parked because human scope, protected-action, stale-target, or gate evidence is required. |
| `ready_for_human` | The Work/EventGraph projection is ready for human review without granting merge or mutation authority. |
| `completed` | The Work-local source-marker stage completed. |
| `abandoned` | The source-marker path was intentionally abandoned without making GitHub canonical. |
| `superseded` | A newer Work task/projection replaces this marker state. |

The `work_ref` field mirrors the Work-owned
`work.issue_scan.source_marker_ref` packet. EventGraph validates that the
projection run, target, and stage match the Work ref, that the Work ref is
`projection_only`, and that the Work ref declares `canonical_source: "work"`.

The `github_marker` field is an output reference only. When present, `system`
must be `github`, its repository and issue number must match the source target,
and it must be marked with `derived_output: true` and `projection_sink: true`.
Consumers must not parse GitHub marker comments or labels back into canonical
Work or EventGraph truth.

Source-marker projections carry `authority_boundary` and
`authority_exclusions` explicitly. Both the projection and embedded Work ref
must include the required exclusion tokens for GitHub projection-only status,
comments/labels not being lifecycle truth, no live GitHub mutation authority,
no production EventGraph write, no Hive write/action/authority API, no deploy,
no Test 001 GREEN, no merge authority, no issue closure, no autonomy increase,
and no value allocation. This conformance contract does not authorize
production EventGraph writes, Hive action APIs, live GitHub mutation, deploy,
value allocation, autonomy increase, Test 001 GREEN, merge, issue closure, or
wiki work.

Source-marker transitions must remain coherent with the embedded Work ref:
`acquired` requires `created`, `running`, `repair_running`, or
`verification_running`, `parked_human_action` requires a blocked Work ref with
`latest_blocker` and a blocked lifecycle state,
`ready_for_human` requires a ready, unblocked Work ref with `ready`, `verified`,
or `repaired` lifecycle state, `completed` requires a certified, unblocked Work
ref with no missing gates/facts and the matching latest gate, `abandoned`
requires `rejected`, and `superseded` requires matching top-level and Work
`superseded_by` values plus `work_ref.lifecycle_state: "superseded"`.
The Work lifecycle state is closed over the current Work status vocabulary:
`created`, `ready`, `running`, `blocked`, `policy_blocked`, `failed`,
`repair_required`, `repair_running`, `repaired`, `verification_running`,
`verified`, `certified`, `rejected`, and `superseded`.
Parked states may cite `needs_human_scope`, `protected_action`,
`stale_target`, `duplicate_chain`, or `missing_gate_evidence`; when
`stale_target` is true, the latest blocker reason must be `stale_target`, and a
`stale_target` blocker must set `stale_target: true`.

Unlike sibling run/stage/blocker projection events, the source-marker
`evidence_refs` field is the structured `IssueScanMarkerEvidenceRefs` object
that mirrors Work marker evidence buckets rather than a flat string list.

## Language Coverage

The Go package registers first-class content structs for these event types and
round-trips them through `UnmarshalContent`. TypeScript, Python, and Rust keep
event content as canonical JSON maps, so they preserve the same JSON field names
through hash-chain and store conformance without requiring per-event content
structs.

## Restart Scenario Fixture

The `docs#172` / `site#115` restart check should emit projection events that
show:

- one selected target issue per run,
- no duplicate canonical stage chains,
- no running workers for parked, blocked, stale, or human-scope runs,
- explicit `duplicate_chain`, `needs_human_scope`, `protected_action`, or
  `stale_target` blockers when those conditions are present,
- a clear `required_action` string for the human or upstream repair step.

Separate source-marker transition fixtures should cover acquired,
parked/human-action, ready-for-human, completed, abandoned, and superseded
marker states.
