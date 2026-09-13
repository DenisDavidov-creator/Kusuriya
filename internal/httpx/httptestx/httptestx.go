// Package httptestx holds helpers for unit-testing chi handlers without a router.
package httptestx

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func WithURLParam(req *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func DoRequest(h http.HandlerFunc, method, id, body string) *httptest.ResponseRecorder {
	var r io.Reader
	if body != "" {
		r = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, "/", r)
	if id != "" {
		req = WithURLParam(req, "id", id)
	}
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

func DecodeBody[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var got T
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	return got
}
