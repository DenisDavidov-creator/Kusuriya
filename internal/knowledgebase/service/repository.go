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
}
