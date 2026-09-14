package dto

import (
	"Kusuriya/internal/domain"
	"Kusuriya/internal/httpx"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineDiagnoseRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     DetermineDiagnoseRequest
		wantErr bool
	}{
		{"symptoms only", DetermineDiagnoseRequest{SymptomCodes: []string{"s1"}}, false},
		{"factors only", DetermineDiagnoseRequest{FactorCodes: []string{"f1"}}, false},
		{"both", DetermineDiagnoseRequest{SymptomCodes: []string{"s1"}, FactorCodes: []string{"f1"}}, false},
		{"both nil", DetermineDiagnoseRequest{}, true},
		{"both empty", DetermineDiagnoseRequest{SymptomCodes: []string{}, FactorCodes: []string{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.ErrorIs(t, err, httpx.ErrValidation)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDetermineDiagnoseRequest_ToDomain(t *testing.T) {
	r := DetermineDiagnoseRequest{SymptomCodes: []string{"s1", "s2"}, FactorCodes: []string{"f1"}}

	assert.Equal(t, domain.DiagnosisRequest{SymptomCodes: []string{"s1", "s2"}, FactorCodes: []string{"f1"}}, r.ToDomain())
}

func TestDomainToCauseWithRemedy(t *testing.T) {
	remedyID, name, desc, notes := int64(7), "antidote", "counter-agent", "keep close"

	t.Run("with remedy: every field is copied", func(t *testing.T) {
		d := domain.CauseWithRemedy{
			ID: 1, Code: "c1", Name: "n", Type: "poison", Prior: 0.5, Description: "d",
			RemedyID: &remedyID, RemedyName: &name, RemedyDescription: &desc, RemedyNotes: &notes,
		}
		assert.Equal(t, CauseResponse{
			ID: 1, Code: "c1", Name: "n", Type: "poison", Prior: 0.5, Description: "d",
			RemedyID: &remedyID, RemedyName: &name, RemedyDescription: &desc, RemedyNotes: &notes,
		}, DomainToCauseWithRemedy(&d))
	})

	t.Run("without remedy: remedy fields serialise as null", func(t *testing.T) {
		d := domain.CauseWithRemedy{ID: 1, Code: "c1"}
		b, err := json.Marshal(DomainToCauseWithRemedy(&d))
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"id":1,"code":"c1","name":"","type":"","prior":0,"description":"",
			"remedy_id":null,"remedy_name":null,"remedy_description":null,"remedy_notes":null
		}`, string(b))
	})
}

func TestDomainToEvidence(t *testing.T) {
	w := 0.8
	in := []domain.Contribution{
		{Kind: "symptom", KindID: 3, Code: "s1", Weight: &w, Matched: true},
		{Kind: "factor", KindID: 8, Code: "f1", Weight: nil, Matched: false},
	}

	got := DomainToEvidence(in)

	assert.Equal(t, []Evidence{
		{Kind: "symptom", KindID: 3, Code: "s1", Weight: &w, Matched: true},
		{Kind: "factor", KindID: 8, Code: "f1", Weight: nil, Matched: false},
	}, got)
	assert.Equal(t, []Evidence{}, DomainToEvidence(nil), "nil in -> empty slice out, so JSON is [] not null")
}

func TestDomainsToDetermineDiagnose(t *testing.T) {
	w := 0.8
	in := []domain.DiagnosisResponse{
		{
			P:               0.9,
			CauseWithRemedy: domain.CauseWithRemedy{ID: 1, Code: "c1"},
			Evidence:        []domain.Contribution{{Kind: "symptom", KindID: 3, Code: "s1", Weight: &w, Matched: true}},
		},
		{P: 0.1, CauseWithRemedy: domain.CauseWithRemedy{ID: 2, Code: "c2"}},
	}

	got := DomainsToDetermineDiagnose(in)

	assert.Equal(t, []DetermineDiagnoseResponse{
		{
			P:        0.9,
			Cause:    CauseResponse{ID: 1, Code: "c1"},
			Evidence: []Evidence{{Kind: "symptom", KindID: 3, Code: "s1", Weight: &w, Matched: true}},
		},
		{P: 0.1, Cause: CauseResponse{ID: 2, Code: "c2"}, Evidence: []Evidence{}},
	}, got, "order preserved; a cause without evidence gets an empty list, not nil")
	assert.Equal(t, []DetermineDiagnoseResponse{}, DomainsToDetermineDiagnose(nil))
}

func TestEvidenceJSONFieldNames(t *testing.T) {
	w := 0.8
	b, err := json.Marshal(Evidence{Kind: "symptom", KindID: 3, Code: "s1", Weight: &w, Matched: true})
	assert.NoError(t, err)
	assert.JSONEq(t, `{"kind":"symptom","kind_id":3,"code":"s1","weight":0.8,"matched":true}`, string(b))
}
