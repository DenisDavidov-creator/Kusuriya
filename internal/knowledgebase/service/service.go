package service

import (
	"Kusuriya/internal/domain"
	"context"
)

type KnowledgebaseService struct {
	repo KnowledgebaseRepositoryInterface
}

func NewKnowledgebaseService(repo KnowledgebaseRepositoryInterface) *KnowledgebaseService {
	return &KnowledgebaseService{
		repo: repo,
	}
}

func (s *KnowledgebaseService) CreateCause(ctx context.Context, d domain.Cause) (*domain.Cause, error) {
	dom, err := s.repo.CreateCause(ctx, d)
	if err != nil {
		return nil, err
	}
	return dom, nil
}
