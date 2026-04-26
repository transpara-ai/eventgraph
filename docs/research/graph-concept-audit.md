# Graph Concept Audit: Ethics Graph, Identity Graph, Community Graph

> **Date:** 2026-04-03
> **Scope:** All transpara-ai repositories
> **Purpose:** Locate every reference to three named product graphs across the codebase

---

## Executive Summary

All three concepts are **fully codified named entities** — not aspirational terms, not buried comments, but formally defined product graphs with layer specifications, composition grammars, multi-language implementations, integration tests, and extensive documentation.

| Concept | Canonical Name | Layer | Grammar Type | Primitives | Status |
|---------|---------------|-------|-------------|------------|--------|
| **Ethics Graph** | Layer 7: Ethics | Alignment | `AlignmentGrammar` | 12 | Spec'd + implemented in eventgraph |
| **Identity Graph** | Layer 8: Identity | Identity | `IdentityGrammar` | 12 | Spec'd + implemented in eventgraph |
| **Community Graph** | Layer 10: Community | Belonging | `BelongingGrammar` | 12 | Spec'd + implemented in eventgraph |

### Reference Density by Repository

| Repository | Ethics Graph | Identity Graph | Community Graph | Total Refs |
|-----------|:----------:|:-------------:|:--------------:|:----------:|
| **eventgraph** | 8+ files | 10+ files | 8+ files | ~26 |
| **site** | ~17 files | ~13 files | ~12 files | ~42 |
| **hive** | 3+ files | 3+ files | 3+ files | ~9 |
| **agent** | Indirect | Indirect | Indirect | ~6 |
| **work** | None | None | None | 0 |

---

## Part I: transpara-ai-eventgraph

The authoritative source. All three graphs are defined here as product layers with full specifications, primitives, grammars, and implementations.

### Ethics Graph (Layer 7: Alignment)

> *"Is becomes ought. Values, constraints, harm, accountability."*

#### Canonical Definition

**`docs/layers/07-ethics.md`** line 108:
```
Layer 7 maps to the **Ethics Graph** — AI alignment with transparent moral reasoning.
```

**`docs/product-layers.md`** lines 140-150:
Full product graph specification documenting key event flows: decision audit, harm detection, value constraints, and accountability chains.

#### Primitives

**`go/pkg/primitive/layer7/primitives.go`** — 12 primitives in 3 groups:

| Group | Primitives |
|-------|-----------|
| Moral Standing | `MoralStatus`, `Dignity`, `Autonomy`, `Flourishing` |
| Moral Obligation | `Duty`, `Harm`, `Care`, `Justice` |
| Moral Agency | `Conscience`, `Virtue`, `Responsibility`, `Motive` |

#### Composition Grammar

**`docs/compositions/07-alignment.md`** — Alignment Grammar specification:

| Operations (10) | Named Functions (5) |
|-----------------|-------------------|
| Constrain, Detect-Harm, Assess-Fairness, Flag-Dilemma, Weigh, Explain, Assign, Repair, Care, Grow | Ethics-Audit, Guardrail, Restorative-Justice, Impact-Assessment, Whistleblow |

#### Implementations

| Language | File |
|----------|------|
| Go | `go/pkg/compositions/alignment.go` |
| Rust | `rust/src/compositions/alignment.rs` |
| TypeScript | `ts/src/compositions.ts` |
| Python | `python/eventgraph/compositions.py` |

#### Integration Test

**`go/pkg/integration/scenario10_ai_ethics_audit_test.go`**
Exercises: fairness audit &rarr; harm detection &rarr; authority escalation &rarr; responsibility assignment &rarr; redress

---

### Identity Graph (Layer 8: Identity)

> *"Doing becomes being. Self-knowledge, continuity, authenticity, expression."*

#### Canonical Definition

**`docs/layers/08-identity.md`** line 108:
```
Layer 8 maps to the **Identity Graph** — self-sovereign identity with narrative and purpose.
```

**`docs/product-layers.md`** lines 151-170:
Product specification documenting key event flows: introspection, narrative construction, alignment verification, self-disclosure. Identity emerges from verifiable action history, not self-reported claims.

#### Primitives

**`go/pkg/primitive/layer8/primitives.go`** — 12 primitives in 3 groups:

| Group | Primitives |
|-------|-----------|
| Self-Knowledge | `Narrative`, `SelfConcept`, `Reflection`, `Memory` |
| Self-Direction | `Purpose`, `Aspiration`, `Authenticity`, `Expression` |
| Self-Becoming | `Growth`, `Continuity`, `Integration`, `Crisis` |

#### Composition Grammar

**`docs/compositions/08-identity.md`** — Identity Grammar specification:

| Operations (10) | Named Functions (5) |
|-----------------|-------------------|
| Introspect, Narrate, Align, Bound, Aspire, Transform, Disclose, Recognize, Distinguish, Memorialize | Credential, Identity-Audit, Retirement, Reinvention, Introduction |

#### Implementations

| Language | File |
|----------|------|
| Go | `go/pkg/compositions/identity.go` |
| Rust | `rust/src/compositions/identity.rs` |
| TypeScript | `ts/src/compositions.ts` |
| Python | `python/eventgraph/compositions.py` |

#### Integration Test

**`go/pkg/integration/scenario11_agent_identity_lifecycle_test.go`**
Exercises: self-model &rarr; authenticity check &rarr; aspiration &rarr; boundary &rarr; work summary &rarr; transformation &rarr; narrative &rarr; dignity

---

### Community Graph (Layer 10: Belonging)

> *"Relationship becomes belonging. Shared meaning, living practice, communal experience."*

#### Canonical Definition

**`docs/layers/10-community.md`** line 108:
```
Layer 10 maps to the **Community Graph** — shared meaning, living practice, and belonging.
```

**`docs/product-layers.md`** lines 180+:
Product specification documenting key event flows: onboarding, traditions, stewardship, succession, milestone celebrations. Belonging as gradient, not binary member/non-member.

#### Primitives

**`go/pkg/primitive/layer10/primitives.go`** — 12 primitives in 3 groups:

| Group | Primitives |
|-------|-----------|
| Shared Meaning | `Culture`, `SharedNarrative`, `Ethos`, `Sacred` |
| Living Practice | `Tradition`, `Ritual`, `Practice`, `Place` |
| Communal Experience | `Belonging`, `Solidarity`, `Voice`, `Welcome` |

#### Composition Grammar

**`docs/compositions/10-belonging.md`** — Belonging Grammar specification:

| Operations (10) | Named Functions (5) |
|-----------------|-------------------|
| Settle, Contribute, Include, Practice, Steward, Sustain, Pass-On, Celebrate, Tell, Gift | Onboard, Festival, Succession, Commons-Governance, Renewal |

#### Implementations

| Language | File |
|----------|------|
| Go | `go/pkg/compositions/belonging.go` |
| Rust | `rust/src/compositions/belonging.rs` |
| TypeScript | `ts/src/compositions.ts` |
| Python | `python/eventgraph/compositions.py` |

#### Integration Test

**`go/pkg/integration/scenario12_community_lifecycle_test.go`**
Exercises: invitation &rarr; settling &rarr; contributions &rarr; tradition participation &rarr; stewardship transfer &rarr; celebration

---

### Additional Integration Scenarios

Beyond the primary three, related scenarios exist:

| Scenario | File | Relationship |
|----------|------|-------------|
| Scenario 4: Community Governance | `scenario04_*_test.go` | Layer 10 primitives in governance context |
| Scenario 16: Community Evolution | `scenario16_*_test.go` | Long-term community dynamics |
| Scenario 17: Agent Lifecycle | `scenario17_*_test.go` | Identity lifecycle patterns |
| Scenario 18: Whistleblower Recall | `scenario18_*_test.go` | Layer 7 ethics in action |

### Roadmap Status

**`ROADMAP.md`** explicitly marks all three as complete:
- Layer 7 — Ethics (Is &rarr; Ought) — **DONE**
- Layer 8 — Identity (Doing &rarr; Being) — **DONE**
- Layer 10 — Community (Relationship &rarr; Belonging) — **DONE**

---

## Part II: transpara-ai-site

The public-facing documentation layer. Contains the richest narrative context for all three concepts.

### Dedicated Reference Documentation

| Graph | Layer Spec | Grammar Spec | Fundamentals |
|-------|-----------|-------------|-------------|
| Ethics | `content/reference/layers/07-ethics.md` | `content/reference/grammars/07-alignment.md` | `content/reference/fundamentals/layer-7-ethics.md` |
| Identity | `content/reference/layers/08-identity.md` | `content/reference/grammars/08-identity.md` | `content/reference/fundamentals/layer-8-identity.md` |
| Community | `content/reference/layers/10-community.md` | `content/reference/grammars/10-belonging.md` | `content/reference/fundamentals/layer-10-community.md` |

**Master Specification:** `content/reference/product-layers.md` — all 13 product graphs defined, including explicit sections for Ethics (line 140), Identity (line 162), and Community (line 209).

### Blog Post References

#### Ethics Graph

| Post | File | Context |
|------|------|---------|
| Post 11 | `post11-thirteen-graphs-one-infrastructure.md` | Framework overview, Ethics Graph as cross-cutting layer |
| Post 16 | `post16-the-social-graph.md` | Ethics Graph monitors Social Graph for harm patterns |
| Post 19 | `post19-the-knowledge-graph.md` | Teaser: "Next deep dive: the Ethics Graph" |
| Post 21 | `post21-the-governance-graph.md` | Ethics Graph in AI governance — constitution on the chain |
| Post 24 | `post24-the-map-complete.md` | Listed as "Post 20 -- The Ethics Graph (Layer 7: Ethics)" |
| Post 27 | `post27-the-weight.md` | "The Ethics Graph doesn't exist. Patterns repeat: 'harm identified, evidence present, accountability absent.'" |
| Post 28 | `post28-the-transition.md` | "The Ethics Graph begins as a monitoring layer across the lower graphs." |
| Post 31 | `post31-what-you-could-build.md` | Market manipulation flagged by Ethics Graph |
| Post 41 | `post41-the-hive.md` | Referenced in identity context |
| Post 42 | `post42-flesh-is-weak.md` | Layer 7 Ethics primitives in personal narrative |

#### Identity Graph

| Post | File | Context |
|------|------|---------|
| Post 11 | `post11-thirteen-graphs-one-infrastructure.md` | Identity as hash-chained record, behaviour-first identity |
| Post 17 | `post17-the-justice-graph.md` | Dispute resolution generates Identity Graph events |
| Post 19 | `post19-the-knowledge-graph.md` | Claims linked to "their Identity Graph" |
| Post 20 | `post20-the-relationship-graph.md` | "Previous: The Identity Graph (Layer 8 deep dive)" |
| Post 24 | `post24-the-map-complete.md` | Listed as "Post 21 -- The Identity Graph (Layer 8: Identity)" |
| Post 28 | `post28-the-transition.md` | "Your Identity Graph is derived from your Work Graph activity, your Market Graph transactions, your Social Graph participation, your Community Graph contributions." |
| Post 31 | `post31-what-you-could-build.md` | Portable reputation on Identity Graph across platforms |
| Post 41 | `post41-the-hive.md` | "Identity that emerges from verifiable action history, not self-reported claims." |

#### Community Graph

| Post | File | Context |
|------|------|---------|
| Post 11 | `post11-thirteen-graphs-one-infrastructure.md` | Community Graph in freelancer scenario |
| Post 20 | `post20-the-relationship-graph.md` | "Next deep dive: the Community Graph" |
| Post 21 | `post21-the-governance-graph.md` | Ten layers of infrastructure including community |
| Post 24 | `post24-the-map-complete.md` | Listed as "Post 23 -- The Community Graph (Layer 10: Population)" |
| Post 27 | `post27-the-weight.md` | "The Community Graph doesn't exist, so belonging is binary — you're in or you're out." |
| Post 28 | `post28-the-transition.md` | Identity derived from "Community Graph contributions" |

---

## Part III: transpara-ai-hive

The operational layer. Contains layer specifications and governance planning, but no shipped implementations of these three graphs.

### Layer Specifications

**`loop/layers-general-spec.md`** defines all three as product layers with concrete entities:

| Layer | Domain | Entities | Key Relationships |
|-------|--------|----------|-------------------|
| 7: Alignment | Transparency & accountability | `analytics`, `audit`, `disclosure` | Makes all other layers visible |
| 8: Identity | Selfhood & credentials | `role`, `credential`, `reputation`, `badge` | Earned through Work, endorsed via Social, governed by Policy |
| 10: Belonging | Group membership & lifecycle | `team`, `department`, `organization`, `invitation` | Determines access to Work, governed via Policy, earned through Identity |

### Governance Context

**`loop/council-full-civilization-2026-03-25.md`** line 1513:
> "12 of 13 product layers remain unbuilt. We shipped Work (Layer 1) and fragments of Social (Layer 3). The Market Graph spec is sitting in `loop/market-graph-spec.md` unimplemented. Governance, Justice, Knowledge, **Alignment**, **Identity** — all spec'd, none built."

### Related Documents

| File | Relevance |
|------|-----------|
| `docs/VISION.md` | Strategic vision grouping layers into five aspects of collective existence |
| `docs/ROLES.md` | Agent role architecture |
| `loop/unified-spec.md` | 18 entity kinds across all layers, unified grammar |
| `loop/product-map.md` | Product ecosystem across layers |
| Agent roles: `philosopher.md`, `advocate.md`, `harmony.md` | Ethics-adjacent agent definitions |

### Implementation Status in Hive

All three are **spec'd but not built**. Only Layer 1 (Work) is shipped; fragments of Layer 3 (Social) exist.

---

## Part IV: transpara-ai-agent

The agent runtime. Contains no explicit graph references but surfaces the concepts indirectly through agent capabilities.

### Ethics Graph Surface

| File | Line | Manifestation |
|------|------|--------------|
| `agent.go` | 21 | `SoulValues: []string{"Take care of your human..."}` |
| `agent.go` | 79-80 | `SoulValues []string` field — core ethical constraints |
| `agent.go` | 92 | Boot sequence: "(identity, soul, model, authority, state &rarr; Idle)" |
| `operations.go` | 191 | `Refuse()` — "soul-protected refusal" |

### Identity Graph Surface

| File | Line | Manifestation |
|------|------|--------------|
| `agent.go` | 104-111 | Deterministic Ed25519 signing key from agent name |
| `agent.go` | 120 | Actor store registration with `ActorTypeAI` |
| `agent.go` | 184 | "withIdentity=false, Spawner handles identity" |

### Community Graph Surface

| File | Line | Manifestation |
|------|------|--------------|
| `operations.go` | 128-145 | `Communicate()` — agent-to-agent messaging via named channels |
| `operations.go` | 129 | Channel-based communication ("general", "alerts") |
| `retire.go` | 70-76 | Farewell communication on "lifecycle" channel |
| `agent.go` | 82-85 | `ConversationID` enables shared conversation threads |

### Architectural Note

The agent package is deliberately **role-agnostic** and **identity-agnostic**. It provides the execution envelope; graph-level concerns are delegated:
- Ethics &rarr; SoulValues (emitted at boot)
- Identity &rarr; Spawner (manages identity assignment)
- Community &rarr; Shared graph membership + channels

---

## Part V: transpara-ai-work

**No matches found.** Zero references to Ethics Graph, Identity Graph, or Community Graph in any form.

This repository is strictly Layer 1 (Work Graph) — task management only. The sole contextual reference is in `events.go` line 25:

```go
// Work Graph event types — Layer 1 of the thirteen-product roadmap.
```

---

## Part VI: Naming Conventions

An important distinction exists between **product names** and **code names**:

| Product Name | Code Type Name | Layer Name |
|-------------|---------------|------------|
| Ethics Graph | `AlignmentGrammar` | Layer 7: Ethics |
| Identity Graph | `IdentityGrammar` | Layer 8: Identity |
| Community Graph | `BelongingGrammar` | Layer 10: Belonging |

The term "X Graph" is the **product-facing name** used in documentation and blog posts. In code, the grammar types use the **layer's compositional name** (Alignment, Identity, Belonging). No Go/Rust struct named `EthicsGraph`, `IdentityGraph`, or `CommunityGraph` exists — these are conceptual projections of the unified event graph, not separate data structures.

---

## Part VII: Cross-Graph Relationships

The three graphs are not isolated. Documentation explicitly describes their interdependencies:

```
                    ┌─────────────────────┐
                    │   Ethics Graph (7)  │
                    │   Monitors all      │
                    │   layers for harm   │
                    └────────┬────────────┘
                             │ constrains
                             ▼
┌─────────────────┐   ┌─────────────────┐   ┌─────────────────────┐
│  Work Graph (1) │──▶│ Identity Graph   │◀──│ Community Graph (10) │
│  Actions earn   │   │      (8)         │   │ Belonging earned     │
│  reputation     │   │ Emerges from     │   │ through identity     │
└─────────────────┘   │ verifiable       │   └──────────────────────┘
                      │ action history   │
                      └──────────────────┘
```

Key relationships from site blog posts:
- **Identity derives from all graphs**: "Your Identity Graph is derived from your Work Graph activity, your Market Graph transactions, your Social Graph participation, your Community Graph contributions." (Post 28)
- **Ethics monitors all graphs**: "The Ethics Graph begins as a monitoring layer across the lower graphs. Pattern detection for harm." (Post 28)
- **Community requires Identity**: Membership earned through identity, governed via policy (hive layer spec)

---

## Appendix: Files Not Found

The following governance documents were searched for but do not exist in any repository:

| Document | Searched In |
|----------|------------|
| `RIGHTS.md` | All 5 repos |
| `SOUL.md` | All 5 repos |
| `ROLES.md` | Only exists in hive (`docs/ROLES.md`) |

No files use the terms "named views," "projections," or "sub-graphs" to describe the relationship between product graphs and the underlying event graph. The architecture uses "product graphs" or "product layers" consistently.

---

*Generated by automated codebase audit across 5 repositories.*
