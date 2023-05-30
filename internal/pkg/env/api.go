package env

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/config"
	"github.com/naivary/instance/internal/pkg/database"
	"github.com/naivary/instance/internal/pkg/log"
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
	svcs       []service.Service
	k          *koanf.Koanf
	http       chi.Router
	monAgent   monitor.Agent
	logManager log.Manager
	db         *dbx.DB
	cfgFile    string
	isInited   bool
}

func NewAPI(cfgFile string, svcs []service.Service) API {
	return API{
		cfgFile: cfgFile,
		svcs:    svcs,
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

func (a *API) initHTTP() {
	root := chi.NewRouter()
	root.Use(middleware.SetHeader("Content-Type", jsonapi.MediaType))
	root.Use(middleware.RequestID)
	root.Use(middleware.CleanPath)
	root.Use(middleware.Timeout(reqTimeout))
	root.Mount("/sys", a.Monitor().HTTP())
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

func (a API) Serve() error {
	srv, err := server.New(a.k, a.HTTP())
	if err != nil {
		return err
	}
	return srv.ListenAndServeTLS(a.k.String("server.crt"), a.k.String("server.key"))
}

func (a *API) Init() error {
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

	a.monAgent = monitor.New(a.svcs)

	a.initHTTP()
	a.isInited = true
	return nil
}

func (a *API) Join(svcs ...service.Service) error {
	err := a.Init()
	if err != nil {
		return err
	}
	for _, svc := range svcs {
		err = svc.Init(a.k, a.db)
		if err != nil {
			return err
		}
		a.http.Mount(svc.Pattern(), svc.HTTP())
	}
	return nil
}
