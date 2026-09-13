package handler

import (
	"Kusuriya/internal/domain"
	"Kusuriya/internal/knowledgebase/handler/mocks"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	validFactorBody = `{"code":"f1","name":"f1_name","description":"desc"}`
	validFactor     = domain.Factor{Code: "f1", Name: "f1_name", Description: "desc"}
	savedFactor     = domain.Factor{ID: 1, Code: "f1", Name: "f1_name", Description: "desc"}
	savedFactorResp = FactorResponse{ID: 1, Code: "f1", Name: "f1_name", Description: "desc"}
)

func TestCreateFactor(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   *FactorResponse
	}{
		{
			name: "OK",
			body: validFactorBody,
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("CreateFactor", mock.Anything, validFactor).Return(&savedFactor, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &savedFactorResp,
		},
		{
			name:       "invalid json",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation fail: empty code",
			body:       `{"code":"","name":"f1_name"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate code",
			body: validFactorBody,
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("CreateFactor", mock.Anything, validFactor).Return(nil, domain.ErrDuplicateCode)
			},
			wantStatus: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewKnowledgebaseServiceInterface(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			h := NewKnowledgebaseHandler(mockService)

			rec := doRequest(h.CreateFactor, http.MethodPost, "", tt.body)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, *tt.wantBody, decodeBody[FactorResponse](t, rec))
			}
		})
	}
}

func TestGetFactorByID(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   *FactorResponse
	}{
		{
			name: "OK",
			id:   "1",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetFactorByID", mock.Anything, 1).Return(&savedFactor, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &savedFactorResp,
		},
		{
			name:       "wrong URLParam",
			id:         "meow",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "not found",
			id:   "1",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetFactorByID", mock.Anything, 1).Return(nil, domain.ErrNotFound)
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

			rec := doRequest(h.GetFactorByID, http.MethodGet, tt.id, "")

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, *tt.wantBody, decodeBody[FactorResponse](t, rec))
			}
		})
	}
}

func TestGetFactors(t *testing.T) {
	second := domain.Factor{ID: 2, Code: "f2", Name: "f2_name"}
	secondResp := FactorResponse{ID: 2, Code: "f2", Name: "f2_name"}

	tests := []struct {
		name       string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   []FactorResponse
	}{
		{
			name: "OK",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetFactors", mock.Anything).Return([]domain.Factor{savedFactor, second}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []FactorResponse{savedFactorResp, secondResp},
		},
		{
			name: "empty list is [] not null",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetFactors", mock.Anything).Return(nil, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []FactorResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewKnowledgebaseServiceInterface(t)
			tt.setupMock(mockService)
			h := NewKnowledgebaseHandler(mockService)

			rec := doRequest(h.GetFactors, http.MethodGet, "", "")

			assert.Equal(t, tt.wantStatus, rec.Code)
			raw := strings.TrimSpace(rec.Body.String())
			assert.Equal(t, tt.wantBody, decodeBody[[]FactorResponse](t, rec))
			if len(tt.wantBody) == 0 {
				assert.Equal(t, "[]", raw)
			}
		})
	}
}
