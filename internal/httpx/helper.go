package httpx

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("write json response", "error", err)
	}

}

func RespondError(w http.ResponseWriter, r *http.Request, err error, details any) {
	status, level, message := MapErr(err)
	slog.Log(r.Context(), level, "request failed", "error", err)
	WriteJSON(w, status, NewErrorResponse(status, message, details))
}

func Decode(r io.ReadCloser, req any) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	return dec.Decode(req)
}
