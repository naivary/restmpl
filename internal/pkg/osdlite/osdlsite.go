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

func (o OSDLite) CreateObject(obj *object) error {
	return o.store.Model(obj).Insert()
}

func (o OSDLite) CreateBucket(b *bucket) error {
	return o.store.Model(b).Insert()
}

func (o OSDLite) GetObject(bucketID, objectID string) (*object, error) {
	obj := object{}
	q := o.store.Select().From("objects").Where(dbx.HashExp{
		"bucket_id": bucketID,
		"id":        objectID,
	})
	if err := q.One(&obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (o OSDLite) GetBucket(bucketID string) (*bucket, error) {
	b := bucket{}
	q := o.store.Select().From("buckets").Where(dbx.HashExp{
		"id": bucketID,
	})

	if err := q.One(&b); err != nil {
		return nil, err
	}
	return &b, nil
}
