# Native Evidence Events

Civilization closeout monitors must consume typed evidence events instead of
parsing PR prose, logs, or agent narration. These events make `TestRun`,
`GateResult`, and `AuditReport` records visible as durable EventGraph content
without authorizing production writes, runtime execution, or protected actions.

## Event Types

| Event type | Purpose |
|---|---|
| `evidence.testrun.recorded` | A validation command or test run result, with source issue and PR refs. |
| `evidence.gateresult.recorded` | A gate verdict such as CFAR or trace completeness, with evidence and waiver refs. |
| `evidence.auditreport.recorded` | Final closeout traceability for source issues, PRs, validation, CFAR, authority boundary, and residual-risk refs. |

## Outcome Vocabulary

Native evidence events use the shared `outcome` vocabulary:

- `passed`
- `failed`
- `blocked`
- `skipped`
- `errored`
- `waived`
- `partial`

Consumers must treat missing or unknown outcomes as invalid evidence, not as a
successful or complete result.

## Closeout Contract

An issue-to-PR closeout packet should project at least:

- `source_issue_refs` for the GitHub issue source-of-intent records,
- `pr_refs` for the pull requests or merge records being closed out,
- `validation_refs` for deterministic local or CI validation,
- `cfar_refs` for exact-head cross-family adversarial review artifacts,
- `authority_boundary_refs` for protected-action or scope boundaries,
- `residual_risk_refs` for unresolved, waived, closed, or explicitly absent
  residual risks.

`evidence.auditreport.recorded` must use `trace_score_basis_points` as an
integer from `0` through `10000`. This avoids cross-language floating-point
canonicalization drift in hash-chain content. A `trace_score_basis_points` value
of `10000` means the packet has no missing links. Lower scores or non-empty
`missing_links` must be rendered as incomplete or partial by dashboards and
monitors.

Identity fields such as `test_run_id`, `gate_result_id`, `factory_order_id`,
`audit_report_id`, `target_type`, and `target_id` are required. When present,
`source_issue_refs` entries must include a non-empty `repo` and a positive
`number`.

## Boundaries

These events are content contracts only. They do not create a production
projection store, schema migration, EventGraph write path, Hive action API, Work
mutation path, Site replacement, deploy path, merge authority, autonomy
increase, value allocation, residual-risk closure, or Test 001 status change.

Consumers may render these events in monitoring surfaces. They must not wake
Hive runtime work or treat the events as authority for protected actions.
