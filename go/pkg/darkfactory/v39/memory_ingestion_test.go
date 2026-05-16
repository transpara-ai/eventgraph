package v39

import (
	"errors"
	"strings"
	"testing"
)

func TestMemoryIngestionQuarantinesSecretBearingMemory(t *testing.T) {
	store := NewInMemoryStore()

	result, err := store.IngestMemory(MemoryIngestionInput{
		MemoryID:     "mem_secret_ingest",
		SourceSystem: "sessiondb",
		SourceRef:    "memory://session/secret",
		Content:      "operator note with password=super-secret",
		Scope:        "factory-order:fo_001",
		CreatedAt:    fixedTime,
		CreatedBy:    "act_001",
	})
	if err != nil {
		t.Fatalf("ingest secret memory: %v", err)
	}
	if result.Classification != MemoryClassificationSecret || result.RedactionState != MemoryRedactionQuarantined {
		t.Fatalf("secret memory was not classified and quarantined: %+v", result)
	}
	if got := len(store.ByType(TypeMemoryIngested)); got != 1 {
		t.Fatalf("expected MemoryIngested event, got %d", got)
	}
	if got := len(store.ByType(TypeMemoryIndexed)); got != 1 {
		t.Fatalf("expected MemoryIndexed event, got %d", got)
	}
	if got := len(store.ByType(TypeMemoryScopeAssigned)); got != 1 {
		t.Fatalf("expected MemoryScopeAssigned event, got %d", got)
	}
	if got := len(store.ByType(TypeMemoryRedactionApplied)); got != 1 {
		t.Fatalf("expected MemoryRedactionApplied event, got %d", got)
	}
	if got := store.RetrieveMemory("factory-order:fo_001", fixedTime); len(got) != 0 {
		t.Fatalf("quarantined memory should not be retrievable, got %+v", got)
	}
}

func TestMemoryIngestionScopedRetrieval(t *testing.T) {
	store := NewInMemoryStore()
	for _, input := range []MemoryIngestionInput{
		{MemoryID: "mem_order", SourceSystem: "sessiondb", SourceRef: "memory://order/context", Content: "public order context", Classification: MemoryClassificationPublic, Scope: "factory-order:fo_001", CreatedAt: fixedTime, CreatedBy: "act_001"},
		{MemoryID: "mem_other", SourceSystem: "sessiondb", SourceRef: "memory://order/other", Content: "public other context", Classification: MemoryClassificationPublic, Scope: "factory-order:fo_002", CreatedAt: fixedTime, CreatedBy: "act_001"},
	} {
		if _, err := store.IngestMemory(input); err != nil {
			t.Fatalf("ingest %s: %v", input.MemoryID, err)
		}
	}

	got := store.RetrieveMemory("factory-order:fo_001", fixedTime)
	if len(got) != 1 || got[0].MemoryID != "mem_order" {
		t.Fatalf("scoped retrieval returned wrong memories: %+v", got)
	}
}

func TestRawMemoryIngestionDoesNotSatisfyReleaseEvidence(t *testing.T) {
	store := stage6BaseStoreWithTaskSourceRefs(t, []string{"memory://planning/context"})
	if _, err := store.IngestMemory(MemoryIngestionInput{
		MemoryID:     "mem_planning_raw",
		SourceSystem: "sessiondb",
		SourceRef:    "memory://planning/context",
		Content:      "planning context that was ingested but not referenced",
		Scope:        "factory-order:fo_001",
		CreatedAt:    fixedTime,
		CreatedBy:    "act_001",
	}); err != nil {
		t.Fatalf("ingest raw memory: %v", err)
	}
	recordStage6ReleaseCandidate(t, store, "rc_raw_memory")

	_, err := store.CertifyReleaseCandidate(stage6Certification("cert_raw_memory", "rc_raw_memory"))
	if !errors.Is(err, ErrRequiredPathMissing) {
		t.Fatalf("expected raw memory without MemoryReference to block certification, got %v", err)
	}
	if !strings.Contains(err.Error(), "MemoryReference for memory://planning/context") {
		t.Fatalf("expected missing MemoryReference in error, got %v", err)
	}
}

func TestRetrievedMemoryInfluenceCreatesReferenceOnlyForMaterialUse(t *testing.T) {
	store := stage6BaseStoreWithTaskSourceRefs(t, []string{"memory://planning/context"})
	if _, err := store.IngestMemory(MemoryIngestionInput{
		MemoryID:     "mem_planning_context",
		SourceSystem: "sessiondb",
		SourceRef:    "memory://planning/context",
		Content:      "reviewed planning context",
		Scope:        "task:tsk_001",
		CreatedAt:    fixedTime,
		CreatedBy:    "act_001",
	}); err != nil {
		t.Fatalf("ingest planning memory: %v", err)
	}
	retrieved := store.RetrieveMemory("task:tsk_001", fixedTime)
	if len(retrieved) != 1 {
		t.Fatalf("expected one retrieved memory, got %+v", retrieved)
	}
	if _, err := store.RecordRetrievedMemoryInfluence(retrieved[0], "tsk_001", "act_001", "ambient_context", "non-material glance", "medium"); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expected non-material influence to be rejected, got %v", err)
	}
	if _, err := store.RecordRetrievedMemoryInfluence(retrieved[0], "tsk_001", "act_001", MemoryInfluencePlanning, "selected repair plan from retrieved context", "medium"); err != nil {
		t.Fatalf("record material memory influence: %v", err)
	}
	recordStage6ReleaseCandidate(t, store, "rc_material_memory")

	cert, err := store.CertifyReleaseCandidate(stage6Certification("cert_material_memory", "rc_material_memory"))
	if err != nil {
		t.Fatalf("certify with material memory reference: %v", err)
	}
	if !containsString(cert.EvidenceRefs, "memory:mem_planning_context:reference:tsk_001:planning") {
		t.Fatalf("certification missing memory reference evidence: %+v", cert.EvidenceRefs)
	}
}
