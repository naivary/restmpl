package fs

import (
	"net/http"

	"github.com/google/jsonapi"
	"github.com/naivary/instance/internal/pkg/japi"
)

func (e Env) ForceFilepath(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") == "multipart/form-data" {
			err := r.ParseMultipartForm(e.K.Int64("fs.maxSize"))
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
