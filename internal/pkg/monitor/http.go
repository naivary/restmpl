package monitor

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/jsonapi"
	"github.com/naivary/instance/internal/pkg/japi"
	"golang.org/x/exp/slog"
)

func (m manager) HTTP() chi.Router {
	r := chi.NewRouter()
	r.Get("/health", m.health)
	r.Get("/metrics", m.metrics)
	r.Get("/services", m.services)
	return r
}

func (m manager) health(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	for _, svc := range m.svcs {
		info, err := svc.Health()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			jerr := japi.NewError(err, http.StatusServiceUnavailable, reqID)
			jsonapi.MarshalErrors(w, japi.Errors(&jerr))
			return
		}
		slog.InfoCtx(r.Context(), "service info", "info", info)
	}
}

func (m manager) metrics(w http.ResponseWriter, r *http.Request) {}

func (m manager) services(w http.ResponseWriter, r *http.Request) {}
