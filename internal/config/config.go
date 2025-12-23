package config

import (
	"time"
)

type ServerConfig struct {
	Address string
}

type Backend struct {
	Name string
	URL  string
}

type Config struct {
	Server   ServerConfig
	Backends []Backend
	Timeout  time.Duration
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{Address: ":8080"},
		Backends: []Backend{
			{Name: "users", URL: "http://localhost:8081"},
			{Name: "orders", URL: "http://localhost:8082"},
		},
		Timeout: 10 * time.Second,
	}, nil
}
