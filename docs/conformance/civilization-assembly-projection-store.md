# Civilization Assembly Projection Store Record

The Civilization Assembly projection-store record makes a derived
Civilization Assembly projection durable as typed EventGraph state. Consumers
must treat it as projection truth for dashboards, issue-intake monitors, and
dry-run replay evidence, not as permission to execute protected actions.

## Record Type

`CivilizationAssemblyProjectionStoreRecord` stores one embedded
`CivilizationAssemblyProjection` plus the metadata required to replay, verify,
and reject unsafe projection-store materialization:

- `projection_id`
- `projection_schema_version`
- `projection_subject`
- `generated_at`
- `source_eventgraph_head_or_state_version`
- `source_event_ids_or_query_window`
- `derivation_status`
- `authority_decision_ref`
- `projection`
- `provenance_refs`
- `validation_refs`
- `boundary_flags`

## Required Authority

The store helper requires `authority_decision_ref` to point at an existing
active `AuthorityDecision` record whose `decision` is `Autonomous`, whose
`scope` includes `eventgraph.civilization_assembly.projection_store.record`,
and whose `expires_at` has not passed at the projection-store record time.
Missing, unscoped, pending, forbidden, or expired authority is invalid input and
must fail closed without appending a projection-store record.

The current governed scope comes from the docs#207 Event 17 / Gate AA authority
packet and eventgraph#59. That scope authorizes only local/testable
projection-store behavior in the EventGraph repository.

## Validation

A valid record must:

- use the current Civilization Assembly projection schema version,
- use the `civilization_assembly` projection subject,
- embed matching top-level projection metadata,
- have derivation status `complete` or `partial`,
- include non-empty source ids, provenance refs, and validation refs,
- include the authority decision ref in provenance,
- require the authority decision scope
  `eventgraph.civilization_assembly.projection_store.record`,
- include local-only and no-production boundary flags.

Failed or unavailable projections are not durable projection truth. The append
path rejects them and leaves the store unchanged.

## Replay Behavior

Projection-store records are excluded from subsequent Civilization Assembly
source snapshots. This prevents recursive self-inclusion, keeps the source
state hash deterministic across idempotent replays, and prevents one projection
record from creating another chain of derived truth.

## Boundaries

This contract does not create or authorize:

- production EventGraph writes,
- EventGraph schema migration,
- RuntimeBroker or Hive runtime execution,
- external adapter enablement,
- service restart or deployment,
- protected settings, secrets, branch protection, or repository permission
  changes,
- autonomous merge authority,
- Test 001 GREEN, docs#172 closure, autonomy increase, value allocation, or
  residual-risk closure.

Consumers may render these records. They must not wake runtime work or treat
the records as authorization for protected actions.
