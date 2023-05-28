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
	// Health(http.ResponseWriter, *http.Request) error

	// Metrics returns the service specific
	// collected metrics. Probably Prometheus in our case
	// Metrics()
}
