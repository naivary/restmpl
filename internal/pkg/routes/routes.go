package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/ctrl"
)

func New(views *ctrl.Views) chi.Router {
	r := chi.NewRouter()
	r.Mount("/sys", sys(views))

	return r
}

func sys(views *ctrl.Views) chi.Router {
	r := chi.NewRouter()
	r.Get("/health", views.Sys.Health)

	return r
}
