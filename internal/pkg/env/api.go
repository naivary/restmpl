package env

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/dependency"
	"github.com/naivary/instance/internal/pkg/monitor"
	"github.com/naivary/instance/internal/pkg/server"
	"github.com/naivary/instance/internal/pkg/service"
	"github.com/pocketbase/dbx"
)

const (
	reqTimeout = 20 * time.Second
)

var _ Env = (*API)(nil)

type API struct {
	svcs     []service.Service
	k        *koanf.Koanf
	http     chi.Router
	monAgent monitor.Agent
	db       *dbx.DB
}

func NewAPI(svcs []service.Service, k *koanf.Koanf, db *dbx.DB, deps []dependency.Pinger) API {
	return API{
		svcs:     svcs,
		k:        k,
		monAgent: monitor.New(svcs, deps),
		db:       db,
	}
}

func (a API) Monitor() monitor.Agent {
	return a.monAgent
}

func (a API) ID() string {
	return uuid.NewString()
}

func (a API) Version() string {
	return a.k.String("version")
}

func (a *API) HTTP() chi.Router {
	if a.http != nil {
		return a.http
	}
	root := chi.NewRouter()
	root.Use(middleware.SetHeader("Content-Type", jsonapi.MediaType))
	root.Use(middleware.RequestID)
	root.Use(middleware.CleanPath)
	root.Use(middleware.Timeout(reqTimeout))
	for _, svc := range a.svcs {
		svc.Register(root)
	}
	root.Mount("/sys", a.Monitor().HTTP())
	a.http = root
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

func (a API) Serve() error {
	srv, err := server.New(a.k, a.HTTP())
	if err != nil {
		return err
	}
	return srv.ListenAndServeTLS(a.k.String("server.crt"), a.k.String("server.key"))
}
