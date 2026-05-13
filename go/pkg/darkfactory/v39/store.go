package v39

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNotFound            = errors.New("dark factory v3.9 record not found")
	ErrIdempotencyConflict = errors.New("dark factory v3.9 idempotency conflict")
	ErrDuplicateRecordID   = errors.New("dark factory v3.9 duplicate record id")
)

type InMemoryStore struct {
	mu            sync.RWMutex
	records       map[string]Record
	canonicalByID map[string][]byte
	byType        map[string][]string
	byIdem        map[string]string
	edges         map[string]CommonEdge
	outEdges      map[string][]string
	inEdges       map[string][]string
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		records:       map[string]Record{},
		canonicalByID: map[string][]byte{},
		byType:        map[string][]string{},
		byIdem:        map[string]string{},
		edges:         map[string]CommonEdge{},
		outEdges:      map[string][]string{},
		inEdges:       map[string][]string{},
	}
}

func (s *InMemoryStore) AppendRecord(r Record) (Record, error) {
	if err := ValidateRecord(r); err != nil {
		return nil, err
	}
	common := r.GetCommon()
	canonical, err := CanonicalJSON(r)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if existingID, ok := s.byIdem[common.IdempotencyKey]; ok {
		existingBytes := s.canonicalByID[existingID]
		if !bytes.Equal(existingBytes, canonical) {
			return nil, fmt.Errorf("%w: key %s", ErrIdempotencyConflict, common.IdempotencyKey)
		}
		return s.records[existingID], nil
	}
	if existing, ok := s.records[common.ID]; ok {
		if !bytes.Equal(s.canonicalByID[common.ID], canonical) {
			return nil, fmt.Errorf("%w: %s: %w", ErrDuplicateRecordID, common.ID, ErrImmutable)
		}
		return existing, nil
	}

	s.records[common.ID] = r
	s.canonicalByID[common.ID] = append([]byte(nil), canonical...)
	s.byType[common.Type] = append(s.byType[common.Type], common.ID)
	s.byIdem[common.IdempotencyKey] = common.ID
	return r, nil
}

func (s *InMemoryStore) AppendEdge(e CommonEdge) (CommonEdge, error) {
	if e.ID == "" || e.Type == "" || e.FromID == "" || e.ToID == "" || e.CreatedAt.IsZero() || e.CreatedBy == "" || e.CorrelationID == "" || e.IdempotencyKey == "" {
		return CommonEdge{}, fmt.Errorf("%w: edge missing required field", ErrInvalidRecord)
	}
	canonical, err := CanonicalJSON(e)
	if err != nil {
		return CommonEdge{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if existingID, ok := s.byIdem[e.IdempotencyKey]; ok {
		existingEdge, ok := s.edges[existingID]
		if !ok {
			return CommonEdge{}, fmt.Errorf("%w: key %s", ErrIdempotencyConflict, e.IdempotencyKey)
		}
		existingBytes, _ := CanonicalJSON(existingEdge)
		if !bytes.Equal(existingBytes, canonical) {
			return CommonEdge{}, fmt.Errorf("%w: key %s", ErrIdempotencyConflict, e.IdempotencyKey)
		}
		return existingEdge, nil
	}
	if _, ok := s.edges[e.ID]; ok {
		return CommonEdge{}, fmt.Errorf("%w: %s: %w", ErrDuplicateRecordID, e.ID, ErrImmutable)
	}
	if _, ok := s.records[e.FromID]; !ok {
		return CommonEdge{}, fmt.Errorf("%w: from_id %s", ErrNotFound, e.FromID)
	}
	if _, ok := s.records[e.ToID]; !ok {
		return CommonEdge{}, fmt.Errorf("%w: to_id %s", ErrNotFound, e.ToID)
	}
	s.edges[e.ID] = e
	s.outEdges[e.FromID] = append(s.outEdges[e.FromID], e.ID)
	s.inEdges[e.ToID] = append(s.inEdges[e.ToID], e.ID)
	s.byIdem[e.IdempotencyKey] = e.ID
	return e, nil
}

func (s *InMemoryStore) Get(id string) (Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.records[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, id)
	}
	return r, nil
}

func (s *InMemoryStore) ByType(typ string) []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := s.byType[typ]
	out := make([]Record, 0, len(ids))
	for _, id := range ids {
		out = append(out, s.records[id])
	}
	return out
}

func (s *InMemoryStore) EdgesFrom(id string) []CommonEdge {
	s.mu.RLock()
	defer s.mu.RUnlock()
	edgeIDs := s.outEdges[id]
	out := make([]CommonEdge, 0, len(edgeIDs))
	for _, edgeID := range edgeIDs {
		out = append(out, s.edges[edgeID])
	}
	return out
}

func (s *InMemoryStore) EdgesTo(id string) []CommonEdge {
	s.mu.RLock()
	defer s.mu.RUnlock()
	edgeIDs := s.inEdges[id]
	out := make([]CommonEdge, 0, len(edgeIDs))
	for _, edgeID := range edgeIDs {
		out = append(out, s.edges[edgeID])
	}
	return out
}

func (s *InMemoryStore) CanonicalRecord(id string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.canonicalByID[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, id)
	}
	return append([]byte(nil), b...), nil
}
