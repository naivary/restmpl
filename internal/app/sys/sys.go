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

type Env struct {
	K  *koanf.Koanf
	DB *sql.DB

	M metadata.Metadata
}

func (e *Env) Health(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	err := e.DB.Ping()
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}
	e.M.DBRunning = err == nil

	err = jsonapi.MarshalPayload(w, &e.M)
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}

	w.WriteHeader(http.StatusOK)
}
