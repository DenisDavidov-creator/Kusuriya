package service

import (
	"Kusuriya/internal/domain"
	"context"
)

type KnowledgebaseRepositoryInterface interface {
	GetCauses(context.Context) ([]domain.Cause, error)
	GetCausesWithRemediesByIDs(context.Context, []int64) ([]domain.CauseWithRemedy, error)

	GetSymptomsByCodes(ctx context.Context, codes []string) ([]domain.SymptomEvidence, error)
	GetFactorsByCodes(ctx context.Context, codes []string) ([]domain.FactorEvidence, error)
}
