package monitor

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/register"
	"github.com/naivary/instance/internal/pkg/service"
)

type Agent interface {
	Health() error
	Metrics() error

	// HTTP access to the
	// agent's services
	HTTP() chi.Router

	// Services up on which
	// monitoring checks will be apllied
	// TODO(naivary): return information
	// about the services not the services on their own
	Services() []service.Service
}

var _ Agent = (*agent)(nil)

type agent struct {
	// svcs which has to be checked
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
