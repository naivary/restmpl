package env

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/restmpl/internal/pkg/service"
)

type Env interface {
	ID() string
	Version() string

	// http router to serve public request to the services
	HTTP() chi.Router

	// Config return the configuration
	// of your env. It can be in any
	// form like your favorite config
	// manager (e.g. viper, koanf) or
	// a simple map.
	Config() *koanf.Koanf

	// Serve the env and its services for public traffic
	Serve() error

	Join(svcs ...service.Service) error

	// Graceful shutdown of the env
	Shutdown() error

	// Health of the env
	Health() error

	// Context of the environment
	Context() context.Context

	Init() error
}
