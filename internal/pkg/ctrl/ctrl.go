package ctrl

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/services"
)

type API struct {
	// Services contains all handler
	// for the corresponding endpoints.
	// Every Handler in the View is represented
	// by a directory in the /internal/app/<handler>
	// and the needed Env of the handler.
	Services services.Services

	// Router contains all the endpoints of
	// which define the REST-API.
	Router chi.Router
}
