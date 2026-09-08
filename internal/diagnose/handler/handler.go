package handler

import (
	"Kusuriya/internal/httpx"
	"net/http"
)

type DiagnoseHandler struct {
	service DiagnoseServiceInterface
}

func NewDiagnoseHandler(service DiagnoseServiceInterface) *DiagnoseHandler {
	return &DiagnoseHandler{service: service}
}

func (h *DiagnoseHandler) DeterminDiagnose(w http.ResponseWriter, r *http.Request) {
	var req DetermineDiagnoseRequest
	err := httpx.Decode(r.Body, &req)
	if err != nil {
		httpx.RespondError(w, r, err, err.Error())
		return
	}
	err = req.Validate()
	if err != nil {
		httpx.RespondError(w, r, err, err.Error())
		return
	}

	result, err := h.service.DeterminDiagnose(r.Context(), req.ToDomain())
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainsToDetermineDiagnose(result))
}
