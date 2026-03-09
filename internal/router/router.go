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

	pool := lb.NewRoundRobin(cfg.UpstreamURLs)
	r.Mount("/", proxy.NewReverseProxy(pool))

	return r, nil
}
