package filestore

import (
	"net/http"

	"github.com/knadh/koanf/v2"
)

type Filestore struct {
	// base represents the base directory
	// where to store the files.
	Base string

	Handler http.Handler
}

func New(k *koanf.Koanf) Filestore {
	return Filestore{
		Base: k.String("fs.base"),
		Handler:    http.FileServer(http.Dir(k.String("fs.base"))),
	}
}
