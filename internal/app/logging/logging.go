package logging

import (
	"net/http"

	"github.com/knadh/koanf/v2"
	"golang.org/x/exp/slog"
)

type Logging struct {
	K    *koanf.Koanf
	Info slog.Handler
}

func (l Logging) defaultLogger(w http.ResponseWriter, r *http.Request) {}

// Logger is the default logger for any http request
func (l Logging) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(l.defaultLogger)
}
