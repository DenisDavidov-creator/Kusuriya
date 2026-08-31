package handler

import (
	"Kusuriya/internal/domain"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

func NewErrorResponse(code int, message string, details any) ErrorResponse {
	return ErrorResponse{
		ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

func MapErr(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrDuplicateCode):
		return http.StatusConflict, "duplicate unique value"
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, "value not found"
	default:
		return http.StatusInternalServerError, "unexpected error"
	}
}
