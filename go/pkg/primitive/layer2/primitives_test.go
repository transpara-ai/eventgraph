package layer2_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive"
	"github.com/lovyou-ai/eventgraph/go/pkg/primitive/layer2"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
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
		// Group A: Common Ground
		layer2.NewTermPrimitive(),
		layer2.NewProtocolPrimitive(),
		layer2.NewOfferPrimitive(),
		layer2.NewAcceptancePrimitive(),
		// Group B: Mutual Binding
		layer2.NewAgreementPrimitive(),
		layer2.NewObligationPrimitive(),
		layer2.NewFulfillmentPrimitive(),
		layer2.NewBreachPrimitive(),
		// Group C: Value Transfer
		layer2.NewExchangePrimitive(),
		layer2.NewAccountabilityPrimitive(),
		layer2.NewDebtPrimitive(),
		layer2.NewReciprocityPrimitive(),
	}

	for _, p := range prims {
		if err := reg.Register(p); err != nil {
			t.Errorf("Register %q: %v", p.ID().Value(), err)
		}
		if p.Layer().Value() != 2 {
			t.Errorf("%q: Layer = %d, want 2", p.ID().Value(), p.Layer().Value())
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

func TestTermProcess(t *testing.T) {
	s, bootstrap := bootstrapStore(t)
	ev := chainEvent(t, s, []types.EventID{bootstrap.ID()})
	p := layer2.NewTermPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{ev}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}

func TestBreachProcess(t *testing.T) {
	s, bootstrap := bootstrapStore(t)
	ev := chainEvent(t, s, []types.EventID{bootstrap.ID()})
	p := layer2.NewBreachPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{ev}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}

func TestReciprocityProcess(t *testing.T) {
	s, bootstrap := bootstrapStore(t)
	ev := chainEvent(t, s, []types.EventID{bootstrap.ID()})
	p := layer2.NewReciprocityPrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{ev}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}

func TestAcceptanceProcess(t *testing.T) {
	s, bootstrap := bootstrapStore(t)
	ev := chainEvent(t, s, []types.EventID{bootstrap.ID()})
	p := layer2.NewAcceptancePrimitive()

	mutations, err := p.Process(types.MustTick(1), []event.Event{ev}, primitive.Snapshot{})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(mutations) == 0 {
		t.Fatal("expected mutations")
	}
}
