package fs

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/filestore"
	"github.com/naivary/instance/internal/pkg/service"
	"github.com/pocketbase/dbx"
	"github.com/spf13/afero"
)

var _ service.Service = (*Fs)(nil)

type Fs struct {
	k *koanf.Koanf

	store filestore.Store[afero.File]
}

func (f Fs) Metrics() error {
	return nil
}

func (f Fs) Health() (*service.Info, error) {
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

func (f Fs) Pattern() string {
	return "/fs"
}

func (f Fs) Name() string {
	return "filestore"
}

func (f Fs) ID() string {
	return uuid.NewString()
}

func (f Fs) Description() string {
	return "a simple filestore which uses the host filesystem as a sotrage"
}

func (f *Fs) Init(k *koanf.Koanf, db *dbx.DB) error {
	f.k = k
	fstore, err := filestore.New(k)
	if err != nil {
		return err
	}
	f.store = fstore
	return nil
}
