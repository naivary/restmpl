package sys

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/models"
)

type Env struct {
	K  *koanf.Koanf
	DB *sql.DB

	m models.Metadata
}

func (e *Env) Health(w http.ResponseWriter, r *http.Request) {
	// TODO(naivary) should be static and calculated to compile time
	e.m = models.Metadata{
		Version: e.K.String("version"),
	}

	err := e.DB.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(&e.m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (e *Env) Metrics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("the promotheus metrics!\n"))
}
