namespace EventGraph.CodeGraph;

// All 61 Code Graph primitives operate at Layer 5 (Code Graph), Cadence 1.

// ════════════════════════════════════════════════════════════════════════
// BASE CLASS — Shared implementation for all Code Graph primitives
// ════════════════════════════════════════════════════════════════════════

/// <summary>Base class for all Code Graph primitives. Layer 5, Cadence 1, common Process logic.</summary>
public abstract class CodeGraphPrimitive : IPrimitive
{
    public abstract PrimitiveId Id { get; }
    public Layer Layer { get; } = new(5);
    public abstract List<SubscriptionPattern> Subscriptions { get; }
    public Cadence Cadence { get; } = new(1);

    public List<Mutation> Process(int tick, List<Event> events, Snapshot snapshot)
    {
        return new List<Mutation>
        {
            new UpdateStateMutation(Id, "eventsProcessed", events.Count),
            new UpdateStateMutation(Id, "lastTick", tick),
        };
    }
}

// ════════════════════════════════════════════════════════════════════════
// DATA PRIMITIVES (6) — Core data structures
// ════════════════════════════════════════════════════════════════════════

/// <summary>Entity — a named, typed, identity-bearing node in the code graph.</summary>
public sealed class EntityPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGEntity");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.entity.*"),
        new SubscriptionPattern("codegraph.io.command.*"),
    };
}

/// <summary>Property — a key-value pair attached to an entity.</summary>
public sealed class PropertyPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGProperty");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.entity.*"),
    };
}

/// <summary>Relation — a typed, directed edge between entities.</summary>
public sealed class RelationPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGRelation");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.entity.*"),
    };
}

/// <summary>Collection — an ordered group of entities with query semantics.</summary>
public sealed class CollectionPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGCollection");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.entity.*"),
        new SubscriptionPattern("codegraph.io.query.*"),
    };
}

/// <summary>State — mutable state attached to an entity or scope.</summary>
public sealed class CGStatePrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGState");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.state.*"),
        new SubscriptionPattern("codegraph.io.command.*"),
    };
}

/// <summary>Event — a domain event within the code graph.</summary>
public sealed class CGEventPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGEvent");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.*"),
    };
}

// ════════════════════════════════════════════════════════════════════════
// LOGIC PRIMITIVES (6) — Control flow and computation
// ════════════════════════════════════════════════════════════════════════

/// <summary>Transform — map data from one shape to another.</summary>
public sealed class TransformPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGTransform");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.logic.transform.*"),
    };
}

/// <summary>Condition — boolean predicate evaluation.</summary>
public sealed class ConditionPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGCondition");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.logic.condition.*"),
    };
}

/// <summary>Sequence — ordered execution of steps.</summary>
public sealed class SequencePrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGSequence");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.logic.sequence.*"),
    };
}

/// <summary>Loop — repeated execution over a collection or condition.</summary>
public sealed class LoopPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGLoop");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.logic.loop.*"),
    };
}

/// <summary>Trigger — event-driven activation of logic.</summary>
public sealed class TriggerPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGTrigger");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.logic.trigger.*"),
        new SubscriptionPattern("codegraph.entity.*"),
    };
}

/// <summary>Constraint — validation rules that must hold.</summary>
public sealed class ConstraintPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGConstraint");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.logic.constraint.*"),
        new SubscriptionPattern("codegraph.io.command.*"),
    };
}

// ════════════════════════════════════════════════════════════════════════
// IO PRIMITIVES (6) — Input/output and external interaction
// ════════════════════════════════════════════════════════════════════════

/// <summary>Query — read data from the graph or external sources.</summary>
public sealed class QueryPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGQuery");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.io.query.*"),
    };
}

/// <summary>Command — write data to the graph or external systems.</summary>
public sealed class CommandPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGCommand");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.io.command.*"),
    };
}

/// <summary>Subscribe — reactive event stream subscription.</summary>
public sealed class SubscribePrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGSubscribe");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.io.subscribe.*"),
    };
}

/// <summary>Authorize — permission checks and authority delegation.</summary>
public sealed class AuthorizePrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGAuthorize");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.io.authorize.*"),
        new SubscriptionPattern("authority.*"),
    };
}

/// <summary>Search — full-text or structured search across entities.</summary>
public sealed class SearchPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGSearch");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.io.search.*"),
    };
}

/// <summary>Interop — communication with external systems and APIs.</summary>
public sealed class InteropPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGInterop");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.io.interop.*"),
    };
}

// ════════════════════════════════════════════════════════════════════════
// UI PRIMITIVES (19) — User interface components
// ════════════════════════════════════════════════════════════════════════

/// <summary>Display — render data to the user.</summary>
public sealed class DisplayPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGDisplay");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.*"),
    };
}

/// <summary>Input — capture user input.</summary>
public sealed class InputPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGInput");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.*"),
    };
}

/// <summary>Layout — arrange components spatially.</summary>
public sealed class LayoutPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGLayout");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.*"),
    };
}

/// <summary>List — display a collection of items.</summary>
public sealed class ListPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGList");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.*"),
        new SubscriptionPattern("codegraph.io.query.*"),
    };
}

/// <summary>Form — structured data entry with validation.</summary>
public sealed class FormPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGForm");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.*"),
        new SubscriptionPattern("codegraph.io.command.*"),
    };
}

/// <summary>Action — a user-triggerable operation.</summary>
public sealed class ActionPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGAction");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.action.*"),
    };
}

/// <summary>Navigation — movement between views and contexts.</summary>
public sealed class NavigationPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGNavigation");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.navigation.*"),
    };
}

/// <summary>View — a composed screen or page.</summary>
public sealed class ViewPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGView");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.view.*"),
    };
}

/// <summary>Feedback — user feedback display (success, error, info).</summary>
public sealed class FeedbackPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGFeedback");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.feedback.*"),
    };
}

/// <summary>Alert — urgent notifications requiring attention.</summary>
public sealed class AlertPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGAlert");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.alert.*"),
    };
}

/// <summary>Thread — conversational or temporal sequence of items.</summary>
public sealed class ThreadPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGThread");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.*"),
        new SubscriptionPattern("codegraph.entity.*"),
    };
}

/// <summary>Avatar — visual representation of an actor or entity.</summary>
public sealed class AvatarPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGAvatar");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.*"),
    };
}

/// <summary>Audit — audit trail display for any graph activity.</summary>
public sealed class AuditPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGAudit");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.*"),
    };
}

/// <summary>Drag — drag-and-drop interaction.</summary>
public sealed class DragPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGDrag");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.drag.*"),
    };
}

/// <summary>Selection — multi-select and batch operations.</summary>
public sealed class SelectionPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGSelection");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.selection.*"),
    };
}

/// <summary>Confirmation — user confirmation dialogs for destructive actions.</summary>
public sealed class ConfirmationPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGConfirmation");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.confirmation.*"),
    };
}

/// <summary>Empty — empty state display when no data is available.</summary>
public sealed class EmptyPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGEmpty");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.*"),
    };
}

/// <summary>Loading — loading/skeleton states during async operations.</summary>
public sealed class LoadingPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGLoading");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.*"),
    };
}

/// <summary>Pagination — paginated data display with cursor support.</summary>
public sealed class PaginationPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGPagination");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.*"),
        new SubscriptionPattern("codegraph.io.query.*"),
    };
}

// ════════════════════════════════════════════════════════════════════════
// AESTHETIC PRIMITIVES (7) — Visual design tokens
// ════════════════════════════════════════════════════════════════════════

/// <summary>Palette — colour scheme and theming.</summary>
public sealed class PalettePrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGPalette");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.aesthetic.*"),
    };
}

/// <summary>Typography — font families, sizes, weights.</summary>
public sealed class TypographyPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGTypography");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.aesthetic.*"),
    };
}

/// <summary>Spacing — margins, padding, gaps.</summary>
public sealed class SpacingPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGSpacing");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.aesthetic.*"),
    };
}

/// <summary>Elevation — shadow and depth layers.</summary>
public sealed class ElevationPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGElevation");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.aesthetic.*"),
    };
}

/// <summary>Motion — animation timing and easing.</summary>
public sealed class MotionPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGMotion");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.aesthetic.*"),
    };
}

/// <summary>Density — compact vs comfortable layout density.</summary>
public sealed class DensityPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGDensity");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.aesthetic.*"),
    };
}

/// <summary>Shape — border radius, corner style, component shapes.</summary>
public sealed class ShapePrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGShape");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.aesthetic.*"),
    };
}

// ════════════════════════════════════════════════════════════════════════
// ACCESSIBILITY PRIMITIVES (4) — Inclusive design support
// ════════════════════════════════════════════════════════════════════════

/// <summary>Announce — screen reader announcements and ARIA live regions.</summary>
public sealed class AnnouncePrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGAnnounce");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.*"),
        new SubscriptionPattern("codegraph.aesthetic.*"),
    };
}

/// <summary>Focus — keyboard focus management and tab order.</summary>
public sealed class FocusPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGFocus");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.*"),
    };
}

/// <summary>Contrast — colour contrast ratios and high-contrast modes.</summary>
public sealed class ContrastPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGContrast");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.aesthetic.*"),
    };
}

/// <summary>Simplify — cognitive load reduction and simplified views.</summary>
public sealed class SimplifyPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGSimplify");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.*"),
        new SubscriptionPattern("codegraph.aesthetic.*"),
    };
}

// ════════════════════════════════════════════════════════════════════════
// TEMPORAL PRIMITIVES (3) — Time-aware operations
// ════════════════════════════════════════════════════════════════════════

/// <summary>Recency — time-based sorting and freshness indicators.</summary>
public sealed class RecencyPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGRecency");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.entity.*"),
        new SubscriptionPattern("codegraph.temporal.*"),
    };
}

/// <summary>History — version history and change tracking.</summary>
public sealed class HistoryPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGHistory");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.entity.*"),
        new SubscriptionPattern("codegraph.temporal.*"),
    };
}

/// <summary>Liveness — real-time presence and activity indicators.</summary>
public sealed class LivenessPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGLiveness");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.io.subscribe.*"),
        new SubscriptionPattern("codegraph.social.presence.*"),
    };
}

// ════════════════════════════════════════════════════════════════════════
// RESILIENCE PRIMITIVES (4) — Error recovery and offline support
// ════════════════════════════════════════════════════════════════════════

/// <summary>Undo — reversible action support.</summary>
public sealed class UndoPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGUndo");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.temporal.undo.*"),
        new SubscriptionPattern("codegraph.io.command.*"),
    };
}

/// <summary>Retry — automatic retry with backoff for failed operations.</summary>
public sealed class RetryPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGRetry");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.temporal.retry.*"),
        new SubscriptionPattern("codegraph.io.command.*"),
    };
}

/// <summary>Fallback — graceful degradation when components fail.</summary>
public sealed class FallbackPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGFallback");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.resilience.*"),
        new SubscriptionPattern("codegraph.ui.*"),
    };
}

/// <summary>Offline — offline-first data and sync support.</summary>
public sealed class OfflinePrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGOffline");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.resilience.offline.*"),
    };
}

// ════════════════════════════════════════════════════════════════════════
// STRUCTURAL PRIMITIVES (3) — Organization and formatting
// ════════════════════════════════════════════════════════════════════════

/// <summary>Scope — hierarchical containment and context boundaries.</summary>
public sealed class ScopePrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGScope");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.structural.*"),
    };
}

/// <summary>Format — data formatting and serialization.</summary>
public sealed class FormatPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGFormat");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.io.*"),
        new SubscriptionPattern("codegraph.structural.*"),
    };
}

/// <summary>Gesture — touch and pointer gesture recognition.</summary>
public sealed class GesturePrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGGesture");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.*"),
    };
}

// ════════════════════════════════════════════════════════════════════════
// SOCIAL PRIMITIVES (3) — Social awareness
// ════════════════════════════════════════════════════════════════════════

/// <summary>Presence — online/offline/away status of actors.</summary>
public sealed class PresencePrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGPresence");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.social.presence.*"),
    };
}

/// <summary>Salience — importance ranking and attention prioritization.</summary>
public sealed class SaliencePrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGSalience");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.social.salience.*"),
        new SubscriptionPattern("codegraph.entity.*"),
    };
}

/// <summary>ConsequencePreview — preview the effects of an action before committing.</summary>
public sealed class ConsequencePreviewPrimitive : CodeGraphPrimitive
{
    public override PrimitiveId Id { get; } = new("CGConsequencePreview");
    public override List<SubscriptionPattern> Subscriptions { get; } = new()
    {
        new SubscriptionPattern("codegraph.ui.confirmation.*"),
        new SubscriptionPattern("codegraph.io.command.*"),
    };
}

// ════════════════════════════════════════════════════════════════════════
// REGISTRATION
// ════════════════════════════════════════════════════════════════════════

/// <summary>Factory and registration for all 61 Code Graph primitives.</summary>
public static class CodeGraphPrimitiveFactory
{
    /// <summary>Returns all 61 Code Graph primitives.</summary>
    public static List<IPrimitive> All() => new()
    {
        // Data (6)
        new EntityPrimitive(),
        new PropertyPrimitive(),
        new RelationPrimitive(),
        new CollectionPrimitive(),
        new CGStatePrimitive(),
        new CGEventPrimitive(),
        // Logic (6)
        new TransformPrimitive(),
        new ConditionPrimitive(),
        new SequencePrimitive(),
        new LoopPrimitive(),
        new TriggerPrimitive(),
        new ConstraintPrimitive(),
        // IO (6)
        new QueryPrimitive(),
        new CommandPrimitive(),
        new SubscribePrimitive(),
        new AuthorizePrimitive(),
        new SearchPrimitive(),
        new InteropPrimitive(),
        // UI (19)
        new DisplayPrimitive(),
        new InputPrimitive(),
        new LayoutPrimitive(),
        new ListPrimitive(),
        new FormPrimitive(),
        new ActionPrimitive(),
        new NavigationPrimitive(),
        new ViewPrimitive(),
        new FeedbackPrimitive(),
        new AlertPrimitive(),
        new ThreadPrimitive(),
        new AvatarPrimitive(),
        new AuditPrimitive(),
        new DragPrimitive(),
        new SelectionPrimitive(),
        new ConfirmationPrimitive(),
        new EmptyPrimitive(),
        new LoadingPrimitive(),
        new PaginationPrimitive(),
        // Aesthetic (7)
        new PalettePrimitive(),
        new TypographyPrimitive(),
        new SpacingPrimitive(),
        new ElevationPrimitive(),
        new MotionPrimitive(),
        new DensityPrimitive(),
        new ShapePrimitive(),
        // Accessibility (4)
        new AnnouncePrimitive(),
        new FocusPrimitive(),
        new ContrastPrimitive(),
        new SimplifyPrimitive(),
        // Temporal (3)
        new RecencyPrimitive(),
        new HistoryPrimitive(),
        new LivenessPrimitive(),
        // Resilience (4)
        new UndoPrimitive(),
        new RetryPrimitive(),
        new FallbackPrimitive(),
        new OfflinePrimitive(),
        // Structural (3)
        new ScopePrimitive(),
        new FormatPrimitive(),
        new GesturePrimitive(),
        // Social (3)
        new PresencePrimitive(),
        new SaliencePrimitive(),
        new ConsequencePreviewPrimitive(),
    };

    /// <summary>Registers all 61 Code Graph primitives and activates them.</summary>
    public static void RegisterAll(PrimitiveRegistry registry)
    {
        foreach (var p in All())
        {
            registry.Register(p);
            registry.Activate(p.Id);
        }
    }

    /// <summary>Returns true if the primitive ID belongs to the Code Graph layer.</summary>
    public static bool IsCodeGraph(PrimitiveId id) => id.Value.StartsWith("CG");
}
