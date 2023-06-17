package env

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/restmpl/internal/pkg/service"
)

type Env interface {
	// Unique ID of the environment
	ID() string

	// version of the environment
	Version() string

	// http router to serve public request to the services
	HTTP() chi.Router

	// Configuration of the environment represented
	// by the koanf config manager.
	Config() *koanf.Koanf

	// Serve the env and its services for public traffic
	Serve() error

	Join(svcs ...service.Service) error

	// Graceful shutdown of the env
	// and its associated services
	Shutdown() error

	// Health of the env
	Health() error

	// Context of the environment
	Context() context.Context

	// Initiliaze the environment and all its dependencies
	Init() error
}
