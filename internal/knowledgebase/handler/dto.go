package handler

import "Kusuriya/internal/domain"

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

func (r CreateCauseRequest) Validate() []string {
	var problems = make([]string, 0, 2)

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
	return problems
}
