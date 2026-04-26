// Package mysqlstore implements a MySQL-backed Store.
package mysqlstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/store"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

const eventsSchema = `
CREATE TABLE IF NOT EXISTS events (
	position BIGINT AUTO_INCREMENT PRIMARY KEY,
	event_id VARCHAR(100) NOT NULL UNIQUE,
	event_type VARCHAR(200) NOT NULL,
	version INT NOT NULL,
	timestamp_nanos BIGINT NOT NULL,
	source VARCHAR(200) NOT NULL,
	content LONGTEXT NOT NULL,
	causes LONGTEXT NOT NULL,
	conversation_id VARCHAR(200) NOT NULL,
	hash VARCHAR(100) NOT NULL,
	prev_hash VARCHAR(100) NOT NULL,
	signature VARBINARY(64) NOT NULL,
	INDEX idx_events_type (event_type),
	INDEX idx_events_source (source),
	INDEX idx_events_conversation (conversation_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`

const causesSchema = `
CREATE TABLE IF NOT EXISTS event_causes (
	event_id VARCHAR(100) NOT NULL,
	cause_id VARCHAR(100) NOT NULL,
	PRIMARY KEY (event_id, cause_id),
	INDEX idx_event_causes_cause (cause_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`

const edgesSchema = `
CREATE TABLE IF NOT EXISTS edges (
	id VARCHAR(100) PRIMARY KEY,
	from_actor VARCHAR(200) NOT NULL,
	to_actor VARCHAR(200) NOT NULL,
	edge_type VARCHAR(100) NOT NULL,
	weight DOUBLE NOT NULL,
	direction VARCHAR(50) NOT NULL,
	scope VARCHAR(200),
	created_at_nanos BIGINT NOT NULL,
	expires_at_nanos BIGINT,
	INDEX idx_edges_from_type (from_actor, edge_type),
	INDEX idx_edges_to_type (to_actor, edge_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`

// MySQLStore implements store.Store backed by MySQL.
type MySQLStore struct {
	mu sync.Mutex
	db *sql.DB
}

// New creates a MySQLStore connected to the given MySQL instance.
// The dsn should be in go-sql-driver/mysql format, e.g. "user:pass@tcp(localhost:3306)/eventgraph?parseTime=true".
// It creates the schema if it doesn't exist.
func New(dsn string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, &store.StoreUnavailableError{Reason: fmt.Sprintf("open: %v", err)}
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, &store.StoreUnavailableError{Reason: fmt.Sprintf("ping: %v", err)}
	}

	for _, ddl := range []string{eventsSchema, causesSchema, edgesSchema} {
		if _, err := db.Exec(ddl); err != nil {
			db.Close()
			return nil, &store.StoreUnavailableError{Reason: fmt.Sprintf("schema: %v", err)}
		}
	}

	return &MySQLStore{db: db}, nil
}

func (s *MySQLStore) Append(ev event.Event) (event.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return event.Event{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("begin tx: %v", err)}
	}
	defer tx.Rollback()

	// Idempotency: if same ID exists, verify hash matches and return it.
	var existingHash string
	err = tx.QueryRowContext(ctx, "SELECT hash FROM events WHERE event_id = ?", ev.ID().Value()).Scan(&existingHash)
	if err == nil {
		if existingHash != ev.Hash().Value() {
			return event.Event{}, &store.HashMismatchError{
				EventID:  ev.ID(),
				Computed: ev.Hash(),
				Stored:   types.MustHash(existingHash),
			}
		}
		tx.Commit()
		return s.getEventUnlocked(ctx, ev.ID())
	}
	if err != sql.ErrNoRows {
		return event.Event{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("check existing: %v", err)}
	}

	// Verify PrevHash matches chain head.
	var headHash string
	var headExists bool
	err = tx.QueryRowContext(ctx, "SELECT hash FROM events ORDER BY position DESC LIMIT 1").Scan(&headHash)
	if err == sql.ErrNoRows {
		headExists = false
	} else if err != nil {
		return event.Event{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("head query: %v", err)}
	} else {
		headExists = true
	}

	if !headExists {
		if ev.PrevHash() != types.ZeroHash() {
			return event.Event{}, &store.ChainIntegrityViolationError{
				Position: 0,
				Expected: types.ZeroHash(),
				Actual:   ev.PrevHash(),
			}
		}
	} else {
		if ev.PrevHash().Value() != headHash {
			var count int
			tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM events").Scan(&count)
			return event.Event{}, &store.ChainIntegrityViolationError{
				Position: count,
				Expected: types.MustHash(headHash),
				Actual:   ev.PrevHash(),
			}
		}
	}

	// Recompute hash and verify.
	canonical := event.CanonicalForm(ev)
	computed, err := event.ComputeHash(canonical)
	if err != nil {
		return event.Event{}, err
	}
	if computed != ev.Hash() {
		return event.Event{}, &store.HashMismatchError{
			EventID:  ev.ID(),
			Computed: computed,
			Stored:   ev.Hash(),
		}
	}

	// Verify causal predecessors exist.
	if !ev.IsBootstrap() {
		for _, causeID := range ev.Causes() {
			var exists bool
			err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM events WHERE event_id = ?)", causeID.Value()).Scan(&exists)
			if err != nil {
				return event.Event{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("cause check: %v", err)}
			}
			if !exists {
				return event.Event{}, &store.CausalLinkMissingError{
					EventID:      ev.ID(),
					MissingCause: causeID,
				}
			}
		}
	}

	// Validate edge content before inserting.
	var edge event.Edge
	var hasEdge bool
	if ev.Type() == event.EventTypeEdgeCreated {
		ec, ok := ev.Content().(event.EdgeCreatedContent)
		if !ok {
			return event.Event{}, &store.EdgeIndexError{
				EventID: ev.ID(),
				Reason:  fmt.Sprintf("wrong content type %T", ev.Content()),
			}
		}
		edgeID, edgeIDErr := types.NewEdgeID(ev.ID().Value())
		if edgeIDErr != nil {
			return event.Event{}, &store.EdgeIndexError{
				EventID: ev.ID(),
				Reason:  fmt.Sprintf("derive edge ID: %v", edgeIDErr),
			}
		}
		var newEdgeErr error
		edge, newEdgeErr = event.NewEdge(
			edgeID, ec.From, ec.To, ec.EdgeType, ec.Weight, ec.Direction,
			ec.Scope, nil, ev.Timestamp(), ec.ExpiresAt,
		)
		if newEdgeErr != nil {
			return event.Event{}, &store.EdgeIndexError{
				EventID: ev.ID(),
				Reason:  fmt.Sprintf("construct edge: %v", newEdgeErr),
			}
		}
		hasEdge = true
	}

	var supersededPrevID string
	if ev.Type() == event.EventTypeEdgeSuperseded {
		ec, ok := ev.Content().(event.EdgeSupersededContent)
		if !ok {
			return event.Event{}, &store.EdgeIndexError{
				EventID: ev.ID(),
				Reason:  fmt.Sprintf("wrong content type for edge.superseded: %T", ev.Content()),
			}
		}
		supersededPrevID = ec.PreviousEdge.Value()
	}

	// Serialize content to JSON.
	contentJSON, err := json.Marshal(ev.Content())
	if err != nil {
		return event.Event{}, fmt.Errorf("marshal content: %w", err)
	}

	// Serialize causes as JSON array of strings.
	causeIDs := ev.Causes()
	causeStrings := make([]string, len(causeIDs))
	for i, c := range causeIDs {
		causeStrings[i] = c.Value()
	}
	causesJSON, err := json.Marshal(causeStrings)
	if err != nil {
		return event.Event{}, fmt.Errorf("marshal causes: %w", err)
	}

	// Insert the event.
	_, err = tx.ExecContext(ctx,
		`INSERT INTO events (event_id, event_type, version, timestamp_nanos, source, content, causes, conversation_id, hash, prev_hash, signature)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		ev.ID().Value(), ev.Type().Value(), ev.Version(), ev.Timestamp().UnixNano(),
		ev.Source().Value(), string(contentJSON), string(causesJSON),
		ev.ConversationID().Value(), ev.Hash().Value(), ev.PrevHash().Value(),
		ev.Signature().Bytes(),
	)
	if err != nil {
		return event.Event{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("insert event: %v", err)}
	}

	// Insert causes into the join table.
	for _, causeID := range causeIDs {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO event_causes (event_id, cause_id) VALUES (?, ?)",
			ev.ID().Value(), causeID.Value(),
		)
		if err != nil {
			return event.Event{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("insert cause: %v", err)}
		}
	}

	// Insert edge if applicable.
	if hasEdge {
		var scopeVal *string
		if edge.Scope().IsSome() {
			s := edge.Scope().Unwrap().Value()
			scopeVal = &s
		}
		var expiresNanos *int64
		if edge.ExpiresAt().IsSome() {
			n := edge.ExpiresAt().Unwrap().UnixNano()
			expiresNanos = &n
		}
		_, err = tx.ExecContext(ctx,
			`INSERT INTO edges (id, from_actor, to_actor, edge_type, weight, direction, scope, created_at_nanos, expires_at_nanos)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			edge.ID().Value(), edge.From().Value(), edge.To().Value(),
			string(edge.Type()), edge.Weight().Value(), string(edge.Direction()),
			scopeVal, edge.CreatedAt().UnixNano(), expiresNanos,
		)
		if err != nil {
			return event.Event{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("insert edge: %v", err)}
		}
	}

	// Handle edge supersession.
	if supersededPrevID != "" {
		_, err = tx.ExecContext(ctx, "DELETE FROM edges WHERE id = ?", supersededPrevID)
		if err != nil {
			return event.Event{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("delete superseded edge: %v", err)}
		}
	}

	if err := tx.Commit(); err != nil {
		return event.Event{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("commit: %v", err)}
	}
	return ev, nil
}

func (s *MySQLStore) Get(id types.EventID) (event.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getEventUnlocked(context.Background(), id)
}

func (s *MySQLStore) Head() (types.Option[event.Event], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := context.Background()
	row := s.db.QueryRowContext(ctx,
		`SELECT event_id, version, event_type, timestamp_nanos, source, conversation_id,
		        hash, prev_hash, signature, content, causes
		 FROM events ORDER BY position DESC LIMIT 1`)
	ev, err := scanEvent(row)
	if err == sql.ErrNoRows {
		return types.None[event.Event](), nil
	}
	if err != nil {
		return types.None[event.Event](), err
	}
	return types.Some(ev), nil
}

func (s *MySQLStore) Recent(limit int, after types.Option[types.Cursor]) (types.Page[event.Event], error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit <= 0 {
		limit = 100
	}
	return s.paginateReverse(context.Background(), "", nil, limit, after)
}

func (s *MySQLStore) ByType(eventType types.EventType, limit int, after types.Option[types.Cursor]) (types.Page[event.Event], error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit <= 0 {
		limit = 100
	}
	return s.paginateReverse(context.Background(), "event_type = ?", []any{eventType.Value()}, limit, after)
}

func (s *MySQLStore) BySource(source types.ActorID, limit int, after types.Option[types.Cursor]) (types.Page[event.Event], error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit <= 0 {
		limit = 100
	}
	return s.paginateReverse(context.Background(), "source = ?", []any{source.Value()}, limit, after)
}

func (s *MySQLStore) ByConversation(id types.ConversationID, limit int, after types.Option[types.Cursor]) (types.Page[event.Event], error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit <= 0 {
		limit = 100
	}
	return s.paginateReverse(context.Background(), "conversation_id = ?", []any{id.Value()}, limit, after)
}

func (s *MySQLStore) Since(afterID types.EventID, limit int) (types.Page[event.Event], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := context.Background()
	if limit <= 0 {
		limit = 100
	}

	// Find the position of the afterID event.
	var afterPos int64
	err := s.db.QueryRowContext(ctx, "SELECT position FROM events WHERE event_id = ?", afterID.Value()).Scan(&afterPos)
	if err == sql.ErrNoRows {
		return types.Page[event.Event]{}, &store.EventNotFoundError{ID: afterID}
	}
	if err != nil {
		return types.Page[event.Event]{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("since lookup: %v", err)}
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT event_id, version, event_type, timestamp_nanos, source, conversation_id,
		        hash, prev_hash, signature, content, causes
		 FROM events WHERE position > ? ORDER BY position ASC LIMIT ?`,
		afterPos, limit)
	if err != nil {
		return types.Page[event.Event]{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("since query: %v", err)}
	}
	defer rows.Close()

	var items []event.Event
	for rows.Next() {
		ev, err := scanEventFromRows(rows)
		if err != nil {
			return types.Page[event.Event]{}, err
		}
		items = append(items, ev)
	}

	// Check hasMore.
	var totalAfter int
	s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM events WHERE position > ?", afterPos).Scan(&totalAfter)
	hasMore := totalAfter > len(items)

	return types.NewPage(items, types.None[types.Cursor](), hasMore), nil
}

func (s *MySQLStore) Ancestors(id types.EventID, maxDepth int) ([]event.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := context.Background()

	// Verify the event exists.
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM events WHERE event_id = ?)", id.Value()).Scan(&exists)
	if err != nil {
		return nil, &store.StoreUnavailableError{Reason: fmt.Sprintf("ancestors check: %v", err)}
	}
	if !exists {
		return nil, &store.EventNotFoundError{ID: id}
	}

	// MySQL 8+ supports recursive CTEs.
	rows, err := s.db.QueryContext(ctx,
		`WITH RECURSIVE ancestors AS (
			SELECT cause_id AS id, 1 AS depth
			FROM event_causes WHERE event_id = ?
		  UNION
			SELECT ec.cause_id, a.depth + 1
			FROM ancestors a
			JOIN event_causes ec ON ec.event_id = a.id
			WHERE a.depth < ?
		)
		SELECT DISTINCT e.event_id, e.version, e.event_type, e.timestamp_nanos, e.source,
		       e.conversation_id, e.hash, e.prev_hash, e.signature, e.content, e.causes
		FROM ancestors a
		JOIN events e ON e.event_id = a.id`,
		id.Value(), maxDepth)
	if err != nil {
		return nil, &store.StoreUnavailableError{Reason: fmt.Sprintf("ancestors query: %v", err)}
	}
	defer rows.Close()

	var result []event.Event
	for rows.Next() {
		ev, err := scanEventFromRows(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, ev)
	}
	return result, nil
}

func (s *MySQLStore) Descendants(id types.EventID, maxDepth int) ([]event.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := context.Background()

	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM events WHERE event_id = ?)", id.Value()).Scan(&exists)
	if err != nil {
		return nil, &store.StoreUnavailableError{Reason: fmt.Sprintf("descendants check: %v", err)}
	}
	if !exists {
		return nil, &store.EventNotFoundError{ID: id}
	}

	rows, err := s.db.QueryContext(ctx,
		`WITH RECURSIVE descendants AS (
			SELECT event_id AS id, 1 AS depth
			FROM event_causes WHERE cause_id = ?
		  UNION
			SELECT ec.event_id, d.depth + 1
			FROM descendants d
			JOIN event_causes ec ON ec.cause_id = d.id
			WHERE d.depth < ?
		)
		SELECT DISTINCT e.event_id, e.version, e.event_type, e.timestamp_nanos, e.source,
		       e.conversation_id, e.hash, e.prev_hash, e.signature, e.content, e.causes
		FROM descendants d
		JOIN events e ON e.event_id = d.id`,
		id.Value(), maxDepth)
	if err != nil {
		return nil, &store.StoreUnavailableError{Reason: fmt.Sprintf("descendants query: %v", err)}
	}
	defer rows.Close()

	var result []event.Event
	for rows.Next() {
		ev, err := scanEventFromRows(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, ev)
	}
	return result, nil
}

func (s *MySQLStore) EdgesFrom(id types.ActorID, edgeType event.EdgeType) ([]event.Edge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.QueryContext(context.Background(),
		`SELECT id, from_actor, to_actor, edge_type, weight, direction, scope, created_at_nanos, expires_at_nanos
		 FROM edges WHERE from_actor = ? AND edge_type = ?`,
		id.Value(), string(edgeType))
	if err != nil {
		return nil, &store.StoreUnavailableError{Reason: fmt.Sprintf("edges from: %v", err)}
	}
	defer rows.Close()
	return scanEdges(rows)
}

func (s *MySQLStore) EdgesTo(id types.ActorID, edgeType event.EdgeType) ([]event.Edge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.QueryContext(context.Background(),
		`SELECT id, from_actor, to_actor, edge_type, weight, direction, scope, created_at_nanos, expires_at_nanos
		 FROM edges WHERE to_actor = ? AND edge_type = ?`,
		id.Value(), string(edgeType))
	if err != nil {
		return nil, &store.StoreUnavailableError{Reason: fmt.Sprintf("edges to: %v", err)}
	}
	defer rows.Close()
	return scanEdges(rows)
}

func (s *MySQLStore) EdgeBetween(from types.ActorID, to types.ActorID, edgeType event.EdgeType) (types.Option[event.Edge], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.QueryContext(context.Background(),
		`SELECT id, from_actor, to_actor, edge_type, weight, direction, scope, created_at_nanos, expires_at_nanos
		 FROM edges WHERE from_actor = ? AND to_actor = ? AND edge_type = ?
		 ORDER BY created_at_nanos DESC LIMIT 1`,
		from.Value(), to.Value(), string(edgeType))
	if err != nil {
		return types.None[event.Edge](), &store.StoreUnavailableError{Reason: fmt.Sprintf("edge between: %v", err)}
	}
	defer rows.Close()

	edges, err := scanEdges(rows)
	if err != nil {
		return types.None[event.Edge](), err
	}
	if len(edges) == 0 {
		return types.None[event.Edge](), nil
	}
	return types.Some(edges[0]), nil
}

func (s *MySQLStore) Count() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var count int
	err := s.db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM events").Scan(&count)
	if err != nil {
		return 0, &store.StoreUnavailableError{Reason: fmt.Sprintf("count: %v", err)}
	}
	return count, nil
}

func (s *MySQLStore) VerifyChain() (event.ChainVerifiedContent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := context.Background()
	start := time.Now()

	rows, err := s.db.QueryContext(ctx,
		`SELECT event_id, version, event_type, timestamp_nanos, source, conversation_id,
		        hash, prev_hash, signature, content, causes
		 FROM events ORDER BY position ASC`)
	if err != nil {
		return event.ChainVerifiedContent{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("verify chain: %v", err)}
	}
	defer rows.Close()

	var prevHash string
	i := 0
	for rows.Next() {
		ev, err := scanEventFromRows(rows)
		if err != nil {
			return event.ChainVerifiedContent{Valid: false, Length: i}, nil
		}

		if i == 0 {
			if !ev.IsBootstrap() {
				return event.ChainVerifiedContent{Valid: false, Length: i}, nil
			}
			if ev.PrevHash() != types.ZeroHash() {
				return event.ChainVerifiedContent{Valid: false, Length: i}, nil
			}
		} else {
			if ev.PrevHash().Value() != prevHash {
				return event.ChainVerifiedContent{Valid: false, Length: i}, nil
			}
		}

		canonical := event.CanonicalForm(ev)
		computed, cerr := event.ComputeHash(canonical)
		if cerr != nil {
			return event.ChainVerifiedContent{Valid: false, Length: i}, nil
		}
		if computed != ev.Hash() {
			return event.ChainVerifiedContent{Valid: false, Length: i}, nil
		}

		prevHash = ev.Hash().Value()
		i++
	}

	ns := time.Since(start).Nanoseconds()
	if ns < 0 {
		ns = 0
	}
	dur := types.MustDuration(ns)
	return event.ChainVerifiedContent{
		Valid:    true,
		Length:   i,
		Duration: dur,
	}, nil
}

func (s *MySQLStore) Close() error {
	return s.db.Close()
}

// Truncate removes all data from the store tables. Used for testing.
func (s *MySQLStore) Truncate(ctx context.Context) error {
	// MySQL requires disabling FK checks to truncate tables with foreign key constraints.
	if _, err := s.db.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		return err
	}
	for _, table := range []string{"event_causes", "edges", "events"} {
		if _, err := s.db.ExecContext(ctx, "TRUNCATE TABLE "+table); err != nil {
			s.db.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 1")
			return err
		}
	}
	_, err := s.db.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 1")
	return err
}

// --- internal helpers ---

// getEventUnlocked fetches a single event by ID. Caller must hold the mutex.
func (s *MySQLStore) getEventUnlocked(ctx context.Context, id types.EventID) (event.Event, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT event_id, version, event_type, timestamp_nanos, source, conversation_id,
		        hash, prev_hash, signature, content, causes
		 FROM events WHERE event_id = ?`, id.Value())
	ev, err := scanEvent(row)
	if err == sql.ErrNoRows {
		return event.Event{}, &store.EventNotFoundError{ID: id}
	}
	return ev, err
}

// paginateReverse returns events in reverse chronological order (most recent first).
func (s *MySQLStore) paginateReverse(ctx context.Context, filterClause string, filterArgs []any, limit int, after types.Option[types.Cursor]) (types.Page[event.Event], error) {
	args := make([]any, 0, len(filterArgs)+2)

	var whereParts []string

	if filterClause != "" {
		whereParts = append(whereParts, filterClause)
		args = append(args, filterArgs...)
	}

	if after.IsSome() {
		cursor := after.Unwrap()
		var cursorPos int64
		err := s.db.QueryRowContext(ctx, "SELECT position FROM events WHERE event_id = ?", cursor.Value()).Scan(&cursorPos)
		if err == sql.ErrNoRows {
			return types.NewPage[event.Event](nil, types.None[types.Cursor](), false),
				&store.InvalidCursorError{Cursor: cursor.Value()}
		}
		if err != nil {
			return types.Page[event.Event]{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("cursor lookup: %v", err)}
		}
		whereParts = append(whereParts, "position < ?")
		args = append(args, cursorPos)
	}

	whereSQL := ""
	if len(whereParts) > 0 {
		whereSQL = "WHERE "
		for i, part := range whereParts {
			if i > 0 {
				whereSQL += " AND "
			}
			whereSQL += part
		}
	}

	// Query limit+1 to determine hasMore.
	query := fmt.Sprintf(
		`SELECT event_id, version, event_type, timestamp_nanos, source, conversation_id,
		        hash, prev_hash, signature, content, causes
		 FROM events %s ORDER BY position DESC LIMIT ?`, whereSQL)
	args = append(args, limit+1)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return types.Page[event.Event]{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("paginate query: %v", err)}
	}
	defer rows.Close()

	var items []event.Event
	for rows.Next() {
		ev, err := scanEventFromRows(rows)
		if err != nil {
			return types.Page[event.Event]{}, err
		}
		items = append(items, ev)
	}

	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}

	var cursorOpt types.Option[types.Cursor]
	if hasMore && len(items) > 0 {
		c := types.MustCursor(items[len(items)-1].ID().Value())
		cursorOpt = types.Some(c)
	}

	return types.NewPage(items, cursorOpt, hasMore), nil
}

// scanner abstracts over *sql.Row and *sql.Rows for scanning.
type scanner interface {
	Scan(dest ...any) error
}

// scanEvent scans a single row into an Event.
func scanEvent(row *sql.Row) (event.Event, error) {
	var (
		id             string
		version        int
		eventType      string
		timestampNanos int64
		source         string
		conversationID string
		hash           string
		prevHash       string
		signature      []byte
		contentStr     string
		causesStr      string
	)
	err := row.Scan(&id, &version, &eventType, &timestampNanos, &source,
		&conversationID, &hash, &prevHash, &signature, &contentStr, &causesStr)
	if err != nil {
		return event.Event{}, err
	}
	return reconstructEvent(id, version, eventType, timestampNanos, source,
		conversationID, hash, prevHash, signature, []byte(contentStr), causesStr)
}

// scanEventFromRows scans the current row from sql.Rows into an Event.
func scanEventFromRows(rows *sql.Rows) (event.Event, error) {
	var (
		id             string
		version        int
		eventType      string
		timestampNanos int64
		source         string
		conversationID string
		hash           string
		prevHash       string
		signature      []byte
		contentStr     string
		causesStr      string
	)
	err := rows.Scan(&id, &version, &eventType, &timestampNanos, &source,
		&conversationID, &hash, &prevHash, &signature, &contentStr, &causesStr)
	if err != nil {
		return event.Event{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("scan event: %v", err)}
	}
	return reconstructEvent(id, version, eventType, timestampNanos, source,
		conversationID, hash, prevHash, signature, []byte(contentStr), causesStr)
}

// reconstructEvent rebuilds an Event from database columns.
func reconstructEvent(
	id string, version int, eventType string, timestampNanos int64,
	source, conversationID, hash, prevHash string,
	signature, contentJSON []byte, causesStr string,
) (event.Event, error) {
	evID := types.MustEventID(id)
	evType := types.MustEventType(eventType)
	ts := types.NewTimestamp(time.Unix(0, timestampNanos))
	src := types.MustActorID(source)
	convID := types.MustConversationID(conversationID)
	h := types.MustHash(hash)
	ph := types.MustHash(prevHash)
	sig := types.MustSignature(signature)

	content, err := unmarshalContent(eventType, contentJSON)
	if err != nil {
		return event.Event{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("unmarshal content: %v", err)}
	}

	// Deserialize causes from JSON array of strings.
	var causeStrings []string
	if err := json.Unmarshal([]byte(causesStr), &causeStrings); err != nil {
		return event.Event{}, &store.StoreUnavailableError{Reason: fmt.Sprintf("unmarshal causes: %v", err)}
	}
	causes := make([]types.EventID, len(causeStrings))
	for i, cs := range causeStrings {
		causes[i] = types.MustEventID(cs)
	}

	if evType == event.EventTypeSystemBootstrapped {
		bc, ok := content.(event.BootstrapContent)
		if !ok {
			return event.Event{}, &store.StoreUnavailableError{Reason: "bootstrap content type mismatch"}
		}
		return event.NewBootstrapEvent(version, evID, evType, ts, src, bc, convID, h, sig), nil
	}

	return event.NewEvent(version, evID, evType, ts, src, content, causes, convID, h, ph, sig), nil
}

// unmarshalContent deserializes JSON into the correct EventContent type.
func unmarshalContent(eventType string, data []byte) (event.EventContent, error) {
	switch eventType {
	case "system.bootstrapped":
		var c event.BootstrapContent
		return c, json.Unmarshal(data, &c)
	case "trust.updated":
		var c event.TrustUpdatedContent
		return c, json.Unmarshal(data, &c)
	case "trust.score":
		var c event.TrustScoreContent
		return c, json.Unmarshal(data, &c)
	case "trust.decayed":
		var c event.TrustDecayedContent
		return c, json.Unmarshal(data, &c)
	case "authority.requested":
		var c event.AuthorityRequestContent
		return c, json.Unmarshal(data, &c)
	case "authority.resolved":
		var c event.AuthorityResolvedContent
		return c, json.Unmarshal(data, &c)
	case "authority.delegated":
		var c event.AuthorityDelegatedContent
		return c, json.Unmarshal(data, &c)
	case "authority.revoked":
		var c event.AuthorityRevokedContent
		return c, json.Unmarshal(data, &c)
	case "authority.timeout":
		var c event.AuthorityTimeoutContent
		return c, json.Unmarshal(data, &c)
	case "actor.registered":
		var c event.ActorRegisteredContent
		return c, json.Unmarshal(data, &c)
	case "actor.suspended":
		var c event.ActorSuspendedContent
		return c, json.Unmarshal(data, &c)
	case "actor.memorial":
		var c event.ActorMemorialContent
		return c, json.Unmarshal(data, &c)
	case "edge.created":
		var c event.EdgeCreatedContent
		return c, json.Unmarshal(data, &c)
	case "edge.superseded":
		var c event.EdgeSupersededContent
		return c, json.Unmarshal(data, &c)
	case "violation.detected":
		var c event.ViolationDetectedContent
		return c, json.Unmarshal(data, &c)
	case "chain.verified":
		var c event.ChainVerifiedContent
		return c, json.Unmarshal(data, &c)
	case "chain.broken":
		var c event.ChainBrokenContent
		return c, json.Unmarshal(data, &c)
	case "clock.tick":
		var c event.ClockTickContent
		return c, json.Unmarshal(data, &c)
	case "health.report":
		var c event.HealthReportContent
		return c, json.Unmarshal(data, &c)
	case "decision.branch.proposed":
		var c event.BranchProposedContent
		return c, json.Unmarshal(data, &c)
	case "decision.branch.inserted":
		var c event.BranchInsertedContent
		return c, json.Unmarshal(data, &c)
	case "decision.cost.report":
		var c event.CostReportContent
		return c, json.Unmarshal(data, &c)
	case "grammar.emit":
		var c event.GrammarEmitContent
		return c, json.Unmarshal(data, &c)
	case "grammar.respond":
		var c event.GrammarRespondContent
		return c, json.Unmarshal(data, &c)
	case "grammar.derive":
		var c event.GrammarDeriveContent
		return c, json.Unmarshal(data, &c)
	case "grammar.extend":
		var c event.GrammarExtendContent
		return c, json.Unmarshal(data, &c)
	case "grammar.retract":
		var c event.GrammarRetractContent
		return c, json.Unmarshal(data, &c)
	case "grammar.annotate":
		var c event.GrammarAnnotateContent
		return c, json.Unmarshal(data, &c)
	case "grammar.merge":
		var c event.GrammarMergeContent
		return c, json.Unmarshal(data, &c)
	case "grammar.consent":
		var c event.GrammarConsentContent
		return c, json.Unmarshal(data, &c)
	case "egip.hello.sent":
		var c event.EGIPHelloSentContent
		return c, json.Unmarshal(data, &c)
	case "egip.hello.received":
		var c event.EGIPHelloReceivedContent
		return c, json.Unmarshal(data, &c)
	case "egip.message.sent":
		var c event.EGIPMessageSentContent
		return c, json.Unmarshal(data, &c)
	case "egip.message.received":
		var c event.EGIPMessageReceivedContent
		return c, json.Unmarshal(data, &c)
	case "egip.receipt.sent":
		var c event.EGIPReceiptSentContent
		return c, json.Unmarshal(data, &c)
	case "egip.receipt.received":
		var c event.EGIPReceiptReceivedContent
		return c, json.Unmarshal(data, &c)
	case "egip.proof.requested":
		var c event.EGIPProofRequestedContent
		return c, json.Unmarshal(data, &c)
	case "egip.proof.received":
		var c event.EGIPProofReceivedContent
		return c, json.Unmarshal(data, &c)
	case "egip.treaty.proposed":
		var c event.EGIPTreatyProposedContent
		return c, json.Unmarshal(data, &c)
	case "egip.treaty.active":
		var c event.EGIPTreatyActiveContent
		return c, json.Unmarshal(data, &c)
	case "egip.trust.updated":
		var c event.EGIPTrustUpdatedContent
		return c, json.Unmarshal(data, &c)
	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}
}

// scanEdges scans rows into Edge slices.
func scanEdges(rows *sql.Rows) ([]event.Edge, error) {
	var result []event.Edge
	for rows.Next() {
		var (
			id           string
			fromActor    string
			toActor      string
			edgeType     string
			weight       float64
			direction    string
			scope        *string
			createdNanos int64
			expiresNanos *int64
		)
		if err := rows.Scan(&id, &fromActor, &toActor, &edgeType, &weight, &direction, &scope, &createdNanos, &expiresNanos); err != nil {
			return nil, &store.StoreUnavailableError{Reason: fmt.Sprintf("scan edge: %v", err)}
		}

		edgeID := types.MustEdgeID(id)
		from := types.MustActorID(fromActor)
		to := types.MustActorID(toActor)
		et := event.EdgeType(edgeType)
		w := types.MustWeight(weight)
		dir := event.EdgeDirection(direction)

		var scopeOpt types.Option[types.DomainScope]
		if scope != nil {
			scopeOpt = types.Some(types.MustDomainScope(*scope))
		}

		createdAt := types.NewTimestamp(time.Unix(0, createdNanos))
		var expiresOpt types.Option[types.Timestamp]
		if expiresNanos != nil {
			expiresOpt = types.Some(types.NewTimestamp(time.Unix(0, *expiresNanos)))
		}

		edge, err := event.NewEdge(edgeID, from, to, et, w, dir, scopeOpt, nil, createdAt, expiresOpt)
		if err != nil {
			return nil, &store.StoreUnavailableError{Reason: fmt.Sprintf("reconstruct edge: %v", err)}
		}
		result = append(result, edge)
	}
	return result, nil
}
