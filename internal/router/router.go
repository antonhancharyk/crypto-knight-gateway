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
		"http://10.10.0.3:8080",
	})
	r.Mount("/api",
		http.StripPrefix("/api",
			proxy.NewReverseProxy(apiPool),
		),
	)

	rabbitmqPool := lb.NewRoundRobin([]string{
		"http://10.10.0.3:15672",
	})
	r.Mount("/rabbitmq",
		http.StripPrefix("/rabbitmq",
			proxy.NewReverseProxy(rabbitmqPool),
		),
	)

	kibanaPool := lb.NewRoundRobin([]string{
		"http://10.10.0.3:5601",
	})
	r.Mount("/kibana",
		http.StripPrefix("/kibana",
			proxy.NewReverseProxy(kibanaPool),
		),
	)

	grafanaPool := lb.NewRoundRobin([]string{
		"http://10.10.0.3:3000",
	})
	r.Mount("/grafana",
		http.StripPrefix("/grafana",
			proxy.NewReverseProxy(grafanaPool),
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
