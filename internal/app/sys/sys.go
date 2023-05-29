package sys

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/naivary/instance/internal/pkg/register"
	"github.com/naivary/instance/internal/pkg/service"
)

var _ service.Service = (*Sys)(nil)

type Sys struct {
	// svcs which has to be checked
	Svcs []service.Service
}

func (s Sys) Metrics() error {
	return nil
}

func (s Sys) health(reg register.Register) (service.Info, error) {
	return service.Info{}, nil
}

func (s Sys) Register(root chi.Router) {
	r := chi.NewRouter()
	r.Get("/health", s.Health)
	root.Mount("/sys", r)
}

func (e Sys) ID() string {
	return uuid.NewString()
}

func (e Sys) Name() string {
	return "sys"
}

func (e Sys) Description() string {
	return "system checks like metrics, health. For the full list of endpoints see the OpenApi definition."
}
