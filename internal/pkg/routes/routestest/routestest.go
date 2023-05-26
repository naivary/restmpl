package routestest

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/jsonapi"
)

const (
	reqTimeout = 20 * time.Second
)

// NewTestRouter returns a chi.Router
// with all the root middleware attached.
// Only use this for testing.
func New() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.SetHeader("Content-Type", jsonapi.MediaType))
	r.Use(middleware.RequestID)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Timeout(reqTimeout))
	return r
}
