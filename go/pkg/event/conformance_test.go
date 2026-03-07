package event

import (
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
	expected := `1||019462a0-0000-7000-8000-000000000001|system.bootstrapped|actor_00000000000000000000000000000001|conv_00000000000000000000000000000001|1700000000000000000|{"ActorID":"actor_00000000000000000000000000000001","ChainGenesis":"0000000000000000000000000000000000000000000000000000000000000000","Timestamp":"2023-11-14T22:13:20Z"}`
	if canonical != expected {
		t.Errorf("canonical form mismatch:\n  got:  %s\n  want: %s", canonical, expected)
	}

	hash, err := ComputeHash(canonical)
	if err != nil {
		t.Fatalf("ComputeHash: %v", err)
	}
	expectedHash := "88a1c89ffad29455acaa28c66aef8c34db2532fdb47e37f16f09ab6fd91ceeb4"
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
	expectedHash := "04bfffee1d6d192302856af63d02d92f137083872a7775e4c778829021712bc5"
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
	expectedHash := "1abc67e3790b10bd954951d0954ae6ccbe5e9f4e04a174e0090f7b003798f846"
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
	// Should start with "1||" (empty prev_hash between pipes)
	if canonical[:3] != "1||" {
		t.Errorf("bootstrap canonical should start with '1||', got: %s", canonical[:10])
	}
}
