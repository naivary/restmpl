package env

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/service"
)

type Env[T any] interface {
	UUID() string
	Version() string

	Services() map[string]service.Service[chi.Router]
	Router() http.Handler

	// Config return the config Manager of
	// your application.
	Config() T
}
