package fs

import (
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/filestore"
)

type Env struct {
	K *koanf.Koanf

	Store filestore.Store
}
