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
	destroy(false)
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
	p := testRandPayload()
	if _, err := o.Write(p); err != nil {
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
}

func TestGetBucket(t *testing.T) {
	b := testBucket()
	if err := testLite.CreateBucket(b); err != nil {
		t.Error(err)
	}
	bG, err := testLite.GetBucket(b.ID)
	if err != nil {
		t.Error(err)
	}
	if bG.ID != b.ID {
		t.Fatalf("id's of the buckets are not equal. Got: %s. Expected: %s", bG.ID, b.ID)
	}
}

func TestCompositeUniqueness(t *testing.T) {
	b1 := testBucket()
	b2 := testBucket()
	o1 := testObj(b1)
	o2 := testObj(b2)
	o3 := testObj(b2)
	o3.Name = o2.Name
	if err := testLite.CreateBucket(b1); err != nil {
		t.Error(err)
	}
	if err := testLite.CreateBucket(b2); err != nil {
		t.Error(err)
	}
	if err := testLite.CreateObject(o1); err != nil {
		t.Error(err)
	}
	if err := testLite.CreateObject(o2); err != nil {
		t.Error(err)
	}
	err := testLite.CreateObject(o3)
	if err == nil {
		t.Fatalf("composite uniqueness failed upon bucket_id and name in objects.")
	}
}
