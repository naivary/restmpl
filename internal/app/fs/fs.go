package fs

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/middleware"
	"github.com/google/jsonapi"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/filestore"
	"github.com/naivary/instance/internal/pkg/japi"
)

type Env struct {
	K *koanf.Koanf

	Store filestore.Store
}

func (e Env) Create(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	err := r.ParseMultipartForm(e.K.Int64("fs.maxSize"))
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}

	file, h, err := r.FormFile(e.K.String("fs.formKey"))
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}

	path := filepath.Join(r.FormValue("filepath"), h.Filename)
	_, err = e.Store.Create(path, file)
	if errors.Is(err, &filestore.ErrWrongNaming) {
		w.WriteHeader(http.StatusBadRequest)
		jsonapi.MarshalErrors(w, japi.Errors(&filestore.ErrWrongNaming))
		return
	}
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (e Env) Remove(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	err := r.ParseForm()
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}

	err = e.Store.Remove(r.Form.Get("filepath"))
	if err != nil {
		jerr := japi.NewError(err, http.StatusBadRequest, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
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

	data, err := e.Store.Read(r.Form.Get("filepath"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
