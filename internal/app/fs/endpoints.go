package fs

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/middleware"
	"github.com/google/jsonapi"
	"github.com/naivary/instance/internal/pkg/filestore"
	"github.com/naivary/instance/internal/pkg/japi"
)

func (f Fs) create(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	err := r.ParseMultipartForm(f.K.Int64("fs.maxSize"))
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}

	file, h, err := r.FormFile(f.K.String("fs.formKey"))
	if err != nil {
		jerr := japi.NewError(err, http.StatusInternalServerError, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}

	path := filepath.Join(r.FormValue("filepath"), h.Filename)
	_, err = f.store.Create(path, file)
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

func (f Fs) remove(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jerr := japi.NewError(err, http.StatusBadRequest, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
		return
	}
	err = f.store.Remove(r.Form.Get("filepath"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jerr := japi.NewError(err, http.StatusBadRequest, reqID)
		jsonapi.MarshalErrors(w, japi.Errors(&jerr))
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
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}
