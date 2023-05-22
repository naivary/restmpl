package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/fs"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/configs"
)

const (
	keyPathDelimiter = "."
)

func New() (*koanf.Koanf, error) {
	k := koanf.New(keyPathDelimiter)
	err := k.Load(fs.Provider(configs.Fs, "instance.yaml"), yaml.Parser())
	return k, err
}

