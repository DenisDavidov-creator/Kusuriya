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

// CreateCase godoc
// @Summary      Open a case
// @Description  Creates a case in status "open".
// @Tags         case
// @Accept       json
// @Produce      json
// @Param        request  body      CreateCaseRequest  true  "case to open"
// @Success      200      {object}  CaseResponse
// @Failure      400      {object}  httpx.ErrorResponse  "malformed JSON or validation failed"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /case [post]
func (h *CaseHandler) CreateCase(w http.ResponseWriter, r *http.Request) {
	var req CreateCaseRequest

	err := httpx.Decode(r.Body, &req)
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
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

// GetCaseByID godoc
// @Summary      Get a case with its diagnosis
// @Description  Returns the case together with the five most likely causes, ranked from
// @Description  the observations recorded on it. A case with no observations is ranked by priors alone.
// @Tags         case
// @Produce      json
// @Param        id   path      int  true  "case id"
// @Success      200  {object}  CaseResponseWithDiagnoses
// @Failure      400  {object}  httpx.ErrorResponse  "id is not a number"
// @Failure      404  {object}  httpx.ErrorResponse  "no such case"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /case/{id} [get]
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

// CloseCaseByID godoc
// @Summary      Close a case
// @Description  Sets the case status to "closed". Closing an already closed case is a no-op.
// @Tags         case
// @Produce      json
// @Param        id   path      int  true  "case id"
// @Success      200  {object}  CaseResponse
// @Failure      400  {object}  httpx.ErrorResponse  "id is not a number"
// @Failure      404  {object}  httpx.ErrorResponse  "no such case"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /case/{id}/close [patch]
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

// CreateObservation godoc
// @Summary      Record an observation on a case
// @Description  Attaches an observed symptom or factor to the case. ref_id points at a
// @Description  symptom id when kind is "symptom", at a factor id when kind is "factor".
// @Tags         case
// @Accept       json
// @Produce      json
// @Param        id       path      int                       true  "case id"
// @Param        request  body      CreateObservationRequest  true  "observation to record"
// @Success      200      {object}  ObservationsResponse
// @Failure      400      {object}  httpx.ErrorResponse  "id is not a number, malformed JSON, or validation failed"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /case/{id}/observation [post]
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

// GetObservationByID godoc
// @Summary      List observations of a case
// @Description  Returns every observation recorded on the case, ordered by id.
// @Description  An unknown case id yields an empty list rather than 404.
// @Tags         case
// @Produce      json
// @Param        id   path      int  true  "case id"
// @Success      200  {array}   ObservationsResponse
// @Failure      400  {object}  httpx.ErrorResponse  "id is not a number"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /case/{id}/observation [get]
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
