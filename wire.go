///go:build wireinject
//// +build wireinject

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
	esvcs  = wire.NewSet(wire.Value([]service.Service(nil)))
	db     = wire.NewSet(database.Connect)
	svcs   = wire.NewSet(wire.Struct(new(sys.Env), "*"), wire.Struct(new(fs.Env), "*"))
	api    = wire.Struct(new(ctrl.API), "*")
	httpFs = wire.NewSet(filestore.New, wire.Bind(new(filestore.Store), new(filestore.Filestore)))
	k      = wire.NewSet(config.New)
	m      = wire.NewSet(metadata.New)
)

func NewApp() (*ctrl.API, error) {
	wire.Build(db, svcs, routes.New, api, k, m, httpFs, esvcs)
	return &ctrl.API{}, nil
}
