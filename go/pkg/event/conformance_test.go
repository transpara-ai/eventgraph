package event

import (
	"fmt"
	"testing"
	"time"

	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// Conformance tests validate canonical form and hash computation.
// These produce the reference values that other language implementations
// must match exactly.

func mustTime(s string) types.Timestamp {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return types.NewTimestamp(t)
}

func TestConformanceBootstrapCanonicalForm(t *testing.T) {
	ev := NewBootstrapEvent(
		1,
		types.MustEventID("019462a0-0000-7000-8000-000000000001"),
		EventTypeSystemBootstrapped,
		mustTime("2023-11-14T22:13:20Z"),
		types.MustActorID("actor_00000000000000000000000000000001"),
		BootstrapContent{
			ActorID:      types.MustActorID("actor_00000000000000000000000000000001"),
			ChainGenesis: types.ZeroHash(),
			Timestamp:    mustTime("2023-11-14T22:13:20Z"),
		},
		types.MustConversationID("conv_00000000000000000000000000000001"),
		types.ZeroHash(),
		types.MustSignature(make([]byte, 64)),
	)

	canonical := CanonicalForm(ev)
	expected := `1|||019462a0-0000-7000-8000-000000000001|system.bootstrapped|actor_00000000000000000000000000000001|conv_00000000000000000000000000000001|1700000000000000000|{"ActorID":"actor_00000000000000000000000000000001","ChainGenesis":"0000000000000000000000000000000000000000000000000000000000000000","Timestamp":"2023-11-14T22:13:20Z"}`
	if canonical != expected {
		t.Errorf("canonical form mismatch:\n  got:  %s\n  want: %s", canonical, expected)
	}

	hash, err := ComputeHash(canonical)
	if err != nil {
		t.Fatalf("ComputeHash: %v", err)
	}
	expectedHash := "f7cae7ae11c1232a932c64f2302432c0e304dffce80f3935e688980dfbafeb75"
	if hash.Value() != expectedHash {
		t.Errorf("hash mismatch:\n  got:  %s\n  want: %s", hash.Value(), expectedHash)
	}
}

func TestConformanceTrustUpdatedCanonicalForm(t *testing.T) {
	prevHash := types.MustHash("a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2")
	ev := NewEvent(
		1,
		types.MustEventID("019462a0-0000-7000-8000-000000000002"),
		EventTypeTrustUpdated,
		mustTime("2023-11-14T22:13:21Z"),
		types.MustActorID("actor_00000000000000000000000000000001"),
		TrustUpdatedContent{
			Actor:    types.MustActorID("actor_00000000000000000000000000000002"),
			Cause:    types.MustEventID("019462a0-0000-7000-8000-000000000001"),
			Current:  types.MustScore(0.85),
			Domain:   types.MustDomainScope("code_review"),
			Previous: types.MustScore(0.8),
		},
		[]types.EventID{types.MustEventID("019462a0-0000-7000-8000-000000000001")},
		types.MustConversationID("conv_00000000000000000000000000000001"),
		types.ZeroHash(),
		prevHash,
		types.MustSignature(make([]byte, 64)),
	)

	canonical := CanonicalForm(ev)

	// Verify key ordering: Actor < Cause < Current < Domain < Previous
	hash, err := ComputeHash(canonical)
	if err != nil {
		t.Fatalf("ComputeHash: %v", err)
	}
	expectedHash := "b2fbcd2684868f0b0d07d2f5136b52f14b8e749da7b4b7bae2a22f67147152b7"
	if hash.Value() != expectedHash {
		t.Errorf("hash mismatch:\n  got:  %s\n  want: %s", hash.Value(), expectedHash)
	}
}

func TestConformanceEdgeCreatedKeyOrdering(t *testing.T) {
	prevHash := types.MustHash("b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3")
	ev := NewEvent(
		1,
		types.MustEventID("019462a0-0000-7000-8000-000000000003"),
		EventTypeEdgeCreated,
		mustTime("2023-11-14T22:13:22Z"),
		types.MustActorID("actor_00000000000000000000000000000001"),
		EdgeCreatedContent{
			From:      types.MustActorID("actor_00000000000000000000000000000001"),
			To:        types.MustActorID("actor_00000000000000000000000000000002"),
			EdgeType:  EdgeTypeTrust,
			Weight:    types.MustWeight(0.5),
			Direction: EdgeDirectionCentripetal,
		},
		[]types.EventID{types.MustEventID("019462a0-0000-7000-8000-000000000001")},
		types.MustConversationID("conv_00000000000000000000000000000001"),
		types.ZeroHash(),
		prevHash,
		types.MustSignature(make([]byte, 64)),
	)

	canonical := CanonicalForm(ev)

	// Verify keys are sorted: Direction < EdgeType < From < To < Weight
	hash, err := ComputeHash(canonical)
	if err != nil {
		t.Fatalf("ComputeHash: %v", err)
	}
	expectedHash := "4e5c6710ca9325676663b4a66d2e82114fcd8fb49dbe5705795051e0b0be374c"
	if hash.Value() != expectedHash {
		t.Errorf("hash mismatch:\n  got:  %s\n  want: %s", hash.Value(), expectedHash)
	}
}

func TestConformanceHashDeterminism(t *testing.T) {
	canonical := "1||test|system.bootstrapped|actor_test|conv_test|1000|{}"
	h1, _ := ComputeHash(canonical)
	h2, _ := ComputeHash(canonical)
	if h1 != h2 {
		t.Error("hash computation is not deterministic")
	}
}

func TestConformanceHashLength(t *testing.T) {
	canonical := "test"
	hash, err := ComputeHash(canonical)
	if err != nil {
		t.Fatalf("ComputeHash: %v", err)
	}
	if len(hash.Value()) != 64 {
		t.Errorf("hash length = %d, want 64", len(hash.Value()))
	}
}

func TestConformanceBootstrapZeroPrevHash(t *testing.T) {
	ev := NewBootstrapEvent(
		1,
		types.MustEventID("019462a0-0000-7000-8000-000000000001"),
		EventTypeSystemBootstrapped,
		mustTime("2023-11-14T22:13:20Z"),
		types.MustActorID("actor_00000000000000000000000000000001"),
		BootstrapContent{},
		types.MustConversationID("conv_00000000000000000000000000000001"),
		types.ZeroHash(),
		types.MustSignature(make([]byte, 64)),
	)

	if ev.PrevHash() != types.ZeroHash() {
		t.Error("bootstrap should have zero prev_hash")
	}
	if !ev.IsBootstrap() {
		t.Error("should be identified as bootstrap")
	}
}

func TestConformanceCanonicalFormBootstrapEmptyPrevHash(t *testing.T) {
	// Bootstrap canonical form has empty prev_hash field (between the two pipes)
	ev := NewBootstrapEvent(
		1,
		types.MustEventID("019462a0-0000-7000-8000-000000000001"),
		EventTypeSystemBootstrapped,
		mustTime("2023-11-14T22:13:20Z"),
		types.MustActorID("actor_00000000000000000000000000000001"),
		BootstrapContent{},
		types.MustConversationID("conv_00000000000000000000000000000001"),
		types.ZeroHash(),
		types.MustSignature(make([]byte, 64)),
	)

	canonical := CanonicalForm(ev)
	// Should start with "1|||" (empty prev_hash and empty causes between pipes)
	if canonical[:4] != "1|||" {
		t.Errorf("bootstrap canonical should start with '1|||', got: %s", canonical[:10])
	}
}

// --- Mock implementations for factory tests ---

// mockSigner implements Signer for testing.
type mockSigner struct{}

func (m *mockSigner) Sign(data []byte) (types.Signature, error) {
	return types.MustSignature(make([]byte, 64)), nil
}

// mockSignerError implements Signer that always returns an error.
type mockSignerError struct{}

func (m *mockSignerError) Sign(data []byte) (types.Signature, error) {
	return types.Signature{}, fmt.Errorf("signing failed")
}

// mockHeadProviderEmpty implements HeadProvider with no head (empty chain).
type mockHeadProviderEmpty struct{}

func (m *mockHeadProviderEmpty) Head() (types.Option[Event], error) {
	return types.None[Event](), nil
}

// mockHeadProviderWithHead implements HeadProvider that returns an existing event as head.
type mockHeadProviderWithHead struct {
	head Event
}

func (m *mockHeadProviderWithHead) Head() (types.Option[Event], error) {
	return types.Some(m.head), nil
}

// --- Factory tests ---

func TestFactoryCreateBasic(t *testing.T) {
	registry := DefaultRegistry()
	factory := NewEventFactory(registry)

	source := types.MustActorID("actor_00000000000000000000000000000001")
	causeID := types.MustEventID("019462a0-0000-7000-8000-000000000001")
	convID := types.MustConversationID("conv_00000000000000000000000000000001")

	content := TrustUpdatedContent{
		Actor:    types.MustActorID("actor_00000000000000000000000000000002"),
		Cause:    causeID,
		Current:  types.MustScore(0.85),
		Domain:   types.MustDomainScope("code_review"),
		Previous: types.MustScore(0.8),
	}

	ev, err := factory.Create(
		EventTypeTrustUpdated,
		source,
		content,
		[]types.EventID{causeID},
		convID,
		&mockHeadProviderEmpty{},
		&mockSigner{},
	)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	// Verify fields
	if ev.Version() != 1 {
		t.Errorf("Version = %d, want 1", ev.Version())
	}
	if ev.Type() != EventTypeTrustUpdated {
		t.Errorf("Type = %s, want %s", ev.Type().Value(), EventTypeTrustUpdated.Value())
	}
	if ev.Source() != source {
		t.Errorf("Source = %s, want %s", ev.Source().Value(), source.Value())
	}
	if ev.ConversationID() != convID {
		t.Errorf("ConversationID = %s, want %s", ev.ConversationID().Value(), convID.Value())
	}
	causes := ev.Causes()
	if len(causes) != 1 || causes[0] != causeID {
		t.Errorf("Causes = %v, want [%s]", causes, causeID.Value())
	}
	// Bootstrap has zero prevHash; empty head means prevHash should be zero
	if ev.PrevHash() != types.ZeroHash() {
		t.Errorf("PrevHash should be zero when head is empty, got %s", ev.PrevHash().Value())
	}
	// Hash should be non-zero (computed from canonical form)
	if ev.Hash() == types.ZeroHash() {
		t.Error("Hash should not be zero")
	}
	// Verify hash matches canonical form
	canonical := CanonicalForm(ev)
	expectedHash, err := ComputeHash(canonical)
	if err != nil {
		t.Fatalf("ComputeHash: %v", err)
	}
	if ev.Hash() != expectedHash {
		t.Errorf("Hash mismatch: event has %s, canonical produces %s", ev.Hash().Value(), expectedHash.Value())
	}
	// Should not be bootstrap
	if ev.IsBootstrap() {
		t.Error("factory-created event should not be bootstrap")
	}
}

func TestFactoryCreateWithHead(t *testing.T) {
	registry := DefaultRegistry()
	factory := NewEventFactory(registry)

	source := types.MustActorID("actor_00000000000000000000000000000001")
	causeID := types.MustEventID("019462a0-0000-7000-8000-000000000001")
	convID := types.MustConversationID("conv_00000000000000000000000000000001")

	// Create a head event with a known hash
	headEvent := NewBootstrapEvent(
		1,
		types.MustEventID("019462a0-0000-7000-8000-000000000099"),
		EventTypeSystemBootstrapped,
		mustTime("2023-11-14T22:13:20Z"),
		source,
		BootstrapContent{
			ActorID:      source,
			ChainGenesis: types.ZeroHash(),
			Timestamp:    mustTime("2023-11-14T22:13:20Z"),
		},
		convID,
		types.MustHash("abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"),
		types.MustSignature(make([]byte, 64)),
	)

	content := TrustUpdatedContent{
		Actor:    types.MustActorID("actor_00000000000000000000000000000002"),
		Cause:    causeID,
		Current:  types.MustScore(0.85),
		Domain:   types.MustDomainScope("code_review"),
		Previous: types.MustScore(0.8),
	}

	ev, err := factory.Create(
		EventTypeTrustUpdated,
		source,
		content,
		[]types.EventID{causeID},
		convID,
		&mockHeadProviderWithHead{head: headEvent},
		&mockSigner{},
	)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	// PrevHash should be the head event's hash
	if ev.PrevHash() != headEvent.Hash() {
		t.Errorf("PrevHash = %s, want %s (head hash)", ev.PrevHash().Value(), headEvent.Hash().Value())
	}
}

func TestFactoryCreateNoCauses(t *testing.T) {
	registry := DefaultRegistry()
	factory := NewEventFactory(registry)

	source := types.MustActorID("actor_00000000000000000000000000000001")
	convID := types.MustConversationID("conv_00000000000000000000000000000001")

	content := TrustUpdatedContent{
		Actor:    types.MustActorID("actor_00000000000000000000000000000002"),
		Cause:    types.MustEventID("019462a0-0000-7000-8000-000000000001"),
		Current:  types.MustScore(0.85),
		Domain:   types.MustDomainScope("code_review"),
		Previous: types.MustScore(0.8),
	}

	_, err := factory.Create(
		EventTypeTrustUpdated,
		source,
		content,
		[]types.EventID{}, // empty causes
		convID,
		&mockHeadProviderEmpty{},
		&mockSigner{},
	)
	if err == nil {
		t.Fatal("Create with no causes should return error")
	}

	// Also test with nil causes
	_, err = factory.Create(
		EventTypeTrustUpdated,
		source,
		content,
		nil, // nil causes
		convID,
		&mockHeadProviderEmpty{},
		&mockSigner{},
	)
	if err == nil {
		t.Fatal("Create with nil causes should return error")
	}
}

func TestFactoryCreateInvalidType(t *testing.T) {
	registry := DefaultRegistry()
	factory := NewEventFactory(registry)

	source := types.MustActorID("actor_00000000000000000000000000000001")
	causeID := types.MustEventID("019462a0-0000-7000-8000-000000000001")
	convID := types.MustConversationID("conv_00000000000000000000000000000001")

	// Use valid content but with a mismatched/unregistered event type
	unregisteredType := types.MustEventType("nonexistent.type")
	content := TrustUpdatedContent{
		Actor:    types.MustActorID("actor_00000000000000000000000000000002"),
		Cause:    causeID,
		Current:  types.MustScore(0.85),
		Domain:   types.MustDomainScope("code_review"),
		Previous: types.MustScore(0.8),
	}

	_, err := factory.Create(
		unregisteredType,
		source,
		content,
		[]types.EventID{causeID},
		convID,
		&mockHeadProviderEmpty{},
		&mockSigner{},
	)
	if err == nil {
		t.Fatal("Create with unregistered event type should return error")
	}

	// Also test content type mismatch: registered type but wrong content
	_, err = factory.Create(
		EventTypeEdgeCreated, // registered but doesn't match TrustUpdatedContent
		source,
		content, // TrustUpdatedContent, not EdgeCreatedContent
		[]types.EventID{causeID},
		convID,
		&mockHeadProviderEmpty{},
		&mockSigner{},
	)
	if err == nil {
		t.Fatal("Create with mismatched content type should return error")
	}
}

func TestFactoryBootstrapInit(t *testing.T) {
	registry := DefaultRegistry()
	factory := NewBootstrapFactory(registry)

	systemActor := types.MustActorID("actor_00000000000000000000000000000001")

	ev, err := factory.Init(systemActor, &mockSigner{})
	if err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	// Verify bootstrap properties
	if ev.Version() != 1 {
		t.Errorf("Version = %d, want 1", ev.Version())
	}
	if ev.Type() != EventTypeSystemBootstrapped {
		t.Errorf("Type = %s, want %s", ev.Type().Value(), EventTypeSystemBootstrapped.Value())
	}
	if ev.Source() != systemActor {
		t.Errorf("Source = %s, want %s", ev.Source().Value(), systemActor.Value())
	}
	if !ev.IsBootstrap() {
		t.Error("bootstrap event should be identified as bootstrap")
	}
	if ev.PrevHash() != types.ZeroHash() {
		t.Errorf("PrevHash should be zero, got %s", ev.PrevHash().Value())
	}
	if len(ev.Causes()) != 0 {
		t.Errorf("Causes should be empty, got %v", ev.Causes())
	}
	// Hash should be non-zero
	if ev.Hash() == types.ZeroHash() {
		t.Error("Hash should not be zero")
	}
	// Verify hash matches canonical form
	canonical := CanonicalForm(ev)
	expectedHash, err := ComputeHash(canonical)
	if err != nil {
		t.Fatalf("ComputeHash: %v", err)
	}
	if ev.Hash() != expectedHash {
		t.Errorf("Hash mismatch: event has %s, canonical produces %s", ev.Hash().Value(), expectedHash.Value())
	}
	// Verify content is BootstrapContent
	bc, ok := ev.Content().(BootstrapContent)
	if !ok {
		t.Fatalf("Content should be BootstrapContent, got %T", ev.Content())
	}
	if bc.ActorID != systemActor {
		t.Errorf("Content.ActorID = %s, want %s", bc.ActorID.Value(), systemActor.Value())
	}
	if bc.ChainGenesis != types.ZeroHash() {
		t.Errorf("Content.ChainGenesis should be zero hash")
	}
	// ConversationID should be the bootstrap conversation
	expectedConvID, _ := types.NewConversationID("conv_bootstrap_00000000000000000001")
	if ev.ConversationID() != expectedConvID {
		t.Errorf("ConversationID = %s, want %s", ev.ConversationID().Value(), expectedConvID.Value())
	}
}
