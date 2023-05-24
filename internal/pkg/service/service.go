package service

import (
	"net/http"
)

type Service interface {
	UUID() string
	Name() string

	// Detailed description of the service
	Description() string

	// Pattern of the domain like /sys
	Pattern() string

	// Handler to mount to the root router.
	Router() http.Handler
}
