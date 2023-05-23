package metadata

import (
	"github.com/knadh/koanf/v2"
)

func New(k *koanf.Koanf) Metadata {
	return Metadata{
		ID:      "26e4a9ae-67e4-430f-9263-de9a18d6160b",
		Version: k.String("version"),
	}
}

type Metadata struct {
	// Ressource ID so it is jsonapi compatible
	ID string `jsonapi:"primary,metadata"`

	Version   string `jsonapi:"attr,metadata"`
	DBRunning bool   `jsonapi:"attr,dbRunning"`
}
