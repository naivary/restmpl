package sys

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/models/metadata"
	"github.com/naivary/instance/internal/pkg/service"
)

var _ service.Service = (*Env)(nil)

type Env struct {
	K  *koanf.Koanf
	DB *sql.DB

	M metadata.Metadata
}

func (e Env) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/health", e.Health)
	return r
}

func (e Env) Routername() string {
	return "/sys"
}

func (e Env) UUID() string {
	return uuid.NewString()
}

func (e Env) Name() string {
	return "sys"
}

func (e Env) Description() string {
	return "system checks like metrics, health. For the full list of endpoints see the OpenApi definition."
}
