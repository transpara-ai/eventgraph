# Model Catalog Currency + Explicit Selection-Mode Projection — Design Packet

- **doc_id:** EG-MODELCONFIG-CATALOG-MODE-DESIGN-001
- **version:** v0.4.0 (CFADA PASS)
- **status:** CFADA PASS (round 3, codex, 0 blockers) — building under TDD
- **issues:** eventgraph#74 (root), hive#243, site#206
- **repos/bases:** eventgraph main @ (pin at build), hive main @ 02ae3d4, site main @ b68e214
- **home:** this packet lives in eventgraph `docs/designs/`; the hive and site PRs reference it (single source of truth — same convention as the recent-intakes rail pair).

## 1. Problem

1. **Catalog currency.** `go/pkg/modelconfig/defaults_catalog.yaml` tops out at Opus 4.6 / Sonnet 4.6 / Haiku 4.5. The local claude CLI (2.1.198) verifiably serves `claude-opus-4-8` and `claude-fable-5` (headless probes 2026-07-02, both `OK`, exit 0). Neither is reachable via resolver id/alias/tier selection. hive's `catalog-mixed.yaml` has no current Claude entries at all.
2. **Provenance dishonesty in mode labels.** Site's mirror types (`OpsHiveModelRoleAssignment.SelectionMode/OverrideMode/EffectiveMode`, `OpsHiveModelSelection.GlobalMode/SelectionMode` — graph/ops.go) are **declared but never populated**: hive's `OperatorModelRoleAssignment` (operator_model_projection.go:236-254) carries no mode field, so site's `obsAssignmentModelModeState` (observatory.go:823-843) always falls through to inference — "Manual · inferred" from `Model`/`PolicyModel` presence, "Auto · inferred" from `SelectionStrategy`/`PreferredTier` presence. The resolver KNOWS how each selection was made (its `Trace` records every layer, resolve.go:17) but that knowledge is discarded at the projection boundary.

## 2. Decisions

### D1 — Catalog additions (eventgraph + hive), sourced pricing, corrections declared

Add to `defaults_catalog.yaml` `models:` (claude-cli provider, mirroring existing schema):

- `claude-opus-4-8` — aliases `[opus-4-8]`; tier `judgment`; capabilities as opus-4-6 (tools, reasoning, coding, vision, operate, large-context, structured-output); context_window **1000000**; max_output_tokens 16384 (operational choice, matches file convention); pricing input 5.00 / output 25.00 / cache_read 0.50 / cache_write 6.25.
- `claude-fable-5` — aliases `[fable, fable-5]`; tier `judgment`; same capability set; context_window **1000000**; max_output_tokens 16384; pricing input 10.00 / output 50.00 / cache_read 1.00 / cache_write 12.50.

Pricing source: platform.claude.com/docs/en/about-claude/pricing (fetched 2026-07-02): Opus 4.8 = $5/$25, cache $0.50/$6.25; Fable 5 = $10/$50, cache $1/$12.50; both list the full 1M context window at standard pricing. **No numbers invented.**

**Declared correction (D1a):** the existing `claude-opus-4-6` entry prices at 15.00/75.00 with 200000 context; published pricing for Opus 4.6 is $5/$25 (cache $0.50/$6.25) and the 1M window is standard. Correct the entry (and its `api-` variant if present) in the same PR, declared in the PR body with the source. Same class, fixed in the same pass. Other models' entries are checked against the same fetched table; any further discrepancies are corrected and declared, or explicitly listed as left-as-is with a reason.

**NOT changed:** `tier_defaults` (judgment stays `claude-opus-4-6`), `role_defaults`, `profiles` — selection behavior is Michael's operational dial; this packet only makes the models *available*. The **existing alias `opus` stays on 4-6** (no silent rebinding of a live alias).

hive `catalog-mixed.yaml`: add the same two entries (claude-cli, subscription). Its role defaults are untouched.

**Fail-closed guard:** aliases must not collide with existing ids/aliases — `NewCatalog` already rejects duplicates (catalog_test.go pins this); tests assert both new ids and all new aliases resolve, and that `opus` still resolves to `claude-opus-4-6`.

### D2 — SelectionMode: a closed vocabulary type in modelconfig, derived at resolution time

New type in `go/pkg/modelconfig` (types.go or resolve.go):

```go
type SelectionMode string
const (
    SelectionModeManualExplicit SelectionMode = "manual-explicit" // the winning token resolved as a concrete catalog id or alias (incl. ModelAliases remap)
    SelectionModeAutoTier       SelectionMode = "auto-tier"       // the winning token fell back through tier resolution (tier-name token or PreferredTier)
    SelectionModeSystemDefault  SelectionMode = "system-default"  // no caller input won: catalog/tier/role defaults chose the model
    SelectionModeUnknown        SelectionMode = ""                // absent — consumers MUST treat as not-projected, never coerce
)
```

`ResolvedConfig` gains `Mode SelectionMode` (json `mode,omitempty`). **Derivation happens at the RESOLUTION POINT of the winning token, not at the layer that supplied it (CFADA1-1):** `Policy.Model`/`TaskOverride.Model`/`AgentDefModel` may hold a model id, an alias, **or a tier name** (policy.go:7), and resolve.go:127 falls back from an unresolved token to tier defaults — so "a policy supplied the token" proves nothing about HOW it resolved. The rules, allowlist-form:

- winning token resolves as a concrete catalog id or alias (resolve.go:118 ModelAliases remap included, resolve.go:123 catalog lookup) → `manual-explicit` (a remapped alias is still an explicit pin);
- winning token falls through tier resolution (resolve.go:127), or `PreferredTier` selects → `auto-tier`;
- resolution completed purely from system/tier/role defaults with no caller-supplied token winning (resolve.go:88) → `system-default`;
- `RequiredCapabilities`, `SelectionStrategy` (advisory, policy.go:33), `AllowDowngrade`, `MaxCostPerCallUSD` NEVER set mode — they filter/veto, they do not choose;
- no rule fires → `SelectionModeUnknown` (no default branch assigns a mode).

Tests: the existing `TestResolve` precedence table extended with mode assertions per row, PLUS the five discriminating cases from CFADA round 1: defaults-only → `system-default`; concrete id/alias pin → `manual-explicit`; **tier-name token via Policy.Model (e.g. "judgment") → `auto-tier`**; PreferredTier-driven selection → `auto-tier`; ModelAliases remap → `manual-explicit`. Trace remains the cross-check oracle (mode must be consistent with the winning trace entry).

`RoleModelPolicy` is unchanged (SelectionStrategy stays advisory). No behavior change to resolution itself — Mode is observational output only.

### D3 — hive projection: carry mode; presence IS the gate (no new schema-version machinery)

`OperatorModelRoleAssignment` gains `SelectionMode string` (json `selection_mode,omitempty`), populated from the eventgraph `ResolvedConfig.Mode` computed when the assignment is built (`operatorModelRoleAssignment()`, operator_model_projection.go:341+). Where an assignment is built WITHOUT running Resolve (error paths, policy-event-only rows), the field stays empty — **never synthesized**. `OperatorModelSelection` gains no global mode (there is no global resolver-level mode fact today; projecting one would be invention — the site's global chip keeps its current derivation).

**Gating decision:** the model-selection projection has no `schema_version` field today (only the civilization assembly projection does). Rather than introducing version machinery for one optional string, **field presence is the gate**: old hive → field absent → site's explicit branch never fires → existing inference behavior, byte-identical. New hive → field present → explicit wins. `omitempty` keeps old-payload byte-compatibility. This is fail-closed by construction (absence degrades to today's behavior, never to a wrong explicit label).

### D4 — site: explicit wins; inference retained and relabeled honestly; unknown fails closed

`OpsHiveModelRoleAssignment.SelectionMode` (already declared) now actually receives the hive value. `obsCanonicalModelMode` (observatory.go:879-887) extends its allowlist: `manual-explicit` → `Manual`; `auto-tier` → `Auto`; `system-default` → `Auto` (a default IS automatic selection; the binary chip stays Manual/Auto).

**Present-but-invalid fails closed BEFORE inference (CFADA1-2):** the current fall-through would let an unrecognized non-empty `selection_mode` drop into the legacy inference branches (or be masked by `OverrideMode`, which is checked first) and render a confident `Manual · override` / `Auto · inferred`. `obsAssignmentModelModeState` is restructured so the NEW projected field is evaluated FIRST with three distinct outcomes: (a) valid vocabulary → mode with provenance `explicit`, done; (b) **non-empty but not allowlisted → return `unknown` / `not projected` immediately — no legacy fallback, no inference** (present-invalid ≠ absent); (c) empty/absent → the legacy chain (`OverrideMode` → `EffectiveMode` → policy-event → strategy → global → model-presence inference) runs unchanged, byte-identical to today. This also resolves the conflict-precedence question: the projected resolver mode, when present, always precedes the legacy site-side fields — `OverrideMode` can never mask it.

**Invalid-unknown is STICKY through global inheritance (CFADA2-1):** the roster assembly at observatory.go:588 replaces an assignment-level `unknown` with the inherited GLOBAL mode whenever global provenance is projected — which would launder a present-invalid `selection_mode` into an inherited `Auto`. The inheritance branch becomes allowlist-gated: it fires ONLY when the assignment carried no projected mode data at all (genuinely-absent legacy case, byte-identical to today). Present-invalid yields a distinct provenance sentinel (`invalid projection`, rendered as `not projected`) that the inheritance branch never overwrites. Build-level test: an assignment with invalid `SelectionMode` PLUS model/policy/global metadata present must still render `unknown · not projected` end-to-end (through buildObs/roster assembly, not just the state function); empty `SelectionMode` must keep today's inherit behavior byte-identical.

Table tests: absent, each valid value (raw `selection_mode` preserved in the projection so the binary chip never erases provenance — CFADA2-adv3), unknown-with-model-present, unknown-with-policy-event, conflicting override/selection pairs, and the end-to-end sticky-invalid case above. D2's test table additionally gains role-default and `Policy.Profile` mode rows (real model-writing branches at resolve.go:95/:217 — CFADA2-adv2), and the implementation note that mode derivation must track the winning SOURCE CLASS through Resolve/applyPolicy rather than inferring from catalog lookup alone, since ModelAliases remap applies to every source (CFADA2-adv1).

Provenance: when the projected `SelectionMode` decides the label, provenance is `explicit` (the existing `EffectiveMode/SelectionMode` branch at observatory.go:827-829 already does exactly this — it starts firing for real). Inference branches keep their `inferred` labels for legacy upstreams. Templ chip copy already renders mode + provenance — expected visible change on live data: "Manual · inferred" → "Manual · explicit" wherever hive actually resolved the assignment. Visual evidence (before/after screenshots) in the site PR.

### D5 — Verification

- eventgraph: catalog tests (both ids + all aliases resolve; `opus` alias unchanged; duplicate-guard still green; cheapest-with-capabilities TIE note: opus-4-6 and opus-4-8 both land at output 25.00 and catalog.go:96 uses strict `<`, so first-in-catalog wins ties — asserted explicitly, CFADA1-adv1); Mode-derivation table asserted against Trace for the full existing precedence table + the five CFADA discriminating cases + unknown-mode fail-closed case; `claude --model <id> -p` probe evidence quoted in the PR (not executed in tests).
- hive: projection round-trip test (assignment built via Resolve carries the mode; error-path assignment carries empty); JSON `omitempty` shape test (absent when empty); a merged-catalog test (catalog-mixed.yaml entries merge over embedded defaults by ID via MergeCatalogs, defaults.go:217 — new entries must survive the merge and resolve, CFADA1-adv3); existing policy tests untouched and green.
- site: mode-state table test extended: explicit-wins (all vocabulary values), absent-field → inference fallback (byte-identical to today's labels), unknown-value → "unknown · not projected"; `templ generate` + `make verify`; live screenshots.
- Cross-repo: site builds against the unmerged hive branch, hive against the unmerged eventgraph branch (replace directives already wire eventgraph); dependency ordering stated in ALL THREE PR bodies (eventgraph → hive → site).

## 3. Non-goals

No tier_default/role_default/profile changes; no alias rebinding; no resolution-behavior change; no schema-version framework for the model-selection projection; no site global-mode chip change; no removal of the inference fallback (legacy upstreams stay honest under `inferred`).

## 4. TDD plan

1. eventgraph: catalog entries + tests (RED: resolve-by-id fails; GREEN after YAML). Opus-4-6 correction with source note.
2. eventgraph: SelectionMode type + Resolve derivation + trace-consistency tests (RED first).
3. hive: catalog-mixed entries; projection field + round-trip/omitempty tests (RED first).
4. site: ops mirror population + obsCanonicalModelMode allowlist + mode-state table tests (RED first); templ; screenshots.
5. Each repo: full test suite + vet; conventional commits with Claude Fable 5 trailer.

## 5. IADA record (v0.1.0, 2026-07-02)

- **IADA-1 (mode must follow the WINNING layer):** a policy can set BOTH `Model` and `PreferredTier`; deriving mode from "any layer that ran" would mislabel. Resolved: `Mode` is set exactly by the precedence layer that produced the final model (the same point that writes the winning `Trace` entry), and the trace-consistency test is the oracle — for every row of the existing `TestResolve` precedence table, mode must match the winning layer's class. No layer wins ⇒ unknown.
- **IADA-2 (D1a pricing correction is NOT purely observational):** catalog pricing feeds `CheapestWithCapabilities` and `MaxCostPerCallUSD` cost gates — correcting opus-4-6 from 15/75 to 5/25 can change which models pass a cost gate (a policy whose budget excluded opus-4-6 at $75/MTok output may now include it). Resolved: declared in the eventgraph PR body as a behavior-relevant correction; tests enumerate the affected selection paths (cheapest-with-capabilities ordering, cost-gate pass/fail at the old vs new price) so the flip is visible, not silent. Cheapness ORDER among tiers is unchanged (haiku < sonnet < opus 4.6 < fable) — asserted.
- **IADA-3 (correction scope kept minimal):** opus-4-6's `context_window` (200000) is NOT bumped even though docs list 1M for Opus 4.6 — pricing was the flagged dishonesty; silently changing large-context capability semantics for an in-use model is a separate decision. Noted in the PR body as known-stale with the source. The two NEW entries carry 1000000 (docs-backed, no in-use behavior to disturb).
- **IADA-4 (alias hygiene):** new aliases (`opus-4-8`, `fable`, `fable-5`) checked against every existing id/alias; `NewCatalog`'s duplicate guard is the fail-closed backstop and the tests assert `opus` still resolves to `claude-opus-4-6` (no silent rebinding of a live alias).
- **IADA-5 (presence-gate blindspot):** a middlebox/old-site pairing (new hive, old site) is safe — old site's mirror already declares the field; unmarshalling an extra JSON key is a no-op even where the mirror lacked it. New-site/old-hive degrades to inference (the point of the gate). Both pairings enumerated in the site test table.

## 6. CFADA record

### Round 3 (codex, 2026-07-02) — VERDICT: PASS (0 blockers)

Codex verified the CFADA2-1 resolution sound and implementable against observatory.go:583-600 with legacy behavior byte-identical. Display caveat adopted: templates print provenance directly, so the sticky `invalid projection` sentinel needs a display mapping to render as `not projected` (helper at the render boundary, sentinel preserved internally for the inheritance gate).

### Round 2 (codex, 2026-07-02) — VERDICT: BLOCKERS (1) → resolved in v0.3.0

- **CFADA2-1 (global inheritance laundered invalid modes):** observatory.go:588 replaces assignment-level `unknown` with the inherited global mode — a present-invalid `selection_mode` could re-enter as `Auto` despite the three-way split. Resolved: inheritance allowlist-gated on genuinely-absent; present-invalid carries a sticky `invalid projection` sentinel that inheritance never overwrites; end-to-end build-level test mandated (see D4).
- Advisories adopted: mode derivation tracks winning source class through Resolve/applyPolicy, never catalog-lookup-alone (adv1); role-default + Policy.Profile mode test rows added (adv2); raw `selection_mode` preserved through the projection so the Auto/Manual chip never erases provenance (adv3); hive population only from successful `resolved.Mode`, error paths stay empty — confirmed against operator_model_projection.go:369 (adv4).

### Round 1 (codex, 2026-07-02) — VERDICT: BLOCKERS (2) → both resolved in v0.2.0

- **CFADA1-1 (mode cannot be derived from the supplying layer):** `Policy.Model` et al. may hold a tier name that wins through tier fallback (resolve.go:127), and defaults-only resolution had no vocabulary value. Resolved: derivation moved to the resolution point of the winning token; vocabulary gained `system-default`; capability/strategy/cost fields never set mode; five discriminating test cases mandated (see D2).
- **CFADA1-2 (present-but-invalid mode fell through to inference):** an unrecognized non-empty `selection_mode` could render `Manual · override` (OverrideMode masking) or `Auto/Manual · inferred` via legacy fallbacks. Resolved: projected mode is evaluated first with a three-way absent/valid/invalid split; present-invalid returns `unknown / not projected` immediately with no fallback; projected mode always precedes legacy fields (see D4).
- Advisories adopted: cheapest-with-capabilities tie behavior asserted (strict `<`, first-in-catalog wins — adv1); context_window confirmed display-only, not a resolver gate (adv2); merged-catalog test for hive added to D5 (adv3); hive normal-path projection confirmed Resolve-sourced, error paths return before model set (adv4); old/new decode compatibility confirmed tested at ops_test.go:1729, the risk is semantic not decode (adv5).
