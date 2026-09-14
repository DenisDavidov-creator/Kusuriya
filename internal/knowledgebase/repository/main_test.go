//go:build integration

package repository

import (
	"Kusuriya/internal/dbx/dbtest"
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMain(m *testing.M) { os.Exit(dbtest.Run(m)) }

var ctx = context.Background()

// newRepo returns a repository over an empty knowledge base.
func newRepo(t *testing.T) (*KnowledgebaseRepository, *pgxpool.Pool) {
	t.Helper()
	pool := dbtest.Pool(t)
	dbtest.Truncate(t, pool, "remedies", "causes", "factors", "categories", "symptoms")
	return NewKnowledgebaseRepository(pool), pool
}

func f(v float64) *float64 { return &v }
func i64(v int64) *int64   { return &v }
func str(v string) *string { return &v }
