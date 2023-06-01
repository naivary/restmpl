package env

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/config"
	"github.com/naivary/instance/internal/pkg/database"
	"github.com/naivary/instance/internal/pkg/server"
	"github.com/naivary/instance/internal/pkg/service"
	"github.com/pocketbase/dbx"
	"golang.org/x/exp/slog"
)

const (
	reqTimeout = 20 * time.Second
)

var _ Env = (*API)(nil)

type API struct {
	cfgFile  string
	isInited bool

	// global dependencies
	db *dbx.DB
	k  *koanf.Koanf

	svcs []service.Service
	http chi.Router
	srv  *http.Server
	ctx  context.Context
}

func NewAPI(cfgFile string) (*API, error) {
	a := &API{
		cfgFile: cfgFile,
		ctx:     context.Background(),
	}
	if err := a.init(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a API) DB() *dbx.DB {
	return a.db
}

func (a API) ID() string {
	return uuid.NewString()
}

func (a API) Version() string {
	return a.k.String("version")
}

func (a *API) initHTTP() {
	root := chi.NewRouter()
	root.Use(middleware.SetHeader("Content-Type", jsonapi.MediaType))
	root.Use(middleware.RequestID)
	root.Use(middleware.CleanPath)
	root.Use(middleware.Timeout(reqTimeout))
	a.http = root
}

func (a *API) HTTP() chi.Router {
	if a.http != nil {
		return a.http
	}
	a.initHTTP()
	return a.http
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

// Serve will start the http server for public
// requests. It also handles the graceful shutdown
// of OS Interrupts signals.
func (a *API) Serve() error {
	srv, err := server.New(a.k, a.HTTP())
	if err != nil {
		return err
	}
	go func() {
		if err := srv.ListenAndServeTLS(a.k.String("server.crt"), a.k.String("server.key")); err != nil {
			return
		}
	}()
	slog.InfoCtx(
		a.ctx,
		"Starting the http server",
		slog.String("api_name", a.k.String("name")),
		slog.String("version", a.k.String("version")),
		slog.String("used_config_file", a.cfgFile),
	)
	a.srv = srv
	return nil
}

func (a *API) init() error {
	if a.isInited {
		return nil
	}

	k, err := config.New(a.cfgFile)
	if err != nil {
		return err
	}
	a.k = k

	db, err := database.Connect(k)
	if err != nil {
		return err
	}
	a.db = db

	a.initHTTP()
	a.isInited = true
	return nil
}

func (a *API) Join(svcs ...service.Service) error {
	for _, svc := range svcs {
		if err := svc.Init(); err != nil {
			return err
		}
		if _, err := svc.Health(); err != nil {
			return err
		}
		a.http.Mount(svc.Pattern(), svc.HTTP())
	}
	return nil
}

// Shutdown will gracefully shutdown the env.
// This includes the http server, the services
// and the global dependencies db, koanf.
func (a *API) Shutdown() error {
	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer cancel()
	if err := a.db.Close(); err != nil {
		return err
	}
	if err := a.srv.Shutdown(ctx); err != nil {
		return err
	}
	for _, svc := range a.svcs {
		if err := svc.Shutdown(); err != nil {
			return err
		}
	}
	return nil
}

func (a API) Context() context.Context {
	return a.ctx
}
