package osdlite

import (
	"log"
	"testing"
)

var (
	testLite = setup()
)

func setup() *OSDLite {
	o, err := New()
	if err != nil {
		log.Fatal(err)
	}
	return o
}

func TestCreate(t *testing.T) {
	o := object{}
	if err := testLite.Create(o); err != nil {
		t.Error(err)
	}
}
