package sys

import (
	"encoding/json"
	"net/http"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/models"
)

type Env struct {
	K *koanf.Koanf
}

func (e *Env) Health(w http.ResponseWriter, r *http.Request) {
	// TODO(naivary) should be static and calculated to compile time
	m := models.Metadata{
		Version: e.K.String("version"),
	}

	err := json.NewEncoder(w).Encode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
