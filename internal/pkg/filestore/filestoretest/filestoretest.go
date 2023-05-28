package filestoretest

import (
	"os"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/filestore"
	"github.com/spf13/afero"
)

// NewTestStore provides a in memory fs which
// caching. The caching has no benefits but
// simulates the "real" filestore.
func New(k *koanf.Koanf) (filestore.Store[afero.File], error) {
	base := k.String("fs.basepath")
	err := os.MkdirAll(base, os.ModePerm)
	if err != nil {
		return nil, err
	}
	firstLayer := afero.NewBasePathFs(afero.NewMemMapFs(), base)
	secLayer := afero.NewMemMapFs()
	return &filestore.Filestore{
		Store: afero.Afero{
			Fs: afero.NewCacheOnReadFs(firstLayer, secLayer, k.Duration("fs.ttl")),
		},
	}, nil

}
