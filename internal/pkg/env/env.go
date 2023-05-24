package env

import "net/http"

type Env interface {
	// Name of the environment
	Name() string

	// ID of the environment
	ID() string

	// All needed middlewares that should be applied to the router
	Middlewares() []func(http.Handler) http.Handler

	// Router of the env to mount onto the main router
	Router() http.Handler

	// Name of the main router: /<Routername>/<Router>
	Routername() string
}
