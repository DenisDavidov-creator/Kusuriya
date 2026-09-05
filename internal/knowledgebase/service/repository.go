package service

import (
	"Kusuriya/internal/domain"
	"context"
)

type KnowledgebaseRepositoryInterface interface {
	CreateCause(context.Context, domain.Cause) (*domain.Cause, error)
	GetCauseByID(context.Context, int) (*domain.Cause, error)
	GetCauses(context.Context) ([]domain.Cause, error)
	UpdateCause(context.Context, domain.Cause) (*domain.Cause, error)
	DeleteCause(context.Context, int) error

	CreateRemedy(context.Context, domain.Remedy) (*domain.Remedy, error)
	GetRemedyByID(context.Context, int) (*domain.Remedy, error)
	GetRemedies(context.Context) ([]domain.Remedy, error)

	CreateFactor(context.Context, domain.Factor) (*domain.Factor, error)
	GetFactorByID(context.Context, int) (*domain.Factor, error)
	GetFactors(context.Context) ([]domain.Factor, error)

	CreateCategory(context.Context, domain.Category) (*domain.Category, error)
	GetCategoryByID(context.Context, int) (*domain.Category, error)
	GetCategories(context.Context) ([]domain.Category, error)

	CreateSymptom(context.Context, domain.Symptom) (*domain.Symptom, error)
	GetSymptomByID(context.Context, int) (*domain.Symptom, error)
	GetSymptoms(context.Context) ([]domain.Symptom, error)

	CreateSymptomProfiles(context.Context, []domain.SymptomProfile) ([]domain.SymptomProfile, error)
	DeleteSymptomProfiles(context.Context, int64) error

	CreateFactorProfiles(context.Context, []domain.FactorProfile) ([]domain.FactorProfile, error)
	DeleteFactorProfiles(context.Context, int64) error

	WithTx(context.Context) (context.Context, error)
	Commit(context.Context) error
	Rollback(ctx context.Context) error
}
