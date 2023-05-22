package fs

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Env struct{}

func (e Env) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, h, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// TODO(naivary): associate with the file sever which is accessible
	// using the /fs endpoint
	f, err := os.Create(filepath.Join(home, h.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer f.Close()

	_, err = io.Copy(f, file)
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
