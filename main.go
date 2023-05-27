package main

import (
	"errors"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/app/fs"
	"github.com/naivary/instance/internal/app/sys"
	"github.com/naivary/instance/internal/pkg/config"
	"github.com/naivary/instance/internal/pkg/database"
	"github.com/naivary/instance/internal/pkg/env"
	"github.com/naivary/instance/internal/pkg/filestore"
	"github.com/naivary/instance/internal/pkg/models/metadata"
	"github.com/naivary/instance/internal/pkg/server"
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

func newEnv(cfgFile string) (env.API, error) {
	k, err := config.New(cfgFile)
	if err != nil {
		return env.API{}, err
	}
	db, err := database.Connect(k)
	if err != nil {
		return env.API{}, err
	}
	m := metadata.New(k, db)
	s := &sys.Sys{
		K:  k,
		DB: db,
		M:  m,
	}
	fstore, err := filestore.New(k)
	if err != nil {
		return env.API{}, err
	}
	f := &fs.Fs{
		K:     k,
		Store: fstore,
	}
	svcs := []service.Service[chi.Router]{s, f}
	api := env.NewAPI(svcs, k)
	return api, nil
}
