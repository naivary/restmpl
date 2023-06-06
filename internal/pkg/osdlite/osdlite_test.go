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

func destroy() {
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
	destroy()
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
	if err := testLite.CreateObj(obj); err != nil {
		t.Error(err)
	}
}
