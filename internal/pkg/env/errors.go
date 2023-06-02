package env

import "errors"

var (
	ErrNotInited  = errors.New("env is not inited. Call `init` first")
	ErrNoServices = errors.New("no services found in the env. Use join to add new services")
)
