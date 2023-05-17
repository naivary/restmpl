package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/ctrl"
)

func New(views *ctrl.Views) chi.Router {
	r := chi.NewRouter()
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
