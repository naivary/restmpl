package fs

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/filestore"
)

type Env struct {
	Fs filestore.Filestore
	K  *koanf.Koanf
}

func (e Env) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(e.K.Int64("fs.maxSize"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	src, h, err := r.FormFile(e.K.String("fs.formKey"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dest, err := os.Create(filepath.Join(e.Fs.Base, h.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	res := struct {
		Filename string `json:"filename"`
		Endpoint string `json:"endpoint"`
	}{
		Filename: h.Filename,
		Endpoint: filepath.Join("fs", h.Filename),
	}

	err = json.NewEncoder(w).Encode(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (e Env) Delete(w http.ResponseWriter, r *http.Request) {}
