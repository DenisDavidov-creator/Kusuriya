package handler

import (
	"Kusuriya/internal/domain"
	"context"
)

type KnowledgebaseServiceInterface interface {
	CreateCause(context.Context, domain.Cause) (*domain.Cause, error)
}
