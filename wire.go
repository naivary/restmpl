package main

import (
	"github.com/naivary/instance/internal/app/fs"
	"github.com/naivary/instance/internal/app/sys"
	"github.com/naivary/instance/internal/pkg/env"
	"github.com/naivary/instance/internal/pkg/models"
	"github.com/naivary/instance/internal/pkg/service"
)

func newEnv(cfgFile string) (env.Env, error) {
	api, err := env.NewAPI(cfgFile)
	if err != nil {
		return nil, err
	}
	svcs, err := createServices(api)
	if err != nil {
		return nil, err
	}
	if err := api.Join(svcs...); err != nil {
		return nil, err
	}
	return api, nil
}

func createServices(env *env.API) ([]service.Service, error) {
	svcs := make([]service.Service, 0)
	k := env.Config()
	db := env.DB()
	// services
	f := new(fs.Fs)
	s := new(sys.Sys)

	f.K = k

	svcs = append(svcs, f, s)
	s.Svcs = svcs
	s.K = k
	m := models.NewMeta(k, db, env)
	s.Meta = m

	return svcs, nil
}
