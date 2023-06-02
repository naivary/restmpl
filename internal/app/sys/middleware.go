package sys

import "net/http"

func (s Sys) Middlewares() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		s.requestCounter,
	}
}

func (s Sys) requestCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter := s.metManager.GetCounter("requestCounter")
		counter.Inc()
		next.ServeHTTP(w, r)
	})
}
