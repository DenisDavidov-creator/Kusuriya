package dto

import (
	"Kusuriya/internal/domain"
	"Kusuriya/internal/httpx"
)

type DetermineDiagnoseRequest struct {
	SymptomCodes []string `json:"symptom_codes"`
	FactorCodes  []string `json:"factor_codes"`
}

func (r DetermineDiagnoseRequest) ToDomain() domain.DiagnosisRequest {
	return domain.DiagnosisRequest{
		SymptomCodes: r.SymptomCodes,
		FactorCodes:  r.FactorCodes,
	}
}

func (r DetermineDiagnoseRequest) Validate() error {
	if len(r.FactorCodes)+len(r.SymptomCodes) == 0 {
		return httpx.ErrValidation
	}
	return nil
}

type DetermineDiagnoseResponse struct {
	P        float64       `json:"p"`
	Cause    CauseResponse `json:"cause"`
	Evidence []Evidence    `json:"evidence"`
}

type Evidence struct {
	Kind    string   `json:"kind"`
	KindID  int64    `json:"kind_id"`
	Code    string   `json:"code"`
	Weight  *float64 `json:"weight"`
	Matched bool     `json:"matched"`
}

func DomainToEvidence(d []domain.Contribution) []Evidence {
	var evidence = make([]Evidence, 0, len(d))
	for _, e := range d {
		ev := Evidence{
			Kind:    e.Kind,
			KindID:  e.KindID,
			Code:    e.Code,
			Weight:  e.Weight,
			Matched: e.Matched,
		}
		evidence = append(evidence, ev)
	}
	return evidence
}

type CauseResponse struct {
	ID                int64   `json:"id"`
	Code              string  `json:"code"`
	Name              string  `json:"name"`
	Type              string  `json:"type"`
	Prior             float64 `json:"prior"`
	Description       string  `json:"description"`
	RemedyID          *int64  `json:"remedy_id"`
	RemedyName        *string `json:"remedy_name"`
	RemedyDescription *string `json:"remedy_description"`
	RemedyNotes       *string `json:"remedy_notes"`
}

func DomainToCauseWithRemedy(d *domain.CauseWithRemedy) CauseResponse {
	return CauseResponse{
		ID:                d.ID,
		Code:              d.Code,
		Name:              d.Name,
		Type:              d.Type,
		Prior:             d.Prior,
		Description:       d.Description,
		RemedyID:          d.RemedyID,
		RemedyName:        d.RemedyName,
		RemedyDescription: d.RemedyDescription,
		RemedyNotes:       d.RemedyNotes,
	}
}

func DomainsToDetermineDiagnose(doms []domain.DiagnosisResponse) []DetermineDiagnoseResponse {
	var results = make([]DetermineDiagnoseResponse, 0, len(doms))
	for _, dom := range doms {
		result := DetermineDiagnoseResponse{
			P:        dom.P,
			Cause:    DomainToCauseWithRemedy(&dom.CauseWithRemedy),
			Evidence: DomainToEvidence(dom.Evidence),
		}
		results = append(results, result)
	}

	return results
}
