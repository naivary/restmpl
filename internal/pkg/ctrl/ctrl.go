package ctrl

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/app/sys"
)

func New() App {
	return App{
		Version: "0.1.0",
		Views: Views{
			Sys: sys.Env{},
		},
	}
}

type App struct {
	Version string

	Views  Views
	Router chi.Router
}

type Views struct {
	Sys sys.Env
}
