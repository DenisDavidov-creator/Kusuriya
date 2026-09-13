package handler

import (
	"Kusuriya/internal/domain"
	"Kusuriya/internal/httpx/httptestx"
	"Kusuriya/internal/knowledgebase/handler/mocks"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSymptomProfiles(t *testing.T) {
	body := `{"profiles":[{"symptom_id":10,"likelihood":0.8},{"symptom_id":11,"likelihood":0.2}]}`
	in := []domain.SymptomProfile{
		{CauseID: 5, SymptomID: 10, Likelihood: 0.8},
		{CauseID: 5, SymptomID: 11, Likelihood: 0.2},
	}
	saved := []domain.SymptomProfile{
		{ID: 1, CauseID: 5, SymptomID: 10, Likelihood: 0.8},
		{ID: 2, CauseID: 5, SymptomID: 11, Likelihood: 0.2},
	}
	savedResp := []SymptomProfileResponse{
		{ID: 1, CauseID: 5, SymptomID: 10, Likelihood: 0.8},
		{ID: 2, CauseID: 5, SymptomID: 11, Likelihood: 0.2},
	}

	tests := []struct {
		name       string
		id         string
		body       string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   []SymptomProfileResponse
	}{
		{
			name: "OK: cause id from URL goes into every profile and the id arg",
			id:   "5",
			body: body,
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("CreateSymptomProfiles", mock.Anything, in, 5).Return(saved, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   savedResp,
		},
		{
			name:       "wrong URLParam",
			id:         "meow",
			body:       body,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			id:         "5",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation fail: likelihood out of range",
			id:         "5",
			body:       `{"profiles":[{"symptom_id":10,"likelihood":1.5}]}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "not found",
			id:   "5",
			body: body,
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("CreateSymptomProfiles", mock.Anything, in, 5).Return(nil, domain.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewKnowledgebaseServiceInterface(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			h := NewKnowledgebaseHandler(mockService)

			rec := httptestx.DoRequest(h.CreateSymptomProfiles, http.MethodPost, tt.id, tt.body)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, tt.wantBody, httptestx.DecodeBody[[]SymptomProfileResponse](t, rec))
			}
		})
	}
}

func TestCreateFactorProfiles(t *testing.T) {
	body := `{"profiles":[{"factor_id":20,"modifier":0.7},{"factor_id":21,"modifier":0.3}]}`
	in := []domain.FactorProfile{
		{CauseID: 5, FactorID: 20, Modifier: 0.7},
		{CauseID: 5, FactorID: 21, Modifier: 0.3},
	}
	saved := []domain.FactorProfile{
		{ID: 1, CauseID: 5, FactorID: 20, Modifier: 0.7},
		{ID: 2, CauseID: 5, FactorID: 21, Modifier: 0.3},
	}
	savedResp := []FactorProfileResponse{
		{ID: 1, CauseID: 5, FactorID: 20, Modifier: 0.7},
		{ID: 2, CauseID: 5, FactorID: 21, Modifier: 0.3},
	}

	tests := []struct {
		name       string
		id         string
		body       string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   []FactorProfileResponse
	}{
		{
			name: "OK",
			id:   "5",
			body: body,
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("CreateFactorProfiles", mock.Anything, in, 5).Return(saved, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   savedResp,
		},
		{
			name:       "wrong URLParam",
			id:         "meow",
			body:       body,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			id:         "5",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation fail: missing factor_id",
			id:         "5",
			body:       `{"profiles":[{"modifier":0.7}]}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "not found",
			id:   "5",
			body: body,
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("CreateFactorProfiles", mock.Anything, in, 5).Return(nil, domain.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewKnowledgebaseServiceInterface(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			h := NewKnowledgebaseHandler(mockService)

			rec := httptestx.DoRequest(h.CreateFactorProfiles, http.MethodPost, tt.id, tt.body)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, tt.wantBody, httptestx.DecodeBody[[]FactorProfileResponse](t, rec))
			}
		})
	}
}
