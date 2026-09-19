package handler

import (
	"Kusuriya/internal/diagnose/dto"
	"Kusuriya/internal/domain"
	"Kusuriya/internal/httpx"
	"time"
)

type CreateCaseRequest struct {
	Title       string `json:"title" example:"Consort Lihua's chambermaid" binding:"required"`
	Description string `json:"description" example:"Collapsed after the spring banquet."`
}

func (r CreateCaseRequest) ToDomain() domain.Case {
	return domain.Case{
		Title:       r.Title,
		Description: r.Description,
		Status:      "open",
	}
}

func (r CreateCaseRequest) Validate() ([]string, error) {
	var problems = make([]string, 0, 1)

	if r.Title == "" {
		problems = append(problems, "title is empty")
	}
	if len(problems) != 0 {
		return problems, httpx.ErrValidation
	}
	return nil, nil
}

type CaseResponse struct {
	ID          int64  `json:"id" example:"1"`
	Title       string `json:"title" example:"Consort Lihua's chambermaid"`
	Description string `json:"description" example:"Collapsed after the spring banquet."`
	Status      string `json:"status" example:"open" enums:"open,closed"`
}

type CaseResponseWithDiagnoses struct {
	Case   CaseResponse                    `json:"case"`
	Causes []dto.DetermineDiagnoseResponse `json:"causes"`
}

func DomainToCaseDetail(diagnoses []domain.DiagnosisResponse, caseDetail *domain.Case) CaseResponseWithDiagnoses {
	return CaseResponseWithDiagnoses{
		Case:   DomainToCase(caseDetail),
		Causes: dto.DomainsToDetermineDiagnose(diagnoses),
	}
}

func DomainToCase(d *domain.Case) CaseResponse {
	return CaseResponse{
		ID:          d.ID,
		Title:       d.Title,
		Description: d.Description,
		Status:      d.Status,
	}
}

type CreateObservationRequest struct {
	Kind  string `json:"kind" example:"symptom" enums:"symptom,factor" binding:"required"`
	RefID int64  `json:"ref_id" example:"38" binding:"required"`
	Note  string `json:"note" example:"Pupils still wide by evening."`
}

func (r CreateObservationRequest) ToDomain(CaseID int) domain.Observation {
	return domain.Observation{
		CaseID: int64(CaseID),
		Kind:   r.Kind,
		RefID:  r.RefID,
		Note:   r.Note,
	}
}

func (r CreateObservationRequest) Validate() ([]string, error) {
	var problems = make([]string, 0, 2)

	if r.Kind != "symptom" && r.Kind != "factor" {
		problems = append(problems, "wrong kind check")
	}
	if r.RefID == 0 {
		problems = append(problems, "RefID is empty")
	}

	if len(problems) != 0 {
		return problems, httpx.ErrValidation
	}
	return nil, nil
}

type ObservationsResponse struct {
	ID         int64     `json:"id" example:"1"`
	CaseID     int64     `json:"case_id" example:"1"`
	Kind       string    `json:"kind" example:"symptom" enums:"symptom,factor"`
	RefID      int64     `json:"ref_id" example:"38"`
	ObservedAt time.Time `json:"observed_at" example:"2026-09-19T10:00:00Z"`
	Note       string    `json:"note" example:"Pupils still wide by evening."`
}

func DomainToObservation(d *domain.Observation) ObservationsResponse {
	return ObservationsResponse{
		ID:         d.ID,
		CaseID:     d.CaseID,
		Kind:       d.Kind,
		RefID:      d.RefID,
		ObservedAt: d.ObservedAt,
		Note:       d.Note,
	}
}

func DomainsToObservations(doms []domain.Observation) []ObservationsResponse {
	var results = make([]ObservationsResponse, 0, len(doms))
	for _, dom := range doms {
		result := DomainToObservation(&dom)
		results = append(results, result)
	}

	return results
}
