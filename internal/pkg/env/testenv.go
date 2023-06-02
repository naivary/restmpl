package env

import "github.com/naivary/apitmpl/internal/pkg/config"

func NewTestAPI() (*API, error) {
	a, err := NewAPI(config.DefaultCfgFile)
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
