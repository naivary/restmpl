package sys

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/jsonapi"
	"github.com/naivary/instance/internal/pkg/japi"
	"github.com/naivary/instance/internal/pkg/service"
)

func (s Sys) health(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	infos := []*service.Info{}
	for _, svc := range s.Svcs {
		info, err := svc.Health()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			jerr := japi.NewError(err, http.StatusServiceUnavailable, reqID)
			jsonapi.MarshalErrors(w, japi.Errors(&jerr))
			return
		}
		infos = append(infos, info)
	}
	s.Meta.Svcs = infos

	if err := jsonapi.MarshalPayload(w, &s.Meta); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		jerr := japi.NewError(err, http.StatusServiceUnavailable, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}
}

func (s Sys) metrics(w http.ResponseWriter, r *http.Request) {}
