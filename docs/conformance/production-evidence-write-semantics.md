# Production Evidence Write Semantics

Event 17 production EventGraph wiring needs a deterministic path from runtime,
issue, authority, TestRun, GateResult, and AuditReport evidence into append-only
EventGraph truth. This contract defines the local write-candidate semantics
only. It does not authorize or execute a live production EventGraph write.

## Local Candidate Contract

`PlanProductionEvidenceWrite` accepts a `ProductionEvidenceWriteCandidate` and
returns an ordered append plan when all authority and evidence checks pass.

The required local candidate fields are:

- `candidate_id`
- `repo`
- `actor_ref`
- `action_class`
- `schema_version`
- `source_state_ref`
- `current_state_ref`
- `source_issue_refs`
- `runtime_evidence_refs`
- `issue_evidence_refs`
- one `authority.decision.recorded`
- one `authority.storegovernance.recorded`
- at least one `evidence.testrun.recorded`
- at least one `evidence.gateresult.recorded`
- at least one `evidence.auditreport.recorded`

The only accepted action class is:

```text
eventgraph.production_evidence.write
```

The only accepted schema version is:

```text
production_evidence_write_v0
```

The only accepted store-governance store name is:

```text
production_eventgraph_evidence
```

## Authority Requirements

The authority decision must:

- use `subject_type: eventgraph_write_path`;
- use `subject_ref` equal to the candidate repo;
- use `outcome: autonomous`;
- use `actor_ref` equal to the candidate actor;
- cite the same source issue refs as the candidate;
- pass the normal `authority.decision.recorded` validator.

The store-governance record must:

- use `write_status: write_path_authorized`;
- use the candidate schema version;
- cite the same source issue refs as the candidate;
- include the authority decision ID in `authority_refs`;
- pass the normal `authority.storegovernance.recorded` validator.

`approval_required`, `forbidden`, missing authority refs, wrong actors, wrong
repos, wrong action classes, invalid schemas, and stale source-state refs fail
closed before append planning.

`runtime_evidence_refs` and `issue_evidence_refs` are metadata refs carried by
the local append plan. They make the runtime and GitHub issue evidence sources
explicit without introducing a live runtime writer or a production issue-sync
adapter.

`source_state_ref` and `current_state_ref` are caller-supplied local snapshot
refs. Matching values prove only local candidate self-consistency. They are not
live-head freshness checks and cannot authorize production writes without a
separate exact-head authority packet or adapter-owned freshness proof.

## Evidence Requirements

Native evidence records must pass their existing validators and cite the same
source issue refs as the candidate. The comparison is exact by source issue ref
count; duplicate refs cannot stand in for a missing issue.

Gate results must include at least one TestRun ID in `evidence_refs`.

Audit reports must include validation refs, CFAR refs, authority-boundary refs,
at least one TestRun ID, and at least one GateResult ID in `evidence_refs`. A
passed audit report must have no `missing_links` and must use
`trace_score_basis_points: 10000`.

The write plan records evidence outcomes. It does not certify that every
TestRun, GateResult, or AuditReport passed. Consumers that require green
closeout must enforce that gate separately before approving merge, deployment,
runtime execution, or any other protected action.

Duplicate evidence IDs inside one plan fail closed. Applying a plan to a local
or fixture entry set is idempotent for exact duplicate entries and fails closed
for conflicting entries with the same event type and evidence ID.

## Boundaries

This contract returns an in-memory ordered append plan and a local fixture apply
helper. It does not:

- open a production store connection;
- execute a live EventGraph write;
- run RuntimeBroker;
- start Hive;
- deploy;
- access secrets;
- change protected settings;
- mutate branch protection;
- close Test 001;
- close `docs#172`;
- increase autonomy;
- allocate value;
- close residual risk.

Production enablement requires a separate exact-head authority packet and
cannot be inferred from this local semantics contract.
