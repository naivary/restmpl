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
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Timeout(reqTimeout))
	r.Use(middleware.SetHeader("content-type", "application/json"))

	r.Mount("/api/v1", apiv1(views))
	r.Mount("/sys", sys(views))
	r.Mount("/fs", fs(views))
	return r
}

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

func sys(views *ctrl.Views) chi.Router {
	r := chi.NewRouter()
	r.Get("/health", views.Sys.Health)
	r.Get("/metrics", views.Sys.Metrics)
	return r
}

func fs(views *ctrl.Views) chi.Router {
	r := chi.NewRouter()
	r.Post("/upload", views.Fs.Upload)
	r.Delete("/delete", views.Fs.Upload)
	return r
}
