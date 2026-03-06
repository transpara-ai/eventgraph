# Task Management System

> **Product layer** — built ON the event graph, not part of it. This is a stub specification.

## Overview

Hierarchical task decomposition with model-tier routing. Task management is a product layer that uses the event graph for persistence, audit, and causal tracking. All task operations emit events on the graph.

## Task Structure

This is a product-layer type. It uses event graph types where applicable.

```
TaskID              // value object, validated at construction
TaskStatus { Pending, InProgress, Completed, Blocked }  // enum, not bare string
ModelTier  { Small, Medium, Large }                     // enum for routing

Task {
    ID:          TaskID
    Subject:     string             // brief description — genuinely freeform text
    Description: string             // detailed requirements — genuinely freeform text
    Status:      TaskStatus
    Priority:    int                // urgency (higher = more urgent)
    Source:      ActorID            // who created the task
    Assignee:    Option<ActorID>    // who is working on it (None = unassigned)
    ParentID:    Option<TaskID>     // parent task (None = top-level)
    BlockedBy:   []TaskID           // task IDs blocking this task
    ModelTier:   ModelTier          // which model tier handles this
    Metadata:    map[string]any     // flexible key-value data
    CreatedAt:   time
    UpdatedAt:   time
    CompletedAt: Option<time>       // None until completed
}
```

## Hierarchical Decomposition

When the system receives a task, it decomposes it into subtasks:

1. **Planning** — An LLM analyses the task and proposes 1-8 ordered subtasks, each with a model tier (small, medium, large)
2. **Ordering** — Foundational work first (schema, data models), dependent work later (logic, tests)
3. **Execution** — Sequential, with incremental build/test verification
4. **Review** — LLM reviews the complete diff, creates fix subtasks if needed (up to 2 rounds)
5. **Completion** — Task marked complete, changes committed

## Reference

See patent specification Section 8 for the full task management system definition.
