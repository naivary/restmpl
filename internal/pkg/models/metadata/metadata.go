package metadata

import "github.com/knadh/koanf/v2"

func New(k *koanf.Koanf) Metadata {
	return Metadata{
		Version: k.String("version"),
	}
}

type Metadata struct {
	Version   string `json:"version"`
	DBRunning bool   `json:"dbRunning"`
}
