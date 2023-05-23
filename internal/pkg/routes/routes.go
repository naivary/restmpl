package routes

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/jsonapi"
	"github.com/naivary/instance/internal/pkg/services"
)

const (
	reqTimeout = 20 * time.Second
)

func New(svcs *services.Services) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.SetHeader("Content-Type", jsonapi.MediaType))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Timeout(reqTimeout))

	r.Mount("/api/v1", apiv1(svcs))
	r.Mount("/sys", sys(svcs))
	r.Mount("/fs", fs(svcs))
	return r
}
