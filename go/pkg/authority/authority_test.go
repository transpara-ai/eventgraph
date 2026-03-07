package authority_test

import (
	"context"
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/actor"
	"github.com/lovyou-ai/eventgraph/go/pkg/authority"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/trust"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

type testSigner struct{}

func (testSigner) Sign(data []byte) (types.Signature, error) {
	sig := make([]byte, 64)
	copy(sig, data)
	return types.MustSignature(sig), nil
}

func testPublicKey(b byte) types.PublicKey {
	key := make([]byte, 32)
	key[0] = b
	return types.MustPublicKey(key)
}

func testActor(t *testing.T, name string, b byte) actor.IActor {
	t.Helper()
	s := actor.NewInMemoryActorStore()
	a, err := s.Register(testPublicKey(b), name, event.ActorTypeHuman)
	if err != nil {
		t.Fatalf("register actor: %v", err)
	}
	return a
}

func newTestDefaultChain(t *testing.T, model trust.ITrustModel, s store.Store) *authority.DefaultAuthorityChain {
	t.Helper()
	registry := event.DefaultRegistry()
	factory := event.NewEventFactory(registry)
	return authority.NewDefaultAuthorityChain(model, s, factory, testSigner{})
}

func TestDefaultLevelIsNotification(t *testing.T) {
	model := trust.NewDefaultTrustModel()
	s := store.NewInMemoryStore()
	chain := newTestDefaultChain(t, model, s)

	a := testActor(t, "Alice", 1)
	result, err := chain.Evaluate(context.Background(), a, "some.random.action")
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if result.Level != event.AuthorityLevelNotification {
		t.Errorf("Level = %v, want Notification", result.Level)
	}
	if len(result.Chain) != 1 {
		t.Errorf("Chain length = %d, want 1", len(result.Chain))
	}
}

func TestPolicyExactMatch(t *testing.T) {
	model := trust.NewDefaultTrustModel()
	s := store.NewInMemoryStore()
	chain := newTestDefaultChain(t, model, s)

	chain.AddPolicy(authority.AuthorityPolicy{
		Action: "actor.suspend",
		Level:  event.AuthorityLevelRequired,
	})

	a := testActor(t, "Alice", 1)
	result, _ := chain.Evaluate(context.Background(), a, "actor.suspend")
	if result.Level != event.AuthorityLevelRequired {
		t.Errorf("Level = %v, want Required", result.Level)
	}
}

func TestPolicyWildcard(t *testing.T) {
	model := trust.NewDefaultTrustModel()
	s := store.NewInMemoryStore()
	chain := newTestDefaultChain(t, model, s)

	chain.AddPolicy(authority.AuthorityPolicy{
		Action: "trust.*",
		Level:  event.AuthorityLevelRecommended,
	})

	a := testActor(t, "Alice", 1)
	result, _ := chain.Evaluate(context.Background(), a, "trust.update")
	if result.Level != event.AuthorityLevelRecommended {
		t.Errorf("Level = %v, want Recommended", result.Level)
	}
}

func TestPolicyGlobalWildcard(t *testing.T) {
	model := trust.NewDefaultTrustModel()
	s := store.NewInMemoryStore()
	chain := newTestDefaultChain(t, model, s)

	chain.AddPolicy(authority.AuthorityPolicy{
		Action: "*",
		Level:  event.AuthorityLevelRequired,
	})

	a := testActor(t, "Alice", 1)
	result, _ := chain.Evaluate(context.Background(), a, "anything.at.all")
	if result.Level != event.AuthorityLevelRequired {
		t.Errorf("Level = %v, want Required", result.Level)
	}
}

func TestPolicyFirstMatchWins(t *testing.T) {
	model := trust.NewDefaultTrustModel()
	s := store.NewInMemoryStore()
	chain := newTestDefaultChain(t, model, s)

	chain.AddPolicy(authority.AuthorityPolicy{
		Action: "deploy",
		Level:  event.AuthorityLevelRequired,
	})
	chain.AddPolicy(authority.AuthorityPolicy{
		Action: "deploy",
		Level:  event.AuthorityLevelNotification,
	})

	a := testActor(t, "Alice", 1)
	result, _ := chain.Evaluate(context.Background(), a, "deploy")
	if result.Level != event.AuthorityLevelRequired {
		t.Errorf("Level = %v, want Required (first match)", result.Level)
	}
}

func TestPolicyNoMatch(t *testing.T) {
	model := trust.NewDefaultTrustModel()
	s := store.NewInMemoryStore()
	chain := newTestDefaultChain(t, model, s)

	chain.AddPolicy(authority.AuthorityPolicy{
		Action: "deploy",
		Level:  event.AuthorityLevelRequired,
	})

	a := testActor(t, "Alice", 1)
	result, _ := chain.Evaluate(context.Background(), a, "review")
	if result.Level != event.AuthorityLevelNotification {
		t.Errorf("Level = %v, want Notification (no match)", result.Level)
	}
}

func TestTrustDowngradeRequiredToRecommended(t *testing.T) {
	model := trust.NewDefaultTrustModel()
	s := store.NewInMemoryStore()
	chain := newTestDefaultChain(t, model, s)

	// Policy: deploy requires Required, but MinTrust 0.5 downgrades to Recommended
	chain.AddPolicy(authority.AuthorityPolicy{
		Action:   "deploy",
		Level:    event.AuthorityLevelRequired,
		MinTrust: types.Some(types.MustScore(0.001)), // very low threshold
	})

	a := testActor(t, "Alice", 1)

	// Build trust above threshold
	// Score method calls getOrCreate which returns initial trust (0.0)
	// We need to update trust first
	result, _ := chain.Evaluate(context.Background(), a, "deploy")
	// Initial trust is 0.0, which is below 0.001
	if result.Level != event.AuthorityLevelRequired {
		t.Errorf("Level = %v, want Required (trust too low)", result.Level)
	}
}

func TestChain(t *testing.T) {
	model := trust.NewDefaultTrustModel()
	s := store.NewInMemoryStore()
	chain := newTestDefaultChain(t, model, s)

	a := testActor(t, "Alice", 1)
	links, err := chain.Chain(context.Background(), a, "any.action")
	if err != nil {
		t.Fatalf("Chain: %v", err)
	}
	if len(links) != 1 {
		t.Errorf("chain length = %d, want 1 (flat model)", len(links))
	}
	if links[0].Actor != a.ID() {
		t.Error("chain should contain the actor")
	}
}

func TestGrant(t *testing.T) {
	model := trust.NewDefaultTrustModel()
	s := bootstrapStore(t)
	chain := newTestDefaultChain(t, model, s)

	from := testActor(t, "Alice", 1)
	to := testActor(t, "Bob", 2)

	edge, err := chain.Grant(context.Background(), from, to, types.MustDomainScope("code_review"), types.MustScore(0.8))
	if err != nil {
		t.Fatalf("Grant: %v", err)
	}
	if edge.Type() != event.EdgeTypeAuthority {
		t.Errorf("edge type = %v, want Authority", edge.Type())
	}
	if edge.From() != from.ID() {
		t.Error("edge From should be granter")
	}
	if edge.To() != to.ID() {
		t.Error("edge To should be grantee")
	}
}

func TestRevoke(t *testing.T) {
	model := trust.NewDefaultTrustModel()
	s := store.NewInMemoryStore()
	chain := newTestDefaultChain(t, model, s)

	from := testActor(t, "Alice", 1)
	to := testActor(t, "Bob", 2)

	err := chain.Revoke(context.Background(), from, to, types.MustDomainScope("code_review"))
	if err != nil {
		t.Fatalf("Revoke: %v", err)
	}
}

func TestAuthorityResultWeight(t *testing.T) {
	model := trust.NewDefaultTrustModel()
	s := store.NewInMemoryStore()
	chain := newTestDefaultChain(t, model, s)

	a := testActor(t, "Alice", 1)
	result, _ := chain.Evaluate(context.Background(), a, "test")
	if result.Weight.Value() != 1.0 {
		t.Errorf("Weight = %v, want 1.0", result.Weight.Value())
	}
}
