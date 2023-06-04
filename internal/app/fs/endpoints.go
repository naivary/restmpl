package fs

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/naivary/apitmpl/internal/pkg/filestore"
)

func (f Fs) create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(f.K.Int64("fs.maxSize")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	file, h, err := r.FormFile(f.K.String("fs.formKey"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	path := filepath.Join(r.FormValue("filepath"), h.Filename)
	_, err = f.store.Create(path, file)
	if errors.Is(err, filestore.ErrWrongNaming) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (f Fs) remove(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = f.store.Remove(r.Form.Get("filepath"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (f Fs) read(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := f.store.Read(r.Form.Get("filepath"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
