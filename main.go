package main

import (
	"errors"
	"os"

	"golang.org/x/exp/slog"
)

func main() {
	if err := run(); err != nil {
		slog.Error("something went wrong while starting the sevrer", slog.String("err", err.Error()))
	}
}

func getCfgFile() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("missing config file as the first argument")
	}
	return os.Args[1], nil
}

func run() error {
	cfgFile, err := getCfgFile()
	if err != nil {
		return err
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
