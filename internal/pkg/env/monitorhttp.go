package env

import (
	"encoding/json"
	"net/http"

	"github.com/naivary/apitmpl/internal/pkg/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (a *API) health(w http.ResponseWriter, r *http.Request) {
	infos := []*service.Info{}
	// check env health
	if err := a.Health(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// check services health
	for _, svc := range a.svcs {
		info, err := svc.Health()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		if info != nil {
			infos = append(infos, info)
		}
	}
	a.meta.Svcs = infos

	if err := json.NewEncoder(w).Encode(&a.meta); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
}

func (a API) metrics() http.Handler {
	opts := promhttp.HandlerOpts{
		Registry: a.metric.Registry(),
	}
	return promhttp.HandlerFor(a.metric.Registry(), opts)
}
