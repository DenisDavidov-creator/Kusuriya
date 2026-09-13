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
	validCategory     = domain.Category{Name: "cat1"}
	savedCategory     = domain.Category{ID: 1, Name: "cat1"}
	savedCategoryResp = CategoryResponse{ID: 1, Name: "cat1"}
)

func TestCreateCategory(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   *CategoryResponse
	}{
		{
			name: "OK",
			body: `{"name":"cat1"}`,
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("CreateCategory", mock.Anything, validCategory).Return(&savedCategory, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &savedCategoryResp,
		},
		{
			name:       "invalid json",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation fail: empty name",
			body:       `{"name":""}`,
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

			rec := doRequest(h.CreateCategory, http.MethodPost, "", tt.body)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, *tt.wantBody, decodeBody[CategoryResponse](t, rec))
			}
		})
	}
}

func TestGetCategoryByID(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   *CategoryResponse
	}{
		{
			name: "OK",
			id:   "1",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetCategoryByID", mock.Anything, 1).Return(&savedCategory, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &savedCategoryResp,
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
				m.On("GetCategoryByID", mock.Anything, 1).Return(nil, domain.ErrNotFound)
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

			rec := doRequest(h.GetCategoryByID, http.MethodGet, tt.id, "")

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, *tt.wantBody, decodeBody[CategoryResponse](t, rec))
			}
		})
	}
}

func TestGetCategories(t *testing.T) {
	second := domain.Category{ID: 2, Name: "cat2"}
	secondResp := CategoryResponse{ID: 2, Name: "cat2"}

	tests := []struct {
		name       string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   []CategoryResponse
	}{
		{
			name: "OK",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetCategories", mock.Anything).Return([]domain.Category{savedCategory, second}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []CategoryResponse{savedCategoryResp, secondResp},
		},
		{
			name: "empty list is [] not null",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetCategories", mock.Anything).Return(nil, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []CategoryResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewKnowledgebaseServiceInterface(t)
			tt.setupMock(mockService)
			h := NewKnowledgebaseHandler(mockService)

			rec := doRequest(h.GetCategories, http.MethodGet, "", "")

			assert.Equal(t, tt.wantStatus, rec.Code)
			raw := strings.TrimSpace(rec.Body.String())
			assert.Equal(t, tt.wantBody, decodeBody[[]CategoryResponse](t, rec))
			if len(tt.wantBody) == 0 {
				assert.Equal(t, "[]", raw)
			}
		})
	}
}
