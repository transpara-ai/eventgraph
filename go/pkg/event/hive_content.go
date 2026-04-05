package event

// hiveContent is embedded in all hive content types to satisfy the
// EventContent interface's Accept method. Hive content types use their
// own visitor rather than the base EventContentVisitor.
type hiveContent struct{}

func (hiveContent) Accept(EventContentVisitor) {}

// GapDetectedContent records a detected gap in the hive's composition.
type GapDetectedContent struct {
	hiveContent
	Category    GapCategory   `json:"Category"`
	MissingRole string        `json:"MissingRole"`
	Evidence    string        `json:"Evidence"`
	Severity    SeverityLevel `json:"Severity"`
}

func (c GapDetectedContent) EventTypeName() string { return "hive.gap.detected" }

// NewGapDetectedContent creates a GapDetectedContent.
func NewGapDetectedContent(category GapCategory, missingRole, evidence string, severity SeverityLevel) GapDetectedContent {
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
	Target   string            `json:"Target"`
	Action   string            `json:"Action"`
	Reason   string            `json:"Reason"`
	Priority DirectivePriority `json:"Priority"`
}

func (c DirectiveIssuedContent) EventTypeName() string { return "hive.directive.issued" }

// NewDirectiveIssuedContent creates a DirectiveIssuedContent.
func NewDirectiveIssuedContent(target, action, reason string, priority DirectivePriority) DirectiveIssuedContent {
	return DirectiveIssuedContent{
		Target:   target,
		Action:   action,
		Reason:   reason,
		Priority: priority,
	}
}
