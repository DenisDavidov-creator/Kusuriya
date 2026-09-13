package handler

import (
	"Kusuriya/internal/domain"
	"context"
)

//go:generate mockery --name=DiagnoseServiceInterface
type DiagnoseServiceInterface interface {
	DetermineDiagnose(context.Context, domain.DiagnosisRequest) ([]domain.DiagnosisResponse, error)
}
