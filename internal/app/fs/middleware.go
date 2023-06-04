package fs

import (
	"net/http"
	"strings"

	"github.com/naivary/apitmpl/internal/pkg/logging/builder"
	"golang.org/x/exp/slog"
)

func (f Fs) Middlewares() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		f.info,
		f.countReq,
		f.forceFilepath,
	}
}

func (f Fs) countReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := f.m.GetCounter("req")
		c.Inc()
		next.ServeHTTP(w, r)
	})
}

func (f Fs) info(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := builder.NewSvcBuilder(r.Context(), slog.LevelInfo, "handling request").IncomingRequest(r)
		f.l.Log(rec)
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
	})
}
