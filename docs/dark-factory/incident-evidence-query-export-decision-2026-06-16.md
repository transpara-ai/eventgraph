# Incident Evidence Query And Export Decision

Date: 2026-06-16

Status: accepted on merge

Repo: `transpara-ai/eventgraph`

Source contract: `transpara-ai/civilization-operation` `docs/operations/eventgraph-incident-record-semantics.md`

## Purpose

This document defines the EventGraph-side query and export behavior that cross-repository incident records may cite during pre-live operation.

The decision is intentionally narrow. It accepts a bounded evidence packet over existing EventGraph query, causal traversal, and hash-chain verification APIs. It does not add a new incident product workflow, mutate runtime behavior, define source-repository truth, or grant authority to close incidents.

## Decision

EventGraph accepts the following existing APIs as the source-repo query/export surface for incident evidence packets:

- exact event lookup with `Store.Get`
- recent, type, source, conversation, and since queries with `Store.Recent`, `Store.ByType`, `Store.BySource`, `Store.ByConversation`, and `Store.Since`
- causal context with `Store.Ancestors` and `Store.Descendants`
- chain integrity evidence with `Store.VerifyChain`
- graph facade queries that wrap the same store behavior through `Query.Recent`, `Query.ByType`, `Query.BySource`, `Query.ByConversation`, `Query.Ancestors`, `Query.Descendants`, and `Query.EventCount`

Incident records may cite an EventGraph evidence packet when the packet records the exact query or export behavior used, the exact EventGraph commit or deployment being queried, and the exact event references returned.

This decision accepts the packet shape below as the minimum source-repo export format:

```text
packet_id:
created_at:
eventgraph_repo:
eventgraph_commit:
store_implementation:
query_surface:
query_method:
query_parameters:
result_limit_or_depth:
pagination_cursor_if_supported:
chain_verification:
event_refs:
causal_context:
classification_requested_by_incident:
source_artifact_refs:
missing_evidence:
operator_notes:
```

Each `event_refs` item must include:

```text
event_id:
event_type:
source:
conversation_id:
timestamp:
causes:
hash:
prev_hash:
signature:
content_fields_used:
source_artifact_refs:
```

Each `causal_context` item must state whether it came from `Ancestors`, `Descendants`, or a direct event result.

`classification_requested_by_incident` must use the civilization-operation vocabulary:

- `DISPOSITIVE`
- `ADVISORY`
- `INSUFFICIENT`

EventGraph may provide the event data needed for the incident record to classify a record, but the incident record owns the classification. The canonical classification definitions live in the source contract. The incident record must say which incident question the event answers and why the evidence is dispositive, advisory, or insufficient under that contract.

## Evidence Boundary

EventGraph can prove:

- an event with a specific id exists in the queried store
- the event has a specific type, source, conversation, timestamp, causal predecessor list, hash, previous hash, signature value, and typed content
- a query by type, source, conversation, recency, or event id returned a specific bounded result set
- causal ancestors or descendants exist for a cited event within the requested depth
- the queried store's hash chain verifies at the time of export

EventGraph does not by itself prove:

- active `hive` runtime behavior or deployment state
- `work` task execution state
- rendered `site` or `civilization-wiki` user-facing state
- canonical doctrine truth in `docs`
- human authorization for a protected action
- incident closure

Those claims require source-repository artifacts, deployment evidence, rendered user-facing evidence, or human authorization artifacts outside EventGraph.

## Incident Evidence Rules

The rules below are EventGraph-side preconditions layered on the canonical source-contract definitions. They do not redefine the `DISPOSITIVE`, `ADVISORY`, or `INSUFFICIENT` vocabulary.

An EventGraph evidence packet may be cited as `DISPOSITIVE` by the incident record only when the event content, event metadata, causal context, source artifact references, and chain verification together settle that specific question without relying on chat memory or unstated interpretation.

An EventGraph evidence packet does not meet the EventGraph-side dispositive precondition when it narrows investigation but lacks a source artifact, deployment observation, user-facing surface, required event content field, causal path, or authorization artifact needed to settle the incident question. The incident record classifies that packet under the source contract.

An EventGraph evidence packet does not meet the EventGraph-side dispositive precondition when the relevant query returns no result, the chain cannot be verified, the event content does not contain the required fields, or the event is not tied to the source artifact named by the incident question. The incident record classifies that packet under the source contract.

No packet may be cited as dispositive unless `chain_verification` records a successful `VerifyChain` result. A packet without successful chain verification may be cited only as `ADVISORY` or `INSUFFICIENT` unless a later EventGraph decision, reviewed and merged in `transpara-ai/eventgraph`, explicitly defines a different authority path.

## Evidence

The repository evidence below is from `origin/main` at `b7fca161fc2cdb57513317fbb6b711497f247bfd`.

| Evidence | File and symbol |
| --- | --- |
| Store query/export boundary | `go/pkg/store/store.go` defines `Store` as the event and edge persistence interface. |
| Exact event lookup | `go/pkg/store/store.go` defines `Get(id types.EventID)`. |
| Bounded query methods | `go/pkg/store/store.go` defines `Recent`, `ByType`, `BySource`, `ByConversation`, and `Since`. |
| Causal context methods | `go/pkg/store/store.go` defines `Ancestors` and `Descendants`. |
| Query bounds and cursors | `go/pkg/store/store.go` defines `limit` and `after` cursor parameters for paged queries, `limit` for `Since`, and `maxDepth` for causal traversal. |
| Chain integrity method | `go/pkg/store/store.go` defines `VerifyChain`. |
| Graph query facade | `go/pkg/graph/query.go` wraps `Recent`, `ByType`, `BySource`, `ByConversation`, `Ancestors`, `Descendants`, and `EventCount`. |
| Event metadata accessors | `go/pkg/event/event.go` exposes `ID`, `Type`, `Timestamp`, `Source`, `Content`, `Causes`, `ConversationID`, `Hash`, `PrevHash`, and `Signature`. |
| Canonical hash input | `go/pkg/event/event.go` defines `CanonicalForm` and `ComputeHash`. |

## Relationship To Civilization Operation Test 001

This source-repo adoption slice resolves the first-pass EventGraph query/export behavior gap for Test 001:

```text
Incident records may cite an exact EventGraph evidence packet built from accepted query, causal traversal, event metadata, and chain-verification behavior.
```

It does not make Test 001 `GREEN`. The tabletop still needs a real or simulated incident to cite actual event records or to record the accepted absence of incident-dispositive event records. Other source repositories still need to supply runtime, work-state, user-facing, roster, authority, or closure evidence when those truth types are at issue.

## Scope Boundary

This decision introduces no new event type, store implementation, CLI command, deployment, runtime mutation, authority policy, or incident-closure authority.

This decision does not make EventGraph the source of truth for source repositories that own runtime behavior, work execution, public surfaces, doctrine, or human authorization.

This decision does not require persistent-store schema changes because it accepts the existing generic Store and Query surfaces as the source-repo evidence boundary.

## Validation

Validation for this decision PR should include:

```text
git diff --check
cd go && go test ./pkg/store ./pkg/graph ./pkg/event
```

Full repository CI remains the merge gate for cross-language regressions.

## Closure Condition

This decision is accepted only after review and merge in `transpara-ai/eventgraph`.

After merge, `civilization-operation` may cite the merge commit as EventGraph source-repo adoption evidence for the incident query/export behavior slice. A later civilization-operation PR must still reconcile Test 001 without claiming `GREEN` until actual event records or accepted missing-evidence findings are cited.
