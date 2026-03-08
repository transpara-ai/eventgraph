// Package sqlitestore provides a SQLite-backed implementation of store.Store.
//
// Uses modernc.org/sqlite for a pure-Go SQLite driver (no CGo required).
package sqlitestore

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS events (
    position        INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id        TEXT NOT NULL UNIQUE,
    event_type      TEXT NOT NULL,
    version         INTEGER NOT NULL,
    timestamp_nanos INTEGER NOT NULL,
    source          TEXT NOT NULL,
    content         TEXT NOT NULL,
    causes          TEXT NOT NULL,
    conversation_id TEXT NOT NULL,
    hash            TEXT NOT NULL,
    prev_hash       TEXT NOT NULL,
    signature       BLOB NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_events_type ON events(event_type);
CREATE INDEX IF NOT EXISTS idx_events_source ON events(source);
CREATE INDEX IF NOT EXISTS idx_events_conversation ON events(conversation_id);

CREATE TABLE IF NOT EXISTS edges (
    id          TEXT PRIMARY KEY,
    from_actor  TEXT NOT NULL,
    to_actor    TEXT NOT NULL,
    edge_type   TEXT NOT NULL,
    weight      REAL NOT NULL,
    direction   TEXT NOT NULL,
    scope       TEXT,
    created_at  INTEGER NOT NULL,
    expires_at  INTEGER
);
CREATE INDEX IF NOT EXISTS idx_edges_from ON edges(from_actor, edge_type);
CREATE INDEX IF NOT EXISTS idx_edges_to ON edges(to_actor, edge_type);
`

// SQLiteStore implements store.Store backed by a SQLite database.
type SQLiteStore struct {
	mu sync.Mutex
	db *sql.DB
}

// New opens a SQLite store at the given path. Use ":memory:" for testing.
func New(dsn string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, &store.StoreUnavailableError{Reason: err.Error()}
	}
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, &store.StoreUnavailableError{Reason: err.Error()}
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, &store.StoreUnavailableError{Reason: err.Error()}
	}
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Append(ev event.Event) (event.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Idempotency: check if event already exists
	var existingHash string
	err := s.db.QueryRow("SELECT hash FROM events WHERE event_id = ?", ev.ID().Value()).Scan(&existingHash)
	if err == nil {
		if existingHash != ev.Hash().Value() {
			return event.Event{}, &store.HashMismatchError{
				EventID:  ev.ID(),
				Computed: ev.Hash(),
				Stored:   types.MustHash(existingHash),
			}
		}
		return s.getUnlocked(ev.ID())
	}

	// Verify chain continuity
	var headHash string
	err = s.db.QueryRow("SELECT hash FROM events ORDER BY position DESC LIMIT 1").Scan(&headHash)
	if err == sql.ErrNoRows {
		if ev.PrevHash() != types.ZeroHash() {
			return event.Event{}, &store.ChainIntegrityViolationError{
				Position: 0,
				Expected: types.ZeroHash(),
				Actual:   ev.PrevHash(),
			}
		}
	} else if err == nil {
		if ev.PrevHash().Value() != headHash {
			return event.Event{}, &store.ChainIntegrityViolationError{
				Position: 0,
				Expected: types.MustHash(headHash),
				Actual:   ev.PrevHash(),
			}
		}
	} else {
		return event.Event{}, &store.StoreUnavailableError{Reason: err.Error()}
	}

	// Recompute and verify hash
	canonical := event.CanonicalForm(ev)
	computed, computeErr := event.ComputeHash(canonical)
	if computeErr != nil {
		return event.Event{}, computeErr
	}
	if computed != ev.Hash() {
		return event.Event{}, &store.HashMismatchError{
			EventID:  ev.ID(),
			Computed: computed,
			Stored:   ev.Hash(),
		}
	}

	// Verify causal predecessors exist (CAUSALITY invariant)
	if !ev.IsBootstrap() {
		for _, causeID := range ev.Causes() {
			var exists int
			err := s.db.QueryRow("SELECT 1 FROM events WHERE event_id = ?", causeID.Value()).Scan(&exists)
			if err != nil {
				return event.Event{}, &store.CausalLinkMissingError{
					EventID:      ev.ID(),
					MissingCause: causeID,
				}
			}
		}
	}

	// Serialize causes and content
	causeStrs := make([]string, 0, len(ev.Causes()))
	for _, c := range ev.Causes() {
		causeStrs = append(causeStrs, c.Value())
	}
	causesJSON, _ := json.Marshal(causeStrs)
	contentJSON := event.ContentJSON(ev)

	_, err = s.db.Exec(
		`INSERT INTO events (event_id, event_type, version, timestamp_nanos, source,
		 content, causes, conversation_id, hash, prev_hash, signature)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		ev.ID().Value(), ev.Type().Value(), ev.Version(),
		ev.Timestamp().Nanoseconds(), ev.Source().Value(),
		contentJSON, string(causesJSON),
		ev.ConversationID().Value(), ev.Hash().Value(), ev.PrevHash().Value(),
		ev.Signature().Bytes(),
	)
	if err != nil {
		return event.Event{}, &store.StoreUnavailableError{Reason: err.Error()}
	}

	// Handle edge creation
	if ev.Type() == event.EventTypeEdgeCreated {
		ec, ok := ev.Content().(event.EdgeCreatedContent)
		if ok {
			edgeID, _ := types.NewEdgeID(ev.ID().Value())
			edge, edgeErr := event.NewEdge(
				edgeID, ec.From, ec.To, ec.EdgeType, ec.Weight, ec.Direction,
				ec.Scope, nil, ev.Timestamp(), ec.ExpiresAt,
			)
			if edgeErr == nil {
				s.insertEdge(edge)
			}
		}
	}

	// Handle edge supersession
	if ev.Type() == event.EventTypeEdgeSuperseded {
		ec, ok := ev.Content().(event.EdgeSupersededContent)
		if ok {
			s.db.Exec("DELETE FROM edges WHERE id = ?", ec.PreviousEdge.Value())
		}
	}

	return ev, nil
}

func (s *SQLiteStore) insertEdge(e event.Edge) {
	var scope sql.NullString
	if e.Scope().IsSome() {
		scope.String = e.Scope().Unwrap().Value()
		scope.Valid = true
	}
	s.db.Exec(
		`INSERT OR REPLACE INTO edges (id, from_actor, to_actor, edge_type, weight, direction, scope, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID().Value(), e.From().Value(), e.To().Value(), string(e.Type()),
		e.Weight().Value(), string(e.Direction()), scope, e.CreatedAt().Nanoseconds(),
	)
}

func (s *SQLiteStore) getUnlocked(id types.EventID) (event.Event, error) {
	row := s.db.QueryRow("SELECT * FROM events WHERE event_id = ?", id.Value())
	return s.scanEvent(row)
}

func (s *SQLiteStore) scanEvent(row *sql.Row) (event.Event, error) {
	var (
		position      int
		eventID       string
		eventType     string
		version       int
		timestampNs   int64
		source        string
		contentJSON   string
		causesJSON    string
		conversationID string
		hash          string
		prevHash      string
		sig           []byte
	)
	err := row.Scan(&position, &eventID, &eventType, &version, &timestampNs,
		&source, &contentJSON, &causesJSON, &conversationID, &hash, &prevHash, &sig)
	if err != nil {
		return event.Event{}, &store.EventNotFoundError{ID: types.MustEventID(eventID)}
	}

	var causeStrs []string
	json.Unmarshal([]byte(causesJSON), &causeStrs)
	causes := make([]types.EventID, 0, len(causeStrs))
	for _, c := range causeStrs {
		causes = append(causes, types.MustEventID(c))
	}

	content := event.UnmarshalContent(types.MustEventType(eventType), contentJSON)

	signature := types.MustSignature(sig)
	ts := types.MustTimestamp(timestampNs)

	return event.NewEvent(
		version,
		types.MustEventID(eventID),
		types.MustEventType(eventType),
		ts,
		types.MustActorID(source),
		content,
		causes,
		types.MustConversationID(conversationID),
		types.MustHash(hash),
		types.MustHash(prevHash),
		signature,
	), nil
}

func (s *SQLiteStore) Get(id types.EventID) (event.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.getUnlocked(id)
}

func (s *SQLiteStore) Head() (types.Option[event.Event], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	row := s.db.QueryRow("SELECT * FROM events ORDER BY position DESC LIMIT 1")
	ev, err := s.scanEvent(row)
	if err != nil {
		return types.None[event.Event](), nil
	}
	return types.Some(ev), nil
}

func (s *SQLiteStore) Recent(limit int, after types.Option[types.Cursor]) (types.Page[event.Event], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 100
	}

	var rows *sql.Rows
	var err error
	if after.IsSome() {
		cursor := after.Unwrap()
		var startPos int
		scanErr := s.db.QueryRow("SELECT position FROM events WHERE event_id = ?", cursor.Value()).Scan(&startPos)
		if scanErr != nil {
			return types.Page[event.Event]{}, &store.InvalidCursorError{Cursor: cursor.Value()}
		}
		rows, err = s.db.Query("SELECT * FROM events WHERE position < ? ORDER BY position DESC LIMIT ?", startPos, limit)
	} else {
		rows, err = s.db.Query("SELECT * FROM events ORDER BY position DESC LIMIT ?", limit)
	}
	if err != nil {
		return types.Page[event.Event]{}, &store.StoreUnavailableError{Reason: err.Error()}
	}
	defer rows.Close()

	items := s.scanRows(rows)
	hasMore := len(items) == limit
	var cursor types.Option[types.Cursor]
	if hasMore && len(items) > 0 {
		c := types.MustCursor(items[len(items)-1].ID().Value())
		cursor = types.Some(c)
	}
	return types.NewPage(items, cursor, hasMore), nil
}

func (s *SQLiteStore) scanRows(rows *sql.Rows) []event.Event {
	var result []event.Event
	for rows.Next() {
		var (
			position       int
			eventID        string
			eventType      string
			version        int
			timestampNs    int64
			source         string
			contentJSON    string
			causesJSON     string
			conversationID string
			hash           string
			prevHash       string
			sig            []byte
		)
		if err := rows.Scan(&position, &eventID, &eventType, &version, &timestampNs,
			&source, &contentJSON, &causesJSON, &conversationID, &hash, &prevHash, &sig); err != nil {
			continue
		}

		var causeStrs []string
		json.Unmarshal([]byte(causesJSON), &causeStrs)
		causes := make([]types.EventID, 0, len(causeStrs))
		for _, c := range causeStrs {
			causes = append(causes, types.MustEventID(c))
		}

		content := event.UnmarshalContent(types.MustEventType(eventType), contentJSON)
		signature := types.MustSignature(sig)
		ts := types.MustTimestamp(timestampNs)

		ev := event.NewEvent(version, types.MustEventID(eventID), types.MustEventType(eventType),
			ts, types.MustActorID(source), content, causes,
			types.MustConversationID(conversationID), types.MustHash(hash),
			types.MustHash(prevHash), signature)
		result = append(result, ev)
	}
	return result
}

func (s *SQLiteStore) ByType(eventType types.EventType, limit int, after types.Option[types.Cursor]) (types.Page[event.Event], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.Query("SELECT * FROM events WHERE event_type = ? ORDER BY position DESC LIMIT ?",
		eventType.Value(), limit)
	if err != nil {
		return types.Page[event.Event]{}, &store.StoreUnavailableError{Reason: err.Error()}
	}
	defer rows.Close()

	items := s.scanRows(rows)
	return types.NewPage(items, types.None[types.Cursor](), false), nil
}

func (s *SQLiteStore) BySource(source types.ActorID, limit int, after types.Option[types.Cursor]) (types.Page[event.Event], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.Query("SELECT * FROM events WHERE source = ? ORDER BY position DESC LIMIT ?",
		source.Value(), limit)
	if err != nil {
		return types.Page[event.Event]{}, &store.StoreUnavailableError{Reason: err.Error()}
	}
	defer rows.Close()

	items := s.scanRows(rows)
	return types.NewPage(items, types.None[types.Cursor](), false), nil
}

func (s *SQLiteStore) ByConversation(id types.ConversationID, limit int, after types.Option[types.Cursor]) (types.Page[event.Event], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.Query("SELECT * FROM events WHERE conversation_id = ? ORDER BY position DESC LIMIT ?",
		id.Value(), limit)
	if err != nil {
		return types.Page[event.Event]{}, &store.StoreUnavailableError{Reason: err.Error()}
	}
	defer rows.Close()

	items := s.scanRows(rows)
	return types.NewPage(items, types.None[types.Cursor](), false), nil
}

func (s *SQLiteStore) Since(afterID types.EventID, limit int) (types.Page[event.Event], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 100
	}

	var startPos int
	err := s.db.QueryRow("SELECT position FROM events WHERE event_id = ?", afterID.Value()).Scan(&startPos)
	if err != nil {
		return types.Page[event.Event]{}, &store.EventNotFoundError{ID: afterID}
	}

	rows, err := s.db.Query("SELECT * FROM events WHERE position > ? ORDER BY position ASC LIMIT ?", startPos, limit)
	if err != nil {
		return types.Page[event.Event]{}, &store.StoreUnavailableError{Reason: err.Error()}
	}
	defer rows.Close()

	items := s.scanRows(rows)

	// Check hasMore
	var total int
	s.db.QueryRow("SELECT COUNT(*) FROM events WHERE position > ?", startPos).Scan(&total)
	hasMore := total > len(items)

	return types.NewPage(items, types.None[types.Cursor](), hasMore), nil
}

func (s *SQLiteStore) Ancestors(id types.EventID, maxDepth int) ([]event.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	row := s.db.QueryRow("SELECT causes FROM events WHERE event_id = ?", id.Value())
	var causesJSON string
	if err := row.Scan(&causesJSON); err != nil {
		return nil, &store.EventNotFoundError{ID: id}
	}

	visited := map[string]bool{id.Value(): true}
	var result []event.Event

	var causeStrs []string
	json.Unmarshal([]byte(causesJSON), &causeStrs)

	frontier := make([]string, 0)
	for _, c := range causeStrs {
		if c != id.Value() {
			frontier = append(frontier, c)
		}
	}

	for d := 0; d < maxDepth && len(frontier) > 0; d++ {
		nextFrontier := make([]string, 0)
		for _, eid := range frontier {
			if visited[eid] {
				continue
			}
			visited[eid] = true

			ev, err := s.getUnlocked(types.MustEventID(eid))
			if err != nil {
				continue
			}
			result = append(result, ev)
			for _, c := range ev.Causes() {
				if !visited[c.Value()] {
					nextFrontier = append(nextFrontier, c.Value())
				}
			}
		}
		frontier = nextFrontier
	}
	return result, nil
}

func (s *SQLiteStore) Descendants(id types.EventID, maxDepth int) ([]event.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var exists int
	err := s.db.QueryRow("SELECT 1 FROM events WHERE event_id = ?", id.Value()).Scan(&exists)
	if err != nil {
		return nil, &store.EventNotFoundError{ID: id}
	}

	// Build reverse index by querying all events
	rows, err := s.db.Query("SELECT event_id, causes FROM events ORDER BY position ASC")
	if err != nil {
		return nil, &store.StoreUnavailableError{Reason: err.Error()}
	}
	defer rows.Close()

	children := map[string][]string{}
	for rows.Next() {
		var eid, causesJSON string
		rows.Scan(&eid, &causesJSON)
		var causeStrs []string
		json.Unmarshal([]byte(causesJSON), &causeStrs)
		for _, c := range causeStrs {
			if c != eid {
				children[c] = append(children[c], eid)
			}
		}
	}

	visited := map[string]bool{id.Value(): true}
	var result []event.Event
	frontier := children[id.Value()]

	for d := 0; d < maxDepth && len(frontier) > 0; d++ {
		nextFrontier := make([]string, 0)
		for _, eid := range frontier {
			if visited[eid] {
				continue
			}
			visited[eid] = true

			ev, err := s.getUnlocked(types.MustEventID(eid))
			if err != nil {
				continue
			}
			result = append(result, ev)
			for _, child := range children[eid] {
				if !visited[child] {
					nextFrontier = append(nextFrontier, child)
				}
			}
		}
		frontier = nextFrontier
	}
	return result, nil
}

func (s *SQLiteStore) EdgesFrom(id types.ActorID, edgeType event.EdgeType) ([]event.Edge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Simplified: scan edges table
	rows, err := s.db.Query("SELECT * FROM edges WHERE from_actor = ? AND edge_type = ?",
		id.Value(), string(edgeType))
	if err != nil {
		return nil, nil
	}
	defer rows.Close()
	return s.scanEdges(rows), nil
}

func (s *SQLiteStore) EdgesTo(id types.ActorID, edgeType event.EdgeType) ([]event.Edge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query("SELECT * FROM edges WHERE to_actor = ? AND edge_type = ?",
		id.Value(), string(edgeType))
	if err != nil {
		return nil, nil
	}
	defer rows.Close()
	return s.scanEdges(rows), nil
}

func (s *SQLiteStore) EdgeBetween(from types.ActorID, to types.ActorID, edgeType event.EdgeType) (types.Option[event.Edge], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query(
		"SELECT * FROM edges WHERE from_actor = ? AND to_actor = ? AND edge_type = ? ORDER BY created_at DESC LIMIT 1",
		from.Value(), to.Value(), string(edgeType))
	if err != nil {
		return types.None[event.Edge](), nil
	}
	defer rows.Close()

	edges := s.scanEdges(rows)
	if len(edges) == 0 {
		return types.None[event.Edge](), nil
	}
	return types.Some(edges[0]), nil
}

func (s *SQLiteStore) scanEdges(rows *sql.Rows) []event.Edge {
	var result []event.Edge
	for rows.Next() {
		var (
			id        string
			fromActor string
			toActor   string
			edgeType  string
			weight    float64
			direction string
			scope     sql.NullString
			createdAt int64
			expiresAt sql.NullInt64
		)
		if err := rows.Scan(&id, &fromActor, &toActor, &edgeType, &weight,
			&direction, &scope, &createdAt, &expiresAt); err != nil {
			continue
		}

		edgeID, _ := types.NewEdgeID(id)
		w := types.MustWeight(weight)
		ts := types.MustTimestamp(createdAt)

		var scopeOpt types.Option[types.DomainScope]
		if scope.Valid {
			scopeOpt = types.Some(types.MustDomainScope(scope.String))
		}

		var expiresOpt types.Option[types.Timestamp]
		if expiresAt.Valid {
			expiresOpt = types.Some(types.MustTimestamp(expiresAt.Int64))
		}

		dir := event.EdgeDirection(direction)
		et := event.EdgeType(edgeType)
		e, err := event.NewEdge(edgeID, types.MustActorID(fromActor), types.MustActorID(toActor),
			et, w, dir, scopeOpt, nil, ts, expiresOpt)
		if err == nil {
			result = append(result, e)
		}
	}
	return result
}

func (s *SQLiteStore) Count() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var count int
	s.db.QueryRow("SELECT COUNT(*) FROM events").Scan(&count)
	return count, nil
}

func (s *SQLiteStore) VerifyChain() (event.ChainVerifiedContent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	start := time.Now()

	rows, err := s.db.Query("SELECT hash, prev_hash, event_id, event_type, version, timestamp_nanos, source, content, causes, conversation_id, signature FROM events ORDER BY position ASC")
	if err != nil {
		return event.ChainVerifiedContent{Valid: false}, &store.StoreUnavailableError{Reason: err.Error()}
	}
	defer rows.Close()

	var prevHash string
	count := 0
	for rows.Next() {
		var (
			hash, ph, eid, et string
			version           int
			tsNanos           int64
			source, content   string
			causesJSON, cid   string
			sig               []byte
		)
		if err := rows.Scan(&hash, &ph, &eid, &et, &version, &tsNanos, &source, &content, &causesJSON, &cid, &sig); err != nil {
			return event.ChainVerifiedContent{Valid: false, Length: count}, nil
		}

		if count == 0 {
			if ph != types.ZeroHash().Value() && ph != "" {
				return event.ChainVerifiedContent{Valid: false, Length: count}, nil
			}
		} else if ph != prevHash {
			return event.ChainVerifiedContent{Valid: false, Length: count}, nil
		}

		prevHash = hash
		count++
	}

	ns := time.Since(start).Nanoseconds()
	if ns < 0 {
		ns = 0
	}
	dur := types.MustDuration(ns)
	return event.ChainVerifiedContent{Valid: true, Length: count, Duration: dur}, nil
}

func (s *SQLiteStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Close()
}

// Ensure SQLiteStore implements store.Store at compile time.
var _ store.Store = (*SQLiteStore)(nil)

// ContentJSON returns the JSON representation of the event's content.
// This is a helper exposed for the store package to use.
func contentJSON(ev event.Event) string {
	data, _ := json.Marshal(ev.Content())
	return string(data)
}
