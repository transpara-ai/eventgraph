namespace EventGraph.CodeGraph;

/// <summary>All 35 Code Graph event type constants. All use the "codegraph." prefix.</summary>
public static class CodeGraphEventTypes
{
    // ── Data events ──────────────────────────────────────────────────────
    public static readonly EventType EntityCreated = new("codegraph.entity.created");
    public static readonly EventType EntityUpdated = new("codegraph.entity.updated");
    public static readonly EventType EntityDeleted = new("codegraph.entity.deleted");
    public static readonly EventType PropertySet = new("codegraph.entity.property.set");
    public static readonly EventType RelationCreated = new("codegraph.entity.relation.created");
    public static readonly EventType StateChanged = new("codegraph.state.changed");

    // ── Logic events ─────────────────────────────────────────────────────
    public static readonly EventType TransformApplied = new("codegraph.logic.transform.applied");
    public static readonly EventType ConditionEvaluated = new("codegraph.logic.condition.evaluated");
    public static readonly EventType SequenceStarted = new("codegraph.logic.sequence.started");
    public static readonly EventType SequenceCompleted = new("codegraph.logic.sequence.completed");
    public static readonly EventType LoopIterated = new("codegraph.logic.loop.iterated");
    public static readonly EventType TriggerFired = new("codegraph.logic.trigger.fired");
    public static readonly EventType ConstraintViolated = new("codegraph.logic.constraint.violated");

    // ── IO events ────────────────────────────────────────────────────────
    public static readonly EventType QueryExecuted = new("codegraph.io.query.executed");
    public static readonly EventType CommandExecuted = new("codegraph.io.command.executed");
    public static readonly EventType SubscriptionCreated = new("codegraph.io.subscribe.created");
    public static readonly EventType AuthorizeChecked = new("codegraph.io.authorize.checked");
    public static readonly EventType SearchExecuted = new("codegraph.io.search.executed");
    public static readonly EventType InteropCalled = new("codegraph.io.interop.called");

    // ── UI events ────────────────────────────────────────────────────────
    public static readonly EventType ViewRendered = new("codegraph.ui.view.rendered");
    public static readonly EventType ActionTriggered = new("codegraph.ui.action.triggered");
    public static readonly EventType NavigationChanged = new("codegraph.ui.navigation.changed");
    public static readonly EventType FeedbackShown = new("codegraph.ui.feedback.shown");
    public static readonly EventType AlertRaised = new("codegraph.ui.alert.raised");
    public static readonly EventType ConfirmationRequested = new("codegraph.ui.confirmation.requested");

    // ── Aesthetic events ─────────────────────────────────────────────────
    public static readonly EventType ThemeApplied = new("codegraph.aesthetic.theme.applied");

    // ── Temporal events ──────────────────────────────────────────────────
    public static readonly EventType UndoRequested = new("codegraph.temporal.undo.requested");
    public static readonly EventType RetryAttempted = new("codegraph.temporal.retry.attempted");

    // ── Resilience events ────────────────────────────────────────────────
    public static readonly EventType FallbackActivated = new("codegraph.resilience.fallback.activated");
    public static readonly EventType OfflineSynced = new("codegraph.resilience.offline.synced");

    // ── Structural events ────────────────────────────────────────────────
    public static readonly EventType ScopeEntered = new("codegraph.structural.scope.entered");
    public static readonly EventType FormatApplied = new("codegraph.structural.format.applied");

    // ── Social events ────────────────────────────────────────────────────
    public static readonly EventType PresenceChanged = new("codegraph.social.presence.changed");
    public static readonly EventType SalienceScored = new("codegraph.social.salience.scored");
    public static readonly EventType ConsequencePreviewed = new("codegraph.ui.confirmation.previewed");

    /// <summary>Returns all 35 registered Code Graph event types.</summary>
    public static List<EventType> AllCodeGraphEventTypes() => new()
    {
        // Data
        EntityCreated, EntityUpdated, EntityDeleted,
        PropertySet, RelationCreated, StateChanged,
        // Logic
        TransformApplied, ConditionEvaluated,
        SequenceStarted, SequenceCompleted,
        LoopIterated, TriggerFired, ConstraintViolated,
        // IO
        QueryExecuted, CommandExecuted,
        SubscriptionCreated, AuthorizeChecked,
        SearchExecuted, InteropCalled,
        // UI
        ViewRendered, ActionTriggered, NavigationChanged,
        FeedbackShown, AlertRaised, ConfirmationRequested,
        // Aesthetic
        ThemeApplied,
        // Temporal
        UndoRequested, RetryAttempted,
        // Resilience
        FallbackActivated, OfflineSynced,
        // Structural
        ScopeEntered, FormatApplied,
        // Social
        PresenceChanged, SalienceScored, ConsequencePreviewed,
    };
}
