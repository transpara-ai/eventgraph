# Authority System

Three-tier approval for significant actions. Authority ensures that no significant action occurs without appropriate oversight — the CONSENT and AUTHORITY invariants.

## Authority Levels

```
AuthorityLevel {
    Required        // blocks until human approves or rejects
    Recommended     // auto-approves after timeout (default: 15 min)
    Notification    // auto-approves immediately, logged for audit
}
```

| Level | Behaviour | Use Case |
|---|---|---|
| **Required** | Blocks execution until a human actor explicitly approves or rejects. No timeout. No auto-approval. | Irreversible actions, financial transactions, actor suspension, data deletion |
| **Recommended** | Proceeds after `GraphConfig.RecommendedTimeout` (default: 15 min) if no human responds. Human can approve or reject before timeout. | Moderate-impact actions, automated decisions with low confidence, new actor registration |
| **Notification** | Proceeds immediately. Human is notified but not required to act. | Low-impact actions, routine operations, read-only queries, metric emissions |

---

## Authority Data Structures

### Request

```
AuthorityRequest {
    ID:             EventID             // the event ID of this request
    Action:         string              // what the actor wants to do
    Actor:          ActorID             // who is requesting
    Level:          AuthorityLevel      // what level of approval is needed
    Justification:  string              // why this action is being taken
    Causes:         NonEmpty<EventID>   // the events that led to this request
    Context:        map[string]any      // domain-specific context for the approver
    CreatedAt:      time
    ExpiresAt:      Option<time>        // only for Recommended (timeout)
}
```

### Resolution

```
AuthorityResolution {
    RequestID:      EventID             // which request this resolves
    Approved:       bool                // true = approved, false = rejected
    Resolver:       ActorID             // who approved/rejected (must be Human for Required)
    Reason:         Option<string>      // optional explanation
    ResolvedAt:     time
    AutoApproved:   bool                // true if timeout/notification auto-approved
}
```

### Policy

```
AuthorityPolicy {
    Action:         string              // or pattern — "actor.suspend", "trust.*", "*"
    Level:          AuthorityLevel      // what level this action requires
    Approvers:      []ActorID           // who can approve (empty = any human)
    MinTrust:       Option<Score>       // minimum trust to bypass Required → Recommended
    Scope:          Option<DomainScope> // domain-specific policy
}
```

---

## Authority Flow

### 1. Evaluation

When `IDecisionMaker.Decide()` processes an action, it evaluates authority:

```
fn EvaluateAuthority(actor IActor, action string, trust TrustMetrics) → AuthorityLevel:
    policy = FindPolicy(action)         // most specific matching policy

    if policy.MinTrust is Some(threshold):
        if trust.Overall >= threshold:
            // High-trust actors get lower authority requirements
            return DemoteLevel(policy.Level)

    return policy.Level
```

Level demotion based on trust:

| Policy Level | Trust ≥ threshold | Effective Level |
|---|---|---|
| Required | Yes | Recommended |
| Recommended | Yes | Notification |
| Notification | Yes | Notification (no change) |

Trust does not eliminate authority checks — it reduces them. Even high-trust actors are still logged (Notification minimum).

### 2. Request

When authority is needed, an authority request event is emitted:

```
"authority.requested" → AuthorityRequestContent {
    Action:         string
    Actor:          ActorID
    Level:          AuthorityLevel
    Justification:  string
    Causes:         NonEmpty<EventID>
}
```

For **Required** level, execution blocks until resolution. For **Recommended** level, a timer starts. For **Notification** level, execution proceeds immediately.

### 3. Resolution

Authority requests are resolved by:

- **Human approval/rejection** — a human actor emits an `authority.resolved` event
- **Timeout auto-approval** — for Recommended level, after `GraphConfig.RecommendedTimeout`
- **Immediate auto-approval** — for Notification level, immediately after the request

```
"authority.resolved" → AuthorityResolvedContent {
    RequestID:  EventID
    Approved:   bool
    Resolver:   ActorID
    Reason:     Option<string>
}
```

### 4. Action Execution

After approval, the original action proceeds. The decision event references both the authority request and resolution in its `Causes`, creating a complete causal chain:

```
original cause → authority.requested → authority.resolved → action.executed
```

This chain is fully traversable — you can trace any action back to its approval decision and the events that triggered the request.

---

## IAuthorityChain Interface

The strategy interface for authority evaluation. Override for hierarchical orgs, role-based access, or custom delegation.

```
IAuthorityChain {
    Evaluate(context, actor IActor, action string) → Result<AuthorityResult, DecisionError>
    Chain(context, actor IActor, action string) → Result<[]AuthorityLink, DecisionError>
    Grant(context, from IActor, to IActor, scope DomainScope, weight Score) → Result<Edge, StoreError>
    Revoke(context, from IActor, to IActor, scope DomainScope) → Result<Edge, StoreError>
}

AuthorityResult {
    Level:      AuthorityLevel
    Weight:     Score                    // [0.0, 1.0] — authority strength
    Chain:      []AuthorityLink          // delegation chain
    Delegated:  bool                     // true if authority was delegated, not direct
    ExpiresAt:  Option<time>             // delegation expiry
}

AuthorityLink {
    Actor:    ActorID
    Level:    AuthorityLevel
    Weight:   Score                      // [0.0, 1.0]
}
```

---

## Delegation

Authority can be delegated from one actor to another via `Authority` edges:

```
Edge {
    From:      ActorID                   // the delegator
    To:        ActorID                   // the delegate
    Type:      EdgeType.Authority
    Weight:    Weight                    // delegation strength — caps the delegate's authority
    Scope:     Some(DomainScope)         // delegation must be scoped
    ExpiresAt: Option<time>             // delegations can expire
}
```

### Delegation Rules

1. **Scoped** — delegation must specify a `DomainScope`. No blanket authority delegation.
2. **Weighted** — the delegate's authority is capped by the delegation weight. A `Weight(0.5)` delegation means the delegate has at most 50% of the delegator's authority in that scope.
3. **Non-transitive** — A delegates to B, B delegates to C. C does NOT inherit A's authority. Each delegation is a direct relationship.
4. **Revocable** — delegation edges can be superseded by `edge.superseded` events.
5. **Expiring** — delegations can have an `ExpiresAt`. After expiry, the edge is no longer effective.

### Delegation Chain Walk

```
fn EvaluateChain(actor IActor, action string, scope DomainScope) → []AuthorityLink:
    chain = []
    current = actor

    // Direct authority
    directEdges = store.EdgesTo(current.ID(), EdgeType.Authority)
        .filter(e → e.Scope matches scope)

    for edge in directEdges:
        chain.append(AuthorityLink {
            Actor:  edge.From,
            Level:  AuthorityLevel from policy,
            Weight: Score(edge.Weight.Value()),
        })

    // Effective weight is the minimum weight in the chain
    effectiveWeight = chain.map(l → l.Weight).min()

    return chain
```

---

## Trust-Authority Interaction

Trust and authority work together:

| Trust Level | Effect on Authority |
|---|---|
| Score < 0.2 | Authority level escalated (Notification → Recommended → Required) |
| Score 0.2-0.5 | Policy level applied as-is |
| Score 0.5-0.8 | Policy level applied, auto-approval faster |
| Score > 0.8 | Authority level demoted (Required → Recommended → Notification) |

The interaction is configurable via `AuthorityPolicy.MinTrust`. The defaults above are guidelines — product layers set their own thresholds.

---

## Authority Events

All authority activity is recorded:

```
"authority.requested" → AuthorityRequestContent { ... }
"authority.resolved" → AuthorityResolvedContent { ... }
"authority.delegated" → AuthorityDelegatedContent {
    From:   ActorID
    To:     ActorID
    Scope:  DomainScope
    Weight: Score
    ExpiresAt: Option<time>
}
"authority.revoked" → AuthorityRevokedContent {
    From:   ActorID
    To:     ActorID
    Scope:  DomainScope
    Reason: EventID
}
"authority.timeout" → AuthorityTimeoutContent {
    RequestID: EventID
    Level:     AuthorityLevel           // Recommended — the only level that times out
    Duration:  duration
}
```

---

## Default Policies

The system ships with sensible default policies. Product layers override for their domain.

| Action Pattern | Default Level | Rationale |
|---|---|---|
| `actor.suspend` | Required | Suspending an actor is a significant action |
| `actor.memorial` | Required | Memorialising is irreversible |
| `authority.delegated` | Required | Granting authority is itself an authority action |
| `trust.updated` (negative) | Recommended | Trust reductions should be reviewed |
| `trust.updated` (positive) | Notification | Trust increases are routine |
| `edge.created` | Notification | Creating relationships is routine |
| `event.*` (standard) | Notification | Most events are routine |
| `*` (catch-all) | Notification | Default to least friction |

---

## Configuration

From `GraphConfig`:

| Setting | Default | Purpose |
|---|---|---|
| `RecommendedTimeout` | 15 min | Auto-approve timeout for Recommended level |

---

## Invariants

1. **AUTHORITY** — significant actions require approval at the appropriate level. No bypassing.
2. **CONSENT** — no significant action without appropriate approval.
3. **OBSERVABLE** — all authority requests and resolutions are events on the graph.
4. **CAUSALITY** — authority events are causally linked to the events that triggered them and the actions they authorise.
5. **TRANSPARENT** — authority requests include the actor type, so humans know if an AI is requesting authority.

---

## Reference

- `docs/interfaces.md` — `IAuthorityChain`, `AuthorityResult`, `AuthorityLink`, `AuthorityLevel`, `AuthorityPolicy`
- `docs/trust.md` — Trust-authority interaction
- `docs/layers/00-foundation.md` — Layer 0 authority primitives
- `ROADMAP.md` Phase 1 — Implementation status
