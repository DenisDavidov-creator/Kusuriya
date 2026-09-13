package handler

import (
	"Kusuriya/internal/diagnose/handler/mocks"
	"Kusuriya/internal/domain"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDetermineDiagnose(t *testing.T) {
	remedyID := int64(7)
	remedyName := "antidote"

	fullBody := `{"symptom_codes":["s1","s2"],"factor_codes":["f1"]}`
	fullReq := domain.DiagnosisRequest{SymptomCodes: []string{"s1", "s2"}, FactorCodes: []string{"f1"}}

	results := []domain.DiagnosisResponse{
		{P: 0.9, CauseWithRemedy: domain.CauseWithRemedy{
			ID: 1, Code: "c1", Name: "c1_name", Type: "poison", Prior: 0.5,
			RemedyID: &remedyID, RemedyName: &remedyName,
		}},
		{P: 0.1, CauseWithRemedy: domain.CauseWithRemedy{
			ID: 2, Code: "c2", Name: "c2_name", Type: "disease", Prior: 0.2,
		}},
	}
	resultsResp := []DetermineDiagnoseResponse{
		{P: 0.9, Cause: CauseResponse{
			ID: 1, Code: "c1", Name: "c1_name", Type: "poison", Prior: 0.5,
			RemedyID: &remedyID, RemedyName: &remedyName,
		}},
		{P: 0.1, Cause: CauseResponse{
			ID: 2, Code: "c2", Name: "c2_name", Type: "disease", Prior: 0.2,
		}},
	}

	tests := []struct {
		name       string
		body       string
		setupMock  func(m *mocks.DiagnoseServiceInterface)
		wantStatus int
		wantBody   []DetermineDiagnoseResponse
	}{
		{
			name: "OK: symptoms and factors",
			body: fullBody,
			setupMock: func(m *mocks.DiagnoseServiceInterface) {
				m.On("DeterminDiagnose", mock.Anything, fullReq).Return(results, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   resultsResp,
		},
		{
			name: "OK: symptoms only is enough",
			body: `{"symptom_codes":["s1"]}`,
			setupMock: func(m *mocks.DiagnoseServiceInterface) {
				m.On("DeterminDiagnose", mock.Anything, domain.DiagnosisRequest{SymptomCodes: []string{"s1"}}).
					Return(nil, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []DetermineDiagnoseResponse{},
		},
		{
			name:       "invalid json",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation fail: no codes at all",
			body:       `{"symptom_codes":[],"factor_codes":[]}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "unexpected error",
			body: fullBody,
			setupMock: func(m *mocks.DiagnoseServiceInterface) {
				m.On("DeterminDiagnose", mock.Anything, fullReq).Return(nil, errors.New("Meow"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewDiagnoseServiceInterface(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			h := NewDiagnoseHandler(mockService)

			req := httptest.NewRequest(http.MethodPost, "/diagnose", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			h.DetermineDiagnose(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != nil {
				raw := strings.TrimSpace(rec.Body.String())
				var got []DetermineDiagnoseResponse
				require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
				assert.Equal(t, tt.wantBody, got)
				if len(tt.wantBody) == 0 {
					assert.Equal(t, "[]", raw)
				}
			}
		})
	}
}
