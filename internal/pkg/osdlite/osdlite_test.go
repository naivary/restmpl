package osdlite

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/pocketbase/dbx"
)

var (
	testLite    = setup()
	randPayload = []byte("some random payload")
)

func setup() *OSDLite {
	o, err := New()
	if err != nil {
		log.Fatal(err)
	}
	return o
}

func TestCreateBucket(t *testing.T) {
	owner := uuid.NewString()
	b := NewBucket("create_bucket", owner, "test")
	if err := testLite.CreatBucket(b); err != nil {
		t.Error(err)
		return
	}
	q := testLite.fs.Select("name").From("buckets").Where(dbx.HashExp{"id": b.ID})
	bu := bucket{}
	if err := q.One(&bu); err != nil {
		t.Error(err)
		return
	}
	if bu.Name != "create_bucket" {
		t.Log(bu)
		t.Fatalf("Name not equal. Expected: create_bucket. Got: %s", bu.Name)
	}
}

func TestCreateObject(t *testing.T) {
	b := NewBucket("create_object", uuid.NewString(), "test")
	if err := testLite.CreatBucket(b); err != nil {
		t.Error(err)
	}
	o := NewObject(b, "create_obj")
	if _, err := o.Write(randPayload); err != nil {
		t.Error(err)
		return
	}
	if err := testLite.CreateObject(o); err != nil {
		t.Error(err)
	}
}

func TestRemoveBucket(t *testing.T) {
	b := NewBucket("remove_bucket", uuid.NewString(), "test")
	if err := testLite.CreatBucket(b); err != nil {
		t.Error(err)
	}
	o := NewObject(b, "remove_bucket_obj")
	if _, err := o.Write(randPayload); err != nil {
		t.Error(err)
		return
	}
	if err := testLite.CreateObject(o); err != nil {
		t.Error(err)
	}
	if err := testLite.RemoveBucket(b); err != nil {
		t.Error(err)
	}
}

func TestRemoveObject(t *testing.T) {
	b := NewBucket("remove_object", uuid.NewString(), "test")
	if err := testLite.CreatBucket(b); err != nil {
		t.Error(err)
		return
	}
	o := NewObject(b, "remove_object_obj")
	if _, err := o.Write(randPayload); err != nil {
		t.Error(err)
		return
	}
	if err := testLite.CreateObject(o); err != nil {
		t.Error(err)
		return
	}
	if err := testLite.RemoveObject(o); err != nil {
		t.Error(err)
		return
	}
}
