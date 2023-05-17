package ctrl

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/app/sys"
)

func New() App {
	return App{
		Views: Views{
			Sys: sys.Env{},
		},
	}
}

type App struct {
	// Views contains all handler
	// for the corresponding endpoints.
	// Every Handler in the View is represented
	// by a directory in the /internal/app/<handler>
	// and the needed Env of the handler.
	Views Views

	// Router contains all the endpoints of
	// which define the REST-API.
	Router chi.Router
}

type Views struct {
	Sys sys.Env
}
