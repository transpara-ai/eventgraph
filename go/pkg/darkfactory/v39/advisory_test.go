package v39

import (
	"errors"
	"strings"
	"testing"
)

func TestMemoryReferenceRejectsQuarantinedSecretBearingMemory(t *testing.T) {
	store := NewInMemoryStore()
	ref := &MemoryReference{AdvisoryReference: advisory("mem_secret", TypeMemoryReference, "tsk_001")}
	ref.SourceSystem = "sessiondb"
	ref.SourceRef = "memory://secret/session"
	ref.TrustLevel = "quarantined"
	ref.RedactionState = "quarantined"

	if _, err := store.AppendRecord(ref); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expected quarantined memory to be rejected, got %v", err)
	}
}

func TestKnowledgeReferenceRejectsStaleKnowledgeForHighRiskUse(t *testing.T) {
	store := NewInMemoryStore()
	ref := &KnowledgeReference{AdvisoryReference: advisory("know_stale", TypeKnowledgeReference, "tsk_001")}
	ref.SourceRef = "knowledge://recipe/security"
	ref.RiskScope = "high-risk planning"
	ref.TrustLevel = "stale"
	ref.FreshnessStatus = "stale"

	if _, err := store.AppendRecord(ref); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expected stale high-risk knowledge to be rejected, got %v", err)
	}
}

func TestKnowledgeReferenceBlocksOpenHighContradictionForHighRiskUse(t *testing.T) {
	store := stage6BaseStore(t)
	appendRecord(t, store, &ContradictionLog{
		CommonNode:      common("contradiction_high", TypeContradictionLog, "open"),
		ContradictionID: "contradiction_high",
		ClaimARef:       "knowledge://recipe/security",
		ClaimBRef:       "eventgraph://gate/security",
		Severity:        "high",
	})
	ref := &KnowledgeReference{AdvisoryReference: advisory("know_contradicted", TypeKnowledgeReference, "tsk_001")}
	ref.SourceRef = "knowledge://recipe/security"
	ref.RiskScope = "high"
	ref.ContradictionRefs = []string{"contradiction_high"}

	if _, err := store.RecordKnowledgeReference(ref); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("expected open high contradiction to block high-risk knowledge use, got %v", err)
	}
}

func TestCertificationFailsWhenMaterialMemoryInfluenceLacksReference(t *testing.T) {
	store := stage6BaseStoreWithTaskSourceRefs(t, []string{"memory://planning/context"})
	recordStage6ReleaseCandidate(t, store, "rc_missing_memory_ref")

	_, err := store.CertifyReleaseCandidate(stage6Certification("cert_missing_memory_ref", "rc_missing_memory_ref"))
	if !errors.Is(err, ErrRequiredPathMissing) {
		t.Fatalf("expected missing memory reference to block certification, got %v", err)
	}
	if !strings.Contains(err.Error(), "MemoryReference for memory://planning/context") {
		t.Fatalf("expected missing memory reference in error, got %v", err)
	}
}

func TestCertificationPassesWithMaterialMemoryInfluenceReference(t *testing.T) {
	store := stage6BaseStoreWithTaskSourceRefs(t, []string{"memory://planning/context"})
	ref := &MemoryReference{AdvisoryReference: advisory("mem_planning_context", TypeMemoryReference, "tsk_001")}
	ref.SourceRef = "memory://planning/context"
	ref.RiskScope = "medium"
	if _, err := store.RecordMemoryReference(ref); err != nil {
		t.Fatalf("record memory reference: %v", err)
	}
	recordStage6ReleaseCandidate(t, store, "rc_with_memory_ref")

	cert, err := store.CertifyReleaseCandidate(stage6Certification("cert_with_memory_ref", "rc_with_memory_ref"))
	if err != nil {
		t.Fatalf("certify with memory reference: %v", err)
	}
	if !containsString(cert.EvidenceRefs, "mem_planning_context") {
		t.Fatalf("certification missing memory reference evidence: %+v", cert.EvidenceRefs)
	}
}

func stage6BaseStoreWithTaskSourceRefs(t *testing.T, sourceRefs []string) *InMemoryStore {
	t.Helper()
	store := NewInMemoryStore()
	for _, record := range completeTier0Records() {
		switch typed := record.(type) {
		case *ReleaseCandidate, *Certification, *Rejection, *AuditReport:
			continue
		case *Task:
			if typed.CommonNode.ID == "tsk_001" {
				typed.SourceRefs = sourceRefs
			}
		}
		appendRecord(t, store, record)
	}
	appendStage6TraceEdges(t, store)
	return store
}