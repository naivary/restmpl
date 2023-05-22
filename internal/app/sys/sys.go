package sys

import (
	"database/sql"
	"encoding/json"
	"net/http"

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

	err = json.NewEncoder(w).Encode(&e.M)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
