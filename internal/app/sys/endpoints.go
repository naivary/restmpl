package sys

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/jsonapi"
	"github.com/naivary/instance/internal/pkg/japi"
)

func (s *Sys) Health(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)
	err := s.DB.DB().PingContext(ctx)
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}
	err = jsonapi.MarshalPayload(w, &s.M)
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}
	w.WriteHeader(http.StatusOK)
}
