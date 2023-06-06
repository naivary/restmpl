package osdlite

import (
	"errors"

	"github.com/pocketbase/dbx"
	_ "modernc.org/sqlite"
)

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

func (o OSDLite) CreateObject(obj *object) error {
	if len(obj.Payload.Bytes()) <= 0 {
		return errors.New("object has not payload")
	}
	return o.fs.Model(obj).Insert()
}

func (o OSDLite) CreatBucket(b *bucket) error {
	return o.fs.Model(b).Insert()
}

func (o OSDLite) RemoveBucket(b *bucket) error {
	return o.fs.Model(b).Delete()
}

func (o OSDLite) RemoveObject(obj *object) error {
	return o.fs.Model(obj).Delete()
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
