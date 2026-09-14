//go:build integration

package repository

import (
	"Kusuriya/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemedyCRUD(t *testing.T) {
	repo, _ := newRepo(t)

	full, err := repo.CreateRemedy(ctx, domain.Remedy{Name: "antidote", Description: "d", Notes: "n"})
	require.NoError(t, err)
	assert.Equal(t, &domain.Remedy{ID: 1, Name: "antidote", Description: "d", Notes: "n"}, full)

	bare, err := repo.CreateRemedy(ctx, domain.Remedy{Name: "charcoal"})
	require.NoError(t, err)
	assert.Equal(t, &domain.Remedy{ID: 2, Name: "charcoal"}, bare, "null description/notes -> empty strings")

	got, err := repo.GetRemedyByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, full, got)

	all, err := repo.GetRemedies(ctx)
	require.NoError(t, err)
	assert.Equal(t, []domain.Remedy{*full, *bare}, all)

	_, err = repo.GetRemedyByID(ctx, 999)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = repo.CreateRemedy(ctx, domain.Remedy{Name: "antidote"})
	assert.ErrorIs(t, err, domain.ErrDuplicateCode, "remedy name is unique (migration 000002)")
}

func TestFactorCRUD(t *testing.T) {
	repo, _ := newRepo(t)

	f1, err := repo.CreateFactor(ctx, domain.Factor{Code: "f1", Name: "n1", Description: "d"})
	require.NoError(t, err)
	assert.Equal(t, &domain.Factor{ID: 1, Code: "f1", Name: "n1", Description: "d"}, f1)

	f2, err := repo.CreateFactor(ctx, domain.Factor{Code: "f2", Name: "n2"})
	require.NoError(t, err)

	got, err := repo.GetFactorByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, f1, got)

	all, err := repo.GetFactors(ctx)
	require.NoError(t, err)
	assert.Equal(t, []domain.Factor{*f1, *f2}, all)

	_, err = repo.GetFactorByID(ctx, 999)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = repo.CreateFactor(ctx, domain.Factor{Code: "f1", Name: "dup"})
	assert.ErrorIs(t, err, domain.ErrDuplicateCode)
}

func TestCategoryAndSymptomCRUD(t *testing.T) {
	repo, _ := newRepo(t)

	cat, err := repo.CreateCategory(ctx, domain.Category{Name: "fever"})
	require.NoError(t, err)
	assert.Equal(t, &domain.Category{ID: 1, Name: "fever"}, cat)

	gotCat, err := repo.GetCategoryByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, cat, gotCat)

	cats, err := repo.GetCategories(ctx)
	require.NoError(t, err)
	assert.Equal(t, []domain.Category{*cat}, cats)

	_, err = repo.GetCategoryByID(ctx, 999)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	s1, err := repo.CreateSymptom(ctx, domain.Symptom{Code: "s1", Name: "n1", Description: "d", CategoryID: cat.ID})
	require.NoError(t, err)
	assert.Equal(t, &domain.Symptom{ID: 1, Code: "s1", Name: "n1", Description: "d", CategoryID: 1}, s1)

	s2, err := repo.CreateSymptom(ctx, domain.Symptom{Code: "s2", Name: "n2", CategoryID: cat.ID})
	require.NoError(t, err)

	gotS, err := repo.GetSymptomByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, s1, gotS)

	all, err := repo.GetSymptoms(ctx)
	require.NoError(t, err)
	assert.Equal(t, []domain.Symptom{*s1, *s2}, all)

	_, err = repo.GetSymptomByID(ctx, 999)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = repo.CreateSymptom(ctx, domain.Symptom{Code: "s1", Name: "dup", CategoryID: cat.ID})
	assert.ErrorIs(t, err, domain.ErrDuplicateCode)

	t.Run("symptom with unknown category", func(t *testing.T) {
		// FK violation (23503) is not mapped today, so it surfaces as a raw pg error
		_, err := repo.CreateSymptom(ctx, domain.Symptom{Code: "s3", Name: "n", CategoryID: 999})
		require.Error(t, err)
		assert.NotErrorIs(t, err, domain.ErrNotFound)
		assert.NotErrorIs(t, err, domain.ErrCheckPass)
	})
}

func TestEmptyListsAreEmptyNotNil(t *testing.T) {
	repo, _ := newRepo(t)

	causes, err := repo.GetCauses(ctx)
	require.NoError(t, err)
	assert.NotNil(t, causes)
	assert.Empty(t, causes)

	remedies, err := repo.GetRemedies(ctx)
	require.NoError(t, err)
	assert.NotNil(t, remedies)

	factors, err := repo.GetFactors(ctx)
	require.NoError(t, err)
	assert.NotNil(t, factors)

	symptoms, err := repo.GetSymptoms(ctx)
	require.NoError(t, err)
	assert.NotNil(t, symptoms)
}
