package main

import (
	"errors"
	"os"

	"github.com/naivary/instance/internal/app/fs"
	"github.com/naivary/instance/internal/pkg/env"
	"github.com/naivary/instance/internal/pkg/service"
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
	e, err := newEnv(cfgFile)
	if err != nil {
		return err
	}
	err = e.Init()
	if err != nil {
		return err
	}
	slog.Info("running the env", "usedCfgFile", cfgFile)
	return e.Serve()
}

func newEnv(cfgFile string) (env.Env, error) {
	f := &fs.Fs{}
	svcs := []service.Service{f}
	api := env.NewAPI(cfgFile, svcs)
	return &api, nil
}
