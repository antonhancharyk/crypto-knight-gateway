package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/antonhancharyk/crypto-knight-gateway/internal/config"
	"go.uber.org/zap"
)

func TestNew_Healthz(t *testing.T) {
	cfg := &config.Config{
		Server:       config.ServerConfig{Address: ":8080"},
		Timeout:      10 * time.Second,
		UpstreamURLs: []string{"http://localhost:9999"},
	}
	logger := zap.NewNop()

	r, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("New() err = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GET /healthz status = %d, want %d", rec.Code, http.StatusOK)
	}
	if body := rec.Body.String(); body != "ok" {
		t.Errorf("GET /healthz body = %q, want %q", body, "ok")
	}
}
