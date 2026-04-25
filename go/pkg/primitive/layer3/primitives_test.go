package layer3_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/primitive"
	"github.com/transpara-ai/eventgraph/go/pkg/primitive/layer3"
	"github.com/transpara-ai/eventgraph/go/pkg/store"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

var (
	systemActor = types.MustActorID("actor_00000000000000000000000000000001")
	actor2      = types.MustActorID("actor_00000000000000000000000000000002")
	convID      = types.MustConversationID("conv_00000000000000000000000000000001")
)

type testSigner struct{}

func (testSigner) Sign(data []byte) (types.Signature, error) {
	sig := make([]byte, 64)
	copy(sig, data)
	return types.MustSignature(sig), nil
}

type headFromStore struct{ s store.Store }

func (h headFromStore) Head() (types.Option[event.Event], error) { return h.s.Head() }

func bootstrapStore(t *testing.T) (store.Store, event.Event) {
	t.Helper()
	s := store.NewInMemoryStore()
	factory := event.NewBootstrapFactory(event.DefaultRegistry())
	ev, err := factory.Init(systemActor, testSigner{})
	if err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	if _, err := s.Append(ev); err != nil {
		t.Fatalf("append: %v", err)
	}
	return s, ev
}

func chainEvent(t *testing.T, s store.Store, causes []types.EventID) event.Event {
	t.Helper()
	factory := event.NewEventFactory(event.DefaultRegistry())
	ev, err := factory.Create(
		event.EventTypeTrustUpdated, systemActor,
		event.TrustUpdatedContent{
			Actor: actor2, Previous: types.MustScore(0.5),
			Current: types.MustScore(0.6), Domain: types.MustDomainScope("test"),
			Cause: causes[0],
		},
		causes, convID, headFromStore{s}, testSigner{},
	)
	if err != nil {
		t.Fatalf("create event: %v", err)
	}
	if _, err := s.Append(ev); err != nil {
		t.Fatalf("append: %v", err)
	}
	return ev
}

func TestAllPrimitivesRegister(t *testing.T) {
	reg := primitive.NewRegistry()

	prims := []primitive.Primitive{
		// Group A: Collective Identity
		layer3.NewGroupPrimitive(),
		layer3.NewMembershipPrimitive(),
		layer3.NewRolePrimitive(),
		layer3.NewConsentPrimitive(),
		// Group B: Social Order
		layer3.NewNormPrimitive(),
		layer3.NewReputationPrimitive(),
		layer3.NewSanctionPrimitive(),
		layer3.NewAuthorityPrimitive(),
		// Group C: Collective Agency
		layer3.NewPropertyPrimitive(),
		layer3.NewCommonsPrimitive(),
		layer3.NewGovernancePrimitive(),
		layer3.NewCollectiveActPrimitive(),
	}

	for _, p := range prims {
		if err := reg.Register(p); err != nil {
			t.Errorf("Register %q: %v", p.ID().Value(), err)
		}
		if p.Layer().Value() != 3 {
			t.Errorf("%q: Layer = %d, want 3", p.ID().Value(), p.Layer().Value())
		}
		if p.Lifecycle() != types.LifecycleActive {
			t.Errorf("%q: Lifecycle = %v, want Active", p.ID().Value(), p.Lifecycle())
		}
		if len(p.Subscriptions()) == 0 {
			t.Errorf("%q: no subscriptions", p.ID().Value())
		}
	}

	if reg.Count() != 12 {
		t.Errorf("registered %d primitives, want 12", reg.Count())
	}
}

func TestGroupProcess(t *testing.T) {
	s, bootstrap := bootstrapStore(t)
	ev := chainEvent(t, s, []types.EventID{bootstrap.ID()})
	p := layer3.NewGroupPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{ev}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}

func TestSanctionProcess(t *testing.T) {
	s, bootstrap := bootstrapStore(t)
	ev := chainEvent(t, s, []types.EventID{bootstrap.ID()})
	p := layer3.NewSanctionPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{ev}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}

func TestCollectiveActProcess(t *testing.T) {
	s, bootstrap := bootstrapStore(t)
	ev := chainEvent(t, s, []types.EventID{bootstrap.ID()})
	p := layer3.NewCollectiveActPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{ev}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}

func TestReputationProcess(t *testing.T) {
	s, bootstrap := bootstrapStore(t)
	ev := chainEvent(t, s, []types.EventID{bootstrap.ID()})
	p := layer3.NewReputationPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{ev}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}
