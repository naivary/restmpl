package routes

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/services"
)

func apiv1(services *services.Services) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))
	r.Mount("/example", example(services))
	return r
}

func example(services *services.Services) chi.Router {
	r := chi.NewRouter()
	return r
}
