//go:build integration

package repository

import (
	"Kusuriya/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// seedEvidence builds the smallest knowledge base that exercises every join branch:
// two causes, two symptoms, two factors, and profiles only on cause 1.
func seedEvidence(t *testing.T, repo *KnowledgebaseRepository) {
	t.Helper()
	remedy, err := repo.CreateRemedy(ctx, domain.Remedy{Name: "antidote", Description: "rd", Notes: "rn"})
	require.NoError(t, err)
	_, err = repo.CreateCause(ctx, domain.Cause{Code: "c1", Name: "n", Type: "poison", Prior: 0.5, RemedyID: &remedy.ID})
	require.NoError(t, err)
	_, err = repo.CreateCause(ctx, domain.Cause{Code: "c2", Name: "n", Type: "disease", Prior: 0.5})
	require.NoError(t, err)

	cat, err := repo.CreateCategory(ctx, domain.Category{Name: "cat"})
	require.NoError(t, err)
	_, err = repo.CreateSymptom(ctx, domain.Symptom{Code: "s1", Name: "n", CategoryID: cat.ID})
	require.NoError(t, err)
	_, err = repo.CreateSymptom(ctx, domain.Symptom{Code: "s2", Name: "n", CategoryID: cat.ID})
	require.NoError(t, err)
	_, err = repo.CreateFactor(ctx, domain.Factor{Code: "f1", Name: "n"})
	require.NoError(t, err)
	_, err = repo.CreateFactor(ctx, domain.Factor{Code: "f2", Name: "n"})
	require.NoError(t, err)

	_, err = repo.CreateSymptomProfiles(ctx, []domain.SymptomProfile{
		{CauseID: 1, SymptomID: 1, Likelihood: 0.8},
		{CauseID: 1, SymptomID: 2, Likelihood: 0.3},
	})
	require.NoError(t, err)
	_, err = repo.CreateFactorProfiles(ctx, []domain.FactorProfile{
		{CauseID: 1, FactorID: 1, Modifier: 0.7},
	})
	require.NoError(t, err)
}

func TestGetSymptomsByCodes(t *testing.T) {
	repo, _ := newRepo(t)
	seedEvidence(t, repo)

	t.Run("one row per cause x symptom; missing profile is nil", func(t *testing.T) {
		got, err := repo.GetSymptomsByCodes(ctx, []string{"s1", "s2"})
		require.NoError(t, err)
		assert.ElementsMatch(t, []domain.SymptomEvidence{
			{CauseID: 1, SymptomID: 1, SymptomCode: "s1", Likelihood: f(0.8)},
			{CauseID: 1, SymptomID: 2, SymptomCode: "s2", Likelihood: f(0.3)},
			{CauseID: 2, SymptomID: 1, SymptomCode: "s1", Likelihood: nil},
			{CauseID: 2, SymptomID: 2, SymptomCode: "s2", Likelihood: nil},
		}, got)
	})
	t.Run("unknown code contributes nothing", func(t *testing.T) {
		got, err := repo.GetSymptomsByCodes(ctx, []string{"s1", "nope"})
		require.NoError(t, err)
		assert.Len(t, got, 2, "only s1 x {c1, c2}")
	})
	t.Run("no codes -> no rows", func(t *testing.T) {
		got, err := repo.GetSymptomsByCodes(ctx, nil)
		require.NoError(t, err)
		assert.Empty(t, got)
	})
	t.Run("by ids gives the same rows", func(t *testing.T) {
		byCodes, err := repo.GetSymptomsByCodes(ctx, []string{"s1"})
		require.NoError(t, err)
		byIDs, err := repo.GetSymptomsByIDs(ctx, []int64{1})
		require.NoError(t, err)
		assert.ElementsMatch(t, byCodes, byIDs)
	})
}

func TestGetFactorsByCodes(t *testing.T) {
	repo, _ := newRepo(t)
	seedEvidence(t, repo)

	got, err := repo.GetFactorsByCodes(ctx, []string{"f1", "f2"})
	require.NoError(t, err)
	assert.ElementsMatch(t, []domain.FactorEvidence{
		{CauseID: 1, FactorID: 1, FactorCode: "f1", Modifier: f(0.7)},
		{CauseID: 1, FactorID: 2, FactorCode: "f2", Modifier: nil},
		{CauseID: 2, FactorID: 1, FactorCode: "f1", Modifier: nil},
		{CauseID: 2, FactorID: 2, FactorCode: "f2", Modifier: nil},
	}, got)

	byIDs, err := repo.GetFactorsByIDs(ctx, []int64{1, 2})
	require.NoError(t, err)
	assert.ElementsMatch(t, got, byIDs)
}

func TestGetCausesWithRemediesByIDs(t *testing.T) {
	repo, _ := newRepo(t)
	seedEvidence(t, repo)

	got, err := repo.GetCausesWithRemediesByIDs(ctx, []int64{2, 1, 999})
	require.NoError(t, err)
	assert.Equal(t, []domain.CauseWithRemedy{
		{ID: 1, Code: "c1", Name: "n", Type: "poison", Prior: 0.5,
			RemedyID: i64(1), RemedyName: str("antidote"), RemedyDescription: str("rd"), RemedyNotes: str("rn")},
		{ID: 2, Code: "c2", Name: "n", Type: "disease", Prior: 0.5},
	}, got, "ordered by id regardless of input order; unknown ids skipped; LEFT JOIN nulls -> nil")

	empty, err := repo.GetCausesWithRemediesByIDs(ctx, nil)
	require.NoError(t, err)
	assert.Empty(t, empty)
}

func TestProfiles(t *testing.T) {
	repo, pool := newRepo(t)
	seedEvidence(t, repo)

	t.Run("likelihood outside CHECK", func(t *testing.T) {
		_, err := repo.CreateSymptomProfiles(ctx, []domain.SymptomProfile{{CauseID: 2, SymptomID: 1, Likelihood: 1.5}})
		assert.ErrorIs(t, err, domain.ErrCheckPass)
	})
	t.Run("same symptom twice for one cause", func(t *testing.T) {
		_, err := repo.CreateSymptomProfiles(ctx, []domain.SymptomProfile{{CauseID: 1, SymptomID: 1, Likelihood: 0.5}})
		assert.ErrorIs(t, err, domain.ErrDuplicateCode)
	})
	t.Run("delete removes all profiles of a cause", func(t *testing.T) {
		require.NoError(t, repo.DeleteSymptomProfiles(ctx, 1))
		require.NoError(t, repo.DeleteFactorProfiles(ctx, 1))
		var n int
		require.NoError(t, pool.QueryRow(ctx, "SELECT count(*) FROM symptom_profiles WHERE cause_id = 1").Scan(&n))
		assert.Zero(t, n)
		require.NoError(t, pool.QueryRow(ctx, "SELECT count(*) FROM factor_profiles WHERE cause_id = 1").Scan(&n))
		assert.Zero(t, n)
	})
	t.Run("deleting a cause cascades to its profiles", func(t *testing.T) {
		_, err := repo.CreateSymptomProfiles(ctx, []domain.SymptomProfile{{CauseID: 2, SymptomID: 1, Likelihood: 0.5}})
		require.NoError(t, err)
		require.NoError(t, repo.DeleteCause(ctx, 2))
		var n int
		require.NoError(t, pool.QueryRow(ctx, "SELECT count(*) FROM symptom_profiles WHERE cause_id = 2").Scan(&n))
		assert.Zero(t, n)
	})
}
