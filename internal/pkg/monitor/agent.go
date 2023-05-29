package monitor

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/register"
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
}

func New(svcs []service.Service) agent {
	return agent{
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
	reg := register.New()
	for _, svc := range a.svcs {
		_, err := svc.Health(reg)
		if err != nil {
			return err
		}
	}
	return nil
}
