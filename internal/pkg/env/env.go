package env

import (
	"github.com/naivary/instance/internal/pkg/service"
)

// TODO(naivary): second generic Parameter
// for chi.Router
type Env[T any, R any] interface {
	UUID() string
	Version() string

	Services() map[string]service.Service[R]

	// Router returns the root
	// router of the application
	// for public traffic
	Router() R

	// Config return the config Manager of
	// your application.
	Config() T
}
