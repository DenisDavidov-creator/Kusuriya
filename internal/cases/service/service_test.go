package service

import (
	"Kusuriya/internal/cases/service/mocks"
	"Kusuriya/internal/domain"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func f(v float64) *float64 { return &v }

func TestGetCaseByID(t *testing.T) {
	openCase := domain.Case{ID: 5, Title: "t", Status: "open"}
	twoCauses := []domain.Cause{{ID: 1, Prior: 0.5}, {ID: 2, Prior: 0.5}}
	boom := errors.New("db on fire")


	tests := []struct {
		name    string
		setup   func(cr *mocks.CaseRepositoryInterface, kb *mocks.KnowledgebaseRepositoryInterface)
		wantErr error
		check   func(t *testing.T, diagnoses []domain.DiagnosisResponse, c *domain.Case)
	}{
		{
			name: "observations are split by kind into symptom and factor ids",
			setup: func(cr *mocks.CaseRepositoryInterface, kb *mocks.KnowledgebaseRepositoryInterface) {
				cr.On("GetCaseByID", mock.Anything, 5).Return(&openCase, nil)
				cr.On("GetObservationByID", mock.Anything, 5).Return([]domain.Observation{
					{Kind: "symptom", RefID: 10},
					{Kind: "factor", RefID: 20},
					{Kind: "symptom", RefID: 11},
					{Kind: "banana", RefID: 99}, // unknown kind is ignored
				}, nil)
				kb.On("GetCauses", mock.Anything).Return(twoCauses, nil)
				kb.On("GetFactorsByIDs", mock.Anything, []int64{20}).Return(nil, nil)
				kb.On("GetSymptomsByIDs", mock.Anything, []int64{10, 11}).Return([]domain.SymptomEvidence{
					{CauseID: 1, SymptomID: 10, SymptomCode: "s10", Likelihood: f(0.8)},
					{CauseID: 2, SymptomID: 10, SymptomCode: "s10", Likelihood: nil},
				}, nil)
				kb.On("GetCausesWithRemediesByIDs", mock.Anything, []int64{1, 2}).Return([]domain.CauseWithRemedy{
					{ID: 1, Code: "c1"}, {ID: 2, Code: "c2"},
				}, nil)
			},
			check: func(t *testing.T, diagnoses []domain.DiagnosisResponse, c *domain.Case) {
				assert.Equal(t, &openCase, c)
				require.Len(t, diagnoses, 2)
				assert.Equal(t, "c1", diagnoses[0].Code, "matched cause ranks first")
				assert.InDelta(t, 0.4/0.425, diagnoses[0].P, 1e-9)
				assert.Equal(t, []domain.Contribution{
					{Kind: "symptom", KindID: 10, Code: "s10", Weight: f(0.8), Matched: true},
				}, diagnoses[0].Evidence)
			},
		},
		{
			name: "case without observations still returns the case, ranked by priors only",
			setup: func(cr *mocks.CaseRepositoryInterface, kb *mocks.KnowledgebaseRepositoryInterface) {
				cr.On("GetCaseByID", mock.Anything, 5).Return(&openCase, nil)
				cr.On("GetObservationByID", mock.Anything, 5).Return(nil, nil)
				kb.On("GetCauses", mock.Anything).Return(twoCauses, nil)
				kb.On("GetFactorsByIDs", mock.Anything, []int64(nil)).Return(nil, nil)
				kb.On("GetSymptomsByIDs", mock.Anything, []int64(nil)).Return(nil, nil)
				kb.On("GetCausesWithRemediesByIDs", mock.Anything, mock.Anything).Return([]domain.CauseWithRemedy{{ID: 1}, {ID: 2}}, nil)
			},
			check: func(t *testing.T, diagnoses []domain.DiagnosisResponse, c *domain.Case) {
				assert.Equal(t, &openCase, c)
				require.Len(t, diagnoses, 2)
				assert.InDelta(t, 0.5, diagnoses[0].P, 1e-9)
				assert.Empty(t, diagnoses[0].Evidence)
			},
		},
		{
			name: "only the top 5 causes are loaded",
			setup: func(cr *mocks.CaseRepositoryInterface, kb *mocks.KnowledgebaseRepositoryInterface) {
				cr.On("GetCaseByID", mock.Anything, 5).Return(&openCase, nil)
				cr.On("GetObservationByID", mock.Anything, 5).Return(nil, nil)
				kb.On("GetCauses", mock.Anything).Return([]domain.Cause{
					{ID: 1, Prior: 0.1}, {ID: 2, Prior: 0.6}, {ID: 3, Prior: 0.3},
					{ID: 4, Prior: 0.05}, {ID: 5, Prior: 0.5}, {ID: 6, Prior: 0.2}, {ID: 7, Prior: 0.4},
				}, nil)
				kb.On("GetFactorsByIDs", mock.Anything, mock.Anything).Return(nil, nil)
				kb.On("GetSymptomsByIDs", mock.Anything, mock.Anything).Return(nil, nil)
				kb.On("GetCausesWithRemediesByIDs", mock.Anything, []int64{2, 5, 7, 3, 6}).Return(nil, nil)
			},
			check: func(t *testing.T, diagnoses []domain.DiagnosisResponse, c *domain.Case) {},
		},

		// errors: each call's failure is returned as is and stops the pipeline
		{
			name: "case not found: knowledgebase is never touched",
			setup: func(cr *mocks.CaseRepositoryInterface, kb *mocks.KnowledgebaseRepositoryInterface) {
				cr.On("GetCaseByID", mock.Anything, 5).Return(nil, domain.ErrNotFound)
			},
			wantErr: domain.ErrNotFound,
		},
		{
			name: "GetObservationByID fails",
			setup: func(cr *mocks.CaseRepositoryInterface, kb *mocks.KnowledgebaseRepositoryInterface) {
				cr.On("GetCaseByID", mock.Anything, 5).Return(&openCase, nil)
				cr.On("GetObservationByID", mock.Anything, 5).Return(nil, boom)
			},
			wantErr: boom,
		},
		{
			name: "GetCauses fails",
			setup: func(cr *mocks.CaseRepositoryInterface, kb *mocks.KnowledgebaseRepositoryInterface) {
				cr.On("GetCaseByID", mock.Anything, 5).Return(&openCase, nil)
				cr.On("GetObservationByID", mock.Anything, 5).Return(nil, nil)
				kb.On("GetCauses", mock.Anything).Return(nil, boom)
			},
			wantErr: boom,
		},
		{
			name: "GetCausesWithRemediesByIDs fails",
			setup: func(cr *mocks.CaseRepositoryInterface, kb *mocks.KnowledgebaseRepositoryInterface) {
				cr.On("GetCaseByID", mock.Anything, 5).Return(&openCase, nil)
				cr.On("GetObservationByID", mock.Anything, 5).Return(nil, nil)
				kb.On("GetCauses", mock.Anything).Return(twoCauses, nil)
				kb.On("GetFactorsByIDs", mock.Anything, mock.Anything).Return(nil, nil)
				kb.On("GetSymptomsByIDs", mock.Anything, mock.Anything).Return(nil, nil)
				kb.On("GetCausesWithRemediesByIDs", mock.Anything, mock.Anything).Return(nil, boom)
			},
			wantErr: boom,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := mocks.NewCaseRepositoryInterface(t)
			kb := mocks.NewKnowledgebaseRepositoryInterface(t)
			tt.setup(cr, kb)
			s := NewCaseService(cr, kb)

			diagnoses, c, err := s.GetCaseByID(context.Background(), 5)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, diagnoses)
				assert.Nil(t, c)
				return
			}
			require.NoError(t, err)
			tt.check(t, diagnoses, c)
		})
	}
}
