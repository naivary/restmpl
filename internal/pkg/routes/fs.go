package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/ctrl"
)

func fs(views *ctrl.Views) chi.Router {
	r := chi.NewRouter()
	r.Post("/upload", views.Fs.Upload)
	r.Delete("/delete", views.Fs.Upload)
	return r
}
