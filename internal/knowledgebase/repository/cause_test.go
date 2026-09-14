//go:build integration

package repository

import (
	"Kusuriya/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCauseCRUD(t *testing.T) {
	repo, _ := newRepo(t)

	remedy, err := repo.CreateRemedy(ctx, domain.Remedy{Name: "antidote"})
	require.NoError(t, err)

	// create: every column round-trips, including the nullable ones
	created, err := repo.CreateCause(ctx, domain.Cause{
		Code: "c1", Name: "n1", Type: "poison", Prior: 0.25, Description: "d", RemedyID: &remedy.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, &domain.Cause{ID: 1, Code: "c1", Name: "n1", Type: "poison", Prior: 0.25, Description: "d", RemedyID: &remedy.ID}, created)

	// nullable columns come back as zero values / nil, not as errors
	bare, err := repo.CreateCause(ctx, domain.Cause{Code: "c2", Name: "n2", Type: "disease", Prior: 0.5})
	require.NoError(t, err)
	assert.Equal(t, "", bare.Description)
	assert.Nil(t, bare.RemedyID)

	got, err := repo.GetCauseByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, created, got)

	all, err := repo.GetCauses(ctx)
	require.NoError(t, err)
	assert.Equal(t, []domain.Cause{*created, *bare}, all, "ordered by id")

	updated, err := repo.UpdateCause(ctx, domain.Cause{ID: 1, Code: "c1", Name: "renamed", Type: "disease", Prior: 0.9})
	require.NoError(t, err)
	assert.Equal(t, &domain.Cause{ID: 1, Code: "c1", Name: "renamed", Type: "disease", Prior: 0.9}, updated,
		"update clears remedy and description when not provided")

	require.NoError(t, repo.DeleteCause(ctx, 1))
	_, err = repo.GetCauseByID(ctx, 1)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestCauseConstraints(t *testing.T) {
	repo, _ := newRepo(t)
	_, err := repo.CreateCause(ctx, domain.Cause{Code: "c1", Name: "n", Type: "poison", Prior: 0.5})
	require.NoError(t, err)

	t.Run("duplicate code", func(t *testing.T) {
		_, err := repo.CreateCause(ctx, domain.Cause{Code: "c1", Name: "other", Type: "poison", Prior: 0.5})
		assert.ErrorIs(t, err, domain.ErrDuplicateCode)
	})
	t.Run("type outside CHECK", func(t *testing.T) {
		_, err := repo.CreateCause(ctx, domain.Cause{Code: "c2", Name: "n", Type: "banana", Prior: 0.5})
		assert.ErrorIs(t, err, domain.ErrCheckPass)
	})
	t.Run("get missing id", func(t *testing.T) {
		_, err := repo.GetCauseByID(ctx, 999)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
	t.Run("update missing id", func(t *testing.T) {
		_, err := repo.UpdateCause(ctx, domain.Cause{ID: 999, Code: "x", Name: "n", Type: "poison", Prior: 0.5})
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
	t.Run("delete missing id is silently ok", func(t *testing.T) {
		// :exec does not report affected rows, so this cannot be a 404 today
		assert.NoError(t, repo.DeleteCause(ctx, 999))
	})
	t.Run("deleting a remedy sets cause.remedy_id to null", func(t *testing.T) {
		remedy, err := repo.CreateRemedy(ctx, domain.Remedy{Name: "r"})
		require.NoError(t, err)
		c, err := repo.CreateCause(ctx, domain.Cause{Code: "c3", Name: "n", Type: "poison", Prior: 0.5, RemedyID: &remedy.ID})
		require.NoError(t, err)
		_, err = repo.pool.Exec(ctx, "DELETE FROM remedies WHERE id = $1", remedy.ID)
		require.NoError(t, err)
		got, err := repo.GetCauseByID(ctx, int(c.ID))
		require.NoError(t, err)
		assert.Nil(t, got.RemedyID)
	})
}
