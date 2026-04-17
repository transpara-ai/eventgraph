package event

// GapCategory represents the domain category of a detected gap.
type GapCategory string

const (
	GapCategoryLeadership  GapCategory = "Leadership"
	GapCategoryTechnical   GapCategory = "Technical"
	GapCategoryProcess     GapCategory = "Process"
	GapCategoryStaffing    GapCategory = "Staffing"
	GapCategoryCapability  GapCategory = "Capability"
)

var validGapCategories = map[GapCategory]bool{
	GapCategoryLeadership: true, GapCategoryTechnical: true,
	GapCategoryProcess: true, GapCategoryStaffing: true,
	GapCategoryCapability: true,
}

// IsValid returns true if the gap category is a known category.
func (c GapCategory) IsValid() bool { return validGapCategories[c] }

// DirectivePriority represents the urgency of a directive.
type DirectivePriority string

const (
	DirectivePriorityCritical DirectivePriority = "Critical"
	DirectivePriorityHigh     DirectivePriority = "High"
	DirectivePriorityMedium   DirectivePriority = "Medium"
	DirectivePriorityLow      DirectivePriority = "Low"
)

var validDirectivePriorities = map[DirectivePriority]bool{
	DirectivePriorityCritical: true, DirectivePriorityHigh: true,
	DirectivePriorityMedium: true, DirectivePriorityLow: true,
}

// IsValid returns true if the directive priority is a known priority.
func (p DirectivePriority) IsValid() bool { return validDirectivePriorities[p] }

// SpecOutcome represents the terminal state of a parsed spec.
type SpecOutcome string

const (
	SpecOutcomeSuccess SpecOutcome = "success"
	SpecOutcomePartial SpecOutcome = "partial"
	SpecOutcomeFailed  SpecOutcome = "failed"
)

var validSpecOutcomes = map[SpecOutcome]bool{
	SpecOutcomeSuccess: true,
	SpecOutcomePartial: true,
	SpecOutcomeFailed:  true,
}

// IsValid returns true if the spec outcome is a known outcome.
func (o SpecOutcome) IsValid() bool { return validSpecOutcomes[o] }
