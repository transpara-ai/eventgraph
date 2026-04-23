package codegraph

import "github.com/transpara-ai/eventgraph/go/pkg/types"

// CompositionDescriptor names a product pattern and the code graph primitives it wires together.
type CompositionDescriptor struct {
	Name       string
	Purpose    string
	Primitives []types.PrimitiveID
}

var (
	// BoardComposition is the kanban/sprint board pattern.
	BoardComposition = CompositionDescriptor{
		Name:    "Board",
		Purpose: "Kanban/sprint board: Loop(State) → List(Query(Entity)) with Drag(Command(transition))",
		Primitives: ids("CGLayout", "CGList", "CGQuery", "CGDisplay", "CGDrag",
			"CGCommand", "CGEmpty", "CGAction", "CGLoop", "CGState"),
	}

	// DetailComposition is the single entity view with full context.
	DetailComposition = CompositionDescriptor{
		Name:    "Detail",
		Purpose: "Entity detail: properties, thread, audit trail, history, related entities",
		Primitives: ids("CGLayout", "CGDisplay", "CGForm", "CGThread", "CGAudit",
			"CGHistory", "CGList", "CGAction", "CGNavigation"),
	}

	// FeedComposition is the activity stream pattern.
	FeedComposition = CompositionDescriptor{
		Name:    "Feed",
		Purpose: "Activity stream: chronological event list with subscribe for live updates",
		Primitives: ids("CGList", "CGQuery", "CGDisplay", "CGAvatar", "CGSubscribe",
			"CGPagination", "CGRecency"),
	}

	// DashboardComposition is the metrics/aggregate view.
	DashboardComposition = CompositionDescriptor{
		Name:    "Dashboard",
		Purpose: "Metrics dashboard: aggregate queries, counts, charts, KPIs",
		Primitives: ids("CGLayout", "CGDisplay", "CGQuery", "CGTransform", "CGSalience"),
	}

	// InboxComposition is the attention queue pattern.
	InboxComposition = CompositionDescriptor{
		Name:    "Inbox",
		Purpose: "Attention queue: filtered by relevance, prioritized by urgency",
		Primitives: ids("CGList", "CGQuery", "CGDisplay", "CGAvatar", "CGSalience",
			"CGAction", "CGSelection", "CGEmpty"),
	}

	// WizardComposition is the multi-step creation pattern.
	WizardComposition = CompositionDescriptor{
		Name:    "Wizard",
		Purpose: "Multi-step creation: Sequence of Forms with progress and review",
		Primitives: ids("CGSequence", "CGForm", "CGInput", "CGDisplay", "CGAction",
			"CGNavigation", "CGConsequencePreview", "CGConstraint"),
	}

	// SkinComposition is the aesthetic identity pattern.
	SkinComposition = CompositionDescriptor{
		Name:    "Skin",
		Purpose: "Visual identity: Palette + Typography + Spacing + Elevation + Motion + Density + Shape",
		Primitives: ids("CGPalette", "CGTypography", "CGSpacing", "CGElevation",
			"CGMotion", "CGDensity", "CGShape"),
	}
)

// AllCompositions returns all named code graph compositions.
func AllCompositions() []CompositionDescriptor {
	return []CompositionDescriptor{
		BoardComposition,
		DetailComposition,
		FeedComposition,
		DashboardComposition,
		InboxComposition,
		WizardComposition,
		SkinComposition,
	}
}

func ids(names ...string) []types.PrimitiveID {
	out := make([]types.PrimitiveID, len(names))
	for i, n := range names {
		out[i] = types.MustPrimitiveID(n)
	}
	return out
}
