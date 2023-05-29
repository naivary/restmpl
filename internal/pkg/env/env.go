package env

import (
	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/monitor"
	"github.com/naivary/instance/internal/pkg/service"
)

type Env interface {
	ID() string
	Version() string

	// Running services of the environment keyed by ID.
	Services() map[string]service.Service

	// Public router to serve
	// http request to the services
	HTTP() chi.Router

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

	// Monitor agent of the
	// env for all the services
	Monitor() monitor.Agent
}
