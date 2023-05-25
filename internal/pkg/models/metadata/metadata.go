package metadata

import (
	"github.com/knadh/koanf/v2"
)

func New(k *koanf.Koanf) Metadata {
	return Metadata{
		// ID is static so it will be uniquely identified
		// on every deployment version.
		ID:      "26e4a9ae-67e4-430f-9263-de9a18d6160b",
		Version: k.String("version"),
	}
}

type Metadata struct {
	ID string `jsonapi:"primary,metadata"`

	Version   string `jsonapi:"attr,version"`
	DBRunning bool   `jsonapi:"attr,dbRunning"`
}
