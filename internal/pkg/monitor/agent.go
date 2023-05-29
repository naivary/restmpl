package monitor

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/dependency"
	"github.com/naivary/instance/internal/pkg/service"
)

type Agent interface {
	// TODO(naivary): health should provide info about the services checed
	Health() error
	Metrics() error

	// HTTP access to the
	// agent's services
	HTTP() chi.Router

	// services which the agent
	// is checking
	Services() []service.Service
}

var _ Agent = (*agent)(nil)

type agent struct {
	// svcs which will be checked
	svcs []service.Service

	deps []dependency.Pinger
}

func New(svcs []service.Service, deps []dependency.Pinger) agent {
	return agent{
		svcs: svcs,
		deps: deps,
	}
}

func (a agent) Services() []service.Service {
	return a.svcs
}

func (a agent) Metrics() error {
	return nil
}

func (a agent) Health() error {
	for _, dep := range a.deps {
		err := dep.Ping()
		if err != nil {
			return err
		}
	}
	return nil
}
