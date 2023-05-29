package main

import (
	"errors"
	"os"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/app/fs"
	"github.com/naivary/instance/internal/pkg/config"
	"github.com/naivary/instance/internal/pkg/database"
	"github.com/naivary/instance/internal/pkg/env"
	"github.com/naivary/instance/internal/pkg/filestore"
	"github.com/naivary/instance/internal/pkg/service"
	"github.com/pocketbase/dbx"
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
	slog.Info("running the env", "usedCfgFile", cfgFile)
	return e.Serve()
}

func newEnv(cfgFile string) (env.Env, error) {
	k, err := config.New(cfgFile)
	if err != nil {
		return nil, err
	}
	db, err := database.Connect(k)
	if err != nil {
		return nil, err
	}
	fstore, err := filestore.New(k)
	if err != nil {
		return nil, err
	}
	f := &fs.Fs{
		K:     nil,
		Store: fstore,
	}
	svcs := []service.Service{f}
	api := env.NewAPI(svcs, k, db)
	return &api, nil
}

func depsReg(db *dbx.DB) {
	kPing := func(k *koanf.Koanf) error {
		if k == nil {
			return errors.New("config manager is not present")
		}
		return nil
	}
	dbPing := func(db *dbx.DB) error {
		return db.DB().Ping()
	}
	fsStore := func(fs *filestore.Filestore) error {
		if fs == nil {
			return errors.New("filestore missong")
		}
		return nil
	}

}
