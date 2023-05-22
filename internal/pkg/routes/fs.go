package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/naivary/instance/internal/pkg/services"
)

func fs(svcs *services.Services) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.SetHeader("content-type", "image/png"))

	r.Post("/create", svcs.Fs.Create)
	r.Delete("/remove", svcs.Fs.Remove)
	r.Get("/read", svcs.Fs.Read)
	return r
}
