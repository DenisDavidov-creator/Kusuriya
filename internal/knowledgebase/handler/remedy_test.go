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
	validRemedyBody = `{"name":"r1","description":"desc","notes":"note"}`
	validRemedy     = domain.Remedy{Name: "r1", Description: "desc", Notes: "note"}
	savedRemedy     = domain.Remedy{ID: 1, Name: "r1", Description: "desc", Notes: "note"}
	savedRemedyResp = RemedyResponse{ID: 1, Name: "r1", Description: "desc", Notes: "note"}
)

func TestCreateRemedy(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   *RemedyResponse
	}{
		{
			name: "OK",
			body: validRemedyBody,
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("CreateRemedy", mock.Anything, validRemedy).Return(&savedRemedy, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &savedRemedyResp,
		},
		{
			name:       "invalid json",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation fail: empty name",
			body:       `{"name":"","description":"desc"}`,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewKnowledgebaseServiceInterface(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			h := NewKnowledgebaseHandler(mockService)

			rec := doRequest(h.CreateRemedy, http.MethodPost, "", tt.body)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, *tt.wantBody, decodeBody[RemedyResponse](t, rec))
			}
		})
	}
}

func TestGetRemedyByID(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   *RemedyResponse
	}{
		{
			name: "OK",
			id:   "1",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetRemedyByID", mock.Anything, 1).Return(&savedRemedy, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &savedRemedyResp,
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
				m.On("GetRemedyByID", mock.Anything, 1).Return(nil, domain.ErrNotFound)
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

			rec := doRequest(h.GetRemedyByID, http.MethodGet, tt.id, "")

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, *tt.wantBody, decodeBody[RemedyResponse](t, rec))
			}
		})
	}
}

func TestGetRemedies(t *testing.T) {
	second := domain.Remedy{ID: 2, Name: "r2"}
	secondResp := RemedyResponse{ID: 2, Name: "r2"}

	tests := []struct {
		name       string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   []RemedyResponse
	}{
		{
			name: "OK",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetRemedies", mock.Anything).Return([]domain.Remedy{savedRemedy, second}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []RemedyResponse{savedRemedyResp, secondResp},
		},
		{
			name: "empty list is [] not null",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetRemedies", mock.Anything).Return(nil, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []RemedyResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewKnowledgebaseServiceInterface(t)
			tt.setupMock(mockService)
			h := NewKnowledgebaseHandler(mockService)

			rec := doRequest(h.GetRemedies, http.MethodGet, "", "")

			assert.Equal(t, tt.wantStatus, rec.Code)
			raw := strings.TrimSpace(rec.Body.String())
			assert.Equal(t, tt.wantBody, decodeBody[[]RemedyResponse](t, rec))
			if len(tt.wantBody) == 0 {
				assert.Equal(t, "[]", raw)
			}
		})
	}
}
