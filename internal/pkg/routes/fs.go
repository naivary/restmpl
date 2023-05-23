package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/services"
)

func fs(svcs *services.Services) chi.Router {
	r := chi.NewRouter()
	for _, mw := range svcs.Fs.Middlewares() {
		r.Use(mw)
	}
	r.Post("/create", svcs.Fs.Create)
	r.Delete("/remove", svcs.Fs.Remove)
	r.Get("/read", svcs.Fs.Read)
	return r
}
