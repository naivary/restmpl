package sys

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/jsonapi"
	"github.com/naivary/instance/internal/pkg/japi"
)

func (s *Sys) Health(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	err := s.DB.Ping()
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}
	s.M.DBRunning = err == nil
	s.M.ReqID = reqID
	err = jsonapi.MarshalPayload(w, &s.M)
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}

	w.WriteHeader(http.StatusOK)
}
