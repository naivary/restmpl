package monitor

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (a agent) HTTP() chi.Router {
	r := chi.NewRouter()
	r.Get("/health", a.health)
	r.Get("/metrics", a.metrics)
	r.Get("/services", a.services)
	return r
}

func (a agent) health(w http.ResponseWriter, r *http.Request) {}

func (a agent) metrics(w http.ResponseWriter, r *http.Request) {}

func (a agent) services(w http.ResponseWriter, r *http.Request) {}
