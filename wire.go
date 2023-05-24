//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/naivary/instance/internal/app/fs"
	"github.com/naivary/instance/internal/app/sys"
	"github.com/naivary/instance/internal/pkg/config"
	"github.com/naivary/instance/internal/pkg/ctrl"
	"github.com/naivary/instance/internal/pkg/database"
	"github.com/naivary/instance/internal/pkg/filestore"
	"github.com/naivary/instance/internal/pkg/models/metadata"
	"github.com/naivary/instance/internal/pkg/routes"
	"github.com/naivary/instance/internal/pkg/service"
)

var (
	db         = wire.NewSet(database.Connect)
	svcs       = wire.NewSet(wire.Struct(new(sys.Env), "*"), wire.Struct(new(fs.Env), "*"))
	api        = wire.Struct(new(ctrl.API), "*")
	httpFs     = wire.NewSet(filestore.New, wire.Bind(new(filestore.Store), new(filestore.Filestore)))
	rootRouter = wire.NewSet(routes.New)
	k          = wire.NewSet(config.New)
	m          = wire.NewSet(metadata.New)
)

func allSvcs(sys *sys.Env, fs *fs.Env) []service.Service {
	return []service.Service{
		sys,
		fs,
	}
}

func NewApp() (*ctrl.API, error) {
	wire.Build(db, svcs, rootRouter, api, k, m, httpFs, allSvcs)
	return &ctrl.API{}, nil
}
