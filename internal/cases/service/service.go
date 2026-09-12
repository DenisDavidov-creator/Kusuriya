package service

import (
	"Kusuriya/internal/domain"
	"context"
)

type CaseService struct {
	repoCase          CaseRepositoryInterface
	repoKnowledgebase KnowledgebaseRepositoryInterface
}

func NewCaseService(repoCase CaseRepositoryInterface, repoKnowledgebase KnowledgebaseRepositoryInterface) *CaseService {
	return &CaseService{
		repoCase:          repoCase,
		repoKnowledgebase: repoKnowledgebase,
	}
}

func (s *CaseService) CreateCase(ctx context.Context, dom domain.Case) (*domain.Case, error) {
	result, err := s.repoCase.CreateCase(ctx, dom)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CaseService) GetCaseByID(ctx context.Context, id int) ([]domain.DiagnosisResponse, *domain.Case, error) {
	result, err := s.repoCase.GetCaseByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	observations, err := s.repoCase.GetObservationByID(ctx, int(id))
	if err != nil {
		return nil, nil, err
	}
	var symptomIDs []int64
	var factorIDs []int64
	for _, obs := range observations {
		if obs.Kind == "symptom" {
			symptomIDs = append(symptomIDs, obs.RefID)
		}
		if obs.Kind == "factor" {
			factorIDs = append(factorIDs, obs.RefID)
		}
	}

	causesAll, err := s.repoKnowledgebase.GetCauses(ctx)
	if err != nil {
		return nil, nil, err
	}
	factors, err := s.repoKnowledgebase.GetFactorsByIDs(ctx, factorIDs)
	if err != nil {
		return nil, nil, err
	}
	symptom, err := s.repoKnowledgebase.GetSymptomsByIDs(ctx, symptomIDs)
	if err != nil {
		return nil, nil, err
	}

	scores := domain.EngineCause(causesAll, factors, symptom)
	if len(scores) > 5 {
		scores = scores[:5]
	}

	var causesIDs []int64
	for _, value := range scores {
		causesIDs = append(causesIDs, int64(value.CauseID))
	}

	causes, err := s.repoKnowledgebase.GetCausesWithRemediesByIDs(ctx, causesIDs)
	if err != nil {
		return nil, nil, err
	}

	finalScores := domain.CombineCausesWithScores(scores, causes)

	return finalScores, result, nil
}

func (s *CaseService) CloseCaseByID(ctx context.Context, id int) (*domain.Case, error) {
	result, err := s.repoCase.CloseCaseByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *CaseService) CreateObservation(ctx context.Context, d domain.Observation) (*domain.Observation, error) {

	result, err := s.repoCase.CreateObservation(ctx, d)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CaseService) GetObservationByID(ctx context.Context, ID int) ([]domain.Observation, error) {
	result, err := s.repoCase.GetObservationByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return result, nil
}
