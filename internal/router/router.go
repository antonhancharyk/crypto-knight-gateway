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
	r.Use(middleware.Metrics())
	// auth can be added per-route

	r.Get("/healthz", health.Handler)

	for _, b := range cfg.Backends {
		pool := lb.NewRoundRobin([]string{b.URL})
		r.Mount("/"+b.Name, http.StripPrefix("/"+b.Name, proxy.NewReverseProxy(pool)))
	}

	r.Route("/api", func(r chi.Router) {
		r.Method("GET", "/v1/users/*", proxy.NewReverseProxy(lb.NewRoundRobin([]string{"http://localhost:8081"})))
	})

	r.Handle("/metrics", middleware.PrometheusHandler())

	return r, nil
}
