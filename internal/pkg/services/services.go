package services

import (
	"github.com/naivary/instance/internal/app/fs"
	"github.com/naivary/instance/internal/app/sys"
)

type Services struct {
	Sys sys.Env
	Fs  fs.Env
}
