package env

import (
	"github.com/go-chi/chi"
	"github.com/naivary/instance/internal/pkg/service"
)

type API struct {
	Services []service.Service

	// Router contains all the endpoints of
	// which define the REST-API.
	Router chi.Router
}
