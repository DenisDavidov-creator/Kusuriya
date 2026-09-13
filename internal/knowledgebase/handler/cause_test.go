package handler

import (
	"Kusuriya/internal/domain"
	"Kusuriya/internal/httpx/httptestx"
	"Kusuriya/internal/knowledgebase/handler/mocks"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	validCauseBody  = `{"code":"c1","name":"c1_name","type":"poison","prior":0.5}`
	validCause      = domain.Cause{Code: "c1", Name: "c1_name", Type: "poison", Prior: 0.5}
	savedCause      = domain.Cause{ID: 1, Code: "c1", Name: "c1_name", Type: "poison", Prior: 0.5}
	savedCauseResp  = CauseResponse{ID: 1, Code: "c1", Name: "c1_name", Type: "poison", Prior: 0.5}
	secondCause     = domain.Cause{ID: 2, Code: "c2", Name: "c2_name", Type: "disease", Prior: 0.1}
	secondCauseResp = CauseResponse{ID: 2, Code: "c2", Name: "c2_name", Type: "disease", Prior: 0.1}
)

func TestCreateCause(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   *CauseResponse
	}{
		{
			name: "OK",
			body: validCauseBody,
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("CreateCause", mock.Anything, validCause).Return(&savedCause, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &savedCauseResp,
		},
		{
			name:       "invalid json",
			body:       `{"code":"c1","name":"c1_name","type":"poison","prior":0.5`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation fail",
			body:       `{"code":"c1","name":"c1_name","type":"banana","prior":0.5}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate code",
			body: validCauseBody,
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("CreateCause", mock.Anything, validCause).Return(nil, domain.ErrDuplicateCode)
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "unexpected error",
			body: validCauseBody,
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("CreateCause", mock.Anything, validCause).Return(nil, errors.New("Meow"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewKnowledgebaseServiceInterface(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			h := NewKnowledgebaseHandler(mockService)

			rec := httptestx.DoRequest(h.CreateCause, http.MethodPost, "", tt.body)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, *tt.wantBody, httptestx.DecodeBody[CauseResponse](t, rec))
			}
		})
	}
}

func TestGetCauseByID(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   *CauseResponse
	}{
		{
			name: "OK",
			id:   "1",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetCauseByID", mock.Anything, 1).Return(&savedCause, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &savedCauseResp,
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
				m.On("GetCauseByID", mock.Anything, 1).Return(nil, domain.ErrNotFound)
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

			rec := httptestx.DoRequest(h.GetCauseByID, http.MethodGet, tt.id, "")

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, *tt.wantBody, httptestx.DecodeBody[CauseResponse](t, rec))
			}
		})
	}
}

func TestGetCauses(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   []CauseResponse
	}{
		{
			name: "OK",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetCauses", mock.Anything).Return([]domain.Cause{savedCause, secondCause}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []CauseResponse{savedCauseResp, secondCauseResp},
		},
		{
			name: "empty list is [] not null",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetCauses", mock.Anything).Return(nil, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []CauseResponse{},
		},
		{
			name: "unexpected error",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("GetCauses", mock.Anything).Return(nil, errors.New("Meow"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewKnowledgebaseServiceInterface(t)
			tt.setupMock(mockService)
			h := NewKnowledgebaseHandler(mockService)

			rec := httptestx.DoRequest(h.GetCauses, http.MethodGet, "", "")

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				raw := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, tt.wantBody, httptestx.DecodeBody[[]CauseResponse](t, rec))
				if len(tt.wantBody) == 0 {
					assert.Equal(t, "[]", raw)
				}
			}
		})
	}
}

func TestUpdateCause(t *testing.T) {
	updated := domain.Cause{ID: 1, Code: "c1", Name: "new_name", Type: "disease", Prior: 0.9}
	updatedResp := CauseResponse{ID: 1, Code: "c1", Name: "new_name", Type: "disease", Prior: 0.9}
	updateBody := `{"code":"c1","name":"new_name","type":"disease","prior":0.9}`

	tests := []struct {
		name       string
		id         string
		body       string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
		wantBody   *CauseResponse
	}{
		{
			name: "OK: id from URL is put into domain",
			id:   "1",
			body: updateBody,
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("UpdateCause", mock.Anything, updated).Return(&updated, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &updatedResp,
		},
		{
			name:       "wrong URLParam",
			id:         "meow",
			body:       updateBody,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			id:         "1",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation fail",
			id:         "1",
			body:       `{"code":"","name":"new_name","type":"disease","prior":0.9}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "not found",
			id:   "1",
			body: updateBody,
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("UpdateCause", mock.Anything, updated).Return(nil, domain.ErrNotFound)
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

			rec := httptestx.DoRequest(h.UpdateCause, http.MethodPut, tt.id, tt.body)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, *tt.wantBody, httptestx.DecodeBody[CauseResponse](t, rec))
			}
		})
	}
}

func TestDeleteCause(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setupMock  func(m *mocks.KnowledgebaseServiceInterface)
		wantStatus int
	}{
		{
			name: "OK",
			id:   "1",
			setupMock: func(m *mocks.KnowledgebaseServiceInterface) {
				m.On("DeleteCause", mock.Anything, 1).Return(nil)
			},
			wantStatus: http.StatusAccepted,
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
				m.On("DeleteCause", mock.Anything, 1).Return(domain.ErrNotFound)
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

			rec := httptestx.DoRequest(h.DeleteCause, http.MethodDelete, tt.id, "")

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
