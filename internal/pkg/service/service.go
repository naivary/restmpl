package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/register"
	"github.com/pocketbase/dbx"
)

type Service interface {
	// Unique identifier of the service.
	ID() string

	// Human friendly name of the service.
	Name() string

	// Detailed description of the service
	Description() string

	// HTTP router to serve public request
	HTTP() chi.Router

	// Recommended pattern to use router mountage
	Pattern() string

	// Health returns the health status
	// of the service. If the error is
	// non nil the service is considered unhealthy.
	Health(register.Register) (*Info, error)

	// Metrics returns the service specific
	// collected metrics. Probably Prometheus in our case
	Metrics() error

	// Initialize the service given the global dependencies
	Init(*koanf.Koanf, *dbx.DB) error
}

type Info struct {
	ID   string `jsonapi:"attr,id"`
	Name string `jsonapi:"attr,name"`
	Desc string `jsonapi:"attr,desc"`
}
