package config

import (
	"time"
)

type ServerConfig struct {
	Address string
}

type Config struct {
	Server  ServerConfig
	Timeout time.Duration
}

func Load() (*Config, error) {
	return &Config{
		Server:  ServerConfig{Address: ":8080"},
		Timeout: 10 * time.Second,
	}, nil
}
