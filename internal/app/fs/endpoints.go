package fs

import (
	"errors"
	"io"
	"net/http"

	"github.com/naivary/objst"
	"github.com/naivary/restmpl/internal/pkg/jwtauth"
)

func (f Fs) create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(f.maxSize); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	file, h, err := r.FormFile(f.formKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	claims, err := jwtauth.GetClaims(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	owner, err := claims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	obj := objst.NewObject(h.Filename, owner)
	if _, err := obj.ReadFrom(file); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !r.Form.Has(objst.ContentType) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	obj.SetMeta(objst.ContentType, r.Form.Get(objst.ContentType))
	if err := f.b.Create(obj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (f Fs) remove(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !r.Form.Has("name") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	claims, err := jwtauth.GetClaims(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	owner, err := claims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = f.b.IsAuthorizedByName(owner, r.Form.Get("name"))
	if errors.Is(err, objst.ErrUnauthorized) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := f.b.DeleteByName(r.Form.Get("name")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
	if !r.Form.Has("name") {
		http.Error(w, "paramter name must be set", http.StatusBadRequest)
		return
	}
	claims, err := jwtauth.GetClaims(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	owner, err := claims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	obj, err := f.b.IsAuthorizedByName(owner, r.Form.Get("name"))
	if errors.Is(err, objst.ErrUnauthorized) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && !errors.Is(err, objst.ErrUnauthorized) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ct, _ := obj.GetMeta(objst.ContentType)
	w.Header().Set("Content-Type", ct)
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
