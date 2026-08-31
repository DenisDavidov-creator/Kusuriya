package service

import (
	"Kusuriya/internal/domain"
	"context"
)

type KnowledgebaseRepositoryInterface interface {
	CreateCause(context.Context, domain.Cause) (*domain.Cause, error)
}
