package service

import (
	"net/http"
)

type Service interface {
	UUID() string
	Name() string

	// Detailed description of the Service
	Description() string

	// Name of the router to mount it
	// to at the root level.
	Routername() string

	// Router to mount to the root router.
	Router() http.Handler
}
