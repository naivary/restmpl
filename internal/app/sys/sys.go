package sys

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/google/jsonapi"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/japi"
	"github.com/naivary/instance/internal/pkg/models/metadata"
)

type Sys struct {
	K  *koanf.Koanf
	DB *sql.DB

	M metadata.Metadata
}

func (s *Sys) Health(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	err := s.DB.Ping()
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}
	s.M.DBRunning = err == nil

	w.Header().Add("Content-Type", jsonapi.MediaType)
	err = jsonapi.MarshalPayload(w, &s.M)
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}

	w.WriteHeader(http.StatusOK)
}
