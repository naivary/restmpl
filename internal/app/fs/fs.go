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
	K *koanf.Koanf

	logManager logging.Manager
	store      filestore.Store[afero.File]
}

func (f Fs) Health() (*service.Info, error) {
	if f.K == nil {
		return nil, errors.New("missing config manager")
	}

	if f.store == nil {
		return nil, errors.New("missing filestore")
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
	r.Post("/", f.create)
	r.Delete("/remove", f.remove)
	r.Get("/read", f.read)
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

func (f *Fs) Init() error {
	fstore, err := filestore.New(f.K)
	if err != nil {
		return err
	}
	f.store = fstore

	mngr, err := logging.NewSvcManager(f.K, f)
	if err != nil {
		return err
	}
	f.logManager = mngr
	return nil
}

func (f Fs) Shutdown() error {
	f.logManager.Shutdown()
	return nil
}

func (f Fs) Metrics() error {
	return nil
}

func (f Fs) Description() string {
	return "simple file storage based upon sqlite"
}
