package service

import "github.com/go-chi/chi/v5"

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
	// Health() (Info, error)

	// Metrics returns the service specific
	// collected metrics. Probably Prometheus in our case
	// Metrics()
}

type Info struct {
	ID   string
	Name string
	Desc string
	// Dependencies of the Service keyed
	// by the name of the dependency
	Deps map[string]string
}
