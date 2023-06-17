package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/naivary/restmpl/internal/app/auth"
	"github.com/naivary/restmpl/internal/app/fs"
	"github.com/naivary/restmpl/internal/app/users"
	"github.com/naivary/restmpl/internal/pkg/env"
	"github.com/naivary/restmpl/internal/pkg/service"
	"golang.org/x/exp/slog"
)

func newEnv(cfgFile string) (env.Env, error) {
	api, err := env.NewAPI(cfgFile)
	if err != nil {
		return nil, err
	}
	if err := api.Init(); err != nil {
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

	a := new(auth.Auth)
	a.K = env.Config()
	a.DB = env.DB()

	u := new(users.Users)
	u.DB = env.DB()
	u.K = env.Config()

	f := new(fs.Fs)
	f.K = env.Config()
	f.DB = env.DB()

	svcs = append(svcs, a, u, f)
	return svcs, nil
}

func initGracefulShutdown(e env.Env) error {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	if err := e.Shutdown(); err != nil {
		return err
	}
	slog.InfoCtx(e.Context(), "Graceful shutdown succeeded! Exiting.")
	return nil
}
