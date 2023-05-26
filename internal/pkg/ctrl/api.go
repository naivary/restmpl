package ctrl

import (
	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/service"
)

type API struct {
	Services []service.Service[chi.Router]
	Router   chi.Router
	K        *koanf.Koanf
}
