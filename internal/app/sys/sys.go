package sys

import (
	"database/sql"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/models/metadata"
)

type Env struct {
	K  *koanf.Koanf
	DB *sql.DB

	M metadata.Metadata
}

func (e *Env) Health(w http.ResponseWriter, r *http.Request) {
	err := e.DB.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	e.M.DBRunning = err == nil

	w.WriteHeader(http.StatusOK)
	err = jsonapi.MarshalPayload(w, &e.M)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
