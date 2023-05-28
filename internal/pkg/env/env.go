package env

import (
	"github.com/naivary/instance/internal/pkg/service"
)

type Env[T any, R any] interface {
	// Unique ID of the environment
	ID() string

	// Current version of the environment
	Version() string

	// Services returns the running services
	// of the environment keyed by the ID
	// of the services.
	Services() map[string]service.Service[R]

	// Router returns the current
	// router which serves possible
	// public traffic.
	Router() R

	// Config return the configuration
	// of your env. It can be in any
	// form like your favorite config
	// manager (e.g. viper, koanf) or
	// a simple map.
	Config() T

	// Run will prepare and
	// run the env to accept
	// possible public traffic.
	Run() error
}
