package fs

import (
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/env"
	"github.com/naivary/instance/internal/pkg/filestore"
)

var _ env.Env = (*Env)(nil)

type Env struct {
	K *koanf.Koanf

	Store filestore.Store
}

func (e Env) Name() string {
	return "fs"
}

func (e Env) ID() string {
	return uuid.NewString()
}
