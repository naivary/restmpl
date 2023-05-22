package routes

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/ctrl"
)

func apiv1(views *ctrl.Views) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))
	r.Mount("/example", example(views))
	return r
}

func example(views *ctrl.Views) chi.Router {
	r := chi.NewRouter()
	return r
}
