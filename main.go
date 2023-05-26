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

	srv, err := server.New(api.K, api.Router)
	if err != nil {
		return err
	}

	return srv.ListenAndServeTLS(api.K.String("server.crt"), api.K.String("server.key"))
}
