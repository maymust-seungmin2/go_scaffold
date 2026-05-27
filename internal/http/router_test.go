package http

import (
	"io"
	"log/slog"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRouter_Healthz(t *testing.T) {
	injector := do.New()
	do.ProvideValue(injector, slog.New(slog.NewTextHandler(io.Discard, nil)))

	router := NewRouter(injector)

	req := httptest.NewRequest(nethttp.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, nethttp.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())
}

func TestNewRouter_NotFound(t *testing.T) {
	injector := do.New()
	do.ProvideValue(injector, slog.New(slog.NewTextHandler(io.Discard, nil)))

	router := NewRouter(injector)

	req := httptest.NewRequest(nethttp.MethodGet, "/does-not-exist", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, nethttp.StatusNotFound, rec.Code)
}
