package httpx

import (
	"Kusuriya/internal/domain"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// jsonErr runs a real Decode so the test sees the exact error types
// encoding/json produces, instead of hand-built structs that may drift.
func jsonErr(body string) error {
	var v struct {
		Name  string  `json:"name"`
		Prior float64 `json:"prior"`
	}
	return Decode(io.NopCloser(strings.NewReader(body)), &v)
}

func TestMapErr(t *testing.T) {
	_, atoiErr := strconv.Atoi("meow")

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantLevel  slog.Level
		wantMsg    string
	}{
		// domain errors
		{"duplicate code", domain.ErrDuplicateCode, http.StatusConflict, slog.LevelInfo, "duplicate unique value"},
		{"not found", domain.ErrNotFound, http.StatusNotFound, slog.LevelInfo, "value not found"},
		{"validation", ErrValidation, http.StatusBadRequest, slog.LevelInfo, "validation not passed"},
		{"check constraint is a client error", domain.ErrCheckPass, http.StatusBadRequest, slog.LevelInfo, "check constraint not passed"},

		// wrapped errors must still be recognised (errors.Is, not ==)
		{"wrapped not found", fmt.Errorf("get cause: %w", domain.ErrNotFound), http.StatusNotFound, slog.LevelInfo, "value not found"},
		{"wrapped validation", fmt.Errorf("cause: %w", ErrValidation), http.StatusBadRequest, slog.LevelInfo, "validation not passed"},

		// request parsing
		{"empty body", jsonErr(``), http.StatusBadRequest, slog.LevelInfo, "bad request"},
		{"truncated json", jsonErr(`{"name":"x"`), http.StatusBadRequest, slog.LevelInfo, "bad request"},
		{"malformed json", jsonErr(`{"name": x}`), http.StatusBadRequest, slog.LevelInfo, "wrong JSON format"},
		{"wrong field type", jsonErr(`{"prior":"high"}`), http.StatusBadRequest, slog.LevelInfo, "wrong JSON format"},
		{"unknown field echoes the field name", jsonErr(`{"nmae":"x"}`), http.StatusBadRequest, slog.LevelInfo, `json: unknown field "nmae"`},
		{"bad url param", atoiErr, http.StatusBadRequest, slog.LevelInfo, "bad request"},

		// anything else is a 500 and logged as an error
		{"unknown error", errors.New("db on fire"), http.StatusInternalServerError, slog.LevelError, "unexpected error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, level, msg := MapErr(tt.err)

			assert.Equal(t, tt.wantStatus, status, "status")
			assert.Equal(t, tt.wantLevel, level, "log level")
			assert.Equal(t, tt.wantMsg, msg, "message")
		})
	}
}
