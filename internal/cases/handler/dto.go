package handler

import (
	"Kusuriya/internal/diagnose/dto"
	"Kusuriya/internal/domain"
	"Kusuriya/internal/httpx"
	"time"
)

type CreateCaseRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
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
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
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
	Kind  string `json:"kind"`
	RefID int64  `json:"ref_id"`
	Note  string `json:"note"`
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
	ID         int64     `json:"id"`
	CaseID     int64     `json:"case_id"`
	Kind       string    `json:"kind"`
	RefID      int64     `json:"ref_id"`
	ObservedAt time.Time `json:"observed_at"`
	Note       string    `json:"note"`
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
