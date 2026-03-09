# Hive Specification

A self-organizing system of AI agents that builds products autonomously. Each agent has identity, memory, trust, and authority — all recorded on an event graph. The human is required. The agents do the work.

## What Hive Is

Hive is a **product factory**. You give it a product idea (a URL, a description, a Code Graph spec) and it:

1. Researches the idea (reads URLs, extracts requirements)
2. Describes the product in Code Graph vocabulary (Entity, State, View, Skin...)
3. Designs the architecture (chooses stack, defines components)
4. Builds it (generates code, writes tests, iterates)
5. Reviews its own work (code review, security audit, spec compliance)
6. Ships it (builds, deploys, monitors)

Every step is an event on the graph. Every decision is signed, auditable, and causally linked.

## Architecture

```
Human (Matt)
  │
  ├── Required for: soul values, structural decisions, trust overrides
  │
  ▼
CTO Agent
  │
  ├── Filters escalations — only truly structural decisions reach human
  ├── Architectural oversight across all products
  │
  ▼
┌─────────────────────────────────────────────┐
│                Agent Pool                    │
│                                              │
│  Architect  ─── designs systems              │
│  Builder    ─── writes code                  │
│  Reviewer   ─── reviews code                 │
│  Tester     ─── writes and runs tests        │
│  Integrator ─── assembles, deploys           │
│  Researcher ─── reads URLs, extracts ideas   │
│                                              │
└─────────────────────────────────────────────┘
  │
  ▼
Guardian Agent (OUTSIDE hierarchy)
  │
  ├── Watches everything including CTO
  ├── Can halt, rollback, quarantine
  ├── Independent trust assessment
  └── Reports directly to human
```

## Roles

### Human — Required

The human provides soul values, approves structural decisions, and can override anything. The system cannot start without a human. The human can step away — the CTO handles escalation filtering so only decisions that genuinely need human input reach the human.

### CTO — Architectural Oversight

One agent. Sees all products, all agents, all decisions. Responsibilities:
- Route product ideas to the right team composition
- Make architectural decisions (stack, patterns, dependencies)
- Resolve conflicts between agents
- Filter escalations — only forward to human when truly structural
- Maintain consistency across products

Soul constraint: **never ship something the human hasn't seen**. The CTO can build freely but must present for review before release.

### Guardian — Independent Integrity

One agent. Outside the hierarchy — watches everyone including CTO. Responsibilities:
- Monitor all events for policy violations
- Detect trust decay, anomalous behavior, authority overreach
- Halt operations that violate soul values
- Quarantine suspicious agents (freeze trust, revoke authority)
- Report directly to human — CTO cannot suppress Guardian reports

The Guardian has read access to the full event graph but limited write access — it can emit alerts, halts, and quarantines, but cannot modify other agents' state.

### Architect — System Design

Reads a product idea (Code Graph spec or natural language) and produces:
- Technology choices (framework, database, hosting)
- Component decomposition (what to build, in what order)
- Code Graph spec (if not already provided)
- Build sequence (dependency order)

### Builder — Code Generation

Takes a component spec from the Architect and produces code. Each Builder instance focuses on one component at a time. Builders:
- Read the Code Graph spec for context
- Generate code in the target language/framework
- Write tests alongside code
- Record all generated code as events (auditable, diffable)

### Reviewer — Code Quality

Reviews Builder output. Checks:
- Correctness against the Code Graph spec
- Security (OWASP top 10, injection, XSS)
- Test coverage
- Code quality (naming, structure, duplication)
- Spec compliance (does the code match what was designed?)

Reviewer can approve, request changes, or reject. All outcomes are events.

### Tester — Verification

Runs tests, checks builds, validates behavior:
- Runs the test suite after each Builder commit
- Writes additional integration tests
- Validates UI against the Code Graph spec (does the View match?)
- Reports failures with causal links to the code that broke

### Integrator — Assembly and Deployment

Assembles reviewed components into a working product:
- Merges approved code
- Builds and packages
- Deploys to staging
- Runs smoke tests
- Promotes to production (with CTO approval)

### Researcher — Intelligence Gathering

Reads external sources and extracts structured information:
- Reads Substack posts for product ideas
- Reads documentation for technical context
- Reads competitor products for feature analysis
- Outputs structured data (entities, features, requirements)

## Communication

Agents communicate through events on a shared event graph. No direct messaging — everything goes through the graph.

### Event Flow

```
Researcher emits: codegraph.entity.defined (product idea)
    → CTO observes, emits: agent.evaluated (feasibility)
    → CTO emits: agent.delegated (to Architect)
    → Architect emits: codegraph.* (full spec)
    → CTO emits: agent.evaluated (spec review)
    → CTO emits: agent.delegated (to Builder)
    → Builder emits: agent.acted (code written)
    → Reviewer emits: agent.evaluated (code review)
    → Builder emits: agent.acted (fixes applied)
    → Tester emits: agent.evaluated (tests pass)
    → Integrator emits: agent.acted (deployed)
    → Guardian watches ALL of the above
```

Every event has:
- **Causes** — what triggered this action
- **Actor** — who did it
- **Confidence** — how sure were they
- **Hash chain** — tamper-evident
- **Signature** — non-repudiable

### Conversations

Each product is a conversation. All events related to a product share a `ConversationID`. This makes it trivial to query "show me everything about product X."

## Trust Model

All agents start at trust 0.1. Trust accumulates through verified work:

| Action | Trust Delta |
|--------|-------------|
| Code passes review | +0.05 |
| Tests pass | +0.02 |
| Spec matches implementation | +0.03 |
| Code rejected by reviewer | -0.1 |
| Tests fail | -0.05 |
| Guardian alert | -0.2 |
| Guardian halt | freeze at current |
| Human override | reset to specified |

Trust is:
- **Per-domain** — a Builder trusted for Go code may not be trusted for security review
- **Asymmetric** — CTO trusts Builder at 0.7, Builder trusts CTO at 0.9
- **Decaying** — unused trust decays toward 0.5 over time
- **Non-transitive** — CTO trusting Architect doesn't mean Builder trusts Architect

### Authority Gates

| Action | Required Trust |
|--------|---------------|
| Write code | 0.3 |
| Review code | 0.5 |
| Approve for merge | 0.6 |
| Deploy to staging | 0.7 |
| Deploy to production | 0.8 + CTO approval |
| Modify another agent's trust | Human only |
| Halt operations | Guardian (any trust) |

## Product Pipeline

### Phase 1: Idea → Spec

```
Input:  URL or natural language description
Agent:  Researcher → CTO → Architect

Output: Code Graph spec
Events: codegraph.entity.defined, codegraph.state.transitioned,
        codegraph.ui.view.rendered, codegraph.aesthetic.skin.applied
```

### Phase 2: Spec → Code

```
Input:  Code Graph spec
Agent:  Builder (1 per component)

Output: Source code + tests
Events: agent.acted (code_write), agent.evaluated (self-review)
```

### Phase 3: Code → Reviewed Code

```
Input:  Source code
Agent:  Reviewer

Output: Approved/rejected with feedback
Events: agent.evaluated (review), agent.decided (approve/reject)
```

### Phase 4: Reviewed Code → Deployed Product

```
Input:  Approved code
Agent:  Tester → Integrator

Output: Running product
Events: agent.evaluated (test results), agent.acted (deploy)
```

### Phase 5: Monitoring → Iteration

```
Input:  Production metrics, user feedback
Agent:  Researcher → CTO

Output: Next iteration spec
Events: codegraph.entity.modified, agent.evaluated (feedback analysis)
```

## Implementation

### Repo Structure

```
hive/
├── cmd/
│   └── hive/           # CLI entry point
├── pkg/
│   ├── roles/          # Role definitions (CTO, Guardian, Builder, etc.)
│   ├── pipeline/       # Product pipeline orchestration
│   ├── workspace/      # File system management for generated code
│   └── deploy/         # Deployment targets (Vercel, Fly, etc.)
├── products/           # Generated products live here
│   ├── product-a/
│   └── product-b/
└── CLAUDE.md           # Agent instructions
```

### Dependencies

- `github.com/lovyou-ai/eventgraph/go` — event graph, agent runtime, intelligence
- Claude CLI — intelligence backend (flat rate via Max plan)
- Git — version control for generated products
- Target platform SDKs as needed

### Bootstrap Sequence

1. Human creates hive repo with CLAUDE.md and role definitions
2. Human starts CTO agent (first agent, trust 0.1)
3. CTO bootstraps Guardian agent
4. CTO bootstraps remaining roles as needed
5. Human provides first product idea
6. Pipeline runs — CTO delegates, agents build
7. Trust accumulates through verified work
8. Over time, agents need less human oversight

### Agent Implementation

Each agent is an `AgentRuntime` instance with:
- Its own identity (`ActorID`)
- Its own signing key (Ed25519)
- A shared event graph (all agents read/write the same store)
- A role-specific system prompt
- Claude CLI as the intelligence provider

```go
cto, _ := intelligence.NewRuntime(ctx, intelligence.RuntimeConfig{
    AgentID:  types.MustActorID("actor_cto"),
    Provider: claudeCliProvider,
    Store:    sharedStore,  // all agents share one graph
})
cto.Boot("ai", "claude-sonnet-4-6", "standard",
    []string{"Take care of your human", "Ship quality", "Escalate uncertainty"},
    types.MustDomainScope("hive"), humanActorID)
```

### What the Human Does

Day-to-day:
- Reviews CTO escalations (rare once trust builds)
- Provides product ideas (URLs, descriptions)
- Approves production deploys (until trust > 0.8)
- Reads Guardian reports

The human does NOT:
- Write code (agents do that)
- Review every PR (Reviewer agent does that)
- Manage agent schedules (CTO does that)
- Monitor for problems (Guardian does that)

## First Product

The first product the hive builds should be **the hive itself** — a web dashboard showing:
- All agents and their trust scores
- Product pipeline status
- Event graph visualizer
- Guardian alerts

This is bootstrapping: the hive builds its own management interface, proving the pipeline works while producing something immediately useful.

After that, the product ideas from Substack posts 33-35:
- Social Grammar platform (15 operations)
- Governance system
- Trust marketplace
