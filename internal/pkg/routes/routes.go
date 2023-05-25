package routes

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/jsonapi"
	"github.com/naivary/instance/internal/pkg/service"
)

const (
	reqTimeout = 20 * time.Second
)

func newRoot() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.SetHeader("Content-Type", jsonapi.MediaType))
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Timeout(reqTimeout))
	return r
}

// NewTestRouter returns a chi.Router
// with all the root middleware attached.
// Only use this for testing.
func NewTestRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.SetHeader("Content-Type", jsonapi.MediaType))
	r.Use(middleware.RequestID)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Timeout(reqTimeout))
	return r
}

func New(svcs []service.Service[chi.Router]) chi.Router {
	root := newRoot()
	for _, svc := range svcs {
		svc.Register(root)
	}
	return root
}
