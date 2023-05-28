package sys

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/models/metadata"
	"github.com/naivary/instance/internal/pkg/service"
	"github.com/pocketbase/dbx"
)

var _ service.Service[chi.Router] = (*Sys)(nil)

type Sys struct {
	K  *koanf.Koanf
	DB *dbx.DB

	M metadata.Metadata
}

func (s Sys) RegisterRootMiddleware(root chi.Router) {}

func (s Sys) Register(root chi.Router) {
	r := chi.NewRouter()
	r.Get("/health", s.Health)
	root.Mount("/sys", r)
	fmt.Println("djaksldjaskdjsa")
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
