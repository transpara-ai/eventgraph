# Decision Trees

The mechanical-to-intelligent continuum. Every primitive can have an evolving decision tree that routes decisions through deterministic branches first, falling through to `IIntelligence` only when the tree can't handle it.

## Why Decision Trees

The insight: most decisions that start as expensive AI calls become cheap deterministic rules once patterns emerge. A trust primitive might call Claude the first 100 times it encounters a "code review completed" event. By the 101st time, the pattern is clear: +0.05 trust, every time. That rule should be a branch in a decision tree, not an LLM invocation.

Decision trees are the bridge between mechanical and intelligent processing. They evolve over time — growing deterministic branches as patterns are observed, progressively reducing the proportion of decisions that require expensive model calls.

---

## Tree Structure

```
DecisionTree {
    Root:       DecisionNode
    Version:    int                     // incremented on each evolution step
    Stats:      TreeStats               // tracking for evolution
}

DecisionNode = InternalNode | LeafNode

InternalNode {
    Condition:  Condition               // what to test
    Branches:   NonEmpty<Branch>        // at least one branch (otherwise it's a leaf)
    Default:    DecisionNode            // fallthrough if no branch matches
}

Branch {
    Match:      MatchValue              // what the condition result must equal
    Child:      DecisionNode            // where to go on match
}

LeafNode {
    Outcome:      Option<DecisionOutcome>  // Some = deterministic. None = needs intelligence.
    NeedsLLM:     bool                     // true if this leaf requires IIntelligence
    Confidence:   Score                    // how confident this leaf is [0.0, 1.0]
    Stats:        LeafStats                // tracking for evolution
}
```

### Conditions

Conditions extract a value from the decision context and test it:

```
Condition {
    Field:      string                  // dot-path into DecisionInput context (e.g., "event.type", "actor.trust")
    Operator:   ConditionOperator       // how to compare
    Threshold:  Option<Score>           // for numeric and Semantic comparisons
    Prompt:     Option<string>          // only for Semantic — what to ask IIntelligence
}

ConditionOperator {
    Equals              // exact match
    GreaterThan         // numeric comparison
    LessThan            // numeric comparison
    InRange             // between two values
    Matches             // pattern match (SubscriptionPattern)
    Exists              // field is present (not None)
    Semantic            // delegates to IIntelligence — returns Score, branches on threshold
}

MatchValue {
    String:     Option<string>
    Number:     Option<float64>
    Boolean:    Option<bool>
    EventType:  Option<EventType>
}
```

---

## Evaluation

```
fn Evaluate(tree DecisionTree, input DecisionInput, intelligence Option<IIntelligence>) → Result<Decision, DecisionError>:

    path = []                           // track the path taken for auditing
    node = tree.Root

    while node is InternalNode:
        matched = false

        if node.Condition.Operator == Semantic:
            // Semantic condition — delegates to IIntelligence
            if intelligence is None:
                // Can't evaluate — fall through to default
                path.append(PathStep { condition: node.Condition, branch: "default" })
                node = node.Default
                continue

            response = intelligence.Reason(context, node.Condition.Prompt.unwrap(), input.history)
            value = response.Confidence    // IIntelligence returns a Score

            for branch in node.Branches:
                if value >= node.Condition.Threshold.unwrap():
                    path.append(PathStep { condition: node.Condition, branch: branch.Match })
                    node = branch.Child
                    matched = true
                    break
        else:
            // Mechanical condition — deterministic evaluation
            value = ExtractField(input, node.Condition.Field)

            for branch in node.Branches:
                if TestCondition(value, node.Condition.Operator, branch.Match):
                    path.append(PathStep { condition: node.Condition, branch: branch.Match })
                    node = branch.Child
                    matched = true
                    break

        if not matched:
            path.append(PathStep { condition: node.Condition, branch: "default" })
            node = node.Default

    // Reached a leaf
    leaf = node as LeafNode
    leaf.Stats.HitCount += 1

    if leaf.NeedsLLM:
        if intelligence is None:
            return Err(DecisionError.InsufficientAuthority { ... })

        response = intelligence.Reason(context, FormatPrompt(input, path), input.history)
        leaf.Stats.LastLLMResponse = response
        leaf.Stats.LLMCallCount += 1

        return Ok(Decision {
            Outcome:    ParseOutcome(response.Content),
            Confidence: response.Confidence,
            Evidence:   input.Causes,
            Path:       path,
            ...
        })

    return Ok(Decision {
        Outcome:    leaf.Outcome.unwrap(),
        Confidence: leaf.Confidence,
        Evidence:   input.Causes,
        Path:       path,
        ...
    })
```

### Path Tracking

Every decision records the path taken through the tree. This serves two purposes:

1. **Auditing** — you can trace exactly why a decision was made, which conditions were tested, and which branches were taken
2. **Evolution** — the path data feeds the pattern recognition system

```
PathStep {
    Condition:  Condition
    Branch:     MatchValue              // which branch was taken (or "default")
}
```

---

## Evolution

The key mechanism. Decision trees grow smarter over time by observing patterns in LLM responses and extracting deterministic branches.

### Phase 1: Observation

Every LLM-requiring leaf tracks:

```
LeafStats {
    HitCount:       int                 // total times this leaf was reached
    LLMCallCount:   int                 // times IIntelligence was invoked
    LastLLMResponse: Option<Response>
    ResponseHistory: []ResponseRecord   // last N responses for pattern detection
    PatternScore:    Score              // [0.0, 1.0] — how patterned the responses are
}

ResponseRecord {
    Input:      DecisionInput
    Output:     DecisionOutcome
    Confidence: Score
    Timestamp:  time
}
```

### Phase 2: Pattern Recognition

When a leaf accumulates enough data (`HitCount >= PatternThreshold`, default: 50), the evolution engine analyses response history for extractable patterns:

```
fn DetectPatterns(leaf LeafNode) → []ProposedBranch:

    // Group responses by outcome
    groups = leaf.Stats.ResponseHistory.groupBy(r → r.Output)

    proposedBranches = []

    for (outcome, records) in groups:
        // Find common conditions across records with this outcome
        commonConditions = FindCommonConditions(records)

        if commonConditions is not empty:
            accuracy = TestAccuracy(commonConditions, records)

            if accuracy >= EvolutionThreshold:   // default: 0.95
                proposedBranches.append(ProposedBranch {
                    Conditions:  commonConditions,
                    Outcome:     outcome,
                    Confidence:  accuracy,
                    SampleSize:  len(records),
                })

    return proposedBranches
```

### Phase 3: Branch Insertion

Proposed branches are inserted above the LLM-requiring leaf. The leaf becomes the default fallthrough for inputs that don't match the new deterministic branches.

```
fn InsertBranch(tree DecisionTree, leafPath []PathStep, proposed ProposedBranch) → DecisionTree:

    // Create new internal node
    newNode = InternalNode {
        Condition:  proposed.Conditions[0],
        Branches:   [Branch {
            Match: proposed.MatchValue,
            Child: LeafNode {
                Outcome:    Some(proposed.Outcome),
                NeedsLLM:   false,
                Confidence: proposed.Confidence,
            }
        }],
        Default:    originalLeaf,           // LLM fallthrough preserved
    }

    // Replace the leaf with the new internal node
    return tree.replaceNode(leafPath, newNode)
```

The original LLM leaf is preserved as the default path — new inputs that don't match any extracted pattern still go to intelligence. Over time, more patterns are extracted, and the proportion of LLM calls decreases.

### Phase 4: Cost Demotion

Periodic review of the tree's economics:

```
fn ReviewCosts(tree DecisionTree) → []CostReport:

    reports = []

    for leaf in tree.AllLeaves():
        if leaf.NeedsLLM and leaf.Stats.HitCount > 0:
            llmRate = leaf.Stats.LLMCallCount / leaf.Stats.HitCount
            reports.append(CostReport {
                Path:       leaf.Path,
                HitCount:   leaf.Stats.HitCount,
                LLMRate:    llmRate,            // proportion still going to LLM
                TokenCost:  leaf.Stats.TotalTokens,
                Extractable: leaf.Stats.PatternScore > EvolutionThreshold,
            })

    return reports
```

### Semantic Condition Evolution

`Semantic` conditions follow the same evolution path as LLM-requiring leaves. The pattern recognition system observes the relationship between `Semantic` condition inputs and the resulting `Score`, looking for mechanical correlations:

- "Is this message hostile?" might correlate with `SeverityLevel >= Serious` and the presence of certain event types in the causal chain
- "Does this conflict with stated values?" might correlate with a contradiction score above `0.7`

When a mechanical correlation is found with sufficient accuracy (`>= EvolutionThreshold`), the `Semantic` condition is proposed for replacement with a mechanical equivalent:

```
fn EvolveSemantic(condition Condition, history []SemanticEvalRecord) → Option<Condition>:
    // Find mechanical predictor of the Semantic result
    for candidate in MechanicalCandidates(history):
        accuracy = TestAccuracy(candidate, history)
        if accuracy >= EvolutionThreshold:
            return Some(Condition {
                Field:     candidate.Field,
                Operator:  candidate.Operator,      // mechanical operator
                Threshold: candidate.Threshold,
                Prompt:    None,                     // no longer needs intelligence
            })
    return None     // can't mechanise yet — keep using Semantic

SemanticEvalRecord {
    Input:      DecisionInput
    Score:      Score                   // what IIntelligence returned
    Branch:     MatchValue              // which branch was taken
}
```

This is the continuum in action: a `Semantic` condition that asks "is this proportional?" 500 times, and observes that it always maps to `trust_delta / evidence_count < 0.1`, gets replaced with a `LessThan` condition. The intelligence call disappears. The tree gets cheaper. The behaviour is preserved.

---

## Evolution Events

All evolution steps are recorded as events on the graph:

```
"decision.branch.proposed" → BranchProposedContent {
    PrimitiveID:  PrimitiveID
    TreeVersion:  int
    Condition:    Condition
    Outcome:      DecisionOutcome
    Accuracy:     Score
    SampleSize:   int
}

"decision.branch.inserted" → BranchInsertedContent {
    PrimitiveID:  PrimitiveID
    TreeVersion:  int               // new version after insertion
    Path:         []PathStep
    Outcome:      DecisionOutcome
    Confidence:   Score
}

"decision.cost.report" → CostReportContent {
    PrimitiveID:  PrimitiveID
    TreeVersion:  int
    TotalLeaves:  int
    LLMLeaves:    int
    MechanicalRate: Score           // proportion of decisions handled without LLM
    TotalTokens:  int
}
```

---

## Configuration

From `GraphConfig`:

| Setting | Default | Purpose |
|---|---|---|
| `MaxTokensPerTick` | — | Token budget per tick across all primitives (shared by `Semantic` conditions and LLM leaves) |
| `FallbackToMechanical` | false | If true, never invoke `IIntelligence` — `Semantic` conditions fall through to `Default`, LLM leaves return `Defer` |

Per-primitive configuration (in primitive state):

| Setting | Default | Purpose |
|---|---|---|
| `PatternThreshold` | 50 | Minimum hits before pattern detection runs |
| `EvolutionThreshold` | 0.95 | Minimum accuracy for branch extraction |
| `MaxTreeDepth` | 20 | Prevent unbounded tree growth |
| `ResponseHistorySize` | 200 | How many LLM responses to retain for pattern detection |

---

## Integration with Primitives

Not every primitive needs a decision tree. The tree is an optional component:

- **Purely mechanical primitives** (Hash, Clock, Event) — no tree, no intelligence. Deterministic logic in `Process()`.
- **Hybrid primitives** (TrustScore, Confidence) — start with a tree that has mostly LLM leaves, grow deterministic branches over time.
- **Intelligence-heavy primitives** (higher layers: Ethics, Identity, Wonder) — deep trees with many LLM leaves. These may never fully mechanise — some decisions genuinely require reasoning.

The primitive's `Process()` function calls the tree:

```
fn Process(tick, events, snapshot) → Result<[]Mutation, StoreError>:
    for event in events:
        input = DecisionInput { ... }
        decision = self.tree.Evaluate(input, self.intelligence)
        mutations.append(decisionToMutations(decision))
    return mutations
```

---

## Invariants

1. **Auditable** — every decision path is recorded. You can trace from outcome back through conditions to the original input.
2. **Monotonically improving** — branches are only added, never removed (unless explicitly reset). The tree only gets better.
3. **LLM fallthrough preserved** — extracted branches never replace the LLM path. They intercept known patterns *before* reaching the LLM. Unknown patterns still get intelligent handling.
4. **Evolution is observable** — all tree mutations are events on the graph. The system's learning is visible and auditable.

---

## Reference

- `docs/interfaces.md` — `IDecisionMaker`, `IIntelligence`, `DecisionInput`, `Decision`
- `docs/tick-engine.md` — How primitives are invoked (Process function)
- `docs/primitives.md` — Which primitives use decision trees
- `ROADMAP.md` Phase 1 (Decision Tree Engine) — Implementation status
