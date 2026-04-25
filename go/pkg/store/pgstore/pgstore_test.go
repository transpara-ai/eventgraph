package pgstore_test

import (
	"context"
	"os"
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/store"
	"github.com/transpara-ai/eventgraph/go/pkg/store/pgstore"
	"github.com/transpara-ai/eventgraph/go/pkg/store/storetest"
)

func TestPostgresConformance(t *testing.T) {
	connStr := os.Getenv("EVENTGRAPH_POSTGRES_URL")
	if connStr == "" {
		t.Skip("EVENTGRAPH_POSTGRES_URL not set; skipping PostgresStore conformance tests")
	}

	ctx := context.Background()

	storetest.RunConformanceSuite(t, func() store.Store {
		s, err := pgstore.NewPostgresStore(ctx, connStr)
		if err != nil {
			t.Fatalf("NewPostgresStore: %v", err)
		}
		// Truncate to start fresh for each sub-test.
		if err := s.Truncate(ctx); err != nil {
			t.Fatalf("Truncate: %v", err)
		}
		return s
	})
}
