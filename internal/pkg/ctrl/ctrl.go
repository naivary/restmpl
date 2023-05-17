package ctrl

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/app/sys"
)

type App struct {
	Views  Views
	Router chi.Router
}

type Views struct {
	Sys sys.Env
}
