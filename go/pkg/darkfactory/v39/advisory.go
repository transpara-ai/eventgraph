package v39

import "fmt"

type advisoryReferenceKind string

const (
	advisoryReferenceMemory    advisoryReferenceKind = "memory"
	advisoryReferenceKnowledge advisoryReferenceKind = "knowledge"
)

func (s *InMemoryStore) RecordMemoryReference(reference *MemoryReference) (*MemoryReference, error) {
	if reference == nil {
		return nil, fmt.Errorf("%w: nil MemoryReference", ErrInvalidRecord)
	}
	if err := s.validateAdvisoryReferenceUse(&reference.AdvisoryReference, TypeMemoryReference); err != nil {
		return nil, err
	}
	stored, err := s.AppendRecord(reference)
	if err != nil {
		return nil, err
	}
	memoryReference, ok := stored.(*MemoryReference)
	if !ok {
		return nil, fmt.Errorf("%w: MemoryReference append returned %T", ErrInvalidRecord, stored)
	}
	if _, err := s.AppendEdge(derivedEdge(EdgeReferencedMemory, memoryReference.UsedInTask, memoryReference.CommonNode.ID, memoryReference.CommonNode)); err != nil {
		return nil, err
	}
	return memoryReference, nil
}

func (s *InMemoryStore) RecordKnowledgeReference(reference *KnowledgeReference) (*KnowledgeReference, error) {
	if reference == nil {
		return nil, fmt.Errorf("%w: nil KnowledgeReference", ErrInvalidRecord)
	}
	if err := s.validateAdvisoryReferenceUse(&reference.AdvisoryReference, TypeKnowledgeReference); err != nil {
		return nil, err
	}
	stored, err := s.AppendRecord(reference)
	if err != nil {
		return nil, err
	}
	knowledgeReference, ok := stored.(*KnowledgeReference)
	if !ok {
		return nil, fmt.Errorf("%w: KnowledgeReference append returned %T", ErrInvalidRecord, stored)
	}
	if _, err := s.AppendEdge(derivedEdge(EdgeReferencedKnowledge, knowledgeReference.UsedInTask, knowledgeReference.CommonNode.ID, knowledgeReference.CommonNode)); err != nil {
		return nil, err
	}
	return knowledgeReference, nil
}

func (s *InMemoryStore) validateAdvisoryReferenceUse(reference *AdvisoryReference, typ string) error {
	if _, ok := s.mustGetTask(reference.UsedInTask); !ok {
		return fmt.Errorf("%w: Task %s", ErrNotFound, reference.UsedInTask)
	}
	for _, contradictionRef := range reference.ContradictionRefs {
		contradiction, ok := s.mustGetContradictionLog(contradictionRef)
		if !ok {
			return fmt.Errorf("%w: ContradictionLog %s", ErrNotFound, contradictionRef)
		}
		if contradictionBlocksHighRiskUse(contradiction) && isHighRiskScope(reference.RiskScope) {
			return fieldError(typ, "contradiction_refs", "open high or critical contradiction blocks high-risk use")
		}
	}
	return nil
}

func (s *InMemoryStore) AdvisoryReferenceEvidencePath(releaseCandidateID string) (RequiredPath, error) {
	path := RequiredPath{Name: "Material advisory influence -> MemoryReference / KnowledgeReference", NodeIDs: []string{releaseCandidateID}}
	rc, ok := s.mustGetReleaseCandidate(releaseCandidateID)
	if !ok {
		path.Missing = append(path.Missing, "ReleaseCandidate "+releaseCandidateID)
		return path, path.Err()
	}
	orderPath, _ := s.FactoryOrderRequirementAcceptanceTask(rc.FactoryOrderID)
	path.EdgeIDs = append(path.EdgeIDs, orderPath.EdgeIDs...)
	path.Missing = append(path.Missing, orderPath.Missing...)

	for _, taskID := range taskIDsFromPath(s, orderPath) {
		task, ok := s.mustGetTask(taskID)
		if !ok {
			path.Missing = append(path.Missing, "Task "+taskID)
			continue
		}
		for _, sourceRef := range advisorySourceRefs(task.SourceRefs, advisoryReferenceMemory) {
			if !s.addMatchingAdvisoryReference(&path, taskID, sourceRef, advisoryReferenceMemory) {
				path.Missing = append(path.Missing, "MemoryReference for "+sourceRef+" used in Task "+taskID)
			}
		}
		for _, sourceRef := range advisorySourceRefs(task.SourceRefs, advisoryReferenceKnowledge) {
			if !s.addMatchingAdvisoryReference(&path, taskID, sourceRef, advisoryReferenceKnowledge) {
				path.Missing = append(path.Missing, "KnowledgeReference for "+sourceRef+" used in Task "+taskID)
			}
		}
		s.addLinkedAdvisoryReferences(&path, taskID, advisoryReferenceMemory)
		s.addLinkedAdvisoryReferences(&path, taskID, advisoryReferenceKnowledge)
	}

	path.Completed = len(path.Missing) == 0
	return path, path.Err()
}

func (s *InMemoryStore) addMatchingAdvisoryReference(path *RequiredPath, taskID, sourceRef string, kind advisoryReferenceKind) bool {
	for _, record := range s.advisoryReferencesForTask(taskID, kind) {
		if record.CommonNode.ID == sourceRef || record.SourceRef == sourceRef || record.SourceHashOrImmutableLocator == sourceRef {
			path.NodeIDs = appendUniqueStrings(path.NodeIDs, taskID, record.CommonNode.ID)
			path.EdgeIDs = appendUniqueStrings(path.EdgeIDs, edgeIDsBetween(s, taskID, record.CommonNode.ID, edgeTypeForAdvisoryKind(kind))...)
			return true
		}
	}
	return false
}

func (s *InMemoryStore) addLinkedAdvisoryReferences(path *RequiredPath, taskID string, kind advisoryReferenceKind) {
	for _, edge := range s.outgoingEdges(taskID, edgeTypeForAdvisoryKind(kind)) {
		for _, record := range s.advisoryReferencesForTask(taskID, kind) {
			if edge.ToID == record.CommonNode.ID {
				path.NodeIDs = appendUniqueStrings(path.NodeIDs, taskID, record.CommonNode.ID)
				path.EdgeIDs = appendUniqueStrings(path.EdgeIDs, edge.ID)
			}
		}
	}
}

func (s *InMemoryStore) advisoryReferencesForTask(taskID string, kind advisoryReferenceKind) []AdvisoryReference {
	var out []AdvisoryReference
	recordType := TypeMemoryReference
	if kind == advisoryReferenceKnowledge {
		recordType = TypeKnowledgeReference
	}
	for _, record := range s.ByType(recordType) {
		switch typed := record.(type) {
		case *MemoryReference:
			if typed.UsedInTask == taskID {
				out = append(out, typed.AdvisoryReference)
			}
		case *KnowledgeReference:
			if typed.UsedInTask == taskID {
				out = append(out, typed.AdvisoryReference)
			}
		}
	}
	return out
}

func advisorySourceRefs(sourceRefs []string, kind advisoryReferenceKind) []string {
	var out []string
	for _, sourceRef := range sourceRefs {
		switch kind {
		case advisoryReferenceMemory:
			if hasAnyPrefix(sourceRef, "memory:", "mem:") {
				out = append(out, sourceRef)
			}
		case advisoryReferenceKnowledge:
			if hasAnyPrefix(sourceRef, "knowledge:", "know:") {
				out = append(out, sourceRef)
			}
		}
	}
	return out
}

func hasAnyPrefix(value string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if len(value) >= len(prefix) && value[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

func contradictionBlocksHighRiskUse(contradiction *ContradictionLog) bool {
	if contradiction == nil || contradiction.CommonNode.Status == nil || *contradiction.CommonNode.Status != "open" {
		return false
	}
	return contradiction.Severity == "high" || contradiction.Severity == "critical"
}

func edgeTypeForAdvisoryKind(kind advisoryReferenceKind) string {
	if kind == advisoryReferenceKnowledge {
		return EdgeReferencedKnowledge
	}
	return EdgeReferencedMemory
}

func edgeIDsBetween(s *InMemoryStore, fromID, toID, edgeType string) []string {
	var ids []string
	for _, edge := range s.outgoingEdges(fromID, edgeType) {
		if edge.ToID == toID {
			ids = append(ids, edge.ID)
		}
	}
	return ids
}