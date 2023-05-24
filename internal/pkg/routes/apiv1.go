package routes

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/jsonapi"
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

func new(svcs *services.Services) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.SetHeader("Content-Type", jsonapi.MediaType))
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Timeout(reqTimeout))

	r.Mount("/api/v1", apiv1(svcs))
	r.Mount("/sys", sys(svcs))
	r.Mount("/fs", fs(svcs))
	return r
}
