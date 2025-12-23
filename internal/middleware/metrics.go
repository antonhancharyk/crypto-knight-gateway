package middleware

import (
	"net/http"

	prom "github.com/prometheus/client_golang/prometheus/promhttp"
)

func Metrics() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

func PrometheusHandler() http.Handler {
	return prom.Handler()
}
