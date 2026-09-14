//go:build integration

package repository

import (
	"Kusuriya/internal/dbx/dbtest"
	"Kusuriya/internal/domain"
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) { os.Exit(dbtest.Run(m)) }

var ctx = context.Background()

func newRepo(t *testing.T) (*CaseRepository, *pgxpool.Pool) {
	t.Helper()
	pool := dbtest.Pool(t)
	dbtest.Truncate(t, pool, "cases")
	return NewCaseRepository(pool), pool
}

func TestCaseLifecycle(t *testing.T) {
	repo, _ := newRepo(t)

	created, err := repo.CreateCase(ctx, domain.Case{Title: "t", Description: "d", Status: "open"})
	require.NoError(t, err)
	assert.Equal(t, &domain.Case{ID: 1, Title: "t", Description: "d", Status: "open"}, created)

	bare, err := repo.CreateCase(ctx, domain.Case{Title: "t2", Status: "open"})
	require.NoError(t, err)
	assert.Equal(t, "", bare.Description, "null description -> empty string")

	got, err := repo.GetCaseByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, created, got)

	closed, err := repo.CloseCaseByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, &domain.Case{ID: 1, Title: "t", Description: "d", Status: "closed"}, closed,
		"only status changes")

	again, err := repo.CloseCaseByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "closed", again.Status, "closing twice is idempotent")

	t.Run("errors", func(t *testing.T) {
		_, err := repo.GetCaseByID(ctx, 999)
		assert.ErrorIs(t, err, domain.ErrNotFound)

		_, err = repo.CloseCaseByID(ctx, 999)
		assert.ErrorIs(t, err, domain.ErrNotFound)

		_, err = repo.CreateCase(ctx, domain.Case{Title: "t", Status: "banana"})
		assert.ErrorIs(t, err, domain.ErrCheckPass, "status outside CHECK")
	})
}

func TestObservations(t *testing.T) {
	repo, pool := newRepo(t)
	c, err := repo.CreateCase(ctx, domain.Case{Title: "t", Status: "open"})
	require.NoError(t, err)

	before := time.Now().Add(-time.Second)
	first, err := repo.CreateObservation(ctx, domain.Observation{CaseID: c.ID, Kind: "symptom", RefID: 3, Note: "n"})
	require.NoError(t, err)
	assert.Equal(t, int64(1), first.ID)
	assert.Equal(t, c.ID, first.CaseID)
	assert.Equal(t, "symptom", first.Kind)
	assert.Equal(t, int64(3), first.RefID)
	assert.Equal(t, "n", first.Note)
	assert.True(t, first.ObservedAt.After(before), "observed_at defaults to now()")

	second, err := repo.CreateObservation(ctx, domain.Observation{CaseID: c.ID, Kind: "factor", RefID: 8})
	require.NoError(t, err)
	assert.Equal(t, "", second.Note)

	got, err := repo.GetObservationByID(ctx, int(c.ID))
	require.NoError(t, err)
	assert.Equal(t, []domain.Observation{*first, *second}, got, "ordered by id")

	t.Run("case without observations returns empty, not nil", func(t *testing.T) {
		other, err := repo.CreateCase(ctx, domain.Case{Title: "empty", Status: "open"})
		require.NoError(t, err)
		got, err := repo.GetObservationByID(ctx, int(other.ID))
		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.Empty(t, got)
	})

	t.Run("unknown case id is not an error, just empty", func(t *testing.T) {
		got, err := repo.GetObservationByID(ctx, 999)
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("kind outside CHECK", func(t *testing.T) {
		_, err := repo.CreateObservation(ctx, domain.Observation{CaseID: c.ID, Kind: "banana", RefID: 1})
		assert.ErrorIs(t, err, domain.ErrCheckPass)
	})

	t.Run("observation for a missing case is an unmapped FK error", func(t *testing.T) {
		_, err := repo.CreateObservation(ctx, domain.Observation{CaseID: 999, Kind: "symptom", RefID: 1})
		require.Error(t, err)
		assert.NotErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("deleting the case cascades to observations", func(t *testing.T) {
		_, err := pool.Exec(ctx, "DELETE FROM cases WHERE id = $1", c.ID)
		require.NoError(t, err)
		got, err := repo.GetObservationByID(ctx, int(c.ID))
		require.NoError(t, err)
		assert.Empty(t, got)
	})
}
