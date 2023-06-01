package sys

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/log"
	"github.com/naivary/instance/internal/pkg/models"
	"github.com/naivary/instance/internal/pkg/service"
)

var _ service.Service = (*Sys)(nil)

type Sys struct {
	Svcs []service.Service
	K    *koanf.Koanf
	Meta models.Meta

	logManager log.Manager
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

func (s Sys) Metrics() error {
	return nil
}

func (s *Sys) Init() error {
	mngr, err := log.New(s.K, s)
	if err != nil {
		return err
	}
	s.logManager = mngr
	return nil
}

func (s *Sys) Shutdown() error {
	s.logManager.Shutdown()
	return nil
}
