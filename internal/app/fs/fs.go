package fs

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/filestore"
	"github.com/naivary/instance/internal/pkg/service"
	"github.com/spf13/afero"
)

var _ service.Service = (*Fs)(nil)

type Fs struct {
	K *koanf.Koanf

	Store filestore.Store[afero.File]
}

func (f Fs) Register(root chi.Router) {
	r := chi.NewRouter()
	for _, mw := range f.Middlewares() {
		r.Use(mw)
	}
	r.Post("/", f.Create)
	r.Delete("/remove", f.Remove)
	r.Get("/read", f.Read)
	root.Mount("/fs", r)
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
