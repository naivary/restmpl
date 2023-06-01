package env

import "errors"

var (
	ErrNotInited = errors.New("env is not inited. Call `init` first")
)
