package sys

import "net/http"

type Env struct {
}

func (e *Env) Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Its healthy!"))
	w.WriteHeader(http.StatusOK)
}
