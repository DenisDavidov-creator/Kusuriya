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

func (s *KnowledgebaseService) GetCauseByID(ctx context.Context, id int) (*domain.Cause, error) {
	dom, err := s.repo.GetCauseByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (s *KnowledgebaseService) GetCauses(ctx context.Context) ([]domain.Cause, error) {

	doms, err := s.repo.GetCauses(ctx)
	if err != nil {
		return nil, err
	}
	return doms, nil
}

func (s *KnowledgebaseService) UpdateCause(ctx context.Context, d domain.Cause) (*domain.Cause, error) {
	dom, err := s.repo.UpdateCause(ctx, d)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (s *KnowledgebaseService) DeleteCause(ctx context.Context, id int) error {
	err := s.repo.DeleteCause(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
