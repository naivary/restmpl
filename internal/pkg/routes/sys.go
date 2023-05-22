package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/services"
)

func sys(services *services.Services) chi.Router {
	r := chi.NewRouter()
	r.Get("/health", services.Sys.Health)
	return r
}
