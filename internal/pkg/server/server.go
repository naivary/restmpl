package server

import (
	"net/http"

	"github.com/knadh/koanf/v2"
)

func New(k *koanf.Koanf, r http.Handler) (*http.Server, error) {
	srv := &http.Server{
		Addr:              k.String("server.addr"),
		ReadHeaderTimeout: k.Duration("server.timeout.readHeader"),
		WriteTimeout:      k.Duration("server.timeout.write"),
		IdleTimeout:       k.Duration("server.timeout.idle"),
		MaxHeaderBytes:    k.Int("server.maxHeaderBytes"),
		Handler:           r,
	}
	return srv, nil
}
