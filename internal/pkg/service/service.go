package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/register"
)

type Service interface {
	// Unique identifier of the service.
	ID() string

	// Human friendly name of the service.
	Name() string

	// Detailed description of the service
	Description() string

	// Register registers the service
	// to the public router of type T
	Register(chi.Router)

	// Health returns the health status
	// of the service. If the error is
	// non nil the service is considered unhealthy.
	// If the service is healthy, some information
	// about the service will be provided.
	Health(register.Register) (*Info, error)

	// Metrics returns the service specific
	// collected metrics. Probably Prometheus in our case
	Metrics() error
}

type Info struct {
	ID   string `jsonapi:"attr,id"`
	Name string `jsonapi:"attr,name"`
	Desc string `jsonapi:"attr,desc"`
}
