package env

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/apitmpl/internal/pkg/config"
	"github.com/naivary/apitmpl/internal/pkg/database"
	"github.com/naivary/apitmpl/internal/pkg/logging"
	"github.com/naivary/apitmpl/internal/pkg/logging/builder"
	"github.com/naivary/apitmpl/internal/pkg/metrics"
	"github.com/naivary/apitmpl/internal/pkg/server"
	"github.com/naivary/apitmpl/internal/pkg/service"
	"github.com/pocketbase/dbx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
)

var _ Env = (*API)(nil)

type API struct {
	// global dependencies
	db *dbx.DB
	k  *koanf.Koanf

	// internal
	cfgFile    string
	isInited   bool
	svcs       []service.Service
	http       chi.Router
	metrics    metrics.Manager
	srv        *http.Server
	logManager logging.Manager
	ctx        context.Context
}

// NewAPI creates the an API env provided
// with the given config file. All dependencies
// beside the config of the env will not be
// inited. This can be accomplished using
// the init function.
func NewAPI(cfgFile string) (*API, error) {
	a := &API{
		cfgFile: cfgFile,
		ctx:     context.Background(),
	}
	k, err := config.New(cfgFile)
	if err != nil {
		return nil, err
	}
	a.k = k
	if err := a.k.Set("cfgFile", cfgFile); err != nil {
		return nil, err
	}
	a.logManager = logging.NewEnvManager(os.Stdout)
	return a, nil
}

// Init will initialze the env by setting up
// the remaining dependencies not initiliazed by `NewAPI`
// and run a health check using `Health`.
func (a *API) Init() error {
	if a.isInited {
		return nil
	}

	db, err := database.Connect(a.k)
	if err != nil {
		return err
	}
	a.db = db
	a.metrics = metrics.New()
	a.initHTTP()
	a.isInited = true
	return a.Health()
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

// Serve will start the http server for public
// requests. It also handles the graceful shutdown
// of OS Interrupts signals.
func (a *API) Serve() error {
	if !a.isInited {
		return ErrNotInited
	}
	srv, err := server.New(a.k, a.HTTP())
	if err != nil {
		return err
	}
	go func() {
		if err := srv.ListenAndServeTLS(a.k.String("server.crt"), a.k.String("server.key")); err != nil {
			return
		}
	}()
	b := builder.NewEnvBuilder(a.ctx, slog.LevelInfo, "Successfully started API server!")
	b.APIServerStart(a.k, srv)
	a.logManager.Log(b)

	a.srv = srv
	return nil
}

func (a *API) Join(svcs ...service.Service) error {
	if !a.isInited {
		return ErrNotInited
	}
	for _, svc := range svcs {
		if err := svc.Init(); err != nil {
			return err
		}
		if _, err := svc.Health(); err != nil {
			return err
		}
		if err := a.metrics.Register(svc.Metrics()...); err != nil {
			return err
		}
		a.http.Mount(svc.Pattern(), svc.HTTP())
		rec := builder.NewEnvBuilder(a.ctx, slog.LevelInfo, "Service successfully started!").ServiceInit(svc)
		a.logManager.Log(rec)
	}
	a.svcs = append(a.svcs, svcs...)
	return nil
}

// Shutdown will gracefully shutdown the env.
// This includes the http server, the services
// and the global dependencies db, koanf.
func (a *API) Shutdown() error {
	if !a.isInited {
		return nil
	}
	slog.InfoCtx(a.ctx, "Gracefully shutting down API server")
	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer cancel()
	if err := a.db.Close(); err != nil {
		return err
	}
	if err := a.srv.Shutdown(ctx); err != nil {
		return err
	}
	for _, svc := range a.svcs {
		b := builder.NewEnvBuilder(a.ctx, slog.LevelInfo, "service shutdown").ServiceShutdown(svc)
		a.logManager.Log(b)
		if err := svc.Shutdown(); err != nil {
			return err
		}
	}
	return nil
}

func (a API) Context() context.Context {
	return a.ctx
}

func (a API) Health() error {
	if err := a.db.DB().PingContext(a.ctx); err != nil {
		return err
	}
	if a.k == nil {
		return errors.New("config manager is nil")
	}
	return nil
}

func (a *API) initHTTP() {
	root := chi.NewRouter()
	root.Use(middleware.SetHeader("Content-Type", jsonapi.MediaType))
	root.Use(middleware.RequestID)
	root.Use(middleware.CleanPath)
	root.Use(middleware.Timeout(a.k.Duration("server.timeout.request")))
	root.Mount("/metrics", promhttp.HandlerFor(a.metrics.Registry(), promhttp.HandlerOpts{Registry: a.metrics.Registry()}))
	a.http = root
}
