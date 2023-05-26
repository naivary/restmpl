package logging

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// Logger is the default logger for any http request
func (l Logging) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		l.Info.InfoCtx(r.Context(), "request", "method", r.Method, "reqID", reqID, "endpoint", r.URL, "host", r.Host)
		next.ServeHTTP(w, r)
	})
}
