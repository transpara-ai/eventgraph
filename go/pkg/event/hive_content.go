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

// RoleProposedContent records a role proposal submitted for governance review.
type RoleProposedContent struct {
	hiveContent
	Name          string   `json:"name"`
	Model         string   `json:"model"`
	WatchPatterns []string `json:"watch_patterns"`
	CanOperate    bool     `json:"can_operate"`
	MaxIterations int      `json:"max_iterations"`
	Prompt        string   `json:"prompt"`
	Reason        string   `json:"reason"`
	ProposedBy    string   `json:"proposed_by"`
}

func (c RoleProposedContent) EventTypeName() string { return "hive.role.proposed" }

// RoleApprovedContent records the approval of a proposed role.
type RoleApprovedContent struct {
	hiveContent
	Name       string `json:"name"`
	ApprovedBy string `json:"approved_by"`
	Reason     string `json:"reason"`
}

func (c RoleApprovedContent) EventTypeName() string { return "hive.role.approved" }

// RoleRejectedContent records the rejection of a proposed role.
type RoleRejectedContent struct {
	hiveContent
	Name       string `json:"name"`
	RejectedBy string `json:"rejected_by"`
	Reason     string `json:"reason"`
}

func (c RoleRejectedContent) EventTypeName() string { return "hive.role.rejected" }
