package handler

import (
	"Kusuriya/internal/domain"
	"Kusuriya/internal/httpx/httptestx"
	"Kusuriya/internal/knowledgebase/handler/mocks"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	validSymptomBody = `{"code":"s1","name":"s1_name","description":"desc","category_id":3}`
	validSymptom     = domain.Symptom{Code: "s1", Name: "s1_name", Description: "desc", CategoryID: 3}
	savedSymptom     = domain.Symptom{ID: 1, Code: "s1", Name: "s1_name", Description: "desc", CategoryID: 3}
	savedSymptomResp = SymptomResponse{ID: 1, Code: "s1", Name: "s1_name", Description: "desc", CategoryID: 3}
)

func TestCreateSymptom(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   *SymptomResponse
	}{
		{
			name: "OK",
			body: validSymptomBody,
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("CreateSymptom", mock.Anything, validSymptom).Return(&savedSymptom, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &savedSymptomResp,
		},
		{
			name:       "invalid json",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation fail: missing category",
			body:       `{"code":"s1","name":"s1_name"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate code",
			body: validSymptomBody,
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("CreateSymptom", mock.Anything, validSymptom).Return(nil, domain.ErrDuplicateCode)
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

			rec := httptestx.DoRequest(h.CreateSymptom, http.MethodPost, "", tt.body)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, *tt.wantBody, httptestx.DecodeBody[SymptomResponse](t, rec))
			}
		})
	}
}

func TestGetSymptomByID(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   *SymptomResponse
	}{
		{
			name: "OK",
			id:   "1",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetSymptomByID", mock.Anything, 1).Return(&savedSymptom, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &savedSymptomResp,
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
				m.On("GetSymptomByID", mock.Anything, 1).Return(nil, domain.ErrNotFound)
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

			rec := httptestx.DoRequest(h.GetSymptomByID, http.MethodGet, tt.id, "")

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, *tt.wantBody, httptestx.DecodeBody[SymptomResponse](t, rec))
			}
		})
	}
}

func TestGetSymptoms(t *testing.T) {
	second := domain.Symptom{ID: 2, Code: "s2", Name: "s2_name", CategoryID: 3}
	secondResp := SymptomResponse{ID: 2, Code: "s2", Name: "s2_name", CategoryID: 3}

	tests := []struct {
		name       string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   []SymptomResponse
	}{
		{
			name: "OK",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetSymptoms", mock.Anything).Return([]domain.Symptom{savedSymptom, second}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []SymptomResponse{savedSymptomResp, secondResp},
		},
		{
			name: "empty list is [] not null",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetSymptoms", mock.Anything).Return(nil, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []SymptomResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewKnowledgebaseServiceInterface(t)
			tt.setupMock(mockService)
			h := NewKnowledgebaseHandler(mockService)

			rec := httptestx.DoRequest(h.GetSymptoms, http.MethodGet, "", "")

			assert.Equal(t, tt.wantStatus, rec.Code)
			raw := strings.TrimSpace(rec.Body.String())
			assert.Equal(t, tt.wantBody, httptestx.DecodeBody[[]SymptomResponse](t, rec))
			if len(tt.wantBody) == 0 {
				assert.Equal(t, "[]", raw)
			}
		})
	}
}
