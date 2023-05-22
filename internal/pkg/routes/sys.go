package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/ctrl"
)

func sys(views *ctrl.Views) chi.Router {
	r := chi.NewRouter()
	r.Get("/health", views.Sys.Health)
	r.Get("/metrics", views.Sys.Metrics)
	return r
}
