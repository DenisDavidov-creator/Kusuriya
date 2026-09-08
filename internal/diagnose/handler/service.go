package handler

import (
	"Kusuriya/internal/domain"
	"context"
)

type DiagnoseServiceInterface interface {
	DeterminDiagnose(context.Context, domain.DiagnosisRequest) ([]domain.DiagnosisResponse, error)
}
