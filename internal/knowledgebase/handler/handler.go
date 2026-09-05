package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type KnowledgebaseHandler struct {
	service KnowledgebaseServiceInterface
}

func NewKnowledgebaseHandler(service KnowledgebaseServiceInterface) *KnowledgebaseHandler {
	return &KnowledgebaseHandler{service: service}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, r *http.Request, err error, details any) {
	status, level, message := MapErr(err)
	slog.Log(r.Context(), level, "request failed", "error", err)
	writeJSON(w, status, NewErrorResponse(status, message, details))
}

func decode(r io.ReadCloser, req any) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	return dec.Decode(req)
}

func (h *KnowledgebaseHandler) CreateCause(w http.ResponseWriter, r *http.Request) {
	var req CreateCauseRequest
	if err := decode(r.Body, &req); err != nil {
		respondError(w, r, err, err.Error())
		return
	}

	problems, err := req.Validate()
	if len(problems) != 0 {
		respondError(w, r, err, problems)
		return
	}

	dom, err := h.service.CreateCause(r.Context(), req.ToDomain())
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainToCause(dom))
}

func (h *KnowledgebaseHandler) GetCauseByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, r, err, "wrong URLParam")
		return
	}

	dom, err := h.service.GetCauseByID(r.Context(), id)
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainToCause(dom))
}

func (h *KnowledgebaseHandler) GetCauses(w http.ResponseWriter, r *http.Request) {
	doms, err := h.service.GetCauses(r.Context())
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainsToCauses(doms))
}

func (h *KnowledgebaseHandler) UpdateCause(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, r, err, "wrong URLParam")
		return
	}

	var req UpdateCauseRequest
	if err := decode(r.Body, &req); err != nil {
		respondError(w, r, err, err.Error())
		return
	}

	problems, err := req.Validate()
	if len(problems) != 0 {
		respondError(w, r, err, problems)
		return
	}

	dom, err := h.service.UpdateCause(r.Context(), req.ToDomain(id))
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainToCause(dom))
}

func (h *KnowledgebaseHandler) DeleteCause(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, r, err, "wrong URLParam")
		return
	}

	err = h.service.DeleteCause(r.Context(), id)
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusAccepted, "cause deleted")
}

func (h *KnowledgebaseHandler) CreateRemedy(w http.ResponseWriter, r *http.Request) {
	var req CreateRemedyRequest
	if err := decode(r.Body, &req); err != nil {
		respondError(w, r, err, err.Error())
		return
	}
	problems, err := req.Validate()
	if err != nil {
		respondError(w, r, err, problems)
		return
	}

	dom, err := h.service.CreateRemedy(r.Context(), req.ToDomain())
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainToRemedy(dom))
}

func (h *KnowledgebaseHandler) GetRemedyByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, r, err, "wrong URLParam")
		return
	}

	dom, err := h.service.GetRemedyByID(r.Context(), id)
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainToRemedy(dom))
}

func (h *KnowledgebaseHandler) GetRemedies(w http.ResponseWriter, r *http.Request) {

	doms, err := h.service.GetRemedies(r.Context())
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainsToRemedies(doms))
}

func (h *KnowledgebaseHandler) CreateFactor(w http.ResponseWriter, r *http.Request) {
	var req CreateFactorRequest
	if err := decode(r.Body, &req); err != nil {
		respondError(w, r, err, err.Error())
		return
	}
	problems, err := req.Validate()
	if err != nil {
		respondError(w, r, err, problems)
		return
	}

	dom, err := h.service.CreateFactor(r.Context(), req.ToDomain())
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainToFactor(dom))
}

func (h *KnowledgebaseHandler) GetFactorByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, r, err, "wrong URLParam")
		return
	}

	dom, err := h.service.GetFactorByID(r.Context(), id)
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainToFactor(dom))
}

func (h *KnowledgebaseHandler) GetFactors(w http.ResponseWriter, r *http.Request) {

	doms, err := h.service.GetFactors(r.Context())
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainsToFactors(doms))
}

func (h *KnowledgebaseHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	err := decode(r.Body, &req)

	if err != nil {
		respondError(w, r, err, err.Error())
		return
	}

	problmes, err := req.Validate()
	if err != nil {
		respondError(w, r, err, problmes)
		return
	}

	dom, err := h.service.CreateCategory(r.Context(), req.ToDomain())
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainToCategory(dom))
}

func (h *KnowledgebaseHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, r, err, "wrong URLParam")
		return
	}

	dom, err := h.service.GetCategoryByID(r.Context(), id)
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainToCategory(dom))
}

func (h *KnowledgebaseHandler) GetCategories(w http.ResponseWriter, r *http.Request) {

	doms, err := h.service.GetCategories(r.Context())
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainsToCategories(doms))
}

func (h *KnowledgebaseHandler) CreateSymptom(w http.ResponseWriter, r *http.Request) {
	var req CreateSymptomRequest
	err := decode(r.Body, &req)

	if err != nil {
		respondError(w, r, err, err.Error())
		return
	}

	problmes, err := req.Validate()
	if err != nil {
		respondError(w, r, err, problmes)
		return
	}

	dom, err := h.service.CreateSymptom(r.Context(), req.ToDomain())
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainToSymptom(dom))
}

func (h *KnowledgebaseHandler) GetSymptomByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, r, err, "wrong URLParam")
		return
	}

	dom, err := h.service.GetSymptomByID(r.Context(), id)
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainToSymptom(dom))
}

func (h *KnowledgebaseHandler) GetSymptoms(w http.ResponseWriter, r *http.Request) {

	doms, err := h.service.GetSymptoms(r.Context())
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, DomainsToSymptoms(doms))
}

func (h *KnowledgebaseHandler) CreateSymptomProfiles(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, r, err, "wrongURLParam")
		return
	}
	var req CreateSymptomsProfileRequests
	err = decode(r.Body, &req)
	if err != nil {
		respondError(w, r, err, err.Error())
		return
	}
	problems, err := req.Validate()
	if err != nil {
		respondError(w, r, err, problems)
		return
	}
	doms, err := h.service.CreateSymptomProfiles(r.Context(), req.ToDomains(int64(id)), id)
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, doms)
}

func (h *KnowledgebaseHandler) CreateFactorProfiles(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, r, err, "wrongURLParam")
		return
	}
	var req CreateFactorsProfileRequests
	err = decode(r.Body, &req)
	if err != nil {
		respondError(w, r, err, err.Error())
		return
	}
	problems, err := req.Validate()
	if err != nil {
		respondError(w, r, err, problems)
		return
	}
	doms, err := h.service.CreateFactorProfiles(r.Context(), req.ToDomains(int64(id)), id)
	if err != nil {
		respondError(w, r, err, "")
		return
	}
	writeJSON(w, http.StatusOK, doms)
}
