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

Consumers must treat `blocked`, `parked`, `human_action`, `ready_for_human`,
`superseded`, and `projection_only` as terminal or non-running dashboard states
unless a later projection event changes the state. A monitor or dashboard must
not wake Hive runtime work from these projection events.

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
