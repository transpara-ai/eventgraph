# Docs#226 EventGraph Evidence-Source Proposal Packet - 2026-07-09

- **doc_id:** EG-DOCS-226-EVIDENCE-SOURCE-PROPOSAL-2026-07-09
- **status:** proposal/design only
- **issue:** eventgraph#79
- **source authority:** docs#226 AuthorityDecision for `AUTHREQ-DOCS-226-EVENTGRAPH-EVIDENCE-SOURCE-2026-07-09-001`
- **primary repo:** `transpara-ai/eventgraph`
- **coordination repo:** `transpara-ai/operation`

## 1. Purpose

This packet defines the EventGraph evidence-source proposal lane authorized by
the closed `transpara-ai/docs#226` authority-routing record and tracked by
`transpara-ai/eventgraph#79`.

The packet is intentionally proposal-only. It defines fields, classifications,
validation expectations, and audit posture for a future EventGraph evidence
source. It does not authorize implementation, runtime execution, production
EventGraph reads or queries, production EventGraph writes, EventGraph schema or
migration changes, proof collection, issue disposition, Test 001 state change,
production go-live, value allocation, autonomy increase, or wiki work.

## 2. Source Records

| Record | Role |
| --- | --- |
| `transpara-ai/docs#226` | Closed authority-routing tracker for the EventGraph/runtime operation path. |
| `docs#226` issuecomment-4921675575 | Source-of-intent for the proposal-only EventGraph evidence-source AuthorityRequest lane, listed in the AuthorityDecision source records; authored by `MichaelSaucier`, created and last updated at `2026-07-09T04:51:18Z`, body hash `sha256:4838e29a69d2a0ace9e597d7dbe9fbf547a20d77e56cd7115954f3b399098d92`. |
| `docs#226` issuecomment-4922119091 | Human AuthorityDecision approving one future, separately reviewed EventGraph evidence-source child proposal packet; authored by `MichaelSaucier`, created and last updated at `2026-07-09T06:15:50Z`, body hash `sha256:e2f7c968c1efe5f18fa30b87cb933f50af165a2d82f41c8a06dd1d3dcbfd1517`. |
| `docs#226` issuecomment-4922122423 | Successor handoff naming `transpara-ai/eventgraph#79` as the child tracker; authored by `MichaelSaucier`, created and last updated at `2026-07-09T06:16:26Z`, body hash `sha256:3a49590ca2615457fa1d9a61e0dd04f4747556d9a9d22bb6b80b06ca87016552`. |
| `docs#226` issuecomment-4922462961 | Closure receipt for `docs#226`; remaining proposal/design work is owned by `eventgraph#79`; authored by `MichaelSaucier`, created and last updated at `2026-07-09T06:58:53Z`, body hash `sha256:2d4c7799b9ac35a2925350023bfd01de6fcf4131c80a109804a8a9c31ea277e8`. |
| `operation#61` | Merged proposal-only authority-path packet; reviewed head `1b06968dfcfb518d4022ca15d5fe9833ce490992`, merge commit `ea739d82494ab09df2d02378e0b53b00b4ba30fb`. |
| `operation#62` | Merged EventGraph evidence-source AuthorityRequest packet; reviewed head `9d99accb0c2a7f916b4888d5d12b01192c3aa57b`, merge commit `28b48c429bb9c07346ccebce7d1272de5c24c8de`. |
| `operation#26` | Open Test 001 YELLOW/live-evidence tracker; no closure or GREEN state is authorized here. |
| `operation#45` | Closed MVP private/firewalled proof tracker; no disposition change is authorized here. |

The GitHub issue comments above are source records. This reviewed packet is the
canonical EventGraph-side proposal-shape record for the field and audit shape.
`operation#61` and `operation#62` are authority-routing and AuthorityRequest
source records; neither consumed the one child-proposal slot authorized by
`docs#226` issuecomment-4922119091.

## 3. Authority Boundary

Approved action from the AuthorityDecision:

```text
Authorize one future, separately reviewed EventGraph evidence-source child
proposal packet.
```

On merge, this packet consumes that one authorized child-proposal slot for
`eventgraph#79` and satisfies that proposal/design lane only. If this PR is
closed unmerged, the slot remains unconsumed until a successor `eventgraph#79`
proposal packet is separately reviewed. Any downstream EventGraph
evidence-source proposal, implementation, live evidence collection, or closure
claim requires a fresh AuthorityRequest and human AuthorityDecision. This packet
authorizes no later action by implication.

Non-authorizations preserved:

- no production EventGraph read or query;
- no production EventGraph write;
- no EventGraph schema or migration change;
- no RuntimeBroker execution;
- no external adapter invocation;
- no Hive wake/start/action API;
- no Work write or Site write;
- no deploy or service restart;
- no private fetch, authentication, or production connection use;
- no protected settings, secrets, ruleset, CODEOWNERS, or branch-protection
  change;
- no Test 001 `GREEN`;
- no `operation#26` closure;
- no `operation#45` disposition change;
- no production go-live;
- no value allocation;
- no autonomy increase;
- no wiki work.

Any future request that needs one of those actions must start from a separate
AuthorityRequest and a separate human AuthorityDecision.

## 4. Evidence-Source Proposal Fields

This authorized child proposal defines the required field shape below. Any
downstream EventGraph evidence-source implementation or live-evidence packet must
define every field below before it can claim readiness.

| Field | Required value |
| --- | --- |
| Proposal id | Stable id scoped to `docs#226`, `eventgraph#79`, date, and target child lane. |
| Source authority | Exact issue comments, merged packets, reviewed heads, and merge commits used as authority inputs. |
| Source/store classification | One of the classifications in section 5. Missing classification means deferred. |
| Evidence source | Exact repo, path, fixture, export, service, query target, or source-record reference. |
| Environment | Documentation/source-record environment by default; any local, private/firewalled, staging, or production environment must be separately authorized. |
| Actor and supervision | Human decider, supervised proposal author, reviewer, evidence recorder, and rollback/correction owner. |
| Credential boundary | `none` by default; any credential, secret, authentication, or production connection must be separately authorized. |
| Data class | Public-safe governance source records by default; private/firewalled or production data must be separately authorized. |
| Schema state | Exact schema/version/migration state, or `not applicable for source-record-only`. |
| Event identifiers | Exact event ids, or `not produced` for proposal-only/source-record-only work. |
| Chain root | Exact root/predecessor/anchor, or `not available`; missing live chain data cannot be converted into evidence by summary. |
| Chain verification method | Exact command or manual method, expected output, and failure interpretation. |
| Export/query evidence | Exact command, artifact path, checksum where available, data class, and freshness rule. |
| Dispositive classification | `source-record-only`, `advisory`, `incident-dispositive`, `insufficient`, or `unavailable`. |
| Rollback/correction posture | How stale, bad, or wrongly displayed evidence is corrected, superseded, or marked unavailable. |
| Audit storage | Exact path, PR comment, issue comment, or artifact location for the AuditReport. |
| Validation plan | Local validation target, negative checks, CFADA/CFAR gates, and exact-head approval requirements. |
| Non-authorizations | Every surviving forbidden action from section 3 repeated in the child proposal. |
| Expiry | Date, condition, successor decision, or reason the proposal remains pending. |

No missing field may be inferred as satisfied.

`source-record-only` evidence may not select `incident-dispositive` for any
incident. Incident-dispositive status requires separately authorized live or
fixture chain proof, because source records do not establish EventGraph
reachability, chain continuity, store query results, or write truth.

## 5. Source/Store Classification

The child proposal must select exactly one current classification for each
evidence source.

| Classification | Proposal use | Readiness rule |
| --- | --- | --- |
| Source-record only | Merged docs, Operation, EventGraph records; issue comments; PR bodies and review comments. | Allowed in this packet. Cannot claim live EventGraph evidence. |
| Local fixture/export | A local fixture, generated export, or checked-in artifact. | Requires exact path, command, checksum where available, and reason it is not production evidence. |
| Private/firewalled read/export | A read from an allowlisted private or firewalled operator surface. | Not authorized here. Requires exact target, actor, capture method, credential boundary, storage path, and explicit future authority. |
| Production read/query | A query or export from a production EventGraph store. | Not authorized here. Requires exact store, query/export command, credential source, data class, time window, output path, and explicit future production-read authority. |
| Production write | A new production EventGraph event or chain mutation. | Not authorized here. Requires separate AuthorityRequest and AuthorityDecision naming the exact write, schema, chain root, rollback/correction process, and production-write authority. |
| Schema or migration change | EventGraph schema, migration, or storage contract mutation. | Not authorized here. Requires separate implementation issue, validation plan, migration/rollback plan, and exact human approval. |

If a proposal needs multiple classifications, it must either split the work or
make the cross-class boundary explicit and keep unauthorized classifications
blocked.

## 6. Chain-Verification Expectations

A future EventGraph evidence source must distinguish source-record provenance
from EventGraph chain proof.

Source-record-only proof may cite:

- issue comments and timestamps;
- merged PR heads and merge commits;
- PR-visible CFADA/CFAR evidence;
- local validation outputs;
- OpenBrain checkpoint references only when a cited GitHub source record already
  memorializes the checkpoint boundary. An OpenBrain checkpoint remains a
  source-record citation; it is not EventGraph chain proof.

It must not claim:

- event ids were produced;
- a chain root was observed;
- a production store was queried;
- a production write occurred;
- chain continuity was verified.

Live or fixture chain proof must name:

- event ids and predecessor ids;
- chain root or anchor;
- schema version and store implementation;
- producer actor and authority;
- verification command or manual method;
- expected pass/fail output;
- artifact path and checksum where available;
- freshness or expiry rule.

Missing chain fields make the evidence `insufficient` or `unavailable`; they do
not become proof through narrative summary.

## 7. Export/Query Evidence Requirements

Any future export/query evidence must include:

- exact command or manual capture method;
- target store, fixture, file, or source record;
- environment and network boundary;
- credential source or explicit `none`;
- data class and redaction rule;
- output path;
- checksum where feasible;
- timestamp and freshness window;
- actor;
- failure interpretation;
- storage path for the AuditReport.

For this proposal packet, every export/query action is `not performed`.

## 8. Missing-Field Behavior

The lane fails closed when any of these are missing, ambiguous, stale, or
contradicted by a surviving non-authorization:

- source authority;
- classification;
- actor or supervision boundary;
- environment;
- credential boundary;
- data class;
- schema state;
- event ids or explicit `not produced`;
- chain root or explicit `not available`;
- verification command or manual method;
- export/query artifact path or explicit `not performed`;
- rollback/correction posture;
- audit storage;
- validation plan;
- exact human approval point.

The only valid readiness outcome for a child proposal with a missing required
field is `deferred`, `blocked`, `unavailable`, or `source-record-only`, as
appropriate. Missing fields must not be treated as operational proof.

## 9. Rollback And Correction Posture

For this proposal/design packet, rollback is limited to:

- reverting or superseding the EventGraph design PR;
- correcting this packet with a later reviewed EventGraph packet;
- recording a corrective issue or PR comment that names the stale field and the
  superseding source record.

No runtime rollback command exists because no runtime action is authorized.

For any future live or fixture evidence packet, rollback/correction must state:

- who can suspend or supersede the evidence;
- how bad or stale EventGraph evidence is marked non-dispositive;
- how display or export claims are corrected;
- how successor records preserve the original audit trail.

## 10. AuditReport Requirements

A future PR that implements this proposal/design lane must include an
AuditReport or PR-visible audit section that records:

- source issue, successor tracker, and AuthorityDecision refs;
- branch, exact head SHA, changed files, and scope comparison;
- source/store classification review;
- evidence-source field completeness;
- chain-verification field completeness;
- local validation and negative checks;
- read-only change-control aggregation scan;
- CFADA and CFAR findings/dispositions;
- issue posture for `eventgraph#79`, `docs#226`, and `operation#26`;
- statement that Test 001 remains `YELLOW` unless separately authorized;
- statement that no production EventGraph read/query/write, schema/migration,
  runtime execution, deploy, protected settings change, value allocation,
  autonomy increase, or wiki work occurred.

## 11. Validation Expectations

Minimum validation for a proposal/design-only PR:

- inspect live `eventgraph#79`, closed `docs#226`, and open `operation#26`;
- run the read-only change-control aggregation scan;
- run `make verify` when feasible, or the narrow validation target justified by
  touched files;
- when validating from a linked worktree where Go VCS stamping cannot read VCS
  status, rerun the same validation as `GOFLAGS=-buildvcs=false make verify`
  and record the original failure and workaround in PR evidence;
- run CFADA because the packet records authority and evidence boundaries;
- run draft-state CFAR;
- after ready transition, run ready-state exact-head CFAR;
- post PR-visible evidence for every gate;
- request exact-head human merge approval before merge.

PR bodies and commit messages must not use closing keywords for `docs#226`,
`operation#26`, or `operation#45`. `eventgraph#79` may close only through a
separately approved exact-head merge that resolves this proposal/design issue.

## 12. Current Disposition

After this packet:

- `docs#226` remains closed as completed; its remaining EventGraph proposal work
  is owned by `eventgraph#79`;
- `eventgraph#79` remains open until the proposal/design PR is merged with
  exact-head approval;
- `operation#26` remains open and Test 001 remains `YELLOW`;
- `operation#45` disposition remains unchanged;
- no implementation issue is made PR-ready beyond proposal/design work;
- no production EventGraph read/query/write, schema/migration, RuntimeBroker
  execution, adapter invocation, Hive action, Work write, Site write, deploy,
  service restart, private fetch, authentication, production connection,
  protected settings change, production go-live, value allocation, autonomy
  increase, or wiki work is authorized.
