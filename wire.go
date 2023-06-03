package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/naivary/apitmpl/internal/app/auth"
	"github.com/naivary/apitmpl/internal/app/fs"
	"github.com/naivary/apitmpl/internal/app/users"
	"github.com/naivary/apitmpl/internal/pkg/env"
	"github.com/naivary/apitmpl/internal/pkg/service"
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
	f := new(fs.Fs)
	f.K = env.Config()

	a := new(auth.Auth)
	a.K = env.Config()
	a.DB = env.DB()

	u := new(users.Users)
	u.DB = env.DB()
	u.K = env.Config()

	svcs = append(svcs, f, a, u)
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
