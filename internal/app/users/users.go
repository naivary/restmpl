package users

import (
	"errors"

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
	u.l = logging.NewSvcManager(u)
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
	u.m.AddCounter("incomingReq", metrics.IncomingHTTPRequest(&u))
	return u.m.All()
}

func (u Users) HTTP() chi.Router {
	r := chi.NewRouter()
	for _, mw := range u.middlewares() {
		r.Use(mw)
	}
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
	return nil
}
