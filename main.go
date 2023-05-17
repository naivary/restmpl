package main

import (
	"log"

	"github.com/naivary/instance/internal/app/sys"
	"github.com/naivary/instance/internal/pkg/ctrl"
	"github.com/naivary/instance/internal/pkg/routes"
	"github.com/naivary/instance/internal/pkg/server"
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

	srv, err := server.New(":8080", app.Router)
	if err != nil {
		return err
	}

	return srv.ListenAndServe()
}
