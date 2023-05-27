//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/wire"
	"github.com/knadh/koanf/v2"

	"github.com/naivary/instance/internal/app/fs"
	"github.com/naivary/instance/internal/app/sys"
	"github.com/naivary/instance/internal/pkg/config"
	"github.com/naivary/instance/internal/pkg/database"
	"github.com/naivary/instance/internal/pkg/env"
	"github.com/naivary/instance/internal/pkg/filestore"
	"github.com/naivary/instance/internal/pkg/service"
)

var (
	db     = wire.NewSet(database.Connect)
	svcs   = wire.NewSet(wire.Struct(new(sys.Sys), "K", "DB"), wire.Struct(new(fs.Fs), "*"))
	httpFs = wire.NewSet(filestore.New, wire.Bind(new(filestore.Store), new(filestore.Filestore)))
	e      = wire.NewSet(env.NewAPI, wire.Bind(new(env.Env[*koanf.Koanf]), new(env.API)))
	k      = wire.NewSet(config.New)
)

func allSvcs(sys *sys.Sys, fs *fs.Fs) []service.Service[chi.Router] {
	return []service.Service[chi.Router]{
		sys,
		fs,
	}
}

func newEnv(cfgFile string) (env.API, error) {
	wire.Build(db, svcs, k, httpFs, allSvcs, e)
	return env.API{}, nil
}
