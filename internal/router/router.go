package router

import (
	"net/http"

	"github.com/antonhancharyk/crypto-knight-gateway/internal/config"
	"github.com/antonhancharyk/crypto-knight-gateway/internal/health"
	"github.com/antonhancharyk/crypto-knight-gateway/internal/lb"
	"github.com/antonhancharyk/crypto-knight-gateway/internal/middleware"
	"github.com/antonhancharyk/crypto-knight-gateway/internal/proxy"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func New(cfg *config.Config, logger *zap.Logger) (http.Handler, error) {
	r := chi.NewRouter()

	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.Timeout(cfg.Timeout))

	r.Get("/healthz", health.Handler)

	apiPool := lb.NewRoundRobin([]string{
		"http://backend:8081",
	})
	r.Mount("/api",
		http.StripPrefix("/api",
			proxy.NewReverseProxy(apiPool),
		),
	)

	frontendPool := lb.NewRoundRobin([]string{
		"http://frontend:80",
	})
	r.Mount("/",
		proxy.NewReverseProxy(frontendPool),
	)

	return r, nil
}
