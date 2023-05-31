package main

import (
	"errors"
	"os"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/app/fs"
	"github.com/naivary/instance/internal/pkg/env"
	"github.com/naivary/instance/internal/pkg/log"
	"github.com/naivary/instance/internal/pkg/service"
	"github.com/pocketbase/dbx"
	"golang.org/x/exp/slog"
)

func main() {
	if err := run(); err != nil {
		slog.Error("something went wrong while starting the sevrer", slog.String("err", err.Error()))
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
	slog.Info("serving the api", "used_config_file", cfgFile)
	return e.Serve()
}

func newEnv(cfgFile string) (env.Env, error) {
	api, err := env.NewAPI(cfgFile)
	if err != nil {
		return nil, err
	}
	svcs, err := createServices(api.Config(), api.DB())
	if err != nil {
		return nil, err
	}
	if err := api.Join(svcs...); err != nil {
		return nil, err
	}
	return api, nil
}

func createServices(k *koanf.Koanf, db *dbx.DB) ([]service.Service, error) {
	svcs := make([]service.Service, 0)
	f := new(fs.Fs)
	f.K = k
	svcs = append(svcs, f)

	for _, svc := range svcs {
		mgr, err := log.New(k, svc)
		if err != nil {
			return nil, err
		}
		f.LogManager = mgr
	}
	return svcs, nil
}
