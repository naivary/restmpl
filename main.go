package main

import (
	"log"

	"github.com/naivary/instance/internal/pkg/ctrl"
	"github.com/naivary/instance/internal/pkg/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	api, err := ctrl.New("instance.yaml")
	if err != nil {
		return err
	}

	srv, err := server.New(":8080", api.Router)
	if err != nil {
		return err
	}

	return srv.ListenAndServe()
}
