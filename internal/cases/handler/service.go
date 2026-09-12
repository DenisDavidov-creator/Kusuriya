package handler

import (
	"Kusuriya/internal/domain"
	"context"
)

type CaseServiceInterface interface {
	CreateCase(context.Context, domain.Case) (*domain.Case, error)

	GetCaseByID(context.Context, int) ([]domain.DiagnosisResponse, *domain.Case, error)

	CloseCaseByID(context.Context, int) (*domain.Case, error)
	CreateObservation(context.Context, domain.Observation) (*domain.Observation, error)
	GetObservationByID(context.Context, int) ([]domain.Observation, error)
}
