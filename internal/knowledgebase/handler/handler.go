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
