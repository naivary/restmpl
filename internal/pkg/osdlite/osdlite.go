package osdlite

import (
	"errors"

	"github.com/pocketbase/dbx"
	_ "modernc.org/sqlite"
)

type OSD interface {
	Health() error
	Create(object) error
	Read(string) error
	Remove(string) error
}

var _ OSD = (*OSDLite)(nil)

type OSDLite struct {
	fs *dbx.DB
}

func New() (*OSDLite, error) {
	fs, err := dbx.Open("sqlite", "file:fs.db")
	if err != nil {
		return nil, err
	}
	o := &OSDLite{
		fs: fs,
	}

	if err := o.initPragmas(); err != nil {
		return nil, err
	}
	o.initOptions()
	if err := o.initSchema(); err != nil {
		return nil, err
	}
	return o, o.Health()
}

func (o OSDLite) Create(obj object) error {
	return o.fs.Model(&obj).Insert()
}

func (o OSDLite) Health() error {
	if o.fs == nil {
		return errors.New("object storage is nil")
	}
	return o.fs.DB().Ping()
}

func (o OSDLite) Read(path string) error {
	return nil
}

func (o OSDLite) Remove(path string) error {
	return nil
}
