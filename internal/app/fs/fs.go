package fs

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/filestore"
	"github.com/naivary/instance/internal/pkg/service"
)

var _ service.Service = (*Env)(nil)

type Env struct {
	K *koanf.Koanf

	Store filestore.Store
}

func (e Env) Name() string {
	return "filestore"
}

func (e Env) UUID() string {
	return uuid.NewString()
}

func (e Env) Pattern() string {
	return "/fs"
}

func (e Env) Router() http.Handler {
	r := chi.NewRouter()
	for _, mw := range e.Middlewares() {
		r.Use(mw)
	}
	r.Post("/", e.Create)
	r.Delete("/remove", e.Remove)
	r.Get("/read", e.Read)
	return r
}

func (e Env) Description() string {
	return "a simple filestore which uses the host filesystem as a sotrage"
}
