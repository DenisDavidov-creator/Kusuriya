package handler

import (
	"Kusuriya/internal/domain"
	"Kusuriya/internal/httpx"
	"fmt"
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
		return problems, httpx.ErrValidation
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

func DomainsToCauses(doms []domain.Cause) []CauseResponse {
	res := make([]CauseResponse, 0, len(doms))
	for _, dom := range doms {
		r := DomainToCause(&dom)
		res = append(res, r)
	}
	return res
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
		return problems, httpx.ErrValidation
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

type CreateRemedyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
}

type RemedyResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
}

func (r *CreateRemedyRequest) ToDomain() domain.Remedy {
	return domain.Remedy{
		Name:        r.Name,
		Description: r.Description,
		Notes:       r.Notes,
	}
}

func (r *CreateRemedyRequest) Validate() ([]string, error) {
	problems := make([]string, 0, 1)

	if r.Name == "" {
		problems = append(problems, "name is empty")
	}
	if len(problems) != 0 {
		return problems, httpx.ErrValidation
	}
	return nil, nil
}

func DomainToRemedy(d *domain.Remedy) RemedyResponse {
	return RemedyResponse{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Notes:       d.Notes,
	}
}
func DomainsToRemedies(doms []domain.Remedy) []RemedyResponse {
	var res = make([]RemedyResponse, 0, len(doms))
	for _, dom := range doms {
		r := DomainToRemedy(&dom)
		res = append(res, r)
	}
	return res
}

type CreateFactorRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r CreateFactorRequest) ToDomain() domain.Factor {
	return domain.Factor{
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
	}
}
func (r CreateFactorRequest) Validate() ([]string, error) {
	problems := make([]string, 0, 2)
	if r.Code == "" {
		problems = append(problems, "code is empty")
	}
	if r.Name == "" {
		problems = append(problems, "name is empty")
	}
	if len(problems) != 0 {
		return problems, httpx.ErrValidation
	}
	return nil, nil
}

type FactorResponse struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func DomainToFactor(d *domain.Factor) FactorResponse {
	return FactorResponse{
		ID:          d.ID,
		Code:        d.Code,
		Name:        d.Name,
		Description: d.Description,
	}
}

func DomainsToFactors(doms []domain.Factor) []FactorResponse {
	var res = make([]FactorResponse, 0, len(doms))
	for _, dom := range doms {
		r := DomainToFactor(&dom)
		res = append(res, r)
	}
	return res
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func (r CreateCategoryRequest) ToDomain() domain.Category {
	return domain.Category{
		Name: r.Name,
	}
}
func (r CreateCategoryRequest) Validate() ([]string, error) {
	problems := make([]string, 0, 1)
	if r.Name == "" {
		problems = append(problems, "name is empty")
	}
	if len(problems) != 0 {
		return problems, httpx.ErrValidation
	}
	return nil, nil
}

type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func DomainToCategory(d *domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:   d.ID,
		Name: d.Name,
	}
}
func DomainsToCategories(doms []domain.Category) []CategoryResponse {
	var res = make([]CategoryResponse, 0, len(doms))
	for _, dom := range doms {
		r := DomainToCategory(&dom)
		res = append(res, r)
	}
	return res
}

type CreateSymptomRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryID  int64  `json:"category_id"`
}

func (r CreateSymptomRequest) ToDomain() domain.Symptom {
	return domain.Symptom{
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		CategoryID:  r.CategoryID,
	}
}

func (r CreateSymptomRequest) Validate() ([]string, error) {
	var problems = make([]string, 0, 4)
	if r.Code == "" {
		problems = append(problems, "code is empty")
	}
	if r.Name == "" {
		problems = append(problems, "name is empty")
	}
	if r.CategoryID == 0 {
		problems = append(problems, "category is empty")
	}
	if len(problems) != 0 {
		return problems, httpx.ErrValidation
	}
	return nil, nil
}

type SymptomResponse struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryID  int64  `json:"category_id"`
}

func DomainToSymptom(d *domain.Symptom) SymptomResponse {
	return SymptomResponse{
		ID:          d.ID,
		Code:        d.Code,
		Name:        d.Name,
		Description: d.Description,
		CategoryID:  d.CategoryID,
	}
}

func DomainsToSymptoms(doms []domain.Symptom) []SymptomResponse {
	var res = make([]SymptomResponse, 0, len(doms))
	for _, dom := range doms {
		r := DomainToSymptom(&dom)
		res = append(res, r)
	}
	return res
}

type CreateSymptomProfileRequest struct {
	SymptomID  int64   `json:"symptom_id"`
	Likelihood float64 `json:"likelihood"`
}
type CreateSymptomsProfileRequests struct {
	Profiles []CreateSymptomProfileRequest `json:"profiles"`
}

func (rs CreateSymptomsProfileRequests) ToDomains(id int64) []domain.SymptomProfile {
	var doms []domain.SymptomProfile
	for _, r := range rs.Profiles {
		dom := domain.SymptomProfile{
			CauseID:    id,
			SymptomID:  r.SymptomID,
			Likelihood: r.Likelihood,
		}
		doms = append(doms, dom)
	}
	return doms
}

func (rs CreateSymptomsProfileRequests) Validate() ([]string, error) {
	var problems = make([]string, 0, 3)
	for _, r := range rs.Profiles {

		if r.SymptomID == 0 {
			problems = append(problems, fmt.Sprintf("%d: symptom_id is empty", r.SymptomID))
		}
		if r.Likelihood > 1 || r.Likelihood <= 0 {
			problems = append(problems, fmt.Sprintf("%d: likelihood is empty", r.SymptomID))
		}
		if len(problems) != 0 {
			return problems, httpx.ErrValidation
		}
	}
	return nil, nil
}

type SymptomProfileResponse struct {
	ID         int64   `json:"id"`
	CauseID    int64   `json:"cause_id"`
	SymptomID  int64   `json:"symptom_id"`
	Likelihood float64 `json:"likelihood"`
}

func DomainsToSymptomsProfiles(doms []domain.SymptomProfile) []SymptomProfileResponse {
	var res = make([]SymptomProfileResponse, 0, len(doms))
	for _, d := range doms {
		r := SymptomProfileResponse{
			ID:         d.ID,
			CauseID:    d.CauseID,
			SymptomID:  d.SymptomID,
			Likelihood: d.Likelihood,
		}
		res = append(res, r)
	}
	return res
}

type CreateFactorProfileRequest struct {
	FactorID int64   `json:"factor_id"`
	Modifier float64 `json:"modifier"`
}
type CreateFactorsProfileRequests struct {
	Profiles []CreateFactorProfileRequest `json:"profiles"`
}

func (rs CreateFactorsProfileRequests) ToDomains(id int64) []domain.FactorProfile {
	var doms []domain.FactorProfile
	for _, r := range rs.Profiles {
		dom := domain.FactorProfile{
			CauseID:  id,
			FactorID: r.FactorID,
			Modifier: r.Modifier,
		}
		doms = append(doms, dom)
	}
	return doms
}

func (rs CreateFactorsProfileRequests) Validate() ([]string, error) {
	var problems = make([]string, 0, 3)
	for _, r := range rs.Profiles {

		if r.FactorID == 0 {
			problems = append(problems, fmt.Sprintf("%d: Factor_id is empty", r.FactorID))
		}
		if r.Modifier > 1 || r.Modifier <= 0 {
			problems = append(problems, fmt.Sprintf("%d: modifier is empty", r.FactorID))
		}
		if len(problems) != 0 {
			return problems, httpx.ErrValidation
		}
	}
	return nil, nil
}

type FactorProfileResponse struct {
	ID       int64   `json:"id"`
	CauseID  int64   `json:"cause_id"`
	FactorID int64   `json:"Factor_id"`
	Modifier float64 `json:"modifier"`
}

func DomainsToFactorsProfiles(doms []domain.FactorProfile) []FactorProfileResponse {
	var res = make([]FactorProfileResponse, 0, len(doms))
	for _, d := range doms {
		r := FactorProfileResponse{
			ID:       d.ID,
			CauseID:  d.CauseID,
			FactorID: d.FactorID,
			Modifier: d.Modifier,
		}
		res = append(res, r)
	}
	return res
}
