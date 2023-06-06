package osdlite

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/naivary/apitmpl/internal/pkg/random"
)

var testLite *OSDLite

func testBucket() *bucket {
	return NewBucket(fmt.Sprintf("test_name_%s", uuid.NewString()), "test", uuid.NewString())
}

func testObj(b *bucket) *object {
	return NewObject(b, fmt.Sprintf("test_obj_%s", uuid.NewString()))
}

func testRandPayload() []byte {
	return []byte(random.String(5))
}

func setup() {
	o, err := New()
	if err != nil {
		log.Fatal(err)
	}
	testLite = o
}

func destroy(ok bool) {
	if !ok {
		return
	}
	defer testLite.store.Close()
	files := []string{"osdlite.db", "osdlite.db-shm", "osdlite.db-wal"}
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			log.Fatal(err)
		}
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	destroy(true)
	os.Exit(code)
}

func TestCreateBucket(t *testing.T) {
	b := testBucket()
	if err := testLite.CreateBucket(b); err != nil {
		t.Error(err)
	}
}

func TestCreatObject(t *testing.T) {
	b := testBucket()
	if err := testLite.CreateBucket(b); err != nil {
		t.Error(err)
	}
	obj := testObj(b)
	if _, err := obj.Write(testRandPayload()); err != nil {
		t.Error(err)
	}
	if err := testLite.CreateObject(obj); err != nil {
		t.Error(err)
	}
}

func TestGetObject(t *testing.T) {
	b := testBucket()
	if err := testLite.CreateBucket(b); err != nil {
		t.Error(err)
	}
	o := testObj(b)
	if _, err := o.Write(testRandPayload()); err != nil {
		t.Error(err)
	}
	if err := testLite.CreateObject(o); err != nil {
		t.Error(err)
	}
	oG, err := testLite.GetObject(b.ID, o.ID)
	if err != nil {
		t.Error(err)
	}

	if oG.ID != o.ID {
		t.Fatalf("id's of the objects are not equal. Got: %s. Expected: %s", oG.ID, o.ID)
	}
	fmt.Println(oG)
}
