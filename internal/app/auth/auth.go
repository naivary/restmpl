package auth

import (
	"errors"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/restmpl/internal/pkg/jwtauth"
	"github.com/naivary/restmpl/internal/pkg/logging"
	"github.com/naivary/restmpl/internal/pkg/metrics"
	"github.com/naivary/restmpl/internal/pkg/service"
	"github.com/pocketbase/dbx"
	"github.com/prometheus/client_golang/prometheus"
)

var _ service.Service = (*Auth)(nil)

type Auth struct {
	K  *koanf.Koanf
	DB *dbx.DB

	l   logging.Manager
	m   metrics.Managee
	jwt *jwtauth.JWTAuth
}

func (a Auth) Name() string {
	return "authentication"
}

func (a Auth) ID() string {
	return uuid.NewString()
}

func (a Auth) Description() string {
	return "jwt authentication"
}

func (a Auth) Pattern() string {
	return "/auth"
}

func (a *Auth) Init() error {
	l, err := logging.NewSvcManager(a.K, a)
	if err != nil {
		return err
	}
	a.l = l

	a.m = metrics.NewManagee()

	secret := os.Getenv("API_JWT_SECRET")
	if secret == "" {
		return errors.New("API_JWT_SECRET env variable is not set")
	}
	a.jwt = jwtauth.New(a.K)
	return nil
}

func (a Auth) Health() (*service.Info, error) {
	if a.jwt == nil {
		return nil, errors.New("missing jwt auth instanec")
	}
	if a.l == nil {
		return nil, errors.New("missing logger manager")
	}
	if a.m == nil {
		return nil, errors.New("missing metrics managee")
	}
	return nil, nil
}

func (a Auth) Metrics() []prometheus.Collector {
	return nil
}

func (a Auth) HTTP() chi.Router {
	r := chi.NewRouter()
	r.Post("/signin", a.signin)
	return r
}

func (a Auth) Shutdown() error {
	return nil
}
