# Epic 1A CapabilityArtifact Persistence And Backfill Decision

Date: 2026-05-28

Status: proposed resolution for review

Repo: `transpara-ai/eventgraph`

Issue: `transpara-ai/eventgraph#46`

## Purpose

This document records the EventGraph-side persistence-scope decision for `transpara-ai/eventgraph#46`.

The decision is intentionally narrow. It defines the current v3.9 Dark Factory `CapabilityArtifact` usage logging inventory boundary and the operator backfill guidance for legacy or externally loaded records before material capability use.

## Issue

Issue: `transpara-ai/eventgraph#46`

Title: `Dark Factory v3.9: define persistent-store CapabilityArtifact usage logging inventory`

The issue asks EventGraph to resolve the residual persistence question left after PR #45 and issue #43:

- Decide whether v3.9 Dark Factory records remain in-memory only for this surface, or identify persistent stores that must expose equivalent inventory behavior.
- If persistent-store parity is required, add adapter or query behavior and tests for those stores.
- Document the operator backfill path for legacy `CapabilityArtifact` records before material capability use.

## Decision

Current v3.9 Dark Factory capability-evidence support is `InMemoryStore`-only.

Persistent-store parity is not required now because current EventGraph does not contain a persistent v3.9 Dark Factory `Record` store. The supported v3.9 Dark Factory store surface is the package-local `InMemoryStore` in `go/pkg/darkfactory/v39`.

Generic EventGraph persistent stores persist `event.Event` through the generic `store.Store` interface. They must not be represented as supporting v3.9 Dark Factory `Record` inventory unless a future implementation explicitly bridges those models.

This decision does not claim persistent adapter or query parity for `CapabilityArtifactUsageLoggingFindings`.

Future persistent v3.9 `Record` support must add all of the following before claiming parity:

- an explicit v3.9 Dark Factory persistent store abstraction or adapter boundary
- an inventory query equivalent to `InMemoryStore.CapabilityArtifactUsageLoggingFindings`
- tests covering false and missing `usage_logging_required` records
- migration or operator backfill guidance for existing persisted records
- validation that material capability influence still requires both `CapabilityArtifact` and `USED_CAPABILITY` evidence

## Evidence

The repository evidence below is from `origin/main` at `2a924e9cfd1163638ec18cd2ce05b5a3ec99f862`.

| Evidence | File and symbol |
|---|---|
| v3.9 typed record model | `go/pkg/darkfactory/v39/schema.go:61` defines `type Record interface`. |
| v3.9 `CapabilityArtifact` schema | `go/pkg/darkfactory/v39/schema.go:605` defines `type CapabilityArtifact struct`; `go/pkg/darkfactory/v39/schema.go:619` defines `UsageLoggingRequired bool` with JSON field `usage_logging_required`. |
| v3.9 in-memory store | `go/pkg/darkfactory/v39/store.go:17` defines `type InMemoryStore struct`; `go/pkg/darkfactory/v39/store.go:28` defines `NewInMemoryStore`; `go/pkg/darkfactory/v39/store.go:40` defines `AppendRecord`; `go/pkg/darkfactory/v39/store.go:254` defines `ByType`. |
| v3.9 usage logging inventory | `go/pkg/darkfactory/v39/capability_evolution.go:479` documents `CapabilityArtifactUsageLoggingFindings`; `go/pkg/darkfactory/v39/capability_evolution.go:484` implements it on `*InMemoryStore`. |
| v3.9 material capability usage recording | `go/pkg/darkfactory/v39/capability_evolution.go:46` defines `RecordCapabilityUsage`; it rejects capability use when the target `CapabilityArtifact` lacks `usage_logging_required=true`. |
| v3.9 material capability usage path | `go/pkg/darkfactory/v39/capability_evolution.go:60` defines `CapabilityUsageEvidencePath`; the path name requires `CapabilityArtifact / USED_CAPABILITY`. |
| Certification requires capability usage evidence | `go/pkg/darkfactory/v39/query.go:376` defines `EvaluateCertificationEligibility`; `go/pkg/darkfactory/v39/query.go:404` includes `CapabilityUsageEvidencePath` in certification eligibility. |
| Generic event store abstraction | `go/pkg/store/store.go:8` documents `Store` as the event and edge persistence interface; `go/pkg/store/store.go:15` appends `event.Event`. |
| Generic Postgres store | `go/pkg/store/pgstore/pgstore.go:18` creates an `events` table with `content_json`; `go/pkg/store/pgstore/pgstore.go:61` defines `PostgresStore` as a `store.Store` backed by PostgreSQL; `go/pkg/store/pgstore/pgstore.go:156` appends `event.Event`. |
| Generic SQLite store | `go/pkg/store/sqlitestore/sqlitestore.go:21` creates an `events` table with serialized content; `go/pkg/store/sqlitestore/sqlitestore.go:55` defines `SQLiteStore` as a `store.Store` backed by SQLite; `go/pkg/store/sqlitestore/sqlitestore.go:78` appends `event.Event`. |
| Generic MySQL store | `go/pkg/store/mysqlstore/mysqlstore.go:19` creates an `events` table with serialized content; `go/pkg/store/mysqlstore/mysqlstore.go:61` defines `MySQLStore` as a `store.Store` backed by MySQL; `go/pkg/store/mysqlstore/mysqlstore.go:90` appends `event.Event`. |

The evidence shows two separate storage surfaces:

- `go/pkg/darkfactory/v39` stores typed Dark Factory `Record` values in its package-local `InMemoryStore`.
- `go/pkg/store` and its persistent implementations store generic `event.Event` values.

Because there is no current persistent v3.9 Dark Factory `Record` store, there is no persistent `CapabilityArtifact` inventory adapter or query surface to implement in this step.

## Operator Backfill Guidance

Before material capability use, operators should run or use `InMemoryStore.CapabilityArtifactUsageLoggingFindings` against the current v3.9 `InMemoryStore` surface.

Treat any `CapabilityArtifact` with missing or false `usage_logging_required` as requiring operator review and backfill before material capability influence.

The v3.9 store is append-only. Backfill must re-emit a new compliant `CapabilityArtifact` record with `usage_logging_required=true` according to v3.9 invariants before capability use; existing records are not mutated in place. The re-emitted record must still preserve the artifact identity, source, content hash or immutable locator, owner, risk class, activation scope, review reference, rollback reference, and any evidence references required by the current v3.9 schema.

Material capability influence remains invalid without both:

- a valid `CapabilityArtifact` record
- `USED_CAPABILITY` evidence recorded through the v3.9 capability usage path

Do not treat raw `CapabilityArtifact` presence as sufficient evidence of material use.

If a future persistent v3.9 `Record` store is introduced, this in-memory procedure is not sufficient. The persistent store must expose equivalent inventory and backfill support before persistent parity can be claimed.

## Scope Boundary

This decision introduces no new persistent v3.9 `Record` store.

No adapter or query parity is implemented because there is no current persistent v3.9 `Record` store surface.

This decision changes no capability activation, promotion, rollback automation, external runtime, Hive, Work, Site, or Agent behavior.

This decision does not implement or authorize Epic 3.

This decision does not close Gate B globally unless `transpara-ai/eventgraph#46` is accepted as resolved by review and merge, and docs are reconciled later.

## Validation

Validation results for this decision PR:

```text
cd go
go test ./pkg/darkfactory/v39 -- PASS
go test ./... -- PASS
go build ./... -- PASS
go vet ./... -- PASS
go run honnef.co/go/tools/cmd/staticcheck@latest ./pkg/darkfactory/v39/... -- PASS
cd ..
git diff --check -- PASS
```

## Closure Condition

`transpara-ai/eventgraph#46` may be closed as completed only if reviewers accept the `InMemoryStore`-only boundary decision and operator backfill guidance as satisfying the issue acceptance criteria.

If reviewers require persistent-store parity, this decision must instead create or reference a narrower implementation issue for the missing persistent v3.9 `Record` store abstraction, inventory query, tests, and migration or backfill plan.
