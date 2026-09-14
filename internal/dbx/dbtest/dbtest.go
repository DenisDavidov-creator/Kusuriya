//go:build integration

// Package dbtest starts a throwaway Postgres in Docker for repository tests.
// Use from a package's TestMain:
//
//	func TestMain(m *testing.M) { os.Exit(dbtest.Run(m)) }
//
// then in each test: pool := dbtest.Pool(t).
// Run with: go test -tags integration ./...
package dbtest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var connString string

// Run starts Postgres with all *.up.sql migrations applied, runs the tests, and stops it.
func Run(m *testing.M) int {
	ctx := context.Background()

	migrations, err := filepath.Glob(filepath.Join(repoRoot(), "db", "migrations", "*.up.sql"))
	if err != nil || len(migrations) == 0 {
		fmt.Fprintln(os.Stderr, "dbtest: no migrations found:", err)
		return 1
	}

	c, err := postgres.Run(ctx, "postgres:16",
		postgres.WithDatabase("kusuriya_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.WithInitScripts(migrations...), // applied in filename order
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dbtest: start postgres:", err)
		return 1
	}
	defer func() { _ = testcontainers.TerminateContainer(c) }()

	connString, err = c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintln(os.Stderr, "dbtest: connection string:", err)
		return 1
	}
	return m.Run()
}

// Pool returns a connection pool to the test database, closed when the test ends.
func Pool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	require.NotEmpty(t, connString, "dbtest.Pool called without dbtest.Run in TestMain")
	pool, err := pgxpool.New(context.Background(), connString)
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	return pool
}

// Truncate empties the given tables (and anything referencing them) and resets identities,
// so every test starts from ID 1 with no rows.
func Truncate(t *testing.T, pool *pgxpool.Pool, tables ...string) {
	t.Helper()
	_, err := pool.Exec(context.Background(),
		"TRUNCATE "+strings.Join(tables, ", ")+" RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

// repoRoot walks up from the working directory until it finds go.mod.
func repoRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}
