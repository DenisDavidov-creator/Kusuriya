package service

import (
	"Kusuriya/internal/domain"
	"context"
)

//go:generate mockery --name=CaseRepositoryInterface
type CaseRepositoryInterface interface {
	CreateCase(context.Context, domain.Case) (*domain.Case, error)
	GetCaseByID(context.Context, int) (*domain.Case, error)
	CloseCaseByID(context.Context, int) (*domain.Case, error)
	CreateObservation(context.Context, domain.Observation) (*domain.Observation, error)
	GetObservationByID(context.Context, int) ([]domain.Observation, error)
}

//go:generate mockery --name=KnowledgebaseRepositoryInterface
type KnowledgebaseRepositoryInterface interface {
	GetCauses(context.Context) ([]domain.Cause, error)
	GetCausesWithRemediesByIDs(context.Context, []int64) ([]domain.CauseWithRemedy, error)

	GetFactorsByIDs(ctx context.Context, IDs []int64) ([]domain.FactorEvidence, error)
	GetSymptomsByIDs(ctx context.Context, IDs []int64) ([]domain.SymptomEvidence, error)
}
