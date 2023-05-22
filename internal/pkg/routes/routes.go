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
