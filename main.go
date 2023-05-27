package main

import (
	"errors"
	"os"

	"github.com/naivary/instance/internal/pkg/server"
	"golang.org/x/exp/slog"
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
	cfgFile, err := getCfgFile()
	if err != nil {
		return err
	}
	api, err := newEnv(cfgFile)
	if err != nil {
		return err
	}
	srv, err := server.New(api.Config(), api.Router())
	if err != nil {
		return err
	}
	slog.Info("starting the server", "usedCfgFile", cfgFile)
	return srv.ListenAndServeTLS(api.Config().String("server.crt"), api.Config().String("server.key"))
}
