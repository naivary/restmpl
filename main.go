package main

import (
	"errors"
	"log"
	"os"

	"github.com/naivary/instance/internal/pkg/ctrl"
	"github.com/naivary/instance/internal/pkg/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func getCfgFile() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("missing config file as the first argument")
	}
	return os.Args[1], nil
}

func run() error {
	cfgFile, err := getCfgFile()
	if err != nil {
		return err
	}

	api, err := ctrl.New(cfgFile)
	if err != nil {
		return err
	}

	srv, err := server.New(api.K, api.Router)
	if err != nil {
		return err
	}
	return srv.ListenAndServeTLS(api.K.String("server.crt"), api.K.String("server.key"))
}
