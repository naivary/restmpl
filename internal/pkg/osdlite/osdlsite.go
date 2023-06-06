package osdlite

import (
	"github.com/pocketbase/dbx"
	_ "modernc.org/sqlite"
)

type OSDLite struct {
	store *dbx.DB
}

func New() (*OSDLite, error) {
	store, err := dbx.Open("sqlite", "file:osdlite.db")
	if err != nil {
		return nil, err
	}
	o := &OSDLite{
		store: store,
	}
	o.initOptions()
	if err := o.initPragmas(); err != nil {
		return nil, err
	}
	if err := o.initSchema(); err != nil {
		return nil, err
	}
	return o, nil
}

func (o OSDLite) CreateObj(obj *object) error {
	return o.store.Model(obj).Insert()
}

func (o OSDLite) CreateBucket(b *bucket) error {
	return o.store.Model(b).Insert()
}
