package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/services"
)

func sys(svcs *services.Services) chi.Router {
	r := chi.NewRouter()
	r.Get("/health", svcs.Sys.Health)
	return r
}
