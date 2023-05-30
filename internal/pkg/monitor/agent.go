package monitor

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/service"
)

type Agent interface {
	Health() error
	Metrics() error

	// HTTP access to the
	// agent's services
	HTTP() chi.Router

	// services which the agent
	// is checking
	Services() []service.Service

	Join(...service.Service)
}

var _ Agent = (*agent)(nil)

type agent struct {
	// svcs which will be checked
	svcs []service.Service
}

func New(svcs []service.Service) *agent {
	return &agent{
		svcs: svcs,
	}
}

func (a agent) Services() []service.Service {
	return a.svcs
}

func (a agent) Metrics() error {
	return nil
}

func (a agent) Health() error {
	return nil
}

func (a *agent) Join(svcs ...service.Service) {
	a.svcs = append(a.svcs, svcs...)
}
