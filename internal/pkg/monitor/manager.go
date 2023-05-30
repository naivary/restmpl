package monitor

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/service"
)

type Manager interface {
	Health() error
	Metrics() error

	// HTTP access to the
	// manager's services
	HTTP() chi.Router

	// services which the manager
	// is checking
	Services() []service.Service

	Join(...service.Service)
}

var _ Manager = (*manager)(nil)

type manager struct {
	// svcs which will be checked
	svcs []service.Service
}

func New(svcs []service.Service) *manager {
	return &manager{
		svcs: svcs,
	}
}

func (m manager) Services() []service.Service {
	return m.svcs
}

func (m manager) Metrics() error {
	return nil
}

func (m manager) Health() error {
	return nil
}

func (m *manager) Join(svcs ...service.Service) {
	m.svcs = append(m.svcs, svcs...)
}
