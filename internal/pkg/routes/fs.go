package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/services"
)

func fs(services *services.Services) chi.Router {
	r := chi.NewRouter()
	r.Post("/upload", services.Fs.Upload)
	r.Delete("/delete", services.Fs.Delete)
	r.Handle("/", services.Fs.Fs.Handler)
	return r
}
