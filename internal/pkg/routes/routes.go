package routes

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/naivary/instance/internal/pkg/services"
)

const (
	reqTimeout = 20 * time.Second
)

func New(services *services.Services) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Timeout(reqTimeout))
	r.Use(middleware.SetHeader("content-type", "application/json"))

	r.Mount("/api/v1", apiv1(services))
	r.Mount("/sys", sys(services))
	r.Mount("/fs", fs(services))
	return r
}
