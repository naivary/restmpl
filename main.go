package main

import (
	"log"
	"net/http"

	"github.com/naivary/instance/internal/app/sys"
	"github.com/naivary/instance/internal/pkg/ctrl"
	"github.com/naivary/instance/internal/pkg/routes"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	app := ctrl.App{
		Views: ctrl.Views{
			Sys: sys.Env{},
		},
	}
	app.Router = routes.New(&app.Views)

	return http.ListenAndServe(":8080", app.Router)
}
