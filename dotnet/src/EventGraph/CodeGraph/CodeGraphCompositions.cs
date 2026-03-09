namespace EventGraph.CodeGraph;

/// <summary>A named sequence of Code Graph primitive operations.</summary>
public sealed record CodeGraphComposition(string Name, List<string> Primitives, List<EventType> Events);

/// <summary>7 named compositions built from the 61 Code Graph primitives.</summary>
public static class CodeGraphCompositions
{
    /// <summary>Kanban/project board with drag-and-drop columns.
    /// Layout + List + Query + Display + Drag + Command + Empty + Action + Loop + State</summary>
    public static CodeGraphComposition Board() => new(
        "Board",
        new() { "CGLayout", "CGList", "CGQuery", "CGDisplay", "CGDrag", "CGCommand", "CGEmpty", "CGAction", "CGLoop", "CGState" },
        new()
        {
            CodeGraphEventTypes.ViewRendered,
            CodeGraphEventTypes.QueryExecuted,
            CodeGraphEventTypes.CommandExecuted,
            CodeGraphEventTypes.ActionTriggered,
        }
    );

    /// <summary>Detail view for a single entity with history and activity.
    /// Layout + Display + Form + Thread + Audit + History + List + Action + Navigation</summary>
    public static CodeGraphComposition Detail() => new(
        "Detail",
        new() { "CGLayout", "CGDisplay", "CGForm", "CGThread", "CGAudit", "CGHistory", "CGList", "CGAction", "CGNavigation" },
        new()
        {
            CodeGraphEventTypes.ViewRendered,
            CodeGraphEventTypes.CommandExecuted,
            CodeGraphEventTypes.ActionTriggered,
            CodeGraphEventTypes.NavigationChanged,
        }
    );

    /// <summary>Scrollable feed of items with real-time updates.
    /// List + Query + Display + Avatar + Subscribe + Pagination + Recency</summary>
    public static CodeGraphComposition Feed() => new(
        "Feed",
        new() { "CGList", "CGQuery", "CGDisplay", "CGAvatar", "CGSubscribe", "CGPagination", "CGRecency" },
        new()
        {
            CodeGraphEventTypes.QueryExecuted,
            CodeGraphEventTypes.SubscriptionCreated,
            CodeGraphEventTypes.ViewRendered,
        }
    );

    /// <summary>Summary dashboard with key metrics and visualizations.
    /// Layout + Display + Query + Transform + Salience</summary>
    public static CodeGraphComposition Dashboard() => new(
        "Dashboard",
        new() { "CGLayout", "CGDisplay", "CGQuery", "CGTransform", "CGSalience" },
        new()
        {
            CodeGraphEventTypes.ViewRendered,
            CodeGraphEventTypes.QueryExecuted,
            CodeGraphEventTypes.TransformApplied,
            CodeGraphEventTypes.SalienceScored,
        }
    );

    /// <summary>Inbox/notification centre with priority sorting.
    /// List + Query + Display + Avatar + Salience + Action + Selection + Empty</summary>
    public static CodeGraphComposition Inbox() => new(
        "Inbox",
        new() { "CGList", "CGQuery", "CGDisplay", "CGAvatar", "CGSalience", "CGAction", "CGSelection", "CGEmpty" },
        new()
        {
            CodeGraphEventTypes.QueryExecuted,
            CodeGraphEventTypes.ActionTriggered,
            CodeGraphEventTypes.SalienceScored,
            CodeGraphEventTypes.ViewRendered,
        }
    );

    /// <summary>Multi-step wizard with validation and consequence preview.
    /// Sequence + Form + Input + Display + Action + Navigation + ConsequencePreview + Constraint</summary>
    public static CodeGraphComposition Wizard() => new(
        "Wizard",
        new() { "CGSequence", "CGForm", "CGInput", "CGDisplay", "CGAction", "CGNavigation", "CGConsequencePreview", "CGConstraint" },
        new()
        {
            CodeGraphEventTypes.SequenceStarted,
            CodeGraphEventTypes.SequenceCompleted,
            CodeGraphEventTypes.CommandExecuted,
            CodeGraphEventTypes.ConstraintViolated,
            CodeGraphEventTypes.ConsequencePreviewed,
            CodeGraphEventTypes.NavigationChanged,
        }
    );

    /// <summary>Visual theming skin with all aesthetic tokens.
    /// Palette + Typography + Spacing + Elevation + Motion + Density + Shape</summary>
    public static CodeGraphComposition Skin() => new(
        "Skin",
        new() { "CGPalette", "CGTypography", "CGSpacing", "CGElevation", "CGMotion", "CGDensity", "CGShape" },
        new()
        {
            CodeGraphEventTypes.ThemeApplied,
        }
    );

    /// <summary>Returns all 7 named compositions.</summary>
    public static List<CodeGraphComposition> All() => new()
    {
        Board(),
        Detail(),
        Feed(),
        Dashboard(),
        Inbox(),
        Wizard(),
        Skin(),
    };
}
