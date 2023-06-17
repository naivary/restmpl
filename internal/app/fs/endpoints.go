package fs

import (
	"fmt"
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
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	claims, err := jwtauth.GetClaims(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	owner, err := claims.GetSubject()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	obj := objst.NewObject(h.Filename, owner)
	if _, err := obj.ReadFrom(file); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !r.Form.Has(objst.ContentType) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	obj.SetMeta(objst.ContentType, r.Form.Get(objst.ContentType))
	if err := f.b.Create(obj); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
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
}

func (f Fs) read(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	obj, err := f.b.GetByName(r.Form.Get("name"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ct, _ := obj.GetMeta(objst.ContentType)
	w.Header().Set("Content-Type", ct)
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, obj); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
