package filestore

import (
	"net/http"
	"os"

	"github.com/knadh/koanf/v2"
)

type Filestore struct {
	// base represents the base directory
	// where to store the files.
	Base string

	Store http.Handler
}

func New(k *koanf.Koanf) (Filestore, error) {
	base := k.String("fs.base")
	dir := http.Dir(base)
	h := http.FileServer(dir)
	fs := Filestore{
		Base:  base,
		Store: http.StripPrefix("/fs", h),
	}

	err := os.MkdirAll(fs.Base, os.ModePerm)
	if err != nil {
		return Filestore{}, err
	}
	return fs, nil
}
