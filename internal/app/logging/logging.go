package logging

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/service"
	"golang.org/x/exp/slog"
)

var _ service.Service[chi.Router] = (*Logging)(nil)

type Logging struct {
	K    *koanf.Koanf
	Info *slog.Logger
}

func (l Logging) UUID() string {
	return uuid.NewString()
}

func (l Logging) Name() string {
	return "logger"
}

func (l Logging) Description() string {
	return "logger provides middleware for standard structured logging using slog"
}

func (l Logging) middlewares() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		l.Logger,
	}
}

func (l Logging) RegisterRootMiddleware(root chi.Router) {
	for _, mw := range l.middlewares() {
		root.Use(mw)
	}
}

func (l Logging) Register(root chi.Router) {}
