# Layer 0: Foundation

## Overview

44 primitives in 11 groups. The irreducible computational foundations — everything above is derived from gaps in this layer.

## Groups

### Group 0 — Core

The event graph itself.

| Primitive | Purpose | Mechanical/Intelligent |
|-----------|---------|----------------------|
| **Event** | The fundamental unit. Creates and validates events. | Mechanical |
| **EventStore** | Persistence. Append, query, verify. | Mechanical |
| **Clock** | Temporal ordering. Tick counting. Timestamps. | Mechanical |
| **Hash** | Cryptographic hashing. Chain integrity. | Mechanical |
| **Self** | System identity and message routing. | Intelligent |

### Group 1 — Causality

Why things happened.

| Primitive | Purpose | Mechanical/Intelligent |
|-----------|---------|----------------------|
| **CausalLink** | Establishing and validating causal edges between events | Mechanical |
| **Ancestry** | Traversing causal chains upward (what caused this?) | Mechanical |
| **Descendancy** | Traversing causal chains downward (what did this cause?) | Mechanical |
| **FirstCause** | Finding root causes. Identifying bootstrap events. | Mostly mechanical, intelligent for complex causal analysis |

### Group 2 — Identity

Who did things.

| Primitive | Purpose | Mechanical/Intelligent |
|-----------|---------|----------------------|
| **ActorID** | Actor identity management. Keypair association. | Mechanical |
| **ActorRegistry** | Actor registration, lookup, lifecycle. | Mechanical |
| **Signature** | Ed25519 signing of events and messages. | Mechanical |
| **Verify** | Signature verification. Authentication. | Mechanical |

### Group 3 — Expectations

What should happen.

| Primitive | Purpose | Mechanical/Intelligent |
|-----------|---------|----------------------|
| **Expectation** | Defining what the system expects to happen after an event | Intelligent |
| **Timeout** | Tracking when expectations expire | Mechanical |
| **Violation** | Detecting when expectations are not met | Intelligent |
| **Severity** | Assessing how serious a violation is | Intelligent |

### Group 4 — Trust

Who to believe.

| Primitive | Purpose | Mechanical/Intelligent |
|-----------|---------|----------------------|
| **TrustScore** | Maintaining trust scores for actors and systems | Both |
| **TrustUpdate** | Processing events that affect trust | Intelligent |
| **Corroboration** | Detecting when multiple sources agree | Intelligent |
| **Contradiction** | Detecting when sources disagree | Intelligent |

### Group 5 — Confidence

How sure we are.

| Primitive | Purpose | Mechanical/Intelligent |
|-----------|---------|----------------------|
| **Confidence** | Maintaining confidence levels for beliefs and states | Both |
| **Evidence** | Tracking evidence for and against propositions | Intelligent |
| **Revision** | Updating beliefs when new evidence arrives | Intelligent |
| **Uncertainty** | Quantifying and communicating what the system doesn't know | Intelligent |

### Group 6 — Instrumentation

What we're watching.

| Primitive | Purpose | Mechanical/Intelligent |
|-----------|---------|----------------------|
| **InstrumentationSpec** | Defining what should be monitored | Intelligent |
| **CoverageCheck** | Verifying that instrumentation covers all critical paths | Both |
| **Gap** | Identifying what's not being monitored | Intelligent |
| **Blind** | Detecting blind spots — things the system can't see | Intelligent |

### Group 7 — Query

How we find things.

| Primitive | Purpose | Mechanical/Intelligent |
|-----------|---------|----------------------|
| **PathQuery** | Finding paths between events or actors in the graph | Mechanical |
| **SubgraphExtract** | Extracting relevant subgraphs for analysis | Both |
| **Annotate** | Adding metadata to existing events | Both |
| **Timeline** | Constructing temporal views of event sequences | Mechanical |

### Group 8 — Integrity

Is the graph intact.

| Primitive | Purpose | Mechanical/Intelligent |
|-----------|---------|----------------------|
| **HashChain** | Managing the linear hash chain | Mechanical |
| **ChainVerify** | Verifying chain integrity | Mechanical |
| **Witness** | External witnessing of chain state (for cross-system trust) | Mechanical |
| **IntegrityViolation** | Detecting and responding to integrity breaches | Both |

### Group 9 — Deception

Is someone lying.

| Primitive | Purpose | Mechanical/Intelligent |
|-----------|---------|----------------------|
| **Pattern** | Recognising behavioural patterns across events | Intelligent |
| **DeceptionIndicator** | Identifying potential deception signals | Intelligent |
| **Suspicion** | Maintaining and updating suspicion levels | Intelligent |
| **Quarantine** | Isolating suspected-compromised events or actors | Both |

### Group 10 — Health

Is the system well.

| Primitive | Purpose | Mechanical/Intelligent |
|-----------|---------|----------------------|
| **GraphHealth** | Overall health assessment of the event graph | Both |
| **Invariant** | Defining system invariants that must hold | Intelligent |
| **InvariantCheck** | Verifying invariants are maintained | Both |
| **Bootstrap** | System initialisation and recovery | Mechanical |

## Layer Activation

Layer 0 primitives activate on system start. They must be stable before any Layer 1 primitive activates. "Stable" means: in Active lifecycle state with consistent outputs for at least one full tick cycle.

## Implementation Notes

- Mechanical primitives should start with mostly deterministic decision tree branches
- Intelligent primitives should start with `NeedsLLM = true` leaves
- All primitives must emit events for their significant actions
- All primitives must declare event subscriptions
- Default cadence for Layer 0: 1 (every tick)
