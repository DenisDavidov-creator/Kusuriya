package dto

import (
	"Kusuriya/internal/domain"
	"Kusuriya/internal/httpx"
)

type DetermineDiagnoseRequest struct {
	SymptomCodes []string `json:"symptom_codes" example:"SYM-DILA,SYM-HALL,SYM-DELI"`
	FactorCodes  []string `json:"factor_codes" example:"FAC-FORAGED"`
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
	P        float64       `json:"p" example:"0.5167882008816377"`
	Cause    CauseResponse `json:"cause"`
	Evidence []Evidence    `json:"evidence"`
}

// Evidence is one term of the score: a symptom or factor the cause was judged on.
// Weight is the likelihood/modifier that was multiplied in, or null when the cause
// has no profile for that code (Matched == false).
type Evidence struct {
	Kind    string   `json:"kind" example:"symptom" enums:"symptom,factor"`
	KindID  int64    `json:"kind_id" example:"38"`
	Code    string   `json:"code" example:"SYM-DILA"`
	Weight  *float64 `json:"weight" example:"0.85"`
	Matched bool     `json:"matched" example:"true"`
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
	ID                int64   `json:"id" example:"8"`
	Code              string  `json:"code" example:"POI-DATU"`
	Name              string  `json:"name" example:"Datura (thorn-apple)"`
	Type              string  `json:"type" example:"poison" enums:"poison,disease"`
	Prior             float64 `json:"prior" example:"0.008"`
	Description       string  `json:"description" example:"Delirium-inducing seeds and flowers."`
	RemedyID          *int64  `json:"remedy_id" example:"11"`
	RemedyName        *string `json:"remedy_name" example:"Palace Physician's Antidote"`
	RemedyDescription *string `json:"remedy_description" example:"A prepared counter-agent for slowed heart and clouded mind."`
	RemedyNotes       *string `json:"remedy_notes" example:"Kept close at every banquet."`
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
