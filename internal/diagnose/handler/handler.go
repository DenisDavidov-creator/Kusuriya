package handler

import (
	"Kusuriya/internal/diagnose/dto"
	"Kusuriya/internal/httpx"
	"net/http"
)

type DiagnoseHandler struct {
	service DiagnoseServiceInterface
}

func NewDiagnoseHandler(service DiagnoseServiceInterface) *DiagnoseHandler {
	return &DiagnoseHandler{service: service}
}

// DetermineDiagnose godoc
// @Summary      Rank causes for the given evidence
// @Description  Scores every cause against the supplied symptom and factor codes and
// @Description  returns the five most likely, each with the evidence it was judged on.
// @Description  At least one symptom or factor code is required. Unknown codes are ignored.
// @Tags         diagnose
// @Accept       json
// @Produce      json
// @Param        request  body      dto.DetermineDiagnoseRequest  true  "observed symptom and factor codes"
// @Success      200      {array}   dto.DetermineDiagnoseResponse
// @Failure      400      {object}  httpx.ErrorResponse  "malformed JSON or no codes given"
// @Failure      500      {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /diagnose [post]
func (h *DiagnoseHandler) DetermineDiagnose(w http.ResponseWriter, r *http.Request) {
	var req dto.DetermineDiagnoseRequest
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

	result, err := h.service.DetermineDiagnose(r.Context(), req.ToDomain())
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto.DomainsToDetermineDiagnose(result))
}
