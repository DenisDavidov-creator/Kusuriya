package service

import (
	"Kusuriya/internal/domain"
	"context"
)

type DiagnoseService struct {
	repo KnowledgebaseRepositoryInterface
}

func NewDiagnoseService(repo KnowledgebaseRepositoryInterface) *DiagnoseService {
	return &DiagnoseService{repo: repo}
}

func (s *DiagnoseService) DeterminDiagnose(ctx context.Context, dom domain.DiagnosisRequest) ([]domain.DiagnosisResponse, error) {
	causeDoms, err := s.repo.GetCauses(ctx)

	if err != nil {
		return nil, err
	}
	factors, err := s.repo.GetFactorsByCodes(ctx, dom.FactorCodes)
	if err != nil {
		return nil, err
	}
	symptoms, err := s.repo.GetSymptomsByCodes(ctx, dom.SymptomCodes)
	if err != nil {
		return nil, err
	}

	scores := domain.EngineCause(causeDoms, factors, symptoms)
	if len(scores) > 5 {
		scores = scores[:5]
	}

	var causesIDs []int64
	for _, value := range scores {
		causesIDs = append(causesIDs, int64(value.CauseID))
	}

	causes, err := s.repo.GetCausesWithRemediesByIDs(ctx, causesIDs)
	if err != nil {
		return nil, err
	}

	var causesMap = make(map[int64]domain.CauseWithRemedy, len(causes))
	for _, value := range causes {
		causesMap[value.ID] = value
	}

	var finalScores = make([]domain.DiagnosisResponse, 0, len(causes))
	for _, value := range scores {
		cause, ok := causesMap[value.CauseID]
		if ok {
			score := domain.DiagnosisResponse{
				P:               value.P,
				CauseWithRemedy: cause,
			}
			finalScores = append(finalScores, score)
		}
	}

	return finalScores, nil
}
