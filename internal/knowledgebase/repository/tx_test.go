//go:build integration

package repository

import (
	"Kusuriya/internal/dbx"
	"Kusuriya/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func countProfiles(t *testing.T, repo *KnowledgebaseRepository, causeID int64) int {
	t.Helper()
	var n int
	require.NoError(t, repo.pool.QueryRow(ctx, "SELECT count(*) FROM symptom_profiles WHERE cause_id = $1", causeID).Scan(&n))
	return n
}

func TestTransactions(t *testing.T) {
	repo, _ := newRepo(t)
	seedEvidence(t, repo)
	require.NoError(t, repo.DeleteSymptomProfiles(ctx, 1))

	t.Run("commit without WithTx", func(t *testing.T) {
		assert.ErrorIs(t, repo.Commit(ctx), dbx.ErrTransaction)
		assert.ErrorIs(t, repo.Rollback(ctx), dbx.ErrTransaction)
	})

	t.Run("rollback leaves no rows", func(t *testing.T) {
		txCtx, err := repo.WithTx(ctx)
		require.NoError(t, err)
		_, err = repo.CreateSymptomProfiles(txCtx, []domain.SymptomProfile{{CauseID: 1, SymptomID: 1, Likelihood: 0.5}})
		require.NoError(t, err)
		require.NoError(t, repo.Rollback(txCtx))

		assert.Zero(t, countProfiles(t, repo, 1))
	})

	t.Run("commit persists rows", func(t *testing.T) {
		txCtx, err := repo.WithTx(ctx)
		require.NoError(t, err)
		_, err = repo.CreateSymptomProfiles(txCtx, []domain.SymptomProfile{{CauseID: 1, SymptomID: 1, Likelihood: 0.5}})
		require.NoError(t, err)
		require.NoError(t, repo.Commit(txCtx))

		assert.Equal(t, 1, countProfiles(t, repo, 1))
	})

	t.Run("delete + failing insert inside one tx rolls back the delete too", func(t *testing.T) {
		// this is the replace-profiles flow the service uses
		txCtx, err := repo.WithTx(ctx)
		require.NoError(t, err)
		require.NoError(t, repo.DeleteSymptomProfiles(txCtx, 1))
		_, err = repo.CreateSymptomProfiles(txCtx, []domain.SymptomProfile{
			{CauseID: 1, SymptomID: 2, Likelihood: 0.4},
			{CauseID: 1, SymptomID: 999, Likelihood: 0.4}, // FK violation
		})
		require.Error(t, err)
		require.NoError(t, repo.Rollback(txCtx))

		assert.Equal(t, 1, countProfiles(t, repo, 1), "the row committed earlier survived")
	})
}
