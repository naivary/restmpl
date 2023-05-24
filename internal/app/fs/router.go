package fs

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (e Env) Router() http.Handler {
	r := chi.NewRouter()
	for _, mw := range e.Middlewares() {
		r.Use(mw)
	}
	r.Post("/create", e.Create)
	r.Delete("/remove", e.Remove)
	r.Get("/read", e.Read)
	return r
}
