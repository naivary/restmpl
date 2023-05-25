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

<<<<<<< HEAD
=======
var _ service.Service[chi.Router] = (*Sys)(nil)

>>>>>>> svc-interface
type Sys struct {
	K  *koanf.Koanf
	DB *sql.DB

	M metadata.Metadata
}

<<<<<<< HEAD
func (s *Sys) Health(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	err := s.DB.Ping()
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}
	s.M.DBRunning = err == nil

	w.Header().Add("Content-Type", jsonapi.MediaType)
	err = jsonapi.MarshalPayload(w, &s.M)
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}

	w.WriteHeader(http.StatusOK)
=======
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
>>>>>>> svc-interface
}
