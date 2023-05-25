package fs

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/filestore"
	"github.com/naivary/instance/internal/pkg/service"
)

<<<<<<< HEAD
=======
var _ service.Service[chi.Router] = (*Fs)(nil)

>>>>>>> svc-interface
type Fs struct {
	K *koanf.Koanf

	Store filestore.Store
}

func (f Fs) Register(root chi.Router) {
	root.Mount(f.pattern(), f.router())
}

func (f Fs) Name() string {
	return "filestore"
}

func (f Fs) UUID() string {
	return uuid.NewString()
}

func (f Fs) pattern() string {
	return "/fs"
}

func (f Fs) router() http.Handler {
	r := chi.NewRouter()
	for _, mw := range f.Middlewares() {
		r.Use(mw)
	}
	r.Post("/", f.Create)
	r.Delete("/remove", f.Remove)
	r.Get("/read", f.Read)
	return r
}

func (f Fs) Description() string {
	return "a simple filestore which uses the host filesystem as a sotrage"
}
