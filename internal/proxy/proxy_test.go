package proxy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type mockPool struct {
	next string
}

func (m *mockPool) Next() string { return m.next }

func TestNewReverseProxy_NoUpstream(t *testing.T) {
	pool := &mockPool{next: ""}
	handler := NewReverseProxy(pool)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadGateway)
	}
	if body := rec.Body.String(); body != "no upstream" {
		t.Errorf("body = %q, want %q", body, "no upstream")
	}
}

func TestNewReverseProxy_BadUpstreamURL(t *testing.T) {
	pool := &mockPool{next: "://invalid"}
	handler := NewReverseProxy(pool)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if body := rec.Body.String(); body != "bad upstream" {
		t.Errorf("body = %q, want %q", body, "bad upstream")
	}
}

func TestNewReverseProxy_ProxiesToUpstream(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("backend response"))
	}))
	defer backend.Close()

	u, _ := url.Parse(backend.URL)
	pool := &mockPool{next: u.String()}
	handler := NewReverseProxy(pool)

	req := httptest.NewRequest(http.MethodGet, "/path", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if body := rec.Body.String(); body != "backend response" {
		t.Errorf("body = %q, want %q", body, "backend response")
	}
}

func TestNewReverseProxy_ErrorHandlerDoesNotLeakUpstreamError(t *testing.T) {
	// Backend that closes connection immediately so proxy sees an error
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Skip("hijacker not supported")
		}
		conn, _, _ := hj.Hijack()
		conn.Close()
	}))
	defer backend.Close()

	u, _ := url.Parse(backend.URL)
	pool := &mockPool{next: u.String()}
	handler := NewReverseProxy(pool)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadGateway)
	}
	body := rec.Body.String()
	if body != "upstream error" {
		t.Errorf("body = %q, want %q (must not leak upstream error details)", body, "upstream error")
	}
}
