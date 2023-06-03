package env

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/google/jsonapi"
	"github.com/naivary/apitmpl/internal/pkg/japi"
	"github.com/naivary/apitmpl/internal/pkg/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (a *API) health(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	infos := []*service.Info{}
	// check env health
	if err := a.Health(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		jerr := japi.NewError(err, http.StatusServiceUnavailable, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}

	// check services health
	for _, svc := range a.svcs {
		info, err := svc.Health()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			jerr := japi.NewError(err, http.StatusServiceUnavailable, reqID)
			jsonapi.MarshalErrors(w, japi.Errors(&jerr))
			return
		}
		if info != nil {
			infos = append(infos, info)
		}
	}
	a.meta.Svcs = infos

	if err := jsonapi.MarshalPayload(w, &a.meta); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		jerr := japi.NewError(err, http.StatusServiceUnavailable, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}
}

func (a API) metrics() http.Handler {
	opts := promhttp.HandlerOpts{
		Registry: a.metric.Registry(),
	}
	return promhttp.HandlerFor(a.metric.Registry(), opts)
}
