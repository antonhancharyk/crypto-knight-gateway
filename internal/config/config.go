package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultAddress     = ":8080"
	defaultTimeout     = 10 * time.Second
	defaultUpstreamURL = "http://frontend:80"
)

type ServerConfig struct {
	Address string
}

type Config struct {
	Server       ServerConfig
	Timeout      time.Duration
	UpstreamURLs []string
}

// Load reads configuration from the environment with fallbacks to defaults.
// Env vars: PORT (e.g. "8080" -> ":8080"), UPSTREAM_URLS (comma-separated), REQUEST_TIMEOUT (e.g. "10s").
func Load() (*Config, error) {
	addr := defaultAddress
	if port := os.Getenv("PORT"); port != "" {
		if _, err := strconv.Atoi(port); err == nil {
			addr = ":" + port
		}
	}
	if a := os.Getenv("SERVER_ADDRESS"); a != "" {
		addr = a
	}

	timeout := defaultTimeout
	if s := os.Getenv("REQUEST_TIMEOUT"); s != "" {
		if d, err := time.ParseDuration(s); err == nil && d > 0 {
			timeout = d
		}
	}

	upstreams := []string{defaultUpstreamURL}
	if s := os.Getenv("UPSTREAM_URLS"); s != "" {
		parts := strings.Split(s, ",")
		var trimmed []string
		for _, p := range parts {
			if u := strings.TrimSpace(p); u != "" {
				trimmed = append(trimmed, u)
			}
		}
		if len(trimmed) > 0 {
			upstreams = trimmed
		}
	}

	return &Config{
		Server:       ServerConfig{Address: addr},
		Timeout:      timeout,
		UpstreamURLs: upstreams,
	}, nil
}
