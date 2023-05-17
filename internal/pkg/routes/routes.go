package routes

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/naivary/instance/internal/pkg/ctrl"
)

const (
	reqTimeout = 20 * time.Second
)

func New(views *ctrl.Views) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Timeout(reqTimeout))
	r.Use(middleware.SetHeader("content-type", "application/json"))

	r.Route("/v1", func(r chi.Router) {
		// All REST services for v1
	})
	r.Mount("/sys", sys(views))
	return r
}

func sys(views *ctrl.Views) chi.Router {
	r := chi.NewRouter()
	r.Get("/health", views.Sys.Health)
	return r
}
