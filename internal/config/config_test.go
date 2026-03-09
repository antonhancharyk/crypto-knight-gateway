package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("REQUEST_TIMEOUT")
	os.Unsetenv("UPSTREAM_URLS")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if cfg.Server.Address != ":8080" {
		t.Errorf("Server.Address = %q, want :8080", cfg.Server.Address)
	}
	if cfg.Timeout != 10*time.Second {
		t.Errorf("Timeout = %v, want 10s", cfg.Timeout)
	}
	if len(cfg.UpstreamURLs) != 1 || cfg.UpstreamURLs[0] != "http://frontend:80" {
		t.Errorf("UpstreamURLs = %v, want [http://frontend:80]", cfg.UpstreamURLs)
	}
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("REQUEST_TIMEOUT", "5s")
	os.Setenv("UPSTREAM_URLS", "http://a:80, http://b:80")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("REQUEST_TIMEOUT")
		os.Unsetenv("UPSTREAM_URLS")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if cfg.Server.Address != ":9090" {
		t.Errorf("Server.Address = %q, want :9090", cfg.Server.Address)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("Timeout = %v, want 5s", cfg.Timeout)
	}
	if len(cfg.UpstreamURLs) != 2 {
		t.Errorf("len(UpstreamURLs) = %d, want 2", len(cfg.UpstreamURLs))
	}
	if cfg.UpstreamURLs[0] != "http://a:80" || cfg.UpstreamURLs[1] != "http://b:80" {
		t.Errorf("UpstreamURLs = %v", cfg.UpstreamURLs)
	}
}

func TestLoad_SERVER_ADDRESS_Overrides_PORT(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("SERVER_ADDRESS", "0.0.0.0:3000")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("SERVER_ADDRESS")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if cfg.Server.Address != "0.0.0.0:3000" {
		t.Errorf("Server.Address = %q, want 0.0.0.0:3000", cfg.Server.Address)
	}
}
