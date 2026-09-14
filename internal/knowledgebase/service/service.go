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

func (s *KnowledgebaseService) CreateRemedy(ctx context.Context, d domain.Remedy) (*domain.Remedy, error) {
	dom, err := s.repo.CreateRemedy(ctx, d)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (s *KnowledgebaseService) GetRemedyByID(ctx context.Context, id int) (*domain.Remedy, error) {
	dom, err := s.repo.GetRemedyByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (s *KnowledgebaseService) GetRemedies(ctx context.Context) ([]domain.Remedy, error) {

	doms, err := s.repo.GetRemedies(ctx)
	if err != nil {
		return nil, err
	}
	return doms, nil
}

func (s *KnowledgebaseService) CreateFactor(ctx context.Context, d domain.Factor) (*domain.Factor, error) {
	dom, err := s.repo.CreateFactor(ctx, d)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (s *KnowledgebaseService) GetFactorByID(ctx context.Context, id int) (*domain.Factor, error) {
	dom, err := s.repo.GetFactorByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (s *KnowledgebaseService) GetFactors(ctx context.Context) ([]domain.Factor, error) {

	doms, err := s.repo.GetFactors(ctx)
	if err != nil {
		return nil, err
	}
	return doms, nil
}

func (s *KnowledgebaseService) CreateCategory(ctx context.Context, d domain.Category) (*domain.Category, error) {
	dom, err := s.repo.CreateCategory(ctx, d)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (s *KnowledgebaseService) GetCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	dom, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (s *KnowledgebaseService) GetCategories(ctx context.Context) ([]domain.Category, error) {

	doms, err := s.repo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}
	return doms, nil
}

func (s *KnowledgebaseService) CreateSymptom(ctx context.Context, d domain.Symptom) (*domain.Symptom, error) {
	dom, err := s.repo.CreateSymptom(ctx, d)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (s *KnowledgebaseService) GetSymptomByID(ctx context.Context, id int) (*domain.Symptom, error) {
	dom, err := s.repo.GetSymptomByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (s *KnowledgebaseService) GetSymptoms(ctx context.Context) ([]domain.Symptom, error) {

	doms, err := s.repo.GetSymptoms(ctx)
	if err != nil {
		return nil, err
	}
	return doms, nil
}

func (s *KnowledgebaseService) CreateSymptomProfiles(ctx context.Context, ds []domain.SymptomProfile, id int) ([]domain.SymptomProfile, error) {
	ctx, err := s.repo.WithTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = s.repo.Rollback(ctx) }()
	err = s.repo.DeleteSymptomProfiles(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	doms, err := s.repo.CreateSymptomProfiles(ctx, ds)
	if err != nil {
		return nil, err
	}
	err = s.repo.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return doms, nil
}

func (s *KnowledgebaseService) CreateFactorProfiles(ctx context.Context, ds []domain.FactorProfile, id int) ([]domain.FactorProfile, error) {
	ctx, err := s.repo.WithTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = s.repo.Rollback(ctx) }()
	err = s.repo.DeleteFactorProfiles(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	doms, err := s.repo.CreateFactorProfiles(ctx, ds)
	if err != nil {
		return nil, err
	}
	err = s.repo.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return doms, nil
}
