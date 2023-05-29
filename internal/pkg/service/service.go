package service

import (
	"github.com/go-chi/chi/v5"
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

	// Metrics returns the service specific
	// collected metrics. Probably Prometheus in our case
	Metrics() error
}

type Info struct {
	ID   string `jsonapi:"attr,id"`
	Name string `jsonapi:"attr,name"`
	Desc string `jsonapi:"attr,desc"`
}
