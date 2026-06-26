# Authority Evidence Governance Events

Civilization authority monitors must consume typed authority evidence events
instead of inferring authority from markdown prose, issue comments, PR text, or
agent narration. These events make authority decisions, protected-action
boundaries, residuals, and projection-store governance visible without
authorizing production writes, runtime execution, or protected actions.

## Event Types

| Event type | Purpose |
|---|---|
| `authority.decision.recorded` | A durable authority decision outcome for a GitHub issue, PR, protected action, or governance subject. |
| `authority.boundary.recorded` | A typed protected-action or human-scope boundary, including required authority refs and protected action names. |
| `authority.residual.recorded` | A carried, waived, closed, open, or not-applicable residual risk or process residual. |
| `authority.storegovernance.recorded` | Projection-store schema, migration, validation, and write-path posture for authority/evidence state. |

## Authority Outcome Vocabulary

`authority.decision.recorded` uses lowercase JSON values corresponding to the
Dark Factory authority outcomes:

- `autonomous`
- `notify`
- `approval_required`
- `forbidden`

Consumers must treat a missing or unknown outcome as invalid evidence, not as
approval. `autonomous` and `approval_required` decisions must include
`authority_refs`; otherwise they are invalid evidence.

## Boundary Vocabulary

`authority.boundary.recorded` uses these `boundary_type` values:

- `protected_action`
- `human_scope`
- `runtime_execution`
- `eventgraph_write`
- `merge`
- `deployment`
- `autonomy_increase`
- `value_allocation`
- `residual_risk_closure`
- `protected_settings`

`state` is one of:

- `blocked`
- `approval_required`
- `authorized`
- `not_applicable`

If `boundary_type` is `protected_action`, `protected_actions` must contain at
least one protected-action name registered by the EventGraph implementation.
Aliases such as `deploy.production` are invalid. This PR does not expand the
existing protected-action vocabulary or change the existing
`authority.requested` helper path; full vocabulary alignment with
`docs/dark-factory/authority-vocabulary.md` remains separate governance work.

If `state` is `authorized`, `required_authority_refs` must contain at least one
durable authority reference. This prevents dashboards and Hive intake from
treating a bare status word as authority.

## Residual Vocabulary

`authority.residual.recorded` uses these `status` values:

- `open`
- `carried`
- `waived`
- `closed`
- `not_applicable`

`severity` is one of:

- `low`
- `medium`
- `high`
- `critical`

Open and carried residuals must include `required_action`. Monitors should park
work instead of continuing token-consuming runtime loops when a residual is
open or carried and the required action is missing.

## Store Governance

`authority.storegovernance.recorded` describes projection-store posture with
`write_status`:

- `schema_only`
- `projection_only`
- `migration_required`
- `write_path_blocked`
- `write_path_authorized`

`write_path_authorized` requires both `authority_refs` and
`required_validation_refs`. Any store consumer that sees missing authority or
validation for an authorized write path must fail closed and render the state as
blocked.

The current `eventgraph#62` contract is schema/projection governance only. It
does not create a production projection store, execute a migration, enable a
write path, or authorize runtime truth mutation.

## Dashboard And Hive Contract

Dashboards and Hive issue-intake flows may render or reason over:

- `source_issue_refs`
- `authority_refs`
- `required_authority_refs`
- `protected_actions`
- `evidence_refs`
- `non_claim_refs`
- `migration_refs`
- `required_validation_refs`

They must not infer authority from prose fields. Unknown enum values, missing
required refs, empty protected-action lists on protected boundaries, or
authorized write paths without authority and validation refs are invalid
evidence and must park the work.

## Boundaries

These events are content contracts only. They do not create a production
projection store, schema migration, EventGraph write path, Hive action API, Work
mutation path, Site replacement, deployment path, merge authority, autonomy
increase, value allocation, residual-risk closure, or Test 001 status change.

Consumers may render these events in monitoring surfaces. They must not wake
Hive runtime work or treat the events as authority for protected actions.
