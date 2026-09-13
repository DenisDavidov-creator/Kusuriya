package handler

import (
	"Kusuriya/internal/cases/handler/mocks"
	"Kusuriya/internal/diagnose/dto"
	"Kusuriya/internal/domain"
	"Kusuriya/internal/httpx"
	"Kusuriya/internal/httpx/httptestx"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	openCase     = domain.Case{ID: 1, Title: "fever in the east wing", Description: "desc", Status: "open"}
	openCaseResp = CaseResponse{ID: 1, Title: "fever in the east wing", Description: "desc", Status: "open"}
)

func TestCreateCase(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		setupMock  func(m *mocks.CaseServiceInterface)
		wantStatus int
		wantMsg    string // error message expected in the body, "" = don't check
		wantBody   *CaseResponse
	}{
		{
			name: "OK: status is forced to open",
			body: `{"title":"fever in the east wing","description":"desc"}`,
			setupMock: func(m *mocks.CaseServiceInterface) {
				m.On("CreateCase", mock.Anything, domain.Case{
					Title: "fever in the east wing", Description: "desc", Status: "open",
				}).Return(&openCase, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &openCaseResp,
		},
		{
			name:       "invalid json",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "bad request",
		},
		{
			name:       "validation fail: empty title",
			body:       `{"title":"","description":"desc"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "validation not passed",
		},
		{
			name: "unexpected error",
			body: `{"title":"fever in the east wing"}`,
			setupMock: func(m *mocks.CaseServiceInterface) {
				m.On("CreateCase", mock.Anything, mock.Anything).Return(nil, errors.New("Meow"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewCaseServiceInterface(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			h := NewCaseHandler(mockService)

			rec := httptestx.DoRequest(h.CreateCase, http.MethodPost, "", tt.body)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, *tt.wantBody, httptestx.DecodeBody[CaseResponse](t, rec))
			}
			if tt.wantMsg != "" {
				assert.Equal(t, tt.wantMsg, httptestx.DecodeBody[httpx.ErrorResponse](t, rec).Error.Message)
			}
		})
	}
}

func TestGetCaseByID(t *testing.T) {
	weight := 0.8
	diagnoses := []domain.DiagnosisResponse{
		{
			P:               0.9,
			CauseWithRemedy: domain.CauseWithRemedy{ID: 7, Code: "c1", Name: "c1_name", Type: "poison", Prior: 0.5},
			Evidence: []domain.Contribution{
				{Kind: "symptom", KindID: 3, Code: "s1", Weight: &weight, Matched: true},
			},
		},
	}
	want := CaseResponseWithDiagnoses{
		Case: openCaseResp,
		Causes: []dto.DetermineDiagnoseResponse{
			{
				P:     0.9,
				Cause: dto.CauseResponse{ID: 7, Code: "c1", Name: "c1_name", Type: "poison", Prior: 0.5},
				Evidence: []dto.Evidence{
					{Kind: "symptom", KindID: 3, Code: "s1", Weight: &weight, Matched: true},
				},
			},
		},
	}

	tests := []struct {
		name       string
		id         string
		setupMock  func(m *mocks.CaseServiceInterface)
		wantStatus int
		wantBody   *CaseResponseWithDiagnoses
	}{
		{
			name: "OK: case with diagnoses and evidence",
			id:   "1",
			setupMock: func(m *mocks.CaseServiceInterface) {
				m.On("GetCaseByID", mock.Anything, 1).Return(diagnoses, &openCase, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &want,
		},
		{
			name: "OK: case without observations has empty causes, not null",
			id:   "1",
			setupMock: func(m *mocks.CaseServiceInterface) {
				m.On("GetCaseByID", mock.Anything, 1).Return(nil, &openCase, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &CaseResponseWithDiagnoses{Case: openCaseResp, Causes: []dto.DetermineDiagnoseResponse{}},
		},
		{
			name:       "wrong URLParam",
			id:         "meow",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "not found",
			id:   "1",
			setupMock: func(m *mocks.CaseServiceInterface) {
				m.On("GetCaseByID", mock.Anything, 1).Return(nil, nil, domain.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewCaseServiceInterface(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			h := NewCaseHandler(mockService)

			rec := httptestx.DoRequest(h.GetCaseByID, http.MethodGet, tt.id, "")

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				raw := rec.Body.String()
				assert.Equal(t, *tt.wantBody, httptestx.DecodeBody[CaseResponseWithDiagnoses](t, rec))
				if len(tt.wantBody.Causes) == 0 {
					assert.Contains(t, raw, `"causes":[]`)
				}
			}
		})
	}
}

func TestCloseCaseByID(t *testing.T) {
	closedCase := domain.Case{ID: 1, Title: "fever in the east wing", Description: "desc", Status: "closed"}
	closedCaseResp := CaseResponse{ID: 1, Title: "fever in the east wing", Description: "desc", Status: "closed"}

	tests := []struct {
		name       string
		id         string
		setupMock  func(m *mocks.CaseServiceInterface)
		wantStatus int
		wantBody   *CaseResponse
	}{
		{
			name: "OK",
			id:   "1",
			setupMock: func(m *mocks.CaseServiceInterface) {
				m.On("CloseCaseByID", mock.Anything, 1).Return(&closedCase, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &closedCaseResp,
		},
		{
			name:       "wrong URLParam",
			id:         "meow",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "not found",
			id:   "1",
			setupMock: func(m *mocks.CaseServiceInterface) {
				m.On("CloseCaseByID", mock.Anything, 1).Return(nil, domain.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewCaseServiceInterface(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			h := NewCaseHandler(mockService)

			rec := httptestx.DoRequest(h.CloseCaseByID, http.MethodPost, tt.id, "")

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, *tt.wantBody, httptestx.DecodeBody[CaseResponse](t, rec))
			}
		})
	}
}

func TestCreateObservation(t *testing.T) {
	observedAt := time.Date(2026, 9, 13, 10, 0, 0, 0, time.UTC)
	in := domain.Observation{CaseID: 5, Kind: "symptom", RefID: 3, Note: "high fever"}
	saved := domain.Observation{ID: 1, CaseID: 5, Kind: "symptom", RefID: 3, ObservedAt: observedAt, Note: "high fever"}
	savedResp := ObservationsResponse{ID: 1, CaseID: 5, Kind: "symptom", RefID: 3, ObservedAt: observedAt, Note: "high fever"}

	tests := []struct {
		name       string
		id         string
		body       string
		setupMock  func(m *mocks.CaseServiceInterface)
		wantStatus int
		wantBody   *ObservationsResponse
	}{
		{
			name: "OK: case id from URL goes into the observation",
			id:   "5",
			body: `{"kind":"symptom","ref_id":3,"note":"high fever"}`,
			setupMock: func(m *mocks.CaseServiceInterface) {
				m.On("CreateObservation", mock.Anything, in).Return(&saved, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &savedResp,
		},
		{
			name:       "wrong URLParam",
			id:         "meow",
			body:       `{"kind":"symptom","ref_id":3}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			id:         "5",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation fail: unknown kind",
			id:         "5",
			body:       `{"kind":"banana","ref_id":3}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation fail: missing ref_id",
			id:         "5",
			body:       `{"kind":"factor"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "case not found",
			id:   "5",
			body: `{"kind":"symptom","ref_id":3,"note":"high fever"}`,
			setupMock: func(m *mocks.CaseServiceInterface) {
				m.On("CreateObservation", mock.Anything, in).Return(nil, domain.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewCaseServiceInterface(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			h := NewCaseHandler(mockService)

			rec := httptestx.DoRequest(h.CreateObservation, http.MethodPost, tt.id, tt.body)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				assert.Equal(t, *tt.wantBody, httptestx.DecodeBody[ObservationsResponse](t, rec))
			}
		})
	}
}

func TestGetObservationByID(t *testing.T) {
	observedAt := time.Date(2026, 9, 13, 10, 0, 0, 0, time.UTC)
	obs := []domain.Observation{
		{ID: 1, CaseID: 5, Kind: "symptom", RefID: 3, ObservedAt: observedAt, Note: "high fever"},
		{ID: 2, CaseID: 5, Kind: "factor", RefID: 8, ObservedAt: observedAt},
	}
	obsResp := []ObservationsResponse{
		{ID: 1, CaseID: 5, Kind: "symptom", RefID: 3, ObservedAt: observedAt, Note: "high fever"},
		{ID: 2, CaseID: 5, Kind: "factor", RefID: 8, ObservedAt: observedAt},
	}

	tests := []struct {
		name       string
		id         string
		setupMock  func(m *mocks.CaseServiceInterface)
		wantStatus int
		wantBody   []ObservationsResponse
	}{
		{
			name: "OK",
			id:   "5",
			setupMock: func(m *mocks.CaseServiceInterface) {
				m.On("GetObservationByID", mock.Anything, 5).Return(obs, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   obsResp,
		},
		{
			name: "empty list is [] not null",
			id:   "5",
			setupMock: func(m *mocks.CaseServiceInterface) {
				m.On("GetObservationByID", mock.Anything, 5).Return(nil, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []ObservationsResponse{},
		},
		{
			name:       "wrong URLParam",
			id:         "meow",
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewCaseServiceInterface(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			h := NewCaseHandler(mockService)

			rec := httptestx.DoRequest(h.GetObservationByID, http.MethodGet, tt.id, "")

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				raw := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, tt.wantBody, httptestx.DecodeBody[[]ObservationsResponse](t, rec))
				if len(tt.wantBody) == 0 {
					assert.Equal(t, "[]", raw)
				}
			}
		})
	}
}
