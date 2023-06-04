package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

type Service interface {
	// Unique identifier of the service.
	ID() string

	// Unique Human friendly name of the service.
	Name() string

	// Detailed description of the service
	Description() string

	// HTTP router to serve public request
	HTTP() chi.Router

	// Recommended pattern to mount the router to.
	Pattern() string

	// Health returns the health status
	// of the service. If the error is
	// non nil the service is considered unhealthy.
	// The information providede will be served
	// to the requester.
	Health() (*Info, error)

	// Metrics returns the service specific
	// collected metrics. Probably Prometheus in our case
	Metrics() []prometheus.Collector

	// Initialize the service given the global dependencies
	Init() error

	// Graceful shutdown of the service
	Shutdown() error
}

type Info struct {
	ID   string `json:"service"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}
