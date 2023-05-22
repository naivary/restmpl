package main

import (
	"log"

	"github.com/naivary/instance/internal/pkg/config"
	"github.com/naivary/instance/internal/pkg/ctrl"
	"github.com/naivary/instance/internal/pkg/database"
	"github.com/naivary/instance/internal/pkg/routes"
	"github.com/naivary/instance/internal/pkg/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	app := ctrl.New()
	app.Router = routes.New(&app.Views)

	srv, err := server.New(":8080", app.Router)
	if err != nil {
		return err
	}

	k, err := config.New()
	if err != nil {
		return err
	}
	app.Views.Sys.K = k

	db, err := database.Connect()
	if err != nil {
		return err
	}
	app.Views.Sys.DB = db

	return srv.ListenAndServe()
}
