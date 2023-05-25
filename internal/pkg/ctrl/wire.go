package ctrl

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
	"github.com/naivary/instance/internal/pkg/services"
)

var (
	db     = wire.NewSet(database.Connect)
	svc    = wire.NewSet(wire.Struct(new(sys.Sys), "*"), wire.Struct(new(fs.Fs), "*"), wire.Struct(new(services.Services), "*"))
	app    = wire.Struct(new(ctrl.API), "*")
	httpFs = wire.NewSet(filestore.New, wire.Bind(new(filestore.Store), new(filestore.Filestore)))
	k      = wire.NewSet(config.New)
	m      = wire.NewSet(metadata.New)
)

func New() (*ctrl.API, error) {
	wire.Build(db, svc, routes.New, app, k, m, httpFs)
	return &ctrl.API{}, nil
}
