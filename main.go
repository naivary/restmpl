package main

import (
	"errors"
	"os"

	"github.com/naivary/instance/internal/pkg/ctrl"
	"github.com/naivary/instance/internal/pkg/server"
	"golang.org/x/exp/slog"
)

const (
	cfgFile = "instance.yaml"
)

func main() {
	if err := run(); err != nil {
		slog.Error("something went wrong while starting the sevrer", "err", err.Error())
	}
}

func getCfgFile() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("missing config file as the first argument")
	}
	return os.Args[1], nil
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
