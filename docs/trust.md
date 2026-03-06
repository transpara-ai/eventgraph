# Trust Model

Trust is the system's memory of reliability. It is continuous, contextual, evidence-based, and decaying. Trust is not a binary gate — it is a weighted signal that feeds into every decision.

## Properties

| Property | Description |
|---|---|
| **Bounded** | `Score` [0.0, 1.0]. Zero is complete distrust. One is complete trust. Both extremes are theoretical — real trust lives in the middle. |
| **Asymmetric** | A trusting B does not imply B trusting A. Trust is a directed edge. |
| **Non-transitive** | A trusting B and B trusting C does not imply A trusting C. Trust is earned directly, not inherited. |
| **Domain-specific** | Trust in domain X says nothing about domain Y. An actor trusted for `code_review` may not be trusted for `financial`. Domain is a `DomainScope` value object. |
| **Evidence-based** | Trust changes because of observed behaviour, not because someone declared it. Every trust change references the `EventID` that caused it. |
| **Decaying** | Trust decays over time without reinforcement. An actor not seen for months has lower effective trust than one actively contributing. |
| **Append-only** | Trust history is never rewritten. Trust scores are derived from the full event trail. |

---

## Trust Data Structures

From `interfaces.md`:

```
TrustMetrics {
    Actor:       ActorID
    Overall:     Score                    // [0.0, 1.0]
    ByDomain:    map[DomainScope]Score    // per-domain trust
    Confidence:  Score                    // [0.0, 1.0] — how confident we are in this score
    Trend:       Weight                   // [-1.0, 1.0] — direction of recent changes
    Evidence:    []EventID               // events that contributed to current score
    LastUpdated: time
    DecayRate:   Score                    // [0.0, 1.0] — how fast trust decays
}
```

Trust is stored as edges on the graph:

```
Edge {
    From:      ActorID                   // the truster
    To:        ActorID                   // the trusted
    Type:      EdgeType.Trust
    Weight:    Weight                    // current trust level (mapped from Score)
    Scope:     Option<DomainScope>       // None = global trust, Some = domain-specific
    ...
}
```

---

## ITrustModel Interface

The strategy interface. Override for domain-specific trust mechanics.

```
ITrustModel {
    Score(context, actor IActor) → Result<TrustMetrics, StoreError>
    ScoreInDomain(context, actor IActor, domain DomainScope) → Result<TrustMetrics, StoreError>
    Update(context, actor IActor, evidence Event) → Result<TrustMetrics, StoreError>
    Decay(context, actor IActor, elapsed duration) → Result<TrustMetrics, StoreError>
    Between(context, from IActor, to IActor) → Result<TrustMetrics, StoreError>
}
```

### Default Implementation

The default `ITrustModel` uses linear decay and equal evidence weighting:

```
DefaultTrustModel:
    InitialTrust:   Score(0.0)          // all actors start at zero
    DecayRate:      Score(0.01)         // per day
    MaxAdjustment:  Weight(0.1)         // single event can't change trust by more than ±0.1
```

---

## Trust Lifecycle

### 1. Initial State

New actors start at `Score(0.0)`. Trust is earned, never assumed.

### 2. Accumulation

Trust increases through positive evidence — observed behaviour that meets or exceeds expectations:

| Evidence Type | Typical Adjustment |
|---|---|
| Task completed on time | `Weight(+0.05)` |
| Review approved | `Weight(+0.03)` |
| Commitment honoured | `Weight(+0.04)` |
| Corroboration of a claim | `Weight(+0.02)` |

Adjustments are constrained by `MaxAdjustment` — a single event can't cause a trust swing larger than the configured maximum.

### 3. Degradation

Trust decreases through negative evidence — violations, contradictions, missed commitments:

| Evidence Type | Typical Adjustment |
|---|---|
| Deadline missed | `Weight(-0.10)` |
| Expectation violated | `Weight(-0.15)` |
| Contradiction detected | `Weight(-0.20)` |
| Deception indicator | `Weight(-0.30)` |

Negative adjustments are typically larger than positive ones — trust is easier to lose than to gain. This is by design: it models real-world trust dynamics.

### 4. Decay

Trust decays over time without reinforcement. The decay function:

```
fn Decay(score Score, elapsed duration, rate Score) → Score:
    days = elapsed.TotalDays()
    decayed = score.Value() * (1.0 - rate.Value()) ^ days
    return Score(max(0.0, decayed))
```

An actor with `Score(0.8)` and `DecayRate(0.01)`:
- After 1 day: `0.792`
- After 7 days: `0.745`
- After 30 days: `0.592`
- After 90 days: `0.326`
- After 365 days: `0.021`

Decay ensures the system doesn't blindly trust actors who have been absent. Active participants maintain trust through ongoing positive evidence.

### 5. Recovery

After trust degradation, recovery is slower than initial accumulation. The further trust has fallen, the harder it is to recover:

```
fn AdjustedGain(currentScore Score, rawGain Weight) → Weight:
    // Recovery factor: slower recovery at lower trust
    if rawGain.Value() > 0 and currentScore.Value() < 0.3:
        recoveryFactor = currentScore.Value() / 0.3     // 0.0 to 1.0
        return Weight(rawGain.Value() * recoveryFactor)
    return rawGain
```

An actor at `Score(0.1)` gains trust at 33% the rate of an actor at `Score(0.5)`. This prevents rapid trust recovery after significant violations.

---

## Domain-Specific Trust

Trust is scoped by `DomainScope`. An actor's overall trust is the weighted average of their domain-specific scores, but domain scores are tracked independently.

```
fn OverallTrust(actor IActor, domains map[DomainScope]Score) → Score:
    if domains is empty:
        return Score(0.0)
    total = sum(domains.values().map(s → s.Value()))
    return Score(total / len(domains))
```

Domain trust is independent — a violation in `financial` doesn't affect `code_review`. The Layer 0 TrustScore primitive maintains the per-domain bookkeeping.

### Custom Trust Mechanics

Override `ITrustModel` for domain-specific behaviour:

```
type ProjectTrust struct { trust.DefaultModel }

var (
    TaskCompleted  = MustEventType("task.completed")
    DeadlineMissed = MustEventType("deadline.missed")
    ReviewApproved = MustEventType("review.approved")
)

func (t *ProjectTrust) Update(ctx, actor, evidence) Result<TrustMetrics, StoreError> {
    switch evidence.Type {
    case TaskCompleted:
        return t.Adjust(actor, Weight(+0.05))
    case DeadlineMissed:
        return t.Adjust(actor, Weight(-0.10))
    case ReviewApproved:
        return t.Adjust(actor, Weight(+0.03))
    }
    return t.DefaultModel.Update(ctx, actor, evidence)
}
```

---

## Trust Events

All trust changes are recorded as events:

```
"trust.updated" → TrustUpdatedContent {
    Actor:    ActorID
    Previous: Score
    Current:  Score
    Domain:   DomainScope
    Cause:    EventID               // the evidence that triggered this change
}

"trust.score" → TrustScoreContent {
    Actor:   ActorID
    Metrics: TrustMetrics           // full snapshot at this point
}

"trust.decayed" → TrustDecayedContent {
    Actor:    ActorID
    Previous: Score
    Current:  Score
    Elapsed:  duration
    Rate:     Score
}
```

---

## Trust in Decisions

Trust feeds into every decision via the `Decision` structure:

```
Decision {
    ...
    TrustWeights:   []TrustWeight        // trust scores of relevant actors
    ...
}

TrustWeight {
    Actor:    ActorID
    Score:    Score
    Domain:   DomainScope
}
```

The `IDecisionMaker` considers trust when evaluating actions:
- Low trust → higher authority requirement (Recommended → Required)
- High trust → lower authority requirement (Required → Recommended)
- Below threshold → action denied with `DecisionError.TrustBelowThreshold`

---

## Inter-System Trust (EGIP)

Trust between systems follows the same model, applied across graph boundaries:

- Systems start at `Score(0.0)` — unknown systems are untrusted
- Trust accumulates through treaty compliance and consistent behaviour
- Trust decays faster for systems (higher default `DecayRate`) — machines can change rapidly
- Trust is verified through EGIP PROOF messages — integrity proofs that demonstrate a system's chain is valid
- Treaty `TrustRequired` field sets the minimum trust for bilateral operations

```
"egip.trust.updated" → EGIPTrustUpdatedContent {
    System:   SystemURI
    Previous: Score
    Current:  Score
    Evidence: EnvelopeID            // the EGIP interaction that triggered the change
}
```

---

## Corroboration and Contradiction

Layer 0 primitives handle trust evidence aggregation:

- **Corroboration** — when multiple independent sources confirm the same claim, trust in the claim (and the sources) increases. The Corroboration primitive detects convergent evidence.
- **Contradiction** — when sources disagree, trust in at least one source must decrease. The Contradiction primitive detects divergent claims and emits events for the trust system to process.

These feed into trust updates:
```
corroboration detected → trust.updated (+0.02 per corroborating source)
contradiction detected → trust.updated (-0.10 for the contradicting source, pending resolution)
```

---

## Invariants

1. **Evidence-based** — every trust change references the event that caused it. No trust change without a `Cause` event.
2. **Append-only** — trust history is never modified. Current trust is derived from the full trail.
3. **Bounded** — trust never exceeds [0.0, 1.0]. Constrained by `Score` type.
4. **Asymmetric** — trust edges are directional. A→B is independent of B→A.
5. **Observable** — every trust change emits a `trust.updated` event. No silent trust changes.

---

## Reference

- `docs/interfaces.md` — `ITrustModel`, `TrustMetrics`, `TrustWeight`, `Score`, `DomainScope`
- `docs/layers/00-foundation.md` — Layer 0 Group 4 (Trust) and Group 5 (Confidence) primitives
- `docs/primitives.md` — TrustScore, TrustUpdate, Corroboration, Contradiction primitives
- `docs/protocol.md` — Inter-system trust via EGIP
