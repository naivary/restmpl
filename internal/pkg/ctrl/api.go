package ctrl

import (
	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/service"
)

type API struct {
	Services []service.Service
	Router   chi.Router
}
