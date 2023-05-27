package fs

import (
	"net/http"

	"github.com/naivary/instance/internal/pkg/log"
	"golang.org/x/exp/slog"
)

func (f Fs) info(r *http.Request) {
	wrs := log.NewWriters(f.K)
	file, _ := wrs.Get(slog.LevelInfo)
	slog.New(slog.NewTextHandler(file, nil)).InfoCtx(r.Context(),
		"message",
		slog.Group("request",
			slog.String("method", r.Method),
			slog.String("host", r.Host),
			slog.String("endpoint", r.URL.Path),
			slog.String("remoteAddr", r.RemoteAddr),
		),
		slog.Group("service",
			slog.String("name", f.Name()),
			slog.String("id", f.UUID()),
		),
	)
}
