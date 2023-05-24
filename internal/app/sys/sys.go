package sys

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/models/metadata"
)

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
	return "system controller to check the status of the application using health checks, getting metrics etc."
}
