package handler

import (
	"Kusuriya/internal/httpx"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CaseHandler struct {
	service CaseServiceInterface
}

func NewCaseHandler(service CaseServiceInterface) *CaseHandler {
	return &CaseHandler{
		service: service,
	}
}

func (h *CaseHandler) CreateCase(w http.ResponseWriter, r *http.Request) {
	var req CreateCaseRequest

	httpx.Decode(r.Body, &req)

	problmes, err := req.Validate()
	if err != nil {
		httpx.RespondError(w, r, err, problmes)
		return
	}

	result, err := h.service.CreateCase(r.Context(), req.ToDomain())
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainToCase(result))
}

func (h *CaseHandler) GetCaseByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondError(w, r, err, "wrong URLParam")
		return
	}

	causes, caseDetail, err := h.service.GetCaseByID(r.Context(), id)
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, DomainToCaseDetail(causes, caseDetail))
}

func (h *CaseHandler) CloseCaseByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondError(w, r, err, "wrong URLParam")
		return
	}

	result, err := h.service.CloseCaseByID(r.Context(), id)
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, DomainToCase(result))
}

func (h *CaseHandler) CreateObservation(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondError(w, r, err, "wrong URLParam")
		return
	}
	var req CreateObservationRequest
	err = httpx.Decode(r.Body, &req)
	if err != nil {
		httpx.RespondError(w, r, err, "wrong JSON format")
		return
	}
	problems, err := req.Validate()
	if err != nil {
		httpx.RespondError(w, r, err, problems)
		return
	}
	result, err := h.service.CreateObservation(r.Context(), req.ToDomain(id))
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, DomainToObservation(result))
}

func (h *CaseHandler) GetObservationByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondError(w, r, err, "wrong URLParam")
		return
	}

	result, err := h.service.GetObservationByID(r.Context(), id)
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, DomainsToObservations(result))
}
