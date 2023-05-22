package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/naivary/instance/internal/pkg/services"
)

func fs(svcs *services.Services) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.SetHeader("content-type", "image/png"))

	r.Post("/upload", svcs.Fs.Upload)
	r.Delete("/delete", svcs.Fs.Delete)
	r.Method(http.MethodGet, "/", svcs.Fs.Fs.Store)
	return r
}
