// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/wire"
	"github.com/naivary/instance/internal/app/fs"
	"github.com/naivary/instance/internal/app/sys"
	"github.com/naivary/instance/internal/pkg/config"
	"github.com/naivary/instance/internal/pkg/database"
	"github.com/naivary/instance/internal/pkg/filestore"
	"github.com/naivary/instance/internal/pkg/models/metadata"
	"github.com/naivary/instance/internal/pkg/routes"
	"github.com/naivary/instance/internal/pkg/services"
)

// Injectors from wire.go:

func NewApp() (*App, error) {
	koanf, err := config.New()
	if err != nil {
		return nil, err
	}
	sqlDB, err := database.Connect()
	if err != nil {
		return nil, err
	}
	metadataMetadata := metadata.New(koanf)
	env := sys.Env{
		K:  koanf,
		DB: sqlDB,
		M:  metadataMetadata,
	}
	filestoreFilestore := filestore.New(koanf)
	fsEnv := fs.Env{
		Fs: filestoreFilestore,
		K:  koanf,
	}
	servicesServices := services.Services{
		Sys: env,
		Fs:  fsEnv,
	}
	services2 := &services.Services{
		Sys: env,
		Fs:  fsEnv,
	}
	router := routes.New(services2)
	mainApp := &App{
		Services: servicesServices,
		Router:   router,
	}
	return mainApp, nil
}

// wire.go:

type App struct {
	// Services contains all handler
	// for the corresponding endpoints.
	// Every Handler in the View is represented
	// by a directory in the /internal/app/<handler>
	// and the needed Env of the handler.
	Services services.Services

	// Router contains all the endpoints of
	// which define the REST-API.
	Router chi.Router
}

var (
	db     = wire.NewSet(database.Connect)
	svc    = wire.NewSet(wire.Struct(new(sys.Env), "*"), wire.Struct(new(fs.Env), "*"), wire.Struct(new(services.Services), "*"))
	app    = wire.Struct(new(App), "*")
	httpFs = wire.NewSet(filestore.New)
	k      = wire.NewSet(config.New)
	m      = wire.NewSet(metadata.New)
)
