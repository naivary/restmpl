package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func New(addr string, r chi.Router) (*http.Server, error) {
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
	}

	return srv, nil
}
