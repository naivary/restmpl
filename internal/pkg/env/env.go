package env

import (
	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/service"
)

type Env interface {
	// Unique ID of the environment
	ID() string

	// Current version of the environment
	Version() string

	// Services returns the running services
	// of the environment keyed by ID.
	Services() map[string]service.Service

	// Router returns the public
	// router to serve public traffic
	Router() chi.Router

	// Config return the configuration
	// of your env. It can be in any
	// form like your favorite config
	// manager (e.g. viper, koanf) or
	// a simple map.
	Config() *koanf.Koanf

	// Run will prepare and
	// run the env to accept
	// possible public traffic.
	Run() error
}
