package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Pool selects the next upstream address. Used for testing with mocks.
type Pool interface {
	Next() string
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(message))
}

// NewReverseProxy returns an http.Handler that reverse-proxies to the next upstream from pool.
func NewReverseProxy(pool Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		up := pool.Next()
		if up == "" {
			writeError(w, http.StatusBadGateway, "no upstream")
			return
		}

		u, err := url.Parse(up)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "bad upstream")
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(u)
		proxy.ErrorHandler = func(w http.ResponseWriter, _ *http.Request, _ error) {
			writeError(w, http.StatusBadGateway, "upstream error")
		}

		proxy.ServeHTTP(w, r)
	})
}
