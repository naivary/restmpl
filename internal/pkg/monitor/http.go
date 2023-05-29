package monitor

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/jsonapi"
	"github.com/naivary/instance/internal/pkg/japi"
	"golang.org/x/exp/slog"
)

func (a agent) HTTP() chi.Router {
	r := chi.NewRouter()
	r.Get("/health", a.health)
	r.Get("/metrics", a.metrics)
	r.Get("/services", a.services)
	return r
}

func (a agent) health(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	for _, svc := range a.svcs {
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

func (a agent) metrics(w http.ResponseWriter, r *http.Request) {}

func (a agent) services(w http.ResponseWriter, r *http.Request) {}
