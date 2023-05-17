package sys

import (
	"net/http"

	"github.com/knadh/koanf/v2"
)

type Env struct {
	K *koanf.Koanf
}

func (e *Env) Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Its healthy!"))
	w.WriteHeader(http.StatusOK)
}
