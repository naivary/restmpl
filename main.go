package main

import (
	"log"

	"github.com/naivary/instance/internal/pkg/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	app, err := NewApp()
	if err != nil {
		return err
	}

	srv, err := server.New(":8080", app.Router)
	if err != nil {
		return err
	}
	return srv.ListenAndServe()
}
