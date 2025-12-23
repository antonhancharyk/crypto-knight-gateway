package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/antonhancharyk/crypto-knight-gateway/internal/lb"
)

func NewReverseProxy(pool *lb.RoundRobin) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		up := pool.Next()
		if up == "" {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("no upstream"))
			return
		}

		u, err := url.Parse(up)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("bad upstream"))
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(u)
		proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("upstream error: " + err.Error()))
		}

		proxy.ServeHTTP(w, r)
	})
}
