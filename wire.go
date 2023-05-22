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
	"github.com/naivary/instance/internal/pkg/models/metadata"
	"github.com/naivary/instance/internal/pkg/routes"
)

var (
	db    = wire.NewSet(database.Connect)
	views = wire.NewSet(wire.Struct(new(sys.Env), "*"), wire.Struct(new(fs.Env), "*"), wire.Struct(new(ctrl.Views), "*"))
	k     = wire.NewSet(config.New)
	app   = wire.Struct(new(ctrl.App), "*")
	m     = wire.NewSet(metadata.New)
)

func NewApp() (*ctrl.App, error) {
	wire.Build(db, views, routes.New, app, k, m)
	return &ctrl.App{}, nil
}
