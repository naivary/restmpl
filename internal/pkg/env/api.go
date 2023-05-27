package env

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/service"
)

const (
	reqTimeout = 20 * time.Second
)

var _ Env[*koanf.Koanf] = (*API)(nil)

// API is API env and
// implements the Env
// interface
type API struct {
	services   []service.Service[chi.Router]
	rootRouter chi.Router
	k          *koanf.Koanf
}

func NewAPI(svcs []service.Service[chi.Router], k *koanf.Koanf) API {
	return API{
		services: svcs,
		k:        k,
	}
}

func (a API) UUID() string {
	return uuid.NewString()
}

func (a API) Version() string {
	return a.k.String("version")
}

func (a API) Router() http.Handler {
	root := chi.NewRouter()
	root.Use(middleware.SetHeader("Content-Type", jsonapi.MediaType))
	root.Use(middleware.RequestID)
	root.Use(middleware.CleanPath)
	root.Use(middleware.Timeout(reqTimeout))
	for _, svc := range a.services {
		svc.RegisterRootMiddleware(root)
	}
	for _, svc := range a.services {
		svc.Register(root)
	}

	a.rootRouter = root
	return a.rootRouter
}

func (a API) Config() *koanf.Koanf {
	return a.k
}

func (a API) Services() map[string]service.Service[chi.Router] {
	m := make(map[string]service.Service[chi.Router], len(a.services))
	for _, svc := range a.services {
		m[svc.UUID()] = svc
	}
	return m
}
