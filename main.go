package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

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

	// graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	slog.InfoCtx(e.Context(), "Gracefully shutting down env and services")
	if err := e.Shutdown(); err != nil {
		return err
	}
	slog.InfoCtx(e.Context(), "Services and env shutdown. Exiting.")
	return nil
}
