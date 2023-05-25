// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ctrl

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
	"github.com/naivary/instance/internal/pkg/service"
)

// Injectors from wire.go:

func New(cfgFile string) (*API, error) {
	koanf, err := config.New(cfgFile)
	if err != nil {
		return nil, err
	}
	sqlDB, err := database.Connect(koanf)
	if err != nil {
		return nil, err
	}
	metadataMetadata := metadata.New(koanf)
	sysSys := &sys.Sys{
		K:  koanf,
		DB: sqlDB,
		M:  metadataMetadata,
	}
	filestoreFilestore, err := filestore.New(koanf)
	if err != nil {
		return nil, err
	}
	fsFs := &fs.Fs{
		K:     koanf,
		Store: filestoreFilestore,
	}
	v := allSvcs(sysSys, fsFs)
	router := routes.New(v)
	ctrlAPI := &API{
		Services: v,
		Router:   router,
	}
	return ctrlAPI, nil
}

// wire.go:

var (
	db         = wire.NewSet(database.Connect)
	svcs       = wire.NewSet(wire.Struct(new(sys.Sys), "*"), wire.Struct(new(fs.Fs), "*"))
	api        = wire.Struct(new(API), "*")
	httpFs     = wire.NewSet(filestore.New, wire.Bind(new(filestore.Store), new(filestore.Filestore)))
	rootRouter = wire.NewSet(routes.New)
	k          = wire.NewSet(config.New)
	m          = wire.NewSet(metadata.New)
)

func allSvcs(sys2 *sys.Sys, fs2 *fs.Fs) []service.Service[chi.Router] {
	return []service.Service[chi.Router]{sys2, fs2}
}