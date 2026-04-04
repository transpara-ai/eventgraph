package event

// hiveContent is embedded in all hive content types to satisfy the
// EventContent interface's Accept method. Hive content types use their
// own visitor rather than the base EventContentVisitor.
type hiveContent struct{}

func (hiveContent) Accept(EventContentVisitor) {}

// GapDetectedContent records a detected gap in the hive's composition.
type GapDetectedContent struct {
	hiveContent
	Category    string `json:"category"`
	MissingRole string `json:"missing_role"`
	Evidence    string `json:"evidence"`
	Severity    string `json:"severity"`
}

func (c GapDetectedContent) EventTypeName() string { return "hive.gap.detected" }

// NewGapDetectedContent creates a GapDetectedContent.
func NewGapDetectedContent(category, missingRole, evidence, severity string) GapDetectedContent {
	return GapDetectedContent{
		Category:    category,
		MissingRole: missingRole,
		Evidence:    evidence,
		Severity:    severity,
	}
}

// DirectiveIssuedContent records a directive issued by the CTO agent.
type DirectiveIssuedContent struct {
	hiveContent
	Target   string `json:"target"`
	Action   string `json:"action"`
	Reason   string `json:"reason"`
	Priority string `json:"priority"`
}

func (c DirectiveIssuedContent) EventTypeName() string { return "hive.directive.issued" }

// NewDirectiveIssuedContent creates a DirectiveIssuedContent.
func NewDirectiveIssuedContent(target, action, reason, priority string) DirectiveIssuedContent {
	return DirectiveIssuedContent{
		Target:   target,
		Action:   action,
		Reason:   reason,
		Priority: priority,
	}
}
