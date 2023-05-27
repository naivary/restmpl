package configtest

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/fs"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/configs"
)

func New() (*koanf.Koanf, error) {
	k := koanf.New(".")
	err := k.Load(fs.Provider(configs.Fs, "instance.yaml"), yaml.Parser())
	if err != nil {
		return nil, err
	}
	return k, nil
}
