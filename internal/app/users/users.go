package users

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/apitmpl/internal/pkg/jwtauth"
	"github.com/naivary/apitmpl/internal/pkg/logging"
	"github.com/naivary/apitmpl/internal/pkg/metrics"
	"github.com/naivary/apitmpl/internal/pkg/service"
	"github.com/pocketbase/dbx"
	"github.com/prometheus/client_golang/prometheus"
)

var _ service.Service = (*Users)(nil)

type Users struct {
	K  *koanf.Koanf
	DB *dbx.DB

	l logging.Manager
	m metrics.Managee
}

func (u Users) ID() string {
	return uuid.NewString()
}

func (u Users) Name() string {
	return "user"
}

func (u Users) Description() string {
	return "create and manage identities and access to ressources"
}

func (u *Users) Init() error {
	mngr, err := logging.NewSvcManager(u.K, u)
	if err != nil {
		return err
	}
	u.l = mngr
	u.m = metrics.NewManagee()
	return nil
}

func (u Users) Health() (*service.Info, error) {
	if u.l == nil {
		return nil, errors.New("missing loggin manager")
	}
	if u.m == nil {
		return nil, errors.New("missing metrics")
	}
	return nil, nil
}

func (u Users) Metrics() []prometheus.Collector {
	return nil
}

func (u Users) HTTP() chi.Router {
	r := chi.NewRouter()
	r.Post("/", u.create)
	r.Get("/{userID}", u.single)
	r.Get("/list", u.list)
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verify)
		r.Delete("/delete", u.delete)
	})
	return r
}

func (u Users) Pattern() string {
	return "/users"
}

func (u Users) Shutdown() error {
	u.l.Shutdown()
	return nil
}
