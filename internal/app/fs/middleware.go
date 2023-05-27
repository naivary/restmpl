package fs

import (
	"net/http"
	"strings"

	"github.com/google/jsonapi"
	"github.com/naivary/instance/internal/pkg/japi"
	"github.com/naivary/instance/internal/pkg/log"
	"golang.org/x/exp/slog"
)

func (f Fs) Middlewares() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		f.forceFilepath,
		f.infoLog,
	}
}

func (f Fs) infoLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		next.ServeHTTP(w, r)
	})
}

func (f Fs) forceFilepath(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			err := r.ParseMultipartForm(f.K.Int64("fs.maxSize"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if r.Form.Get("filepath") != "" {
				next.ServeHTTP(w, r)
				return
			}
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Form.Get("filepath") != "" {
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		jsonapi.MarshalErrors(w, japi.Errors(&errEmptyFilepath))
	})
}
