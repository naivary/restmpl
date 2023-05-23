package fs

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/filestore"
)

type Env struct {
	K *koanf.Koanf

	Store filestore.Filestore
}

var (
	errEmptyFilepath = errors.New("query parameter filepath must be set")
)

func (e Env) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(e.K.Int64("fs.maxSize"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, h, err := r.FormFile(e.K.String("fs.formKey"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	path := filepath.Join(r.FormValue("filepath"), h.Filename)
	_, err = e.Store.Create(path, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (e Env) Remove(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = e.Store.Remove(r.Form.Get("filepath"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (e Env) Read(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Form.Get("filepath") == "" {
		http.Error(w, errEmptyFilepath.Error(), http.StatusBadRequest)
	}
	data, err := e.Store.Read(r.Form.Get("filepath"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
