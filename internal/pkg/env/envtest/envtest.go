package envtest

import (
	"github.com/naivary/apitmpl/internal/pkg/config"
	"github.com/naivary/apitmpl/internal/pkg/env"
)

func NewAPI() (*env.API, error) {
	a, err := env.NewAPI(config.DefaultCfgFile)
	if err != nil {
		return nil, err
	}
	if err := a.Config().Set("testing", true); err != nil {
		return nil, err
	}
	if err := a.Init(); err != nil {
		return nil, err
	}
	return a, nil
}
