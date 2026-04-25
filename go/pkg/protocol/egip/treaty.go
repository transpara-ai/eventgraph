package egip

import (
	"fmt"
	"time"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

// Treaty represents a bilateral governance agreement between two systems.
type Treaty struct {
	ID        types.TreatyID
	SystemA   types.SystemURI
	SystemB   types.SystemURI
	Status    event.TreatyStatus
	Terms     []TreatyTerm
	CreatedAt time.Time
	UpdatedAt time.Time
}

// validTreatyTransitions defines which status transitions are valid.
var validTreatyTransitions = map[event.TreatyStatus][]event.TreatyStatus{
	event.TreatyStatusProposed:   {event.TreatyStatusActive, event.TreatyStatusTerminated},
	event.TreatyStatusActive:     {event.TreatyStatusSuspended, event.TreatyStatusTerminated},
	event.TreatyStatusSuspended:  {event.TreatyStatusActive, event.TreatyStatusTerminated},
	event.TreatyStatusTerminated: {}, // terminal state
}

// Transition attempts to move the treaty to a new status.
// Returns an error if the transition is invalid.
func (t *Treaty) Transition(to event.TreatyStatus) error {
	allowed := validTreatyTransitions[t.Status]
	for _, s := range allowed {
		if s == to {
			t.Status = to
			t.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("invalid treaty transition: %s → %s", t.Status, to)
}

// ApplyAction applies a treaty action and returns the resulting status.
func (t *Treaty) ApplyAction(action event.TreatyAction) error {
	switch action {
	case event.TreatyActionAccept:
		return t.Transition(event.TreatyStatusActive)
	case event.TreatyActionSuspend:
		return t.Transition(event.TreatyStatusSuspended)
	case event.TreatyActionTerminate:
		return t.Transition(event.TreatyStatusTerminated)
	case event.TreatyActionModify:
		// Modify doesn't change status — terms are updated separately.
		if t.Status != event.TreatyStatusActive {
			return fmt.Errorf("can only modify active treaties, current status: %s", t.Status)
		}
		t.UpdatedAt = time.Now()
		return nil
	case event.TreatyActionPropose:
		return fmt.Errorf("cannot apply Propose to existing treaty")
	default:
		return fmt.Errorf("unknown treaty action: %s", action)
	}
}

// NewTreaty creates a new treaty in Proposed status.
func NewTreaty(id types.TreatyID, systemA, systemB types.SystemURI, terms []TreatyTerm) *Treaty {
	now := time.Now()
	return &Treaty{
		ID:        id,
		SystemA:   systemA,
		SystemB:   systemB,
		Status:    event.TreatyStatusProposed,
		Terms:     terms,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
