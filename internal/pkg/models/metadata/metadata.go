package metadata

import (
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
)

func New(k *koanf.Koanf) Metadata {
	return Metadata{
		ID:      uuid.NewString(),
		Version: k.String("version"),
	}
}

type Metadata struct {
	// Ressource ID so it is jsonapi compatible
	ID string `jsonapi:"primary,metadata"`

	Version   string `jsonapi:"attr,metadata"`
	DBRunning bool   `jsonapi:"attr,dbRunning"`
}
