package env

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/server"
	"github.com/naivary/instance/internal/pkg/service"
	"github.com/pocketbase/dbx"
)

const (
	reqTimeout = 20 * time.Second
)

var _ Env = (*API)(nil)

// API environment which
// implements the Env interface.
type API struct {
	svcs   []service.Service
	k      *koanf.Koanf
	router chi.Router
	db     *dbx.DB
}

func NewAPI(svcs []service.Service, k *koanf.Koanf, db *dbx.DB) API {
	return API{
		svcs: svcs,
		k:    k,
		db:   db,
	}
}

func (a API) ID() string {
	return uuid.NewString()
}

func (a API) Version() string {
	return a.k.String("version")
}

func (a *API) Router() chi.Router {
	if a.router != nil {
		return a.router
	}
	root := chi.NewRouter()
	root.Use(middleware.SetHeader("Content-Type", jsonapi.MediaType))
	root.Use(middleware.RequestID)
	root.Use(middleware.CleanPath)
	root.Use(middleware.Timeout(reqTimeout))
	for _, svc := range a.svcs {
		svc.Register(root)
	}
	a.router = root
	return root
}

func (a API) Config() *koanf.Koanf {
	return a.k
}

func (a API) Services() map[string]service.Service {
	m := make(map[string]service.Service, len(a.svcs))
	for _, svc := range a.svcs {
		m[svc.ID()] = svc
	}
	return m
}

func (a API) Run() error {
	srv, err := server.New(a.k, a.Router())
	if err != nil {
		return err
	}
	return srv.ListenAndServeTLS(a.k.String("server.crt"), a.k.String("server.key"))
}
