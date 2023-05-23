package filestore

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/knadh/koanf/v2"
	"github.com/spf13/afero"
)

type Filestore struct {
	Basepath string

	Store afero.Afero
}

var (
	ErrWrongNaming = errors.New("file name must follow the following regex pattern: [a-z._0-9-]+")
)

func New(k *koanf.Koanf) (Filestore, error) {
	base := k.String("fs.base")
	err := os.MkdirAll(base, os.ModePerm)
	if err != nil {
		return Filestore{}, err
	}
	osFs := afero.NewBasePathFs(afero.NewOsFs(), base)
	memFs := afero.NewMemMapFs()

	return Filestore{
		Store: afero.Afero{
			Fs: afero.NewCacheOnReadFs(osFs, memFs, k.Duration("fs.ttl")),
		},
	}, nil
}

func (f Filestore) followsNamingConvention(name string) bool {
	r, _ := regexp.Compile("^[a-z._0-9-]+$")
	return r.MatchString(name)
}

func (f Filestore) Create(path string, r io.Reader) (afero.File, error) {
	// assure that it is following a proper naming convention.
	if !f.followsNamingConvention(filepath.Base(path)) {
		return nil, ErrWrongNaming
	}

	// dont create the file, if exiting
	isExisting, err := f.Store.Exists(path)
	if err != nil {
		return nil, err
	}
	if isExisting {
		return nil, os.ErrExist
	}

	err = f.Store.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return nil, err
	}
	file, err := f.Store.Create(path)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(file, r)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f Filestore) Remove(path string) error {
	return f.Store.Remove(path)
}

func (f Filestore) Read(path string) ([]byte, error) {
	return f.Store.ReadFile(path)
}