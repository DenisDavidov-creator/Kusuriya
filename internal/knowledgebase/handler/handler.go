package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
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

func (h *KnowledgebaseHandler) CreateCause(w http.ResponseWriter, r *http.Request) {
	var req CreateCauseRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "Wrong json format", ""))
		return
	}
	problems := req.Validate()
	if len(problems) != 0 {
		writeJSON(w, http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "wrong fields", problems))
		return
	}

	dom, err := h.service.CreateCause(r.Context(), req.ToDomain())
	if err != nil {
		status, detail := MapErr(err)
		slog.ErrorContext(r.Context(), "cause couldn't create", "error", err, "code", req.Code)
		writeJSON(w, status, NewErrorResponse(status, "cause couldn't create", detail))
		return
	}
	writeJSON(w, http.StatusOK, dom)
}
