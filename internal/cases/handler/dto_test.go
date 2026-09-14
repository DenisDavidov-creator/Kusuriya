package handler

import (
	"Kusuriya/internal/diagnose/dto"
	"Kusuriya/internal/domain"
	"Kusuriya/internal/httpx"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func checkValidate(t *testing.T, wantProblems []string, problems []string, err error) {
	t.Helper()
	if len(wantProblems) == 0 {
		assert.Nil(t, problems)
		assert.NoError(t, err)
		return
	}
	assert.Equal(t, wantProblems, problems)
	assert.ErrorIs(t, err, httpx.ErrValidation)
}

func TestCreateCaseRequest(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		p, err := CreateCaseRequest{Title: "t"}.Validate()
		checkValidate(t, nil, p, err)

		p, err = CreateCaseRequest{Description: "no title"}.Validate()
		checkValidate(t, []string{"title is empty"}, p, err)
	})

	t.Run("ToDomain forces status open and never sets an id", func(t *testing.T) {
		got := CreateCaseRequest{Title: "t", Description: "d"}.ToDomain()
		assert.Equal(t, domain.Case{Title: "t", Description: "d", Status: "open"}, got)
	})
}

func TestCreateObservationRequest_Validate(t *testing.T) {
	tests := []struct {
		name string
		req  CreateObservationRequest
		want []string
	}{
		{"symptom", CreateObservationRequest{Kind: "symptom", RefID: 1}, nil},
		{"factor", CreateObservationRequest{Kind: "factor", RefID: 1}, nil},
		{"note is optional", CreateObservationRequest{Kind: "factor", RefID: 1, Note: ""}, nil},
		{"unknown kind", CreateObservationRequest{Kind: "banana", RefID: 1}, []string{"wrong kind check"}},
		{"kind is case sensitive", CreateObservationRequest{Kind: "Symptom", RefID: 1}, []string{"wrong kind check"}},
		{"missing ref_id", CreateObservationRequest{Kind: "symptom"}, []string{"RefID is empty"}},
		{"everything wrong", CreateObservationRequest{}, []string{"wrong kind check", "RefID is empty"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := tt.req.Validate()
			checkValidate(t, tt.want, p, err)
		})
	}
}

func TestCreateObservationRequest_ToDomain(t *testing.T) {
	got := CreateObservationRequest{Kind: "symptom", RefID: 3, Note: "n"}.ToDomain(5)

	assert.Equal(t, domain.Observation{CaseID: 5, Kind: "symptom", RefID: 3, Note: "n"}, got,
		"case id comes from the argument; ID and ObservedAt are left for the DB")
}

func TestDomainToCaseDetail(t *testing.T) {
	c := domain.Case{ID: 1, Title: "t", Description: "d", Status: "open"}
	w := 0.8
	diagnoses := []domain.DiagnosisResponse{{
		P:               0.9,
		CauseWithRemedy: domain.CauseWithRemedy{ID: 7, Code: "c1"},
		Evidence:        []domain.Contribution{{Kind: "symptom", KindID: 3, Code: "s1", Weight: &w, Matched: true}},
	}}

	got := DomainToCaseDetail(diagnoses, &c)

	assert.Equal(t, CaseResponseWithDiagnoses{
		Case: CaseResponse{ID: 1, Title: "t", Description: "d", Status: "open"},
		Causes: []dto.DetermineDiagnoseResponse{{
			P:        0.9,
			Cause:    dto.CauseResponse{ID: 7, Code: "c1"},
			Evidence: []dto.Evidence{{Kind: "symptom", KindID: 3, Code: "s1", Weight: &w, Matched: true}},
		}},
	}, got)

	assert.Equal(t, []dto.DetermineDiagnoseResponse{}, DomainToCaseDetail(nil, &c).Causes,
		"a case with no observations has an empty causes list, not nil")
}

func TestDomainsToObservations(t *testing.T) {
	at := time.Date(2026, 9, 14, 10, 0, 0, 0, time.UTC)
	in := []domain.Observation{
		{ID: 1, CaseID: 5, Kind: "symptom", RefID: 3, ObservedAt: at, Note: "n"},
		{ID: 2, CaseID: 5, Kind: "factor", RefID: 8, ObservedAt: at},
	}

	got := DomainsToObservations(in)

	assert.Equal(t, []ObservationsResponse{
		{ID: 1, CaseID: 5, Kind: "symptom", RefID: 3, ObservedAt: at, Note: "n"},
		{ID: 2, CaseID: 5, Kind: "factor", RefID: 8, ObservedAt: at},
	}, got)
	assert.Equal(t, []ObservationsResponse{}, DomainsToObservations(nil))
}

func TestObservationJSONFieldNames(t *testing.T) {
	at := time.Date(2026, 9, 14, 10, 0, 0, 0, time.UTC)
	b, err := json.Marshal(ObservationsResponse{ID: 1, CaseID: 5, Kind: "symptom", RefID: 3, ObservedAt: at, Note: "n"})
	assert.NoError(t, err)
	assert.JSONEq(t, `{"id":1,"case_id":5,"kind":"symptom","ref_id":3,"observed_at":"2026-09-14T10:00:00Z","note":"n"}`, string(b))
}
