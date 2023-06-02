package main

import (
	"os"

	"golang.org/x/exp/slog"
)

func main() {
	if err := run(); err != nil {
		slog.Error("something went wrong while starting the server", slog.String("err", err.Error()))
	}
}

func getCfgFile() string {
	if len(os.Args) < 2 {
		return ""
	}
	return os.Args[1]
}

func run() error {
	cfgFile := getCfgFile()
	if cfgFile == "" {
		slog.Info("no custom config file provided. Using default.yaml")
	}
	e, err := newEnv(cfgFile)
	if err != nil {
		return err
	}

	go func() {
		if err := e.Serve(); err != nil {
			slog.ErrorCtx(e.Context(), "error serving the env", "err", err.Error())
			return
		}
	}()

	return initGracefulShutdown(e)
}
