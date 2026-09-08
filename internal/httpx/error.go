package httpx

import (
	"Kusuriya/internal/domain"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
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

var ErrValidation = errors.New("validation not passed")

func MapErr(err error) (int, slog.Level, string) {
	switch {
	case errors.Is(err, domain.ErrDuplicateCode):
		return http.StatusConflict, slog.LevelInfo, "duplicate unique value"
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, slog.LevelInfo, "value not found"
	case errors.Is(err, ErrValidation):
		return http.StatusBadRequest, slog.LevelInfo, "validation not passed"
	case errors.Is(err, io.EOF),
		errors.Is(err, io.ErrUnexpectedEOF),
		errors.Is(err, strconv.ErrSyntax):
		return http.StatusBadRequest, slog.LevelInfo, "bad request"
	default:
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &syntaxErr) ||
			errors.As(err, &typeErr) {
			return http.StatusBadRequest, slog.LevelInfo, "wrong JSON format"
		}
		if strings.HasPrefix(err.Error(), "json: unknown field ") {
			return http.StatusBadRequest, slog.LevelInfo, err.Error()
		}

		return http.StatusInternalServerError, slog.LevelError, "unexpected error"
	}
}
