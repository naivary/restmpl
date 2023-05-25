package sys

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/models/metadata"
	"github.com/naivary/instance/internal/pkg/service"
)

var _ service.Service[chi.Router] = (*Sys)(nil)

type Sys struct {
	K  *koanf.Koanf
	DB *sql.DB

	M metadata.Metadata
}

func (s Sys) Register(root chi.Router) {
	root.Mount(s.pattern(), s.router())
}

func (s Sys) router() http.Handler {
	r := chi.NewRouter()
	r.Get("/health", s.Health)
	return r
}

func (s Sys) pattern() string {
	return "/sys"
}

func (e Sys) UUID() string {
	return uuid.NewString()
}

func (e Sys) Name() string {
	return "sys"
}

func (e Sys) Description() string {
	return "system checks like metrics, health. For the full list of endpoints see the OpenApi definition."
}
