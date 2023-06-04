package env

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/apitmpl/internal/pkg/config"
	"github.com/naivary/apitmpl/internal/pkg/database"
	"github.com/naivary/apitmpl/internal/pkg/jwtauth"
	"github.com/naivary/apitmpl/internal/pkg/logging"
	"github.com/naivary/apitmpl/internal/pkg/logging/builder"
	"github.com/naivary/apitmpl/internal/pkg/metrics"
	"github.com/naivary/apitmpl/internal/pkg/server"
	"github.com/naivary/apitmpl/internal/pkg/service"
	"github.com/pocketbase/dbx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"golang.org/x/exp/slog"
)

var _ Env = (*API)(nil)

type API struct {
	// global dependencies
	db   *dbx.DB
	k    *koanf.Koanf
	meta meta

	// internal
	cfgFile  string
	isInited bool
	svcs     []service.Service
	http     chi.Router
	metric   metrics.Manager
	srv      *http.Server
	logger   logging.Manager
	ctx      context.Context
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
	if err := a.initConfig(); err != nil {
		return nil, err
	}
	a.logger = logging.NewEnvManager(os.Stdout)
	return a, nil
}

func (a *API) initConfig() error {
	k, err := config.New(a.cfgFile)
	if err != nil {
		return err
	}
	a.k = k
	if err := a.k.Set("cfgFile", a.cfgFile); err != nil {
		return err
	}

	secret := os.Getenv("API_JWT_SECRET")
	if secret == "" {
		return errors.New("API_JWT_SECRET env variable not set")
	}
	jwtauth.SetSecret(secret)
	return nil
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
	a.metric = metrics.NewManager()
	a.meta = a.newMeta()
	a.initHTTP()
	// registrining env metrics
	if err := a.metric.Register(a.Metrics()...); err != nil {
		return err
	}
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
	if len(a.svcs) == 0 {
		return ErrNoServices
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
	b := builder.NewEnvBuilder(a.ctx, slog.LevelInfo, "Successfully started API server!").APIServerStart(a.k, srv)
	a.logger.Log(b)
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
		if err := a.metric.Register(svc.Metrics()...); err != nil {
			return err
		}
		a.http.Mount(svc.Pattern(), svc.HTTP())
		rec := builder.NewEnvBuilder(a.ctx, slog.LevelInfo, "Service successfully started!").ServiceInit(svc)
		a.logger.Log(rec)
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
		b := builder.NewEnvBuilder(a.ctx, slog.LevelInfo, "service shutdown").ServiceInfo(svc)
		a.logger.Log(b)
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
	if a.metric == nil {
		return errors.New("missing metric manager")
	}
	if a.logger == nil {
		return errors.New("missing logger manager")
	}
	return nil
}

func (a *API) initHTTP() {
	root := chi.NewRouter()
	root.Use(middleware.SetHeader("Content-Type", "application/json"))
	root.Use(middleware.RequestID)
	root.Use(middleware.CleanPath)
	root.Use(middleware.Timeout(a.k.Duration("server.timeout.request")))
	for _, mw := range a.middlewares() {
		root.Use(mw)
	}
	root.Mount("/sys", a.initMonitorHTTP())
	a.http = root
}

func (a *API) initMonitorHTTP() chi.Router {
	r := chi.NewRouter()
	r.Use(jwtauth.Verify)
	r.Mount("/metrics", a.metrics())
	r.Get("/health", a.health)
	return r
}

func (a *API) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		collectors.NewGoCollector(),
		collectors.NewDBStatsCollector(a.db.DB(), "main_db"),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	}
}
