package sys

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/apitmpl/internal/pkg/logging"
	"github.com/naivary/apitmpl/internal/pkg/metrics"
	"github.com/naivary/apitmpl/internal/pkg/models"
	"github.com/naivary/apitmpl/internal/pkg/service"
	"github.com/prometheus/client_golang/prometheus"
)

var _ service.Service = (*Sys)(nil)

type Sys struct {
	Svcs []service.Service
	K    *koanf.Koanf
	Meta models.Meta

	metManager metrics.Manager
	logManager logging.Manager
}

func (s Sys) ID() string {
	return uuid.NewString()
}

func (s Sys) Name() string {
	return "sys"
}

func (s Sys) Pattern() string {
	return "/sys"
}

func (s Sys) HTTP() chi.Router {
	r := chi.NewRouter()
	for _, mw := range s.Middlewares() {
		r.Use(mw)
	}
	r.Get("/health", s.health)
	r.Get("/metrics", s.metrics)
	return r
}

func (s Sys) Description() string {
	return "status and information about the api"
}

func (s Sys) Health() (*service.Info, error) {
	if s.K == nil {
		return nil, errors.New("missing config manager")
	}
	if len(s.Svcs) == 0 {
		return nil, errors.New("no services provided to check for health")
	}
	return &service.Info{
		ID:   s.ID(),
		Name: s.Name(),
		Desc: s.Description(),
	}, nil
}

func (s Sys) Metrics() []prometheus.Collector {
	s.metManager.AddCounter("requestCounter", metrics.IncomingHTTPRequest(&s))
	s.metManager.AddCounter("errorCounter", metrics.NumberOfErrors(&s))
	return s.metManager.All()
}

func (s *Sys) Init() error {
	mngr, err := logging.NewSvcManager(s.K, s)
	if err != nil {
		return err
	}
	s.logManager = mngr
	s.metManager = metrics.New()
	return nil
}

func (s *Sys) Shutdown() error {
	s.logManager.Shutdown()
	return nil
}
