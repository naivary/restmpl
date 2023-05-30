package fs

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/filestore"
	"github.com/naivary/instance/internal/pkg/logging"
	"github.com/naivary/instance/internal/pkg/service"
	"github.com/spf13/afero"
)

var _ service.Service = (*Fs)(nil)

type Fs struct {
	service.Service
	k       *koanf.Koanf
	logMngr logging.Manager
	store   filestore.Store[afero.File]
}

func (f Fs) Health() (*service.Info, error) {
	f.k = koanf.New(".")
	if f.k == nil {
		return nil, errors.New("missing config manager")
	}
	return &service.Info{
		ID:   f.ID(),
		Name: f.Name(),
		Desc: f.Description(),
	}, nil
}

func (f Fs) HTTP() chi.Router {
	r := chi.NewRouter()
	for _, mw := range f.Middlewares() {
		r.Use(mw)
	}
	r.Post("/", f.Create)
	r.Delete("/remove", f.Remove)
	r.Get("/read", f.Read)
	return r
}

func (f Fs) ID() string {
	return uuid.NewString()
}

func (f Fs) Name() string {
	return "filestore"
}

func (f Fs) Pattern() string {
	return "/fs"
}
