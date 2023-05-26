package main

import (
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

func run() error {
	api, err := ctrl.New(cfgFile)
	if err != nil {
		return err
	}
	srv, err := server.New(api.K, api.Router)
	if err != nil {
		return err
	}
	slog.Info("starting server", "cfgFile", cfgFile)
	return srv.ListenAndServeTLS(api.K.String("server.crt"), api.K.String("server.key"))
}
