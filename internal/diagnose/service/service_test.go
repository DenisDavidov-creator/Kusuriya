package service

import (
	"Kusuriya/internal/diagnose/service/mocks"
	"Kusuriya/internal/domain"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func f(v float64) *float64 { return &v }

func TestDetermineDiagnose(t *testing.T) {
	req := domain.DiagnosisRequest{SymptomCodes: []string{"s1"}, FactorCodes: []string{"f1"}}
	twoCauses := []domain.Cause{{ID: 1, Prior: 0.5}, {ID: 2, Prior: 0.5}}
	boom := errors.New("db on fire")

	tests := []struct {
		name      string
		setupMock func(m *mocks.KnowledgebaseRepositoryInterface)
		wantErr   error
		check     func(t *testing.T, got []domain.DiagnosisResponse)
	}{
		{
			name: "codes are passed to the repo as given; evidence and remedy data reach the result",
			setupMock: func(m *mocks.KnowledgebaseRepositoryInterface) {
				m.On("GetCauses", mock.Anything).Return(twoCauses, nil)
				m.On("GetFactorsByCodes", mock.Anything, []string{"f1"}).Return(nil, nil)
				m.On("GetSymptomsByCodes", mock.Anything, []string{"s1"}).Return([]domain.SymptomEvidence{
					{CauseID: 1, SymptomID: 3, SymptomCode: "s1", Likelihood: f(0.8)},
					{CauseID: 2, SymptomID: 3, SymptomCode: "s1", Likelihood: nil},
				}, nil)
				m.On("GetCausesWithRemediesByIDs", mock.Anything, []int64{1, 2}).Return([]domain.CauseWithRemedy{
					{ID: 1, Code: "c1"}, {ID: 2, Code: "c2"},
				}, nil)
			},
			check: func(t *testing.T, got []domain.DiagnosisResponse) {
				require.Len(t, got, 2)
				assert.Equal(t, "c1", got[0].Code, "matched cause ranks first")
				assert.Equal(t, "c2", got[1].Code)
				assert.InDelta(t, 0.4/0.425, got[0].P, 1e-9)
				assert.Equal(t, []domain.Contribution{
					{Kind: "symptom", KindID: 3, Code: "s1", Weight: f(0.8), Matched: true},
				}, got[0].Evidence)
				assert.Equal(t, []domain.Contribution{
					{Kind: "symptom", KindID: 3, Code: "s1", Weight: nil, Matched: false},
				}, got[1].Evidence)
			},
		},
		{
			name: "only the top 5 causes are loaded and returned, in rank order",
			setupMock: func(m *mocks.KnowledgebaseRepositoryInterface) {
				// priors chosen so rank order != id order
				m.On("GetCauses", mock.Anything).Return([]domain.Cause{
					{ID: 1, Prior: 0.1}, {ID: 2, Prior: 0.6}, {ID: 3, Prior: 0.3},
					{ID: 4, Prior: 0.05}, {ID: 5, Prior: 0.5}, {ID: 6, Prior: 0.2}, {ID: 7, Prior: 0.4},
				}, nil)
				m.On("GetFactorsByCodes", mock.Anything, mock.Anything).Return(nil, nil)
				m.On("GetSymptomsByCodes", mock.Anything, mock.Anything).Return(nil, nil)
				// ids 1 and 4 (the two lowest priors) must not be requested
				m.On("GetCausesWithRemediesByIDs", mock.Anything, []int64{2, 5, 7, 3, 6}).Return([]domain.CauseWithRemedy{
					{ID: 2}, {ID: 5}, {ID: 7}, {ID: 3}, {ID: 6},
				}, nil)
			},
			check: func(t *testing.T, got []domain.DiagnosisResponse) {
				ids := make([]int64, 0, len(got))
				for _, g := range got {
					ids = append(ids, g.ID)
				}
				assert.Equal(t, []int64{2, 5, 7, 3, 6}, ids)
			},
		},
		{
			name: "empty knowledge base gives an empty result",
			setupMock: func(m *mocks.KnowledgebaseRepositoryInterface) {
				m.On("GetCauses", mock.Anything).Return(nil, nil)
				m.On("GetFactorsByCodes", mock.Anything, mock.Anything).Return(nil, nil)
				m.On("GetSymptomsByCodes", mock.Anything, mock.Anything).Return(nil, nil)
				m.On("GetCausesWithRemediesByIDs", mock.Anything, []int64(nil)).Return(nil, nil)
			},
			check: func(t *testing.T, got []domain.DiagnosisResponse) {
				assert.Empty(t, got)
			},
		},

		// every repo error is returned as is and stops the pipeline
		{
			name: "GetCauses fails",
			setupMock: func(m *mocks.KnowledgebaseRepositoryInterface) {
				m.On("GetCauses", mock.Anything).Return(nil, boom)
			},
			wantErr: boom,
		},
		{
			name: "GetFactorsByCodes fails",
			setupMock: func(m *mocks.KnowledgebaseRepositoryInterface) {
				m.On("GetCauses", mock.Anything).Return(twoCauses, nil)
				m.On("GetFactorsByCodes", mock.Anything, mock.Anything).Return(nil, boom)
			},
			wantErr: boom,
		},
		{
			name: "GetSymptomsByCodes fails",
			setupMock: func(m *mocks.KnowledgebaseRepositoryInterface) {
				m.On("GetCauses", mock.Anything).Return(twoCauses, nil)
				m.On("GetFactorsByCodes", mock.Anything, mock.Anything).Return(nil, nil)
				m.On("GetSymptomsByCodes", mock.Anything, mock.Anything).Return(nil, boom)
			},
			wantErr: boom,
		},
		{
			name: "GetCausesWithRemediesByIDs fails",
			setupMock: func(m *mocks.KnowledgebaseRepositoryInterface) {
				m.On("GetCauses", mock.Anything).Return(twoCauses, nil)
				m.On("GetFactorsByCodes", mock.Anything, mock.Anything).Return(nil, nil)
				m.On("GetSymptomsByCodes", mock.Anything, mock.Anything).Return(nil, nil)
				m.On("GetCausesWithRemediesByIDs", mock.Anything, mock.Anything).Return(nil, boom)
			},
			wantErr: boom,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewKnowledgebaseRepositoryInterface(t)
			tt.setupMock(repo)
			s := NewDiagnoseService(repo)

			got, err := s.DetermineDiagnose(context.Background(), req)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}
			require.NoError(t, err)
			tt.check(t, got)
		})
	}
}
