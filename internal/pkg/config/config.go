package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	keyPathDelimiter = "."
)

func New(path string) (*koanf.Koanf, error) {
	k := koanf.New(keyPathDelimiter)
	err := k.Load(file.Provider(path), yaml.Parser())
	return k, err
}
