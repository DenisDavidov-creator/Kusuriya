package handler

import (
	"Kusuriya/internal/domain"
)

type CreateCauseRequest struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Prior       float64 `json:"prior"`
	Description string  `json:"description"`
	RemedyID    *int64  `json:"remedy_id"`
}

func (r CreateCauseRequest) ToDomain() domain.Cause {
	return domain.Cause{
		Code:        r.Code,
		Name:        r.Name,
		Type:        r.Type,
		Prior:       r.Prior,
		Description: r.Description,
		RemedyID:    r.RemedyID,
	}
}
func (r CreateCauseRequest) Validate() ([]string, error) {
	var problems = make([]string, 0, 4)

	if r.Code == "" {
		problems = append(problems, "code is empty")
	}
	if r.Name == "" {
		problems = append(problems, "name is empty")
	}
	if r.Type != "poison" && r.Type != "disease" {
		problems = append(problems, "type is incorrect")
	}
	if r.Prior > 1 || r.Prior <= 0 {
		problems = append(problems, "prior out of range (0,1]")
	}
	if len(problems) != 0 {
		return problems, ErrValidation
	}
	return nil, nil
}

type CauseResponse struct {
	ID          int64   `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Prior       float64 `json:"prior"`
	Description string  `json:"description"`
	RemedyID    *int64  `json:"remedy_id"`
}

func DomainToCause(d *domain.Cause) CauseResponse {
	return CauseResponse{
		ID:          d.ID,
		Code:        d.Code,
		Name:        d.Name,
		Type:        d.Type,
		Prior:       d.Prior,
		Description: d.Description,
		RemedyID:    d.RemedyID,
	}
}

func DomainsToCauses(d []domain.Cause) []CauseResponse {

	crs := make([]CauseResponse, 0, len(d))
	for _, dom := range d {
		cr := DomainToCause(&dom)
		crs = append(crs, cr)
	}
	return crs
}

type UpdateCauseRequest struct {
	ID          int64   `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Prior       float64 `json:"prior"`
	Description string  `json:"description"`
	RemedyID    *int64  `json:"remedy_id"`
}

func (r UpdateCauseRequest) Validate() ([]string, error) {
	var problems = make([]string, 0, 4)
	if r.Code == "" {
		problems = append(problems, "code is empty")
	}
	if r.Name == "" {
		problems = append(problems, "name is empty")
	}
	if r.Type != "poison" && r.Type != "disease" {
		problems = append(problems, "type is incorrect")
	}
	if r.Prior > 1 || r.Prior <= 0 {
		problems = append(problems, "prior out of range (0,1]")
	}
	if len(problems) != 0 {
		return problems, ErrValidation
	}
	return nil, nil
}

func (r *UpdateCauseRequest) ToDomain(id int) domain.Cause {
	return domain.Cause{
		ID:          int64(id),
		Code:        r.Code,
		Name:        r.Name,
		Type:        r.Type,
		Prior:       r.Prior,
		Description: r.Description,
		RemedyID:    r.RemedyID,
	}
}
